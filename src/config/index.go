package config

import (
	"Sgrid/src/utils"
	"fmt"
	"os"
	"reflect"

	jsonToYaml "github.com/ghodss/yaml"
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

func ResetConfig(yamlContent string, filePath string) error {
	err := os.WriteFile(filePath, []byte(yamlContent), 0644)
	if err != nil {
		return fmt.Errorf("error writing YAML to file: %v", err)
	}
	return nil
}

func CoverConfig(content string, filePath string) error {
	utils.IFExistThenRemove(filePath, false)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func ParseConfig(yamlString string) (string, error) {
	yml, err := jsonToYaml.JSONToYAML([]byte(yamlString))
	if err != nil {
		fmt.Println("JSON TO YamlError")
	}
	fmt.Println("Cover yml \n", string(yml))

	return string(yml), nil
}
