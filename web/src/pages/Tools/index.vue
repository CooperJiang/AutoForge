<template>
  <div class="flex flex-col h-full">
    <!-- Content Area with Background -->
    <div class="flex-1 flex flex-col bg-bg-elevated rounded-xl border border-border-primary overflow-hidden shadow-sm">
      <!-- Tool Grid - Scrollable -->
      <div class="flex-1 overflow-y-auto p-6">
        <div
          v-if="!loading && tools.length > 0"
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4"
        >
          <ToolCard
            v-for="tool in tools"
            :key="tool.code"
            :tool="tool"
            @click="handleToolClick(tool)"
          />
        </div>

        <!-- Empty State -->
        <div v-else-if="!loading && tools.length === 0" class="flex items-center justify-center h-full">
          <div class="text-text-placeholder text-center">
            <p class="text-lg">暂无可用工具</p>
          </div>
        </div>

        <!-- Loading -->
        <div v-else class="flex items-center justify-center h-full">
          <div class="text-text-tertiary">加载工具中...</div>
        </div>
      </div>
    </div>

    <ToolDetailDialog v-model="showDetailDialog" :tool="selectedTool" @use-tool="handleUseTool" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import ToolCard from './components/ToolCard.vue'
import ToolDetailDialog from './components/ToolDetailDialog.vue'
import * as toolApi from '@/api/tool'
import { message } from '@/utils/message'

interface Tool {
  code: string
  name: string
  description: string
  category: string
  icon: string
  tags: string[]
  version: string
  author: string
}

const router = useRouter()
const loading = ref(true)
const tools = ref<Tool[]>([])
const selectedTool = ref<Tool | null>(null)
const showDetailDialog = ref(false)

// 加载工具列表
const loadTools = async () => {
  try {
    loading.value = true
    const response = await toolApi.getToolList()
    // Parse tags from string to array
    tools.value = response.map((tool) => ({
      ...tool,
      tags: typeof tool.tags === 'string' ? JSON.parse(tool.tags) : tool.tags,
    }))
  } catch (error: any) {
    message.error('加载工具列表失败: ' + (error.message || '未知错误'))
  } finally {
    loading.value = false
  }
}

// 点击工具卡片
const handleToolClick = (tool: Tool) => {
  selectedTool.value = tool
  showDetailDialog.value = true
}

// 使用工具
const handleUseTool = (toolCode: string) => {
  showDetailDialog.value = false
  // 跳转到任务页面并选中该工具
  router.push({
    path: '/',
    query: { tool: toolCode },
  })
}

onMounted(() => {
  loadTools()
})
</script>
