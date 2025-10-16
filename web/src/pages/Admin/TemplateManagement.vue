<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between mb-6">
      <div class="flex gap-2 items-center">
        <BaseSelect
          v-model="filterCategory"
          :options="categoryFilterOptions"
          @update:modelValue="loadTemplates"
          style="width: 260px"
        />
        <BaseSelect
          v-model="filterFeatured"
          :options="featuredFilterOptions"
          @update:modelValue="loadTemplates"
          style="width: 260px"
        />
        <BaseInput
          v-model="searchKeyword"
          placeholder="搜索模板..."
          style="width: 260px"
          @keyup.enter="loadTemplates"
        />
        <BaseButton @click="loadTemplates" variant="primary">
          搜索
        </BaseButton>
      </div>
      <BaseButton @click="createTemplate">
        <Plus class="w-4 h-4 mr-1.5" />
        新建模板
      </BaseButton>
    </div>

    <div v-if="loading" class="text-center py-8 text-text-secondary">加载中...</div>

    <div v-else-if="templates.length === 0" class="text-center py-8 text-text-secondary">
      暂无模板
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4">
      <div
        v-for="template in templates"
        :key="template.id"
        class="bg-white dark:bg-gray-800 border-2 border-gray-200 dark:border-gray-700 rounded-xl overflow-hidden hover:border-primary hover:shadow-xl transition-all shadow-md"
      >
        <div
          v-if="template.cover_image"
          class="h-28 bg-cover bg-center"
          :style="`background-image: url(${template.cover_image});`"
        ></div>
        <div v-else class="h-28 bg-gradient-to-br from-green-400 to-primary"></div>
        <div class="p-4">
          <div class="flex items-start justify-between mb-2">
            <h4 class="text-sm font-semibold text-text-primary flex-1 line-clamp-1">{{ template.name }}</h4>
            <div class="flex items-center gap-1 ml-2">
              <span
                v-if="template.is_official"
                class="text-xs px-1.5 py-0.5 bg-blue-500/10 text-blue-600 dark:text-blue-400 rounded"
              >
                官方
              </span>
              <span
                v-if="template.is_featured"
                class="text-xs px-1.5 py-0.5 bg-yellow-500/10 text-yellow-600 dark:text-yellow-400 rounded"
              >
                精选
              </span>
            </div>
          </div>
          <p class="text-xs text-text-secondary mb-3 line-clamp-2 min-h-[2.5rem]">
            {{ template.description }}
          </p>
          <div class="flex items-center justify-between text-xs mb-3 pb-3 border-b border-gray-200 dark:border-gray-700">
            <span class="inline-block px-2 py-1 rounded-md text-xs font-medium bg-primary/10 text-primary border border-primary/20">{{ template.category }}</span>
            <div class="flex items-center gap-3 text-text-tertiary">
              <div class="flex items-center gap-1">
                <Download class="w-3.5 h-3.5" />
                <span>{{ template.install_count }}</span>
              </div>
              <div class="flex items-center gap-1">
                <Eye class="w-3.5 h-3.5" />
                <span>{{ template.view_count }}</span>
              </div>
            </div>
          </div>
          <div class="flex items-center gap-2">
            <button
              @click="editTemplate(template)"
              class="flex-1 px-3 py-2 bg-blue-500/10 hover:bg-blue-500/20 text-blue-600 dark:text-blue-400 rounded-lg transition-colors border border-blue-500/20 flex items-center justify-center gap-1.5 text-xs font-medium"
            >
              <Edit2 class="w-3.5 h-3.5" />
              编辑
            </button>
            <button
              @click="deleteTemplateConfirm(template)"
              class="flex-1 px-3 py-2 bg-red-500/10 hover:bg-red-500/20 text-red-600 dark:text-red-400 rounded-lg transition-colors border border-red-500/20 flex items-center justify-center gap-1.5 text-xs font-medium"
            >
              <Trash2 class="w-3.5 h-3.5" />
              删除
            </button>
          </div>
        </div>
      </div>
    </div>

    <div v-if="total > 0" class="flex justify-between items-center mt-6">
      <div class="text-sm text-text-secondary">共 {{ total }} 条记录</div>
      <div class="flex gap-2">
        <button
          @click="prevPage"
          :disabled="currentPage === 1"
          class="px-4 py-2 bg-bg-tertiary text-text-secondary text-sm font-medium rounded-lg hover:bg-bg-tertiary disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          上一页
        </button>
        <span class="px-4 py-2 text-sm text-text-secondary">
          {{ currentPage }} / {{ totalPages }}
        </span>
        <button
          @click="nextPage"
          :disabled="currentPage >= totalPages"
          class="px-4 py-2 bg-bg-tertiary text-text-secondary text-sm font-medium rounded-lg hover:bg-bg-tertiary disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          下一页
        </button>
      </div>
    </div>

    <Dialog
      v-model="showEditDialog"
      :title="`编辑模板 - ${editingTemplate?.name || ''}`"
      @confirm="saveTemplate"
      @cancel="cancelEdit"
      confirm-text="保存"
      cancel-text="取消"
      max-width="max-w-2xl"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-text-primary mb-1">模板名称</label>
          <BaseInput v-model="form.name" placeholder="请输入模板名称" />
        </div>
        <div>
          <label class="block text-sm font-medium text-text-primary mb-1">模板描述</label>
          <textarea
            v-model="form.description"
            placeholder="请输入模板描述"
            rows="3"
            class="w-full px-3 py-2 bg-bg-primary border border-border-primary rounded-lg text-text-primary placeholder-text-tertiary focus:outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>
        <div>
          <label class="block text-sm font-medium text-text-primary mb-1">分类</label>
          <BaseSelect v-model="form.category" :options="categoryOptions" />
        </div>
        <div>
          <label class="block text-sm font-medium text-text-primary mb-1">封面图片 URL</label>
          <BaseInput v-model="form.cover_image" placeholder="请输入封面图片 URL（可选）" />
          <p class="text-xs text-text-tertiary mt-1">建议尺寸: 512x512 或 1:1 比例</p>
        </div>
        <div>
          <label class="block text-sm font-medium text-text-primary mb-1">使用指南</label>
          <textarea
            v-model="form.usage_guide"
            placeholder="请输入使用指南（可选）"
            rows="4"
            class="w-full px-3 py-2 bg-bg-primary border border-border-primary rounded-lg text-text-primary placeholder-text-tertiary focus:outline-none focus:ring-2 focus:ring-green-500"
          />
        </div>
        <div class="flex items-center space-x-2">
          <input
            type="checkbox"
            id="is_featured"
            v-model="form.is_featured"
            class="w-4 h-4 rounded border-border-primary text-green-600 focus:ring-green-500"
          />
          <label for="is_featured" class="text-sm font-medium text-text-primary">
            设为精选
          </label>
        </div>
        <div>
          <label class="block text-sm font-medium text-text-primary mb-2">状态</label>
          <RadioGroup v-model="form.status" :options="statusOptions" />
        </div>
      </div>
    </Dialog>

    <Dialog
      v-model="showDeleteDialog"
      title="确认删除"
      :message="`确定要删除模板 &quot;${templateToDelete?.name}&quot; 吗？此操作不可恢复！`"
      confirm-text="删除"
      cancel-text="取消"
      confirm-variant="danger"
      @confirm="deleteTemplate"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { Search, Edit2, Trash2, Plus, Download, Eye } from 'lucide-vue-next'
import { templateApi } from '@/api/template'
import type {
  TemplateBasicInfo,
  TemplateCategory,
  UpdateTemplateDto,
} from '@/api/template'
import { message } from '@/utils/message'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import RadioGroup from '@/components/RadioGroup/index.vue'
import Dialog from '@/components/Dialog'

const router = useRouter()
const loading = ref(false)
const templates = ref<TemplateBasicInfo[]>([])
const categories = ref<TemplateCategory[]>([])
const showEditDialog = ref(false)
const showDeleteDialog = ref(false)
const editingTemplate = ref<TemplateBasicInfo | null>(null)
const templateToDelete = ref<TemplateBasicInfo | null>(null)

const searchKeyword = ref('')
const filterCategory = ref('')
const filterFeatured = ref('')
const currentPage = ref(1)
const totalPages = ref(1)
const total = ref(0)

const form = ref({
  name: '',
  description: '',
  category: '',
  cover_image: '',
  usage_guide: '',
  is_featured: false,
  status: 'active',
})

// 分类筛选选项
const categoryFilterOptions = computed(() => {
  const options = [{ label: '全部分类', value: '' }]
  categories.value.forEach((cat) => {
    options.push({
      label: cat.name,
      value: cat.name,
    })
  })
  return options
})

// 精选筛选选项
const featuredFilterOptions = [
  { label: '全部', value: '' },
  { label: '精选', value: 'true' },
  { label: '非精选', value: 'false' },
]

// 分类选项（用于编辑对话框）
const categoryOptions = computed(() => {
  const options = [{ label: '请选择分类', value: '' }]
  categories.value.forEach((cat) => {
    options.push({
      label: cat.name,
      value: cat.name,
    })
  })
  return options
})

// 状态选项
const statusOptions = [
  { label: '启用', value: 'active' },
  { label: '禁用', value: 'inactive' },
]

const loadCategories = async () => {
  try {
    const res = await templateApi.listCategories({ page_size: 100, is_active: true })
    categories.value = res.items
  } catch (error: any) {
    console.error('加载分类失败:', error)
  }
}

const loadTemplates = async () => {
  loading.value = true
  try {
    const params: any = {
      page: currentPage.value,
      page_size: 12,
    }
    if (searchKeyword.value) params.search = searchKeyword.value
    if (filterCategory.value) params.category = filterCategory.value
    if (filterFeatured.value) params.is_featured = filterFeatured.value === 'true'

    const res = await templateApi.list(params)
    templates.value = res.items
    total.value = res.total
    totalPages.value = res.total_pages
  } catch (error: any) {
    message.error(error.response?.data?.message || '加载模板列表失败')
  } finally {
    loading.value = false
  }
}

const createTemplate = () => {
  router.push('/workflows/new')
}

const editTemplate = (template: TemplateBasicInfo) => {
  editingTemplate.value = template
  form.value = {
    name: template.name,
    description: template.description,
    category: template.category,
    cover_image: template.cover_image || '',
    usage_guide: '',
    is_featured: template.is_featured,
    status: 'active',
  }
  showEditDialog.value = true
}

const cancelEdit = () => {
  editingTemplate.value = null
  form.value = {
    name: '',
    description: '',
    category: '',
    cover_image: '',
    usage_guide: '',
    is_featured: false,
    status: 'active',
  }
}

const saveTemplate = async () => {
  if (!editingTemplate.value) return
  if (!form.value.name.trim()) {
    message.error('请输入模板名称')
    return
  }

  try {
    const updateData: UpdateTemplateDto = {
      name: form.value.name,
      description: form.value.description,
      category: form.value.category,
      cover_image: form.value.cover_image,
      usage_guide: form.value.usage_guide,
      is_featured: form.value.is_featured,
      status: form.value.status,
    }
    await templateApi.update(editingTemplate.value.id, updateData)
    message.success('模板更新成功')
    showEditDialog.value = false
    cancelEdit()
    loadTemplates()
  } catch (error: any) {
    message.error(error.response?.data?.message || '保存失败')
  }
}

const deleteTemplateConfirm = (template: TemplateBasicInfo) => {
  templateToDelete.value = template
  showDeleteDialog.value = true
}

const deleteTemplate = async () => {
  if (!templateToDelete.value) return

  try {
    await templateApi.delete(templateToDelete.value.id)
    message.success('模板已删除')
    loadTemplates()
  } catch (error: any) {
    message.error(error.response?.data?.message || '删除失败')
  }
}

const prevPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
    loadTemplates()
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
    loadTemplates()
  }
}

onMounted(() => {
  loadCategories()
  loadTemplates()
})
</script>
