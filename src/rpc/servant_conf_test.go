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

const errorStringPackageServer1 = "server.SgridPackageServer@grpc      -h  127.0.0.1-p 14938"
const successStringPacakgeServer2 = "server.SgridPackageServer@grpc -h  127.0.0.1  -p 14938"
const successStringPacakgeServer3 = "server.SgridPackageServer@grpc  -h  127.0.0.1   -p 14938"

func TestStringToProxy(t *testing.T) {
	_, err := rpc.StringToProxy(successStringPacakgeServer2)
	assert.Nil(t, err)
	_, err = rpc.StringToProxy(successStringPacakgeServer3)
	assert.Nil(t, err)
	_, err = rpc.StringToProxy(errorStringPackageServer1)
	assert.NotNil(t, err)
}

func TestConnTrace(t *testing.T) {
	proxys := make([]string, 0)
	proxys = append(proxys, successStringPacakgeServer2)
	_, err := rpc.NewSgridGrpcClient[any](proxys, rpc.WithDiaoptions[any](grpc.WithTransportCredentials(insecure.NewCredentials())))
	assert.NoError(t, err)
}

func TestInvoke(t *testing.T) {
	proxys := make([]string, 0)
	proxys = append(proxys, successStringPacakgeServer2)
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
	proxys = append(proxys, successStringPacakgeServer2, successStringPacakgeServer2)
	_, err := rpc.NewSgridGrpcClient[protocol.FileTransferServiceClient](proxys,
		rpc.WithRequestPrefix[protocol.FileTransferServiceClient]("/SgridProtocol.FileTransferService/"),
		rpc.WithDiaoptions[protocol.FileTransferServiceClient](
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
	)
	assert.NoError(t, err)
}
