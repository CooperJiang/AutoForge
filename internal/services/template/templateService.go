package template

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/dto/response"
	"auto-forge/internal/models"
	"auto-forge/pkg/database"
	log "auto-forge/pkg/logger"
	"auto-forge/pkg/utils"
	"errors"
	"fmt"
	"math"
	"time"

	"gorm.io/gorm"
)

type TemplateService struct{}

func NewTemplateService() *TemplateService {
	return &TemplateService{}
}

// CreateTemplate 从现有工作流创建模板
func (s *TemplateService) CreateTemplate(userID string, req *request.CreateTemplateRequest) (*models.WorkflowTemplate, error) {
	db := database.GetDB()

	// 验证工作流是否存在且属于该用户
	var workflow models.Workflow
	if err := db.Where("id = ? AND user_id = ?", req.WorkflowID, userID).First(&workflow).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("工作流不存在")
		}
		return nil, err
	}

	// 构建模板数据
	templateData := models.TemplateData{
		Nodes:   workflow.Nodes,
		Edges:   workflow.Edges,
		EnvVars: workflow.EnvVars,
	}

	// 提取工作流使用的工具
	requiredTools := s.extractRequiredTools(workflow.Nodes)
	if len(req.RequiredTools) > 0 {
		requiredTools = req.RequiredTools
	}

	template := &models.WorkflowTemplate{
		Name:          req.Name,
		Description:   req.Description,
		Category:      req.Category,
		TemplateData:  templateData,
		CoverImage:    req.CoverImage,
		Icon:          req.Icon,
		RequiredTools: requiredTools,
		Status:        "published",
		IsOfficial:    true,
		IsFeatured:    req.IsFeatured,
		InstallCount:  0,
		ViewCount:     0,
		AuthorID:      userID,
		AuthorName:    "Admin",
		UsageGuide:    req.UsageGuide,
	}

	if err := db.Create(template).Error; err != nil {
		log.Error("创建模板失败: %v", err)
		return nil, err
	}

	log.Info("用户 %s 创建模板: %s (ID: %s)", userID, template.Name, template.ID)
	return template, nil
}

// GetTemplateList 获取模板列表
func (s *TemplateService) GetTemplateList(req *request.ListTemplatesRequest) (*response.TemplateListResponse, error) {
	db := database.GetDB()

	// 设置默认分页参数
	if req.Page == 0 {
		req.Page = 1
	}
	if req.PageSize == 0 {
		req.PageSize = 20
	}

	// 构建查询
	queryDB := db.Model(&models.WorkflowTemplate{}).Where("status = ?", "published")

	// 分类筛选
	if req.Category != "" {
		queryDB = queryDB.Where("category = ?", req.Category)
	}

	// 精选筛选
	if req.IsFeatured != nil {
		queryDB = queryDB.Where("is_featured = ?", *req.IsFeatured)
	}

	// 搜索
	if req.Search != "" {
		searchPattern := "%" + req.Search + "%"
		queryDB = queryDB.Where("name LIKE ? OR description LIKE ?", searchPattern, searchPattern)
	}

	// 获取总数
	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, err
	}

	// 查询数据
	var templates []models.WorkflowTemplate
	offset := (req.Page - 1) * req.PageSize
	if err := queryDB.Order("is_featured DESC, install_count DESC, created_at DESC").
		Offset(offset).
		Limit(req.PageSize).
		Find(&templates).Error; err != nil {
		return nil, err
	}

	// 转换为响应格式
	items := make([]response.TemplateBasicInfo, len(templates))
	for i, tpl := range templates {
		items[i] = response.TemplateBasicInfo{
			ID:            tpl.ID,
			Name:          tpl.Name,
			Description:   tpl.Description,
			Category:      tpl.Category,
			CoverImage:    tpl.CoverImage,
			Icon:          tpl.Icon,
			InstallCount:  tpl.InstallCount,
			ViewCount:     tpl.ViewCount,
			IsOfficial:    tpl.IsOfficial,
			IsFeatured:    tpl.IsFeatured,
			AuthorName:    tpl.AuthorName,
			RequiredTools: tpl.RequiredTools,
			CreatedAt:     tpl.CreatedAt,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(req.PageSize)))

	return &response.TemplateListResponse{
		Items:      items,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// GetTemplateDetail 获取模板详情
func (s *TemplateService) GetTemplateDetail(templateID string) (*response.TemplateDetailResponse, error) {
	db := database.GetDB()

	var template models.WorkflowTemplate
	if err := db.Where("id = ? AND status = ?", templateID, "published").First(&template).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("模板不存在")
		}
		return nil, err
	}

	// 增加浏览量
	go func() {
		db.Model(&models.WorkflowTemplate{}).Where("id = ?", templateID).
			Update("view_count", gorm.Expr("view_count + 1"))
	}()

	return &response.TemplateDetailResponse{
		ID:            template.ID,
		Name:          template.Name,
		Description:   template.Description,
		Category:      template.Category,
		CoverImage:    template.CoverImage,
		Icon:          template.Icon,
		InstallCount:  template.InstallCount,
		ViewCount:     template.ViewCount,
		IsOfficial:    template.IsOfficial,
		IsFeatured:    template.IsFeatured,
		AuthorName:    template.AuthorName,
		RequiredTools: template.RequiredTools,
		UsageGuide:    template.UsageGuide,
		TemplateData:  template.TemplateData,
		Status:        template.Status,
		CreatedAt:     template.CreatedAt,
		UpdatedAt:     template.UpdatedAt,
	}, nil
}

// InstallTemplate 安装模板到用户工作流
func (s *TemplateService) InstallTemplate(userID string, req *request.InstallTemplateRequest) (*response.InstallTemplateResponse, error) {
	db := database.GetDB()

	// 获取模板
	var template models.WorkflowTemplate
	if err := db.Where("id = ? AND status = ?", req.TemplateID, "published").First(&template).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("模板不存在")
		}
		return nil, err
	}

	// 准备工作流名称
	workflowName := req.WorkflowName
	if workflowName == "" {
		workflowName = template.Name + " - " + time.Now().Format("2006-01-02 15:04:05")
	}

	// 应用环境变量
	envVars := template.TemplateData.EnvVars
	if len(req.EnvVars) > 0 {
		for i := range envVars {
			if value, ok := req.EnvVars[envVars[i].Key]; ok {
				envVars[i].Value = value
			}
		}
	}

	// 创建工作流
	workflow := &models.Workflow{
		UserID:        userID,
		Name:          workflowName,
		Description:   template.Description,
		Nodes:         models.WorkflowNodes(template.TemplateData.Nodes),
		Edges:         models.WorkflowEdges(template.TemplateData.Edges),
		EnvVars:       models.WorkflowEnvVars(envVars),
		ScheduleType:  "manual",
		ScheduleValue: "",
		Enabled:       false,
	}

	// 生成 API Key
	if apiKey, err := utils.GenerateWorkflowAPIKey(); err == nil {
		workflow.APIKey = apiKey
	}

	if err := db.Create(workflow).Error; err != nil {
		log.Error("安装模板失败: %v", err)
		return nil, fmt.Errorf("创建工作流失败")
	}

	// 记录安装
	install := &models.TemplateInstall{
		TemplateID: req.TemplateID,
		UserID:     userID,
		WorkflowID: workflow.GetID(),
	}

	if err := db.Create(install).Error; err != nil {
		log.Error("记录模板安装失败: %v", err)
	}

	// 增加安装计数
	go func() {
		db.Model(&models.WorkflowTemplate{}).Where("id = ?", req.TemplateID).
			Update("install_count", gorm.Expr("install_count + 1"))
	}()

	log.Info("用户 %s 安装模板: %s (WorkflowID: %s)", userID, template.Name, workflow.GetID())

	return &response.InstallTemplateResponse{
		WorkflowID:   workflow.GetID(),
		WorkflowName: workflow.Name,
		Message:      "模板安装成功",
	}, nil
}

// UpdateTemplate 更新模板
func (s *TemplateService) UpdateTemplate(templateID, userID string, req *request.UpdateTemplateRequest) (*models.WorkflowTemplate, error) {
	db := database.GetDB()

	// 查找模板
	var template models.WorkflowTemplate
	if err := db.Where("id = ?", templateID).First(&template).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("模板不存在")
		}
		return nil, err
	}

	// 构建更新
	updates := make(map[string]interface{})

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Description != "" {
		updates["description"] = req.Description
	}
	if req.Category != "" {
		updates["category"] = req.Category
	}
	if req.CoverImage != "" {
		updates["cover_image"] = req.CoverImage
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	if req.UsageGuide != "" {
		updates["usage_guide"] = req.UsageGuide
	}
	if req.Status != "" {
		updates["status"] = req.Status
	}

	updates["is_featured"] = req.IsFeatured

	if err := db.Model(&template).Updates(updates).Error; err != nil {
		return nil, err
	}

	// 重新查询更新后的模板
	if err := db.First(&template, "id = ?", templateID).Error; err != nil {
		return nil, err
	}

	log.Info("用户 %s 更新模板: %s (ID: %s)", userID, template.Name, template.ID)
	return &template, nil
}

// DeleteTemplate 删除模板
func (s *TemplateService) DeleteTemplate(templateID, userID string) error {
	db := database.GetDB()

	var template models.WorkflowTemplate
	if err := db.Where("id = ?", templateID).First(&template).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("模板不存在")
		}
		return err
	}

	if err := db.Delete(&template).Error; err != nil {
		return err
	}

	log.Info("用户 %s 删除模板: %s (ID: %s)", userID, template.Name, template.ID)
	return nil
}

// GetUserInstallHistory 获取用户安装历史
func (s *TemplateService) GetUserInstallHistory(userID string) ([]models.TemplateInstall, error) {
	db := database.GetDB()

	var installs []models.TemplateInstall
	if err := db.Where("user_id = ?", userID).
		Order("installed_at DESC").
		Find(&installs).Error; err != nil {
		return nil, err
	}

	return installs, nil
}

// extractRequiredTools 从节点中提取所需工具
func (s *TemplateService) extractRequiredTools(nodes []models.WorkflowNode) []string {
	toolSet := make(map[string]bool)
	tools := []string{}

	for _, node := range nodes {
		if node.Type == "tool" && node.ToolCode != "" {
			if !toolSet[node.ToolCode] {
				toolSet[node.ToolCode] = true
				tools = append(tools, node.ToolCode)
			}
		}
	}

	return tools
}
