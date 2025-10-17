<template>
  <div class="w-64 bg-bg-elevated border-r border-border-primary flex flex-col overflow-hidden">
    <div class="px-4 py-3 border-b border-border-primary">
      <h3 class="text-sm font-semibold text-text-primary">工具箱</h3>
      <p class="text-xs text-text-tertiary mt-1">点击添加到画布</p>
    </div>

    <div class="flex-1 overflow-y-auto p-3 space-y-2">
      <div class="mb-4">
        <div class="text-xs font-semibold text-text-secondary mb-2 px-2">触发器</div>
        <div class="space-y-2">
          <button
            @click="handleAddExternalTrigger"
            draggable="true"
            @dragstart="handleDragStartExternalTrigger($event)"
            class="w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg border border-border-primary hover:border-blue-500 hover:bg-blue-500/10 transition-all group cursor-move"
          >
            <div
              class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500 to-blue-600 flex items-center justify-center text-white shadow-sm"
            >
              <Globe class="w-4 h-4" />
            </div>
            <div class="flex-1 text-left min-w-0">
              <div class="text-sm font-medium text-text-primary truncate">外部 API 触发</div>
              <div class="text-xs text-text-tertiary truncate">接收外部参数</div>
            </div>
          </button>

          <button
            @click="handleAddTrigger"
            draggable="true"
            @dragstart="handleDragStartTrigger($event)"
            class="w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg border border-border-primary hover:border-primary hover:bg-primary-light transition-all group cursor-move"
          >
            <div
              class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-primary to-accent flex items-center justify-center text-white shadow-sm"
            >
              <Clock class="w-4 h-4" />
            </div>
            <div class="flex-1 text-left min-w-0">
              <div class="text-sm font-medium text-text-primary truncate">定时触发</div>
              <div class="text-xs text-text-tertiary truncate">按计划执行</div>
            </div>
          </button>
        </div>
      </div>

      <div class="mb-4">
        <div class="text-xs font-semibold text-text-secondary mb-2 px-2">流程控制</div>
        <div class="space-y-2">
          <button
            @click="handleAddCondition"
            draggable="true"
            @dragstart="handleDragStartCondition($event)"
            class="w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg border border-border-primary hover:border-warning hover:bg-warning-light transition-all group cursor-move"
          >
            <div
              class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-amber-400 to-orange-500 flex items-center justify-center text-white shadow-sm"
            >
              <GitBranch class="w-4 h-4" />
            </div>
            <div class="flex-1 text-left min-w-0">
              <div class="text-sm font-medium text-text-primary truncate">条件判断</div>
              <div class="text-xs text-text-tertiary truncate">IF 分支</div>
            </div>
          </button>

          <button
            @click="handleAddSwitch"
            draggable="true"
            @dragstart="handleDragStartSwitch($event)"
            class="w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg border border-border-primary hover:border-accent hover:bg-accent/10 transition-all group cursor-move"
          >
            <div
              class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-accent to-accent-hover flex items-center justify-center text-white shadow-sm"
            >
              <Split class="w-4 h-4" />
            </div>
            <div class="flex-1 text-left min-w-0">
              <div class="text-sm font-medium text-text-primary truncate">开关分支</div>
              <div class="text-xs text-text-tertiary truncate">Switch 多路</div>
            </div>
          </button>

          <button
            @click="handleAddDelay"
            draggable="true"
            @dragstart="handleDragStartDelay($event)"
            class="w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg border border-border-primary hover:border-accent hover:bg-accent-light transition-all group cursor-move"
          >
            <div
              class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-accent to-accent-hover flex items-center justify-center text-white shadow-sm"
            >
              <Timer class="w-4 h-4" />
            </div>
            <div class="flex-1 text-left min-w-0">
              <div class="text-sm font-medium text-text-primary truncate">延迟等待</div>
              <div class="text-xs text-text-tertiary truncate">等待指定时间</div>
            </div>
          </button>
        </div>
      </div>

      <div>
        <div class="text-xs font-semibold text-text-secondary mb-2 px-2">工具</div>
        <div class="space-y-2">
          <button
            v-for="tool in tools"
            :key="tool.code"
            @click="handleAddTool(tool)"
            draggable="true"
            @dragstart="handleDragStart($event, tool)"
            class="w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg border border-border-primary hover:border-primary hover:bg-primary-light transition-all group cursor-move"
          >
            <div
              :class="[
                'flex-shrink-0 w-8 h-8 rounded-lg flex items-center justify-center text-white shadow-sm',
                getToolIconBg(tool.code),
              ]"
            >
              <component
                v-if="isLucideIcon(tool.code)"
                :is="getToolIcon(tool.code)"
                class="w-4 h-4"
              />
              <img v-else :src="getToolIcon(tool.code)" alt="" class="w-4 h-4 object-contain" />
            </div>
            <div class="flex-1 text-left min-w-0">
              <div class="text-sm font-medium text-text-primary truncate">{{ tool.name }}</div>
              <div class="text-xs text-text-tertiary truncate">{{ tool.description }}</div>
            </div>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Clock, Globe, GitBranch, Split, Timer } from 'lucide-vue-next'
import * as toolApi from '@/api/tool'
import { getToolIcon, getToolIconBg } from '@/config/tools'

const emit = defineEmits<{
  addNode: [toolCode: string, toolName: string, nodeType?: string]
}>()

const tools = ref<any[]>([])

// 加载工具列表
const loadTools = async () => {
  try {
    tools.value = await toolApi.getToolList()
  } catch (error) {
    console.error('Failed to load tools:', error)
  }
}

// 添加外部 API 触发器
const handleAddExternalTrigger = () => {
  emit('addNode', 'external_trigger', '外部 API 触发', 'external_trigger')
}

// 添加触发器
const handleAddTrigger = () => {
  emit('addNode', 'trigger', '定时触发', 'trigger')
}

// 添加条件判断
const handleAddCondition = () => {
  emit('addNode', 'condition', '条件判断', 'condition')
}

// 添加开关分支
const handleAddSwitch = () => {
  emit('addNode', 'switch', '开关分支', 'switch')
}

// 添加延迟
const handleAddDelay = () => {
  emit('addNode', 'delay', '延迟等待', 'delay')
}

// 添加工具
const handleAddTool = (tool: any) => {
  emit('addNode', tool.code, tool.name, 'tool')
}

// 拖拽开始 - 外部 API 触发器
const handleDragStartExternalTrigger = (event: DragEvent) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData(
      'application/vueflow',
      JSON.stringify({
        toolCode: 'external_trigger',
        toolName: '外部 API 触发',
        nodeType: 'external_trigger',
      })
    )
  }
}

// 拖拽开始 - 触发器
const handleDragStartTrigger = (event: DragEvent) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData(
      'application/vueflow',
      JSON.stringify({
        toolCode: 'trigger',
        toolName: '定时触发',
        nodeType: 'trigger',
      })
    )
  }
}

// 拖拽开始 - 条件判断
const handleDragStartCondition = (event: DragEvent) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData(
      'application/vueflow',
      JSON.stringify({
        toolCode: 'condition',
        toolName: '条件判断',
        nodeType: 'condition',
      })
    )
  }
}

// 拖拽开始 - 开关分支
const handleDragStartSwitch = (event: DragEvent) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData(
      'application/vueflow',
      JSON.stringify({
        toolCode: 'switch',
        toolName: '开关分支',
        nodeType: 'switch',
      })
    )
  }
}

// 拖拽开始 - 延迟
const handleDragStartDelay = (event: DragEvent) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData(
      'application/vueflow',
      JSON.stringify({
        toolCode: 'delay',
        toolName: '延迟等待',
        nodeType: 'delay',
      })
    )
  }
}

// 拖拽开始 - 工具
const handleDragStart = (event: DragEvent, tool: any) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData(
      'application/vueflow',
      JSON.stringify({
        toolCode: tool.code,
        toolName: tool.name,
        nodeType: 'tool',
      })
    )
  }
}

// 判断是否为 Lucide 图标
const isLucideIcon = (code: string) => {
  const icon = getToolIcon(code)
  return typeof icon !== 'string'
}

onMounted(() => {
  loadTools()
})
</script>
