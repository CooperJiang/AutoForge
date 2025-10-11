# AutoForge 架构设计文档

> 从定时任务工具到智能自动化平台的演进

---

## 📋 目录

- [设计理念](#设计理念)
- [核心架构](#核心架构)
- [工具系统设计](#工具系统设计)
- [数据模型](#数据模型)
- [工作流引擎](#工作流引擎)
- [实施路线图](#实施路线图)

---

## 🎯 设计理念

### 当前问题

**v1.0 架构局限**：
- ❌ 定时任务 + HTTP 请求耦合在一起
- ❌ 只能调用 API 接口，功能单一
- ❌ 无法扩展其他类型的操作（邮件、短信、爬虫等）
- ❌ 不支持多步骤工作流

**痛点**：
```
用户需求：定时发送邮件
当前方案：无法实现，只能调用 HTTP API
```

### 设计目标

**v2.0 架构愿景**：
- ✅ **解耦**：定时调度 ↔️ 执行动作（工具）
- ✅ **可扩展**：插件化工具系统，易于添加新工具
- ✅ **组合性**：支持多步骤工作流（类似 n8n）
- ✅ **标准化**：统一的工具接口，支持 AI Agent 调用
- ✅ **可观测**：完整的执行日志和监控

---

## 🏗️ 核心架构

### 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                         用户界面                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │ 任务管理 │  │ 工具库   │  │ 工作流   │  │ 执行记录 │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                        API 层 (Gin)                          │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐    │
│  │  任务API │  │ 工具API  │  │ 工作流API│  │ 执行API  │    │
│  └──────────┘  └──────────┘  └──────────┘  └──────────┘    │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       业务逻辑层                              │
│  ┌────────────────┐  ┌────────────────┐  ┌──────────────┐  │
│  │  任务调度器    │  │  工具注册表    │  │  工作流引擎  │  │
│  │  (Cron)        │  │  (Registry)    │  │  (Engine)    │  │
│  └────────────────┘  └────────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                       工具执行层                              │
│  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐         │
│  │ HTTP │  │ Email│  │ SMS  │  │Crawler│  │ DB  │  ...    │
│  └──────┘  └──────┘  └──────┘  └──────┘  └──────┘         │
└─────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────┐
│                      数据存储层                               │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐      │
│  │   MySQL      │  │    Redis     │  │  File System │      │
│  └──────────────┘  └──────────────┘  └──────────────┘      │
└─────────────────────────────────────────────────────────────┘
```

### 核心概念

#### 1. **任务 (Task)**
- 定时触发器
- 定义"何时"执行
- 关联一个工具实例或工作流

#### 2. **工具 (Tool)**
- 可执行的原子操作
- 定义"做什么"
- 标准化接口，可被 AI Agent 调用

#### 3. **工具实例 (Tool Instance)**
- 用户配置的具体工具
- 包含工具的参数配置

#### 4. **工作流 (Workflow)**
- 多个工具的组合
- 支持条件、循环、分支
- 类似 n8n 的节点流程

---

## 🔧 工具系统设计

### 设计原则

1. **标准化接口**：所有工具实现统一接口
2. **自描述**：工具提供 Schema 描述自己的配置
3. **可组合**：工具输出可作为下一个工具的输入
4. **AI 友好**：接口设计支持 LLM/Agent 调用
5. **可观测**：完整的执行日志和错误信息

### 工具接口定义

```go
// pkg/utools/types.go

// Tool 工具接口（所有工具必须实现）
type Tool interface {
    // GetMetadata 获取工具元数据
    GetMetadata() *ToolMetadata

    // GetSchema 获取配置 Schema（JSON Schema 格式）
    GetSchema() *ConfigSchema

    // Validate 验证配置是否合法
    Validate(config map[string]interface{}) error

    // Execute 执行工具
    Execute(ctx *ExecutionContext, config map[string]interface{}) (*ExecutionResult, error)
}

// ToolMetadata 工具元数据
type ToolMetadata struct {
    Code        string   `json:"code"`         // 唯一标识：http_request
    Name        string   `json:"name"`         // 显示名称：HTTP 请求
    Description string   `json:"description"`  // 描述
    Category    string   `json:"category"`     // 分类：network/notification/data
    Version     string   `json:"version"`      // 版本：1.0.0
    Author      string   `json:"author"`       // 作者
    Icon        string   `json:"icon"`         // 图标（emoji 或 URL）
    Tags        []string `json:"tags"`         // 标签
    AICallable  bool     `json:"ai_callable"`  // 是否可被 AI 调用
}

// ConfigSchema 配置 Schema（JSON Schema 标准）
type ConfigSchema struct {
    Type       string                 `json:"type"`       // object
    Required   []string               `json:"required"`   // 必填字段
    Properties map[string]PropertyDef `json:"properties"` // 字段定义
}

// PropertyDef 字段定义
type PropertyDef struct {
    Type        string      `json:"type"`        // string/number/boolean/object/array
    Title       string      `json:"title"`       // 显示标题
    Description string      `json:"description"` // 描述
    Default     interface{} `json:"default"`     // 默认值
    Enum        []string    `json:"enum"`        // 枚举值
    Format      string      `json:"format"`      // 格式：email/url/date-time
    MinLength   int         `json:"minLength"`   // 最小长度
    MaxLength   int         `json:"maxLength"`   // 最大长度
    Pattern     string      `json:"pattern"`     // 正则表达式
}

// ExecutionContext 执行上下文
type ExecutionContext struct {
    TaskID      string                 // 任务ID
    UserID      string                 // 用户ID
    ExecutionID string                 // 执行ID
    Variables   map[string]interface{} // 上下文变量（工作流用）
    Logger      Logger                 // 日志记录器
    Timeout     time.Duration          // 超时时间
}

// ExecutionResult 执行结果
type ExecutionResult struct {
    Success    bool                   `json:"success"`     // 是否成功
    Output     map[string]interface{} `json:"output"`      // 输出数据
    Error      string                 `json:"error"`       // 错误信息
    StatusCode int                    `json:"status_code"` // HTTP 状态码（如适用）
    Duration   int64                  `json:"duration"`    // 执行时长（毫秒）
    Logs       []string               `json:"logs"`        // 执行日志
    Metadata   map[string]interface{} `json:"metadata"`    // 元数据
}
```

### 工具目录结构

```
pkg/utools/
├── types.go              # 接口定义
├── registry.go           # 工具注册表
├── base.go               # 基础工具类
├── http/                 # HTTP 请求工具
│   ├── http_tool.go      # 工具实现
│   ├── http_tool_test.go # 单元测试
│   └── README.md         # 工具文档
├── email/                # 邮件工具
│   ├── email_tool.go
│   ├── email_tool_test.go
│   └── README.md
├── sms/                  # 短信工具
├── database/             # 数据库工具
├── crawler/              # 爬虫工具
├── webhook/              # Webhook 工具
├── json/                 # JSON 处理工具
├── file/                 # 文件操作工具
└── ai/                   # AI 工具（OpenAI/Claude）
```

### HTTP 工具示例

```go
// pkg/utools/http/http_tool.go
package http

import "auto-forge/pkg/utools"

type HTTPTool struct {
    utools.BaseTool // 继承基础实现
}

func (t *HTTPTool) GetMetadata() *utools.ToolMetadata {
    return &utools.ToolMetadata{
        Code:        "http_request",
        Name:        "HTTP 请求",
        Description: "向指定 URL 发送 HTTP 请求，支持 GET/POST/PUT/DELETE 等方法",
        Category:    "network",
        Version:     "1.0.0",
        Author:      "AutoForge Team",
        Icon:        "🌐",
        Tags:        []string{"http", "api", "network", "web"},
        AICallable:  true,
    }
}

func (t *HTTPTool) GetSchema() *utools.ConfigSchema {
    return &utools.ConfigSchema{
        Type:     "object",
        Required: []string{"url", "method"},
        Properties: map[string]utools.PropertyDef{
            "url": {
                Type:        "string",
                Title:       "请求地址",
                Description: "目标 API 的完整 URL",
                Format:      "uri",
            },
            "method": {
                Type:        "string",
                Title:       "请求方法",
                Description: "HTTP 请求方法",
                Enum:        []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
                Default:     "GET",
            },
            "headers": {
                Type:        "object",
                Title:       "请求头",
                Description: "自定义 HTTP 请求头",
            },
            "body": {
                Type:        "object",
                Title:       "请求体",
                Description: "POST/PUT 请求的 Body 数据",
            },
            "timeout": {
                Type:        "number",
                Title:       "超时时间",
                Description: "请求超时时间（秒）",
                Default:     30,
            },
        },
    }
}

func (t *HTTPTool) Validate(config map[string]interface{}) error {
    // 验证 URL 格式
    // 验证 method 是否在枚举值内
    // ...
    return nil
}

func (t *HTTPTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
    startTime := time.Now()

    // 1. 提取配置
    url := config["url"].(string)
    method := config["method"].(string)

    // 2. 构造请求
    req, err := http.NewRequest(method, url, nil)
    if err != nil {
        return &utools.ExecutionResult{
            Success: false,
            Error:   err.Error(),
        }, err
    }

    // 3. 添加请求头
    if headers, ok := config["headers"].(map[string]interface{}); ok {
        for k, v := range headers {
            req.Header.Set(k, fmt.Sprint(v))
        }
    }

    // 4. 发送请求
    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return &utools.ExecutionResult{
            Success: false,
            Error:   err.Error(),
        }, err
    }
    defer resp.Body.Close()

    // 5. 读取响应
    body, _ := io.ReadAll(resp.Body)

    // 6. 返回结果
    return &utools.ExecutionResult{
        Success:    resp.StatusCode >= 200 && resp.StatusCode < 300,
        StatusCode: resp.StatusCode,
        Output: map[string]interface{}{
            "status_code": resp.StatusCode,
            "headers":     resp.Header,
            "body":        string(body),
        },
        Duration: time.Since(startTime).Milliseconds(),
        Logs: []string{
            fmt.Sprintf("Request: %s %s", method, url),
            fmt.Sprintf("Response: %d", resp.StatusCode),
        },
    }, nil
}

// 注册工具
func init() {
    utools.Register(&HTTPTool{})
}
```

### 工具注册表

```go
// pkg/utools/registry.go
package utools

import (
    "fmt"
    "sync"
)

var (
    registry     = make(map[string]Tool)
    registryLock sync.RWMutex
)

// Register 注册工具
func Register(tool Tool) {
    registryLock.Lock()
    defer registryLock.Unlock()

    metadata := tool.GetMetadata()
    if metadata.Code == "" {
        panic("tool code cannot be empty")
    }

    registry[metadata.Code] = tool
}

// GetTool 获取工具
func GetTool(code string) (Tool, error) {
    registryLock.RLock()
    defer registryLock.RUnlock()

    if tool, ok := registry[code]; ok {
        return tool, nil
    }
    return nil, fmt.Errorf("tool not found: %s", code)
}

// ListTools 列出所有工具
func ListTools() []Tool {
    registryLock.RLock()
    defer registryLock.RUnlock()

    tools := make([]Tool, 0, len(registry))
    for _, tool := range registry {
        tools = append(tools, tool)
    }
    return tools
}

// ListToolsByCategory 按分类列出工具
func ListToolsByCategory(category string) []Tool {
    tools := ListTools()
    filtered := make([]Tool, 0)

    for _, tool := range tools {
        if tool.GetMetadata().Category == category {
            filtered = append(filtered, tool)
        }
    }
    return filtered
}
```

---

## 📊 数据模型

### 数据库表设计

#### 1. tools 表（工具定义）

```sql
CREATE TABLE `tools` (
  `id` varchar(36) PRIMARY KEY,
  `code` varchar(50) UNIQUE NOT NULL,           -- 工具唯一标识
  `name` varchar(100) NOT NULL,                 -- 工具名称
  `description` text,                           -- 描述
  `category` varchar(50),                       -- 分类
  `version` varchar(20),                        -- 版本
  `author` varchar(100),                        -- 作者
  `icon` varchar(255),                          -- 图标
  `tags` json,                                  -- 标签数组
  `ai_callable` boolean DEFAULT false,          -- 是否可被 AI 调用
  `config_schema` json NOT NULL,                -- 配置 Schema
  `enabled` boolean DEFAULT true,               -- 是否启用
  `created_at` bigint,
  `updated_at` bigint,
  INDEX `idx_category` (`category`),
  INDEX `idx_enabled` (`enabled`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 2. tool_instances 表（工具实例）

```sql
CREATE TABLE `tool_instances` (
  `id` varchar(36) PRIMARY KEY,
  `user_id` varchar(36) NOT NULL,               -- 所属用户
  `tool_code` varchar(50) NOT NULL,             -- 引用工具
  `name` varchar(200) NOT NULL,                 -- 实例名称（用户自定义）
  `description` text,                           -- 描述
  `config` json NOT NULL,                       -- 配置数据（加密存储敏感信息）
  `enabled` boolean DEFAULT true,               -- 是否启用
  `created_at` bigint,
  `updated_at` bigint,
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_tool_code` (`tool_code`),
  FOREIGN KEY (`tool_code`) REFERENCES `tools`(`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 3. tasks 表（任务 - 重构）

```sql
CREATE TABLE `tasks` (
  `id` varchar(36) PRIMARY KEY,
  `user_id` varchar(36) NOT NULL,
  `name` varchar(200) NOT NULL,
  `description` text,

  -- 调度配置
  `schedule_type` varchar(20) NOT NULL,         -- daily/weekly/monthly/interval/cron
  `schedule_rule` varchar(200) NOT NULL,        -- 调度规则
  `timezone` varchar(50) DEFAULT 'Asia/Shanghai',
  `enabled` boolean DEFAULT true,

  -- 执行目标（核心变化）
  `execution_type` varchar(20) NOT NULL,        -- 'tool' 或 'workflow'
  `execution_id` varchar(36) NOT NULL,          -- tool_instance.id 或 workflow.id

  -- 执行状态
  `next_run_time` bigint,                       -- 下次执行时间
  `last_run_time` bigint,                       -- 上次执行时间
  `last_run_status` varchar(20),                -- success/failed

  `created_at` bigint,
  `updated_at` bigint,
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_enabled` (`enabled`),
  INDEX `idx_next_run_time` (`next_run_time`),
  INDEX `idx_execution` (`execution_type`, `execution_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 4. task_executions 表（执行记录 - 扩展）

```sql
CREATE TABLE `task_executions` (
  `id` varchar(36) PRIMARY KEY,
  `task_id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `execution_type` varchar(20) NOT NULL,        -- 'tool' 或 'workflow'
  `execution_id` varchar(36) NOT NULL,

  -- 执行结果
  `status` varchar(20) NOT NULL,                -- success/failed/timeout
  `output` json,                                -- 输出数据
  `error` text,                                 -- 错误信息
  `logs` json,                                  -- 执行日志数组

  -- 性能指标
  `started_at` bigint NOT NULL,
  `completed_at` bigint,
  `duration` int,                               -- 执行时长（毫秒）

  `created_at` bigint,
  INDEX `idx_task_id` (`task_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_status` (`status`),
  INDEX `idx_started_at` (`started_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

#### 5. workflows 表（工作流 - 未来）

```sql
CREATE TABLE `workflows` (
  `id` varchar(36) PRIMARY KEY,
  `user_id` varchar(36) NOT NULL,
  `name` varchar(200) NOT NULL,
  `description` text,
  `nodes` json NOT NULL,                        -- 节点定义（工具、条件、循环）
  `edges` json NOT NULL,                        -- 连接关系
  `variables` json,                             -- 全局变量
  `version` int DEFAULT 1,                      -- 版本号
  `enabled` boolean DEFAULT true,
  `created_at` bigint,
  `updated_at` bigint,
  INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```

### Go Model 定义

```go
// internal/models/tool.go
type Tool struct {
    BaseModel
    Code         string         `gorm:"size:50;uniqueIndex" json:"code"`
    Name         string         `gorm:"size:100" json:"name"`
    Description  string         `gorm:"type:text" json:"description"`
    Category     string         `gorm:"size:50;index" json:"category"`
    Version      string         `gorm:"size:20" json:"version"`
    Author       string         `gorm:"size:100" json:"author"`
    Icon         string         `gorm:"size:255" json:"icon"`
    Tags         datatypes.JSON `json:"tags"`
    AICallable   bool           `json:"ai_callable"`
    ConfigSchema datatypes.JSON `gorm:"type:json" json:"config_schema"`
    Enabled      bool           `gorm:"index" json:"enabled"`
}

// internal/models/tool_instance.go
type ToolInstance struct {
    BaseModel
    UserID      UUID           `gorm:"index" json:"user_id"`
    ToolCode    string         `gorm:"size:50;index" json:"tool_code"`
    Name        string         `gorm:"size:200" json:"name"`
    Description string         `gorm:"type:text" json:"description"`
    Config      datatypes.JSON `gorm:"type:json" json:"config"`
    Enabled     bool           `json:"enabled"`

    // 关联
    Tool        *Tool          `gorm:"foreignKey:ToolCode;references:Code" json:"tool,omitempty"`
}

// internal/models/task.go (重构)
type Task struct {
    BaseModel
    UserID       UUID   `gorm:"index" json:"user_id"`
    Name         string `gorm:"size:200" json:"name"`
    Description  string `gorm:"type:text" json:"description"`

    // 调度配置
    ScheduleType string `gorm:"size:20" json:"schedule_type"`
    ScheduleRule string `gorm:"size:200" json:"schedule_rule"`
    Timezone     string `gorm:"size:50" json:"timezone"`
    Enabled      bool   `gorm:"index" json:"enabled"`

    // 执行目标
    ExecutionType string `gorm:"size:20;index:idx_execution" json:"execution_type"`
    ExecutionID   UUID   `gorm:"index:idx_execution" json:"execution_id"`

    // 执行状态
    NextRunTime    int64  `gorm:"index" json:"next_run_time"`
    LastRunTime    int64  `json:"last_run_time"`
    LastRunStatus  string `gorm:"size:20" json:"last_run_status"`
}
```

---

## 🔄 工作流引擎（第二阶段）

### 节点类型

```go
type NodeType string

const (
    NodeTypeTool      NodeType = "tool"       // 工具节点
    NodeTypeCondition NodeType = "condition"  // 条件判断
    NodeTypeLoop      NodeType = "loop"       // 循环
    NodeTypeMerge     NodeType = "merge"      // 合并
    NodeTypeDelay     NodeType = "delay"      // 延迟
    NodeTypeTrigger   NodeType = "trigger"    // 触发器
)

type WorkflowNode struct {
    ID       string                 `json:"id"`
    Type     NodeType               `json:"type"`
    Name     string                 `json:"name"`
    ToolCode string                 `json:"tool_code,omitempty"` // 工具节点
    Config   map[string]interface{} `json:"config"`
    Position map[string]int         `json:"position"` // x, y 坐标
}

type WorkflowEdge struct {
    ID      string `json:"id"`
    Source  string `json:"source"`  // 源节点ID
    Target  string `json:"target"`  // 目标节点ID
    Label   string `json:"label"`   // 连接标签（条件）
}
```

### 工作流示例

**场景**：爬取新闻 → 判断关键词 → 发送邮件/短信

```json
{
  "nodes": [
    {
      "id": "node-1",
      "type": "tool",
      "name": "爬取新闻",
      "tool_code": "web_crawler",
      "config": {
        "url": "https://news.example.com",
        "selector": ".news-title"
      }
    },
    {
      "id": "node-2",
      "type": "condition",
      "name": "判断关键词",
      "config": {
        "expression": "output.title.contains('AI')"
      }
    },
    {
      "id": "node-3",
      "type": "tool",
      "name": "发送邮件",
      "tool_code": "send_email",
      "config": {
        "to": "user@example.com",
        "subject": "发现AI新闻",
        "body": "{{node-1.output.title}}"
      }
    },
    {
      "id": "node-4",
      "type": "tool",
      "name": "发送短信",
      "tool_code": "send_sms",
      "config": {
        "phone": "13800138000",
        "message": "新闻已更新"
      }
    }
  ],
  "edges": [
    {"source": "node-1", "target": "node-2"},
    {"source": "node-2", "target": "node-3", "label": "true"},
    {"source": "node-2", "target": "node-4", "label": "false"}
  ]
}
```

---

## 🗓️ 实施路线图

### 阶段一：工具抽象层（1-2周）

**Week 1: 核心架构**
- [ ] 创建 `pkg/utools` 目录结构
- [ ] 定义工具接口和类型
- [ ] 实现工具注册表
- [ ] 创建 Tool、ToolInstance 数据模型
- [ ] 数据库迁移脚本

**Week 2: HTTP 工具重构**
- [ ] HTTP 工具实现（迁移现有逻辑）
- [ ] 数据迁移：Task → ToolInstance
- [ ] 更新调度器支持工具执行
- [ ] API 接口适配
- [ ] 前端工具选择界面

**验收标准**：
- ✅ 现有 HTTP 任务正常运行
- ✅ 可以通过工具库创建新任务
- ✅ 向后兼容（旧任务不受影响）

---

### 阶段二：扩展工具库（2-3周）

**Week 3-4: 基础工具**
- [ ] 邮件工具（复用现有邮件服务）
- [ ] Webhook 工具
- [ ] JSON 处理工具
- [ ] 数据库查询工具
- [ ] 文件操作工具

**Week 5: 高级工具**
- [ ] 爬虫工具（集成 colly）
- [ ] 短信工具
- [ ] AI 工具（OpenAI/Claude API）
- [ ] RSS 订阅工具

**验收标准**：
- ✅ 至少 8 个工具可用
- ✅ 工具文档完整
- ✅ 单元测试覆盖 80%+

---

### 阶段三：工作流引擎（3-4周）

**Week 6-7: 工作流核心**
- [ ] Workflow 数据模型
- [ ] 工作流执行引擎
- [ ] 节点调度算法
- [ ] 变量传递机制
- [ ] 错误处理和重试

**Week 8-9: 工作流编辑器**
- [ ] 可视化拖拽编辑器（React Flow）
- [ ] 节点连接和配置
- [ ] 工作流测试运行
- [ ] 工作流版本管理

**验收标准**：
- ✅ 支持线性工作流
- ✅ 支持条件分支
- ✅ 可视化编辑和调试

---

### 阶段四：AI Agent 集成（2周）

**Week 10-11: AI 能力**
- [ ] 工具描述符生成（OpenAI Function Calling 格式）
- [ ] AI Agent 调用接口
- [ ] 自然语言任务创建
- [ ] 智能推荐工具

**验收标准**：
- ✅ AI 可以调用工具库
- ✅ 用户可以用自然语言创建任务

---

## 📚 技术选型

### 后端技术

| 组件 | 技术 | 说明 |
|------|------|------|
| 工具执行 | Goroutine Pool | 并发执行工具 |
| JSON Schema | gojsonschema | 配置验证 |
| 工作流引擎 | DAG 算法 | 有向无环图调度 |
| 加密存储 | AES-256 | 敏感配置加密 |

### 前端技术

| 组件 | 技术 | 说明 |
|------|------|------|
| 动态表单 | JSON Schema Form | 根据 Schema 生成表单 |
| 工作流编辑器 | React Flow / Vue Flow | 可视化拖拽 |
| 代码编辑器 | Monaco Editor | JSON/代码编辑 |

---

## 🔐 安全考虑

### 1. 配置加密
```go
// 敏感字段加密存储
type ToolInstance struct {
    Config datatypes.JSON `gorm:"type:json" json:"config"`
    // Config 存储前加密，读取时解密
}
```

### 2. 权限控制
```go
// 工具权限
type ToolPermission struct {
    UserID   UUID
    ToolCode string
    CanUse   bool
}
```

### 3. 执行限制
- 单个工具执行超时：30s
- 工作流总执行时间：5min
- 用户每日执行次数限制
- 并发执行数限制

---

## 📊 监控和观测

### 执行指标
- 工具执行成功率
- 平均执行时长
- 错误类型统计
- 用户活跃度

### 日志
```go
type ExecutionLog struct {
    Level     string // INFO/WARN/ERROR
    Message   string
    Timestamp int64
    Metadata  map[string]interface{}
}
```

---

## 🎯 成功指标

### 技术指标
- [ ] 工具执行成功率 > 95%
- [ ] API 响应时间 < 200ms
- [ ] 单元测试覆盖率 > 80%
- [ ] 系统可用性 > 99.9%

### 业务指标
- [ ] 工具库数量 > 10 个
- [ ] 工作流创建数 > 100
- [ ] 用户留存率 > 60%
- [ ] NPS 评分 > 50

---

## 📝 开发规范

### 工具开发规范

1. **目录结构**
```
pkg/utools/your_tool/
├── your_tool.go          # 工具实现
├── your_tool_test.go     # 单元测试
├── README.md             # 工具文档
└── examples/             # 使用示例
```

2. **命名规范**
- 工具 Code：小写下划线 `send_email`
- Go 结构体：驼峰命名 `EmailTool`
- 配置字段：小写下划线 `smtp_host`

3. **文档要求**
- 工具用途说明
- 配置参数说明
- 使用示例
- 常见问题

4. **测试要求**
- 单元测试覆盖核心逻辑
- Mock 外部依赖
- 测试用例包含正常和异常场景

---

## 🔗 相关资源

- [JSON Schema 规范](https://json-schema.org/)
- [n8n 架构设计](https://docs.n8n.io/hosting/architecture/)
- [OpenAI Function Calling](https://platform.openai.com/docs/guides/function-calling)
- [LangChain Tools](https://python.langchain.com/docs/modules/tools/)

---

## 📧 联系方式

如有疑问或建议，请联系：
- GitHub Issues: https://github.com/CooperJiang/AutoForge/issues
- 架构讨论：创建 Discussion

---

**文档版本**: v1.0
**最后更新**: 2025-01-11
**维护者**: AutoForge Team
