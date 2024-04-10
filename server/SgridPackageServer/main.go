package main

import (
	file_gen "Sgrid/src/proto/file.gen"
	"Sgrid/src/public"
	"errors"
	"io"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type fileTransferServer struct {
	file_gen.UnimplementedFileTransferServiceServer
}

const (
	App = "application"
)

func (s *fileTransferServer) StreamFile(stream file_gen.FileTransferService_StreamFileServer) error {
	// 创建文件来存储接收到的数据
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return errors.New("无法获取元数据")
	}
	// 获取文件名和哈希值
	filename := md.Get("filename")[0]
	targetFilePath := public.Join(App, filename)
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
		response := &file_gen.FileResp{
			Msg:  "Chunk received",
			Code: 200,
		}
		if err := stream.Send(response); err != nil {
			return err
		}
	}

	// 发送文件接收完成的响应
	finalResponse := &file_gen.FileResp{
		Msg:  "File received successfully",
		Code: 200,
	}
	if err := stream.Send(finalResponse); err != nil {
		return err
	}

	return nil
}

func main() {
	sc, err := public.NewConfig()
	if err != nil {
		log.Fatalf("failed to NewConfig: %v", err)
	}
	lis, err := net.Listen(sc.Server.Type, ":"+string(rune(sc.Server.Port)))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	file_gen.RegisterFileTransferServiceServer(grpcServer, &fileTransferServer{})
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
