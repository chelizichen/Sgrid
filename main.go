package main

import (
	"Sgrid/server/SgridPackageServer"
	h "Sgrid/src/http"
	"Sgrid/src/public"
	service "Sgrid/src/service"
	"fmt"
	"os"
	"path/filepath"
)

// func init() {
// 	fmt.Println(`
// ********* info **********
//          开始验证
// ********* info **********
// 		`)
// 	b, err := saas.UsesaasPerm.CheckAuth()
// 	if !b || err != nil {
// 		fmt.Println(`
// ********* error **********
//          验证失败
// ********* error **********
// 			`)
// 		return
// 	}
// }

func init() {
	cwd, _ := os.Getwd()
	public.CheckDirectoryOrCreate(filepath.Join(cwd, "server"))
	public.CheckDirectoryOrCreate(filepath.Join(cwd, "server", "SgridPackageServer"))
}

func main() {
	ctx := h.NewSgridServerCtx(
		h.WithSgridServerType(public.PROTOCOL_HTTP),
		h.WithSgridGinStatic([2]string{"/web", "dist"}),
		h.WithCors(),
	)
	ctx.RegistryHttpRouter(service.InitService)
	ctx.RegistrySubServer(SgridPackageServer.SgridPackageInstance)
	h.NewSgridServer(ctx, func(port string) {
		ctx.Engine.Run(port)
		fmt.Println("Sgrid Gin Http Server started on " + port)
	})
}
