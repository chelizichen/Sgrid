package http

import (
	"Sgrid/src/public"
	"Sgrid/src/utils"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/robfig/cron"
)

type SimpHttpServerCtx struct {
	Port        string
	Name        string
	Engine      *gin.Engine
	Storage     *sqlx.DB
	isMain      bool
	StoragePath string
	StaticPath  string
	Host        string
	MapConf     map[string]interface{}
}

func Resp(code int, message string, data interface{}) *gin.H {
	return &gin.H{
		"Code":    code,
		"Message": message,
		"Data":    data,
	}
}

func (c *SimpHttpServerCtx) Use(callback func(engine *SimpHttpServerCtx)) {
	callback(c)
}

// 主控服务需要做日志系统与监控
func (c *SimpHttpServerCtx) DefineMain() {

	cwd, err := os.Getwd()
	if err != nil {
		Err_Message := "Error To Get Wd" + err.Error()
		fmt.Println(Err_Message)
		panic(Err_Message)
	}
	staticPath := filepath.Join(cwd, "static")
	s := utils.IsExist(filepath.Join(cwd, "static/main"))
	if !s {
		err = os.Mkdir(staticPath, os.ModePerm)
		if err != nil {
			Err_Message := "Error To Mkdir" + err.Error()
			fmt.Println(Err_Message)
			panic(Err_Message)
		}
	}
	mainPath := filepath.Join(cwd, "static/main")
	b := utils.IsExist(filepath.Join(cwd, "static/main"))
	if !b {
		err = os.Mkdir(mainPath, os.ModePerm)
		if err != nil {
			Err_Message := "Error To Mkdir" + err.Error()
			fmt.Println(Err_Message)
			panic(Err_Message)
		}
	}
	go func() {
		c := cron.New()

		// 4小时执行一次，更换日志文件指定目录
		spec := "0 0 0 * * *"
		utils.AutoSetLogWriter()
		// 添加定时任务
		err := c.AddFunc(spec, func() {
			utils.AutoSetLogWriter()
		})
		if err != nil {
			fmt.Println("AddFuncErr", err)
		}
		// 启动Cron调度器
		go c.Start()
	}()
	c.isMain = true
}

func (c *SimpHttpServerCtx) UseSPA(path string, root string) {
	wd, _ := os.Getwd()
	SIMP_PRODUCTION := os.Getenv("SIMP_PRODUCTION")
	s := c.Name
	pre := strings.ToLower(s)
	f := utils.Join(pre)

	// 设置缓存
	setCacheHeaders := func(ctx *gin.Context, fileInfo os.FileInfo) {
		ctx.Header("Cache-Control", "public, max-age=2592000")
		expires := time.Now().Add(time.Hour * 24 * 30)
		ctx.Header("Expires", expires.Format(time.RFC1123))
		lastModified := fileInfo.ModTime()
		ctx.Header("Last-Modified", lastModified.Format(time.RFC1123))
	}

	c.Engine.GET(f(path)+"/*path", func(ctx *gin.Context) {
		requestPath := ctx.Param("path")
		var webRoot, targetPath string

		if SIMP_PRODUCTION == "Yes" {
			SIMP_SERVER_PATH := os.Getenv("SIMP_SERVER_PATH")
			webRoot = filepath.Join(SIMP_SERVER_PATH, root)
			targetPath = filepath.Join(SIMP_SERVER_PATH, root, requestPath)
		} else {
			webRoot = filepath.Join(wd, root)
			targetPath = filepath.Join(wd, root, requestPath)
		}

		if _, err := os.Stat(targetPath); os.IsNotExist(err) {
			targetPath = filepath.Join(webRoot, "index.html")
		} else if fileInfo, err := os.Stat(targetPath); err == nil {
			if strings.HasSuffix(targetPath, ".js") || strings.HasSuffix(targetPath, ".css") {
				setCacheHeaders(ctx, fileInfo)
			}
		}

		ctx.File(targetPath)
	})
}

func (c *SimpHttpServerCtx) Static(realPath string, args ...string) {
	s := c.Name
	pre := strings.ToLower(s)
	f := utils.Join(pre)
	target := f(realPath)

	wd, _ := os.Getwd()
	SIMP_PRODUCTION := os.Getenv("SIMP_PRODUCTION")
	var staticPath string
	if SIMP_PRODUCTION == "Yes" {
		SIMP_SERVER_PATH := os.Getenv("SIMP_SERVER_PATH")
		if len(args) > 0 {
			otherPath := filepath.Join(args...)
			staticPath = filepath.Join(SIMP_SERVER_PATH, otherPath)
		} else {
			staticPath = filepath.Join(SIMP_SERVER_PATH, c.StaticPath)
		}
	} else {
		if len(args) > 0 {
			otherPath := filepath.Join(args...)
			staticPath = filepath.Join(wd, otherPath)
		} else {
			staticPath = filepath.Join(wd, c.StaticPath)
		}
	}
	c.Engine.Static(target, staticPath)
}

func NewSgridHttpServerCtx(G *gin.Engine) (ctx *SimpHttpServerCtx) {
	conf, err := public.NewConfig()
	if err != nil {
		fmt.Println("NewConfig Error:", err.Error())
		panic(err.Error())
	}
	if err != nil {
		fmt.Println("get Config Error :", err.Error())
	}
	database, err := sqlx.Open("mysql", conf.Server.Storage)
	if err != nil {
		fmt.Println("init db error", err.Error())
	}
	Storage := database
	err = database.Ping()
	if err != nil {
		fmt.Println("Error! database ping ", err.Error())
	}
	ctx = &SimpHttpServerCtx{
		Name:        conf.Server.Name,
		Port:        ":" + strconv.Itoa(conf.Server.Port),
		Engine:      G,
		StoragePath: conf.Server.Storage,
		StaticPath:  conf.Server.StaticPath,
		Storage:     Storage,
		MapConf:     conf.Server.MapConf,
		Host:        conf.Server.Host,
	}
	return ctx
}

func NewSgridHttpServer(ctx *SimpHttpServerCtx, callback func(port string)) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	for _, addr := range addrs {
		// Check if the address is not a loopback address
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				fmt.Println("IPv4 Address:", ipnet.IP.String())
			} else {
				fmt.Println("IPv6 Address:", ipnet.IP.String())
			}
		}
	}
	SIMP_PRODUCTION := os.Getenv("SIMP_PRODUCTION")
	fmt.Println("SIMP_PRODUCTION", SIMP_PRODUCTION)
	// 子服务生产时需要提供对应的API路由
	if SIMP_PRODUCTION == "Yes" {
		fmt.Println("CreateAPIFile |", ctx.Name)
		utils.CreateAPIFile(ctx.Engine, ctx.Name)
	}
	SIMP_TARGET_PORT := os.Getenv("SIMP_TARGET_PORT")
	SIMP_CONF_PORT := ctx.Port
	var CallBackPort string = ""
	if SIMP_TARGET_PORT != "" {
		CallBackPort = ":" + SIMP_TARGET_PORT
	} else {
		CallBackPort = SIMP_CONF_PORT
	}
	callback(CallBackPort)
}
