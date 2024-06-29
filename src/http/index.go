package http

import (
	"Sgrid/src/config"
	"Sgrid/src/public"
	"Sgrid/src/rpc"
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
)

type SgridServerCtx struct {
	Engine    *gin.Engine
	SgridConf *config.SgridConf
	Context   context.Context
}

func (c *SgridServerCtx) GetServerName() string {
	return c.SgridConf.Server.Name
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

func (c *SgridServerCtx) RegistrySubServer(registry rpc.SgridSubServer) {
	pkgPath := reflect.TypeOf(registry).Elem().PkgPath()
	packagePath := strings.ReplaceAll(pkgPath, "/", ".")
	if !strings.HasSuffix(packagePath, registry.NameSpace()) {
		panic("sgrid/error package path is not equal with server path")
	}
	go registry.Registry(config.GlobalConf.Servant[registry.NameSpace()])
}

// 静态资源托管
func (c *SgridServerCtx) Static(realPath string, filePath string) {
	target := join(strings.ToLower(c.GetServerName()))(realPath)
	staticPath := public.Join(filePath)
	c.Engine.Static(target, staticPath)
}

type sgridConf struct {
	SgridController    bool
	ServerType         string
	SgridConfPath      string
	SgridGinStaticPath [2]string
	SgridGinWithCors   bool
}
type NewSgrid func(conf *sgridConf)

func WithSgridController() NewSgrid {
	return func(conf *sgridConf) {
		conf.SgridController = true
	}
}

func WithSgridServerType(T string) NewSgrid {
	return func(conf *sgridConf) {
		conf.ServerType = T
	}
}

func WithConfPath(P string) NewSgrid {
	return func(conf *sgridConf) {
		conf.SgridConfPath = P
	}
}

func WithSgridGinStatic(paths [2]string) NewSgrid {
	return func(conf *sgridConf) {
		conf.SgridGinStaticPath = paths
	}
}

func WithCors() NewSgrid {
	return func(conf *sgridConf) {
		conf.SgridGinWithCors = true
	}
}

func NewSgridServerCtx(opt ...NewSgrid) *SgridServerCtx {
	sgridConf := &sgridConf{}
	ctx := &SgridServerCtx{
		Context: context.Background(),
	}
	for _, fn := range opt {
		fn(sgridConf)
	}
	conf, err := public.NewConfig(public.WithTargetPath(sgridConf.SgridConfPath))
	config.GlobalConf = conf
	ctx.SgridConf = conf
	if err != nil {
		fmt.Println("NewConfig Error:", err.Error())
		panic(err.Error())
	}
	if sgridConf.ServerType == public.PROTOCOL_HTTP {
		ctx.Engine = gin.Default()
		if len(sgridConf.SgridGinStaticPath) != 0 {
			GROUP := ctx.Engine.Group(strings.ToLower(ctx.GetServerName()))
			staticRoot := public.Join(sgridConf.SgridGinStaticPath[1])
			GROUP.Static(sgridConf.SgridGinStaticPath[0], staticRoot)
		}
		if sgridConf.SgridGinWithCors {
			GROUP := ctx.Engine.Group(strings.ToLower(ctx.GetServerName()))
			GROUP.Use(withCORSMiddleware())
		}
	}
	if sgridConf.ServerType == public.PROTOCOL_GRPC {
		ctx.Engine = nil
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
	SIMP_CONF_PORT := fmt.Sprintf(":%v", ctx.SgridConf.Server.Port)
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
	c.AbortWithStatusJSON(http.StatusOK, &gin.H{
		"code":    -1,
		"message": err,
		"data":    nil,
	})
}

// Done
var AbortWithSucc = func(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, &gin.H{
		"code":    0,
		"message": "ok",
		"data":    data,
	})
}

// List
var AbortWithSuccList = func(c *gin.Context, data interface{}, total int64) {
	c.AbortWithStatusJSON(http.StatusOK, &gin.H{
		"code":    0,
		"message": "ok",
		"data":    data,
		"total":   total,
	})
}

func join(pre string) func(t string) string {
	return func(target string) string {
		return path.Join(pre, target)
	}
}
