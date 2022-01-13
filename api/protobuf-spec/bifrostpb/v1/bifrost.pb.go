// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.13.0
// source: api/protobuf-spec/bifrostpb/v1/bifrost.proto

package v1

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Null struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Null) Reset() {
	*x = Null{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Null) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Null) ProtoMessage() {}

func (x *Null) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Null.ProtoReflect.Descriptor instead.
func (*Null) Descriptor() ([]byte, []int) {
	return file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescGZIP(), []int{0}
}

type ServerNames struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Names []*ServerName `protobuf:"bytes,1,rep,name=Names,proto3" json:"Names,omitempty"`
}

func (x *ServerNames) Reset() {
	*x = ServerNames{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerNames) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerNames) ProtoMessage() {}

func (x *ServerNames) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerNames.ProtoReflect.Descriptor instead.
func (*ServerNames) Descriptor() ([]byte, []int) {
	return file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescGZIP(), []int{1}
}

func (x *ServerNames) GetNames() []*ServerName {
	if x != nil {
		return x.Names
	}
	return nil
}

type ServerName struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=Name,proto3" json:"Name,omitempty"`
}

func (x *ServerName) Reset() {
	*x = ServerName{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerName) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerName) ProtoMessage() {}

func (x *ServerName) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerName.ProtoReflect.Descriptor instead.
func (*ServerName) Descriptor() ([]byte, []int) {
	return file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescGZIP(), []int{2}
}

func (x *ServerName) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type ServerConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServerName string `protobuf:"bytes,1,opt,name=ServerName,proto3" json:"ServerName,omitempty"`
	JsonData   []byte `protobuf:"bytes,2,opt,name=JsonData,proto3" json:"JsonData,omitempty"`
}

func (x *ServerConfig) Reset() {
	*x = ServerConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerConfig) ProtoMessage() {}

func (x *ServerConfig) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerConfig.ProtoReflect.Descriptor instead.
func (*ServerConfig) Descriptor() ([]byte, []int) {
	return file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescGZIP(), []int{3}
}

func (x *ServerConfig) GetServerName() string {
	if x != nil {
		return x.ServerName
	}
	return ""
}

func (x *ServerConfig) GetJsonData() []byte {
	if x != nil {
		return x.JsonData
	}
	return nil
}

type Response struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Msg []byte `protobuf:"bytes,1,opt,name=Msg,proto3" json:"Msg,omitempty"`
}

func (x *Response) Reset() {
	*x = Response{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Response) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Response) ProtoMessage() {}

func (x *Response) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Response.ProtoReflect.Descriptor instead.
func (*Response) Descriptor() ([]byte, []int) {
	return file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescGZIP(), []int{4}
}

func (x *Response) GetMsg() []byte {
	if x != nil {
		return x.Msg
	}
	return nil
}

type Statistics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JsonData []byte `protobuf:"bytes,1,opt,name=JsonData,proto3" json:"JsonData,omitempty"`
}

func (x *Statistics) Reset() {
	*x = Statistics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Statistics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Statistics) ProtoMessage() {}

func (x *Statistics) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Statistics.ProtoReflect.Descriptor instead.
func (*Statistics) Descriptor() ([]byte, []int) {
	return file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescGZIP(), []int{5}
}

func (x *Statistics) GetJsonData() []byte {
	if x != nil {
		return x.JsonData
	}
	return nil
}

type Metrics struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	JsonData []byte `protobuf:"bytes,1,opt,name=JsonData,proto3" json:"JsonData,omitempty"`
}

func (x *Metrics) Reset() {
	*x = Metrics{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metrics) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metrics) ProtoMessage() {}

func (x *Metrics) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metrics.ProtoReflect.Descriptor instead.
func (*Metrics) Descriptor() ([]byte, []int) {
	return file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescGZIP(), []int{6}
}

func (x *Metrics) GetJsonData() []byte {
	if x != nil {
		return x.JsonData
	}
	return nil
}

type LogWatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ServerName string `protobuf:"bytes,1,opt,name=ServerName,proto3" json:"ServerName,omitempty"`
	LogName    string `protobuf:"bytes,2,opt,name=LogName,proto3" json:"LogName,omitempty"`
}

func (x *LogWatchRequest) Reset() {
	*x = LogWatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LogWatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LogWatchRequest) ProtoMessage() {}

func (x *LogWatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LogWatchRequest.ProtoReflect.Descriptor instead.
func (*LogWatchRequest) Descriptor() ([]byte, []int) {
	return file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescGZIP(), []int{7}
}

func (x *LogWatchRequest) GetServerName() string {
	if x != nil {
		return x.ServerName
	}
	return ""
}

func (x *LogWatchRequest) GetLogName() string {
	if x != nil {
		return x.LogName
	}
	return ""
}

var File_api_protobuf_spec_bifrostpb_v1_bifrost_proto protoreflect.FileDescriptor

var file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x61, 0x70, 0x69, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2d, 0x73,
	0x70, 0x65, 0x63, 0x2f, 0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x70, 0x62, 0x2f, 0x76, 0x31,
	0x2f, 0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x70, 0x62, 0x22, 0x06, 0x0a, 0x04, 0x4e, 0x75, 0x6c,
	0x6c, 0x22, 0x3a, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x73,
	0x12, 0x2b, 0x0a, 0x05, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x52, 0x05, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x22, 0x20, 0x0a,
	0x0a, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x4e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x4e, 0x61, 0x6d, 0x65, 0x22,
	0x4a, 0x0a, 0x0c, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12,
	0x1e, 0x0a, 0x0a, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1a, 0x0a, 0x08, 0x4a, 0x73, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x08, 0x4a, 0x73, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x22, 0x1c, 0x0a, 0x08, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x4d, 0x73, 0x67, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x4d, 0x73, 0x67, 0x22, 0x28, 0x0a, 0x0a, 0x53, 0x74, 0x61,
	0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x4a, 0x73, 0x6f, 0x6e, 0x44,
	0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x4a, 0x73, 0x6f, 0x6e, 0x44,
	0x61, 0x74, 0x61, 0x22, 0x25, 0x0a, 0x07, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x12, 0x1a,
	0x0a, 0x08, 0x4a, 0x73, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c,
	0x52, 0x08, 0x4a, 0x73, 0x6f, 0x6e, 0x44, 0x61, 0x74, 0x61, 0x22, 0x4b, 0x0a, 0x0f, 0x4c, 0x6f,
	0x67, 0x57, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a,
	0x0a, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x4c, 0x6f, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x4c, 0x6f, 0x67, 0x4e, 0x61, 0x6d, 0x65, 0x32, 0xc5, 0x01, 0x0a, 0x0f, 0x57, 0x65, 0x62, 0x53,
	0x65, 0x72, 0x76, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x12, 0x3b, 0x0a, 0x0e, 0x47,
	0x65, 0x74, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x12, 0x0f, 0x2e,
	0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x4e, 0x75, 0x6c, 0x6c, 0x1a, 0x16,
	0x2e, 0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x4e, 0x61, 0x6d, 0x65, 0x73, 0x22, 0x00, 0x12, 0x39, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12,
	0x15, 0x2e, 0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x72, 0x76,
	0x65, 0x72, 0x4e, 0x61, 0x6d, 0x65, 0x1a, 0x17, 0x2e, 0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74,
	0x70, 0x62, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x22,
	0x00, 0x30, 0x01, 0x12, 0x3a, 0x0a, 0x06, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x12, 0x17, 0x2e,
	0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x1a, 0x13, 0x2e, 0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74,
	0x70, 0x62, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x32,
	0x4e, 0x0a, 0x13, 0x57, 0x65, 0x62, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x12, 0x37, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x15, 0x2e,
	0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72,
	0x4e, 0x61, 0x6d, 0x65, 0x1a, 0x15, 0x2e, 0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x70, 0x62,
	0x2e, 0x53, 0x74, 0x61, 0x74, 0x69, 0x73, 0x74, 0x69, 0x63, 0x73, 0x22, 0x00, 0x30, 0x01, 0x32,
	0x41, 0x0a, 0x0f, 0x57, 0x65, 0x62, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x2e, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x0f, 0x2e, 0x62, 0x69, 0x66, 0x72,
	0x6f, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x4e, 0x75, 0x6c, 0x6c, 0x1a, 0x12, 0x2e, 0x62, 0x69, 0x66,
	0x72, 0x6f, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x4d, 0x65, 0x74, 0x72, 0x69, 0x63, 0x73, 0x22, 0x00,
	0x30, 0x01, 0x32, 0x55, 0x0a, 0x13, 0x57, 0x65, 0x62, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x4c,
	0x6f, 0x67, 0x57, 0x61, 0x74, 0x63, 0x68, 0x65, 0x72, 0x12, 0x3e, 0x0a, 0x05, 0x57, 0x61, 0x74,
	0x63, 0x68, 0x12, 0x1a, 0x2e, 0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x4c,
	0x6f, 0x67, 0x57, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x13,
	0x2e, 0x62, 0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x70, 0x62, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x20, 0x5a, 0x1e, 0x61, 0x70, 0x69,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2d, 0x73, 0x70, 0x65, 0x63, 0x2f, 0x62,
	0x69, 0x66, 0x72, 0x6f, 0x73, 0x74, 0x70, 0x62, 0x2f, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescOnce sync.Once
	file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescData = file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDesc
)

func file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescGZIP() []byte {
	file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescOnce.Do(func() {
		file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescData)
	})
	return file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDescData
}

var file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_goTypes = []interface{}{
	(*Null)(nil),            // 0: bifrostpb.Null
	(*ServerNames)(nil),     // 1: bifrostpb.ServerNames
	(*ServerName)(nil),      // 2: bifrostpb.ServerName
	(*ServerConfig)(nil),    // 3: bifrostpb.ServerConfig
	(*Response)(nil),        // 4: bifrostpb.Response
	(*Statistics)(nil),      // 5: bifrostpb.Statistics
	(*Metrics)(nil),         // 6: bifrostpb.Metrics
	(*LogWatchRequest)(nil), // 7: bifrostpb.LogWatchRequest
}
var file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_depIdxs = []int32{
	2, // 0: bifrostpb.ServerNames.Names:type_name -> bifrostpb.ServerName
	0, // 1: bifrostpb.WebServerConfig.GetServerNames:input_type -> bifrostpb.Null
	2, // 2: bifrostpb.WebServerConfig.Get:input_type -> bifrostpb.ServerName
	3, // 3: bifrostpb.WebServerConfig.Update:input_type -> bifrostpb.ServerConfig
	2, // 4: bifrostpb.WebServerStatistics.Get:input_type -> bifrostpb.ServerName
	0, // 5: bifrostpb.WebServerStatus.Get:input_type -> bifrostpb.Null
	7, // 6: bifrostpb.WebServerLogWatcher.Watch:input_type -> bifrostpb.LogWatchRequest
	1, // 7: bifrostpb.WebServerConfig.GetServerNames:output_type -> bifrostpb.ServerNames
	3, // 8: bifrostpb.WebServerConfig.Get:output_type -> bifrostpb.ServerConfig
	4, // 9: bifrostpb.WebServerConfig.Update:output_type -> bifrostpb.Response
	5, // 10: bifrostpb.WebServerStatistics.Get:output_type -> bifrostpb.Statistics
	6, // 11: bifrostpb.WebServerStatus.Get:output_type -> bifrostpb.Metrics
	4, // 12: bifrostpb.WebServerLogWatcher.Watch:output_type -> bifrostpb.Response
	7, // [7:13] is the sub-list for method output_type
	1, // [1:7] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_init() }
func file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_init() {
	if File_api_protobuf_spec_bifrostpb_v1_bifrost_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Null); i {
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
		file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerNames); i {
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
		file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerName); i {
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
		file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ServerConfig); i {
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
		file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Response); i {
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
		file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Statistics); i {
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
		file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Metrics); i {
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
		file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LogWatchRequest); i {
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
			RawDescriptor: file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   4,
		},
		GoTypes:           file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_goTypes,
		DependencyIndexes: file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_depIdxs,
		MessageInfos:      file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_msgTypes,
	}.Build()
	File_api_protobuf_spec_bifrostpb_v1_bifrost_proto = out.File
	file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_rawDesc = nil
	file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_goTypes = nil
	file_api_protobuf_spec_bifrostpb_v1_bifrost_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// WebServerConfigClient is the client API for WebServerConfig service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WebServerConfigClient interface {
	GetServerNames(ctx context.Context, in *Null, opts ...grpc.CallOption) (*ServerNames, error)
	Get(ctx context.Context, in *ServerName, opts ...grpc.CallOption) (WebServerConfig_GetClient, error)
	Update(ctx context.Context, opts ...grpc.CallOption) (WebServerConfig_UpdateClient, error)
}

type webServerConfigClient struct {
	cc grpc.ClientConnInterface
}

func NewWebServerConfigClient(cc grpc.ClientConnInterface) WebServerConfigClient {
	return &webServerConfigClient{cc}
}

func (c *webServerConfigClient) GetServerNames(ctx context.Context, in *Null, opts ...grpc.CallOption) (*ServerNames, error) {
	out := new(ServerNames)
	err := c.cc.Invoke(ctx, "/bifrostpb.WebServerConfig/GetServerNames", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *webServerConfigClient) Get(ctx context.Context, in *ServerName, opts ...grpc.CallOption) (WebServerConfig_GetClient, error) {
	stream, err := c.cc.NewStream(ctx, &_WebServerConfig_serviceDesc.Streams[0], "/bifrostpb.WebServerConfig/Get", opts...)
	if err != nil {
		return nil, err
	}
	x := &webServerConfigGetClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type WebServerConfig_GetClient interface {
	Recv() (*ServerConfig, error)
	grpc.ClientStream
}

type webServerConfigGetClient struct {
	grpc.ClientStream
}

func (x *webServerConfigGetClient) Recv() (*ServerConfig, error) {
	m := new(ServerConfig)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *webServerConfigClient) Update(ctx context.Context, opts ...grpc.CallOption) (WebServerConfig_UpdateClient, error) {
	stream, err := c.cc.NewStream(ctx, &_WebServerConfig_serviceDesc.Streams[1], "/bifrostpb.WebServerConfig/Update", opts...)
	if err != nil {
		return nil, err
	}
	x := &webServerConfigUpdateClient{stream}
	return x, nil
}

type WebServerConfig_UpdateClient interface {
	Send(*ServerConfig) error
	CloseAndRecv() (*Response, error)
	grpc.ClientStream
}

type webServerConfigUpdateClient struct {
	grpc.ClientStream
}

func (x *webServerConfigUpdateClient) Send(m *ServerConfig) error {
	return x.ClientStream.SendMsg(m)
}

func (x *webServerConfigUpdateClient) CloseAndRecv() (*Response, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// WebServerConfigServer is the server API for WebServerConfig service.
type WebServerConfigServer interface {
	GetServerNames(context.Context, *Null) (*ServerNames, error)
	Get(*ServerName, WebServerConfig_GetServer) error
	Update(WebServerConfig_UpdateServer) error
}

// UnimplementedWebServerConfigServer can be embedded to have forward compatible implementations.
type UnimplementedWebServerConfigServer struct {
}

func (*UnimplementedWebServerConfigServer) GetServerNames(context.Context, *Null) (*ServerNames, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetServerNames not implemented")
}
func (*UnimplementedWebServerConfigServer) Get(*ServerName, WebServerConfig_GetServer) error {
	return status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (*UnimplementedWebServerConfigServer) Update(WebServerConfig_UpdateServer) error {
	return status.Errorf(codes.Unimplemented, "method Update not implemented")
}

func RegisterWebServerConfigServer(s *grpc.Server, srv WebServerConfigServer) {
	s.RegisterService(&_WebServerConfig_serviceDesc, srv)
}

func _WebServerConfig_GetServerNames_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Null)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(WebServerConfigServer).GetServerNames(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bifrostpb.WebServerConfig/GetServerNames",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(WebServerConfigServer).GetServerNames(ctx, req.(*Null))
	}
	return interceptor(ctx, in, info, handler)
}

func _WebServerConfig_Get_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ServerName)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WebServerConfigServer).Get(m, &webServerConfigGetServer{stream})
}

type WebServerConfig_GetServer interface {
	Send(*ServerConfig) error
	grpc.ServerStream
}

type webServerConfigGetServer struct {
	grpc.ServerStream
}

func (x *webServerConfigGetServer) Send(m *ServerConfig) error {
	return x.ServerStream.SendMsg(m)
}

func _WebServerConfig_Update_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(WebServerConfigServer).Update(&webServerConfigUpdateServer{stream})
}

type WebServerConfig_UpdateServer interface {
	SendAndClose(*Response) error
	Recv() (*ServerConfig, error)
	grpc.ServerStream
}

type webServerConfigUpdateServer struct {
	grpc.ServerStream
}

func (x *webServerConfigUpdateServer) SendAndClose(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *webServerConfigUpdateServer) Recv() (*ServerConfig, error) {
	m := new(ServerConfig)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _WebServerConfig_serviceDesc = grpc.ServiceDesc{
	ServiceName: "bifrostpb.WebServerConfig",
	HandlerType: (*WebServerConfigServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetServerNames",
			Handler:    _WebServerConfig_GetServerNames_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Get",
			Handler:       _WebServerConfig_Get_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "Update",
			Handler:       _WebServerConfig_Update_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "api/protobuf-spec/bifrostpb/v1/bifrost.proto",
}

// WebServerStatisticsClient is the client API for WebServerStatistics service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WebServerStatisticsClient interface {
	Get(ctx context.Context, in *ServerName, opts ...grpc.CallOption) (WebServerStatistics_GetClient, error)
}

type webServerStatisticsClient struct {
	cc grpc.ClientConnInterface
}

func NewWebServerStatisticsClient(cc grpc.ClientConnInterface) WebServerStatisticsClient {
	return &webServerStatisticsClient{cc}
}

func (c *webServerStatisticsClient) Get(ctx context.Context, in *ServerName, opts ...grpc.CallOption) (WebServerStatistics_GetClient, error) {
	stream, err := c.cc.NewStream(ctx, &_WebServerStatistics_serviceDesc.Streams[0], "/bifrostpb.WebServerStatistics/Get", opts...)
	if err != nil {
		return nil, err
	}
	x := &webServerStatisticsGetClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type WebServerStatistics_GetClient interface {
	Recv() (*Statistics, error)
	grpc.ClientStream
}

type webServerStatisticsGetClient struct {
	grpc.ClientStream
}

func (x *webServerStatisticsGetClient) Recv() (*Statistics, error) {
	m := new(Statistics)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// WebServerStatisticsServer is the server API for WebServerStatistics service.
type WebServerStatisticsServer interface {
	Get(*ServerName, WebServerStatistics_GetServer) error
}

// UnimplementedWebServerStatisticsServer can be embedded to have forward compatible implementations.
type UnimplementedWebServerStatisticsServer struct {
}

func (*UnimplementedWebServerStatisticsServer) Get(*ServerName, WebServerStatistics_GetServer) error {
	return status.Errorf(codes.Unimplemented, "method Get not implemented")
}

func RegisterWebServerStatisticsServer(s *grpc.Server, srv WebServerStatisticsServer) {
	s.RegisterService(&_WebServerStatistics_serviceDesc, srv)
}

func _WebServerStatistics_Get_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ServerName)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WebServerStatisticsServer).Get(m, &webServerStatisticsGetServer{stream})
}

type WebServerStatistics_GetServer interface {
	Send(*Statistics) error
	grpc.ServerStream
}

type webServerStatisticsGetServer struct {
	grpc.ServerStream
}

func (x *webServerStatisticsGetServer) Send(m *Statistics) error {
	return x.ServerStream.SendMsg(m)
}

var _WebServerStatistics_serviceDesc = grpc.ServiceDesc{
	ServiceName: "bifrostpb.WebServerStatistics",
	HandlerType: (*WebServerStatisticsServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Get",
			Handler:       _WebServerStatistics_Get_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/protobuf-spec/bifrostpb/v1/bifrost.proto",
}

// WebServerStatusClient is the client API for WebServerStatus service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WebServerStatusClient interface {
	Get(ctx context.Context, in *Null, opts ...grpc.CallOption) (WebServerStatus_GetClient, error)
}

type webServerStatusClient struct {
	cc grpc.ClientConnInterface
}

func NewWebServerStatusClient(cc grpc.ClientConnInterface) WebServerStatusClient {
	return &webServerStatusClient{cc}
}

func (c *webServerStatusClient) Get(ctx context.Context, in *Null, opts ...grpc.CallOption) (WebServerStatus_GetClient, error) {
	stream, err := c.cc.NewStream(ctx, &_WebServerStatus_serviceDesc.Streams[0], "/bifrostpb.WebServerStatus/Get", opts...)
	if err != nil {
		return nil, err
	}
	x := &webServerStatusGetClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type WebServerStatus_GetClient interface {
	Recv() (*Metrics, error)
	grpc.ClientStream
}

type webServerStatusGetClient struct {
	grpc.ClientStream
}

func (x *webServerStatusGetClient) Recv() (*Metrics, error) {
	m := new(Metrics)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// WebServerStatusServer is the server API for WebServerStatus service.
type WebServerStatusServer interface {
	Get(*Null, WebServerStatus_GetServer) error
}

// UnimplementedWebServerStatusServer can be embedded to have forward compatible implementations.
type UnimplementedWebServerStatusServer struct {
}

func (*UnimplementedWebServerStatusServer) Get(*Null, WebServerStatus_GetServer) error {
	return status.Errorf(codes.Unimplemented, "method Get not implemented")
}

func RegisterWebServerStatusServer(s *grpc.Server, srv WebServerStatusServer) {
	s.RegisterService(&_WebServerStatus_serviceDesc, srv)
}

func _WebServerStatus_Get_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Null)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(WebServerStatusServer).Get(m, &webServerStatusGetServer{stream})
}

type WebServerStatus_GetServer interface {
	Send(*Metrics) error
	grpc.ServerStream
}

type webServerStatusGetServer struct {
	grpc.ServerStream
}

func (x *webServerStatusGetServer) Send(m *Metrics) error {
	return x.ServerStream.SendMsg(m)
}

var _WebServerStatus_serviceDesc = grpc.ServiceDesc{
	ServiceName: "bifrostpb.WebServerStatus",
	HandlerType: (*WebServerStatusServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Get",
			Handler:       _WebServerStatus_Get_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/protobuf-spec/bifrostpb/v1/bifrost.proto",
}

// WebServerLogWatcherClient is the client API for WebServerLogWatcher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type WebServerLogWatcherClient interface {
	Watch(ctx context.Context, opts ...grpc.CallOption) (WebServerLogWatcher_WatchClient, error)
}

type webServerLogWatcherClient struct {
	cc grpc.ClientConnInterface
}

func NewWebServerLogWatcherClient(cc grpc.ClientConnInterface) WebServerLogWatcherClient {
	return &webServerLogWatcherClient{cc}
}

func (c *webServerLogWatcherClient) Watch(ctx context.Context, opts ...grpc.CallOption) (WebServerLogWatcher_WatchClient, error) {
	stream, err := c.cc.NewStream(ctx, &_WebServerLogWatcher_serviceDesc.Streams[0], "/bifrostpb.WebServerLogWatcher/Watch", opts...)
	if err != nil {
		return nil, err
	}
	x := &webServerLogWatcherWatchClient{stream}
	return x, nil
}

type WebServerLogWatcher_WatchClient interface {
	Send(*LogWatchRequest) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type webServerLogWatcherWatchClient struct {
	grpc.ClientStream
}

func (x *webServerLogWatcherWatchClient) Send(m *LogWatchRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *webServerLogWatcherWatchClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// WebServerLogWatcherServer is the server API for WebServerLogWatcher service.
type WebServerLogWatcherServer interface {
	Watch(WebServerLogWatcher_WatchServer) error
}

// UnimplementedWebServerLogWatcherServer can be embedded to have forward compatible implementations.
type UnimplementedWebServerLogWatcherServer struct {
}

func (*UnimplementedWebServerLogWatcherServer) Watch(WebServerLogWatcher_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "method Watch not implemented")
}

func RegisterWebServerLogWatcherServer(s *grpc.Server, srv WebServerLogWatcherServer) {
	s.RegisterService(&_WebServerLogWatcher_serviceDesc, srv)
}

func _WebServerLogWatcher_Watch_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(WebServerLogWatcherServer).Watch(&webServerLogWatcherWatchServer{stream})
}

type WebServerLogWatcher_WatchServer interface {
	Send(*Response) error
	Recv() (*LogWatchRequest, error)
	grpc.ServerStream
}

type webServerLogWatcherWatchServer struct {
	grpc.ServerStream
}

func (x *webServerLogWatcherWatchServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *webServerLogWatcherWatchServer) Recv() (*LogWatchRequest, error) {
	m := new(LogWatchRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _WebServerLogWatcher_serviceDesc = grpc.ServiceDesc{
	ServiceName: "bifrostpb.WebServerLogWatcher",
	HandlerType: (*WebServerLogWatcherServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Watch",
			Handler:       _WebServerLogWatcher_Watch_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "api/protobuf-spec/bifrostpb/v1/bifrost.proto",
}
