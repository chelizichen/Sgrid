package main

import (
	"Sgrid/server/SgridLogTraceServer"
	"Sgrid/server/SgridPackageServer"
	h "Sgrid/src/http"
	"Sgrid/src/public"
	"Sgrid/src/saas"
	service "Sgrid/src/service"
	"fmt"
)

func init() {
	fmt.Println(`
********* info **********
		开始验证
********* info **********
		`)
	b, err := saas.UsesaasPerm.CheckAuth()
	if !b || err != nil {
		fmt.Println(`
********* error **********
		验证失败
********* error **********
			`)
		return
	}
}

func main() {
	ctx := h.NewSgridServerCtx(
		h.WithSgridServerType(public.PROTOCOL_HTTP),
		h.WithSgridGinStatic([2]string{"/web", "dist"}),
		h.WithCors(),
	)
	ctx.RegistryHttpRouter(service.InitService)
	ctx.RegistrySubServer(SgridLogTraceServer.SgridLogTraceInstance)
	ctx.RegistrySubServer(SgridPackageServer.SgridPackageInstance)
	h.NewSgridServer(ctx, func(port string) {
		ctx.Engine.Run(port)
		fmt.Println("Sgrid Gin Http Server started on " + port)
	})
}
