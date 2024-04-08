package pojo

// 节点
type Node struct {
	Id         int
	Ip         string // IP地址
	Status     int    // 状态
	CreateTime string `gorm:"autoCreateTime"` // 创建时间
	PlatForm   string // 平台
	Main       string // 是否为主机
	UploadPath string // 上传路径
}

// 服务组
type ServantGroup struct {
	Id             int
	TagName        string // 服务标签
	TagEnglishName string // 英文
	CreateTime     string `gorm:"autoCreateTime"` // 创建时间
}

// 服务
type Servant struct {
	Id             int
	ServerName     string // 服务名称
	CreateTime     string `gorm:"autoCreateTime"` // 创建时间
	Language       string // 语言
	UpStreamName   string // 转发名称
	Location       string // 路径
	ServantGroupId int    `gorm:"foreignKey:ServantGroupId"` // 服务组ID
}

// 服务网格 用于查看所有节点信息
type Grid struct {
	Id         int
	Port       int // 网格端口
	Status     int // 网格状态
	Pid        int // 网格Pid
	UpdateTime int `gorm:"autoCreateTime"`            // 更新时间
	NodeId     int `gorm:"foreignKey:NodeId"`         // 网格所属节点ID
	ServantId  int `gorm:"foreignKey:ServantGroupId"` // 网格容纳服务ID
}

// 服务包
type ServantPackage struct {
	Id         int
	ServantId  int    // 服务Id
	Hash       string // Hash值
	FilePath   string // 文件路径
	Content    string // 上传内容
	CreateTime string `gorm:"autoCreateTime"` // 创建时间
}

type Properties struct {
	Id         int
	Key        string
	Value      string
	CreateTime string `gorm:"autoCreateTime"` // 创建时间
}
