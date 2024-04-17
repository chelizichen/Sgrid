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
	Name     string `yaml:"name" `     // ServerName
	Host     string `yaml:"host"`      // Host
	Port     int    `yaml:"port" `     // Port
	Protocol string `yaml:"protoccol"` // Protocol Http Grpc
	Language string `yaml:"language"`  // Language Java Node Go
}

type SgridConf struct {
	Server  server                 `yaml:"server"`
	Conf    map[string]interface{} `yaml:"config"`  // define Conf
	Servant map[string]*SgridConf  `yaml:"servant"` // define Conf
}

func (s *SgridConf) Get(key string) interface{} {
	return s.Conf[key]
}

func (s *SgridConf) GetString(key string) string {
	v := s.Conf[key]
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
