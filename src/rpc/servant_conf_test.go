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
		rpc.WithDiaoptions[protocol.FileTransferServiceClient](
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	assert.NoError(t, err)
	var rsp protocol.GetLogFileByHostResp
	err = prx.Request(rpc.RequestPack{
		Method: "/SgridProtocol.FileTransferService/GetLogFileByHost",
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
