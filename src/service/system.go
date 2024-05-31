package service

import (
	handlers "Sgrid/src/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SystemService(ctx *handlers.SgridServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
	// devops component
	GROUP.GET("/system/user/get", getUser)
	GROUP.POST("/system/user/set", getUser)
	GROUP.GET("/system/user/one", getUser)
	GROUP.GET("/system/user/delete", getUser)
}

func getUser(ctx *gin.Context) {

}
