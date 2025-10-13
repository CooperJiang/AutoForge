<template>
  <div class="relative">
    <div class="relative">
      <input
        ref="inputRef"
        v-model="localValue"
        :type="type"
        :placeholder="placeholder"
        @input="handleInput"
        @keydown="handleKeydown"
        class="w-full px-3 py-1.5 pr-20 text-sm text-text-primary bg-bg-primary border-2 border-border-primary rounded-md transition-all duration-200 focus:border-border-focus focus:ring-2 focus:ring-primary-light focus:outline-none hover:border-border-secondary placeholder:text-text-placeholder font-mono"
      />
      <button
        type="button"
        @click="showVariablePicker = !showVariablePicker"
        class="absolute right-2 top-1/2 -translate-y-1/2 px-2 py-0.5 text-xs font-medium text-text-secondary hover:text-primary border border-border-primary rounded hover:border-primary transition-colors"
        title="插入变量"
      >
        <Variable class="w-3.5 h-3.5 inline-block mr-1" />
        变量
      </button>
    </div>

    <!-- 变量选择器下拉框 -->
    <div
      v-if="showVariablePicker"
      class="absolute z-50 mt-1 w-full bg-bg-elevated border-2 border-border-primary rounded-lg shadow-lg max-h-80 overflow-y-auto"
    >
      <div class="p-2 border-b border-border-primary">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="搜索变量..."
          class="w-full px-3 py-1.5 text-sm text-text-primary bg-bg-primary border-2 border-border-primary rounded-md transition-all duration-200 focus:border-border-focus focus:ring-2 focus:ring-primary-light focus:outline-none hover:border-border-secondary placeholder:text-text-placeholder"
        />
      </div>

      <!-- 环境变量 -->
      <div v-if="filteredEnvVars.length > 0" class="border-b border-border-primary">
        <div class="px-3 py-2 text-xs font-semibold text-text-secondary bg-bg-hover">环境变量</div>
        <button
          v-for="envVar in filteredEnvVars"
          :key="envVar.key"
          type="button"
          @click="insertVariable(`env.${envVar.key}`)"
          class="w-full px-3 py-2 text-left hover:bg-success-light transition-colors group"
        >
          <div class="flex items-center justify-between">
            <div class="flex-1 min-w-0">
              <div class="text-sm font-mono text-text-primary truncate">&#123;&#123;env.{{ envVar.key }}&#125;&#125;</div>
              <div class="text-xs text-text-tertiary truncate">{{ envVar.description }}</div>
            </div>
          </div>
        </button>
      </div>

      <!-- 触发器数据 -->
      <div v-if="showTriggerData" class="border-b border-border-primary">
        <div class="px-3 py-2 text-xs font-semibold text-text-secondary bg-bg-hover">触发器</div>
        <button
          v-for="field in filteredTriggerFields"
          :key="field.key"
          type="button"
          @click="insertVariable(`trigger.${field.key}`)"
          class="w-full px-3 py-2 text-left hover:bg-success-light transition-colors"
        >
          <div class="text-sm font-mono text-text-primary truncate">&#123;&#123;trigger.{{ field.key }}&#125;&#125;</div>
          <div class="text-xs text-text-tertiary truncate">{{ field.description }}</div>
        </button>
      </div>

      <!-- 前置节点输出 -->
      <div v-if="filteredPreviousNodes.length > 0">
        <div class="px-3 py-2 text-xs font-semibold text-text-secondary bg-bg-hover">前置节点输出</div>
        <div v-for="node in filteredPreviousNodes" :key="node.id" class="border-b border-border-primary last:border-0">
          <div class="px-3 py-2 text-xs font-semibold text-text-secondary bg-bg-hover/50">
            {{ node.name }} ({{ node.id }})
          </div>
          <button
            v-for="(desc, field) in node.outputs"
            :key="field"
            type="button"
            @click="insertVariable(`${node.id}.${field}`)"
            class="w-full px-3 py-2 pl-6 text-left hover:bg-success-light transition-colors"
          >
            <div class="text-sm font-mono text-text-primary truncate">&#123;&#123;{{ node.id }}.{{ field }}&#125;&#125;</div>
            <div class="text-xs text-text-tertiary truncate">{{ desc }}</div>
          </button>
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="filteredEnvVars.length === 0 && filteredPreviousNodes.length === 0 && !showTriggerData" class="px-3 py-8 text-center text-sm text-text-tertiary">
        暂无可用变量
      </div>
    </div>

    <!-- 点击外部关闭 -->
    <div
      v-if="showVariablePicker"
      class="fixed inset-0 z-40"
      @click="showVariablePicker = false"
    ></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { Variable } from 'lucide-vue-next'
import { getNodeOutputSchema } from '@/utils/variableParser'

interface Props {
  modelValue: string
  placeholder?: string
  type?: string
  previousNodes?: Array<{
    id: string
    name: string
    type: string
    toolCode?: string
  }>
  showTriggerData?: boolean
  envVars?: Array<{ key: string; description: string }>
}

const props = withDefaults(defineProps<Props>(), {
  placeholder: '输入值或选择变量',
  type: 'text',
  previousNodes: () => [],
  showTriggerData: false,
  envVars: () => []
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
}>()

const localValue = ref(props.modelValue)
const showVariablePicker = ref(false)
const searchQuery = ref('')
const inputRef = ref<HTMLInputElement | null>(null)

watch(() => props.modelValue, (val) => {
  localValue.value = val
})

const handleInput = () => {
  emit('update:modelValue', localValue.value)
}

// 触发器字段
const triggerFields = [
  { key: 'timestamp', description: '触发时间戳' },
  { key: 'type', description: '触发类型' }
]

// 过滤环境变量
const filteredEnvVars = computed(() => {
  if (!searchQuery.value) return props.envVars
  return props.envVars.filter(v =>
    v.key.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    v.description.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// 过滤触发器字段
const filteredTriggerFields = computed(() => {
  if (!searchQuery.value) return triggerFields
  return triggerFields.filter(f =>
    f.key.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    f.description.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// 过滤前置节点
const filteredPreviousNodes = computed(() => {
  const nodes = props.previousNodes.map(node => ({
    ...node,
    outputs: getNodeOutputSchema(node.type, node.toolCode)
  }))

  if (!searchQuery.value) return nodes

  return nodes.filter(node =>
    node.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    node.id.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
    Object.keys(node.outputs).some(key =>
      key.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  )
})

// 插入变量
const insertVariable = (variable: string) => {
  const input = inputRef.value
  if (!input) return

  const cursorPos = input.selectionStart || 0
  const textBefore = localValue.value.substring(0, cursorPos)
  const textAfter = localValue.value.substring(cursorPos)

  localValue.value = `${textBefore}{{${variable}}}${textAfter}`
  emit('update:modelValue', localValue.value)

  showVariablePicker.value = false
  searchQuery.value = ''

  // 恢复焦点并移动光标
  setTimeout(() => {
    input.focus()
    const newPos = cursorPos + variable.length + 4 // {{variable}}
    input.setSelectionRange(newPos, newPos)
  }, 0)
}

// 快捷键支持
const handleKeydown = (e: KeyboardEvent) => {
  // Cmd/Ctrl + K 打开变量选择器
  if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
    e.preventDefault()
    showVariablePicker.value = !showVariablePicker.value
  }
}
</script>
