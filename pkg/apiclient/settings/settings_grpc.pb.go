// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.2
// source: server/settings/settings.proto

// Settings Service
//
// Settings Service API retrieves Argo CD settings

package settings

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
	SettingsService_Get_FullMethodName        = "/cluster.SettingsService/Get"
	SettingsService_GetPlugins_FullMethodName = "/cluster.SettingsService/GetPlugins"
)

// SettingsServiceClient is the client API for SettingsService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// SettingsService
type SettingsServiceClient interface {
	// Get returns Argo CD settings
	Get(ctx context.Context, in *SettingsQuery, opts ...grpc.CallOption) (*Settings, error)
	// Get returns Argo CD plugins
	GetPlugins(ctx context.Context, in *SettingsQuery, opts ...grpc.CallOption) (*SettingsPluginsResponse, error)
}

type settingsServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSettingsServiceClient(cc grpc.ClientConnInterface) SettingsServiceClient {
	return &settingsServiceClient{cc}
}

func (c *settingsServiceClient) Get(ctx context.Context, in *SettingsQuery, opts ...grpc.CallOption) (*Settings, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Settings)
	err := c.cc.Invoke(ctx, SettingsService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *settingsServiceClient) GetPlugins(ctx context.Context, in *SettingsQuery, opts ...grpc.CallOption) (*SettingsPluginsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SettingsPluginsResponse)
	err := c.cc.Invoke(ctx, SettingsService_GetPlugins_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SettingsServiceServer is the server API for SettingsService service.
// All implementations must embed UnimplementedSettingsServiceServer
// for forward compatibility.
//
// SettingsService
type SettingsServiceServer interface {
	// Get returns Argo CD settings
	Get(context.Context, *SettingsQuery) (*Settings, error)
	// Get returns Argo CD plugins
	GetPlugins(context.Context, *SettingsQuery) (*SettingsPluginsResponse, error)
	mustEmbedUnimplementedSettingsServiceServer()
}

// UnimplementedSettingsServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedSettingsServiceServer struct{}

func (UnimplementedSettingsServiceServer) Get(context.Context, *SettingsQuery) (*Settings, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedSettingsServiceServer) GetPlugins(context.Context, *SettingsQuery) (*SettingsPluginsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlugins not implemented")
}
func (UnimplementedSettingsServiceServer) mustEmbedUnimplementedSettingsServiceServer() {}
func (UnimplementedSettingsServiceServer) testEmbeddedByValue()                         {}

// UnsafeSettingsServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SettingsServiceServer will
// result in compilation errors.
type UnsafeSettingsServiceServer interface {
	mustEmbedUnimplementedSettingsServiceServer()
}

func RegisterSettingsServiceServer(s grpc.ServiceRegistrar, srv SettingsServiceServer) {
	// If the following call pancis, it indicates UnimplementedSettingsServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&SettingsService_ServiceDesc, srv)
}

func _SettingsService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SettingsQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingsServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SettingsService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingsServiceServer).Get(ctx, req.(*SettingsQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _SettingsService_GetPlugins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SettingsQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SettingsServiceServer).GetPlugins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: SettingsService_GetPlugins_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SettingsServiceServer).GetPlugins(ctx, req.(*SettingsQuery))
	}
	return interceptor(ctx, in, info, handler)
}

// SettingsService_ServiceDesc is the grpc.ServiceDesc for SettingsService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SettingsService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "cluster.SettingsService",
	HandlerType: (*SettingsServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Get",
			Handler:    _SettingsService_Get_Handler,
		},
		{
			MethodName: "GetPlugins",
			Handler:    _SettingsService_GetPlugins_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server/settings/settings.proto",
}
