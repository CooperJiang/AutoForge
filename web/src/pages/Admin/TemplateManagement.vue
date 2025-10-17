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
        class="border-2 rounded-xl overflow-hidden hover:shadow-xl transition-all shadow-md"
        :style="{
          backgroundColor: 'var(--color-bg-elevated)',
          borderColor: 'var(--color-border-primary)',
        }"
        @mouseenter="(e) => (e.currentTarget as HTMLElement).style.borderColor = 'var(--color-primary)'"
        @mouseleave="(e) => (e.currentTarget as HTMLElement).style.borderColor = 'var(--color-border-primary)'"
      >
        <div
          v-if="template.cover_image"
          class="h-40 bg-cover bg-center"
          :style="`background-image: url(${template.cover_image});`"
        ></div>
        <div
          v-else
          class="h-40 flex items-center justify-center p-4"
          :style="{
            background: 'var(--gradient-primary)',
          }"
        >
          <h3
            class="text-xl font-bold text-center line-clamp-3 max-w-[200px] drop-shadow-lg"
            :style="{ color: 'var(--color-primary-text)' }"
          >
            {{ template.name }}
          </h3>
        </div>
        <div class="p-3">
          <div class="flex items-start justify-between mb-1.5 gap-2">
            <h4
              class="text-sm font-semibold line-clamp-1 flex-1 min-w-0"
              :style="{ color: 'var(--color-text-primary)' }"
            >
              {{ template.name }}
            </h4>
            <div class="flex items-center gap-2 flex-shrink-0">
              <div class="flex items-center gap-2 text-xs">
                <div class="flex items-center gap-1" :style="{ color: 'var(--color-success)' }">
                  <Download class="w-3.5 h-3.5" />
                  <span>{{ template.install_count }}</span>
                </div>
                <div class="flex items-center gap-1" :style="{ color: 'var(--color-info)' }">
                  <Eye class="w-3.5 h-3.5" />
                  <span>{{ template.view_count }}</span>
                </div>
              </div>
              <div class="flex items-center gap-1">
                <span
                  v-if="template.status === 'draft'"
                  class="text-xs px-1.5 py-0.5 rounded"
                  :style="{
                    backgroundColor: 'var(--color-bg-tertiary)',
                    color: 'var(--color-text-secondary)',
                  }"
                >
                  草稿
                </span>
                <span
                  v-if="template.status === 'archived'"
                  class="text-xs px-1.5 py-0.5 rounded"
                  :style="{
                    backgroundColor: 'var(--color-error-light)',
                    color: 'var(--color-error-text)',
                  }"
                >
                  已下架
                </span>
              </div>
            </div>
          </div>
          <p class="text-xs mb-2 line-clamp-2 min-h-[2rem]" :style="{ color: 'var(--color-text-secondary)' }">
            {{ template.description }}
          </p>
          <div class="flex items-center justify-between text-xs">
            <div class="flex items-center gap-1.5">
              <span
                v-if="template.is_official"
                class="text-xs px-1.5 py-0.5 rounded"
                :style="{
                  backgroundColor: 'var(--color-info-light)',
                  color: 'var(--color-info-text)',
                }"
              >
                官方
              </span>
              <span
                v-if="template.is_featured"
                class="text-xs px-1.5 py-0.5 rounded"
                :style="{
                  backgroundColor: 'var(--color-warning-light)',
                  color: 'var(--color-warning-text)',
                }"
              >
                精选
              </span>
              <span
                class="inline-block px-2 py-0.5 rounded-md text-xs font-medium border"
                :style="{
                  backgroundColor: 'var(--color-primary-light)',
                  color: 'var(--color-primary)',
                  borderColor: 'var(--color-primary)',
                  opacity: '0.8',
                }"
              >
                {{ template.category }}
              </span>
            </div>
            <div class="flex items-center gap-2">
              <button
                @click="editTemplate(template)"
                class="p-1 rounded transition-colors"
                :style="{ color: 'var(--color-info)' }"
                @mouseenter="(e) => (e.currentTarget as HTMLElement).style.backgroundColor = 'var(--color-info-light)'"
                @mouseleave="(e) => (e.currentTarget as HTMLElement).style.backgroundColor = 'transparent'"
                title="编辑"
              >
                <Settings class="w-3.5 h-3.5" />
              </button>
              <button
                @click="deleteTemplateConfirm(template)"
                class="p-1 rounded transition-colors"
                :style="{ color: 'var(--color-error)' }"
                @mouseenter="(e) => (e.currentTarget as HTMLElement).style.backgroundColor = 'var(--color-error-light)'"
                @mouseleave="(e) => (e.currentTarget as HTMLElement).style.backgroundColor = 'transparent'"
                title="删除"
              >
                <Trash2 class="w-3.5 h-3.5" />
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <Pagination
      v-if="total > 0"
      :current="currentPage"
      :page-size="12"
      :total="total"
      :bordered="false"
      @change="handlePageChange"
    />

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
        <div>
          <label class="block text-sm font-medium text-text-primary mb-2">
            案例展示图片
          </label>
          <div class="space-y-2">
            <div
              v-for="(url, index) in form.case_images"
              :key="index"
              class="flex items-center gap-2"
            >
              <BaseInput
                v-model="form.case_images[index]"
                placeholder="请输入案例图片 URL"
                class="flex-1"
              />
              <BaseButton
                variant="outline"
                size="sm"
                @click="removeCaseImage(index)"
              >
                删除
              </BaseButton>
            </div>
            <BaseButton
              variant="outline"
              size="sm"
              @click="addCaseImage"
              class="w-full"
            >
              <Plus class="w-4 h-4 mr-1" />
              添加案例图片
            </BaseButton>
          </div>
          <p class="text-xs text-text-tertiary mt-1">
            添加工作流运行结果的案例图片，帮助用户了解效果
          </p>
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
import { Search, Settings, Trash2, Plus, Download, Eye } from 'lucide-vue-next'
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
import Pagination from '@/components/Pagination'
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
  case_images: [] as string[],
  is_featured: false,
  status: 'published',
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
  { label: '已发布', value: 'published' },
  { label: '草稿', value: 'draft' },
  { label: '已下架', value: 'archived' },
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
      show_all: true, // 管理后台显示所有状态的模板
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

const editTemplate = async (template: TemplateBasicInfo) => {
  try {
    // 获取模板详情以获得 usage_guide 和 case_images
    const detail = await templateApi.getById(template.id)

    editingTemplate.value = template
    form.value = {
      name: template.name,
      description: template.description,
      category: template.category,
      cover_image: template.cover_image || '',
      usage_guide: detail.usage_guide || '',
      case_images: detail.case_images || [],
      is_featured: template.is_featured,
      status: template.status || 'published',
    }
    showEditDialog.value = true
  } catch (error: any) {
    message.error(error.response?.data?.message || '获取模板详情失败')
  }
}

const cancelEdit = () => {
  editingTemplate.value = null
  form.value = {
    name: '',
    description: '',
    category: '',
    cover_image: '',
    usage_guide: '',
    case_images: [],
    is_featured: false,
    status: 'published',
  }
}

const addCaseImage = () => {
  form.value.case_images.push('')
}

const removeCaseImage = (index: number) => {
  form.value.case_images.splice(index, 1)
}

const saveTemplate = async () => {
  if (!editingTemplate.value) return
  if (!form.value.name.trim()) {
    message.error('请输入模板名称')
    return
  }

  try {
    // 过滤掉空的案例图片 URL
    const caseImages = form.value.case_images.filter((url) => url.trim() !== '')

    const updateData: UpdateTemplateDto = {
      name: form.value.name,
      description: form.value.description,
      category: form.value.category,
      cover_image: form.value.cover_image,
      usage_guide: form.value.usage_guide,
      case_images: caseImages,
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

// 页码变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  loadTemplates()
}

onMounted(() => {
  loadCategories()
  loadTemplates()
})
</script>
