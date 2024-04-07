package pojo

// 节点
type Node struct {
	Id         int
	Ip         string // IP地址
	Status     string // 状态
	CreateTime string `gorm:"autoCreateTime"` // 创建时间
	PlatForm   string // 平台
	Main       string // 是否为主机
}

// 服务组
type ServantGroup struct {
	Id         int
	TagName    string // 服务标签
	CreateTime string `gorm:"autoCreateTime"` // 创建时间
}

// 服务
type Servant struct {
	Id             int
	ServerName     string // 服务名称
	CreateTime     string `gorm:"autoCreateTime"` // 创建时间
	Language       string // 语言
	ServantGroupId int    // 服务组ID
}

// 服务网格 用于查看所有节点信息
type Grid struct {
	Id         int
	ServantId  int // 网格容纳服务ID
	Port       int // 网格端口
	NodeId     int // 网格所属节点ID
	Status     int // 网格状态
	Pid        int // 网格Pid
	UpdateTime int `gorm:"autoCreateTime"` // 更新时间
}

// 扩容服务
type ExpansionGrid struct {
	Id           int
	ServantId    int    // 服务Id
	Location     string // Nginx Location
	ProxyPass    string // 转发地址
	UpStreamName string // Nginx UpStreamName
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
