// Copyright 2023 Google LLC
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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.21.12
// source: internal/tool/ssh/impl/ssh.proto

package impl

import (
	protos "github.com/ServiceWeaver/weaver/runtime/protos"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// BabysitterInfo contains app deployment information that is needed by a
// babysitter started using SSH to manage a colocation group.
type BabysitterInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Deployment  *protos.Deployment      `protobuf:"bytes,1,opt,name=deployment,proto3" json:"deployment,omitempty"`
	Group       *protos.ColocationGroup `protobuf:"bytes,2,opt,name=group,proto3" json:"group,omitempty"`
	ReplicaId   int32                   `protobuf:"varint,3,opt,name=replica_id,json=replicaId,proto3" json:"replica_id,omitempty"`
	ManagerAddr string                  `protobuf:"bytes,4,opt,name=manager_addr,json=managerAddr,proto3" json:"manager_addr,omitempty"`
	LogDir      string                  `protobuf:"bytes,5,opt,name=logDir,proto3" json:"logDir,omitempty"`
}

func (x *BabysitterInfo) Reset() {
	*x = BabysitterInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tool_ssh_impl_ssh_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BabysitterInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BabysitterInfo) ProtoMessage() {}

func (x *BabysitterInfo) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tool_ssh_impl_ssh_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BabysitterInfo.ProtoReflect.Descriptor instead.
func (*BabysitterInfo) Descriptor() ([]byte, []int) {
	return file_internal_tool_ssh_impl_ssh_proto_rawDescGZIP(), []int{0}
}

func (x *BabysitterInfo) GetDeployment() *protos.Deployment {
	if x != nil {
		return x.Deployment
	}
	return nil
}

func (x *BabysitterInfo) GetGroup() *protos.ColocationGroup {
	if x != nil {
		return x.Group
	}
	return nil
}

func (x *BabysitterInfo) GetReplicaId() int32 {
	if x != nil {
		return x.ReplicaId
	}
	return 0
}

func (x *BabysitterInfo) GetManagerAddr() string {
	if x != nil {
		return x.ManagerAddr
	}
	return ""
}

func (x *BabysitterInfo) GetLogDir() string {
	if x != nil {
		return x.LogDir
	}
	return ""
}

// When a babysitter receives a GetComponentsToStart request from a weavelet, it
// forwards the request to the manager. But first, it wraps the request in a
// GetComponents message to include its group name.
type GetComponents struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Group         string                       `protobuf:"bytes,1,opt,name=group,proto3" json:"group,omitempty"`
	GetComponents *protos.GetComponentsToStart `protobuf:"bytes,2,opt,name=get_components,json=getComponents,proto3" json:"get_components,omitempty"`
}

func (x *GetComponents) Reset() {
	*x = GetComponents{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tool_ssh_impl_ssh_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetComponents) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetComponents) ProtoMessage() {}

func (x *GetComponents) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tool_ssh_impl_ssh_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetComponents.ProtoReflect.Descriptor instead.
func (*GetComponents) Descriptor() ([]byte, []int) {
	return file_internal_tool_ssh_impl_ssh_proto_rawDescGZIP(), []int{1}
}

func (x *GetComponents) GetGroup() string {
	if x != nil {
		return x.Group
	}
	return ""
}

func (x *GetComponents) GetGetComponents() *protos.GetComponentsToStart {
	if x != nil {
		return x.GetComponents
	}
	return nil
}

// BabysitterMetrics is a snapshot of a deployment's metrics as collected by a
// babysitter for a given colocation group.
type BabysitterMetrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	GroupName string                   `protobuf:"bytes,1,opt,name=group_name,json=groupName,proto3" json:"group_name,omitempty"`
	ReplicaId int32                    `protobuf:"varint,2,opt,name=replica_id,json=replicaId,proto3" json:"replica_id,omitempty"`
	Metrics   []*protos.MetricSnapshot `protobuf:"bytes,3,rep,name=metrics,proto3" json:"metrics,omitempty"`
}

func (x *BabysitterMetrics) Reset() {
	*x = BabysitterMetrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_tool_ssh_impl_ssh_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BabysitterMetrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BabysitterMetrics) ProtoMessage() {}

func (x *BabysitterMetrics) ProtoReflect() protoreflect.Message {
	mi := &file_internal_tool_ssh_impl_ssh_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BabysitterMetrics.ProtoReflect.Descriptor instead.
func (*BabysitterMetrics) Descriptor() ([]byte, []int) {
	return file_internal_tool_ssh_impl_ssh_proto_rawDescGZIP(), []int{2}
}

func (x *BabysitterMetrics) GetGroupName() string {
	if x != nil {
		return x.GroupName
	}
	return ""
}

func (x *BabysitterMetrics) GetReplicaId() int32 {
	if x != nil {
		return x.ReplicaId
	}
	return 0
}

func (x *BabysitterMetrics) GetMetrics() []*protos.MetricSnapshot {
	if x != nil {
		return x.Metrics
	}
	return nil
}

var File_internal_tool_ssh_impl_ssh_proto protoreflect.FileDescriptor

var file_internal_tool_ssh_impl_ssh_proto_rawDesc = []byte{
	0x0a, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x6f, 0x6f, 0x6c, 0x2f,
	0x73, 0x73, 0x68, 0x2f, 0x69, 0x6d, 0x70, 0x6c, 0x2f, 0x73, 0x73, 0x68, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x04, 0x69, 0x6d, 0x70, 0x6c, 0x1a, 0x1c, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcf, 0x01, 0x0a, 0x0e, 0x42, 0x61, 0x62, 0x79, 0x73,
	0x69, 0x74, 0x74, 0x65, 0x72, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x33, 0x0a, 0x0a, 0x64, 0x65, 0x70,
	0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e,
	0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x44, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65,
	0x6e, 0x74, 0x52, 0x0a, 0x64, 0x65, 0x70, 0x6c, 0x6f, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x2e,
	0x0a, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x18, 0x2e,
	0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x43, 0x6f, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x1d,
	0x0a, 0x0a, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x09, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x49, 0x64, 0x12, 0x21, 0x0a,
	0x0c, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0b, 0x6d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x72, 0x41, 0x64, 0x64, 0x72,
	0x12, 0x16, 0x0a, 0x06, 0x6c, 0x6f, 0x67, 0x44, 0x69, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x6c, 0x6f, 0x67, 0x44, 0x69, 0x72, 0x22, 0x6b, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x43,
	0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x14, 0x0a, 0x05, 0x67, 0x72, 0x6f,
	0x75, 0x70, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x12,
	0x44, 0x0a, 0x0e, 0x67, 0x65, 0x74, 0x5f, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74,
	0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x72, 0x75, 0x6e, 0x74, 0x69, 0x6d,
	0x65, 0x2e, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x73, 0x54,
	0x6f, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x0d, 0x67, 0x65, 0x74, 0x43, 0x6f, 0x6d, 0x70, 0x6f,
	0x6e, 0x65, 0x6e, 0x74, 0x73, 0x22, 0x84, 0x01, 0x0a, 0x11, 0x42, 0x61, 0x62, 0x79, 0x73, 0x69,
	0x74, 0x74, 0x65, 0x72, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x1d, 0x0a, 0x0a, 0x67,
	0x72, 0x6f, 0x75, 0x70, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x09, 0x67, 0x72, 0x6f, 0x75, 0x70, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1d, 0x0a, 0x0a, 0x72, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09,
	0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x07, 0x6d, 0x65, 0x74,
	0x72, 0x69, 0x63, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x72, 0x75, 0x6e,
	0x74, 0x69, 0x6d, 0x65, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x53, 0x6e, 0x61, 0x70, 0x73,
	0x68, 0x6f, 0x74, 0x52, 0x07, 0x6d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x42, 0x38, 0x5a, 0x36,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x57, 0x65, 0x61, 0x76, 0x65, 0x72, 0x2f, 0x77, 0x65, 0x61, 0x76, 0x65, 0x72, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x74, 0x6f, 0x6f, 0x6c, 0x2f, 0x73, 0x73,
	0x68, 0x2f, 0x69, 0x6d, 0x70, 0x6c, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_tool_ssh_impl_ssh_proto_rawDescOnce sync.Once
	file_internal_tool_ssh_impl_ssh_proto_rawDescData = file_internal_tool_ssh_impl_ssh_proto_rawDesc
)

func file_internal_tool_ssh_impl_ssh_proto_rawDescGZIP() []byte {
	file_internal_tool_ssh_impl_ssh_proto_rawDescOnce.Do(func() {
		file_internal_tool_ssh_impl_ssh_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_tool_ssh_impl_ssh_proto_rawDescData)
	})
	return file_internal_tool_ssh_impl_ssh_proto_rawDescData
}

var file_internal_tool_ssh_impl_ssh_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_internal_tool_ssh_impl_ssh_proto_goTypes = []interface{}{
	(*BabysitterInfo)(nil),              // 0: impl.BabysitterInfo
	(*GetComponents)(nil),               // 1: impl.GetComponents
	(*BabysitterMetrics)(nil),           // 2: impl.BabysitterMetrics
	(*protos.Deployment)(nil),           // 3: runtime.Deployment
	(*protos.ColocationGroup)(nil),      // 4: runtime.ColocationGroup
	(*protos.GetComponentsToStart)(nil), // 5: runtime.GetComponentsToStart
	(*protos.MetricSnapshot)(nil),       // 6: runtime.MetricSnapshot
}
var file_internal_tool_ssh_impl_ssh_proto_depIdxs = []int32{
	3, // 0: impl.BabysitterInfo.deployment:type_name -> runtime.Deployment
	4, // 1: impl.BabysitterInfo.group:type_name -> runtime.ColocationGroup
	5, // 2: impl.GetComponents.get_components:type_name -> runtime.GetComponentsToStart
	6, // 3: impl.BabysitterMetrics.metrics:type_name -> runtime.MetricSnapshot
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_internal_tool_ssh_impl_ssh_proto_init() }
func file_internal_tool_ssh_impl_ssh_proto_init() {
	if File_internal_tool_ssh_impl_ssh_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_tool_ssh_impl_ssh_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BabysitterInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_tool_ssh_impl_ssh_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetComponents); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_internal_tool_ssh_impl_ssh_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BabysitterMetrics); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_tool_ssh_impl_ssh_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_internal_tool_ssh_impl_ssh_proto_goTypes,
		DependencyIndexes: file_internal_tool_ssh_impl_ssh_proto_depIdxs,
		MessageInfos:      file_internal_tool_ssh_impl_ssh_proto_msgTypes,
	}.Build()
	File_internal_tool_ssh_impl_ssh_proto = out.File
	file_internal_tool_ssh_impl_ssh_proto_rawDesc = nil
	file_internal_tool_ssh_impl_ssh_proto_goTypes = nil
	file_internal_tool_ssh_impl_ssh_proto_depIdxs = nil
}
