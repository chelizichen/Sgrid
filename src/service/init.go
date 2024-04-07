package service

import (
	"Sgrid/src/configuration"
	handlers "Sgrid/src/http"
)

func InitService(ctx *handlers.SimpHttpServerCtx) {
	configuration.InitStorage(ctx)
}
