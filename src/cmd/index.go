package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// SgridCmd 命令行工具结构体
type SgridCmd struct {
	cmd         string
	description string
	runMap      map[string]func(string)
	helpMap     map[string]string
}

// NewSgridCmd 创建新的命令行工具实例
func NewSgridCmd() *SgridCmd {
	return &SgridCmd{
		cmd:         "sgrid.cmd",
		description: "sgrid 命令行工具",
		runMap:      make(map[string]func(string)),
		helpMap:     make(map[string]string),
	}
}

// Run 运行命令
func (s *SgridCmd) Run(fnkey string, args string) {
	if fn, exists := s.runMap[fnkey]; exists {
		fn(args)
	} else {
		fmt.Printf("%s 命令不存在\n", fnkey)
	}
}

// Registry 注册命令
func (s *SgridCmd) Registry(fnkey string, fn func(string), reason string) {
	s.runMap[fnkey] = fn
	s.helpMap[fnkey] = reason
}

// Help 打印帮助信息
func (s *SgridCmd) Help() {
	fmt.Printf("%s 命令列表:\n", s.cmd)
	for key, desc := range s.helpMap {
		fmt.Printf("  %s - %s\n", key, desc)
	}
}

// Start 启动命令行工具
func (s *SgridCmd) Start() {
	fmt.Println("start")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		parts := strings.Split(input, "|")

		cmd := parts[0]
		args := ""
		if len(parts) > 1 {
			args = parts[1]
		}

		fmt.Printf("data >>> %s\n", input)
		fmt.Printf("cmd >>> %s\n", cmd)
		fmt.Printf("args >>> %s\n", args)

		s.Run(cmd, args)
	}
}

// func main() {
//     cmd := NewSgridCmd()

//     // 注册命令示例
//     cmd.Registry("hello", func(args string) {
//         fmt.Println("Hello,", args)
//     }, "打印问候语")

//     cmd.Start()
// }
