package servant

import "Sgrid/src/config"

type SgridRegistryServiceInf interface {
	Registry(conf *config.SgridConf)
	NameSpace() string
	ServerPath() string
	JoinPath(args ...string) string
}
