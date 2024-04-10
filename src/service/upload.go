package service

import (
	handlers "Sgrid/src/http"
	file_gen "Sgrid/src/proto/file.gen"
	"Sgrid/src/storage"
	"Sgrid/src/storage/pojo"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	"github.com/gin-gonic/gin"
)

func UploadService(ctx *handlers.SimpHttpServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
	addresses := []string{"http://47.98.174.10:24283", "http://150.158.120.244/:24283"}
	clients := []*file_gen.FileTransferServiceClient{}
	for _, v := range addresses {
		add := v
		conn, err := grpc.Dial(add, grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			log.Fatalf("无法连接: %v", err)
		}
		client := file_gen.NewFileTransferServiceClient(conn)
		clients = append(clients, &client)
		defer conn.Close()
	}

	GROUP.POST("/upload/uploadServer", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")
		// 从请求中获取服务器名、文件哈希和文件
		serverName := c.PostForm("serverName")
		fileHash := c.PostForm("fileHash")
		servantId, _ := strconv.Atoi(c.PostForm("servantId"))
		F, err := c.FormFile("file")
		contetnt := c.PostForm("content")
		if err != nil {
			c.Writer.WriteString(string(err.Error()))
			return
		}

		// 打开文件
		file, err := F.Open()
		if err != nil {
			c.Writer.WriteString(string(err.Error()))
			return
		}
		defer file.Close()
		dateTime := time.Now().Format(time.DateTime)
		fileName := fmt.Sprintf("%v_%v_%v.tgz", serverName, dateTime, fileHash)
		META := metadata.Pairs("filename", fileName)
		ctx := metadata.NewOutgoingContext(context.Background(), META)
		for _, client := range clients {
			CLIENT := client
			// 设置 gRPC 客户端
			stream, err := (*CLIENT).StreamFile(ctx)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			defer func() {
				if err := stream.CloseSend(); err != nil {
					log.Fatalf("无法关闭流: %v", err)
				}
			}()
			// 缓冲区大小，即每次发送的文件块大小
			bufferSize := 1024 * 10
			buffer := make([]byte, bufferSize)

			// 逐个发送文件块
			for {
				// 从文件中读取数据到缓冲区
				n, err := file.Read(buffer)
				if err != nil {
					// 如果读取到文件末尾，则跳出循环
					if err == io.EOF {
						break
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}

				// 构造文件块
				chunk := &file_gen.FileChunk{
					Data:       buffer[:n],
					Offset:     int64(n), // 当前文件块在文件中的偏移量
					ServerName: serverName,
					FileHash:   fileHash,
				}

				// 发送文件块
				if err := stream.Send(chunk); err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				log.Printf("已发送带有偏移量的文件块: %d", chunk.Offset)
			}
			response, err := stream.Recv()
			if err != nil {
				c.Writer.WriteString(string(err.Error()))
				continue
			}
			b, err := json.Marshal(response)
			if err != nil {
				c.Writer.WriteString(string(b))
				continue
			}
			c.Writer.WriteString(string(b))
		}

		// 接收服务器的响应
		if err != nil {
			c.Writer.WriteString(string(err.Error()))
			return
		}
		storage.SaveHashPackage(pojo.ServantPackage{
			ServantId:  servantId,
			Hash:       fileHash,
			FilePath:   fileName,
			Content:    contetnt,
			CreateTime: dateTime,
		})
		c.Writer.WriteString("done!")

	})

	ctx.Engine.Use(GROUP.Handlers...)
}
