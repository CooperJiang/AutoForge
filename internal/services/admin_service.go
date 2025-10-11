package services

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"auto-forge/internal/cron"
	"auto-forge/internal/models"
	"auto-forge/pkg/config"
	"auto-forge/pkg/database"
	"auto-forge/pkg/errors"
)

// AdminSession 管理员会话
type AdminSession struct {
	Token     string
	ExpiresAt time.Time
}

// AdminService 管理员服务
type AdminService struct {
	sessions sync.Map // token -> AdminSession
	cfg      *config.Config
}

var (
	adminServiceInstance *AdminService
	adminServiceOnce     sync.Once
)

// NewAdminService 创建管理员服务实例（单例）
func NewAdminService() *AdminService {
	adminServiceOnce.Do(func() {
		adminServiceInstance = &AdminService{
			cfg: config.GetConfig(),
		}
	})
	return adminServiceInstance
}

// Login 管理员登录
func (s *AdminService) Login(password string) (string, error) {
	// 验证密码
	if password != s.cfg.Admin.Password {
		return "", errors.New(errors.CodeUnauthorized, "密码错误")
	}

	// 生成 token
	token, err := s.generateToken()
	if err != nil {
		return "", err
	}

	// 计算过期时间
	expiresAt := time.Now().Add(time.Duration(s.cfg.Admin.SessionExpires) * time.Second)

	// 存储会话
	s.sessions.Store(token, AdminSession{
		Token:     token,
		ExpiresAt: expiresAt,
	})

	return token, nil
}

// ValidateToken 验证 token
func (s *AdminService) ValidateToken(token string) bool {
	value, ok := s.sessions.Load(token)
	if !ok {
		return false
	}

	session := value.(AdminSession)

	// 检查是否过期
	if time.Now().After(session.ExpiresAt) {
		s.sessions.Delete(token)
		return false
	}

	return true
}

// Logout 登出
func (s *AdminService) Logout(token string) {
	s.sessions.Delete(token)
}

// generateToken 生成随机 token
func (s *AdminService) generateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GetAllTasks 获取所有任务
func (s *AdminService) GetAllTasks(page, pageSize int, userID, status, keyword string) ([]models.Task, int64, error) {
	db := database.GetDB()

	// 构建查询
	query := db.Model(&models.Task{})

	// 筛选条件
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if status == "enabled" {
		query = query.Where("enabled = ?", true)
	} else if status == "disabled" {
		query = query.Where("enabled = ?", false)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var tasks []models.Task
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

// UpdateTaskStatus 更新任务状态
func (s *AdminService) UpdateTaskStatus(taskID string, enabled bool) error {
	db := database.GetDB()

	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return errors.New(errors.CodeNotFound, "任务不存在")
	}

	task.Enabled = enabled
	if err := db.Save(&task).Error; err != nil {
		return err
	}

	// 重新加载所有任务（会自动处理启用/禁用状态）
	scheduler := cron.GetTaskScheduler()
	if scheduler != nil {
		scheduler.ReloadTasks()
	}

	return nil
}

// DeleteTask 删除任务
func (s *AdminService) DeleteTask(taskID string) error {
	db := database.GetDB()

	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return errors.New(errors.CodeNotFound, "任务不存在")
	}

	// 删除任务执行记录
	if err := db.Where("task_id = ?", taskID).Delete(&models.TaskExecution{}).Error; err != nil {
		return err
	}

	// 删除任务
	if err := db.Delete(&task).Error; err != nil {
		return err
	}

	// 重新加载所有任务
	scheduler := cron.GetTaskScheduler()
	if scheduler != nil {
		scheduler.ReloadTasks()
	}

	return nil
}

// GetStats 获取统计数据
func (s *AdminService) GetStats() (map[string]interface{}, error) {
	db := database.GetDB()

	stats := make(map[string]interface{})

	// 总用户数
	var totalUsers int64
	db.Model(&models.Task{}).Distinct("user_id").Count(&totalUsers)
	stats["total_users"] = totalUsers

	// 总任务数
	var totalTasks int64
	db.Model(&models.Task{}).Count(&totalTasks)
	stats["total_tasks"] = totalTasks

	// 启用的任务数
	var activeTasks int64
	db.Model(&models.Task{}).Where("enabled = ?", true).Count(&activeTasks)
	stats["active_tasks"] = activeTasks

	// 今日执行次数
	todayStart := time.Now().Truncate(24 * time.Hour).Unix()
	var todayExecutions int64
	db.Model(&models.TaskExecution{}).Where("started_at >= ?", todayStart).Count(&todayExecutions)
	stats["today_executions"] = todayExecutions

	// 成功率（最近100次执行）
	var successCount int64
	var totalCount int64
	db.Model(&models.TaskExecution{}).Order("id DESC").Limit(100).Count(&totalCount)
	db.Model(&models.TaskExecution{}).Where("status = ?", "success").Order("id DESC").Limit(100).Count(&successCount)

	successRate := 0.0
	if totalCount > 0 {
		successRate = float64(successCount) / float64(totalCount) * 100
	}
	stats["success_rate"] = successRate

	// 最近活跃用户
	type UserActivity struct {
		UserID      string
		TaskCount   int64
		LastActive  time.Time
	}

	var recentUsers []UserActivity
	db.Raw(`
		SELECT user_id, COUNT(*) as task_count, MAX(updated_at) as last_active
		FROM tasks
		GROUP BY user_id
		ORDER BY last_active DESC
		LIMIT 5
	`).Scan(&recentUsers)

	stats["recent_users"] = recentUsers

	return stats, nil
}

// ExecuteTask 立即执行任务
func (s *AdminService) ExecuteTask(taskID string) error {
	db := database.GetDB()

	var task models.Task
	if err := db.Where("id = ?", taskID).First(&task).Error; err != nil {
		return errors.New(errors.CodeNotFound, "任务不存在")
	}

	// 使用 TaskScheduler 异步执行任务
	scheduler := cron.GetTaskScheduler()
	if scheduler == nil {
		return errors.New(errors.CodeInternal, "任务调度器未初始化")
	}

	scheduler.ExecuteTaskNow(&task)
	return nil
}

// GetAllExecutions 获取所有执行记录
func (s *AdminService) GetAllExecutions(page, pageSize int, userID, taskID, status string) ([]models.TaskExecution, int64, error) {
	db := database.GetDB()

	// 构建查询
	query := db.Model(&models.TaskExecution{})

	// 筛选条件
	if userID != "" {
		query = query.Where("user_id = ?", userID)
	}
	if taskID != "" {
		query = query.Where("task_id = ?", taskID)
	}
	if status == "success" {
		query = query.Where("status = ?", "success")
	} else if status == "failed" {
		query = query.Where("status != ?", "success")
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，预加载任务信息
	var executions []models.TaskExecution
	offset := (page - 1) * pageSize
	if err := query.Preload("Task").Order("started_at DESC").Offset(offset).Limit(pageSize).Find(&executions).Error; err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// UserWithStats 用户信息及任务统计
type UserWithStats struct {
	models.User
	TotalTasks   int64 `json:"total_tasks"`
	EnabledTasks int64 `json:"enabled_tasks"`
}

// GetUsers 获取用户列表
func (s *AdminService) GetUsers(page, pageSize int, keyword string, status int) ([]UserWithStats, int64, error) {
	db := database.GetDB()
	if db == nil {
		return nil, 0, errors.New(errors.CodeDBConnectionFailed, "数据库连接失败")
	}

	query := db.Model(&models.User{})

	// 关键词搜索（用户名或邮箱）
	if keyword != "" {
		query = query.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 状态筛选
	if status > 0 {
		query = query.Where("status = ?", status)
	}

	// 统计总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	var users []models.User
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 为每个用户统计任务数
	usersWithStats := make([]UserWithStats, 0, len(users))
	for _, user := range users {
		userStats := UserWithStats{User: user}

		// 统计总任务数
		db.Model(&models.Task{}).Where("user_id = ?", user.ID.String()).Count(&userStats.TotalTasks)

		// 统计启用的任务数
		db.Model(&models.Task{}).Where("user_id = ? AND enabled = ?", user.ID.String(), true).Count(&userStats.EnabledTasks)

		usersWithStats = append(usersWithStats, userStats)
	}

	return usersWithStats, total, nil
}

// UpdateUserStatus 更新用户状态（启用/禁用）
func (s *AdminService) UpdateUserStatus(userID string, status int) error {
	db := database.GetDB()
	if db == nil {
		return errors.New(errors.CodeDBConnectionFailed, "数据库连接失败")
	}

	// 查找用户
	var user models.User
	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		return errors.New(errors.CodeUserNotFound, "用户不存在")
	}

	// 不能禁用超级管理员（role = 1）
	if user.Role == 1 && status == 2 {
		return errors.New(errors.CodeForbidden, "不能禁用超级管理员")
	}

	// 更新状态
	if err := db.Model(&user).Update("status", status).Error; err != nil {
		return errors.New(errors.CodeQueryFailed, "更新用户状态失败")
	}

	return nil
}
