package SgridPackageServer

import (
	protocol "Sgrid/server/SgridPackageServer/proto"
	"Sgrid/src/config"
	"Sgrid/src/pool"
	"Sgrid/src/public"
	"Sgrid/src/storage"
	"Sgrid/src/storage/pojo"
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	pk "Sgrid/src/public/process"

	"github.com/panjf2000/ants"
	"github.com/robfig/cron"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	p "github.com/shirou/gopsutil/process"

	"google.golang.org/grpc"
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
	language 	 string
}

func NewSgridMonitor(cmd *exec.Cmd, serverName string, id int, port int, t time.Duration) *SgridMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	monitor := &SgridMonitor{
		ctx:    ctx,
		cancel: cancel,
		interval:  t,
		cmd:cmd,
		serverName: serverName,
		gridId: id,
		port: port,
		cronInstance: cron.New(), // 初始化 
	}
	stdin, err := monitor.cmd.StdinPipe()
	if err != nil {
		fmt.Println("Failed to get StdinPipe:", err)
		return nil
	}
	monitor.stdin = stdin
	return monitor
}
func (s *SgridMonitor) Start() error{
	return s.cmd.Start()
}

func (s *SgridMonitor) kill() {
	s.cancel()
    if s.cronInstance != nil {
        s.cronInstance.Stop()
    }
	fmt.Println("system/log/server.kill |", s.serverName)
	storage.SaveStatLog(&pojo.StatLog{
		GridId: s.gridId,
		Stat:   BEHAVIOR_KILL,
		Pid:    s.getPid(),
	})
	var err error
	if s.cmd == nil || s.cmd.Process == nil {
        fmt.Println("system/warn/process.kill | process is already terminated")
        return
    }
	if(s.language == public.RELEASE_PYTHON_EXE || s.language == public.RELEASE_PYTHON_TAR){
		err = s.cmd.Process.Signal(syscall.SIGINT)
	}else{
		err = s.cmd.Process.Kill()
	}
	if err != nil {
		fmt.Println("system/err/process.kill | ", err.Error())
		storage.PushErr(&pojo.SystemErr{
			Type: "system/error/SgridPackageServer/s.cmd.Process.Kill()",
			Info: err.Error(),
		})
	}
}

func (s *SgridMonitor) getPid() int {
	if process := s.cmd.Process; s.cmd != nil && process != nil {
		return process.Pid
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
    for _, grid := range req.GetReq() {
        // 检查主机是否匹配
        if grid.GetHost() != globalConf.Server.Host {
            continue
        }
        
        gridId := int(grid.GetGridId())
        
        // 尝试获取并关闭已存在的进程
        if sm, ok := globalGrids.Load(gridId); ok && sm != nil {
            monitor := sm.(*SgridMonitor)
            fmt.Printf("Common Down [%d]\n", gridId)
            storage.SaveStatLog(&pojo.StatLog{
                GridId: gridId,
                Stat:   BEHAVIOR_DOWN,
                Pid:    monitor.getPid(),
            })
            monitor.kill()
            globalGrids.Delete(gridId)
        } else {
            // 强制终止进程
            fmt.Printf("Force Down [%d]\n", gridId)
            pk.SgridProcessUtil.QueryProcessPidThenKill(int(grid.GetPort()))
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
			err = fmt.Errorf("recovered from panic: %v", r)
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
	logDir := SgridPackageInstance.JoinPath(Logger,serverName) // 路径
	servantId := req.ServantId                                              // 服务ID
	err = public.CheckDirectoryOrCreate(startDir)                           // 检查并创建目录
	fmt.Println("packageFile", packageFile)
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
        fmt.Printf("globalGrids item: %v, ok: %v\n", item, ok)
		if ok && item != nil { // 终止
			item.(*SgridMonitor).kill()
			globalGrids.Delete(int(id))
		}
        fmt.Printf("invoke >> CreateCommand")
		cmd,err := CreateCommand(
			serverProtocol,
			serverName,
			serverLanguage,
			startDir,
			logDir,
			servantConf,
			execFilePath,
			int(grid.Port),
			ProcessIndex,
		)
		if err != nil{
			return &protocol.BasicResp{
				Code:    -1,
				Message: fmt.Sprintf("ReleaseServerByPackage.json.Unmarshal([]byte(execFilePath), &parseExecArgs).error %v", err.Error()),
			}, err
		}
		monitor := NewSgridMonitor(cmd,serverName,int(id),int(grid.Port),time.Second * 60)
		monitor.language = serverLanguage
		globalGrids.LoadOrStore(int(id), monitor)
		err = monitor.Start()
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
			logDir := SgridPackageInstance.JoinPath(Logger,serverName) // 路径
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
			cmd,err := CreateCommand(
				serverProtocol,
				serverName,
				serverLanguage,
				startDir,
				logDir,
				servantConf,
				execFilePath,
				int(req.Port),
				processIndex,
			)
			if err != nil{
				return &protocol.BasicResp{
					Code:    -1,
					Message: fmt.Sprintf("ReleaseServerByPackage.json.Unmarshal([]byte(execFilePath), &parseExecArgs).error %v", err.Error()),
				}, err
			}
			monitor := NewSgridMonitor(cmd,serverName,int(id),int(req.Port),time.Second * 60)
			globalGrids.LoadOrStore(int(id), monitor)
			err = monitor.Start()
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
		}
	}

	return &protocol.BasicResp{
		Code:    0,
		Message: "ok",
	}, nil
}

func (s *fileTransferServer) GetLogFileByHost(ctx context.Context, in *protocol.GetLogFileByHostReq) (*protocol.GetLogFileByHostResp, error) {
	logDir := SgridPackageInstance.JoinPath(Logger, in.ServerName)
	local_host := fmt.Sprintf("%v://%v", "grpc", in.Host)
	p,err := url.Parse(local_host)
	if err != nil{
		return nil,err
	}
	fmt.Println("p.Hostname()", p.Hostname())
	fmt.Println("globalConf.Server.Host", globalConf.Server.Host)
	if p.Hostname()!= globalConf.Server.Host {
		return nil, nil
	}
	logList,err := GetLogFileList(logDir)
	return &protocol.GetLogFileByHostResp{
		Data:    logList,
		Code:    0,
		Message: "ok",
	}, err
}

func (s *fileTransferServer) GetLogByFile(ctx context.Context, in *protocol.GetLogByFileReq) (*protocol.GetLogByFileResp, error) {
	logTargetPath := SgridPackageInstance.JoinPath(Logger, in.ServerName, in.LogFileName)
	fmt.Println("logTargetPath", logTargetPath)
	local_host := fmt.Sprintf("%v://%v", "grpc", in.Host)
	p,err := url.Parse(local_host)
	if err != nil{
		return nil,err
	}
	fmt.Println("p.Hostname()", p.Hostname())
	fmt.Println("globalConf.Server.Host", globalConf.Server.Host)
	if p.Hostname()!= globalConf.Server.Host {
		return nil, nil
	}
	if in.LogType == 0{
		in.LogType = 2
	}
	fmt.Println("len",in.Len)
	logContent ,err := SearchLog(logTargetPath,in.GetLogType(),in.GetKeyword(),in.GetLen())
	return &protocol.GetLogByFileResp{
		Data:  logContent,
		Total: int64(len(logContent)),
	}, err
}

func (s *fileTransferServer) GetPidInfo(ctx context.Context, in *protocol.GetPidInfoReq) (resp *protocol.GetPidInfoResp, err error) {
	if len(in.HostPids) == 0 {
		return nil, nil
	}
	ret := &protocol.GetPidInfoResp{
		Data: []*protocol.HostPidInfo{},
	}
	for _, v := range in.HostPids {
		if v.Host != globalConf.Server.Host {
			continue
		}
		item, ok := globalGrids.Load(int(v.GridId))
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
		pid := item.(*SgridMonitor).getPid()
		statInfo := getStat(item.(*SgridMonitor).getPid())
		if statInfo == nil {
			fmt.Println("error process", pid)
			now := time.Now()
			storage.UpdateGrid(&pojo.Grid{
				Id:         int(v.GridId),
				Status:     0,
				Pid:        pid,
				UpdateTime: &now,
			})
			storage.SaveStatLog(&pojo.StatLog{
				GridId: int(v.GridId),
				Stat:   BEHAVIOR_CHECK,
				Pid:    pid,
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
		running := SERVER_NOT_RUN
		if isRuning {
			running = SERVER_RUNING
		}
		MemoryData := mis.Data
		ret.Data = append(ret.Data, &protocol.HostPidInfo{
			Pid:         int32(pid),
			MemoryStack: stack,
			MemoryData:  MemoryData,
			Threads:     int64(threads),
			IsRuning:    running,
			Cpu:         float32(cpu),
			Name:        name,
			Stat:        Stat,
		})
		var gridStat int = 0
		if Stat == PROCESS_DEAD {
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



type SgridPackage struct{}

func (s *SgridPackage) Registry(sc *config.SgridConf) {
	globalConf = sc
	initDir()
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
