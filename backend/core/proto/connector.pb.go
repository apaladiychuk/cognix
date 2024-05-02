// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: core/proto/connector.proto

package proto

import (
	empty "github.com/golang/protobuf/ptypes/empty"
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

type ResponseStep int32

const (
	ResponseStep_NO_ACTION ResponseStep = 0
	ResponseStep_LOAD      ResponseStep = 1
	ResponseStep_EMBEDDING ResponseStep = 2
	ResponseStep_FINISH    ResponseStep = 3
)

// Enum value maps for ResponseStep.
var (
	ResponseStep_name = map[int32]string{
		0: "NO_ACTION",
		1: "LOAD",
		2: "EMBEDDING",
		3: "FINISH",
	}
	ResponseStep_value = map[string]int32{
		"NO_ACTION": 0,
		"LOAD":      1,
		"EMBEDDING": 2,
		"FINISH":    3,
	}
)

func (x ResponseStep) Enum() *ResponseStep {
	p := new(ResponseStep)
	*p = x
	return p
}

func (x ResponseStep) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResponseStep) Descriptor() protoreflect.EnumDescriptor {
	return file_core_proto_connector_proto_enumTypes[0].Descriptor()
}

func (ResponseStep) Type() protoreflect.EnumType {
	return &file_core_proto_connector_proto_enumTypes[0]
}

func (x ResponseStep) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResponseStep.Descriptor instead.
func (ResponseStep) EnumDescriptor() ([]byte, []int) {
	return file_core_proto_connector_proto_rawDescGZIP(), []int{0}
}

type ConnectorRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id int64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ConnectorRequest) Reset() {
	*x = ConnectorRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_connector_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectorRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectorRequest) ProtoMessage() {}

func (x *ConnectorRequest) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_connector_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectorRequest.ProtoReflect.Descriptor instead.
func (*ConnectorRequest) Descriptor() ([]byte, []int) {
	return file_core_proto_connector_proto_rawDescGZIP(), []int{0}
}

func (x *ConnectorRequest) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

type ConnectorStepResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id         int64        `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DocumentId int64        `protobuf:"varint,2,opt,name=document_id,json=documentId,proto3" json:"document_id,omitempty"`
	Step       ResponseStep `protobuf:"varint,3,opt,name=step,proto3,enum=proto.ResponseStep" json:"step,omitempty"`
	Content    string       `protobuf:"bytes,4,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *ConnectorStepResponse) Reset() {
	*x = ConnectorStepResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_core_proto_connector_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConnectorStepResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConnectorStepResponse) ProtoMessage() {}

func (x *ConnectorStepResponse) ProtoReflect() protoreflect.Message {
	mi := &file_core_proto_connector_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConnectorStepResponse.ProtoReflect.Descriptor instead.
func (*ConnectorStepResponse) Descriptor() ([]byte, []int) {
	return file_core_proto_connector_proto_rawDescGZIP(), []int{1}
}

func (x *ConnectorStepResponse) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ConnectorStepResponse) GetDocumentId() int64 {
	if x != nil {
		return x.DocumentId
	}
	return 0
}

func (x *ConnectorStepResponse) GetStep() ResponseStep {
	if x != nil {
		return x.Step
	}
	return ResponseStep_NO_ACTION
}

func (x *ConnectorStepResponse) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

var File_core_proto_connector_proto protoreflect.FileDescriptor

var file_core_proto_connector_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x63, 0x6f, 0x6e,
	0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0x22, 0x0a, 0x10, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x02, 0x69, 0x64, 0x22, 0x8b, 0x01, 0x0a, 0x15, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x53, 0x74, 0x65, 0x70, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1f,
	0x0a, 0x0b, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x0a, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12,
	0x27, 0x0a, 0x04, 0x73, 0x74, 0x65, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x13, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74,
	0x65, 0x70, 0x52, 0x04, 0x73, 0x74, 0x65, 0x70, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65,
	0x6e, 0x74, 0x2a, 0x42, 0x0a, 0x0c, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74,
	0x65, 0x70, 0x12, 0x0d, 0x0a, 0x09, 0x4e, 0x4f, 0x5f, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10,
	0x00, 0x12, 0x08, 0x0a, 0x04, 0x4c, 0x4f, 0x41, 0x44, 0x10, 0x01, 0x12, 0x0d, 0x0a, 0x09, 0x45,
	0x4d, 0x42, 0x45, 0x44, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x12, 0x0a, 0x0a, 0x06, 0x46, 0x49,
	0x4e, 0x49, 0x53, 0x48, 0x10, 0x03, 0x32, 0x44, 0x0a, 0x09, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63,
	0x74, 0x6f, 0x72, 0x12, 0x37, 0x0a, 0x03, 0x52, 0x75, 0x6e, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x42, 0x12, 0x5a, 0x10,
	0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_core_proto_connector_proto_rawDescOnce sync.Once
	file_core_proto_connector_proto_rawDescData = file_core_proto_connector_proto_rawDesc
)

func file_core_proto_connector_proto_rawDescGZIP() []byte {
	file_core_proto_connector_proto_rawDescOnce.Do(func() {
		file_core_proto_connector_proto_rawDescData = protoimpl.X.CompressGZIP(file_core_proto_connector_proto_rawDescData)
	})
	return file_core_proto_connector_proto_rawDescData
}

var file_core_proto_connector_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_core_proto_connector_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_core_proto_connector_proto_goTypes = []interface{}{
	(ResponseStep)(0),             // 0: proto.ResponseStep
	(*ConnectorRequest)(nil),      // 1: proto.ConnectorRequest
	(*ConnectorStepResponse)(nil), // 2: proto.ConnectorStepResponse
	(*empty.Empty)(nil),           // 3: google.protobuf.Empty
}
var file_core_proto_connector_proto_depIdxs = []int32{
	0, // 0: proto.ConnectorStepResponse.step:type_name -> proto.ResponseStep
	3, // 1: proto.Connector.Run:input_type -> google.protobuf.Empty
	3, // 2: proto.Connector.Run:output_type -> google.protobuf.Empty
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_core_proto_connector_proto_init() }
func file_core_proto_connector_proto_init() {
	if File_core_proto_connector_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_core_proto_connector_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectorRequest); i {
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
		file_core_proto_connector_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ConnectorStepResponse); i {
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
			RawDescriptor: file_core_proto_connector_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_core_proto_connector_proto_goTypes,
		DependencyIndexes: file_core_proto_connector_proto_depIdxs,
		EnumInfos:         file_core_proto_connector_proto_enumTypes,
		MessageInfos:      file_core_proto_connector_proto_msgTypes,
	}.Build()
	File_core_proto_connector_proto = out.File
	file_core_proto_connector_proto_rawDesc = nil
	file_core_proto_connector_proto_goTypes = nil
	file_core_proto_connector_proto_depIdxs = nil
}
