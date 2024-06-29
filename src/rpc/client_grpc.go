package rpc

import (
	set "Sgrid/src/public/set"
	"fmt"
	"regexp"
	"strconv"

	"google.golang.org/grpc"
)

// grpc client pool
// support
// 1. balance request
// 2. singal req , multiple req
type SgridGrpcClient[T any] struct {
	_opts    []grpc.DialOption  // grpc统一可选配置
	_targets []string           // 链接地址
	_servant string             // 服务名
	_conns   []*grpc.ClientConn // grpc 客户端链接
}

func (s *SgridGrpcClient[T]) GetServantName() string {
	return s._servant
}

type sgridOptCb[T any] func(s *SgridGrpcClient[T])

func WithDiaoptions[T any](opts ...grpc.DialOption) sgridOptCb[T] {
	return func(s *SgridGrpcClient[T]) {
		s._opts = opts
	}
}

func WithTargets[T any](targets []string) sgridOptCb[T] {
	return func(s *SgridGrpcClient[T]) {
		sets := new(set.SgridSet[string])
		for _, v := range targets {
			sets.Add(v)
		}
		s._targets = sets.GetAll()
	}

}

func NewSgridGrpcClient[T any](clients []string, opts ...sgridOptCb[T]) (*SgridGrpcClient[T], error) {
	confs := []*SgridGrpcServantConfig{}
	targets := make([]string, 0)
	for _, v := range clients {
		prx, err := stringToProxy(v)
		if err != nil {
			return nil, fmt.Errorf("stringToProxy.error %v", err.Error())
		}
		confs = append(confs, prx)
		targets = append(targets, prx.Host)
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
	}
	c._servant = confs[0].ServiceName
	return c, nil
}

// server.SgridSubServer @grpc -h 127.0.0.1 -p 15996
func stringToProxy(str string) (*SgridGrpcServantConfig, error) {
	re := regexp.MustCompile(`(?P<ServiceName>\w+)\.\w+ @(?P<Protocol>\w+) -h (?P<Host>\d+\.\d+\.\d+\.\d+) -p (?P<Port>\d+)`)
	matches := re.FindStringSubmatch(str)

	if matches == nil {
		return nil, fmt.Errorf("input string does not match expected format")
	}

	result := &SgridGrpcServantConfig{}

	for i, name := range re.SubexpNames() {
		switch name {
		case "ServiceName":
			result.ServiceName = matches[i]
		case "Protocol":
			result.Protocol = matches[i]
		case "Host":
			result.Host = matches[i]
		case "Port":
			port, err := strconv.Atoi(matches[i])
			if err != nil {
				return nil, err
			}
			result.Port = port
		}
	}

	return result, nil
}
