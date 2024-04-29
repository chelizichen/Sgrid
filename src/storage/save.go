package storage

import (
	c "Sgrid/src/configuration"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/pojo"
	"fmt"
	"time"
)

func SaveHashPackage(pkg pojo.ServantPackage) int {
	c.GORM.Create(&pkg)
	return pkg.Id
}

func SaveServant(svr *pojo.Servant) int {
	c.GORM.Create(&svr)
	return svr.Id
}

func SaveServantGroup(group *pojo.ServantGroup) int {
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

func DeleteGrid(id int) {
	c.GORM.Debug().
		Model(&pojo.Grid{}).
		Delete(&pojo.Grid{
			Id: id,
		})
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

func UpdateNode(d *dto.NodeDTO) int {
	fmt.Println("d", d)
	t := time.Now()
	obj := &pojo.Node{
		PlatForm:   d.PlatForm,
		UploadPath: d.UploadPath,
		Status:     d.Status,
		Main:       d.Main,
		Ip:         d.Ip,
		CreateTime: &t,
	}
	c.GORM.Debug().Create(obj)
	return obj.Id
}

func PushErr(d *pojo.SystemErr) {
	c.GORM.Debug().Create(d)
}
