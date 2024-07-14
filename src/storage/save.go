package storage

import (
	protocol "Sgrid/server/SgridLogTraceServer/proto"
	"Sgrid/src/pool"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/pojo"
	"errors"
	"fmt"
	"time"
)

func SaveHashPackage(pkg pojo.ServantPackage) int {
	pool.GORM.Create(&pkg)
	return pkg.Id
}

func GetServantByNameAndGroup(name string, groupId int) bool {
	var svr pojo.Servant
	pool.GORM.Where("server_name = ? and servant_group_id = ?", name, groupId).First(&svr)
	fmt.Println("svr", svr)
	return svr.Id != 0
}

func SaveServant(svr *pojo.Servant) int {
	pool.GORM.Create(&svr)
	return svr.Id
}

func UpdateServant(svr *pojo.Servant) int {
	pool.GORM.
		Model(&svr).
		Where("id = ?", svr.Id).
		Updates(&pojo.Servant{
			ServantGroupId: svr.ServantGroupId,
			ExecPath:       svr.ExecPath,
			Protocol:       svr.Protocol,
			Language:       svr.Language,
			Preview:        svr.Preview,
		})
	return svr.Id
}

func DelServant(svr *pojo.Servant) int {
	fmt.Println("svr", svr.Stat)
	fmt.Println("svr", svr.Id)
	pool.GORM.Debug().
		Model(&svr).
		Where("id = ?", svr.Id).
		Update("stat", svr.Stat)
	return svr.Id
}

func SaveServantGroup(group *pojo.ServantGroup) int {
	if group.Id == 0 {
		pool.GORM.Create(&group)
		return group.Id
	} else {
		pool.GORM.Debug().
			Model(&group).
			Where("id = ?", group.Id).
			Updates(&pojo.ServantGroup{
				TagName:        group.TagName,
				TagEnglishName: group.TagEnglishName,
			})
		return group.Id
	}
}

func SaveStatLog(stat *pojo.StatLog) {
	pool.GORM.Debug().Create(&stat)
}

func UpdateGrid(grid *pojo.Grid) int {
	if grid.Id == 0 {
		pool.GORM.Debug().Create(&grid)
		return (grid.Id)
	} else {
		pool.GORM.Debug().
			Model(&grid).
			Select("status", "pid").
			Where("id = ?", grid.Id).
			Updates(&pojo.Grid{
				Status: grid.Status,
				Pid:    grid.Pid,
			})
		return (grid.Id)
	}
}

func DeleteGrid(id int) {
	pool.GORM.Debug().
		Model(&pojo.Grid{}).
		Delete(&pojo.Grid{
			Id: id,
		})
}

func DeletePackage(id int) error {
	return pool.GORM.Debug().
		Model(&pojo.ServantPackage{}).
		Select("status").
		Where("id = ?", id).
		Updates(&pojo.Grid{
			Status: -1,
		}).Error
}

// 判断服务组下有无存在的服务，如果没有，则删除服务组
func DelGroup(id int) error {
	var count int64
	pool.GORM.Model(&pojo.Servant{}).Where("servant_group_id = ? and stat != -1", id).Count(&count)
	if count > 0 {
		st := fmt.Sprintf("服务组下存在%v个服务，不能删除", count)
		return errors.New(st)
	}
	pool.GORM.Debug().
		Model(&pojo.ServantGroup{}).
		Delete(&pojo.ServantGroup{
			Id: id,
		})

	return nil
}

func UpdatePackageVersion(dto *pojo.ServantPackage) {
	pool.GORM.Debug().
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
	pool.GORM.Debug().Create(obj)
	return obj.Id
}

func PushErr(d *pojo.SystemErr) {
	pool.GORM.Debug().Create(d)
}

func UpdateConf(d *pojo.ServantConf) {
	if d.Id == 0 {
		pool.GORM.Debug().Create(d)
		return
	}
	pool.GORM.
		Model(&pojo.ServantConf{}).
		Where(&pojo.ServantConf{
			ServantId: d.ServantId,
		}).
		Updates(d)
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
	err = pool.GORM.Debug().Create(pj).Error
	return err
}

func UpsertProperty(p *pojo.Properties) {
	if p.Id == 0 {
		pool.GORM.Model(&pojo.Properties{}).Create(p)
		return
	}
	pool.GORM.Model(&pojo.Properties{}).
		Where("id = ?", p.Id).
		Updates(&pojo.Properties{
			Key:   p.Key,
			Value: p.Value,
		})
}

func DelProperty(id int) {
	pool.GORM.Model(&pojo.Properties{}).Delete(&pojo.Properties{
		Id: id,
	})
}

func UpsertAssets(d *pojo.AssetsAdmin) error {
	fmt.Println("d", d.GridId)
	err := DelAssetById(d.GridId)
	if err != nil {
		return err
	}
	return pool.GORM.
		Debug().
		Model(&pojo.AssetsAdmin{}).
		Create(&d).Error
}

func DelAssetById(id int) error {
	return pool.GORM.
		Model(&pojo.AssetsAdmin{}).
		Delete("grid_id = ?", id).Error
}

func GetAssetById(id int) (resp *pojo.AssetsAdmin, err error) {
	err = pool.GORM.
		Model(&pojo.AssetsAdmin{}).
		Where(&pojo.AssetsAdmin{
			GridId: id,
		}).Find(&resp).Error
	return resp, err
}
