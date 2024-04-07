package storage

import (
	c "Sgrid/src/configuration"
	"Sgrid/src/storage/pojo"
)

func SaveHashPackage(pkg pojo.ServantPackage) int {
	c.GORM.Create(&pkg)
	return pkg.Id
}

func SaveServant(svr pojo.Servant) int {
	c.GORM.Create(&svr)
	return svr.Id
}

func SaveServantGroup(group pojo.ServantGroup) int {
	c.GORM.Create(&group)
	return group.Id
}

func UpdateGrid(group pojo.Grid) int {
	if group.Id != 0 {
		c.GORM.Debug().
			Model(&group).
			Select("content", "title").
			Where("id = ?", group.Id).
			Updates(&pojo.Grid{
				Status: group.Status,
				Pid:    group.Status,
			})
		return (group.Id)
	}
	c.GORM.Create(&group)
	return (group.Id)
}
