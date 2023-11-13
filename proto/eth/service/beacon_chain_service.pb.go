// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        v4.23.3
// source: proto/eth/service/beacon_chain_service.proto

package service

import (
	context "context"
	reflect "reflect"

	v1 "github.com/prysmaticlabs/prysm/v4/proto/eth/v1"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	_ "google.golang.org/protobuf/types/descriptorpb"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

var File_proto_eth_service_beacon_chain_service_proto protoreflect.FileDescriptor

var file_proto_eth_service_beacon_chain_service_proto_rawDesc = []byte{
	0x0a, 0x2c, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2f, 0x62, 0x65, 0x61, 0x63, 0x6f, 0x6e, 0x5f, 0x63, 0x68, 0x61, 0x69, 0x6e,
	0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x14,
	0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69,
	0x2f, 0x61, 0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x20, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x62, 0x75, 0x66, 0x2f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x1a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x76, 0x31, 0x2f,
	0x62, 0x65, 0x61, 0x63, 0x6f, 0x6e, 0x5f, 0x62, 0x6c, 0x6f, 0x63, 0x6b, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x1a, 0x1f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x76, 0x31,
	0x2f, 0x62, 0x65, 0x61, 0x63, 0x6f, 0x6e, 0x5f, 0x63, 0x68, 0x61, 0x69, 0x6e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x32, 0x87, 0x03, 0x0a, 0x0b, 0x42, 0x65, 0x61, 0x63, 0x6f, 0x6e, 0x43, 0x68,
	0x61, 0x69, 0x6e, 0x12, 0x8e, 0x01, 0x0a, 0x13, 0x47, 0x65, 0x74, 0x57, 0x65, 0x61, 0x6b, 0x53,
	0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x12, 0x16, 0x2e, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d,
	0x70, 0x74, 0x79, 0x1a, 0x29, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65,
	0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x57, 0x65, 0x61, 0x6b, 0x53, 0x75, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x34,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x2b, 0x12, 0x29, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x62, 0x65, 0x61, 0x63, 0x6f, 0x6e, 0x2f,
	0x77, 0x65, 0x61, 0x6b, 0x5f, 0x73, 0x75, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74,
	0x79, 0x88, 0x02, 0x01, 0x12, 0x7f, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x46, 0x6f, 0x72, 0x6b, 0x53,
	0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a,
	0x25, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x76,
	0x31, 0x2e, 0x46, 0x6f, 0x72, 0x6b, 0x53, 0x63, 0x68, 0x65, 0x64, 0x75, 0x6c, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x2d, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x27, 0x12, 0x25,
	0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x76, 0x31,
	0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x66, 0x6f, 0x72, 0x6b, 0x5f, 0x73, 0x63, 0x68,
	0x65, 0x64, 0x75, 0x6c, 0x65, 0x12, 0x66, 0x0a, 0x07, 0x47, 0x65, 0x74, 0x53, 0x70, 0x65, 0x63,
	0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x1d, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72,
	0x65, 0x75, 0x6d, 0x2e, 0x65, 0x74, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x70, 0x65, 0x63, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x24, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x1e, 0x12,
	0x1c, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x65, 0x74, 0x68, 0x2f, 0x76,
	0x31, 0x2f, 0x63, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x2f, 0x73, 0x70, 0x65, 0x63, 0x42, 0x98, 0x01,
	0x0a, 0x18, 0x6f, 0x72, 0x67, 0x2e, 0x65, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x65,
	0x74, 0x68, 0x2e, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x42, 0x17, 0x42, 0x65, 0x61, 0x63,
	0x6f, 0x6e, 0x43, 0x68, 0x61, 0x69, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x50, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x33, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x70, 0x72, 0x79, 0x73, 0x6d, 0x61, 0x74, 0x69, 0x63, 0x6c, 0x61, 0x62, 0x73, 0x2f,
	0x70, 0x72, 0x79, 0x73, 0x6d, 0x2f, 0x76, 0x34, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65,
	0x74, 0x68, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0xaa, 0x02, 0x14, 0x45, 0x74, 0x68,
	0x65, 0x72, 0x65, 0x75, 0x6d, 0x2e, 0x45, 0x74, 0x68, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0xca, 0x02, 0x14, 0x45, 0x74, 0x68, 0x65, 0x72, 0x65, 0x75, 0x6d, 0x5c, 0x45, 0x74, 0x68,
	0x5c, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_proto_eth_service_beacon_chain_service_proto_goTypes = []interface{}{
	(*emptypb.Empty)(nil),               // 0: google.protobuf.Empty
	(*v1.WeakSubjectivityResponse)(nil), // 1: ethereum.eth.v1.WeakSubjectivityResponse
	(*v1.ForkScheduleResponse)(nil),     // 2: ethereum.eth.v1.ForkScheduleResponse
	(*v1.SpecResponse)(nil),             // 3: ethereum.eth.v1.SpecResponse
}
var file_proto_eth_service_beacon_chain_service_proto_depIdxs = []int32{
	0, // 0: ethereum.eth.service.BeaconChain.GetWeakSubjectivity:input_type -> google.protobuf.Empty
	0, // 1: ethereum.eth.service.BeaconChain.GetForkSchedule:input_type -> google.protobuf.Empty
	0, // 2: ethereum.eth.service.BeaconChain.GetSpec:input_type -> google.protobuf.Empty
	1, // 3: ethereum.eth.service.BeaconChain.GetWeakSubjectivity:output_type -> ethereum.eth.v1.WeakSubjectivityResponse
	2, // 4: ethereum.eth.service.BeaconChain.GetForkSchedule:output_type -> ethereum.eth.v1.ForkScheduleResponse
	3, // 5: ethereum.eth.service.BeaconChain.GetSpec:output_type -> ethereum.eth.v1.SpecResponse
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_eth_service_beacon_chain_service_proto_init() }
func file_proto_eth_service_beacon_chain_service_proto_init() {
	if File_proto_eth_service_beacon_chain_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_eth_service_beacon_chain_service_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_eth_service_beacon_chain_service_proto_goTypes,
		DependencyIndexes: file_proto_eth_service_beacon_chain_service_proto_depIdxs,
	}.Build()
	File_proto_eth_service_beacon_chain_service_proto = out.File
	file_proto_eth_service_beacon_chain_service_proto_rawDesc = nil
	file_proto_eth_service_beacon_chain_service_proto_goTypes = nil
	file_proto_eth_service_beacon_chain_service_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// BeaconChainClient is the client API for BeaconChain service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BeaconChainClient interface {
	// Deprecated: Do not use.
	GetWeakSubjectivity(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1.WeakSubjectivityResponse, error)
	GetForkSchedule(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1.ForkScheduleResponse, error)
	GetSpec(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1.SpecResponse, error)
}

type beaconChainClient struct {
	cc grpc.ClientConnInterface
}

func NewBeaconChainClient(cc grpc.ClientConnInterface) BeaconChainClient {
	return &beaconChainClient{cc}
}

// Deprecated: Do not use.
func (c *beaconChainClient) GetWeakSubjectivity(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1.WeakSubjectivityResponse, error) {
	out := new(v1.WeakSubjectivityResponse)
	err := c.cc.Invoke(ctx, "/ethereum.eth.service.BeaconChain/GetWeakSubjectivity", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beaconChainClient) GetForkSchedule(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1.ForkScheduleResponse, error) {
	out := new(v1.ForkScheduleResponse)
	err := c.cc.Invoke(ctx, "/ethereum.eth.service.BeaconChain/GetForkSchedule", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *beaconChainClient) GetSpec(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*v1.SpecResponse, error) {
	out := new(v1.SpecResponse)
	err := c.cc.Invoke(ctx, "/ethereum.eth.service.BeaconChain/GetSpec", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BeaconChainServer is the server API for BeaconChain service.
type BeaconChainServer interface {
	// Deprecated: Do not use.
	GetWeakSubjectivity(context.Context, *emptypb.Empty) (*v1.WeakSubjectivityResponse, error)
	GetForkSchedule(context.Context, *emptypb.Empty) (*v1.ForkScheduleResponse, error)
	GetSpec(context.Context, *emptypb.Empty) (*v1.SpecResponse, error)
}

// UnimplementedBeaconChainServer can be embedded to have forward compatible implementations.
type UnimplementedBeaconChainServer struct {
}

func (*UnimplementedBeaconChainServer) GetWeakSubjectivity(context.Context, *emptypb.Empty) (*v1.WeakSubjectivityResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWeakSubjectivity not implemented")
}
func (*UnimplementedBeaconChainServer) GetForkSchedule(context.Context, *emptypb.Empty) (*v1.ForkScheduleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetForkSchedule not implemented")
}
func (*UnimplementedBeaconChainServer) GetSpec(context.Context, *emptypb.Empty) (*v1.SpecResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetSpec not implemented")
}

func RegisterBeaconChainServer(s *grpc.Server, srv BeaconChainServer) {
	s.RegisterService(&_BeaconChain_serviceDesc, srv)
}

func _BeaconChain_GetWeakSubjectivity_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeaconChainServer).GetWeakSubjectivity(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.eth.service.BeaconChain/GetWeakSubjectivity",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeaconChainServer).GetWeakSubjectivity(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeaconChain_GetForkSchedule_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeaconChainServer).GetForkSchedule(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.eth.service.BeaconChain/GetForkSchedule",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeaconChainServer).GetForkSchedule(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _BeaconChain_GetSpec_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(emptypb.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BeaconChainServer).GetSpec(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ethereum.eth.service.BeaconChain/GetSpec",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BeaconChainServer).GetSpec(ctx, req.(*emptypb.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _BeaconChain_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ethereum.eth.service.BeaconChain",
	HandlerType: (*BeaconChainServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetWeakSubjectivity",
			Handler:    _BeaconChain_GetWeakSubjectivity_Handler,
		},
		{
			MethodName: "GetForkSchedule",
			Handler:    _BeaconChain_GetForkSchedule_Handler,
		},
		{
			MethodName: "GetSpec",
			Handler:    _BeaconChain_GetSpec_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/eth/service/beacon_chain_service.proto",
}
