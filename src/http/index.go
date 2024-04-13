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

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron"
)

type SgridServerCtx struct {
	Port        string
	Name        string
	Engine      *gin.Engine
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

func (c *SgridServerCtx) Use(callback func(engine *SgridServerCtx)) {
	callback(c)
}

// 主控服务需要做日志系统与监控
func (c *SgridServerCtx) InitController() {

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
}

func (c *SgridServerCtx) Static(realPath string, args ...string) {
	s := c.Name
	pre := strings.ToLower(s)
	f := utils.Join(pre)
	target := f(realPath)

	staticPath := public.Join(args...)
	c.Engine.Static(target, staticPath)
}

type InitConf struct {
	SgridController    bool
	ServerType         string
	SgridConfPath      string
	SgridGinStaticPath string
	SgridGinWithCors   bool
}
type NewSgrid func(conf *InitConf)

func WithSgridController() NewSgrid {
	return func(conf *InitConf) {
		conf.SgridController = true
	}
}

func WithSgridServerType(T string) NewSgrid {
	return func(conf *InitConf) {
		conf.ServerType = T
	}
}

func WithConfPath(P string) NewSgrid {
	return func(conf *InitConf) {
		conf.SgridConfPath = P
	}
}

func WithSgridGinStatic(P string) NewSgrid {
	return func(conf *InitConf) {
		conf.SgridGinStaticPath = P
	}
}

func WithCors() NewSgrid {
	return func(conf *InitConf) {
		conf.SgridGinWithCors = true
	}
}

func NewSgridServerCtx(opt ...NewSgrid) *SgridServerCtx {
	initConf := &InitConf{}
	ctx := &SgridServerCtx{}
	for _, fn := range opt {
		fn(initConf)
	}
	conf, err := public.NewConfig(public.WithTargetPath(initConf.SgridConfPath))
	if err != nil {
		fmt.Println("NewConfig Error:", err.Error())
		panic(err.Error())
	}
	if initConf.ServerType == public.PROTOCOL_HTTP {
		ctx.Engine = gin.Default()
		ctx.Port = ":" + strconv.Itoa(conf.Server.Port)
		ctx.StoragePath = conf.Server.Storage
		ctx.StaticPath = conf.Server.StaticPath
		ctx.MapConf = conf.Server.MapConf
		ctx.Host = conf.Server.Host
		ctx.Name = conf.Server.Name
	}
	// 初始化控制器
	if initConf.SgridController {
		ctx.InitController()
	}
	if len(initConf.SgridGinStaticPath) != 0 {
		GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
		staticRoot := public.Join(ctx.StaticPath)
		GROUP.Static(initConf.SgridGinStaticPath, staticRoot)
		ctx.Engine.Use(GROUP.Handlers...)
	}

	if initConf.SgridGinWithCors {
		GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
		GROUP.Use(withCORSMiddleware())
		ctx.Engine.Use(GROUP.Handlers...)
	}

	if err != nil {
		fmt.Println("init db error", err.Error())
	}
	if err != nil {
		fmt.Println("Error! database ping ", err.Error())
	}
	return ctx
}

func NewSgridServer(ctx *SgridServerCtx, callback func(port string)) {
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

	SIMP_TARGET_PORT := os.Getenv(public.ENV_TARGET_PORT)
	SIMP_CONF_PORT := ctx.Port
	var CallBackPort string = ""
	if SIMP_TARGET_PORT != "" {
		CallBackPort = ":" + SIMP_TARGET_PORT
	} else {
		CallBackPort = SIMP_CONF_PORT
	}

	callback(CallBackPort)
}

func withCORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
