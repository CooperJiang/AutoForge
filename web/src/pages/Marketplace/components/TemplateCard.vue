<template>
  <div
    class="bg-bg-elevated rounded-xl border-2 border-border-primary hover:border-primary hover:shadow-xl transition-all cursor-pointer overflow-hidden shadow-sm"
    @click="$emit('view', template)"
  >
    <!-- Cover Image -->
    <div
      v-if="template.cover_image"
      class="h-28 bg-cover bg-center"
      :style="`background-image: url(${template.cover_image});`"
    ></div>
    <div
      v-else
      class="h-28 bg-gradient-to-br from-green-400 via-primary to-green-600 flex items-center justify-center p-2.5"
    >
      <h3 class="text-white text-base font-bold text-center line-clamp-2 max-w-[160px] drop-shadow-lg">
        {{ template.name }}
      </h3>
    </div>

    <!-- Content -->
    <div class="p-3">
      <!-- Title and Badges -->
      <div class="flex items-start justify-between mb-2 gap-2">
        <h3 class="text-sm font-semibold text-text-primary line-clamp-1 flex-1 min-w-0">
          {{ template.name }}
        </h3>
        <div class="flex items-center gap-1 flex-shrink-0">
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

      <!-- Description -->
      <p class="text-xs text-text-secondary mb-2 line-clamp-2 min-h-[2rem]">
        {{ template.description || '暂无描述' }}
      </p>

      <!-- Category and Stats -->
      <div class="flex items-center justify-between mb-2 pb-2 border-b border-border-primary">
        <span
          class="inline-block px-2 py-0.5 rounded text-xs font-medium bg-primary/10 text-primary border border-primary/20"
        >
          {{ getCategoryName(template.category) }}
        </span>
        <div class="flex items-center gap-2.5 text-xs">
          <div class="flex items-center gap-1 text-green-600 dark:text-green-400">
            <Download class="w-3.5 h-3.5" />
            <span>{{ template.install_count }}</span>
          </div>
          <div class="flex items-center gap-1 text-blue-600 dark:text-blue-400">
            <Eye class="w-3.5 h-3.5" />
            <span>{{ template.view_count }}</span>
          </div>
        </div>
      </div>

      <!-- Required Tools -->
      <div v-if="template.required_tools && template.required_tools.length > 0" class="mb-2">
        <div class="flex items-center gap-1 mb-1">
          <Wrench class="w-3 h-3 text-text-tertiary flex-shrink-0" />
          <span class="text-xs font-medium text-text-tertiary">所需工具:</span>
          <span class="text-xs text-text-secondary truncate flex-1">
            {{ template.required_tools.join(', ') }}
          </span>
        </div>
      </div>

      <!-- Install Button -->
      <BaseButton class="w-full text-sm py-1.5" @click.stop="$emit('install', template)">
        <Download class="w-3.5 h-3.5 mr-1" />
        安装工作流
      </BaseButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Shield, Star, Wrench, Download, Eye } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import type { TemplateBasicInfo } from '@/api/template'

defineProps<{
  template: TemplateBasicInfo
}>()

defineEmits<{
  view: [template: TemplateBasicInfo]
  install: [template: TemplateBasicInfo]
}>()

const getCategoryName = (category: string): string => {
  // 直接返回分类名称，因为现在从API获取的就是中文名称
  return category || '其他'
}

const formatDate = (dateStr: string): string => {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days}天前`
  if (days < 30) return `${Math.floor(days / 7)}周前`
  if (days < 365) return `${Math.floor(days / 30)}个月前`
  return `${Math.floor(days / 365)}年前`
}
</script>
