package tencentcos

import (
	"auto-forge/pkg/config"
	"auto-forge/pkg/utools"
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type TencentCOSTool struct {
	*utools.BaseTool
}

func NewTencentCOSTool() *TencentCOSTool {
	metadata := &utools.ToolMetadata{
		Code:        "tencent_cos",
		Name:        "腾讯云 COS 上传",
		Description: "上传文件到腾讯云对象存储服务",
		Category:    "data",
		Version:     "1.0.0",
		Author:      "Cooper Team",
		AICallable:  true,
		Tags:        []string{"storage", "upload", "tencent", "cos"},
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
				Description: "COS 中的存储路径（可选，默认为文件名）",
			},
		},
		Required: []string{"file"},
	}

	return &TencentCOSTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}

func (t *TencentCOSTool) Execute(ctx *utools.ExecutionContext, toolConfig map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()

	// 从配置文件加载 COS 配置
	cfg := config.GetConfig()

	if !cfg.TencentCOS.Enabled {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "腾讯云 COS 未启用",
			Error:      "tencent_cos.enabled is false in config",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("腾讯云 COS 未启用")
	}

	// 验证配置
	if cfg.TencentCOS.SecretID == "" || cfg.TencentCOS.SecretKey == "" || cfg.TencentCOS.Bucket == "" || cfg.TencentCOS.Region == "" {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "腾讯云 COS 配置不完整",
			Error:      "missing required COS configuration",
			DurationMs: time.Since(startTime).Milliseconds(),
		}, fmt.Errorf("腾讯云 COS 配置不完整")
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

	// 确定 COS 路径
	cosPath := filepath.Base(filePath)
	if customPath, ok := toolConfig["path"].(string); ok && customPath != "" {
		cosPath = customPath
	}

	// 上传到 COS
	fileURL, err := t.uploadToCOS(&cfg.TencentCOS, cosPath, fileData)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "上传到腾讯云 COS 失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	// 构建响应对象
	responseObj := map[string]interface{}{
		"url":         fileURL,
		"filename":    filepath.Base(filePath),
		"size":        fileInfo.Size(),
		"bucket":      cfg.TencentCOS.Bucket,
		"uploaded_at": time.Now().Format(time.RFC3339),
	}

	return &utools.ExecutionResult{
		Success: true,
		Message: fmt.Sprintf("文件已上传到腾讯云 COS: %s", fileURL),
		Output: map[string]interface{}{
			"response": responseObj,
			"url":      fileURL,
			"filename": filepath.Base(filePath),
		},
		DurationMs: time.Since(startTime).Milliseconds(),
	}, nil
}

func (t *TencentCOSTool) uploadToCOS(cosConfig *config.TencentCOSConfig, objectKey string, data []byte) (string, error) {
	// 构建 COS URL
	host := fmt.Sprintf("%s.cos.%s.myqcloud.com", cosConfig.Bucket, cosConfig.Region)
	uploadURL := fmt.Sprintf("https://%s/%s", host, objectKey)

	// 创建请求
	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	// 设置请求头
	contentType := http.DetectContentType(data)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-Length", fmt.Sprintf("%d", len(data)))
	req.Header.Set("Host", host)

	// 生成签名
	authorization := t.generateAuthorization(cosConfig, req, objectKey)
	req.Header.Set("Authorization", authorization)

	// 发送请求
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("COS upload failed with status %d: %s", resp.StatusCode, string(body))
	}

	return uploadURL, nil
}

func (t *TencentCOSTool) generateAuthorization(cosConfig *config.TencentCOSConfig, req *http.Request, objectKey string) string {
	// 腾讯云 COS 签名算法
	now := time.Now()
	expiredTime := now.Add(time.Hour)

	// q-sign-time
	signTime := fmt.Sprintf("%d;%d", now.Unix(), expiredTime.Unix())

	// q-key-time (与 q-sign-time 相同)
	keyTime := signTime

	// 计算 SignKey
	signKey := t.hmacSha1(cosConfig.SecretKey, keyTime)

	// HttpString
	httpMethod := strings.ToLower(req.Method)
	uriPathname := "/" + objectKey
	httpParameters := t.getCanonicalQueryString(req.URL.Query())
	httpHeaders := t.getCanonicalHeaders(req.Header)
	httpString := fmt.Sprintf("%s\n%s\n%s\n%s\n", httpMethod, uriPathname, httpParameters, httpHeaders)

	// StringToSign
	httpStringSha1 := t.sha1Hash(httpString)
	stringToSign := fmt.Sprintf("sha1\n%s\n%s\n", signTime, httpStringSha1)

	// Signature
	signature := t.hmacSha1(signKey, stringToSign)

	// Authorization
	authorization := fmt.Sprintf("q-sign-algorithm=sha1&q-ak=%s&q-sign-time=%s&q-key-time=%s&q-header-list=%s&q-url-param-list=%s&q-signature=%s",
		cosConfig.SecretID,
		signTime,
		keyTime,
		t.getHeaderList(req.Header),
		t.getParamList(req.URL.Query()),
		signature,
	)

	return authorization
}

func (t *TencentCOSTool) hmacSha1(key, data string) string {
	h := hmac.New(sha1.New, []byte(key))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (t *TencentCOSTool) sha1Hash(data string) string {
	h := sha1.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func (t *TencentCOSTool) getCanonicalQueryString(params url.Values) string {
	if len(params) == 0 {
		return ""
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, strings.ToLower(k))
	}
	sort.Strings(keys)

	var result []string
	for _, k := range keys {
		result = append(result, fmt.Sprintf("%s=%s", url.QueryEscape(k), url.QueryEscape(params.Get(k))))
	}
	return strings.Join(result, "&")
}

func (t *TencentCOSTool) getCanonicalHeaders(headers http.Header) string {
	keys := make([]string, 0)
	for k := range headers {
		lowerKey := strings.ToLower(k)
		if lowerKey == "content-type" || lowerKey == "host" {
			keys = append(keys, lowerKey)
		}
	}
	sort.Strings(keys)

	var result []string
	for _, k := range keys {
		result = append(result, fmt.Sprintf("%s=%s", k, url.QueryEscape(headers.Get(k))))
	}
	return strings.Join(result, "&")
}

func (t *TencentCOSTool) getHeaderList(headers http.Header) string {
	keys := make([]string, 0)
	for k := range headers {
		lowerKey := strings.ToLower(k)
		if lowerKey == "content-type" || lowerKey == "host" {
			keys = append(keys, lowerKey)
		}
	}
	sort.Strings(keys)
	return strings.Join(keys, ";")
}

func (t *TencentCOSTool) getParamList(params url.Values) string {
	if len(params) == 0 {
		return ""
	}

	keys := make([]string, 0, len(params))
	for k := range params {
		keys = append(keys, strings.ToLower(k))
	}
	sort.Strings(keys)
	return strings.Join(keys, ";")
}

func init() {
	tool := NewTencentCOSTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register TencentCOS tool: %v", err))
	}
}
