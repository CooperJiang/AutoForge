<template>
  <div class="tree-node">
    <div
      class="flex items-center gap-2 px-2 py-1.5 rounded hover:bg-bg-hover transition-colors cursor-pointer group"
      :style="{ paddingLeft: `${level * 16 + 8}px` }"
      @click="handleClick"
    >
      <ChevronRight
        v-if="isExpandable"
        class="w-3.5 h-3.5 text-text-secondary transition-transform flex-shrink-0"
        :class="{ 'rotate-90': expanded }"
      />
      <div v-else class="w-3.5" />

      <component :is="getIcon()" class="w-3.5 h-3.5 flex-shrink-0" :class="getIconColor()" />

      <span class="font-mono text-xs font-medium text-text-primary flex-shrink-0">
        {{ name }}
      </span>

      <span class="text-xs text-text-tertiary truncate flex-1">
        {{ getDescription() }}
      </span>

      <button
        v-if="!isExpandable"
        @click.stop="handleCopy"
        class="opacity-0 group-hover:opacity-100 transition-opacity p-1 hover:bg-bg-elevated rounded"
        title="点击复制"
      >
        <Copy class="w-3.5 h-3.5 text-primary" />
      </button>
    </div>

    <div v-if="isExpandable && expanded" class="mt-0.5">
      <VariableTreeNode
        v-for="(childValue, childKey) in getChildren()"
        :key="childKey"
        :name="String(childKey)"
        :value="childValue"
        :path="`${path}.${childKey}`"
        :level="level + 1"
        @copy="$emit('copy', $event)"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ChevronRight, Copy, FileText, Folder, Hash, CheckSquare } from 'lucide-vue-next'

// 显式定义组件名称以支持递归
defineOptions({
  name: 'VariableTreeNode',
})

interface Props {
  name: string
  value: any
  path: string
  level: number
}

const props = defineProps<Props>()

const emit = defineEmits<{
  copy: [path: string]
}>()

const expanded = ref(false)

// 是否可展开（对象类型）
const isExpandable = computed(() => {
  return typeof props.value === 'object' && props.value !== null && props.name !== '_description'
})

// 获取图标
const getIcon = () => {
  if (isExpandable.value) {
    return Folder
  }

  if (typeof props.value === 'number') {
    return Hash
  }

  if (typeof props.value === 'boolean') {
    return CheckSquare
  }

  return FileText
}

// 获取图标颜色
const getIconColor = () => {
  if (isExpandable.value) {
    return 'text-amber-500'
  }

  if (typeof props.value === 'number') {
    return 'text-blue-500'
  }

  if (typeof props.value === 'boolean') {
    return 'text-green-500'
  }

  return 'text-purple-500'
}

// 获取描述
const getDescription = () => {
  if (typeof props.value === 'string') {
    return props.value
  }

  if (typeof props.value === 'number') {
    return `数字`
  }

  if (typeof props.value === 'boolean') {
    return `布尔值`
  }

  if (typeof props.value === 'object' && props.value !== null) {
    const desc = props.value._description
    if (desc) {
      return desc
    }
    return `对象`
  }

  return ''
}

// 获取子节点（过滤掉 _description）
const getChildren = () => {
  if (typeof props.value !== 'object' || props.value === null) {
    return {}
  }

  const children: Record<string, any> = {}
  for (const [key, val] of Object.entries(props.value)) {
    if (key !== '_description') {
      children[key] = val
    }
  }
  return children
}

// 处理点击
const handleClick = () => {
  if (isExpandable.value) {
    expanded.value = !expanded.value
  } else {
    handleCopy()
  }
}

// 处理复制
const handleCopy = () => {
  emit('copy', props.path)
}
</script>

<style scoped>
.tree-node {
  user-select: none;
}
</style>
