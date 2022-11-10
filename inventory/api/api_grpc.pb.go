// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: inventory/api/api.proto

package api

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

// InventoryServiceClient is the client API for InventoryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InventoryServiceClient interface {
	AddLocations(ctx context.Context, in *AddLocationsReq, opts ...grpc.CallOption) (*AddLocationsResp, error)
	AddProducts(ctx context.Context, in *AddProductsReq, opts ...grpc.CallOption) (*AddProductsResp, error)
	ListLocations(ctx context.Context, in *ListLocationsReq, opts ...grpc.CallOption) (*ListLocationsResp, error)
	MoveLocation(ctx context.Context, in *MoveLocationReq, opts ...grpc.CallOption) (*MoveLocationResp, error)
	UpdateInventory(ctx context.Context, in *UpdateInventoryReq, opts ...grpc.CallOption) (*UpdateInventoryResp, error)
	GetLocInventory(ctx context.Context, in *GetLocInventoryReq, opts ...grpc.CallOption) (*GetLocInventoryResp, error)
	Reserve(ctx context.Context, in *ReserveReq, opts ...grpc.CallOption) (*ReserveResp, error)
}

type inventoryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewInventoryServiceClient(cc grpc.ClientConnInterface) InventoryServiceClient {
	return &inventoryServiceClient{cc}
}

func (c *inventoryServiceClient) AddLocations(ctx context.Context, in *AddLocationsReq, opts ...grpc.CallOption) (*AddLocationsResp, error) {
	out := new(AddLocationsResp)
	err := c.cc.Invoke(ctx, "/protos.InventoryService/AddLocations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryServiceClient) AddProducts(ctx context.Context, in *AddProductsReq, opts ...grpc.CallOption) (*AddProductsResp, error) {
	out := new(AddProductsResp)
	err := c.cc.Invoke(ctx, "/protos.InventoryService/AddProducts", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryServiceClient) ListLocations(ctx context.Context, in *ListLocationsReq, opts ...grpc.CallOption) (*ListLocationsResp, error) {
	out := new(ListLocationsResp)
	err := c.cc.Invoke(ctx, "/protos.InventoryService/ListLocations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryServiceClient) MoveLocation(ctx context.Context, in *MoveLocationReq, opts ...grpc.CallOption) (*MoveLocationResp, error) {
	out := new(MoveLocationResp)
	err := c.cc.Invoke(ctx, "/protos.InventoryService/MoveLocation", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryServiceClient) UpdateInventory(ctx context.Context, in *UpdateInventoryReq, opts ...grpc.CallOption) (*UpdateInventoryResp, error) {
	out := new(UpdateInventoryResp)
	err := c.cc.Invoke(ctx, "/protos.InventoryService/UpdateInventory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryServiceClient) GetLocInventory(ctx context.Context, in *GetLocInventoryReq, opts ...grpc.CallOption) (*GetLocInventoryResp, error) {
	out := new(GetLocInventoryResp)
	err := c.cc.Invoke(ctx, "/protos.InventoryService/GetLocInventory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *inventoryServiceClient) Reserve(ctx context.Context, in *ReserveReq, opts ...grpc.CallOption) (*ReserveResp, error) {
	out := new(ReserveResp)
	err := c.cc.Invoke(ctx, "/protos.InventoryService/Reserve", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InventoryServiceServer is the server API for InventoryService service.
// All implementations must embed UnimplementedInventoryServiceServer
// for forward compatibility
type InventoryServiceServer interface {
	AddLocations(context.Context, *AddLocationsReq) (*AddLocationsResp, error)
	AddProducts(context.Context, *AddProductsReq) (*AddProductsResp, error)
	ListLocations(context.Context, *ListLocationsReq) (*ListLocationsResp, error)
	MoveLocation(context.Context, *MoveLocationReq) (*MoveLocationResp, error)
	UpdateInventory(context.Context, *UpdateInventoryReq) (*UpdateInventoryResp, error)
	GetLocInventory(context.Context, *GetLocInventoryReq) (*GetLocInventoryResp, error)
	Reserve(context.Context, *ReserveReq) (*ReserveResp, error)
	mustEmbedUnimplementedInventoryServiceServer()
}

// UnimplementedInventoryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedInventoryServiceServer struct {
}

func (UnimplementedInventoryServiceServer) AddLocations(context.Context, *AddLocationsReq) (*AddLocationsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddLocations not implemented")
}
func (UnimplementedInventoryServiceServer) AddProducts(context.Context, *AddProductsReq) (*AddProductsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddProducts not implemented")
}
func (UnimplementedInventoryServiceServer) ListLocations(context.Context, *ListLocationsReq) (*ListLocationsResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListLocations not implemented")
}
func (UnimplementedInventoryServiceServer) MoveLocation(context.Context, *MoveLocationReq) (*MoveLocationResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MoveLocation not implemented")
}
func (UnimplementedInventoryServiceServer) UpdateInventory(context.Context, *UpdateInventoryReq) (*UpdateInventoryResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateInventory not implemented")
}
func (UnimplementedInventoryServiceServer) GetLocInventory(context.Context, *GetLocInventoryReq) (*GetLocInventoryResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetLocInventory not implemented")
}
func (UnimplementedInventoryServiceServer) Reserve(context.Context, *ReserveReq) (*ReserveResp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reserve not implemented")
}
func (UnimplementedInventoryServiceServer) mustEmbedUnimplementedInventoryServiceServer() {}

// UnsafeInventoryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InventoryServiceServer will
// result in compilation errors.
type UnsafeInventoryServiceServer interface {
	mustEmbedUnimplementedInventoryServiceServer()
}

func RegisterInventoryServiceServer(s grpc.ServiceRegistrar, srv InventoryServiceServer) {
	s.RegisterService(&InventoryService_ServiceDesc, srv)
}

func _InventoryService_AddLocations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddLocationsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServiceServer).AddLocations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.InventoryService/AddLocations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServiceServer).AddLocations(ctx, req.(*AddLocationsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _InventoryService_AddProducts_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddProductsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServiceServer).AddProducts(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.InventoryService/AddProducts",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServiceServer).AddProducts(ctx, req.(*AddProductsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _InventoryService_ListLocations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListLocationsReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServiceServer).ListLocations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.InventoryService/ListLocations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServiceServer).ListLocations(ctx, req.(*ListLocationsReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _InventoryService_MoveLocation_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MoveLocationReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServiceServer).MoveLocation(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.InventoryService/MoveLocation",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServiceServer).MoveLocation(ctx, req.(*MoveLocationReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _InventoryService_UpdateInventory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateInventoryReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServiceServer).UpdateInventory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.InventoryService/UpdateInventory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServiceServer).UpdateInventory(ctx, req.(*UpdateInventoryReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _InventoryService_GetLocInventory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetLocInventoryReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServiceServer).GetLocInventory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.InventoryService/GetLocInventory",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServiceServer).GetLocInventory(ctx, req.(*GetLocInventoryReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _InventoryService_Reserve_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReserveReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InventoryServiceServer).Reserve(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.InventoryService/Reserve",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InventoryServiceServer).Reserve(ctx, req.(*ReserveReq))
	}
	return interceptor(ctx, in, info, handler)
}

// InventoryService_ServiceDesc is the grpc.ServiceDesc for InventoryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var InventoryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protos.InventoryService",
	HandlerType: (*InventoryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddLocations",
			Handler:    _InventoryService_AddLocations_Handler,
		},
		{
			MethodName: "AddProducts",
			Handler:    _InventoryService_AddProducts_Handler,
		},
		{
			MethodName: "ListLocations",
			Handler:    _InventoryService_ListLocations_Handler,
		},
		{
			MethodName: "MoveLocation",
			Handler:    _InventoryService_MoveLocation_Handler,
		},
		{
			MethodName: "UpdateInventory",
			Handler:    _InventoryService_UpdateInventory_Handler,
		},
		{
			MethodName: "GetLocInventory",
			Handler:    _InventoryService_GetLocInventory_Handler,
		},
		{
			MethodName: "Reserve",
			Handler:    _InventoryService_Reserve_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "inventory/api/api.proto",
}

// SpecServiceClient is the client API for SpecService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SpecServiceClient interface {
	Spec(ctx context.Context, in *SpecRequest, opts ...grpc.CallOption) (*SpecResponse, error)
}

type specServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSpecServiceClient(cc grpc.ClientConnInterface) SpecServiceClient {
	return &specServiceClient{cc}
}

func (c *specServiceClient) Spec(ctx context.Context, in *SpecRequest, opts ...grpc.CallOption) (*SpecResponse, error) {
	out := new(SpecResponse)
	err := c.cc.Invoke(ctx, "/protos.SpecService/Spec", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SpecServiceServer is the server API for SpecService service.
// All implementations must embed UnimplementedSpecServiceServer
// for forward compatibility
type SpecServiceServer interface {
	Spec(context.Context, *SpecRequest) (*SpecResponse, error)
	mustEmbedUnimplementedSpecServiceServer()
}

// UnimplementedSpecServiceServer must be embedded to have forward compatible implementations.
type UnimplementedSpecServiceServer struct {
}

func (UnimplementedSpecServiceServer) Spec(context.Context, *SpecRequest) (*SpecResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Spec not implemented")
}
func (UnimplementedSpecServiceServer) mustEmbedUnimplementedSpecServiceServer() {}

// UnsafeSpecServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SpecServiceServer will
// result in compilation errors.
type UnsafeSpecServiceServer interface {
	mustEmbedUnimplementedSpecServiceServer()
}

func RegisterSpecServiceServer(s grpc.ServiceRegistrar, srv SpecServiceServer) {
	s.RegisterService(&SpecService_ServiceDesc, srv)
}

func _SpecService_Spec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SpecRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpecServiceServer).Spec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.SpecService/Spec",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpecServiceServer).Spec(ctx, req.(*SpecRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// SpecService_ServiceDesc is the grpc.ServiceDesc for SpecService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SpecService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protos.SpecService",
	HandlerType: (*SpecServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Spec",
			Handler:    _SpecService_Spec_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "inventory/api/api.proto",
}
