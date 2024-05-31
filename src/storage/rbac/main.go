package rbac

import (
	"time"
)

// User Role Menu
type User struct {
	Id            int        `json:"id,omitempty"`        // id
	UserName      string     `json:"userName,omitempty"`  // 用户名
	Password      string     `json:"password,omitempty"`  // 密码
	TurthName     string     `json:"turthName,omitempty"` // 真实姓名
	CreateTime    *time.Time `gorm:"autoCreateTime" json:"createTime,omitempty"`
	LastLoginTime *time.Time `json:"lastLoginTime,omitempty"` // 上次登陆时间
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
