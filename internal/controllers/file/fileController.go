package file

import (
	"os"
	"strconv"
	"auto-forge/internal/models"
	"auto-forge/pkg/common"

	"github.com/gin-gonic/gin"
)


func DownloadFile(c *gin.Context) {
	fileInfo, exists := c.Get("file_info")
	if !exists {
		common.ServerError(c, "无法获取文件信息")
		return
	}

	file := fileInfo.(*models.UploadFile)


	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		common.NotFound(c, "文件不存在")
		return
	}


	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Disposition", "attachment; filename=\""+file.Filename+"\"")
	c.Header("Content-Length", strconv.FormatInt(file.FileSize, 10))


	c.File(file.FilePath)
}


func PreviewFile(c *gin.Context) {
	fileInfo, exists := c.Get("file_info")
	if !exists {
		common.ServerError(c, "无法获取文件信息")
		return
	}

	file := fileInfo.(*models.UploadFile)


	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		common.NotFound(c, "文件不存在")
		return
	}


	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Disposition", "inline; filename=\""+file.Filename+"\"")


	c.File(file.FilePath)
}
