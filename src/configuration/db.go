package configuration

import (
	"Sgrid/src/config"
	"Sgrid/src/storage/pojo"
	"context"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var GORM *gorm.DB
var RDBContext = context.Background()

func InitStorage(ctx *config.SgridConf) {

	db, err := gorm.Open(mysql.Open(ctx.Server.Storage), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "grid_",
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Println("Error To init gorm", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Error To init db.DB", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	db.Debug().AutoMigrate(&pojo.Node{})
	db.Debug().AutoMigrate(&pojo.Grid{})
	db.Debug().AutoMigrate(&pojo.ServantPackage{})
	db.Debug().AutoMigrate(&pojo.ServantGroup{})
	db.Debug().AutoMigrate(&pojo.Servant{})
	db.Debug().AutoMigrate(&pojo.Properties{})
	db.Debug().AutoMigrate(&pojo.StatLog{})
	GORM = db

}
