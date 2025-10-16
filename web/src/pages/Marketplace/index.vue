<template>
  <div>
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h1 class="text-2xl font-bold text-text-primary mb-1">模板市场</h1>
        <p class="text-sm text-text-secondary">浏览和安装官方工作流模板</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex gap-4 mb-6">
      <!-- Category Filter -->
      <BaseSelect
        v-model="selectedCategory"
        :options="categoryOptions"
        style="width: 260px"
        placeholder="选择分类"
        @update:modelValue="handleFilterChange"
      />

      <!-- Featured Filter -->
      <BaseButton :variant="showFeatured ? 'primary' : 'outline'" @click="toggleFeatured">
        <Star :class="['w-4 h-4 mr-1', showFeatured && 'fill-current']" />
        精选模板
      </BaseButton>

      <!-- Search -->
      <BaseInput
        v-model="searchKeyword"
        placeholder="搜索模板..."
        style="width: 260px"
        @keyup.enter="handleFilterChange"
      >
        <template #prefix>
          <Search class="w-4 h-4" />
        </template>
      </BaseInput>
    </div>

    <!-- Template Grid -->
    <div
      v-if="!loading && templates.length > 0"
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4"
    >
      <TemplateCard
        v-for="template in templates"
        :key="template.id"
        :template="template"
        @view="handleViewDetail"
        @install="handleInstall"
      />
    </div>

    <!-- Empty State -->
    <div v-else-if="!loading && templates.length === 0" class="text-center py-20">
      <div class="text-text-placeholder mb-4">
        <Package class="w-16 h-16 mx-auto mb-4" />
        <p class="text-lg">暂无模板</p>
        <p class="text-sm">暂时没有符合条件的模板</p>
      </div>
    </div>

    <!-- Loading -->
    <div v-else class="flex justify-center items-center py-20">
      <div class="text-text-tertiary">加载中...</div>
    </div>

    <!-- Pagination -->
    <div v-if="!loading && templates.length > 0" class="mt-8">
      <Pagination
        :current="currentPage"
        :total="totalItems"
        :page-size="pageSize"
        @change="handlePageChange"
      />
    </div>

    <!-- Template Detail Dialog -->
    <TemplateDetailDialog
      :visible="detailDialogVisible"
      :template-id="selectedTemplateId"
      @close="detailDialogVisible = false"
      @install="handleInstallFromDetail"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Star, Search, Package } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import BaseSelect from '@/components/BaseSelect'
import BaseInput from '@/components/BaseInput'
import Pagination from '@/components/Pagination'
import TemplateCard from './components/TemplateCard.vue'
import TemplateDetailDialog from './components/TemplateDetailDialog.vue'
import { templateApi } from '@/api/template'
import type { TemplateBasicInfo } from '@/api/template'
import { message } from '@/utils/message'

const router = useRouter()
const loading = ref(false)
const templates = ref<TemplateBasicInfo[]>([])
const selectedCategory = ref('')
const showFeatured = ref(false)
const searchKeyword = ref('')
const currentPage = ref(1)
const pageSize = ref(20)
const totalItems = ref(0)
const detailDialogVisible = ref(false)
const selectedTemplateId = ref('')
const categories = ref<any[]>([])

// Computed category options
const categoryOptions = computed(() => {
  const options = [{ label: '全部分类', value: '' }]
  categories.value.forEach((cat) => {
    options.push({
      label: cat.name,
      value: cat.name,
    })
  })
  return options
})

// Load categories
const loadCategories = async () => {
  try {
    const data = await templateApi.listCategories({ page_size: 100, is_active: true })
    categories.value = data.items || []
  } catch (error) {
    console.error('Failed to load categories:', error)
    message.error('加载分类列表失败')
  }
}

// Load templates
const loadTemplates = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value,
    }

    if (selectedCategory.value) {
      params.category = selectedCategory.value
    }

    if (showFeatured.value) {
      params.is_featured = true
    }

    if (searchKeyword.value) {
      params.search = searchKeyword.value
    }

    const data = await templateApi.list(params)
    templates.value = data.items || []
    totalItems.value = data.total || 0
  } catch (error) {
    console.error('Failed to load templates:', error)
    message.error('加载模板列表失败')
  } finally {
    loading.value = false
  }
}

// Toggle featured
const toggleFeatured = () => {
  showFeatured.value = !showFeatured.value
  currentPage.value = 1
  loadTemplates()
}

// Filter change
const handleFilterChange = () => {
  currentPage.value = 1
  loadTemplates()
}

// Page change
const handlePageChange = (page: number) => {
  currentPage.value = page
  loadTemplates()
}

// View detail
const handleViewDetail = (template: TemplateBasicInfo) => {
  selectedTemplateId.value = template.id
  detailDialogVisible.value = true
}

// Install template
const handleInstall = async (template: TemplateBasicInfo) => {
  try {
    const result = await templateApi.install({
      template_id: template.id,
    })
    message.success('模板安装成功')
    router.push(`/workflows/${result.workflow_id}/edit`)
  } catch (error: any) {
    console.error('Install template failed:', error)
    message.error(error.response?.data?.message || '安装失败')
  }
}

// Install from detail dialog
const handleInstallFromDetail = async (templateId: string) => {
  try {
    const result = await templateApi.install({
      template_id: templateId,
    })
    detailDialogVisible.value = false
    message.success('模板安装成功')
    router.push(`/workflows/${result.workflow_id}/edit`)
  } catch (error: any) {
    console.error('Install template failed:', error)
    message.error(error.response?.data?.message || '安装失败')
  }
}

onMounted(() => {
  loadCategories()
  loadTemplates()
})
</script>
