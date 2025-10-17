import { ref, computed } from 'vue'
import type { WorkflowTemplate, SelectOption } from '../types'

/**
 * 模板管理状态
 */
export function useTemplateState() {
  const templates = ref<WorkflowTemplate[]>([])
  const loading = ref(false)
  const filterCategory = ref('')
  const filterFeatured = ref('')
  const searchKeyword = ref('')
  const showEditDialog = ref(false)
  const showDeleteConfirm = ref(false)
  const currentTemplate = ref<WorkflowTemplate | null>(null)

  // 分类选项
  const categoryFilterOptions: SelectOption[] = [
    { label: '全部分类', value: '' },
    { label: '自动化', value: 'automation' },
    { label: '通知', value: 'notification' },
    { label: '数据处理', value: 'data' },
    { label: 'AI/LLM', value: 'ai' },
    { label: '监控', value: 'monitoring' },
    { label: '其他', value: 'other' },
  ]

  // 精选过滤选项
  const featuredFilterOptions: SelectOption[] = [
    { label: '全部', value: '' },
    { label: '精选', value: 'true' },
    { label: '非精选', value: 'false' },
  ]

  return {
    templates,
    loading,
    filterCategory,
    filterFeatured,
    searchKeyword,
    showEditDialog,
    showDeleteConfirm,
    currentTemplate,
    categoryFilterOptions,
    featuredFilterOptions,
  }
}

