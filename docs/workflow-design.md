# AutoForge 工作流系统设计方案

> 版本：v1.0
> 作者：AutoForge Team
> 日期：2025-01-12

## 📋 目录

- [1. 概述](#1-概述)
- [2. 核心概念](#2-核心概念)
- [3. 技术架构](#3-技术架构)
- [4. 数据结构设计](#4-数据结构设计)
- [5. 执行引擎设计](#5-执行引擎设计)
- [6. 前端可视化编辑器](#6-前端可视化编辑器)
- [7. 实现路线图](#7-实现路线图)
- [8. 使用示例](#8-使用示例)

---

## 1. 概述

### 1.1 背景

当前 AutoForge 支持单个工具的定时执行，但实际业务场景往往需要多个工具协同工作。例如：
- 监控网站健康状态，失败时发送邮件告警
- 定期备份数据库，成功后上传到云存储
- 爬取数据后进行处理，再发送通知

### 1.2 目标

构建类似 n8n/Zapier 的工作流系统，支持：
- ✅ 多工具串联执行
- ✅ 条件分支（if-else）
- ✅ 数据传递和变量替换
- ✅ 可视化流程编辑
- ✅ 错误处理和重试
- ✅ 执行日志和调试

### 1.3 设计原则

1. **渐进式实现**：从简单到复杂，分阶段实施
2. **向后兼容**：保持现有单工具任务的功能
3. **易用性优先**：提供可视化编辑，降低使用门槛
4. **可扩展性**：架构支持未来的高级特性

---

## 2. 核心概念

### 2.1 工作流（Workflow）

工作流是由多个节点和连接组成的有向图（DAG），定义了工具的执行顺序和条件。

```
[节点A] --条件--> [节点B] --always--> [节点C]
```

### 2.2 节点（Node）

节点是工作流的基本执行单元，每个节点对应一个工具的配置。

**节点属性**：
- `id`：唯一标识符
- `name`：节点名称（用户自定义）
- `tool_code`：工具代码（http_request, email_sender 等）
- `config`：工具配置参数
- `position`：画布上的位置（x, y）

### 2.3 边（Edge）

边连接两个节点，定义执行流向和条件。

**边属性**：
- `id`：唯一标识符
- `source`：源节点 ID
- `target`：目标节点 ID
- `condition`：执行条件（可选）

### 2.4 条件（Condition）

条件决定是否执行目标节点。

**条件类型**：
- `always`：始终执行
- `never`：从不执行
- `on_success`：前置节点成功时执行
- `on_failure`：前置节点失败时执行
- `expression`：自定义表达式

### 2.5 变量传递

使用模板语法在节点间传递数据：

```
{{node_id.field}}
{{node_id.output.nested.field}}
```

**示例**：
```json
{
  "subject": "告警 - {{health_check.output.url}}",
  "body": "错误信息：{{health_check.message}}"
}
```

---

## 3. 技术架构

### 3.1 架构图

```
┌─────────────────────────────────────────────────────┐
│                     前端层                           │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐ │
│  │ 可视化编辑器  │  │  JSON 编辑器 │  │ 执行监控  │ │
│  └──────────────┘  └──────────────┘  └───────────┘ │
└─────────────────────────────────────────────────────┘
                         │ REST API
┌─────────────────────────────────────────────────────┐
│                     后端层                           │
│  ┌──────────────┐  ┌──────────────┐  ┌───────────┐ │
│  │ 工作流控制器  │  │  任务调度器  │  │ 执行记录  │ │
│  └──────────────┘  └──────────────┘  └───────────┘ │
│                         │                            │
│  ┌──────────────────────────────────────────────┐  │
│  │           工作流执行引擎                      │  │
│  │  ┌────────────┐  ┌────────────┐  ┌────────┐ │  │
│  │  │ 拓扑排序   │  │ 条件评估   │  │ 变量解析│ │  │
│  │  └────────────┘  └────────────┘  └────────┘ │  │
│  └──────────────────────────────────────────────┘  │
│                         │                            │
│  ┌──────────────────────────────────────────────┐  │
│  │              工具系统 (utools)                │  │
│  │  [HTTP] [Email] [Health] [Database] [...]   │  │
│  └──────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────┘
                         │
┌─────────────────────────────────────────────────────┐
│                   数据持久层                         │
│  ┌──────────┐  ┌──────────┐  ┌──────────────────┐ │
│  │ 任务表   │  │ 工作流表 │  │ 工作流执行记录    │ │
│  └──────────┘  └──────────┘  └──────────────────┘ │
└─────────────────────────────────────────────────────┘
```

### 3.2 技术选型

#### 后端
- **执行引擎**：Go（自研）
- **表达式解析**：govaluate 或自实现
- **JSON 处理**：encoding/json
- **并发控制**：goroutines + channels

#### 前端
- **流程图库**：React Flow 或 Vue Flow（推荐）
- **UI 组件**：现有 Vue 3 + Tailwind CSS
- **状态管理**：Pinia
- **图形渲染**：SVG（通过 Vue Flow）

### 3.3 关键技术点

1. **拓扑排序**：确定节点执行顺序
2. **DAG 检测**：防止循环依赖
3. **变量插值**：正则表达式匹配和替换
4. **条件评估**：表达式解析引擎
5. **执行隔离**：Context 传递和超时控制

---

## 4. 数据结构设计

### 4.1 工作流定义

```json
{
  "version": "1.0",
  "name": "网站健康监控工作流",
  "description": "检查网站状态，失败时发送邮件告警",
  "nodes": [
    {
      "id": "node_1",
      "name": "健康检查",
      "tool_code": "health_checker",
      "position": { "x": 100, "y": 100 },
      "config": {
        "url": "https://example.com",
        "method": "GET",
        "timeout": 10,
        "expected_status": 200
      }
    },
    {
      "id": "node_2",
      "name": "发送告警邮件",
      "tool_code": "email_sender",
      "position": { "x": 400, "y": 100 },
      "config": {
        "to": "admin@company.com",
        "subject": "⚠️ 网站告警 - {{node_1.output.url}}",
        "body": "网站检查失败！\n\n状态：{{node_1.message}}\n响应时间：{{node_1.output.response_time}}ms\n检查时间：{{timestamp}}"
      }
    }
  ],
  "edges": [
    {
      "id": "edge_1",
      "source": "node_1",
      "target": "node_2",
      "label": "失败时",
      "condition": {
        "type": "on_failure"
      }
    }
  ],
  "settings": {
    "continue_on_error": false,
    "timeout": 300,
    "retry_on_failure": false,
    "max_retries": 3
  }
}
```

### 4.2 数据库表设计

#### 方案 A：简单方案（推荐第一阶段）

```sql
-- 扩展现有 tasks 表
ALTER TABLE tasks ADD COLUMN task_type VARCHAR(20) DEFAULT 'single_tool';
-- task_type: 'single_tool' 或 'workflow'

ALTER TABLE tasks ADD COLUMN workflow_config TEXT;
-- 存储工作流 JSON（当 task_type='workflow' 时使用）

-- 示例数据
INSERT INTO tasks (name, task_type, workflow_config, schedule_type, schedule_value)
VALUES (
  '网站监控工作流',
  'workflow',
  '{"version":"1.0","nodes":[...],"edges":[...]}',
  'interval',
  '300'
);
```

#### 方案 B：完整方案（第二阶段）

```sql
-- 工作流表
CREATE TABLE workflows (
    id VARCHAR(36) PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    version VARCHAR(20) DEFAULT '1.0',
    config JSON NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_user_id (user_id)
);

-- 工作流执行记录
CREATE TABLE workflow_executions (
    id VARCHAR(36) PRIMARY KEY,
    workflow_id VARCHAR(36) NOT NULL,
    task_id VARCHAR(36),
    trigger_type VARCHAR(20),  -- schedule/manual/webhook
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    duration_ms BIGINT,
    status VARCHAR(20),  -- running/success/failed/partial
    node_results JSON,   -- 每个节点的执行结果
    execution_order TEXT, -- 节点执行顺序
    error TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    INDEX idx_workflow_id (workflow_id),
    INDEX idx_task_id (task_id),
    INDEX idx_status (status)
);

-- 工作流节点执行记录（详细日志）
CREATE TABLE workflow_node_executions (
    id VARCHAR(36) PRIMARY KEY,
    workflow_execution_id VARCHAR(36) NOT NULL,
    node_id VARCHAR(50) NOT NULL,
    node_name VARCHAR(100),
    tool_code VARCHAR(50),
    start_time TIMESTAMP,
    end_time TIMESTAMP,
    duration_ms BIGINT,
    status VARCHAR(20),
    input JSON,
    output JSON,
    error TEXT,
    INDEX idx_workflow_execution_id (workflow_execution_id)
);
```

### 4.3 Go 数据模型

```go
// pkg/workflow/types.go
package workflow

import "time"

// WorkflowDefinition 工作流定义
type WorkflowDefinition struct {
    Version     string           `json:"version"`
    Name        string           `json:"name"`
    Description string           `json:"description"`
    Nodes       []Node           `json:"nodes"`
    Edges       []Edge           `json:"edges"`
    Settings    WorkflowSettings `json:"settings"`
}

// Node 工作流节点
type Node struct {
    ID       string                 `json:"id"`
    Name     string                 `json:"name"`
    ToolCode string                 `json:"tool_code"`
    Position Position               `json:"position"`
    Config   map[string]interface{} `json:"config"`
}

// Position 节点位置
type Position struct {
    X int `json:"x"`
    Y int `json:"y"`
}

// Edge 节点连接
type Edge struct {
    ID        string    `json:"id"`
    Source    string    `json:"source"`
    Target    string    `json:"target"`
    Label     string    `json:"label,omitempty"`
    Condition Condition `json:"condition"`
}

// Condition 执行条件
type Condition struct {
    Type       string `json:"type"` // always/never/on_success/on_failure/expression
    Expression string `json:"expression,omitempty"`
}

// WorkflowSettings 工作流设置
type WorkflowSettings struct {
    ContinueOnError bool `json:"continue_on_error"`
    Timeout         int  `json:"timeout"`
    RetryOnFailure  bool `json:"retry_on_failure"`
    MaxRetries      int  `json:"max_retries"`
}

// WorkflowExecution 工作流执行记录
type WorkflowExecution struct {
    ID             string                            `json:"id"`
    WorkflowID     string                            `json:"workflow_id"`
    TaskID         string                            `json:"task_id"`
    TriggerType    string                            `json:"trigger_type"`
    StartTime      time.Time                         `json:"start_time"`
    EndTime        time.Time                         `json:"end_time"`
    DurationMs     int64                             `json:"duration_ms"`
    Status         string                            `json:"status"`
    NodeResults    map[string]*NodeExecutionResult   `json:"node_results"`
    ExecutionOrder []string                          `json:"execution_order"`
    Error          string                            `json:"error,omitempty"`
}

// NodeExecutionResult 节点执行结果
type NodeExecutionResult struct {
    NodeID     string                 `json:"node_id"`
    NodeName   string                 `json:"node_name"`
    ToolCode   string                 `json:"tool_code"`
    StartTime  time.Time              `json:"start_time"`
    EndTime    time.Time              `json:"end_time"`
    DurationMs int64                  `json:"duration_ms"`
    Status     string                 `json:"status"`
    Success    bool                   `json:"success"`
    Message    string                 `json:"message"`
    Output     map[string]interface{} `json:"output"`
    Error      string                 `json:"error,omitempty"`
}
```

---

## 5. 执行引擎设计

### 5.1 核心算法

#### 拓扑排序（Topological Sort）

确定节点的执行顺序，检测循环依赖。

```go
func (e *WorkflowEngine) topologicalSort(workflow *WorkflowDefinition) ([]string, error) {
    // 构建入度表和邻接表
    inDegree := make(map[string]int)
    adjacency := make(map[string][]string)

    for _, node := range workflow.Nodes {
        inDegree[node.ID] = 0
        adjacency[node.ID] = []string{}
    }

    for _, edge := range workflow.Edges {
        inDegree[edge.Target]++
        adjacency[edge.Source] = append(adjacency[edge.Source], edge.Target)
    }

    // Kahn 算法
    queue := []string{}
    for nodeID, degree := range inDegree {
        if degree == 0 {
            queue = append(queue, nodeID)
        }
    }

    result := []string{}
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]
        result = append(result, current)

        for _, neighbor := range adjacency[current] {
            inDegree[neighbor]--
            if inDegree[neighbor] == 0 {
                queue = append(queue, neighbor)
            }
        }
    }

    // 检测循环依赖
    if len(result) != len(workflow.Nodes) {
        return nil, errors.New("circular dependency detected in workflow")
    }

    return result, nil
}
```

#### 变量解析

```go
func (e *WorkflowEngine) resolveVariables(
    config map[string]interface{},
    nodeResults map[string]*NodeExecutionResult,
) map[string]interface{} {
    resolved := make(map[string]interface{})

    for key, value := range config {
        switch v := value.(type) {
        case string:
            resolved[key] = e.resolveString(v, nodeResults)
        case map[string]interface{}:
            resolved[key] = e.resolveVariables(v, nodeResults)
        default:
            resolved[key] = value
        }
    }

    return resolved
}

func (e *WorkflowEngine) resolveString(
    str string,
    nodeResults map[string]*NodeExecutionResult,
) string {
    // 匹配 {{node_id.field}} 或 {{node_id.output.nested}}
    re := regexp.MustCompile(`\{\{([^}]+)\}\}`)

    return re.ReplaceAllStringFunc(str, func(match string) string {
        expr := strings.TrimSpace(match[2 : len(match)-2])
        parts := strings.Split(expr, ".")

        if len(parts) < 2 {
            return match
        }

        nodeID := parts[0]
        result, exists := nodeResults[nodeID]
        if !exists {
            return match
        }

        return e.getNestedValue(result, parts[1:])
    })
}
```

#### 条件评估

```go
func (e *WorkflowEngine) evaluateCondition(
    condition Condition,
    sourceResult *NodeExecutionResult,
    nodeResults map[string]*NodeExecutionResult,
) bool {
    switch condition.Type {
    case "always":
        return true

    case "never":
        return false

    case "on_success":
        return sourceResult.Success

    case "on_failure":
        return !sourceResult.Success

    case "expression":
        // 解析变量
        expr := e.resolveString(condition.Expression, nodeResults)
        // 评估表达式（使用 govaluate 或简单实现）
        return e.evaluateExpression(expr)

    default:
        return true
    }
}
```

### 5.2 执行流程

```go
func (e *WorkflowEngine) Execute(
    ctx context.Context,
    workflow *WorkflowDefinition,
    taskID string,
    userID string,
) (*WorkflowExecution, error) {
    execution := &WorkflowExecution{
        ID:          generateID(),
        TaskID:      taskID,
        StartTime:   time.Now(),
        Status:      "running",
        NodeResults: make(map[string]*NodeExecutionResult),
    }

    // 1. 拓扑排序
    order, err := e.topologicalSort(workflow)
    if err != nil {
        return nil, err
    }

    // 2. 按顺序执行节点
    for _, nodeID := range order {
        node := e.getNode(workflow, nodeID)

        // 检查前置条件
        if !e.shouldExecute(node, workflow, execution.NodeResults) {
            continue
        }

        // 解析变量
        config := e.resolveVariables(node.Config, execution.NodeResults)

        // 执行工具
        result := e.executeNode(ctx, node, config, userID)
        execution.NodeResults[nodeID] = result
        execution.ExecutionOrder = append(execution.ExecutionOrder, nodeID)

        // 错误处理
        if !result.Success && !workflow.Settings.ContinueOnError {
            execution.Status = "failed"
            execution.Error = result.Error
            break
        }
    }

    // 3. 汇总结果
    execution.EndTime = time.Now()
    execution.DurationMs = execution.EndTime.Sub(execution.StartTime).Milliseconds()

    if execution.Status == "running" {
        execution.Status = e.determineStatus(execution.NodeResults)
    }

    return execution, nil
}
```

---

## 6. 前端可视化编辑器

### 6.1 技术选型

推荐使用 **Vue Flow**（React Flow 的 Vue 版本）：

- 官网：https://vueflow.dev/
- GitHub：https://github.com/bcakmakoglu/vue-flow
- 特性：拖拽节点、连线、缩放、自定义节点样式

### 6.2 替代方案

| 库名 | 优点 | 缺点 | 推荐度 |
|-----|------|------|--------|
| **Vue Flow** | Vue 3 原生、性能好、文档完善 | 社区较小 | ⭐⭐⭐⭐⭐ |
| Rete.js | 功能强大、插件丰富 | Vue 3 支持一般 | ⭐⭐⭐⭐ |
| GoJS | 企业级、功能完整 | 商业收费 | ⭐⭐⭐ |
| jsPlumb | 轻量级 | 需要自己实现很多功能 | ⭐⭐ |
| 自研 SVG | 完全可控 | 开发成本高 | ⭐ |

### 6.3 Vue Flow 实现方案

#### 安装依赖

```bash
cd web
pnpm add @vue-flow/core @vue-flow/background @vue-flow/controls @vue-flow/minimap
```

#### 组件结构

```
web/src/pages/Workflows/
├── index.vue                    # 工作流列表页
├── components/
│   ├── WorkflowEditor.vue       # 主编辑器
│   ├── NodePalette.vue          # 工具面板（左侧）
│   ├── CustomToolNode.vue       # 自定义节点组件
│   ├── NodeConfigPanel.vue      # 节点配置面板（右侧）
│   └── EdgeConfigModal.vue      # 边配置弹窗
```

#### 核心代码示例

```vue
<!-- WorkflowEditor.vue -->
<template>
  <div class="workflow-editor h-screen flex">
    <!-- 左侧工具面板 -->
    <NodePalette
      :tools="availableTools"
      @drag-start="handleDragStart"
    />

    <!-- 中间画布 -->
    <div class="flex-1 relative">
      <VueFlow
        v-model:nodes="nodes"
        v-model:edges="edges"
        :node-types="nodeTypes"
        @node-click="handleNodeClick"
        @edge-click="handleEdgeClick"
        @connect="handleConnect"
        @drop="handleDrop"
        @drag-over="handleDragOver"
        fit-view-on-init
      >
        <!-- 背景网格 -->
        <Background variant="dots" :gap="20" />

        <!-- 控制按钮 -->
        <Controls />

        <!-- 小地图 -->
        <MiniMap />
      </VueFlow>

      <!-- 顶部工具栏 -->
      <div class="absolute top-4 left-1/2 transform -translate-x-1/2 z-10">
        <div class="bg-white shadow-lg rounded-lg px-4 py-2 flex gap-2">
          <button @click="handleSave" class="btn-primary">
            保存工作流
          </button>
          <button @click="handleTest" class="btn-secondary">
            测试执行
          </button>
          <button @click="handleZoomIn" class="btn-ghost">
            <ZoomIn :size="20" />
          </button>
          <button @click="handleZoomOut" class="btn-ghost">
            <ZoomOut :size="20" />
          </button>
        </div>
      </div>
    </div>

    <!-- 右侧配置面板 -->
    <NodeConfigPanel
      v-if="selectedNode"
      :node="selectedNode"
      :tools="availableTools"
      @update="handleNodeUpdate"
      @close="selectedNode = null"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, markRaw } from 'vue'
import { VueFlow, useVueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { MiniMap } from '@vue-flow/minimap'
import { ZoomIn, ZoomOut } from 'lucide-vue-next'
import CustomToolNode from './CustomToolNode.vue'
import NodePalette from './NodePalette.vue'
import NodeConfigPanel from './NodeConfigPanel.vue'
import { message } from '@/utils/message'
import * as workflowApi from '@/api/workflow'

// Vue Flow 实例
const { addNodes, addEdges, fitView, zoomIn, zoomOut, project } = useVueFlow()

// 节点和边数据
const nodes = ref([])
const edges = ref([])

// 自定义节点类型
const nodeTypes = {
  customTool: markRaw(CustomToolNode)
}

// 可用工具列表
const availableTools = ref([
  { code: 'http_request', name: 'HTTP 请求', icon: 'globe', color: '#3b82f6' },
  { code: 'email_sender', name: '邮件发送', icon: 'mail', color: '#10b981' },
  { code: 'health_checker', name: '健康检查', icon: 'heart-pulse', color: '#f59e0b' }
])

// 选中的节点
const selectedNode = ref(null)

// 拖拽开始
const handleDragStart = (event: DragEvent, tool: any) => {
  event.dataTransfer!.effectAllowed = 'move'
  event.dataTransfer!.setData('application/vueflow', JSON.stringify(tool))
}

// 拖拽释放
const handleDrop = (event: DragEvent) => {
  event.preventDefault()

  const tool = JSON.parse(event.dataTransfer!.getData('application/vueflow'))
  const position = project({ x: event.clientX, y: event.clientY })

  const newNode = {
    id: `node_${Date.now()}`,
    type: 'customTool',
    position,
    data: {
      tool_code: tool.code,
      name: tool.name,
      icon: tool.icon,
      color: tool.color,
      config: {}
    }
  }

  addNodes([newNode])
  message.success(`已添加 ${tool.name} 节点`)
}

const handleDragOver = (event: DragEvent) => {
  event.preventDefault()
  event.dataTransfer!.dropEffect = 'move'
}

// 连接节点
const handleConnect = (connection: any) => {
  const newEdge = {
    id: `edge_${Date.now()}`,
    source: connection.source,
    target: connection.target,
    type: 'smoothstep',
    animated: true,
    label: '始终执行',
    data: {
      condition: { type: 'always' }
    }
  }

  addEdges([newEdge])
}

// 点击节点
const handleNodeClick = (event: any) => {
  selectedNode.value = nodes.value.find(n => n.id === event.node.id)
}

// 更新节点配置
const handleNodeUpdate = (updatedNode: any) => {
  const index = nodes.value.findIndex(n => n.id === updatedNode.id)
  if (index !== -1) {
    nodes.value[index] = { ...nodes.value[index], ...updatedNode }
  }
}

// 保存工作流
const handleSave = async () => {
  try {
    const workflow = {
      version: '1.0',
      nodes: nodes.value.map(n => ({
        id: n.id,
        name: n.data.name,
        tool_code: n.data.tool_code,
        position: n.position,
        config: n.data.config
      })),
      edges: edges.value.map(e => ({
        id: e.id,
        source: e.source,
        target: e.target,
        label: e.label,
        condition: e.data?.condition || { type: 'always' }
      })),
      settings: {
        continue_on_error: false,
        timeout: 300
      }
    }

    await workflowApi.saveWorkflow(workflow)
    message.success('工作流保存成功')
  } catch (error: any) {
    message.error('保存失败：' + error.message)
  }
}

// 测试执行
const handleTest = async () => {
  message.info('测试功能开发中...')
}

// 缩放控制
const handleZoomIn = () => zoomIn()
const handleZoomOut = () => zoomOut()
</script>

<style>
/* Vue Flow 样式 */
@import '@vue-flow/core/dist/style.css';
@import '@vue-flow/core/dist/theme-default.css';
@import '@vue-flow/controls/dist/style.css';
@import '@vue-flow/minimap/dist/style.css';

.workflow-editor {
  background: #f3f4f6;
}
</style>
```

#### 自定义节点组件

```vue
<!-- CustomToolNode.vue -->
<template>
  <div
    :class="[
      'custom-tool-node',
      'bg-white rounded-lg shadow-lg border-2 p-3 min-w-[180px]',
      selected ? 'border-blue-500' : 'border-gray-200'
    ]"
    :style="{ borderColor: data.color }"
  >
    <!-- 节点头部 -->
    <div class="flex items-center gap-2 mb-2">
      <div
        class="w-8 h-8 rounded flex items-center justify-center text-white"
        :style="{ backgroundColor: data.color }"
      >
        <component :is="getIcon(data.icon)" :size="18" />
      </div>
      <div class="flex-1">
        <div class="font-semibold text-sm text-gray-800">
          {{ data.name }}
        </div>
        <div class="text-xs text-gray-500">
          {{ data.tool_code }}
        </div>
      </div>
    </div>

    <!-- 配置状态 -->
    <div class="flex items-center gap-1 text-xs">
      <CheckCircle
        v-if="isConfigured"
        :size="14"
        class="text-green-500"
      />
      <AlertCircle
        v-else
        :size="14"
        class="text-amber-500"
      />
      <span class="text-gray-600">
        {{ isConfigured ? '已配置' : '未配置' }}
      </span>
    </div>

    <!-- 连接点 -->
    <Handle type="target" position="left" />
    <Handle type="source" position="right" />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Handle } from '@vue-flow/core'
import {
  Globe, Mail, HeartPulse,
  CheckCircle, AlertCircle
} from 'lucide-vue-next'

interface Props {
  id: string
  data: {
    tool_code: string
    name: string
    icon: string
    color: string
    config: any
  }
  selected?: boolean
}

const props = defineProps<Props>()

const isConfigured = computed(() => {
  return Object.keys(props.data.config).length > 0
})

const getIcon = (iconName: string) => {
  const icons: Record<string, any> = {
    'globe': Globe,
    'mail': Mail,
    'heart-pulse': HeartPulse
  }
  return icons[iconName] || Globe
}
</script>

<style scoped>
.custom-tool-node {
  cursor: pointer;
  transition: all 0.2s;
}

.custom-tool-node:hover {
  box-shadow: 0 10px 25px -5px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}
</style>
```

#### 工具面板

```vue
<!-- NodePalette.vue -->
<template>
  <div class="w-64 bg-white border-r border-gray-200 p-4">
    <h3 class="text-lg font-semibold mb-4 text-gray-800">
      工具箱
    </h3>

    <div class="space-y-2">
      <div
        v-for="tool in tools"
        :key="tool.code"
        :draggable="true"
        @dragstart="handleDragStart($event, tool)"
        :class="[
          'p-3 rounded-lg cursor-move',
          'border-2 border-gray-200',
          'hover:border-gray-300 hover:shadow-md',
          'transition-all duration-200'
        ]"
      >
        <div class="flex items-center gap-2">
          <div
            class="w-8 h-8 rounded flex items-center justify-center text-white"
            :style="{ backgroundColor: tool.color }"
          >
            <component :is="getIcon(tool.icon)" :size="18" />
          </div>
          <div>
            <div class="font-medium text-sm text-gray-800">
              {{ tool.name }}
            </div>
            <div class="text-xs text-gray-500">
              拖拽到画布添加
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 帮助提示 -->
    <div class="mt-6 p-3 bg-blue-50 rounded-lg text-xs text-blue-800">
      <p class="font-semibold mb-1">💡 使用提示</p>
      <ul class="space-y-1 list-disc list-inside">
        <li>拖拽工具到画布创建节点</li>
        <li>连接节点创建工作流</li>
        <li>点击节点配置参数</li>
      </ul>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Globe, Mail, HeartPulse } from 'lucide-vue-next'

interface Tool {
  code: string
  name: string
  icon: string
  color: string
}

interface Props {
  tools: Tool[]
}

defineProps<Props>()

const emit = defineEmits<{
  'drag-start': [event: DragEvent, tool: Tool]
}>()

const handleDragStart = (event: DragEvent, tool: Tool) => {
  emit('drag-start', event, tool)
}

const getIcon = (iconName: string) => {
  const icons: Record<string, any> = {
    'globe': Globe,
    'mail': Mail,
    'heart-pulse': HeartPulse
  }
  return icons[iconName] || Globe
}
</script>
```

### 6.4 界面效果

```
┌────────────────────────────────────────────────────────────┐
│  [保存工作流] [测试执行] [+] [-]                            │
├──────────┬────────────────────────────────────┬────────────┤
│          │                                    │            │
│ 工具箱   │          画布区域                  │ 配置面板   │
│          │                                    │            │
│ ┌──────┐│   ┌────────┐      ┌────────┐      │┌─────────┐│
│ │ HTTP ││   │ 健康   │─────→│ 邮件   │      ││ 节点名称││
│ │ 请求 ││   │ 检查   │ 失败 │ 发送   │      ││         ││
│ └──────┘│   └────────┘      └────────┘      ││ 工具类型││
│          │                                    ││         ││
│ ┌──────┐│                                    ││ 配置参数││
│ │ 邮件 ││                                    ││ [...]   ││
│ │ 发送 ││                                    │└─────────┘│
│ └──────┘│                                    │            │
│          │                                    │            │
│ ┌──────┐│                                    │            │
│ │ 健康 ││                                    │            │
│ │ 检查 ││                                    │            │
│ └──────┘│                                    │            │
└──────────┴────────────────────────────────────┴────────────┘
```

---

## 7. 实现路线图

### 7.1 第一阶段：MVP（2-3 周）

**目标**：基本工作流功能可用

#### 后端（1.5 周）
- ✅ 工作流数据模型和数据库表
- ✅ WorkflowEngine 核心执行引擎
- ✅ 变量解析器（支持 `{{node.field}}`）
- ✅ 条件评估器（支持 always/on_success/on_failure）
- ✅ 修改 Task 执行逻辑，支持工作流
- ✅ 工作流 CRUD API

#### 前端（1 周）
- ✅ 安装和配置 Vue Flow
- ✅ 工作流编辑器基础界面
- ✅ 节点拖拽和连接
- ✅ 自定义节点样式
- ✅ 节点配置面板（复用 ToolConfigDrawer）
- ✅ 保存和加载工作流

#### 测试（0.5 周）
- ✅ 单元测试（执行引擎、变量解析）
- ✅ 集成测试（完整工作流执行）
- ✅ 前端 E2E 测试

**交付物**：
- 可以创建简单的顺序工作流
- 支持条件分支（成功/失败）
- 可以保存和执行工作流

### 7.2 第二阶段：增强功能（2-3 周）

#### 新增功能
- ✅ 复杂条件表达式（使用 govaluate）
- ✅ 并行执行支持
- ✅ 错误重试机制
- ✅ 工作流执行详细日志
- ✅ 可视化执行流程（高亮当前执行节点）
- ✅ 工作流模板功能

#### 前端增强
- ✅ 边条件配置弹窗
- ✅ 工作流缩略图预览
- ✅ 执行历史和调试
- ✅ 节点搜索和分类
- ✅ 快捷键支持

### 7.3 第三阶段：高级特性（长期）

- ✅ 循环和迭代节点
- ✅ 子工作流支持
- ✅ Webhook 触发器
- ✅ 工作流版本管理
- ✅ 工作流市场（分享和导入）
- ✅ A/B 测试功能
- ✅ 智能推荐（AI 辅助）

---

## 8. 使用示例

### 8.1 场景 1：网站监控 + 邮件告警

```json
{
  "version": "1.0",
  "name": "网站健康监控",
  "nodes": [
    {
      "id": "check",
      "name": "检查网站",
      "tool_code": "health_checker",
      "config": {
        "url": "https://example.com",
        "method": "GET",
        "timeout": 10,
        "expected_status": 200
      }
    },
    {
      "id": "alert",
      "name": "发送告警",
      "tool_code": "email_sender",
      "config": {
        "to": "admin@company.com",
        "subject": "⚠️ 网站告警",
        "body": "网站 {{check.output.url}} 检查失败\n错误：{{check.message}}"
      }
    }
  ],
  "edges": [
    {
      "source": "check",
      "target": "alert",
      "condition": { "type": "on_failure" }
    }
  ]
}
```

### 8.2 场景 2：API 监控 + 多级告警

```json
{
  "version": "1.0",
  "name": "API 多级告警",
  "nodes": [
    {
      "id": "api_check",
      "name": "检查 API",
      "tool_code": "http_request",
      "config": {
        "url": "https://api.example.com/health",
        "method": "GET"
      }
    },
    {
      "id": "email_admin",
      "name": "邮件通知管理员",
      "tool_code": "email_sender",
      "config": {
        "to": "admin@company.com",
        "subject": "严重告警 - API 不可用"
      }
    },
    {
      "id": "email_team",
      "name": "邮件通知团队",
      "tool_code": "email_sender",
      "config": {
        "to": "team@company.com",
        "subject": "API 监控告警"
      }
    }
  ],
  "edges": [
    {
      "source": "api_check",
      "target": "email_admin",
      "condition": {
        "type": "expression",
        "expression": "{{api_check.output.status_code}} >= 500"
      }
    },
    {
      "source": "api_check",
      "target": "email_team",
      "condition": {
        "type": "expression",
        "expression": "{{api_check.output.status_code}} >= 400 && {{api_check.output.status_code}} < 500"
      }
    }
  ]
}
```

### 8.3 场景 3：定期检查 + 成功通知

```json
{
  "version": "1.0",
  "name": "每日健康报告",
  "nodes": [
    {
      "id": "check",
      "name": "检查所有服务",
      "tool_code": "health_checker",
      "config": {
        "url": "https://status.example.com",
        "check_ssl": true
      }
    },
    {
      "id": "report",
      "name": "发送日报",
      "tool_code": "email_sender",
      "config": {
        "to": "team@company.com",
        "subject": "✅ 每日健康报告",
        "body": "所有服务运行正常\n\nSSL 证书剩余：{{check.output.ssl.days_to_expiry}} 天\n响应时间：{{check.output.response_time}}ms"
      }
    }
  ],
  "edges": [
    {
      "source": "check",
      "target": "report",
      "condition": { "type": "on_success" }
    }
  ]
}
```

---

## 9. 技术难点和解决方案

### 9.1 循环依赖检测

**问题**：用户可能创建循环引用的工作流

**解决方案**：
- 在保存时进行拓扑排序验证
- 前端实时检测循环（Vue Flow 提供工具）

### 9.2 变量作用域

**问题**：节点变量命名冲突

**解决方案**：
- 使用节点 ID 作为命名空间（`{{node_id.field}}`）
- 提供全局变量（`{{timestamp}}`, `{{user_id}}`）

### 9.3 长时间运行

**问题**：工作流可能运行很久

**解决方案**：
- 使用 Goroutine 异步执行
- 支持超时控制
- 提供取消机制

### 9.4 错误传播

**问题**：节点失败如何影响后续节点

**解决方案**：
- 提供 `continue_on_error` 配置
- 条件分支支持 `on_failure`
- 完整的错误堆栈记录

---

## 10. 后续扩展

### 10.1 高级节点类型

- **条件节点**：if-else 分支
- **循环节点**：遍历数组
- **聚合节点**：等待多个分支
- **延迟节点**：等待一段时间
- **转换节点**：数据格式转换

### 10.2 触发器扩展

- Webhook 触发
- 文件监控触发
- 数据库变更触发
- 消息队列触发

### 10.3 AI 增强

- 智能工作流推荐
- 自动优化执行路径
- 异常检测和预警

---

## 11. 总结

这个工作流系统设计：

✅ **渐进式**：从简单到复杂，分阶段实施
✅ **易用性**：可视化编辑，拖拽式操作
✅ **可扩展**：架构支持未来高级功能
✅ **向后兼容**：不影响现有单工具任务

### 关键优势

1. **Vue Flow 成熟稳定**：开箱即用的拖拽、连线功能
2. **插件化架构**：工具系统无缝集成
3. **模板语法简单**：用户易于理解和使用
4. **执行引擎高效**：Go 语言性能保证

### 实施建议

**推荐从第一阶段 MVP 开始**：
1. 后端实现工作流引擎核心功能
2. 前端使用 Vue Flow 实现可视化编辑
3. 先支持简单的顺序和条件分支
4. 逐步添加高级功能

---

**文档版本**：v1.0
**最后更新**：2025-01-12
**反馈渠道**：GitHub Issues
