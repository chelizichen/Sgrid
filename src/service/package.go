package service

import (
	protocol "Sgrid/server/SgridPackageServer/proto"
	handlers "Sgrid/src/http"
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

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

const SgridPackageServerHosts = "SgridPackageServerHosts"

func PackageService(ctx *handlers.SgridServerCtx) {
	router := ctx.Engine.Group(strings.ToLower(ctx.Name))
	clients := ctx.Context.Value(public.GRPC_CLIENT_PROXYS{}).([]*clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient])

	router.POST("/upload/uploadServer", func(c *gin.Context) {
		// 从请求中获取服务器名、文件哈希和文件
		serverName := c.PostForm("serverName")
		fileHash := c.PostForm("fileHash")
		servantId, _ := strconv.Atoi(c.PostForm("servantId"))
		F, err := c.FormFile("file")
		content := c.PostForm("content")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "file error"+string(err.Error()), nil))
			return
		}
		file, err := F.Open()
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "open error"+string(err.Error()), nil))
			return
		}
		defer file.Close()
		now := time.Now()
		dateTime := strings.ReplaceAll(now.Format(time.DateTime), " ", "")
		fileName := fmt.Sprintf("%v_%v_%v.tgz", serverName, dateTime, fileHash)
		META := metadata.Pairs("filename", fileName, "serverName", serverName)
		ctx := metadata.NewOutgoingContext(context.Background(), META)
		var syncPackage sync.WaitGroup
		for _, client := range clients {
			syncPackage.Add(1)
			go func(client *clientgrpc.SgridGrpcClient[protocol.FileTransferServiceClient]) {
				stream, _ := client.GetClient().StreamFile(ctx)
				defer func() {
					if err := stream.CloseSend(); err != nil {
						fmt.Println("err.Error()", err.Error())
						log.Fatalf("无法关闭流: %v", err)
						c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "StreamFile error"+string(err.Error()), nil))
					}
				}()
				buffer := make([]byte, public.ChunkFileSize)
				var of int64 = 0
				for {
					// 构造文件块
					n, readErr := file.ReadAt(buffer, of)
					if readErr == io.EOF && n != 0 {
						fmt.Println("读取完成至最后一个Chunk", readErr.Error(), n)
						of += int64(n)
					} else if readErr == io.EOF && n == 0 {
						fmt.Println("读取完成", readErr.Error())
						break
					} else if readErr != nil && readErr != io.EOF {
						fmt.Println("读取错误", readErr.Error())
						break
					} else {
						of += int64(n)

					}
					fmt.Println("read offset", n)
					chunk := &protocol.FileChunk{
						Data:       buffer[:n],
						Offset:     int64(n), // 当前文件块在文件中的偏移量
						ServerName: serverName,
						FileHash:   fileHash,
					}
					if sendError := stream.Send(chunk); sendError != nil {
						if sendError == io.EOF {
							fmt.Println("Send EOF " + string(sendError.Error()))
						} else {
							break
						}
					}
				}
				for {
					res, recvErr := stream.Recv()
					if recvErr != nil {
						fmt.Println("recvErr", recvErr.Error())
						break
					}
					if res != nil {
						fmt.Println("fr.Message Success", res)
						if res.Code == 200 {
							syncPackage.Done()
							break
						}
					}

				}
			}(client)
		}
		syncPackage.Wait()
		// 接收服务器的响应
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "server error"+string(err.Error()), nil))
			return
		}
		storage.SaveHashPackage(pojo.ServantPackage{
			ServantId:  servantId,
			Hash:       fileHash,
			FilePath:   fileName,
			Content:    content,
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
		fmt.Println("req", req)
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

	router.GET("/release/updatePackageVersion", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.Query("id"))
		version := c.Query("version")
		s := &pojo.ServantPackage{
			Id:      id,
			Version: version,
		}
		storage.UpdatePackageVersion(s)
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
		grid, _ := strconv.Atoi(c.Query("gridId"))
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
						GridId:     int64(grid),
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
		c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(int(resp.Code), resp.Message, resp.Data))
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
						Offset:     req.Offset,
						Size:       req.Size,
						GridId:     req.GridId,
						LogType:    req.LogType,
						DateTime:   req.DateTime,
					})
					if err != nil {
						fmt.Println("error", err.Error())
						c.AbortWithStatusJSON(http.StatusOK, handlers.Resp(-1, "err"+err.Error(), nil))
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
}
