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
	"os/exec"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
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
	// 上传服务，并且判断是不是集群请求

	GROUP.POST("/restartServer", func(c *gin.Context) {
		// event-stream 返回数据
		c.Request.Response.Header.Set("Content-Type", "text/event-stream")
		c.Request.Response.Header.Set("Cache-Control", "no-cache")
		c.Request.Response.Header.Set("Connection", "keep-alive")
		releaseType := c.DefaultPostForm("releaseType", utils.RELEASE_SINGLENODE) // 集群模式下需要指定Type 默认普通单节点发布
		fileName := c.PostForm("fileName")                                        // 文件名称
		serverName := c.PostForm("serverName")                                    // 服务名称
		targetPort := c.PostForm("targetPort")                                    // 集群模式下需要指定端口
		isSame := utils.ConfirmFileName(serverName, fileName)
		if !isSame {
			msg := "Error File!" + fileName + "  | " + serverName
			fmt.Println(msg)
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, msg, nil))
			return
		}

		svr := utils.GetServant(serverName, targetPort)
		err := svr.StopServant()
		if err != nil {
			msg := "Error StopServant!" + fileName + "  | " + serverName
			fmt.Println(msg)
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, msg, nil))
			return
		}

		cwd, err := os.Getwd()
		if err != nil {
			msg := "Error To GetWd :" + err.Error()
			fmt.Printf("msg: %v\n", msg)
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, msg, nil))
			return
		}
		storageYmlEPath := utils.GetFilePath(cwd, serverName, utils.DevConfEntry)
		storageYmlProdPath := utils.GetFilePath(cwd, serverName, utils.ProdConfEntry)
		storagePath := utils.GetFilePath(cwd, serverName, fileName)
		dest := filepath.Join(cwd, utils.PublishPath, serverName)

		isFirstRelease := !utils.IsExist(storageYmlProdPath)
		var runScript func() *exec.Cmd
		var confPort int
		if isFirstRelease {
			err = utils.Unzip(storagePath, dest) // 直接解压
			if err != nil {
				fmt.Println("Error To Unzip", err.Error())
			}
			err = public.CopyProdYml(storageYmlEPath, storageYmlProdPath) // 拷贝配置文件
			if err != nil {
				fmt.Println("Error To CopyProdYml", err.Error())
			}
			sc, err := public.NewConfig(public.WithTargetPath(storageYmlProdPath)) // 引入配置文件
			if err != nil {
				fmt.Println("Error To NewConfig", err.Error())
			}
			confPort = sc.Server.Port
			svr.Language = sc.Server.Type
			if sc.Server.Type == utils.RELEASE_TYPE_NODEJS {
				runScript = func() *exec.Cmd {
					storageNodePath := utils.GetFilePath(cwd, serverName, utils.NodeJsEntry)
					var cmd *exec.Cmd = exec.Command("node", storageNodePath)
					return cmd
				}
			} else if sc.Server.Type == utils.RELEASE_TYPE_GO {
				runScript = func() *exec.Cmd {
					storageExEPath := utils.GetFilePath(cwd, serverName, utils.GoEntry)
					var cmd *exec.Cmd = exec.Command(storageExEPath)
					return cmd
				}
			} else if sc.Server.Type == utils.RELEASE_TYPE_JAVA {
				// java -jar your-application.jar -Dspring.config.location=file:/path/to/application.yml,file:/path/to/another-config.yaml
				runScript = func() *exec.Cmd {
					storageJavaPath := utils.GetFilePath(cwd, serverName, utils.SpringEntry)
					var cmd *exec.Cmd = exec.Command("java", "-jar", storageJavaPath, "-Dspring.config.location="+"file:"+storageYmlProdPath)
					return cmd
				}
			}
		} else if releaseType == utils.RELEASE_SINGLENODE {
			var clearScript func(sc config.SgridConf)
			sc, err := public.NewConfig(public.WithTargetPath(storageYmlProdPath))
			svr.Language = sc.Server.Type

			if svr.Language == utils.RELEASE_TYPE_NODEJS {
				storageNodePath := utils.GetFilePath(cwd, serverName, utils.NodeJsEntry)
				clearScript = func(sc config.SgridConf) {
					confPort = sc.Server.Port
					storageStaticPath := utils.GetFilePath(cwd, serverName, sc.Server.StaticPath)
					err = utils.IFExistThenRemove(storageStaticPath, true)
					if err != nil {
						fmt.Println("remove File Error storageStaticPath "+storageStaticPath, err.Error())
					}
					err = utils.IFExistThenRemove(storageNodePath, false)
					if err != nil {
						fmt.Println("remove File Error storageNodePath "+storageNodePath, err.Error())
					}
					err = utils.Unzip(storagePath, dest) // 直接解压
					if err != nil {
						fmt.Println("Error To Unzip", err.Error())
					}
				}
				runScript = func() *exec.Cmd {
					var cmd *exec.Cmd = exec.Command("node", storageNodePath)
					return cmd
				}
			} else if svr.Language == utils.RELEASE_TYPE_GO {
				storageExEPath := utils.GetFilePath(cwd, serverName, utils.GoEntry)
				clearScript = func(sc config.SgridConf) {
					confPort = sc.Server.Port
					storageStaticPath := utils.GetFilePath(cwd, serverName, sc.Server.StaticPath)
					err = utils.IFExistThenRemove(storageStaticPath, true)
					if err != nil {
						fmt.Println("remove File Error storageStaticPath "+storageStaticPath, err.Error())
					}
					err = utils.IFExistThenRemove(storageExEPath, false)
					if err != nil {
						fmt.Println("remove File Error storageExEPath "+storageExEPath, err.Error())
					}
					err = utils.Unzip(storagePath, dest) // 直接解压
					if err != nil {
						fmt.Println("Error To Unzip", err.Error())
					}
				}
				runScript = func() *exec.Cmd {
					var cmd *exec.Cmd = exec.Command(storageExEPath)
					return cmd
				}
			} else if svr.Language == utils.RELEASE_TYPE_JAVA {
				clearScript = func(sc config.SgridConf) {
					confPort = sc.Server.Port
					storageStaticPath := utils.GetFilePath(cwd, serverName, sc.Server.StaticPath)
					err = utils.IFExistThenRemove(storageStaticPath, true)
					if err != nil {
						fmt.Println("remove File Error storageStaticPath "+storageStaticPath, err.Error())
					}
					err = utils.Unzip(storagePath, dest) // 直接解压
					if err != nil {
						fmt.Println("Error To Unzip", err.Error())
					}
				}
				// java -jar your-application.jar -Dspring.config.location=file:/path/to/application.yml,file:/path/to/another-config.yaml
				runScript = func() *exec.Cmd {
					storageJavaPath := utils.GetFilePath(cwd, serverName, utils.SpringEntry)
					var cmd *exec.Cmd = exec.Command("java", "-jar", storageJavaPath, "-D", "spring.config.location="+"file:"+storageYmlProdPath)
					return cmd
				}
			}
			confPort = sc.Server.Port
			if err != nil {
				fmt.Println("Error To NewConfig", err.Error())
			}
			err = utils.IFExistThenRemove(storageYmlEPath, false)
			if err != nil {
				fmt.Println("remove File Error storageYmlEPath "+storageYmlEPath, err.Error())
			}
			clearScript(*sc) // 执行清除
		} else if releaseType == utils.RELEASE_CLUSTER {
			sc, err := public.NewConfig(public.WithTargetPath(storageYmlProdPath))
			if err != nil {
				fmt.Println("err NewConfig", err.Error())
			}
			svr.Language = sc.Server.Type
			if svr.Language == utils.RELEASE_TYPE_NODEJS {
				runScript = func() *exec.Cmd {
					storageNodePath := utils.GetFilePath(cwd, serverName, utils.NodeJsEntry)
					var cmd *exec.Cmd = exec.Command("node", storageNodePath)
					return cmd
				}
			} else if svr.Language == utils.RELEASE_TYPE_GO {
				runScript = func() *exec.Cmd {
					storageExEPath := utils.GetFilePath(cwd, serverName, utils.GoEntry)
					var cmd *exec.Cmd = exec.Command(storageExEPath)
					return cmd
				}
			} else if svr.Language == utils.RELEASE_TYPE_JAVA {
				runScript = func() *exec.Cmd {
					storageJavaPath := utils.GetFilePath(cwd, serverName, utils.SpringEntry)
					var cmd *exec.Cmd = exec.Command("java", "-jar", storageJavaPath, "-Dspring.config.location="+"file:"+storageYmlProdPath)
					return cmd
				}
			}
		}

		script := runScript()
		stdoutPipe, err := script.StdoutPipe()
		if err != nil {
			fmt.Println("Error Get StdoutPiper", err.Error())
		}
		stderrPipe, err := script.StderrPipe()
		if err != nil {
			fmt.Println("Error Get stderrPipe", err.Error())
		}
		// 设置环境变量
		script.Env = append(os.Environ(), "SIMP_PRODUCTION=Yes", "SIMP_SERVER_PATH="+dest)
		if svr.Language == utils.RELEASE_TYPE_JAVA {
			script.Env = append(script.Env, "SIMP_PROD_YAML="+storageYmlProdPath)
		}
		if releaseType == utils.RELEASE_CLUSTER {
			script.Env = append(script.Env, "SIMP_TARGET_PORT="+targetPort, "SIMP_SERVER_INDEX="+targetPort)
		} else {
			script.Env = append(script.Env, "SIMP_SERVER_INDEX=1", fmt.Sprintf("SIMP_TARGET_PORT=%v", confPort))
		}
		if err != nil {
			msg := "Error To New Monitor" + err.Error()
			fmt.Printf("msg: %v\n", msg)
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, msg, nil))
			return
		}
		err = script.Start()

		if err != nil {
			msg := "Error To EXEC Cmd Start ：" + err.Error()
			fmt.Println(msg)
			fmt.Println("Cmd", script.Args)
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, msg, nil))
		}
		var exit atomic.Bool
		cronSvr := cron.New()
		svr.Process = script
		svr.Pid = script.Process.Pid
		svr.ExitSignal = exit
		svr.Cron = cronSvr

		// 启动一个协程，用于读取并打印命令的输出
		go func() {
			spec := "0 0 0 * * *"
			sm, err := utils.NewSimpMonitor(serverName, "", targetPort)
			if err != nil {
				fmt.Println("Error To New Monitor", err.Error())
			}
			// 添加定时任务
			err = cronSvr.AddFunc(spec, func() {

				newSM, err := utils.NewSimpMonitor(serverName, "", targetPort)
				if err != nil {
					fmt.Println("Error To New Monitor", err.Error())
					return
				}
				sm = newSM
			})
			if err != nil {
				fmt.Println("AddFuncErr", err)
			}

			// 启动Cron调度器
			go cronSvr.Start()

			go func() {
				for {
					if exit.Load() {
						return
					}
					// 读取输出
					buf := make([]byte, 1024)
					s := time.Now().Format(time.DateTime)
					n, err := stdoutPipe.Read(buf)
					if err != nil {
						break
					}
					// 打印输出
					content := s + " ServerName " + serverName + " || " + string(buf[:n]) + "\n"
					sm.AppendLogger(content)
				}
			}()
			go func() {
				for {
					if exit.Load() {
						return
					}
					// 读取输出
					buf := make([]byte, 1024)
					s := time.Now().Format(time.DateTime)
					n, err := stderrPipe.Read(buf)
					if err != nil {
						break
					}
					// 打印输出
					content := s + " Error : ServerName " + serverName + " || " + string(buf[:n]) + "\n"
					fmt.Println(content)
					sm.AppendLogger(content)
				}
			}()

			go func() {
				for {
					content, isAlive := svr.ServantMonitor()
					if !isAlive {
						break
					}
					sm.AppendLogger(content)
				}
			}()
		}()
		v := make(map[string]interface{}, 2)
		v["pid"] = svr.Pid
		v["status"] = true
		utils.SubServants[svr.GetContextName()] = svr
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", v))
	})

	GROUP.POST("/getServers", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverPath := filepath.Join(cwd, utils.PublishPath)
		fmt.Println("serverPath", serverPath)
		subdirectories, err := utils.GetSubdirectories(serverPath)
		if err != nil {
			fmt.Println("Error To GetSubdirectories")
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", subdirectories))
	})

	GROUP.POST("/createServer", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To GetWd", nil))
			return
		}
		value := c.PostForm("serverName")
		fmt.Println("createServer | serverName ", value)
		serverPath := filepath.Join(cwd, utils.PublishPath, value)
		utils.AutoCreateLoggerFile(value)
		err = os.Mkdir(serverPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error To Mkdir", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To Mkdir", nil))
			return
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	GROUP.POST("/getServerPackageList", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "Error To GetWd", nil))
			return
		}
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils.PublishPath, serverName)
		fmt.Println("serverPath", serverPath)
		var packages []utils.ReleasePackageVo
		err = filepath.Walk(serverPath, utils.VisitTgzS(&packages, serverName))
		if err != nil {
			fmt.Printf("error walking the path %v: %v\n", serverPath, err)
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", packages))
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

	GROUP.POST("/checkServer", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		port := c.PostForm("port")
		s := utils.GetServant(serverName, port)
		b := utils.IsPidAlive(s.Pid)
		v := make(map[string]interface{}, 10)
		v["status"] = false
		if b {
			v["pid"] = s.Pid
			v["status"] = true
		}
		c.JSON(http.StatusOK, handlers.Resp(0, "ok", v))
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

	GROUP.POST("/deleteServer", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		ErrorToRemoveAll := "Error To Remove All"
		ErrorToGetServerName := "Error To Get ServerName"
		if serverName == "" {
			c.AbortWithStatusJSON(200, handlers.Resp(-1, ErrorToGetServerName, nil))
			return
		}
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverPath := filepath.Join(cwd, utils.PublishPath, serverName)
		fmt.Println("DeleteDirectory", serverPath)
		err = utils.DeleteDirectory(serverPath)
		if err != nil {
			fmt.Println(ErrorToRemoveAll, err.Error())
			c.AbortWithStatusJSON(200, handlers.Resp(-1, ErrorToRemoveAll, nil))
			return
		}
		c.JSON(200, handlers.Resp(0, "ok", nil))
	})

	GROUP.POST("/getChildStats", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		if serverName == "" {
			c.AbortWithStatusJSON(200, handlers.Resp(-1, "missing params serverName", nil))
			return
		}
		m := make(map[string]map[string]interface{})
		for sName, ctx := range utils.SubServants {
			context := &ctx
			if strings.HasPrefix(sName, serverName) {
				pid := context.Pid
				name := sName
				b := utils.IsPidAlive(pid)
				v := make(map[string]interface{}, 10)
				v["status"] = false
				if b {
					v["pid"] = pid
					v["status"] = true
				}

				m[name] = v
			}
		}
		c.AbortWithStatusJSON(200, handlers.Resp(0, "ok", m))
	})

	GROUP.POST("/shutdownServer", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		port := c.PostForm("port")
		s := utils.GetServant(serverName, port)

		if s.Pid == 0 {
			c.AbortWithStatusJSON(200, handlers.Resp(-1, "暂无PID", nil))
			return
		}
		err := s.StopServant()
		if err != nil {
			fmt.Println("x:", err)
			c.AbortWithStatusJSON(200, handlers.Resp(-1, "关闭服务异常", err.Error()))
			return
		}
		delete(utils.SubServants, s.GetContextName())
		c.AbortWithStatusJSON(200, handlers.Resp(0, "ok", nil))
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

	GROUP.POST("/getApiJson", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils.PublishPath, serverName, "API.json")
		Content, err := os.ReadFile(serverPath)
		if err != nil {
			fmt.Println("Error To ReadFile", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", string(Content)))
	})

	GROUP.POST("/getDoc", func(c *gin.Context) {
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		serverName := c.PostForm("serverName")
		serverPath := filepath.Join(cwd, utils.PublishPath, serverName, "doc.txt")
		Content, err := os.ReadFile(serverPath)
		if err != nil {
			fmt.Println("Error To ReadFile", err.Error())
		}
		c.JSON(200, handlers.Resp(0, "ok", string(Content)))
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
		return

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
