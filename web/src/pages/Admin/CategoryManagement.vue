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
        <BaseButton @click="loadCategories" variant="primary"> 搜索 </BaseButton>
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

    <!-- 分类卡片网格 -->
    <div
      v-else
      class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 2xl:grid-cols-5 gap-4"
    >
      <div
        v-for="category in categories"
        :key="category.id"
        class="border-2 rounded-xl overflow-hidden hover:shadow-xl transition-all shadow-md relative"
        :style="{
          backgroundColor: 'var(--color-bg-elevated)',
          borderColor: category.is_active
            ? 'var(--color-border-primary)'
            : 'var(--color-border-secondary)',
          opacity: category.is_active ? '1' : '0.6',
        }"
        @mouseenter="
          (e) => ((e.currentTarget as HTMLElement).style.borderColor = 'var(--color-primary)')
        "
        @mouseleave="
          (e) =>
            ((e.currentTarget as HTMLElement).style.borderColor = category.is_active
              ? 'var(--color-border-primary)'
              : 'var(--color-border-secondary)')
        "
      >
        <!-- 分类头部 -->
        <div class="p-4 border-b-2 border-border-primary">
          <div class="flex items-start justify-between">
            <div class="flex-1 min-w-0">
              <div class="flex items-center gap-2 mb-1">
                <h3 class="text-base font-bold text-text-primary truncate">
                  {{ category.name }}
                </h3>
                <span
                  v-if="!category.is_active"
                  class="px-2 py-0.5 text-xs font-semibold rounded-full bg-bg-tertiary text-text-disabled flex-shrink-0"
                >
                  已禁用
                </span>
              </div>
              <p class="text-xs text-text-tertiary line-clamp-2 min-h-[2.5rem]">
                {{ category.description || '暂无描述' }}
              </p>
            </div>
          </div>
        </div>

        <!-- 分类信息 -->
        <div class="p-4 space-y-2">
          <div class="flex items-center justify-between text-xs">
            <span class="text-text-tertiary">排序值</span>
            <span class="font-semibold text-text-primary">{{ category.sort_order }}</span>
          </div>
          <div class="flex items-center justify-between text-xs">
            <span class="text-text-tertiary">创建时间</span>
            <span class="font-medium text-text-secondary">{{
              formatDate(category.created_at)
            }}</span>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div
          class="p-3 border-t-2 border-border-primary flex items-center justify-between bg-bg-secondary"
        >
          <div class="flex items-center gap-1.5">
            <button
              @click="editCategory(category)"
              class="p-1.5 rounded-lg transition-all hover:bg-info-light"
              :style="{ color: 'var(--color-info)' }"
              title="编辑"
            >
              <Edit2 class="w-4 h-4" />
            </button>
            <button
              @click="deleteCategoryConfirm(category)"
              class="p-1.5 rounded-lg transition-all hover:bg-error-light"
              :style="{ color: 'var(--color-error)' }"
              title="删除"
            >
              <Trash2 class="w-4 h-4" />
            </button>
          </div>
          <div
            class="text-xs font-medium"
            :style="{
              color: category.is_active ? 'var(--color-success)' : 'var(--color-text-disabled)',
            }"
          >
            {{ category.is_active ? '已启用' : '已禁用' }}
          </div>
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
          <BaseInput v-model.number="form.sort_order" type="number" placeholder="100" />
        </div>
        <div v-if="editingCategory" class="flex items-center space-x-2">
          <BaseCheckbox v-model="form.is_active" label="启用此分类" />
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
import { ref, onMounted } from 'vue'
import { Plus, Edit2, Trash2 } from 'lucide-vue-next'
import { templateApi } from '@/api/template'
import type { TemplateCategory, CreateCategoryDto, UpdateCategoryDto } from '@/api/template'
import { message } from '@/utils/message'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import Pagination from '@/components/Pagination'
import Dialog from '@/components/Dialog'
import BaseCheckbox from '@/components/BaseCheckbox/index.vue'

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
