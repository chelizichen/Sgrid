package vo

import "Sgrid/src/storage/pojo"

// 服务
type ServantVo struct {
	Id           int
	ServerName   string            // 服务名称
	CreateTime   string            `gorm:"autoCreateTime"` // 创建时间
	Language     string            // 语言
	ServantGroup pojo.ServantGroup // 服务组ID
}

// 服务网格 用于查看所有节点信息
type GridVo struct {
	Id         int
	Servant    pojo.Servant // 网格容纳服务
	Port       int          // 网格端口
	Node       pojo.Node    // 网格所属节点
	Status     int          // 网格状态
	Pid        int          // 网格Pid
	UpdateTime int          `gorm:"autoCreateTime"` // 更新时间
}

type ServantPackageVo struct {
	Id         int
	Servant    pojo.Servant // 服务Id
	Hash       string       // Hash值
	FilePath   string       // 文件路径
	Content    string       // 上传内容
	CreateTime string       `gorm:"autoCreateTime"` // 创建时间
}
