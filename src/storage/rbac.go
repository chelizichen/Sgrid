package storage

import (
	c "Sgrid/src/configuration"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/rbac"
	"Sgrid/src/utils"
	"fmt"
)

func GetUserList(req *dto.PageBasicReq) ([]rbac.User, int64) {
	var respList []rbac.User
	var count int64
	args := make([]interface{}, 10)
	where := "1 = 1"
	if req.Keyword != "" {
		where += " and user_name like ?"
		args = append(args, "%"+req.Keyword+"%")
	}
	c.GORM.
		Model(&rbac.User{}).
		Offset(req.Offset).
		Limit(req.Size).
		Count(&count).
		Where(
			where,
			utils.Removenullvalue(args)...,
		).
		Find(&respList)
	return respList, count
}

func GetMenuList() []rbac.RoleMenu {
	var respList []rbac.RoleMenu
	c.GORM.
		Model(&rbac.RoleMenu{}).
		Find(&respList)
	return respList
}

func GetRoleList() []rbac.UserRole {
	var respList []rbac.UserRole
	c.GORM.
		Model(&rbac.UserRole{}).
		Find(&respList)

	return respList
}

// 通过角色ID 拿到菜单列表
func GetMenuListByRoleId(roleId int) []rbac.RoleToMenu {
	var respList []rbac.RoleToMenu
	c.GORM.
		Model(&rbac.RoleToMenu{}).
		Where("role_id = ?", roleId).
		Find(&respList)
	return respList
}

func DeleteMenu(id int) {
	c.GORM.Model(&rbac.RoleMenu{}).Delete(&rbac.RoleMenu{
		Id: id,
	})
	c.GORM.Model(&rbac.RoleToMenu{}).Delete(&rbac.RoleToMenu{
		MenuId: id,
	})
}

func DeleteRole(id int) {
	c.GORM.Model(&rbac.UserRole{}).Delete(&rbac.UserRole{
		Id: id,
	})
	c.GORM.Model(&rbac.UserToRole{}).Delete(&rbac.UserToRole{
		RoleId: id,
	})
}

func SetUserToRole(userId int, roleIds []int) {
	c.GORM.Delete(&rbac.UserToRole{}, "user_id = ?", userId)
	var userToRoles []*rbac.UserToRole
	for _, v := range roleIds {
		userToRoles = append(userToRoles, &rbac.UserToRole{
			UserId: userId,
			RoleId: v,
		})
	}
	c.GORM.Create(userToRoles)
}

func SetRoleToMenu(roleId int, menuIds []int) {
	c.GORM.Delete(&rbac.RoleToMenu{}, "role_id = ?", roleId)
	var userToRoles []*rbac.RoleToMenu
	for _, v := range menuIds {
		userToRoles = append(userToRoles, &rbac.RoleToMenu{
			RoleId: roleId,
			MenuId: v,
		})
	}
	c.GORM.Create(userToRoles)
}

func CreateRole(role *rbac.UserRole) {
	if role.Id == 0 {
		c.GORM.Create(role)
	} else {
		c.GORM.Model(&rbac.UserRole{}).
			Where("id = ?", role.Id).
			Updates(&rbac.UserRole{
				Name:        role.Name,
				Description: role.Description,
			})
	}
}

func CreateUser(user *rbac.User) {
	fmt.Println("user", user)
	if user.Id == 0 {
		user.Password = "e10adc3949ba59abbe56e057f20f883e" // 123456

		c.GORM.Create(user)
	} else {
		c.GORM.Model(&rbac.User{}).
			Where("id = ?", user.Id).
			Updates(&rbac.User{
				UserName:  user.UserName,
				TurthName: user.TurthName,
			})
	}
}

func CreateMenu(menu *rbac.RoleMenu) {
	if menu.Id == 0 {
		c.GORM.Create(menu)
	} else {
		c.GORM.Model(&rbac.RoleMenu{}).
			Where("id = ?", menu.Id).
			Updates(&rbac.RoleMenu{
				Title:     menu.Title,
				Path:      menu.Path,
				Name:      menu.Name,
				Component: menu.Component,
			})
	}
}

// relation
type RelationUserToRole struct {
	ID   uint   `gorm:"id" json:"id,omitempty"`
	Name string `gorm:"name" json:"name,omitempty"`
}

func GetUserToRoleRelation(id int) []RelationUserToRole {
	var findList []RelationUserToRole
	c.GORM.Debug().Raw(` 
	select gsr.id,gsr.name from grid_user_to_role gstr
	left join grid_user_role gsr on gstr.role_id = gsr.id
	left join grid_user gu on gu.id = gstr.user_id
	where gstr.user_id = ?
	`, id).Scan(&findList)
	return findList
}

func GetUserMenusByUserId(id int) []rbac.RoleMenu {
	var findList []rbac.RoleMenu
	c.GORM.Raw(` 
	select
	grm.*
from
	grid_role_to_menu grtm
left join grid_role_menu grm on
	grtm.menu_id = grm.id
where
	grtm.role_id  in (
	select
		gutr.role_id
	from
		grid_user_to_role gutr
	left join grid_user gu on
		gutr.user_id = gu.id
	where
		gu.id = ?
	)
	`, id).Scan(&findList)
	return findList
}
