package pixelpunk

import (
	toolConfigService "auto-forge/internal/services/tool_config"
	"auto-forge/pkg/agent/tooling"
	log "auto-forge/pkg/logger"
	"auto-forge/pkg/utools"
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// PixelPunkUploadTool PixelPunk 图床上传工具
type PixelPunkUploadTool struct {
	*utools.BaseTool
}

// UploadResponse 上传响应结构（单文件上传）
type UploadResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		Uploaded       UploadedFile `json:"uploaded"` // 单文件上传返回对象，不是数组
		OversizedFiles []string     `json:"oversized_files,omitempty"`
		SizeLimit      string       `json:"size_limit,omitempty"`
		UploadErrors   []string     `json:"upload_errors,omitempty"`
	} `json:"data"`
}

// UploadedFile 上传成功的文件信息
type UploadedFile struct {
	ID           string `json:"id"`
	URL          string `json:"url"`
	ThumbURL     string `json:"thumb_url"`
	OriginalName string `json:"original_name"`
	Size         int64  `json:"size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Format       string `json:"format"`
	AccessLevel  string `json:"access_level"`
	CreatedAt    string `json:"created_at"`
}

// NewPixelPunkUploadTool 创建 PixelPunk 图床上传工具实例
func NewPixelPunkUploadTool() *PixelPunkUploadTool {
	metadata := &utools.ToolMetadata{
		Code:        "pixelpunk_upload",
		Name:        "PixelPunk 图床上传",
		Description: "上传图片到 PixelPunk 图床，支持文件对象或 URL，返回 CDN URL",
		Category:    utools.CategoryStorage,
		Version:     "1.1.0",
		Author:      "Cooper Team",
		AICallable:  true,
		Tags:        []string{"image", "upload", "cdn", "storage", "pixelpunk"},
		OutputFieldsSchema: map[string]utools.OutputFieldDef{
			"response": {
				Type:  "object",
				Label: "PixelPunk 完整响应",
				Children: map[string]utools.OutputFieldDef{
					"url": {
						Type:  "string",
						Label: "图片 CDN 地址",
					},
					"thumb_url": {
						Type:  "string",
						Label: "缩略图地址",
					},
					"id": {
						Type:  "string",
						Label: "图片唯一 ID",
					},
					"original_name": {
						Type:  "string",
						Label: "原始文件名",
					},
					"size": {
						Type:  "number",
						Label: "文件大小（字节）",
					},
					"width": {
						Type:  "number",
						Label: "图片宽度",
					},
					"height": {
						Type:  "number",
						Label: "图片高度",
					},
					"format": {
						Type:  "string",
						Label: "图片格式（png/jpg/webp 等）",
					},
					"access_level": {
						Type:  "string",
						Label: "访问级别",
					},
					"created_at": {
						Type:  "string",
						Label: "创建时间",
					},
				},
			},
			// 快捷访问字段（平铺常用字段）
			"url": {
				Type:  "string",
				Label: "图片 CDN 地址（快捷访问）",
			},
			"thumb_url": {
				Type:  "string",
				Label: "缩略图地址（快捷访问）",
			},
			"id": {
				Type:  "string",
				Label: "图片 ID（快捷访问）",
			},
		},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"file": {
				Type:        "object",
				Title:       "文件对象",
				Description: "要上传的图片文件对象（与 url 二选一）",
			},
			"url": {
				Type:        "string",
				Title:       "图片 URL",
				Description: "要上传的图片 URL，工具会自动下载后上传（与 file 二选一）",
				Format:      "uri",
			},
			"access_level": {
				Type:        "string",
				Title:       "访问级别",
				Description: "图片访问权限",
				Enum:        []interface{}{"public", "private", "protected"},
				Default:     "public",
			},
			"optimize": {
				Type:        "boolean",
				Title:       "优化图片",
				Description: "是否对图片进行优化压缩",
				Default:     true,
			},
			"file_path": {
				Type:        "string",
				Title:       "虚拟路径",
				Description: "图片存储的虚拟路径（如：projects/website）",
				Default:     "",
			},
			"folder_id": {
				Type:        "string",
				Title:       "文件夹ID",
				Description: "目标文件夹ID（优先级高于虚拟路径）",
				Default:     "",
			},
		},
		Required: []string{},
	}

	return &PixelPunkUploadTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

// Execute 执行工具
func (t *PixelPunkUploadTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()

	// 1. 从数据库加载 PixelPunk 配置
	dbConfig, err := toolConfigService.GetToolConfigForExecution("pixelpunk_upload")
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "PixelPunk 配置错误",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 解析配置字段
	baseURL, _ := dbConfig["base_url"].(string)
	apiKey, _ := dbConfig["api_key"].(string)

	// 验证配置
	if baseURL == "" || apiKey == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "PixelPunk 配置不完整",
			Error:      "missing required PixelPunk configuration",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("PixelPunk 配置不完整")
	}

	// 2. 获取文件路径（支持 file 对象或 url）
	var filePath string
	var tempFile string // 用于记录临时文件，最后需要清理

	// 优先检查 url 参数
	if imageURL, ok := config["url"].(string); ok && imageURL != "" {
		log.Info("从 URL 下载图片: %s", imageURL)
		downloadedPath, err := t.downloadFromURL(imageURL)
		if err != nil {
			return &utools.ExecutionResult{
				Success:    false,
				Message:    "下载图片失败",
				Error:      err.Error(),
				DurationMs: time.Since(startTime).Milliseconds(),
			}, err
		}
		filePath = downloadedPath
		tempFile = downloadedPath // 标记为临时文件
		log.Info("图片下载成功: %s", filePath)
	} else if fileObj, ok := config["file"].(map[string]interface{}); ok {
		// 使用 file 对象
		path, ok := fileObj["path"].(string)
		if !ok || path == "" {
			return &utools.ExecutionResult{
				Success:    false,
				Message:    "文件路径不存在",
				Error:      "无法获取文件路径",
				DurationMs: time.Since(startTime).Milliseconds(),
			}, fmt.Errorf("文件路径为空")
		}
		filePath = path
	} else {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "缺少必需参数",
			Error:      "必须提供 file 对象或 url 参数",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("缺少 file 或 url 参数")
	}

	// 确保在函数结束时清理临时文件
	if tempFile != "" {
		defer func() {
			if err := os.Remove(tempFile); err != nil {
				log.Info("清理临时文件失败: %s, error: %v", tempFile, err)
			} else {
				log.Info("已清理临时文件: %s", tempFile)
			}
		}()
	}

	// 4. 获取可选参数
	accessLevel := "public"
	if val, ok := config["access_level"].(string); ok && val != "" {
		accessLevel = val
	}

	optimize := true
	if val, ok := config["optimize"].(bool); ok {
		optimize = val
	}

	filePaths := ""
	if val, ok := config["file_path"].(string); ok {
		filePaths = val
	}

	folderID := ""
	if val, ok := config["folder_id"].(string); ok {
		folderID = val
	}

	// 5. 上传文件到 PixelPunk
	log.Info("开始上传文件到 PixelPunk: %s", filePath)
	uploadedFile, err := t.uploadFile(baseURL, apiKey, filePath, accessLevel, optimize, filePaths, folderID)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "上传失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 6. 构建响应对象
	responseObj := map[string]interface{}{
		"url":           uploadedFile.URL,
		"thumb_url":     uploadedFile.ThumbURL,
		"id":            uploadedFile.ID,
		"original_name": uploadedFile.OriginalName,
		"size":          uploadedFile.Size,
		"width":         uploadedFile.Width,
		"height":        uploadedFile.Height,
		"format":        uploadedFile.Format,
		"access_level":  uploadedFile.AccessLevel,
		"created_at":    uploadedFile.CreatedAt,
	}

	// 7. 返回成功结果（包含完整响应对象 + 快捷访问字段）
	log.Info("文件上传成功: %s -> %s", uploadedFile.OriginalName, uploadedFile.URL)
	return &utools.ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("文件已上传到 PixelPunk: %s", uploadedFile.OriginalName),
		Output: map[string]interface{}{
			"response":  responseObj,      // 完整响应对象
			"url":       uploadedFile.URL, // 快捷访问
			"thumb_url": uploadedFile.ThumbURL,
			"id":        uploadedFile.ID,
		},
		DurationMs: time.Since(startTime).Milliseconds(),
	}, nil
}

// uploadFile 上传文件到 PixelPunk
func (t *PixelPunkUploadTool) uploadFile(baseURL, apiKey, filePath, accessLevel string, optimize bool, virtualPath, folderID string) (*UploadedFile, error) {
	// 1. 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("无法打开文件: %w", err)
	}
	defer file.Close()

	// 2. 创建 multipart form
	var requestBody bytes.Buffer
	writer := multipart.NewWriter(&requestBody)

	// 3. 添加文件字段
	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return nil, fmt.Errorf("创建文件字段失败: %w", err)
	}

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("拷贝文件内容失败: %w", err)
	}

	// 4. 添加其他字段
	if accessLevel != "" {
		writer.WriteField("access_level", accessLevel)
	}
	if optimize {
		writer.WriteField("optimize", "true")
	} else {
		writer.WriteField("optimize", "false")
	}
	if virtualPath != "" {
		writer.WriteField("filePath", virtualPath)
	}
	if folderID != "" {
		writer.WriteField("folderId", folderID)
	}

	// 5. 关闭 writer
	err = writer.Close()
	if err != nil {
		return nil, fmt.Errorf("关闭 writer 失败: %w", err)
	}

	// 6. 创建 HTTP 请求
	url := fmt.Sprintf("%s/api/v1/external/upload", baseURL)
	req, err := http.NewRequest("POST", url, &requestBody)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 7. 设置请求头
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("x-pixelpunk-key", apiKey)

	// 8. 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 9. 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 10. 解析响应
	var uploadResp UploadResponse
	err = json.Unmarshal(body, &uploadResp)
	if err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, body: %s", err, string(body))
	}

	// 11. 检查业务错误码（0 表示成功，200 也表示成功）
	if uploadResp.Code != 0 && uploadResp.Code != 200 {
		return nil, t.handleError(uploadResp.Code, uploadResp.Message)
	}

	// 12. 检查文件 URL 是否存在
	if uploadResp.Data.Uploaded.URL == "" {
		return nil, fmt.Errorf("上传失败: 没有返回文件 URL")
	}

	log.Info("文件上传成功: URL=%s, ID=%s", uploadResp.Data.Uploaded.URL, uploadResp.Data.Uploaded.ID)
	return &uploadResp.Data.Uploaded, nil
}

// handleError 处理 PixelPunk 错误码
func (t *PixelPunkUploadTool) handleError(code int, message string) error {
	switch code {
	case 102:
		return fmt.Errorf("未授权: 请检查 API Key 是否有效 (%s)", message)
	case 4000:
		return fmt.Errorf("文件大小超限: 请压缩文件或检查 API Key 限制 (%s)", message)
	case 4001:
		return fmt.Errorf("文件格式不支持: 请转换为支持的格式 (%s)", message)
	case 4008:
		return fmt.Errorf("存储容量已用尽: 请清理旧文件或升级容量 (%s)", message)
	case 4009:
		return fmt.Errorf("上传次数已用尽: 请等待下个周期或升级配额 (%s)", message)
	default:
		return fmt.Errorf("上传失败 (错误码: %d): %s", code, message)
	}
}

// downloadFromURL 从 URL 下载文件到临时目录
func (t *PixelPunkUploadTool) downloadFromURL(urlStr string) (string, error) {
	// 1. 验证 URL
	parsedURL, err := neturl.Parse(urlStr)
	if err != nil {
		return "", fmt.Errorf("无效的 URL: %w", err)
	}
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return "", fmt.Errorf("不支持的协议: %s (仅支持 http/https)", parsedURL.Scheme)
	}

	// 2. 创建 HTTP 客户端（60 秒超时，跟随重定向，可选禁用 SSL 验证）
	client := &http.Client{
		Timeout: 60 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, // 允许自签名证书
		},
	}

	// 3. 发送 GET 请求
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %w", err)
	}
	req.Header.Set("User-Agent", "AutoForge-PixelPunk/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("下载失败: %w", err)
	}
	defer resp.Body.Close()

	// 4. 检查状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败: HTTP %d", resp.StatusCode)
	}

	// 5. 推断文件名
	filename := t.inferFilename(urlStr, resp.Header.Get("Content-Type"))

	// 6. 创建临时文件
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, filename)

	outFile, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("创建临时文件失败: %w", err)
	}
	defer outFile.Close()

	// 7. 写入文件
	written, err := io.Copy(outFile, resp.Body)
	if err != nil {
		os.Remove(tempFile) // 清理失败的文件
		return "", fmt.Errorf("保存文件失败: %w", err)
	}

	log.Info("文件下载完成: %s (%d bytes)", tempFile, written)
	return tempFile, nil
}

// inferFilename 从 URL 或 Content-Type 推断文件名
func (t *PixelPunkUploadTool) inferFilename(urlStr, contentType string) string {
	// 1. 尝试从 URL 路径提取文件名
	parsedURL, _ := neturl.Parse(urlStr)
	if parsedURL != nil {
		pathParts := strings.Split(parsedURL.Path, "/")
		if len(pathParts) > 0 {
			lastPart := pathParts[len(pathParts)-1]
			// 检查是否有文件扩展名
			if strings.Contains(lastPart, ".") {
				// 清理文件名（移除查询参数等）
				cleanName := regexp.MustCompile(`[?#].*`).ReplaceAllString(lastPart, "")
				if cleanName != "" {
					return cleanName
				}
			}
		}
	}

	// 2. 根据 Content-Type 生成文件名
	ext := ".png" // 默认扩展名
	if strings.Contains(contentType, "image/jpeg") || strings.Contains(contentType, "image/jpg") {
		ext = ".jpg"
	} else if strings.Contains(contentType, "image/png") {
		ext = ".png"
	} else if strings.Contains(contentType, "image/gif") {
		ext = ".gif"
	} else if strings.Contains(contentType, "image/webp") {
		ext = ".webp"
	}

	// 3. 生成唯一文件名
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("download_%d%s", timestamp, ext)
}

// GetExecutionConfig 返回工具执行配置（实现 ConfigurableTool 接口）
func (t *PixelPunkUploadTool) GetExecutionConfig() *tooling.ExecutionConfig {
	return &tooling.ExecutionConfig{
		// 超时配置：上传可能需要较长时间
		TimeoutSeconds: 120, // 2 分钟

		// 重试配置：网络问题可以重试
		Retry: &tooling.RetryConfig{
			MaxRetries:        2,
			InitialBackoff:    2000,  // 2 秒
			MaxBackoff:        10000, // 10 秒
			BackoffMultiplier: 2.0,
			RetryableErrors: []string{
				"timeout",
				"connection",
				"network",
				"503",
				"504",
				"EOF",
			},
		},

		// 依赖配置
		Dependencies: &tooling.DependencyConfig{
			// 需要图片 URL 或文件对象
			Requires: []string{"image_url", "file_object"},

			// 提供 CDN URL
			Provides: []string{"cdn_url", "image_info"},

			// 建议的前置工具
			SuggestedPredecessors: []string{
				"openai_image",    // AI 生成图片
				"file_downloader", // 下载文件
			},
		},

		// 缓存配置：相同的图片可以缓存
		Cache: &tooling.CacheConfig{
			Enabled: false, // 暂不启用（每次上传都会生成新的 URL）
			TTL:     0,
		},
	}
}

// init 自动注册工具
func init() {
	tool := NewPixelPunkUploadTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register PixelPunk upload tool: %v", err))
	}
}
