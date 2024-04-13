package public

import (
	"Sgrid/src/config"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/c4milo/unpackit"
	"gopkg.in/yaml.v2"
)

const (
	ENV_PRODUCTION  = "SGRID_PRODUCTION"
	ENV_TARGET_PORT = "SGRID_TARGET_PORT"
	DEV_CONF_NAME   = "sgrid.yml"
	PROD_CONF_NAME  = "sgridProd.yml"
)

const (
	RELEASE_GO   = "go"
	RELEASE_NODE = "node"
	RELEASE_JAVA = "java"
)

const (
	PROTOCOL_HTTP = "http"
	PROTOCOL_GRPC = "grpc"
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

type ConfOpt func(*withConf)

func WithTargetPath(targetPath string) ConfOpt {
	return func(conf *withConf) {
		conf.targetPath = targetPath
	}
}

type withConf struct {
	targetPath string
}

func NewConfig(opts ...ConfOpt) (conf *config.SgridConf, err error) {
	wc := &withConf{}
	for _, opt := range opts {
		opt(wc)
	}
	var path string
	if len(wc.targetPath) != 0 {
		path = wc.targetPath
	} else if SgridProduction() {
		path = Join(PROD_CONF_NAME)
	} else {
		path = Join(DEV_CONF_NAME)
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

func CheckDirectoryOrCreate(directoryPath string) error {
	_, err := os.Stat(directoryPath)
	if os.IsNotExist(err) {
		err := os.Mkdir(directoryPath, os.ModePerm)
		if err != nil {
			fmt.Println("err", err)
		}
		fmt.Printf("Path %s does not exist.\n", directoryPath)
		return err
	}
	return err
}

func Tar2Dest(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		fmt.Println("Open Error", err.Error())
		return err
	}
	defer file.Close()
	err = unpackit.Unpack(file, dest)
	if err != nil {
		fmt.Println("Unpackit Error", err.Error())
		return err
	}
	return nil
}

func CopyProdYml(storageYmlEPath, storageYmlProdPath string) (err error) {
	_, err = os.Stat(storageYmlProdPath)
	if err != nil {
		fmt.Println("os.Stat ", err.Error())
	}
	if os.IsNotExist(err) {
		err = CopyFile(storageYmlEPath, storageYmlProdPath)
		if err != nil {
			fmt.Println("utils.CopyFile ", storageYmlEPath, err.Error())
		}
	}
	return err
}

func CopyFile(src string, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}
	return nil
}

func IsExist(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	} else {
		return true
	}
}
