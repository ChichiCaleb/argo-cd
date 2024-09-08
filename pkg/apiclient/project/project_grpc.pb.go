// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.2
// source: server/project/project.proto

// Project Service
//
// Project Service API performs CRUD actions against project resources

package project

import (
	context "context"
	application "github.com/argoproj/argo-cd/v2/pkg/apiclient/application"
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
	ProjectService_CreateToken_FullMethodName         = "/project.ProjectService/CreateToken"
	ProjectService_DeleteToken_FullMethodName         = "/project.ProjectService/DeleteToken"
	ProjectService_Create_FullMethodName              = "/project.ProjectService/Create"
	ProjectService_List_FullMethodName                = "/project.ProjectService/List"
	ProjectService_GetDetailedProject_FullMethodName  = "/project.ProjectService/GetDetailedProject"
	ProjectService_Get_FullMethodName                 = "/project.ProjectService/Get"
	ProjectService_GetGlobalProjects_FullMethodName   = "/project.ProjectService/GetGlobalProjects"
	ProjectService_Update_FullMethodName              = "/project.ProjectService/Update"
	ProjectService_Delete_FullMethodName              = "/project.ProjectService/Delete"
	ProjectService_ListEvents_FullMethodName          = "/project.ProjectService/ListEvents"
	ProjectService_GetSyncWindowsState_FullMethodName = "/project.ProjectService/GetSyncWindowsState"
	ProjectService_ListLinks_FullMethodName           = "/project.ProjectService/ListLinks"
)

// ProjectServiceClient is the client API for ProjectService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// ProjectService
type ProjectServiceClient interface {
	// Create a new project token
	CreateToken(ctx context.Context, in *ProjectTokenCreateRequest, opts ...grpc.CallOption) (*ProjectTokenResponse, error)
	// Delete a new project token
	DeleteToken(ctx context.Context, in *ProjectTokenDeleteRequest, opts ...grpc.CallOption) (*EmptyResponse, error)
	// Create a new project
	Create(ctx context.Context, in *ProjectCreateRequest, opts ...grpc.CallOption) (*v1alpha1.AppProject, error)
	// List returns list of projects
	List(ctx context.Context, in *ProjectQuery, opts ...grpc.CallOption) (*v1alpha1.AppProjectList, error)
	// GetDetailedProject returns a project that include project, global project and scoped resources by name
	GetDetailedProject(ctx context.Context, in *ProjectQuery, opts ...grpc.CallOption) (*DetailedProjectsResponse, error)
	// Get returns a project by name
	Get(ctx context.Context, in *ProjectQuery, opts ...grpc.CallOption) (*v1alpha1.AppProject, error)
	// Get returns a virtual project by name
	GetGlobalProjects(ctx context.Context, in *ProjectQuery, opts ...grpc.CallOption) (*GlobalProjectsResponse, error)
	// Update updates a project
	Update(ctx context.Context, in *ProjectUpdateRequest, opts ...grpc.CallOption) (*v1alpha1.AppProject, error)
	// Delete deletes a project
	Delete(ctx context.Context, in *ProjectQuery, opts ...grpc.CallOption) (*EmptyResponse, error)
	// ListEvents returns a list of project events
	ListEvents(ctx context.Context, in *ProjectQuery, opts ...grpc.CallOption) (*EventListWrapper, error)
	// GetSchedulesState returns true if there are any active sync syncWindows
	GetSyncWindowsState(ctx context.Context, in *SyncWindowsQuery, opts ...grpc.CallOption) (*SyncWindowsResponse, error)
	// ListLinks returns all deep links for the particular project
	ListLinks(ctx context.Context, in *ListProjectLinksRequest, opts ...grpc.CallOption) (*application.LinksResponse, error)
}

type projectServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewProjectServiceClient(cc grpc.ClientConnInterface) ProjectServiceClient {
	return &projectServiceClient{cc}
}

func (c *projectServiceClient) CreateToken(ctx context.Context, in *ProjectTokenCreateRequest, opts ...grpc.CallOption) (*ProjectTokenResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ProjectTokenResponse)
	err := c.cc.Invoke(ctx, ProjectService_CreateToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) DeleteToken(ctx context.Context, in *ProjectTokenDeleteRequest, opts ...grpc.CallOption) (*EmptyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, ProjectService_DeleteToken_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) Create(ctx context.Context, in *ProjectCreateRequest, opts ...grpc.CallOption) (*v1alpha1.AppProject, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(v1alpha1.AppProject)
	err := c.cc.Invoke(ctx, ProjectService_Create_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) List(ctx context.Context, in *ProjectQuery, opts ...grpc.CallOption) (*v1alpha1.AppProjectList, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(v1alpha1.AppProjectList)
	err := c.cc.Invoke(ctx, ProjectService_List_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetDetailedProject(ctx context.Context, in *ProjectQuery, opts ...grpc.CallOption) (*DetailedProjectsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DetailedProjectsResponse)
	err := c.cc.Invoke(ctx, ProjectService_GetDetailedProject_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) Get(ctx context.Context, in *ProjectQuery, opts ...grpc.CallOption) (*v1alpha1.AppProject, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(v1alpha1.AppProject)
	err := c.cc.Invoke(ctx, ProjectService_Get_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetGlobalProjects(ctx context.Context, in *ProjectQuery, opts ...grpc.CallOption) (*GlobalProjectsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GlobalProjectsResponse)
	err := c.cc.Invoke(ctx, ProjectService_GetGlobalProjects_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) Update(ctx context.Context, in *ProjectUpdateRequest, opts ...grpc.CallOption) (*v1alpha1.AppProject, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(v1alpha1.AppProject)
	err := c.cc.Invoke(ctx, ProjectService_Update_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) Delete(ctx context.Context, in *ProjectQuery, opts ...grpc.CallOption) (*EmptyResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, ProjectService_Delete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) ListEvents(ctx context.Context, in *ProjectQuery, opts ...grpc.CallOption) (*EventListWrapper, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(EventListWrapper)
	err := c.cc.Invoke(ctx, ProjectService_ListEvents_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) GetSyncWindowsState(ctx context.Context, in *SyncWindowsQuery, opts ...grpc.CallOption) (*SyncWindowsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SyncWindowsResponse)
	err := c.cc.Invoke(ctx, ProjectService_GetSyncWindowsState_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *projectServiceClient) ListLinks(ctx context.Context, in *ListProjectLinksRequest, opts ...grpc.CallOption) (*application.LinksResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(application.LinksResponse)
	err := c.cc.Invoke(ctx, ProjectService_ListLinks_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProjectServiceServer is the server API for ProjectService service.
// All implementations must embed UnimplementedProjectServiceServer
// for forward compatibility.
//
// ProjectService
type ProjectServiceServer interface {
	// Create a new project token
	CreateToken(context.Context, *ProjectTokenCreateRequest) (*ProjectTokenResponse, error)
	// Delete a new project token
	DeleteToken(context.Context, *ProjectTokenDeleteRequest) (*EmptyResponse, error)
	// Create a new project
	Create(context.Context, *ProjectCreateRequest) (*v1alpha1.AppProject, error)
	// List returns list of projects
	List(context.Context, *ProjectQuery) (*v1alpha1.AppProjectList, error)
	// GetDetailedProject returns a project that include project, global project and scoped resources by name
	GetDetailedProject(context.Context, *ProjectQuery) (*DetailedProjectsResponse, error)
	// Get returns a project by name
	Get(context.Context, *ProjectQuery) (*v1alpha1.AppProject, error)
	// Get returns a virtual project by name
	GetGlobalProjects(context.Context, *ProjectQuery) (*GlobalProjectsResponse, error)
	// Update updates a project
	Update(context.Context, *ProjectUpdateRequest) (*v1alpha1.AppProject, error)
	// Delete deletes a project
	Delete(context.Context, *ProjectQuery) (*EmptyResponse, error)
	// ListEvents returns a list of project events
	ListEvents(context.Context, *ProjectQuery) (*EventListWrapper, error)
	// GetSchedulesState returns true if there are any active sync syncWindows
	GetSyncWindowsState(context.Context, *SyncWindowsQuery) (*SyncWindowsResponse, error)
	// ListLinks returns all deep links for the particular project
	ListLinks(context.Context, *ListProjectLinksRequest) (*application.LinksResponse, error)
	mustEmbedUnimplementedProjectServiceServer()
}

// UnimplementedProjectServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedProjectServiceServer struct{}

func (UnimplementedProjectServiceServer) CreateToken(context.Context, *ProjectTokenCreateRequest) (*ProjectTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateToken not implemented")
}
func (UnimplementedProjectServiceServer) DeleteToken(context.Context, *ProjectTokenDeleteRequest) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteToken not implemented")
}
func (UnimplementedProjectServiceServer) Create(context.Context, *ProjectCreateRequest) (*v1alpha1.AppProject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedProjectServiceServer) List(context.Context, *ProjectQuery) (*v1alpha1.AppProjectList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedProjectServiceServer) GetDetailedProject(context.Context, *ProjectQuery) (*DetailedProjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDetailedProject not implemented")
}
func (UnimplementedProjectServiceServer) Get(context.Context, *ProjectQuery) (*v1alpha1.AppProject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedProjectServiceServer) GetGlobalProjects(context.Context, *ProjectQuery) (*GlobalProjectsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGlobalProjects not implemented")
}
func (UnimplementedProjectServiceServer) Update(context.Context, *ProjectUpdateRequest) (*v1alpha1.AppProject, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedProjectServiceServer) Delete(context.Context, *ProjectQuery) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedProjectServiceServer) ListEvents(context.Context, *ProjectQuery) (*EventListWrapper, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListEvents not implemented")
}
func (UnimplementedProjectServiceServer) GetSyncWindowsState(context.Context, *SyncWindowsQuery) (*SyncWindowsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSyncWindowsState not implemented")
}
func (UnimplementedProjectServiceServer) ListLinks(context.Context, *ListProjectLinksRequest) (*application.LinksResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLinks not implemented")
}
func (UnimplementedProjectServiceServer) mustEmbedUnimplementedProjectServiceServer() {}
func (UnimplementedProjectServiceServer) testEmbeddedByValue()                        {}

// UnsafeProjectServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProjectServiceServer will
// result in compilation errors.
type UnsafeProjectServiceServer interface {
	mustEmbedUnimplementedProjectServiceServer()
}

func RegisterProjectServiceServer(s grpc.ServiceRegistrar, srv ProjectServiceServer) {
	// If the following call pancis, it indicates UnimplementedProjectServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&ProjectService_ServiceDesc, srv)
}

func _ProjectService_CreateToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectTokenCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).CreateToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectService_CreateToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).CreateToken(ctx, req.(*ProjectTokenCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_DeleteToken_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectTokenDeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).DeleteToken(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectService_DeleteToken_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).DeleteToken(ctx, req.(*ProjectTokenDeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectCreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectService_Create_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).Create(ctx, req.(*ProjectCreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectService_List_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).List(ctx, req.(*ProjectQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetDetailedProject_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetDetailedProject(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectService_GetDetailedProject_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetDetailedProject(ctx, req.(*ProjectQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectService_Get_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).Get(ctx, req.(*ProjectQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetGlobalProjects_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetGlobalProjects(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectService_GetGlobalProjects_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetGlobalProjects(ctx, req.(*ProjectQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectService_Update_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).Update(ctx, req.(*ProjectUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectService_Delete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).Delete(ctx, req.(*ProjectQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_ListEvents_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProjectQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).ListEvents(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectService_ListEvents_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).ListEvents(ctx, req.(*ProjectQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_GetSyncWindowsState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SyncWindowsQuery)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).GetSyncWindowsState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectService_GetSyncWindowsState_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).GetSyncWindowsState(ctx, req.(*SyncWindowsQuery))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProjectService_ListLinks_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListProjectLinksRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProjectServiceServer).ListLinks(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: ProjectService_ListLinks_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProjectServiceServer).ListLinks(ctx, req.(*ListProjectLinksRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ProjectService_ServiceDesc is the grpc.ServiceDesc for ProjectService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ProjectService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "project.ProjectService",
	HandlerType: (*ProjectServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateToken",
			Handler:    _ProjectService_CreateToken_Handler,
		},
		{
			MethodName: "DeleteToken",
			Handler:    _ProjectService_DeleteToken_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _ProjectService_Create_Handler,
		},
		{
			MethodName: "List",
			Handler:    _ProjectService_List_Handler,
		},
		{
			MethodName: "GetDetailedProject",
			Handler:    _ProjectService_GetDetailedProject_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _ProjectService_Get_Handler,
		},
		{
			MethodName: "GetGlobalProjects",
			Handler:    _ProjectService_GetGlobalProjects_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _ProjectService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _ProjectService_Delete_Handler,
		},
		{
			MethodName: "ListEvents",
			Handler:    _ProjectService_ListEvents_Handler,
		},
		{
			MethodName: "GetSyncWindowsState",
			Handler:    _ProjectService_GetSyncWindowsState_Handler,
		},
		{
			MethodName: "ListLinks",
			Handler:    _ProjectService_ListLinks_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server/project/project.proto",
}
