# 上下文管理器工具使用指南

## 🎯 概述

上下文管理器（Context Manager）是一个通用的对话记忆管理工具，支持任意 LLM 模型（OpenAI、Gemini、Claude 等）。

**核心功能**：
- ✅ **Prepare 模式**：读取历史对话 + 拼接当前用户消息 → 输出完整消息数组
- ✅ **Persist 模式**：保存 AI 回复到历史记录
- ✅ **会话隔离**：通过 session_key 区分不同用户的对话
- ✅ **自动裁剪**：保留最近 N 条消息，防止上下文过长
- ✅ **通用设计**：不绑定任何特定 LLM，可与任意对话工具配合使用

---

## 📦 工作流程

```
外部 API 触发
    ↓
    传入参数：session_id, user_message
    ↓
┌───────────────────────────────────┐
│  上下文管理器 (Prepare)            │
│  • 读取 Redis 历史记录             │
│  • 拼接 system 消息（可选）         │
│  • 拼接当前 user 消息              │
│  • 裁剪到窗口大小                  │
│  • 输出 messages_json             │
└───────────────────────────────────┘
    ↓
┌───────────────────────────────────┐
│  OpenAI Chat / Gemini / 其他 LLM   │
│  • 接收 messages_json             │
│  • 调用 AI 模型                    │
│  • 输出 content (AI 回复)          │
└───────────────────────────────────┘
    ↓
┌───────────────────────────────────┐
│  上下文管理器 (Persist)            │
│  • 读取历史记录                    │
│  • 追加 user 消息（可选）           │
│  • 追加 assistant 消息             │
│  • 保存到 Redis (TTL=7天)          │
└───────────────────────────────────┘
```

---

## 🔧 配置说明

### Prepare 模式配置

| 参数 | 类型 | 必填 | 说明 | 示例 |
|-----|------|------|------|------|
| `mode` | string | ✅ | 工作模式 | `"prepare"` |
| `session_key` | string | ✅ | 会话标识 | `"{{external.session_id}}"` |
| `user_input` | string | ✅ | 当前用户消息 | `"{{external.user_message}}"` |
| `system_message` | string | ❌ | 系统提示词 | `"你是一个友好的AI助手..."` |
| `window_size` | number | ❌ | 保留消息数 | `10`（默认） |
| `ttl_seconds` | number | ❌ | 过期时间（秒）| `604800`（7天，默认） |
| `clear_history` | boolean | ❌ | 清空历史 | `false`（默认） |

**输出字段**：
- `messages_json`：标准消息数组的 JSON 字符串（传给 LLM）
- `messages`：消息数组对象
- `session_key`：会话标识
- `message_count`：消息总数
- `preview`：对话预览文本

---

### Persist 模式配置

| 参数 | 类型 | 必填 | 说明 | 示例 |
|-----|------|------|------|------|
| `mode` | string | ✅ | 工作模式 | `"persist"` |
| `session_key` | string | ✅ | 会话标识 | `"{{external.session_id}}"` |
| `assistant_output` | string | ✅ | AI 回复内容 | `"{{nodes.openai_chat.content}}"` |
| `user_input` | string | ❌ | 用户消息（通常 Prepare 已保存）| - |
| `window_size` | number | ❌ | 保留消息数 | `10`（默认） |
| `ttl_seconds` | number | ❌ | 过期时间（秒）| `604800`（默认） |

**输出字段**：
- `updated`：是否已更新（`true`）
- `session_key`：会话标识
- `message_count`：保存后的消息总数
- `ttl_seconds`：过期时间

---

## 📝 完整示例

### 1. 导入工作流

将 `web/chat_with_memory_v3.json` 导入到系统中。

### 2. 工作流结构

```json
{
  "nodes": [
    {
      "id": "trigger_start",
      "type": "external_trigger",
      "config": {
        "params": [
          {"key": "session_id", "type": "string", "required": true},
          {"key": "user_message", "type": "string", "required": true}
        ]
      }
    },
    {
      "id": "context_prepare",
      "type": "tool",
      "toolCode": "context_manager",
      "config": {
        "mode": "prepare",
        "session_key": "{{external.session_id}}",
        "user_input": "{{external.user_message}}",
        "system_message": "你是一个友好的AI助手...",
        "window_size": 10
      }
    },
    {
      "id": "openai_chat",
      "type": "tool",
      "toolCode": "openai_chatgpt",
      "config": {
        "model": "gpt-3.5-turbo",
        "messages_json": "{{nodes.context_prepare.messages_json}}",
        "temperature": 0.7
      }
    },
    {
      "id": "context_persist",
      "type": "tool",
      "toolCode": "context_manager",
      "config": {
        "mode": "persist",
        "session_key": "{{external.session_id}}",
        "assistant_output": "{{nodes.openai_chat.content}}"
      }
    }
  ],
  "edges": [
    {"source": "trigger_start", "target": "context_prepare"},
    {"source": "context_prepare", "target": "openai_chat"},
    {"source": "openai_chat", "target": "context_persist"}
  ]
}
```

### 3. 调用 API

**首次对话**：
```bash
curl -X POST http://localhost:7777/api/v1/workflows/{workflow_id}/execute \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "session_id": "user_alice_001",
      "user_message": "你好，我叫 Alice，我是一名软件工程师"
    }
  }'
```

**响应**：
```json
{
  "success": true,
  "output": {
    "context_persist": {
      "updated": true,
      "message_count": 3
    },
    "openai_chat": {
      "content": "你好 Alice！很高兴认识你。作为一名软件工程师，你一定对技术充满热情..."
    }
  }
}
```

**第二轮对话（测试记忆）**：
```bash
curl -X POST http://localhost:7777/api/v1/workflows/{workflow_id}/execute \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "session_id": "user_alice_001",
      "user_message": "我叫什么？我的职业是什么？"
    }
  }'
```

**响应**：
```json
{
  "success": true,
  "output": {
    "openai_chat": {
      "content": "你叫 Alice，你的职业是软件工程师。"
    }
  }
}
```

✅ **成功标志**：AI 准确回答了你的姓名和职业！

---

## 🔍 验证记忆功能

### 检查 Redis 存储

```bash
# 查看所有会话
redis-cli KEYS "chat:context:*"

# 查看特定会话的历史
redis-cli GET "chat:context:user_alice_001"
```

**预期输出**：
```json
[
  {
    "role": "system",
    "content": "你是一个友好的AI助手..."
  },
  {
    "role": "user",
    "content": "你好，我叫 Alice，我是一名软件工程师"
  },
  {
    "role": "assistant",
    "content": "你好 Alice！很高兴认识你..."
  },
  {
    "role": "user",
    "content": "我叫什么？我的职业是什么？"
  },
  {
    "role": "assistant",
    "content": "你叫 Alice，你的职业是软件工程师。"
  }
]
```

---

## 🧪 会话隔离测试

```bash
# 用户 Bob 的对话（新会话）
curl -X POST http://localhost:7777/api/v1/workflows/{workflow_id}/execute \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "session_id": "user_bob_002",
      "user_message": "我叫什么？"
    }
  }'
```

**预期响应**：
```json
{
  "output": {
    "openai_chat": {
      "content": "抱歉，我们刚开始对话，您还没有告诉我您的名字..."
    }
  }
}
```

✅ **成功标志**：不同 session_id 的对话完全隔离！

---

## 🛠️ 高级用法

### 1. 清空对话历史

在 Prepare 节点配置中设置：
```json
{
  "clear_history": true
}
```

这会在执行前清空当前会话的历史记录，重新开始对话。

---

### 2. 调整窗口大小

```json
{
  "window_size": 20
}
```

保留最近 20 条消息（包括 system、user、assistant）。

**注意**：
- system 消息始终保留（不计入窗口）
- 裁剪时保留最近的 N 条消息
- 建议根据模型 token 限制调整窗口大小

---

### 3. 自定义过期时间

```json
{
  "ttl_seconds": 86400
}
```

设置对话历史的存储时长（24 小时）。

---

### 4. 与其他 LLM 配合使用

**Gemini 示例**：
```json
{
  "id": "gemini_chat",
  "type": "tool",
  "toolCode": "gemini_chatgpt",
  "config": {
    "model": "gemini-pro",
    "messages_json": "{{nodes.context_prepare.messages_json}}"
  }
}
```

只需修改工具代码和模型参数，其他配置完全相同！

---

## 📊 变量引用说明

### 外部参数变量

**格式**：`{{external.参数名}}`

示例：
- `{{external.session_id}}`
- `{{external.user_message}}`
- `{{external.user_id}}`

### 上游节点输出

**格式**：`{{nodes.节点ID.字段名}}`

示例：
- `{{nodes.context_prepare.messages_json}}` - Prepare 输出的消息数组
- `{{nodes.openai_chat.content}}` - OpenAI 回复内容
- `{{nodes.context_prepare.message_count}}` - 消息数量

---

## ❓ 常见问题

### Q1: 变量未替换怎么办？

**症状**：Redis 中 key 为 `chat:context:{{external.session_id}}`

**原因**：引擎未正确替换变量

**解决**：
1. 确认使用 `{{external.xxx}}` 格式
2. 检查 API 调用时是否传入了参数
3. 查看日志确认变量替换

---

### Q2: AI 没有记住历史怎么办？

**排查步骤**：

1. **检查 Redis 是否有数据**：
```bash
redis-cli KEYS "chat:context:*"
```

2. **检查 Prepare 节点输出**：
确认 `messages_json` 包含历史消息

3. **检查 LLM 节点配置**：
确认使用了 `{{nodes.context_prepare.messages_json}}`

4. **检查日志**：
```bash
tail -f logs/app.log | grep "上下文管理器"
```

---

### Q3: 如何支持多轮复杂对话？

**建议配置**：
- `window_size`: 20-30（根据模型 token 限制）
- `ttl_seconds`: 604800（7天）
- `system_message`: 详细说明 AI 的角色和记忆能力

**示例 system 消息**：
```
你是一个专业的客服助手。你需要记住用户在对话中提供的所有信息，包括：
- 姓名、联系方式等个人信息
- 之前提出的问题和需求
- 偏好和特殊要求

在后续对话中，你要主动引用这些信息，提供个性化的服务。
```

---

### Q4: 如何实现多租户隔离？

**方案 1：session_key 包含租户 ID**
```json
{
  "session_key": "{{external.tenant_id}}_{{external.user_id}}"
}
```

**方案 2：使用环境变量**
```json
{
  "session_key": "tenant_{{env.TENANT_ID}}_user_{{external.user_id}}"
}
```

---

## 🎯 最佳实践

1. **session_key 命名规范**：使用有意义的标识，如 `user_{user_id}` 或 `session_{timestamp}`

2. **窗口大小设置**：
   - GPT-3.5: 10-15 条
   - GPT-4: 20-30 条
   - Gemini Pro: 15-20 条

3. **system 消息编写**：
   - 明确告知 AI 需要记住用户信息
   - 说明如何使用历史信息
   - 定义 AI 的角色和行为

4. **错误处理**：
   - 在 Persist 节点前添加条件节点，检查 LLM 执行是否成功
   - 如果失败，跳过 Persist，避免保存错误信息

5. **性能优化**：
   - 合理设置 TTL，避免 Redis 存储过多数据
   - 根据业务需求调整 window_size
   - 定期清理过期会话

---

## 🚀 总结

**上下文管理器工具的优势**：
- ✅ **通用性**：不绑定任何 LLM，可与任意对话工具配合
- ✅ **可见性**：在画布上可见，便于调试和理解
- ✅ **灵活性**：可在 Prepare 和 Persist 之间插入任意节点
- ✅ **可控性**：完全控制何时读取、何时保存
- ✅ **可扩展**：支持自定义 system 消息、窗口大小、TTL 等

**与旧的中间件方案对比**：
- 旧方案：自动化，但不可见，不够通用
- 新方案：可视化，通用，灵活，易于理解和调试

**适用场景**：
- 客服对话系统
- 个人 AI 助手
- 教育辅导系统
- 知识问答系统
- 任何需要多轮对话记忆的场景
