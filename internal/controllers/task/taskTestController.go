package task

import (
	"auto-forge/internal/models"
	"auto-forge/pkg/common"
	"auto-forge/pkg/errors"
	"auto-forge/pkg/utools"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)


type TestTaskRequest struct {
	ToolCode string                 `json:"tool_code"`
	URL      string                 `json:"url"`
	Method   string                 `json:"method"`
	Headers  models.KeyValueArray   `json:"headers"`
	Params   models.KeyValueArray   `json:"params"`
	Body     string                 `json:"body"`
	Config   map[string]interface{} `json:"config"`
}


type TestTaskResponse struct {
	Success      bool                   `json:"success"`
	StatusCode   int                    `json:"status_code"`
	ResponseBody string                 `json:"response_body"`
	DurationMs   int64                  `json:"duration_ms"`
	ErrorMessage string                 `json:"error_message,omitempty"`
	Output       map[string]interface{} `json:"output,omitempty"`
	Message      string                 `json:"message,omitempty"`
}


func TestTask(c *gin.Context) {
	req, err := common.ValidateRequest[TestTaskRequest](c)
	if err != nil {
		errors.HandleError(c, err)
		return
	}


	if req.ToolCode != "" || req.Config != nil {
		testWithToolSystem(c, req)
		return
	}


	testHTTPRequest(c, req)
}


func testWithToolSystem(c *gin.Context, req *TestTaskRequest) {

	toolCode := req.ToolCode
	if toolCode == "" {

		if req.URL != "" {
			toolCode = "http_request"
		}
	}

	if toolCode == "" {
		errors.HandleError(c, errors.NewValidationError("tool_code", "工具代码不能为空"))
		return
	}


	tool, err := utools.Get(toolCode)
	if err != nil {
		errors.HandleError(c, errors.NewValidationError("tool_code", "工具不存在: "+toolCode))
		return
	}


	var config map[string]interface{}
	if req.Config != nil {
		config = req.Config
	} else {

		config = map[string]interface{}{
			"url":    req.URL,
			"method": req.Method,
			"body":   req.Body,
		}


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


	ctx := &utools.ExecutionContext{
		Context: context.Background(),
		TaskID:  "test",
		UserID:  c.GetString("user_id"),
	}

	result, err := tool.Execute(ctx, config)
	if err != nil {
		resp := TestTaskResponse{
			Success:      false,
			ErrorMessage: "执行失败: " + err.Error(),
		}
		if result != nil {
			resp.StatusCode = result.StatusCode
			resp.DurationMs = result.DurationMs
			resp.Output = result.Output
			resp.Message = result.Message
		}
		errors.ResponseSuccess(c, resp, "测试完成")
		return
	}


	response := TestTaskResponse{
		Success:      result.Success,
		StatusCode:   result.StatusCode,
		ResponseBody: result.Message + "\n\n详细信息:\n" + result.ResponseBody,
		DurationMs:   result.DurationMs,
		ErrorMessage: result.Error,
		Output:       result.Output,
		Message:      result.Message,
	}
	errors.ResponseSuccess(c, response, "测试完成")
}


func testHTTPRequest(c *gin.Context, req *TestTaskRequest) {

	targetURL := strings.TrimSpace(req.URL)


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


	startTime := time.Now()


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


	for _, header := range req.Headers {
		if header.Key != "" {
			httpReq.Header.Set(header.Key, header.Value)
		}
	}


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


	bodyStr := string(bodyBytes)


	success := resp.StatusCode >= 200 && resp.StatusCode < 300

	errors.ResponseSuccess(c, TestTaskResponse{
		Success:      success,
		StatusCode:   resp.StatusCode,
		ResponseBody: bodyStr,
		DurationMs:   durationMs,
	}, "测试完成")
}
