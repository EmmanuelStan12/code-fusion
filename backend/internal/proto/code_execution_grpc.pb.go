// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.6
// source: internal/proto/code_execution.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	CodeExecutionService_ExecuteCode_FullMethodName = "/CodeExecutionService/ExecuteCode"
)

// CodeExecutionServiceClient is the client API for CodeExecutionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CodeExecutionServiceClient interface {
	ExecuteCode(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[ExecuteCodeRequest, ExecuteCodeResponse], error)
}

type codeExecutionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCodeExecutionServiceClient(cc grpc.ClientConnInterface) CodeExecutionServiceClient {
	return &codeExecutionServiceClient{cc}
}

func (c *codeExecutionServiceClient) ExecuteCode(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[ExecuteCodeRequest, ExecuteCodeResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &CodeExecutionService_ServiceDesc.Streams[0], CodeExecutionService_ExecuteCode_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ExecuteCodeRequest, ExecuteCodeResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type CodeExecutionService_ExecuteCodeClient = grpc.BidiStreamingClient[ExecuteCodeRequest, ExecuteCodeResponse]

// CodeExecutionServiceServer is the server API for CodeExecutionService service.
// All implementations must embed UnimplementedCodeExecutionServiceServer
// for forward compatibility.
type CodeExecutionServiceServer interface {
	ExecuteCode(grpc.BidiStreamingServer[ExecuteCodeRequest, ExecuteCodeResponse]) error
	mustEmbedUnimplementedCodeExecutionServiceServer()
}

// UnimplementedCodeExecutionServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCodeExecutionServiceServer struct{}

func (UnimplementedCodeExecutionServiceServer) ExecuteCode(grpc.BidiStreamingServer[ExecuteCodeRequest, ExecuteCodeResponse]) error {
	return status.Errorf(codes.Unimplemented, "method ExecuteCode not implemented")
}
func (UnimplementedCodeExecutionServiceServer) mustEmbedUnimplementedCodeExecutionServiceServer() {}
func (UnimplementedCodeExecutionServiceServer) testEmbeddedByValue()                              {}

// UnsafeCodeExecutionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CodeExecutionServiceServer will
// result in compilation errors.
type UnsafeCodeExecutionServiceServer interface {
	mustEmbedUnimplementedCodeExecutionServiceServer()
}

func RegisterCodeExecutionServiceServer(s grpc.ServiceRegistrar, srv CodeExecutionServiceServer) {
	// If the following call pancis, it indicates UnimplementedCodeExecutionServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CodeExecutionService_ServiceDesc, srv)
}

func _CodeExecutionService_ExecuteCode_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CodeExecutionServiceServer).ExecuteCode(&grpc.GenericServerStream[ExecuteCodeRequest, ExecuteCodeResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type CodeExecutionService_ExecuteCodeServer = grpc.BidiStreamingServer[ExecuteCodeRequest, ExecuteCodeResponse]

// CodeExecutionService_ServiceDesc is the grpc.ServiceDesc for CodeExecutionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CodeExecutionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "CodeExecutionService",
	HandlerType: (*CodeExecutionServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ExecuteCode",
			Handler:       _CodeExecutionService_ExecuteCode_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "internal/proto/code_execution.proto",
}