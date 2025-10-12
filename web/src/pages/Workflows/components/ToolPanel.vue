<template>
  <div class="w-64 bg-white border-r border-slate-200 flex flex-col overflow-hidden">
    <!-- 面板标题 -->
    <div class="px-4 py-3 border-b border-slate-200">
      <h3 class="text-sm font-semibold text-slate-900">工具箱</h3>
      <p class="text-xs text-slate-500 mt-1">点击添加到画布</p>
    </div>

    <!-- 工具列表 -->
    <div class="flex-1 overflow-y-auto p-3 space-y-2">
      <!-- 触发器分类 -->
      <div class="mb-4">
        <div class="text-xs font-semibold text-slate-600 mb-2 px-2">触发器</div>
        <button
          @click="handleAddTrigger"
          draggable="true"
          @dragstart="handleDragStartTrigger($event)"
          class="w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg border border-slate-200 hover:border-blue-400 hover:bg-blue-50 transition-all group cursor-move"
        >
          <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center text-white shadow-sm">
            <Clock class="w-4 h-4" />
          </div>
          <div class="flex-1 text-left min-w-0">
            <div class="text-sm font-medium text-slate-900 truncate">定时触发</div>
            <div class="text-xs text-slate-500 truncate">按计划执行</div>
          </div>
        </button>
      </div>

      <!-- 流程控制 -->
      <div class="mb-4">
        <div class="text-xs font-semibold text-slate-600 mb-2 px-2">流程控制</div>
        <div class="space-y-2">
          <!-- 条件判断 -->
          <button
            @click="handleAddCondition"
            draggable="true"
            @dragstart="handleDragStartCondition($event)"
            class="w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg border border-slate-200 hover:border-amber-400 hover:bg-amber-50 transition-all group cursor-move"
          >
            <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-amber-400 to-orange-500 flex items-center justify-center text-white shadow-sm">
              <GitBranch class="w-4 h-4" />
            </div>
            <div class="flex-1 text-left min-w-0">
              <div class="text-sm font-medium text-slate-900 truncate">条件判断</div>
              <div class="text-xs text-slate-500 truncate">IF 分支</div>
            </div>
          </button>

          <!-- 开关分支 -->
          <button
            @click="handleAddSwitch"
            draggable="true"
            @dragstart="handleDragStartSwitch($event)"
            class="w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg border border-slate-200 hover:border-indigo-400 hover:bg-indigo-50 transition-all group cursor-move"
          >
            <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-indigo-400 to-indigo-600 flex items-center justify-center text-white shadow-sm">
              <Split class="w-4 h-4" />
            </div>
            <div class="flex-1 text-left min-w-0">
              <div class="text-sm font-medium text-slate-900 truncate">开关分支</div>
              <div class="text-xs text-slate-500 truncate">Switch 多路</div>
            </div>
          </button>

          <!-- 延迟 -->
          <button
            @click="handleAddDelay"
            draggable="true"
            @dragstart="handleDragStartDelay($event)"
            class="w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg border border-slate-200 hover:border-purple-400 hover:bg-purple-50 transition-all group cursor-move"
          >
            <div class="flex-shrink-0 w-8 h-8 rounded-lg bg-gradient-to-br from-purple-400 to-purple-600 flex items-center justify-center text-white shadow-sm">
              <Timer class="w-4 h-4" />
            </div>
            <div class="flex-1 text-left min-w-0">
              <div class="text-sm font-medium text-slate-900 truncate">延迟等待</div>
              <div class="text-xs text-slate-500 truncate">等待指定时间</div>
            </div>
          </button>
        </div>
      </div>

      <!-- 工具分类 -->
      <div>
        <div class="text-xs font-semibold text-slate-600 mb-2 px-2">工具</div>
        <div class="space-y-2">
          <button
            v-for="tool in tools"
            :key="tool.code"
            @click="handleAddTool(tool)"
            draggable="true"
            @dragstart="handleDragStart($event, tool)"
            class="w-full flex items-center gap-2.5 px-3 py-2.5 rounded-lg border border-slate-200 hover:border-blue-400 hover:bg-blue-50 transition-all group cursor-move"
          >
            <div
              :class="[
                'flex-shrink-0 w-8 h-8 rounded-lg flex items-center justify-center text-white shadow-sm',
                getToolBgClass(tool.code)
              ]"
            >
              <component :is="getToolIcon(tool.code)" class="w-4 h-4" />
            </div>
            <div class="flex-1 text-left min-w-0">
              <div class="text-sm font-medium text-slate-900 truncate">{{ tool.name }}</div>
              <div class="text-xs text-slate-500 truncate">{{ tool.description }}</div>
            </div>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Clock, Globe, Mail, HeartPulse, GitBranch, Split, Timer } from 'lucide-vue-next'
import * as toolApi from '@/api/tool'

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

// 拖拽开始 - 触发器
const handleDragStartTrigger = (event: DragEvent) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData('application/vueflow', JSON.stringify({
      toolCode: 'trigger',
      toolName: '定时触发',
      nodeType: 'trigger'
    }))
  }
}

// 拖拽开始 - 条件判断
const handleDragStartCondition = (event: DragEvent) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData('application/vueflow', JSON.stringify({
      toolCode: 'condition',
      toolName: '条件判断',
      nodeType: 'condition'
    }))
  }
}

// 拖拽开始 - 开关分支
const handleDragStartSwitch = (event: DragEvent) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData('application/vueflow', JSON.stringify({
      toolCode: 'switch',
      toolName: '开关分支',
      nodeType: 'switch'
    }))
  }
}

// 拖拽开始 - 延迟
const handleDragStartDelay = (event: DragEvent) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData('application/vueflow', JSON.stringify({
      toolCode: 'delay',
      toolName: '延迟等待',
      nodeType: 'delay'
    }))
  }
}

// 拖拽开始 - 工具
const handleDragStart = (event: DragEvent, tool: any) => {
  if (event.dataTransfer) {
    event.dataTransfer.effectAllowed = 'copy'
    event.dataTransfer.setData('application/vueflow', JSON.stringify({
      toolCode: tool.code,
      toolName: tool.name,
      nodeType: 'tool'
    }))
  }
}

// 获取工具图标
const getToolIcon = (code: string) => {
  const iconMap: Record<string, any> = {
    'http_request': Globe,
    'email_sender': Mail,
    'health_checker': HeartPulse
  }
  return iconMap[code] || Globe
}

// 获取工具背景色
const getToolBgClass = (code: string) => {
  const colorMap: Record<string, string> = {
    'http_request': 'bg-gradient-to-br from-blue-500 to-purple-600',
    'email_sender': 'bg-gradient-to-br from-purple-500 to-pink-600',
    'health_checker': 'bg-gradient-to-br from-indigo-500 to-blue-600'
  }
  return colorMap[code] || 'bg-gradient-to-br from-blue-500 to-purple-600'
}

onMounted(() => {
  loadTools()
})
</script>
