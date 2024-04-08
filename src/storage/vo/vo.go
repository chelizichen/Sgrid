package vo

import "Sgrid/src/storage/pojo"

// 节点

type ServantVo struct {
	Id           int
	ServerName   string            // 服务名称
	CreateTime   string            `gorm:"autoCreateTime"` // 创建时间
	Language     string            // 语言
	ServantGroup pojo.ServantGroup // 服务组ID
	Node         pojo.Node         // 服务节点
}

type ServantPackageVo struct {
	Id         int
	Servant    pojo.Servant // 服务Id
	Hash       string       // Hash值
	FilePath   string       // 文件路径
	Content    string       // 上传内容
	CreateTime string       `gorm:"autoCreateTime"` // 创建时间
}

type Grid struct {
	ID          int `gorm:"column:id" json:"id,omitempty"`
	ServantID   int `gorm:"column:servant_id" json:"servantId,omitempty"`
	Port        int `gorm:"column:port" json:"port,omitempty"`
	NodeID      int `gorm:"column:node_id" json:"nodeId,omitempty"`
	Status      int `gorm:"column:status" json:"status,omitempty"`
	Pid         int `gorm:"column:pid" json:"pid,omitempty"`
	UpdateTime  int `gorm:"column:update_time" json:"updateTime,omitempty"`
	GridServant `gorm:"embedded" json:"gridServant,omitempty"`
	GridNode    `gorm:"embedded" json:"gridNode,omitempty"`
}

type GridServant struct {
	ID                int    `gorm:"column:gs_id" json:"servantId,omitempty"`
	Language          string `gorm:"column:gs_language" json:"language,omitempty"`
	ServantGroupID    int    `gorm:"column:gs_servant_group_id" json:"servantGroupId,omitempty"`
	ServerName        string `gorm:"column:gs_server_name" json:"serverName,omitempty"`
	ServantCreateTime string `gorm:"column:gs_create_time" json:"servantCreateTime,omitempty"`
}

type GridNode struct {
	NodeID         int    `gorm:"column:gn_id" json:"nodeId,omitempty"`
	NodeIP         string `gorm:"column:gn_ip" json:"nodeIp,omitempty"`
	Main           string `gorm:"column:gn_main" json:"main,omitempty"`
	Platform       string `gorm:"column:gn_pslat_form" json:"platform,omitempty"`
	NodeStatus     int    `gorm:"column:gn_status" json:"nodeStatus,omitempty"`
	NodeCreateTime string `gorm:"column:gn_create_time" json:"nodeCreateTime,omitempty"`
}
