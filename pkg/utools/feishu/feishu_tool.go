package feishu

import (
	"auto-forge/pkg/utools"
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"go.uber.org/zap"
)


type FeishuTool struct {
	*utools.BaseTool
	logger *zap.Logger
}


func NewFeishuTool() *FeishuTool {
    metadata := &utools.ToolMetadata{
        Code:        "feishu_bot",
        Name:        "飞书机器人",
        Description: "通过飞书机器人 Webhook 发送消息通知，支持文本、富文本、图片和卡片消息",
        Category:    utools.CategoryNotification,
        Version:     "1.0.0",
        Author:      "AutoForge",
        AICallable:  true,
        Tags:        []string{"feishu", "lark", "notification", "bot", "webhook"},
        OutputFieldsSchema: map[string]utools.OutputFieldDef{
            "success": {Type: "boolean", Label: "是否发送成功"},
            "message": {Type: "string", Label: "执行消息"},
            "response": {Type: "object", Label: "飞书返回原始响应"},
        },
    }

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"webhook_url": {
				Type:        "string",
				Title:       "Webhook 地址",
				Description: "飞书机器人的 Webhook URL，格式：https://...",
			},
			"sign_secret": {
				Type:        "string",
				Title:       "签名密钥",
				Description: "机器人的签名密钥（可选，用于安全验证），留空表示不使用签名验证",
				Secret:      true,
			},
			"app_id": {
				Type:        "string",
				Title:       "应用 ID",
				Description: "飞书应用的 App ID（可选，用于上传图片），如果要发送图片消息必须填写",
			},
			"app_secret": {
				Type:        "string",
				Title:       "应用密钥",
				Description: "飞书应用的 App Secret（可选，用于上传图片），如果要发送图片消息必须填写",
				Secret:      true,
			},
			"msg_type": {
				Type:        "string",
				Title:       "消息类型",
				Description: "选择要发送的消息类型",
				Default:     "text",
				Enum:        []interface{}{"text", "post", "image", "interactive"},
			},
			"content": {
				Type:        "string",
				Title:       "消息内容",
				Description: "文本消息的内容（当消息类型为 text 时使用）",
			},
			"title": {
				Type:        "string",
				Title:       "标题",
				Description: "富文本或卡片消息的标题",
			},
			"post_content": {
				Type:        "string",
				Title:       "富文本内容",
				Description: "富文本消息的 content 数组（JSON格式），支持 text/a/at/img/md/hr/code_block 等标签。示例：[[{\"tag\":\"text\",\"text\":\"内容\"}],[{\"tag\":\"a\",\"href\":\"http://xxx\",\"text\":\"链接\"}]]",
			},
			"image_url": {
				Type:        "string",
				Title:       "图片 URL",
				Description: "图片的公网访问地址（当消息类型为 image 时使用），格式：https://...",
			},
			"card_template": {
				Type:        "string",
				Title:       "卡片模板",
				Description: "选择预设的卡片模板（当消息类型为 interactive 时使用）",
				Default:     "notification",
				Enum:        []interface{}{"notification", "alert", "report", "custom"},
			},
			"card_content": {
				Type:        "string",
				Title:       "卡片内容",
				Description: "卡片的主要内容描述",
			},
			"card_status": {
				Type:        "string",
				Title:       "卡片状态",
				Description: "卡片的状态标识",
				Default:     "info",
				Enum:        []interface{}{"success", "warning", "error", "info"},
			},
			"card_fields": {
				Type:        "string",
				Title:       "卡片字段",
				Description: "卡片字段列表，JSON 格式数组，示例：[{\"key\":\"任务名称\",\"value\":\"数据同步\"}]",
			},
			"card_buttons": {
				Type:        "string",
				Title:       "卡片按钮",
				Description: "卡片按钮列表，JSON 格式数组",
			},
			"card_custom_json": {
				Type:        "string",
				Title:       "自定义卡片 JSON",
				Description: "自定义卡片消息的完整 JSON（当卡片模板为 custom 时使用）",
			},
		},
		Required: []string{"webhook_url", "msg_type"},
	}

	logger, _ := zap.NewProduction()

	return &FeishuTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
		logger:   logger,
	}
}


func (t *FeishuTool) Execute(ctx *utools.ExecutionContext, toolConfig map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()


	webhookURL, _ := toolConfig["webhook_url"].(string)
	signSecret, _ := toolConfig["sign_secret"].(string)
	msgType, _ := toolConfig["msg_type"].(string)

	if webhookURL == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "Webhook URL 不能为空",
			Error:      "webhook_url is required",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("webhook_url is required")
	}


	filteredConfig := t.filterConfigByType(toolConfig, msgType)


	var messageBody map[string]interface{}
	var err error

	switch msgType {
	case "text":
		messageBody, err = t.buildTextMessage(filteredConfig)
	case "post":
		messageBody, err = t.buildPostMessage(filteredConfig)
	case "image":
		messageBody, err = t.buildImageMessage(filteredConfig)
	case "interactive":
		messageBody, err = t.buildInteractiveMessage(filteredConfig)
	default:
		return &utools.ExecutionResult{
			Success:    false,
			Message:    fmt.Sprintf("不支持的消息类型: %s", msgType),
			Error:      "unsupported message type",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("unsupported message type: %s", msgType)
	}

	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "构建消息失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}


	if signSecret != "" {
		timestamp := time.Now().Unix()
		sign := t.generateSign(signSecret, timestamp)
		messageBody["timestamp"] = strconv.FormatInt(timestamp, 10)
		messageBody["sign"] = sign
	}


	jsonData, err := json.Marshal(messageBody)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "序列化消息失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "发送消息失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}
	defer resp.Body.Close()


	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "读取响应失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}


	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "解析响应失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}


	code, _ := result["code"].(float64)
	if code != 0 {
		msg, _ := result["msg"].(string)
		return &utools.ExecutionResult{
			Success:    false,
			Message:    fmt.Sprintf("飞书返回错误: %s", msg),
			Error:      fmt.Sprintf("feishu error code: %v, msg: %s", code, msg),
			Output:     result,
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("feishu error: %s", msg)
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: "飞书消息发送成功",
		Output: map[string]interface{}{
			"msg_type": msgType,
			"response": result,
		},
		DurationMs: time.Since(startTime).Milliseconds(),
	}, nil
}


func (t *FeishuTool) buildTextMessage(config map[string]interface{}) (map[string]interface{}, error) {
	content, _ := config["content"].(string)
	if content == "" {
		return nil, fmt.Errorf("文本消息内容不能为空")
	}

	return map[string]interface{}{
		"msg_type": "text",
		"content": map[string]interface{}{
			"text": content,
		},
	}, nil
}


func (t *FeishuTool) buildPostMessage(config map[string]interface{}) (map[string]interface{}, error) {
	title, _ := config["title"].(string)
	postContentRaw := config["post_content"]

	if postContentRaw == nil {
		return nil, fmt.Errorf("富文本内容不能为空")
	}

	// 情况1: post_content 已经是数组（workflow 变量解析后的结果）
	if contentArray, ok := postContentRaw.([]interface{}); ok {
		// 检查是否是二维数组 [[{...}], [{...}]]
		contentLines := make([][]interface{}, 0, len(contentArray))
		for _, line := range contentArray {
			if lineArray, ok := line.([]interface{}); ok {
				contentLines = append(contentLines, lineArray)
			} else {
				// 如果不是二维数组，包装成一行
				contentLines = append(contentLines, []interface{}{line})
			}
		}

		if len(contentLines) == 0 {
			return nil, fmt.Errorf("富文本内容不能为空")
		}

		if title == "" {
			title = "通知"
		}
		return map[string]interface{}{
			"msg_type": "post",
			"content": map[string]interface{}{
				"post": map[string]interface{}{
					"zh_cn": map[string]interface{}{
						"title":   title,
						"content": contentLines,
					},
				},
			},
		}, nil
	}

	// 情况2: post_content 是 map（完整格式 {"zh_cn": {...}}）
	if fullFormat, ok := postContentRaw.(map[string]interface{}); ok {
		// 检查是否是完整格式（包含 zh_cn 或 en_us）
		if zhCn, ok := fullFormat["zh_cn"].(map[string]interface{}); ok {
			if _, hasContent := zhCn["content"]; hasContent {
				return map[string]interface{}{
					"msg_type": "post",
					"content": map[string]interface{}{
						"post": fullFormat,
					},
				}, nil
			}
		}
		if enUs, ok := fullFormat["en_us"].(map[string]interface{}); ok {
			if _, hasContent := enUs["content"]; hasContent {
				return map[string]interface{}{
					"msg_type": "post",
					"content": map[string]interface{}{
						"post": fullFormat,
					},
				}, nil
			}
		}
	}

	// 情况3: post_content 是字符串（需要 JSON 解析）
	postContent, ok := postContentRaw.(string)
	if !ok || postContent == "" {
		return nil, fmt.Errorf("富文本内容不能为空")
	}

	// 处理字符串中的实际换行符（JSON 不允许字符串中有实际换行符）
	postContent = strings.ReplaceAll(postContent, "\n", "")
	postContent = strings.ReplaceAll(postContent, "\r", "")

	// 先尝试解析为完整的飞书富文本格式 {"zh_cn": {"title": "...", "content": [[...]]}}
	var fullFormat map[string]interface{}
	if err := json.Unmarshal([]byte(postContent), &fullFormat); err == nil {
		// 检查是否是完整格式（包含 zh_cn 或 en_us）
		if zhCn, ok := fullFormat["zh_cn"].(map[string]interface{}); ok {
			if _, hasContent := zhCn["content"]; hasContent {
				// 完整格式，直接使用
				return map[string]interface{}{
					"msg_type": "post",
					"content": map[string]interface{}{
						"post": fullFormat,
					},
				}, nil
			}
		}
		if enUs, ok := fullFormat["en_us"].(map[string]interface{}); ok {
			if _, hasContent := enUs["content"]; hasContent {
				// 完整格式，直接使用
				return map[string]interface{}{
					"msg_type": "post",
					"content": map[string]interface{}{
						"post": fullFormat,
					},
				}, nil
			}
		}
	}

	// 再尝试解析为纯 content 数组格式 [[{...}], [{...}]]
	var contentLines [][]interface{}
	if err := json.Unmarshal([]byte(postContent), &contentLines); err == nil {
		if title == "" {
			title = "通知"
		}
		return map[string]interface{}{
			"msg_type": "post",
			"content": map[string]interface{}{
				"post": map[string]interface{}{
					"zh_cn": map[string]interface{}{
						"title":   title,
						"content": contentLines,
					},
				},
			},
		}, nil
	}

	// 都失败则当作纯文本处理
	if title == "" {
		title = "通知"
	}
	return map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title": title,
					"content": [][]interface{}{
						{
							map[string]interface{}{
								"tag":  "text",
								"text": postContent,
							},
						},
					},
				},
			},
		},
	}, nil
}


func (t *FeishuTool) buildImageMessage(config map[string]interface{}) (map[string]interface{}, error) {
	imageURL, _ := config["image_url"].(string)
	if imageURL == "" {
		return nil, fmt.Errorf("图片 URL 不能为空")
	}

	appID, _ := config["app_id"].(string)
	appSecret, _ := config["app_secret"].(string)


	if appID != "" && appSecret != "" {
		t.logger.Info("开始上传图片到飞书",
			zap.String("app_id", appID),
			zap.String("image_url", imageURL))

		imageKey, err := t.uploadImage(appID, appSecret, imageURL)
		if err != nil {
			t.logger.Error("上传图片失败，使用链接方案",
				zap.Error(err),
				zap.String("image_url", imageURL))
		} else if imageKey != "" {
			t.logger.Info("图片上传成功",
				zap.String("image_key", imageKey))

			return map[string]interface{}{
				"msg_type": "image",
				"content": map[string]interface{}{
					"image_key": imageKey,
				},
			}, nil
		}

	} else {
		t.logger.Info("未提供 App ID/Secret，使用链接方案显示图片")
	}


	title, _ := config["title"].(string)
	if title == "" {
		title = "图片"
	}

	content, _ := config["content"].(string)
	if content == "" {
		content = "点击查看图片"
	}


	return map[string]interface{}{
		"msg_type": "post",
		"content": map[string]interface{}{
			"post": map[string]interface{}{
				"zh_cn": map[string]interface{}{
					"title": title,
					"content": [][]interface{}{
						{
							map[string]interface{}{
								"tag":  "text",
								"text": content + " ",
							},
							map[string]interface{}{
								"tag":  "a",
								"text": "🔗 点击查看图片",
								"href": imageURL,
							},
						},
					},
				},
			},
		},
	}, nil
}


func (t *FeishuTool) buildInteractiveMessage(config map[string]interface{}) (map[string]interface{}, error) {
	cardTemplate, _ := config["card_template"].(string)
	if cardTemplate == "" {
		cardTemplate = "notification"
	}

	var card map[string]interface{}
	var err error

	if cardTemplate == "custom" {

		customJSON, _ := config["card_custom_json"].(string)
		if customJSON == "" {
			return nil, fmt.Errorf("自定义卡片 JSON 不能为空")
		}
		if err := json.Unmarshal([]byte(customJSON), &card); err != nil {
			return nil, fmt.Errorf("自定义卡片 JSON 格式错误: %v", err)
		}
	} else {

		card, err = t.buildCardByTemplate(cardTemplate, config)
		if err != nil {
			return nil, err
		}
	}

	return map[string]interface{}{
		"msg_type": "interactive",
		"card":     card,
	}, nil
}


func (t *FeishuTool) buildCardByTemplate(template string, config map[string]interface{}) (map[string]interface{}, error) {
	title, _ := config["title"].(string)
	content, _ := config["card_content"].(string)
	status, _ := config["card_status"].(string)
	fieldsJSON, _ := config["card_fields"].(string)
	buttonsJSON, _ := config["card_buttons"].(string)

	if title == "" {
		title = "通知"
	}
	if status == "" {
		status = "info"
	}


	var fields []map[string]string
	if fieldsJSON != "" {
		if err := json.Unmarshal([]byte(fieldsJSON), &fields); err != nil {
			return nil, fmt.Errorf("卡片字段 JSON 格式错误: %v", err)
		}
	}


	var buttons []map[string]string
	if buttonsJSON != "" {
		if err := json.Unmarshal([]byte(buttonsJSON), &buttons); err != nil {
			return nil, fmt.Errorf("卡片按钮 JSON 格式错误: %v", err)
		}
	}


	elements := make([]interface{}, 0)


	statusEmoji := t.getStatusEmoji(status)
	statusColor := t.getStatusColor(status)

	elements = append(elements, map[string]interface{}{
		"tag": "div",
		"text": map[string]interface{}{
			"tag":     "lark_md",
			"content": fmt.Sprintf("%s **%s**", statusEmoji, title),
		},
	})


	elements = append(elements, map[string]interface{}{
		"tag": "hr",
	})


	if content != "" {
		elements = append(elements, map[string]interface{}{
			"tag": "div",
			"text": map[string]interface{}{
				"tag":     "lark_md",
				"content": content,
			},
		})
	}


	if len(fields) > 0 {
		for _, field := range fields {
			elements = append(elements, map[string]interface{}{
				"tag": "div",
				"fields": []interface{}{
					map[string]interface{}{
						"is_short": true,
						"text": map[string]interface{}{
							"tag":     "lark_md",
							"content": fmt.Sprintf("**%s**\n%s", field["key"], field["value"]),
						},
					},
				},
			})
		}
	}


	if len(buttons) > 0 {
		actions := make([]interface{}, 0)
		for _, button := range buttons {
			actions = append(actions, map[string]interface{}{
				"tag": "button",
				"text": map[string]interface{}{
					"tag":     "plain_text",
					"content": button["text"],
				},
				"type": "primary",
				"url":  button["url"],
			})
		}

		elements = append(elements, map[string]interface{}{
			"tag":     "action",
			"actions": actions,
		})
	}


	config_card := map[string]interface{}{
		"wide_screen_mode": true,
	}


	header := map[string]interface{}{
		"title": map[string]interface{}{
			"tag":     "plain_text",
			"content": title,
		},
		"template": statusColor,
	}

	return map[string]interface{}{
		"config":   config_card,
		"header":   header,
		"elements": elements,
	}, nil
}


func (t *FeishuTool) getStatusEmoji(status string) string {
	switch status {
	case "success":
		return "✅"
	case "warning":
		return "⚠️"
	case "error":
		return "❌"
	case "info":
		return "ℹ️"
	default:
		return "📋"
	}
}


func (t *FeishuTool) getStatusColor(status string) string {
	switch status {
	case "success":
		return "green"
	case "warning":
		return "yellow"
	case "error":
		return "red"
	case "info":
		return "blue"
	default:
		return "blue"
	}
}


func (t *FeishuTool) filterConfigByType(config map[string]interface{}, msgType string) map[string]interface{} {

	filtered := map[string]interface{}{
		"webhook_url": config["webhook_url"],
		"sign_secret": config["sign_secret"],
		"app_id":      config["app_id"],
		"app_secret":  config["app_secret"],
		"msg_type":    config["msg_type"],
	}


	switch msgType {
	case "text":

		filtered["content"] = config["content"]

	case "post":

		filtered["title"] = config["title"]
		filtered["post_content"] = config["post_content"]

	case "image":

		filtered["image_url"] = config["image_url"]
		filtered["title"] = config["title"]
		filtered["content"] = config["content"]

	case "interactive":

		filtered["title"] = config["title"]
		filtered["card_template"] = config["card_template"]
		filtered["card_content"] = config["card_content"]
		filtered["card_status"] = config["card_status"]
		filtered["card_fields"] = config["card_fields"]
		filtered["card_buttons"] = config["card_buttons"]
		filtered["card_custom_json"] = config["card_custom_json"]
	}

	return filtered
}



func (t *FeishuTool) generateSign(secret string, timestamp int64) string {

	stringToSign := fmt.Sprintf("%v", timestamp) + "\n" + secret


	h := sha256.New()
	h.Write([]byte(stringToSign))


	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))
	return signature
}


func (t *FeishuTool) getTenantAccessToken(appID, appSecret string) (string, error) {
	url := "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal"

	payload := map[string]string{
		"app_id":     appID,
		"app_secret": appSecret,
	}

	t.logger.Info("正在获取飞书 tenant_access_token",
		zap.String("app_id", appID),
		zap.String("url", url))

	jsonData, _ := json.Marshal(payload)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.logger.Error("HTTP 请求失败", zap.Error(err))
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	t.logger.Info("获取 token 响应", zap.String("response", string(body)))

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		t.logger.Error("解析响应失败", zap.Error(err))
		return "", err
	}

	code, _ := result["code"].(float64)
	if code != 0 {
		msg, _ := result["msg"].(string)
		t.logger.Error("获取 token 失败", zap.Float64("code", code), zap.String("msg", msg))
		return "", fmt.Errorf("获取 token 失败: %s", msg)
	}

	token, _ := result["tenant_access_token"].(string)
	t.logger.Info("成功获取 tenant_access_token", zap.String("token", token[:20]+"..."))
	return token, nil
}


func (t *FeishuTool) uploadImage(appID, appSecret, imageURL string) (string, error) {

	t.logger.Info("步骤 1: 获取 tenant_access_token")
	token, err := t.getTenantAccessToken(appID, appSecret)
	if err != nil {
		t.logger.Error("获取 token 失败", zap.Error(err))
		return "", err
	}


	t.logger.Info("步骤 2: 下载图片", zap.String("image_url", imageURL))
	resp, err := http.Get(imageURL)
	if err != nil {
		t.logger.Error("下载图片失败", zap.Error(err))
		return "", fmt.Errorf("下载图片失败: %v", err)
	}
	defer resp.Body.Close()

	imageData, err := io.ReadAll(resp.Body)
	if err != nil {
		t.logger.Error("读取图片数据失败", zap.Error(err))
		return "", fmt.Errorf("读取图片失败: %v", err)
	}
	t.logger.Info("图片下载成功", zap.Int("size", len(imageData)))


	uploadURL := "https://open.feishu.cn/open-apis/im/v1/images"
	t.logger.Info("步骤 3: 上传图片到飞书", zap.String("upload_url", uploadURL))


	body := &bytes.Buffer{}


	boundary := "----WebKitFormBoundary7MA4YWxkTrZu0gW"
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Disposition: form-data; name=\"image_type\"\r\n\r\n")
	body.WriteString("message\r\n")
	body.WriteString(fmt.Sprintf("--%s\r\n", boundary))
	body.WriteString("Content-Disposition: form-data; name=\"image\"; filename=\"image.jpg\"\r\n")
	body.WriteString("Content-Type: image/jpeg\r\n\r\n")
	body.Write(imageData)
	body.WriteString(fmt.Sprintf("\r\n--%s--\r\n", boundary))

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		t.logger.Error("创建上传请求失败", zap.Error(err))
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", fmt.Sprintf("multipart/form-data; boundary=%s", boundary))

	client := &http.Client{}
	uploadResp, err := client.Do(req)
	if err != nil {
		t.logger.Error("执行上传请求失败", zap.Error(err))
		return "", fmt.Errorf("上传图片失败: %v", err)
	}
	defer uploadResp.Body.Close()

	uploadBody, _ := io.ReadAll(uploadResp.Body)
	t.logger.Info("上传响应", zap.String("response", string(uploadBody)))

	var uploadResult map[string]interface{}
	if err := json.Unmarshal(uploadBody, &uploadResult); err != nil {
		t.logger.Error("解析上传响应失败", zap.Error(err))
		return "", err
	}

	code, _ := uploadResult["code"].(float64)
	if code != 0 {
		msg, _ := uploadResult["msg"].(string)
		t.logger.Error("飞书上传接口返回错误", zap.Float64("code", code), zap.String("msg", msg))
		return "", fmt.Errorf("上传图片失败: %s", msg)
	}

	data, _ := uploadResult["data"].(map[string]interface{})
	imageKey, _ := data["image_key"].(string)

	t.logger.Info("图片上传成功", zap.String("image_key", imageKey))
	return imageKey, nil
}


func init() {
	tool := NewFeishuTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register Feishu tool: %v", err))
	}
}
