package aliyunoss

import (
	toolConfigService "auto-forge/internal/services/tool_config"
	"auto-forge/pkg/utools"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type AliyunOSSTool struct {
	*utools.BaseTool
}

func NewAliyunOSSTool() *AliyunOSSTool {
	metadata := &utools.ToolMetadata{
		Code:        "aliyun_oss",
		Name:        "阿里云 OSS 上传",
		Description: "上传文件到阿里云对象存储服务",
		Category:    utools.CategoryStorage,
		Version:     "1.0.0",
		Author:      "Cooper Team",
		AICallable:  true,
		Tags:        []string{"storage", "upload", "aliyun", "oss"},
		OutputFieldsSchema: map[string]utools.OutputFieldDef{
			"response": {
				Type:  "object",
				Label: "完整响应",
				Children: map[string]utools.OutputFieldDef{
					"url": {
						Type:  "string",
						Label: "文件 URL",
					},
					"filename": {
						Type:  "string",
						Label: "文件名",
					},
					"size": {
						Type:  "number",
						Label: "文件大小（字节）",
					},
					"bucket": {
						Type:  "string",
						Label: "存储桶",
					},
					"uploaded_at": {
						Type:  "string",
						Label: "上传时间",
					},
				},
			},
			"url": {
				Type:  "string",
				Label: "文件 URL（快捷访问）",
			},
			"filename": {
				Type:  "string",
				Label: "文件名（快捷访问）",
			},
		},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"file": {
				Type:        "string",
				Title:       "文件",
				Description: "要上传的文件路径",
			},
			"path": {
				Type:        "string",
				Title:       "存储路径",
				Description: "OSS 中的存储路径（可选，默认为文件名）",
			},
		},
		Required: []string{"file"},
	}

	return &AliyunOSSTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

func (t *AliyunOSSTool) Execute(ctx *utools.ExecutionContext, toolConfig map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()

	// 从数据库加载 OSS 配置
	dbConfig, err := toolConfigService.GetToolConfigForExecution("aliyun_oss")
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "阿里云 OSS 配置错误",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 解析配置字段
	endpoint, _ := dbConfig["endpoint"].(string)
	accessKeyID, _ := dbConfig["access_key_id"].(string)
	accessKeySecret, _ := dbConfig["access_key_secret"].(string)
	bucket, _ := dbConfig["bucket"].(string)

	// 验证配置
	if endpoint == "" || accessKeyID == "" || accessKeySecret == "" || bucket == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "阿里云 OSS 配置不完整",
			Error:      "missing required OSS configuration",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("阿里云 OSS 配置不完整")
	}

	// 解析文件参数（优先处理文件对象）
	var filePath string

	// 1. 检查是否为文件对象（从外部API触发器传入）
	if fileObj, ok := toolConfig["file"].(map[string]interface{}); ok {
		if path, ok := fileObj["path"].(string); ok && path != "" {
			filePath = path
		}
	}

	// 2. 如果不是文件对象,则尝试作为字符串路径
	if filePath == "" {
		if strPath, ok := toolConfig["file"].(string); ok && strPath != "" {
			filePath = strPath
		}
	}

	// 3. 最终验证
	if filePath == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "文件参数无效",
			Error:      "invalid file parameter: must be file object or file path string",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("文件参数无效")
	}

	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "文件不存在",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 读取文件内容
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "读取文件失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 确定 OSS 路径
	ossPath := filepath.Base(filePath)
	if customPath, ok := toolConfig["path"].(string); ok && customPath != "" {
		ossPath = customPath
	}

	// 上传到 OSS
	fileURL, err := t.uploadToOSS(endpoint, accessKeyID, accessKeySecret, bucket, ossPath, fileData)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "上传到阿里云 OSS 失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 构建响应对象
	responseObj := map[string]interface{}{
		"url":         fileURL,
		"filename":    filepath.Base(filePath),
		"size":        fileInfo.Size(),
		"bucket":      bucket,
		"uploaded_at": time.Now().Format(time.RFC3339),
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("文件已上传到阿里云 OSS: %s", fileURL),
		Output: map[string]interface{}{
			"response": responseObj,
			"url":      fileURL,
			"filename": filepath.Base(filePath),
		},
		DurationMs: time.Since(startTime).Milliseconds(),
	}, nil
}

func (t *AliyunOSSTool) uploadToOSS(endpoint, accessKeyID, accessKeySecret, bucket, objectKey string, data []byte) (string, error) {
	// 构建 OSS URL
	host := fmt.Sprintf("%s.%s", bucket, endpoint)
	url := fmt.Sprintf("https://%s/%s", host, objectKey)

	// 创建请求
	req, err := http.NewRequest("PUT", url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	// 设置请求头
	contentType := http.DetectContentType(data)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))

	// 生成签名
	date := time.Now().UTC().Format(http.TimeFormat)
	req.Header.Set("Date", date)

	signature := t.generateSignature(accessKeySecret, bucket, "PUT", objectKey, contentType, date)
	req.Header.Set("Authorization", fmt.Sprintf("OSS %s:%s", accessKeyID, signature))

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("OSS upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return url, nil
}

func (t *AliyunOSSTool) generateSignature(accessKeySecret, bucket, method, objectKey, contentType, date string) string {
	// 构建签名字符串
	canonicalizedResource := fmt.Sprintf("/%s/%s", bucket, objectKey)
	stringToSign := fmt.Sprintf("%s\n\n%s\n%s\n%s",
		method,
		contentType,
		date,
		canonicalizedResource,
	)

	// 使用 HMAC-SHA1 生成签名
	h := hmac.New(sha1.New, []byte(accessKeySecret))
	h.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	return signature
}

func init() {
	tool := NewAliyunOSSTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register AliyunOSS tool: %v", err))
	}
}
