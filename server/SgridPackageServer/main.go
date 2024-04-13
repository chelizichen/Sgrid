package main

import (
	"Sgrid/src/config"
	protocol "Sgrid/src/proto"
	"Sgrid/src/public"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/exec"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type WithSgridMonitorConfFunc func() func(*SgridMonitor)

type SgridMonitor struct {
	interval int // 上报时间
}

func WithMonitorInterval(interval int) func(*SgridMonitor) {
	return func(monitor *SgridMonitor) {
		monitor.interval = interval
	}
}

func NewSgridMonitor() *SgridMonitor {
	return &SgridMonitor{}
}

type fileTransferServer struct {
	protocol.UnimplementedFileTransferServiceServer
}

const (
	App      = "application"
	Servants = "servants"
)

var db *sql.DB
var globalConf *config.SgridConf

func initSgridConf() *config.SgridConf {
	sc, err := public.NewConfig()
	if err != nil {
		fmt.Println("Error To NewConfig", err)
	}
	fmt.Println("sc", sc.Server.Storage)
	S, err := gorm.Open(mysql.Open(sc.Server.Storage), &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "grid_",
			SingularTable: true,
		},
	})
	if err != nil {
		fmt.Println("Error To init gorm", err)
	}
	db, err = S.DB()
	if err != nil {
		fmt.Println("Error To DB", err)
	}
	err = db.Ping()
	if err != nil {
		fmt.Println("Error To Ping", err)
	}
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

// 发布 -> 上报给主控
func (s *fileTransferServer) ReleaseServerByPackage(ctx context.Context, req *protocol.ReleaseServerReq) (res *protocol.BasicResp, err error) {
	filePath := req.FilePath             // 服务路径
	serverLanguage := req.ServerLanguage // 服务语言
	serverName := req.ServerName         // 服务名称
	serverProtocol := req.ServerProtocol // 协议
	execFilePath := req.ExecPath         // 执行路径
	if len(req.ServantGrids) == 0 {
		return
	}
	servantGrid := req.ServantGrids // 服务
	// 通过Host过滤拿到IP，然后进行服务启动
	for _, grid := range servantGrid {
		GRID := grid
		fmt.Println("GRID", GRID)
		if GRID.Ip != globalConf.Server.Host {
			fmt.Println("server is not equal")
			return
		} else {
			var startFile string
			var packageFile string
			startDir := public.Join(Servants, serverName)
			err = public.CheckDirectoryOrCreate(startDir)
			if err != nil {
				fmt.Println("error", err.Error())
			}
			var cmd *exec.Cmd
			if serverProtocol == public.PROTOCOL_GRPC {
				if serverLanguage == public.RELEASE_GO {
					startFile = public.Join(Servants, serverName, execFilePath) // 启动文件
					packageFile = public.Join(App, serverName, filePath)
					public.Tar2Dest(packageFile, startDir)
					cmd = exec.Command(startFile)
				}
				if serverLanguage == public.RELEASE_NODE {
					startFile = public.Join(Servants, serverName, execFilePath) // 启动文件
					packageFile = public.Join(App, serverName, filePath)
					public.Tar2Dest(packageFile, startDir)
					cmd = exec.Command("node", startFile)
				}
			}

			if serverProtocol == public.PROTOCOL_HTTP {
				if serverLanguage == public.RELEASE_GO {
					startFile = public.Join(Servants, serverName, execFilePath) // 启动文件
					packageFile = public.Join(App, serverName, filePath)
					public.Tar2Dest(packageFile, startDir)
					cmd = exec.Command(startFile)
				}
				if serverLanguage == public.RELEASE_NODE {
					startFile = public.Join(Servants, serverName, execFilePath) // 启动文件
					packageFile = public.Join(App, serverName, filePath)
					public.Tar2Dest(packageFile, startDir)
					cmd = exec.Command("node", startFile)
				}
			}
			fmt.Println("startFile", startFile)
			fmt.Println("packageFile", packageFile)
			cmd.Env = append(cmd.Env, fmt.Sprintf("%v=%v", public.ENV_TARGET_PORT, grid.Port), fmt.Sprintf("%v=%v", public.ENV_PRODUCTION, "SgridPackageServer"))
			err = cmd.Start()
			if err != nil {
				fmt.Println("error", err.Error())
			}
		}
	}
	return &protocol.BasicResp{
		Code:    0,
		Message: "ok",
	}, nil

}

func main() {
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
