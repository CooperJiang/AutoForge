# 🔧 Cooper 工具开发指南

> 快速开发高质量工具的核心指南

## 📋 目录

- [开发流程](#-开发流程)
- [后端开发](#-后端开发)
- [前端开发](#-前端开发)
- [开发规范](#-开发规范)
- [检查清单](#-检查清单)

---

## 🎯 开发流程

### 三步完成工具开发

1. **后端实现** (`pkg/utools/your_tool/`) - 实现 Tool 接口，定义配置和输出
2. **前端图标** (`web/src/config/tools.ts`) - 配置图标、标签、使用说明
3. **前端组件** (`web/src/components/tools/`) - 创建配置表单并注册

---

## 🔨 后端开发

### 核心代码结构

```go
package your_tool

import "auto-forge/pkg/utools"

type YourTool struct {
    *utools.BaseTool
}

func NewYourTool() *YourTool {
    metadata := &utools.ToolMetadata{
        Code:        "your_tool",
        Name:        "你的工具",
        Description: "功能描述",
        Category:    "automation",  // network/notification/data/automation
        Version:     "1.0.0",
        OutputFieldsSchema: map[string]utools.OutputFieldDef{
            "response": {  // ⚠️ 必须包含
                Type:  "object",
                Label: "完整响应",
                Children: map[string]utools.OutputFieldDef{
                    "url": {Type: "string", Label: "URL"},
                },
            },
            "url": {Type: "string", Label: "URL（快捷访问）"},
        },
    }

    schema := &utools.ConfigSchema{
        Type: "object",
        Properties: map[string]utools.PropertySchema{
            "param": {
                Type:    "string",
                Title:   "参数名",
                Secret:  false,  // true 表示敏感信息
            },
        },
        Required: []string{"param"},
    }

    return &YourTool{BaseTool: utools.NewBaseTool(metadata, schema)}
}

func (t *YourTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
    // 1. 解析配置
    param := config["param"].(string)
    
    // 2. 执行逻辑
    result := doSomething(param)
    
    // 3. 返回结果（必须包含 response 字段）
    return &utools.ExecutionResult{
        Success: true,
        Message: "成功",
        Output: map[string]interface{}{
            "response": result,  // 完整响应
            "url": result.URL,   // 快捷访问
        },
    }, nil
}

func init() {
    utools.Register(NewYourTool())
}
```

### 注册工具

⚠️ **在 `cmd/main.go` 中导入**（不是 `init.go`）：

```go
import (
    // ... 其他导入
    _ "auto-forge/pkg/utools/your_tool"
)
```

### 后端配置管理（可选）

如果工具需要后端配置（如 API Key、Endpoint），需要更新两个文件：

#### 1. `pkg/config/config.go` - 定义配置结构

```go
type Config struct {
    // ... 其他配置
    YourTool YourToolConfig `yaml:"your_tool" env:"YOUR_TOOL"`
}

type YourToolConfig struct {
    APIKey  string `yaml:"api_key" env:"API_KEY"`
    BaseURL string `yaml:"base_url" env:"BASE_URL"`
    Enabled bool   `yaml:"enabled" env:"ENABLED"`
}
```

#### 2. `config.yaml` / `config.example.yaml` - 添加配置项

```yaml
# 你的工具配置
your_tool:
  api_key: ""
  base_url: "https://api.example.com"
  enabled: false
```

#### 3. Execute 方法中读取配置

```go
func (t *YourTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
    cfg := config.GetConfig()
    
    if !cfg.YourTool.Enabled {
        return &utools.ExecutionResult{
            Success: false,
            Message: "工具未启用",
        }, fmt.Errorf("工具未启用")
    }
    
    // 使用 cfg.YourTool.APIKey 等
}
```

---

## 🎨 前端开发

### 1. 配置工具图标和元数据

在 `web/src/config/tools.ts` 中添加：

```typescript
import { YourIcon } from 'lucide-vue-next'

export const TOOL_CONFIGS: Record<string, ToolConfig> = {
  your_tool: {
    code: 'your_tool',
    title: '你的工具',
    description: '功能描述',
    icon: YourIcon,
    iconBg: 'bg-gradient-to-br from-blue-500 to-indigo-600',
    tags: ['tag1', 'tag2'],
  },
}
```

**常用图标：** `Globe`(网络)、`Mail`(邮件)、`Activity`(监控)、`Image`(图片) - [完整列表](https://lucide.dev/icons/)

### 2. 创建配置组件

在 `web/src/components/tools/` 创建 `YourToolConfig/` 文件夹并在其中创建 `index.vue`：

**目录结构：**
```
web/src/components/tools/
  YourToolConfig/
    index.vue          # 主配置组件（必需）
    types.ts           # 类型定义（可选）
    composables/       # 复杂逻辑拆分（可选）
      useYourToolState.ts
      useYourToolActions.ts
    components/        # 子组件（可选）
```

**`YourToolConfig/index.vue`：**

```vue
<template>
  <div class="space-y-4">
    <h3 class="text-sm font-semibold text-text-primary mb-3">你的工具配置</h3>

    <!-- 使用全局组件，不要用原生 input/select -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">
        参数名 <span class="text-error">*</span>
      </label>
      <BaseInput
        v-model="localConfig.param"
        placeholder="请输入参数"
      />
    </div>

    <!-- 下拉选择使用 BaseSelect -->
    <div>
      <label class="block text-xs font-medium text-text-secondary mb-1.5">选项</label>
      <BaseSelect
        v-model="localConfig.option"
        :options="[
          { label: '选项 1', value: 'option1' },
          { label: '选项 2', value: 'option2' },
        ]"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
// ⚠️ 重要：不要带 .vue 后缀！（项目使用文件夹结构）
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'

interface Props {
  config: Record<string, any>
}

const props = defineProps<Props>()
const emit = defineEmits<{
  (e: 'update:config', value: Record<string, any>): void
}>()

const localConfig = ref({
  param: props.config.param || '',
  option: props.config.option || 'option1',
})

watch(localConfig, (newConfig) => {
    emit('update:config', newConfig)
}, { deep: true })
</script>
```

**⚠️ 关键点：**
- ✅ 使用 `BaseInput`、`BaseSelect` 等全局组件，不要用原生 `<input>`、`<select>`
- ✅ 导入时**不带 `.vue` 后缀**：`import BaseInput from '@/components/BaseInput'`
- ✅ `BaseSelect` 的 `options` 格式：`[{ label: string, value: any }]`

### 3. 注册配置组件（⚠️ 重要：两个位置）

#### 任务编辑器

`web/src/pages/Tasks/components/ToolConfigDrawer/index.vue`：

**1. 添加导入（在 `<script setup>` 中）：**
```typescript
import YourToolConfig from '@/components/tools/YourToolConfig/index.vue'
```

**2. 添加组件使用：**
```vue
<YourToolConfig 
  v-else-if="props.toolCode === 'your_tool'"
  :config="localConfig"
  @update:config="syncConfig"
/>
```

#### 工作流编辑器

⚠️ **注意：这里有两个文件都需要注册**

**A. `web/src/pages/Workflows/components/NodeConfigDrawer.vue`**

**1. 添加导入：**
```typescript
import YourToolConfig from '@/components/tools/YourToolConfig/index.vue'
```

**2. 添加组件使用：**
```vue
<YourToolConfig
  v-else-if="selectedNode?.toolCode === 'your_tool'"
  :config="localNode.config"
  :previous-nodes="props.previousNodes"
  :env-vars="formattedEnvVars"
  @update:config="handleConfigUpdate"
/>
```

**B. `web/src/pages/Workflows/editor.vue`（旧版编辑器）**

在 `NodeConfigDrawer` 内嵌部分添加：
```vue
<YourToolConfig
  v-else-if="selectedNode?.toolCode === 'your_tool'"
  v-model:config="selectedNode.config"
/>
```

**关键点：**
- ⚠️ **必须使用完整路径** `@/components/tools/YourToolConfig/index.vue`（Vite 要求）
- 使用 `toolCode`（不是 `tool_code`）
- `NodeConfigDrawer.vue` 使用 `@update:config`
- `editor.vue` 使用 `v-model:config`

---

## 🧩 前端配置组件最佳实践（避免递归与编译错误）

在实现工具的配置组件时，容易踩到两个坑：

1) 递归更新（Maximum recursive updates）

症状：输入框每次输入就报错 `Maximum recursive updates exceeded in component <BaseTransition>`。

根因：子组件 `emit('update:config', ...)` → 父组件回填 `props.config` → 子组件 `watch(props.config)` 又立即覆盖本地 `localConfig` → 再触发 `emit`，形成闭环。

解决方案（推荐模版）：

```vue
<script setup lang="ts">
import { ref, watch } from 'vue'

interface Props { config: Record<string, any> }
const props = defineProps<Props>()
const emit = defineEmits<{ 'update:config': [Record<string, any>] }>()

// 本地副本
const localConfig = ref({
  // ...用 props.config 初始化
  url: props.config.url || '',
  timeout: props.config.timeout ?? 60,
})

// 防抖旗标防止递归
const updatingFromProps = ref(false)

// 子改父
watch(localConfig, (v) => {
  if (!updatingFromProps.value) emit('update:config', { ...v })
}, { deep: true })

// 父改子
watch(() => props.config, (cfg) => {
  updatingFromProps.value = true
  localConfig.value = {
    url: cfg?.url || '',
    timeout: cfg?.timeout ?? 60,
  }
  setTimeout(() => { updatingFromProps.value = false }, 0)
}, { deep: true })
</script>
```

2) Mustache 示例被当成表达式解析

症状：Vite 报错 `Error parsing JavaScript expression: Unterminated string constant.`

根因：在模板中直接写 `{{external.xxx}}` / `{{nodes.xxx.yyy}}` 作为说明文本，被 Vue 编译器当作表达式。

解决方案：为展示用 Mustache 包上 `v-pre` 或使用转义写法：

```vue
<p class="text-xs" >
  支持变量：<code v-pre>{{external.xxx}}</code> / <code v-pre>{{nodes.xxx.yyy}}</code>
</p>
```

> 小结：配置组件应始终维护“本地副本 + 更新守卫”的模式，并用 v-pre 展示变量占位，避免编译器解析。

---

## 🧾 Headers 数据结构规范与前后端兼容

- 推荐前端在 UI 中用数组形态维护 Headers：`[{ key: string, value: string }]`，便于增删改。
- 发送到后端之前，可以转换成对象：

```ts
const headersObj = Object.fromEntries(localHeaders.map(h => [h.key, h.value]))
```

- 后端工具实现建议兼容两种形态（对象或数组），示例（Go）：

```go
// 优先对象
if hdrs, ok := config["headers"].(map[string]interface{}); ok {
    for k, v := range hdrs { if s, ok := v.(string); ok { req.Header.Set(k, s) } }
} else if arr, ok := config["headers"].([]interface{}); ok {
    for _, item := range arr {
        if m, ok := item.(map[string]interface{}); ok {
            k, _ := m["key"].(string)
            v, _ := m["value"].(string)
            if strings.TrimSpace(k) != "" { req.Header.Set(k, v) }
        }
    }
}
```

---

## 📦 文件类工具规范（输出对象与临时文件）

为与平台的"文件传递链路（下载 → 上传到图床/OSS/COS → 渲染）"契合，文件类工具输出应包含：

- `response`（object）：完整响应或元信息（用于调试）
- `file`（object）：标准文件对象
  - `path`: 本地绝对路径（供后续工具读取）
  - `filename`: 文件名
  - `size`: 字节大小（int64）
  - `mime_type`: MIME 类型

**标准文件对象示例：**
```go
fileObject := map[string]interface{}{
    "path":      "/tmp/autoforge-xxx/file.png",
    "filename":  "file.png",
    "size":      int64(12345),
    "mime_type": "image/png",
}
```

**临时文件管理：**
- 工具应将临时文件保存到 `/tmp/autoforge-<tool-name>/` 或 `/tmp/workflow-files/<execution_id>/`
- 引擎在执行完成后会自动清理 `/tmp/workflow-files/<execution_id>/` 目录
- 其他临时目录需要工具自行清理或依赖系统定时清理

---

## 🔀 多模式输出工具最佳实践（如二维码生成）

当工具支持多种输出模式时（如 Base64 字符串 vs 文件对象），应遵循以下规范：

### 1. 配置参数设计
```go
Properties: map[string]utools.PropertySchema{
    "output_format": {
        Type:        "string",
        Title:       "输出格式",
        Description: "base64: 返回字符串; file: 返回文件对象",
        Default:     "base64",
        Enum:        []interface{}{"base64", "file"},
    },
}
```

### 2. OutputFieldsSchema 设计
```go
OutputFieldsSchema: map[string]utools.OutputFieldDef{
    "response": {
        Type:  "object",
        Label: "完整响应",
        Children: map[string]utools.OutputFieldDef{
            "data": {
                Type:  "string",
                Label: "Base64 数据（仅 base64 模式）",
            },
            "file": {
                Type:  "object",
                Label: "文件对象（仅 file 模式）",
                Children: map[string]utools.OutputFieldDef{
                    "path":      {Type: "string", Label: "文件路径"},
                    "filename":  {Type: "string", Label: "文件名"},
                    "size":      {Type: "integer", Label: "文件大小"},
                    "mime_type": {Type: "string", Label: "MIME 类型"},
                },
            },
        },
    },
    "data": {Type: "string", Label: "Base64 数据（快捷访问，仅 base64 模式）"},
    "file": {
        Type:  "object",
        Label: "文件对象（快捷访问，仅 file 模式）",
        Children: map[string]utools.OutputFieldDef{
            "path": {Type: "string", Label: "文件路径"},
        },
    },
}
```

### 3. Execute 实现
```go
var output map[string]interface{}

if outputFormat == "file" {
    // 生成临时文件
    tempDir := filepath.Join(os.TempDir(), "autoforge-yourtool")
    os.MkdirAll(tempDir, 0755)
    filePath := filepath.Join(tempDir, "file.ext")
    os.WriteFile(filePath, data, 0644)
    
    fileObject := map[string]interface{}{
        "path":      filePath,
        "filename":  "file.ext",
        "size":      int64(len(data)),
        "mime_type": "image/png",
    }
    
    output = map[string]interface{}{
        "response": map[string]interface{}{
            "file": fileObject,
            // 其他字段...
        },
        "file": fileObject,  // 快捷访问
    }
} else {
    // Base64 模式
    base64Data := base64.StdEncoding.EncodeToString(data)
    
    output = map[string]interface{}{
        "response": map[string]interface{}{
            "data": base64Data,
            // 其他字段...
        },
        "data": base64Data,  // 快捷访问
    }
}
```

### 4. 前端使用说明
**⚠️ 重要：不同模式输出的字段不同！**

**Base64 模式：**
- ✅ 可访问：`{{nodes.xxx.data}}` - Base64 字符串
- ❌ 不可访问：`{{nodes.xxx.file}}` - 此字段不存在

**File 模式：**
- ✅ 可访问：`{{nodes.xxx.file}}` - 文件对象（可传递给上传工具）
- ✅ 可访问：`{{nodes.xxx.file.path}}` - 文件路径
- ❌ 不可访问：`{{nodes.xxx.data}}` - 此字段不存在

**配置说明中应明确标注：**
```typescript
usageItems: [
  { text: 'Base64 模式：输出字符串，可用于直接显示' },
  { text: 'File 模式：输出文件对象，可传递给上传工具' },
  { text: '注意：两种模式输出字段不同，请根据后续节点选择' },
]
```

---

## 🧭 工作流编辑器接入 Checklist

1) 工具元数据（必填）：`web/src/config/tools.ts`
   - code/title/description/icon/tags/usageItems 等
2) 配置组件（必填）：`web/src/components/tools/YourToolConfig/index.vue`
   - 采用“本地副本 + 更新守卫”模式，避免递归更新
   - 变量说明用 `v-pre` 展示
3) 侧边配置面板注册（必填）：`web/src/pages/Workflows/components/NodeConfigDrawer.vue`
   - `import YourToolConfig ...`
   - 在模板里：`v-else-if="node.toolCode === 'your_tool'"`
4)（如适用）任务编辑器注册：`web/src/pages/Tasks/components/ToolConfigDrawer.vue`
5) 工具前端列表与图标渲染可用 `getToolIcon/getToolIconBg`（已由工具元数据驱动）

---

## 🧰 后端接入 Checklist（复盘）

1) 工具类实现：`pkg/utools/<code>/<code>_tool.go`
   - `ToolMetadata` + `ConfigSchema` + `Execute`
   - `OutputFieldsSchema` 必须包含 `response`（object）
2) 工具注册：在 `cmd/main.go` 以空导入 `_ "auto-forge/pkg/utools/<code>"`
3)（如需）后端全局配置：更新 `pkg/config/config.go` 与 `config.yaml`
4) 工具配置（敏感信息）统一用工具配置中心：`internal/services/tool_config`

---

## 🧪 调试与排错（新增）

- 工具未显示在列表：
  - 检查 `cmd/main.go` 是否空导入了你的工具；重启后端。
  - 调用管理员 API 同步工具定义（或查看启动日志）。
- 配置组件不显示/报错：
  - 确认 `NodeConfigDrawer.vue` 已注册 `v-else-if="node.toolCode === 'xxx'"`
  - 避免在模板中直接写 `{{external.xxx}}`，改用 `v-pre`
  - 若出现递归更新报错，按“本地副本 + 更新守卫”模式改造
- 变量不生效：
  - 使用 `VariableSelector` 并传入 `previousNodes` 与 `envVars`
  - 检查执行详情中的 `resolved_config`
- 文件未清理：
  - 确认工具写入了 `/tmp/workflow-files/<execution_id>/`；执行结束后引擎会清理该目录

---

## 🌰 参考：文件下载器（file_downloader）最小表单

```vue
<FileDownloaderConfig
  v-model:config="localNode.config"
  :previous-nodes="props.previousNodes"
  :env-vars="props.envVars"
/>
```

后端输出（关键字段）：

```json
{
  "response": { "url": "...", "status_code": 200, "headers": {"Content-Type": "..."} },
  "file": { "type": "file", "path": "/tmp/workflow-files/<execID>/...", "filename": "...", "size": 123, "mime_type": "..." }
}
```

> 以上规范已在近期开发中验证可用，按此模板开发可避免常见前后端联调问题。

## 📐 开发规范

### 后端规范

#### 必须遵守

1. **OutputFieldsSchema 必须包含 `response` 字段**
   ```go
   OutputFieldsSchema: map[string]utools.OutputFieldDef{
       "response": {
           Type: "object",      // ✅ 正确
           Label: "完整响应",    // ✅ 使用 Label，不是 Description
           Children: map[string]utools.OutputFieldDef{  // ✅ 使用 Children，不是 Properties
               "field": {Type: "string", Label: "字段"},
           },
       },
       "field": {Type: "string", Label: "快捷访问"},  // 可选：快捷访问字段
   }
   ```
   
   **⚠️ 常见错误：**
   - ❌ 使用 `Properties` → ✅ 应使用 `Children`
   - ❌ 使用 `Description` → ✅ 应使用 `Label`

2. **ConfigSchema 数值范围必须使用指针**
   ```go
   Properties: map[string]utools.PropertySchema{
       "size": {
           Type:    "integer",
           Title:   "尺寸",
           Default: 256,
           Minimum: float64Ptr(64),   // ✅ 必须使用指针
           Maximum: float64Ptr(2048), // ✅ 必须使用指针
       },
   }
   
   // 辅助函数
   func float64Ptr(v float64) *float64 {
       return &v
   }
   ```

3. **Enum 必须是 []interface{} 类型**
   ```go
   Properties: map[string]utools.PropertySchema{
       "level": {
           Type:  "string",
           Title: "级别",
           Enum:  []interface{}{"Low", "Medium", "High"},  // ✅ []interface{}
           // ❌ 错误：Enum: []string{...}
       },
   }
   ```

4. **Required 字段在 ConfigSchema 中，不在 PropertySchema 中**
   ```go
   schema := &utools.ConfigSchema{
       Type: "object",
       Properties: map[string]utools.PropertySchema{
           "param": {
               Type:  "string",
               Title: "参数",
               // ❌ 错误：Required: true  // PropertySchema 没有 Required 字段
           },
       },
       Required: []string{"param"},  // ✅ 正确：在 ConfigSchema 层级
   }
   ```

5. **敏感信息必须标记**
   ```go
   Properties: map[string]utools.PropertySchema{
       "api_key": {Type: "string", Title: "API Key", Secret: true},
   }
   ```

6. **返回结果必须包含 response**
   ```go
   Output: map[string]interface{}{
       "response": fullResponse,  // 完整对象
       "field":    quickAccess,   // 快捷字段
   }
   ```

4. **文件参数处理（重要！）**
   
   当工具接收文件参数时，必须按以下顺序解析：

```go
   var filePath string
   
   // 1. 优先检查文件对象（从外部API/工作流传入）
   if fileObj, ok := toolConfig["file"].(map[string]interface{}); ok {
       if path, ok := fileObj["path"].(string); ok && path != "" {
           filePath = path
       }
   }
   
   // 2. 如果不是文件对象,再尝试字符串路径
   if filePath == "" {
       if strPath, ok := toolConfig["file"].(string); ok && strPath != "" {
           filePath = strPath
       }
   }
   
   // 3. 最终验证
   if filePath == "" {
    return &utools.ExecutionResult{
           Success: false,
           Message: "文件参数无效",
       }, fmt.Errorf("文件参数无效")
   }
   ```
   
   **为什么这样做？**
   - 外部API触发器传入的是文件对象：`{"path": "/tmp/...", "filename": "...", "size": 123}`
   - 用户手动输入的可能是字符串路径：`"/path/to/file"`
   - 必须先检查对象，否则会误报"参数无效"

### 前端规范

#### 组件接口标准

```typescript
// ✅ 正确
const localConfig = ref({...props.config})
watch(localConfig, (v) => emit('update:config', v), {deep: true})

// ❌ 错误
props.config.param = value  // 禁止直接修改 props
```

#### 样式规范（支持主题切换）

```vue
<!-- ✅ 使用语义化变量 -->
<div class="text-text-primary bg-bg-primary border-border-primary">

<!-- ❌ 禁止硬编码颜色 -->
<div class="text-gray-900 bg-white border-gray-300">
```

**关键变量：**
- `text-text-primary/secondary/tertiary` - 文本颜色
- `bg-bg-primary/elevated/hover` - 背景色
- `border-border-primary` - 边框色
- `text-error` / `text-primary` - 状态色

#### 表单元素标准

```vue
<!-- 标准输入框 -->
      <input
  v-model="localConfig.param"
  class="w-full px-3 py-2 text-sm 
         bg-bg-primary text-text-primary
         border border-border-primary rounded-md
         focus:ring-2 focus:ring-primary"
/>

<!-- 敏感信息 -->
<input type="password" />

<!-- 下拉选择 -->
<select class="...同上">
  <option value="">请选择</option>
      </select>
```

---

## 📋 检查清单

### 后端检查

- [ ] 工具实现了 Tool 接口
- [ ] **OutputFieldsSchema 已定义且包含 `response` 字段**
- [ ] ConfigSchema 定义完整
- [ ] 敏感信息标记为 Secret: true
- [ ] **在 `cmd/main.go` 中导入（不是 `init.go`）**
- [ ] Execute 方法有错误处理
- [ ] **如果接收文件参数，先检查对象再检查字符串**
- [ ] 如果需要后端配置：
  - [ ] 在 `pkg/config/config.go` 中定义结构体
  - [ ] 在 `config.yaml` 和 `config.example.yaml` 中添加配置
  - [ ] Execute 中正确读取配置并验证

### 前端检查

- [ ] 在 `web/src/config/tools.ts` 配置图标和元数据
- [ ] 配置组件遵循标准接口（Props/Emits）
- [ ] 使用本地状态 + watch 模式
- [ ] 使用语义化 CSS 变量
- [ ] **在两个位置都注册了组件**（⚠️ 最容易忘记）
  - [ ] 任务编辑器: `pages/Tasks/components/ToolConfigDrawer/index.vue`
  - [ ] **工作流编辑器: `pages/Workflows/editor.vue`**（使用 `toolCode`）
- [ ] 必填字段有 `*` 标记
- [ ] 敏感字段使用 `type="password"`
- [ ] 代码通过 ESLint 和 TypeScript 检查
- [ ] 在亮色和暗色主题下都测试过

### 测试检查

- [ ] 后端启动成功，日志无报错
- [ ] 工具出现在工具列表中
- [ ] 拖入画布后配置面板正常显示
- [ ] 配置项双向绑定正常工作
- [ ] 执行工作流成功，输出符合预期
- [ ] 变量引用正常工作（如 `{{nodes.xxx.response.field}}`）
- [ ] 如果涉及文件上传，测试外部API触发器场景

---

## 🧪 测试与调试

### 重启后端服务

```bash
# 方法1: 杀掉旧进程并重启
lsof -ti:7777 | xargs kill -9 && sleep 1 && nohup go run cmd/main.go > /tmp/cooper-backend.log 2>&1 &

# 方法2: 使用 pkill
pkill -9 -f "cmd/main.go" && sleep 1 && nohup go run cmd/main.go > /tmp/cooper-backend.log 2>&1 &
```

### 查看后端日志

```bash
# 查看启动日志
tail -20 /tmp/cooper-backend.log | grep "INFO.*服务启动成功"

# 实时监控
tail -f /tmp/cooper-backend.log

# 查看错误
tail -100 /tmp/cooper-backend.log | grep -i error
```

### 调试步骤

1. **后端注册检查**
   ```bash
   # 启动后应该看到工具注册日志
   tail /tmp/cooper-backend.log | grep "your_tool"
   ```

2. **前端工具列表检查**
   - 打开浏览器控制台
   - 进入工作流编辑页面
   - 查看 Network 中的 `/api/v1/tools` 响应
   - 确认你的工具在列表中

3. **配置组件检查**
   - 在工作流中拖入工具
   - 点击节点打开右侧配置面板
   - 如果不显示：
     * 打开控制台查看是否有 Vue 报错
     * 检查 `selectedNode.toolCode` 的值
     * 确认条件 `v-else-if` 中的 `toolCode` 是否匹配

4. **执行失败排查**
   - 查看执行详情中的 `resolved_config` 字段
   - 检查变量是否正确替换
   - 查看后端日志中的错误堆栈

---

## ❓ 常见问题

**Q: 工具未出现在列表中？**
- ✅ 检查是否在 `cmd/main.go` 中导入（不是 `init.go`）
- ✅ 重启后端：`lsof -ti:7777 | xargs kill -9 && go run cmd/main.go`
- ✅ 查看日志确认工具已注册

**Q: 配置组件不显示？**
- ✅ 确认在 `editor.vue` 中注册（使用 `toolCode` 不是 `tool_code`）
- ✅ 确认在 `ToolConfigDrawer/index.vue` 中注册
- ✅ 检查浏览器控制台是否有 Vue 报错
- ✅ 打印 `selectedNode.toolCode` 确认值是否匹配

**Q: 文件上传失败，提示"文件参数无效"？**
- ✅ 检查是否先解析文件对象再解析字符串（见"文件参数处理"章节）
- ✅ 查看执行详情中的 `resolved_config.file` 是对象还是字符串
- ✅ 确认文件路径存在且可读

**Q: 外部服务返回 403/401？**
- ✅ 检查 `config.yaml` 中的 API Key 是否正确
- ✅ 检查 Endpoint/Region 是否匹配（如 OSS）
- ✅ 查看后端日志中的完整错误响应
- ✅ 使用 Postman 测试该 API 是否正常

**Q: 变量引用不工作？**
- ✅ 确认后端 OutputFieldsSchema 中定义了该字段
- ✅ 使用 `{{nodes.xxx.response.field}}` 访问嵌套字段
- ✅ 检查执行详情中前置节点的 `output` 字段

---

**文档版本**: v2.3  
**最后更新**: 2025-01-17  
**维护者**: Cooper Team

## 📝 更新日志

### v2.3 (2025-01-17) - 📁 组件结构规范化
**统一组件组织结构，提升可维护性和扩展性**

- 🏗️ **工具配置组件统一为文件夹结构**：所有工具配置从单文件 `.vue` 改为 `ToolName/index.vue` 结构
- 📂 **标准目录结构**：每个工具一个文件夹，支持 `types.ts`、`composables/`、`components/` 子目录
- 🔄 **更新导入路径**：**必须使用完整路径** `@/components/tools/XxxConfig/index.vue`（Vite 要求）
- ✨ **提升可扩展性**：便于后续拆分复杂组件、添加类型定义和 composables

**目录结构示例：**
```
web/src/components/tools/
  AliyunOSSConfig/
    index.vue
  FeishuBotConfig/
    index.vue
    types.ts
    composables/
      useFeishuState.ts
      useFeishuActions.ts
```

**导入示例：**
```typescript
// ✅ 正确：使用完整路径
import AliyunOSSConfig from '@/components/tools/AliyunOSSConfig/index.vue'

// ❌ 错误：Vite 无法解析
import AliyunOSSConfig from '@/components/tools/AliyunOSSConfig'
```

### v2.2 (2025-01-17) - 🔥 重要更新
**基于阿里云 OSS 工具开发的实战经验优化**

- ⚠️ **修正致命错误**：工具注册位置应为 `cmd/main.go`（不是 `init.go`）
- ⚠️ **修正前端注册路径**：工作流编辑器应在 `editor.vue` 中注册（使用 `toolCode`）
- 🆕 **新增后端配置管理章节**：详细说明如何管理工具的后端配置
- 🆕 **新增文件参数处理规范**：必须先检查对象再检查字符串，避免"参数无效"错误
- 🆕 **新增测试与调试章节**：完整的重启命令、日志查看、调试步骤
- 🆕 **新增常见问题排查**：覆盖配置不显示、文件上传失败、403/401 错误等实际问题
- ✨ **扩充检查清单**：新增后端配置、文件参数、测试验证等检查项

**关键修复点：**
1. 工具注册必须在 `cmd/main.go` 导入
2. 工作流编辑器注册条件使用 `selectedNode?.toolCode`（不是 `tool_code`）
3. 文件参数解析必须先检查对象（`map[string]interface{}`）再检查字符串
4. 后端配置需要同步更新 `config.go` 和 `config.yaml`
5. OSS 等云服务需确保 Endpoint 与 Region 匹配

### v2.1 (2025-01-17)
- 🎯 大幅精简文档，聚焦核心要点
- 删除冗长示例代码，保留最精简模板
- 移除测试相关内容
- 保留最关键的规范和检查清单

### v2.0 (2025-01-17)
- 新增完整的前端配置组件开发规范
- 明确双编辑器注册要求
- 添加完整检查清单
