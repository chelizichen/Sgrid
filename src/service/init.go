package service

import (
	protocol "Sgrid/server/SgridPackageServer/proto"
	handlers "Sgrid/src/http"
	"Sgrid/src/pool"
	"Sgrid/src/public"
	"Sgrid/src/rpc"
	"Sgrid/src/storage"
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type PackageServantProxy struct{}

func InitService(ctx *handlers.SgridServerCtx) {
	sc, err := public.NewConfig()
	if err != nil {
		fmt.Println("Error To NewConfig", err)
	}
	pool.InitStorage(sc)

	gn := storage.QueryPropertiesByKey(SgridPackageServerHosts)
	addresses := []string{}
	for _, v := range gn {
		addresses = append(addresses, v.Value)
	}

	packageServant, err := rpc.NewSgridGrpcClient[protocol.FileTransferServiceClient](
		addresses,
		rpc.WithDiaoptions[protocol.FileTransferServiceClient](
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
		rpc.WithClientService[protocol.FileTransferServiceClient](protocol.NewFileTransferServiceClient),
		rpc.WithRequestPrefix[protocol.FileTransferServiceClient]("/SgridProtocol.FileTransferService/"),
	)
	if err != nil {
		fmt.Println("Error To NewSgridGrpcClient ", err.Error())
	}
	ctx.Context = context.WithValue(ctx.Context, PackageServantProxy{}, packageServant)
	ctx.RegistryHttpRouter(PackageService)
	ctx.RegistryHttpRouter(Registry)
	ctx.RegistryHttpRouter(DevopsService)
	ctx.RegistryHttpRouter(SystemService)
	ctx.RegistryHttpRouter(SystemStatisticsRegisty)
	ctx.RegistryHttpRouter(AssetsService)
}
