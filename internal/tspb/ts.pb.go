// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: ts.proto

package tspb

import (
	"reflect"
	"sync"

	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type TestRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *TestRequest) Reset() {
	*x = TestRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ts_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestRequest) ProtoMessage() {}

func (x *TestRequest) ProtoReflect() protoreflect.Message {
	mi := &file_ts_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestRequest.ProtoReflect.Descriptor instead.
func (*TestRequest) Descriptor() ([]byte, []int) {
	return file_ts_proto_rawDescGZIP(), []int{0}
}

func (x *TestRequest) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type TestResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value string `protobuf:"bytes,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *TestResponse) Reset() {
	*x = TestResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_ts_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestResponse) ProtoMessage() {}

func (x *TestResponse) ProtoReflect() protoreflect.Message {
	mi := &file_ts_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestResponse.ProtoReflect.Descriptor instead.
func (*TestResponse) Descriptor() ([]byte, []int) {
	return file_ts_proto_rawDescGZIP(), []int{1}
}

func (x *TestResponse) GetValue() string {
	if x != nil {
		return x.Value
	}
	return ""
}

var File_ts_proto protoreflect.FileDescriptor

var file_ts_proto_rawDesc = []byte{
	0x0a, 0x08, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x26, 0x63, 0x6f, 0x6d, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x61, 0x70, 0x70, 0x75, 0x63, 0x63, 0x69, 0x6e,
	0x6f, 0x74, 0x6d, 0x2e, 0x67, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70,
	0x6c, 0x65, 0x22, 0x1f, 0x0a, 0x0b, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x22, 0x24, 0x0a, 0x0c, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x32, 0x80, 0x01, 0x0a, 0x0b, 0x54, 0x65,
	0x73, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x71, 0x0a, 0x04, 0x54, 0x65, 0x73,
	0x74, 0x12, 0x33, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x61, 0x70, 0x70, 0x75, 0x63, 0x63, 0x69, 0x6e, 0x6f, 0x74, 0x6d, 0x2e, 0x67, 0x63, 0x61, 0x63,
	0x68, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x34, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x61, 0x70, 0x70, 0x75, 0x63, 0x63, 0x69, 0x6e, 0x6f, 0x74, 0x6d,
	0x2e, 0x67, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2e, 0x65, 0x78, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x2e,
	0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x25, 0x5a, 0x23,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x61, 0x70, 0x70, 0x75,
	0x63, 0x63, 0x69, 0x6e, 0x6f, 0x74, 0x6d, 0x2f, 0x67, 0x63, 0x61, 0x63, 0x68, 0x65, 0x2f, 0x74,
	0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_ts_proto_rawDescOnce sync.Once
	file_ts_proto_rawDescData = file_ts_proto_rawDesc
)

func file_ts_proto_rawDescGZIP() []byte {
	file_ts_proto_rawDescOnce.Do(func() {
		file_ts_proto_rawDescData = protoimpl.X.CompressGZIP(file_ts_proto_rawDescData)
	})
	return file_ts_proto_rawDescData
}

var file_ts_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_ts_proto_goTypes = []interface{}{
	(*TestRequest)(nil),  // 0: com.github.cappuccinotm.gcache.example.TestRequest
	(*TestResponse)(nil), // 1: com.github.cappuccinotm.gcache.example.TestResponse
}
var file_ts_proto_depIdxs = []int32{
	0, // 0: com.github.cappuccinotm.gcache.example.TestService.Test:input_type -> com.github.cappuccinotm.gcache.example.TestRequest
	1, // 1: com.github.cappuccinotm.gcache.example.TestService.Test:output_type -> com.github.cappuccinotm.gcache.example.TestResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_ts_proto_init() }
func file_ts_proto_init() {
	if File_ts_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_ts_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestRequest); i {
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
		file_ts_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestResponse); i {
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
			RawDescriptor: file_ts_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_ts_proto_goTypes,
		DependencyIndexes: file_ts_proto_depIdxs,
		MessageInfos:      file_ts_proto_msgTypes,
	}.Build()
	File_ts_proto = out.File
	file_ts_proto_rawDesc = nil
	file_ts_proto_goTypes = nil
	file_ts_proto_depIdxs = nil
}
