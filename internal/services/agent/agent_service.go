package agent

import (
	"auto-forge/internal/models"
	"auto-forge/pkg/agent/executor"
	"auto-forge/pkg/agent/llm"
	"auto-forge/pkg/agent/registry"
	"auto-forge/pkg/common"
	"auto-forge/pkg/config"
	"auto-forge/pkg/database"
	log "auto-forge/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type AgentService struct{}

func NewAgentService() *AgentService {
	return &AgentService{}
}

// CreateConversation 创建新对话
func (s *AgentService) CreateConversation(userID, title string) (*models.AgentConversation, error) {
	db := database.GetDB()

	conversation := &models.AgentConversation{
		ID:     common.NewUUID().String(),
		UserID: userID,
		Title:  title,
	}

	if err := db.Create(conversation).Error; err != nil {
		return nil, fmt.Errorf("创建对话失败: %w", err)
	}

	return conversation, nil
}

// GetConversations 获取用户的对话列表
func (s *AgentService) GetConversations(userID string, page, pageSize int) ([]models.AgentConversation, int64, error) {
	db := database.GetDB()

	var conversations []models.AgentConversation
	var total int64

	// 统计总数
	if err := db.Model(&models.AgentConversation{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("统计对话数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := db.Where("user_id = ?", userID).
		Order("updated_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&conversations).Error; err != nil {
		return nil, 0, fmt.Errorf("查询对话列表失败: %w", err)
	}

	return conversations, total, nil
}

// GetConversationByID 根据ID获取对话
func (s *AgentService) GetConversationByID(conversationID, userID string) (*models.AgentConversation, error) {
	db := database.GetDB()

	var conversation models.AgentConversation
	if err := db.Where("id = ? AND user_id = ?", conversationID, userID).First(&conversation).Error; err != nil {
		return nil, fmt.Errorf("对话不存在: %w", err)
	}

	return &conversation, nil
}

// UpdateConversation 更新对话信息
func (s *AgentService) UpdateConversation(conversationID, userID, title string) error {
	db := database.GetDB()

	result := db.Model(&models.AgentConversation{}).
		Where("id = ? AND user_id = ?", conversationID, userID).
		Updates(map[string]interface{}{
			"title":      title,
			"updated_at": time.Now().Unix(),
		})

	if result.Error != nil {
		return fmt.Errorf("更新对话失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("对话不存在")
	}

	return nil
}

// DeleteConversation 删除对话（会级联删除所有消息）
func (s *AgentService) DeleteConversation(conversationID, userID string) error {
	db := database.GetDB()

	// 先删除所有消息
	if err := db.Where("conversation_id = ?", conversationID).Delete(&models.AgentMessage{}).Error; err != nil {
		return fmt.Errorf("删除对话消息失败: %w", err)
	}

	// 再删除对话
	result := db.Where("id = ? AND user_id = ?", conversationID, userID).Delete(&models.AgentConversation{})
	if result.Error != nil {
		return fmt.Errorf("删除对话失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("对话不存在")
	}

	return nil
}

// CreateMessage 创建新消息（用户消息）
func (s *AgentService) CreateMessage(conversationID, role, content string, files []models.AgentFile) (*models.AgentMessage, error) {
	db := database.GetDB()

	message := &models.AgentMessage{
		ID:             common.NewUUID().String(),
		ConversationID: conversationID,
		Role:           role,
		Content:        content,
		Files:          files,
		Status:         "completed", // 用户消息直接完成
	}

	if err := db.Create(message).Error; err != nil {
		return nil, fmt.Errorf("创建消息失败: %w", err)
	}

	// 更新对话的 updated_at
	db.Model(&models.AgentConversation{}).Where("id = ?", conversationID).Update("updated_at", time.Now().Unix())

	return message, nil
}

// CreateAgentMessage 创建 Agent 消息（初始状态为 pending）
func (s *AgentService) CreateAgentMessage(conversationID string, config *models.AgentConfig) (*models.AgentMessage, error) {
	db := database.GetDB()

	// 确保 Agent 消息的时间戳比用户消息晚
	// 获取最新的用户消息时间戳
	var lastUserMessage models.AgentMessage
	db.Where("conversation_id = ? AND role = ?", conversationID, "user").
		Order("created_at DESC").
		First(&lastUserMessage)

	// Agent 消息的时间戳至少比最后一条用户消息晚 1 秒
	agentCreatedAt := time.Now().Unix()
	if lastUserMessage.CreatedAt >= agentCreatedAt {
		agentCreatedAt = lastUserMessage.CreatedAt + 1
	}

	message := &models.AgentMessage{
		ID:             common.NewUUID().String(),
		ConversationID: conversationID,
		Role:           "agent",
		Content:        "",
		Config:         config,
		Status:         "pending",
		CreatedAt:      agentCreatedAt,
	}

	if err := db.Create(message).Error; err != nil {
		return nil, fmt.Errorf("创建 Agent 消息失败: %w", err)
	}

	return message, nil
}

// UpdateMessageStatus 更新消息状态
func (s *AgentService) UpdateMessageStatus(messageID, status string, errorMsg string) error {
	db := database.GetDB()

	updates := map[string]interface{}{
		"status": status,
	}

	if errorMsg != "" {
		updates["error"] = errorMsg
	}

	return db.Model(&models.AgentMessage{}).Where("id = ?", messageID).Updates(updates).Error
}

// UpdateMessagePlan 更新消息的执行计划
func (s *AgentService) UpdateMessagePlan(messageID string, plan *models.AgentPlan) error {
	db := database.GetDB()

	planJSON, err := json.Marshal(plan)
	if err != nil {
		return fmt.Errorf("序列化计划失败: %w", err)
	}

	return db.Model(&models.AgentMessage{}).
		Where("id = ?", messageID).
		Update("plan", planJSON).Error
}

// UpdateMessageTrace 更新消息的执行轨迹
func (s *AgentService) UpdateMessageTrace(messageID string, trace *models.AgentTrace) error {
	db := database.GetDB()

	traceJSON, err := json.Marshal(trace)
	if err != nil {
		return fmt.Errorf("序列化轨迹失败: %w", err)
	}

	return db.Model(&models.AgentMessage{}).
		Where("id = ?", messageID).
		Updates(map[string]interface{}{
			"trace":   traceJSON,
			"content": trace.FinalAnswer,
		}).Error
}

// UpdateMessageTokenUsage 更新消息的 Token 使用情况
func (s *AgentService) UpdateMessageTokenUsage(messageID string, usage *models.TokenUsage) error {
	db := database.GetDB()

	usageJSON, err := json.Marshal(usage)
	if err != nil {
		return fmt.Errorf("序列化 Token 使用情况失败: %w", err)
	}

	return db.Model(&models.AgentMessage{}).
		Where("id = ?", messageID).
		Update("token_usage", usageJSON).Error
}

// GetMessages 获取对话的所有消息
func (s *AgentService) GetMessages(conversationID string) ([]models.AgentMessage, error) {
	db := database.GetDB()

	var messages []models.AgentMessage
	if err := db.Where("conversation_id = ?", conversationID).
		Order("created_at ASC").
		Find(&messages).Error; err != nil {
		return nil, fmt.Errorf("查询消息列表失败: %w", err)
	}

	return messages, nil
}

// GetMessageByID 根据ID获取消息
func (s *AgentService) GetMessageByID(messageID string) (*models.AgentMessage, error) {
	db := database.GetDB()

	var message models.AgentMessage
	if err := db.Where("id = ?", messageID).First(&message).Error; err != nil {
		return nil, fmt.Errorf("消息不存在: %w", err)
	}

	return &message, nil
}

// BuildContextFromMessages 从历史消息构建上下文（用于对话记忆）
func (s *AgentService) BuildContextFromMessages(messages []models.AgentMessage, maxMessages int) string {
	if len(messages) == 0 {
		return ""
	}

	// 只取最近的 N 条消息
	start := 0
	if len(messages) > maxMessages {
		start = len(messages) - maxMessages
	}

	var contextParts []string
	for _, msg := range messages[start:] {
		if msg.Role == "user" {
			contextParts = append(contextParts, fmt.Sprintf("User: %s", msg.Content))
		} else if msg.Role == "agent" && msg.Content != "" {
			contextParts = append(contextParts, fmt.Sprintf("Assistant: %s", msg.Content))
		}
	}

	return strings.Join(contextParts, "\n\n")
}

// ExecuteAgent 执行 Agent（核心方法）
func (s *AgentService) ExecuteAgent(
	ctx context.Context,
	messageID string,
	userMessage string,
	files []models.AgentFile,
	config *models.AgentConfig,
	conversationHistory string,
	streamCallback func(event AgentStreamEvent) error,
) error {
	// 获取配置（使用默认值）
	model := "gpt-4o-mini"
	mode := "direct" // direct / plan
	maxSteps := 10
	temperature := 0.7
	var allowedTools []string

	if config != nil {
		if config.Model != "" {
			model = config.Model
		}
		if config.Mode != "" {
			mode = config.Mode
		}
		if config.MaxSteps > 0 {
			maxSteps = config.MaxSteps
		}
		if config.Temperature > 0 {
			temperature = config.Temperature
		}
		if len(config.AllowedTools) > 0 {
			allowedTools = config.AllowedTools
		}
	}

	log.Info("Agent 开始执行: MessageID=%s, Model=%s, Mode=%s, MaxSteps=%d", messageID, model, mode, maxSteps)

	// 更新状态为 running
	if err := s.UpdateMessageStatus(messageID, "running", ""); err != nil {
		return err
	}

	// 使用新的执行器
	err := s.executeWithNewEngine(ctx, messageID, userMessage, model, mode, maxSteps, temperature, allowedTools, conversationHistory, streamCallback)

	// 更新最终状态
	if err != nil {
		s.UpdateMessageStatus(messageID, "failed", err.Error())
		return err
	}

	s.UpdateMessageStatus(messageID, "completed", "")
	return nil
}

// 辅助函数：从配置中获取字符串
func getStringFromConfig(config map[string]interface{}, key, defaultValue string) string {
	if val, ok := config[key].(string); ok {
		return val
	}
	return defaultValue
}

// 辅助函数：从配置中获取整数
func getIntFromConfig(config map[string]interface{}, key string, defaultValue int) int {
	if val, ok := config[key].(float64); ok {
		return int(val)
	}
	if val, ok := config[key].(int); ok {
		return val
	}
	return defaultValue
}

// 辅助函数：从配置中获取浮点数
func getFloatFromConfig(config map[string]interface{}, key string, defaultValue float64) float64 {
	if val, ok := config[key].(float64); ok {
		return val
	}
	return defaultValue
}

// 辅助函数：从配置中获取字符串数组
func getStringArrayFromConfig(config map[string]interface{}, key string) []string {
	if val, ok := config[key].([]interface{}); ok {
		var result []string
		for _, v := range val {
			if str, ok := v.(string); ok {
				result = append(result, str)
			}
		}
		return result
	}
	return nil
}

// executeWithNewEngine 使用新的执行引擎
func (s *AgentService) executeWithNewEngine(
	ctx context.Context,
	messageID string,
	userMessage string,
	model string,
	mode string,
	maxSteps int,
	temperature float64,
	allowedTools []string,
	conversationHistory string,
	streamCallback func(event AgentStreamEvent) error,
) error {
	// 1. 初始化 LLM 客户端
	cfg := config.GetConfig()
	apiKey := cfg.Agent.OpenAI.APIKey
	baseURL := cfg.Agent.OpenAI.BaseURL

	// 如果配置文件中没有设置，尝试从环境变量读取
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	if baseURL == "" {
		baseURL = os.Getenv("OPENAI_BASE_URL")
	}
	if baseURL == "" {
		baseURL = "https://api.openai.com/v1"
	}

	if apiKey == "" {
		return fmt.Errorf("OpenAI API Key 未配置，请在 config.yaml 中设置 agent.openai.api_key 或设置环境变量 OPENAI_API_KEY")
	}

	llmClient := llm.NewOpenAIClient(model, apiKey, baseURL)

	// 2. 初始化工具注册表
	toolRegistry := registry.NewToolRegistry()
	if err := toolRegistry.RegisterFromUTools(); err != nil {
		return fmt.Errorf("注册工具失败: %w", err)
	}

	log.Info("工具注册完成，共 %d 个工具", len(toolRegistry.ListTools()))

	// 3. 创建执行器回调适配器
	executorCallback := func(event executor.StreamEvent) {
		if streamCallback == nil {
			return
		}

		// 转换事件格式
		agentEvent := AgentStreamEvent{
			Type: event.Type,
			Data: event.Data,
		}

		streamCallback(agentEvent)
	}

	// 4. 根据模式选择执行器
	var result *executor.ExecutionResult
	var err error

	if mode == "plan" {
		// Plan 模式
		planExecutor := executor.NewPlanExecutor(llmClient, toolRegistry, temperature)
		result, err = planExecutor.Execute(ctx, userMessage, conversationHistory, allowedTools, maxSteps, executorCallback)
	} else {
		// ReAct 模式（默认）
		reactExecutor := executor.NewReActExecutor(llmClient, toolRegistry, maxSteps, temperature)
		result, err = reactExecutor.Execute(ctx, userMessage, conversationHistory, allowedTools, executorCallback)
	}

	if err != nil {
		return fmt.Errorf("执行失败: %w", err)
	}

	// 5. 保存执行结果
	if result.Trace != nil {
		// 更新消息内容为最终答案
		if err := s.UpdateMessageContent(messageID, result.Trace.FinalAnswer); err != nil {
			log.Error("更新消息内容失败: %v", err)
		}

		// 保存执行轨迹
		if err := s.UpdateMessageTrace(messageID, result.Trace); err != nil {
			log.Error("保存执行轨迹失败: %v", err)
		}

		// 保存 Token 使用情况
		if result.Trace.TokenUsage != nil {
			if err := s.UpdateMessageTokenUsage(messageID, result.Trace.TokenUsage); err != nil {
				log.Error("保存 Token 使用情况失败: %v", err)
			}
		}
	}

	if !result.Success {
		return fmt.Errorf("执行未成功完成: %s", result.Error)
	}

	return nil
}

// UpdateMessageContent 更新消息内容
func (s *AgentService) UpdateMessageContent(messageID, content string) error {
	db := database.GetDB()
	return db.Model(&models.AgentMessage{}).Where("id = ?", messageID).Update("content", content).Error
}
