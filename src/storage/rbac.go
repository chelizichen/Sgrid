package storage

import (
	"Sgrid/src/pool"
	"Sgrid/src/public"
	"Sgrid/src/public/replace"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/rbac"
	"fmt"
)

func GetUserGroupList(req *dto.PageBasicReq) ([]rbac.UserGroupVo, *int64, error) {
	var respList []rbac.UserGroupVo
	var selects = `
		gug.*,
		count(gutug.user_id) as total
	`
	where := "where 1 = 1 "
	args := make([]interface{}, 10)
	if req.Keyword != "" {
		where += " and gug.name like ? "
		args = append(args, public.BuildKeyword(req.Keyword))
	}
	var sql, err = replace.BuildReplaceChain(`
SELECT
	${SELECTS}
from
	grid_user_group gug
left join grid_user_to_user_group gutug on
	gutug.user_group_id = gug.id
	${WHERE}
group by
	gug.id
	${PAGINATION}
`)
	querySql := sql.ReplaceSelects(selects).ReplaceWhere(where).ReplacePagination(req.Offset, req.Size)
	err = pool.GORM.Debug().Raw(querySql.Get(), public.Removenullvalue(args)...).Scan(&respList).Error
	countSql := sql.Reset().ReplaceAsCount().ReplaceWhere(where).ReplaceWithNoPagination()
	pool.GORM.Debug().Raw(countSql.Get(), public.Removenullvalue(args)...).Scan(countSql.GetCountVo())
	return respList, countSql.GetCountVo(), err
}

func GetUsersByUserGroup(req *dto.PageBasicReq) ([]rbac.UserToUserGroupVo, *int64, error) {
	var respList []rbac.UserToUserGroupVo
	where := " where 1 = 1 "
	args := make([]interface{}, 10)
	var selects = `
		gu.user_name,gug.name,gutug.user_id,gutug.user_group_id
	`
	if req.Id != 0 {
		where += " and gutug.user_group_id  = ? "
		args = append(args, req.Id)
	}
	var sql, _ = replace.BuildReplaceChain(`
SELECT
	${SELECTS}
from
		grid_user_to_user_group gutug
left join grid_user gu on
		gu.id = gutug.user_id
left JOIN grid_user_group gug on
		gug.id = gutug.user_group_id
	${WHERE}
	${PAGINATION}
`)
	querySql := sql.ReplaceSelects(selects).ReplaceWhere(where).ReplacePagination(req.Offset, req.Size)
	err := pool.GORM.Debug().Raw(querySql.Get(), public.Removenullvalue(args)...).Scan(&respList).Error
	countSql := sql.Reset().ReplaceAsCount().ReplaceWhere(where).ReplaceWithNoPagination()
	pool.GORM.Debug().Raw(countSql.Get(), public.Removenullvalue(args)...).Scan(countSql.GetCountVo())
	return respList, countSql.GetCountVo(), err
}

func GetServantGroupsByUserGroupId(req *dto.PageBasicReq) ([]rbac.UserGroupToServantGroupVo, *int64, error) {
	var respList []rbac.UserGroupToServantGroupVo
	where := " where 1 = 1 "
	args := make([]interface{}, 10)
	var selects = `
	gugtsg.user_group_id  		as user_group_id,
	gugtsg.servant_group_id 	as servant_group_id,
	gug.name 					as user_group_name,
	gsg.tag_name  				as servant_group_name
	`
	if req.Id != 0 {
		where += " and gugtsg.user_group_id  = ? "
		args = append(args, req.Id)
	}
	var sql, _ = replace.BuildReplaceChain(`
SELECT
	${SELECTS}
	from grid_user_group_to_servant_group gugtsg
left join grid_servant_group gsg on gsg.id  = gugtsg.servant_group_id
left join grid_user_group gug on gug.id  = gugtsg.user_group_id
	${WHERE}
	${PAGINATION}
`)
	querySql := sql.ReplaceSelects(selects).ReplaceWhere(where).ReplacePagination(req.Offset, req.Size)
	err := pool.GORM.Debug().Raw(querySql.Get(), public.Removenullvalue(args)...).Scan(&respList).Error
	countSql := sql.Reset().ReplaceAsCount().ReplaceWhere(where).ReplaceWithNoPagination()
	pool.GORM.Debug().Raw(countSql.Get(), public.Removenullvalue(args)...).Scan(countSql.GetCountVo())
	return respList, countSql.GetCountVo(), err
}

func GetUserList(req *dto.PageBasicReq) ([]rbac.User, int64) {
	var respList []rbac.User
	var count int64
	args := make([]interface{}, 10)
	where := "1 = 1"
	if req.Keyword != "" {
		where += " and user_name like ?"
		args = append(args, "%"+req.Keyword+"%")
	}
	pool.GORM.
		Model(&rbac.User{}).
		Offset(req.Offset).
		Limit(req.Size).
		Count(&count).
		Where(
			where,
			public.Removenullvalue(args)...,
		).
		Find(&respList)
	return respList, count
}

func GetMenuList() []rbac.RoleMenu {
	var respList []rbac.RoleMenu
	pool.GORM.
		Model(&rbac.RoleMenu{}).
		Find(&respList)
	return respList
}

func GetRoleList() []rbac.UserRole {
	var respList []rbac.UserRole
	pool.GORM.
		Model(&rbac.UserRole{}).
		Find(&respList)

	return respList
}

// 通过角色ID 拿到菜单列表
func GetMenuListByRoleId(roleId int) []rbac.RoleToMenu {
	var respList []rbac.RoleToMenu
	pool.GORM.
		Model(&rbac.RoleToMenu{}).
		Where("role_id = ?", roleId).
		Find(&respList)
	return respList
}

func DeleteMenu(id int) {
	pool.GORM.Model(&rbac.RoleMenu{}).Delete(&rbac.RoleMenu{
		Id: id,
	})
	pool.GORM.Model(&rbac.RoleToMenu{}).Delete(&rbac.RoleToMenu{
		MenuId: id,
	})
}

func DeleteRole(id int) {
	pool.GORM.Model(&rbac.UserRole{}).Delete(&rbac.UserRole{
		Id: id,
	})
	pool.GORM.Model(&rbac.UserToRole{}).Delete(&rbac.UserToRole{
		RoleId: id,
	})
}

func DeleteUserGroup(id int) {
	pool.GORM.Model(&rbac.UserGroup{}).Delete(&rbac.UserGroup{
		Id: id,
	})
	pool.GORM.Model(&rbac.UserToUserGroup{}).Delete(&rbac.UserToUserGroup{
		UserGroupId: id,
	})
}

func SetUserToRole(userId int, roleIds []int) {
	pool.GORM.Delete(&rbac.UserToRole{}, "user_id = ?", userId)
	var userToRoles []*rbac.UserToRole
	for _, v := range roleIds {
		userToRoles = append(userToRoles, &rbac.UserToRole{
			UserId: userId,
			RoleId: v,
		})
	}
	pool.GORM.Create(userToRoles)
}

func SetRoleToMenu(roleId int, menuIds []int) {
	pool.GORM.Delete(&rbac.RoleToMenu{}, "role_id = ?", roleId)
	var userToRoles []*rbac.RoleToMenu
	for _, v := range menuIds {
		userToRoles = append(userToRoles, &rbac.RoleToMenu{
			RoleId: roleId,
			MenuId: v,
		})
	}
	pool.GORM.Create(userToRoles)
}

func SetUserToGroup(group_id int, user_ids []int) {
	fmt.Println("group_id", group_id)
	fmt.Println("user_ids", user_ids)
	pool.GORM.Delete(&rbac.UserToUserGroup{}, "user_group_id = ?", group_id)
	var userToGroupList []*rbac.UserToUserGroup
	for _, v := range user_ids {
		userToGroupList = append(userToGroupList, &rbac.UserToUserGroup{
			UserId:      v,
			UserGroupId: group_id,
		})
	}
	pool.GORM.Create(userToGroupList)
}

func SetUserGroupToServantGroup(user_group_id int, servant_group_ids []int) {
	fmt.Println("user_group_id", user_group_id)
	fmt.Println("servant_group_ids", servant_group_ids)
	pool.GORM.Delete(&rbac.UserGroupToServantGroup{}, "user_group_id = ?", user_group_id)
	var userToGroupList []*rbac.UserGroupToServantGroup
	for _, v := range servant_group_ids {
		userToGroupList = append(userToGroupList, &rbac.UserGroupToServantGroup{
			UserGroupId:    user_group_id,
			ServantGroupId: v,
		})
	}
	pool.GORM.Create(userToGroupList)
}

func CreateRole(role *rbac.UserRole) {
	if role.Id == 0 {
		pool.GORM.Create(role)
	} else {
		pool.GORM.Model(&rbac.UserRole{}).
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

		pool.GORM.Create(user)
	} else {
		pool.GORM.Model(&rbac.User{}).
			Where("id = ?", user.Id).
			Updates(&rbac.User{
				UserName:  user.UserName,
				TurthName: user.TurthName,
			})
	}
}
func CreateGroup(g *rbac.UserGroup) {
	if g.Id == 0 {
		pool.GORM.Create(g)
	} else {
		pool.GORM.Model(&rbac.UserGroup{}).
			Where("id = ?", g.Id).
			Updates(&rbac.UserGroup{
				Name:         g.Name,
				Description:  g.Description,
				CreateUserId: g.CreateUserId,
			})
	}
}

func CreateMenu(menu *rbac.RoleMenu) {
	if menu.Id == 0 {
		pool.GORM.Create(menu)
	} else {
		pool.GORM.Model(&rbac.RoleMenu{}).
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
	pool.GORM.Debug().Raw(`
	select gsr.id,gsr.name from grid_user_to_role gstr
	left join grid_user_role gsr on gstr.role_id = gsr.id
	left join grid_user gu on gu.id = gstr.user_id
	where gstr.user_id = ?
	`, id).Scan(&findList)
	return findList
}

func GetUserMenusByUserId(id int) []rbac.RoleMenu {
	var findList []rbac.RoleMenu
	pool.GORM.Raw(`
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

// 根据用户查找所属服务组 ids
func GetUserGroupsFromUserId(userId int) []int {
	var ids []int
	var sql = `
	select servant_group_id from grid_user_group_to_servant_group
	where user_group_id in (
		select user_group_id from grid_user_to_user_group
		where user_id = ?
	)
	`
	pool.GORM.Raw(sql, userId).Scan(&ids)
	return ids
}
