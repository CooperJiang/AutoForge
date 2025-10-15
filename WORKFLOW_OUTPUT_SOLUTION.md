# 工作流输出格式解决方案

## 🎯 问题总结

### 当前问题
1. ✅ **上下文管理器已通用**（支持所有 LLM 模型）
2. ❌ **工作流执行是异步的**，接口立即返回
3. ❌ **没有明确的最终输出**，所有节点输出混在一起
4. ❌ **无法自定义输出格式**

### 你的需求
1. 控制接口返回的内容
2. 返回最后一个节点的输出
3. 或自定义输出格式

---

## 💡 解决方案

### 方案 1: 使用现有的 Output Formatter 工具（推荐）

**原理**：在工作流最后添加一个 Output Formatter 节点，定义输出格式

**示例工作流**：
```
外部触发
  ↓
上下文管理器 (Prepare)
  ↓
OpenAI Chat
  ↓
上下文管理器 (Persist)
  ↓
Output Formatter ← 定义最终输出格式
```

**配置**：
```json
{
  "id": "format_output",
  "type": "tool",
  "toolCode": "output_formatter",
  "config": {
    "type": "json",
    "content": {
      "reply": "{{nodes.ai_chat.content}}",
      "session_id": "{{external.session_id}}",
      "message_count": "{{nodes.save_context.message_count}}"
    }
  }
}
```

**API 返回示例**：
```json
{
  "success": true,
  "data": {
    "execution_id": "exec_xxx",
    "node_logs": [...],
    "final_output": {  // ← 最后一个节点的输出
      "reply": "你好张三！很高兴为你服务...",
      "session_id": "user_001",
      "message_count": 5
    }
  }
}
```

---

### 方案 2: 使用 JSON Transform 工具（灵活性更高）

**原理**：使用 JS 表达式动态构建输出

**配置**：
```json
{
  "id": "transform_output",
  "type": "tool",
  "toolCode": "json_transform",
  "config": {
    "input_data": "{{nodes}}",
    "expression": "{ reply: data.ai_chat.content, timestamp: new Date().toISOString() }"
  }
}
```

---

### 方案 3: 添加同步执行模式（需开发）

**需要修改的文件**：

#### 1. 添加请求参数
```go
// internal/dto/request/workflowRequest.go

type ExecuteWorkflowRequest struct {
    EnvVars map[string]string      `json:"env_vars"`
    Params  map[string]interface{} `json:"params"`
    Sync    bool                   `json:"sync"` // ← 新增：是否同步执行
    Timeout int                    `json:"timeout"` // ← 新增：超时时间（秒）
}
```

#### 2. 修改执行控制器
```go
// internal/controllers/workflow/executionController.go

func ExecuteWorkflow(c *gin.Context) {
    // ... 现有代码 ...

    if req.Sync {
        // 同步执行
        errChan := make(chan error, 1)
        go func() {
            errChan <- engineService.ExecuteWorkflow(execution.GetID(), req.EnvVars, req.Params)
        }()

        // 等待结果（最多 timeout 秒）
        timeout := time.Duration(req.Timeout) * time.Second
        if timeout == 0 {
            timeout = 60 * time.Second
        }

        select {
        case err := <-errChan:
            if err != nil {
                errors.HandleError(c, errors.New(errors.CodeInternal, "执行失败: "+err.Error()))
                return
            }
        case <-time.After(timeout):
            errors.HandleError(c, errors.New(errors.CodeTimeout, "执行超时"))
            return
        }

        // 获取最终结果
        execution, _ = executionService.GetExecutionByID(execution.GetID(), userID)

        // 提取最后一个节点的输出
        var finalOutput map[string]interface{}
        if len(execution.NodeLogs) > 0 {
            lastNode := execution.NodeLogs[len(execution.NodeLogs)-1]
            finalOutput = lastNode.Output
        }

        errors.ResponseSuccess(c, response.ExecuteWorkflowResponse{
            ExecutionID: execution.GetID(),
            Status:      execution.Status,
            Output:      finalOutput, // ← 返回最终输出
            Message:     "执行成功",
        }, "执行成功")
    } else {
        // 异步执行（现有逻辑）
        go func() {
            if err := engineService.ExecuteWorkflow(execution.GetID(), req.EnvVars, req.Params); err != nil {
                log.Error("工作流执行失败: ExecutionID=%s, Error=%v", execution.GetID(), err)
            }
        }()

        errors.ResponseSuccess(c, response.ExecuteWorkflowResponse{
            ExecutionID: execution.GetID(),
            Status:      execution.Status,
            Message:     "工作流已开始执行",
        }, "工作流已开始执行")
    }
}
```

#### 3. 修改响应结构
```go
// internal/dto/response/workflowResponse.go

type ExecuteWorkflowResponse struct {
    ExecutionID string                 `json:"execution_id"`
    Status      string                 `json:"status"`
    Output      map[string]interface{} `json:"output,omitempty"` // ← 新增：最终输出
    Message     string                 `json:"message"`
}
```

---

## 🚀 快速实现（推荐使用方案 1）

### 完整工作流示例

```json
{
  "name": "智能对话（带格式化输出）",
  "nodes": [
    {
      "id": "trigger",
      "type": "external_trigger",
      "config": {
        "params": [
          {"key": "session_id", "type": "string", "required": true},
          {"key": "user_message", "type": "string", "required": true}
        ]
      }
    },
    {
      "id": "prepare",
      "type": "tool",
      "toolCode": "context_manager",
      "config": {
        "mode": "prepare",
        "session_key": "{{external.session_id}}",
        "user_input": "{{external.user_message}}",
        "system_message": "你是一个友好的AI助手"
      }
    },
    {
      "id": "chat",
      "type": "tool",
      "toolCode": "openai_chatgpt",
      "config": {
        "model": "gpt-3.5-turbo",
        "messages_json": "{{nodes.prepare.messages_json}}"
      }
    },
    {
      "id": "persist",
      "type": "tool",
      "toolCode": "context_manager",
      "config": {
        "mode": "persist",
        "session_key": "{{external.session_id}}",
        "assistant_output": "{{nodes.chat.content}}"
      }
    },
    {
      "id": "format_response",
      "type": "tool",
      "toolCode": "output_formatter",
      "config": {
        "type": "json",
        "content": {
          "success": true,
          "data": {
            "reply": "{{nodes.chat.content}}",
            "session_id": "{{external.session_id}}",
            "conversation_length": "{{nodes.persist.message_count}}"
          }
        }
      }
    }
  ],
  "edges": [
    {"source": "trigger", "target": "prepare"},
    {"source": "prepare", "target": "chat"},
    {"source": "chat", "target": "persist"},
    {"source": "persist", "target": "format_response"}
  ]
}
```

### 调用 API

```bash
curl -X POST http://localhost:7777/api/v1/workflows/{id}/execute \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "session_id": "user_001",
      "user_message": "你好"
    }
  }'
```

### 等待并获取结果

```bash
# 1. 立即返回 execution_id
{
  "execution_id": "exec_xxx",
  "status": "pending"
}

# 2. 轮询结果（或使用 WebSocket）
curl http://localhost:7777/api/v1/executions/exec_xxx

# 3. 获取最后一个节点的输出
{
  "node_logs": [...],
  "status": "success",
  "final_node_output": {  // ← 最后一个节点（format_response）的输出
    "success": true,
    "data": {
      "reply": "你好！我能帮你什么？",
      "session_id": "user_001",
      "conversation_length": 3
    }
  }
}
```

---

## 📊 三种方案对比

| 方案 | 优点 | 缺点 | 适用场景 |
|-----|------|------|---------|
| **Output Formatter** | 简单、无需修改代码、配置灵活 | 需要轮询结果 | 大多数场景 |
| **JSON Transform** | 超级灵活、支持复杂逻辑 | 需要写 JS 表达式 | 复杂数据转换 |
| **同步执行模式** | 一次调用得到结果 | 需要修改代码、可能超时 | 快速原型、简单流程 |

---

## ✅ 推荐做法

1. **立即可用**：使用 Output Formatter（方案 1）
2. **复杂场景**：使用 JSON Transform（方案 2）
3. **未来优化**：开发同步执行模式（方案 3）

---

## 🔧 Output Formatter 工具使用详解

### 基本配置

```json
{
  "type": "json",  // 输出类型：json, text, html, markdown, image等
  "content": {     // 输出内容（支持变量）
    "field1": "{{nodes.xxx.yyy}}",
    "field2": "固定值"
  }
}
```

### 变量引用

- `{{external.xxx}}` - 外部参数
- `{{nodes.节点ID.字段名}}` - 节点输出
- `{{env.xxx}}` - 环境变量

### 示例 1: 简单输出

```json
{
  "type": "json",
  "content": {
    "reply": "{{nodes.chat.content}}"
  }
}
```

### 示例 2: 复杂输出

```json
{
  "type": "json",
  "content": {
    "status": "success",
    "data": {
      "user_message": "{{external.user_message}}",
      "ai_reply": "{{nodes.chat.content}}",
      "model": "{{nodes.chat.model}}",
      "tokens": "{{nodes.chat.usage.total_tokens}}",
      "session_info": {
        "id": "{{external.session_id}}",
        "messages_count": "{{nodes.persist.message_count}}"
      }
    },
    "metadata": {
      "execution_time": "{{execution.duration_ms}}ms"
    }
  }
}
```

---

## 💡 最佳实践

1. **标准化输出格式**：统一使用 `{success, data, error}` 结构
2. **添加元数据**：包含 session_id、timestamp 等信息
3. **错误处理**：添加条件节点，区分成功/失败输出
4. **日志记录**：保留原始节点输出用于调试

---

## 🎯 下一步

选择你喜欢的方案：
1. **我给你生成一个带 Output Formatter 的完整示例 JSON**
2. **我帮你开发同步执行模式（方案 3）**
3. **我给你更多 Output Formatter 的示例**

请告诉我你的选择！
