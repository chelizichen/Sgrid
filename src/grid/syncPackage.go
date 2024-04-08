package grid

import (
	c "Sgrid/src/configuration"
	h "Sgrid/src/http"
	"Sgrid/src/storage/pojo"
	"bytes"
	"fmt"
	"io"
	"net/http"
)

func SyncPackage(body []byte, ctx h.SimpHttpServerCtx) error {
	var uploadList []pojo.Node
	fmt.Println("ctx", ctx.Host)
	c.GORM.Model(&pojo.Node{}).
		Where("ip != ?", ctx.Host).
		Find(&uploadList)
	fmt.Println("uploadList", uploadList)
	for _, v := range uploadList {
		V := v
		fmt.Println("upload", V)
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
	return nil
}
