package downloader

import (
    "auto-forge/pkg/utools"
    "crypto/tls"
    "fmt"
    "io"
    "net/http"
    neturl "net/url"
    "os"
    "path/filepath"
    "regexp"
    "strings"
    "time"
)

type FileDownloaderTool struct {
    *utools.BaseTool
}

func NewFileDownloaderTool() *FileDownloaderTool {
    metadata := &utools.ToolMetadata{
        Code:        "file_downloader",
        Name:        "文件下载器",
        Description: "从 URL 下载文件并生成文件对象，便于后续上传到对象存储等",
        Category:    "data",
        Version:     "1.0.0",
        Author:      "AutoForge",
        AICallable:  true,
        Tags:        []string{"download", "file", "http", "storage"},
        OutputFieldsSchema: map[string]utools.OutputFieldDef{
            "response": {
                Type:  "object",
                Label: "完整响应",
                Children: map[string]utools.OutputFieldDef{
                    "url":         {Type: "string", Label: "请求 URL"},
                    "status_code": {Type: "number", Label: "HTTP 状态码"},
                    "headers":     {Type: "object", Label: "响应头"},
                },
            },
            "file": {
                Type:  "object",
                Label: "下载后的文件对象",
                Children: map[string]utools.OutputFieldDef{
                    "type":      {Type: "string", Label: "对象类型"},
                    "path":      {Type: "string", Label: "本地路径"},
                    "filename":  {Type: "string", Label: "文件名"},
                    "size":      {Type: "number", Label: "字节大小"},
                    "mime_type": {Type: "string", Label: "MIME 类型"},
                },
            },
            "status_code": {Type: "number", Label: "HTTP 状态码"},
            "headers":     {Type: "object", Label: "响应头"},
            "path":        {Type: "string", Label: "文件路径（快捷访问）"},
            "filename":    {Type: "string", Label: "文件名（快捷访问）"},
            "mime_type":   {Type: "string", Label: "MIME 类型（快捷访问）"},
            "size":        {Type: "number", Label: "大小（快捷访问）"},
        },
    }

    schema := &utools.ConfigSchema{
        Type: "object",
        Properties: map[string]utools.PropertySchema{
            "url": {
                Type:        "string",
                Title:       "下载链接",
                Description: "目标文件的 URL",
                Format:      "uri",
            },
            "headers": {
                Type:        "object",
                Title:       "请求头",
                Description: "可选的 HTTP 请求头，键值对",
            },
            "filename": {
                Type:        "string",
                Title:       "保存文件名",
                Description: "自定义文件名，留空则自动推断",
            },
            "timeout": {
                Type:        "number",
                Title:       "超时时间(秒)",
                Description: "HTTP 请求超时，默认 60 秒",
                Default:     60.0,
            },
            "verify_ssl": {
                Type:        "boolean",
                Title:       "验证 SSL",
                Description: "是否验证 SSL 证书",
                Default:     true,
            },
            "follow_redirects": {
                Type:        "boolean",
                Title:       "跟随重定向",
                Description: "是否自动跟随 HTTP 重定向",
                Default:     true,
            },
        },
        Required: []string{"url"},
    }

    return &FileDownloaderTool{BaseTool: utools.NewBaseTool(metadata, schema)}
}

func (t *FileDownloaderTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
    start := time.Now()

    urlStr, _ := config["url"].(string)
    if strings.TrimSpace(urlStr) == "" {
        return &utools.ExecutionResult{Success: false, Message: "URL 不能为空", Error: "missing url", DurationMs: time.Since(start).Milliseconds()}, fmt.Errorf("url is required")
    }

    timeout := 60
    if v, ok := config["timeout"].(float64); ok && v > 0 {
        timeout = int(v)
    }
    verifySSL := true
    if v, ok := config["verify_ssl"].(bool); ok {
        verifySSL = v
    }
    followRedirects := true
    if v, ok := config["follow_redirects"].(bool); ok {
        followRedirects = v
    }

    client := &http.Client{
        Timeout: time.Duration(timeout) * time.Second,
        Transport: &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: !verifySSL}},
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            if !followRedirects {
                return http.ErrUseLastResponse
            }
            return nil
        },
    }

    req, err := http.NewRequest("GET", urlStr, nil)
    if err != nil {
        return &utools.ExecutionResult{Success: false, Message: "创建请求失败", Error: err.Error(), DurationMs: time.Since(start).Milliseconds()}, err
    }
    // 支持两种 headers 形态：对象或数组[{key,value}]
    if hdrs, ok := config["headers"].(map[string]interface{}); ok {
        for k, v := range hdrs {
            if s, ok := v.(string); ok {
                req.Header.Set(k, s)
            }
        }
    } else if arr, ok := config["headers"].([]interface{}); ok {
        for _, item := range arr {
            if m, ok := item.(map[string]interface{}); ok {
                k, _ := m["key"].(string)
                v, _ := m["value"].(string)
                if strings.TrimSpace(k) != "" {
                    req.Header.Set(k, v)
                }
            }
        }
    }

    resp, err := client.Do(req)
    if err != nil {
        return &utools.ExecutionResult{Success: false, Message: "请求失败", Error: err.Error(), DurationMs: time.Since(start).Milliseconds()}, err
    }
    defer resp.Body.Close()

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        body, _ := io.ReadAll(resp.Body)
        return &utools.ExecutionResult{
            Success:    false,
            Message:    fmt.Sprintf("HTTP 状态码: %d", resp.StatusCode),
            Error:      string(body),
            Output:     map[string]interface{}{"status_code": resp.StatusCode, "headers": resp.Header, "body": string(body)},
            DurationMs: time.Since(start).Milliseconds(),
        }, fmt.Errorf("unexpected status: %d", resp.StatusCode)
    }

    data, err := io.ReadAll(resp.Body)
    if err != nil {
        return &utools.ExecutionResult{Success: false, Message: "读取响应失败", Error: err.Error(), DurationMs: time.Since(start).Milliseconds()}, err
    }

    filename := strings.TrimSpace(getFilenameFromConfigOrHeader(config, resp.Header, urlStr))
    if filename == "" {
        filename = fmt.Sprintf("download-%d.bin", time.Now().UnixNano())
    }

    baseDir := "/tmp/workflow-files"
    execID := ""
    if ctx != nil && ctx.Metadata != nil {
        if v, ok := ctx.Metadata["execution_id"].(string); ok && v != "" {
            execID = v
        }
    }
    dir := baseDir
    if execID != "" {
        dir = filepath.Join(baseDir, execID)
    } else {
        dir = filepath.Join(baseDir, "shared")
    }
    if err := os.MkdirAll(dir, 0755); err != nil {
        return &utools.ExecutionResult{Success: false, Message: "创建目录失败", Error: err.Error(), DurationMs: time.Since(start).Milliseconds()}, err
    }

    filePath := filepath.Join(dir, filename)
    if err := os.WriteFile(filePath, data, 0644); err != nil {
        return &utools.ExecutionResult{Success: false, Message: "保存文件失败", Error: err.Error(), DurationMs: time.Since(start).Milliseconds()}, err
    }

    size := int64(len(data))
    mime := resp.Header.Get("Content-Type")
    if mime == "" {
        sniff := 512
        if len(data) < sniff {
            sniff = len(data)
        }
        mime = http.DetectContentType(data[:sniff])
    }

    fileObj := map[string]interface{}{
        "type":      "file",
        "path":      filePath,
        "filename":  filename,
        "size":      size,
        "mime_type": mime,
    }

    output := map[string]interface{}{
        "response": map[string]interface{}{
            "url":         urlStr,
            "status_code": resp.StatusCode,
            "headers":     resp.Header,
        },
        "file":        fileObj,
        "status_code": resp.StatusCode,
        "headers":     resp.Header,
        "path":        filePath,
        "filename":    filename,
        "mime_type":   mime,
        "size":        size,
    }

    return &utools.ExecutionResult{Success: true, Message: "下载成功", Output: output, DurationMs: time.Since(start).Milliseconds()}, nil
}

func getFilenameFromConfigOrHeader(config map[string]interface{}, headers http.Header, urlStr string) string {
    if v, ok := config["filename"].(string); ok && strings.TrimSpace(v) != "" {
        return strings.TrimSpace(v)
    }
    cd := headers.Get("Content-Disposition")
    if cd != "" {
        re := regexp.MustCompile(`(?i)filename\*=UTF-8''([^;]+)`)
        if m := re.FindStringSubmatch(cd); len(m) == 2 {
            if name, err := neturl.QueryUnescape(m[1]); err == nil && name != "" {
                return name
            }
        }
        re2 := regexp.MustCompile(`(?i)filename="?([^";]+)`)
        if m := re2.FindStringSubmatch(cd); len(m) == 2 {
            return m[1]
        }
    }
    if parsed, err := neturl.Parse(urlStr); err == nil {
        base := filepath.Base(parsed.Path)
        if base != "/" && base != "." && base != "" {
            return base
        }
    }
    return ""
}

func init() {
    tool := NewFileDownloaderTool()
    if err := utools.Register(tool); err != nil {
        panic(fmt.Sprintf("Failed to register File Downloader tool: %v", err))
    }
}
