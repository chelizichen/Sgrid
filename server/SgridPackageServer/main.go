// ************************ SgridCloud **********************
// SgridPackageServer created at 2024.4.8
// Author @chelizichen
// Operations and Deployment Services
// ************************ SgridCloud **********************

package SgridPackageServer

import (
	protocol "Sgrid/server/SgridPackageServer/proto"
	"Sgrid/src/config"
	"Sgrid/src/configuration"
	"Sgrid/src/public"
	"Sgrid/src/public/pool"
	"Sgrid/src/storage"
	"Sgrid/src/storage/pojo"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync/atomic"
	"time"

	"github.com/robfig/cron"
	p "github.com/shirou/gopsutil/process"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	App      = "application"
	Servants = "servants"
	Logger   = "logger"
)

const (
	BEHAVIOR_PULL  = "pull"
	BEHAVIOR_KILL  = "kill"
	BEHAVIOR_DOWN  = "down"
	BEHAVIOR_ALIVE = "alive"
	BEHAVIOR_CHECK = "check"
)

const CONSTANT_MONITOR_INTERVAL = 30

var globalConf *config.SgridConf
var globalPool *pool.RoutinePool
var globalGrids map[int]*SgridMonitor = make(map[int]*SgridMonitor)

var SgridPackageInstance = &SgridPackage{}

func getStat(pid int) *p.Process {
	process, err := p.NewProcess(int32(pid))
	if err != nil {
		fmt.Println("Error creating new process:", err)
		return nil
	}
	return process
}

func CheckProdConf(devConf, prodConf string) {
	fmt.Println("CheckProdConf", devConf, prodConf)
	if !public.IsExist(prodConf) {
		err := public.CopyFile(devConf, prodConf)
		if err != nil {
			fmt.Println("CheckProdConf", err.Error())
		}
	}
}

type WithSgridMonitorConfFunc func(*SgridMonitor)

type SgridMonitor struct {
	interval   time.Duration // 上报时间
	cmd        *exec.Cmd
	next       atomic.Bool
	serverName string
	cron       *cron.Cron
	dataLog    *os.File
	errLog     *os.File
	statLog    *os.File
	gridId     int
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

func WithMonitorGridID(id int) func(*SgridMonitor) {
	return func(monitor *SgridMonitor) {
		monitor.gridId = id
	}
}

func WithMonitorServerName(serverName string) func(*SgridMonitor) {
	return func(monitor *SgridMonitor) {
		monitor.serverName = serverName
	}
}

func NewSgridMonitor(opt ...WithSgridMonitorConfFunc) *SgridMonitor {
	monitor := &SgridMonitor{
		next: atomic.Bool{},
	}
	for _, v := range opt {
		fn := v
		fn(monitor)
	}
	monitor.getFile()
	return monitor
}

func (s *SgridMonitor) Report() {
	for {
		fmt.Println("Next Load Report", s.next.Load())
		time.Sleep(s.interval)
		if s.next.Load() {
			break
		}
		globalPool.Add(func() {
			id := s.getPid()
			statInfo := getStat(id)
			status, _ := statInfo.Status()
			var gridStat int = 0
			if status == "Z" { // down 了 进行物理kill
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
			content := fmt.Sprintf("time:%v |serverName:%v | pid:%v | cpu:%v | thread:%v | status:%v \n", now.Format(time.DateTime), s.serverName, s.getPid(), cpu, threads, status)
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
			s.statLog.Write([]byte(content))
		})
	}
}

func (s *SgridMonitor) PrintLogger() {
	op, err := s.cmd.StdoutPipe()
	if err != nil {
		fmt.Println("GetPipeError", err.Error())
	}
	ep, err := s.cmd.StderrPipe()
	if err != nil {
		fmt.Println("GetPipeError", err.Error())
	}
	go func() {
		for {
			// 读取输出
			buf := make([]byte, 1024)
			now := time.Now().Format(time.DateTime)
			n, err := op.Read(buf)
			if err != nil {
				time.Sleep(time.Millisecond * 500)
			} else {
				// 打印输出
				content := now + "|ServerName|" + s.serverName + "|" + string(buf[:n]) + "\n"
				nn, err := s.dataLog.Write([]byte(content))
				if err != nil {
					fmt.Println("s.datalog.write error", err.Error())
				}
				fmt.Println("nn", nn)
			}

		}
	}()
	go func() {
		for {
			// 读取输出
			buf := make([]byte, 1024)
			now := time.Now().Format(time.DateTime)
			n, err := ep.Read(buf)
			if err != nil {
				time.Sleep(time.Millisecond * 500)
			} else {
				// 打印输出
				content := now + "|ServerName|" + s.serverName + "|" + string(buf[:n]) + "\n"
				nn, err := s.errLog.Write([]byte(content))
				if err != nil {
					fmt.Println("s.datalog.write error", err.Error())
				}
				fmt.Println("nn", nn)
			}

		}
	}()

	go func() {
		spec := public.CRON_EVERY_DAY
		s.cron = cron.New()
		err := s.cron.AddFunc(spec, func() {
			s.getFile()
		})
		if err != nil {
			fmt.Println("Error", err.Error())
		}
		go s.cron.Start()
	}()
}

func (s *SgridMonitor) kill() {
	s.cron.Stop()
	storage.SaveStatLog(&pojo.StatLog{
		GridId: s.gridId,
		Stat:   BEHAVIOR_KILL,
		Pid:    s.getPid(),
	})
	s.dataLog.Close()
	s.errLog.Close()
	s.statLog.Close()
	err := s.cmd.Process.Kill()
	if err != nil {
		fmt.Println("kill error", err.Error())
	}
	s.next.Store(true)
}

func (s *SgridMonitor) getPid() int {
	if s.cmd != nil {
		if s.cmd.Process != nil {
			return s.cmd.Process.Pid
		}
	}
	return 0
}

func (s *SgridMonitor) getFile() {

	today := time.Now().Format(time.DateOnly)
	directoryPath := SgridPackageInstance.JoinPath(Logger, s.serverName)
	err := public.CheckDirectoryOrCreate(directoryPath)
	if err != nil {
		fmt.Println("CheckDirectoryOrCreate Error", err.Error())
	}
	logDataPath := SgridPackageInstance.JoinPath(Logger, s.serverName, fmt.Sprintf("log-data-%v.log", today))
	logErrorPath := SgridPackageInstance.JoinPath(Logger, s.serverName, fmt.Sprintf("log-error-%v.log", today))
	logStatPath := SgridPackageInstance.JoinPath(Logger, s.serverName, fmt.Sprintf("log-stat-%v.log", today))

	opf, err := os.OpenFile(logDataPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("OpenFile Error", logDataPath)
	}
	s.dataLog = opf

	epf, err := os.OpenFile(logErrorPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("OpenFile Error", logErrorPath)
	}
	s.errLog = epf

	spf, err := os.OpenFile(logStatPath, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0777)
	if err != nil {
		fmt.Println("OpenFile Error", logStatPath)
	}
	s.statLog = spf

	fmt.Println("Then PrintLogger")
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

func (s *fileTransferServer) DeletePackage(ctx context.Context, req *protocol.DeletePackageReq) (res *protocol.BasicResp, err error) {
	if len(req.FilePath) == 0 || len(req.ServerName) == 0 {
		fmt.Println("invoke Error", req)
		return &protocol.BasicResp{
			Code:    -1,
			Message: "invoke Error ",
		}, err
	}
	f := req.FilePath
	svr := req.ServerName
	t := public.Join(svr, f)
	err = os.Remove(t)
	if err != nil {
		return &protocol.BasicResp{
			Code:    -1,
			Message: "error" + err.Error(),
		}, err
	}
	return &protocol.BasicResp{
		Code:    0,
		Message: "ok",
	}, nil
}

func (s *fileTransferServer) ShutdownGrid(ctx context.Context, req *protocol.ShutdownGridReq) (res *protocol.BasicResp, err error) {
	for _, _grid := range req.GetReq() {
		grid := _grid
		h := grid.GetHost()

		if globalConf.Server.Host == h {
			i := grid.GetGridId()
			sm, ok := globalGrids[int(i)]
			if ok && sm != nil {
				storage.SaveStatLog(&pojo.StatLog{
					GridId: int(i),
					Stat:   BEHAVIOR_DOWN,
					Pid:    sm.getPid(),
				})
				sm.kill()
				delete(globalGrids, int(i))
			}
		}
	}
	return &protocol.BasicResp{
		Code:    0,
		Message: "ok",
	}, nil
}

// 发布 -> 上报给主控
func (s *fileTransferServer) ReleaseServerByPackage(ctx context.Context, req *protocol.ReleaseServerReq) (res *protocol.BasicResp, err error) {
	if len(req.ServantGrids) == 0 {
		return
	}
	ReleaseServerByPackageReq, err := json.Marshal(req)
	if err != nil {
		fmt.Println("error", err.Error())
	}
	fmt.Println("ReleaseServerByPackage Req ||", string(ReleaseServerByPackageReq))
	filePath := req.FilePath                                                // 服务路径
	serverLanguage := req.ServerLanguage                                    // 服务语言
	serverName := req.ServerName                                            // 服务名称
	serverProtocol := req.ServerProtocol                                    // 协议
	execFilePath := req.ExecPath                                            // 执行路径
	startDir := SgridPackageInstance.JoinPath(Servants, serverName)         // 解压目录
	packageFile := SgridPackageInstance.JoinPath(App, serverName, filePath) // 路径
	public.Tar2Dest(packageFile, startDir)                                  // 解压
	servantGrid := req.ServantGrids                                         // 服务列表  通过Host过滤拿到IP，然后进行服务启动
	var startFile string                                                    // 启动文件
	CheckProdConf(path.Join(startDir, public.DEV_CONF_NAME), path.Join(startDir, public.PROD_CONF_NAME))
	for _, grid := range servantGrid { // 通过Host过滤拿到IP，然后进行服务启动
		GRID := grid
		id := GRID.GridId
		fmt.Println("GRID.IP", GRID.Ip)
		fmt.Println("globalConf.Server.Host", globalConf.Server.Host)
		if GRID.Ip != globalConf.Server.Host {
			fmt.Println("server is not equal")
		} else {
			item, ok := globalGrids[int(id)]
			if ok && item != nil { // 终止
				item.kill()
				delete(globalGrids, int(id))
			}
			fmt.Println("GRID", GRID)
			fmt.Println("serverProtocol", serverProtocol)
			err = public.CheckDirectoryOrCreate(startDir)
			if err != nil {
				fmt.Println("error", err.Error())
			}
			var cmd *exec.Cmd
			if serverProtocol == public.PROTOCOL_GRPC {
				if serverLanguage == public.RELEASE_GO {
					startFile = SgridPackageInstance.JoinPath(Servants, serverName, execFilePath) // 启动文件
					cmd = exec.Command(startFile)
				}
				if serverLanguage == public.RELEASE_NODE {
					startFile = SgridPackageInstance.JoinPath(Servants, serverName, execFilePath) // 启动文件
					cmd = exec.Command("node", startFile)
				}
				if serverLanguage == public.RELEASE_JAVA {
					startFile = SgridPackageInstance.JoinPath(Servants, serverName, execFilePath) // 启动文件
					prodConf := path.Join(startDir, public.PROD_CONF_NAME)
					cmd = exec.Command("java", "-jar", startFile, fmt.Sprintf("-Dspring.config.location=file:%v", prodConf))
					cmd.Env = append(cmd.Env, fmt.Sprintf("SGRID_PROD_CONF_PATH=%v", prodConf))
				}
			}

			if serverProtocol == public.PROTOCOL_HTTP {
				if serverLanguage == public.RELEASE_GO {
					startFile = SgridPackageInstance.JoinPath(Servants, serverName, execFilePath) // 启动文件
					cmd = exec.Command(startFile)
				}
				if serverLanguage == public.RELEASE_NODE {
					startFile = SgridPackageInstance.JoinPath(Servants, serverName, execFilePath) // 启动文件
					cmd = exec.Command("node", startFile)
				}
				if serverLanguage == public.RELEASE_JAVA {
					startFile = SgridPackageInstance.JoinPath(Servants, serverName, execFilePath) // 启动文件
					exec.Command("java", "-jar", startFile, "-Dspring.config.location="+"file:"+path.Join(startDir, public.PROD_CONF_NAME))
				}
			}
			fmt.Println("startFile", startFile)
			cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", public.ENV_TARGET_PORT, grid.Port), fmt.Sprintf("%v=%v", public.ENV_PRODUCTION, startDir))
			fmt.Println("cmd.Env", cmd.Env)

			monitor := NewSgridMonitor(
				WithMonitorInterval(time.Second*5),
				WithMonitorSetCmd(cmd),
				WithMonitorServerName(serverName),
				WithMonitorGridID(int(id)),
			)

			globalGrids[int(id)] = monitor

			monitor.PrintLogger()
			err = cmd.Start()
			fmt.Println("*************服务启动**************")
			if err != nil {
				fmt.Println("服务启动失败 ｜", err.Error())
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
	dir := SgridPackageInstance.JoinPath(Logger, in.ServerName)
	fmt.Println("dir", dir)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fs []string
	for _, v := range files {
		fs = append(fs, v.Name())
	}
	fmt.Println("fs", fs)
	return &protocol.GetLogFileByHostResp{
		FileName: fs,
	}, nil
}

func (s *fileTransferServer) GetLogByFile(ctx context.Context, in *protocol.GetLogByFileReq) (*protocol.GetLogByFileResp, error) {
	filePath := SgridPackageInstance.JoinPath(Logger, in.ServerName, in.LogFile)
	s2, err := public.GetLogger(filePath, in.Pattern, int(in.Rows))
	if err != nil {
		return &protocol.GetLogByFileResp{
			Data: s2,
		}, err
	}
	return &protocol.GetLogByFileResp{
		Data: s2,
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
		if v.Host == globalConf.Server.Host {
			fmt.Println("int(v.Pid)", int(v.Pid))
			statInfo := getStat(int(v.Pid))
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

	}
	return ret, nil
}

func initDir() {
	public.CheckDirectoryOrCreate(SgridPackageInstance.JoinPath(Logger))
	public.CheckDirectoryOrCreate(SgridPackageInstance.JoinPath(App))
	public.CheckDirectoryOrCreate(SgridPackageInstance.JoinPath(Servants))
}

type SgridPackage struct{}

func (s *SgridPackage) Registry(sc *config.SgridConf) {
	initDir()
	globalPool = pool.NewRoutinePool(1000)
	globalConf = sc
	go globalPool.Run()
	configuration.InitStorage(sc)
	port := fmt.Sprintf(":%v", sc.Server.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	protocol.RegisterFileTransferServiceServer(grpcServer, &fileTransferServer{})
	fmt.Println("Sgrid svr started on", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
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
