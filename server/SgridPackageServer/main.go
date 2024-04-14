// ************************ SgridCloud **********************
// SgridPackageServer created at 2024.4.8
// Author @chelizichen
// Operations and Deployment Services
// ************************ SgridCloud **********************

package main

import (
	"Sgrid/src/config"
	"Sgrid/src/configuration"
	protocol "Sgrid/src/proto"
	"Sgrid/src/public"
	"Sgrid/src/public/pool"
	"Sgrid/src/storage"
	"Sgrid/src/storage/pojo"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"sync/atomic"
	"time"

	"github.com/robfig/cron"
	p "github.com/shirou/gopsutil/process"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

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

func getStat(pid int) *p.Process {
	process, err := p.NewProcess(int32(pid))
	if err != nil {
		fmt.Println("Error creating new process:", err)
		return nil
	}
	return process
}

func (s *SgridMonitor) Report() {
	for {
		fmt.Println("Next Load Report", s.next.Load())
		if s.next.Load() {
			break
		}
		time.Sleep(s.interval)
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
			now := time.Now()
			content := fmt.Sprintf("time:%v |serverName:%v | pid:%v | cpu:%v | thread:%v | status:%v \n", now.Format(time.DateTime), s.serverName, s.getPid(), cpu, threads, status)
			storage.UpdateGrid(&pojo.Grid{
				Id:         s.gridId,
				Status:     gridStat,
				Pid:        s.getPid(),
				UpdateTime: &now,
			})
			s.statLog.Write([]byte(content))
		})
	}
}

func (s *SgridMonitor) PrintLogger() {
	op, err := s.cmd.StdoutPipe()
	if err != nil {
		fmt.Println("StdoutPipe | Error", err.Error())
	}
	ep, err := s.cmd.StderrPipe()
	if err != nil {
		fmt.Println("StderrPipe | Error", err.Error())
	}
	s.getFile()
	go func() {
		for {
			// 读取输出
			buf := make([]byte, 1024)
			time := time.Now().Format(time.DateTime)
			n, err := op.Read(buf)
			if err != nil {
				break
			}
			// 打印输出
			content := time + "|ServerName|" + s.serverName + "|" + string(buf[:n]) + "\n"
			s.dataLog.Write([]byte(content))
		}
	}()
	go func() {
		for {
			// 读取输出
			buf := make([]byte, 1024)
			time := time.Now().Format(time.DateTime)
			n, err := ep.Read(buf)
			if err != nil {
				break
			}
			// 打印输出
			content := time + "|ServerName|" + s.serverName + "|" + string(buf[:n]) + "\n"
			s.errLog.Write([]byte(content))
		}
	}()

	go func() {
		spec := public.CRON_EVERY_DAY
		s.cron = cron.New()
		err = s.cron.AddFunc(spec, func() {
			s.getFile()
		})
		go s.cron.Start()
	}()
}

func (s *SgridMonitor) kill() {
	s.cmd.Process.Kill()
	s.next.Store(true)
}

func (s *SgridMonitor) getPid() int {
	return s.cmd.Process.Pid
}

func (s *SgridMonitor) getFile() {
	today := time.Now().Format(time.DateOnly)
	directoryPath := public.Join(Logger, s.serverName)
	err := public.CheckDirectoryOrCreate(directoryPath)
	if err != nil {
		fmt.Println("CheckDirectoryOrCreate Error", err.Error())
	}
	logDataPath := public.Join(Logger, s.serverName, fmt.Sprintf("log-data-%v.log", today))
	logErrorPath := public.Join(Logger, s.serverName, fmt.Sprintf("log-error-%v.log", today))
	logStatPath := public.Join(Logger, s.serverName, fmt.Sprintf("log-stat-%v.log", today))
	opf, err := os.OpenFile(logDataPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("OpenFile Error", logDataPath)
	}
	epf, err := os.OpenFile(logErrorPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("OpenFile Error", logErrorPath)
	}
	spf, err := os.OpenFile(logStatPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("OpenFile Error", logStatPath)
	}
	s.dataLog = opf
	s.errLog = epf
	s.statLog = spf
}

func WithMonitorInterval(interval time.Duration) func(*SgridMonitor) {
	return func(monitor *SgridMonitor) {
		if interval.Seconds() < 5 { // min 5s
			interval = time.Second * 5
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

func NewSgridMonitor(opt ...WithSgridMonitorConfFunc) *SgridMonitor {
	monitor := &SgridMonitor{
		next: atomic.Bool{},
	}
	for _, v := range opt {
		fn := v
		fn(monitor)
	}
	return monitor
}

func WithMonitorServerName(serverName string) func(*SgridMonitor) {
	return func(monitor *SgridMonitor) {
		monitor.serverName = serverName
	}
}

type fileTransferServer struct {
	protocol.UnimplementedFileTransferServiceServer
}

const (
	App      = "application"
	Servants = "servants"
	Logger   = "logger"
)

var globalConf *config.SgridConf
var globalPool *pool.RoutinePool
var globalGrids map[int]*SgridMonitor = make(map[int]*SgridMonitor)

func initDir() {
	public.CheckDirectoryOrCreate(public.Join(Logger))
	public.CheckDirectoryOrCreate(public.Join(App))
	public.CheckDirectoryOrCreate(public.Join(Servants))
}

func initSgridConf() *config.SgridConf {
	initDir()
	sc, err := public.NewConfig()
	if err != nil {
		fmt.Println("Error To NewConfig", err)
	}
	configuration.InitStorage(sc)
	return sc
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
	directoryPath := public.Join(App, servername)
	err := public.CheckDirectoryOrCreate(directoryPath)
	if err != nil {
		fmt.Println("check directory error")
	}
	targetFilePath := public.Join(App, servername, filename)
	file, err := os.Create(targetFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 循环接收文件块，直到流结束
	for {
		fileChunk, err := stream.Recv()
		if err == io.EOF {
			// 流结束，退出循环
			break
		}
		if err != nil {
			return err
		}

		// 写入文件块到文件
		_, err = file.Write(fileChunk.Data)
		if err != nil {
			return err
		}

		// Respond to the client (optional)
		response := &protocol.FileResp{
			Msg:  "Chunk received",
			Code: 200,
		}
		if err := stream.Send(response); err != nil {
			return err
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
	i := req.GetGridId()
	p := req.GetPid()

	sm := globalGrids[int(i)]
	if sm.getPid() == int(p) {
		sm.kill()
		delete(globalGrids, int(i))
		return &protocol.BasicResp{
			Code:    0,
			Message: "ok",
		}, nil
	}
	return &protocol.BasicResp{
		Code:    -1,
		Message: "error",
	}, nil
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

// 发布 -> 上报给主控
func (s *fileTransferServer) ReleaseServerByPackage(ctx context.Context, req *protocol.ReleaseServerReq) (res *protocol.BasicResp, err error) {
	if len(req.ServantGrids) == 0 {
		return
	}
	filePath := req.FilePath                              // 服务路径
	serverLanguage := req.ServerLanguage                  // 服务语言
	serverName := req.ServerName                          // 服务名称
	serverProtocol := req.ServerProtocol                  // 协议
	execFilePath := req.ExecPath                          // 执行路径
	startDir := public.Join(Servants, serverName)         // 解压目录
	packageFile := public.Join(App, serverName, filePath) // 路径
	public.Tar2Dest(packageFile, startDir)                // 解压
	servantGrid := req.ServantGrids                       // 服务列表  通过Host过滤拿到IP，然后进行服务启动
	var startFile string                                  // 启动文件
	CheckProdConf(path.Join(startDir, public.DEV_CONF_NAME), path.Join(startDir, public.PROD_CONF_NAME))
	for _, grid := range servantGrid { // 通过Host过滤拿到IP，然后进行服务启动
		GRID := grid
		id := GRID.GridId
		fmt.Println("GRID", GRID)
		if GRID.Ip != globalConf.Server.Host {
			fmt.Println("server is not equal")
			return
		} else {
			err = public.CheckDirectoryOrCreate(startDir)
			if err != nil {
				fmt.Println("error", err.Error())
			}
			var cmd *exec.Cmd
			if serverProtocol == public.PROTOCOL_GRPC {
				if serverLanguage == public.RELEASE_GO {
					startFile = public.Join(Servants, serverName, execFilePath) // 启动文件
					cmd = exec.Command(startFile)
				}
				if serverLanguage == public.RELEASE_NODE {
					startFile = public.Join(Servants, serverName, execFilePath) // 启动文件
					cmd = exec.Command("node", startFile)
				}
			}

			if serverProtocol == public.PROTOCOL_HTTP {
				if serverLanguage == public.RELEASE_GO {
					startFile = public.Join(Servants, serverName, execFilePath) // 启动文件
					cmd = exec.Command(startFile)
				}
				if serverLanguage == public.RELEASE_NODE {
					startFile = public.Join(Servants, serverName, execFilePath) // 启动文件
					cmd = exec.Command("node", startFile)
				}
			}
			cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", public.ENV_TARGET_PORT, grid.Port), fmt.Sprintf("%v=%v", public.ENV_PRODUCTION, startDir))
			fmt.Println("cmd.Env", cmd.Env)
			monitor := NewSgridMonitor(
				WithMonitorInterval(time.Second*5),
				WithMonitorSetCmd(cmd),
				WithMonitorServerName(serverName),
				WithMonitorGridID(int(id)),
			)

			delete(globalGrids, int(id))
			globalGrids[int(id)] = monitor

			monitor.PrintLogger()
			go monitor.Report()

			go func() {
				err = cmd.Start()
				if err != nil {
					fmt.Println("error", err.Error())
				}
			}()
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

func main() {
	globalPool = pool.NewRoutinePool(1000)
	go globalPool.Run()
	sc := initSgridConf()
	fmt.Println("sc", sc)
	port := fmt.Sprintf(":%v", sc.Server.Port)
	globalConf = sc
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
