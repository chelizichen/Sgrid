// ************************ SgridCloud **********************
// grpc_client created at 2024.4.15
// Author @chelizichen
// grpc client invoke easily
// ************************ SgridCloud **********************

package clientgrpc

import (
	protocol "Sgrid/server/SgridPackageServer/proto"
	"Sgrid/src/storage"
	"Sgrid/src/storage/pojo"
	"context"
	"fmt"
	"net/url"
	"reflect"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type WithSgridGrpcOptFunc[T any] func(*SgridGrpcClient[T])

func WithSgridGrpcClientAddress[T any](address string) WithSgridGrpcOptFunc[T] {
	return func(c *SgridGrpcClient[T]) {
		c.serverAddress = address
	}
}

func WithSgridGrpcConn[T any](conn *grpc.ClientConn) WithSgridGrpcOptFunc[T] {
	return func(c *SgridGrpcClient[T]) {
		c.conn = conn
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
	conn          *grpc.ClientConn
}

func (c *SgridGrpcClient[T]) GetClient() T {
	return *c.client
}

func (c *SgridGrpcClient[T]) GetConn() *grpc.ClientConn {
	return c.conn
}

func (c *SgridGrpcClient[T]) GetAddress() string {
	return c.serverAddress
}

func (c *SgridGrpcClient[T]) ParseHost(pre string) (*url.URL, error) {
	s := fmt.Sprintf("%v://%v", pre, c.GetAddress())
	fmt.Println("s", s)
	return url.Parse(s)
}

type WithSgridGrpcGroupOptFunc[T any] func(*SgridGrpcClientGroup[T])

type SgridGrpcClientGroup[T any] struct {
	clients   []*SgridGrpcClient[T]
	nameSpace string
	address   []string
	grpcOpt   []grpc.DialOption
	newClient NewClient[T]
	context   *context.Context
}

func (s *SgridGrpcClientGroup[T]) ReqRandom() {
	for _, v := range s.clients {
		go func(client SgridGrpcClient[T]) {
			client.GetClient()
		}(*v)
	}
}

func (s *SgridGrpcClientGroup[T]) ReqAll(methodName string, req interface{}) []reflect.Value {
	var wg sync.WaitGroup
	var invokeResponse []reflect.Value
	for _, v := range s.clients {
		wg.Add(1)
		go func(client SgridGrpcClient[T]) {
			fn := reflect.ValueOf(client.GetClient()).MethodByName(methodName)
			resp := fn.Call([]reflect.Value{reflect.ValueOf(*s.context), reflect.ValueOf(req)})
			invokeResponse = append(invokeResponse, resp...)
			wg.Done()
		}(*v)
	}
	wg.Wait()
	return invokeResponse
}

func NewSgridGrpcClientGroup[T any](ctx context.Context, clientConn NewClient[T], opts ...WithSgridGrpcGroupOptFunc[T]) *SgridGrpcClientGroup[T] {
	opts = append(opts, withSgridGrpcClientGroupNewFn[T](clientConn), withSgridGrpcClientGroupCtx[T](&ctx))
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

func withSgridGrpcClientGroupNewFn[T any](fn NewClient[T]) WithSgridGrpcGroupOptFunc[T] {
	return func(c *SgridGrpcClientGroup[T]) {
		c.newClient = fn
	}
}

func WithSgridGrpcClientGroupAddress[T any](address []string) WithSgridGrpcGroupOptFunc[T] {

	return func(c *SgridGrpcClientGroup[T]) {
		c.address = address
	}
}

func WithSgridGrpcClientGroupNameSpace[T any](nameSpace string) WithSgridGrpcGroupOptFunc[T] {
	return func(c *SgridGrpcClientGroup[T]) {
		c.nameSpace = nameSpace
	}
}

func WithSgridGrpcClientGroupNewConnOpt[T any](opts ...grpc.DialOption) WithSgridGrpcGroupOptFunc[T] {
	return func(c *SgridGrpcClientGroup[T]) {
		c.grpcOpt = opts
	}
}

func withSgridGrpcClientGroupCtx[T any](ctx *context.Context) WithSgridGrpcGroupOptFunc[T] {
	return func(c *SgridGrpcClientGroup[T]) {
		c.context = ctx
	}
}

// todo SampleInvoke
// reflect to invoke
func SampleInvoke() {
	addresses := []string{"localhost:14938"}
	ctx := context.Background()
	g := NewSgridGrpcClientGroup[protocol.FileTransferServiceClient](
		ctx,
		protocol.NewFileTransferServiceClient,
		WithSgridGrpcClientGroupAddress[protocol.FileTransferServiceClient](addresses), // 请求
		WithSgridGrpcClientGroupNewConnOpt[protocol.FileTransferServiceClient](grpc.WithTransportCredentials(insecure.NewCredentials())),
		WithSgridGrpcClientGroupNameSpace[protocol.FileTransferServiceClient]("server.SgridPackageServer"),
	)
	all := g.ReqAll("DeletePackage", &protocol.DeletePackageReq{})
	for i, v := range all {
		fmt.Println(i, v)
	}
}

func ProxyInvoke(g *SgridGrpcClientGroup[any], methodName string, args interface{}) interface{} {
	v := g.ReqAll(methodName, args)
	resu := make([]interface{}, len(g.address))
	for index, resp := range v {
		resu[index] = resp
	}
	return resu
}

func NewSgridGrpcProxyConn[T any](Key string, ClientConn func(cc grpc.ClientConnInterface) T) []*SgridGrpcClient[T] {
	gn := storage.QueryPropertiesByKey(Key)
	addresses := []string{}
	for _, v := range gn {
		addresses = append(addresses, v.Value)
	}
	fmt.Println("NewSgridGrpcProxyConn address", addresses)
	clients := []*SgridGrpcClient[T]{}
	for _, v := range addresses {
		add := v
		conn, err := grpc.Dial(add,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                60 * time.Second, // 发送ping消息的时间间隔
				Timeout:             20 * time.Second, // 等待ping响应的超时时间
				PermitWithoutStream: true,             // 即使没有活动的流，也允许发送ping消息
			}),
			grpc.WithConnectParams(grpc.ConnectParams{
				MinConnectTimeout: 5 * time.Second,
			}),
		)
		fmt.Println("client_grpc.getstate", conn.GetState().String())
		if err != nil {
			info := fmt.Sprintf("NewSgridGrpcProxyConn address %v %v", addresses, err)
			storage.PushErr(&pojo.SystemErr{
				Type: "system/error/grpc/dial",
				Info: info,
			})
		}
		client := NewSgridClient[T](
			ClientConn(conn),
			WithSgridGrpcClientAddress[T](add),
			WithSgridGrpcConn[T](conn),
		)
		defer func() {
			if err := recover(); err != nil {
				fmt.Println("err", err)
			}
		}()
		clients = append(clients, client)
	}
	return clients
}
