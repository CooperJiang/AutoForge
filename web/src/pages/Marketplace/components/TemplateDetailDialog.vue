<template>
  <Dialog
    :modelValue="visible"
    title="模板详情"
    max-width="max-w-3xl"
    hide-footer
    @update:modelValue="handleClose"
  >
    <div v-if="loading" class="flex justify-center items-center py-20">
      <div class="text-text-tertiary">加载中...</div>
    </div>

    <div v-else-if="template" class="space-y-6">
      <!-- Header -->
      <div class="flex items-start gap-4">
        <div
          class="w-16 h-16 rounded-lg bg-gradient-to-br from-primary/20 to-accent/20 flex items-center justify-center flex-shrink-0"
        >
          <span class="text-4xl">{{ template.icon || '📦' }}</span>
        </div>
        <div class="flex-1">
          <div class="flex items-center gap-2 mb-2">
            <h2 class="text-xl font-bold text-text-primary">{{ template.name }}</h2>
            <span
              v-if="template.is_official"
              class="inline-flex items-center px-2 py-1 rounded text-xs bg-primary/10 text-primary"
            >
              <Shield class="w-3 h-3 mr-1" />
              官方
            </span>
            <span
              v-if="template.is_featured"
              class="inline-flex items-center px-2 py-1 rounded text-xs bg-accent/10 text-accent"
            >
              <Star class="w-3 h-3 mr-1 fill-current" />
              精选
            </span>
          </div>
          <p class="text-sm text-text-secondary">{{ template.description }}</p>
        </div>
      </div>

      <!-- Stats -->
      <div class="flex items-center gap-6 text-sm">
        <div class="flex items-center gap-2 text-text-tertiary">
          <Download class="w-4 h-4" />
          <span>{{ template.install_count }} 次安装</span>
        </div>
        <div class="flex items-center gap-2 text-text-tertiary">
          <Eye class="w-4 h-4" />
          <span>{{ template.view_count }} 次浏览</span>
        </div>
        <div class="flex items-center gap-2 text-text-tertiary">
          <Calendar class="w-4 h-4" />
          <span>创建于 {{ formatDate(template.created_at) }}</span>
        </div>
      </div>

      <!-- Category -->
      <div>
        <h3 class="text-sm font-semibold text-text-primary mb-2">分类</h3>
        <span
          class="inline-block px-3 py-1 rounded bg-surface-tertiary text-text-secondary text-sm"
        >
          {{ getCategoryName(template.category) }}
        </span>
      </div>

      <!-- Required Tools -->
      <div v-if="template.required_tools && template.required_tools.length > 0">
        <h3 class="text-sm font-semibold text-text-primary mb-2">所需工具</h3>
        <div class="flex flex-wrap gap-2">
          <span
            v-for="tool in template.required_tools"
            :key="tool"
            class="inline-flex items-center px-3 py-1 rounded bg-surface-tertiary text-text-secondary text-sm"
          >
            <Wrench class="w-3 h-3 mr-1" />
            {{ tool }}
          </span>
        </div>
      </div>

      <!-- Usage Guide -->
      <div v-if="template.usage_guide">
        <h3 class="text-sm font-semibold text-text-primary mb-2">使用指南</h3>
        <div
          class="bg-surface-tertiary rounded p-4 text-sm text-text-secondary whitespace-pre-wrap"
        >
          {{ template.usage_guide }}
        </div>
      </div>

      <!-- Workflow Preview -->
      <div>
        <h3 class="text-sm font-semibold text-text-primary mb-2">工作流结构</h3>
        <div class="bg-surface-tertiary rounded p-4">
          <div class="text-sm text-text-secondary space-y-2">
            <div class="flex items-center gap-2">
              <span class="font-medium">节点数:</span>
              <span>{{ template.template_data?.nodes?.length || 0 }}</span>
            </div>
            <div class="flex items-center gap-2">
              <span class="font-medium">连接数:</span>
              <span>{{ template.template_data?.edges?.length || 0 }}</span>
            </div>
            <div v-if="template.template_data?.env_vars?.length" class="flex items-center gap-2">
              <span class="font-medium">环境变量:</span>
              <span>{{ template.template_data.env_vars.length }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <template v-if="template" #footer>
      <div class="flex justify-end gap-3">
        <BaseButton variant="outline" @click="handleClose">取消</BaseButton>
        <BaseButton @click="handleInstall">
          <Download class="w-4 h-4 mr-1" />
          安装模板
        </BaseButton>
      </div>
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Shield, Star, Download, Eye, Calendar, Wrench } from 'lucide-vue-next'
import Dialog from '@/components/Dialog'
import BaseButton from '@/components/BaseButton'
import { templateApi } from '@/api/template'
import type { TemplateDetail } from '@/api/template'
import { message } from '@/utils/message'

const props = defineProps<{
  visible: boolean
  templateId: string
}>()

const emit = defineEmits<{
  close: []
  install: [templateId: string]
}>()

const loading = ref(false)
const template = ref<TemplateDetail | null>(null)

const loadTemplate = async () => {
  if (!props.templateId) return

  loading.value = true
  try {
    template.value = await templateApi.getById(props.templateId)
  } catch (error) {
    console.error('Failed to load template:', error)
    message.error('加载模板详情失败')
  } finally {
    loading.value = false
  }
}

const handleInstall = () => {
  emit('install', props.templateId)
}

const handleClose = () => {
  emit('close')
}

const getCategoryName = (category: string): string => {
  const categoryMap: Record<string, string> = {
    automation: '自动化',
    notification: '通知',
    data: '数据处理',
    integration: '集成',
    other: '其他',
  }
  return categoryMap[category] || category
}

const formatDate = (dateStr: string): string => {
  return new Date(dateStr).toLocaleDateString('zh-CN')
}

watch(
  () => props.visible,
  (newVal) => {
    if (newVal && props.templateId) {
      loadTemplate()
    }
  }
)
</script>
