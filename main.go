package main

import (
	h "Sgrid/src/http"
	service "Sgrid/src/service"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	ctx := h.NewSimpHttpCtx(gin.Default())
	ctx.DefineMain()
	ctx.Use(service.InitService)
	ctx.Use(service.UploadService)
	ctx.Use(service.Registry)
	ctx.Use(service.Static)
	ctx.Use(service.TestRoute)
	h.NewSimpHttpServer(ctx, func(port string) {
		ctx.Engine.Run(port)
		fmt.Println("Server started on " + port)
	})
}
