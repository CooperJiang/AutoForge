# 🔨 AutoForge

<p align="center">
  <img src="web/public/logo.png" alt="AutoForge Logo" width="120">
</p>

<p align="center">
  <strong>强大的自动化工坊 - 让任务自动流动，让时间为你工作</strong>
</p>

<p align="center">
  <a href="https://github.com/CooperJiang/AutoForge">
    <img src="https://img.shields.io/github/stars/CooperJiang/AutoForge?style=social" alt="GitHub stars">
  </a>
  <a href="https://github.com/CooperJiang/AutoForge/blob/main/LICENSE">
    <img src="https://img.shields.io/github/license/CooperJiang/AutoForge" alt="License">
  </a>
  <a href="https://golang.org">
    <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go" alt="Go Version">
  </a>
  <a href="https://vuejs.org">
    <img src="https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js" alt="Vue Version">
  </a>
</p>

<p align="center">
  <a href="#-特性">特性</a> •
  <a href="#-技术栈">技术栈</a> •
  <a href="#-快速开始">快速开始</a> •
  <a href="#-部署指南">部署指南</a> •
  <a href="#-贡献指南">贡献</a>
</p>

---

## ✨ 特性

### 核心功能
- ⏰ **灵活的定时任务** - 支持 Cron、间隔、每日/周/月等多种调度方式
- 🔧 **多工具系统** - 内置 HTTP 请求、邮件发送、健康检查等实用工具
- 📊 **可视化管理** - 现代化的 Web 界面，实时查看任务状态和执行记录
- 🎯 **精准触发** - 秒级精度的任务调度引擎
- 📝 **执行日志** - 完整记录任务执行历史、响应结果、错误信息

### 工具能力
- 🔗 **HTTP 请求** - 支持所有 HTTP 方法、自定义 Headers/Body/Params，一键解析 cURL 命令
- 📧 **邮件发送** - SMTP 邮件发送，支持 HTML/文本格式、多收件人、抄送
- 🏥 **健康检查** - 网站/API 健康监控，SSL 证书检查，响应时间统计，支持复杂鉴权

### 用户体验
- 🔐 **安全认证** - JWT + OAuth2 (Linux.do) 双重登录方式
- 👥 **多用户系统** - 支持用户注册、权限管理、管理员后台
- 🎨 **现代UI** - 响应式设计、一键测试、ESC 快捷键支持
- 🚀 **高性能** - Go 后端 + Vue3 前端，极速响应
- 🐳 **容器化** - 支持 Docker 一键部署

---

## 🛠 技术栈

### 后端
- **Go 1.21+** - 高性能后端服务
- **Gin** - 轻量级 Web 框架
- **GORM** - 强大的 ORM 数据库操作
- **Cron v3** - 可靠的定时任务调度
- **JWT** - 安全的身份认证
- **OAuth2** - Linux.do 社区登录集成

### 前端
- **Vue 3** - 渐进式前端框架
- **TypeScript** - 类型安全的 JavaScript
- **Vite** - 极速构建工具
- **Tailwind CSS** - 原子化 CSS 框架
- **Pinia** - 轻量级状态管理
- **Lucide** - 精美的图标库

### 数据库 & 缓存
- **MySQL 8.0+** - 主数据库
- **SQLite** - 开发/轻量级部署
- **Redis** (可选) - 高性能缓存

---

## 🎯 功能亮点

### 🔧 插件化工具系统
AutoForge 采用插件化架构，每个工具都是独立的插件，易于扩展。内置三大核心工具：

#### 1. HTTP 请求工具
- ✨ 支持所有 HTTP 方法（GET/POST/PUT/DELETE/PATCH）
- 📋 一键粘贴 cURL 命令自动解析
- 🎨 可视化编辑 Headers、Params、Body
- 🧪 实时测试功能，查看完整响应

#### 2. 邮件发送工具
- 📧 SMTP 协议发送邮件
- 🔐 支持 SSL/TLS 加密（端口 465/587）
- 📝 支持 HTML 和纯文本格式
- 👥 多收件人、抄送支持
- ⚡ 系统统一配置，用户无需提供 SMTP 信息

#### 3. 健康检查工具
- 🏥 网站/API 可用性监控
- 🔒 SSL 证书到期检查和告警
- ⏱️ 响应时间统计
- 🔍 支持正则匹配响应内容
- 🔐 支持复杂鉴权（Headers/Body）
- 📊 清晰的状态报告（网站状态 + SSL 证书状态）

### ⏰ 灵活的调度系统
- **每天**：固定时间执行（如每天 8:00）
- **每周**：指定星期执行（如每周一 9:00）
- **每月**：指定日期执行（如每月 1 号）
- **间隔**：固定间隔执行（如每 5 分钟）
- **Cron**：自定义 Cron 表达式（最灵活）

### 🎨 现代化用户体验
- 🚀 Vue 3 + TypeScript + Tailwind CSS
- 📱 响应式设计，完美支持移动端
- ⌨️ 快捷键支持（ESC 关闭弹窗）
- 🧪 任务一键测试，实时查看结果
- 🎯 直观的任务状态管理（启用/禁用）
- 📊 详细的执行记录和日志

---

## 🚀 快速开始

### 前置要求
- Go 1.21+
- Node.js 18+
- pnpm (推荐) 或 npm
- MySQL 8.0+ (可选，默认 SQLite)

### 1. 克隆项目
```bash
git clone https://github.com/CooperJiang/AutoForge.git
cd AutoForge
```

### 2. 配置数据库 (可选)
```bash
# 如果使用 MySQL
mysql -u root -p
CREATE DATABASE autoforge CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 3. 配置文件
```bash
# 复制示例配置
cp config.example.yaml config.yaml

# 编辑配置文件
vim config.yaml
```

主要配置项：
```yaml
app:
  name: "AutoForge"
  port: 7777
  mode: "debug"  # debug 或 release

database:
  driver: "mysql"  # mysql 或 sqlite
  host: "127.0.0.1"
  port: 3306
  username: "root"
  password: "your_password"
  name: "autoforge"

jwt:
  secret_key: "your-secret-key-here"  # 修改为强密钥
  expires_in: 24  # 小时

mail:
  enabled: true  # 启用邮件功能
  host: "smtp.example.com"
  port: 465  # 465 使用 SSL, 587 使用 STARTTLS
  username: "your-email@example.com"
  password: "your-password"
  from: "noreply@example.com"
  from_name: "AutoForge"
  ssl: true  # 是否使用 TLS/SSL
```

### 4. 启动开发环境

**方式一：使用 Makefile（推荐）**
```bash
# 一键启动（自动打开新终端窗口）
make start

# 或者后台启动
make start-bg

# 查看服务状态
make status

# 查看日志
make logs

# 停止服务
make stop
```

**方式二：手动启动**
```bash
# 终端 1: 启动后端
go mod download
go run cmd/main.go

# 终端 2: 启动前端
cd web
pnpm install
pnpm dev
```

### 5. 访问应用
- **前端界面**: http://localhost:3200
- **后端 API**: http://localhost:7777
- **管理后台**: http://localhost:3200/admin

**首次启动**：
- 系统会自动创建管理员账号
- 默认密码在控制台输出，请及时修改

---

## 📖 使用指南

### 创建定时任务

1. **注册/登录账号**
   - 普通注册：邮箱 + 密码
   - OAuth2 登录：使用 Linux.do 账号

2. **选择工具类型**
   - **HTTP 请求**：发送 HTTP 请求到指定 URL
   - **邮件发送**：自动发送邮件通知
   - **健康检查**：监控网站/API 可用性和 SSL 证书

3. **配置工具参数**

   **HTTP 请求工具**：
   - 请求方式：GET/POST/PUT/DELETE/PATCH
   - 接口地址：目标 API URL
   - 请求头/参数/请求体：自定义配置
   - 💡 支持直接粘贴 cURL 命令自动解析

   **邮件发送工具**：
   - 收件人：多个邮箱用逗号分隔
   - 邮件主题和正文
   - 内容类型：纯文本/HTML
   - 系统自动使用配置的 SMTP 服务器

   **健康检查工具**：
   - 检查 URL 和请求方法
   - 超时时间、期望状态码
   - SSL 证书检查和到期告警
   - 支持正则匹配响应内容
   - 💡 支持复杂鉴权（Headers/Body）

4. **配置调度规则**
   - **每天**：每天固定时间执行
   - **每周**：每周特定星期执行
   - **每月**：每月特定日期执行
   - **间隔**：按固定间隔执行
   - **Cron**：使用 Cron 表达式（最灵活）

5. **测试和启用**
   - 点击"测试配置"按钮验证工具配置
   - 启用任务，自动按计划执行

### 查看执行记录

- **任务列表**：查看所有任务和状态
- **执行记录**：查看详细执行日志
- **响应内容**：查看 API 返回结果
- **错误信息**：排查失败原因

### 管理员功能

访问 `/admin` 进入管理后台：
- **任务管理**：查看/编辑/删除所有用户的任务
- **执行记录**：全局执行记录查询
- **用户管理**：启用/禁用用户账号
- **系统统计**：任务数、执行次数、成功率

---

## 📦 部署指南

### 方式一：直接部署

```bash
# 1. 构建生产版本
make build

# 2. 部署包位于 ./release/
ls release/

# 3. 上传到服务器
scp release/autoforge_prod_package.tar.gz user@server:/opt/

# 4. 解压运行
tar -xzf autoforge_prod_package.tar.gz
./autoforge
```

### 方式二：Docker 部署

```bash
# 构建镜像
docker build -t autoforge:latest .

# 运行容器
docker run -d \
  --name autoforge \
  -p 7777:7777 \
  -v $(pwd)/config.yaml:/app/config.yaml \
  -v $(pwd)/data:/app/data \
  -v $(pwd)/logs:/app/logs \
  --restart unless-stopped \
  autoforge:latest
```

### 方式三：Docker Compose

```yaml
version: '3.8'
services:
  autoforge:
    image: autoforge:latest
    container_name: autoforge
    ports:
      - "7777:7777"
    volumes:
      - ./config.yaml:/app/config.yaml
      - ./data:/app/data
      - ./logs:/app/logs
    environment:
      - TZ=Asia/Shanghai
    restart: unless-stopped

  mysql:
    image: mysql:8.0
    container_name: autoforge-mysql
    environment:
      MYSQL_ROOT_PASSWORD: your_password
      MYSQL_DATABASE: autoforge
    volumes:
      - mysql_data:/var/lib/mysql
    restart: unless-stopped

volumes:
  mysql_data:
```

### Nginx 反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:7777;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

---

## 🔧 开发指南

### 项目结构

```
AutoForge/
├── cmd/                    # 应用入口
│   └── main.go            # 主程序
├── internal/              # 内部代码
│   ├── controllers/       # HTTP 控制器
│   │   ├── task/          # 任务相关接口
│   │   ├── auth/          # 认证相关接口
│   │   └── admin/         # 管理员接口
│   ├── services/          # 业务逻辑层
│   │   ├── taskService/   # 任务服务
│   │   ├── authService/   # 认证服务
│   │   └── cronService/   # 调度服务
│   ├── models/            # 数据模型
│   ├── routes/            # 路由定义
│   ├── middleware/        # 中间件
│   └── cron/              # 定时任务调度器
├── pkg/                   # 公共包
│   ├── config/            # 配置管理
│   ├── database/          # 数据库连接
│   ├── logger/            # 日志工具
│   ├── errors/            # 错误处理
│   ├── common/            # 公共工具
│   └── utools/            # 工具系统
│       ├── http/          # HTTP 请求工具
│       ├── email/         # 邮件发送工具
│       └── health/        # 健康检查工具
├── web/                   # 前端代码
│   ├── src/
│   │   ├── pages/         # 页面组件
│   │   │   ├── Tasks/     # 任务管理页面
│   │   │   ├── Auth/      # 登录注册页面
│   │   │   └── Admin/     # 管理后台页面
│   │   ├── components/    # 通用组件
│   │   │   ├── BaseInput.vue
│   │   │   ├── BaseSelect.vue
│   │   │   ├── Drawer.vue
│   │   │   └── Dialog.vue
│   │   ├── api/           # API 接口
│   │   ├── utils/         # 工具函数
│   │   │   ├── curlParser.ts  # cURL 解析
│   │   │   └── message.ts     # 消息提示
│   │   ├── composables/   # 组合式函数
│   │   ├── router/        # 路由配置
│   │   └── layouts/       # 布局组件
│   ├── public/            # 静态资源
│   └── package.json       # 前端依赖
├── config.yaml            # 配置文件
├── config.example.yaml    # 配置示例
├── go.mod                 # Go 模块定义
├── Makefile               # 构建脚本
├── Dockerfile             # Docker 镜像
└── README.md              # 项目文档
```

### 可用命令

```bash
# 开发相关
make start              # 启动开发环境（新终端）
make start-bg           # 启动开发环境（后台模式）
make dev                # 只启动后端
make web-dev            # 只启动前端
make stop               # 停止所有服务
make status             # 查看服务状态
make logs               # 查看日志

# 构建相关
make build              # 构建生产版本
make build-backend      # 只构建后端
make web-build          # 只构建前端

# 测试相关
make test               # 运行测试
make test-coverage      # 测试覆盖率

# 代码质量
make fmt                # 格式化代码
make lint               # 代码检查
make web-lint           # 前端代码检查

# 清理
make clean              # 清理构建文件

# 帮助
make help               # 查看所有命令
```

### 添加新工具

AutoForge 使用插件化的工具系统，添加新工具非常简单：

#### 1. 创建工具实现 (`pkg/utools/your_tool/`)

```go
package your_tool

import (
    "auto-forge/pkg/utools"
    "time"
)

type YourTool struct {
    *utools.BaseTool
}

func NewYourTool() *YourTool {
    metadata := &utools.ToolMetadata{
        Code:        "your_tool",
        Name:        "工具名称",
        Description: "工具描述",
        Category:    "工具分类",
        Version:     "1.0.0",
        Author:      "作者",
        Icon:        "图标名称",
        AICallable:  true,
        Tags:        []string{"标签1", "标签2"},
    }

    schema := &utools.ConfigSchema{
        Type: "object",
        Properties: map[string]utools.PropertySchema{
            "param1": {
                Type:        "string",
                Title:       "参数1",
                Description: "参数说明",
            },
        },
        Required: []string{"param1"},
    }

    return &YourTool{
        BaseTool: utools.NewBaseTool(metadata, schema),
    }
}

func (t *YourTool) Execute(ctx *utools.ExecutionContext, config map[string]interface{}) (*utools.ExecutionResult, error) {
    startTime := time.Now()

    // 解析配置
    param1, _ := config["param1"].(string)

    // 执行工具逻辑
    // ...

    return &utools.ExecutionResult{
        Success:    true,
        Message:    "执行成功",
        Output:     map[string]interface{}{"result": "结果"},
        DurationMs: time.Since(startTime).Milliseconds(),
    }, nil
}

func init() {
    tool := NewYourTool()
    if err := utools.Register(tool); err != nil {
        panic(err)
    }
}
```

#### 2. 添加前端配置界面 (`web/src/pages/Tasks/components/ToolConfigDrawer.vue`)

在工具配置抽屉中添加对应的配置表单，参考现有的 HTTP、邮件、健康检查工具配置。

#### 3. 其他开发流程

**后端 API**：
- 数据模型：`internal/models/`
- 业务逻辑：`internal/services/`
- 控制器：`internal/controllers/`
- 路由：`internal/routes/`

**前端开发**：
- API 接口：`web/src/api/`
- 页面组件：`web/src/pages/`
- 通用组件：`web/src/components/`
- 工具函数：`web/src/utils/`

---

## 🤝 贡献指南

欢迎所有形式的贡献！

### 贡献方式

1. Fork 本仓库
2. 创建特性分支
   ```bash
   git checkout -b feature/AmazingFeature
   ```
3. 提交更改
   ```bash
   git commit -m 'Add some AmazingFeature'
   ```
4. 推送到分支
   ```bash
   git push origin feature/AmazingFeature
   ```
5. 开启 Pull Request

### 代码规范

- **Go**: 遵循 [Effective Go](https://golang.org/doc/effective_go) 和 [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- **TypeScript/Vue**: 遵循 [Vue 3 Style Guide](https://vuejs.org/style-guide/) 和 [TypeScript Style Guide](https://www.typescriptlang.org/docs/handbook/declaration-files/do-s-and-don-ts.html)
- 提交前运行 `make fmt` 和 `make lint`

---

## 📝 路线图

### 已完成 ✅
- [x] 定时任务调度系统（Cron/间隔/每日/每周/每月）
- [x] HTTP 请求工具（支持所有方法，cURL 解析）
- [x] 邮件发送工具（SMTP，HTML/文本）
- [x] 健康检查工具（网站监控，SSL 证书检查）
- [x] 用户认证系统（JWT + OAuth2）
- [x] OAuth2 登录（Linux.do）
- [x] 管理员后台
- [x] 执行日志记录
- [x] 任务一键测试
- [x] 快捷键支持（ESC 关闭抽屉）

### 计划中 🚧
- [ ] 更多工具插件（数据库备份、文件同步等）
- [ ] Webhook 触发器
- [ ] 任务依赖关系和工作流
- [ ] API 监控和智能告警
- [ ] 执行统计图表和仪表板
- [ ] 更多 OAuth2 登录方式（GitHub、Google）
- [ ] 移动端 App
- [ ] 国际化支持（i18n）

---

## 📄 License

本项目采用 [MIT License](LICENSE) 开源协议。

---

## 🙏 致谢

- [Gin](https://github.com/gin-gonic/gin) - 高性能 HTTP Web 框架
- [GORM](https://github.com/go-gorm/gorm) - Go ORM 库
- [Vue.js](https://vuejs.org/) - 渐进式 JavaScript 框架
- [Tailwind CSS](https://tailwindcss.com/) - 实用优先的 CSS 框架
- [Linux.do](https://linux.do) - OAuth2 登录支持
- 所有贡献者和使用者 ❤️

---

## 📧 联系方式

- **GitHub Issues**: [提交问题](https://github.com/CooperJiang/AutoForge/issues)
- **Pull Requests**: [贡献代码](https://github.com/CooperJiang/AutoForge/pulls)

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给个 Star！⭐**

Made with ❤️ by [CooperJiang](https://github.com/CooperJiang)

[⬆ 回到顶部](#-autoforge)

</div>
