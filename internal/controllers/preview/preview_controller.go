package preview

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// PreviewController 预览控制器
type PreviewController struct{}

// NewPreviewController 创建预览控制器
func NewPreviewController() *PreviewController {
	return &PreviewController{}
}

// GetHtmlPreview 获取 HTML 预览
func (c *PreviewController) GetHtmlPreview(ctx *gin.Context) {
	id := ctx.Param("id")
	if id == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "预览 ID 不能为空",
		})
		return
	}

	// 读取 HTML 文件
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
