package grid

import (
	c "Sgrid/src/configuration"
	h "Sgrid/src/http"
	"Sgrid/src/storage/pojo"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// 将服务包同步给其他的服务
func SyncPackage(body []byte, ctx h.SimpHttpServerCtx) error {
	var uploadList []pojo.Node
	c.GORM.Model(&pojo.Node{}).Where("ip != ?", ctx.Host).Find(&uploadList)
	body = bytes.Replace(body, []byte("CLUSTER_REQUEST"), []byte("SINGLE_REQUEST"), -1)
	fmt.Println("uploadList", uploadList)
	fmt.Println("ctx", ctx.Host)
	for _, v := range uploadList {
		V := v
		if V.UploadPath != "" {
			// 创建一个新的请求，将 body 数据发送给其他服务器
			req, err := http.NewRequest("POST", V.UploadPath, bytes.NewBuffer(body))
			if err != nil {
				return err
			}
			req.Header.Set("Content-Type", "application/json")

			// 发送请求
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			// 读取响应
			respBody, err := io.ReadAll(resp.Body)
			if err != nil {
				return err
			}
			fmt.Println("Response from other server:", string(respBody))
		}
	}
	return nil
}

// 同步请求给其他端口
func SyncPostRequest(body []byte, ctx h.SimpHttpServerCtx, g *gin.Context) error {
	var uploadList []pojo.Node
	c.GORM.Model(&pojo.Node{}).Where("ip != ?", ctx.Host).Find(&uploadList)
	body = bytes.Replace(body, []byte("CLUSTER_REQUEST"), []byte("SINGLE_REQUEST"), -1)
	for _, v := range uploadList {
		V := v
		URL := "http://" + V.Ip + "/" + g.Request.URL.Path
		// 创建一个新的请求，将 body 数据发送给其他服务器
		req, err := http.NewRequest("POST", URL, bytes.NewBuffer(body))
		if err != nil {
			return err
		}
		req.Header.Set("Content-Type", "application/json")
		// 发送请求
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// 读取响应
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		fmt.Println("responseBody", string(respBody))
		g.Writer.WriteString(string(respBody))
		g.Writer.Flush()
	}
	return nil
}
