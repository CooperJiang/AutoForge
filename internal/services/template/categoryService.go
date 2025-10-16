package template

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/dto/response"
	"auto-forge/internal/models"
	"auto-forge/pkg/database"
	log "auto-forge/pkg/logger"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type CategoryService struct{}

func NewCategoryService() *CategoryService {
	return &CategoryService{}
}

// CreateCategory 创建模板分类
func (s *CategoryService) CreateCategory(req *request.CreateTemplateCategoryRequest) (*models.TemplateCategory, error) {
	db := database.GetDB()

	// 检查分类名称是否已存在
	var count int64
	if err := db.Model(&models.TemplateCategory{}).Where("name = ?", req.Name).Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, fmt.Errorf("分类名称已存在")
	}

	// 设置默认排序值
	sortOrder := req.SortOrder
	if sortOrder == 0 {
		sortOrder = 100
	}

	category := &models.TemplateCategory{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   sortOrder,
		IsActive:    true,
	}

	if err := db.Create(category).Error; err != nil {
		log.Error("创建模板分类失败: %v", err)
		return nil, err
	}

	log.Info("创建模板分类: %s (ID: %s)", category.Name, category.ID)
	return category, nil
}

// GetCategoryList 获取模板分类列表
func (s *CategoryService) GetCategoryList(query *request.TemplateCategoryListQuery) (*response.TemplateCategoryListResponse, error) {
	db := database.GetDB()

	page := query.Page
	if page < 1 {
		page = 1
	}
	pageSize := query.PageSize
	if pageSize < 1 {
		pageSize = 20
	}

	// 构建查询
	dbQuery := db.Model(&models.TemplateCategory{})

	if query.IsActive != nil {
		dbQuery = dbQuery.Where("is_active = ?", *query.IsActive)
	}

	// 获取总数
	var total int64
	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, err
	}

	// 获取列表（按排序字段排序）
	var categories []models.TemplateCategory
	offset := (page - 1) * pageSize
	if err := dbQuery.Order("sort_order ASC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&categories).Error; err != nil {
		return nil, err
	}

	// 转换响应
	items := make([]response.TemplateCategoryResponse, len(categories))
	for i, cat := range categories {
		items[i] = response.TemplateCategoryResponse{
			ID:          cat.GetID(),
			Name:        cat.Name,
			Description: cat.Description,
			SortOrder:   cat.SortOrder,
			IsActive:    cat.IsActive,
			CreatedAt:   cat.GetCreatedAt().Unix(),
			UpdatedAt:   cat.GetUpdatedAt().Unix(),
		}
	}

	return &response.TemplateCategoryListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// GetCategoryByID 根据ID获取分类
func (s *CategoryService) GetCategoryByID(categoryID string) (*models.TemplateCategory, error) {
	db := database.GetDB()

	var category models.TemplateCategory
	if err := db.Where("id = ?", categoryID).First(&category).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("分类不存在")
		}
		return nil, err
	}

	return &category, nil
}

// UpdateCategory 更新模板分类
func (s *CategoryService) UpdateCategory(categoryID string, req *request.UpdateTemplateCategoryRequest) (*models.TemplateCategory, error) {
	db := database.GetDB()

	category, err := s.GetCategoryByID(categoryID)
	if err != nil {
		return nil, err
	}

	// 如果更新名称，检查是否与其他分类重名
	if req.Name != nil && *req.Name != category.Name {
		var count int64
		if err := db.Model(&models.TemplateCategory{}).
			Where("name = ? AND id != ?", *req.Name, categoryID).
			Count(&count).Error; err != nil {
			return nil, err
		}
		if count > 0 {
			return nil, fmt.Errorf("分类名称已存在")
		}
	}

	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	if err := db.Model(category).Updates(updates).Error; err != nil {
		return nil, err
	}

	if err := db.First(category, "id = ?", categoryID).Error; err != nil {
		return nil, err
	}

	log.Info("更新模板分类: %s (ID: %s)", category.Name, category.ID)
	return category, nil
}

// DeleteCategory 删除模板分类
func (s *CategoryService) DeleteCategory(categoryID string) error {
	db := database.GetDB()

	category, err := s.GetCategoryByID(categoryID)
	if err != nil {
		return err
	}

	// 检查是否有模板使用该分类
	var templateCount int64
	if err := db.Model(&models.WorkflowTemplate{}).
		Where("category = ?", category.Name).
		Count(&templateCount).Error; err != nil {
		return err
	}

	if templateCount > 0 {
		return fmt.Errorf("该分类下有 %d 个模板正在使用，无法删除", templateCount)
	}

	if err := db.Delete(category).Error; err != nil {
		log.Error("删除模板分类失败: %v", err)
		return err
	}

	log.Info("删除模板分类: %s (ID: %s)", category.Name, category.ID)
	return nil
}
