package preview

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)


type PreviewController struct{}


func NewPreviewController() *PreviewController {
	return &PreviewController{}
}


func (c *PreviewController) GetHtmlPreview(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "预览 ID 不能为空",
		})
		return
	}


	filePath := filepath.Join("./data/html-preview", id+".html")
	content, err := os.ReadFile(filePath)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "预览内容不存在或已过期",
		})
		return
	}

	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.Header("X-Content-Type-Options", "nosniff")
	ctx.Data(http.StatusOK, "text/html; charset=utf-8", content)
}
