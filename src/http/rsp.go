package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

var AbortWithError = func(c *gin.Context, err string) {
	c.AbortWithStatusJSON(http.StatusOK, &gin.H{
		"code":    -1,
		"message": err,
		"data":    nil,
	})
}

// Done
var AbortWithSucc = func(c *gin.Context, data interface{}) {
	c.AbortWithStatusJSON(http.StatusOK, &gin.H{
		"code":    0,
		"message": "ok",
		"data":    data,
	})
}

// List
var AbortWithSuccList = func(c *gin.Context, data interface{}, total int64) {
	c.AbortWithStatusJSON(http.StatusOK, &gin.H{
		"code":    0,
		"message": "ok",
		"data":    data,
		"total":   total,
	})
}
