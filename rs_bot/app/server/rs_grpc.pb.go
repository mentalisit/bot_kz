// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.0--rc2
// source: rs.proto

package server

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
	LogicService_LogicRs_FullMethodName = "/rs.LogicService/LogicRs"
)

// LogicServiceClient is the client API for LogicService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Определение сервиса и функции LogicRs
type LogicServiceClient interface {
	LogicRs(ctx context.Context, in *InMessage, opts ...grpc.CallOption) (*Empty, error)
}

type logicServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLogicServiceClient(cc grpc.ClientConnInterface) LogicServiceClient {
	return &logicServiceClient{cc}
}

func (c *logicServiceClient) LogicRs(ctx context.Context, in *InMessage, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, LogicService_LogicRs_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LogicServiceServer is the server API for LogicService service.
// All implementations must embed UnimplementedLogicServiceServer
// for forward compatibility.
//
// Определение сервиса и функции LogicRs
type LogicServiceServer interface {
	LogicRs(context.Context, *InMessage) (*Empty, error)
	mustEmbedUnimplementedLogicServiceServer()
}

// UnimplementedLogicServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedLogicServiceServer struct{}

func (UnimplementedLogicServiceServer) LogicRs(context.Context, *InMessage) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogicRs not implemented")
}
func (UnimplementedLogicServiceServer) mustEmbedUnimplementedLogicServiceServer() {}
func (UnimplementedLogicServiceServer) testEmbeddedByValue()                      {}

// UnsafeLogicServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LogicServiceServer will
// result in compilation errors.
type UnsafeLogicServiceServer interface {
	mustEmbedUnimplementedLogicServiceServer()
}

func RegisterLogicServiceServer(s grpc.ServiceRegistrar, srv LogicServiceServer) {
	// If the following call pancis, it indicates UnimplementedLogicServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&LogicService_ServiceDesc, srv)
}

func _LogicService_LogicRs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LogicServiceServer).LogicRs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LogicService_LogicRs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LogicServiceServer).LogicRs(ctx, req.(*InMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// LogicService_ServiceDesc is the grpc.ServiceDesc for LogicService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LogicService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "rs.LogicService",
	HandlerType: (*LogicServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LogicRs",
			Handler:    _LogicService_LogicRs_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rs.proto",
}
