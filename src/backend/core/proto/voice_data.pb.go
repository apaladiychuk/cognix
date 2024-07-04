// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.12.4
// source: voice_data.proto

package proto

import (
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

type VoiceData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// This is the url where the file is located.
	// Based on the chunking type it will be a WEB URL (HTML type)
	// Will be an S3/MINIO link with a proper authentication in case of a file
	Url            string   `protobuf:"bytes,1,opt,name=url,proto3" json:"url,omitempty"`
	DocumentId     int64    `protobuf:"varint,2,opt,name=document_id,json=documentId,proto3" json:"document_id,omitempty"`
	ConnectorId    int64    `protobuf:"varint,3,opt,name=connector_id,json=connectorId,proto3" json:"connector_id,omitempty"`
	FileType       FileType `protobuf:"varint,4,opt,name=file_type,json=fileType,proto3,enum=com.cognix.FileType" json:"file_type,omitempty"`
	CollectionName string   `protobuf:"bytes,5,opt,name=collection_name,json=collectionName,proto3" json:"collection_name,omitempty"`
	ModelName      string   `protobuf:"bytes,6,opt,name=model_name,json=modelName,proto3" json:"model_name,omitempty"`
	ModelDimension int32    `protobuf:"varint,7,opt,name=model_dimension,json=modelDimension,proto3" json:"model_dimension,omitempty"`
}

func (x *VoiceData) Reset() {
	*x = VoiceData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_voice_data_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoiceData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoiceData) ProtoMessage() {}

func (x *VoiceData) ProtoReflect() protoreflect.Message {
	mi := &file_voice_data_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoiceData.ProtoReflect.Descriptor instead.
func (*VoiceData) Descriptor() ([]byte, []int) {
	return file_voice_data_proto_rawDescGZIP(), []int{0}
}

func (x *VoiceData) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *VoiceData) GetDocumentId() int64 {
	if x != nil {
		return x.DocumentId
	}
	return 0
}

func (x *VoiceData) GetConnectorId() int64 {
	if x != nil {
		return x.ConnectorId
	}
	return 0
}

func (x *VoiceData) GetFileType() FileType {
	if x != nil {
		return x.FileType
	}
	return FileType_UNKNOWN
}

func (x *VoiceData) GetCollectionName() string {
	if x != nil {
		return x.CollectionName
	}
	return ""
}

func (x *VoiceData) GetModelName() string {
	if x != nil {
		return x.ModelName
	}
	return ""
}

func (x *VoiceData) GetModelDimension() int32 {
	if x != nil {
		return x.ModelDimension
	}
	return 0
}

var File_voice_data_proto protoreflect.FileDescriptor

var file_voice_data_proto_rawDesc = []byte{
	0x0a, 0x10, 0x76, 0x6f, 0x69, 0x63, 0x65, 0x5f, 0x64, 0x61, 0x74, 0x61, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x67, 0x6e, 0x69, 0x78, 0x1a, 0x0f,
	0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x85, 0x02, 0x0a, 0x09, 0x56, 0x6f, 0x69, 0x63, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x10, 0x0a,
	0x03, 0x75, 0x72, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12,
	0x1f, 0x0a, 0x0b, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0a, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x49, 0x64,
	0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0b, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x49, 0x64, 0x12, 0x31, 0x0a, 0x09, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x67,
	0x6e, 0x69, 0x78, 0x2e, 0x46, 0x69, 0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x52, 0x08, 0x66, 0x69,
	0x6c, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0e, 0x63, 0x6f, 0x6c, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x1d, 0x0a, 0x0a, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x27,
	0x0a, 0x0f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x5f, 0x64, 0x69, 0x6d, 0x65, 0x6e, 0x73, 0x69, 0x6f,
	0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0e, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x44, 0x69,
	0x6d, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x42, 0x1a, 0x5a, 0x18, 0x62, 0x61, 0x63, 0x6b, 0x65,
	0x6e, 0x64, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x3b, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_voice_data_proto_rawDescOnce sync.Once
	file_voice_data_proto_rawDescData = file_voice_data_proto_rawDesc
)

func file_voice_data_proto_rawDescGZIP() []byte {
	file_voice_data_proto_rawDescOnce.Do(func() {
		file_voice_data_proto_rawDescData = protoimpl.X.CompressGZIP(file_voice_data_proto_rawDescData)
	})
	return file_voice_data_proto_rawDescData
}

var file_voice_data_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_voice_data_proto_goTypes = []interface{}{
	(*VoiceData)(nil), // 0: com.cognix.VoiceData
	(FileType)(0),     // 1: com.cognix.FileType
}
var file_voice_data_proto_depIdxs = []int32{
	1, // 0: com.cognix.VoiceData.file_type:type_name -> com.cognix.FileType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_voice_data_proto_init() }
func file_voice_data_proto_init() {
	if File_voice_data_proto != nil {
		return
	}
	file_file_type_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_voice_data_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoiceData); i {
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
			RawDescriptor: file_voice_data_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_voice_data_proto_goTypes,
		DependencyIndexes: file_voice_data_proto_depIdxs,
		MessageInfos:      file_voice_data_proto_msgTypes,
	}.Build()
	File_voice_data_proto = out.File
	file_voice_data_proto_rawDesc = nil
	file_voice_data_proto_goTypes = nil
	file_voice_data_proto_depIdxs = nil
}
