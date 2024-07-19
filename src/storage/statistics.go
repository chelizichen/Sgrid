package storage

import "Sgrid/src/pool"

type KVI struct {
	Label string `json:"label" 	gorm:"column:label"`
	Value string `json:"value" 	gorm:"column:value"`
	Id    string `json:"id" 	gorm:"column:id"`
}

type StatisticsFunc = func() (rsp []KVI, err error)

var StatisticsMap = make(map[string]StatisticsFunc)

func init() {
	StatisticsMap["1"] = StatisticsGetServerPackage
	StatisticsMap["2"] = StatisticsGetLatestLog
	StatisticsMap["3"] = StatisticsGetStatus
	StatisticsMap["4"] = StatisticsGetServerType
	StatisticsMap["5"] = StatisticsGetServants
}

// 实时统计发包次数
func StatisticsGetServerPackage() (rsp []KVI, err error) {
	var sql = `
	select
		gs.server_name as label,
		count(gs.server_name) as value,
		gs.id as id
from
		grid_servant_package gsp
left join grid_servant gs on
		gsp.servant_id = gs.id
		GROUP BY  gsp.servant_id
	`
	err = pool.GORM.Raw(sql).Scan(&rsp).Error
	return rsp, err
}

// 最近20条日志更新
func StatisticsGetLatestLog() (rsp []KVI, err error) {
	var sql = `
SELECT
	concat(gtl.log_server_name,"(",gtl.log_host,")") as label,
	gtl.log_content as value,
	id
from
	grid_trace_log gtl
where
	gtl.log_type != "service-stat"
order by
	id desc
limit 0,20
	`
	err = pool.GORM.Raw(sql).Scan(&rsp).Error
	return rsp, err
}

// 服务状态实时更新
func StatisticsGetStatus() (rsp []KVI, err error) {
	var sql = `
	SELECT
		CONCAT(gs.server_name,"(",gg.pid ,")") as label,
		CASE gs.stat
		when 0 then '停止'
		when 1 then '运行'
		END as 'value',
		gg.id
from
		grid_grid gg
left join grid_servant gs on
		gg.servant_id = gs.id
where
		gs.stat != -1
	`
	err = pool.GORM.Raw(sql).Scan(&rsp).Error
	return rsp, err
}

// 服务类型
func StatisticsGetServerType() (rsp []KVI, err error) {
	var sql = `
	select gs.language as label ,count(gs.language) as value ,count(gs.language) as id
	from grid_servant gs group by gs.language
	`
	err = pool.GORM.Raw(sql).Scan(&rsp).Error
	return rsp, err
}

// 服务统计图
func StatisticsGetServants() (rsp []KVI, err error) {
	var sql = `
	select
		CONCAT(gs.server_name, '(', gs.protocol, ')') as label,
		gs.create_time as value,
		gs.id
from
		grid_servant gs
where
		gs.stat != -1
	`
	err = pool.GORM.Raw(sql).Scan(&rsp).Error
	return rsp, err
}
