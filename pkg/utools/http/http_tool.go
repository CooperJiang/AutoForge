package http

import (
	"auto-forge/pkg/utools"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)


type HTTPTool struct {
	*utools.BaseTool
}


func NewHTTPTool() *HTTPTool {
    metadata := &utools.ToolMetadata{
        Code:        "http_request",
        Name:        "HTTP 请求",
        Description: "发送 HTTP/HTTPS 请求，支持 GET、POST、PUT、DELETE 等方法",
        Category:    utools.CategoryNetwork,
        Version:     "1.0.0",
        Author:      "AutoForge",
        AICallable:  true,
        Tags:        []string{"http", "api", "request", "network"},
        OutputFieldsSchema: map[string]utools.OutputFieldDef{
            "status_code": {Type: "number", Label: "HTTP 状态码"},
            "headers":     {Type: "object", Label: "响应头"},
            "body":        {Type: "string", Label: "原始响应体"},
            "json":        {Type: "object", Label: "解析后的 JSON（若可解析）"},
        },
    }

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"url": {
				Type:        "string",
				Title:       "请求 URL",
				Description: "目标 API 地址",
				Format:      "uri",
			},
			"method": {
				Type:        "string",
				Title:       "请求方法",
				Description: "HTTP 请求方法",
				Default:     "GET",
				Enum:        []interface{}{"GET", "POST", "PUT", "DELETE", "PATCH", "HEAD", "OPTIONS"},
			},
			"headers": {
				Type:        "object",
				Title:       "请求头",
				Description: "HTTP 请求头，键值对形式",
			},
			"params": {
				Type:        "object",
				Title:       "查询参数",
				Description: "URL 查询参数，键值对形式",
			},
			"body": {
				Type:        "string",
				Title:       "请求体",
				Description: "请求体内容（JSON 字符串或其他格式）",
			},
			"timeout": {
				Type:        "number",
				Title:       "超时时间",
				Description: "请求超时时间（秒）",
				Default:     30.0,
				Minimum:     func() *float64 { v := 1.0; return &v }(),
				Maximum:     func() *float64 { v := 300.0; return &v }(),
			},
			"follow_redirects": {
				Type:        "boolean",
				Title:       "跟随重定向",
				Description: "是否自动跟随 HTTP 重定向",
				Default:     true,
			},
		},
		Required: []string{"url", "method"},
	}

	return &HTTPTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}


func (t *HTTPTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()


	url, _ := config["url"].(string)
	method, _ := config["method"].(string)


	timeout := 30.0
	if timeoutVal, ok := config["timeout"].(float64); ok {
		timeout = timeoutVal
	}


	var bodyReader io.Reader
	if body, ok := config["body"].(string); ok && body != "" {
		bodyReader = bytes.NewBufferString(body)
	}


	req, err := http.NewRequestWithContext(ctx.Context, method, url, bodyReader)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "创建 HTTP 请求失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}


	if headers, ok := config["headers"].(map[string]interface{}); ok {
		for key, value := range headers {
			if strValue, ok := value.(string); ok {
				req.Header.Set(key, strValue)
			}
		}
	}


	if params, ok := config["params"].(map[string]interface{}); ok {
		q := req.URL.Query()
		for key, value := range params {
			q.Add(key, fmt.Sprintf("%v", value))
		}
		req.URL.RawQuery = q.Encode()
	}


	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if followRedirects, ok := config["follow_redirects"].(bool); ok && !followRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}


	resp, err := client.Do(req)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "HTTP 请求失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}
	defer resp.Body.Close()


	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &utools.ExecutionResult{
			Success:      false,
			Message:      "读取响应体失败",
			Error:        err.Error(),
			StatusCode:   resp.StatusCode,
			DurationMs:   time.Since(startTime).Milliseconds(),
		}, err
	}

	responseBody := string(bodyBytes)


	var jsonResponse map[string]interface{}
	_ = json.Unmarshal(bodyBytes, &jsonResponse)


	success := resp.StatusCode >= 200 && resp.StatusCode < 300

	output := map[string]interface{}{
		"status_code": resp.StatusCode,
		"headers":     resp.Header,
		"body":        responseBody,
	}

	if jsonResponse != nil {
		output["json"] = jsonResponse
	}

	message := fmt.Sprintf("HTTP 请求完成，状态码: %d", resp.StatusCode)
	if !success {
		message = fmt.Sprintf("HTTP 请求失败，状态码: %d", resp.StatusCode)
	}

	return &utools.ExecutionResult{
		Success:      success,
		Message:      message,
		Output:       output,
		StatusCode:   resp.StatusCode,
		ResponseBody: responseBody,
		DurationMs:   time.Since(startTime).Milliseconds(),
	}, nil
}


func init() {
	tool := NewHTTPTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register HTTP tool: %v", err))
	}
}
