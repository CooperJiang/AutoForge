# Release Notes - v0.2.0

## 🎉 新功能

### 模板市场
- 全新的工作流模板市场，支持浏览、搜索和安装模板
- 一键将工作流发布为模板，方便团队协作和复用
- Lucide 图标库集成，80+ 精美图标可选
- 支持自定义模板分类和图片 URL
- 模板详情展示，包含使用指南和参数说明
- 完整的模板安装历史记录

### 工作流编辑器增强
- **画布状态持久化**：自动保存画布位置和缩放级别，下次打开时恢复到上次状态
- **智能离开确认**：改进的离开确认对话框，支持三种操作
  - 保存并离开
  - 放弃更改
  - 继续编辑
- **更准确的更改检测**：修复了进入编辑页面立即提示有未保存更改的问题

## 🐛 Bug 修复
- 修复工作流编辑器中"已禁用"按钮无法切换的问题
- 修复 Dialog 组件 Vue 警告（modelValue prop）
- 修复保存后导航被阻塞的问题
- 优化视口恢复，消除页面闪烁

## 🔧 技术改进

### 后端
- 新增 `WorkflowTemplate` 模型，支持模板管理
- 新增 `TemplateInstall` 模型，记录安装历史
- 新增 `WorkflowViewport` 字段，支持画布状态持久化
- 完整的模板 API 接口和服务层
- 规范的请求和响应 DTO 结构

### 前端
- 新增 Marketplace 页面和相关组件
- ConfirmDialog 组件支持第三个可选按钮
- PublishTemplateDialog 支持图标库选择
- VueFlow 渲染优化，避免视口闪烁
- 代码质量提升，清理调试日志

## 📦 文件变更统计

### 新增文件
- `internal/controllers/template/` - 模板控制器
- `internal/models/workflowTemplate.go` - 模板模型
- `internal/models/templateInstall.go` - 安装记录模型
- `internal/services/template/` - 模板服务层
- `web/src/pages/Marketplace/` - 模板市场页面
- `web/src/components/tools/` - 新增多个工具配置组件
- `CHANGELOG.md` - 变更日志

### 修改文件
- `internal/models/workflow.go` - 添加 viewport 字段
- `internal/dto/request/workflowRequest.go` - 支持 viewport
- `internal/dto/response/workflowResponse.go` - 返回 viewport
- `web/src/pages/Workflows/editor.vue` - 视口持久化和离开确认
- `web/src/components/ConfirmDialog/index.vue` - 支持第三按钮
- 多个其他文件的优化和改进

## 🚀 升级指南

1. **数据库迁移**：新增 `workflow_template` 和 `template_install` 表，`workflow` 表新增 `viewport` 字段，启动时会自动迁移

2. **前端依赖**：无新增依赖，运行 `pnpm install` 确保依赖最新

3. **配置变更**：无配置变更

## 📝 使用提示

### 发布模板
1. 创建或编辑一个工作流
2. 保存工作流
3. 点击顶部的"发布为模板"按钮（仅管理员可见）
4. 填写模板信息，选择图标和分类
5. 发布成功后，模板将出现在模板市场

### 使用模板市场
1. 访问"模板市场"页面
2. 浏览或搜索需要的模板
3. 点击模板查看详情
4. 点击"安装模板"即可创建基于该模板的工作流
5. 安装历史可在个人中心查看

### 画布状态保存
- 编辑工作流时，画布的位置和缩放会自动保存
- 下次打开该工作流时，会自动恢复到上次的视图状态
- 无需任何额外操作，体验更流畅

---

**发布日期**: 2025-10-16
**版本**: v0.2.0
**构建**: Auto-Forge Workflow Engine
