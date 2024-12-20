package dto

import (
	"Sgrid/src/storage/pojo"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

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
	UserId         int        `json:"userId,omitempty"`         // 用户ID
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
	UserId         int    `json:"userId,omitempty"`         // 用户ID
	Preview        string `json:"preview,omitempty"`        // 预览地址
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

// 节点
type NodeDTO struct {
	Id         int        `json:"id,omitempty"`
	Ip         string     `json:"ip,omitempty"`         // IP地址
	Status     int        `json:"nodeStatus,omitempty"` // 状态
	CreateTime *time.Time `json:"createTime,omitempty"` // 创建时间
	PlatForm   string     `json:"platForm,omitempty"`   // 平台
	Main       string     `json:"main,omitempty"`       // 是否为主机
	UploadPath string     `json:"uploadPath,omitempty"` // 上传路径
}

type TraceLogDto struct {
	Keyword    string `json:"keyword,omitempty"`
	Offset     int    `json:"offset,omitempty"`
	Size       int    `json:"size,omitempty"`
	SearchTime string `json:"searchTime,omitempty"`
	pojo.TraceLog
}

func NewPageBaiscReq(ctx *gin.Context) *PageBasicReq {
	size, _ := strconv.Atoi(ctx.DefaultQuery("size", "10"))
	offset, _ := strconv.Atoi(ctx.DefaultQuery("offset", "0"))
	id, _ := strconv.Atoi(ctx.DefaultQuery("id", "0"))

	keyword := ctx.DefaultQuery("keyword", "")
	return &PageBasicReq{
		Size:    size,
		Offset:  offset,
		Keyword: keyword,
		Id:      id,
	}
}

type LogTypeDto struct {
	DateTime string `json:"dateTime,omitempty"`
	LogType  string `json:"logType,omitempty"`
}
