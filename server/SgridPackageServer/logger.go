package SgridPackageServer

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	HEAD = 1
	TAIL = 2
)

// tail -50000 ${logFile}|tail -${len} | iconv -c -f UTF-8 -t UTF-8|sed 's/[\cA-\cZ]//g'
// tail -50000 ${logFile} | grep -a ${keyword}|tail -${len} | iconv -c -f UTF-8 -t UTF-8|sed 's/[\cA-\cZ]//g'
func SearchLog(logFile string, logType uint32, keyword string, len uint32) ([]string, error) {
	log_type := ""
	log_cmd := ""
	if logType == HEAD {
		log_type = "head"
	}
	if logType == TAIL {
		log_type = "tail"
	}
	var cmd *exec.Cmd
	if keyword == "" {
		log_cmd = fmt.Sprintf("%s -500000 %s|%s -%d | iconv -c -f UTF-8 -t UTF-8|sed 's/[cA-cZ]//g'", log_type, logFile, log_type, len)
		// 如果没有提供关键词，只截取文件末尾的内容
		cmd = exec.Command("sh", "-c", log_cmd)
	} else {
		log_cmd = fmt.Sprintf("%s -500000 %s |%s -a %s|tail -%d | iconv -c -f UTF-8 -t UTF-8|sed 's/[\\cA-\\cZ]//g'", log_type, logFile, "grep", keyword, len)
		// 如果提供了关键词，先截取文件末尾内容，再筛选包含关键词的行
		cmd = exec.Command("sh", "-c", log_cmd)
	}

	// 创建一个字节缓冲区来存储命令执行的输出
	var out bytes.Buffer
	cmd.Stdout = &out

	// 执行命令
	err := cmd.Run()
	if err != nil {
		fmt.Printf("执行命令时出错: %v\n", err)
		return nil, err
	}
	
	// 将输出按行分割并过滤空行
	output := out.String()
	lines := strings.Split(output, "\n")
	result := make([]string, 0)
	
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			result = append(result, line)
		}
	}
	
	fmt.Println(log_cmd)
	return result, nil
}

func GetLogFileList(logDir string) ([]string, error) {
	var logFiles []string
    fmt.Println("logDir >> ", logDir)
    // 检查目录是否存在
    if _, err := os.Stat(logDir); os.IsNotExist(err) {
        return nil, fmt.Errorf("目录不存在: %s", logDir)
    }
    
    // 遍历目录
    err := filepath.Walk(logDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        // 检查是否为文件且后缀为.log
        if !info.IsDir() && filepath.Ext(path) == ".log" {
            logFiles = append(logFiles, path)
        }
        return nil
    })
    
    if err != nil {
        return nil, fmt.Errorf("遍历目录失败: %v", err)
    }
    
    return logFiles, nil
}

// head -500000 /Users/leemulus/Desktop/临时/DXZQ.TgH5WebServer_副本.log |head -a '东兴金蟾'|tail -100 | iconv -c -f UTF-8 -t UTF-8|sed 's/[\cA-\cZ]//g'
// func test() {
// 	cwd, _ := os.Getwd()
// 	fp := filepath.Join(cwd, "DXZQ.TgH5WebServer_副本.log")
// 	logRsp, err := searchLog(fp, HEAD, "2025-01-15 15:12:53", 100)
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	fmt.Println("logRsp >> ", logRsp)
// }
