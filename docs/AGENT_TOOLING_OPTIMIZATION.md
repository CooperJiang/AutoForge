# Agent 工具系统优化文档

## 📋 概述

本文档描述了 Agent 工具系统的模块化重构和高优先级优化功能。

## 🏗️ 架构设计

### 模块划分

```
pkg/agent/
├── executor/
│   ├── plan.go              # 原有执行器（保留兼容）
│   ├── plan_refactored.go   # 重构版执行器（推荐使用）
│   ├── step_executor.go     # 步骤执行器
│   └── context.go           # 执行上下文
├── tooling/
│   ├── metadata.go          # 工具元数据定义
│   ├── executor.go          # 工具执行器（带超时和重试）
│   └── validator.go         # 计划验证器
└── registry/
    └── registry.go          # 工具注册表
```

### 设计原则

1. **单一职责**：每个模块只负责一个明确的功能
2. **可测试性**：所有模块都可以独立测试
3. **可扩展性**：易于添加新功能而不影响现有代码
4. **向后兼容**：保留原有 API，新功能通过新接口提供

## ✨ 核心功能

### 1. 工具执行超时控制

**位置**：`pkg/agent/tooling/metadata.go`

**功能**：为每个工具设置执行超时时间，防止长时间卡死。

**配置示例**：

```go
config := &tooling.ExecutionConfig{
    TimeoutSeconds: 300, // 5 分钟超时
}
```

**默认值**：300 秒（5 分钟）

**使用方式**：

```go
// 方式 1：工具实现 ConfigurableTool 接口
type MyTool struct {
    // ...
}

func (t *MyTool) GetExecutionConfig() *tooling.ExecutionConfig {
    return &tooling.ExecutionConfig{
        TimeoutSeconds: 60, // 1 分钟
    }
}

// 方式 2：使用默认配置
// 如果工具未实现接口，会自动使用默认配置
```

### 2. 智能重试机制

**位置**：`pkg/agent/tooling/metadata.go`

**功能**：工具执行失败时自动重试，支持指数退避。

**配置示例**：

```go
config := &tooling.ExecutionConfig{
    Retry: &tooling.RetryConfig{
        MaxRetries:        2,              // 最多重试 2 次
        InitialBackoff:    1000,           // 初始退避 1 秒
        MaxBackoff:        10000,          // 最大退避 10 秒
        BackoffMultiplier: 2.0,            // 退避倍数
        RetryableErrors: []string{         // 可重试的错误关键词
            "timeout",
            "connection",
            "network",
            "rate limit",
            "503",
            "504",
        },
    },
}
```

**重试逻辑**：

1. 第 1 次失败：等待 1 秒后重试
2. 第 2 次失败：等待 2 秒后重试
3. 第 3 次失败：放弃，返回错误

**错误判断**：

- 只有包含 `RetryableErrors` 中关键词的错误才会重试
- 其他错误（如参数错误）直接失败，不重试

### 3. 工具依赖验证

**位置**：`pkg/agent/tooling/validator.go`

**功能**：在执行前验证计划的合理性，检查工具依赖关系。

**配置示例**：

```go
// 在工具的 ExecutionConfig 中声明依赖
config := &tooling.ExecutionConfig{
    Dependencies: &tooling.DependencyConfig{
        // 输入要求
        Requires: []string{"image_url"},
        
        // 输出提供
        Provides: []string{"cdn_url"},
        
        // 建议的前置工具
        SuggestedPredecessors: []string{"openai_image", "file_downloader"},
        
        // 互斥工具
        ConflictsWith: []string{"another_upload_tool"},
    },
}
```

**验证规则**：

1. **依赖检查**：如果工具需要 `image_url`，前面的步骤必须提供它
2. **冲突检查**：如果两个工具互斥，不能同时出现在计划中
3. **建议检查**：如果缺少建议的前置工具，会发出警告（不阻止执行）

**验证结果**：

```go
type ValidationResult struct {
    Valid    bool     // 是否有效
    Errors   []string // 错误列表（阻止执行）
    Warnings []string // 警告列表（不阻止执行）
}
```

### 4. 模块化步骤执行器

**位置**：`pkg/agent/executor/step_executor.go`

**功能**：将步骤执行逻辑从主执行器中分离，提高可维护性。

**职责**：

- 生成工具参数
- 执行工具（调用 tooling.ToolExecutor）
- 处理执行结果
- 发送 SSE 事件

**使用示例**：

```go
stepExecutor := NewStepExecutor(toolRegistry)

result := stepExecutor.ExecuteStep(ctx, &ExecuteStepRequest{
    PlanStep:      planStep,
    StepIndex:     1,
    UserMessage:   "用户的原始请求",
    PreviousSteps: []models.AgentStep{},
    StreamCallback: func(event StreamEvent) {
        // 处理 SSE 事件
    },
})

if result.Error != nil {
    // 处理错误
}
```

## 🔄 迁移指南

### 从旧版本迁移到新版本

**步骤 1：更新导入**

```go
// 旧版本
import "auto-forge/pkg/agent/executor"

// 新版本（保持不变，但使用新的构造函数）
import "auto-forge/pkg/agent/executor"
```

**步骤 2：更新执行器创建**

```go
// 旧版本
executor := executor.NewPlanExecutor(llmClient, toolRegistry, temperature)

// 新版本
executor := executor.NewPlanExecutorV2(llmClient, toolRegistry, temperature)
```

**步骤 3：API 保持不变**

```go
// Execute 方法签名完全相同
result, err := executor.Execute(
    ctx,
    userMessage,
    conversationHistory,
    allowedTools,
    maxSteps,
    streamCallback,
)
```

### 兼容性

- ✅ **完全向后兼容**：旧的 `PlanExecutor` 仍然可用
- ✅ **API 不变**：`Execute` 方法签名完全相同
- ✅ **渐进式迁移**：可以逐步切换到新版本

## 🛠️ 工具开发指南

### 如何让工具支持新功能

**步骤 1：实现 ConfigurableTool 接口**

```go
package mytools

import (
    "auto-forge/pkg/agent/tooling"
    "auto-forge/pkg/utools"
)

type MyTool struct {
    utools.BaseTool
}

// 实现 ConfigurableTool 接口
func (t *MyTool) GetExecutionConfig() *tooling.ExecutionConfig {
    return &tooling.ExecutionConfig{
        // 超时配置
        TimeoutSeconds: 60, // 1 分钟
        
        // 重试配置
        Retry: &tooling.RetryConfig{
            MaxRetries:        2,
            InitialBackoff:    1000,
            MaxBackoff:        5000,
            BackoffMultiplier: 2.0,
            RetryableErrors: []string{
                "timeout",
                "connection failed",
            },
        },
        
        // 依赖配置
        Dependencies: &tooling.DependencyConfig{
            Requires: []string{"image_url"},
            Provides: []string{"processed_image"},
            SuggestedPredecessors: []string{"image_generator"},
        },
        
        // 缓存配置（可选）
        Cache: &tooling.CacheConfig{
            Enabled: false,
            TTL:     5 * time.Minute,
        },
    }
}
```

**步骤 2：注册工具**

```go
// 工具会自动继承配置
registry.Register(tool)
```

### 配置优先级

1. **工具级别配置**：`tool.GetExecutionConfig()`
2. **默认配置**：`tooling.DefaultExecutionConfig()`

## 📊 性能优化

### 超时控制的影响

- **优点**：防止单个工具卡死整个流程
- **开销**：几乎无性能开销（使用 context.WithTimeout）

### 重试机制的影响

- **优点**：提高成功率，减少临时错误
- **开销**：失败时会增加总执行时间（但比手动重试更高效）

### 依赖验证的影响

- **优点**：提前发现问题，避免无效执行
- **开销**：计划生成后增加 < 10ms 的验证时间

## 🧪 测试

### 单元测试

```go
// 测试工具执行器
func TestToolExecutor(t *testing.T) {
    config := &tooling.ExecutionConfig{
        TimeoutSeconds: 1,
        Retry: &tooling.RetryConfig{
            MaxRetries: 2,
        },
    }
    
    executor := tooling.NewToolExecutor(config)
    
    // 测试超时
    result := executor.Execute(ctx, slowTool, args)
    assert.Error(t, result.Error)
    assert.Contains(t, result.Error.Error(), "timeout")
    
    // 测试重试
    result = executor.Execute(ctx, flakeyTool, args)
    assert.NoError(t, result.Error)
    assert.Equal(t, 2, result.Attempts) // 第一次失败，第二次成功
}
```

### 集成测试

```go
// 测试完整流程
func TestPlanExecutorV2(t *testing.T) {
    executor := executor.NewPlanExecutorV2(llmClient, toolRegistry, 0.7)
    
    result, err := executor.Execute(
        ctx,
        "生成图片并上传",
        "",
        []string{"openai_image", "pixelpunk_upload"},
        10,
        nil,
    )
    
    assert.NoError(t, err)
    assert.True(t, result.Success)
    assert.Len(t, result.Trace.Steps, 2)
}
```

## 📈 监控和调试

### 日志输出

新版本增加了详细的日志输出：

```
[INFO] 执行工具: pixelpunk_upload, 参数: {...}
[INFO] 工具执行成功（尝试 2 次），结果长度: 1234
[WARN] 计划验证失败: [步骤 2 需要 image_url，但前面未提供]
[ERROR] 工具执行失败（尝试 3 次): timeout
```

### SSE 事件

新增 `tool_progress` 事件：

```json
{
  "type": "tool_progress",
  "data": {
    "step": 1,
    "tool": "openai_image",
    "attempt": 2,
    "message": "重试中（第 2 次）"
  }
}
```

## 🚀 未来扩展

### 已规划的功能

1. **工具输出缓存**：避免重复执行相同的工具调用
2. **并行执行**：多个独立步骤可以并行执行
3. **条件执行**：根据结果选择不同的执行路径
4. **工具组合模板**：预定义常见的工具组合

### 扩展点

- `tooling.ExecutionConfig`：可以添加新的配置项
- `tooling.PlanValidator`：可以添加新的验证规则
- `executor.StepExecutor`：可以添加新的执行策略

## 📝 最佳实践

### 1. 合理设置超时时间

```go
// ❌ 不好：超时时间过短
TimeoutSeconds: 5  // AI 生成可能需要更长时间

// ✅ 好：根据工具特性设置
TimeoutSeconds: 60  // AI 生成
TimeoutSeconds: 10  // 简单的 HTTP 请求
TimeoutSeconds: 300 // 大文件上传
```

### 2. 精确定义可重试错误

```go
// ❌ 不好：重试所有错误
RetryableErrors: []string{"error"}

// ✅ 好：只重试临时错误
RetryableErrors: []string{
    "timeout",
    "connection",
    "503",
    "504",
}
```

### 3. 明确声明依赖关系

```go
// ❌ 不好：不声明依赖
Dependencies: nil

// ✅ 好：明确声明
Dependencies: &tooling.DependencyConfig{
    Requires: []string{"image_url"},
    Provides: []string{"cdn_url"},
}
```

## 🔗 相关文档

- [Agent 架构文档](./AGENT_ARCHITECTURE.md)
- [工具开发指南](./TOOL_DEVELOPMENT_GUIDE.md)
- [Agent 聊天开发](./AGENT_CHAT_DEVELOPMENT.md)

## 📞 支持

如有问题，请查看：

1. 日志输出（包含详细的执行信息）
2. SSE 事件（包含实时进度）
3. ValidationResult（包含验证错误和警告）

