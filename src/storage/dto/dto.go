package dto

import "time"

type PageBasicReq struct {
	Size    int    `json:"size,omitempty"`
	Offset  int    `json:"offset,omitempty"`
	Keyword string `json:"keyword,omitempty"`
	Id      int    `json:"id,omitempty"`
}

type QueryPackageDto struct {
	Size    int    `json:"size,omitempty"`
	Offset  int    `json:"offset,omitempty"`
	Version string `json:"version,omitempty"`
	Id      int    `json:"id,omitempty"`
}

type SaveServantGroupDto struct {
	Id             int        `json:"id,omitempty"`
	TagName        string     `json:"tagName,omitempty"`        // 服务标签
	TagEnglishName string     `json:"tagEnglishName,omitempty"` // 英文
	CreateTime     *time.Time `json:"createTime,omitempty"`     // 创建时间
}

type SaveServantDto struct {
	Id             int    `json:"id,omitempty"`
	ServerName     string `json:"serverName,omitempty"`     // 服务名称
	CreateTime     string `json:"createTime,omitempty"`     // 创建时间
	Language       string `json:"language,omitempty"`       // 语言
	UpStreamName   string `json:"upStreamName,omitempty"`   // 转发名称
	Location       string `json:"location,omitempty"`       // 路径
	Protocol       string `json:"protocol,omitempty"`       // 协议
	ExecPath       string `json:"execPath,omitempty"`       // 可执行路径
	ServantGroupId int    `json:"servantGroupId,omitempty"` // 服务组ID
}

type GridDTO struct {
	Id         int        `json:"id,omitempty"`
	Port       int        `json:"port,omitempty"`        // 网格端口
	Status     int        `json:"status,omitempty"`      // 网格状态
	Pid        int        `json:"pid,omitempty"`         // 网格Pid
	UpdateTime *time.Time `json:"update_time,omitempty"` // 更新时间
	NodeId     int        `json:"nodeId,omitempty"`      // 网格所属节点ID
	ServantId  int        `json:"servantId,omitempty"`   // 网格容纳服务ID
}
