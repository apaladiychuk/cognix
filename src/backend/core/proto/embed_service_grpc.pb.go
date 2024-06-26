// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: embed_service.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// EmbedServiceClient is the client API for EmbedService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmbedServiceClient interface {
	GetEmbeding(ctx context.Context, in *EmbedRequest, opts ...grpc.CallOption) (*EmbedResponse, error)
}

type embedServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEmbedServiceClient(cc grpc.ClientConnInterface) EmbedServiceClient {
	return &embedServiceClient{cc}
}

func (c *embedServiceClient) GetEmbeding(ctx context.Context, in *EmbedRequest, opts ...grpc.CallOption) (*EmbedResponse, error) {
	out := new(EmbedResponse)
	err := c.cc.Invoke(ctx, "/com.embedd.EmbedService/GetEmbeding", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmbedServiceServer is the server API for EmbedService service.
// All implementations must embed UnimplementedEmbedServiceServer
// for forward compatibility
type EmbedServiceServer interface {
	GetEmbeding(context.Context, *EmbedRequest) (*EmbedResponse, error)
	mustEmbedUnimplementedEmbedServiceServer()
}

// UnimplementedEmbedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEmbedServiceServer struct {
}

func (UnimplementedEmbedServiceServer) GetEmbeding(context.Context, *EmbedRequest) (*EmbedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEmbeding not implemented")
}
func (UnimplementedEmbedServiceServer) mustEmbedUnimplementedEmbedServiceServer() {}

// UnsafeEmbedServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmbedServiceServer will
// result in compilation errors.
type UnsafeEmbedServiceServer interface {
	mustEmbedUnimplementedEmbedServiceServer()
}

func RegisterEmbedServiceServer(s grpc.ServiceRegistrar, srv EmbedServiceServer) {
	s.RegisterService(&EmbedService_ServiceDesc, srv)
}

func _EmbedService_GetEmbeding_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmbedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmbedServiceServer).GetEmbeding(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/com.cognix.EmbedService/GetEmbeding",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmbedServiceServer).GetEmbeding(ctx, req.(*EmbedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EmbedService_ServiceDesc is the grpc.ServiceDesc for EmbedService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EmbedService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "com.cognix.EmbedService",
	HandlerType: (*EmbedServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEmbeding",
			Handler:    _EmbedService_GetEmbeding_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "embed_service.proto",
}
