// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: proto/power.proto

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

// NotificationClient is the client API for Notification service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NotificationClient interface {
	SubscribeCustomer(ctx context.Context, in *Customer, opts ...grpc.CallOption) (*Customer, error)
	UnsubscribeCustomer(ctx context.Context, in *CustomerId, opts ...grpc.CallOption) (*UnsubscribeResponse, error)
}

type notificationClient struct {
	cc grpc.ClientConnInterface
}

func NewNotificationClient(cc grpc.ClientConnInterface) NotificationClient {
	return &notificationClient{cc}
}

func (c *notificationClient) SubscribeCustomer(ctx context.Context, in *Customer, opts ...grpc.CallOption) (*Customer, error) {
	out := new(Customer)
	err := c.cc.Invoke(ctx, "/Notification/SubscribeCustomer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *notificationClient) UnsubscribeCustomer(ctx context.Context, in *CustomerId, opts ...grpc.CallOption) (*UnsubscribeResponse, error) {
	out := new(UnsubscribeResponse)
	err := c.cc.Invoke(ctx, "/Notification/UnsubscribeCustomer", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NotificationServer is the server API for Notification service.
// All implementations must embed UnimplementedNotificationServer
// for forward compatibility
type NotificationServer interface {
	SubscribeCustomer(context.Context, *Customer) (*Customer, error)
	UnsubscribeCustomer(context.Context, *CustomerId) (*UnsubscribeResponse, error)
	mustEmbedUnimplementedNotificationServer()
}

// UnimplementedNotificationServer must be embedded to have forward compatible implementations.
type UnimplementedNotificationServer struct {
}

func (UnimplementedNotificationServer) SubscribeCustomer(context.Context, *Customer) (*Customer, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SubscribeCustomer not implemented")
}
func (UnimplementedNotificationServer) UnsubscribeCustomer(context.Context, *CustomerId) (*UnsubscribeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnsubscribeCustomer not implemented")
}
func (UnimplementedNotificationServer) mustEmbedUnimplementedNotificationServer() {}

// UnsafeNotificationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NotificationServer will
// result in compilation errors.
type UnsafeNotificationServer interface {
	mustEmbedUnimplementedNotificationServer()
}

func RegisterNotificationServer(s grpc.ServiceRegistrar, srv NotificationServer) {
	s.RegisterService(&Notification_ServiceDesc, srv)
}

func _Notification_SubscribeCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Customer)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServer).SubscribeCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Notification/SubscribeCustomer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServer).SubscribeCustomer(ctx, req.(*Customer))
	}
	return interceptor(ctx, in, info, handler)
}

func _Notification_UnsubscribeCustomer_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CustomerId)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NotificationServer).UnsubscribeCustomer(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Notification/UnsubscribeCustomer",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NotificationServer).UnsubscribeCustomer(ctx, req.(*CustomerId))
	}
	return interceptor(ctx, in, info, handler)
}

// Notification_ServiceDesc is the grpc.ServiceDesc for Notification service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Notification_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Notification",
	HandlerType: (*NotificationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SubscribeCustomer",
			Handler:    _Notification_SubscribeCustomer_Handler,
		},
		{
			MethodName: "UnsubscribeCustomer",
			Handler:    _Notification_UnsubscribeCustomer_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/power.proto",
}

// ScrappingClient is the client API for Scrapping service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ScrappingClient interface {
	FilterOutages(ctx context.Context, in *OutageFilter, opts ...grpc.CallOption) (*OutageFilterResponse, error)
	GetLocationsUnder(ctx context.Context, in *Cordinate, opts ...grpc.CallOption) (*LocationsUnderCords, error)
}

type scrappingClient struct {
	cc grpc.ClientConnInterface
}

func NewScrappingClient(cc grpc.ClientConnInterface) ScrappingClient {
	return &scrappingClient{cc}
}

func (c *scrappingClient) FilterOutages(ctx context.Context, in *OutageFilter, opts ...grpc.CallOption) (*OutageFilterResponse, error) {
	out := new(OutageFilterResponse)
	err := c.cc.Invoke(ctx, "/Scrapping/FilterOutages", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *scrappingClient) GetLocationsUnder(ctx context.Context, in *Cordinate, opts ...grpc.CallOption) (*LocationsUnderCords, error) {
	out := new(LocationsUnderCords)
	err := c.cc.Invoke(ctx, "/Scrapping/GetLocationsUnder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ScrappingServer is the server API for Scrapping service.
// All implementations must embed UnimplementedScrappingServer
// for forward compatibility
type ScrappingServer interface {
	FilterOutages(context.Context, *OutageFilter) (*OutageFilterResponse, error)
	GetLocationsUnder(context.Context, *Cordinate) (*LocationsUnderCords, error)
	mustEmbedUnimplementedScrappingServer()
}

// UnimplementedScrappingServer must be embedded to have forward compatible implementations.
type UnimplementedScrappingServer struct {
}

func (UnimplementedScrappingServer) FilterOutages(context.Context, *OutageFilter) (*OutageFilterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FilterOutages not implemented")
}
func (UnimplementedScrappingServer) GetLocationsUnder(context.Context, *Cordinate) (*LocationsUnderCords, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLocationsUnder not implemented")
}
func (UnimplementedScrappingServer) mustEmbedUnimplementedScrappingServer() {}

// UnsafeScrappingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ScrappingServer will
// result in compilation errors.
type UnsafeScrappingServer interface {
	mustEmbedUnimplementedScrappingServer()
}

func RegisterScrappingServer(s grpc.ServiceRegistrar, srv ScrappingServer) {
	s.RegisterService(&Scrapping_ServiceDesc, srv)
}

func _Scrapping_FilterOutages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OutageFilter)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScrappingServer).FilterOutages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Scrapping/FilterOutages",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScrappingServer).FilterOutages(ctx, req.(*OutageFilter))
	}
	return interceptor(ctx, in, info, handler)
}

func _Scrapping_GetLocationsUnder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Cordinate)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ScrappingServer).GetLocationsUnder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Scrapping/GetLocationsUnder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ScrappingServer).GetLocationsUnder(ctx, req.(*Cordinate))
	}
	return interceptor(ctx, in, info, handler)
}

// Scrapping_ServiceDesc is the grpc.ServiceDesc for Scrapping service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Scrapping_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Scrapping",
	HandlerType: (*ScrappingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FilterOutages",
			Handler:    _Scrapping_FilterOutages_Handler,
		},
		{
			MethodName: "GetLocationsUnder",
			Handler:    _Scrapping_GetLocationsUnder_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/power.proto",
}