package vo

import (
	"Sgrid/src/config"
	"Sgrid/src/storage/pojo"
	"time"
)

type VoGroupByServant struct {
	Id             int         `json:"id,omitempty"`
	TagName        string      `json:"tagName,omitempty"`        // 服务标签
	TagEnglishName string      `json:"tagEnglishName,omitempty"` // 英文
	CreateTime     string      `json:"createTime,omitempty"`     // 创建时间
	Servants       []VoServant `json:"servants"`                 // 服务组
}

// 节点
type VoServantGroup struct {
	Id             int       `gorm:"column:id" json:"id,omitempty"`
	TagName        string    `gorm:"column:tag_name" json:"tagName,omitempty"`                // 服务标签
	TagEnglishName string    `gorm:"column:tag_english_name" json:"tagEnglishName,omitempty"` // 英文
	CreateTime     string    `gorm:"autoCreateTime" json:"creatTime,omitempty"`               // 创建时间
	VoServant      VoServant `gorm:"embedded" json:"servantGroup,omitempty"`
}

type VoGroupObj struct {
	Id             int    `gorm:"column:id" json:"id,omitempty"`
	TagName        string `gorm:"column:tag_name" json:"tagName,omitempty"`                // 服务标签
	TagEnglishName string `gorm:"column:tag_english_name" json:"tagEnglishName,omitempty"` // 英文
	CreateTime     string `gorm:"column:create_time" json:"creatTime,omitempty"`           // 创建时间
	UserId         int    `gorm:"column:user_id" json:"userId,omitempty"`                  // 用户ID
	UserGroupId    int    `gorm:"column:user_group_id" json:"user_group_id,omitempty"`     // 用户服务组ID
}

// 联表查用的
type VoServant struct {
	Id             int    `gorm:"column:gs_id" json:"id,omitempty"`
	ServerName     string `gorm:"column:gs_server_name" json:"serverName,omitempty"`      // 服务名称
	CreateTime     string `gorm:"column:gs_create_time" json:"createTime,omitempty"`      // 创建时间
	Language       string `gorm:"column:gs_language" json:"language,omitempty"`           // 语言
	UpStreamName   string `gorm:"column:gs_up_stream_name" json:"upStreamName,omitempty"` // 转发名称
	Location       string `gorm:"column:gs_location"  json:"location,omitempty"`          // 路径
	ServantGroupId int    `gorm:"column:gs_groupId" json:"servantGroupId,omitempty"`      // 服务组ID
	Stat           int    `gorm:"column:gs_stat" json:"stat,omitempty"`                   // 服务组ID
}

// 纯Obj
type VoServantObj struct {
	Id             int        `json:"id,omitempty"`
	ServerName     string     `json:"serverName,omitempty"`     // 服务名称
	CreateTime     *time.Time `json:"createTime,omitempty"`     // 创建时间
	Language       string     `json:"language,omitempty"`       // 语言
	UpStreamName   string     `json:"upStreamName,omitempty"`   // 转发名称
	Location       string     `json:"location,omitempty"`       // 路径
	Protocol       string     `json:"protocol,omitempty"`       // 协议
	ExecPath       string     `json:"execPath,omitempty"`       // 可执行路径
	ServantGroupId int        `json:"servantGroupId,omitempty"` // 服务组ID
	Stat           int        `json:"stat,omitempty"`           // 服务组ID
	Preview        string     `json:"preview,omitempty"`        // 预览地址
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
	ID          int        `gorm:"column:id" json:"id,omitempty"`
	ServantID   int        `gorm:"column:servant_id" json:"servantId,omitempty"`
	Port        int        `gorm:"column:port" json:"port,omitempty"`
	NodeID      int        `gorm:"column:node_id" json:"nodeId,omitempty"`
	Status      int        `gorm:"column:status" json:"status,omitempty"`
	Pid         int        `gorm:"column:pid" json:"pid,omitempty"`
	UpdateTime  *time.Time `gorm:"column:update_time" json:"updateTime,omitempty"`
	GridServant `gorm:"embedded" json:"gridServant,omitempty"`
	GridNode    `gorm:"embedded" json:"gridNode,omitempty"`
}

type GridServant struct {
	ID                int        `gorm:"column:gs_id" json:"servantId,omitempty"`
	Language          string     `gorm:"column:gs_language" json:"language,omitempty"`
	ServantGroupID    int        `gorm:"column:gs_servant_group_id" json:"servantGroupId,omitempty"`
	ServerName        string     `gorm:"column:gs_server_name" json:"serverName,omitempty"`
	ServantCreateTime *time.Time `gorm:"column:gs_create_time" json:"servantCreateTime,omitempty"`
	ExecPath          string     `gorm:"column:gs_exec_path" json:"execPath,omitempty"`
	Protocol          string     `gorm:"column:gs_protocol" json:"protocol,omitempty"`
	Preview           string     `gorm:"column:gs_preview" json:"preview,omitempty"`
}

type GridNode struct {
	NodeID         int        `gorm:"column:gn_id" json:"id,omitempty"`
	NodeIP         string     `gorm:"column:gn_ip" json:"ip,omitempty"`
	Main           string     `gorm:"column:gn_main" json:"main,omitempty"`
	Platform       string     `gorm:"column:gn_pslat_form" json:"platform,omitempty"`
	NodeStatus     int        `gorm:"column:gn_status" json:"nodeStatus,omitempty"`
	NodeCreateTime *time.Time `gorm:"column:gn_create_time" json:"nodeCreateTime,omitempty"`
}

type CoverConfigVo struct {
	Conf       config.SgridConf
	ServerName string
}

type VoServantPackage struct {
	ID           uint       `gorm:"id" json:"id,omitempty"`
	ServantID    uint       `gorm:"servant_id" json:"servantId,omitempty"`
	Hash         string     `gorm:"hash" json:"hash,omitempty"`
	FilePath     string     `gorm:"file_path" json:"filePath,omitempty"`
	Content      string     `gorm:"content" json:"content,omitempty"`
	CreateTime   *time.Time `gorm:"create_time" json:"createTime,omitempty"`
	Version      string     `gorm:"version" json:"version,omitempty"`
	Status       string     `gorm:"status" json:"status,omitempty"`
	GSServerName string     `gorm:"gs_server_name" json:"serverName,omitempty"`
	GSCreateTime *time.Time `gorm:"gs_create_time" json:"serverCreateTime,omitempty"`
	GSLanguage   string     `gorm:"gs_language" json:"language,omitempty"`
	GSExecPath   string     `gorm:"gs_exec_path" json:"execPath,omitempty"`
	GSProtocol   string     `gorm:"gs_protocol" json:"protocol,omitempty"`
}

type VoUser struct {
	Id         int        `json:"id,omitempty"`
	UserName   string     `json:"userName,omitempty"`
	Password   string     `json:"password,omitempty"`
	CreateTime *time.Time `json:"createTime,omitempty"`
	Token      string     `json:"token,omitempty"`
}
