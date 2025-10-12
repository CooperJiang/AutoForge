# 工作流系统前端设计文档

## 📋 目录
- [系统概述](#系统概述)
- [技术栈](#技术栈)
- [核心功能](#核心功能)
- [数据结构](#数据结构)
- [节点类型](#节点类型)
- [API设计建议](#api设计建议)
- [前后端交互流程](#前后端交互流程)
- [文件结构](#文件结构)
- [使用示例](#使用示例)

---

## 系统概述

工作流系统是一个可视化的任务编排工具，允许用户通过拖拽节点和连接的方式创建自动化工作流。类似于 n8n、Apache Airflow 等工作流引擎。

### 主要特性
- 🎨 可视化拖拽编辑器
- 🔀 支持条件分支和多路分支
- ⏰ 定时触发器 + Webhook触发器
- 🛠️ 集成现有工具（HTTP请求、邮件发送、健康检查）
- 🔄 数据流和变量系统（节点输出引用、环境变量）
- 🔁 错误重试机制（支持指数退避）
- 🧪 节点测试运行（调试功能）
- 📊 执行历史和日志查看
- 💾 导出/导入 JSON 配置

---

## 技术栈

- **框架**: Vue 3 + TypeScript
- **工作流引擎**: @vue-flow/core v1.47.0
- **路由**: Vue Router
- **状态管理**: Composables (useWorkflow)
- **UI组件**: 自定义组件库
- **图标**: lucide-vue-next

---

## 核心功能

### 1. 工作流列表页 (`/workflows`)
- 展示所有工作流
- 创建新工作流
- 编辑/删除工作流

### 2. 工作流编辑器 (`/workflows/create`, `/workflows/:id/edit`)
- 左侧工具面板：拖拽添加节点
- 中间画布：可视化编辑工作流
- 节点配置抽屉：配置节点参数
- 顶部工具栏：保存、测试、导出

### 3. 节点配置
- 每种节点类型有独立的配置组件
- 支持实时验证
- HTTP工具支持 cURL 粘贴自动解析

---

## 数据结构

### WorkflowNode（工作流节点）

```typescript
export type NodeType = 'trigger' | 'tool' | 'condition' | 'delay' | 'switch' | 'end'

export interface NodeRetryConfig {
  enabled: boolean              // 是否启用重试
  maxRetries: number            // 最大重试次数（1-10）
  retryInterval: number         // 重试间隔（秒）
  exponentialBackoff: boolean   // 是否使用指数退避
}

export interface WorkflowNode {
  id: string                    // 节点唯一ID
  type: NodeType                // 节点类型
  toolCode?: string             // 工具代码（type为tool时必填）
  name: string                  // 节点名称
  config: Record<string, any>   // 节点配置
  retry?: NodeRetryConfig       // 错误重试配置（工具节点可用）
  position: { x: number; y: number }  // 画布位置
}
```

### WorkflowEdge（连接线）

```typescript
export interface WorkflowEdge {
  id: string                    // 连接唯一ID
  source: string                // 源节点ID
  target: string                // 目标节点ID
  sourceHandle?: string         // 源节点输出点ID（用于分支）
  targetHandle?: string         // 目标节点输入点ID
  condition?: string            // 条件标签（可选）
}
```

### Workflow（完整工作流）

```typescript
export interface WorkflowEnvVar {
  key: string                   // 变量名
  value: string                 // 变量值
  description?: string          // 变量描述
  encrypted?: boolean           // 是否加密存储
}

export interface Workflow {
  id?: string                   // 工作流ID（后端生成）
  name: string                  // 工作流名称
  description: string           // 工作流描述
  trigger: WorkflowTrigger      // 触发器配置
  nodes: WorkflowNode[]         // 节点列表
  edges: WorkflowEdge[]         // 连接列表
  envVars?: WorkflowEnvVar[]    // 环境变量
  enabled?: boolean             // 是否启用
  created_at?: string           // 创建时间
  updated_at?: string           // 更新时间
}
```

### WorkflowTrigger（触发器）

```typescript
export interface WorkflowTrigger {
  type: 'schedule' | 'manual' | 'webhook'  // 触发类型
  scheduleType?: string         // 调度类型（daily/weekly/monthly/interval/cron）
  scheduleValue?: string        // 调度值
  webhookPath?: string          // Webhook路径
  webhookMethod?: string        // Webhook请求方法（POST/GET/PUT）
}
```

---

## 节点类型

### 1. 触发器节点 (Trigger)

**类型**: `trigger`
**颜色**: 蓝色/紫色
**配置**:
```typescript
{
  scheduleType: 'daily' | 'weekly' | 'monthly' | 'hourly' | 'interval' | 'cron',
  scheduleValue: string  // 时间值，根据scheduleType不同格式不同
}
```

**示例配置**:
- 每日: `{ scheduleType: 'daily', scheduleValue: '09:00' }`
- 间隔: `{ scheduleType: 'interval', scheduleValue: '300' }` (300秒)
- Cron: `{ scheduleType: 'cron', scheduleValue: '0 0 * * * *' }`

---

### 2. 工具节点 (Tool)

**类型**: `tool`
**工具代码**: `http_request` | `email_sender` | `health_checker`

#### 2.1 HTTP请求 (`http_request`)

**颜色**: 蓝色/紫色
**配置**:
```typescript
{
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH',
  url: string,
  headers: Array<{ key: string, value: string }>,
  params: Array<{ key: string, value: string }>,
  body: string  // JSON字符串
}
```

**特性**: 支持 Cmd/Ctrl+V 粘贴 cURL 命令自动解析

---

#### 2.2 邮件发送 (`email_sender`)

**颜色**: 紫色/粉色
**配置**:
```typescript
{
  to: string,              // 收件人，逗号分隔
  cc?: string,             // 抄送人
  subject: string,         // 邮件主题
  body: string,            // 邮件正文
  content_type: 'text/plain' | 'text/html'
}
```

---

#### 2.3 健康检查 (`health_checker`)

**颜色**: 靛蓝色/蓝色
**配置**:
```typescript
{
  url: string,
  method: 'GET' | 'POST' | 'PUT' | 'DELETE',
  headers: Array<{ key: string, value: string }>,
  body: string,
  timeout: number,                    // 超时时间（秒）
  expected_status: number,            // 期望状态码
  response_pattern?: string,          // 响应内容正则匹配
  ssl_expiry_days?: number           // SSL证书到期警告天数
}
```

**特性**: 支持 Cmd/Ctrl+V 粘贴 cURL 命令自动解析

---

### 3. 条件判断节点 (Condition)

**类型**: `condition`
**颜色**: 黄色/橙色
**输出分支**: 2个（True、False）

**配置**:
```typescript
{
  conditionType: 'simple' | 'expression' | 'script',

  // simple 模式
  field?: string,           // 检查字段
  operator?: string,        // 操作符：equals, not_equals, greater_than, less_than, contains, etc.
  value?: string,           // 期望值

  // expression 模式
  expression?: string,      // 表达式：status == 200 && success == true

  // script 模式
  script?: string          // JavaScript代码，可访问input变量，返回boolean
}
```

**操作符列表**:
- `equals`: 等于 (==)
- `not_equals`: 不等于 (!=)
- `greater_than`: 大于 (>)
- `less_than`: 小于 (<)
- `greater_or_equal`: 大于等于 (>=)
- `less_or_equal`: 小于等于 (<=)
- `contains`: 包含
- `not_contains`: 不包含
- `starts_with`: 以...开始
- `ends_with`: 以...结束
- `is_empty`: 为空
- `is_not_empty`: 不为空

**分支连接**:
- True分支: `sourceHandle: 'true'`
- False分支: `sourceHandle: 'false'`

---

### 4. 开关分支节点 (Switch)

**类型**: `switch`
**颜色**: 靛蓝色
**输出分支**: N+1个（N个Case + 1个Default）

**配置**:
```typescript
{
  field: string,                           // 检查字段
  cases: Array<{
    label: string,                         // 分支标签
    value: string                          // 匹配值
  }>
}
```

**示例配置**:
```json
{
  "field": "status",
  "cases": [
    { "label": "Success", "value": "200" },
    { "label": "Not Found", "value": "404" },
    { "label": "Server Error", "value": "500" }
  ]
}
```

**分支连接**:
- Case分支: `sourceHandle: 'case_0'`, `case_1`, etc.
- Default分支: `sourceHandle: 'default'`

---

### 5. 延迟等待节点 (Delay)

**类型**: `delay`
**颜色**: 紫色

**配置**:
```typescript
{
  duration: number,                 // 等待时长
  unit: 'seconds' | 'minutes' | 'hours'
}
```

**示例**:
- 等待5秒: `{ duration: 5, unit: 'seconds' }`
- 等待10分钟: `{ duration: 10, unit: 'minutes' }`

---

## API设计建议

### 基础接口

#### 1. 获取工作流列表
```
GET /api/v1/workflows
```

**响应**:
```json
{
  "code": 200,
  "data": [
    {
      "id": "wf_123",
      "name": "每日健康检查",
      "description": "检查服务器健康状态并发送报告",
      "enabled": true,
      "created_at": "2025-01-01T00:00:00Z",
      "updated_at": "2025-01-01T00:00:00Z",
      "node_count": 5,
      "last_run_at": "2025-01-01T09:00:00Z",
      "last_run_status": "success"
    }
  ]
}
```

---

#### 2. 获取工作流详情
```
GET /api/v1/workflows/:id
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "id": "wf_123",
    "name": "每日健康检查",
    "description": "检查服务器健康状态并发送报告",
    "enabled": true,
    "trigger": {
      "type": "schedule",
      "scheduleType": "daily",
      "scheduleValue": "09:00"
    },
    "nodes": [...],
    "edges": [...],
    "created_at": "2025-01-01T00:00:00Z",
    "updated_at": "2025-01-01T00:00:00Z"
  }
}
```

---

#### 3. 创建工作流
```
POST /api/v1/workflows
```

**请求体**:
```json
{
  "name": "工作流名称",
  "description": "工作流描述",
  "enabled": true,
  "trigger": {...},
  "nodes": [...],
  "edges": [...]
}
```

---

#### 4. 更新工作流
```
PUT /api/v1/workflows/:id
```

**请求体**: 同创建接口

---

#### 5. 删除工作流
```
DELETE /api/v1/workflows/:id
```

---

#### 6. 测试运行工作流
```
POST /api/v1/workflows/:id/test
```

**请求体**:
```json
{
  "input": {}  // 可选的测试输入数据
}
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "execution_id": "exec_456",
    "status": "running",
    "start_time": "2025-01-01T10:00:00Z",
    "nodes_executed": 0,
    "total_nodes": 5
  }
}
```

---

#### 7. 获取执行历史
```
GET /api/v1/workflows/:id/executions
```

**响应**:
```json
{
  "code": 200,
  "data": [
    {
      "execution_id": "exec_456",
      "workflow_id": "wf_123",
      "status": "success",
      "start_time": "2025-01-01T09:00:00Z",
      "end_time": "2025-01-01T09:00:05Z",
      "duration": 5,
      "nodes_executed": 5,
      "error_message": null
    }
  ]
}
```

---

#### 8. 获取执行详情
```
GET /api/v1/workflows/:id/executions/:execution_id
```

**响应**:
```json
{
  "code": 200,
  "data": {
    "execution_id": "exec_456",
    "workflow_id": "wf_123",
    "status": "success",
    "start_time": "2025-01-01T09:00:00Z",
    "end_time": "2025-01-01T09:00:05Z",
    "nodes": [
      {
        "node_id": "node_1",
        "node_name": "健康检查",
        "status": "success",
        "start_time": "2025-01-01T09:00:01Z",
        "end_time": "2025-01-01T09:00:02Z",
        "input": {...},
        "output": {
          "status": 200,
          "healthy": true
        },
        "error": null
      }
    ]
  }
}
```

---

## 前后端交互流程

### 1. 创建工作流流程

```
用户操作 → 前端编辑器 → 生成JSON → 调用API → 后端存储
```

**前端发送的数据结构**:
```json
{
  "name": "每日健康检查",
  "description": "检查服务器并发送报告",
  "enabled": true,
  "trigger": {
    "type": "schedule",
    "scheduleType": "daily",
    "scheduleValue": "09:00"
  },
  "nodes": [
    {
      "id": "node_1",
      "type": "trigger",
      "name": "定时触发",
      "config": {
        "scheduleType": "daily",
        "scheduleValue": "09:00"
      },
      "position": { "x": 250, "y": 100 }
    },
    {
      "id": "node_2",
      "type": "tool",
      "toolCode": "health_checker",
      "name": "健康检查",
      "config": {
        "url": "https://api.example.com/health",
        "method": "GET",
        "timeout": 10,
        "expected_status": 200
      },
      "position": { "x": 250, "y": 250 }
    },
    {
      "id": "node_3",
      "type": "condition",
      "name": "检查结果",
      "config": {
        "conditionType": "simple",
        "field": "healthy",
        "operator": "equals",
        "value": "true"
      },
      "position": { "x": 250, "y": 400 }
    },
    {
      "id": "node_4",
      "type": "tool",
      "toolCode": "email_sender",
      "name": "发送成功通知",
      "config": {
        "to": "admin@example.com",
        "subject": "健康检查成功",
        "body": "服务器运行正常",
        "content_type": "text/plain"
      },
      "position": { "x": 150, "y": 550 }
    },
    {
      "id": "node_5",
      "type": "tool",
      "toolCode": "email_sender",
      "name": "发送告警",
      "config": {
        "to": "admin@example.com",
        "subject": "健康检查失败",
        "body": "服务器异常，请检查",
        "content_type": "text/plain"
      },
      "position": { "x": 350, "y": 550 }
    }
  ],
  "edges": [
    {
      "id": "edge_1",
      "source": "node_1",
      "target": "node_2"
    },
    {
      "id": "edge_2",
      "source": "node_2",
      "target": "node_3"
    },
    {
      "id": "edge_3",
      "source": "node_3",
      "target": "node_4",
      "sourceHandle": "true"
    },
    {
      "id": "edge_4",
      "source": "node_3",
      "target": "node_5",
      "sourceHandle": "false"
    }
  ]
}
```

---

### 2. 工作流执行流程

```
触发器 → 按照edges顺序执行nodes → 记录每个节点的输入/输出 → 返回执行结果
```

**后端执行逻辑**:
1. 根据trigger配置注册定时任务
2. 触发时，创建执行实例（execution_id）
3. 从trigger节点开始，按照edges顺序执行
4. 每个节点执行时：
   - 获取上一个节点的输出作为输入
   - 执行节点逻辑（调用对应工具）
   - 记录输出结果
   - 根据节点类型决定下一步：
     - 普通节点：执行所有连接的下游节点
     - 条件节点：根据条件结果选择true/false分支
     - 开关节点：根据字段值选择对应case分支
     - 延迟节点：等待指定时间后继续
5. 所有节点执行完成后，更新执行状态

---

### 3. 节点数据流转

**节点输入输出规范**:

每个节点执行时：
- **输入**: 上一个节点的输出（`input`）
- **输出**: 当前节点的执行结果（`output`）

**示例**:

```javascript
// HTTP请求节点输出
{
  "status": 200,
  "statusText": "OK",
  "headers": {...},
  "data": {...}
}

// 健康检查节点输出
{
  "healthy": true,
  "status": 200,
  "response_time": 123,
  "ssl_valid": true,
  "ssl_days_remaining": 60
}

// 条件节点输出
{
  "result": true,
  "field": "healthy",
  "value": true
}
```

---

## 文件结构

```
web/
├── src/
│   ├── types/
│   │   └── workflow.ts                    # 工作流类型定义
│   │
│   ├── utils/
│   │   └── variableParser.ts              # 变量解析工具
│   │
│   ├── composables/
│   │   └── useWorkflow.ts                 # 工作流状态管理
│   │
│   ├── pages/
│   │   └── Workflows/
│   │       ├── index.vue                  # 工作流列表页
│   │       ├── editor.vue                 # 工作流编辑器
│   │       ├── executions.vue             # 执行历史列表
│   │       ├── execution-detail.vue       # 执行详情页
│   │       └── components/
│   │           ├── ToolPanel.vue          # 工具面板
│   │           ├── ToolNode.vue           # 工具节点组件
│   │           ├── TriggerNode.vue        # 触发器节点组件
│   │           ├── ConditionNode.vue      # 条件节点组件
│   │           ├── DelayNode.vue          # 延迟节点组件
│   │           ├── SwitchNode.vue         # 开关节点组件
│   │           ├── NodeConfigDrawer.vue   # 节点配置抽屉
│   │           ├── WorkflowCard.vue       # 工作流卡片
│   │           └── EnvVarManager.vue      # 环境变量管理
│   │
│   └── components/
│       ├── VariableSelector.vue           # 变量选择器
│       ├── RetryConfig.vue                # 重试配置组件
│       └── tools/
│           ├── EmailToolConfig.vue        # 邮件工具配置
│           ├── HealthCheckerConfig.vue    # 健康检查配置
│           ├── TriggerConfig.vue          # 触发器配置
│           ├── ConditionConfig.vue        # 条件配置
│           ├── DelayConfig.vue            # 延迟配置
│           └── SwitchConfig.vue           # 开关配置
```

---

## 变量系统详解

### 变量类型

工作流系统支持三种类型的变量引用：

#### 1. 环境变量（Environment Variables）
**语法**: `{{env.VARIABLE_NAME}}`

环境变量在工作流级别定义，可以在所有节点中使用。适用于存储API密钥、配置参数等敏感信息。

**配置位置**: 工作流编辑器 → 环境变量按钮

**示例**:
```json
{
  "key": "API_KEY",
  "value": "sk-1234567890",
  "description": "OpenAI API密钥",
  "encrypted": true
}
```

**使用**:
```
URL: https://api.openai.com/v1/chat
Headers:
  Authorization: Bearer {{env.API_KEY}}
```

#### 2. 节点输出引用（Node Output Reference）
**语法**: `{{node_id.field}}`

引用前置节点的输出数据。系统会自动解析节点输出结构，只能引用已执行节点的输出。

**示例**:
```
前置节点输出:
{
  "status": 200,
  "data": {
    "user_id": 12345,
    "username": "alice"
  }
}

引用方式:
{{http_node_1.status}}          → 200
{{http_node_1.data.user_id}}    → 12345
{{http_node_1.data.username}}   → "alice"
```

#### 3. 触发器数据（Trigger Data）
**语法**: `{{trigger.field}}`

访问触发工作流时传入的数据（主要用于Webhook触发）。

**Webhook触发示例**:
```bash
curl -X POST https://your-domain.com/api/webhook/my-workflow \
  -H "Content-Type: application/json" \
  -d '{"event": "user.created", "user_id": 12345}'
```

**引用**:
```
{{trigger.event}}      → "user.created"
{{trigger.user_id}}    → 12345
{{trigger.timestamp}}  → 触发时间戳
{{trigger.type}}       → 触发类型
```

### 变量选择器

**快捷键**: `Cmd/Ctrl + K` 打开变量选择器

**功能**:
- 🔍 搜索变量
- 📋 分类显示（环境变量、触发器、前置节点）
- 🎯 点击插入到光标位置
- 💡 显示变量描述和示例

---

## 错误重试机制

### 配置选项

工具节点支持自动重试机制，当节点执行失败时自动重试。

#### 1. 基本配置
- **启用重试**: 开关
- **最大重试次数**: 1-10次（建议3-5次）
- **重试间隔**: 秒（等待时间）
- **指数退避**: 是否启用

#### 2. 指数退避策略

启用后，重试间隔会指数增长：

```
第1次重试: 等待 interval 秒
第2次重试: 等待 interval × 2 秒
第3次重试: 等待 interval × 4 秒
第N次重试: 等待 interval × 2^(N-1) 秒
```

**示例配置**:
```json
{
  "enabled": true,
  "maxRetries": 3,
  "retryInterval": 5,
  "exponentialBackoff": true
}
```

**实际重试时间**:
- 第1次重试：等待5秒
- 第2次重试：等待10秒
- 第3次重试：等待20秒

### 适用场景

✅ **适合使用重试**:
- HTTP请求（网络不稳定）
- 第三方API调用
- 资源暂时不可用
- 速率限制（配合退避）

❌ **不适合使用重试**:
- 邮件发送（避免重复发送）
- 数据写入操作（幂等性问题）
- 长时间运行的任务

---

## 使用示例

### 示例1: 每日健康检查工作流

```
[定时触发: 每天09:00]
  ↓
[健康检查: GET https://api.example.com/health]
  ↓
[条件判断: healthy == true?]
  ├─ True  → [发送邮件: 服务正常]
  └─ False → [发送邮件: 告警通知]
```

---

### 示例2: API监控与分级告警

```
[定时触发: 每5分钟]
  ↓
[HTTP请求: GET https://api.example.com/status]
  ↓
[开关分支: status值]
  ├─ Case 200 → [记录日志: 正常]
  ├─ Case 404 → [发送邮件: 资源不存在]
  ├─ Case 500 → [发送邮件: 服务器错误] → [延迟5分钟] → [重试请求]
  └─ Default  → [发送邮件: 未知错误]
```

---

### 示例3: 带重试的任务执行

```
[定时触发: 每小时]
  ↓
[HTTP请求: POST https://api.example.com/task]
  ↓
[条件判断: status >= 200 && status < 300?]
  ├─ True  → [发送邮件: 任务成功]
  └─ False → [延迟30秒] → [HTTP请求: 重试] → [条件判断]
```

---

## 前端已完成功能

### 核心编辑器
✅ 工作流可视化编辑器
✅ 拖拽添加节点
✅ 节点连接
✅ 节点配置抽屉
✅ 条件分支（IF）
✅ 多路分支（Switch）
✅ 延迟等待
✅ 分支可视化（True/False标签）
✅ 导出JSON

### 触发器
✅ 定时触发器（Daily/Weekly/Monthly/Hourly/Interval/Cron）
✅ 手动触发器
✅ Webhook触发器（支持POST/GET/PUT）

### 工具节点
✅ HTTP请求工具（支持cURL解析）
✅ 邮件发送工具
✅ 健康检查工具（支持cURL解析）

### 数据流和变量
✅ 变量解析引擎（支持 `{{variable}}` 语法）
✅ 环境变量管理（支持加密存储）
✅ 节点输出引用（`{{node_id.field}}`）
✅ 触发器数据引用（`{{trigger.data}}`）
✅ 可视化变量选择器（支持搜索和快捷键）

### 执行和调试
✅ 手动执行按钮（列表页和编辑器）
✅ 启用/禁用开关（工作流状态管理）
✅ 错误重试机制（支持指数退避策略）
✅ 节点测试运行（实时查看输出结构）
✅ 执行历史列表（状态过滤、搜索）
✅ 执行详情页（节点时间线、输入输出数据）

---

## 后端需要实现的功能

### 核心功能
1. **工作流存储**
   - 存储工作流定义（nodes + edges + envVars）
   - 环境变量加密存储
   - 支持版本管理（可选）

2. **调度引擎**
   - 根据trigger配置注册定时任务
   - 支持多种调度类型（daily/weekly/monthly/interval/cron）
   - Webhook触发器路由注册

3. **执行引擎**
   - 解析工作流DAG（有向无环图）
   - 按顺序执行节点
   - 处理条件分支
   - **变量解析和替换**（`{{env.KEY}}`、`{{node_id.field}}`、`{{trigger.data}}`）
   - **错误重试机制**（支持指数退避）
   - 处理异常和重试

4. **工具执行器**
   - HTTP请求执行器
   - 邮件发送执行器
   - 健康检查执行器

5. **监控和日志**
   - 记录执行历史
   - 记录每个节点的输入/输出
   - 错误日志和告警
   - 执行状态实时更新

6. **并发控制**
   - 工作流并发限制
   - 节点超时控制

7. **调试支持**
   - 节点测试运行接口
   - 手动触发执行接口

### 新增API接口

#### 执行相关
```
POST   /api/v1/workflows/:id/execute          # 手动执行工作流
POST   /api/v1/workflows/:id/nodes/:nodeId/test  # 测试节点
GET    /api/v1/workflows/:id/executions       # 执行历史列表
GET    /api/v1/workflows/:id/executions/:executionId  # 执行详情
```

#### Webhook触发
```
POST   /api/webhook/:webhookPath               # Webhook触发入口
GET    /api/webhook/:webhookPath               # Webhook触发入口（GET）
PUT    /api/webhook/:webhookPath               # Webhook触发入口（PUT）
```

#### 环境变量
```
GET    /api/v1/workflows/:id/env-vars          # 获取环境变量列表
PUT    /api/v1/workflows/:id/env-vars          # 更新环境变量
```

### 可选功能
- 工作流模板
- 执行统计和报表
- 工作流导入功能
- 工作流暂停/恢复
- 实时执行状态推送（WebSocket）

---

## 技术建议

### 后端技术选型
- **调度引擎**: Cron / APScheduler / Celery Beat
- **任务队列**: Redis + Bull / Celery
- **工作流引擎**: 自研 / Temporal / Airflow
- **存储**: PostgreSQL / MongoDB

### 执行模型
```
触发器 → 创建执行实例 → 加入任务队列 → Worker执行 → 记录结果
```

### DAG执行算法
```python
def execute_workflow(workflow):
    execution = create_execution(workflow)

    # 找到起始节点（trigger）
    current_nodes = find_trigger_nodes(workflow)

    while current_nodes:
        next_nodes = []

        for node in current_nodes:
            # 执行节点
            output = execute_node(node, execution)

            # 根据节点类型决定下一步
            if node.type == 'condition':
                # 条件分支
                if output.result:
                    next_nodes += find_next_nodes(node, 'true')
                else:
                    next_nodes += find_next_nodes(node, 'false')

            elif node.type == 'switch':
                # 开关分支
                case = match_case(node.config, output)
                next_nodes += find_next_nodes(node, case)

            elif node.type == 'delay':
                # 延迟执行
                sleep(node.config.duration)
                next_nodes += find_next_nodes(node)

            else:
                # 普通节点
                next_nodes += find_next_nodes(node)

        current_nodes = next_nodes

    finish_execution(execution)
```

---

## 数据库设计建议

### workflows 表
```sql
CREATE TABLE workflows (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    enabled BOOLEAN DEFAULT true,
    trigger_config JSON NOT NULL,
    nodes JSON NOT NULL,
    edges JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);
```

### workflow_executions 表
```sql
CREATE TABLE workflow_executions (
    id VARCHAR(50) PRIMARY KEY,
    workflow_id VARCHAR(50) NOT NULL,
    status ENUM('pending', 'running', 'success', 'failed') DEFAULT 'pending',
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    duration INT,
    error_message TEXT,
    nodes_executed JSON,
    FOREIGN KEY (workflow_id) REFERENCES workflows(id)
);
```

### workflow_execution_logs 表
```sql
CREATE TABLE workflow_execution_logs (
    id BIGINT PRIMARY KEY AUTO_INCREMENT,
    execution_id VARCHAR(50) NOT NULL,
    node_id VARCHAR(50) NOT NULL,
    node_name VARCHAR(255),
    status ENUM('pending', 'running', 'success', 'failed'),
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    input JSON,
    output JSON,
    error TEXT,
    FOREIGN KEY (execution_id) REFERENCES workflow_executions(id)
);
```

---

## 前端TODO（后续优化）

### 功能增强
- [ ] 工作流导入功能
- [ ] 撤销/重做功能
- [ ] 节点复制/粘贴
- [ ] 批量操作节点
- [ ] 节点搜索
- [ ] 缩略图导航
- [ ] 实时执行状态显示（WebSocket）
- [ ] 执行日志实时流
- [ ] 工作流模板库

### UI优化
- [ ] 暗色主题
- [ ] 自定义节点颜色
- [ ] 节点图标库
- [ ] 连接线动画
- [ ] 更多节点布局算法
- [ ] 节点分组功能

### 配置增强
- [ ] 条件表达式编辑器增强
- [ ] 变量智能提示
- [ ] 节点输出字段智能补全
- [ ] 更多内置函数（日期、字符串处理等）

---

## 常见问题

### Q1: 如何处理循环依赖？
A: 前端不阻止用户创建循环，但后端执行引擎需要检测DAG中的环，拒绝执行包含环的工作流。

### Q2: 节点执行失败如何处理？
A: 可以配置节点级别的重试策略，或者在工作流中明确添加错误处理分支。

### Q3: 如何传递节点之间的数据？
A: 每个节点的输出会作为下一个节点的输入（`input`变量），条件判断等可以引用这些数据。

### Q4: 并行执行如何支持？
A: 当前版本不支持并行，所有节点串行执行。未来可以添加"并行网关"节点支持并行执行。

---

## 更新日志

### v1.1.0 (2025-01-12) - 数据流和调试功能
- ✅ 变量系统（环境变量、节点输出引用、触发器数据）
- ✅ 可视化变量选择器
- ✅ Webhook触发器
- ✅ 错误重试机制（支持指数退避）
- ✅ 节点测试运行
- ✅ 手动执行按钮
- ✅ 启用/禁用开关
- ✅ 执行历史和详情页

### v1.0.0 (2025-01-10)
- ✅ 基础工作流编辑器
- ✅ 触发器节点（定时、手动）
- ✅ 工具节点（HTTP、邮件、健康检查）
- ✅ 条件判断节点
- ✅ 开关分支节点
- ✅ 延迟等待节点
- ✅ cURL解析功能

---

## 联系方式

如有问题或建议，请联系开发团队。

