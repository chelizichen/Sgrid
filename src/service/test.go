package service

import (
	handlers "Sgrid/src/http"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func TestRoute(ctx *handlers.SimpHttpServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
	GROUP.POST("/test/sync", func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			fmt.Println("err", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}
		// err = grid.SyncPackage(body, *ctx)
		if err != nil {
			fmt.Println("err", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read request body"})
			return
		}
		body = bytes.Replace(body, []byte("CLUSTER_REQUEST"), []byte("SINGLE_REQUEST"), -1)
		c.JSON(http.StatusOK, gin.H{"message": string(body)})

	})
	ctx.Engine.Use(GROUP.Handlers...)
}
