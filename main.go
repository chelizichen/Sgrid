package main

import (
	h "Sgrid/src/http"
	service "Sgrid/src/service"
	"fmt"
)

func main() {
	ctx := h.NewSgridServerCtx(
		h.WithSgridServerType(h.GIN_HTTP_SERVER),
		h.WithSgridGinStatic("/static"),
		// h.WithSgridController(),
		h.WithCors(),
	)
	ctx.Use(service.InitService)
	ctx.Use(service.UploadService)
	ctx.Use(service.Registry)
	h.NewSgridServer(ctx, func(port string) {
		ctx.Engine.Run(port)
		fmt.Println("Server started on " + port)
	})
}
