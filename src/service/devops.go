package service

import (
	handlers "Sgrid/src/http"
	sgridError "Sgrid/src/public/error"
	"Sgrid/src/storage"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/pojo"
	utils "Sgrid/src/utils"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func DevopsService(ctx *handlers.SgridServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
	// devops component
	GROUP.GET("/devops/getGroups", getGroups)
	// 1.选择服务组 ｜ 创建
	GROUP.POST("/devops/saveGroup", saveGroup)
	GROUP.POST("/devops/delGroup", delGroup)
	// 2.创建服务
	GROUP.POST("/devops/saveServant", saveServant)
	GROUP.POST("/devops/updateServant", updateServant)
	GROUP.POST("/devops/delServant", delServant)

	// 2.1 选择服务
	GROUP.GET("/devops/getServants", getServants)

	// 3.选择节点
	GROUP.GET("/devops/queryNodes", queryNodes)
	GROUP.POST("/devops/saveNode", saveNode)
	// 4.添加至服务网格
	GROUP.POST("/devops/saveGrid", saveGrid)
	GROUP.POST("/devops/deleteGrid", deleteGrid)

	// 配置中心
	GROUP.GET("/devops/getConfig", getConfig)
	GROUP.POST("/devops/updateConfig", updateConfig)

	// 属性配置中心
	GROUP.POST("/devops/getPropertys", getPropertys)
	GROUP.POST("/devops/setProperty", setProperty)
	GROUP.GET("/devops/delProperty", delProperty)

	// 服务组
	GROUP.GET("/main/queryGrid", queryGrid)
	GROUP.POST("/main/queryServantGroup", queryServantGroup)
	GROUP.GET("/main/queryNodes", mainQueryNodes)

	GROUP.GET("/main/port/random", getRandomPort)
}

func getGroups(c *gin.Context) {
	user_id := c.DefaultQuery("id", "0")
	userId, err := strconv.Atoi(user_id)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	vsg := storage.QueryGroups(userId)
	handlers.AbortWithSucc(c, vsg)
}

func getServants(c *gin.Context) {
	user_id := c.DefaultQuery("id", "0")
	userId, err := strconv.Atoi(user_id)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	vsg := storage.QueryServants(userId)
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
		Id:             req.Id,
		TagName:        req.TagName,
		TagEnglishName: req.TagEnglishName,
		CreateTime:     &now,
		UserId:         req.UserId,
	}
	vsg := storage.SaveServantGroup(&record)
	handlers.AbortWithSucc(c, vsg)
}

func delGroup(c *gin.Context) {
	group_id := c.DefaultQuery("id", "0")
	groupId, err := strconv.Atoi(group_id)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	err = storage.DelGroup(groupId)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	handlers.AbortWithSucc(c, nil)
}

func saveServant(c *gin.Context) {
	var req *dto.SaveServantDto
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	now := time.Now()
	isExist := storage.GetServantByNameAndGroup(req.ServerName, req.ServantGroupId)
	if isExist {
		handlers.AbortWithError(c, "该服务组下已存在同名服务")
		return
	}
	record := &pojo.Servant{
		ServerName:     req.ServerName,
		Language:       req.Language,
		Protocol:       req.Protocol,
		ServantGroupId: req.ServantGroupId,
		ExecPath:       req.ExecPath,
		CreateTime:     &now,
		UserId:         req.UserId,
	}
	vsg := storage.SaveServant(record)
	handlers.AbortWithSucc(c, vsg)
}

func updateServant(c *gin.Context) {
	var req *dto.SaveServantDto
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	record := &pojo.Servant{
		Language:       req.Language,
		Protocol:       req.Protocol,
		ServantGroupId: req.ServantGroupId,
		ExecPath:       req.ExecPath,
	}
	vsg := storage.UpdateServant(record)
	handlers.AbortWithSucc(c, vsg)
}

func delServant(c *gin.Context) {
	servant_id := c.Query("id")
	stat := c.Query("stat")
	servantId, err := strconv.Atoi(servant_id)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	Stat, err := strconv.Atoi(stat)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	record := &pojo.Servant{
		Id:   servantId,
		Stat: Stat,
	}
	vsg := storage.DelServant(record)
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
	servant_id := c.Query("id")
	servantId, err := strconv.Atoi(servant_id)
	if err != nil {
		fmt.Println("err", err.Error())
		handlers.AbortWithError(c, err.Error())
		return
	}
	sc := storage.GetServantConfById(servantId)
	handlers.AbortWithSucc(c, sc)
}

func updateConfig(c *gin.Context) {
	var req *pojo.ServantConf
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("err", err.Error())
		handlers.AbortWithError(c, err.Error())
	}
	storage.UpdateConf(req)
	handlers.AbortWithSucc(c, nil)
}

// ** property **

func getPropertys(c *gin.Context) {
	// var req *pojo.Properties
	// if err := c.BindJSON(&req); err != nil {
	// 	fmt.Println("err", err.Error())
	// 	handlers.AbortWithError(c, err.Error())
	// }
	p := storage.QueryProperties()
	handlers.AbortWithSucc(c, p)
}

func setProperty(c *gin.Context) {
	var req *pojo.Properties
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("err", err.Error())
		handlers.AbortWithError(c, err.Error())
	}
	storage.UpsertProperty(req)
	handlers.AbortWithSucc(c, nil)
}

func delProperty(c *gin.Context) {
	servant_id, _ := strconv.Atoi(c.Query("id"))
	storage.DelProperty(servant_id)
	handlers.AbortWithSucc(c, nil)
}

func queryGrid(c *gin.Context) {
	pbr := utils.NewPageBaiscReq(c)
	gv := storage.QueryGrid(pbr)
	c.JSON(200, handlers.Resp(0, "ok", gv))
}

func queryServantGroup(c *gin.Context) {
	var req *dto.PageBasicReq
	if err := c.BindJSON(&req); err != nil {
		fmt.Println("err", err.Error())
		handlers.AbortWithError(c, err.Error())
	}
	gv := storage.QueryServantGroup(req)
	vgbs := storage.ConvertToVoGroupByServant(gv)
	c.JSON(200, handlers.Resp(0, "ok", vgbs))
}

func mainQueryNodes(c *gin.Context) {
	nodes := storage.QueryNodes()
	c.JSON(200, handlers.Resp(0, "ok", nodes))
}

func getRandomPort(c *gin.Context) {
	registryPorts := storage.GetAllPort()
	maxPort := 25000
	minPort := 10000
	port := 0
	rand.NewSource(time.Now().UnixNano()) // 确保每次运行时随机种子不同
	for {
		port = rand.Intn(maxPort-minPort+1) + minPort // 在指定范围内生成随机端口
		if port%10 != 0 {                             // 确保端口最后一位不是0
			found := false
			for _, registeredPort := range registryPorts {
				if port == registeredPort {
					found = true
					break
				}
			}
			if !found { // 检查端口是否已注册
				handlers.AbortWithSucc(c, port)
				break
			}
		}
	}
}
