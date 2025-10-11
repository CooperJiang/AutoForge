package email

import (
	"auto-forge/pkg/config"
	"auto-forge/pkg/utools"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
	"time"
)

// EmailTool 邮件发送工具
type EmailTool struct {
	*utools.BaseTool
}

// NewEmailTool 创建邮件工具实例
func NewEmailTool() *EmailTool {
	metadata := &utools.ToolMetadata{
		Code:        "email_sender",
		Name:        "邮件发送",
		Description: "发送邮件通知，支持HTML格式和多个收件人",
		Category:    "notification",
		Version:     "1.0.0",
		Author:      "AutoForge",
		Icon:        "mail",
		AICallable:  true,
		Tags:        []string{"email", "notification", "smtp", "alert"},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"to": {
				Type:        "string",
				Title:       "收件人",
				Description: "收件人邮箱，多个用逗号分隔",
			},
			"cc": {
				Type:        "string",
				Title:       "抄送人",
				Description: "抄送人邮箱，多个用逗号分隔（可选）",
			},
			"subject": {
				Type:        "string",
				Title:       "邮件主题",
				Description: "邮件标题",
			},
			"body": {
				Type:        "string",
				Title:       "邮件正文",
				Description: "邮件内容，支持纯文本或HTML",
			},
			"content_type": {
				Type:        "string",
				Title:       "内容类型",
				Description: "邮件内容格式",
				Default:     "text/plain",
				Enum:        []interface{}{"text/plain", "text/html"},
			},
		},
		Required: []string{"to", "subject", "body"},
	}

	return &EmailTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

// Execute 执行邮件发送
func (t *EmailTool) Execute(ctx *utools.ExecutionContext, toolConfig map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()

	// 获取系统邮件配置
	systemConfig := config.GetConfig()
	mailConfig := systemConfig.Mail

	// 检查邮件是否启用
	if !mailConfig.Enabled {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "系统邮件功能未启用，请联系管理员配置",
			Error:      "email service disabled",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("email service disabled")
	}

	// 使用系统配置
	smtpHost := mailConfig.Host
	smtpPort := mailConfig.Port
	username := mailConfig.Username
	password := mailConfig.Password
	from := mailConfig.From
	fromName := mailConfig.FromName
	useTLS := mailConfig.SSL

	// 解析用户输入的配置
	to, _ := toolConfig["to"].(string)
	cc, _ := toolConfig["cc"].(string)
	subject, _ := toolConfig["subject"].(string)
	body, _ := toolConfig["body"].(string)
	contentType, _ := toolConfig["content_type"].(string)
	if contentType == "" {
		contentType = "text/plain"
	}

	// 构建收件人列表
	toList := parseEmailList(to)
	if len(toList) == 0 {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "收件人列表为空",
			Error:      "no recipients",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("no recipients")
	}

	ccList := parseEmailList(cc)

	// 构建邮件内容
	message := buildEmailMessage(from, fromName, toList, ccList, subject, body, contentType)

	// 发送邮件
	addr := fmt.Sprintf("%s:%d", smtpHost, smtpPort)
	auth := smtp.PlainAuth("", username, password, smtpHost)

	var err error
	// 根据端口和配置选择连接方式
	// 465 端口使用 SSL/TLS (SMTPS)
	// 587 端口使用 STARTTLS
	// 25 端口通常不加密或 STARTTLS
	if useTLS && smtpPort == 465 {
		// 465 端口使用直接 TLS 连接
		err = sendMailWithTLS(addr, auth, from, append(toList, ccList...), []byte(message))
	} else if useTLS {
		// 587 或其他端口使用 STARTTLS
		err = sendMailWithSTARTTLS(addr, auth, from, append(toList, ccList...), []byte(message))
	} else {
		// 不加密
		err = smtp.SendMail(addr, auth, from, append(toList, ccList...), []byte(message))
	}

	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "邮件发送失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("邮件发送成功，收件人: %d 人", len(toList)),
		Output: map[string]interface{}{
			"recipients_count": len(toList),
			"cc_count":         len(ccList),
			"subject":          subject,
		},
		DurationMs: time.Since(startTime).Milliseconds(),
	}, nil
}

// parseEmailList 解析邮箱列表（逗号分隔）
func parseEmailList(emails string) []string {
	if emails == "" {
		return []string{}
	}

	list := strings.Split(emails, ",")
	result := make([]string, 0, len(list))

	for _, email := range list {
		email = strings.TrimSpace(email)
		if email != "" {
			result = append(result, email)
		}
	}

	return result
}

// buildEmailMessage 构建邮件内容
func buildEmailMessage(from, fromName string, to, cc []string, subject, body, contentType string) string {
	var msg strings.Builder

	// From
	if fromName != "" {
		msg.WriteString(fmt.Sprintf("From: %s <%s>\r\n", fromName, from))
	} else {
		msg.WriteString(fmt.Sprintf("From: %s\r\n", from))
	}

	// To
	msg.WriteString(fmt.Sprintf("To: %s\r\n", strings.Join(to, ", ")))

	// CC
	if len(cc) > 0 {
		msg.WriteString(fmt.Sprintf("Cc: %s\r\n", strings.Join(cc, ", ")))
	}

	// Subject
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", subject))

	// MIME
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString(fmt.Sprintf("Content-Type: %s; charset=UTF-8\r\n", contentType))
	msg.WriteString("\r\n")

	// Body
	msg.WriteString(body)

	return msg.String()
}

// sendMailWithTLS 使用直接 TLS 连接发送邮件 (适用于 465 端口)
func sendMailWithTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// 分离主机和端口
	host := strings.Split(addr, ":")[0]

	// 创建 TLS 配置
	tlsConfig := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
	}

	// 连接到服务器
	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 创建 SMTP 客户端
	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer client.Quit()

	// 认证
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return err
		}
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return err
	}

	// 设置收件人
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}

// sendMailWithSTARTTLS 使用 STARTTLS 发送邮件 (适用于 587 端口)
func sendMailWithSTARTTLS(addr string, auth smtp.Auth, from string, to []string, msg []byte) error {
	// 分离主机和端口
	host := strings.Split(addr, ":")[0]

	// 先建立普通 TCP 连接
	client, err := smtp.Dial(addr)
	if err != nil {
		return fmt.Errorf("连接 SMTP 服务器失败: %v", err)
	}
	defer client.Quit()

	// 发送 EHLO/HELO
	if err = client.Hello(host); err != nil {
		return fmt.Errorf("EHLO 失败: %v", err)
	}

	// 启动 TLS
	tlsConfig := &tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
	}

	if err = client.StartTLS(tlsConfig); err != nil {
		return fmt.Errorf("STARTTLS 失败: %v", err)
	}

	// 认证
	if auth != nil {
		if err = client.Auth(auth); err != nil {
			return fmt.Errorf("认证失败: %v", err)
		}
	}

	// 设置发件人
	if err = client.Mail(from); err != nil {
		return fmt.Errorf("设置发件人失败: %v", err)
	}

	// 设置收件人
	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return fmt.Errorf("设置收件人失败: %v", err)
		}
	}

	// 发送邮件内容
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("准备发送数据失败: %v", err)
	}

	_, err = w.Write(msg)
	if err != nil {
		return fmt.Errorf("写入邮件内容失败: %v", err)
	}

	err = w.Close()
	if err != nil {
		return fmt.Errorf("完成发送失败: %v", err)
	}

	return nil
}

// init 自动注册工具
func init() {
	tool := NewEmailTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register Email tool: %v", err))
	}
}
