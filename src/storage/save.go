package storage

import (
	protocol "Sgrid/server/SgridLogTraceServer/proto"
	c "Sgrid/src/configuration"
	"Sgrid/src/public"
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

func UpdateServant(svr *pojo.Servant) int {
	c.GORM.
		Model(&svr).
		Where("id = ?", svr.Id).
		Updates(&pojo.Servant{
			ServantGroupId: svr.ServantGroupId,
			ExecPath:       svr.ExecPath,
			Protocol:       svr.Protocol,
		})
	return svr.Id
}

func DelServant(svr *pojo.Servant) int {
	c.GORM.
		Model(&svr).
		Where("id = ?", svr.Id).
		Updates(&pojo.Servant{
			Stat: public.STAT_SERVANT_DELETE,
		})
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

func UpdatePackageVersion(dto *pojo.ServantPackage) {
	c.GORM.Debug().
		Model(dto).
		Where("id = ?", dto.Id).
		Updates(&pojo.ServantPackage{
			Version: dto.Version,
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

func UpdateConf(d *pojo.ServantConf) {
	if d.Id == 0 {
		CreateConf(d)
		return
	}
	c.GORM.
		Model(&pojo.ServantConf{}).
		Where(&pojo.ServantConf{
			ServantId: d.ServantId,
		}).
		Updates(d)
}

func CreateConf(d *pojo.ServantConf) {
	c.GORM.Debug().Create(d)
}

func SaveLog(d *protocol.LogTraceReq) error {
	t, err := time.Parse(time.DateTime, d.CreateTime)
	if err != nil {
		PushErr(&pojo.SystemErr{
			Type: "system/error/saveLog/time.Parse",
			Info: err.Error(),
		})
		return err
	}
	pj := &pojo.TraceLog{
		CreateTime:    &t,
		LogServerName: d.LogServerName,
		LogHost:       d.LogHost,
		LogType:       d.LogType,
		LogContent:    d.LogContent,
		LogGridId:     d.LogGridId,
		LogBytesLen:   d.LogBytesLen,
	}
	err = c.GORM.Debug().Create(pj).Error
	return err
}

func UpsertProperty(p *pojo.Properties) {
	if p.Id == 0 {
		c.GORM.Model(&pojo.Properties{}).Create(p)
		return
	}
	c.GORM.Model(&pojo.Properties{}).
		Where("id = ?", p.Id).
		Updates(&pojo.Properties{
			Key:   p.Key,
			Value: p.Value,
		})
}

func DelProperty(id int) {
	c.GORM.Model(&pojo.Properties{}).Delete(&pojo.Properties{
		Id: id,
	})
}
