package main

import (
	h "Sgrid/src/http"
	"Sgrid/src/public"
	service "Sgrid/src/service"
	"fmt"
)

func main() {
	ctx := h.NewSgridServerCtx(
		h.WithSgridServerType(public.PROTOCOL_HTTP),
		h.WithSgridGinStatic("/static"),
		h.WithSgridController(),
		h.WithCors(),
	)
	ctx.Use(service.InitService)
	ctx.Use(service.PackageService)
	ctx.Use(service.Registry)
	h.NewSgridServer(ctx, func(port string) {
		ctx.Engine.Run(port)
		fmt.Println("Server started on " + port)
	})
}
