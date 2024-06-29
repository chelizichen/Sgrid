package rpc

import "Sgrid/src/config"

// child rpc server interface
type SgridSubServer interface {
	Registry(conf *config.SgridConf)
	NameSpace() string
	ServerPath() string
	JoinPath(args ...string) string
}
