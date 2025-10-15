package jsontransform

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"auto-forge/pkg/utools"
	"github.com/dop251/goja"
)

type JSONTransformTool struct {
	*utools.BaseTool
}

func NewJSONTransformTool() *JSONTransformTool {
	metadata := &utools.ToolMetadata{
		Code:        "json_transform",
		Name:        "JSON 转换",
		Description: "基于变量系统和 JS 表达式动态转换数据",
		Category:    "data",
		Version:     "2.0.0",
		Author:      "AutoForge",
		AICallable:  true,
		Tags:        []string{"json", "transform", "expression", "mapping"},
	}

	schema := &utools.ConfigSchema{
		Type: "object",
		Properties: map[string]utools.PropertySchema{
			"data_source": {
				Type:        "string",
				Title:       "数据来源",
				Description: "支持变量引用，例如 {{nodes.node_xxx.output}}",
			},
			"expression": {
				Type:        "string",
				Title:       "JS 表达式",
				Description: "可直接编写如 data.map(item => item.url) 的表达式，形参 data/ctx 已注入",
			},
			"output_name": {
				Type:        "string",
				Title:       "输出字段名称",
				Description: "最终结果在节点输出中的键名",
				Default:     "result",
			},
			"timeout_ms": {
				Type:        "number",
				Title:       "执行超时 (ms)",
				Description: "防止表达式长时间执行，默认 1500ms",
				Default:     1500.0,
			},
			"sample_json": {
				Type:        "string",
				Title:       "样例 JSON (可选)",
				Description: "仅用于前端预览，运行时不读取",
			},
		},
		Required: []string{"data_source", "expression"},
	}

	return &JSONTransformTool{BaseTool: utools.NewBaseTool(metadata, schema)}
}

func (t *JSONTransformTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
	start := time.Now()

	if err := t.Validate(config); err != nil {
		return failureResult(start, "配置校验失败", err)
	}

	dataSourceRaw, _ := config["data_source"].(string)
	expressionRaw, _ := config["expression"].(string)

	if strings.TrimSpace(dataSourceRaw) == "" {
		return failureResult(start, "数据来源不能为空", errors.New("data_source empty"))
	}
	if strings.TrimSpace(expressionRaw) == "" {
		return failureResult(start, "表达式不能为空", errors.New("expression empty"))
	}

	outputName := "result"
	if v, ok := config["output_name"].(string); ok && strings.TrimSpace(v) != "" {
		outputName = strings.TrimSpace(v)
	}

	timeout := 1500
	if v, ok := config["timeout_ms"].(float64); ok && v > 0 {
		timeout = int(v)
	}

	resolvedData, err := resolveDataSource(dataSourceRaw, ctx)
	if err != nil {
		return failureResult(start, "解析数据来源失败", err)
	}

	expression := strings.TrimSpace(expressionRaw)
	jsResult, err := executeExpression(resolvedData, ctx, expression, timeout)
	if err != nil {
		return failureResult(start, "表达式执行失败", err)
	}

	output := map[string]interface{}{}
	output[outputName] = jsResult
	if previewBytes, err := json.MarshalIndent(jsResult, "", "  "); err == nil {
		output["preview"] = string(previewBytes)
	}

	return &utools.ExecutionResult{
		Success:    true,
		Message:    "转换成功",
		Output:     output,
		DurationMs: time.Since(start).Milliseconds(),
	}, nil
}

func failureResult(start time.Time, message string, err error) (*utools.ExecutionResult, error) {
	result := &utools.ExecutionResult{
		Success:    false,
		Message:    message,
		Error:      err.Error(),
		DurationMs: time.Since(start).Milliseconds(),
	}
	return result, err
}

func resolveDataSource(raw string, ctx *utools.ExecutionContext) (interface{}, error) {
	trimmed := strings.TrimSpace(raw)

	if strings.HasPrefix(trimmed, "{{") && strings.HasSuffix(trimmed, "}}") {
		path := strings.TrimSpace(trimmed[2 : len(trimmed)-2])
		return lookupContextPath(path, ctx)
	}

	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		var v interface{}
		if err := json.Unmarshal([]byte(trimmed), &v); err == nil {
			return v, nil
		}
	}

	return trimmed, nil
}

func lookupContextPath(path string, ctx *utools.ExecutionContext) (interface{}, error) {
	if ctx == nil {
		return nil, fmt.Errorf("执行上下文为空，无法解析 %s", path)
	}

	tokens, err := tokenizePath(path)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("变量路径 %s 无效", path)
	}

	var current interface{}
	if key, ok := tokens[0].(string); ok {
		if ctx.Variables != nil {
			current = ctx.Variables[key]
		}
		if current == nil && ctx.Metadata != nil {
			current = ctx.Metadata[key]
		}
		if current == nil {
			return nil, fmt.Errorf("未找到变量 %s", path)
		}
	} else {
		return nil, fmt.Errorf("变量路径 %s 无法解析", path)
	}

	for _, token := range tokens[1:] {
		switch seg := token.(type) {
		case string:
			m, ok := current.(map[string]interface{})
			if !ok {
				return nil, fmt.Errorf("无法在 %T 上访问字段 %s", current, seg)
			}
			current = m[seg]
		case int:
			arr, ok := current.([]interface{})
			if !ok {
				return nil, fmt.Errorf("无法在 %T 上访问索引 %d", current, seg)
			}
			if seg < 0 || seg >= len(arr) {
				return nil, fmt.Errorf("数组索引越界: %d", seg)
			}
			current = arr[seg]
		default:
			return nil, fmt.Errorf("不支持的路径片段: %v", token)
		}

		if current == nil {
			return nil, fmt.Errorf("变量 %s 中的片段解析为 nil", path)
		}
	}

	return current, nil
}

func tokenizePath(path string) ([]interface{}, error) {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return nil, fmt.Errorf("变量路径为空")
	}

	rawSegments := strings.Split(trimmed, ".")
	var tokens []interface{}

	for _, segment := range rawSegments {
		if segment == "" {
			continue
		}
		expanded, err := expandToken(segment)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, expanded...)
	}

	return tokens, nil
}

func expandToken(token string) ([]interface{}, error) {
	var segments []interface{}
	remaining := token

	for {
		idx := strings.IndexByte(remaining, '[')
		if idx == -1 {
			if remaining != "" {
				segments = append(segments, remaining)
			}
			break
		}

		if idx > 0 {
			segments = append(segments, remaining[:idx])
		}

		remaining = remaining[idx+1:]
		end := strings.IndexByte(remaining, ']')
		if end == -1 {
			return nil, fmt.Errorf("路径 %s 缺少 ]", token)
		}

		indexToken := strings.TrimSpace(remaining[:end])
		remaining = remaining[end+1:]

		if indexToken == "" {
			return nil, fmt.Errorf("路径 %s 中存在空索引", token)
		}

		idxVal, err := parseIndex(indexToken)
		if err != nil {
			return nil, err
		}
		segments = append(segments, idxVal)
	}

	return segments, nil
}

func parseIndex(raw string) (int, error) {
	var idx int
	_, err := fmt.Sscanf(raw, "%d", &idx)
	if err != nil {
		return 0, fmt.Errorf("索引 %s 不是有效数字", raw)
	}
	return idx, nil
}

func executeExpression(data interface{}, execCtx *utools.ExecutionContext, expression string, timeoutMs int) (interface{}, error) {
	runtime := goja.New()

	if err := runtime.Set("data", data); err != nil {
		return nil, fmt.Errorf("设置 data 失败: %w", err)
	}

	ctxVars := map[string]interface{}{}
	if execCtx != nil {
		if execCtx.Variables != nil {
			for k, v := range execCtx.Variables {
				ctxVars[k] = v
			}
		}
		if execCtx.Metadata != nil {
			meta := map[string]interface{}{}
			for k, v := range execCtx.Metadata {
				meta[k] = v
			}
			if len(meta) > 0 {
				ctxVars["meta"] = meta
			}
		}
	}

	if err := runtime.Set("ctx", ctxVars); err != nil {
		return nil, fmt.Errorf("设置 ctx 失败: %w", err)
	}

	interruption := fmt.Errorf("表达式执行超时 (%dms)", timeoutMs)
	timer := time.AfterFunc(time.Duration(timeoutMs)*time.Millisecond, func() {
		runtime.Interrupt(interruption)
	})
	defer timer.Stop()

	done := make(chan struct{})
	if execCtx != nil && execCtx.Context != nil {
		go func(parent context.Context) {
			select {
			case <-parent.Done():
				runtime.Interrupt(parent.Err())
			case <-done:
			}
		}(execCtx.Context)
	}

	wrapped := fmt.Sprintf("(function(){ return (%s); })()", expression)
	value, err := runtime.RunString(wrapped)
	close(done)
	if err != nil {
		switch e := err.(type) {
		case *goja.Exception:
			return nil, errors.New(e.String())
		case *goja.InterruptedError:
			if val := e.Value(); val != nil {
				switch v := val.(type) {
				case error:
					return nil, v
				case string:
					return nil, errors.New(v)
				default:
					return nil, fmt.Errorf("%v", v)
				}
			}
			return nil, interruption
		default:
			return nil, err
		}
	}

	return value.Export(), nil
}

func init() {
	tool := NewJSONTransformTool()
	if err := utools.Register(tool); err != nil {
		panic(fmt.Sprintf("failed to register json transform tool: %v", err))
	}
}
