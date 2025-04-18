package service

import (
	protocol "Sgrid/server/SgridPackageServer/proto"
	handlers "Sgrid/src/http"
	"Sgrid/src/public"
	"Sgrid/src/rpc"
	"Sgrid/src/storage"
	"Sgrid/src/storage/dto"
	"Sgrid/src/storage/pojo"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/metadata"
)

const SgridPackageServerHosts = "SgridPackageServerHosts"

func PackageService(ctx *handlers.SgridServerCtx) {
	router := ctx.Engine.Group(strings.ToLower(ctx.GetServerName()))
	packageServant := ctx.Context.Value(PackageServantProxy{}).(*rpc.SgridGrpcClient[protocol.FileTransferServiceClient])
	clients := packageServant.GetClients()
	router.POST("/upload/uploadServer", func(c *gin.Context) {
		// 从请求中获取服务器名、文件哈希和文件
		serverName := c.PostForm("serverName")
		fileHash := c.PostForm("fileHash")
		servantId, _ := strconv.Atoi(c.PostForm("servantId"))
		F, err := c.FormFile("file")
		content := c.PostForm("content")
		servantLanguage := c.PostForm("servantLanguage")
		if err != nil {
			handlers.AbortWithError(c, "file error"+string(err.Error()))
			return
		}
		file, err := F.Open()
		if err != nil {
			handlers.AbortWithError(c, "open error"+string(err.Error()))
			return
		}
		defer file.Close()
		now := time.Now()
		dateTime := strings.ReplaceAll(now.Format(time.DateTime), " ", "")
		fileName := ""
		if servantLanguage == public.RELEASE_EXE {
			fileName = fmt.Sprintf("%v_%v_%v", serverName, dateTime, fileHash)
		} else if servantLanguage == public.RELEASE_JAVA_JAR {
			fileName = fmt.Sprintf("%v_%v_%v.jar", serverName, dateTime, fileHash)
		} else {
			fileName = fmt.Sprintf("%v_%v_%v.tgz", serverName, dateTime, fileHash)

		}
		META := metadata.Pairs("filename", fileName, "serverName", serverName)
		ctx := metadata.NewOutgoingContext(context.Background(), META)
		var syncPackage sync.WaitGroup
		for _, client := range clients {
			syncPackage.Add(1)
			go func(clt protocol.FileTransferServiceClient) {
				stream, _ := clt.StreamFile(ctx)
				buffer := make([]byte, public.ChunkFileSize)
				go func() {
					for {
						res, recvErr := stream.Recv()
						if recvErr != nil {
							fmt.Println("src.service.package.uploadServer.recvErr ", recvErr.Error())
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
				}()
				var of int64 = 0
				go func() {
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
							if sendError != io.EOF {
								defer func() {
									if err := stream.CloseSend(); err != nil {
										fmt.Println("service.package.uploadServer.error 无法关闭流", err.Error())
										handlers.AbortWithError(c, "StreamFile error"+string(err.Error()))
										return
									}
								}()
								fmt.Println("service.package.uploadServer.stream.Send,EOF ERROR")
								break
							} else {
								defer func() {
									if err := stream.CloseSend(); err != nil {
										fmt.Println("service.package.uploadServer.error 无法关闭流", err.Error())
										handlers.AbortWithError(c, "StreamFile error"+string(err.Error()))
										return
									}
								}()
								handlers.AbortWithError(c, "service.package.uploadServer.stream.send.error "+string(sendError.Error()))
							}
						}
					}
				}()
			}(client)
		}
		syncPackage.Wait()
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
		fmt.Println("id", id)
		fmt.Println("serverName", serverName)
		if err != nil || id == 0 || len(serverName) == 0 {
			handlers.AbortWithError(c, "params error ! [id] or [serverName]"+err.Error())
			return
		}
		sp := storage.QueryPackageById(id)
		var wg sync.WaitGroup
		for _, client := range clients {
			wg.Add(1)
			go func(clt protocol.FileTransferServiceClient) {
				clt.DeletePacakge(&gin.Context{}, &protocol.DeletePackageReq{
					FilePath:   sp.FilePath,
					ServerName: serverName,
					Id:         int32(sp.Id),
				})
				wg.Done()
			}(client)
		}
		wg.Wait()
		storage.DeletePackage(id)
		handlers.AbortWithSucc(c, nil)
	})

	router.POST("/restart/server", func(c *gin.Context) {
		var req *protocol.ReleaseServerReq
		err := c.BindJSON(&req)
		if err != nil {
			handlers.AbortWithError(c, err.Error())
			return
		}
		var patchServerReq = &protocol.PatchServerReq{
			Req: make([]*protocol.PatchServerDto, 0),
		}
		for _, rt := range req.ServantGrids {
			psd := &protocol.PatchServerDto{
				ServerName:     req.ServerName,
				ServerLanguage: req.ServerLanguage,
				ServantId:      req.ServantId,
				ServerProtocol: req.ServerProtocol,
				ExecPath:       req.ExecPath,
				Port:           rt.Port,
				GridId:         rt.GridId,
				Host:           rt.Ip,
			}
			patchServerReq.Req = append(patchServerReq.Req, psd)
		}
		var wg sync.WaitGroup
		for _, client := range clients {
			wg.Add(1)
			go func(clt protocol.FileTransferServiceClient) {
				clt.PatchServer(&gin.Context{}, patchServerReq)
				wg.Done()
			}(client)
		}
		wg.Wait()
		handlers.AbortWithSucc(c, nil)
	})

	router.POST("/release/server", func(c *gin.Context) {
		var req *protocol.ReleaseServerReq
		err := c.BindJSON(&req)
		if err != nil {
			handlers.AbortWithError(c, err.Error())
			return
		}
		var releaseId string = c.Query("releaseId")
		if releaseId == "" {
			handlers.AbortWithError(c, "PARAMS ERROR >>> releaseId is empty")
			return
		}
		var wg sync.WaitGroup
		fmt.Println("req", req)
		key := fmt.Sprintf("server.version.%d", req.ServantId)
		err = storage.ResetProperty(&pojo.Properties{
			Key:   key,
			Value: releaseId,
		})
		if err != nil {
			handlers.AbortWithError(c, err.Error())
			return
		}
		for _, client := range clients {
			wg.Add(1)
			go func(clt protocol.FileTransferServiceClient) {
				clt.ReleaseServerByPackage(&gin.Context{}, req)
				wg.Done()
			}(client)
		}
		wg.Wait()
		handlers.AbortWithSucc(c, nil)
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
			handlers.AbortWithError(c, err.Error())
			return
		}
		var wg sync.WaitGroup
		for _, client := range clients {
			wg.Add(1)
			go func(c protocol.FileTransferServiceClient) {
				c.ShutdownGrid(&gin.Context{}, req)
				wg.Done()
			}(client)
		}
		wg.Wait()
		handlers.AbortWithSucc(c, nil)
	})

	router.GET("/statlog/getlist", func(c *gin.Context) {
		id, _ := strconv.Atoi(c.DefaultQuery("id", "0"))
		offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
		size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))
		sl := storage.QueryStatLogById(id, offset, size)
		handlers.AbortWithSucc(c, sl)
	})

	router.GET("/statlog/getLogFileList", func(c *gin.Context) {
		host := c.Query("host")
		serverName := c.Query("serverName")
		grid, _ := strconv.Atoi(c.Query("gridId"))
		var wg sync.WaitGroup
		var resp *protocol.GetLogFileByHostResp
		for index, client := range clients {
			idx := index
			wg.Add(1)
			go func(client protocol.FileTransferServiceClient) {
				u, err := packageServant.ParseHost(idx, "grpc")
				fmt.Println("SgridGrpcClient", u.Hostname(), "|", host)
				if err != nil {
					fmt.Println("ParseHost Error", err.Error())
					wg.Done()
					return
				}
				if u.Hostname() == host {
					r, err := client.GetLogFileByHost(&gin.Context{}, &protocol.GetLogFileByHostReq{
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
			}(client)
		}
		wg.Wait()
		handlers.AbortWithSucc(c, resp.Data)
	})

	router.POST("/statlog/getLog", func(c *gin.Context) {
		var req *protocol.GetLogByFileReq
		var resp *protocol.GetLogByFileResp
		c.BindJSON(&req)
		var wg sync.WaitGroup
		for index, client := range clients {
			wg.Add(1)
			idx := index
			go func(clt protocol.FileTransferServiceClient) {
				u, err := packageServant.ParseHost(idx, "grpc")
				if err != nil {
					handlers.AbortWithError(c, err.Error())
					wg.Done()
					return
				}
				if u.Hostname() == req.Host {
					r, err := clt.GetLogByFile(&gin.Context{}, &protocol.GetLogByFileReq{
						Host:       u.Host,
						ServerName: req.ServerName,
						Keyword:    req.Keyword,
						GridId:     req.GridId,
						LogType:    req.LogType,
						LogFileName: req.LogFileName,
						Len:        req.Len,
					})
					if err != nil {
						fmt.Println("error", err.Error())
						handlers.AbortWithError(c, err.Error())
						wg.Done()
						return
					}
					resp = r
					wg.Done()
					return
				}
				wg.Done()

			}(client)
		}
		wg.Wait()
		handlers.AbortWithSuccList(c, resp.Data, resp.Total)
	})

	router.POST("/statlog/check", func(c *gin.Context) {
		var req *protocol.GetPidInfoReq
		var resp []*protocol.GetPidInfoResp
		c.BindJSON(&req)
		var wg sync.WaitGroup
		for _, client := range clients {
			wg.Add(1)
			go func(clt protocol.FileTransferServiceClient) {
				r, err := clt.GetPidInfo(&gin.Context{}, req)
				if err != nil {
					fmt.Println("error", err.Error())
				}
				resp = append(resp, r)
				wg.Done()
			}(client)
		}
		wg.Wait()
		handlers.AbortWithSucc(c, resp)
	})

	router.POST("/statlog/deleteByLogType", func(c *gin.Context) {
		var req *dto.LogTypeDto
		c.BindJSON(&req)
		storage.DeleteLogByLogType(req.DateTime, req.LogType)
		handlers.AbortWithSucc(c, nil)
	})
	router.POST("/notify", func(c *gin.Context) {
		var rpl protocol.BasicResp
		var req *protocol.NotifyReq
		err := c.BindJSON(&req)
		if err != nil {
			handlers.AbortWithError(c, err.Error())
			return
		}
		if req.GridId == 0 || req.ServerName == "" {
			handlers.AbortWithError(c, "input error ")
			return
		}
		packageServant.Request(rpc.RequestPack{
			Method: "NotifyReq",
			Body:   req,
			Reply:  &rpl,
		})
		handlers.AbortWithSucc(c, nil)
	})

	router.GET("/download/serverPackage", func(ctx *gin.Context) {
		fileName := ctx.Query("fileName")
		serverName := ctx.Query("serverName")
		if fileName == "" || serverName == "" {
			handlers.AbortWithError(ctx, errors.New("fileName or serverName is empty").Error())
			return
		}
		fmt.Println("fileName", fileName)
		cwd, _ := os.Getwd()
		targetFilePath := filepath.Join(cwd, "server", "SgridPackageServer", "application", serverName, fileName)
		ctx.Header("Content-Type", "application/octet-stream")
		ctx.Header("Content-Disposition", "attachment; filename="+fileName)
		ctx.FileAttachment(targetFilePath, fileName)
	})

	router.POST("/main/invokeWithCmd", func(c *gin.Context) {
		var req *protocol.InvokeWithCmdReq
		err := c.BindJSON(&req)
		if err != nil {
			handlers.AbortWithError(c, err.Error())
			return
		}
		fmt.Println("debugger >> invokeWithCmd >> req", req)
		var wg sync.WaitGroup
		for _, client := range clients {
			wg.Add(1)
			go func(clt protocol.FileTransferServiceClient) {
				clt.InvokeWithCmd(&gin.Context{}, req)
				wg.Done()
			}(client)
		}
		wg.Wait()
		handlers.AbortWithSucc(c, nil)
	})
}
