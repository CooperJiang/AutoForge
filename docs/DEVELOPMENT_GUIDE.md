# AutoForge 开发指南

> **完整的项目开发文档 - 适合新开发者快速上手**

---

## 📋 目录

- [项目概述](#项目概述)
- [技术架构](#技术架构)
- [项目结构](#项目结构)
- [开发环境搭建](#开发环境搭建)
- [前后端分工](#前后端分工)
- [开发规范](#开发规范)
- [常见任务](#常见任务)
- [部署流程](#部署流程)
- [问题排查](#问题排查)

---

## 项目概述

### 项目定位
AutoForge 是一个**强大的自动化任务调度平台**，支持：
- ⏰ 定时任务调度（Cron、间隔、每日/周/月）
- 🔧 可视化工作流编排（Vue Flow）
- 🛠️ 插件化工具系统（HTTP、邮件、健康检查等）
- 📊 任务执行监控和日志管理
- 👥 多用户系统和权限管理

### 核心特性
- **工作流引擎**：可视化拖拽式工作流编排
- **插件化架构**：易于扩展新工具
- **高性能**：Go 后端 + Vue 3 前端
- **容器化部署**：支持 Docker 一键部署

### 项目历史
- **v1.0** (2024-10): 基础定时任务系统
- **v2.0** (2025-01): 工作流引擎、组件库重构

---

## 技术架构

### 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                     用户界面层 (Web)                      │
│         Vue 3 + TypeScript + Vite + Tailwind CSS        │
└────────────────────┬────────────────────────────────────┘
                     │ HTTP/REST API
┌────────────────────▼────────────────────────────────────┐
│                    应用服务层 (Backend)                   │
│           Gin Web Framework + 业务逻辑                    │
├─────────────────────────────────────────────────────────┤
│                    插件工具层 (Tools)                     │
│       HTTP请求 │ 邮件发送 │ 健康检查 │ 自定义工具          │
├─────────────────────────────────────────────────────────┤
│                    调度引擎层 (Scheduler)                 │
│                  Cron v3 定时调度引擎                     │
├─────────────────────────────────────────────────────────┤
│              工作流引擎层 (Workflow Engine)               │
│            DAG执行 │ 条件判断 │ 延迟控制                  │
└────────────────────┬────────────────────────────────────┘
                     │
┌────────────────────▼────────────────────────────────────┐
│                   数据持久化层 (Storage)                  │
│              MySQL/SQLite + GORM + Redis(可选)           │
└─────────────────────────────────────────────────────────┘
```

### 技术栈详情

#### 后端 (Backend)
| 技术 | 版本 | 用途 | 文档 |
|------|------|------|------|
| **Go** | 1.21+ | 主编程语言 | [golang.org](https://golang.org) |
| **Gin** | v1.9+ | Web 框架 | [gin-gonic.com](https://gin-gonic.com) |
| **GORM** | v1.25+ | ORM 数据库操作 | [gorm.io](https://gorm.io) |
| **Cron v3** | v3.0+ | 定时任务调度 | [robfig/cron](https://github.com/robfig/cron) |
| **JWT** | - | 身份认证 | [jwt.io](https://jwt.io) |
| **Viper** | - | 配置管理 | [spf13/viper](https://github.com/spf13/viper) |
| **Zap** | - | 结构化日志 | [uber-go/zap](https://github.com/uber-go/zap) |

#### 前端 (Frontend)
| 技术 | 版本 | 用途 | 文档 |
|------|------|------|------|
| **Vue 3** | 3.5+ | 前端框架 | [vuejs.org](https://vuejs.org) |
| **TypeScript** | 5.8+ | 类型安全 | [typescriptlang.org](https://www.typescriptlang.org) |
| **Vite** | 6.2+ | 构建工具 | [vitejs.dev](https://vitejs.dev) |
| **Pinia** | 3.0+ | 状态管理 | [pinia.vuejs.org](https://pinia.vuejs.org) |
| **Vue Router** | 4.5+ | 路由管理 | [router.vuejs.org](https://router.vuejs.org) |
| **Tailwind CSS** | 3.4+ | 原子化 CSS | [tailwindcss.com](https://tailwindcss.com) |
| **Vue Flow** | 1.47+ | 工作流可视化 | [vueflow.dev](https://vueflow.dev) |
| **Lucide** | - | 图标库 | [lucide.dev](https://lucide.dev) |
| **Axios** | 1.9+ | HTTP 客户端 | [axios-http.com](https://axios-http.com) |

#### 数据库
| 技术 | 版本 | 用途 |
|------|------|------|
| **MySQL** | 8.0+ | 生产数据库 |
| **SQLite** | 3.x | 开发/轻量级部署 |
| **Redis** | 7.0+ | 缓存（可选） |

---

## 项目结构

### 完整目录树

```
AutoForge/
├── cmd/                          # 应用程序入口
│   └── main.go                   # 主程序入口
│
├── internal/                     # 内部代码（不对外暴露）
│   ├── controllers/              # HTTP 控制器层
│   │   ├── auth/                 # 认证控制器
│   │   │   ├── login.go          # 登录接口
│   │   │   ├── register.go       # 注册接口
│   │   │   └── oauth.go          # OAuth2 第三方登录
│   │   ├── task/                 # 任务控制器
│   │   │   ├── create.go         # 创建任务
│   │   │   ├── list.go           # 任务列表
│   │   │   ├── update.go         # 更新任务
│   │   │   ├── delete.go         # 删除任务
│   │   │   └── execute.go        # 手动执行
│   │   ├── workflow/             # 工作流控制器
│   │   │   ├── create.go         # 创建工作流
│   │   │   ├── editor.go         # 编辑工作流
│   │   │   └── execute.go        # 执行工作流
│   │   ├── admin/                # 管理员控制器
│   │   │   ├── users.go          # 用户管理
│   │   │   ├── tasks.go          # 任务管理
│   │   │   └── stats.go          # 统计数据
│   │   └── tool/                 # 工具控制器
│   │       ├── list.go           # 工具列表
│   │       └── test.go           # 测试工具配置
│   │
│   ├── services/                 # 业务逻辑层
│   │   ├── taskService/          # 任务服务
│   │   │   ├── service.go        # 任务业务逻辑
│   │   │   └── execution.go      # 任务执行逻辑
│   │   ├── workflowService/      # 工作流服务
│   │   │   ├── service.go        # 工作流业务逻辑
│   │   │   ├── engine.go         # 工作流引擎
│   │   │   └── executor.go       # 节点执行器
│   │   ├── authService/          # 认证服务
│   │   │   ├── jwt.go            # JWT 令牌处理
│   │   │   └── oauth.go          # OAuth2 处理
│   │   └── cronService/          # 调度服务
│   │       ├── scheduler.go      # 调度器
│   │       └── job.go            # 任务包装器
│   │
│   ├── models/                   # 数据模型层
│   │   ├── user.go               # 用户模型
│   │   ├── task.go               # 任务模型
│   │   ├── workflow.go           # 工作流模型
│   │   ├── execution.go          # 执行记录模型
│   │   └── tool_config.go        # 工具配置模型
│   │
│   ├── routes/                   # 路由定义
│   │   ├── routes.go             # 路由注册
│   │   ├── api.go                # API 路由组
│   │   └── admin.go              # 管理员路由组
│   │
│   ├── middleware/               # 中间件
│   │   ├── auth.go               # 认证中间件
│   │   ├── cors.go               # CORS 中间件
│   │   ├── logger.go             # 日志中间件
│   │   ├── recovery.go           # 异常恢复
│   │   └── rate_limit.go         # 限流中间件
│   │
│   ├── cron/                     # 定时任务调度器
│   │   ├── scheduler.go          # Cron 调度器实现
│   │   └── manager.go            # 任务管理器
│   │
│   └── dto/                      # 数据传输对象
│       ├── request/              # 请求 DTO
│       │   ├── taskRequest.go    # 任务请求
│       │   └── authRequest.go    # 认证请求
│       └── response/             # 响应 DTO
│           ├── taskResponse.go   # 任务响应
│           └── userResponse.go   # 用户响应
│
├── pkg/                          # 公共包（可对外暴露）
│   ├── config/                   # 配置管理
│   │   ├── config.go             # 配置结构定义
│   │   └── loader.go             # 配置加载器
│   │
│   ├── database/                 # 数据库连接
│   │   ├── mysql.go              # MySQL 连接
│   │   ├── sqlite.go             # SQLite 连接
│   │   └── migrations.go         # 数据库迁移
│   │
│   ├── logger/                   # 日志工具
│   │   └── logger.go             # 日志初始化
│   │
│   ├── errors/                   # 错误处理
│   │   ├── codes.go              # 错误码定义
│   │   └── errors.go             # 错误包装
│   │
│   ├── common/                   # 公共工具
│   │   ├── response.go           # 统一响应格式
│   │   ├── crypto.go             # 加密工具
│   │   └── validator.go          # 参数验证
│   │
│   └── utools/                   # 工具系统核心
│       ├── base.go               # 工具基类
│       ├── registry.go           # 工具注册表
│       ├── schema.go             # 配置 Schema
│       ├── http/                 # HTTP 请求工具
│       │   └── http_tool.go      # HTTP 工具实现
│       ├── email/                # 邮件发送工具
│       │   └── email_tool.go     # 邮件工具实现
│       └── health/               # 健康检查工具
│           └── health_tool.go    # 健康检查实现
│
├── web/                          # 前端项目
│   ├── src/
│   │   ├── main.ts               # 前端入口
│   │   │
│   │   ├── pages/                # 页面组件
│   │   │   ├── Home/             # 首页
│   │   │   │   └── index.vue
│   │   │   ├── Login/            # 登录页
│   │   │   │   └── index.vue
│   │   │   ├── Register/         # 注册页
│   │   │   │   └── index.vue
│   │   │   ├── Tasks/            # 任务管理页面
│   │   │   │   ├── index.vue     # 任务列表
│   │   │   │   └── components/   # 任务相关组件
│   │   │   │       ├── TaskCard.vue
│   │   │   │       ├── TaskDrawer.vue
│   │   │   │       └── ToolConfigDrawer.vue
│   │   │   ├── Workflows/        # 工作流管理页面
│   │   │   │   ├── index.vue     # 工作流列表
│   │   │   │   ├── editor.vue    # 工作流编辑器
│   │   │   │   └── components/   # 工作流组件
│   │   │   │       ├── WorkflowCard.vue
│   │   │   │       ├── NodeConfigDrawer.vue
│   │   │   │       └── EnvVarManager.vue
│   │   │   ├── Admin/            # 管理后台
│   │   │   │   ├── Dashboard.vue # 仪表盘
│   │   │   │   ├── Users.vue     # 用户管理
│   │   │   │   ├── Tasks.vue     # 任务管理
│   │   │   │   └── Executions.vue# 执行记录
│   │   │   ├── Profile/          # 个人中心
│   │   │   │   └── index.vue
│   │   │   └── Settings/         # 设置页面
│   │   │       └── index.vue
│   │   │
│   │   ├── components/           # 通用组件库
│   │   │   ├── index.ts          # 组件统一导出
│   │   │   ├── README.md         # 组件库文档
│   │   │   ├── COMPONENT_DESIGN_SPEC.md  # 设计规范
│   │   │   │
│   │   │   ├── BaseButton/       # 基础按钮
│   │   │   │   ├── index.vue
│   │   │   │   └── index.ts      # 导出文件
│   │   │   ├── BaseInput/        # 基础输入框
│   │   │   │   ├── index.vue
│   │   │   │   └── index.ts
│   │   │   ├── BaseSelect/       # 下拉选择
│   │   │   │   ├── index.vue
│   │   │   │   └── index.ts
│   │   │   ├── Dialog/           # 对话框
│   │   │   │   ├── index.vue
│   │   │   │   └── index.ts
│   │   │   ├── Drawer/           # 侧边抽屉
│   │   │   │   ├── index.vue
│   │   │   │   └── index.ts
│   │   │   ├── Table/            # 数据表格
│   │   │   │   ├── index.vue
│   │   │   │   └── index.ts
│   │   │   ├── Pagination/       # 分页组件
│   │   │   │   ├── index.vue
│   │   │   │   └── index.ts
│   │   │   └── ...               # 其他 21 个组件
│   │   │
│   │   ├── api/                  # API 接口封装
│   │   │   ├── request.ts        # Axios 配置
│   │   │   ├── auth.ts           # 认证接口
│   │   │   ├── task.ts           # 任务接口
│   │   │   ├── workflow.ts       # 工作流接口
│   │   │   ├── admin.ts          # 管理员接口
│   │   │   └── tool.ts           # 工具接口
│   │   │
│   │   ├── stores/               # Pinia 状态管理
│   │   │   ├── auth.ts           # 认证状态
│   │   │   ├── task.ts           # 任务状态
│   │   │   └── workflow.ts       # 工作流状态
│   │   │
│   │   ├── router/               # 路由配置
│   │   │   ├── index.ts          # 路由主文件
│   │   │   └── guards.ts         # 路由守卫
│   │   │
│   │   ├── layouts/              # 布局组件
│   │   │   ├── DefaultLayout.vue # 默认布局
│   │   │   ├── AdminLayout.vue   # 管理后台布局
│   │   │   └── EmptyLayout.vue   # 空白布局
│   │   │
│   │   ├── utils/                # 工具函数
│   │   │   ├── message.ts        # 消息提示
│   │   │   ├── curlParser.ts     # cURL 解析器
│   │   │   ├── dateFormat.ts     # 日期格式化
│   │   │   └── validators.ts     # 表单验证
│   │   │
│   │   ├── composables/          # 组合式函数
│   │   │   ├── useAuth.ts        # 认证逻辑
│   │   │   ├── useTask.ts        # 任务逻辑
│   │   │   └── useWorkflow.ts    # 工作流逻辑
│   │   │
│   │   └── types/                # TypeScript 类型定义
│   │       ├── api.ts            # API 类型
│   │       ├── task.ts           # 任务类型
│   │       ├── workflow.ts       # 工作流类型
│   │       └── user.ts           # 用户类型
│   │
│   ├── public/                   # 静态资源
│   │   ├── logo.png
│   │   └── favicon.ico
│   │
│   ├── index.html                # HTML 模板
│   ├── vite.config.ts            # Vite 配置
│   ├── tailwind.config.js        # Tailwind 配置
│   ├── tsconfig.json             # TypeScript 配置
│   ├── package.json              # 前端依赖
│   └── pnpm-lock.yaml            # 依赖锁文件
│
├── tests/                        # 测试文件
│   ├── unit/                     # 单元测试
│   └── integration/              # 集成测试
│
├── scripts/                      # 脚本工具
│   ├── build.sh                  # 构建脚本
│   ├── deploy.sh                 # 部署脚本
│   └── init_db.sql               # 数据库初始化
│
├── docs/                         # 项目文档
│   ├── API.md                    # API 文档
│   ├── ARCHITECTURE.md           # 架构文档
│   └── DEVELOPMENT_GUIDE.md      # 开发指南（本文档）
│
├── config.yaml                   # 配置文件（运行时）
├── config.example.yaml           # 配置示例
├── config.prod.yaml              # 生产配置
├── go.mod                        # Go 模块定义
├── go.sum                        # Go 依赖锁文件
├── Makefile                      # 构建命令
├── Dockerfile                    # Docker 镜像构建
├── docker-compose.yml            # Docker Compose 配置
├── .gitignore                    # Git 忽略文件
├── .air.toml                     # 热重载配置
└── README.md                     # 项目说明
```

### 核心目录说明

#### 后端核心目录

| 目录 | 职责 | 示例 |
|------|------|------|
| `cmd/` | 应用程序入口 | `main.go` - 启动服务 |
| `internal/controllers/` | HTTP 控制器，处理请求和响应 | `task/create.go` - 创建任务接口 |
| `internal/services/` | 业务逻辑层，核心业务实现 | `taskService/execution.go` - 任务执行逻辑 |
| `internal/models/` | 数据模型，ORM 映射 | `task.go` - 任务表结构 |
| `internal/routes/` | 路由定义，URL 映射 | `routes.go` - 注册所有路由 |
| `internal/middleware/` | 中间件，请求拦截处理 | `auth.go` - JWT 认证 |
| `pkg/utools/` | 工具插件系统 | `http/http_tool.go` - HTTP 请求工具 |
| `pkg/config/` | 配置管理 | `config.go` - 配置结构 |
| `pkg/database/` | 数据库连接 | `mysql.go` - MySQL 初始化 |

#### 前端核心目录

| 目录 | 职责 | 示例 |
|------|------|------|
| `pages/` | 页面组件，路由对应的页面 | `Tasks/index.vue` - 任务列表页 |
| `components/` | 通用组件库，可复用的 UI 组件 | `BaseButton/` - 按钮组件 |
| `api/` | API 接口封装，与后端通信 | `task.ts` - 任务相关 API |
| `stores/` | Pinia 状态管理 | `auth.ts` - 用户认证状态 |
| `router/` | Vue Router 路由配置 | `index.ts` - 路由定义 |
| `layouts/` | 布局组件，页面外层布局 | `DefaultLayout.vue` - 默认布局 |
| `utils/` | 工具函数，通用逻辑 | `curlParser.ts` - cURL 解析 |
| `composables/` | 组合式函数，可复用逻辑 | `useAuth.ts` - 认证逻辑 |
| `types/` | TypeScript 类型定义 | `workflow.ts` - 工作流类型 |

---

## 开发环境搭建

### 1. 前置要求

| 工具 | 版本 | 安装命令 | 验证命令 |
|------|------|----------|----------|
| **Go** | 1.21+ | [官网下载](https://golang.org/dl/) | `go version` |
| **Node.js** | 18+ | [官网下载](https://nodejs.org/) | `node -v` |
| **pnpm** | 最新 | `npm install -g pnpm` | `pnpm -v` |
| **MySQL** | 8.0+ (可选) | [官网下载](https://dev.mysql.com/downloads/) | `mysql --version` |
| **Git** | 最新 | [官网下载](https://git-scm.com/) | `git --version` |

### 2. 克隆项目

```bash
# 克隆仓库
git clone https://github.com/CooperJiang/AutoForge.git
cd AutoForge

# 查看项目结构
tree -L 2
```

### 3. 配置数据库（可选）

#### 使用 SQLite（默认，无需配置）
项目默认使用 SQLite，数据库文件自动创建在 `data/autoforge.db`。

#### 使用 MySQL（推荐生产环境）

```bash
# 登录 MySQL
mysql -u root -p

# 创建数据库
CREATE DATABASE autoforge CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

# 创建用户（可选）
CREATE USER 'autoforge'@'localhost' IDENTIFIED BY 'your_password';
GRANT ALL PRIVILEGES ON autoforge.* TO 'autoforge'@'localhost';
FLUSH PRIVILEGES;
```

### 4. 配置文件

```bash
# 复制配置示例
cp config.example.yaml config.yaml

# 编辑配置文件（使用你喜欢的编辑器）
vim config.yaml  # 或 code config.yaml
```

**重要配置项**：

```yaml
app:
  name: "AutoForge"
  port: 7777                    # 后端端口
  mode: "debug"                 # debug 或 release
  base_url: "http://localhost:7777"

database:
  driver: "sqlite"              # mysql 或 sqlite
  # SQLite 配置
  sqlite_path: "data/autoforge.db"
  # MySQL 配置（如果使用 MySQL）
  host: "127.0.0.1"
  port: 3306
  username: "root"
  password: "your_password"
  name: "autoforge"

jwt:
  secret_key: "CHANGE_THIS_SECRET_KEY"  # ⚠️ 生产环境必须修改
  expires_in: 24                        # Token 过期时间（小时）

mail:
  enabled: true                 # 是否启用邮件功能
  host: "smtp.qq.com"          # SMTP 服务器
  port: 465                     # SMTP 端口（465=SSL, 587=TLS）
  username: "your@email.com"
  password: "your_password"     # SMTP 授权码
  from: "noreply@autoforge.com"
  from_name: "AutoForge"
  ssl: true

oauth:
  linux_do:                     # Linux.do OAuth2 登录
    enabled: false              # 是否启用
    client_id: "your_client_id"
    client_secret: "your_secret"
    redirect_url: "http://localhost:7777/api/v1/auth/oauth/callback"
```

### 5. 安装依赖

#### 后端依赖

```bash
# 下载 Go 模块
go mod download

# 验证依赖
go mod verify

# 整理依赖（可选）
go mod tidy
```

#### 前端依赖

```bash
# 进入前端目录
cd web

# 安装依赖（推荐使用 pnpm）
pnpm install

# 或使用 npm
# npm install
```

### 6. 启动开发环境

#### 方式一：使用 Makefile（推荐）

```bash
# 回到项目根目录
cd ..

# 一键启动前后端（会打开新终端窗口）
make start

# 或者后台启动
make start-bg

# 查看服务状态
make status

# 查看日志
make logs

# 停止所有服务
make stop
```

#### 方式二：手动启动

**终端 1 - 后端**：
```bash
# 启动后端（热重载）
go run cmd/main.go

# 或使用 Air 热重载（需先安装 Air）
air
```

**终端 2 - 前端**：
```bash
cd web
pnpm dev
```

### 7. 访问应用

| 服务 | 地址 | 说明 |
|------|------|------|
| **前端界面** | http://localhost:3200 | 主应用界面 |
| **后端 API** | http://localhost:7777 | RESTful API |
| **管理后台** | http://localhost:3200/admin | 管理员后台 |
| **API 文档** | http://localhost:7777/swagger | Swagger 文档（如果启用） |

### 8. 初始化账号

**首次启动时**，系统会自动创建管理员账号：

```
管理员账号：admin
默认密码：<控制台输出>
```

⚠️ **请立即登录并修改密码**！

---

## 前后端分工

### 前端开发职责

#### 核心职责
1. **页面开发**：实现所有用户界面
2. **组件开发**：构建可复用的 UI 组件
3. **状态管理**：使用 Pinia 管理应用状态
4. **路由管理**：配置页面路由和权限
5. **API 集成**：调用后端接口并处理数据
6. **交互优化**：提升用户体验和动画效果

#### 主要工作内容

##### 1. 页面开发 (`src/pages/`)
- **任务管理页面**：任务列表、创建、编辑、执行历史
- **工作流编辑器**：可视化拖拽编辑器（基于 Vue Flow）
- **管理后台**：用户管理、系统统计、执行监控
- **认证页面**：登录、注册、找回密码
- **个人中心**：个人信息、偏好设置

##### 2. 组件开发 (`src/components/`)
- **基础组件**：Button、Input、Select、Dialog、Drawer
- **业务组件**：TaskCard、WorkflowCard、NodeConfig
- **表单组件**：TimePicker、WeekDayPicker、ParamInput
- **展示组件**：Table、Pagination、JsonViewer

##### 3. API 集成 (`src/api/`)
```typescript
// 示例：任务 API
export const taskApi = {
  // 获取任务列表
  list: (params: TaskListParams) => request.get('/api/v1/tasks', { params }),

  // 创建任务
  create: (data: CreateTaskDto) => request.post('/api/v1/tasks', data),

  // 更新任务
  update: (id: number, data: UpdateTaskDto) =>
    request.put(`/api/v1/tasks/${id}`, data),

  // 删除任务
  delete: (id: number) => request.delete(`/api/v1/tasks/${id}`),

  // 执行任务
  execute: (id: number) => request.post(`/api/v1/tasks/${id}/execute`)
}
```

##### 4. 状态管理 (`src/stores/`)
```typescript
// 示例：认证状态
export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null as User | null,
    token: localStorage.getItem('token') || '',
    isAuthenticated: false
  }),

  actions: {
    async login(credentials: LoginDto) {
      const { token, user } = await authApi.login(credentials)
      this.token = token
      this.user = user
      this.isAuthenticated = true
      localStorage.setItem('token', token)
    },

    logout() {
      this.user = null
      this.token = ''
      this.isAuthenticated = false
      localStorage.removeItem('token')
    }
  }
})
```

#### 前端技术要点

| 技术点 | 说明 | 示例 |
|--------|------|------|
| **Composition API** | 使用 `<script setup>` | `const count = ref(0)` |
| **TypeScript** | 类型安全 | `interface Task { id: number }` |
| **Tailwind CSS** | 原子化样式 | `class="flex items-center gap-2"` |
| **Vue Flow** | 工作流可视化 | 拖拽节点、连线 |
| **响应式设计** | 移动端适配 | `class="md:flex-row flex-col"` |

---

### 后端开发职责

#### 核心职责
1. **API 开发**：提供 RESTful API 接口
2. **业务逻辑**：实现核心业务功能
3. **数据库设计**：设计表结构和关系
4. **任务调度**：实现 Cron 调度引擎
5. **工作流引擎**：实现 DAG 执行引擎
6. **工具开发**：开发新的工具插件

#### 主要工作内容

##### 1. API 开发 (`internal/controllers/`)

```go
// 示例：任务控制器
package task

import (
    "github.com/gin-gonic/gin"
    "auto-forge/internal/services/taskService"
)

// 创建任务
func CreateTask(c *gin.Context) {
    var req CreateTaskRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    task, err := taskService.CreateTask(c, &req)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    c.JSON(200, gin.H{"data": task})
}
```

##### 2. 业务逻辑 (`internal/services/`)

```go
// 示例：任务服务
package taskService

func CreateTask(ctx context.Context, req *CreateTaskRequest) (*models.Task, error) {
    // 1. 验证参数
    if err := validateTaskConfig(req); err != nil {
        return nil, err
    }

    // 2. 创建任务记录
    task := &models.Task{
        Name:        req.Name,
        Description: req.Description,
        ToolType:    req.ToolType,
        Config:      req.Config,
        Schedule:    req.Schedule,
        Enabled:     false,
    }

    if err := db.Create(task).Error; err != nil {
        return nil, err
    }

    // 3. 注册到调度器
    if req.Enabled {
        if err := cronService.AddTask(task); err != nil {
            return nil, err
        }
    }

    return task, nil
}
```

##### 3. 数据库设计 (`internal/models/`)

```go
// 示例：任务模型
package models

type Task struct {
    ID          uint      `gorm:"primaryKey" json:"id"`
    UserID      uint      `gorm:"not null;index" json:"user_id"`
    Name        string    `gorm:"size:100;not null" json:"name"`
    Description string    `gorm:"size:500" json:"description"`
    ToolType    string    `gorm:"size:50;not null" json:"tool_type"`
    Config      JSON      `gorm:"type:json" json:"config"`
    Schedule    JSON      `gorm:"type:json" json:"schedule"`
    Enabled     bool      `gorm:"default:false" json:"enabled"`
    NextRunAt   *time.Time `gorm:"index" json:"next_run_at"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`

    // 关联
    User       User         `gorm:"foreignKey:UserID" json:"user,omitempty"`
    Executions []Execution  `gorm:"foreignKey:TaskID" json:"executions,omitempty"`
}
```

##### 4. 工具开发 (`pkg/utools/`)

```go
// 示例：HTTP 请求工具
package http

type HTTPTool struct {
    *utools.BaseTool
}

func NewHTTPTool() *HTTPTool {
    metadata := &utools.ToolMetadata{
        Code:        "http_request",
        Name:        "HTTP 请求",
        Description: "发送 HTTP 请求",
        Category:    "网络",
        Version:     "1.0.0",
        Author:      "AutoForge",
        Icon:        "globe",
    }

    schema := &utools.ConfigSchema{
        Type: "object",
        Properties: map[string]utools.PropertySchema{
            "url": {
                Type:        "string",
                Title:       "请求地址",
                Description: "HTTP 请求的 URL",
            },
            "method": {
                Type:    "string",
                Title:   "请求方法",
                Enum:    []string{"GET", "POST", "PUT", "DELETE"},
                Default: "GET",
            },
        },
        Required: []string{"url"},
    }

    return &HTTPTool{
        BaseTool: utools.NewBaseTool(metadata, schema),
    }
}

func (t *HTTPTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
    // 解析配置
    url := config["url"].(string)
    method := config["method"].(string)

    // 发送请求
    req, _ := http.NewRequest(method, url, nil)
    client := &http.Client{Timeout: 30 * time.Second}
    resp, err := client.Do(req)
    if err != nil {
        return &utools.ExecutionResult{
            Success: false,
            Message: err.Error(),
        }, err
    }
    defer resp.Body.Close()

    // 读取响应
    body, _ := io.ReadAll(resp.Body)

    return &utools.ExecutionResult{
        Success:    resp.StatusCode >= 200 && resp.StatusCode < 300,
        Message:    fmt.Sprintf("状态码: %d", resp.StatusCode),
        Output:     map[string]interface{}{
            "status_code": resp.StatusCode,
            "body":        string(body),
        },
    }, nil
}
```

#### 后端技术要点

| 技术点 | 说明 | 示例 |
|--------|------|------|
| **Gin 路由** | HTTP 路由和中间件 | `r.POST("/tasks", CreateTask)` |
| **GORM ORM** | 数据库操作 | `db.Find(&tasks).Where("enabled = ?", true)` |
| **Cron 调度** | 定时任务 | `c.AddFunc("0 0 * * *", func() {...})` |
| **JWT 认证** | Token 鉴权 | `token, _ := jwt.ParseWithClaims(...)` |
| **Context 传递** | 请求上下文 | `ctx.Value("user_id")` |

---

### 前后端协作流程

#### 1. 需求阶段
1. **产品需求** → 拆分为前后端任务
2. **接口设计** → 确定 API 接口格式（前后端共同参与）
3. **数据模型** → 确定数据结构和字段

#### 2. 开发阶段
1. **后端优先**：后端先开发 API 并提供接口文档
2. **前端跟进**：前端根据接口文档集成 API
3. **并行开发**：后端可以使用 Mock 数据，前端可以使用 Mock API

#### 3. 联调阶段
1. **本地联调**：前后端在本地环境联调
2. **接口验证**：使用 Postman/Insomnia 验证接口
3. **问题修复**：及时沟通解决接口不一致问题

#### 4. 测试部署
1. **集成测试**：测试完整业务流程
2. **性能测试**：压力测试和性能优化
3. **部署上线**：前后端同步部署

---

## 开发规范

### 代码风格规范

#### Go 代码规范

```go
// ✅ 好的示例

// 1. 包命名：简短、小写、单数
package taskService

// 2. 函数命名：驼峰式，首字母大写表示公开
func CreateTask(ctx context.Context, req *CreateTaskRequest) (*Task, error) {
    // 3. 变量命名：驼峰式，有意义的名称
    userID := ctx.Value("user_id").(uint)

    // 4. 错误处理：立即返回错误
    if req.Name == "" {
        return nil, errors.New("任务名称不能为空")
    }

    // 5. 注释：函数和关键逻辑都要注释
    // 创建任务记录
    task := &Task{
        Name:   req.Name,
        UserID: userID,
    }

    return task, nil
}

// ❌ 不好的示例
func createtask(c context.Context, r *CreateTaskRequest) (*Task, error) {
    u := c.Value("user_id").(uint)  // 变量名太短
    if r.Name == "" {
        // 没有立即返回
    }
    t := &Task{Name: r.Name, UserID: u}
    return t, nil
}
```

**遵循规范**：
- [Effective Go](https://golang.org/doc/effective_go)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- 使用 `gofmt` 格式化代码：`go fmt ./...`
- 使用 `golangci-lint` 静态检查：`golangci-lint run`

#### TypeScript/Vue 代码规范

```vue
<!-- ✅ 好的示例 -->
<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { taskApi } from '@/api/task'
import type { Task } from '@/types/task'

// 1. Props 定义：使用 TypeScript 类型
interface Props {
  taskId: number
  editable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  editable: true
})

// 2. Emits 定义：明确事件类型
const emit = defineEmits<{
  'update:task': [task: Task]
  'delete': [taskId: number]
}>()

// 3. 响应式变量：有意义的命名
const isLoading = ref(false)
const task = ref<Task | null>(null)

// 4. 计算属性：使用 computed
const isEnabled = computed(() => task.value?.enabled ?? false)

// 5. 函数命名：动词开头
const fetchTask = async () => {
  isLoading.value = true
  try {
    const response = await taskApi.getById(props.taskId)
    task.value = response.data
  } catch (error) {
    console.error('获取任务失败:', error)
  } finally {
    isLoading.value = false
  }
}

// 6. 生命周期：清晰的逻辑
onMounted(() => {
  fetchTask()
})
</script>

<template>
  <!-- 7. 模板：语义化标签和类名 -->
  <div class="task-detail">
    <div v-if="isLoading" class="loading-spinner">
      加载中...
    </div>
    <div v-else-if="task" class="task-content">
      <h2 class="task-title">{{ task.name }}</h2>
      <p class="task-description">{{ task.description }}</p>
    </div>
  </div>
</template>

<!-- ❌ 不好的示例 -->
<script setup lang="ts">
const t = ref(null)  // 变量名太短
const l = ref(false)

function f() {  // 函数名不明确
  // ...
}
</script>
```

**遵循规范**：
- [Vue 3 风格指南](https://vuejs.org/style-guide/)
- [TypeScript 风格指南](https://google.github.io/styleguide/tsguide.html)
- 使用 ESLint：`pnpm lint`
- 使用 Prettier：`pnpm format`

---

### Git 提交规范

#### Commit Message 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

**示例**：
```bash
# 功能开发
git commit -m "feat(task): 添加任务批量删除功能"

# Bug 修复
git commit -m "fix(workflow): 修复工作流执行时节点连接丢失问题"

# 文档更新
git commit -m "docs(readme): 更新开发环境搭建说明"

# 样式调整
git commit -m "style(button): 调整按钮间距和圆角"

# 代码重构
git commit -m "refactor(api): 重构 API 请求封装逻辑"

# 性能优化
git commit -m "perf(table): 优化大数据表格渲染性能"

# 测试相关
git commit -m "test(task): 添加任务创建单元测试"
```

#### Type 类型

| Type | 说明 | 示例 |
|------|------|------|
| `feat` | 新功能 | feat(workflow): 添加工作流导入导出功能 |
| `fix` | Bug 修复 | fix(auth): 修复 JWT Token 过期未刷新问题 |
| `docs` | 文档更新 | docs(api): 更新 API 接口文档 |
| `style` | 代码格式（不影响功能） | style: 统一代码缩进为 2 空格 |
| `refactor` | 代码重构 | refactor(service): 重构任务服务层结构 |
| `perf` | 性能优化 | perf(db): 优化数据库查询索引 |
| `test` | 测试相关 | test(controller): 添加控制器集成测试 |
| `chore` | 构建/工具变动 | chore(deps): 更新依赖版本 |
| `revert` | 回滚 | revert: 回滚 feat(workflow) 提交 |

#### 分支命名规范

```bash
# 功能开发分支
feature/workflow-export      # 工作流导出功能
feature/task-batch-delete    # 任务批量删除

# Bug 修复分支
fix/auth-token-refresh       # 修复 Token 刷新问题
hotfix/critical-bug          # 紧急 Bug 修复

# 文档分支
docs/dev-guide               # 开发指南文档

# 重构分支
refactor/api-layer           # API 层重构
```

#### 工作流程

```bash
# 1. 创建新分支
git checkout -b feature/new-feature

# 2. 开发并提交
git add .
git commit -m "feat(module): 添加新功能"

# 3. 推送到远程
git push origin feature/new-feature

# 4. 创建 Pull Request
# 在 GitHub/GitLab 上创建 PR，等待 Code Review

# 5. 合并到主分支
# Code Review 通过后，合并到 main/develop 分支
```

---

### API 设计规范

#### RESTful 规范

| 操作 | 方法 | 路径 | 说明 |
|------|------|------|------|
| 列表 | GET | `/api/v1/tasks` | 获取任务列表 |
| 详情 | GET | `/api/v1/tasks/:id` | 获取单个任务详情 |
| 创建 | POST | `/api/v1/tasks` | 创建新任务 |
| 更新 | PUT | `/api/v1/tasks/:id` | 更新任务（全量） |
| 部分更新 | PATCH | `/api/v1/tasks/:id` | 更新任务（部分） |
| 删除 | DELETE | `/api/v1/tasks/:id` | 删除任务 |
| 批量删除 | DELETE | `/api/v1/tasks` | 批量删除（Body 传 IDs） |
| 执行 | POST | `/api/v1/tasks/:id/execute` | 执行任务 |
| 停止 | POST | `/api/v1/tasks/:id/stop` | 停止任务 |

#### 统一响应格式

```typescript
// 成功响应
{
  "code": 0,
  "message": "success",
  "data": {
    "id": 1,
    "name": "任务名称"
  }
}

// 列表响应
{
  "code": 0,
  "message": "success",
  "data": {
    "items": [...],
    "total": 100,
    "page": 1,
    "page_size": 20
  }
}

// 错误响应
{
  "code": 400,
  "message": "参数错误",
  "errors": {
    "name": "任务名称不能为空"
  }
}
```

#### 错误码规范

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 400 | 请求参数错误 |
| 401 | 未认证 |
| 403 | 无权限 |
| 404 | 资源不存在 |
| 409 | 资源冲突（如重复创建） |
| 500 | 服务器内部错误 |

---

### 数据库设计规范

#### 表命名规范
- 使用**小写字母**和**下划线**
- 使用**复数形式**：`tasks`、`users`、`workflows`
- 关联表使用**两个表名组合**：`user_roles`

#### 字段命名规范
- 主键统一使用 `id`
- 外键使用 `表名_id`：`user_id`、`task_id`
- 时间字段：`created_at`、`updated_at`、`deleted_at`
- 布尔字段：`is_xxx`、`has_xxx`、`enabled`

#### 索引规范
- 主键自动创建索引
- 外键必须创建索引
- 频繁查询的字段创建索引
- 组合索引遵循**最左前缀原则**

#### 示例

```sql
-- 任务表
CREATE TABLE tasks (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    name VARCHAR(100) NOT NULL COMMENT '任务名称',
    description VARCHAR(500) COMMENT '任务描述',
    tool_type VARCHAR(50) NOT NULL COMMENT '工具类型',
    config JSON COMMENT '工具配置',
    schedule JSON COMMENT '调度配置',
    enabled BOOLEAN DEFAULT FALSE COMMENT '是否启用',
    next_run_at DATETIME COMMENT '下次执行时间',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at DATETIME COMMENT '软删除时间',

    INDEX idx_user_id (user_id),
    INDEX idx_enabled_next_run (enabled, next_run_at),
    INDEX idx_created_at (created_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='任务表';
```

---

### 组件开发规范

详细规范请参考：`web/src/components/COMPONENT_DESIGN_SPEC.md` 和 `web/src/components/README.md`

#### 组件结构规范

```
ComponentName/
├── index.vue          # 组件入口文件（必需）
├── index.ts           # 导出文件（必需）
├── types.ts           # 类型定义（可选）
├── hooks.ts           # 组合式函数（可选）
├── constants.ts       # 常量定义（可选）
└── README.md          # 组件文档（推荐）
```

#### 组件代码规范

```vue
<script setup lang="ts">
// ✅ 好的组件结构

// 1. 导入区：按类型分组
// Vue 核心
import { ref, computed, watch } from 'vue'

// 第三方库
import { useRouter } from 'vue-router'

// 本地组件
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'

// 类型定义
import type { Task } from '@/types/task'

// 工具函数
import { formatDate } from '@/utils/dateFormat'

// 2. Props 定义：使用 TypeScript 接口
interface Props {
  task: Task
  editable?: boolean
  size?: 'sm' | 'md' | 'lg'
}

const props = withDefaults(defineProps<Props>(), {
  editable: true,
  size: 'md'
})

// 3. Emits 定义：明确事件类型
const emit = defineEmits<{
  'update:task': [task: Task]
  'delete': []
}>()

// 4. 响应式状态
const isEditing = ref(false)
const localTask = ref<Task>({ ...props.task })

// 5. 计算属性
const displayName = computed(() =>
  localTask.value.name || '未命名任务'
)

// 6. 方法
const handleSave = () => {
  emit('update:task', localTask.value)
  isEditing.value = false
}

// 7. 监听器
watch(() => props.task, (newTask) => {
  localTask.value = { ...newTask }
})
</script>

<template>
  <!-- 清晰的模板结构 -->
  <div class="task-card">
    <div class="task-header">
      <h3 class="task-title">{{ displayName }}</h3>
      <BaseButton
        v-if="editable"
        size="sm"
        @click="isEditing = true"
      >
        编辑
      </BaseButton>
    </div>

    <div v-if="isEditing" class="task-edit">
      <BaseInput v-model="localTask.name" placeholder="任务名称" />
      <div class="actions">
        <BaseButton @click="handleSave">保存</BaseButton>
        <BaseButton variant="ghost" @click="isEditing = false">
          取消
        </BaseButton>
      </div>
    </div>
  </div>
</template>

<style scoped>
/* 使用 Tailwind 优先，复杂样式使用 scoped CSS */
.task-card {
  @apply rounded-lg border border-slate-200 p-4;
}

.task-header {
  @apply flex items-center justify-between mb-3;
}

.task-title {
  @apply text-lg font-semibold text-slate-900;
}
</style>
```

#### 组件命名规范

| 类型 | 命名规则 | 示例 |
|------|----------|------|
| 基础组件 | `Base` 前缀 | `BaseButton`、`BaseInput` |
| 业务组件 | 功能描述 | `TaskCard`、`WorkflowEditor` |
| 布局组件 | `Layout` 后缀 | `DefaultLayout`、`AdminLayout` |
| 页面组件 | 页面名称 | `Home`、`TaskList` |

#### 组件大小限制

- **单个组件不超过 500 行**
- 超过则需要拆分为多个子组件或使用组合式函数

---

## 常见任务

### 1. 添加新页面

#### 步骤

1. **创建页面组件**

```bash
# 在 pages 目录下创建新页面
mkdir -p web/src/pages/NewFeature
touch web/src/pages/NewFeature/index.vue
```

```vue
<!-- web/src/pages/NewFeature/index.vue -->
<script setup lang="ts">
import { ref } from 'vue'

const message = ref('这是新功能页面')
</script>

<template>
  <div class="new-feature-page">
    <h1>{{ message }}</h1>
  </div>
</template>
```

2. **添加路由配置**

```typescript
// web/src/router/index.ts

const routes = [
  // ... 其他路由
  {
    path: '/new-feature',
    name: 'NewFeature',
    component: () => import('@/pages/NewFeature/index.vue'),
    meta: {
      title: '新功能',
      requiresAuth: true  // 需要登录
    }
  }
]
```

3. **添加导航链接**

```vue
<!-- 在布局组件中添加导航 -->
<template>
  <nav>
    <RouterLink to="/new-feature">新功能</RouterLink>
  </nav>
</template>
```

### 2. 添加新 API 接口

#### 后端

1. **定义路由**

```go
// internal/routes/api.go

func RegisterAPIRoutes(r *gin.Engine) {
    api := r.Group("/api/v1")

    // 新功能路由
    feature := api.Group("/features")
    {
        feature.GET("", featureController.List)
        feature.POST("", featureController.Create)
        feature.GET("/:id", featureController.Get)
        feature.PUT("/:id", featureController.Update)
        feature.DELETE("/:id", featureController.Delete)
    }
}
```

2. **创建控制器**

```go
// internal/controllers/feature/create.go

package feature

func Create(c *gin.Context) {
    var req CreateFeatureRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(400, gin.H{"code": 400, "message": err.Error()})
        return
    }

    feature, err := featureService.Create(c, &req)
    if err != nil {
        c.JSON(500, gin.H{"code": 500, "message": err.Error()})
        return
    }

    c.JSON(200, gin.H{"code": 0, "data": feature})
}
```

3. **创建服务**

```go
// internal/services/featureService/service.go

package featureService

func Create(ctx context.Context, req *CreateFeatureRequest) (*models.Feature, error) {
    feature := &models.Feature{
        Name: req.Name,
    }

    if err := db.Create(feature).Error; err != nil {
        return nil, err
    }

    return feature, nil
}
```

#### 前端

1. **创建 API 封装**

```typescript
// web/src/api/feature.ts

import request from './request'

export interface Feature {
  id: number
  name: string
  created_at: string
}

export interface CreateFeatureDto {
  name: string
}

export const featureApi = {
  list: () => request.get<Feature[]>('/api/v1/features'),

  create: (data: CreateFeatureDto) =>
    request.post<Feature>('/api/v1/features', data),

  getById: (id: number) =>
    request.get<Feature>(`/api/v1/features/${id}`),

  update: (id: number, data: Partial<Feature>) =>
    request.put<Feature>(`/api/v1/features/${id}`, data),

  delete: (id: number) =>
    request.delete(`/api/v1/features/${id}`)
}
```

2. **在组件中使用**

```vue
<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { featureApi } from '@/api/feature'
import type { Feature } from '@/api/feature'

const features = ref<Feature[]>([])
const isLoading = ref(false)

const fetchFeatures = async () => {
  isLoading.value = true
  try {
    const { data } = await featureApi.list()
    features.value = data
  } catch (error) {
    console.error('获取功能列表失败:', error)
  } finally {
    isLoading.value = false
  }
}

onMounted(() => {
  fetchFeatures()
})
</script>
```

### 3. 添加新工具插件

详细步骤请参考 README.md 中的"添加新工具"章节。

**简要步骤**：

1. **创建工具目录**

```bash
mkdir -p pkg/utools/your_tool
```

2. **实现工具逻辑**

```go
// pkg/utools/your_tool/your_tool.go

package your_tool

type YourTool struct {
    *utools.BaseTool
}

func NewYourTool() *YourTool {
    metadata := &utools.ToolMetadata{
        Code: "your_tool",
        Name: "工具名称",
    }

    schema := &utools.ConfigSchema{
        // 配置 Schema
    }

    return &YourTool{
        BaseTool: utools.NewBaseTool(metadata, schema),
    }
}

func (t *YourTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
    // 执行逻辑
    return &utools.ExecutionResult{
        Success: true,
        Message: "执行成功",
    }, nil
}

func init() {
    tool := NewYourTool()
    utools.Register(tool)
}
```

3. **添加前端配置界面**

在 `web/src/pages/Tasks/components/ToolConfigDrawer.vue` 中添加对应的配置表单。

### 4. 数据库迁移

```bash
# 1. 修改模型
# 编辑 internal/models/xxx.go

# 2. 运行迁移
# 开发环境会自动迁移（AutoMigrate）
# 生产环境需要手动执行 SQL

# 3. 验证迁移
mysql -u root -p autoforge
SHOW TABLES;
DESC tasks;
```

---

## 部署流程

### 1. 构建生产版本

```bash
# 使用 Makefile 一键构建
make build

# 构建产物在 release/ 目录
ls -lh release/

# 输出：
# autoforge_prod_package.tar.gz  - 完整部署包
# autoforge                      - 二进制文件
# web/dist/                      - 前端静态文件
```

### 2. 部署到服务器

#### 方式一：直接部署

```bash
# 1. 上传到服务器
scp release/autoforge_prod_package.tar.gz user@server:/opt/

# 2. 解压
ssh user@server
cd /opt
tar -xzf autoforge_prod_package.tar.gz
cd autoforge

# 3. 配置
cp config.example.yaml config.yaml
vim config.yaml  # 修改生产配置

# 4. 启动服务
./autoforge

# 5. 配置 Systemd（推荐）
sudo vim /etc/systemd/system/autoforge.service
```

**Systemd 配置**：

```ini
[Unit]
Description=AutoForge Service
After=network.target

[Service]
Type=simple
User=autoforge
WorkingDirectory=/opt/autoforge
ExecStart=/opt/autoforge/autoforge
Restart=on-failure
RestartSec=5s

[Install]
WantedBy=multi-user.target
```

```bash
# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable autoforge
sudo systemctl start autoforge
sudo systemctl status autoforge
```

#### 方式二：Docker 部署

```bash
# 1. 构建镜像
docker build -t autoforge:latest .

# 2. 运行容器
docker run -d \
  --name autoforge \
  -p 7777:7777 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  --restart unless-stopped \
  autoforge:latest

# 3. 查看日志
docker logs -f autoforge

# 4. 停止/重启
docker stop autoforge
docker start autoforge
```

#### 方式三：Docker Compose

```bash
# 1. 编辑 docker-compose.yml

# 2. 启动所有服务
docker-compose up -d

# 3. 查看状态
docker-compose ps

# 4. 查看日志
docker-compose logs -f autoforge

# 5. 停止所有服务
docker-compose down
```

### 3. Nginx 反向代理

```nginx
# /etc/nginx/sites-available/autoforge

server {
    listen 80;
    server_name autoforge.example.com;

    # 前端静态文件（如果前后端分离部署）
    location / {
        root /opt/autoforge/web/dist;
        try_files $uri $uri/ /index.html;
    }

    # API 代理
    location /api/ {
        proxy_pass http://localhost:7777;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

```bash
# 启用配置
sudo ln -s /etc/nginx/sites-available/autoforge /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

### 4. SSL 证书（HTTPS）

```bash
# 使用 Let's Encrypt
sudo apt install certbot python3-certbot-nginx
sudo certbot --nginx -d autoforge.example.com

# 证书会自动配置到 Nginx
# 自动续期
sudo certbot renew --dry-run
```

---

## 问题排查

### 后端问题

#### 1. 服务启动失败

```bash
# 查看日志
tail -f logs/autoforge.log

# 常见原因：
# - 端口被占用
lsof -i :7777
kill -9 <PID>

# - 数据库连接失败
# 检查 config.yaml 中的数据库配置

# - 配置文件格式错误
# 使用 YAML 校验工具检查
```

#### 2. 任务不执行

```bash
# 1. 检查任务是否启用
# 在 Web 界面查看任务状态

# 2. 查看调度器日志
# 搜索关键词 "cron" 或 "scheduler"

# 3. 检查 next_run_at 字段
mysql -u root -p autoforge
SELECT id, name, enabled, next_run_at FROM tasks;

# 4. 手动触发任务测试
# 使用 API 或 Web 界面的"立即执行"按钮
```

#### 3. 内存/CPU 占用高

```bash
# 1. 查看进程状态
top -p $(pgrep autoforge)

# 2. 查看 goroutine 数量
# 访问 /debug/pprof (需在配置中启用)

# 3. 优化建议：
# - 减少并发任务数量
# - 增加任务执行间隔
# - 检查是否有死循环或内存泄漏
```

### 前端问题

#### 1. 页面白屏

```bash
# 1. 打开浏览器控制台查看错误

# 2. 检查 API 请求是否成功
# Network 标签查看请求状态

# 3. 检查路由配置
# 确认路由路径是否正确

# 4. 清除缓存
# 浏览器硬刷新：Ctrl + Shift + R
```

#### 2. API 请求失败

```typescript
// 检查以下配置

// 1. API 基础 URL 是否正确
// web/src/api/request.ts
const request = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || 'http://localhost:7777',
  timeout: 30000
})

// 2. 跨域问题
// 检查后端 CORS 配置
// internal/middleware/cors.go

// 3. Token 是否有效
// 检查 localStorage 中的 token
console.log(localStorage.getItem('token'))
```

#### 3. 组件导入错误

```bash
# 常见错误：
# Failed to resolve import "@/components/BaseButton"

# 解决方案：
# 1. 确认组件文件夹存在
ls -la web/src/components/BaseButton

# 2. 确认有 index.ts 导出文件
cat web/src/components/BaseButton/index.ts

# 3. 确认 index.vue 存在
ls web/src/components/BaseButton/index.vue

# 4. 重启开发服务器
cd web && pnpm dev
```

### 数据库问题

#### 1. 连接失败

```bash
# MySQL 连接测试
mysql -h 127.0.0.1 -P 3306 -u root -p

# 检查配置
cat config.yaml | grep -A 10 database

# 检查 MySQL 服务状态
sudo systemctl status mysql

# 查看 MySQL 错误日志
sudo tail -f /var/log/mysql/error.log
```

#### 2. 迁移失败

```bash
# 查看当前表结构
mysql -u root -p autoforge
SHOW TABLES;
DESC tasks;

# 删除表重新迁移（⚠️ 慎用，会丢失数据）
DROP TABLE tasks;

# 重启服务触发 AutoMigrate
./autoforge
```

---

## 附录

### 常用命令速查

#### Makefile 命令

```bash
make help               # 查看所有可用命令
make start              # 启动开发环境
make start-bg           # 后台启动
make stop               # 停止所有服务
make status             # 查看服务状态
make logs               # 查看日志
make build              # 构建生产版本
make test               # 运行测试
make lint               # 代码检查
make fmt                # 格式化代码
make clean              # 清理构建文件
```

#### Go 命令

```bash
go run cmd/main.go      # 运行程序
go build -o autoforge   # 编译
go test ./...           # 运行测试
go fmt ./...            # 格式化
go mod tidy             # 整理依赖
go mod download         # 下载依赖
```

#### 前端命令

```bash
pnpm install            # 安装依赖
pnpm dev                # 启动开发服务器
pnpm build              # 构建生产版本
pnpm preview            # 预览生产构建
pnpm lint               # 代码检查
pnpm format             # 格式化代码
pnpm type-check         # 类型检查
```

#### Git 命令

```bash
git status              # 查看状态
git add .               # 添加所有文件
git commit -m "msg"     # 提交
git push                # 推送
git pull                # 拉取
git checkout -b feat/x  # 创建新分支
git log --oneline       # 查看提交历史
```

### 环境变量

```bash
# 开发环境
export GO_ENV=development
export GIN_MODE=debug

# 生产环境
export GO_ENV=production
export GIN_MODE=release

# 前端环境变量（.env 文件）
VITE_API_BASE_URL=http://localhost:7777
VITE_APP_TITLE=AutoForge
```

### 端口占用

| 服务 | 默认端口 | 说明 |
|------|----------|------|
| 后端 API | 7777 | 可在 config.yaml 修改 |
| 前端开发服务器 | 3200 | 可在 vite.config.ts 修改 |
| MySQL | 3306 | 标准端口 |
| Redis | 6379 | 标准端口 |

### 日志位置

| 类型 | 位置 | 说明 |
|------|------|------|
| 后端日志 | `logs/autoforge.log` | 应用日志 |
| 访问日志 | `logs/access.log` | HTTP 请求日志 |
| 错误日志 | `logs/error.log` | 错误日志 |
| Nginx 日志 | `/var/log/nginx/` | Nginx 日志 |
| MySQL 日志 | `/var/log/mysql/` | 数据库日志 |

---

## 资源链接

### 官方文档
- [Go 官方文档](https://golang.org/doc/)
- [Vue 3 官方文档](https://vuejs.org/)
- [TypeScript 官方文档](https://www.typescriptlang.org/)
- [Gin 框架文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [Tailwind CSS 文档](https://tailwindcss.com/docs)
- [Vue Flow 文档](https://vueflow.dev/)

### 学习资源
- [Effective Go](https://golang.org/doc/effective_go)
- [Vue 3 风格指南](https://vuejs.org/style-guide/)
- [TypeScript 深入理解](https://www.typescriptlang.org/docs/handbook/intro.html)

### 工具推荐
- [VS Code](https://code.visualstudio.com/) - 代码编辑器
- [GoLand](https://www.jetbrains.com/go/) - Go IDE
- [Postman](https://www.postman.com/) - API 测试
- [TablePlus](https://tableplus.com/) - 数据库管理
- [Docker Desktop](https://www.docker.com/products/docker-desktop) - 容器管理

---

## 联系方式

- **GitHub Issues**: [提交问题](https://github.com/CooperJiang/AutoForge/issues)
- **Pull Requests**: [贡献代码](https://github.com/CooperJiang/AutoForge/pulls)
- **文档反馈**: 如果文档有不清楚的地方，欢迎提 Issue

---

**文档版本**: v2.0
**最后更新**: 2025-01-12
**维护者**: AutoForge Team

---

**祝你开发愉快！🎉**
