// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package babysitter

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net"
	"net/http"
	"sort"
	"sync"
	"syscall"
	"time"

	"github.com/ServiceWeaver/weaver/internal/logtype"
	imetrics "github.com/ServiceWeaver/weaver/internal/metrics"
	"github.com/ServiceWeaver/weaver/internal/proxy"
	"github.com/ServiceWeaver/weaver/internal/status"
	"github.com/ServiceWeaver/weaver/internal/versioned"
	"github.com/ServiceWeaver/weaver/runtime/envelope"
	"github.com/ServiceWeaver/weaver/runtime/logging"
	"github.com/ServiceWeaver/weaver/runtime/metrics"
	"github.com/ServiceWeaver/weaver/runtime/perfetto"
	"github.com/ServiceWeaver/weaver/runtime/protomsg"
	"github.com/ServiceWeaver/weaver/runtime/protos"
	"github.com/ServiceWeaver/weaver/runtime/tool"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/sdk/trace"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// The default number of times a component is replicated.
const defaultReplication = 2

// A Babysitter manages an application deployment.
type Babysitter struct {
	ctx     context.Context
	dep     *protos.Deployment
	done    chan error
	started time.Time
	logger  logtype.Logger

	// logSaver processes log entries generated by the weavelet. The entries
	// either have the timestamp produced by the weavelet, or have a nil Time
	// field. Defaults to a log saver that pretty prints log entries to stderr.
	//
	// logSaver is called concurrently from multiple goroutines, so it should
	// be thread safe.
	logSaver func(*protos.LogEntry)

	// traceSaver processes trace spans generated by the weavelet. If nil,
	// weavelet traces are dropped.
	//
	// traceSaver is called concurrently from multiple goroutines, so it should
	// be thread safe.
	traceSaver func([]trace.ReadOnlySpan) error

	// statsProcessor tracks and computes stats to be rendered on the /statusz page.
	statsProcessor *imetrics.StatsProcessor

	// The babysitter places components into co-location groups based on the
	// colocate stanza in a config. For example, consider the following config.
	//
	//     colocate = [
	//         ["A", "B", "C"],
	//         ["D", "E"],
	//     ]
	//
	// The babysitter creates a co-location group with components "A", "B", and
	// "C" and a co-location group with components "D" and "E". All other
	// components are placed in their own co-location group. We use the first
	// listed component as the co-location group name.
	//
	// colocation maps components listed in the colocate stanza to the name of
	// their group.
	colocation map[string]string

	// guards access to the following maps, but not the contents inside
	// the maps.
	mu      sync.Mutex
	groups  map[string]*group     // groups, by group name
	proxies map[string]*proxyInfo // proxies, by listener name
}

// A group contains information about a co-location group.
type group struct {
	name       string                                // group name
	components *versioned.Versioned[map[string]bool] // started components

	// guards the following data structures, but not their contents.
	mu        sync.Mutex
	addresses map[string]bool                                      // weavelet addresses
	envelopes []*envelope.Envelope                                 // envelopes, one per weavelet
	pids      []int64                                              // weavelet pids
	routings  map[string]*versioned.Versioned[*protos.RoutingInfo] // routing info, by component
}

// A proxyInfo contains information about a proxy.
type proxyInfo struct {
	listener string       // listener associated with the proxy
	proxy    *proxy.Proxy // the proxy
	addr     string       // dialable address of the proxy
}

// handler handles a connection to a weavelet.
type handler struct {
	*Babysitter
	g *group
}

var _ envelope.EnvelopeHandler = &handler{}

// NewBabysitter creates a new babysitter.
func NewBabysitter(ctx context.Context, dep *protos.Deployment, logSaver func(*protos.LogEntry)) (*Babysitter, error) {
	logger := logging.FuncLogger{
		Opts: logging.Options{
			App:       dep.App.Name,
			Component: "babysitter",
			Weavelet:  uuid.NewString(),
			Attrs:     []string{"serviceweaver/system", ""},
		},
		Write: logSaver,
	}

	// Create the trace saver.
	traceDB, err := perfetto.Open(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot open Perfetto database: %w", err)
	}
	traceSaver := func(spans []trace.ReadOnlySpan) error {
		return traceDB.Store(ctx, dep.App.Name, dep.Id, spans)
	}

	// Form co-location.
	colocation := map[string]string{}
	for _, group := range dep.App.SameProcess {
		for _, c := range group.Components {
			colocation[c] = group.Components[0]
		}
	}

	b := &Babysitter{
		ctx:            ctx,
		logger:         logger,
		logSaver:       logSaver,
		traceSaver:     traceSaver,
		statsProcessor: imetrics.NewStatsProcessor(),
		done:           make(chan error, 1),
		dep:            dep,
		started:        time.Now(),
		colocation:     colocation,
		groups:         map[string]*group{},
		proxies:        map[string]*proxyInfo{},
	}
	go b.statsProcessor.CollectMetrics(b.ctx, b.readMetrics)
	return b, nil
}

func (b *Babysitter) Done() chan error { return b.done }

// RegisterStatusPages registers the status pages with the provided mux.
func (b *Babysitter) RegisterStatusPages(mux *http.ServeMux) {
	status.RegisterServer(mux, b, b.logger)
}

// group returns the co-location group containing the provided component.
func (b *Babysitter) group(component string) *group {
	b.mu.Lock()
	defer b.mu.Unlock()

	name, ok := b.colocation[component]
	if !ok {
		name = component
	}

	g, ok := b.groups[name]
	if !ok {
		g = &group{
			name:       name,
			addresses:  map[string]bool{},
			components: versioned.Version(map[string]bool{}),
			routings:   map[string]*versioned.Versioned[*protos.RoutingInfo]{},
		}
		b.groups[name] = g
	}
	return g
}

// routing returns the RoutingInfo for the provided component.
//
// REQUIRES: g.mu is held.
func (g *group) routing(component string) *versioned.Versioned[*protos.RoutingInfo] {
	routing, ok := g.routings[component]
	if !ok {
		routing = versioned.Version(&protos.RoutingInfo{})
		g.routings[component] = routing
	}
	return routing
}

// allGroups returns all of the managed colocation groups.
func (b *Babysitter) allGroups() []*group {
	b.mu.Lock()
	defer b.mu.Unlock()
	return maps.Values(b.groups) // creates a new slice
}

// allProxies returns all of the managed proxies.
func (b *Babysitter) allProxies() []*proxyInfo {
	b.mu.Lock()
	defer b.mu.Unlock()
	return maps.Values(b.proxies)
}

// startColocationGroup starts the colocation group hosting the provided
// component, if it hasn't been started already.
//
// REQUIRES: g.mu is held.
func (b *Babysitter) startColocationGroup(g *group) error {
	if len(g.envelopes) == defaultReplication {
		// Already started.
		return nil
	}

	for r := 0; r < defaultReplication; r++ {
		// Start the weavelet and capture its logs, traces, and metrics.
		wlet := &protos.WeaveletInfo{
			App:           b.dep.App.Name,
			DeploymentId:  b.dep.Id,
			Group:         &protos.ColocationGroup{Name: g.name},
			GroupId:       uuid.New().String(),
			Id:            uuid.New().String(),
			SameProcess:   b.dep.App.SameProcess,
			Sections:      b.dep.App.Sections,
			SingleProcess: b.dep.SingleProcess,
			SingleMachine: true,
		}
		h := &handler{b, g}
		e, err := envelope.NewEnvelope(wlet, b.dep.App, h)
		if err != nil {
			return err
		}
		go func() {
			if err := e.Run(b.ctx); err != nil {
				b.done <- err
			}
		}()
		g.envelopes = append(g.envelopes, e)
	}
	return nil
}

// StartComponent implements the envelope.EnvelopeHandler interface.
func (b *Babysitter) StartComponent(req *protos.ComponentToStart) error {
	g := b.group(req.Component)

	// Hold the lock to avoid concurrent updates to g.addresses.
	//
	// TODO(mwhittaker): Reduce lock scope.
	g.mu.Lock()
	defer g.mu.Unlock()

	// Record the component.
	record := func() bool {
		g.components.Lock()
		defer g.components.Unlock()
		if g.components.Val[req.Component] {
			// Component already started, or is in the process of being started.
			return true
		}
		g.components.Val[req.Component] = true
		return false
	}
	if record() { // already started
		return nil
	}

	// Update the routing info.
	routing := g.routing(req.Component)
	update := func() {
		routing.Lock()
		defer routing.Unlock()

		routing.Val.Replicas = maps.Keys(g.addresses)
		if req.Routed {
			assignment := &protos.Assignment{
				App:          b.dep.App.Name,
				DeploymentId: b.dep.Id,
				Component:    req.Component,
			}
			routing.Val.Assignment = routingAlgo(assignment, routing.Val.Replicas)
		}
	}
	update()

	return b.startColocationGroup(g)
}

// RegisterReplica implements the envelope.EnvelopeHandler interface.
func (b *Babysitter) RegisterReplica(req *protos.ReplicaToRegister) error {
	g := b.group(req.Group)
	g.mu.Lock()
	defer g.mu.Unlock()

	// Update addresses and pids.
	if g.addresses[req.Address] {
		// Replica already registered.
		return nil
	}
	g.addresses[req.Address] = true
	g.pids = append(g.pids, req.Pid)

	// Update routing.
	replicas := maps.Keys(g.addresses)
	for _, routing := range g.routings {
		routing.Lock()
		routing.Val.Replicas = replicas
		if routing.Val.Assignment != nil {
			routing.Val.Assignment = routingAlgo(routing.Val.Assignment, replicas)
		}
		routing.Unlock()
	}
	return nil
}

// GetComponentsToStart implements the envelope.EnvelopeHandler interface.
func (h *handler) GetComponentsToStart(req *protos.GetComponentsToStart) (*protos.ComponentsToStart, error) {
	version := h.g.components.RLock(req.Version)
	defer h.g.components.RUnlock()
	return &protos.ComponentsToStart{
		Version:    version,
		Components: maps.Keys(h.g.components.Val),
	}, nil
}

// RecvLogEntry implements the envelope.EnvelopeHandler interface.
func (b *Babysitter) RecvLogEntry(entry *protos.LogEntry) {
	b.logSaver(entry)
}

// RecvTraceSpans implements the envelope.EnvelopeHandler interface.
func (b *Babysitter) RecvTraceSpans(spans []trace.ReadOnlySpan) error {
	if b.traceSaver == nil {
		return nil
	}
	return b.traceSaver(spans)
}

// ReportLoad implements the envelope.EnvelopeHandler interface.
func (b *Babysitter) ReportLoad(*protos.WeaveletLoadReport) error {
	return nil
}

// GetAddress implements the envelope.EnvelopeHandler interface.
func (b *Babysitter) GetAddress(req *protos.GetAddressRequest) (*protos.GetAddressReply, error) {
	return &protos.GetAddressReply{Address: "localhost:0"}, nil
}

// ExportListener implements the envelope.EnvelopeHandler interface.
func (b *Babysitter) ExportListener(req *protos.ExportListenerRequest) (*protos.ExportListenerReply, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	// Update the proxy.
	if p, ok := b.proxies[req.Listener.Name]; ok {
		p.proxy.AddBackend(req.Listener.Addr)
		return &protos.ExportListenerReply{ProxyAddress: p.addr}, nil
	}

	lis, err := net.Listen("tcp", req.LocalAddress)
	if errors.Is(err, syscall.EADDRINUSE) {
		// Don't retry if this address is already in use.
		return &protos.ExportListenerReply{Error: err.Error()}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("proxy listen: %w", err)
	}
	addr := lis.Addr().String()
	b.logger.Info("Proxy listening", "address", addr)
	proxy := proxy.NewProxy(b.logger)
	proxy.AddBackend(req.Listener.Addr)
	b.proxies[req.Listener.Name] = &proxyInfo{
		listener: req.Listener.Name,
		proxy:    proxy,
		addr:     addr,
	}
	go func() {
		if err := serveHTTP(b.ctx, lis, proxy); err != nil {
			b.logger.Error("proxy", err)
		}
	}()
	return &protos.ExportListenerReply{ProxyAddress: addr}, nil
}

// GetRoutingInfo implements the envelope.EnvelopeHandler interface.
func (b *Babysitter) GetRoutingInfo(req *protos.GetRoutingInfo) (*protos.RoutingInfo, error) {
	g := b.group(req.Component)
	g.mu.Lock()
	routing := g.routing(req.Component)
	g.mu.Unlock()

	version := routing.RLock(req.Version)
	defer routing.RUnlock()
	r := protomsg.Clone(routing.Val)
	r.Version = version
	return r, nil
}

func (b *Babysitter) readMetrics() []*metrics.MetricSnapshot {
	var ms []*metrics.MetricSnapshot
	read := func(g *group) {
		g.mu.Lock()
		defer g.mu.Unlock()
		for _, e := range g.envelopes {
			m, err := e.ReadMetrics()
			if err != nil {
				continue
			}
			ms = append(ms, m...)
		}
	}
	for _, g := range b.allGroups() {
		read(g)
	}
	return append(ms, metrics.Snapshot()...)
}

// Profile implements the status.Server interface.
func (b *Babysitter) Profile(_ context.Context, req *protos.RunProfiling) (*protos.Profile, error) {
	// Make a copy of the envelopes so we can operate on it without holding the
	// lock. A profile can last a long time.
	envelopes := map[string][]*envelope.Envelope{}
	for _, g := range b.allGroups() {
		g.mu.Lock()
		envelopes[g.name] = slices.Clone(g.envelopes)
		g.mu.Unlock()
	}

	profile, err := runProfiling(b.ctx, req, envelopes)
	if err != nil {
		return nil, err
	}
	profile.AppName = b.dep.App.Name
	profile.VersionId = b.dep.Id
	return profile, nil
}

// Status implements the status.Server interface.
func (b *Babysitter) Status(ctx context.Context) (*status.Status, error) {
	stats := b.statsProcessor.GetStatsStatusz()
	var components []*status.Component
	for _, g := range b.allGroups() {
		g.components.Lock()
		cs := maps.Keys(g.components.Val)
		g.components.Unlock()
		g.mu.Lock()
		pids := slices.Clone(g.pids)
		g.mu.Unlock()
		for _, component := range cs {
			c := &status.Component{
				Name:  component,
				Group: g.name,
				Pids:  pids,
			}
			components = append(components, c)

			// TODO(mwhittaker): Unify with ui package and remove duplication.
			s := stats[logging.ShortenComponent(component)]
			if s == nil {
				continue
			}
			for _, methodStats := range s {
				c.Methods = append(c.Methods, &status.Method{
					Name: methodStats.Name,
					Minute: &status.MethodStats{
						NumCalls:     methodStats.Minute.NumCalls,
						AvgLatencyMs: methodStats.Minute.AvgLatencyMs,
						RecvKbPerSec: methodStats.Minute.RecvKBPerSec,
						SentKbPerSec: methodStats.Minute.SentKBPerSec,
					},
					Hour: &status.MethodStats{
						NumCalls:     methodStats.Hour.NumCalls,
						AvgLatencyMs: methodStats.Hour.AvgLatencyMs,
						RecvKbPerSec: methodStats.Hour.RecvKBPerSec,
						SentKbPerSec: methodStats.Hour.SentKBPerSec,
					},
					Total: &status.MethodStats{
						NumCalls:     methodStats.Total.NumCalls,
						AvgLatencyMs: methodStats.Total.AvgLatencyMs,
						RecvKbPerSec: methodStats.Total.RecvKBPerSec,
						SentKbPerSec: methodStats.Total.SentKBPerSec,
					},
				})
			}
		}
	}

	var listeners []*status.Listener
	for _, proxy := range b.allProxies() {
		listeners = append(listeners, &status.Listener{
			Name: proxy.listener,
			Addr: proxy.addr,
		})
	}

	return &status.Status{
		App:            b.dep.App.Name,
		DeploymentId:   b.dep.Id,
		SubmissionTime: timestamppb.New(b.started),
		Components:     components,
		Listeners:      listeners,
		Config:         b.dep.App,
	}, nil
}

// Metrics implements the status.Server interface.
func (b *Babysitter) Metrics(ctx context.Context) (*status.Metrics, error) {
	m := &status.Metrics{}
	for _, snap := range b.readMetrics() {
		m.Metrics = append(m.Metrics, snap.ToProto())
	}
	return m, nil
}

// routingAlgo is an implementation of a routing algorithm that distributes the
// entire key space approximately equally across all healthy resources.
//
// The algorithm is as follows:
//   - split the entire key space in a number of slices that is more likely to
//     spread the key space uniformly among all healthy resources.
//   - distribute the slices round robin across all healthy resources
func routingAlgo(currAssignment *protos.Assignment, candidates []string) *protos.Assignment {
	newAssignment := protomsg.Clone(currAssignment)
	newAssignment.Version++

	// Note that the healthy resources should be sorted. This is required because
	// we want to do a deterministic assignment of slices to resources among
	// different invocations, to avoid unnecessary churn while generating
	// new assignments.
	sort.Strings(candidates)

	if len(candidates) == 0 {
		newAssignment.Slices = nil
		return newAssignment
	}

	const minSliceKey = 0
	const maxSliceKey = math.MaxUint64

	// If there is only one healthy resource, assign the entire key space to it.
	if len(candidates) == 1 {
		newAssignment.Slices = []*protos.Assignment_Slice{
			{Start: minSliceKey, Replicas: candidates},
		}
		return newAssignment
	}

	// Compute the total number of slices in the assignment.
	numSlices := nextPowerOfTwo(len(candidates))

	// Split slices in equal subslices in order to generate numSlices.
	splits := [][]uint64{{minSliceKey, maxSliceKey}}
	var curr []uint64
	for ok := true; ok; ok = len(splits) != numSlices {
		curr, splits = splits[0], splits[1:]
		midPoint := curr[0] + uint64(math.Floor(0.5*float64(curr[1]-curr[0])))
		splitl := []uint64{curr[0], midPoint}
		splitr := []uint64{midPoint, curr[1]}
		splits = append(splits, splitl, splitr)
	}

	// Sort the computed slices in increasing order based on the start key, in
	// order to provide a deterministic assignment across multiple runs, hence to
	// minimize churn.
	sort.Slice(splits, func(i, j int) bool {
		return splits[i][0] <= splits[j][0]
	})

	// Assign the computed slices to resources in a round robin fashion.
	slices := make([]*protos.Assignment_Slice, len(splits))
	rId := 0
	for i, s := range splits {
		slices[i] = &protos.Assignment_Slice{
			Start:    s[0],
			Replicas: []string{candidates[rId]},
		}
		rId = (rId + 1) % len(candidates)
	}
	newAssignment.Slices = slices
	return newAssignment
}

// serveHTTP serves HTTP traffic on the provided listener using the provided
// handler. The server is shut down when then provided context is cancelled.
func serveHTTP(ctx context.Context, lis net.Listener, handler http.Handler) error {
	server := http.Server{Handler: handler}
	errs := make(chan error, 1)
	go func() { errs <- server.Serve(lis) }()
	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return server.Shutdown(ctx)
	}
}

// nextPowerOfTwo returns the next power of 2 that is greater or equal to x.
func nextPowerOfTwo(x int) int {
	// If x is already power of 2, return x.
	if x&(x-1) == 0 {
		return x
	}
	return int(math.Pow(2, math.Ceil(math.Log2(float64(x)))))
}

// runProfiling runs a profiling request on a set of processes.
func runProfiling(ctx context.Context, req *protos.RunProfiling, processes map[string][]*envelope.Envelope) (*protos.Profile, error) {
	// Collect together the groups we want to profile.
	groups := make([][]func() (*protos.Profile, error), 0, len(processes))
	for _, envelopes := range processes {
		group := make([]func() (*protos.Profile, error), 0, len(envelopes))
		for _, e := range envelopes {
			group = append(group, func() (*protos.Profile, error) {
				return e.RunProfiling(ctx, req)
			})
		}
		groups = append(groups, group)
	}
	return tool.ProfileGroups(groups)
}
