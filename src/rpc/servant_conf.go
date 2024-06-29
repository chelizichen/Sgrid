package rpc

type SgridGrpcServantConfig struct {
	ServiceName string `json:"service_name"` // 服务名
	Protocol    string `json:"protocol"`     // 协议
	Host        string `json:"host"`         // 主机
	Port        int    `json:"port"`         // 端口
}
