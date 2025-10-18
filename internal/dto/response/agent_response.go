package response

import "auto-forge/internal/models"

// SendAgentMessageResponse 发送消息响应
type SendAgentMessageResponse struct {
	UserMessage  *models.AgentMessage `json:"user_message"`
	AgentMessage *models.AgentMessage `json:"agent_message"`
}


