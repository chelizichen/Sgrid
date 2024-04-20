// devops component
// 1.选择服务组 ｜ 创建
// 2.创建服务
// 3.选择节点
// 4.添加至服务网格
// 5.选择端口
package service

import (
	"Sgrid/src/config"
	handlers "Sgrid/src/http"
	"Sgrid/src/public"
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
	GROUP.GET("/devops/getGroups", getGroups)
	GROUP.POST("/devops/saveGroup", saveGroup)
	GROUP.POST("/devops/saveServant", saveServant)
	GROUP.GET("/devops/queryNodes", queryNodes)
	GROUP.POST("/devops/saveGrid", saveGrid)

	// todo 中心数据库 or 文件形式
	GROUP.GET("/devops/getConfig", getConfig)
	GROUP.POST("/devops/updateConfig", updateConfig)

}

func getGroups(c *gin.Context) {
	vsg := storage.QueryGroups()
	c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", vsg))
}

func saveGroup(c *gin.Context) {
	var req *dto.SaveServantGroupDto
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "err"+err.Error(), nil))
		return
	}
	now := time.Now()
	record := pojo.ServantGroup{
		TagName:        req.TagName,
		TagEnglishName: req.TagEnglishName,
		CreateTime:     &now,
	}
	vsg := storage.SaveServantGroup(&record)
	c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", vsg))
}

func saveServant(c *gin.Context) {
	var req *dto.SaveServantDto
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "err"+err.Error(), nil))
		return
	}
	now := time.Now()
	record := &pojo.Servant{
		ServerName:     req.ServerName,
		Language:       req.Language,
		Protocol:       req.Protocol,
		ServantGroupId: req.ServantGroupId,
		CreateTime:     &now,
	}
	vsg := storage.SaveServant(record)
	c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", vsg))
}

func queryNodes(c *gin.Context) {
	gn := storage.QueryNodes()
	c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", gn))
}

func saveGrid(c *gin.Context) {
	var req *dto.GridDTO
	err := c.BindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "err"+err.Error(), nil))
		return
	}
	now := time.Now()
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

// **************** Conf ***************

func getConfig(c *gin.Context) {
	serverName := c.PostForm("serverName")
	configProdPath := filepath.Join(utils.PublishPath, serverName, "simpProd.yaml")
	prod, err := public.NewConfig(public.WithTargetPath(configProdPath))
	if err != nil {
		fmt.Println("Error To Get NewConfig", err.Error())
	}
	c.JSON(200, handlers.Resp(0, "ok", prod))
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
