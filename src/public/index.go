package public

import (
	"Sgrid/src/config"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/c4milo/unpackit"
	"github.com/shirou/gopsutil/cpu"
	"gopkg.in/yaml.v2"
)

const (
	ENV_PRODUCTION    = "SGRID_PRODUCTION"
	ENV_TARGET_PORT   = "SGRID_TARGET_PORT"
	ENV_PROCESS_INDEX = "SGRID_PROCESS_INDEX"
	SGRID_CONFIG      = "SGRID_CONFIG"
	DEV_CONF_NAME     = "sgrid.yml"
	PROD_CONF_NAME    = "sgridProd.yml"
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

const (
	CRON_EVERY_DAY = "0 0 0 * * *"
)

const (
	LOG_TYPE_STAT         = "service-stat"
	LOG_TYPE_DATA         = "service-data"
	LOG_TYPE_ERROR        = "service-error"
	LOG_TYPE_SYSTEM_INNER = "system-inner"
)

const (
	STAT_SERVANT_COMMON = 0
	STAT_SERVANT_DELETE = -1
)

type GRPC_CLIENT_PROXYS struct{}

const (
	ChunkFileSize = 1024 * 1024
)

func SgridProduction() bool {
	s := os.Getenv(ENV_PRODUCTION)
	fmt.Println("s", s)
	if len(s) == 0 {
		return false
	} else {
		return true
	}
}

func GetCurrTime() string {
	return time.Now().Format(time.DateTime)
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

func GetCpuPercent() string {
	cpuPercentage, err := cpu.Percent(time.Second, false)
	if err != nil {
		fmt.Println("Error getting CPU usage:", err)
		return "0.00%"
	}
	percent := fmt.Sprintf("%.2f", cpuPercentage[0])
	return percent
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

func WithNameSpace(nameSpace string) ConfOpt {
	return func(conf *withConf) {
		conf.nameSpace = nameSpace
	}
}

type withConf struct {
	targetPath string
	nameSpace  string
}

func NewConfig(opts ...ConfOpt) (conf *config.SgridConf, err error) {
	prod := os.Getenv(SGRID_CONFIG)
	wc := &withConf{}
	if len(prod) > 0 {
		err = yaml.Unmarshal([]byte(prod), &conf)
		if err != nil {
			fmt.Println("err", err.Error())
		}
		fmt.Printf("SGRID_PROD_CONFIG %+v\n", conf)
	} else {
		for _, opt := range opts {
			opt(wc)
		}
		var path string
		if len(wc.targetPath) != 0 {
			path = wc.targetPath
		} else if len(wc.nameSpace) != 0 {
			path = Join(DEV_CONF_NAME)
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
		fmt.Printf("SGRID_DEV_CONFIG %+v\n", conf)

	}
	// 打印解析后的配置
	return conf, nil
}

func CheckDirectoryOrCreate(directoryPath string) error {
	fmt.Println("directoryPath", directoryPath)
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

// 暂时不可用
func ThreadLock() bool {
	if len(os.Getenv(ENV_PRODUCTION)) > 0 {
		return os.Getenv(ENV_PROCESS_INDEX) == "1"
	}
	return false
}

func Removenullvalue(slice []interface{}) []interface{} {
	var output []interface{}
	for _, element := range slice {
		if element != nil { //if condition satisfies add the elements in new slice
			output = append(output, element)
		}
	}
	return output //slice with no nil-values
}
