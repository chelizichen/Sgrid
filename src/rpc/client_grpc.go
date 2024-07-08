package rpc

import (
	set "Sgrid/src/public/set"
	"context"
	"fmt"
	"net/url"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
)

// grpc client pool
// support
// 1. balance request
// 2. singal req , multiple req
// 3. support diy req
type SgridGrpcClient[T any] struct {
	_opts           []grpc.DialOption   // grpc统一可选配置
	_targets        []string            // 链接地址
	_servant        string              // 服务实例
	_service        string              // 服务名
	_conns          []*grpc.ClientConn  // grpc 客户端链接
	_total          uint8               // 总共
	_curr           *atomic.Uint32      // 负载均衡
	_prefix         string              // 服务前缀
	_ctx            context.Context     // 上下文
	_client_service NewClientService[T] // 客户端接口
	_clients        []T                 // 代理服务组
}

func (s *SgridGrpcClient[T]) GetServantName() string {
	return s._servant
}

func (s *SgridGrpcClient[T]) GetContext() context.Context {
	return s._ctx
}

func (s *SgridGrpcClient[T]) GetClients() []T {
	return s._clients
}

func (s *SgridGrpcClient[T]) ParseHost(idx int, pre string) (*url.URL, error) {
	h := fmt.Sprintf("%v://%v", pre, s._targets[idx])
	fmt.Println("h", h)
	return url.Parse(h)
}

func (s *SgridGrpcClient[T]) getFullRequestName(name string) string {
	return s._prefix + name
}

// balance request
func (s *SgridGrpcClient[T]) Request(pack RequestPack) error {
	fmt.Println("SgridGrpcClient.Request |", s)
	u := s._curr.Add(1)
	if u >= 1024 {
		s._curr.Store(0)
	}
	current := u % uint32(s._total)
	err := s._conns[current].Invoke(
		s._ctx,
		s.getFullRequestName(pack.Method),
		pack.Body,
		pack.Reply,
	)
	return err
}

func (s *SgridGrpcClient[T]) RequestAll(pack RequestPack, reply []any) (err error) {
	var wg sync.WaitGroup
	for i, cc := range s._conns {
		err = cc.Invoke(s._ctx, pack.Method, pack.Body, reply[i])
		wg.Add(1)
	}
	wg.Wait()
	return err
}

func WithDiaoptions[T any](opts ...grpc.DialOption) sgridOptCb[T] {
	return func(s *SgridGrpcClient[T]) {
		s._opts = append(s._opts, opts...)
	}
}

func WithRequestPrefix[T any](prefix string) sgridOptCb[T] {
	return func(s *SgridGrpcClient[T]) {
		s._prefix = prefix
	}
}

type NewClientService[T any] func(cc grpc.ClientConnInterface) T

func WithClientService[T any](fn NewClientService[T]) sgridOptCb[T] {
	return func(s *SgridGrpcClient[T]) {
		s._client_service = fn
	}
}

func WithTargets[T any](targets []string) sgridOptCb[T] {
	return func(s *SgridGrpcClient[T]) {
		sets := make(set.SgridSet[string])
		// fmt.Println("targets", targets)
		for _, v := range targets {
			sets.Add(v)
		}
		s._targets = sets.GetAll()
		// fmt.Println("s._targets", s._targets)
	}
}

func NewSgridGrpcClient[T any](clients []string, opts ...sgridOptCb[T]) (*SgridGrpcClient[T], error) {
	confs := []*SgridGrpcServantConfig{}
	targets := make([]string, 0)
	for _, v := range clients {
		prx, err := StringToProxy(v)
		if err != nil {
			return nil, fmt.Errorf("stringToProxy.error %v", err.Error())
		}
		confs = append(confs, prx)
		targets = append(targets, fmt.Sprintf("%s:%d", prx.Host, prx.Port))
	}
	opts = append(opts, WithTargets[T](targets))
	c := new(SgridGrpcClient[T])
	// init
	for _, opt := range opts {
		opt(c)
	}
	for _, v := range c._targets {
		cc, err := grpc.Dial(v, c._opts...)
		if err != nil {
			return nil, fmt.Errorf("grpc.Dial error %v", err.Error())
		}
		cc.Connect()
		c._conns = append(c._conns, cc)
		t := c._client_service(cc)
		c._clients = append(c._clients, t)
		c._total += 1
		fmt.Println("grpc.dial success ", v)
	}
	c._servant = confs[0].ServantName
	c._service = confs[0].ServiceName
	c._ctx = context.Background()
	c._curr = &atomic.Uint32{}
	return c, nil
}
