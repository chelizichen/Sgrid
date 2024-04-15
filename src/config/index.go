package config

import (
	"Sgrid/src/utils"
	"fmt"
	"os"

	jsonToYaml "github.com/ghodss/yaml"
)

type SgridConf struct {
	Server struct {
		Name       string                 `yaml:"name" `       // ServerName
		Host       string                 `yaml:"host"`        // Host
		Port       int                    `yaml:"port" `       // Port
		Protocol   string                 `yaml:"protoccol"`   // Protocol Http Grpc
		Language   string                 `yaml:"language"`    // Language Java Node Go
		StaticPath string                 `yaml:"staticPath" ` //
		Storage    string                 `yaml:"storage" `
		MapConf    map[string]interface{} `yaml:"mapConf"` // define Conf
	} `yaml:"server"`
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
