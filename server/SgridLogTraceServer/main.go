package SgridLogTraceServer

import (
	protocol "Sgrid/server/SgridLogTraceServer/proto"
	"Sgrid/src/config"
	"Sgrid/src/pool"
	"Sgrid/src/public"
	"Sgrid/src/storage"
	"Sgrid/src/storage/pojo"
	"context"
	"fmt"
	"net"
	"path"
	"strings"

	"github.com/panjf2000/ants"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var AntsPool *ants.Pool
var SgridLogTraceInstance = &SgridLog{}

type logTraceServer struct {
	protocol.UnimplementedSgridLogTraceServiceServer
}

func (t *logTraceServer) LogTrace(ctx context.Context, in *protocol.LogTraceReq) (*emptypb.Empty, error) {
	fmt.Println("recive request", in)
	err := AntsPool.Submit(func() {
		err := storage.SaveLog(in)
		if err != nil {
			var SystemErrDto = &protocol.LogTraceReq{
				CreateTime:    public.GetCurrTime(),
				LogContent:    err.Error(),
				LogType:       public.LOG_TYPE_SYSTEM_INNER,
				LogHost:       in.LogHost,
				LogGridId:     in.LogGridId,
				LogServerName: in.LogServerName,
				LogBytesLen:   in.LogBytesLen,
			}
			storage.SaveLog(SystemErrDto)
		}
	})
	if err != nil {
		storage.PushErr(&pojo.SystemErr{
			Type: "system/error/SgridLogTraceServer/AntsPool.Submit",
			Info: err.Error(),
		})
	}
	return nil, err
}

type SgridLog struct{}

func (s *SgridLog) Registry(conf *config.SgridConf) {
	AntsPool, _ = ants.NewPool(300, ants.WithPanicHandler(func(i interface{}) {
		info := fmt.Sprintf("failed to listen: %v", i)
		storage.PushErr(&pojo.SystemErr{
			Type: "system/error/SgridLogTraceServer/WithPanicHandler",
			Info: info,
		})
	}))
	pool.InitStorage(conf)
	port := fmt.Sprintf(":%v", conf.Server.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		info := fmt.Sprintf("failed to listen: %v", err)
		storage.PushErr(&pojo.SystemErr{
			Type: "system/error/SgridLogTraceServer/net.Listen",
			Info: info,
		})
	}

	grpcServer := grpc.NewServer()
	protocol.RegisterSgridLogTraceServiceServer(grpcServer, &logTraceServer{})
	fmt.Println("SgridTrace svr started on", port)
	if err := grpcServer.Serve(lis); err != nil {
		info := fmt.Sprintf("failed to serve: %v", err)
		storage.PushErr(&pojo.SystemErr{
			Type: "system/error/SgridLogTraceServer/grpcServer.Serve",
			Info: info,
		})
	}
}

func (s *SgridLog) NameSpace() string {
	return "server.SgridLogTraceServer"
}

func (s *SgridLog) ServerPath() string {
	return strings.ReplaceAll(s.NameSpace(), ".", "/")
}

func (s *SgridLog) JoinPath(args ...string) string {
	p := path.Join(args...)
	return public.Join(s.ServerPath(), p)
}

// func main() {
// 	ctx := h.NewSgridServerCtx(
// 		h.WithSgridServerType(public.PROTOCOL_GRPC),
// 	)
// 	logServant := new(SgridLog)
// 	fmt.Println("SgridLogTraceServer.ctx.SgridConf,", ctx.SgridConf)
// 	logServant.RegistryServer(ctx.SgridConf)
// }
