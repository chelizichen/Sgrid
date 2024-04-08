package storage

import (
	c "Sgrid/src/configuration"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/pojo"
	"Sgrid/src/storage/vo"
	"Sgrid/src/utils"
)

// 查询标签组，创建时用
func QueryTags() []pojo.ServantGroup {
	var dataList []pojo.ServantGroup
	c.GORM.Model(dataList).Find(&dataList)
	return dataList
}

// 查询节点组
func QueryNodes() []pojo.Node {
	var dataList []pojo.Node
	c.GORM.Model(dataList).Find(&dataList)
	return dataList
}

func QueryServants() []vo.ServantVo {
	var dataList []vo.ServantVo
	c.GORM.Model(dataList).Find(&dataList)
	return dataList
}

func QueryGrid(req *dto.PageBasicReq) (resp []vo.Grid) {
	where := "1 = 1"
	args := make([]interface{}, 10)
	if req.Keyword != "" {
		where += "and ( gn.ip like '%?%' or gs.server_name like '%?%')"
		args = append(args, req.Keyword, req.Keyword)
	}
	c.GORM.Table("grid_grid gg").
		Select(`
	gg.*,
	gs.id AS gs_id,
	gs.language AS gs_language,
	gs.servant_group_id AS gs_servant_group_id,
	gs.server_name as gs_server_name,
	gs.create_time as gs_create_time,
	gn.id as gn_id,
	gn.ip as gn_ip,
	gn.main as gn_main,
	gn.plat_form as gn_plat_form,
	gn.status as gn_status,
	gn.create_time as gn_create_time
	`).
		Joins("LEFT JOIN grid_servant gs ON gs.id = gg.servant_id").
		Joins("LEFT JOIN grid_node gn ON gn.id = gg.node_id").
		Where(where, utils.Removenullvalue(args)...).
		Find(&resp)

	return
}

func ServantPackage(ServantId int) []vo.Grid {
	var dataList []vo.Grid
	c.GORM.Model(dataList).Find(&dataList)
	return dataList
}
