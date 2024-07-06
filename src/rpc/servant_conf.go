package rpc

import (
	"fmt"
	"regexp"
	"strconv"
)

// opt
type sgridOptCb[T any] func(s *SgridGrpcClient[T])

// conf
type SgridGrpcServantConfig struct {
	ServiceName string `json:"service_name"` // 服务名
	ServantName string `json:"servant_name"` // 服务实例名
	Protocol    string `json:"protocol"`     // 协议
	Host        string `json:"host"`         // 主机
	Port        int    `json:"port"`         // 端口
}

// rpc pack
type RequestPack struct {
	Method string
	Body   any
	Reply  any
}

// ************************* utils **************************

// server.SgridSubServer @grpc -h 127.0.0.1 -p 15996
func StringToProxy(str string) (*SgridGrpcServantConfig, error) {
	re := regexp.MustCompile(`(?P<ServiceName>\w+)\.(?P<ServantName>\w+)\w+@(?P<Protocol>\w+)\s+-h\s+(?P<Host>\d+\.\d+\.\d+\.\d+)\s+-p\s+(?P<Port>\d+)`)
	matches := re.FindStringSubmatch(str)

	if matches == nil {
		return nil, fmt.Errorf("input string does not match expected format")
	}

	result := &SgridGrpcServantConfig{}

	for i, name := range re.SubexpNames() {
		switch name {
		case "ServiceName":
			result.ServiceName = matches[i]
		case "ServantName":
			result.ServantName = matches[i]
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
