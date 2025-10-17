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
          placeholder="搜索模板..."
          style="width: 260px"
          @keyup.enter="handleFilterChange"
        >
          <template #prefix>
            <Search class="w-4 h-4" />
          </template>
        </BaseInput>

        <!-- Featured Filter - 放到最后 -->
        <BaseButton :variant="showFeatured ? 'primary' : 'outline'" @click="toggleFeatured">
          <Star :class="['w-4 h-4 mr-1', showFeatured && 'fill-current']" />
          精选工作流
        </BaseButton>
      </div>
    </div>

    <!-- Content Area with Background -->
    <div class="flex-1 flex flex-col bg-bg-elevated rounded-xl border border-border-primary overflow-hidden shadow-sm">
      <!-- Template Grid - Scrollable -->
      <div class="flex-1 overflow-y-auto p-4">
        <div
          v-if="!loading && templates.length > 0"
          class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-3"
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
        <div v-else-if="!loading && templates.length === 0" class="flex items-center justify-center h-full">
          <div class="text-text-placeholder text-center">
            <Package class="w-16 h-16 mx-auto mb-4" />
            <p class="text-lg">暂无模板</p>
            <p class="text-sm">暂时没有符合条件的模板</p>
          </div>
        </div>

        <!-- Loading -->
        <div v-else class="flex items-center justify-center h-full">
          <div class="text-text-tertiary">加载中...</div>
        </div>
      </div>

      <!-- Pagination - Fixed at Bottom -->
      <div
        v-if="!loading && templates.length > 0"
        class="border-t border-border-primary px-4 py-3 flex-shrink-0"
      >
        <Pagination
          :current="currentPage"
          :total="totalItems"
          :page-size="pageSize"
          @change="handlePageChange"
        />
      </div>
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
