package routes

import (
	"auto-forge/internal/controllers/preview"

	"github.com/gin-gonic/gin"
)

// RegisterPreviewRoutes 注册预览相关路由
func RegisterPreviewRoutes(r *gin.Engine) {
	previewController := preview.NewPreviewController()

	// HTML 预览路由
	r.GET("/preview/:id", previewController.GetHtmlPreview)
}
