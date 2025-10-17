import { onMounted, type Ref } from 'vue'
import { templateApi } from '@/api/template'
import { message } from '@/utils/message'
import type { WorkflowTemplate } from '../types'

/**
 * 模板管理操作逻辑
 */
export function useTemplateActions(options: {
  templates: Ref<WorkflowTemplate[]>
  loading: Ref<boolean>
  filterCategory: Ref<string>
  filterFeatured: Ref<string>
  searchKeyword: Ref<string>
  currentTemplate: Ref<WorkflowTemplate | null>
  showEditDialog: Ref<boolean>
  showDeleteConfirm: Ref<boolean>
}) {
  const {
    templates,
    loading,
    filterCategory,
    filterFeatured,
    searchKeyword,
    currentTemplate,
    showEditDialog,
    showDeleteConfirm,
  } = options

  // 加载模板列表
  const loadTemplates = async () => {
    try {
      loading.value = true
      const params: any = {}

      if (filterCategory.value) {
        params.category = filterCategory.value
      }
      if (filterFeatured.value) {
        params.featured = filterFeatured.value === 'true'
      }
      if (searchKeyword.value) {
        params.keyword = searchKeyword.value
      }

      const data = await templateApi.getList(params)
      templates.value = data.list || []
    } catch (error: any) {
      message.error(error.response?.data?.message || '加载模板失败')
    } finally {
      loading.value = false
    }
  }

  // 创建模板
  const createTemplate = () => {
    currentTemplate.value = null
    showEditDialog.value = true
  }

  // 编辑模板
  const editTemplate = (template: WorkflowTemplate) => {
    currentTemplate.value = template
    showEditDialog.value = true
  }

  // 删除模板
  const deleteTemplate = (template: WorkflowTemplate) => {
    currentTemplate.value = template
    showDeleteConfirm.value = true
  }

  // 确认删除
  const confirmDelete = async () => {
    if (!currentTemplate.value) return

    try {
      await templateApi.delete(currentTemplate.value.id)
      message.success('删除成功')
      showDeleteConfirm.value = false
      await loadTemplates()
    } catch (error: any) {
      message.error(error.response?.data?.message || '删除失败')
    }
  }

  // 切换精选状态
  const toggleFeatured = async (template: WorkflowTemplate) => {
    try {
      await templateApi.update(template.id, {
        featured: !template.featured,
      })
      message.success('更新成功')
      await loadTemplates()
    } catch (error: any) {
      message.error(error.response?.data?.message || '更新失败')
    }
  }

  // 切换发布状态
  const toggleStatus = async (template: WorkflowTemplate) => {
    try {
      const newStatus = template.status === 'published' ? 'draft' : 'published'
      await templateApi.update(template.id, { status: newStatus })
      message.success('更新成功')
      await loadTemplates()
    } catch (error: any) {
      message.error(error.response?.data?.message || '更新失败')
    }
  }

  // 保存模板
  const saveTemplate = async (formData: Partial<WorkflowTemplate>) => {
    try {
      if (currentTemplate.value) {
        await templateApi.update(currentTemplate.value.id, formData)
        message.success('更新成功')
      } else {
        await templateApi.create(formData)
        message.success('创建成功')
      }
      showEditDialog.value = false
      await loadTemplates()
    } catch (error: any) {
      message.error(error.response?.data?.message || '保存失败')
    }
  }

  // 初始化
  onMounted(() => {
    loadTemplates()
  })

  return {
    loadTemplates,
    createTemplate,
    editTemplate,
    deleteTemplate,
    confirmDelete,
    toggleFeatured,
    toggleStatus,
    saveTemplate,
  }
}

