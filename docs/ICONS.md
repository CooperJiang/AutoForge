# 图标使用指南

## Lucide Vue Next

项目已安装 `lucide-vue-next` 图标库，提供了1000+精美的图标。

### 使用方法

```vue
<script setup lang="ts">
import { RefreshCw, Trash2, Copy, Check, AlertCircle } from 'lucide-vue-next'
</script>

<template>
  <!-- 基础使用 -->
  <RefreshCw :size="16" />

  <!-- 自定义颜色和大小 -->
  <Trash2 :size="20" color="#ef4444" />

  <!-- 添加class -->
  <Copy :size="16" class="text-blue-500 hover:text-blue-700" />

  <!-- 旋转动画 -->
  <RefreshCw :size="16" class="animate-spin" />
</template>
```

### 常用图标

- **操作类**: `Trash2`, `Edit2`, `Copy`, `Check`, `X`, `Plus`, `Minus`
- **导航类**: `ChevronLeft`, `ChevronRight`, `ChevronUp`, `ChevronDown`, `ArrowLeft`, `ArrowRight`
- **状态类**: `CheckCircle`, `XCircle`, `AlertCircle`, `Info`, `Loader2`
- **刷新类**: `RefreshCw`, `RotateCw`, `RotateCcw`
- **文件类**: `File`, `Folder`, `Upload`, `Download`
- **其他**: `Search`, `Settings`, `Eye`, `EyeOff`, `Calendar`, `Clock`

### 在线查找图标

访问 https://lucide.dev/icons/ 查看所有可用图标

### 示例：替换现有SVG

```vue
<!-- 旧的手写SVG -->
<svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
</svg>

<!-- 新的 Lucide 图标 -->
<RefreshCw :size="16" />
```
