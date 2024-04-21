package service

import (
	"Sgrid/src/configuration"
	handlers "Sgrid/src/http"
	protocol "Sgrid/src/proto"
	"Sgrid/src/public"
	clientgrpc "Sgrid/src/public/client_grpc"
	"Sgrid/src/storage"
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func InitService(ctx *handlers.SgridServerCtx) {
	sc, err := public.NewConfig()
	if err != nil {
		fmt.Println("Error To NewConfig", err)
	}
	configuration.InitStorage(sc)

	gn := storage.QueryPropertiesByKey(SgridPackageServerHosts)
	addresses := []string{}

	for _, v := range gn {
		addresses = append(addresses, v.Value)
	}

	clients := []*clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient]{}
	for _, v := range addresses {
		add := v
		conn, err := grpc.Dial(add, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("无法连接: %v", err)
		}
		// defer conn.Close() // 移动到循环内部
		client := clientgrpc.NewSgridClient[protocol.FileTransferServiceClient](
			protocol.NewFileTransferServiceClient(conn),
			clientgrpc.WithSgridGrpcClientAddress[protocol.FileTransferServiceClient](add),
		)
		clients = append(clients, client)
	}
	ctx.Context = context.WithValue(ctx.Context, public.GRPC_CLIENT_PROXYS{}, clients)
	ctx.RegistryHttpRouter(PackageService)
	ctx.RegistryHttpRouter(Registry)
	ctx.RegistryHttpRouter(DevopsService)
}
