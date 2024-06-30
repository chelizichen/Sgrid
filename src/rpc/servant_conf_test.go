package rpc_test

import (
	"Sgrid/src/rpc"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const traceServant = "server.SgridTraceServer@grpc -h 127.0.0.1 -p 15887"

func TestStringToProxy(t *testing.T) {
	sgsc, _ := rpc.StringToProxy(traceServant)
	fmt.Println(sgsc.Host)
	fmt.Println(sgsc.ServiceName)
	fmt.Println(sgsc.ServantName)
	fmt.Println(sgsc.Port)
	fmt.Println(sgsc.Protocol)
}

func TestConnTrace(t *testing.T) {
	proxys := make([]string, 0)
	proxys = append(proxys, traceServant)
	_, err := rpc.NewSgridGrpcClient[any](proxys, rpc.WithDiaoptions[any](grpc.WithTransportCredentials(insecure.NewCredentials())))
	assert.NoError(t, err)
}
