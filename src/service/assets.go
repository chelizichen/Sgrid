// 设置上下线的时间
// 上线时间可以不填，不填默认为当前时间
// 下限时间可以不填，默认一直上线
// 两者必须填一个
/*
 * @LastEditTime: 2024-06-14
 * @Description: 资源上下线的定时任务
 * @FilePath: /Sgrid/src/service/assets.go
 */

package service

import (
	protocol "Sgrid/server/SgridPackageServer/proto"
	handlers "Sgrid/src/http"
	"Sgrid/src/pool"
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
	CronSepcTime       = "@every 3m"             // 定时任务时间
	CronExpireTime     = time.Minute * 1         // 超时时间自动删
)

func AssetsService(ctx *handlers.SgridServerCtx) {
	cronJob(ctx)
	router := ctx.Engine.Group(strings.ToLower(ctx.GetServerName()))
	router.POST("/assets/admin/getList", getList)
	router.POST("/assets/admin/upsertAsset", upsertAsset)
	router.GET("/assets/admin/delAssert", delAssert)
	router.GET("/assets/admin/getAssert", getAssert)
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
	fmt.Println("total", total)
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

func getAssert(c *gin.Context) {
	id, _ := strconv.Atoi(c.Query("id"))
	resp, err := storage.GetAssetById(id)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	handlers.AbortWithSucc(c, resp)
}

func cronJob(ctx *handlers.SgridServerCtx) {
	rds := pool.GRDB
	rds_ctx := pool.RDBContext
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
		if false != true {
			sq := storage.QueryNeedShutDownAssets()
			var ids []int
			for _, aa := range sq {
				ids = append(ids, aa.GridId)
			}
			fmt.Println("storage.QueryNeedShutDownAssets.ids", ids)
			if len(ids) != 0 {
				grids := storage.BatchQueryGridByStat(ids, 1)
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
		}
		if true != false {
			np := storage.QueryNeedPullAssets()
			var ids []int
			for _, aa := range np {
				ids = append(ids, aa.GridId)
			}
			fmt.Println("storage.QueryNeedPullAssets.ids", ids)
			if len(ids) != 0 {
				grids := storage.BatchQueryGridByStat(ids, 0)
				Req := []*protocol.PatchServerDto{}
				for _, v := range grids {
					Req = append(Req, &protocol.PatchServerDto{
						ServerName:     v.ServerName,
						GridId:         int32(v.GridId),
						ServantId:      int32(v.ServantId),
						ExecPath:       v.ExecPath,
						Host:           v.Host,
						Port:           int32(v.Port),
						ServerLanguage: v.ServerLanguage,
						ServerProtocol: v.ServerProtocol,
					})
				}
				for _, client := range clients {
					client.GetClient().PatchServer(ctx.Context, &protocol.PatchServerReq{
						Req: Req,
					})
				}
			}
		}
	}
	c := cron.New()
	c.AddFunc(CronSepcTime, Job)
	go c.Start()
}
