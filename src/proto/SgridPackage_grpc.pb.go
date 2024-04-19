// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.25.1
// source: SgridPackage.proto

package protocol

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

// FileTransferServiceClient is the client API for FileTransferService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FileTransferServiceClient interface {
	// 双向流式 RPC 用于文件传输
	StreamFile(ctx context.Context, opts ...grpc.CallOption) (FileTransferService_StreamFileClient, error)
	// 删除包
	DeletePackage(ctx context.Context, in *DeletePackageReq, opts ...grpc.CallOption) (*BasicResp, error)
	// 发布服务
	ReleaseServerByPackage(ctx context.Context, in *ReleaseServerReq, opts ...grpc.CallOption) (*BasicResp, error)
	// 关闭指定节点服务
	ShutdownGrid(ctx context.Context, in *ShutdownGridReq, opts ...grpc.CallOption) (*BasicResp, error)
	// 获取服务日志列表
	GetLogFileByHost(ctx context.Context, in *GetLogFileByHostReq, opts ...grpc.CallOption) (*GetLogFileByHostResp, error)
	// 获取服务日志
	GetLogByFile(ctx context.Context, in *GetLogByFileReq, opts ...grpc.CallOption) (*GetLogByFileResp, error)
	// 获取Pid信息 checkAlive
	GetPidInfo(ctx context.Context, in *GetPidInfoReq, opts ...grpc.CallOption) (*GetPidInfoResp, error)
}

type fileTransferServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewFileTransferServiceClient(cc grpc.ClientConnInterface) FileTransferServiceClient {
	return &fileTransferServiceClient{cc}
}

func (c *fileTransferServiceClient) StreamFile(ctx context.Context, opts ...grpc.CallOption) (FileTransferService_StreamFileClient, error) {
	stream, err := c.cc.NewStream(ctx, &FileTransferService_ServiceDesc.Streams[0], "/SgridProtocol.FileTransferService/StreamFile", opts...)
	if err != nil {
		return nil, err
	}
	x := &fileTransferServiceStreamFileClient{stream}
	return x, nil
}

type FileTransferService_StreamFileClient interface {
	Send(*FileChunk) error
	CloseAndRecv() (*FileResp, error)
	grpc.ClientStream
}

type fileTransferServiceStreamFileClient struct {
	grpc.ClientStream
}

func (x *fileTransferServiceStreamFileClient) Send(m *FileChunk) error {
	return x.ClientStream.SendMsg(m)
}

func (x *fileTransferServiceStreamFileClient) CloseAndRecv() (*FileResp, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(FileResp)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *fileTransferServiceClient) DeletePackage(ctx context.Context, in *DeletePackageReq, opts ...grpc.CallOption) (*BasicResp, error) {
	out := new(BasicResp)
	err := c.cc.Invoke(ctx, "/SgridProtocol.FileTransferService/DeletePackage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferServiceClient) ReleaseServerByPackage(ctx context.Context, in *ReleaseServerReq, opts ...grpc.CallOption) (*BasicResp, error) {
	out := new(BasicResp)
	err := c.cc.Invoke(ctx, "/SgridProtocol.FileTransferService/ReleaseServerByPackage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferServiceClient) ShutdownGrid(ctx context.Context, in *ShutdownGridReq, opts ...grpc.CallOption) (*BasicResp, error) {
	out := new(BasicResp)
	err := c.cc.Invoke(ctx, "/SgridProtocol.FileTransferService/ShutdownGrid", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferServiceClient) GetLogFileByHost(ctx context.Context, in *GetLogFileByHostReq, opts ...grpc.CallOption) (*GetLogFileByHostResp, error) {
	out := new(GetLogFileByHostResp)
	err := c.cc.Invoke(ctx, "/SgridProtocol.FileTransferService/GetLogFileByHost", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferServiceClient) GetLogByFile(ctx context.Context, in *GetLogByFileReq, opts ...grpc.CallOption) (*GetLogByFileResp, error) {
	out := new(GetLogByFileResp)
	err := c.cc.Invoke(ctx, "/SgridProtocol.FileTransferService/GetLogByFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *fileTransferServiceClient) GetPidInfo(ctx context.Context, in *GetPidInfoReq, opts ...grpc.CallOption) (*GetPidInfoResp, error) {
	out := new(GetPidInfoResp)
	err := c.cc.Invoke(ctx, "/SgridProtocol.FileTransferService/GetPidInfo", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FileTransferServiceServer is the server API for FileTransferService service.
// All implementations must embed UnimplementedFileTransferServiceServer
// for forward compatibility
type FileTransferServiceServer interface {
	// 双向流式 RPC 用于文件传输
	StreamFile(FileTransferService_StreamFileServer) error
	// 删除包
	DeletePackage(context.Context, *DeletePackageReq) (*BasicResp, error)
	// 发布服务
	ReleaseServerByPackage(context.Context, *ReleaseServerReq) (*BasicResp, error)
	// 关闭指定节点服务
	ShutdownGrid(context.Context, *ShutdownGridReq) (*BasicResp, error)
	// 获取服务日志列表
	GetLogFileByHost(context.Context, *GetLogFileByHostReq) (*GetLogFileByHostResp, error)
	// 获取服务日志
	GetLogByFile(context.Context, *GetLogByFileReq) (*GetLogByFileResp, error)
	// 获取Pid信息 checkAlive
	GetPidInfo(context.Context, *GetPidInfoReq) (*GetPidInfoResp, error)
	mustEmbedUnimplementedFileTransferServiceServer()
}

// UnimplementedFileTransferServiceServer must be embedded to have forward compatible implementations.
type UnimplementedFileTransferServiceServer struct {
}

func (UnimplementedFileTransferServiceServer) StreamFile(FileTransferService_StreamFileServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamFile not implemented")
}
func (UnimplementedFileTransferServiceServer) DeletePackage(context.Context, *DeletePackageReq) (*BasicResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeletePackage not implemented")
}
func (UnimplementedFileTransferServiceServer) ReleaseServerByPackage(context.Context, *ReleaseServerReq) (*BasicResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReleaseServerByPackage not implemented")
}
func (UnimplementedFileTransferServiceServer) ShutdownGrid(context.Context, *ShutdownGridReq) (*BasicResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ShutdownGrid not implemented")
}
func (UnimplementedFileTransferServiceServer) GetLogFileByHost(context.Context, *GetLogFileByHostReq) (*GetLogFileByHostResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLogFileByHost not implemented")
}
func (UnimplementedFileTransferServiceServer) GetLogByFile(context.Context, *GetLogByFileReq) (*GetLogByFileResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLogByFile not implemented")
}
func (UnimplementedFileTransferServiceServer) GetPidInfo(context.Context, *GetPidInfoReq) (*GetPidInfoResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPidInfo not implemented")
}
func (UnimplementedFileTransferServiceServer) mustEmbedUnimplementedFileTransferServiceServer() {}

// UnsafeFileTransferServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FileTransferServiceServer will
// result in compilation errors.
type UnsafeFileTransferServiceServer interface {
	mustEmbedUnimplementedFileTransferServiceServer()
}

func RegisterFileTransferServiceServer(s grpc.ServiceRegistrar, srv FileTransferServiceServer) {
	s.RegisterService(&FileTransferService_ServiceDesc, srv)
}

func _FileTransferService_StreamFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(FileTransferServiceServer).StreamFile(&fileTransferServiceStreamFileServer{stream})
}

type FileTransferService_StreamFileServer interface {
	SendAndClose(*FileResp) error
	Recv() (*FileChunk, error)
	grpc.ServerStream
}

type fileTransferServiceStreamFileServer struct {
	grpc.ServerStream
}

func (x *fileTransferServiceStreamFileServer) SendAndClose(m *FileResp) error {
	return x.ServerStream.SendMsg(m)
}

func (x *fileTransferServiceStreamFileServer) Recv() (*FileChunk, error) {
	m := new(FileChunk)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _FileTransferService_DeletePackage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeletePackageReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServiceServer).DeletePackage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.FileTransferService/DeletePackage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServiceServer).DeletePackage(ctx, req.(*DeletePackageReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransferService_ReleaseServerByPackage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReleaseServerReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServiceServer).ReleaseServerByPackage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.FileTransferService/ReleaseServerByPackage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServiceServer).ReleaseServerByPackage(ctx, req.(*ReleaseServerReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransferService_ShutdownGrid_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ShutdownGridReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServiceServer).ShutdownGrid(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.FileTransferService/ShutdownGrid",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServiceServer).ShutdownGrid(ctx, req.(*ShutdownGridReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransferService_GetLogFileByHost_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLogFileByHostReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServiceServer).GetLogFileByHost(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.FileTransferService/GetLogFileByHost",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServiceServer).GetLogFileByHost(ctx, req.(*GetLogFileByHostReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransferService_GetLogByFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLogByFileReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServiceServer).GetLogByFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.FileTransferService/GetLogByFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServiceServer).GetLogByFile(ctx, req.(*GetLogByFileReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _FileTransferService_GetPidInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPidInfoReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FileTransferServiceServer).GetPidInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SgridProtocol.FileTransferService/GetPidInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FileTransferServiceServer).GetPidInfo(ctx, req.(*GetPidInfoReq))
	}
	return interceptor(ctx, in, info, handler)
}

// FileTransferService_ServiceDesc is the grpc.ServiceDesc for FileTransferService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var FileTransferService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SgridProtocol.FileTransferService",
	HandlerType: (*FileTransferServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeletePackage",
			Handler:    _FileTransferService_DeletePackage_Handler,
		},
		{
			MethodName: "ReleaseServerByPackage",
			Handler:    _FileTransferService_ReleaseServerByPackage_Handler,
		},
		{
			MethodName: "ShutdownGrid",
			Handler:    _FileTransferService_ShutdownGrid_Handler,
		},
		{
			MethodName: "GetLogFileByHost",
			Handler:    _FileTransferService_GetLogFileByHost_Handler,
		},
		{
			MethodName: "GetLogByFile",
			Handler:    _FileTransferService_GetLogByFile_Handler,
		},
		{
			MethodName: "GetPidInfo",
			Handler:    _FileTransferService_GetPidInfo_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamFile",
			Handler:       _FileTransferService_StreamFile_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "SgridPackage.proto",
}
