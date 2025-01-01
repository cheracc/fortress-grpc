// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.1
// source: fortress.proto

package grpc

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
	Auth_Authorize_FullMethodName = "/grpc.Auth/Authorize"
)

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AuthClient interface {
	Authorize(ctx context.Context, in *PlayerInfo, opts ...grpc.CallOption) (*AuthInfo, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) Authorize(ctx context.Context, in *PlayerInfo, opts ...grpc.CallOption) (*AuthInfo, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AuthInfo)
	err := c.cc.Invoke(ctx, Auth_Authorize_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
// All implementations must embed UnimplementedAuthServer
// for forward compatibility.
type AuthServer interface {
	Authorize(context.Context, *PlayerInfo) (*AuthInfo, error)
	mustEmbedUnimplementedAuthServer()
}

// UnimplementedAuthServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedAuthServer struct{}

func (UnimplementedAuthServer) Authorize(context.Context, *PlayerInfo) (*AuthInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}
func (UnimplementedAuthServer) mustEmbedUnimplementedAuthServer() {}
func (UnimplementedAuthServer) testEmbeddedByValue()              {}

// UnsafeAuthServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AuthServer will
// result in compilation errors.
type UnsafeAuthServer interface {
	mustEmbedUnimplementedAuthServer()
}

func RegisterAuthServer(s grpc.ServiceRegistrar, srv AuthServer) {
	// If the following call pancis, it indicates UnimplementedAuthServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Auth_ServiceDesc, srv)
}

func _Auth_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlayerInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Auth_Authorize_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Authorize(ctx, req.(*PlayerInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// Auth_ServiceDesc is the grpc.ServiceDesc for Auth service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Auth_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authorize",
			Handler:    _Auth_Authorize_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fortress.proto",
}

const (
	Command_Command_FullMethodName = "/grpc.Command/Command"
)

// CommandClient is the client API for Command service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommandClient interface {
	Command(ctx context.Context, in *CommandInfo, opts ...grpc.CallOption) (*CommandReturn, error)
}

type commandClient struct {
	cc grpc.ClientConnInterface
}

func NewCommandClient(cc grpc.ClientConnInterface) CommandClient {
	return &commandClient{cc}
}

func (c *commandClient) Command(ctx context.Context, in *CommandInfo, opts ...grpc.CallOption) (*CommandReturn, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CommandReturn)
	err := c.cc.Invoke(ctx, Command_Command_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommandServer is the server API for Command service.
// All implementations must embed UnimplementedCommandServer
// for forward compatibility.
type CommandServer interface {
	Command(context.Context, *CommandInfo) (*CommandReturn, error)
	mustEmbedUnimplementedCommandServer()
}

// UnimplementedCommandServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCommandServer struct{}

func (UnimplementedCommandServer) Command(context.Context, *CommandInfo) (*CommandReturn, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Command not implemented")
}
func (UnimplementedCommandServer) mustEmbedUnimplementedCommandServer() {}
func (UnimplementedCommandServer) testEmbeddedByValue()                 {}

// UnsafeCommandServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommandServer will
// result in compilation errors.
type UnsafeCommandServer interface {
	mustEmbedUnimplementedCommandServer()
}

func RegisterCommandServer(s grpc.ServiceRegistrar, srv CommandServer) {
	// If the following call pancis, it indicates UnimplementedCommandServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Command_ServiceDesc, srv)
}

func _Command_Command_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CommandInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommandServer).Command(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Command_Command_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommandServer).Command(ctx, req.(*CommandInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// Command_ServiceDesc is the grpc.ServiceDesc for Command service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Command_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Command",
	HandlerType: (*CommandServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Command",
			Handler:    _Command_Command_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fortress.proto",
}

const (
	Player_GetPlayerData_FullMethodName = "/grpc.Player/GetPlayerData"
)

// PlayerClient is the client API for Player service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlayerClient interface {
	GetPlayerData(ctx context.Context, in *PlayerInfo, opts ...grpc.CallOption) (*PlayerMessage, error)
}

type playerClient struct {
	cc grpc.ClientConnInterface
}

func NewPlayerClient(cc grpc.ClientConnInterface) PlayerClient {
	return &playerClient{cc}
}

func (c *playerClient) GetPlayerData(ctx context.Context, in *PlayerInfo, opts ...grpc.CallOption) (*PlayerMessage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PlayerMessage)
	err := c.cc.Invoke(ctx, Player_GetPlayerData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlayerServer is the server API for Player service.
// All implementations must embed UnimplementedPlayerServer
// for forward compatibility.
type PlayerServer interface {
	GetPlayerData(context.Context, *PlayerInfo) (*PlayerMessage, error)
	mustEmbedUnimplementedPlayerServer()
}

// UnimplementedPlayerServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedPlayerServer struct{}

func (UnimplementedPlayerServer) GetPlayerData(context.Context, *PlayerInfo) (*PlayerMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlayerData not implemented")
}
func (UnimplementedPlayerServer) mustEmbedUnimplementedPlayerServer() {}
func (UnimplementedPlayerServer) testEmbeddedByValue()                {}

// UnsafePlayerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlayerServer will
// result in compilation errors.
type UnsafePlayerServer interface {
	mustEmbedUnimplementedPlayerServer()
}

func RegisterPlayerServer(s grpc.ServiceRegistrar, srv PlayerServer) {
	// If the following call pancis, it indicates UnimplementedPlayerServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Player_ServiceDesc, srv)
}

func _Player_GetPlayerData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PlayerInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServer).GetPlayerData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Player_GetPlayerData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServer).GetPlayerData(ctx, req.(*PlayerInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// Player_ServiceDesc is the grpc.ServiceDesc for Player service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Player_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Player",
	HandlerType: (*PlayerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPlayerData",
			Handler:    _Player_GetPlayerData_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fortress.proto",
}

const (
	Chat_JoinChannel_FullMethodName = "/grpc.Chat/JoinChannel"
	Chat_SendMessage_FullMethodName = "/grpc.Chat/SendMessage"
)

// ChatClient is the client API for Chat service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ChatClient interface {
	JoinChannel(ctx context.Context, in *ChatRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ChatMessage], error)
	SendMessage(ctx context.Context, in *ChatMessage, opts ...grpc.CallOption) (*Empty, error)
}

type chatClient struct {
	cc grpc.ClientConnInterface
}

func NewChatClient(cc grpc.ClientConnInterface) ChatClient {
	return &chatClient{cc}
}

func (c *chatClient) JoinChannel(ctx context.Context, in *ChatRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[ChatMessage], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &Chat_ServiceDesc.Streams[0], Chat_JoinChannel_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[ChatRequest, ChatMessage]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Chat_JoinChannelClient = grpc.ServerStreamingClient[ChatMessage]

func (c *chatClient) SendMessage(ctx context.Context, in *ChatMessage, opts ...grpc.CallOption) (*Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Empty)
	err := c.cc.Invoke(ctx, Chat_SendMessage_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ChatServer is the server API for Chat service.
// All implementations must embed UnimplementedChatServer
// for forward compatibility.
type ChatServer interface {
	JoinChannel(*ChatRequest, grpc.ServerStreamingServer[ChatMessage]) error
	SendMessage(context.Context, *ChatMessage) (*Empty, error)
	mustEmbedUnimplementedChatServer()
}

// UnimplementedChatServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedChatServer struct{}

func (UnimplementedChatServer) JoinChannel(*ChatRequest, grpc.ServerStreamingServer[ChatMessage]) error {
	return status.Errorf(codes.Unimplemented, "method JoinChannel not implemented")
}
func (UnimplementedChatServer) SendMessage(context.Context, *ChatMessage) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendMessage not implemented")
}
func (UnimplementedChatServer) mustEmbedUnimplementedChatServer() {}
func (UnimplementedChatServer) testEmbeddedByValue()              {}

// UnsafeChatServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ChatServer will
// result in compilation errors.
type UnsafeChatServer interface {
	mustEmbedUnimplementedChatServer()
}

func RegisterChatServer(s grpc.ServiceRegistrar, srv ChatServer) {
	// If the following call pancis, it indicates UnimplementedChatServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Chat_ServiceDesc, srv)
}

func _Chat_JoinChannel_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ChatRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ChatServer).JoinChannel(m, &grpc.GenericServerStream[ChatRequest, ChatMessage]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type Chat_JoinChannelServer = grpc.ServerStreamingServer[ChatMessage]

func _Chat_SendMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChatMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ChatServer).SendMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Chat_SendMessage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ChatServer).SendMessage(ctx, req.(*ChatMessage))
	}
	return interceptor(ctx, in, info, handler)
}

// Chat_ServiceDesc is the grpc.ServiceDesc for Chat service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Chat_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "grpc.Chat",
	HandlerType: (*ChatServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendMessage",
			Handler:    _Chat_SendMessage_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "JoinChannel",
			Handler:       _Chat_JoinChannel_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "fortress.proto",
}
