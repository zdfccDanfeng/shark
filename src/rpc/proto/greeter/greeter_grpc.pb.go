// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package helloworld

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// GreeterClient is the client API for Greeter service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GreeterClient interface {
	// Sends a greeting
	Say2Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error)
}

type greeterClient struct {
	cc grpc.ClientConnInterface
}

func NewGreeterClient(cc grpc.ClientConnInterface) GreeterClient {
	return &greeterClient{cc}
}

var greeterSay2HelloStreamDesc = &grpc.StreamDesc{
	StreamName: "Say2Hello",
}

func (c *greeterClient) Say2Hello(ctx context.Context, in *HelloRequest, opts ...grpc.CallOption) (*HelloReply, error) {
	out := new(HelloReply)
	err := c.cc.Invoke(ctx, "/helloworld.Greeter/Say2Hello", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GreeterService is the service API for Greeter service.
// Fields should be assigned to their respective handler implementations only before
// RegisterGreeterService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type GreeterService struct {
	// Sends a greeting
	Say2Hello func(context.Context, *HelloRequest) (*HelloReply, error)
}

func (s *GreeterService) say2Hello(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.Say2Hello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/helloworld.Greeter/Say2Hello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.Say2Hello(ctx, req.(*HelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterGreeterService registers a service implementation with a gRPC server.
func RegisterGreeterService(s grpc.ServiceRegistrar, srv *GreeterService) {
	srvCopy := *srv
	if srvCopy.Say2Hello == nil {
		srvCopy.Say2Hello = func(context.Context, *HelloRequest) (*HelloReply, error) {
			return nil, status.Errorf(codes.Unimplemented, "method Say2Hello not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "helloworld.Greeter",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "Say2Hello",
				Handler:    srvCopy.say2Hello,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "greeter.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewGreeterService creates a new GreeterService containing the
// implemented methods of the Greeter service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewGreeterService(s interface{}) *GreeterService {
	ns := &GreeterService{}
	if h, ok := s.(interface {
		Say2Hello(context.Context, *HelloRequest) (*HelloReply, error)
	}); ok {
		ns.Say2Hello = h.Say2Hello
	}
	return ns
}

// UnstableGreeterService is the service API for Greeter service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstableGreeterService interface {
	// Sends a greeting
	Say2Hello(context.Context, *HelloRequest) (*HelloReply, error)
}