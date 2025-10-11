package task

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"auto-forge/internal/models"
	"auto-forge/pkg/common"
	"auto-forge/pkg/errors"
	"auto-forge/pkg/utools"
	"time"

	"github.com/gin-gonic/gin"
)

// TestTaskRequest 测试任务请求
type TestTaskRequest struct {
	ToolCode string                `json:"tool_code"`
	URL      string                `json:"url"`
	Method   string                `json:"method"`
	Headers  models.KeyValueArray  `json:"headers"`
	Params   models.KeyValueArray  `json:"params"`
	Body     string                `json:"body"`
	Config   map[string]interface{} `json:"config"`
}

// TestTaskResponse 测试任务响应
type TestTaskResponse struct {
	Success        bool   `json:"success"`
	StatusCode     int    `json:"status_code"`
	ResponseBody   string `json:"response_body"`
	DurationMs     int64  `json:"duration_ms"`
	ErrorMessage   string `json:"error_message,omitempty"`
}

// TestTask 测试任务配置
func TestTask(c *gin.Context) {
	req, err := common.ValidateRequest[TestTaskRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}

	// 判断是否使用新的工具系统（有 tool_code 或 config）
	if req.ToolCode != "" || req.Config != nil {
		testWithToolSystem(c, req)
		return
	}

	// 兼容旧的 HTTP 请求测试
	testHTTPRequest(c, req)
}

// testWithToolSystem 使用工具系统测试
func testWithToolSystem(c *gin.Context, req *TestTaskRequest) {
	// 获取工具代码
	toolCode := req.ToolCode
	if toolCode == "" {
		// 尝试从 config 推断
		if req.URL != "" {
			toolCode = "http_request"
		}
	}

	if toolCode == "" {
		errors.HandleError(c, errors.NewValidationError("tool_code", "工具代码不能为空"))
		return
	}

	// 获取工具
	tool, err := utools.Get(toolCode)
	if err != nil {
		errors.HandleError(c, errors.NewValidationError("tool_code", "工具不存在: "+toolCode))
		return
	}

	// 构建配置
	var config map[string]interface{}
	if req.Config != nil {
		config = req.Config
	} else {
		// 从旧格式转换
		config = map[string]interface{}{
			"url":    req.URL,
			"method": req.Method,
			"body":   req.Body,
		}

		// 处理 headers
		if len(req.Headers) > 0 {
			headers := make(map[string]string)
			for _, h := range req.Headers {
				if h.Key != "" {
					headers[h.Key] = h.Value
				}
			}
			headersJSON, _ := json.Marshal(headers)
			config["headers"] = string(headersJSON)
		}
	}

	// 执行工具
	ctx := &utools.ExecutionContext{
		Context: context.Background(),
		TaskID:  "test",
		UserID:  c.GetString("user_id"),
	}

	result, err := tool.Execute(ctx, config)
	if err != nil {
		errors.ResponseSuccess(c, TestTaskResponse{
			Success:      false,
			ErrorMessage: "执行失败: " + err.Error(),
		}, "测试完成")
		return
	}

	// 返回结果
	errors.ResponseSuccess(c, TestTaskResponse{
		Success:      result.Success,
		StatusCode:   result.StatusCode,
		ResponseBody: result.Message + "\n\n详细信息:\n" + result.ResponseBody,
		DurationMs:   result.DurationMs,
		ErrorMessage: result.Error,
	}, "测试完成")
}

// testHTTPRequest 测试 HTTP 请求（兼容旧接口）
func testHTTPRequest(c *gin.Context, req *TestTaskRequest) {
	// 清理URL
	targetURL := strings.TrimSpace(req.URL)

	// 处理请求参数
	if len(req.Params) > 0 {
		params := url.Values{}
		for _, param := range req.Params {
			if param.Key != "" {
				params.Add(param.Key, param.Value)
			}
		}

		if len(params) > 0 {
			if strings.Contains(targetURL, "?") {
				targetURL = targetURL + "&" + params.Encode()
			} else {
				targetURL = targetURL + "?" + params.Encode()
			}
		}
	}

	// 记录开始时间
	startTime := time.Now()

	// 创建HTTP请求（支持 body）
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, targetURL, bodyReader)
	if err != nil {
		errors.ResponseSuccess(c, TestTaskResponse{
			Success:      false,
			ErrorMessage: "创建请求失败: " + err.Error(),
		}, "测试完成")
		return
	}

	// 添加请求头
	for _, header := range req.Headers {
		if header.Key != "" {
			httpReq.Header.Set(header.Key, header.Value)
		}
	}

	// 发送请求
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(httpReq)
	durationMs := time.Since(startTime).Milliseconds()

	if err != nil {
		errors.ResponseSuccess(c, TestTaskResponse{
			Success:      false,
			DurationMs:   durationMs,
			ErrorMessage: "请求失败: " + err.Error(),
		}, "测试完成")
		return
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		errors.ResponseSuccess(c, TestTaskResponse{
			Success:      false,
			StatusCode:   resp.StatusCode,
			DurationMs:   durationMs,
			ErrorMessage: "读取响应失败: " + err.Error(),
		}, "测试完成")
		return
	}

	// 转换为字符串（不截断）
	bodyStr := string(bodyBytes)

	// 判断是否成功
	success := resp.StatusCode >= 200 && resp.StatusCode < 300

	errors.ResponseSuccess(c, TestTaskResponse{
		Success:      success,
		StatusCode:   resp.StatusCode,
		ResponseBody: bodyStr,
		DurationMs:   durationMs,
	}, "测试完成")
}
