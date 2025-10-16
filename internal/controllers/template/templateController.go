package template

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/services/template"
	"auto-forge/pkg/errors"
	log "auto-forge/pkg/logger"

	"github.com/gin-gonic/gin"
)

var tplService = template.NewTemplateService()

// CreateTemplate 创建模板
func CreateTemplate(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	var req request.CreateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	tpl, err := tplService.CreateTemplate(userID, &req)
	if err != nil {
		log.Error("创建模板失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, tpl, "创建模板成功")
}

// GetTemplateList 获取模板列表
func GetTemplateList(c *gin.Context) {
	var req request.ListTemplatesRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	result, err := tplService.GetTemplateList(&req)
	if err != nil {
		log.Error("获取模板列表失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "获取模板列表失败"))
		return
	}

	errors.ResponseSuccess(c, result, "获取模板列表成功")
}

// GetTemplateDetail 获取模板详情
func GetTemplateDetail(c *gin.Context) {
	templateID := c.Param("id")
	if templateID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "模板ID不能为空"))
		return
	}

	detail, err := tplService.GetTemplateDetail(templateID)
	if err != nil {
		log.Error("获取模板详情失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, detail, "获取模板详情成功")
}

// InstallTemplate 安装模板
func InstallTemplate(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	var req request.InstallTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	result, err := tplService.InstallTemplate(userID, &req)
	if err != nil {
		log.Error("安装模板失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, result, "安装模板成功")
}

// UpdateTemplate 更新模板
func UpdateTemplate(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	templateID := c.Param("id")
	if templateID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "模板ID不能为空"))
		return
	}

	var req request.UpdateTemplateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	tpl, err := tplService.UpdateTemplate(templateID, userID, &req)
	if err != nil {
		log.Error("更新模板失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, tpl, "更新模板成功")
}

// DeleteTemplate 删除模板
func DeleteTemplate(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	templateID := c.Param("id")
	if templateID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "模板ID不能为空"))
		return
	}

	if err := tplService.DeleteTemplate(templateID, userID); err != nil {
		log.Error("删除模板失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, nil, "删除模板成功")
}

// GetInstallHistory 获取安装历史
func GetInstallHistory(c *gin.Context) {
	userID := c.GetString("user_id")
	if userID == "" {
		errors.HandleError(c, errors.New(errors.CodeUnauthorized, "未授权"))
		return
	}

	history, err := tplService.GetUserInstallHistory(userID)
	if err != nil {
		log.Error("获取安装历史失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "获取安装历史失败"))
		return
	}

	errors.ResponseSuccess(c, history, "获取安装历史成功")
}
