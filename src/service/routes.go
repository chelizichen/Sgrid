package service

import (
	"Sgrid/src/config"
	handlers "Sgrid/src/http"
	"Sgrid/src/public"
	"Sgrid/src/storage"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/vo"
	utils "Sgrid/src/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"gopkg.in/yaml.v2"
)

const (
	TOKEN           = "e609d00404645feed1c1733835b8c127"
	CLUSTER_REQUEST = "CLUSTER_REQUEST"
	SINGLE_REQUEST  = "SINGLE_REQUEST"
)

func Registry(ctx *handlers.SgridServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
	GROUP.POST("/login", func(c *gin.Context) {
		token := c.PostForm("token")
		if token == TOKEN {
			c.JSON(http.StatusOK, handlers.Resp(0, "Ok", nil))
			return
		}
		c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error", nil))
	})

	GROUP.POST("/deletePackage", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		F := c.PostForm("fileName")
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		storagePath := filepath.Join(cwd, utils.PublishPath, serverName, F)
		err = os.Remove(storagePath)
		if err != nil {
			fmt.Println("Error To RemoveFile", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To RemoveFile", nil))
			return
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	GROUP.POST("/checkConfig", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		configProdPath := filepath.Join(utils.PublishPath, serverName, "simpProd.yaml")
		prod, err := public.NewConfig(public.WithTargetPath(configProdPath))
		if err != nil {
			fmt.Println("Error To Get NewConfig", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", prod))
	})

	GROUP.POST("/coverConfig", func(c *gin.Context) {
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
	})

	// tail -n rows log_file | grep "pattern"
	GROUP.POST("/getServerLog", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		fileName := c.PostForm("fileName")
		pattern := c.DefaultPostForm("pattern", "")
		rows := c.DefaultPostForm("rows", "100")
		sm, err := utils.NewSearchLogMonitor(serverName, fileName)
		if err != nil {
			fmt.Println("Error To New SimMonitor", err.Error())
		}
		s, err := sm.GetLogger(pattern, rows)
		if err != nil {
			fmt.Println("Error To GetLogger", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", s))
	})

	GROUP.POST("/getLogList", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils.PublishPath, serverName)
		D, err := os.ReadDir(serverPath)
		if err != nil {
			fmt.Println("Error To ReadDir", err.Error())
		}
		var loggers []string
		for i := 0; i < len(D); i++ {
			de := D[i]
			s := de.Name()
			b := strings.HasSuffix(s, ".log")
			if b {
				loggers = append(loggers, s)
			}
		}
		c.JSON(200, handlers.Resp(0, "ok", loggers))
	})

	GROUP.POST("/main/getLogList", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverPath := filepath.Join(cwd, "static/main")
		D, err := os.ReadDir(serverPath)
		if err != nil {
			fmt.Println("Error To ReadDir", err.Error())
		}
		var loggers []string
		for i := 0; i < len(D); i++ {
			de := D[i]
			s := de.Name()
			b := strings.HasSuffix(s, ".log")
			if b {
				loggers = append(loggers, s)
			}
		}
		c.JSON(200, handlers.Resp(0, "ok", loggers))
	})

	GROUP.POST("/main/getServerLog", func(c *gin.Context) {
		logFile := c.PostForm("logFile")
		pattern := c.DefaultPostForm("pattern", "")
		rows := c.DefaultPostForm("rows", "100")
		sm, err := utils.NewMainSearchLogMonitor(logFile)
		if err != nil {
			fmt.Println("Error To New SimMonitor", err.Error())
			c.JSON(200, handlers.Resp(-2, err.Error(), nil))
			return
		}
		s, err := sm.GetLogger(pattern, rows)
		if err != nil {
			fmt.Println("Error To GetLogger", err.Error())
			c.JSON(200, handlers.Resp(-1, err.Error(), nil))
			return
		}
		c.JSON(200, handlers.Resp(0, "ok", s))
	})

	GROUP.GET("/main/queryGrid", func(c *gin.Context) {
		pbr := utils.NewPageBaiscReq(c)
		gv := storage.QueryGrid(pbr)
		c.JSON(200, handlers.Resp(0, "ok", gv))
	})

	GROUP.GET("/main/queryServantGroup", func(c *gin.Context) {
		gv := storage.QueryServantGroup(&dto.PageBasicReq{})
		vgbs := storage.ConvertToVoGroupByServant(gv)
		c.JSON(200, handlers.Resp(0, "ok", vgbs))
	})

	GROUP.GET("/main/queryNodes", func(c *gin.Context) {
		nodes := storage.QueryNodes()
		c.JSON(200, handlers.Resp(0, "ok", nodes))
	})

	ctx.Engine.Use(GROUP.Handlers...)

}
