# 🎨 Cooper 主题系统规范

> 完整的亮色/暗色主题系统，基于语义化设计和 CSS 变量实现

## 📐 设计原则

### 核心理念

1. **语义化优先**：颜色变量按用途命名，不按色值
2. **对比度达标**：确保 WCAG AA 级别可访问性
3. **渐进增强**：优先亮色模式，暗色模式作为增强
4. **用户优先**：尊重系统主题偏好，支持手动切换

---

## 🎨 颜色语义系统

### 1. 背景色 (Background)

#### 页面背景

| 变量名                 | 用途     | 亮色      | 暗色      |
| ---------------------- | -------- | --------- | --------- |
| `--color-bg-primary`   | 主背景   | `#ffffff` | `#0f172a` |
| `--color-bg-secondary` | 次要背景 | `#f8fafc` | `#1e293b` |
| `--color-bg-tertiary`  | 三级背景 | `#f1f5f9` | `#334155` |

#### 组件背景

| 变量名                | 用途         | 亮色              | 暗色               |
| --------------------- | ------------ | ----------------- | ------------------ |
| `--color-bg-elevated` | 卡片、对话框 | `#ffffff`         | `#1e293b`          |
| `--color-bg-overlay`  | 遮罩层       | `rgba(0,0,0,0.5)` | `rgba(0,0,0,0.75)` |
| `--color-bg-hover`    | 悬停背景     | `#f1f5f9`         | `#334155`          |
| `--color-bg-active`   | 激活背景     | `#e2e8f0`         | `#475569`          |

---

### 2. 文本色 (Text)

| 变量名                     | 用途     | 亮色      | 暗色      |
| -------------------------- | -------- | --------- | --------- |
| `--color-text-primary`     | 主要文本 | `#0f172a` | `#f8fafc` |
| `--color-text-secondary`   | 次要文本 | `#475569` | `#cbd5e1` |
| `--color-text-tertiary`    | 三级文本 | `#64748b` | `#94a3b8` |
| `--color-text-disabled`    | 禁用文本 | `#cbd5e1` | `#475569` |
| `--color-text-inverse`     | 反色文本 | `#ffffff` | `#0f172a` |
| `--color-text-placeholder` | 占位符   | `#94a3b8` | `#64748b` |

---

### 3. 边框色 (Border)

| 变量名                     | 用途     | 亮色      | 暗色      |
| -------------------------- | -------- | --------- | --------- |
| `--color-border-primary`   | 主边框   | `#e2e8f0` | `#334155` |
| `--color-border-secondary` | 次要边框 | `#cbd5e1` | `#475569` |
| `--color-border-focus`     | 聚焦边框 | `#3b82f6` | `#60a5fa` |
| `--color-border-error`     | 错误边框 | `#ef4444` | `#f87171` |

---

### 4. 品牌色 (Brand)

#### 主色 - Green

| 变量名                   | 用途     | 亮色      | 暗色      |
| ------------------------ | -------- | --------- | --------- |
| `--color-primary`        | 主色     | `#10b981` | `#34d399` |
| `--color-primary-hover`  | 主色悬停 | `#059669` | `#10b981` |
| `--color-primary-active` | 主色激活 | `#047857` | `#059669` |
| `--color-primary-light`  | 主色浅色 | `#d1fae5` | `#064e3b` |
| `--color-primary-text`   | 主色文本 | `#ffffff` | `#0f172a` |

#### 辅助色 - Cyan

| 变量名                 | 用途       | 亮色      | 暗色      |
| ---------------------- | ---------- | --------- | --------- |
| `--color-accent`       | 强调色     | `#06b6d4` | `#22d3ee` |
| `--color-accent-hover` | 强调色悬停 | `#0891b2` | `#06b6d4` |
| `--color-accent-light` | 强调色浅色 | `#cffafe` | `#164e63` |

---

### 5. 功能色 (Functional)

#### 成功 - Green

| 变量名                  | 用途     | 亮色      | 暗色      |
| ----------------------- | -------- | --------- | --------- |
| `--color-success`       | 成功主色 | `#10b981` | `#34d399` |
| `--color-success-hover` | 成功悬停 | `#059669` | `#10b981` |
| `--color-success-light` | 成功背景 | `#d1fae5` | `#064e3b` |
| `--color-success-text`  | 成功文本 | `#065f46` | `#a7f3d0` |

#### 警告 - Yellow

| 变量名                  | 用途     | 亮色      | 暗色      |
| ----------------------- | -------- | --------- | --------- |
| `--color-warning`       | 警告主色 | `#f59e0b` | `#fbbf24` |
| `--color-warning-hover` | 警告悬停 | `#d97706` | `#f59e0b` |
| `--color-warning-light` | 警告背景 | `#fef3c7` | `#78350f` |
| `--color-warning-text`  | 警告文本 | `#92400e` | `#fde68a` |

#### 错误 - Red

| 变量名                | 用途     | 亮色      | 暗色      |
| --------------------- | -------- | --------- | --------- |
| `--color-error`       | 错误主色 | `#ef4444` | `#f87171` |
| `--color-error-hover` | 错误悬停 | `#dc2626` | `#ef4444` |
| `--color-error-light` | 错误背景 | `#fee2e2` | `#7f1d1d` |
| `--color-error-text`  | 错误文本 | `#991b1b` | `#fecaca` |

#### 信息 - Blue

| 变量名               | 用途     | 亮色      | 暗色      |
| -------------------- | -------- | --------- | --------- |
| `--color-info`       | 信息主色 | `#0ea5e9` | `#38bdf8` |
| `--color-info-hover` | 信息悬停 | `#0284c7` | `#0ea5e9` |
| `--color-info-light` | 信息背景 | `#e0f2fe` | `#075985` |
| `--color-info-text`  | 信息文本 | `#075985` | `#bae6fd` |

---

### 6. 阴影 (Shadow)

| 变量名        | 用途     | 亮色                                | 暗色                                |
| ------------- | -------- | ----------------------------------- | ----------------------------------- |
| `--shadow-sm` | 小阴影   | `0 1px 2px 0 rgb(0 0 0 / 0.05)`     | `0 1px 2px 0 rgb(0 0 0 / 0.5)`      |
| `--shadow-md` | 中阴影   | `0 4px 6px -1px rgb(0 0 0 / 0.1)`   | `0 4px 6px -1px rgb(0 0 0 / 0.6)`   |
| `--shadow-lg` | 大阴影   | `0 10px 15px -3px rgb(0 0 0 / 0.1)` | `0 10px 15px -3px rgb(0 0 0 / 0.7)` |
| `--shadow-xl` | 特大阴影 | `0 20px 25px -5px rgb(0 0 0 / 0.1)` | `0 20px 25px -5px rgb(0 0 0 / 0.8)` |

---

## 🧩 组件专用颜色

### 按钮 (Button)

#### Primary Button

```css
/* 亮色模式 */
--btn-primary-bg: var(--color-primary);
--btn-primary-hover: var(--color-primary-hover);
--btn-primary-text: var(--color-primary-text);

/* 暗色模式 */
--btn-primary-bg: var(--color-primary);
--btn-primary-hover: var(--color-primary-hover);
--btn-primary-text: var(--color-primary-text);
```

#### Secondary Button

```css
/* 亮色模式 */
--btn-secondary-bg: #f1f5f9;
--btn-secondary-hover: #e2e8f0;
--btn-secondary-text: #0f172a;

/* 暗色模式 */
--btn-secondary-bg: #334155;
--btn-secondary-hover: #475569;
--btn-secondary-text: #f8fafc;
```

#### Ghost Button

```css
/* 亮色模式 */
--btn-ghost-hover: #f1f5f9;
--btn-ghost-text: #475569;

/* 暗色模式 */
--btn-ghost-hover: #334155;
--btn-ghost-text: #cbd5e1;
```

### 输入框 (Input)

```css
/* 亮色模式 */
--input-bg: #ffffff;
--input-border: #e2e8f0;
--input-border-focus: #3b82f6;
--input-text: #0f172a;
--input-placeholder: #94a3b8;

/* 暗色模式 */
--input-bg: #1e293b;
--input-border: #334155;
--input-border-focus: #60a5fa;
--input-text: #f8fafc;
--input-placeholder: #64748b;
```

### 卡片 (Card)

```css
/* 亮色模式 */
--card-bg: #ffffff;
--card-border: #e2e8f0;
--card-shadow: var(--shadow-sm);

/* 暗色模式 */
--card-bg: #1e293b;
--card-border: #334155;
--card-shadow: var(--shadow-md);
```

---

## 📱 特殊状态颜色

### 禁用状态

```css
--color-disabled-bg: #f1f5f9; /* 亮色 */
--color-disabled-bg: #334155; /* 暗色 */
--color-disabled-text: #cbd5e1; /* 亮色 */
--color-disabled-text: #475569; /* 暗色 */
```

### 选中状态

```css
--color-selected-bg: #dbeafe; /* 亮色 */
--color-selected-bg: #1e3a8a; /* 暗色 */
--color-selected-border: #3b82f6; /* 亮色 */
--color-selected-border: #60a5fa; /* 暗色 */
```

---

## 🎯 使用规范

### 1. 在 Tailwind 中使用

```html
<!-- 背景色 -->
<div class="bg-bg-primary text-text-primary">内容</div>

<!-- 边框 -->
<div class="border border-border-primary">内容</div>

<!-- 按钮 -->
<button class="bg-primary hover:bg-primary-hover text-primary-text">按钮</button>
```

### 2. 在 CSS 中使用

```css
.custom-component {
  background-color: var(--color-bg-primary);
  color: var(--color-text-primary);
  border: 1px solid var(--color-border-primary);
}

.custom-component:hover {
  background-color: var(--color-bg-hover);
}
```

### 3. 禁止使用硬编码颜色

❌ **错误示例**：

```html
<div class="bg-white text-gray-900">
  <div class="bg-slate-50 border-slate-200">
    <div style="background: #ffffff"></div>
  </div>
</div>
```

✅ **正确示例**：

```html
<div class="bg-bg-primary text-text-primary">
  <div class="bg-bg-secondary border-border-primary">
    <div class="bg-bg-elevated text-text-primary"></div>
  </div>
</div>
```

---

## 🔄 主题切换实现

### 主题检测优先级

1. 用户手动选择（存储在 localStorage）
2. 系统主题偏好（`prefers-color-scheme`）
3. 默认主题（亮色）

### 切换方式

- 通过修改 `<html>` 标签的 `data-theme` 属性
- 值为 `light` 或 `dark`

```html
<html data-theme="light">
  <!-- 亮色模式 -->
  <html data-theme="dark">
    <!-- 暗色模式 -->
  </html>
</html>
```

---

## 📊 对比度检查

所有颜色组合需满足 WCAG AA 级别：

- 正常文本：对比度 ≥ 4.5:1
- 大号文本：对比度 ≥ 3:1
- UI组件：对比度 ≥ 3:1

---

## 🎨 品牌渐变

### 主渐变 (Green to Cyan)

```css
/* 亮色模式 */
--gradient-primary: linear-gradient(135deg, #10b981 0%, #06b6d4 100%);

/* 暗色模式 */
--gradient-primary: linear-gradient(135deg, #34d399 0%, #22d3ee 100%);
```

### 页面背景渐变

```css
/* 亮色模式 */
--gradient-bg: linear-gradient(to bottom right, #f1f5f9 0%, #ecfdf5 50%, #cffafe 100%);

/* 暗色模式 */
--gradient-bg: linear-gradient(to bottom right, #0f172a 0%, #064e3b 50%, #164e63 100%);
```

---

## 📝 文档版本

- **版本**: v1.0
- **创建日期**: 2025-01-13
- **维护者**: Cooper Team
- **状态**: Draft → Review → Active

---

## 🔗 相关文档

- [组件设计规范](../components/COMPONENT_DESIGN_SPEC.md)
- [Tailwind 配置](../../tailwind.config.js)
- [主题 Composable](../composables/useTheme.ts)
