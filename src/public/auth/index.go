package sgirdAuth

import (
	c "Sgrid/src/configuration"
)

// 角色
type Role struct {
	Id   int
	Name string
}

// 接口权限
// 根据RoleId 查询是否有Path
// 如果则可以返回
type ApiToRole struct {
	Id     int
	RoleId int
	ApiId  int
}

type Apis struct {
	Id   int
	Path string
}

// 查询权限列表
func QueryApiList() []Apis {
	var dataList []Apis
	c.GORM.Model(dataList).Find(&dataList)
	return dataList
}

// 查询用户权限
func QueryPassApi(roleId int) []Apis {
	var dataList []Apis
	c.GORM.Exec(`
		select ga.* from 
		grid_api_to_role gatr
		left join grid_apis ga 
		on gatr.role_id = ?
	`, roleId).Find(&dataList)
	return dataList
}

// 重设权限
func SetRoleApis(roleId int, apis []Apis) {
	c.GORM.Exec("delete * from sgrid_api_to_role where role_id = ?", roleId)
	c.GORM.Model(&Apis{}).Save(apis)
}
