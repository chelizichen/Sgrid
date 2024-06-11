package service

import (
	protocol "Sgrid/server/SgridPackageServer/proto"
	"Sgrid/src/configuration"
	handlers "Sgrid/src/http"
	"Sgrid/src/public"
	clientgrpc "Sgrid/src/public/client_grpc"
	"Sgrid/src/storage"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/pojo"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

const (
	CronAssetsAdminKey = "cron|assets|admin"     // rds 管理资源Key
	CronAssetsSetValue = "cron|assets|set|value" // 值
	CronSepcTime       = "* 1 * * * *"           // 定时任务时间
	CronExpireTime     = time.Minute * 10        // 超时时间自动删
)

func AssetsService(ctx *handlers.SgridServerCtx) {
	rds := configuration.GRDB
	rds_ctx := configuration.RDBContext
	clients := ctx.Context.Value(
		public.GRPC_CLIENT_PROXYS{},
	).([]*clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient])
	var Job = func() {
		t := time.Now()
		fmt.Println("AssetsService.Job 开始加锁", t.Format(time.DateTime))
		s := rds.Get(rds_ctx, CronAssetsAdminKey).Val()
		if s == CronAssetsSetValue {
			Info := fmt.Sprintf("加锁失败%v", CronAssetsAdminKey)
			fmt.Println(Info)
			storage.PushErr(&pojo.SystemErr{
				Type: "system/error/AssetsService/c.AddJob",
				Info: Info,
			})
			return
		}
		rds.SetNX(rds_ctx, CronAssetsAdminKey, CronAssetsSetValue, CronExpireTime)
		sq := storage.QueryAssetsTime()
		var ids []int
		for _, aa := range sq {
			ids = append(ids, aa.GridId)
		}
		if len(ids) == 0 {
			return
		}
		grids := storage.BatchQueryGrid(ids)
		Req := []*protocol.ShutdownGridInfo{}
		for _, v := range grids {
			Req = append(Req, &protocol.ShutdownGridInfo{
				GridId: int32(v.GridId),
				Pid:    int32(v.Pid),
				Port:   int32(v.Port),
				Host:   v.Host,
			})
		}
		for _, client := range clients {
			client.GetClient().ShutdownGrid(ctx.Context, &protocol.ShutdownGridReq{
				Req: Req,
			})
		}
	}
	c := cron.Cron{}
	c.AddJob(CronSepcTime, cron.FuncJob(Job))

	router := ctx.Engine.Group(strings.ToLower(ctx.Name))

	router.POST("/getList", getList)

	router.POST("/upsertAsset", upsertAsset)

	router.GET("/delAssert", delAssert)
}

func getList(c *gin.Context) {
	var req *dto.PageBasicReq
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	resp, total, err := storage.QueryAssets(req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	handlers.AbortWithSuccList(c, resp, total)
}

func upsertAsset(c *gin.Context) {
	var req *pojo.AssetsAdmin
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	err = storage.UpsertAssets(req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	handlers.AbortWithSucc(c, nil)
}

func delAssert(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	err := storage.DelAssetById(id)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	handlers.AbortWithSucc(c, nil)
}
