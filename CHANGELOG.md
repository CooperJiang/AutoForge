# Changelog

All notable changes to this project will be documented in this file.

## [Unreleased]

### Added
- **模板市场功能**
  - 新增工作流模板市场，支持模板浏览、搜索和安装
  - 支持将工作流发布为模板，包含图标库选择和自定义分类
  - 模板详情展示，包含使用指南和参数说明
  - 模板安装历史记录

- **工作流编辑器增强**
  - 画布视口状态持久化，保存画布位置和缩放级别
  - 优化视口恢复逻辑，避免页面闪烁
  - 改进离开确认对话框，支持"保存并离开"、"放弃更改"、"继续编辑"三个选项
  - 修复 false positive 的未保存更改检测

- **组件改进**
  - ConfirmDialog 组件支持第三个可选按钮
  - PublishTemplateDialog 支持 Lucide 图标库选择（80+ 图标）
  - 支持自定义模板分类
  - 支持图片 URL 作为模板图标

### Fixed
- 修复工作流编辑器中"已禁用"按钮无法切换的问题
- 修复 Dialog 组件 modelValue 属性警告
- 修复保存后未清除 hasUnsavedChanges 导致的导航阻塞
- 修复工作流列表进入编辑页面立即提示有未保存更改的问题

### Changed
- 清理调试日志，提升代码质量
- 优化视口状态管理，改善用户体验

### Technical
- 后端新增 WorkflowTemplate 和 TemplateInstall 模型
- 新增模板相关的 API 接口和服务层
- 前端新增 Marketplace 页面和相关组件
- 扩展 Workflow 模型支持 viewport 字段
- 完善请求和响应 DTO 结构

---

## Version History

未来版本将在此记录...
