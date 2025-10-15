# 🔧 Cooper 工具开发指南

> 完整的工具开发教程 - 从后端实现到前端配置界面

## 📋 目录

- [概述](#-概述)
- [后端开发](#-后端开发)
- [前端开发](#-前端开发)
- [完整示例](#-完整示例)
- [最佳实践](#-最佳实践)
- [常见问题](#-常见问题)

---

## 🎯 概述

Cooper 使用插件化的工具系统，每个工具都是一个独立的模块。工具开发分为三个部分：

1. **后端工具实现** - 实现工具的核心逻辑（Go）
2. **前端图标配置** - 配置工具的图标、标签和使用说明（TypeScript）
3. **前端配置界面** - 为工具提供用户友好的配置界面（Vue 3）

### 工具系统架构

```
┌─────────────────────────────────────────────┐
│    前端图标配置 (web/src/config/tools.ts)    │
│  定义工具的图标、标签、使用说明等前端元数据     │
└──────────────┬──────────────────────────────┘
               │
               ▼
┌─────────────────────────────────────────────┐
│           前端配置界面 (Vue 3)               │
│  用户通过表单配置工具参数                     │
└──────────────┬──────────────────────────────┘
               │ 工具配置 (JSON)
               ▼
┌─────────────────────────────────────────────┐
│           工具注册表 (Registry)               │
│  管理所有已注册的工具                         │
└──────────────┬──────────────────────────────┘
               │ 工具实例
               ▼
┌─────────────────────────────────────────────┐
│         工具实现 (Tool Interface)            │
│  执行具体的业务逻辑                           │
└─────────────────────────────────────────────┘
```

### 开发流程总结

1. **后端开发** (`pkg/utools/your_tool/`)
   - 创建工具文件并实现 `Tool` 接口
   - 定义 `ToolMetadata`（无需配置 Icon 字段）
   - 定义 `ConfigSchema`
   - 实现 `Execute()` 方法
   - 在 `pkg/utools/init.go` 中导入

2. **前端图标配置** (`web/src/config/tools.ts`)
   - 导入 Lucide 图标或准备自定义图标图片
   - 添加工具配置对象（图标、背景色、标签、使用说明）

3. **前端配置组件** (`web/src/components/tools/`)
   - 创建 Vue 组件实现配置表单
   - 在编辑器中注册配置组件

---

## 🔨 后端开发

### 第一步：创建工具目录

在 `pkg/utools/` 下创建你的工具目录：

```bash
mkdir pkg/utools/your_tool
cd pkg/utools/your_tool
```

### 第二步：实现工具接口

创建 `your_tool.go` 文件，实现 `Tool` 接口：

```go
package your_tool

import (
    "auto-forge/pkg/utools"
    "fmt"
    "time"
)

// YourTool 你的工具实现
type YourTool struct {
    *utools.BaseTool
}

// NewYourTool 创建工具实例
func NewYourTool() *YourTool {
    // 1. 定义工具元数据
    metadata := &utools.ToolMetadata{
        Code:        "your_tool",           // 唯一标识，小写下划线
        Name:        "你的工具",             // 显示名称
        Description: "工具功能描述",         // 详细描述
        Category:    "automation",          // 分类: network, notification, data, automation
        Version:     "1.0.0",               // 版本号
        Author:      "Your Name",           // 作者
        AICallable:  true,                  // 是否可被 AI 调用
        Tags:        []string{"tag1", "tag2"}, // 标签
    }

    // 2. 定义配置 Schema（JSON Schema 格式）
    schema := &utools.ConfigSchema{
        Type: "object",
        Properties: map[string]utools.PropertySchema{
            "param1": {
                Type:        "string",
                Title:       "参数1",
                Description: "参数1的描述",
                Default:     "默认值",
            },
            "param2": {
                Type:        "number",
                Title:       "参数2",
                Description: "数字类型参数",
                Default:     10.0,
                Minimum:     func() *float64 { v := 1.0; return &v }(),
                Maximum:     func() *float64 { v := 100.0; return &v }(),
            },
            "param3": {
                Type:        "string",
                Title:       "选项参数",
                Description: "下拉选择参数",
                Enum:        []interface{}{"option1", "option2", "option3"},
                Default:     "option1",
            },
            "secret_param": {
                Type:        "string",
                Title:       "敏感参数",
                Description: "API Key 等敏感信息",
                Secret:      true, // 标记为敏感信息，会加密存储
            },
        },
        Required: []string{"param1"}, // 必填字段
    }

    return &YourTool{
        BaseTool: utools.NewBaseTool(metadata, schema),
    }
}

// Execute 执行工具
func (t *YourTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
    startTime := time.Now()

    // 1. 解析配置参数
    param1, _ := config["param1"].(string)
    param2, _ := config["param2"].(float64)
    param3, _ := config["param3"].(string)

    // 2. 执行工具逻辑
    // ... 你的业务代码 ...

    // 3. 返回执行结果
    return &utools.ExecutionResult{
        Success:    true,
        Message:    "执行成功",
        Output: map[string]interface{}{
            "result": "执行结果",
            "param1": param1,
            "param2": param2,
        },
        DurationMs: time.Since(startTime).Milliseconds(),
    }, nil
}

// init 自动注册工具
func init() {
    tool := NewYourTool()
    if err := utools.Register(tool); err != nil {
        panic(fmt.Sprintf("Failed to register your tool: %v", err))
    }
}
```

### 第三步：配置 Schema 详解

Schema 定义了工具的配置参数，支持多种数据类型：

#### 字符串类型 (string)

```go
"url": {
    Type:        "string",
    Title:       "URL 地址",
    Description: "目标 URL",
    Format:      "uri",           // 格式验证: uri, email, date-time
    Pattern:     "^https://.*",   // 正则表达式验证
    MinLength:   func() *int { v := 5; return &v }(),
    MaxLength:   func() *int { v := 200; return &v }(),
}
```

#### 数字类型 (number)

```go
"timeout": {
    Type:        "number",
    Title:       "超时时间",
    Description: "请求超时（秒）",
    Default:     30.0,
    Minimum:     func() *float64 { v := 1.0; return &v }(),
    Maximum:     func() *float64 { v := 300.0; return &v }(),
}
```

#### 布尔类型 (boolean)

```go
"enabled": {
    Type:        "boolean",
    Title:       "启用",
    Description: "是否启用此功能",
    Default:     true,
}
```

#### 枚举类型 (enum)

```go
"method": {
    Type:        "string",
    Title:       "请求方法",
    Description: "HTTP 请求方法",
    Enum:        []interface{}{"GET", "POST", "PUT", "DELETE"},
    Default:     "GET",
}
```

#### 对象类型 (object)

```go
"headers": {
    Type:        "object",
    Title:       "请求头",
    Description: "HTTP 请求头",
    Properties: map[string]utools.PropertySchema{
        "Authorization": {
            Type:  "string",
            Title: "认证令牌",
        },
    },
}
```

#### 数组类型 (array)

```go
"tags": {
    Type:        "array",
    Title:       "标签",
    Description: "标签列表",
    Items: &utools.PropertySchema{
        Type: "string",
    },
}
```

#### 敏感信息 (secret)

```go
"api_key": {
    Type:        "string",
    Title:       "API Key",
    Description: "第三方服务的 API 密钥",
    Secret:      true, // 标记为敏感信息，会加密存储
}
```

### 第四步：导入工具模块

在 `pkg/utools/init.go` 中导入你的工具：

```go
package utools

import (
    _ "auto-forge/pkg/utools/http"
    _ "auto-forge/pkg/utools/email"
    _ "auto-forge/pkg/utools/health"
    _ "auto-forge/pkg/utools/your_tool"  // 添加这行
)
```

### 第五步：测试工具

创建单元测试文件 `your_tool_test.go`：

```go
package your_tool

import (
    "auto-forge/pkg/utools"
    "context"
    "testing"
)

func TestYourTool_Execute(t *testing.T) {
    tool := NewYourTool()

    ctx := &utools.ExecutionContext{
        Context: context.Background(),
        TaskID:  "test-task",
        UserID:  "test-user",
    }

    config := map[string]interface{}{
        "param1": "test value",
        "param2": 20.0,
        "param3": "option1",
    }

    result, err := tool.Execute(ctx, config)
    if err != nil {
        t.Fatalf("Execute failed: %v", err)
    }

    if !result.Success {
        t.Errorf("Expected success=true, got %v", result.Success)
    }
}

func TestYourTool_Validate(t *testing.T) {
    tool := NewYourTool()

    // 测试必填字段缺失
    config := map[string]interface{}{
        "param2": 20.0,
    }

    err := tool.Validate(config)
    if err == nil {
        t.Error("Expected validation error for missing required field")
    }
}
```

运行测试：

```bash
go test ./pkg/utools/your_tool -v
```

---

## 🎨 前端开发

### 第一步：添加工具图标和元数据配置

在 `web/src/config/tools.ts` 中为你的工具添加前端配置（图标、标签、使用说明等）：

```typescript
import { YourIcon } from 'lucide-vue-next' // 导入 Lucide 图标

export const TOOL_CONFIGS: Record<string, ToolConfig> = {
  // ... 其他工具配置

  your_tool: {
    code: 'your_tool',
    title: '🔧 你的工具',
    description: '工具功能描述',
    icon: YourIcon,  // 使用 Lucide 图标组件
    // 或使用图片路径: icon: '/icons/your-tool.png',
    iconBg: 'bg-gradient-to-br from-blue-500 to-indigo-600',  // 图标背景渐变色
    usageDescription: '详细的使用说明',
    usageItems: [
      { text: '功能特点 1' },
      { text: '功能特点 2' },
      { text: '适用场景：XXX' },
    ],
    tags: ['tag1', 'tag2', 'tag3'],
  },
}
```

**配置说明：**
- `icon`: 支持两种方式
  - **Lucide 图标**（推荐）：从 `lucide-vue-next` 导入图标组件
  - **自定义图片**：提供图片路径，如 `'/icons/tool.png'`，图片放在 `web/public/icons/` 目录
- `iconBg`: Tailwind CSS 渐变背景类名，用于图标背景色
- `usageDescription`: 工具使用说明
- `usageItems`: 工具特点和使用场景的列表
- `tags`: 工具标签，用于搜索和分类

**常用 Lucide 图标：**
- `Globe` - 网络/HTTP 相关
- `Mail` - 邮件相关
- `Activity` - 监控/健康检查
- `Shuffle` - 转换/处理
- `Zap` - 缓存/性能
- `MessageSquare` - 消息/通知
- `Palette` - 格式化/输出
- `FileJson` - 文件/数据
- `MessageCircle` - 对话/聊天
- `Image` - 图片相关

完整的图标列表：https://lucide.dev/icons/

### 第二步：创建配置组件

在 `web/src/components/tools/` 下创建 `YourToolConfig.vue`：

```vue
<template>
  <div class="space-y-4">
    <h3 class="text-sm font-semibold text-text-primary mb-3">你的工具配置</h3>

    <!-- 参数1 - 文本输入 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        参数1 <span class="text-error">*</span>
      </label>
      <input
        v-model="localConfig.param1"
        type="text"
        placeholder="请输入参数1"
        class="w-full px-3 py-2 text-sm border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary"
      />
    </div>

    <!-- 参数2 - 数字输入 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        参数2
      </label>
      <input
        v-model.number="localConfig.param2"
        type="number"
        :min="1"
        :max="100"
        placeholder="输入数字 (1-100)"
        class="w-full px-3 py-2 text-sm border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary"
      />
    </div>

    <!-- 参数3 - 下拉选择 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        选项参数
      </label>
      <select
        v-model="localConfig.param3"
        class="w-full px-3 py-2 text-sm border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary"
      >
        <option value="option1">选项1</option>
        <option value="option2">选项2</option>
        <option value="option3">选项3</option>
      </select>
    </div>

    <!-- 敏感参数 - 密码输入 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        敏感参数
      </label>
      <input
        v-model="localConfig.secret_param"
        type="password"
        placeholder="输入 API Key"
        class="w-full px-3 py-2 text-sm border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary"
      />
      <p class="mt-1 text-xs text-text-tertiary">此信息将被加密存储</p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'

interface Props {
  config: Record<string, any>
}

interface Emits {
  (e: 'update:config', value: Record<string, any>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

// 本地配置状态
const localConfig = ref({
  param1: props.config.param1 || '',
  param2: props.config.param2 || 10,
  param3: props.config.param3 || 'option1',
  secret_param: props.config.secret_param || '',
})

// 监听配置变化并向父组件发送更新
watch(
  localConfig,
  (newConfig) => {
    emit('update:config', newConfig)
  },
  { deep: true }
)
</script>
```

### 第三步：注册配置组件

在工作流编辑器中注册你的配置组件。

找到 `web/src/pages/Workflows/editor.vue`，在工具配置部分添加你的组件：

```vue
<script setup lang="ts">
import YourToolConfig from '@/components/tools/YourToolConfig.vue'

// 工具配置组件映射
const toolConfigComponents = {
  http_request: HttpRequestConfig,
  send_email: EmailToolConfig,
  health_checker: HealthCheckerConfig,
  your_tool: YourToolConfig,  // 添加这行
}
</script>
```

**注意：** 工具会自动从后端 API 加载并显示在工具列表和工具面板中，图标和元数据会使用 `web/src/config/tools.ts` 中的配置。无需手动在工具面板中添加。

### 组件开发最佳实践

#### 1. 使用受控组件

```vue
<script setup lang="ts">
// ✅ 正确：使用本地状态 + watch
const localConfig = ref({ ...props.config })

watch(
  localConfig,
  (newConfig) => {
    emit('update:config', newConfig)
  },
  { deep: true }
)

// ❌ 错误：直接修改 props
const updateParam = (value) => {
  props.config.param1 = value // 这样做是错误的！
}
</script>
```

#### 2. 提供输入验证

```vue
<template>
  <div>
    <input
      v-model="localConfig.email"
      type="email"
      :class="{ 'border-error': !isValidEmail }"
    />
    <p v-if="!isValidEmail" class="text-xs text-error mt-1">
      请输入有效的邮箱地址
    </p>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const isValidEmail = computed(() => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(localConfig.value.email)
})
</script>
```

#### 3. 添加帮助文本

```vue
<template>
  <div>
    <label class="flex items-center gap-1">
      参数名称
      <Tooltip text="这是参数的详细说明" position="right">
        <HelpCircle class="w-3.5 h-3.5 text-text-tertiary" />
      </Tooltip>
    </label>
    <input v-model="localConfig.param" />
    <p class="text-xs text-text-tertiary mt-1">
      提示：这个参数的使用建议
    </p>
  </div>
</template>
```

---

## 📝 完整示例：短信发送工具

让我们通过一个完整的示例来演示整个开发流程。

### 后端实现

`pkg/utools/sms/sms_tool.go`:

```go
package sms

import (
    "auto-forge/pkg/utools"
    "fmt"
    "time"
)

type SMSTool struct {
    *utools.BaseTool
}

func NewSMSTool() *SMSTool {
    metadata := &utools.ToolMetadata{
        Code:        "send_sms",
        Name:        "发送短信",
        Description: "通过短信服务商发送短信通知",
        Category:    "notification",
        Version:     "1.0.0",
        Author:      "Cooper Team",
        AICallable:  true,
        Tags:        []string{"sms", "notification", "message"},
    }

    schema := &utools.ConfigSchema{
        Type: "object",
        Properties: map[string]utools.PropertySchema{
            "phone": {
                Type:        "string",
                Title:       "手机号码",
                Description: "接收短信的手机号码",
                Pattern:     "^1[3-9]\\d{9}$",
            },
            "message": {
                Type:        "string",
                Title:       "短信内容",
                Description: "要发送的短信内容",
                MinLength:   func() *int { v := 1; return &v }(),
                MaxLength:   func() *int { v := 500; return &v }(),
            },
            "provider": {
                Type:        "string",
                Title:       "服务商",
                Description: "短信服务提供商",
                Enum:        []interface{}{"阿里云", "腾讯云", "华为云"},
                Default:     "阿里云",
            },
            "api_key": {
                Type:        "string",
                Title:       "API Key",
                Description: "短信服务商的 API 密钥",
                Secret:      true,
            },
        },
        Required: []string{"phone", "message", "api_key"},
    }

    return &SMSTool{
        BaseTool: utools.NewBaseTool(metadata, schema),
    }
}

func (t *SMSTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
    startTime := time.Now()

    phone, _ := config["phone"].(string)
    message, _ := config["message"].(string)
    provider, _ := config["provider"].(string)
    apiKey, _ := config["api_key"].(string)

    // 模拟发送短信（实际应该调用短信服务商 API）
    // ... 实际的短信发送代码 ...

    return &utools.ExecutionResult{
        Success: true,
        Message: fmt.Sprintf("短信已发送至 %s", phone),
        Output: map[string]interface{}{
            "phone":    phone,
            "provider": provider,
            "sent_at":  time.Now().Unix(),
        },
        DurationMs: time.Since(startTime).Milliseconds(),
    }, nil
}

func init() {
    tool := NewSMSTool()
    if err := utools.Register(tool); err != nil {
        panic(fmt.Sprintf("Failed to register SMS tool: %v", err))
    }
}
```

### 前端图标配置

`web/src/config/tools.ts` 中添加 SMS 工具配置：

```typescript
import { MessageSquare } from 'lucide-vue-next'

export const TOOL_CONFIGS: Record<string, ToolConfig> = {
  // ... 其他工具

  send_sms: {
    code: 'send_sms',
    title: '📱 发送短信',
    description: '通过短信服务商发送短信通知',
    icon: MessageSquare,
    iconBg: 'bg-gradient-to-br from-green-500 to-emerald-600',
    usageDescription: '支持多家短信服务商发送短信通知',
    usageItems: [
      { text: '支持阿里云、腾讯云、华为云等主流短信服务商' },
      { text: '自动验证手机号格式' },
      { text: 'API Key 加密存储，安全可靠' },
      { text: '适用场景：验证码发送、通知提醒、营销推广等' },
    ],
    tags: ['SMS', 'Notification', 'Message', 'Alert'],
  },
}
```

### 前端配置组件实现

`web/src/components/tools/SMSToolConfig.vue`:

```vue
<template>
  <div class="space-y-4">
    <h3 class="text-sm font-semibold text-text-primary mb-3">短信发送配置</h3>

    <!-- 手机号 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        手机号码 <span class="text-error">*</span>
      </label>
      <input
        v-model="localConfig.phone"
        type="tel"
        placeholder="请输入手机号码"
        :class="[
          'w-full px-3 py-2 text-sm border rounded-md focus:ring-2 focus:ring-primary focus:border-primary',
          !isValidPhone && localConfig.phone ? 'border-error' : 'border-border-primary'
        ]"
      />
      <p v-if="!isValidPhone && localConfig.phone" class="text-xs text-error mt-1">
        请输入有效的手机号码
      </p>
    </div>

    <!-- 短信内容 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        短信内容 <span class="text-error">*</span>
      </label>
      <textarea
        v-model="localConfig.message"
        rows="4"
        placeholder="请输入短信内容"
        maxlength="500"
        class="w-full px-3 py-2 text-sm border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary resize-none"
      ></textarea>
      <p class="text-xs text-text-tertiary mt-1">
        {{ localConfig.message?.length || 0 }} / 500 字
      </p>
    </div>

    <!-- 服务商选择 -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        服务商
      </label>
      <select
        v-model="localConfig.provider"
        class="w-full px-3 py-2 text-sm border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary"
      >
        <option value="阿里云">阿里云</option>
        <option value="腾讯云">腾讯云</option>
        <option value="华为云">华为云</option>
      </select>
    </div>

    <!-- API Key -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        API Key <span class="text-error">*</span>
      </label>
      <input
        v-model="localConfig.api_key"
        type="password"
        placeholder="输入短信服务商的 API Key"
        class="w-full px-3 py-2 text-sm border border-border-primary rounded-md focus:ring-2 focus:ring-primary focus:border-primary"
      />
      <p class="text-xs text-text-tertiary mt-1">
        <Lock class="w-3 h-3 inline" /> 此信息将被加密存储
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { Lock } from 'lucide-vue-next'

interface Props {
  config: Record<string, any>
}

interface Emits {
  (e: 'update:config', value: Record<string, any>): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const localConfig = ref({
  phone: props.config.phone || '',
  message: props.config.message || '',
  provider: props.config.provider || '阿里云',
  api_key: props.config.api_key || '',
})

// 手机号验证
const isValidPhone = computed(() => {
  if (!localConfig.value.phone) return true
  return /^1[3-9]\d{9}$/.test(localConfig.value.phone)
})

watch(
  localConfig,
  (newConfig) => {
    emit('update:config', newConfig)
  },
  { deep: true }
)
</script>
```

---

## 🎯 最佳实践

### 1. 错误处理

```go
func (t *YourTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
    startTime := time.Now()

    // 验证配置
    if err := t.Validate(config); err != nil {
        return &utools.ExecutionResult{
            Success:    false,
            Message:    "配置验证失败",
            Error:      err.Error(),
            DurationMs: time.Since(startTime).Milliseconds(),
        }, err
    }

    // 检查上下文超时
    select {
    case <-ctx.Context.Done():
        return &utools.ExecutionResult{
            Success:    false,
            Message:    "执行超时",
            Error:      "context deadline exceeded",
            DurationMs: time.Since(startTime).Milliseconds(),
        }, ctx.Context.Err()
    default:
    }

    // 执行业务逻辑...
}
```

### 2. 日志记录

```go
import log "auto-forge/pkg/logger"

func (t *YourTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
    log.Info("执行工具: %s, 任务ID: %s", t.GetMetadata().Name, ctx.TaskID)

    // ... 执行逻辑 ...

    log.Info("工具执行成功: %s, 耗时: %dms", t.GetMetadata().Name, result.DurationMs)
    return result, nil
}
```

### 3. 敏感信息处理

后端会自动加密标记为 `Secret: true` 的字段，前端使用 `type="password"` 输入框。

### 4. 测试覆盖

- 单元测试：测试工具的核心逻辑
- 集成测试：测试工具在实际环境中的运行
- 配置验证测试：测试各种配置参数的验证逻辑

---

## ❓ 常见问题

### Q1: 工具未出现在工具列表中

**原因**：工具未正确注册

**解决方案**：
1. 检查 `init()` 函数是否正确调用 `utools.Register()`
2. 检查 `pkg/utools/init.go` 是否导入了你的工具包
3. 重新编译后端：`go build ./cmd/main.go`

### Q2: 配置组件不显示

**原因**：组件未在编辑器中注册

**解决方案**：
在 `web/src/pages/Workflows/editor.vue` 的 `toolConfigComponents` 中添加你的组件映射。

### Q3: 参数验证失败

**原因**：Schema 定义与实际配置不匹配

**解决方案**：
1. 检查 `Required` 字段是否正确
2. 检查字段类型是否匹配
3. 检查枚举值、最小/最大值等约束

### Q4: 如何调试工具执行

**后端调试**：
```go
log.Debug("配置参数: %+v", config)
log.Debug("执行上下文: TaskID=%s, UserID=%s", ctx.TaskID, ctx.UserID)
```

**前端调试**：
```vue
<script setup lang="ts">
watch(localConfig, (newConfig) => {
  console.log('配置更新:', newConfig)
  emit('update:config', newConfig)
}, { deep: true })
</script>
```

---

## 📚 参考资源

- [JSON Schema 规范](https://json-schema.org/)
- [Vue 3 组合式 API](https://cn.vuejs.org/guide/extras/composition-api-faq.html)
- [Go 接口最佳实践](https://golang.org/doc/effective_go#interfaces)
- [Cooper 组件开发规范](./COMPONENT_DEVELOPMENT.md)

---

## 🤝 需要帮助？

如果在工具开发过程中遇到问题：

1. 查看现有工具的实现作为参考（`pkg/utools/http/`, `pkg/utools/email/`）
2. 阅读核心接口定义（`pkg/utools/types.go`）
3. 提交 Issue：[GitHub Issues](https://github.com/CooperJiang/Cooper/issues)

---

**文档版本**: v1.1
**最后更新**: 2025-01-15
**维护者**: Cooper Team

## 📝 更新日志

### v1.1 (2025-01-15)
- 移除后端工具元数据中的 `Icon` 字段
- 新增前端工具图标配置系统 (`web/src/config/tools.ts`)
- 支持 Lucide 图标组件和自定义图片路径
- 添加工具图标配置最佳实践和常用图标列表

### v1.0 (2025-01-13)
- 初始版本发布
