package controllers

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/dto/response"
	"auto-forge/internal/services/tool_config"
	"auto-forge/pkg/common"
	log "auto-forge/pkg/logger"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetAvailableTools 获取可用的工具列表（用户端）
func GetAvailableTools(c *gin.Context) {
	service := tool_config.NewToolConfigService()

	tools, err := service.GetAvailableTools()
	if err != nil {
		log.Error("获取可用工具列表失败: %v", err)
		common.ServerError(c, "获取工具列表失败")
		return
	}

	common.Success(c, tools, "获取工具列表成功")
}

// GetAllTools 获取所有工具列表（管理端）
func GetAllTools(c *gin.Context) {
	service := tool_config.NewToolConfigService()

	tools, err := service.GetAllTools()
	if err != nil {
		log.Error("获取所有工具列表失败: %v", err)
		common.ServerError(c, "获取工具列表失败")
		return
	}

	common.Success(c, tools, "获取工具列表成功")
}

// GetToolDetail 获取工具详情（管理端）
func GetToolDetail(c *gin.Context) {
	toolCode := c.Param("code")

	service := tool_config.NewToolConfigService()

	// 获取工具配置（包含加密的配置）
	toolConfig, err := service.GetToolConfig(toolCode)
	if err != nil {
		log.Error("获取工具详情失败: %v", err)
		common.NotFound(c, "工具不存在")
		return
	}

	// 解密配置
	decryptedConfig, err := service.GetToolConfigDecrypted(toolCode)
	if err != nil {
		log.Error("解密工具配置失败: %v", err)
		common.ServerError(c, "解密配置失败")
		return
	}

	// 构建响应
	resp := response.ToolConfigDetailResponse{
		ToolConfig:      toolConfig,
		DecryptedConfig: decryptedConfig,
	}

	common.Success(c, resp, "获取工具详情成功")
}

// UpdateToolConfig 更新工具配置（管理端）
func UpdateToolConfig(c *gin.Context) {
	toolCode := c.Param("code")

	var req request.UpdateToolConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	service := tool_config.NewToolConfigService()

	// 更新配置
	if err := service.UpdateToolConfig(toolCode, req.Config); err != nil {
		log.Error("更新工具配置失败: %v", err)
		common.ServerError(c, "更新配置失败")
		return
	}

	common.SuccessWithMessage(c, "更新配置成功")
}

// UpdateToolSettings 更新工具设置（管理端）
func UpdateToolSettings(c *gin.Context) {
	toolCode := c.Param("code")

	var req request.UpdateToolSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		common.BadRequest(c, "参数错误")
		return
	}

	service := tool_config.NewToolConfigService()

	// 更新设置
	if err := service.UpdateToolSettings(toolCode, req.Enabled, req.Visible, req.SortOrder); err != nil {
		log.Error("更新工具设置失败: %v", err)
		common.ServerError(c, "更新设置失败")
		return
	}

	common.SuccessWithMessage(c, "更新设置成功")
}

// DeleteTool 删除工具配置（管理端）
func DeleteTool(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		common.BadRequest(c, "无效的ID")
		return
	}

	service := tool_config.NewToolConfigService()

	if err := service.DeleteTool(uint(id)); err != nil {
		log.Error("删除工具配置失败: %v", err)
		common.ServerError(c, "删除工具失败")
		return
	}

	common.SuccessWithMessage(c, "删除工具成功")
}

// SyncTools 同步工具定义（管理端）
func SyncTools(c *gin.Context) {
	service := tool_config.NewToolConfigService()

	if err := service.SyncToolsFromRegistry(); err != nil {
		log.Error("同步工具失败: %v", err)
		common.ServerError(c, "同步工具失败")
		return
	}

	common.SuccessWithMessage(c, "同步工具成功")
}
