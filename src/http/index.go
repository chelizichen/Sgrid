package http

import (
	"Sgrid/src/config"
	"Sgrid/src/public"
	"Sgrid/src/public/servant"
	"Sgrid/src/utils"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type SgridServerCtx struct {
	Port       string
	Name       string
	Engine     *gin.Engine
	Host       string
	ServerConf map[string]interface{}
	SgridConf  *config.SgridConf
	Context    context.Context
}

func Resp(code int, message string, data interface{}) *gin.H {
	return &gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	}
}

func (c *SgridServerCtx) RegistryHttpRouter(callback func(engine *SgridServerCtx)) {
	callback(c)
}

func (c *SgridServerCtx) RegistrySubServer(registry servant.SgridRegistryServiceInf) {
	pkgPath := reflect.TypeOf(registry).Elem().PkgPath()
	packagePath := strings.ReplaceAll(pkgPath, "/", ".")
	if !strings.HasSuffix(packagePath, registry.NameSpace()) {
		panic("sgrid/error package path is not equal with server path")
	}
	go registry.Registry(config.GlobalConf.Servant[registry.NameSpace()])
}

func (c *SgridServerCtx) Static(realPath string, filePath string) {
	s := c.Name
	pre := strings.ToLower(s)
	f := utils.Join(pre)
	target := f(realPath)

	staticPath := public.Join(filePath)
	c.Engine.Static(target, staticPath)
}

type InitConf struct {
	SgridController    bool
	ServerType         string
	SgridConfPath      string
	SgridGinStaticPath [2]string
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

func WithSgridGinStatic(paths [2]string) NewSgrid {
	return func(conf *InitConf) {
		conf.SgridGinStaticPath = paths
	}
}

func WithCors() NewSgrid {
	return func(conf *InitConf) {
		conf.SgridGinWithCors = true
	}
}

func NewSgridServerCtx(opt ...NewSgrid) *SgridServerCtx {
	initConf := &InitConf{}
	ctx := &SgridServerCtx{
		Context: context.Background(),
	}
	for _, fn := range opt {
		fn(initConf)
	}
	conf, err := public.NewConfig(public.WithTargetPath(initConf.SgridConfPath))
	config.GlobalConf = conf
	ctx.SgridConf = conf
	if err != nil {
		fmt.Println("NewConfig Error:", err.Error())
		panic(err.Error())
	}
	ctx.ServerConf = conf.Conf
	ctx.Host = conf.Server.Host
	ctx.Name = conf.Server.Name
	ctx.Port = ":" + strconv.Itoa(conf.Server.Port)

	if initConf.ServerType == public.PROTOCOL_HTTP {
		ctx.Engine = gin.Default()
		if len(initConf.SgridGinStaticPath) != 0 {
			GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
			staticRoot := public.Join(initConf.SgridGinStaticPath[1])
			GROUP.Static(initConf.SgridGinStaticPath[0], staticRoot)
		}
		if initConf.SgridGinWithCors {
			GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
			GROUP.Use(withCORSMiddleware())
		}
	}
	if initConf.ServerType == public.PROTOCOL_GRPC {
		ctx.Engine = nil
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

var AbortWithError = func(c *gin.Context, err string) {
	c.AbortWithStatusJSON(http.StatusOK, Resp(-1, "err ::"+err, nil))
}

// Done
var AbortWithSucc = func(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, Resp(0, "ok", data))
}

// List
var AbortWithSuccList = func(c *gin.Context, data interface{}, total int64) {
	var listVo = make(map[string]interface{}, 0)
	listVo["total"] = total
	listVo["data"] = data
	c.AbortWithStatusJSON(http.StatusOK, Resp(0, "ok", data))
}
