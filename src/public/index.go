package public

import (
	"Sgrid/src/config"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

const (
	ENV_PRODUCTION = "SgridProduction"
)

func SgridProduction() bool {
	s := os.Getenv(ENV_PRODUCTION)
	if len(s) == 0 {
		return false
	} else {
		return true
	}
}

func GetWd() string {
	dir, _ := os.Getwd()
	s := os.Getenv(ENV_PRODUCTION)
	if len(s) == 0 {
		return dir
	} else {
		return s
	}
}

func Join(args ...string) string {
	s := GetWd()
	arr := []string{}
	arr = append(arr, s)
	arr = append(arr, args...)
	return filepath.Join(arr...)
}

func NewConfig() (conf config.SimpConfig, err error) {
	devConfName := "simp.yaml"
	prodConfName := "simpProd.yaml"
	var path string
	if SgridProduction() {
		path = Join(prodConfName)
	} else {
		path = Join(devConfName)
	}
	// 读取 YAML 文件
	yamlFile, err := os.ReadFile(path)
	fmt.Println("Get FilePath from ", path)
	if err != nil {
		fmt.Println("Error reading YAML file:", err)
		return conf, err
	}

	// 解析 YAML 数据
	err = yaml.Unmarshal(yamlFile, &conf)
	if err != nil {
		fmt.Println("Error unmarshalling YAML:", err)
		return conf, err
	}

	// 打印解析后的配置
	fmt.Printf("%+v\n", conf)
	return conf, nil
}
