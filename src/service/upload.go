package service

import (
	handlers "Sgrid/src/http"
	utils "Sgrid/src/utils"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
)

func UploadService(ctx *handlers.SimpHttpServerCtx) {
	GROUP := ctx.Engine.Group(strings.ToLower(ctx.Name))
	// 通过 nginx 代理转发到不同的地址
	// 比如节点 Noded1 Node2 ，服务包上传时直接通过 nginx代理上传 至两个节点上的主控的 uploadServer 即可
	GROUP.POST("/upload/uploadServer", func(c *gin.Context) {
		serverName := c.PostForm("serverName")
		cwd, err := os.Getwd()
		if err != nil {
			fmt.Println("Error To GetWd", err.Error())
		}
		F, err := c.FormFile("file")
		storagePath := filepath.Join(cwd, utils.PublishPath, serverName, F.Filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, handlers.Resp(-1, "接受文件失败", nil))
			return
		}
		// 保存上传的文件到服务器临时目录
		tempPath := filepath.Join(cwd, "temp", F.Filename)
		if err := c.SaveUploadedFile(F, tempPath); err != nil {
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "保存上传的文件到服务器临时目录失败", nil))
			return
		}
		// 校验文件完整性（这里使用MD5哈希值作为示例）
		actualHash, err := utils.CalculateFileHash(tempPath)
		utils.AddHashToPackageName(&storagePath, actualHash)
		if err != nil {
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "计算哈希值失败", nil))
			return
		}
		// 移动文件到目标目录
		fmt.Println("tempPath", tempPath)
		fmt.Println("storagePath", storagePath)

		if err := utils.MoveAndRemove(tempPath, storagePath); err != nil {
			fmt.Println("Error To Rename", err.Error())
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "移动文件失败", nil))
			return
		}
		releaseDoc := c.PostForm("doc")
		storageDocPath := filepath.Join(cwd, utils.PublishPath, serverName, "doc.txt")
		E, err := utils.IFNotExistThenCreate(storageDocPath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, handlers.Resp(-1, "打开或创建文件失败"+err.Error(), nil))
		}
		defer E.Close()
		content := storagePath + "\n" + releaseDoc + "\n"
		E.WriteString(content)
		c.JSON(http.StatusOK, handlers.Resp(0, "上传成功", nil))
	})
	ctx.Engine.Use(GROUP.Handlers...)
}
