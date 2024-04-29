package pojo

import "time"

// 节点
type Node struct {
	Id         int
	Ip         string     // IP地址
	Status     int        // 状态
	CreateTime *time.Time `gorm:"autoCreateTime"` // 创建时间
	PlatForm   string     // 平台
	Main       string     // 是否为主机
	UploadPath string     // 上传路径
}

// 服务组
type ServantGroup struct {
	Id             int
	TagName        string     // 服务标签
	TagEnglishName string     // 英文
	CreateTime     *time.Time `gorm:"autoCreateTime"` // 创建时间
}

// 服务
type Servant struct {
	Id             int
	ServerName     string     // 服务名称
	CreateTime     *time.Time `gorm:"autoCreateTime"` // 创建时间
	Language       string     // 语言
	UpStreamName   string     // 转发名称
	Location       string     // 路径
	Protocol       string     // 协议
	ExecPath       string     // 可执行路径
	ServantGroupId int        `gorm:"foreignKey:ServantGroupId"` // 服务组ID
}

// 服务网格 用于查看所有节点信息
type Grid struct {
	Id         int
	Port       int        // 网格端口
	Status     int        // 网格状态
	Pid        int        // 网格Pid
	UpdateTime *time.Time `gorm:"autoCreateTime"`            // 更新时间
	NodeId     int        `gorm:"foreignKey:NodeId"`         // 网格所属节点ID
	ServantId  int        `gorm:"foreignKey:ServantGroupId"` // 网格容纳服务ID
	CreateTime *time.Time `gorm:"autoCreateTime"`            // 网格容纳服务ID
}

// 服务包
type ServantPackage struct {
	Id         int
	ServantId  int        // 服务Id
	Hash       string     // Hash值
	FilePath   string     // 文件路径
	Content    string     // 上传内容
	Version    string     // 版本号
	CreateTime *time.Time `gorm:"autoCreateTime"` // 创建时间
	Status     int        // 文件状态 -1 为已删除不可用
}

type Properties struct {
	Id         int
	Key        string
	Value      string
	CreateTime *time.Time `gorm:"autoCreateTime"` // 创建时间
}

type StatLog struct {
	Id          int        `json:"id,omitempty"` // id
	GridId      int        `json:"gridId,omitempty"`
	Stat        string     `json:"stat,omitempty"`
	Pid         int        `json:"pid,omitempty"`
	CreateTime  *time.Time `gorm:"autoCreateTime" json:"createTime,omitempty"`
	CPU         float64    `json:"cpu,omitempty"`
	Threads     int32      `json:"threads,omitempty"`
	Name        string     `json:"name,omitempty"`
	IsRunning   string     `json:"isRunning,omitempty"`
	MemoryStack uint64     `json:"memoryStack,omitempty"`
	MemoryData  uint64     `json:"memoryData,omitempty"`
}

type User struct {
	Id         int
	UserName   string
	Password   string
	CreateTime string
	RoleId     int
}

type ServantConf struct {
	Id         int
	ServantId  string
	CreateTime *time.Time
	Conf       string
}

type SystemErr struct {
	Id         int        `json:"id,omitempty"`
	CreateTime *time.Time `gorm:"autoCreateTime" json:"createTime,omitempty"`
	Type       string     `json:"type,omitempty"`
	Info       string     `json:"info,omitempty"`
}
