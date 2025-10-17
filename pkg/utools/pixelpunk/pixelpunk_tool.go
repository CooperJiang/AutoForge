package pixelpunk

import (
	config2 "auto-forge/pkg/config"
	log "auto-forge/pkg/logger"
	"auto-forge/pkg/utools"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
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
		Description: "上传图片到 PixelPunk 图床，返回 CDN URL",
		Category:    "data",
		Version:     "1.0.0",
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
				Title:       "文件",
				Description: "要上传的图片文件（从 external_trigger 接收的文件对象）",
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
		Required: []string{"file"},
	}

	return &PixelPunkUploadTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

// Execute 执行工具
func (t *PixelPunkUploadTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()

	// 1. 检查 PixelPunk 是否启用
	cfg := config2.GetConfig()
	if cfg.PixelPunk.BaseURL == "" || !cfg.PixelPunk.Enabled {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "PixelPunk 图床未配置或未启用",
			Error:      "请在 config.yaml 中配置 pixelpunk 相关参数",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("PixelPunk 未配置")
	}

	// 2. 获取文件对象
	fileObj, ok := config["file"].(map[string]interface{})
	if !ok {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "文件参数格式错误",
			Error:      "file 参数必须是文件对象",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("无效的文件参数")
	}

	// 3. 获取文件路径
	filePath, ok := fileObj["path"].(string)
	if !ok || filePath == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "文件路径不存在",
			Error:      "无法获取文件路径",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("文件路径为空")
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
	uploadedFile, err := t.uploadFile(cfg.PixelPunk.BaseURL, cfg.PixelPunk.APIKey, filePath, accessLevel, optimize, filePaths, folderID)
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

// init 自动注册工具
func init() {
	tool := NewPixelPunkUploadTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register PixelPunk upload tool: %v", err))
	}
}
