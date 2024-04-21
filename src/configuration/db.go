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
	"gorm.io/plugin/dbresolver"
)

const (
	AUTO_MIGRATE = "db_migrate"
)

var GORM *gorm.DB
var RDBContext = context.Background()

func InitStorage(ctx *config.SgridConf) {
	db_master := ctx.GetString("db_master")
	db_slave := ctx.GetString("db_slave")

	db, err := gorm.Open(mysql.Open(db_master), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "grid_",
			SingularTable: true,
		},
	})

	if len(db_slave) != 0 {
		db.Use(dbresolver.Register(dbresolver.Config{
			Sources: []gorm.Dialector{mysql.Open(db_slave)},
		}))
	}

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
	if ctx.GetBool(AUTO_MIGRATE) {
		db.Debug().AutoMigrate(&pojo.Node{})
		db.Debug().AutoMigrate(&pojo.User{})
		db.Debug().AutoMigrate(&pojo.Grid{})
		db.Debug().AutoMigrate(&pojo.ServantPackage{})
		db.Debug().AutoMigrate(&pojo.ServantGroup{})
		db.Debug().AutoMigrate(&pojo.Servant{})
		db.Debug().AutoMigrate(&pojo.Properties{})
		db.Debug().AutoMigrate(&pojo.StatLog{})
	}
	GORM = db
}
