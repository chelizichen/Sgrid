// ************************ SgridCloud **********************
// grpc_client created at 2024.4.15
// Author @chelizichen
// grpc client invoke easily
// ************************ SgridCloud **********************

package clientgrpc

import (
	"fmt"
	"net/url"

	"google.golang.org/grpc"
)

type WithSgridGrpcOptFunc[T any] func(*SgridGrpcClient[T])

func WithSgridGrpcClientAddress[T any](address string) WithSgridGrpcOptFunc[T] {
	return func(c *SgridGrpcClient[T]) {
		c.serverAddress = address
	}
}

func WithSgridGrpcNameSpace[T any](nameSpace string) WithSgridGrpcOptFunc[T] {
	return func(c *SgridGrpcClient[T]) {
		c.nameSpace = nameSpace
	}
}

func NewSgridClient[T any](client T, opt ...WithSgridGrpcOptFunc[T]) *SgridGrpcClient[T] {
	c := &SgridGrpcClient[T]{
		client: &client,
	}
	for _, fn := range opt {
		fn(c)
	}
	return c
}

type SgridGrpcClient[T any] struct {
	serverAddress string
	nameSpace     string
	client        *T
}

func (c *SgridGrpcClient[T]) GetClient() T {
	return *c.client
}

func (c *SgridGrpcClient[T]) GetAddress() string {
	return c.serverAddress
}

func (c *SgridGrpcClient[T]) ParseHost() (*url.URL, error) {
	s := c.GetAddress()
	return url.Parse(s)
}

type WithSgridGrpcGroupOptFunc[T any] func(*SgridGrpcClientGroup[T])

type SgridGrpcClientGroup[T any] struct {
	clients   []*SgridGrpcClient[T]
	nameSpace string
	address   []string
	grpcOpt   []grpc.DialOption
	newClient NewClient[T]
}

func NewSgridGrpcClientGroup[T any](opts ...WithSgridGrpcGroupOptFunc[T]) *SgridGrpcClientGroup[T] {
	s := &SgridGrpcClientGroup[T]{}
	for _, fn := range opts {
		fn(s)
	}
	if len(s.address) == 0 {
		fmt.Println("error! address length is 0")
		return nil
	}
	if s.newClient == nil {
		fmt.Println("error! missing client proxy function")
		return nil
	}
	for _, _v := range s.address {
		v := _v
		conn, err := grpc.Dial(v, s.grpcOpt...)
		if err != nil {
			fmt.Println("err", err.Error())
			break
		}
		client := NewSgridClient[T](
			s.newClient(conn),
			WithSgridGrpcClientAddress[T](v),
		)
		s.clients = append(s.clients, client)
	}
	return s
}

type NewClient[T any] func(cc grpc.ClientConnInterface) T

func WithSgridGrpcNewClientFn[T any](fn NewClient[T]) WithSgridGrpcGroupOptFunc[T] {
	return func(c *SgridGrpcClientGroup[T]) {
		c.newClient = fn
	}
}

func WithSgridGrpcClientGroupAddress[T any](address []string) WithSgridGrpcGroupOptFunc[T] {

	return func(c *SgridGrpcClientGroup[T]) {
		c.address = address
	}
}

func WithSgridGrpcClientNameSpace[T any](nameSpace string) WithSgridGrpcGroupOptFunc[T] {
	return func(c *SgridGrpcClientGroup[T]) {
		c.nameSpace = nameSpace
	}
}

func WithSgridGrpcConnOpt[T any](opts ...grpc.DialOption) WithSgridGrpcGroupOptFunc[T] {
	return func(c *SgridGrpcClientGroup[T]) {
		c.grpcOpt = opts
	}
}
