<template>
  <div>
    
    <div class="text-center mb-12">
      <h1 class="text-4xl font-bold text-text-primary mb-3">🔧 工具箱</h1>
      <p class="text-lg text-text-secondary">选择合适的工具，创建自动化任务</p>
    </div>

    
    <div
      v-if="!loading"
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-6 mb-8"
    >
      <ToolCard
        v-for="tool in tools"
        :key="tool.code"
        :tool="tool"
        @click="handleToolClick(tool)"
      />
    </div>

    
    <div v-else class="flex justify-center items-center py-20">
      <div class="text-text-tertiary">加载工具中...</div>
    </div>

    
    <div v-if="!loading && tools.length === 0" class="text-center py-20">
      <div class="text-text-placeholder text-lg">暂无可用工具</div>
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
