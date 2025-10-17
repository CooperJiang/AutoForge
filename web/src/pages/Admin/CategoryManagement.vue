<template>
  <div class="space-y-4">
    <div class="flex items-center justify-between mb-6">
      <div class="flex gap-2 items-center">
        <BaseInput
          v-model="searchKeyword"
          placeholder="搜索分类名称..."
          style="width: 260px"
          @keyup.enter="loadCategories"
        />
        <BaseSelect
          v-model="filterStatus"
          :options="statusFilterOptions"
          @update:modelValue="loadCategories"
          style="width: 260px"
        />
        <BaseButton @click="loadCategories" variant="primary">
          搜索
        </BaseButton>
      </div>
      <BaseButton @click="showCreateDialog = true">
        <Plus class="w-4 h-4 mr-1.5" />
        新建分类
      </BaseButton>
    </div>

    <div v-if="loading" class="text-center py-8 text-text-secondary">加载中...</div>

    <div v-else-if="categories.length === 0" class="text-center py-8 text-text-secondary">
      暂无分类，点击上方按钮创建第一个分类
    </div>

    <div v-else class="space-y-2">
      <div
        v-for="category in categories"
        :key="category.id"
        class="flex items-center justify-between p-4 bg-bg-primary border border-border-primary rounded-lg hover:bg-bg-hover transition-colors"
      >
        <div class="flex-1">
          <div class="flex items-center space-x-2">
            <h4 class="font-medium text-text-primary">{{ category.name }}</h4>
            <span
              v-if="!category.is_active"
              class="text-xs px-2 py-0.5 bg-gray-500/10 text-gray-600 dark:text-gray-400 rounded"
            >
              已禁用
            </span>
          </div>
          <p class="text-sm text-text-secondary mt-1">{{ category.description || '暂无描述' }}</p>
          <div class="flex items-center space-x-3 mt-2 text-xs text-text-tertiary">
            <span>排序: {{ category.sort_order }}</span>
            <span>创建时间: {{ formatDate(category.created_at) }}</span>
          </div>
        </div>
        <div class="flex items-center space-x-2">
          <button
            @click="editCategory(category)"
            class="p-1.5 bg-blue-500/10 hover:bg-blue-500/20 text-blue-600 dark:text-blue-400 rounded-lg transition-colors border border-blue-500/20"
            title="编辑"
          >
            <Edit2 class="w-4 h-4" />
          </button>
          <button
            @click="deleteCategoryConfirm(category)"
            class="p-1.5 bg-red-500/10 hover:bg-red-500/20 text-red-600 dark:text-red-400 rounded-lg transition-colors border border-red-500/20"
            title="删除"
          >
            <Trash2 class="w-4 h-4" />
          </button>
        </div>
      </div>
    </div>

    <Pagination
      :current="currentPage"
      :page-size="pageSize"
      :total="total"
      :bordered="false"
      @change="handlePageChange"
    />

    <Dialog
      v-model="showCreateDialog"
      :title="editingCategory ? '编辑分类' : '创建分类'"
      @confirm="saveCategory"
      @cancel="cancelEdit"
      confirm-text="保存"
      cancel-text="取消"
    >
      <div class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-text-primary mb-1">分类名称</label>
          <BaseInput v-model="form.name" placeholder="请输入分类名称" />
        </div>
        <div>
          <label class="block text-sm font-medium text-text-primary mb-1">分类描述</label>
          <BaseInput v-model="form.description" placeholder="请输入分类描述（可选）" />
        </div>
        <div>
          <label class="block text-sm font-medium text-text-primary mb-1">
            排序值
            <span class="text-xs text-text-tertiary ml-2">数值越小越靠前，默认 100</span>
          </label>
          <BaseInput
            v-model.number="form.sort_order"
            type="number"
            placeholder="100"
          />
        </div>
        <div v-if="editingCategory" class="flex items-center space-x-2">
          <input
            type="checkbox"
            id="is_active"
            v-model="form.is_active"
            class="w-4 h-4 rounded border-border-primary text-green-600 focus:ring-green-500"
          />
          <label for="is_active" class="text-sm font-medium text-text-primary">启用此分类</label>
        </div>
      </div>
    </Dialog>

    <Dialog
      v-model="showDeleteDialog"
      title="确认删除"
      :message="`确定要删除分类 &quot;${categoryToDelete?.name}&quot; 吗？此操作不可恢复！`"
      confirm-text="删除"
      cancel-text="取消"
      confirm-variant="danger"
      @confirm="deleteCategory"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { Plus, Edit2, Trash2 } from 'lucide-vue-next'
import { templateApi } from '@/api/template'
import type { TemplateCategory, CreateCategoryDto, UpdateCategoryDto } from '@/api/template'
import { message } from '@/utils/message'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import Pagination from '@/components/Pagination'
import Dialog from '@/components/Dialog'

const loading = ref(false)
const categories = ref<TemplateCategory[]>([])
const allCategories = ref<TemplateCategory[]>([])
const showCreateDialog = ref(false)
const showDeleteDialog = ref(false)
const editingCategory = ref<TemplateCategory | null>(null)
const categoryToDelete = ref<TemplateCategory | null>(null)
const searchKeyword = ref('')
const filterStatus = ref('')
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)

const statusFilterOptions = [
  { label: '全部状态', value: '' },
  { label: '已启用', value: 'active' },
  { label: '已禁用', value: 'inactive' },
]

const totalPages = computed(() => Math.ceil(total.value / pageSize.value))

const form = ref({
  name: '',
  description: '',
  sort_order: 100,
  is_active: true,
})

const formatDate = (timestamp: string) => {
  const date = new Date(parseInt(timestamp) * 1000)
  return date.toLocaleString('zh-CN')
}

const loadCategories = async () => {
  loading.value = true
  try {
    const res = await templateApi.listCategories({ page_size: 100 })
    allCategories.value = res.items

    // 应用筛选
    let filtered = allCategories.value

    // 搜索关键词筛选
    if (searchKeyword.value) {
      const keyword = searchKeyword.value.toLowerCase()
      filtered = filtered.filter(
        (cat) =>
          cat.name.toLowerCase().includes(keyword) ||
          cat.description.toLowerCase().includes(keyword)
      )
    }

    // 状态筛选
    if (filterStatus.value === 'active') {
      filtered = filtered.filter((cat) => cat.is_active)
    } else if (filterStatus.value === 'inactive') {
      filtered = filtered.filter((cat) => !cat.is_active)
    }

    total.value = filtered.length

    // 分页
    const start = (currentPage.value - 1) * pageSize.value
    const end = start + pageSize.value
    categories.value = filtered.slice(start, end)
  } catch (error: any) {
    message.error(error.response?.data?.message || '加载分类列表失败')
  } finally {
    loading.value = false
  }
}

// 页码变化
const handlePageChange = (page: number) => {
  currentPage.value = page
  loadCategories()
}

const editCategory = (category: TemplateCategory) => {
  editingCategory.value = category
  form.value = {
    name: category.name,
    description: category.description,
    sort_order: category.sort_order,
    is_active: category.is_active,
  }
  showCreateDialog.value = true
}

const cancelEdit = () => {
  editingCategory.value = null
  form.value = {
    name: '',
    description: '',
    sort_order: 100,
    is_active: true,
  }
}

const saveCategory = async () => {
  if (!form.value.name.trim()) {
    message.error('请输入分类名称')
    return
  }

  try {
    if (editingCategory.value) {
      const updateData: UpdateCategoryDto = {
        name: form.value.name,
        description: form.value.description,
        sort_order: form.value.sort_order,
        is_active: form.value.is_active,
      }
      await templateApi.updateCategory(editingCategory.value.id, updateData)
      message.success('分类更新成功')
    } else {
      const createData: CreateCategoryDto = {
        name: form.value.name,
        description: form.value.description,
        sort_order: form.value.sort_order || 100,
      }
      await templateApi.createCategory(createData)
      message.success('分类创建成功')
    }
    showCreateDialog.value = false
    cancelEdit()
    loadCategories()
  } catch (error: any) {
    message.error(error.response?.data?.message || '保存失败')
  }
}

const deleteCategoryConfirm = (category: TemplateCategory) => {
  categoryToDelete.value = category
  showDeleteDialog.value = true
}

const deleteCategory = async () => {
  if (!categoryToDelete.value) return

  try {
    await templateApi.deleteCategory(categoryToDelete.value.id)
    message.success('分类已删除')
    loadCategories()
  } catch (error: any) {
    message.error(error.response?.data?.message || '删除失败')
  }
}

onMounted(() => {
  loadCategories()
})
</script>
