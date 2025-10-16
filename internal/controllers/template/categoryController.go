package template

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/dto/response"
	"auto-forge/internal/models"
	templateService "auto-forge/internal/services/template"
	"auto-forge/pkg/errors"
	log "auto-forge/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
)

var categoryService = templateService.NewCategoryService()

// toCategoryResponse 转换为响应格式
func toCategoryResponse(category *models.TemplateCategory) response.TemplateCategoryResponse {
	return response.TemplateCategoryResponse{
		ID:          category.GetID(),
		Name:        category.Name,
		Description: category.Description,
		SortOrder:   category.SortOrder,
		IsActive:    category.IsActive,
		CreatedAt:   category.GetCreatedAt().Unix(),
		UpdatedAt:   category.GetUpdatedAt().Unix(),
	}
}

// CreateCategory 创建模板分类
func CreateCategory(c *gin.Context) {
	var req request.CreateTemplateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	category, err := categoryService.CreateCategory(&req)
	if err != nil {
		log.Error("创建模板分类失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "创建分类失败: "+err.Error()))
		return
	}

	resp := toCategoryResponse(category)
	errors.ResponseSuccess(c, resp, "创建分类成功")
}

// GetCategoryList 获取模板分类列表
func GetCategoryList(c *gin.Context) {
	var query request.TemplateCategoryListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	result, err := categoryService.GetCategoryList(&query)
	if err != nil {
		log.Error("获取模板分类列表失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "获取分类列表失败"))
		return
	}

	errors.ResponseSuccess(c, result, "获取分类列表成功")
}

// GetCategoryDetail 获取模板分类详情
func GetCategoryDetail(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "分类ID不能为空"))
		return
	}

	category, err := categoryService.GetCategoryByID(categoryID)
	if err != nil {
		log.Error("获取模板分类失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeNotFound, "获取分类失败: "+err.Error()))
		return
	}

	resp := toCategoryResponse(category)
	errors.ResponseSuccess(c, resp, "获取分类详情成功")
}

// UpdateCategory 更新模板分类
func UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "分类ID不能为空"))
		return
	}

	var req request.UpdateTemplateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	category, err := categoryService.UpdateCategory(categoryID, &req)
	if err != nil {
		log.Error("更新模板分类失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "更新分类失败: "+err.Error()))
		return
	}

	resp := toCategoryResponse(category)
	errors.ResponseSuccess(c, resp, "更新分类成功")
}

// DeleteCategory 删除模板分类
func DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")
	if categoryID == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "分类ID不能为空"))
		return
	}

	if err := categoryService.DeleteCategory(categoryID); err != nil {
		log.Error("删除模板分类失败: %v", err)
		errors.HandleError(c, errors.New(errors.CodeInternal, "删除分类失败: "+err.Error()))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "删除成功",
	})
}
