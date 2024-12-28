package SgridPackageServer

import (
	logProto "Sgrid/server/SgridLogTraceServer/proto"
	protocol "Sgrid/server/SgridPackageServer/proto"
	"Sgrid/src/config"
	"Sgrid/src/pool"
	"Sgrid/src/public"
	"Sgrid/src/rpc"
	"Sgrid/src/storage"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/pojo"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"

	pk "Sgrid/src/public/process"

	"github.com/panjf2000/ants"
	"github.com/robfig/cron"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	p "github.com/shirou/gopsutil/process"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

var globalConf *config.SgridConf
var AntsPool *ants.Pool
var globalGrids sync.Map

var SgridPackageInstance = &SgridPackage{}

func getStat(pid int) *p.Process {
	process, err := p.NewProcess(int32(pid))
	if err != nil {
		fmt.Println("Error creating new process:", err)
		return nil
	}
	return process
}

type WithSgridMonitorConfFunc func(*SgridMonitor)

type SgridMonitor struct {
	interval     time.Duration // 上报时间
	cmd          *exec.Cmd
	ctx          context.Context
	cancel       context.CancelFunc
	serverName   string
	gridId       int
	cronInstance *cron.Cron
	port         int
	stdin        io.WriteCloser // 添加stdin字段
}

func WithMonitorInterval(interval time.Duration) func(*SgridMonitor) {
	return func(monitor *SgridMonitor) {
		if interval.Seconds() < CONSTANT_MONITOR_INTERVAL { // min 5s
			interval = time.Second * CONSTANT_MONITOR_INTERVAL
		}
		monitor.interval = interval
	}
}

func WithMonitorSetCmd(cmd *exec.Cmd) func(*SgridMonitor) {
	return func(monitor *SgridMonitor) {
		monitor.cmd = cmd
	}
}

func WithMonitorGridIDAndPort(id, port int) func(*SgridMonitor) {
	return func(monitor *SgridMonitor) {
		monitor.gridId = id
		monitor.port = port
	}
}

func WithMonitorServerName(serverName string) func(*SgridMonitor) {
	return func(monitor *SgridMonitor) {
		monitor.serverName = serverName
	}
}

func NewSgridMonitor(opt ...WithSgridMonitorConfFunc) *SgridMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	monitor := &SgridMonitor{
		ctx:    ctx,
		cancel: cancel,
	}
	for _, v := range opt {
		fn := v
		fn(monitor)
	}
	if monitor.cmd != nil {
		stdin, err := monitor.cmd.StdinPipe()
		if err != nil {
			fmt.Println("Failed to get StdinPipe:", err)
			return nil
		}
		monitor.stdin = stdin
	}
	return monitor
}

func (s *SgridMonitor) Report() {
	s.cronInstance = cron.New()
	isNeedCheckPortToPid := true
	var job = func() {
		select {
		case <-s.ctx.Done():
			fmt.Printf("%s [%d] recieved done singal by ctx, return Report \n", s.serverName, s.gridId)
			return
		default:
			AntsPool.Submit(func() {
				if isNeedCheckPortToPid {
					err := pk.SgridProcessUtil.ValidatePortToPid(s.port, s.getPid())
					if err != nil {
						storage.PushErr(&pojo.SystemErr{
							Type: "system/error/SgridPackageServer/pk.SgridProcessUtil.ValidatePortToPid",
							Info: err.Error(),
						})
						s.kill()
						return
					}
					isNeedCheckPortToPid = false
				}
				id := s.getPid()
				statInfo := getStat(id)
				if statInfo == nil {
					return
				}
				status, _ := statInfo.Status()
				var gridStat int = 0
				if status == "Z" { // state is down , need kill
					s.kill()
					gridStat = 0
				} else {
					gridStat = 1
				}
				cpu, _ := statInfo.CPUPercent()
				threads, _ := statInfo.NumThreads()
				name, _ := statInfo.Name()
				isRuning, _ := statInfo.IsRunning()
				running := "not run"
				if isRuning {
					running = "running"
				}
				mis, _ := statInfo.MemoryInfo()
				stack := mis.Stack
				MemoryData := mis.Data
				now := time.Now()
				content := fmt.Sprintf("pid:%v | cpu:%v | thread:%v | status:%v \n", s.getPid(), cpu, threads, status)
				storage.UpdateGrid(&pojo.Grid{
					Id:         s.gridId,
					Status:     gridStat,
					Pid:        s.getPid(),
					UpdateTime: &now,
				})
				storage.SaveStatLog(&pojo.StatLog{
					GridId:      s.gridId,
					Stat:        BEHAVIOR_ALIVE,
					Pid:         s.getPid(),
					CPU:         cpu,
					Threads:     threads,
					Name:        name,
					IsRunning:   running,
					MemoryStack: stack,
					MemoryData:  MemoryData,
				})
				var logReq = &logProto.LogTraceReq{
					LogServerName: s.serverName,
					LogHost:       globalConf.Server.Host,
					LogGridId:     int64(s.gridId),
					LogType:       public.LOG_TYPE_STAT,
					LogContent:    content,
					LogBytesLen:   int64(len(content)),
					CreateTime:    public.GetCurrTime(),
				}
				var rpl logProto.BasicResp
				NewSgridLogTraceServant.Request(rpc.RequestPack{
					Method: "LogTrace",
					Body:   logReq,
					Reply:  &rpl,
				})
			})
		}
	}
	s.cronInstance.AddFunc("@every 30s", job)
	s.cronInstance.Start()
}

func (s *SgridMonitor) PrintLogger() {
	op, err := s.cmd.StdoutPipe()
	if err != nil {
		storage.PushErr(&pojo.SystemErr{
			Type: "system/error/SgridPackageServer/s.cmd.StdoutPipe()",
			Info: err.Error(),
		})
	}
	ep, err := s.cmd.StderrPipe()
	if err != nil {
		storage.PushErr(&pojo.SystemErr{
			Type: "system/error/SgridPackageServer/s.cmd.StderrPipe()",
			Info: err.Error(),
		})
	}
	go func() {
		select {
		case <-s.ctx.Done():
			fmt.Printf("%s [%d] recieved done singal by ctx, return PrintLogger.Output \n", s.serverName, s.gridId)
			return
		default:
			buf := make([]byte, 1024)
			n, err := op.Read(buf)
			if err != nil && !errors.Is(err, io.EOF) {
				fmt.Println("system/error/SgridPackageServer/op.Read.error", err.Error())
				storage.PushErr(&pojo.SystemErr{
					Type: "system/error/SgridPackageServer/op.Read(buf)",
					Info: err.Error(),
				})
				break
			} else {
				// 打印输出
				content := string(buf[:n])
				var logReq = &logProto.LogTraceReq{
					LogServerName: s.serverName,
					LogHost:       globalConf.Server.Host,
					LogGridId:     int64(s.gridId),
					LogType:       public.LOG_TYPE_DATA,
					LogContent:    content,
					LogBytesLen:   int64(n),
					CreateTime:    public.GetCurrTime(),
				}
				fmt.Println("SgridPackageServer.read.data.content,", content)
				var rpl logProto.BasicResp
				NewSgridLogTraceServant.Request(rpc.RequestPack{
					Method: "LogTrace",
					Body:   logReq,
					Reply:  &rpl,
				})
			}
		}
	}()
	go func() {
		select {
		case <-s.ctx.Done():
			fmt.Printf("%s [%d] recieved done singal by ctx, return PrintLogger.ErrorPut \n", s.serverName, s.gridId)
			return
		default:
			buf := make([]byte, 1024)
			now := public.GetCurrTime()
			n, err := ep.Read(buf)
			if err != nil && !errors.Is(err, io.EOF) {
				fmt.Println("system/error/SgridPackageServer/ep.Read.error", err.Error())
				storage.PushErr(&pojo.SystemErr{
					Type: "system/error/SgridPackageServer/ep.Read(buf)",
					Info: err.Error(),
				})
				break
			} else {
				// 打印输出
				content := string(buf[:n])
				var logReq = &logProto.LogTraceReq{
					LogServerName: s.serverName,
					LogHost:       globalConf.Server.Host,
					LogGridId:     int64(s.gridId),
					LogType:       public.LOG_TYPE_ERROR,
					LogContent:    content,
					LogBytesLen:   int64(n),
					CreateTime:    now,
				}
				fmt.Println("SgridPackageServer.read.error.content,", content)
				var rpl logProto.BasicResp
				NewSgridLogTraceServant.Request(rpc.RequestPack{
					Method: "LogTrace",
					Body:   logReq,
					Reply:  &rpl,
				})
			}
		}
	}()
}

func (s *SgridMonitor) kill() {
	s.cancel()
	s.cronInstance.Stop()
	fmt.Println("system/log/server.kill |", s.serverName)
	storage.SaveStatLog(&pojo.StatLog{
		GridId: s.gridId,
		Stat:   BEHAVIOR_KILL,
		Pid:    s.getPid(),
	})
	err := s.cmd.Process.Kill()
	if err != nil {
		fmt.Println("system/err/process.kill | ", err.Error())
	}
}

func (s *SgridMonitor) getPid() int {
	if s.cmd != nil {
		if s.cmd.Process != nil {
			return s.cmd.Process.Pid
		}
	}
	return 0
}

type fileTransferServer struct {
	protocol.UnimplementedFileTransferServiceServer
}

func (s *fileTransferServer) StreamFile(stream protocol.FileTransferService_StreamFileServer) error {
	// 创建文件来存储接收到的数据
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return errors.New("无法获取元数据")
	}
	// 获取文件名和哈希值
	filename := md.Get("filename")[0]
	servername := md.Get("serverName")[0]
	directoryPath := SgridPackageInstance.JoinPath(App, servername)
	err := public.CheckDirectoryOrCreate(directoryPath)
	if err != nil {
		fmt.Println("check directory error")
	}
	targetFilePath := SgridPackageInstance.JoinPath(App, servername, filename)
	file, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 循环接收文件块，直到流结束
	for {
		fileChunk, err := stream.Recv()
		if err != nil {
			fmt.Println("err", err.Error())
			if err != io.EOF {
				fmt.Println("error :: server-side :: EOF")
				break
			} else {
				message := fmt.Sprintf("error :: server-side:: %s", err.Error())
				finalResponse := &protocol.FileResp{
					Msg:  message,
					Code: 200,
				}
				if err := stream.RecvMsg(finalResponse); err != nil {
					return err
				}
				// 流结束，退出循环
				return nil
			}

		}
		fmt.Println("chunk", fileChunk.Offset)
		if err != nil {
			return err
		}

		// 写入文件块到文件
		_, err = file.Write(fileChunk.Data)
		if err != nil {
			return err
		}
		if public.ChunkFileSize == fileChunk.Offset {
			// 发送文件接收完成的响应
			finalResponse := &protocol.FileResp{
				Msg:  "Chunk received successfully",
				Code: 100,
			}
			if err := stream.Send(finalResponse); err != nil {
				return err
			}
		} else {
			break
		}

	}

	// 发送文件接收完成的响应
	finalResponse := &protocol.FileResp{
		Msg:  "File received successfully",
		Code: 200,
	}
	if err := stream.Send(finalResponse); err != nil {
		return err
	}

	return nil
}

func (s *fileTransferServer) ShutdownGrid(ctx context.Context, req *protocol.ShutdownGridReq) (res *protocol.BasicResp, err error) {
	for _, _grid := range req.GetReq() {
		grid := _grid
		h := grid.GetHost()
		if globalConf.Server.Host != h {
			continue
		}
		i := grid.GetGridId()
		sm, ok := globalGrids.Load(int(i))
		if ok && sm != nil {
			fmt.Printf("Common Down [%d]", i)
			storage.SaveStatLog(&pojo.StatLog{
				GridId: int(i),
				Stat:   BEHAVIOR_DOWN,
				Pid:    sm.(*SgridMonitor).getPid(),
			})
			fmt.Printf("Killd with gridId[%d] \n", i)
			sm.(*SgridMonitor).kill()
			globalGrids.Delete(int(i))
			continue
		} else {
			fmt.Printf("Force Down [%d]", i)
			pk.SgridProcessUtil.QueryProcessPidThenKill(int(grid.GetPort())) // 有可能需要强制down
			continue
		}
	}
	return &protocol.BasicResp{
		Code:    0,
		Message: "ok",
	}, nil
}

// 发布 -> 上报给主控
func (s *fileTransferServer) ReleaseServerByPackage(ctx context.Context, req *protocol.ReleaseServerReq) (res *protocol.BasicResp, err error) {
	// 异常处理
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("Recovered from panic: %v", r)
			storage.PushErr(&pojo.SystemErr{
				Type: "system/error/SgridPackageServer.ReleaseServerByPackage.recover",
				Info: "recover Error :" + err.Error(),
			})
		}
	}()
	if len(req.ServantGrids) == 0 {
		return &protocol.BasicResp{
			Code:    -2,
			Message: "req.ServantGrids is empty",
		}, err
	}
	fmt.Println("ReleaseServerByPackage Req ||", string(req.String()))      // 日志打印
	filePath := req.FilePath                                                // 服务路径
	serverLanguage := req.ServerLanguage                                    // 服务语言
	serverName := req.ServerName                                            // 服务名称
	serverProtocol := req.ServerProtocol                                    // 协议
	execFilePath := req.ExecPath                                            // 执行路径
	startDir := SgridPackageInstance.JoinPath(Servants, serverName)         // 解压目录
	packageFile := SgridPackageInstance.JoinPath(App, serverName, filePath) // 路径
	servantId := req.ServantId                                              // 服务ID
	err = public.CheckDirectoryOrCreate(startDir)                           // 检查并创建目录
	if err != nil {
		return &protocol.BasicResp{
			Code:    -1,
			Message: fmt.Sprintf("public.CheckDirectoryOrCreate.error %v", err.Error()),
		}, err
	}
	if req.ServerLanguage != public.RELEASE_EXE && req.ServerLanguage != public.RELEASE_JAVA_JAR { // 是 jar 或者 exe 时，不用 解压包
		err = public.Tar2Dest(packageFile, startDir) // 解压
		if err != nil {
			return &protocol.BasicResp{
				Code:    -1,
				Message: fmt.Sprintf("ReleaseServerByPackage.Tar2Dest.error %v", err.Error()),
			}, err
		}
	} else {
		err = public.CopyFile(packageFile, filepath.Join(startDir, req.ExecPath))
		if err != nil {
			var errMsg string = fmt.Sprintf("ReleaseServerByPackage.CopyFile.error %v", err.Error())
			fmt.Println("errMsg", errMsg)
			return &protocol.BasicResp{
				Code:    -1,
				Message: errMsg,
			}, err
		}
	}
	servantGrid := req.ServantGrids // 服务列表  通过Host过滤拿到IP，然后进行服务启动
	var startFile string            // 启动文件
	servantConf := storage.GetServantConfById(int(servantId)).Conf
	for processIndex, grid := range servantGrid { // 通过Host过滤拿到IP，然后进行服务启动
		ProcessIndex := processIndex
		GRID := grid
		id := GRID.GridId
		fmt.Println("GRID.IP", GRID.Ip)
		fmt.Println("globalConf.Server.Host", globalConf.Server.Host)
		// Host 确保与主控配置文件一致
		if GRID.Ip != globalConf.Server.Host {
			fmt.Println("server is not equal")
			continue
		}
		item, ok := globalGrids.Load(int(id))
		if ok && item != nil { // 终止
			item.(*SgridMonitor).kill()
			globalGrids.Delete(int(id))
		}
		if err != nil {
			fmt.Println("error", err.Error())
		}
		var cmd *exec.Cmd
		if serverProtocol == public.PROTOCOL_GRPC {
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
				startFile = filepath.Join(startDir, req.ExecPath) // 启动文件
				cmd = exec.Command(startFile)
			} else if serverLanguage == public.RELEASE_CUSTOM_COMMAND {
				var parseExecArgs []string
				err = json.Unmarshal([]byte(execFilePath), &parseExecArgs)
				if err != nil {
					return &protocol.BasicResp{
						Code:    -1,
						Message: fmt.Sprintf("ReleaseServerByPackage.json.Unmarshal([]byte(execFilePath), &parseExecArgs).error %v", err.Error()),
					}, err
				}
				fmt.Println("parseExecArgs", parseExecArgs)
				cmd = exec.Command(parseExecArgs[0], parseExecArgs[1:]...)
			}
		}

		if serverProtocol == public.PROTOCOL_HTTP {
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
				startFile = filepath.Join(startDir, req.ExecPath) // 启动文件
				cmd = exec.Command(startFile)
			} else if serverLanguage == public.RELEASE_CUSTOM_COMMAND {
				var parseExecArgs []string
				err = json.Unmarshal([]byte(execFilePath), &parseExecArgs)
				if err != nil {
					return &protocol.BasicResp{
						Code:    -1,
						Message: fmt.Sprintf("ReleaseServerByPackage.json.Unmarshal([]byte(execFilePath), &parseExecArgs).error %v", err.Error()),
					}, err
				}
				fmt.Println("parseExecArgs", parseExecArgs)
				cmd = exec.Command(parseExecArgs[0], parseExecArgs[1:]...)
			}
		}
		env := append(os.Environ(),
			fmt.Sprintf("%v=%v", public.ENV_TARGET_PORT, grid.Port),      // 指定端口
			fmt.Sprintf("%v=%v", public.ENV_PRODUCTION, startDir),        // 开启目录
			fmt.Sprintf("%v=%v", public.SGRID_CONFIG, servantConf),       // 配置
			fmt.Sprintf("%v=%v", public.ENV_PROCESS_INDEX, ProcessIndex), // 服务运行索引
		)
		cmd.Dir = startDir // 指定工作目录
		cmd.Env = env      // 指定环境变量
		fmt.Println("startFile", startFile)
		fmt.Println("cmd.Env", cmd.Env)

		monitor := NewSgridMonitor(
			WithMonitorInterval(time.Second*5),
			WithMonitorSetCmd(cmd),
			WithMonitorServerName(serverName),
			WithMonitorGridIDAndPort(int(id), int(grid.Port)),
		)
		globalGrids.LoadOrStore(int(id), monitor)
		monitor.PrintLogger()
		err = cmd.Start()
		fmt.Println("*************服务启动**************")
		if err != nil {
			storage.PushErr(&pojo.SystemErr{
				Type: "system/error/SgridPackageServer/cmd.Start()",
				Info: err.Error(),
			})
			continue
		}
		storage.SaveStatLog(&pojo.StatLog{
			GridId: int(id),
			Stat:   BEHAVIOR_PULL,
			// Pid:    monitor.getPid(),
		})

		fmt.Println("*************开始日志上报**************")
		go monitor.Report()
	}

	//  ********************** debug *********************
	fmt.Println("startFile", startFile)
	fmt.Println("packageFile", packageFile)
	//  ********************** debug *********************

	return &protocol.BasicResp{
		Code:    0,
		Message: "ok",
	}, nil
}

func (s *fileTransferServer) GetLogFileByHost(ctx context.Context, in *protocol.GetLogFileByHostReq) (*protocol.GetLogFileByHostResp, error) {
	files := storage.GetTraceLogFiles(int(in.GridId), in.ServerName)
	var respList []*protocol.GetLogFileByHostVo
	for _, v := range files {
		t, _ := time.Parse(Layout, v.LogTime)
		respList = append(respList, &protocol.GetLogFileByHostVo{
			LogType:  v.LogType,
			DateTime: t.Format(time.DateOnly),
		})
	}
	return &protocol.GetLogFileByHostResp{
		Data:    respList,
		Code:    0,
		Message: "ok",
	}, nil
}

func (s *fileTransferServer) GetLogByFile(ctx context.Context, in *protocol.GetLogByFileReq) (*protocol.GetLogByFileResp, error) {
	var req = &dto.TraceLogDto{
		Keyword:    in.Pattern,
		Size:       int(in.Size),
		SearchTime: in.DateTime,
		Offset:     int(in.Offset),
		TraceLog: pojo.TraceLog{
			LogHost:       in.Host,
			LogServerName: in.ServerName,
			LogGridId:     int64(in.GridId),
			LogType:       in.LogType,
		},
	}
	Data, total, err := storage.GetTraceLog(req)
	if err != nil {
		return nil, err
	}
	return &protocol.GetLogByFileResp{
		Data:  Data,
		Total: total,
	}, nil
}

func (s *fileTransferServer) GetPidInfo(ctx context.Context, in *protocol.GetPidInfoReq) (resp *protocol.GetPidInfoResp, err error) {
	if len(in.HostPids) == 0 {
		return nil, nil
	}
	ret := &protocol.GetPidInfoResp{
		Data: []*protocol.HostPidInfo{},
	}
	fmt.Println("in", in)
	for _, v := range in.HostPids {
		if v.Host != globalConf.Server.Host {
			continue
		}
		statInfo := getStat(int(v.Pid))
		_, ok := globalGrids.Load(int(v.GridId))
		if !ok {
			now := time.Now()
			storage.UpdateGrid(&pojo.Grid{
				Id:         int(v.GridId),
				Status:     0,
				Pid:        0,
				UpdateTime: &now,
			})
			storage.SaveStatLog(&pojo.StatLog{
				GridId:      int(v.GridId),
				Stat:        BEHAVIOR_CHECK,
				Pid:         0,
				CPU:         0,
				Threads:     0,
				Name:        "",
				IsRunning:   "find grid error, code[9329] ",
				MemoryStack: 0,
				MemoryData:  0,
			})
			continue
		}
		fmt.Println("int(v.Pid)", int(v.Pid))
		if statInfo == nil {
			fmt.Println("error process", v.Pid)
			now := time.Now()
			storage.UpdateGrid(&pojo.Grid{
				Id:         int(v.GridId),
				Status:     0,
				Pid:        int(v.Pid),
				UpdateTime: &now,
			})
			storage.SaveStatLog(&pojo.StatLog{
				GridId: int(v.GridId),
				Stat:   BEHAVIOR_CHECK,
				Pid:    int(v.Pid),
			})
			continue
		}
		Stat, _ := statInfo.Status()
		cpu, _ := statInfo.CPUPercent()
		threads, _ := statInfo.NumThreads()
		name, _ := statInfo.Name()
		isRuning, _ := statInfo.IsRunning()
		mis, _ := statInfo.MemoryInfo()
		stack := mis.Stack
		fmt.Println("isRuning", isRuning)
		running := "not run"
		if isRuning {
			running = "running"
		}
		MemoryData := mis.Data
		ret.Data = append(ret.Data, &protocol.HostPidInfo{
			Pid:         v.Pid,
			MemoryStack: stack,
			MemoryData:  MemoryData,
			Threads:     int64(threads),
			IsRuning:    running,
			Cpu:         float32(cpu),
			Name:        name,
			Stat:        Stat,
		})
		var gridStat int = 0
		if Stat == "Z" { // down 了 进行物理kill
			gridStat = 0
		} else {
			gridStat = 1
		}
		now := time.Now()
		storage.UpdateGrid(&pojo.Grid{
			Id:         int(v.GridId),
			Status:     gridStat,
			Pid:        int(statInfo.Pid),
			UpdateTime: &now,
		})
		storage.SaveStatLog(&pojo.StatLog{
			GridId:      int(v.GridId),
			Stat:        BEHAVIOR_CHECK,
			Pid:         int(statInfo.Pid),
			CPU:         cpu,
			Threads:     threads,
			Name:        name,
			IsRunning:   running,
			MemoryStack: stack,
			MemoryData:  MemoryData,
		})

	}
	return ret, nil
}

func (s *fileTransferServer) PatchServer(ctx context.Context, in *protocol.PatchServerReq) (*protocol.BasicResp, error) {
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("Recovered from panic: %v", r)
			storage.PushErr(&pojo.SystemErr{
				Type: "system/error/SgridPackageServer.PatchServer.recover",
				Info: "recover Error :" + err.Error(),
			})
		}
	}()
	gridsInfo := in.Req
	var servantIds = make(map[int]struct{})
	for _, v := range gridsInfo {
		servantIds[int(v.ServantId)] = struct{}{}
	}
	toIds := make([]int, 0, len(servantIds))
	for k := range servantIds {
		toIds = append(toIds, k)
	}
	confs, err := storage.BatchQueryServantConf(toIds)
	if err != nil {
		storage.PushErr(&pojo.SystemErr{
			Type: "system/error/SgridPackageServer.PactchServer.storage.BatchQueryServantConf",
			Info: "getConfError" + err.Error(),
		})
		return nil, err
	}
	servantIds2Grids := make(map[int][]*protocol.PatchServerDto)

	for _, psd := range gridsInfo {
		psd2 := servantIds2Grids[int(psd.ServantId)]
		psd2 = append(psd2, psd)
		servantIds2Grids[int(psd.ServantId)] = psd2
	}
	fmt.Println("servantIds2Grids", servantIds2Grids)

	for servantId, gridList := range servantIds2Grids {
		for processIndex, req := range gridList {
			id := req.GridId
			fmt.Println("confs", confs)
			fmt.Println("servantId", servantId)
			if confs[servantId] == nil {
				storage.PushErr(&pojo.SystemErr{
					Type: "system/error/SgridPackageServer.PactchServer.confs[servantId] == nil",
					Info: "confs[servantId] is nil , please check configuration",
				})
				continue
			}
			servantConf := confs[servantId].Conf
			if servantConf == "" {
				storage.PushErr(&pojo.SystemErr{
					Type: "system/error/SgridPackageServer.PactchServer.confs[servantId] == nil",
					Info: "conf is nil , please check configuration",
				})
				continue
			}

			execFilePath := req.ExecPath                                    // 服务路径
			serverLanguage := req.ServerLanguage                            // 服务语言
			serverName := req.ServerName                                    // 服务名称
			serverProtocol := req.ServerProtocol                            // 服务协议
			startDir := SgridPackageInstance.JoinPath(Servants, serverName) // 工作目录
			isDir := public.IsDir(startDir)
			if !isDir {
				storage.PushErr(&pojo.SystemErr{
					Type: "system/error/SgridPackageServer.PactchServer.startDir == nil",
					Info: fmt.Sprintf("[%v] is empty or not a dir , please check startDir", startDir),
				})
				continue
			}
			host := req.Host
			fmt.Println("servantConf", servantConf)
			fmt.Println("serverProtocol", serverProtocol)
			fmt.Println("serverLanguage", serverLanguage)

			// Host 确保与主控配置文件一致
			if host != globalConf.Server.Host {
				fmt.Println("server is not equal")
				continue
			}
			item, ok := globalGrids.Load(int(id))
			if ok && item != nil {
				item.(*SgridMonitor).kill()
				globalGrids.Delete(int(id))
			}
			var cmd *exec.Cmd
			var startFile string // 启动文件

			if serverProtocol == public.PROTOCOL_GRPC {
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
					startFile = filepath.Join(startDir, req.ExecPath) // 启动文件
					cmd = exec.Command(startFile)
				} else if serverLanguage == public.RELEASE_CUSTOM_COMMAND {
					var parseExecArgs []string
					err = json.Unmarshal([]byte(execFilePath), &parseExecArgs)
					if err != nil {
						return &protocol.BasicResp{
							Code:    -1,
							Message: fmt.Sprintf("ReleaseServerByPackage.json.Unmarshal([]byte(execFilePath), &parseExecArgs).error %v", err.Error()),
						}, err
					}
					fmt.Println("parseExecArgs", parseExecArgs)
					cmd = exec.Command(parseExecArgs[0], parseExecArgs[1:]...)
				}
			}

			if serverProtocol == public.PROTOCOL_HTTP {
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
					startFile = filepath.Join(startDir, req.ExecPath) // 启动文件
					cmd = exec.Command(startFile)
				} else if serverLanguage == public.RELEASE_CUSTOM_COMMAND {
					var parseExecArgs []string
					err = json.Unmarshal([]byte(execFilePath), &parseExecArgs)
					if err != nil {
						return &protocol.BasicResp{
							Code:    -1,
							Message: fmt.Sprintf("ReleaseServerByPackage.json.Unmarshal([]byte(execFilePath), &parseExecArgs).error %v", err.Error()),
						}, err
					}
					fmt.Println("parseExecArgs", parseExecArgs)
					cmd = exec.Command(parseExecArgs[0], parseExecArgs[1:]...)
				}
			}
			env := append(
				os.Environ(),
				fmt.Sprintf("%v=%v", public.ENV_TARGET_PORT, req.Port),          // 指定端口
				fmt.Sprintf("%v=%v", public.ENV_PRODUCTION, startDir),           // 开启目录
				fmt.Sprintf("%v=%v", public.SGRID_CONFIG, servantConf),          // 配置
				fmt.Sprintf("%v=%v", public.ENV_PROCESS_INDEX, processIndex),    // 服务运行索引
				fmt.Sprintf("%v=%v", public.ENV_SGRID_SERVANT_NAME, serverName), // 服务名
			)
			cmd.Dir = startDir // 指定工作目录
			cmd.Env = env      // 指定环境变量
			fmt.Println("startFile", startFile)
			fmt.Println("cmd.Env", cmd.Env)

			monitor := NewSgridMonitor(
				WithMonitorInterval(time.Second*5),
				WithMonitorSetCmd(cmd),
				WithMonitorServerName(serverName),
				WithMonitorGridIDAndPort(int(id), int(req.Port)),
			)
			globalGrids.LoadOrStore(int(id), monitor)
			monitor.PrintLogger()
			err = cmd.Start()
			fmt.Println("*************服务启动**************")
			if err != nil {
				storage.PushErr(&pojo.SystemErr{
					Type: "system/error/SgridPackageServer.PatchServer.cmd.Start()",
					Info: err.Error(),
				})
			}
			storage.SaveStatLog(&pojo.StatLog{
				GridId: int(id),
				Stat:   BEHAVIOR_PULL,
				// Pid:    monitor.getPid(),
			})

			fmt.Println("*************开始日志上报**************")
			go monitor.Report()
		}
	}

	return &protocol.BasicResp{
		Code:    0,
		Message: "ok",
	}, nil
}

func (s *fileTransferServer) DeletePacakge(ctx context.Context, req *protocol.DeletePackageReq) (*protocol.BasicResp, error) {
	packageFile := SgridPackageInstance.JoinPath(App, req.ServerName, req.FilePath) // 路径
	fmt.Println("packageFile : ", packageFile)
	err := os.Remove(packageFile)
	if err != nil {
		return &protocol.BasicResp{
			Code:    -1,
			Message: fmt.Sprintf("os.Remove.error %v", err.Error()),
		}, nil
	}
	err = storage.DeletePackage(int(req.Id))
	if err != nil {
		return &protocol.BasicResp{
			Code:    -2,
			Message: fmt.Sprintf("db.delete.package.error %v", err.Error()),
		}, nil
	}
	return &protocol.BasicResp{
		Code:    0,
		Message: "ok",
	}, nil
}

func (s *fileTransferServer) GetSystemInfo(ctx context.Context, req *emptypb.Empty) (*protocol.GetSystemInfoResp, error) {
	cpuPercent := public.GetCpuPercent()
	memInfo, _ := mem.VirtualMemory()
	cores, _ := cpu.Counts(false)

	rsp := &protocol.SystemInfo{}
	rsp.CpuLength = fmt.Sprintf("%v", cores)
	rsp.CpuPercent = cpuPercent
	rsp.MemoryPercent = fmt.Sprintf("%.2f", memInfo.UsedPercent)
	rsp.Host = config.GlobalConf.Server.Host
	rsp.MemorySize = fmt.Sprintf("%v", memInfo.Total/1024/1024/1024)
	rsp.Nodes = ""
	return &protocol.GetSystemInfoResp{
		Code:    0,
		Message: "ok",
		Data:    rsp,
	}, nil
}

func (c *fileTransferServer) Notify(ctx context.Context, in *protocol.NotifyReq) (*protocol.BasicResp, error) {
	storage.SaveStatLog(&pojo.StatLog{
		GridId:      int(in.GetGridId()),
		Stat:        BEHAVIOR_SERVANT_REPORT,
		Pid:         0,
		CPU:         0,
		Threads:     0,
		Name:        in.GetServerName(),
		IsRunning:   in.GetInfo(),
		MemoryStack: 0,
		MemoryData:  0,
	})
	return nil, nil
}

func (c *fileTransferServer) InvokeWithCmd(ctx context.Context, in *protocol.InvokeWithCmdReq) (*protocol.InvokeWithCmdRsp, error) {
	key := int(in.GetGridId())
	grid, ok := globalGrids.Load(key)
	if grid == nil || !ok {
		return nil, nil
	}

	cmd := in.GetCmd() + "\n"
	_, err := grid.(*SgridMonitor).stdin.Write([]byte(cmd))
	if err != nil {
		fmt.Println("debug.InvokeWithCmd.Write.error", err.Error())
		return nil, err
	}

	return &protocol.InvokeWithCmdRsp{}, nil
}

func initDir() {
	public.CheckDirectoryOrCreate(SgridPackageInstance.JoinPath(Logger))
	public.CheckDirectoryOrCreate(SgridPackageInstance.JoinPath(App))
	public.CheckDirectoryOrCreate(SgridPackageInstance.JoinPath(Servants))
}

var NewSgridLogTraceServant *rpc.SgridGrpcClient[logProto.SgridLogTraceServiceClient]

func initClient() {
	gn := storage.QueryPropertiesByKey(SgridLogTraceServerHosts)
	addresses := []string{}
	for _, v := range gn {
		addresses = append(addresses, v.Value)
	}
	fmt.Println("SgridPackageServer.initClient.address", addresses)
	packageServant, err := rpc.NewSgridGrpcClient[logProto.SgridLogTraceServiceClient](
		addresses,
		rpc.WithDiaoptions[logProto.SgridLogTraceServiceClient](
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		),
		rpc.WithClientService[logProto.SgridLogTraceServiceClient](logProto.NewSgridLogTraceServiceClient),
		rpc.WithRequestPrefix[logProto.SgridLogTraceServiceClient]("/SgridLogTrace.SgridLogTraceService/"),
	)
	if err != nil {
		fmt.Println("Error To NewSgridGrpcClient ", err.Error())
	}
	NewSgridLogTraceServant = packageServant
}

type SgridPackage struct{}

func (s *SgridPackage) Registry(sc *config.SgridConf) {
	globalConf = sc
	initDir()
	initClient()
	AntsPool, _ = ants.NewPool(100, ants.WithPanicHandler(func(i interface{}) {
		info := fmt.Sprintf("failed to listen: %v", i)
		storage.PushErr(&pojo.SystemErr{
			Type: "system/error/SgridPackageServer/WithPanicHandler",
			Info: info,
		})
	}))
	pool.InitStorage(sc)
	port := fmt.Sprintf(":%v", sc.Server.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		info := fmt.Sprintf("failed to listen:%v ", err)
		storage.PushErr(&pojo.SystemErr{
			Type: "system/error/SgridPackageServer/grpc",
			Info: info,
		})
	}

	grpcServer := grpc.NewServer()
	protocol.RegisterFileTransferServiceServer(grpcServer, &fileTransferServer{})
	fmt.Println("SgridPackage svr started on", port)
	if err := grpcServer.Serve(lis); err != nil {
		info := fmt.Sprintf("failed to serve: %v", err)
		storage.PushErr(&pojo.SystemErr{
			Type: "system/error/SgridPackageServer/grpc",
			Info: info,
		})
	}
}

func (s *SgridPackage) NameSpace() string {
	return "server.SgridPackageServer"
}

func (s *SgridPackage) ServerPath() string {
	return strings.ReplaceAll(s.NameSpace(), ".", "/")
}

func (s *SgridPackage) JoinPath(args ...string) string {
	p := path.Join(args...)
	return public.Join(s.ServerPath(), p)
}
