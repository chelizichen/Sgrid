package pool

import (
	"Sgrid/src/config"
	sgridError "Sgrid/src/public/error"
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
	RDS_ADDR     = "redis_addr"
	RDS_PASS     = "redis_pass"
	DB_MASTER    = "db_master"
	DB_SLAVE     = "db_slave"
	DB_PREFIX    = "grid_"
)

var GORM *gorm.DB
var RDBContext = context.Background()
var GRDB *redis.Client

func InitStorage(ctx *config.SgridConf) {
	initDB(ctx)
	initRds(ctx)
}

func initDB(ctx *config.SgridConf) error {
	db_master := ctx.GetString(DB_MASTER)
	db_slave := ctx.GetString(DB_SLAVE)
	if len(db_master) == 0 {
		return sgridError.DB_CONN_ERROR("len(db_master) == 0")
	}
	db, err := gorm.Open(mysql.Open(db_master), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   DB_PREFIX,
			SingularTable: true,
		},
	})
	if err != nil {
		return sgridError.DB_CONN_ERROR(fmt.Sprintf("gorm.Open %v", err.Error()))
	}

	if len(db_slave) != 0 {
		fmt.Println("****** db_slave", db_slave)
		fmt.Println("****** db_master", db_master)
		if e := db.Use(dbresolver.Register(dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(db_master)},
			Replicas: []gorm.Dialector{mysql.Open(db_slave)},
		})); e != nil {
			return sgridError.DB_CONN_ERROR(fmt.Sprintf("use db slave error %v", e.Error()))
		}
	}
	sqlDB, err := db.DB()
	if err != nil {
		return sgridError.DB_CONN_ERROR(fmt.Sprintf("Error To init sdb.DB %v", err.Error()))
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
		db.Debug().AutoMigrate(&pojo.AssetsAdmin{})
		// ****************** rbac ********************
		db.Debug().AutoMigrate(&rbac.User{})
		db.Debug().AutoMigrate(&rbac.UserRole{})
		db.Debug().AutoMigrate(&rbac.UserToRole{})
		db.Debug().AutoMigrate(&rbac.RoleMenu{})
		db.Debug().AutoMigrate(&rbac.RoleToMenu{})
		db.Debug().AutoMigrate(&rbac.VersionUpdateLine{})
		db.Debug().AutoMigrate(&rbac.UserGroup{})
		db.Debug().AutoMigrate(&rbac.UserToUserGroup{})
	}
	GORM = db
	return nil
}

func initRds(ctx *config.SgridConf) error {
	redis_addr := ctx.GetString(RDS_ADDR)
	redis_pass := ctx.GetString(RDS_PASS)
	if len(redis_addr) == 0 {
		return sgridError.RDS_CONN_ERROR("len(redis_addr) == 0")
	}
	GRDB = redis.NewClient(&redis.Options{
		Addr:     redis_addr,
		Password: redis_pass,
		DB:       0,
	})
	pong, err := GRDB.Ping(RDBContext).Result()
	if err != nil {
		return sgridError.RDS_CONN_ERROR("连接redis出错，错误信息：" + err.Error())
	}
	fmt.Println("conn success! ", pong)
	return nil
}
