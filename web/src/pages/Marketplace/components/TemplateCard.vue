<template>
  <div
    class="bg-bg-elevated rounded-xl border-2 border-border-primary hover:border-primary hover:shadow-xl transition-all cursor-pointer p-5 shadow-sm"
    @click="$emit('view', template)"
  >
    <!-- Header with Icon -->
    <div class="flex items-start gap-3 mb-3">
      <div
        class="w-14 h-14 rounded-xl bg-gradient-to-br from-primary/20 to-accent/20 flex items-center justify-center flex-shrink-0 border border-border-primary"
      >
        <span class="text-3xl">{{ template.icon || '📦' }}</span>
      </div>
      <div class="flex-1 min-w-0">
        <h3 class="text-base font-semibold text-text-primary mb-1 truncate">
          {{ template.name }}
        </h3>
        <div class="flex items-center gap-2">
          <span
            v-if="template.is_official"
            class="inline-flex items-center px-2 py-0.5 rounded text-xs bg-primary/10 text-primary"
          >
            <Shield class="w-3 h-3 mr-1" />
            官方
          </span>
          <span
            v-if="template.is_featured"
            class="inline-flex items-center px-2 py-0.5 rounded text-xs bg-accent/10 text-accent"
          >
            <Star class="w-3 h-3 mr-1 fill-current" />
            精选
          </span>
        </div>
      </div>
    </div>

    <!-- Description -->
    <p class="text-sm text-text-secondary mb-3 line-clamp-2">
      {{ template.description || '暂无描述' }}
    </p>

    <!-- Category -->
    <div class="mb-3">
      <span class="inline-block px-2.5 py-1 rounded-md text-xs font-medium bg-primary/10 text-primary border border-primary/20">
        {{ getCategoryName(template.category) }}
      </span>
    </div>

    <!-- Required Tools -->
    <div v-if="template.required_tools && template.required_tools.length > 0" class="mb-3">
      <div class="flex items-center gap-1 mb-2">
        <Wrench class="w-3.5 h-3.5 text-text-tertiary" />
        <span class="text-xs font-medium text-text-tertiary">所需工具</span>
      </div>
      <div class="flex flex-wrap gap-1.5">
        <span
          v-for="tool in template.required_tools.slice(0, 3)"
          :key="tool"
          class="px-2 py-1 rounded-md text-xs bg-surface-secondary text-text-secondary border border-border-primary"
        >
          {{ tool }}
        </span>
        <span
          v-if="template.required_tools.length > 3"
          class="px-2 py-1 rounded-md text-xs text-text-tertiary font-medium"
        >
          +{{ template.required_tools.length - 3 }}
        </span>
      </div>
    </div>

    <!-- Stats -->
    <div
      class="flex items-center justify-between text-xs text-text-tertiary pt-3 border-t border-border-primary"
    >
      <div class="flex items-center gap-3">
        <div class="flex items-center gap-1">
          <Download class="w-3 h-3" />
          <span>{{ template.install_count }}</span>
        </div>
        <div class="flex items-center gap-1">
          <Eye class="w-3 h-3" />
          <span>{{ template.view_count }}</span>
        </div>
      </div>
      <div class="text-text-tertiary">
        {{ formatDate(template.created_at) }}
      </div>
    </div>

    <!-- Install Button -->
    <BaseButton class="w-full mt-3" @click.stop="$emit('install', template)">
      <Download class="w-4 h-4 mr-1.5" />
      安装模板
    </BaseButton>
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
