// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: user.proto

package pb

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

const (
	UserManagement_FetchAllUsers_FullMethodName = "/pb.UserManagement/FetchAllUsers"
	UserManagement_FetchUser_FullMethodName     = "/pb.UserManagement/FetchUser"
	UserManagement_RegisterUser_FullMethodName  = "/pb.UserManagement/RegisterUser"
)

// UserManagementClient is the client API for UserManagement service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserManagementClient interface {
	FetchAllUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserList, error)
	FetchUser(ctx context.Context, in *FetchUserInput, opts ...grpc.CallOption) (*User, error)
	RegisterUser(ctx context.Context, in *RegisterUserInput, opts ...grpc.CallOption) (*User, error)
}

type userManagementClient struct {
	cc grpc.ClientConnInterface
}

func NewUserManagementClient(cc grpc.ClientConnInterface) UserManagementClient {
	return &userManagementClient{cc}
}

func (c *userManagementClient) FetchAllUsers(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*UserList, error) {
	out := new(UserList)
	err := c.cc.Invoke(ctx, UserManagement_FetchAllUsers_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userManagementClient) FetchUser(ctx context.Context, in *FetchUserInput, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserManagement_FetchUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userManagementClient) RegisterUser(ctx context.Context, in *RegisterUserInput, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, UserManagement_RegisterUser_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserManagementServer is the server API for UserManagement service.
// All implementations must embed UnimplementedUserManagementServer
// for forward compatibility
type UserManagementServer interface {
	FetchAllUsers(context.Context, *emptypb.Empty) (*UserList, error)
	FetchUser(context.Context, *FetchUserInput) (*User, error)
	RegisterUser(context.Context, *RegisterUserInput) (*User, error)
	mustEmbedUnimplementedUserManagementServer()
}

// UnimplementedUserManagementServer must be embedded to have forward compatible implementations.
type UnimplementedUserManagementServer struct {
}

func (UnimplementedUserManagementServer) FetchAllUsers(context.Context, *emptypb.Empty) (*UserList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchAllUsers not implemented")
}
func (UnimplementedUserManagementServer) FetchUser(context.Context, *FetchUserInput) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FetchUser not implemented")
}
func (UnimplementedUserManagementServer) RegisterUser(context.Context, *RegisterUserInput) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterUser not implemented")
}
func (UnimplementedUserManagementServer) mustEmbedUnimplementedUserManagementServer() {}

// UnsafeUserManagementServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserManagementServer will
// result in compilation errors.
type UnsafeUserManagementServer interface {
	mustEmbedUnimplementedUserManagementServer()
}

func RegisterUserManagementServer(s grpc.ServiceRegistrar, srv UserManagementServer) {
	s.RegisterService(&UserManagement_ServiceDesc, srv)
}

func _UserManagement_FetchAllUsers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserManagementServer).FetchAllUsers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserManagement_FetchAllUsers_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserManagementServer).FetchAllUsers(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserManagement_FetchUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FetchUserInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserManagementServer).FetchUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserManagement_FetchUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserManagementServer).FetchUser(ctx, req.(*FetchUserInput))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserManagement_RegisterUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterUserInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserManagementServer).RegisterUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: UserManagement_RegisterUser_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserManagementServer).RegisterUser(ctx, req.(*RegisterUserInput))
	}
	return interceptor(ctx, in, info, handler)
}

// UserManagement_ServiceDesc is the grpc.ServiceDesc for UserManagement service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserManagement_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.UserManagement",
	HandlerType: (*UserManagementServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FetchAllUsers",
			Handler:    _UserManagement_FetchAllUsers_Handler,
		},
		{
			MethodName: "FetchUser",
			Handler:    _UserManagement_FetchUser_Handler,
		},
		{
			MethodName: "RegisterUser",
			Handler:    _UserManagement_RegisterUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
