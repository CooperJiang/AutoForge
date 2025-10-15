package context

import (
    gocontext "context"
    "encoding/json"
    "fmt"
    "time"

    "auto-forge/pkg/cache"
    "auto-forge/pkg/utools"
)

// RedisContextTool 提供基于缓存(优先Redis,降级内存)的简单KV读写/删除能力
type RedisContextTool struct {
    *utools.BaseTool
}

// NewRedisContextTool 创建工具实例
func NewRedisContextTool() *RedisContextTool {
    metadata := &utools.ToolMetadata{
        Code:        "redis_context",
        Name:        "Redis 存储器",
        Description: "使用 Redis/内存缓存存取上下文数据，支持 get/set/delete 并可选 TTL",
        Category:    "data",
        Version:     "1.0.0",
        Author:      "AutoForge",
        AICallable:  false,
        Tags:        []string{"redis", "storage", "context", "state"},
    }

    schema := &utools.ConfigSchema{
        Type: "object",
        Properties: map[string]utools.PropertySchema{
            "action": {
                Type:        "string",
                Title:       "操作类型",
                Description: "选择执行的操作：get(读取)、set(写入)、delete(删除)",
                Default:     "get",
                Enum:        []interface{}{"get", "set", "delete"},
            },
            "key": {
                Type:        "string",
                Title:       "键名",
                Description: "上下文键名，例如：session:user_123",
                MinLength:   func() *int { v := 1; return &v }(),
            },
            "value": {
                Type:        "string",
                Title:       "值",
                Description: "当 action=set 时写入的值，支持JSON字符串或普通文本",
            },
            "ttl_seconds": {
                Type:        "number",
                Title:       "过期时间(秒)",
                Description: "0或留空表示不过期",
                Default:     0.0,
                Minimum:     func() *float64 { v := 0.0; return &v }(),
            },
        },
        Required: []string{"action", "key"},
    }

    return &RedisContextTool{BaseTool: utools.NewBaseTool(metadata, schema)}
}

// Validate 额外校验
func (t *RedisContextTool) Validate(config map[string]interface{}) error {
    if err := t.BaseTool.Validate(config); err != nil {
        return err
    }

    if action, _ := config["action"].(string); action == "set" {
        if _, ok := config["value"].(string); !ok {
            return &utools.ValidationError{Field: "value", Message: "value is required for action=set"}
        }
    }
    return nil
}

// Execute 执行操作
func (t *RedisContextTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
    start := time.Now()
    // Ensure cache initialized
    _ = cache.GetCache()

    action, _ := config["action"].(string)
    key, _ := config["key"].(string)

    switch action {
    case "get":
        val, err := cache.Get(key)
        if err != nil {
            return &utools.ExecutionResult{
                Success:    false,
                Message:    "读取失败",
                Error:      err.Error(),
                DurationMs: time.Since(start).Milliseconds(),
                Output: map[string]interface{}{
                    "action": action,
                    "key":    key,
                },
            }, err
        }

        output := map[string]interface{}{
            "action": action,
            "key":    key,
            "value":  val,
        }
        // 尝试解析JSON
        var anyJSON interface{}
        if err := json.Unmarshal([]byte(val), &anyJSON); err == nil {
            output["json"] = anyJSON
        }
        // 读取TTL
        if ttl, err := cache.TTL(key); err == nil {
            // -1 表示无过期，其他为剩余时间
            output["ttl_ms"] = ttl.Milliseconds()
        }

        return &utools.ExecutionResult{
            Success:    true,
            Message:    "读取成功",
            Output:     output,
            DurationMs: time.Since(start).Milliseconds(),
        }, nil

    case "set":
        val, _ := config["value"].(string)
        ttlSeconds := 0.0
        if v, ok := config["ttl_seconds"].(float64); ok {
            ttlSeconds = v
        }
        exp := time.Duration(ttlSeconds) * time.Second
        if err := cache.Set(key, val, exp); err != nil {
            return &utools.ExecutionResult{
                Success:    false,
                Message:    "写入失败",
                Error:      err.Error(),
                DurationMs: time.Since(start).Milliseconds(),
                Output: map[string]interface{}{
                    "action": action,
                    "key":    key,
                },
            }, err
        }
        // 返回TTL（如果有）
        var ttlMs int64 = -1
        if exp > 0 {
            ttlMs = exp.Milliseconds()
        }
        return &utools.ExecutionResult{
            Success:    true,
            Message:    "写入成功",
            Output: map[string]interface{}{
                "action":      action,
                "key":         key,
                "value":       val,
                "ttl_ms":      ttlMs,
                "redis_enabled": cache.IsRedisEnabled(),
            },
            DurationMs: time.Since(start).Milliseconds(),
        }, nil

    case "delete":
        if err := cache.Del(key); err != nil {
            return &utools.ExecutionResult{
                Success:    false,
                Message:    "删除失败",
                Error:      err.Error(),
                DurationMs: time.Since(start).Milliseconds(),
                Output: map[string]interface{}{
                    "action": action,
                    "key":    key,
                },
            }, err
        }
        return &utools.ExecutionResult{
            Success:    true,
            Message:    "删除成功",
            Output: map[string]interface{}{
                "action": action,
                "key":    key,
            },
            DurationMs: time.Since(start).Milliseconds(),
        }, nil
    }

    // 未知action
    return &utools.ExecutionResult{
        Success:    false,
        Message:    fmt.Sprintf("不支持的操作: %s", action),
        Error:      "unsupported action",
        DurationMs: time.Since(start).Milliseconds(),
    }, nil
}

// init 自动注册工具
func init() {
    _ = gocontext.Background() // 避免未使用导入别名的lint
    tool := NewRedisContextTool()
    if err := utools.Register(tool); err != nil {
        panic(fmt.Sprintf("Failed to register Redis Context tool: %v", err))
    }
}

