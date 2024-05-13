package public

import (
	"Sgrid/src/config"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/c4milo/unpackit"
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

func GetLogger(filePath string, pattern string, rows int) (string, error) {
	output, err := TailAndGrep(filePath, rows, pattern)
	if err != nil {
		fmt.Println("执行命令失败:", err)
		return output, err
	}
	return output, err
}

func TailAndGrep(filename string, n int, pattern string) (string, error) {
	// 构造命令
	cmdTail := exec.Command("tail", fmt.Sprintf("-n%d", n), filename)
	cmdGrep := exec.Command("grep", pattern)

	// 创建管道
	r, w := io.Pipe()
	defer r.Close()

	// 将 tail 的输出连接到 grep 的输入
	cmdTail.Stdout = w
	cmdGrep.Stdin = r

	// 创建缓冲区用于存储 grep 的输出
	var output bytes.Buffer
	cmdGrep.Stdout = &output

	// 启动命令
	errTail := cmdTail.Start()
	if errTail != nil {
		return "", errTail
	}

	errGrep := cmdGrep.Start()
	if errGrep != nil {
		return "", errGrep
	}

	// 等待命令执行完成
	errTailWait := cmdTail.Wait()
	if errTailWait != nil {
		return "", errTailWait
	}

	w.Close()

	errGrepWait := cmdGrep.Wait()
	if errGrepWait != nil {
		return "", errGrepWait
	}

	return output.String(), nil
}
