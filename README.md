# AutoForge

<p align="center">
  <img src="web/public/logo.png" alt="AutoForge Logo" width="120">
</p>

<p align="center">
  <strong>现代化的工作流自动化平台</strong>
</p>

<p align="center">
  <a href="#-快速开始">快速开始</a> •
  <a href="#-核心功能">核心功能</a> •
  <a href="#-架构设计">架构设计</a> •
  <a href="#-文档">文档</a>
</p>

---

## 📝 简介

AutoForge 是一个强大而优雅的工作流自动化平台，帮助您轻松构建、管理和执行自动化工作流。

**核心能力**：
- 🎨 可视化工作流编排 - 拖拽式设计，直观清晰
- ⏰ 灵活的调度系统 - Cron、定时、Webhook、外部 API 触发
- 🔗 丰富的工具节点 - HTTP、邮件、飞书、OpenAI、Redis 等 10+ 工具
- 📊 实时执行监控 - 可视化状态追踪，详细执行日志
- 🔐 企业级安全 - JWT 认证、RBAC 权限、数据加密
- 🔌 可扩展架构 - 插件式工具系统，易于开发自定义工具

---

## 🚀 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- pnpm
- MySQL 8.0+ 或 SQLite

### 安装运行

```bash
# 1. 克隆项目
git clone https://github.com/YourOrg/AutoForge.git
cd AutoForge

# 2. 配置环境
cp config.example.yaml config.yaml
# 编辑 config.yaml 配置数据库等信息

# 3. 启动服务
make start

# 或手动启动
# 后端
go run cmd/main.go

# 前端（新终端）
cd web && pnpm install && pnpm dev
```

### 访问应用

- **前端界面**: http://localhost:3200
- **后端 API**: http://localhost:7777
- **默认账号**: root / 88888888

---

## 🎯 核心功能

### 工作流编辑器

- 可视化拖拽编排
- 实时验证和错误提示
- 节点配置和变量管理
- 工作流导入/导出

### 节点类型

**触发器**
- 定时触发（Cron）
- 外部 API 触发
- Webhook 触发

**工具节点**
- **HTTP 请求** - 发送 HTTP 请求，支持所有 HTTP 方法
- **邮件发送** - SMTP 邮件通知，支持 HTML 格式
- **飞书机器人** - 发送飞书消息，支持文本、富文本、卡片
- **健康检查** - 监控网站可用性和 SSL 证书
- **OpenAI 对话** - GPT-3.5/GPT-4/GPT-4o 智能对话
- **OpenAI 图片生成** - DALL-E 文本生成图片
- **JSON 转换** - JavaScript 表达式数据转换
- **Redis 上下文** - Redis 状态存储和读取
- **输出格式化** - 格式化输出为图片、视频、HTML 等
- **HTML 内容保存** - 保存 HTML 并生成预览 URL

**控制节点**
- 条件判断 - If/Else 逻辑分支
- Switch 分支 - 多条件分支控制
- 延迟等待 - 延迟执行下一步

### 调度和执行

- 灵活的 Cron 表达式
- 多种调度规则（每天/每周/每月/间隔）
- 异步执行引擎
- 实时执行状态监控
- 详细的执行日志

### 管理功能

- 用户和权限管理
- 工作流统计分析
- 执行历史查询
- 任务管理面板

---

## 🏗 架构设计

AutoForge 采用现代化的前后端分离架构：

```
┌─────────────────────────────────────────┐
│       前端层 (Vue 3 + TypeScript)        │
│                                         │
│  工作流编辑器 │ 监控面板 │ 管理后台      │
└────────────────┬────────────────────────┘
                 │ RESTful API
                 ▼
┌─────────────────────────────────────────┐
│          后端层 (Go + Gin)               │
│                                         │
│  工作流引擎 │ 调度器 │ 工具系统         │
└────────────────┬────────────────────────┘
                 │
                 ▼
┌─────────────────────────────────────────┐
│         数据层 (MySQL/SQLite)            │
│                                         │
│  工作流 │ 执行记录 │ 用户 │ 配置        │
└─────────────────────────────────────────┘
```

### 核心组件

**工作流引擎**
- 工作流解析和验证
- 异步执行引擎
- 变量系统和上下文管理
- 错误处理和重试机制

**调度系统**
- 基于 Cron v3 的调度器
- 多种触发方式支持
- 任务队列管理
- 执行状态追踪

**工具系统**
- 可扩展的工具架构
- 统一的工具接口
- 内置常用工具
- 支持自定义工具开发

详细架构文档请参考 [工作流设计文档](./docs/workflow-design.md)。

---

## 📖 文档

完整的开发文档帮助您深入了解 AutoForge：

**架构和设计**
- **[工作流设计](./docs/workflow-design.md)** - 工作流引擎架构和设计理念
- **[变量系统](./docs/VARIABLE_SYSTEM.md)** - 变量引用和表达式使用

**开发指南**
- **[开发指南](./docs/DEVELOPMENT_GUIDE.md)** - 环境搭建和开发流程
- **[组件开发规范](./docs/COMPONENT_DEVELOPMENT.md)** - 前端组件开发标准
- **[主题系统规范](./docs/THEME_SPEC.md)** - 主题和样式系统
- **[工具开发指南](./docs/TOOL_DEVELOPMENT_GUIDE.md)** - 自定义工具开发

**使用指南**
- **[飞书机器人使用指南](./docs/FEISHU_BOT_GUIDE.md)** - 飞书机器人配置和使用

---

## 🤝 贡献指南

欢迎贡献代码、报告问题或提出建议！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

开发规范请查看 [开发文档](./docs/DEVELOPMENT_GUIDE.md)。

---

## 📝 路线图

### ✅ v1.0 已完成

- [x] 可视化工作流编辑器（VueFlow）
- [x] 10+ 内置工具节点（HTTP、邮件、飞书、OpenAI 等）
- [x] 灵活的调度系统（Cron、定时、Webhook、外部 API）
- [x] 变量系统和上下文管理
- [x] 工作流导入/导出
- [x] 用户认证和 RBAC 权限
- [x] 执行监控和详细日志
- [x] 管理后台（用户、任务、执行记录）
- [x] 暗色模式支持
- [x] Redis 上下文存储
- [x] 工具配置系统

### 🚧 v1.1 进行中

- [ ] 工作流模板市场
- [ ] 工作流版本管理
- [ ] 更多第三方工具集成
- [ ] 性能优化和监控

### 💡 v2.0 计划中

- [ ] AI 辅助工作流生成
- [ ] 工作流性能分析
- [ ] 协作和分享功能
- [ ] 国际化支持
- [ ] 插件市场

---

## 📄 License

本项目采用 [MIT License](LICENSE) 开源协议。

---

<div align="center">

**Made with ❤️ by AutoForge Team**

[📚 文档](./docs/DEVELOPMENT_GUIDE.md) · [💬 问题反馈](https://github.com/YourOrg/AutoForge/issues)

</div>
