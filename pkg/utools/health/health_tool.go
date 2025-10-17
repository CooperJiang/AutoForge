package health

import (
	"auto-forge/pkg/utools"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"
)


type HealthCheckerTool struct {
	*utools.BaseTool
}


func NewHealthCheckerTool() *HealthCheckerTool {
    metadata := &utools.ToolMetadata{
        Code:        "health_checker",
        Name:        "网站健康检查",
        Description: "检查网站/API的可用性、响应时间、状态码和内容匹配",
        Category:    utools.CategoryMonitoring,
        Version:     "1.0.0",
        Author:      "AutoForge",
        AICallable:  true,
        Tags:        []string{"health", "monitoring", "check", "uptime", "ssl"},
        OutputFieldsSchema: map[string]utools.OutputFieldDef{
            "url":           {Type: "string", Label: "检查 URL"},
            "status_code":   {Type: "number", Label: "HTTP 状态码"},
            "response_time": {Type: "number", Label: "响应时间(ms)"},
            "headers":       {Type: "object", Label: "响应头"},
            "response_body": {Type: "string", Label: "响应体"},
            "ssl": {
                Type:  "object",
                Label: "SSL 证书信息",
                Children: map[string]utools.OutputFieldDef{
                    "subject":       {Type: "string", Label: "主题"},
                    "issuer":        {Type: "string", Label: "签发者"},
                    "not_before":    {Type: "string", Label: "生效时间"},
                    "not_after":     {Type: "string", Label: "过期时间"},
                    "days_to_expiry": {Type: "number", Label: "剩余天数"},
                },
            },
            "issues":   {Type: "array", Label: "问题列表"},
            "warnings": {Type: "array", Label: "警告列表"},
        },
    }

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"url": {
				Type:        "string",
				Title:       "检查 URL",
				Description: "要检查的网站或 API 地址",
				Format:      "uri",
			},
			"method": {
				Type:        "string",
				Title:       "请求方法",
				Description: "HTTP 请求方法",
				Default:     "GET",
				Enum:        []interface{}{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH"},
			},
			"headers": {
				Type:        "string",
				Title:       "请求头",
				Description: "自定义 HTTP 请求头 (JSON 格式)",
			},
			"body": {
				Type:        "string",
				Title:       "请求体",
				Description: "请求体内容 (JSON 或文本)",
			},
			"timeout": {
				Type:        "number",
				Title:       "超时时间",
				Description: "请求超时时间（秒）",
				Default:     10.0,
				Minimum:     func() *float64 { v := 1.0; return &v }(),
				Maximum:     func() *float64 { v := 60.0; return &v }(),
			},
			"expected_status": {
				Type:        "number",
				Title:       "期望状态码",
				Description: "期望的 HTTP 状态码，0 表示任意 2xx 状态码",
				Default:     200.0,
			},
			"expected_content": {
				Type:        "string",
				Title:       "期望内容",
				Description: "响应体中应包含的内容（支持正则表达式）",
			},
			"use_regex": {
				Type:        "boolean",
				Title:       "使用正则匹配",
				Description: "期望内容是否使用正则表达式匹配",
				Default:     false,
			},
			"check_ssl": {
				Type:        "boolean",
				Title:       "检查 SSL 证书",
				Description: "是否检查 SSL 证书有效期",
				Default:     true,
			},
			"ssl_warning_days": {
				Type:        "number",
				Title:       "SSL 到期告警天数",
				Description: "SSL 证书在多少天内到期时告警",
				Default:     30.0,
			},
			"follow_redirects": {
				Type:        "boolean",
				Title:       "跟随重定向",
				Description: "是否自动跟随 HTTP 重定向",
				Default:     true,
			},
			"verify_ssl": {
				Type:        "boolean",
				Title:       "验证 SSL 证书",
				Description: "是否验证 SSL 证书有效性",
				Default:     true,
			},
		},
		Required: []string{"url"},
	}

	return &HealthCheckerTool{
		BaseTool: utools.NewBaseTool(metadata, schema),
	}
}


func (t *HealthCheckerTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	startTime := time.Now()


	url, _ := config["url"].(string)
	method, _ := config["method"].(string)
	if method == "" {
		method = "GET"
	}

	headersStr, _ := config["headers"].(string)
	bodyStr, _ := config["body"].(string)

	timeout := 10.0
	if timeoutVal, ok := config["timeout"].(float64); ok {
		timeout = timeoutVal
	}

	expectedStatus := 200
	if statusVal, ok := config["expected_status"].(float64); ok {
		expectedStatus = int(statusVal)
	}

	expectedContent, _ := config["expected_content"].(string)
	useRegex := false
	if useRegexVal, ok := config["use_regex"].(bool); ok {
		useRegex = useRegexVal
	}

	checkSSL := true
	if checkSSLVal, ok := config["check_ssl"].(bool); ok {
		checkSSL = checkSSLVal
	}

	sslWarningDays := 30
	if sslDaysVal, ok := config["ssl_warning_days"].(float64); ok {
		sslWarningDays = int(sslDaysVal)
	}

	followRedirects := true
	if followRedirectsVal, ok := config["follow_redirects"].(bool); ok {
		followRedirects = followRedirectsVal
	}

	verifySSL := true
	if verifySSLVal, ok := config["verify_ssl"].(bool); ok {
		verifySSL = verifySSLVal
	}


	client := &http.Client{
		Timeout: time.Duration(timeout) * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: !verifySSL,
			},
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if !followRedirects {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}


	var bodyReader io.Reader
	if bodyStr != "" {
		bodyReader = strings.NewReader(bodyStr)
	}


	req, err := http.NewRequestWithContext(ctx.Context, method, url, bodyReader)
	if err != nil {
		return &utools.ExecutionResult{
			Success:    false,
			Message:    "创建请求失败",
			Error:      err.Error(),
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}


	if headersStr != "" {
		var headersMap map[string]interface{}
		if err := json.Unmarshal([]byte(headersStr), &headersMap); err == nil {
			for key, value := range headersMap {
				if strValue, ok := value.(string); ok {
					req.Header.Set(key, strValue)
				}
			}
		}
	}


	if bodyStr != "" && req.Header.Get("Content-Type") == "" {

		var jsonTest interface{}
		if json.Unmarshal([]byte(bodyStr), &jsonTest) == nil {
			req.Header.Set("Content-Type", "application/json")
		} else {
			req.Header.Set("Content-Type", "text/plain")
		}
	}


	reqStartTime := time.Now()
	resp, err := client.Do(req)
	responseTime := time.Since(reqStartTime).Milliseconds()

	if err != nil {
		return &utools.ExecutionResult{
			Success: false,
			Message: fmt.Sprintf("请求失败: %s", err.Error()),
			Error:   err.Error(),
			Output: map[string]interface{}{
				"url":           url,
				"error":         err.Error(),
				"response_time": responseTime,
			},
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}
	defer resp.Body.Close()


	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return &utools.ExecutionResult{
			Success: false,
			Message: "读取响应失败",
			Error:   err.Error(),
			Output: map[string]interface{}{
				"url":           url,
				"status_code":   resp.StatusCode,
				"response_time": responseTime,
			},
			DurationMs: time.Since(startTime).Milliseconds(),
		}, err
	}

	responseBody := string(bodyBytes)


	var issues []string
	var warnings []string


	statusOK := false
	if expectedStatus == 0 {

		statusOK = resp.StatusCode >= 200 && resp.StatusCode < 300
	} else {
		statusOK = resp.StatusCode == expectedStatus
	}

	if !statusOK {
		issues = append(issues, fmt.Sprintf("状态码不符合预期: 期望 %d, 实际 %d", expectedStatus, resp.StatusCode))
	}


	if expectedContent != "" {
		if useRegex {
			matched, err := regexp.MatchString(expectedContent, responseBody)
			if err != nil {
				issues = append(issues, fmt.Sprintf("正则表达式错误: %s", err.Error()))
			} else if !matched {
				issues = append(issues, "响应内容不匹配正则表达式")
			}
		} else {
			if !strings.Contains(responseBody, expectedContent) {
				issues = append(issues, "响应内容中未找到期望的文本")
			}
		}
	}


	var sslInfo map[string]interface{}
	if checkSSL && strings.HasPrefix(url, "https://") {
		cert := resp.TLS.PeerCertificates[0]
		daysUntilExpiry := int(time.Until(cert.NotAfter).Hours() / 24)

		sslInfo = map[string]interface{}{
			"subject":      cert.Subject.CommonName,
			"issuer":       cert.Issuer.CommonName,
			"not_before":   cert.NotBefore.Format(time.RFC3339),
			"not_after":    cert.NotAfter.Format(time.RFC3339),
			"days_to_expiry": daysUntilExpiry,
		}

		if daysUntilExpiry < 0 {
			issues = append(issues, fmt.Sprintf("SSL 证书已过期 %d 天", -daysUntilExpiry))
		} else if daysUntilExpiry < sslWarningDays {
			warnings = append(warnings, fmt.Sprintf("SSL 证书将在 %d 天后过期", daysUntilExpiry))
		}
	}


	output := map[string]interface{}{
		"url":           url,
		"status_code":   resp.StatusCode,
		"response_time": responseTime,
		"body_size":     len(bodyBytes),
		"headers":       resp.Header,
		"response_body": responseBody,
	}

	if sslInfo != nil {
		output["ssl"] = sslInfo
	}

	if len(warnings) > 0 {
		output["warnings"] = warnings
	}


	var statusReport []string


	if len(issues) == 0 {
		statusReport = append(statusReport, fmt.Sprintf("✓ 网站正常 (响应时间: %dms)", responseTime))
	} else {
		statusReport = append(statusReport, fmt.Sprintf("✗ 网站异常: %s", strings.Join(issues, "; ")))
	}


	if sslInfo != nil {
		daysLeft := sslInfo["days_to_expiry"].(int)
		if daysLeft < 0 {
			statusReport = append(statusReport, fmt.Sprintf("⚠️ SSL 证书已过期 %d 天", -daysLeft))
		} else if daysLeft < sslWarningDays {
			statusReport = append(statusReport, fmt.Sprintf("⚠️ SSL 证书剩余 %d 天（已达警戒线 %d 天）", daysLeft, sslWarningDays))
		} else {
			statusReport = append(statusReport, fmt.Sprintf("✓ SSL 证书正常，剩余 %d 天", daysLeft))
		}
	}


	success := len(issues) == 0
	message := strings.Join(statusReport, "\n")

	if len(issues) > 0 {
		output["issues"] = issues
	}
	if len(warnings) > 0 {
		output["warnings"] = warnings
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
	tool := NewHealthCheckerTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("Failed to register HealthChecker tool: %v", err))
	}
}
