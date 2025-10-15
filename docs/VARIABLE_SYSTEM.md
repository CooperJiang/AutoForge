# 变量系统设计文档

> 统一的变量引用系统，支持所有工具节点引用前置节点输出和环境变量

## 📋 功能概述

变量系统是一个可复用的组件系统，允许工作流中的任何节点引用：
- 前置节点的输出数据
- 环境变量
- 触发器数据

## 🏗️ 架构设计

### 核心组件

#### 1. VariableSelector（单行输入）

**位置**：`web/src/components/VariableSelector.vue`

**用途**：用于单行文本输入字段，提供基础变量输入功能

**使用示例**：
```vue
<VariableSelector
  v-model="config.webhook_url"
  placeholder="输入 URL 或使用变量"
  :previous-nodes="previousNodes"
  :env-vars="envVars"
/>
```

**特点**：
- 基于 BaseInput 组件
- 自动支持变量语法高亮（未来可扩展）
- 简单轻量

#### 2. VariableHelper（变量助手面板）

**位置**：`web/src/components/VariableHelper.vue`

**用途**：用于多行文本输入（textarea），提供可视化的变量选择界面

**使用示例**：
```vue
<VariableHelper
  :show="showVariableHelper"
  :previous-nodes="previousNodes"
  :env-vars="envVars"
  @insert-field="(nodeId, fieldName) => insertFieldVariable(nodeId, fieldName, textareaRef)"
  @insert-node="(nodeId) => insertNodeVariable(nodeId, textareaRef)"
  @insert-env="(key) => insertEnvVariable(key, textareaRef)"
/>

<textarea
  ref="textareaRef"
  v-model="config.content"
  placeholder="输入内容...&#10;&#10;示例：{{nodes.node_xxx.message}}"
  class="font-mono text-sm"
/>
```

**特点**：
- 可折叠的辅助面板
- 显示前置节点的常见输出字段
- 点击按钮直接插入变量到 textarea
- 保持光标位置
- 统一的字段定义

### 数据结构

#### PreviousNodes 接口
```typescript
interface PreviousNode {
  id: string          // 节点 ID
  name: string        // 节点显示名称
  type: string        // 节点类型（tool, trigger, condition 等）
  toolCode?: string   // 工具代码（如 http_request, email_sender）
}
```

#### EnvVars 接口
```typescript
interface EnvVar {
  key: string         // 环境变量键名
  value: string       // 环境变量值
  description?: string // 描述信息
}
```

## 📖 变量语法

### 节点输出引用

```
{{nodes.节点ID.字段名}}
```

**示例**：
```
{{nodes.http_001.status}}       // HTTP 状态码
{{nodes.http_001.data.title}}   // 嵌套字段
{{nodes.email_001.message_id}}  // 消息 ID
```

### 环境变量引用

```
{{env.变量名}}
```

**示例**：
```
{{env.API_KEY}}
{{env.DATABASE_URL}}
```

### 部分匹配（用户手动补全）

```
{{nodes.节点ID.
```

用户可以手动输入前缀，系统会在执行时自动补全。

## 🎯 常见字段定义

VariableHelper 组件内置了各工具的常见输出字段定义：

### HTTP Request
- `status` - HTTP 状态码
- `message` - 返回消息
- `success` - 成功标识
- `data` - 返回数据
- `code` - 业务代码
- `error` - 错误信息

### Email Sender
- `success` - 发送成功
- `message_id` - 消息 ID
- `error` - 错误信息

### Health Checker
- `healthy` - 健康状态
- `status` - 检查状态
- `message` - 状态消息
- `latency` - 响应延迟

### Feishu Bot
- `success` - 发送成功
- `message` - 返回消息
- `error` - 错误信息

### Condition Node
- `result` - 条件判断结果 (true/false)
- `message` - 判断说明

### 通用字段
- `success` - 成功标识
- `message` - 消息
- `data` - 数据
- `error` - 错误

## 🔧 工具配置集成指南

### 步骤 1：更新 Props 接口

```typescript
interface Props {
  config: Record<string, any>
  previousNodes?: Array<{ id: string; name: string; type: string; toolCode?: string }>
  envVars?: Array<{ key: string; value: string; description?: string }>
}
```

### 步骤 2：导入组件

```typescript
import VariableSelector from '@/components/VariableSelector'
import VariableHelper from '@/components/VariableHelper'
```

### 步骤 3：添加变量助手状态

```typescript
const showVariableHelper = ref(false)
const textareaRef = ref<HTMLTextAreaElement>()

const formattedEnvVars = computed(() => {
  return props.envVars || []
})
```

### 步骤 4：实现插入函数

```typescript
// 插入字段变量
const insertFieldVariable = (nodeId: string, fieldName: string, targetRef?: { value?: HTMLTextAreaElement }) => {
  insertToTextarea(`{{nodes.${nodeId}.${fieldName}}}`, targetRef)
}

// 插入节点变量
const insertNodeVariable = (nodeId: string, targetRef?: { value?: HTMLTextAreaElement }) => {
  insertToTextarea(`{{nodes.${nodeId}.`, targetRef)
}

// 插入环境变量
const insertEnvVariable = (key: string, targetRef?: { value?: HTMLTextAreaElement }) => {
  insertToTextarea(`{{env.${key}}}`, targetRef)
}

// 插入变量到 textarea
const insertToTextarea = (text: string, targetRef?: { value?: HTMLTextAreaElement }) => {
  const textarea = targetRef?.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const currentValue = textarea.value || ''

  // 更新对应的 config 字段
  localConfig.value.someField = currentValue.substring(0, start) + text + currentValue.substring(end)

  // 恢复光标位置
  setTimeout(() => {
    textarea.focus()
    const newPos = start + text.length
    textarea.setSelectionRange(newPos, newPos)
  }, 0)
}
```

### 步骤 5：在模板中使用

**单行输入字段**：
```vue
<VariableSelector
  v-model="localConfig.title"
  placeholder="输入标题"
  :previous-nodes="previousNodes"
  :env-vars="formattedEnvVars"
/>
```

**多行输入字段**：
```vue
<label class="flex items-center justify-between">
  <span>内容</span>
  <button @click="showVariableHelper = !showVariableHelper">
    {{ showVariableHelper ? '隐藏' : '显示' }}变量助手
  </button>
</label>

<VariableHelper
  :show="showVariableHelper"
  :previous-nodes="previousNodes"
  :env-vars="formattedEnvVars"
  @insert-field="(nodeId, fieldName) => insertFieldVariable(nodeId, fieldName, textareaRef)"
  @insert-node="(nodeId) => insertNodeVariable(nodeId, textareaRef)"
  @insert-env="(key) => insertEnvVariable(key, textareaRef)"
/>

<textarea
  ref="textareaRef"
  v-model="localConfig.content"
  placeholder="输入内容...&#10;&#10;示例：{{nodes.node_xxx.message}}"
  class="font-mono text-sm"
/>
```

### 步骤 6：在 NodeConfigDrawer 中传递 Props

```vue
<YourToolConfig
  v-else-if="node.toolCode === 'your_tool'"
  v-model:config="localNode.config"
  :previous-nodes="props.previousNodes"
  :env-vars="props.envVars"
/>
```

## ✅ 已集成工具

- ✅ HTTP Request（内置于 NodeConfigDrawer）
- ✅ Email Sender
- ✅ Feishu Bot

## 🎨 样式指南

### 文本输入框

使用 `font-mono text-sm` 类以便用户清晰看到变量语法：

```vue
<textarea
  class="w-full px-3 py-2 border-2 border-border-primary rounded-lg
         focus:outline-none focus:border-primary bg-bg-elevated
         text-text-primary font-mono text-sm"
/>
```

### 变量助手按钮

```vue
<button
  @click="showVariableHelper = !showVariableHelper"
  class="text-xs text-primary hover:text-primary"
>
  {{ showVariableHelper ? '隐藏' : '显示' }}变量助手
</button>
```

### Placeholder 提示

在 placeholder 中提供变量使用示例：

```vue
placeholder="输入内容...&#10;&#10;示例：{{nodes.node_xxx.message}}"
```

## 🔄 后端变量解析

变量在工作流执行时由后端解析和替换。

**Go 后端解析逻辑**（位置待定）：
```go
// 解析变量
func ResolveVariables(template string, context *ExecutionContext) string {
    // 1. 解析 {{nodes.xxx.yyy}} 格式
    // 2. 解析 {{env.xxx}} 格式
    // 3. 支持嵌套字段访问
    // 4. 处理不存在的字段（返回空字符串或保留原始变量）
}
```

## 🚀 未来扩展

### 1. 语法高亮
在输入框中为变量语法添加颜色高亮

### 2. 自动补全
输入 `{{` 时弹出补全菜单

### 3. 变量验证
实时验证变量语法和字段是否存在

### 4. 变量预览
鼠标悬停在变量上显示其当前值

### 5. 更多数据源
- 触发器数据
- 全局变量
- 时间函数（如 `{{now()}}`, `{{date()}}`)

---

**维护者**: AutoForge Team
**最后更新**: 2025-01-13
