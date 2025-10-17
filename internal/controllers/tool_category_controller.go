package controllers

import (
	"auto-forge/pkg/common"
	"auto-forge/pkg/utools"

	"github.com/gin-gonic/gin"
)

// GetToolCategories 获取工具分类列表
func GetToolCategories(c *gin.Context) {
	categories := utools.GetToolCategories()
	common.Success(c, categories, "获取工具分类成功")
}
