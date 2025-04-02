package SgridPackageServer

import (
	"Sgrid/src/public"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
)

func CreateCommand(
	serverProtocol,
	serverName,
	serverLanguage,
	startDir,
	logDir,
	servantConf,
	execFilePath string,
	port int,
	processIndex int,
) (*exec.Cmd, error) {
	// print params
	fmt.Printf("CreateCommand| protocol:%s, name:%s, language:%s, startDir:%s, logDir:%s, conf:%s, execPath:%s, port:%d, index:%d\n",
	serverProtocol, serverName, serverLanguage, startDir, logDir, servantConf, execFilePath, port, processIndex)
	err := public.CheckDirectoryOrCreate(logDir)
	if err != nil{
		return nil, err
	}
	var cmd *exec.Cmd
	var startFile string // 启动文件
	if serverLanguage == public.RELEASE_NODE {
		startFile = SgridPackageInstance.JoinPath(Servants, serverName, execFilePath) // 启动文件
		cmd = exec.Command("node", startFile)
	} else if serverLanguage == public.RELEASE_JAVA || serverLanguage == public.RELEASE_JAVA_JAR {
		startFile = SgridPackageInstance.JoinPath(Servants, serverName, execFilePath) // 启动文件
		prodConf := path.Join(startDir, public.PROD_CONF_NAME)
		cmd = exec.Command("java", "-jar", startFile, fmt.Sprintf("-Dspring.config.location=file:%v", prodConf))
		cmd.Env = append(cmd.Env, fmt.Sprintf("SGRID_PROD_CONF_PATH=%v", prodConf))
	} else if serverLanguage == public.RELEASE_GO {
		startFile = SgridPackageInstance.JoinPath(Servants, serverName, execFilePath) // 启动文件
		cmd = exec.Command(startFile)
	} else if serverLanguage == public.RELEASE_EXE {
		startFile = filepath.Join(startDir, execFilePath) // 启动文件
		cmd = exec.Command(startFile)
	} else if serverLanguage == public.RELEASE_PYTHON_TAR{
		startFile = SgridPackageInstance.JoinPath(Servants, serverName, execFilePath) // 启动文件
		cmd = exec.Command(startFile)
    } else if serverLanguage == public.RELEASE_PYTHON_EXE{
		startFile = SgridPackageInstance.JoinPath(Servants, serverName, execFilePath) // 启动文件
		cmd = exec.Command(startFile)
	}else if serverLanguage == public.RELEASE_CUSTOM_COMMAND {
		var parseExecArgs []string
		err := json.Unmarshal([]byte(execFilePath), &parseExecArgs)
		if err != nil {
			return nil, err
		}
		fmt.Println("parseExecArgs", parseExecArgs)
		cmd = exec.Command(parseExecArgs[0], parseExecArgs[1:]...)
	}
	
	env := append(
		os.Environ(),
		fmt.Sprintf("%v=%v", public.ENV_TARGET_PORT, port),              // 指定端口
		fmt.Sprintf("%v=%v", public.ENV_PRODUCTION, startDir),           // 开启目录
		fmt.Sprintf("%v=%v", public.SGRID_CONFIG, servantConf),          // 配置
		fmt.Sprintf("%v=%v", public.ENV_PROCESS_INDEX, processIndex),    // 服务运行索引
		fmt.Sprintf("%v=%v", public.ENV_SGRID_SERVANT_NAME, serverName), // 服务名
		fmt.Sprintf("%v=%v", public.ENV_LOG_DIR, logDir), // 服务协议
	)
	cmd.Dir = startDir // 指定工作目录
	cmd.Env = env      // 指定环境变量
	fmt.Println("startFile", startFile)
	fmt.Println("cmd.Env", cmd.Env)
	return cmd, nil
}
