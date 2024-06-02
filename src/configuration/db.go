package configuration

import (
	"Sgrid/src/config"
	"Sgrid/src/storage/pojo"
	"Sgrid/src/storage/rbac"
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
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
var GRDB *redis.Client

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
		fmt.Println("db_slave", db_slave)
		fmt.Println("db_master", db_master)
		if e := db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(db_master)},
			Replicas: []gorm.Dialector{mysql.Open(db_slave)},
		})); e != nil {
			fmt.Println("e", e.Error())
		}
	}
	// db
	if err != nil {
		fmt.Println("Error To init gorm", err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		fmt.Println("Error To init sdb.DB", err)
	}
	// rds
	redis_addr := ctx.GetString("redis-addr")
	redis_pass := ctx.GetString("redis-pass")
	GRDB = redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_pass,
		DB:       0,
	})
	pong, err := GRDB.Ping(RDBContext).Result()

	if err != nil {
		fmt.Printf("连接redis出错，错误信息：%v", err.Error())
	} else {
		fmt.Println("成功连接redis", pong)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)
	if ctx.GetBool(AUTO_MIGRATE) {
		db.Debug().AutoMigrate(&pojo.Node{})
		db.Debug().AutoMigrate(&pojo.Grid{})
		db.Debug().AutoMigrate(&pojo.ServantPackage{})
		db.Debug().AutoMigrate(&pojo.ServantGroup{})
		db.Debug().AutoMigrate(&pojo.Servant{})
		db.Debug().AutoMigrate(&pojo.Properties{})
		db.Debug().AutoMigrate(&pojo.StatLog{})
		db.Debug().AutoMigrate(&pojo.SystemErr{})
		db.Debug().AutoMigrate(&pojo.ServantConf{})
		db.Debug().AutoMigrate(&pojo.TraceLog{})

		// rbac
		db.Debug().AutoMigrate(&rbac.User{})
		db.Debug().AutoMigrate(&rbac.UserRole{})
		db.Debug().AutoMigrate(&rbac.UserToRole{})
		db.Debug().AutoMigrate(&rbac.RoleMenu{})
		db.Debug().AutoMigrate(&rbac.RoleToMenu{})

	}
	GORM = db
}
