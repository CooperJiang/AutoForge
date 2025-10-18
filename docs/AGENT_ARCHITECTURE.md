# Agent 架构设计文档 v2.0

## 🎯 设计目标

1. **模块化**：清晰的职责分离，易于扩展和维护
2. **强大**：支持复杂的多步骤推理和工具调用
3. **灵活**：支持多种 LLM 和执行模式
4. **可观测**：完整的执行轨迹和调试信息
5. **高性能**：流式响应、并发执行、智能缓存

## 📐 核心架构

```
┌─────────────────────────────────────────────────────────────┐
│                        Agent Service                         │
│  (对话管理、消息存储、执行协调)                                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      Executor (执行器)                        │
│  ┌──────────────┐              ┌──────────────┐             │
│  │  ReAct Mode  │              │  Plan Mode   │             │
│  │  (边思考边做) │              │  (先规划后做) │             │
│  └──────────────┘              └──────────────┘             │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      LLM Client (大模型)                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   OpenAI     │  │    Gemini    │  │   Custom     │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                   Tool Registry (工具注册表)                  │
│  - 工具发现和注册                                              │
│  - 工具描述生成                                                │
│  - 工具参数验证                                                │
│  - 工具执行代理                                                │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                     UTools (实际工具)                         │
│  HTTP | Email | Formatter | Health | Context | ...          │
└─────────────────────────────────────────────────────────────┘
```

## 🧩 核心组件

### 1. LLM Client (大模型客户端)

**职责**：
- 统一的 LLM 调用接口
- 支持多种模型提供商（OpenAI、Gemini、自定义）
- 流式响应支持
- 错误处理和重试
- Token 计数

**接口设计**：
```go
type LLMClient interface {
    // 同步调用
    Call(ctx context.Context, messages []Message, options CallOptions) (*Response, error)
    
    // 流式调用
    Stream(ctx context.Context, messages []Message, options CallOptions) (<-chan StreamChunk, error)
    
    // 获取模型信息
    GetModelInfo() ModelInfo
}

type Message struct {
    Role    string      // system, user, assistant
    Content string      // 文本内容
    Name    string      // 可选：消息发送者名称
    ToolCalls []ToolCall // 可选：工具调用
}

type CallOptions struct {
    Temperature      float64
    MaxTokens        int
    TopP             float64
    Stop             []string
    Tools            []ToolDefinition // 可用工具列表
    ToolChoice       string           // auto, none, required
    ResponseFormat   string           // text, json
}

type Response struct {
    Content      string
    ToolCalls    []ToolCall
    FinishReason string // stop, tool_calls, length, content_filter
    Usage        TokenUsage
}

type ToolCall struct {
    ID       string
    Type     string // function
    Function FunctionCall
}

type FunctionCall struct {
    Name      string
    Arguments string // JSON string
}
```

**实现**：
- `OpenAIClient`: 支持 GPT-4, GPT-3.5 等
- `GeminiClient`: 支持 Gemini Pro 等
- `CustomClient`: 支持自定义 API（兼容 OpenAI 格式）

---

### 2. Tool Registry (工具注册表)

**职责**：
- 自动发现和注册所有 UTools
- 生成工具的 JSON Schema 描述
- 验证工具调用参数
- 执行工具并返回结果
- 工具权限控制

**接口设计**：
```go
type ToolRegistry interface {
    // 注册工具
    Register(tool Tool) error
    
    // 获取所有工具
    GetAll() []Tool
    
    // 根据名称获取工具
    Get(name string) (Tool, error)
    
    // 获取工具的 LLM 描述
    GetToolDefinitions(allowedTools []string) []ToolDefinition
    
    // 执行工具
    Execute(ctx context.Context, name string, args map[string]interface{}) (*ToolResult, error)
    
    // 验证工具参数
    Validate(name string, args map[string]interface{}) error
}

type Tool struct {
    Name        string
    Description string
    Parameters  ParameterSchema
    Handler     ToolHandler
}

type ToolDefinition struct {
    Type     string              `json:"type"`     // "function"
    Function FunctionDefinition  `json:"function"`
}

type FunctionDefinition struct {
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Parameters  map[string]interface{} `json:"parameters"` // JSON Schema
}

type ToolResult struct {
    Success bool
    Output  interface{}
    Error   string
    Metadata map[string]interface{}
}

type ToolHandler func(ctx context.Context, args map[string]interface{}) (*ToolResult, error)
```

**功能**：
1. **自动发现**：扫描 `pkg/utools` 包，自动注册所有工具
2. **Schema 生成**：从工具的配置 Schema 自动生成 LLM 可理解的描述
3. **参数验证**：使用 JSON Schema 验证工具参数
4. **执行代理**：统一的工具执行接口，处理错误和超时
5. **权限控制**：支持工具白名单/黑名单

---

### 3. Executor (执行器)

#### 3.1 ReAct Executor (ReAct 执行器)

**职责**：
- 实现 ReAct (Reasoning + Acting) 循环
- 边思考边执行
- 动态决策下一步动作

**执行流程**：
```
1. 初始化：构建系统提示词 + 用户问题
2. ReAct 循环：
   a. Thought: LLM 思考下一步动作
   b. Action: 决定调用哪个工具及参数
   c. Observation: 执行工具，获取结果
   d. 重复 2a-2c，直到达到终止条件
3. Final Answer: LLM 基于所有观察生成最终答案
```

**提示词模板**：
```
You are a helpful AI agent with access to various tools.

Available Tools:
{tool_definitions}

Instructions:
1. Analyze the user's request carefully
2. Think step by step about what tools you need to use
3. Call tools using the function calling format
4. Based on tool results, decide your next action
5. When you have enough information, provide a final answer

Conversation History:
{conversation_history}

User Request: {user_message}

Let's think step by step.
```

#### 3.2 Plan Executor (计划执行器)

**职责**：
- 先规划完整的执行计划
- 按计划顺序执行
- 支持计划调整

**执行流程**：
```
1. 规划阶段：
   a. LLM 分析任务
   b. 生成完整的执行计划（步骤列表）
   c. 返回计划给用户确认
2. 执行阶段：
   a. 按顺序执行每个步骤
   b. 每步执行后更新状态
   c. 如果某步失败，可选择跳过或终止
3. 总结阶段：
   a. 收集所有步骤结果
   b. LLM 生成最终答案
```

**规划提示词模板**：
```
You are a planning AI. Given a user request, create a detailed execution plan.

Available Tools:
{tool_definitions}

User Request: {user_message}

Create a step-by-step plan to accomplish this task. For each step:
1. Describe what needs to be done
2. Specify which tool to use
3. List the required parameters

Return the plan in JSON format:
{
  "steps": [
    {
      "step": 1,
      "description": "...",
      "tool": "tool_name",
      "parameters": {...}
    }
  ]
}
```

---

### 4. Prompt Templates (提示词模板)

**职责**：
- 管理各种场景的提示词模板
- 支持变量替换
- 多语言支持

**模板类型**：
1. **System Prompt**: 系统角色定义
2. **ReAct Prompt**: ReAct 循环提示词
3. **Plan Prompt**: 规划提示词
4. **Summary Prompt**: 总结提示词
5. **Error Prompt**: 错误处理提示词

**模板引擎**：
```go
type PromptTemplate struct {
    Name     string
    Template string
    Variables []string
}

func (pt *PromptTemplate) Render(vars map[string]string) string {
    result := pt.Template
    for key, value := range vars {
        result = strings.ReplaceAll(result, "{"+key+"}", value)
    }
    return result
}
```

---

## 🔄 执行流程

### ReAct 模式执行流程

```
用户输入 → Agent Service
              ↓
         创建消息记录
              ↓
         ReAct Executor
              ↓
    ┌─────────────────┐
    │  ReAct 循环开始  │
    └─────────────────┘
              ↓
    ┌─────────────────┐
    │  1. LLM 思考     │ ← 系统提示词 + 工具定义 + 历史
    │  (Thought)      │
    └─────────────────┘
              ↓
    ┌─────────────────┐
    │  2. 决定动作     │ ← LLM 返回工具调用
    │  (Action)       │
    └─────────────────┘
              ↓
    ┌─────────────────┐
    │  3. 执行工具     │ ← Tool Registry
    │  (Tool Call)    │
    └─────────────────┘
              ↓
    ┌─────────────────┐
    │  4. 获取结果     │ ← 工具返回结果
    │  (Observation)  │
    └─────────────────┘
              ↓
         是否终止？
         /        \
       否          是
        ↓          ↓
    回到步骤1   生成最终答案
                    ↓
              更新消息记录
                    ↓
              返回给用户
```

### Plan 模式执行流程

```
用户输入 → Agent Service
              ↓
         创建消息记录
              ↓
         Plan Executor
              ↓
    ┌─────────────────┐
    │  1. 生成计划     │ ← LLM 分析任务
    └─────────────────┘
              ↓
    ┌─────────────────┐
    │  2. 返回计划     │ → 前端显示计划
    └─────────────────┘
              ↓
    ┌─────────────────┐
    │  3. 执行循环     │
    │  For each step:  │
    │  - 执行工具      │
    │  - 更新状态      │
    │  - 流式返回      │
    └─────────────────┘
              ↓
    ┌─────────────────┐
    │  4. 生成总结     │ ← LLM 基于所有结果
    └─────────────────┘
              ↓
         更新消息记录
              ↓
         返回给用户
```

---

## 🎨 特性亮点

### 1. 智能工具选择
- LLM 自动分析任务需求
- 根据工具描述选择最合适的工具
- 支持工具链式调用

### 2. 上下文管理
- 对话历史压缩
- 关键信息提取
- 长上下文窗口支持

### 3. 错误处理
- 工具调用失败自动重试
- 参数错误自动修正
- 优雅降级

### 4. 流式响应
- 实时返回思考过程
- 工具执行进度
- 最终答案流式输出

### 5. 可观测性
- 完整的执行轨迹
- Token 使用统计
- 性能指标

---

## 📊 数据结构

### Trace (执行轨迹)

```go
type AgentTrace struct {
    Steps        []AgentStep            // 执行步骤
    FinalAnswer  string                 // 最终答案
    FinishReason string                 // 终止原因
    UsedTools    map[string]ToolStats   // 工具使用统计
    TokenUsage   TokenUsage             // Token 使用
    TotalMs      int64                  // 总耗时
}

type AgentStep struct {
    Step        int                    // 步骤序号
    Thought     string                 // 思考过程（可选）
    Action      *AgentAction           // 动作
    Observation string                 // 观察结果
    ToolOutput  interface{}            // 工具原始输出
    ElapsedMs   int64                  // 耗时
    Timestamp   string                 // 时间戳
    Error       string                 // 错误信息（如果有）
}

type AgentAction struct {
    Type string                 // "action" | "final"
    Tool string                 // 工具名称
    Args map[string]interface{} // 工具参数
}

type ToolStats struct {
    Count   int   // 调用次数
    TotalMs int64 // 总耗时
}
```

---

## 🚀 实现计划

### Phase 1: 基础设施 (Day 1)
- [x] LLM Client 接口定义
- [ ] OpenAI Client 实现
- [ ] Gemini Client 实现
- [ ] Tool Registry 实现
- [ ] Prompt Template 系统

### Phase 2: 执行引擎 (Day 2)
- [ ] ReAct Executor 实现
- [ ] Plan Executor 实现
- [ ] 错误处理机制
- [ ] 流式响应支持

### Phase 3: 集成和优化 (Day 3)
- [ ] 集成到 Agent Service
- [ ] 性能优化
- [ ] 测试用例
- [ ] 文档完善

---

## 🎯 性能目标

- **响应速度**: 首次响应 < 2s
- **工具调用**: 单次调用 < 5s
- **并发支持**: 100+ 并发对话
- **成功率**: > 95%
- **Token 效率**: 优化提示词，减少 30% Token 消耗

---

**版本**: 2.0  
**创建时间**: 2025-10-18  
**作者**: AI Assistant  
**状态**: 设计中 → 实现中

