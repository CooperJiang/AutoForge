# 🚀 Cooper UI 组件设计规范

> 基于现代化设计的组件标准，注重实用性、一致性和可维护性

## 📋 设计原则

### 🎯 核心理念

- **实用性优先**: 功能完整，交互清晰，用户体验流畅
- **统一一致性**: 所有组件遵循相同的设计语言和交互模式
- **高效性能**: 简化不必要的效果，确保在各种设备上流畅运行
- **易于维护**: 代码结构清晰，注释完整，便于团队协作

### 🎨 视觉层次原则

1. **主体突出**: 核心内容区域需要有明确的视觉重点
2. **间距合理**: 使用统一的间距系统，提高可读性
3. **层次分明**: 通过背景、边框、阴影建立清晰的信息层级
4. **细节精致**: 在不影响性能的前提下添加适当的交互细节

---

## 🗂️ 组件文件结构规范

### 📁 标准文件夹结构

每个组件应该是一个独立的文件夹，包含以下文件：

```
ComponentName/
├── index.vue          # 组件入口文件（必需）
├── types.ts           # 类型定义（可选，复杂组件推荐）
├── hooks.ts           # 组合式函数（可选，复杂逻辑推荐）
├── constants.ts       # 常量定义（可选）
└── README.md          # 组件文档（可选，复杂组件推荐）
```

### 📝 组件代码规范

#### 1. 单文件不超过500行

如果组件超过500行，需要拆分：

- 提取类型定义到 `types.ts`
- 提取业务逻辑到 `hooks.ts`
- 提取常量到 `constants.ts`
- 拆分子组件

#### 2. 组件命名规范

- 文件夹名：PascalCase（如 `BaseButton`, `BaseInput`）
- 入口文件：统一使用 `index.vue`
- 类型文件：统一使用 `types.ts`

#### 3. 代码结构顺序

```vue
<template>
  <!-- 模板内容 -->
</template>

<script setup lang="ts">
// 1. 导入语句
import { ref, computed } from 'vue'

// 2. 类型定义
interface Props {
  // ...
}

// 3. Props 和 Emits
const props = defineProps<Props>()
const emit = defineEmits<{
  // ...
}>()

// 4. 响应式数据
const state = ref()

// 5. 计算属性
const computed = computed(() => {})

// 6. 方法
const handleClick = () => {}

// 7. 生命周期钩子
onMounted(() => {})
</script>

<style scoped>
/* 样式 */
</style>
```

---

## 🧩 组件结构规范

### 📐 标准间距系统

#### 组件内边距 (padding)

| 尺寸 | 值           | 适用场景   |
| ---- | ------------ | ---------- |
| xs   | `p-1` (4px)  | 紧凑型组件 |
| sm   | `p-2` (8px)  | 小型组件   |
| md   | `p-3` (12px) | 标准组件   |
| lg   | `p-4` (16px) | 大型组件   |
| xl   | `p-6` (24px) | 特大型组件 |

#### 组件外边距 (margin)

| 尺寸 | 值           | 适用场景 |
| ---- | ------------ | -------- |
| sm   | `m-2` (8px)  | 紧凑布局 |
| md   | `m-4` (16px) | 标准布局 |
| lg   | `m-6` (24px) | 宽松布局 |

#### 元素间距 (gap)

| 尺寸 | 值             | 适用场景 |
| ---- | -------------- | -------- |
| sm   | `gap-2` (8px)  | 紧凑列表 |
| md   | `gap-3` (12px) | 标准列表 |
| lg   | `gap-4` (16px) | 宽松列表 |

### 🎨 颜色系统

#### 主题色彩

```css
/* 主色 - 用于主要按钮、链接 */
--color-primary: #3b82f6;
--color-primary-hover: #2563eb;
--color-primary-active: #1d4ed8;

/* 成功 - 用于成功状态、确认操作 */
--color-success: #10b981;
--color-success-hover: #059669;

/* 警告 - 用于警告提示 */
--color-warning: #f59e0b;
--color-warning-hover: #d97706;

/* 危险 - 用于危险操作、错误状态 */
--color-danger: #ef4444;
--color-danger-hover: #dc2626;

/* 中性 - 用于次要按钮、边框 */
--color-gray-50: #f9fafb;
--color-gray-100: #f3f4f6;
--color-gray-200: #e5e7eb;
--color-gray-300: #d1d5db;
--color-gray-400: #9ca3af;
--color-gray-500: #6b7280;
--color-gray-600: #4b5563;
--color-gray-700: #374151;
--color-gray-800: #1f2937;
--color-gray-900: #111827;
```

### 🔘 按钮变体系统

```typescript
type ButtonVariant =
  | 'primary' // 主要操作
  | 'secondary' // 次要操作
  | 'success' // 成功操作
  | 'danger' // 危险操作
  | 'warning' // 警告操作
  | 'ghost' // 幽灵按钮
  | 'outline' // 描边按钮
  | 'link' // 链接按钮
```

### 📏 尺寸系统

```typescript
type Size = 'xs' | 'sm' | 'md' | 'lg' | 'xl'

// 按钮尺寸
xs: h-6  text-xs  px-2
sm: h-8  text-sm  px-3
md: h-10 text-base px-4
lg: h-12 text-lg  px-6
xl: h-14 text-xl  px-8
```

---

## 🎬 动画效果规范

### ⚡ 动画时长

```css
/* 快速响应 */
transition: all 150ms ease;

/* 标准过渡 */
transition: all 200ms ease;

/* 入场动画 */
animation-duration: 300ms;
animation-timing-function: cubic-bezier(0.16, 1, 0.3, 1);
```

### 🎯 标准动画

```css
/* 淡入淡出 */
@keyframes fadeIn {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

/* 滑入 */
@keyframes slideIn {
  from {
    transform: translateY(-10px);
    opacity: 0;
  }
  to {
    transform: translateY(0);
    opacity: 1;
  }
}

/* 缩放 */
@keyframes scaleIn {
  from {
    transform: scale(0.95);
    opacity: 0;
  }
  to {
    transform: scale(1);
    opacity: 1;
  }
}
```

---

## 📱 响应式设计规范

### 📏 断点系统

```css
/* 小屏幕手机 */
@media (max-width: 640px) {
  /* 移动端优化 */
}

/* 平板 */
@media (max-width: 1024px) {
  /* 平板端优化 */
}

/* 桌面 */
@media (min-width: 1025px) {
  /* 桌面端优化 */
}
```

### 📱 移动端适配原则

1. **触控友好**: 确保所有交互元素至少44px×44px
2. **间距适中**: 移动端适当减少内边距，但保持可读性
3. **字体清晰**: 主要内容字体不小于14px
4. **按钮全宽**: 对话框底部按钮在移动端使用全宽布局

---

## 🧪 组件开发清单

### ✅ 必需功能

- [ ] Props 类型完整定义
- [ ] Emits 事件完整定义
- [ ] 必要的插槽（slot）支持
- [ ] 禁用状态支持
- [ ] 加载状态支持（如适用）
- [ ] 错误状态支持（如适用）

### 🎨 样式要求

- [ ] 支持多种尺寸（xs, sm, md, lg, xl）
- [ ] 支持多种变体（如适用）
- [ ] 响应式设计
- [ ] 悬停和聚焦状态
- [ ] 禁用状态样式

### ♿ 可访问性

- [ ] 适当的 ARIA 属性
- [ ] 键盘导航支持
- [ ] 焦点管理
- [ ] 屏幕阅读器友好

### 📝 文档要求

- [ ] Props 说明注释
- [ ] Emits 说明注释
- [ ] 使用示例（复杂组件）
- [ ] README.md（复杂组件）

---

## 📦 组件导出规范

### 统一导出方式

所有组件通过 `/src/components/index.ts` 统一导出：

```typescript
// 导入组件
import BaseButton from './BaseButton'
import BaseInput from './BaseInput'

// 组件映射表
const componentMap = {
  BaseButton,
  BaseInput,
  // ...
}

// 插件安装函数
export const createUI = (options = {}): Plugin => ({
  install(app: App) {
    Object.entries(componentMap).forEach(([name, component]) => {
      app.component(name, component)
    })
  },
})

// 默认插件
export default createUI()
```

### 全局注册方式

```typescript
// main.ts
import { createApp } from 'vue'
import CooperUI from './components'
import App from './App.vue'

const app = createApp(App)
app.use(CooperUI)
app.mount('#app')
```

### 按需引入方式

```vue
<script setup lang="ts">
import { BaseButton } from '@/components'
</script>
```

---

## 🔄 规范更新

- **版本**: v1.0
- **最后更新**: 2025-01-12
- **维护者**: Cooper Team

---

## 📚 参考示例

参考项目结构：`/Users/lilithgames/Desktop/saiboPic/cyber/web/src/components`

优秀组件示例：

- `BaseButton/index.vue` - 标准按钮组件
- `BaseInput/index.vue` - 标准输入组件
- `Dialog/index.vue` - 对话框组件

遵循本规范可以确保组件的一致性、可维护性和团队协作效率。
