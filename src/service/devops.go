package service

import (
	protocol "Sgrid/server/SgridPackageServer/proto"
	"Sgrid/src/config"
	handlers "Sgrid/src/http"
	"Sgrid/src/public"
	clientgrpc "Sgrid/src/public/client_grpc"
	sgridError "Sgrid/src/public/error"
	"Sgrid/src/storage"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/pojo"
	"Sgrid/src/storage/vo"
	utils "Sgrid/src/utils"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

func DevopsService(ctx *handlers.SgridServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
	// devops component
	GROUP.GET("/devops/getGroups", getGroups)
	// 1.选择服务组 ｜ 创建
	GROUP.POST("/devops/saveGroup", saveGroup)
	// 2.创建服务
	GROUP.POST("/devops/saveServant", saveServant)

	// 2.1 选择服务
	GROUP.GET("/devops/getServants", getServants)

	// 3.选择节点
	GROUP.GET("/devops/queryNodes", queryNodes)
	GROUP.POST("/devops/saveNode", saveNode)
	// 4.添加至服务网格
	GROUP.POST("/devops/saveGrid", saveGrid)
	GROUP.POST("/devops/deleteGrid", deleteGrid)
	// 5.选择端口

	// 中心数据库
	GROUP.GET("/devops/getConfig", getConfig)
	GROUP.POST("/devops/updateConfig", updateConfig)
	// 集群通知变更
	GROUP.POST("/devops/notifyConfig", notifyConfig(ctx.Context.Value(public.GRPC_CLIENT_PROXYS{}).([]*clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient])))

}

func getGroups(c *gin.Context) {
	vsg := storage.QueryGroups()
	handlers.AbortWithSucc(c, vsg)
}

func getServants(c *gin.Context) {
	vsg := storage.QueryServants()
	handlers.AbortWithSucc(c, vsg)
}

func saveGroup(c *gin.Context) {
	var req *dto.SaveServantGroupDto
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	now := time.Now()
	record := pojo.ServantGroup{
		TagName:        req.TagName,
		TagEnglishName: req.TagEnglishName,
		CreateTime:     &now,
	}
	vsg := storage.SaveServantGroup(&record)
	handlers.AbortWithSucc(c, vsg)
}

func saveServant(c *gin.Context) {
	var req *dto.SaveServantDto
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	now := time.Now()
	record := &pojo.Servant{
		ServerName:     req.ServerName,
		Language:       req.Language,
		Protocol:       req.Protocol,
		ServantGroupId: req.ServantGroupId,
		ExecPath:       req.ExecPath,
		CreateTime:     &now,
	}
	vsg := storage.SaveServant(record)
	handlers.AbortWithSucc(c, vsg)
}

func queryNodes(c *gin.Context) {
	gn := storage.QueryNodes()
	c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", gn))
}

func deleteGrid(c *gin.Context) {
	var req *dto.PageBasicReq
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	if req.Id == 0 {
		handlers.AbortWithError(c, sgridError.Request_Error(" missing Id ").Error())
		return
	}
	storage.DeleteGrid(req.Id)
	handlers.AbortWithSucc(c, nil)
}

func saveGrid(c *gin.Context) {
	var req *dto.GridDTO
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())

	}
	now := time.Now()
	count := storage.GetGridByNodePort(req.NodeId, req.Port)
	if count > 0 {
		handlers.AbortWithError(c, sgridError.Request_Error(" port already exist").Error())
		return
	}
	if req.Port == 0 {
		handlers.AbortWithError(c, sgridError.Request_Error(" missing arg port").Error())
		return
	}
	record := &pojo.Grid{
		CreateTime: &now,
		NodeId:     req.NodeId,    // 节点ID
		ServantId:  req.ServantId, // 服务ID
		Port:       req.Port,
		Status:     0,
		Pid:        0,
	}
	i := storage.UpdateGrid(record)
	c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", i))
}

func saveNode(c *gin.Context) {
	var req *dto.NodeDTO
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	i := storage.UpdateNode(req)
	handlers.AbortWithSucc(c, i)
}

// **************** Conf ***************

func getConfig(c *gin.Context) {
	serverName := c.PostForm("serverName")
	configProdPath := filepath.Join(utils.PublishPath, serverName, "simpProd.yaml")
	prod, err := public.NewConfig(public.WithTargetPath(configProdPath))
	if err != nil {
		fmt.Println("Error To Get NewConfig", err.Error())
	}
	handlers.AbortWithSucc(c, prod)
}

func updateConfig(c *gin.Context) {
	var reqVo vo.CoverConfigVo
	if err := c.BindJSON(&reqVo); err != nil {
		c.JSON(http.StatusOK, handlers.Resp(0, "-1", err.Error()))
		return
	}
	serverName := reqVo.ServerName
	uploadConfig := reqVo.Conf
	if serverName == "" {
		fmt.Println("Server Name is Empty")
		c.JSON(http.StatusOK, handlers.Resp(0, "Server Name is Empty", nil))
		return
	}
	marshal, err := yaml.Marshal(uploadConfig)
	if err != nil {
		fmt.Println("Error To Stringify config", err.Error())
		c.JSON(http.StatusOK, handlers.Resp(0, "Error To Stringify config", nil))
		return
	}
	fmt.Println("serverName", serverName)
	fmt.Println("uploadConfig", string(marshal))
	if len(marshal) == 0 {
		fmt.Println("Error To Stringify config", err.Error())
		c.JSON(http.StatusOK, handlers.Resp(0, "Error To Stringify config", nil))
		return
	}
	configPath := filepath.Join(utils.PublishPath, serverName, "simpProd.yaml")
	err = config.CoverConfig(string(marshal), configPath)
	if err != nil {
		fmt.Println("CoverConfig Error", err.Error())
		c.JSON(200, handlers.Resp(-1, "CoverConfig Error", nil))
	}
	c.JSON(200, handlers.Resp(0, "ok", nil))
}

func notifyConfig(clients []*clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient]) func(c *gin.Context) {
	return func(c *gin.Context) {
		var reqVo vo.CoverConfigVo
		if err := c.BindJSON(&reqVo); err != nil {
			c.JSON(http.StatusOK, handlers.Resp(0, "-1", err.Error()))
			return
		}
		serverName := reqVo.ServerName
		uploadConfig := reqVo.Conf
		if serverName == "" {
			fmt.Println("Server Name is Empty")
			c.JSON(http.StatusOK, handlers.Resp(0, "Server Name is Empty", nil))
			return
		}
		marshal, err := yaml.Marshal(uploadConfig)
		if err != nil {
			fmt.Println("Error To Stringify config", err.Error())
			c.JSON(http.StatusOK, handlers.Resp(0, "Error To Stringify config", nil))
			return
		}
		fmt.Println("serverName", serverName)
		fmt.Println("uploadConfig", string(marshal))
		if len(marshal) == 0 {
			fmt.Println("Error To Stringify config", err.Error())
			c.JSON(http.StatusOK, handlers.Resp(0, "Error To Stringify config", nil))
			return
		}
		configPath := filepath.Join(utils.PublishPath, serverName, "simpProd.yaml")
		err = config.CoverConfig(string(marshal), configPath)
		if err != nil {
			fmt.Println("CoverConfig Error", err.Error())
			c.JSON(200, handlers.Resp(-1, "CoverConfig Error", nil))
		}
		c.JSON(200, handlers.Resp(0, "ok", nil))
	}
}
