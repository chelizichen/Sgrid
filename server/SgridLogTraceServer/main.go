package main

import (
	protocol "Sgrid/server/SgridLogTraceServer/proto"
	"Sgrid/src/config"
	"Sgrid/src/configuration"
	h "Sgrid/src/http"
	"Sgrid/src/public"
	"Sgrid/src/storage"
	"Sgrid/src/storage/pojo"
	"context"
	"fmt"
	"net"

	"github.com/panjf2000/ants"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var AntsPool *ants.Pool

type logTraceServer struct {
	protocol.UnimplementedSgridLogTraceServiceServer
}

func (t *logTraceServer) LogTrace(ctx context.Context, in *protocol.LogTraceReq) (*emptypb.Empty, error) {
	AntsPool.Submit(func() {
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
	return nil, nil
}

type SgridLog struct{}

func (s *SgridLog) RegistryServer(conf *config.SgridConf) {
	AntsPool, _ = ants.NewPool(10)
	configuration.InitStorage(conf)
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
	fmt.Println("Sgrid svr started on", port)
	if err := grpcServer.Serve(lis); err != nil {
		info := fmt.Sprintf("failed to serve: %v", err)
		storage.PushErr(&pojo.SystemErr{
			Type: "system/error/SgridLogTraceServer/grpcServer.Serve",
			Info: info,
		})
	}
}

func main() {
	ctx := h.NewSgridServerCtx(
		h.WithSgridServerType(public.PROTOCOL_GRPC),
	)
	logServant := new(SgridLog)
	logServant.RegistryServer(ctx.SgridConf)
}
