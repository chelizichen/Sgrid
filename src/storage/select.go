package storage

import (
	c "Sgrid/src/configuration"
	"Sgrid/src/storage/pojo"
	"Sgrid/src/storage/vo"
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

func QueryGrid(ServantId int) []vo.GridVo {
	var dataList []vo.GridVo
	c.GORM.Model(dataList).Find(&dataList)
	return dataList
}

func ServantPackage(ServantId int) []vo.GridVo {
	var dataList []vo.GridVo
	c.GORM.Model(dataList).Find(&dataList)
	return dataList
}
