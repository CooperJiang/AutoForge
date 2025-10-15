# Cooper UI 组件库

> 基于 Vue 3 + TypeScript 的现代化组件库

## 📁 组件结构

每个组件都是一个独立的文件夹，包含 `index.vue` 入口文件。组件按功能分类如下：

### 🎨 基础组件

| 组件             | 路径                 | 说明                             |
| ---------------- | -------------------- | -------------------------------- |
| BaseButton       | `./BaseButton`       | 基础按钮组件，支持多种变体和尺寸 |
| BaseInput        | `./BaseInput`        | 基础输入框组件，支持标签、验证等 |
| BaseSelect       | `./BaseSelect`       | 下拉选择组件，支持搜索和多选     |
| Dialog           | `./Dialog`           | 模态对话框组件                   |
| Drawer           | `./Drawer`           | 侧边抽屉组件                     |
| Message          | `./Message`          | 消息提示组件                     |
| MessageContainer | `./MessageContainer` | 消息容器组件                     |

### 📋 布局组件

| 组件       | 路径           | 说明           |
| ---------- | -------------- | -------------- |
| AppHeader  | `./AppHeader`  | 应用顶部导航栏 |
| Pagination | `./Pagination` | 分页组件       |
| Table      | `./Table`      | 数据表格组件   |

### 📝 表单组件

| 组件             | 路径                 | 说明                                     |
| ---------------- | -------------------- | ---------------------------------------- |
| TimePicker       | `./TimePicker`       | 时间选择器                               |
| WeekDayPicker    | `./WeekDayPicker`    | 星期选择器                               |
| MonthDayPicker   | `./MonthDayPicker`   | 月份日期选择器                           |
| ParamInput       | `./ParamInput`       | 键值对参数输入组件                       |
| VariableSelector | `./VariableSelector` | 变量选择器（支持环境变量、节点输出引用） |

### 🎭 展示组件

| 组件             | 路径                 | 说明                      |
| ---------------- | -------------------- | ------------------------- |
| JsonViewer       | `./JsonViewer`       | JSON 查看器，支持语法高亮 |
| NextRunCountdown | `./NextRunCountdown` | 倒计时组件                |

### 💬 对话框组件

| 组件                  | 路径                      | 说明           |
| --------------------- | ------------------------- | -------------- |
| TestResultDialog      | `./TestResultDialog`      | 测试结果对话框 |
| TaskDetailDialog      | `./TaskDetailDialog`      | 任务详情对话框 |
| ExecutionDetailDialog | `./ExecutionDetailDialog` | 执行详情对话框 |

### 🔧 高级组件

| 组件        | 路径            | 说明                           |
| ----------- | --------------- | ------------------------------ |
| RetryConfig | `./RetryConfig` | 重试配置组件，支持指数退避策略 |

### 🛠️ 工具组件

工具组件位于 `./tools/` 目录下，主要用于工作流配置：

- `TriggerConfig.vue` - 触发器配置
- `ConditionConfig.vue` - 条件配置
- `DelayConfig.vue` - 延迟配置
- `SwitchConfig.vue` - 开关配置
- `HealthCheckerConfig.vue` - 健康检查配置
- `EmailToolConfig.vue` - 邮件工具配置

---

## 🚀 使用方式

### 1. 全局注册（推荐）

在 `main.ts` 中注册所有组件：

```typescript
import { createApp } from 'vue'
import CooperUI from './components'
import App from './App.vue'

const app = createApp(App)
app.use(CooperUI)
app.mount('#app')
```

使用时无需导入：

```vue
<template>
  <BaseButton size="lg" variant="primary"> 点击我 </BaseButton>
</template>
```

### 2. 按需注册

只注册需要的组件：

```typescript
import { createApp } from 'vue'
import { createCooperUI } from './components'
import App from './App.vue'

const app = createApp(App)
app.use(
  createCooperUI({
    components: ['BaseButton', 'BaseInput', 'Dialog'],
  })
)
app.mount('#app')
```

### 3. 按需引入（局部使用）

在单个组件中按需引入：

```vue
<script setup lang="ts">
import { BaseButton, BaseInput, Dialog } from '@/components'
</script>

<template>
  <BaseButton @click="open">打开对话框</BaseButton>
  <Dialog v-model="visible">
    <BaseInput v-model="value" placeholder="请输入" />
  </Dialog>
</template>
```

---

## 📖 开发规范

### 组件文件结构

```
ComponentName/
├── index.vue          # 组件入口文件（必需）
├── types.ts           # 类型定义（可选）
├── hooks.ts           # 组合式函数（可选）
├── constants.ts       # 常量定义（可选）
└── README.md          # 组件文档（可选）
```

### 代码规范

1. **单文件不超过 500 行**：超过则需拆分
2. **统一命名**：
   - 文件夹名：PascalCase（如 `BaseButton`）
   - 入口文件：统一使用 `index.vue`
3. **类型定义**：使用 TypeScript，定义清晰的 Props 和 Emits
4. **注释完整**：Props、Emits、方法都需要注释

### 导入规范

```typescript
// ✅ 推荐：使用文件夹路径
import BaseButton from '@/components/BaseButton'

// ❌ 避免：使用 .vue 扩展名
import BaseButton from '@/components/BaseButton.vue'
```

---

## 🎨 设计系统

详细的设计规范请参考：[COMPONENT_DESIGN_SPEC.md](./COMPONENT_DESIGN_SPEC.md)

### 快速参考

#### 尺寸系统

```typescript
type Size = 'xs' | 'sm' | 'md' | 'lg' | 'xl'
```

#### 按钮变体

```typescript
type ButtonVariant =
  | 'primary' // 主要操作
  | 'secondary' // 次要操作
  | 'success' // 成功操作
  | 'danger' // 危险操作
  | 'warning' // 警告操作
  | 'ghost' // 幽灵按钮
```

#### 颜色主题

- **主色**：`#3b82f6` - 主要按钮、链接
- **成功**：`#10b981` - 成功状态
- **警告**：`#f59e0b` - 警告提示
- **危险**：`#ef4444` - 危险操作

---

## 🔧 工具函数

### message

消息提示工具函数，支持多种类型：

```typescript
import { message } from '@/components'

// 成功提示
message.success('操作成功')

// 错误提示
message.error('操作失败')

// 警告提示
message.warning('请注意')

// 信息提示
message.info('提示信息')
```

---

## 📦 导出清单

### 组件导出

所有组件通过 `./index.ts` 统一导出：

```typescript
export {
  BaseButton,
  BaseInput,
  BaseSelect,
  Dialog,
  Drawer,
  // ... 其他组件
}
```

### 工具函数导出

```typescript
export { message } from '@/utils/message'
```

---

## 🤝 贡献指南

### 添加新组件

1. 在 `components/` 目录下创建组件文件夹
2. 创建 `index.vue` 入口文件
3. 在 `components/index.ts` 中导出组件
4. 更新本 README 文档

### 代码审查清单

- [ ] 组件结构符合规范
- [ ] 代码不超过 500 行
- [ ] Props 和 Emits 类型定义完整
- [ ] 支持响应式设计
- [ ] 包含必要的注释
- [ ] 已添加到导出文件

---

## 📝 更新日志

### v1.0.0 (2025-01-12)

**重大重构**

- ✅ 将所有单文件组件重构为文件夹结构
- ✅ 创建统一的导出文件 `index.ts`
- ✅ 更新所有组件引用路径
- ✅ 添加组件设计规范文档

**组件清单**

- 21 个基础组件
- 6 个工具配置组件
- 支持全局注册和按需引入

---

## 📚 参考资源

- [组件设计规范](./COMPONENT_DESIGN_SPEC.md)
- [Vue 3 官方文档](https://cn.vuejs.org/)
- [TypeScript 官方文档](https://www.typescriptlang.org/)

---

**维护团队**: Cooper Team
**最后更新**: 2025-01-12
