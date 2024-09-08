// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.2
// source: server/gpgkey/gpgkey.proto

// GPG public key service
//
// GPG public key API performs CRUD actions against GnuPGPublicKey resources

package gpgkey

import (
	context "context"
	v1alpha1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	GPGKeyService_List_FullMethodName   = "/gpgkey.GPGKeyService/List"
	GPGKeyService_Get_FullMethodName    = "/gpgkey.GPGKeyService/Get"
	GPGKeyService_Create_FullMethodName = "/gpgkey.GPGKeyService/Create"
	GPGKeyService_Delete_FullMethodName = "/gpgkey.GPGKeyService/Delete"
)

// GPGKeyServiceClient is the client API for GPGKeyService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// GPGKeyService implements API for managing GPG public keys on the server
type GPGKeyServiceClient interface {
	// List all available repository certificates
	List(ctx context.Context, in *GnuPGPublicKeyQuery, opts ...grpc.CallOption) (*v1alpha1.GnuPGPublicKeyList, error)
	// Get information about specified GPG public key from the server
	Get(ctx context.Context, in *GnuPGPublicKeyQuery, opts ...grpc.CallOption) (*v1alpha1.GnuPGPublicKey, error)
	// Create one or more GPG public keys in the server's configuration
	Create(ctx context.Context, in *GnuPGPublicKeyCreateRequest, opts ...grpc.CallOption) (*GnuPGPublicKeyCreateResponse, error)
	// Delete specified GPG public key from the server's configuration
	Delete(ctx context.Context, in *GnuPGPublicKeyQuery, opts ...grpc.CallOption) (*GnuPGPublicKeyResponse, error)
}

type gPGKeyServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewGPGKeyServiceClient(cc grpc.ClientConnInterface) GPGKeyServiceClient {
	return &gPGKeyServiceClient{cc}
}

func (c *gPGKeyServiceClient) List(ctx context.Context, in *GnuPGPublicKeyQuery, opts ...grpc.CallOption) (*v1alpha1.GnuPGPublicKeyList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(v1alpha1.GnuPGPublicKeyList)
	err := c.cc.Invoke(ctx, GPGKeyService_List_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gPGKeyServiceClient) Get(ctx context.Context, in *GnuPGPublicKeyQuery, opts ...grpc.CallOption) (*v1alpha1.GnuPGPublicKey, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(v1alpha1.GnuPGPublicKey)
	err := c.cc.Invoke(ctx, GPGKeyService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gPGKeyServiceClient) Create(ctx context.Context, in *GnuPGPublicKeyCreateRequest, opts ...grpc.CallOption) (*GnuPGPublicKeyCreateResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GnuPGPublicKeyCreateResponse)
	err := c.cc.Invoke(ctx, GPGKeyService_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gPGKeyServiceClient) Delete(ctx context.Context, in *GnuPGPublicKeyQuery, opts ...grpc.CallOption) (*GnuPGPublicKeyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GnuPGPublicKeyResponse)
	err := c.cc.Invoke(ctx, GPGKeyService_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GPGKeyServiceServer is the server API for GPGKeyService service.
// All implementations must embed UnimplementedGPGKeyServiceServer
// for forward compatibility.
//
// GPGKeyService implements API for managing GPG public keys on the server
type GPGKeyServiceServer interface {
	// List all available repository certificates
	List(context.Context, *GnuPGPublicKeyQuery) (*v1alpha1.GnuPGPublicKeyList, error)
	// Get information about specified GPG public key from the server
	Get(context.Context, *GnuPGPublicKeyQuery) (*v1alpha1.GnuPGPublicKey, error)
	// Create one or more GPG public keys in the server's configuration
	Create(context.Context, *GnuPGPublicKeyCreateRequest) (*GnuPGPublicKeyCreateResponse, error)
	// Delete specified GPG public key from the server's configuration
	Delete(context.Context, *GnuPGPublicKeyQuery) (*GnuPGPublicKeyResponse, error)
	mustEmbedUnimplementedGPGKeyServiceServer()
}

// UnimplementedGPGKeyServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGPGKeyServiceServer struct{}

func (UnimplementedGPGKeyServiceServer) List(context.Context, *GnuPGPublicKeyQuery) (*v1alpha1.GnuPGPublicKeyList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedGPGKeyServiceServer) Get(context.Context, *GnuPGPublicKeyQuery) (*v1alpha1.GnuPGPublicKey, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedGPGKeyServiceServer) Create(context.Context, *GnuPGPublicKeyCreateRequest) (*GnuPGPublicKeyCreateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedGPGKeyServiceServer) Delete(context.Context, *GnuPGPublicKeyQuery) (*GnuPGPublicKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedGPGKeyServiceServer) mustEmbedUnimplementedGPGKeyServiceServer() {}
func (UnimplementedGPGKeyServiceServer) testEmbeddedByValue()                       {}

// UnsafeGPGKeyServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GPGKeyServiceServer will
// result in compilation errors.
type UnsafeGPGKeyServiceServer interface {
	mustEmbedUnimplementedGPGKeyServiceServer()
}

func RegisterGPGKeyServiceServer(s grpc.ServiceRegistrar, srv GPGKeyServiceServer) {
	// If the following call pancis, it indicates UnimplementedGPGKeyServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GPGKeyService_ServiceDesc, srv)
}

func _GPGKeyService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GnuPGPublicKeyQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GPGKeyServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GPGKeyService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GPGKeyServiceServer).List(ctx, req.(*GnuPGPublicKeyQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _GPGKeyService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GnuPGPublicKeyQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GPGKeyServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GPGKeyService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GPGKeyServiceServer).Get(ctx, req.(*GnuPGPublicKeyQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _GPGKeyService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GnuPGPublicKeyCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GPGKeyServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GPGKeyService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GPGKeyServiceServer).Create(ctx, req.(*GnuPGPublicKeyCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _GPGKeyService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GnuPGPublicKeyQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GPGKeyServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GPGKeyService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GPGKeyServiceServer).Delete(ctx, req.(*GnuPGPublicKeyQuery))
	}
	return interceptor(ctx, in, info, handler)
}

// GPGKeyService_ServiceDesc is the grpc.ServiceDesc for GPGKeyService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GPGKeyService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "gpgkey.GPGKeyService",
	HandlerType: (*GPGKeyServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _GPGKeyService_List_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _GPGKeyService_Get_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _GPGKeyService_Create_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _GPGKeyService_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server/gpgkey/gpgkey.proto",
}
