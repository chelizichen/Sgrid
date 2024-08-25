// 系统RBAC角色权限模块
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
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.GetServerName()))
	// list get
	GROUP.POST("/system/user/get", getUser)
	GROUP.POST("/system/role/get", getRole)
	GROUP.POST("/system/menu/get", getMenu)
	GROUP.POST("/system/group/get", getUserGroup)
	// save
	GROUP.POST("/system/user/save", saveUser)
	GROUP.POST("/system/role/save", saveRole)
	GROUP.POST("/system/menu/save", saveMenu)
	GROUP.POST("/system/group/save", saveUserGroup)

	// del
	GROUP.GET("/system/menu/del", delMenu)
	GROUP.GET("/system/role/del", delRole)
	GROUP.GET("/system/group/del", delUserGroup)

	// relation
	GROUP.POST("/system/setUserToRole", setUserToRole)
	GROUP.POST("/system/setRoleToMenu", setRoleToMenu)
	GROUP.GET("/system/getUserToRoleRelation", getUserToRoleRelation)
	GROUP.GET("/system/getMenuListByRoleId", getMenuListByRoleId)

	GROUP.POST("/system/setUserToUserGroup", setUserToGroup)                     // 将用户分配到团队
	GROUP.POST("/system/setUserGroupToServantGroup", setUserGroupToServantGroup) // 为服务组分配团队

	// SPEC
	GROUP.POST("/system/spec/getServantGroupsByUserGroupId", getServantGroupsByUserGroupId)
	GROUP.POST("/system/spec/getUsersByUserGroup", getUsersByUserGroup)

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

func getUserGroup(c *gin.Context) {
	var req *dto.PageBasicReq
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println("err", err.Error())
		handlers.AbortWithError(c, err.Error())
		return
	}
	u, t, err := storage.GetUserGroupList(req)
	if err != nil {
		fmt.Println("err", err.Error())
		handlers.AbortWithError(c, err.Error())
		return
	}
	handlers.AbortWithSuccList(c, u, *t)
}

func getUsersByUserGroup(c *gin.Context) {
	var req *dto.PageBasicReq
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println("err", err.Error())
		handlers.AbortWithError(c, err.Error())
		return
	}
	u, t, err := storage.GetUsersByUserGroup(req)
	if err != nil {
		fmt.Println("err", err.Error())
		handlers.AbortWithError(c, err.Error())
		return
	}
	handlers.AbortWithSuccList(c, u, *t)
}

// SPEC

func getServantGroupsByUserGroupId(c *gin.Context) {
	var req *dto.PageBasicReq
	err := c.BindJSON(&req)
	if err != nil {
		fmt.Println("err", err.Error())
		handlers.AbortWithError(c, err.Error())
		return
	}
	u, t, err := storage.GetServantGroupsByUserGroupId(req)
	if err != nil {
		fmt.Println("err", err.Error())
		handlers.AbortWithError(c, err.Error())
		return
	}
	handlers.AbortWithSuccList(c, u, *t)
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

func saveUserGroup(c *gin.Context) {
	var req *rbac.UserGroup
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	storage.CreateGroup(req)
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

// table struct { 	user_id , user_group_id }
type setUserToGroupDto struct {
	UserIds []int `json:"userIds"`
	GroupId int   `json:"groupId"`
}

func setUserToGroup(c *gin.Context) {
	var req *setUserToGroupDto
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	storage.SetUserToGroup(req.GroupId, req.UserIds)
	handlers.AbortWithSucc(c, nil)
}

type setUserGroupToServantGroupDto struct {
	GroupId         int   `json:"groupId"`
	ServantGroupIds []int `json:"servantGroupIds"`
}

func setUserGroupToServantGroup(c *gin.Context) {
	var req *setUserGroupToServantGroupDto
	err := c.BindJSON(&req)
	if err != nil {
		handlers.AbortWithError(c, err.Error())
		return
	}
	storage.SetUserGroupToServantGroup(req.GroupId, req.ServantGroupIds)
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

func delMenu(c *gin.Context) {
	s, _ := strconv.Atoi(c.Query("id"))
	storage.DeleteMenu(s)
	handlers.AbortWithSucc(c, nil)
}

func delRole(c *gin.Context) {
	s, _ := strconv.Atoi(c.Query("id"))
	storage.DeleteRole(s)
	handlers.AbortWithSucc(c, nil)
}

func delUserGroup(c *gin.Context) {
	s, _ := strconv.Atoi(c.Query("id"))
	storage.DeleteUserGroup(s)
	handlers.AbortWithSucc(c, nil)
}
