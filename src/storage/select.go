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

func QueryServants() []vo.VoServant {
	var dataList []vo.VoServant
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

func QueryServantGroup(req *dto.PageBasicReq) (resp []vo.VoServantGroup) {
	c.GORM.Table("grid_servant_group gsg").
		Select(`
		gsg.*,
		gs.create_time AS gs_create_time,
		gs.id AS gs_id,
		gs.language gs_language,
		gs.servant_group_id AS gs_servant_group_id,
		gs.server_name AS gs_server_name,
		gs.location AS gs_location,
		gs.up_stream_name AS gs_up_stream_name
	`).
		Joins("LEFT JOIN grid_servant gs ON gs.servant_group_id = gsg.id").
		Find(&resp)

	return
}

// 转换函数
func ConvertToVoGroupByServant(voServantGroups []vo.VoServantGroup) []vo.VoGroupByServant {
	resultMap := make(map[int]vo.VoGroupByServant)
	for _, group := range voServantGroups {
		var servants []vo.VoServant
		if existingGroup, ok := resultMap[group.Id]; ok {
			servants = existingGroup.Servants
		}

		servants = append(servants, group.VoServant)

		resultMap[group.Id] = vo.VoGroupByServant{
			Id:             group.Id,
			TagName:        group.TagName,
			TagEnglishName: group.TagEnglishName,
			Servants:       servants,
		}
	}

	// 将 map 中的结果转换为切片
	var result []vo.VoGroupByServant
	for _, value := range resultMap {
		result = append(result, value)
	}

	return result
}
