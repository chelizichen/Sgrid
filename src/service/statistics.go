package service

import (
	protocol "Sgrid/server/SgridPackageServer/proto"
	h "Sgrid/src/http"
	"Sgrid/src/rpc"
	"Sgrid/src/storage"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/types/known/emptypb"
)

func SystemStatisticsRegisty(ctx *h.SgridServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.GetServerName()))
	packageServant := ctx.Context.Value(PackageServantProxy{}).(*rpc.SgridGrpcClient[protocol.FileTransferServiceClient])
	clients := packageServant.GetClients()

	GROUP.GET("/server/statistics/getByType", getByType)
	GROUP.GET("/server/statistics/getNodes", getNodesInfo(clients))
}

func getNodesInfo(clients []protocol.FileTransferServiceClient) func(c *gin.Context) {
	return func(c *gin.Context) {
		var wg sync.WaitGroup
		var result []protocol.SystemInfo
		for _, client := range clients {
			wg.Add(1)
			go func(clt protocol.FileTransferServiceClient) {
				rsp, _ := clt.GetSystemInfo(&gin.Context{}, &emptypb.Empty{})
				result = append(result, *rsp.Data)
				wg.Done()
			}(client)
		}
		wg.Wait()
		h.AbortWithSucc(c, result)
	}
}

func getByType(c *gin.Context) {
	t := c.Query("TYPE")
	f := storage.StatisticsMap[t]
	if f == nil {
		h.AbortWithError(c, "ERROR! 未找到该分类")
	}
	rsp, err := f()
	if err != nil {
		h.AbortWithError(c, "ERROR! "+err.Error())
		return
	}
	h.AbortWithSucc(c, rsp)

}
