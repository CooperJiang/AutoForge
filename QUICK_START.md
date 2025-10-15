# 对话上下文功能 - 快速开始

## 🎯 一分钟快速测试

### 1. 导入测试工作流

在工作流页面 `/workflows` 点击"导入工作流"，粘贴以下 JSON：

<details>
<summary>点击展开 JSON 配置</summary>

```json
{
  "name": "智能对话助手（带上下文记忆）",
  "description": "支持多轮对话的 AI 助手，会记住之前的对话内容",
  "enabled": true,
  "nodes": [
    {
      "id": "trigger_1",
      "type": "external_trigger",
      "position": { "x": 100, "y": 200 },
      "data": { "label": "外部 API 触发" },
      "config": {}
    },
    {
      "id": "openai_chat_1",
      "type": "tool",
      "toolCode": "openai_chatgpt",
      "position": { "x": 400, "y": 200 },
      "data": { "label": "OpenAI Chat（带记忆）" },
      "config": {
        "model": "gpt-3.5-turbo",
        "prompt": "{{params.user_message}}",
        "system_message": "你是一个友好的 AI 助手，擅长记住用户告诉你的信息。",
        "temperature": 0.7,
        "max_tokens": 500,
        "timeout": 60,
        "context_config": {
          "enabled": true,
          "session_key": "{{params.session_id}}",
          "window_size": 10,
          "ttl_seconds": 604800
        }
      }
    }
  ],
  "edges": [
    {
      "id": "edge_1",
      "source": "trigger_1",
      "target": "openai_chat_1",
      "type": "default"
    }
  ],
  "env_vars": []
}
```

</details>

或直接导入文件：
```bash
# 文件位置
test-workflow-chat-with-context.json
```

### 2. 测试对话（替换 YOUR_WORKFLOW_ID）

```bash
# 第一轮：告诉 AI 你的信息
curl -X POST http://localhost:7777/api/v1/workflows/YOUR_WORKFLOW_ID/execute \
  -H "Content-Type: application/json" \
  -d '{"params": {"session_id": "user123", "user_message": "你好，我叫张三"}}'

# 第二轮：测试 AI 是否记住
curl -X POST http://localhost:7777/api/v1/workflows/YOUR_WORKFLOW_ID/execute \
  -H "Content-Type: application/json" \
  -d '{"params": {"session_id": "user123", "user_message": "我叫什么名字？"}}'
```

✅ **成功标志**：AI 回答 "你叫张三"

## 📋 配置说明

在工作流编辑页面，选择 OpenAI Chat 节点，勾选 **"启用对话记忆"**：

- **会话ID**：`{{params.session_id}}` （自动从请求参数获取）
- **窗口条数**：10 （保留最近 10 条消息）
- **过期时间**：604800 秒（7 天）

## 🔍 查看对话历史（可选）

如果配置了 Redis：

```bash
# 查看所有会话
redis-cli KEYS "chat:context:*"

# 查看具体会话内容
redis-cli GET "chat:context:user123"
```

## 📁 文件清单

已生成以下文件供你测试：

- `test-workflow-chat-with-context.json` - 可导入的工作流配置
- `TEST_CHAT_CONTEXT.md` - 详细功能说明和故障排查
- `QUICK_START.md` - 本文件（快速开始）

## ❓ 常见问题

**Q: AI 没有记住历史？**
A: 确保同一用户使用相同的 `session_id`

**Q: 不同用户对话会混吗？**
A: 不会，不同 `session_id` 完全隔离

**Q: 对话会一直保存吗？**
A: 默认 7 天过期，可调整 `ttl_seconds`

---

完整文档请查看 `TEST_CHAT_CONTEXT.md`
