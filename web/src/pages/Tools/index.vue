<template>
  <div class="flex flex-col h-full">
    <!-- Filters -->
    <div class="bg-bg-elevated rounded-xl border border-border-primary p-4 mb-4 shadow-sm flex-shrink-0">
      <div class="flex gap-4">
        <!-- Category Filter -->
        <BaseSelect
          v-model="selectedCategory"
          :options="categoryOptions"
          style="width: 260px"
          placeholder="选择分类"
          @update:modelValue="handleFilterChange"
        />

        <!-- Search -->
        <BaseInput
          v-model="searchKeyword"
          placeholder="搜索工具..."
          style="width: 260px"
          @keyup.enter="handleFilterChange"
        >
          <template #prefix>
            <Search class="w-4 h-4" />
          </template>
        </BaseInput>
      </div>
    </div>

    <!-- Content Area with Background -->
    <div class="flex-1 flex flex-col bg-bg-elevated rounded-xl border border-border-primary overflow-hidden shadow-sm">
      <!-- Tool Grid - Scrollable -->
      <div class="flex-1 overflow-y-auto p-4">
        <div
          v-if="!loading && filteredTools.length > 0"
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-3"
        >
          <ToolCard
            v-for="tool in paginatedTools"
            :key="tool.code"
            :tool="tool"
            @click="handleToolClick(tool)"
          />
        </div>

        <!-- Empty State -->
        <div v-else-if="!loading && filteredTools.length === 0" class="flex items-center justify-center h-full">
          <div class="text-text-placeholder text-center">
            <Package class="w-16 h-16 mx-auto mb-4" />
            <p class="text-lg">暂无工具</p>
            <p class="text-sm">暂时没有符合条件的工具</p>
          </div>
        </div>

        <!-- Loading -->
        <div v-else class="flex items-center justify-center h-full">
          <div class="text-text-tertiary">加载工具中...</div>
        </div>
      </div>

      <!-- Pagination - Fixed at Bottom -->
      <div
        v-if="!loading && filteredTools.length > 0"
        class="border-t border-border-primary px-4 py-3 flex-shrink-0"
      >
        <Pagination
          :current="currentPage"
          :total="filteredTools.length"
          :page-size="pageSize"
          @change="handlePageChange"
        />
      </div>
    </div>

    <ToolDetailDialog v-model="showDetailDialog" :tool="selectedTool" @use-tool="handleUseTool" />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Package } from 'lucide-vue-next'
import BaseSelect from '@/components/BaseSelect'
import BaseInput from '@/components/BaseInput'
import Pagination from '@/components/Pagination'
import ToolCard from './components/ToolCard.vue'
import ToolDetailDialog from './components/ToolDetailDialog.vue'
import * as toolApi from '@/api/tool'
import * as toolConfigApi from '@/api/toolConfig'
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
const categories = ref<toolConfigApi.ToolCategory[]>([])
const selectedTool = ref<Tool | null>(null)
const showDetailDialog = ref(false)

// 筛选和分页
const selectedCategory = ref('all')
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)

// 分类选项（从 API 获取）
const categoryOptions = computed(() => {
  return [
    { label: '全部分类', value: 'all' },
    ...categories.value.map((cat) => ({
      label: cat.name,
      value: cat.code,
    })),
  ]
})

// 筛选后的工具列表
const filteredTools = computed(() => {
  let result = tools.value

  // 分类筛选
  if (selectedCategory.value && selectedCategory.value !== 'all') {
    result = result.filter((tool) => tool.category === selectedCategory.value)
  }

  // 搜索筛选
  if (searchKeyword.value) {
    const keyword = searchKeyword.value.toLowerCase()
    result = result.filter(
      (tool) =>
        tool.name.toLowerCase().includes(keyword) ||
        tool.description.toLowerCase().includes(keyword) ||
        tool.tags?.some((tag) => tag.toLowerCase().includes(keyword))
    )
  }

  return result
})

// 分页后的工具列表
const paginatedTools = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  const end = start + pageSize.value
  return filteredTools.value.slice(start, end)
})

// 加载分类列表
const loadCategories = async () => {
  try {
    const res = await toolConfigApi.getToolCategories()
    categories.value = res.data
  } catch (error: any) {
    console.error('加载分类列表失败:', error)
  }
}

// 加载工具列表
const loadTools = async () => {
  try {
    loading.value = true
    tools.value = await toolApi.getToolList()
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

// 筛选变化
const handleFilterChange = () => {
  currentPage.value = 1
}

// 页码变化
const handlePageChange = (page: number) => {
  currentPage.value = page
}

onMounted(async () => {
  await loadCategories()
  loadTools()
})
</script>
