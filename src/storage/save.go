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

func SaveStatLog(stat *pojo.StatLog) {
	c.GORM.Debug().Create(&stat)
}

func UpdateGrid(group *pojo.Grid) int {
	if group.Id == 0 {
		c.GORM.Debug().Create(&group)
		return (group.Id)
	} else {
		c.GORM.Debug().
			Model(&group).
			Select("status", "pid").
			Where("id = ?", group.Id).
			Updates(&pojo.Grid{
				Status: group.Status,
				Pid:    group.Pid,
			})
		return (group.Id)
	}
}

func DeletePackage(id int) {
	c.GORM.Debug().
		Model(&pojo.ServantPackage{}).
		Select("status").
		Where("id = ?", id).
		Updates(&pojo.Grid{
			Status: -1,
		})
}
