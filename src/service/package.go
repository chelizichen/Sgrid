package service

import (
	handlers "Sgrid/src/http"
	protocol "Sgrid/src/proto"
	"Sgrid/src/public"
	clientgrpc "Sgrid/src/public/client_grpc"
	"Sgrid/src/storage"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/pojo"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/gin-gonic/gin"
)

const SgridPackageServerHosts = "SgridPackageServerHosts"

func PackageService(ctx *handlers.SgridServerCtx) {
	router := ctx.Engine.Group(strings.ToLower(ctx.Name))
	gn := storage.QueryPropertiesByKey(SgridPackageServerHosts)
	addresses := []string{}

	for _, v := range gn {
		addresses = append(addresses, v.Value)
	}

	clients := []*clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient]{}
	for _, v := range addresses {
		add := v
		conn, err := grpc.Dial(add, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("无法连接: %v", err)
		}
		// defer conn.Close() // 移动到循环内部
		client := clientgrpc.NewSgridClient[protocol.FileTransferServiceClient](
			protocol.NewFileTransferServiceClient(conn),
			clientgrpc.WithSgridGrpcClientAddress[protocol.FileTransferServiceClient](add),
		)
		clients = append(clients, client)
	}
	router.POST("/upload/uploadServer", func(c *gin.Context) {
		// 从请求中获取服务器名、文件哈希和文件
		serverName := c.PostForm("serverName")
		fileHash := c.PostForm("fileHash")
		servantId, _ := strconv.Atoi(c.PostForm("servantId"))
		F, err := c.FormFile("file")
		contetnt := c.PostForm("content")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "file error"+string(err.Error()), nil))
			return
		}

		// 打开文件
		file, err := F.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "open error"+string(err.Error()), nil))
			return
		}
		defer file.Close()
		now := time.Now()
		dateTime := strings.ReplaceAll(now.Format(time.DateTime), " ", "")
		fileName := fmt.Sprintf("%v_%v_%v.tgz", serverName, dateTime, fileHash)
		fmt.Println("fileName ", fileName)
		META := metadata.Pairs("filename", fileName, "serverName", serverName)
		ctx := metadata.NewOutgoingContext(context.Background(), META)
		var wg sync.WaitGroup
		for _, client := range clients {
			wg.Add(1)
			go func(client clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient]) {
				// 设置 gRPC 客户端
				stream, err := client.GetClient().StreamFile(ctx)
				if err != nil {
					fmt.Println("err.Error()", err.Error())
					c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "StreamFile error"+string(err.Error()), nil))
					return
				}
				defer func() {
					if err := stream.CloseSend(); err != nil {
						fmt.Println("err.Error()", err.Error())
						log.Fatalf("无法关闭流: %v", err)
						c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "StreamFile error"+string(err.Error()), nil))
					}
				}()
				// 缓冲区大小，即每次发送的文件块大小
				buffer := make([]byte, public.ChunkFileSize)
				// 逐个发送文件块
				for {
					// 从文件中读取数据到缓冲区
					n, err := file.Read(buffer)
					if err != io.EOF {
						break
					}
					if err != nil {
						fmt.Println("error", err.Error())
						break
					}
					if n == 0 {
						break
					}
					// 构造文件块
					chunk := &protocol.FileChunk{
						Data:       buffer[:n],
						Offset:     int64(n), // 当前文件块在文件中的偏移量
						ServerName: serverName,
						FileHash:   fileHash,
					}

					// 发送文件块
					if err := stream.Send(chunk); err != nil {
						c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "Write error"+string(err.Error()), nil))
						return
					}
					log.Printf("已发送带有偏移量的文件块: %d", chunk.Offset)
				}
				response, err := stream.CloseAndRecv()
				if err != nil {
					c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "recv error"+string(err.Error()), nil))
					return
				}
				fmt.Println("response", response)
				if response.Code == 200 {
					wg.Done()
					return
				}
			}(*client)
		}

		wg.Wait()
		// 接收服务器的响应
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "server error"+string(err.Error()), nil))
			return
		}
		storage.SaveHashPackage(pojo.ServantPackage{
			ServantId:  servantId,
			Hash:       fileHash,
			FilePath:   fileName,
			Content:    contetnt,
			CreateTime: &now,
		})
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})
	router.GET("/upload/getList", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
		version := c.DefaultQuery("version", "")

		params := &dto.QueryPackageDto{
			Id:      id,
			Offset:  offset,
			Size:    size,
			Version: version,
		}

		vsp := storage.QueryPackage(params)
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", vsp))
	})

	router.GET("/upload/removePackage", func(c *gin.Context) {
		id, err := strconv.Atoi(c.DefaultQuery("id", "0"))
		serverName := c.DefaultQuery("serverName", "")
		if err != nil || id == 0 || len(serverName) == 0 {
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "params error ! [id] or [serverName]", nil))
			return
		}
		sp := storage.QueryPackageById(id)
		var wg sync.WaitGroup
		for _, client := range clients {
			wg.Add(1)
			go func(client clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient]) {
				client.GetClient().DeletePackage(&gin.Context{}, &protocol.DeletePackageReq{
					FilePath: sp.FilePath,
				})
				wg.Done()
			}(*client)
		}
		wg.Wait()
		storage.DeletePackage(id)
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "params error ! [id] or [serverName]", nil))
	})

	router.POST("/release/server", func(c *gin.Context) {
		var req *protocol.ReleaseServerReq
		err := c.BindJSON(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "err", err.Error()))
		}
		var wg sync.WaitGroup
		for _, client := range clients {
			wg.Add(1)
			go func(client clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient]) {
				client.GetClient().ReleaseServerByPackage(&gin.Context{}, req)
				wg.Done()
			}(*client)
		}
		wg.Wait()
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	router.POST("/release/shutdown", func(c *gin.Context) {
		var req *protocol.ShutdownGridReq
		err := c.BindJSON(&req)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "err", err.Error()))
			return
		}
		var wg sync.WaitGroup
		for _, client := range clients {
			wg.Add(1)
			go func(client clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient]) {
				client.GetClient().ShutdownGrid(&gin.Context{}, req)
				wg.Done()
			}(*client)
		}
		wg.Wait()
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", nil))
	})

	router.GET("/statlog/getlist", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
		sl := storage.QueryStatLogById(id, offset, size)
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", sl))
	})

	router.GET("/statlog/getLogFileList", func(c *gin.Context) {
		host := c.Query("host")
		serverName := c.Query("serverName")
		var wg sync.WaitGroup
		var resp *protocol.GetLogFileByHostResp
		for _, client := range clients {
			wg.Add(1)
			go func(client clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient]) {
				u, err := client.ParseHost("grpc")
				fmt.Println("SgridGrpcClient", u.Hostname(), "|", host)
				if err != nil {
					fmt.Println("ParseHost Error", err.Error())
					wg.Done()
					return
				}
				if u.Hostname() == host {
					r, err := client.GetClient().GetLogFileByHost(&gin.Context{}, &protocol.GetLogFileByHostReq{
						Host:       u.Host,
						ServerName: serverName,
					})
					if err != nil {
						fmt.Println("error", err.Error())
					}
					resp = r
					wg.Done()
					return
				}
				wg.Done()

			}(*client)
		}
		wg.Wait()
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", resp.FileName))
	})

	router.POST("/statlog/getLog", func(c *gin.Context) {
		var req *protocol.GetLogByFileReq
		var resp *protocol.GetLogByFileResp
		c.BindJSON(&req)
		var wg sync.WaitGroup
		for _, client := range clients {
			wg.Add(1)
			go func(client clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient]) {
				u, _ := client.ParseHost("grpc")
				if u.Hostname() == req.Host {
					fmt.Println("req", req)
					r, err := client.GetClient().GetLogByFile(&gin.Context{}, &protocol.GetLogByFileReq{
						Host:       u.Host,
						ServerName: req.ServerName,
						Pattern:    req.Pattern,
						Rows:       req.Rows,
						LogFile:    req.LogFile,
					})
					if err != nil {
						fmt.Println("error", err.Error())
						wg.Done()
						return
					}
					resp = r
					wg.Done()
					return
				}
				wg.Done()

			}(*client)
		}
		wg.Wait()
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", resp.Data))
	})

	router.POST("/statlog/check", func(c *gin.Context) {
		var req *protocol.GetPidInfoReq
		var resp []*protocol.GetPidInfoResp
		c.BindJSON(&req)
		var wg sync.WaitGroup
		for _, client := range clients {
			wg.Add(1)
			go func(client clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient]) {
				r, err := client.GetClient().GetPidInfo(&gin.Context{}, req)
				if err != nil {
					fmt.Println("error", err.Error())
				}
				resp = append(resp, r)
				wg.Done()
			}(*client)
		}
		wg.Wait()
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(0, "ok", resp))
	})
	ctx.Engine.Use(router.Handlers...)
}
