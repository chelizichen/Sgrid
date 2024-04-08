package service

import (
	"Sgrid/src/grid"
	handlers "Sgrid/src/http"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func TestRoute(ctx *handlers.SimpHttpServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
	GROUP.GET("/test/sync", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("err", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}
		err = grid.SyncPackage(body, *ctx)
		if err != nil {
			fmt.Println("err", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "ok"})

	})
	ctx.Engine.Use(GROUP.Handlers...)
}
