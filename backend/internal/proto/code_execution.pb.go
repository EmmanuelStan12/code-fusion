// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        v3.19.6
// source: internal/proto/code_execution.proto

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

type ExecuteCodeRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Code          string                 `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	SessionId     string                 `protobuf:"bytes,2,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExecuteCodeRequest) Reset() {
	*x = ExecuteCodeRequest{}
	mi := &file_internal_proto_code_execution_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExecuteCodeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteCodeRequest) ProtoMessage() {}

func (x *ExecuteCodeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_code_execution_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteCodeRequest.ProtoReflect.Descriptor instead.
func (*ExecuteCodeRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_code_execution_proto_rawDescGZIP(), []int{0}
}

func (x *ExecuteCodeRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *ExecuteCodeRequest) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

type ExecuteCodeResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SessionId     string                 `protobuf:"bytes,1,opt,name=sessionId,proto3" json:"sessionId,omitempty"`
	ContextId     string                 `protobuf:"bytes,2,opt,name=contextId,proto3" json:"contextId,omitempty"`
	Result        string                 `protobuf:"bytes,3,opt,name=result,proto3" json:"result,omitempty"`
	Success       bool                   `protobuf:"varint,4,opt,name=success,proto3" json:"success,omitempty"`
	Stdout        string                 `protobuf:"bytes,5,opt,name=stdout,proto3" json:"stdout,omitempty"`
	Error         string                 `protobuf:"bytes,6,opt,name=error,proto3" json:"error,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ExecuteCodeResponse) Reset() {
	*x = ExecuteCodeResponse{}
	mi := &file_internal_proto_code_execution_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ExecuteCodeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ExecuteCodeResponse) ProtoMessage() {}

func (x *ExecuteCodeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_code_execution_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ExecuteCodeResponse.ProtoReflect.Descriptor instead.
func (*ExecuteCodeResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_code_execution_proto_rawDescGZIP(), []int{1}
}

func (x *ExecuteCodeResponse) GetSessionId() string {
	if x != nil {
		return x.SessionId
	}
	return ""
}

func (x *ExecuteCodeResponse) GetContextId() string {
	if x != nil {
		return x.ContextId
	}
	return ""
}

func (x *ExecuteCodeResponse) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *ExecuteCodeResponse) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *ExecuteCodeResponse) GetStdout() string {
	if x != nil {
		return x.Stdout
	}
	return ""
}

func (x *ExecuteCodeResponse) GetError() string {
	if x != nil {
		return x.Error
	}
	return ""
}

var File_internal_proto_code_execution_proto protoreflect.FileDescriptor

var file_internal_proto_code_execution_proto_rawDesc = []byte{
	0x0a, 0x23, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2f, 0x63, 0x6f, 0x64, 0x65, 0x5f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x46, 0x0a, 0x12, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65,
	0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63,
	0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12,
	0x1c, 0x0a, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0xb1, 0x01,
	0x0a, 0x13, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x1c, 0x0a, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x49, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49, 0x64,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x78, 0x74, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x65,
	0x72, 0x72, 0x6f, 0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x72, 0x72, 0x6f,
	0x72, 0x32, 0x56, 0x0a, 0x14, 0x43, 0x6f, 0x64, 0x65, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x3e, 0x0a, 0x0b, 0x45, 0x78, 0x65,
	0x63, 0x75, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x12, 0x13, 0x2e, 0x45, 0x78, 0x65, 0x63, 0x75,
	0x74, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x14, 0x2e,
	0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x43, 0x6f, 0x64, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x00, 0x28, 0x01, 0x30, 0x01, 0x42, 0x3e, 0x5a, 0x3c, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x45, 0x6d, 0x6d, 0x61, 0x6e, 0x75, 0x65, 0x6c,
	0x53, 0x74, 0x61, 0x6e, 0x31, 0x32, 0x2f, 0x63, 0x6f, 0x64, 0x65, 0x2d, 0x66, 0x75, 0x73, 0x69,
	0x6f, 0x6e, 0x2f, 0x62, 0x61, 0x63, 0x6b, 0x65, 0x6e, 0x64, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_internal_proto_code_execution_proto_rawDescOnce sync.Once
	file_internal_proto_code_execution_proto_rawDescData = file_internal_proto_code_execution_proto_rawDesc
)

func file_internal_proto_code_execution_proto_rawDescGZIP() []byte {
	file_internal_proto_code_execution_proto_rawDescOnce.Do(func() {
		file_internal_proto_code_execution_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_proto_code_execution_proto_rawDescData)
	})
	return file_internal_proto_code_execution_proto_rawDescData
}

var file_internal_proto_code_execution_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_internal_proto_code_execution_proto_goTypes = []any{
	(*ExecuteCodeRequest)(nil),  // 0: ExecuteCodeRequest
	(*ExecuteCodeResponse)(nil), // 1: ExecuteCodeResponse
}
var file_internal_proto_code_execution_proto_depIdxs = []int32{
	0, // 0: CodeExecutionService.ExecuteCode:input_type -> ExecuteCodeRequest
	1, // 1: CodeExecutionService.ExecuteCode:output_type -> ExecuteCodeResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_internal_proto_code_execution_proto_init() }
func file_internal_proto_code_execution_proto_init() {
	if File_internal_proto_code_execution_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_proto_code_execution_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_proto_code_execution_proto_goTypes,
		DependencyIndexes: file_internal_proto_code_execution_proto_depIdxs,
		MessageInfos:      file_internal_proto_code_execution_proto_msgTypes,
	}.Build()
	File_internal_proto_code_execution_proto = out.File
	file_internal_proto_code_execution_proto_rawDesc = nil
	file_internal_proto_code_execution_proto_goTypes = nil
	file_internal_proto_code_execution_proto_depIdxs = nil
}
