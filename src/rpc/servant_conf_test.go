package rpc_test

import (
	protocol "Sgrid/server/SgridPackageServer/proto"
	"Sgrid/src/rpc"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const packageServer = "server.SgridPackageServer@grpc -h 127.0.0.1 -p 14938"

func TestStringToProxy(t *testing.T) {
	sgsc, _ := rpc.StringToProxy(packageServer)
	fmt.Println(sgsc.Host)
	fmt.Println(sgsc.ServiceName)
	fmt.Println(sgsc.ServantName)
	fmt.Println(sgsc.Port)
	fmt.Println(sgsc.Protocol)
}

func TestConnTrace(t *testing.T) {
	proxys := make([]string, 0)
	proxys = append(proxys, packageServer)
	_, err := rpc.NewSgridGrpcClient[any](proxys, rpc.WithDiaoptions[any](grpc.WithTransportCredentials(insecure.NewCredentials())))
	assert.NoError(t, err)
}

func TestInvoke(t *testing.T) {
	proxys := make([]string, 0)
	proxys = append(proxys, packageServer)
	prx, err := rpc.NewSgridGrpcClient[protocol.FileTransferServiceClient](proxys,
		rpc.WithRequestPrefix[protocol.FileTransferServiceClient]("/SgridProtocol.FileTransferService/"),
		rpc.WithDiaoptions[protocol.FileTransferServiceClient](
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	assert.NoError(t, err)
	var rsp protocol.GetLogFileByHostResp
	err = prx.Request(rpc.RequestPack{
		// methodName 在 grpc桩文件里面可以找到
		Method: "GetLogFileByHost",
		Body: &protocol.GetLogFileByHostReq{
			Host:       "127.0.0.1",
			ServerName: "ShopServer",
			GridId:     28,
		},
		Reply: &rsp,
	})
	assert.NoError(t, err)
	fmt.Println("rsp.Code", rsp.Code)
	fmt.Println("rsp.Data", rsp.Data)
	fmt.Println("rsp.Message", rsp.Message)
}

// 验证 相同的Target会不会被Set过滤 Pass
// targets [127.0.0.1:14938 127.0.0.1:14938]
// s._targets [127.0.0.1:14938]
func TestSet(t *testing.T) {
	proxys := make([]string, 0)
	proxys = append(proxys, packageServer, packageServer)
	_, err := rpc.NewSgridGrpcClient[protocol.FileTransferServiceClient](proxys,
		rpc.WithRequestPrefix[protocol.FileTransferServiceClient]("/SgridProtocol.FileTransferService/"),
		rpc.WithDiaoptions[protocol.FileTransferServiceClient](
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	assert.NoError(t, err)
}
