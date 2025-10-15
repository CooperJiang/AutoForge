# 对话上下文功能测试指南

## 功能说明

本次实现了 **对话上下文中间件**，可以让 AI 对话工具（如 OpenAI Chat）自动记住历史对话内容，实现真正的多轮对话。

## 实现方式

- **架构**：引擎层中间件（自动管理，无需用户手动配置多个节点）
- **存储**：Redis 优先，降级到内存（复用现有 cache 包）
- **通用性**：所有标记 `ContextAware: true` 的工具都自动支持
- **配置**：可选启用，默认关闭，不影响现有功能

## 快速测试

### 步骤 1: 导入测试工作流

1. 打开工作流管理页面 `/workflows`
2. 点击"导入工作流"按钮
3. 上传或粘贴 `test-workflow-chat-with-context.json` 文件内容
4. 保存工作流

### 步骤 2: 测试多轮对话

使用以下 cURL 命令进行测试（假设工作流 ID 为 `your-workflow-id`）：

#### 第一轮对话：告诉 AI 你的名字

```bash
curl -X POST http://localhost:7777/api/v1/workflows/your-workflow-id/execute \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "session_id": "test_user_123",
      "user_message": "你好，我叫张三，今年 25 岁"
    }
  }'
```

**预期返回**：
```json
{
  "success": true,
  "message": "工作流执行成功",
  "output": {
    "openai_chat_1": {
      "content": "你好张三！很高兴认识你。25 岁正是充满活力的年纪...",
      "model": "gpt-3.5-turbo",
      ...
    }
  }
}
```

#### 第二轮对话：测试 AI 是否记住你的信息

```bash
curl -X POST http://localhost:7777/api/v1/workflows/your-workflow-id/execute \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "session_id": "test_user_123",
      "user_message": "我叫什么名字？今年多大？"
    }
  }'
```

**预期返回**：
```json
{
  "success": true,
  "output": {
    "openai_chat_1": {
      "content": "你叫张三，今年 25 岁。",
      ...
    }
  }
}
```

✅ **成功标志**：AI 能够准确回答你在第一轮对话中提供的信息！

#### 第三轮对话：继续测试上下文

```bash
curl -X POST http://localhost:7777/api/v1/workflows/your-workflow-id/execute \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "session_id": "test_user_123",
      "user_message": "我明年多大？"
    }
  }'
```

**预期返回**：
```json
{
  "success": true,
  "output": {
    "openai_chat_1": {
      "content": "你明年就 26 岁了！",
      ...
    }
  }
}
```

### 步骤 3: 测试不同会话隔离

使用不同的 `session_id` 测试会话隔离：

```bash
curl -X POST http://localhost:7777/api/v1/workflows/your-workflow-id/execute \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "session_id": "test_user_456",
      "user_message": "我叫什么名字？"
    }
  }'
```

**预期返回**：
```json
{
  "success": true,
  "output": {
    "openai_chat_1": {
      "content": "抱歉，我们刚开始对话，您还没有告诉我您的名字...",
      ...
    }
  }
}
```

✅ **成功标志**：不同 session_id 的对话完全隔离！

## 查看 Redis 中的对话历史

如果你配置了 Redis，可以查看存储的对话数据：

```bash
# 查看所有对话 key
redis-cli KEYS "chat:context:*"

# 查看特定会话的对话内容
redis-cli GET "chat:context:test_user_123"
```

**返回示例**：
```json
[
  {
    "role": "system",
    "content": "你是一个友好的 AI 助手，擅长记住用户告诉你的信息..."
  },
  {
    "role": "user",
    "content": "你好，我叫张三，今年 25 岁"
  },
  {
    "role": "assistant",
    "content": "你好张三！很高兴认识你..."
  },
  {
    "role": "user",
    "content": "我叫什么名字？今年多大？"
  },
  {
    "role": "assistant",
    "content": "你叫张三，今年 25 岁。"
  }
]
```

## 配置说明

### 对话上下文配置项

在工作流编辑页面，选择 OpenAI Chat 节点，会看到以下配置：

| 配置项 | 说明 | 默认值 |
|-------|------|--------|
| **启用对话记忆** | 开启后自动管理上下文 | 关闭 |
| **会话ID** | 区分不同用户，支持变量 | `{{params.session_id}}` |
| **窗口条数** | 保留最近 N 条消息 | 10 |
| **过期时间** | 对话历史 TTL（秒） | 604800（7天） |

### 变量支持

`session_key` 支持以下变量格式：

- `{{params.session_id}}` - 从外部参数获取
- `{{params.user_id}}` - 从外部参数获取
- `{{env.user_id}}` - 从环境变量获取
- `{{external.xxx}}` - 从外部参数获取

如果变量未替换成功，会自动降级到 `global`（全局共享会话）。

## 特性说明

### ✅ 已实现

1. **自动上下文管理**：无需手动拼接 messages，引擎自动处理
2. **会话隔离**：不同 session_id 完全隔离
3. **窗口裁剪**：自动保留最近 N 条消息，防止超出 token 限制
4. **TTL 过期**：支持自定义过期时间，默认 7 天
5. **向后兼容**：不启用时不影响现有功能
6. **通用架构**：未来添加 Gemini、Claude 等工具只需标记 `ContextAware: true`

### 📋 下一步扩展（Phase 2）

1. **Token 智能裁剪**：基于 token 数量动态裁剪（需集成 tiktoken）
2. **摘要策略**：超长对话自动生成摘要
3. **多模态支持**：图像、文件等多模态消息
4. **可视化调试**：独立节点模式（Prepare + Persist）

## 故障排查

### 问题 1：AI 没有记住历史

**可能原因**：
- session_id 变量未正确替换
- Redis 连接失败（降级到内存，重启后丢失）
- 对话已过期（超过 TTL）

**解决方法**：
```bash
# 检查日志
tail -f logs/app.log | grep "对话上下文"

# 查看 Redis 是否有数据
redis-cli KEYS "chat:context:*"
```

### 问题 2：不同请求之间对话丢失

**可能原因**：每次请求使用了不同的 session_id

**解决方法**：确保同一用户的所有请求使用相同的 `session_id`

### 问题 3：工具执行失败

**可能原因**：OpenAI API Key 未配置

**解决方法**：
```yaml
# config.yaml
openai:
  api_key: "your-api-key-here"
  api_base: "https://api.openai.com/v1"  # 可选
```

## 代码修改清单

本次修改涉及以下文件：

### 后端
1. `pkg/utools/types.go` - 添加 `ContextAware` 字段
2. `pkg/utools/openai/openai_tool.go` - 标记支持上下文，添加 context_config schema
3. `internal/services/workflow/engineService.go` - 应用中间件
4. `internal/services/workflow/context_middleware.go` - 中间件实现（新文件）

### 前端
1. `web/src/components/tools/OpenAIConfig.vue` - 更新 UI 配置

### 测试
1. `test-workflow-chat-with-context.json` - 测试工作流
2. `TEST_CHAT_CONTEXT.md` - 本文档

## 联系方式

如有问题，请查看日志或联系开发团队。
