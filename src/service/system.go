package service

import (
	handlers "Sgrid/src/http"
	"Sgrid/src/storage"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/rbac"
	"fmt"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func SystemService(ctx *handlers.SgridServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
	// list get
	GROUP.POST("/system/user/get", getUser)
	GROUP.POST("/system/role/get", getRole)
	GROUP.POST("/system/menu/get", getMenu)
	// save
	GROUP.POST("/system/user/save", saveUser)
	GROUP.POST("/system/role/save", saveRole)
	GROUP.POST("/system/menu/save", saveMenu)
	// relation
	GROUP.POST("/system/setUserToRole", setUserToRole)
	GROUP.POST("/system/setRoleToMenu", setRoleToMenu)
	GROUP.GET("/system/getUserToRoleRelation", getUserToRoleRelation)
	GROUP.GET("/system/getMenuListByRoleId", getMenuListByRoleId)
}

func getUser(c *gin.Context) {
	var req *dto.PageBasicReq
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println("err", err.Error())
		handlers.AbortWithError(c, err.Error())
		return
	}
	u, i := storage.GetUserList(req)
	handlers.AbortWithSuccList(c, u, i)
}

func getRole(c *gin.Context) {
	u := storage.GetRoleList()
	handlers.AbortWithSucc(c, u)
}

func getMenu(c *gin.Context) {
	u := storage.GetMenuList()
	handlers.AbortWithSucc(c, u)
}

func saveUser(c *gin.Context) {
	var req *rbac.User
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	storage.CreateUser(req)
	handlers.AbortWithSucc(c, nil)
}

func saveRole(c *gin.Context) {
	var req *rbac.UserRole
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println("err", err.Error())
		handlers.AbortWithError(c, err.Error())
		return
	}
	storage.CreateRole(req)
	handlers.AbortWithSucc(c, nil)
}

func saveMenu(c *gin.Context) {
	var req *rbac.RoleMenu
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	storage.CreateMenu(req)
	handlers.AbortWithSucc(c, nil)
}

type setUserToRoleDto struct {
	UserId  int   `json:"userId"`
	RoleIds []int `json:"roleIds"`
}

func setUserToRole(c *gin.Context) {
	var req *setUserToRoleDto
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	storage.SetUserToRole(req.UserId, req.RoleIds)
	handlers.AbortWithSucc(c, nil)
}

type setRoleToMenuDto struct {
	RoleId  int   `json:"roleId"`
	MenuIds []int `json:"menuIds"`
}

func setRoleToMenu(c *gin.Context) {
	var req *setRoleToMenuDto
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	storage.SetRoleToMenu(req.RoleId, req.MenuIds)
	handlers.AbortWithSucc(c, nil)
}

func getUserToRoleRelation(c *gin.Context) {
	s, _ := strconv.Atoi(c.Query("id"))
	rutr := storage.GetUserToRoleRelation(s)
	handlers.AbortWithSucc(c, rutr)
}

func getMenuListByRoleId(c *gin.Context) {
	s, _ := strconv.Atoi(c.Query("id"))
	rutr := storage.GetMenuListByRoleId(s)
	handlers.AbortWithSucc(c, rutr)
}
