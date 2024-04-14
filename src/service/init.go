package service

import (
	"Sgrid/src/configuration"
	handlers "Sgrid/src/http"
	"Sgrid/src/public"
	"fmt"
)

func InitService(ctx *handlers.SgridServerCtx) {
	sc, err := public.NewConfig()
	if err != nil {
		fmt.Println("Error To NewConfig", err)
	}
	configuration.InitStorage(sc)
}
