package service

import (
	handlers "Sgrid/src/http"
	"Sgrid/src/public/jwt"
	"Sgrid/src/storage"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/pojo"
	"Sgrid/src/storage/vo"
	utils "Sgrid/src/utils"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	TOKEN           = "e609d00404645feed1c1733835b8c127"
	CLUSTER_REQUEST = "CLUSTER_REQUEST"
	SINGLE_REQUEST  = "SINGLE_REQUEST"
)

func Registry(ctx *handlers.SgridServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
	GROUP.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		v := storage.QueryUser(&pojo.User{
			UserName: username,
			Password: password,
		})
		if len(v.Password) != 0 && len(v.UserName) != 0 {

			s, err := jwt.GenToken(*v)
			if err != nil {
				fmt.Println("gen token error")
				c.JSON(http.StatusOK, handlers.Resp(-1, "Error", v))
			}
			rsp := vo.VoUser{
				UserName:   v.UserName,
				Password:   v.Password,
				CreateTime: v.CreateTime,
				Id:         v.Id,
				Token:      s,
			}
			c.JSON(http.StatusOK, handlers.Resp(0, "ok", rsp))
		} else {
			c.JSON(http.StatusOK, handlers.Resp(-1, "Error", v))
		}
	})

	GROUP.GET("/main/queryGrid", func(c *gin.Context) {
		pbr := utils.NewPageBaiscReq(c)
		gv := storage.QueryGrid(pbr)
		c.JSON(200, handlers.Resp(0, "ok", gv))
	})

	GROUP.GET("/main/queryServantGroup", func(c *gin.Context) {
		gv := storage.QueryServantGroup(&dto.PageBasicReq{})
		vgbs := storage.ConvertToVoGroupByServant(gv)
		c.JSON(200, handlers.Resp(0, "ok", vgbs))
	})

	GROUP.GET("/main/queryNodes", func(c *gin.Context) {
		nodes := storage.QueryNodes()
		c.JSON(200, handlers.Resp(0, "ok", nodes))
	})

	ctx.Engine.Use(GROUP.Handlers...)

}
