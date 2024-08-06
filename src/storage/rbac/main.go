package rbac

import (
	"time"
)

// User Role Menu
type User struct {
	Id            int        `json:"id,omitempty"`                               // id
	UserName      string     `json:"userName,omitempty"`                         // 用户名
	Password      string     `json:"password,omitempty"`                         // 密码
	TurthName     string     `json:"turthName,omitempty"`                        // 真实姓名
	CreateTime    *time.Time `gorm:"autoCreateTime" json:"createTime,omitempty"` // 创建时间
	LastLoginTime *time.Time `json:"lastLoginTime,omitempty"`                    // 上次登陆时间
	UserGroupId   int        `json:"userGroupId"`                                // 用户组ID
}

// 用户组
type UserGroup struct {
	Id           int        `json:"id,omitempty"`            // id
	Name         string     `json:"name,omitempty"`          // 组名
	Description  string     `json:"description,omitempty"`   // 描述
	CreateUserId int        `json:"createUser_id,omitempty"` // 创建人
	Status       int        `json:"status,omitempty"`
	CreateTime   *time.Time `gorm:"autoCreateTime" json:"createTime,omitempty"` // 创建时间
}

// 用户组
type UserGroupVo struct {
	Id           int        `json:"id"`           // id
	Name         string     `json:"name"`         // 组名
	Description  string     `json:"description"`  // 描述
	CreateUserId int        `json:"createUserId"` // 创建人
	CreateTime   *time.Time `json:"createTime"`   // 创建时间
	Total        string     `json:"total"`        // 创建时间
	Status       int        `json:"status"`
}

type UserToUserGroupVo struct {
	UserId      int    `json:"userId"`      // 用户ID
	UserGroupId int    `json:"userGroupId"` // 用户群组ID
	UserName    string `json:"userName"`    // 用户名
	Name        string `json:"name"`        // 小组名
}

// 映射
type UserToUserGroup struct {
	UserId      int // 用户ID
	UserGroupId int // 用户群组ID
}

type UserToRole struct {
	UserId int
	RoleId int
}

type UserRole struct {
	Id          int        `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`                             // 角色名
	Description string     `json:"description,omitempty"`                      // 角色名
	CreateTime  *time.Time `gorm:"autoCreateTime" json:"createTime,omitempty"` // 创建时间
}

type RoleToMenu struct {
	RoleId int `json:"roleId,omitempty"`
	MenuId int `json:"menuId,omitempty"`
}

type RoleMenu struct {
	Id        int    `json:"id,omitempty"`        // id
	Title     string `json:"title,omitempty"`     // 标题
	Path      string `json:"path,omitempty"`      // URL
	Name      string `json:"name,omitempty"`      // 名称
	Component string `json:"component,omitempty"` // 组建路径
	ParentId  int    `json:"parentId,omitempty"`  // 父级id
}

type VersionUpdateLine struct {
	Id         int        `json:"id,omitempty"`
	Title      string     `json:"title,omitempty"`
	Content    string     `gorm:"type:text" json:"content,omitempty" `
	UserId     int        `json:"userId,omitempty"`
	CreateTime *time.Time `gorm:"autoCreateTime" json:"createTime,omitempty"`
}
