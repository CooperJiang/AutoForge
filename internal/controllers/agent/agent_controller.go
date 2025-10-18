package agent

import (
	"auto-forge/internal/dto/request"
	"auto-forge/internal/models"
	"auto-forge/internal/services/agent"
	"auto-forge/pkg/errors"
	log "auto-forge/pkg/logger"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// CreateConversation 创建对话
func CreateConversation(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req request.CreateAgentConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	agentService := agent.NewAgentService()
	conversation, err := agentService.CreateConversation(userID.(string), req.Title)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, conversation, "创建对话成功")
}

// GetConversations 获取对话列表
func GetConversations(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// 直接从查询参数获取
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		fmt.Sscanf(p, "%d", &page)
	}
	if ps := c.Query("page_size"); ps != "" {
		fmt.Sscanf(ps, "%d", &pageSize)
	}

	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	agentService := agent.NewAgentService()
	conversations, total, err := agentService.GetConversations(userID.(string), page, pageSize)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, gin.H{
		"list":      conversations,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	}, "获取对话列表成功")
}

// GetConversationByID 获取对话详情
func GetConversationByID(c *gin.Context) {
	userID, _ := c.Get("user_id")
	conversationID := c.Param("id")

	agentService := agent.NewAgentService()
	conversation, err := agentService.GetConversationByID(conversationID, userID.(string))
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeNotFound, "对话不存在"))
		return
	}

	errors.ResponseSuccess(c, conversation, "获取对话详情成功")
}

// UpdateConversation 更新对话
func UpdateConversation(c *gin.Context) {
	userID, _ := c.Get("user_id")
	conversationID := c.Param("id")

	var req request.UpdateAgentConversationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
		return
	}

	agentService := agent.NewAgentService()
	if err := agentService.UpdateConversation(conversationID, userID.(string), req.Title); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, gin.H{"message": "更新成功"}, "更新成功")
}

// DeleteConversation 删除对话
func DeleteConversation(c *gin.Context) {
	userID, _ := c.Get("user_id")
	conversationID := c.Param("id")

	agentService := agent.NewAgentService()
	if err := agentService.DeleteConversation(conversationID, userID.(string)); err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, gin.H{"message": "删除成功"}, "删除成功")
}

// GetMessages 获取对话消息列表
func GetMessages(c *gin.Context) {
	conversationID := c.Param("id")

	agentService := agent.NewAgentService()
	messages, err := agentService.GetMessages(conversationID)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, err.Error()))
		return
	}

	errors.ResponseSuccess(c, messages, "获取消息列表成功")
}

// SendMessage 发送消息（支持流式响应）
func SendMessage(c *gin.Context) {
	_, _ = c.Get("user_id")
	conversationID := c.Param("id")

	// 解析请求（支持 JSON 和 multipart/form-data）
	var userMessage string
	var files []models.AgentFile
	var config *models.AgentConfig

	contentType := c.GetHeader("Content-Type")

	if strings.Contains(contentType, "multipart/form-data") {
		// 处理文件上传
		userMessage = c.PostForm("message")
		configStr := c.PostForm("config")
		if configStr != "" {
			config = &models.AgentConfig{}
			json.Unmarshal([]byte(configStr), config)
		}

		// TODO: 处理文件上传
		// form, _ := c.MultipartForm()
		// fileHeaders := form.File["files"]
		// ...

	} else {
		// JSON 请求
		var req request.SendAgentMessageRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "参数错误: "+err.Error()))
			return
		}
		userMessage = req.Message
		if req.Config != nil {
			config = &models.AgentConfig{
				Mode:         req.Config.Mode,
				Model:        req.Config.Model,
				MaxSteps:     req.Config.MaxSteps,
				Temperature:  req.Config.Temperature,
				AllowedTools: req.Config.AllowedTools,
			}
		}
	}

	if userMessage == "" {
		errors.HandleError(c, errors.New(errors.CodeInvalidParameter, "消息内容不能为空"))
		return
	}

	agentService := agent.NewAgentService()

	// 1. 创建用户消息
	userMsg, err := agentService.CreateMessage(conversationID, "user", userMessage, files)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "创建用户消息失败: "+err.Error()))
		return
	}

	// 2. 创建 Agent 消息
	agentMsg, err := agentService.CreateAgentMessage(conversationID, config)
	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "创建 Agent 消息失败: "+err.Error()))
		return
	}

	// 3. 获取对话历史
	messages, _ := agentService.GetMessages(conversationID)
	conversationHistory := agentService.BuildContextFromMessages(messages, 10)

	// 4. 检查是否请求流式响应
	acceptHeader := c.GetHeader("Accept")
	if strings.Contains(acceptHeader, "text/event-stream") {
		// 流式响应
		handleStreamResponse(c, agentService, agentMsg.ID, userMessage, files, config, conversationHistory)
	} else {
		// 普通响应
		handleNormalResponse(c, agentService, agentMsg.ID, userMessage, files, config, conversationHistory, userMsg, agentMsg)
	}
}

// handleStreamResponse 处理流式响应（SSE）
func handleStreamResponse(
	c *gin.Context,
	agentService *agent.AgentService,
	messageID string,
	userMessage string,
	files []models.AgentFile,
	config *models.AgentConfig,
	conversationHistory string,
) {
	// 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		errors.HandleError(c, errors.New(errors.CodeInternal, "不支持流式响应"))
		return
	}

	// 流式回调
	streamCallback := func(event agent.AgentStreamEvent) error {
		// 发送 SSE 事件
		eventJSON, err := json.Marshal(event)
		if err != nil {
			return err
		}

		fmt.Fprintf(c.Writer, "data: %s\n\n", eventJSON)
		flusher.Flush()
		return nil
	}

	// 异步执行 Agent
	ctx := context.Background()
	err := agentService.ExecuteAgent(ctx, messageID, userMessage, files, config, conversationHistory, streamCallback)

	if err != nil {
		// 发送错误事件
		errorEvent := agent.AgentStreamEvent{
			Type: "error",
			Data: agent.ErrorEvent{
				Error:   err.Error(),
				Partial: true,
			},
		}
		eventJSON, _ := json.Marshal(errorEvent)
		fmt.Fprintf(c.Writer, "data: %s\n\n", eventJSON)
		flusher.Flush()
	}

	// 发送完成标记
	fmt.Fprintf(c.Writer, "event: done\ndata: {}\n\n")
	flusher.Flush()
}

// handleNormalResponse 处理普通响应
func handleNormalResponse(
	c *gin.Context,
	agentService *agent.AgentService,
	messageID string,
	userMessage string,
	files []models.AgentFile,
	config *models.AgentConfig,
	conversationHistory string,
	userMsg *models.AgentMessage,
	agentMsg *models.AgentMessage,
) {
	// 同步执行 Agent
	ctx := context.Background()
	err := agentService.ExecuteAgent(ctx, messageID, userMessage, files, config, conversationHistory, nil)

	if err != nil {
		errors.HandleError(c, errors.New(errors.CodeInternal, "Agent 执行失败: "+err.Error()))
		return
	}

	// 重新获取 Agent 消息（包含执行结果）
	finalAgentMsg, err := agentService.GetMessageByID(messageID)
	if err != nil {
		log.Error("获取 Agent 消息失败: %v", err)
	}

	errors.ResponseSuccess(c, gin.H{
		"user_message":  userMsg,
		"agent_message": finalAgentMsg,
	}, "发送消息成功")
}
