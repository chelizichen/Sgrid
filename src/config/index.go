package config

import (
	"fmt"
	"reflect"
)

var GlobalConf *SgridConf = &SgridConf{}

type server struct {
	Name     string `yaml:"name" json:"name,omitempty"`          // ServerName
	Host     string `yaml:"host" json:"host,omitempty"`          // Host
	Port     int    `yaml:"port" json:"port,omitempty"`          // Port
	Protocol string `yaml:"protoccol" json:"protocol,omitempty"` // Protocol Http Grpc
	Language string `yaml:"language" json:"language,omitempty"`  // Language Java Node Go
}

type SgridConf struct {
	Server  server                 `yaml:"server" json:"server,omitempty"`
	Conf    map[string]interface{} `yaml:"config" json:"conf,omitempty"`     // define Conf
	Servant map[string]*SgridConf  `yaml:"servant" json:"servant,omitempty"` // define Conf
}

func (s *SgridConf) Get(key string) interface{} {
	return s.Conf[key]
}

func (s *SgridConf) GetString(key string) string {
	v, ok := s.Conf[key]
	if !ok || v == nil {
		return ""
	}
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.String:
		{
			return v.(string)
		}
	default:
		{
			return ""
		}
	}
}

func (s *SgridConf) GetBool(key string) bool {
	v, ok := s.Conf[key]
	if !ok || v == nil {
		return false
	}
	fmt.Println("v", v)
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Bool:
		{
			return v.(bool)
		}
	default:
		{
			return false
		}
	}
}

func (s *SgridConf) GetStringArray(key string) []string {
	v := s.Conf[key]
	t := reflect.TypeOf(v)
	switch t.Kind() {
	case reflect.Slice:
		{
			return v.([]string)
		}
	default:
		{
			return []string{}
		}
	}
}
