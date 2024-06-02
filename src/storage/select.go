package storage

import (
	c "Sgrid/src/configuration"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/pojo"
	"Sgrid/src/storage/rbac"
	"Sgrid/src/storage/vo"
	"Sgrid/src/utils"
	"fmt"
)

// 查询标签组，创建时用
func QueryTags() []pojo.ServantGroup {
	var dataList []pojo.ServantGroup
	c.GORM.Model(dataList).Find(&dataList)
	return dataList
}

// 查询节点组
func QueryNodes() []*vo.GridNode {
	var dataList []pojo.Node
	c.GORM.Model(&pojo.Node{}).Find(&dataList)
	var respList []*vo.GridNode
	for _, v := range dataList {
		V := v
		respList = append(respList, &vo.GridNode{
			NodeID:         V.Id,
			NodeIP:         V.Ip,
			NodeStatus:     V.Status,
			NodeCreateTime: V.CreateTime,
			Platform:       V.PlatForm,
			Main:           V.Main,
		})
	}
	return respList
}

func QueryServants() *[]vo.VoServantObj {
	var dataList []vo.VoServantObj
	c.GORM.Model(&pojo.Servant{}).Find(&dataList)
	return &dataList
}

func QueryGrid(req *dto.PageBasicReq) (resp []vo.Grid) {
	args := make([]interface{}, 10)
	where := "1 = 1"
	if req.Id != 0 {
		where += " and gs.id  = ?"
		args = append(args, req.Id)
	}
	c.GORM.Debug().Table("grid_grid gg").
		Select(`
	gg.*,
	gs.id AS gs_id,
	gs.language AS gs_language,
	gs.servant_group_id AS gs_servant_group_id,
	gs.server_name as gs_server_name,
	gs.create_time as gs_create_time,
	gs.exec_path as gs_exec_path,
	gs.protocol as gs_protocol,
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

// 2024.6.1 pageBasicReq.id as user_id
func QueryServantGroup(req *dto.PageBasicReq) (resp []vo.VoServantGroup) {
	if req.Id != 0 {
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
			Where("gs.user_id = ?", req.Id).
			Find(&resp)
	} else {
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
	}

	return
}

// 转换函数
func ConvertToVoGroupByServant(voServantGroups []vo.VoServantGroup) []vo.VoGroupByServant {
	resultMap := make(map[int]vo.VoGroupByServant)
	ept := vo.VoServant{}
	for _, group := range voServantGroups {
		var servants []vo.VoServant
		if existingGroup, ok := resultMap[group.Id]; ok {
			servants = existingGroup.Servants
		}

		if group.VoServant != ept {
			servants = append(servants, group.VoServant)
		}
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

func QueryPackage(queryPackageDto *dto.QueryPackageDto) []vo.VoServantPackage {
	var queryPackageResp []vo.VoServantPackage
	fmt.Println("queryPackageDto", queryPackageDto)
	where := `1 = 1`
	params := make([]any, 0)
	if queryPackageDto.Id != 0 {
		where += ` and gsp.servant_id = ?`
		params = append(params, queryPackageDto.Id)
	}
	if len(queryPackageDto.Version) != 0 {
		where += ` and gsp.version = ?`
		params = append(params, queryPackageDto.Version)
	}
	c.GORM.Table("grid_servant_package gsp").
		Select(`
	gsp.*,
	gs.server_name as gs_server_name,
	gs.create_time as gs_create_time,
	gs.language  as gs_language,
	gs.exec_path as gs_exec_path,
	gs.protocol as gs_protocol
	`).Joins(`
	left join grid_servant gs on
	gsp.servant_id = gs.id
	`).Where(where, utils.Removenullvalue(params)...).
		Order(" gsp.create_time  DESC").
		Find(&queryPackageResp)
	return queryPackageResp
}

func QueryPackageById(id int) (rsp pojo.ServantPackage) {
	c.GORM.Model(&pojo.ServantPackage{}).
		Where(&pojo.ServantPackage{
			Id: id,
		}).
		Find(&rsp)
	return rsp
}

func QueryStatLogById(id int, offset int, size int) any {
	var total int64
	var rsp []pojo.StatLog
	c.GORM.
		Model(rsp).
		Where(&pojo.StatLog{
			GridId: id,
		}).
		Count(&total).
		Limit(size).
		Offset(offset).
		Order("create_time desc").
		Find(&rsp)
	resp := make(map[string]interface{})
	resp["list"] = rsp
	resp["total"] = total
	return resp
}

func QueryPropertiesByKey(key string) (resp []*pojo.Properties) {
	c.GORM.Model(&pojo.Properties{}).Where(&pojo.Properties{
		Key: key,
	}).Find(&resp)
	return resp
}

func QueryProperties() (resp []*pojo.Properties) {
	c.GORM.Model(&pojo.Properties{}).Find(&resp)
	return
}

func QueryUser(userName string) (resp *rbac.User, err error) {
	err = c.GORM.Model(&rbac.User{}).Where(&rbac.User{
		UserName: userName,
	}).Find(&resp).Error
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func QueryGroups() (resp *[]vo.VoGroupObj) {
	c.GORM.Model(&pojo.ServantGroup{}).Find(&resp)
	fmt.Println("resp", resp)
	return
}

func GetGridByNodePort(nodeId int, port int) int64 {
	var total int64
	c.GORM.Model(&pojo.Grid{}).Where(&pojo.Grid{
		NodeId: nodeId,
		Port:   port,
	}).Count(&total)
	return total
}

func GetServantConfById(ServantId int) (resp *pojo.ServantConf) {
	c.GORM.Model(&pojo.ServantConf{}).Where(&pojo.ServantConf{
		ServantId: &ServantId,
	}).Find(&resp)
	return resp
}

func GetTraceLog(req *dto.TraceLogDto) ([]string, int64, error) {
	var resp []*pojo.TraceLog
	var total int64
	err := c.GORM.Model(&pojo.TraceLog{}).
		Where("log_content like ?", "%"+req.Keyword+"%").
		Where("log_grid_id = ?", req.LogGridId).
		Where("date(create_time) = ?", req.SearchTime).
		Where("log_type = ?", req.LogType).
		Count(&total).
		Offset(req.Offset).
		Limit(req.Size).
		Find(&resp).Error
	var log2String []string
	for _, v := range resp {
		log2String = append(log2String, v.FmtGetLog())
	}
	return log2String, total, err
}

type TraceLogFileVo struct {
	LogType string `gorm:"column:log_type`
	LogTime string `gorm:"column:log_time`
}

func GetTraceLogFiles(gridId int, log_server_name string) []TraceLogFileVo {
	var selectResp []TraceLogFileVo
	where := `1 = 1`
	params := make([]any, 0)
	if len(log_server_name) > 0 {
		where += " AND gtl.log_server_name = ?"
		params = append(params, log_server_name)
	}
	if gridId == 0 {
		where += " AND gtl.log_grid_id = ?"
		params = append(params, gridId)
	}
	c.GORM.Debug().
		Table("grid_trace_log gtl").
		Select(`
	gtl.log_type as log_type,
	date(gtl.create_time) as log_time
	`).
		Where(where, utils.Removenullvalue(params)...).
		Group("log_type").
		Group("log_time").
		Order("log_time").
		Find(&selectResp)
	fmt.Println("selectResp", selectResp)
	return selectResp
}
