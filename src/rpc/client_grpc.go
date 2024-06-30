package rpc

import (
	set "Sgrid/src/public/set"
	"context"
	"fmt"
	"sync"
	"sync/atomic"

	"google.golang.org/grpc"
)

// grpc client pool
// support
// 1. balance request
// 2. singal req , multiple req
type SgridGrpcClient[T any] struct {
	_opts    []grpc.DialOption  // grpc统一可选配置
	_targets []string           // 链接地址
	_servant string             // 服务实例
	_service string             // 服务名
	_conns   []*grpc.ClientConn // grpc 客户端链接
	_total   uint8              // 总共
	_curr    atomic.Uint32
	ctx      context.Context
}

func (s *SgridGrpcClient[T]) GetServantName() string {
	return s._servant
}

// balance request
func (s *SgridGrpcClient[T]) Request(pack RequestPack) error {
	u := s._curr.Add(1)
	if u >= uint32(s._total) {
		s._curr.Store(0)
	}
	curr := s._curr.Load()
	err := s._conns[curr].Invoke(s.ctx, pack.Method, pack.Body, pack.Reply)
	return err
}

func (s *SgridGrpcClient[T]) RequestAll(pack RequestPack, replys []any) (reply []any, err error) {
	var wg sync.WaitGroup
	for i, cc := range s._conns {
		err = cc.Invoke(s.ctx, pack.Method, pack.Body, replys[i])
		wg.Add(1)
	}
	wg.Wait()
	return replys, err
}

func WithDiaoptions[T any](opts ...grpc.DialOption) sgridOptCb[T] {
	return func(s *SgridGrpcClient[T]) {
		s._opts = append(s._opts, opts...)
	}
}

func WithTargets[T any](targets []string) sgridOptCb[T] {
	return func(s *SgridGrpcClient[T]) {
		sets := make(set.SgridSet[string])
		fmt.Println("targets", targets)
		for _, v := range targets {
			sets.Add(v)
		}
		s._targets = sets.GetAll()
		fmt.Println("s._targets", s._targets)
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
		c._total += 1
		fmt.Println("grpc.dial success ", v)
	}
	c._servant = confs[0].ServantName
	c._service = confs[0].ServiceName
	c.ctx = context.Background()
	return c, nil
}
