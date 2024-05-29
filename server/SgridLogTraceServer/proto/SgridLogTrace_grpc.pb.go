// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: SgridLogTrace.proto

package protocol

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SgridLogTraceServiceClient is the client API for SgridLogTraceService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SgridLogTraceServiceClient interface {
	LogTrace(ctx context.Context, in *LogTraceReq, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type sgridLogTraceServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSgridLogTraceServiceClient(cc grpc.ClientConnInterface) SgridLogTraceServiceClient {
	return &sgridLogTraceServiceClient{cc}
}

func (c *sgridLogTraceServiceClient) LogTrace(ctx context.Context, in *LogTraceReq, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/SgridLogTrace.SgridLogTraceService/LogTrace", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SgridLogTraceServiceServer is the server API for SgridLogTraceService service.
// All implementations must embed UnimplementedSgridLogTraceServiceServer
// for forward compatibility
type SgridLogTraceServiceServer interface {
	LogTrace(context.Context, *LogTraceReq) (*emptypb.Empty, error)
	mustEmbedUnimplementedSgridLogTraceServiceServer()
}

// UnimplementedSgridLogTraceServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSgridLogTraceServiceServer struct {
}

func (UnimplementedSgridLogTraceServiceServer) LogTrace(context.Context, *LogTraceReq) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LogTrace not implemented")
}
func (UnimplementedSgridLogTraceServiceServer) mustEmbedUnimplementedSgridLogTraceServiceServer() {}

// UnsafeSgridLogTraceServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SgridLogTraceServiceServer will
// result in compilation errors.
type UnsafeSgridLogTraceServiceServer interface {
	mustEmbedUnimplementedSgridLogTraceServiceServer()
}

func RegisterSgridLogTraceServiceServer(s grpc.ServiceRegistrar, srv SgridLogTraceServiceServer) {
	s.RegisterService(&SgridLogTraceService_ServiceDesc, srv)
}

func _SgridLogTraceService_LogTrace_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogTraceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SgridLogTraceServiceServer).LogTrace(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridLogTrace.SgridLogTraceService/LogTrace",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SgridLogTraceServiceServer).LogTrace(ctx, req.(*LogTraceReq))
	}
	return interceptor(ctx, in, info, handler)
}

// SgridLogTraceService_ServiceDesc is the grpc.ServiceDesc for SgridLogTraceService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SgridLogTraceService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SgridLogTrace.SgridLogTraceService",
	HandlerType: (*SgridLogTraceServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LogTrace",
			Handler:    _SgridLogTraceService_LogTrace_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "SgridLogTrace.proto",
}
