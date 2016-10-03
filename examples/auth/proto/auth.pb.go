// Code generated by protoc-gen-go.
// source: examples/auth/proto/auth.proto
// DO NOT EDIT!

/*
Package auth is a generated protocol buffer package.

It is generated from these files:
	examples/auth/proto/auth.proto

It has these top-level messages:
	Credentials
	Token
*/
package auth

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// The request message containing the new user data.
type Credentials struct {
	Email    string `protobuf:"bytes,1,opt,name=email" json:"email,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password" json:"password,omitempty"`
}

func (m *Credentials) Reset()                    { *m = Credentials{} }
func (m *Credentials) String() string            { return proto.CompactTextString(m) }
func (*Credentials) ProtoMessage()               {}
func (*Credentials) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

// The response message containing the user data
type Token struct {
	Token string `protobuf:"bytes,1,opt,name=token" json:"token,omitempty"`
}

func (m *Token) Reset()                    { *m = Token{} }
func (m *Token) String() string            { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()               {}
func (*Token) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*Credentials)(nil), "auth.Credentials")
	proto.RegisterType((*Token)(nil), "auth.Token")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion3

// Client API for AuthService service

type AuthServiceClient interface {
	Authenticate(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*Token, error)
}

type authServiceClient struct {
	cc *grpc.ClientConn
}

func NewAuthServiceClient(cc *grpc.ClientConn) AuthServiceClient {
	return &authServiceClient{cc}
}

func (c *authServiceClient) Authenticate(ctx context.Context, in *Credentials, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := grpc.Invoke(ctx, "/auth.AuthService/Authenticate", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for AuthService service

type AuthServiceServer interface {
	Authenticate(context.Context, *Credentials) (*Token, error)
}

func RegisterAuthServiceServer(s *grpc.Server, srv AuthServiceServer) {
	s.RegisterService(&_AuthService_serviceDesc, srv)
}

func _AuthService_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Credentials)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServiceServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/auth.AuthService/Authenticate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServiceServer).Authenticate(ctx, req.(*Credentials))
	}
	return interceptor(ctx, in, info, handler)
}

var _AuthService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "auth.AuthService",
	HandlerType: (*AuthServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authenticate",
			Handler:    _AuthService_Authenticate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: fileDescriptor0,
}

func init() { proto.RegisterFile("examples/auth/proto/auth.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 169 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0xad, 0x48, 0xcc,
	0x2d, 0xc8, 0x49, 0x2d, 0xd6, 0x4f, 0x2c, 0x2d, 0xc9, 0xd0, 0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x07,
	0x33, 0xf5, 0xc0, 0x4c, 0x21, 0x16, 0x10, 0x5b, 0xc9, 0x9e, 0x8b, 0xdb, 0xb9, 0x28, 0x35, 0x25,
	0x35, 0xaf, 0x24, 0x33, 0x31, 0xa7, 0x58, 0x48, 0x84, 0x8b, 0x35, 0x35, 0x37, 0x31, 0x33, 0x47,
	0x82, 0x51, 0x81, 0x51, 0x83, 0x33, 0x08, 0xc2, 0x11, 0x92, 0xe2, 0xe2, 0x28, 0x48, 0x2c, 0x2e,
	0x2e, 0xcf, 0x2f, 0x4a, 0x91, 0x60, 0x02, 0x4b, 0xc0, 0xf9, 0x4a, 0xb2, 0x5c, 0xac, 0x21, 0xf9,
	0xd9, 0xa9, 0x79, 0x20, 0xad, 0x25, 0x20, 0x06, 0x4c, 0x2b, 0x98, 0x63, 0x64, 0xcf, 0xc5, 0xed,
	0x58, 0x5a, 0x92, 0x11, 0x9c, 0x5a, 0x54, 0x96, 0x99, 0x9c, 0x2a, 0x64, 0xc0, 0xc5, 0x03, 0xe2,
	0x82, 0xac, 0x4b, 0x4e, 0x2c, 0x49, 0x15, 0x12, 0xd4, 0x03, 0xbb, 0x08, 0xc9, 0x09, 0x52, 0xdc,
	0x10, 0x21, 0xb0, 0xa1, 0x4a, 0x0c, 0x49, 0x6c, 0x60, 0xd7, 0x1a, 0x03, 0x02, 0x00, 0x00, 0xff,
	0xff, 0xf3, 0x2c, 0xf6, 0x33, 0xcf, 0x00, 0x00, 0x00,
}
