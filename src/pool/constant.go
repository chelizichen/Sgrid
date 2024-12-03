package pool

import (
	"time"
)

const (
	AUTO_MIGRATE = "db_migrate"
	RDS_ADDR     = "redis_addr"
	RDS_PASS     = "redis_pass"
	DB_MASTER    = "db_master"
	DB_SLAVE     = "db_slave"
	DB_PREFIX    = "grid_"
)

// 数据库配置常量
const (
	DefaultMaxIdleConns    = 10
	DefaultMaxOpenConns    = 100
	DefaultConnMaxLifetime = time.Hour
	DefaultDBPrefix        = "sgrid_"
)
