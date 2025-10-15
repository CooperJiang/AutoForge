# 工作流节点说明

## 📊 完整流程图

```
外部请求                   自动处理                     返回结果
   ↓                         ↓                           ↓
┌─────────��             ┌──────────┐                ┌──────────┐
│ 触发节点 │────────────→│ OpenAI   │───────────────→│ API 返回 │
│         │  传递参数    │  Chat    │  返回 content   │   结果   │
└─────────┘             └──────────┘                └──────────┘
                             ↕
                        自动读写 Redis
                    （引擎中间件管理）
```

---

## 🔍 节点详解

### 节点 1: 外部 API 触发（trigger_start）

**作用**：接收外部 HTTP 请求

**输入参数**（通过 POST Body 传入）：
```json
{
  "params": {
    "session_id": "用户唯一标识",
    "user_input": "用户发送的消息"
  }
}
```

**输出**：将 `params` 传递给下游节点

---

### 节点 2: OpenAI Chat（chat_with_memory）

**作用**：这是**唯一的核心节点**，完成以下所有工作：

#### 第 1 步：自动读取历史（引擎中间件）
```
if (context_config.enabled) {
  从 Redis 读取 key = "chat:context:{session_id}"
  获取历史消息数组 [
    {role: "system", content: "..."},
    {role: "user", content: "上一轮对话"},
    {role: "assistant", content: "上一轮回复"},
    ...
  ]
}
```

#### 第 2 步：拼接当前消息（引擎中间件）
```
messages.push({
  role: "user",
  content: "{{params.user_input}}"  // 当前用户输入
})

自动裁剪到 window_size = 10 条
```

#### 第 3 步：调用 OpenAI API（工具执行）
```
POST https://api.openai.com/v1/chat/completions
{
  "model": "gpt-3.5-turbo",
  "messages": [...拼接好的完整对话历史],
  "temperature": 0.7
}
```

#### 第 4 步：获取 AI 回复
```
OpenAI 返回：
{
  "choices": [{
    "message": {
      "role": "assistant",
      "content": "AI 的回复内容"
    }
  }]
}
```

#### 第 5 步：自动保存历史（引擎中间件）
```
messages.push({
  role: "assistant",
  content: "AI 的回复内容"
})

保存到 Redis:
key = "chat:context:{session_id}"
value = JSON.stringify(messages)
TTL = 604800 秒（7天）
```

#### 第 6 步：返回结果
```
节点输出：
{
  "content": "AI 的回复内容",
  "model": "gpt-3.5-turbo",
  "finish_reason": "stop",
  "total_tokens": 156,
  ...
}
```

---

## �� API 返回结构

当你调用工作流 API 时，会得到以下返回：

```json
{
  "success": true,
  "message": "工作流执行成功",
  "output": {
    "chat_with_memory": {
      "content": "AI 的完整回复内容",  // ← 这就是你要的聊天内容
      "model": "gpt-3.5-turbo",
      "prompt_tokens": 85,
      "completion_tokens": 71,
      "total_tokens": 156,
      "finish_reason": "stop",
      "response": {
        "id": "chatcmpl-xxx",
        "object": "chat.completion",
        "created": 1234567890,
        "choices": [...],
        "usage": {...}
      }
    }
  },
  "execution_id": "xxx-xxx-xxx",
  "duration_ms": 1234
}
```

**关键字段**：
- `output.chat_with_memory.content` - AI 的回复文本
- `output.chat_with_memory.model` - 使用的模型
- `output.chat_with_memory.total_tokens` - 消耗的 token 数

---

## 🧪 完整测试示例

### 导入工作流

文件：`openai_redis_memory_v2_fixed.json`

### 测试对话（假设工作流 ID = `wf123`）

#### 对话 1：自我介绍

```bash
curl -X POST http://localhost:7777/api/v1/workflows/wf123/execute \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "session_id": "user_alice",
      "user_input": "你好！我叫 Alice，我是一名设计师，喜欢旅游和摄影。"
    }
  }'
```

**返回示例**：
```json
{
  "success": true,
  "output": {
    "chat_with_memory": {
      "content": "你好 Alice！很高兴认识你。设计师的工作一定很有创意，而且旅游和摄影的爱好也能为你的设计带来很多灵感吧！有什么我可以帮助你的吗？",
      "model": "gpt-3.5-turbo",
      "total_tokens": 142
    }
  }
}
```

#### 对话 2：测试记忆

```bash
curl -X POST http://localhost:7777/api/v1/workflows/wf123/execute \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "session_id": "user_alice",
      "user_input": "我的名字是什么？我的职业和爱好是什么？"
    }
  }'
```

**返回示例**：
```json
{
  "success": true,
  "output": {
    "chat_with_memory": {
      "content": "你的名字是 Alice，你是一名设计师，喜欢旅游和摄影。",
      "model": "gpt-3.5-turbo",
      "total_tokens": 98
    }
  }
}
```

✅ **成功！** AI 准确记住了你的信息！

#### 对话 3：基于上下文的推荐

```bash
curl -X POST http://localhost:7777/api/v1/workflows/wf123/execute \
  -H "Content-Type: application/json" \
  -d '{
    "params": {
      "session_id": "user_alice",
      "user_input": "根据我的兴趣，你能推荐一些旅游目的地吗？"
    }
  }'
```

**返回示例**：
```json
{
  "success": true,
  "output": {
    "chat_with_memory": {
      "content": "基于你喜欢摄影和旅游，我推荐以下目的地：\n1. 冰岛 - 极光和壮丽的自然风光\n2. 日本京都 - 传统建筑和四季美景\n3. 新西兰南岛 - 湖光山色，摄影师的天堂\n4. 摩洛哥 - 色彩斑斓的街道和异域风情\n5. 挪威峡湾 - 震撼的自然景观",
      "model": "gpt-3.5-turbo",
      "total_tokens": 215
    }
  }
}
```

✅ **成功！** AI 基于你是设计师 + 喜欢摄影的背景，给出了个性化推荐！

---

## 🔍 Redis 存储查看

```bash
# 查看 Alice 的对话历史
redis-cli GET "chat:context:user_alice"
```

**返回**：
```json
[
  {
    "role": "system",
    "content": "你是一个友好、专业的 AI 助手..."
  },
  {
    "role": "user",
    "content": "你好！我叫 Alice，我是一名设计师，喜欢旅游和摄影。"
  },
  {
    "role": "assistant",
    "content": "你好 Alice！很高兴认识你..."
  },
  {
    "role": "user",
    "content": "我的名字是什么？我的职业和爱好是什么？"
  },
  {
    "role": "assistant",
    "content": "你的名字是 Alice，你是一名设计师，喜欢旅游和摄影。"
  },
  {
    "role": "user",
    "content": "根据我的兴趣，你能推荐一些旅游目的地吗？"
  },
  {
    "role": "assistant",
    "content": "基于你喜欢摄影和旅游，我推荐以下目的地..."
  }
]
```

---

## ❓ 常见疑问解答

### Q: 为什么只有 2 个节点？

**A**: 因为**引擎中间件自动处理了一切**！

旧架构需要 5 个节点：
1. 触发
2. Redis GET（读历史）
3. OpenAI Chat（对话）
4. JSON Transform（拼接历史）
5. Redis SET（保存历史）

新架构只需 2 个节点：
1. 触发
2. OpenAI Chat（**引擎自动完成 2~5 步**）

### Q: OpenAI 节点到底做了什么？

**A**: OpenAI 节点本身只负责**调用 OpenAI API**，但引擎在执行前后会自动：

**执行前**（中间件）：
- 检测到 `context_config.enabled = true`
- 从 Redis 读取历史
- 拼接到 `messages_json`
- 传给 OpenAI 工具

**执行中**（工具）：
- 调用 OpenAI API
- 获取 AI 回复

**执行后**（中间件）：
- 将 AI 回复追加到历史
- 保存到 Redis

### Q: 接口返回的就是聊天内容吗？

**A**: 是的！返回的 `output.chat_with_memory.content` 就是 AI 的完整回复。

### Q: 如果不需要记忆功能呢？

**A**: 在节点配置中，取消勾选"启用对话记忆"即可。此时就是普通的单轮对话。

---

## 📁 文件清单

- ✅ `openai_redis_memory_v2_fixed.json` - 修复后的工作流配置
- ✅ `WORKFLOW_EXPLANATION.md` - 本文档（节点详解）
- ✅ `MIGRATION_GUIDE.md` - 新旧架构对比
- ✅ `TEST_CHAT_CONTEXT.md` - 完整技术文档

---

## 🎯 总结

**新架构只需 2 个节点**：

1. **外部触发** - 接收参数
2. **OpenAI Chat** - 自动记忆 + 对话 + 返回结果

**引擎中间件自动处理**：
- ✅ Redis 读取
- ✅ 消息拼接
- ✅ 窗口裁剪
- ✅ Redis 保存

**用户只需关心**：
- 传入 `session_id` 和 `user_input`
- 获取 `content`（AI 回复）

就是这么简单！���
