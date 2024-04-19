package main

import (
	"Sgrid/server/SgridPackageServer"
	h "Sgrid/src/http"
	"Sgrid/src/public"
	service "Sgrid/src/service"
	"fmt"
)

func main() {
	ctx := h.NewSgridServerCtx(
		h.WithSgridServerType(public.PROTOCOL_HTTP),
		h.WithSgridGinStatic([2]string{"/web", "dist"}),
		// h.WithSgridController(),
		h.WithCors(),
	)
	ctx.RegistryHttpRouter(service.InitService)
	ctx.RegistrySubServer(SgridPackageServer.SgridPackageInstance)
	h.NewSgridServer(ctx, func(port string) {
		ctx.Engine.Run(port)
		fmt.Println("Sgrid Gin Http Server started on " + port)
	})
}
