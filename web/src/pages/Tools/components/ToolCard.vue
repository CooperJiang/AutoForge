<template>
  <div
    class="group bg-gradient-to-br from-bg-elevated to-bg-secondary rounded-xl shadow-md hover:shadow-2xl transition-all duration-300 p-4 cursor-pointer border-2 border-border-primary hover:border-primary hover:-translate-y-1 flex flex-col min-h-[180px]"
    @click="$emit('click')"
  >
    <div class="flex items-start gap-3 mb-3">
      <div
        :class="[
          'flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center text-white',
          'shadow-lg group-hover:scale-110 transition-transform duration-300',
          toolIconBg,
        ]"
      >
        <component v-if="isLucideIcon" :is="toolIcon" class="w-6 h-6" />
        <img v-else :src="toolIcon" alt="Tool Icon" class="w-6 h-6 object-contain" />
      </div>
      <div class="flex-1 min-w-0">
        <h3
          class="text-base font-bold text-text-primary mb-1 group-hover:text-primary transition-colors truncate"
        >
          {{ tool.name }}
        </h3>
        <p class="text-xs text-text-tertiary truncate">{{ tool.version }} · {{ tool.category }}</p>
      </div>
    </div>

    <p class="text-sm text-text-secondary mb-3 line-clamp-2 flex-shrink-0 leading-relaxed">
      {{ tool.description }}
    </p>

    <div class="flex flex-wrap gap-1.5 mb-3 flex-shrink-0">
      <span
        v-for="tag in displayTags.slice(0, 3)"
        :key="tag"
        class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-primary-light text-primary"
      >
        {{ tag }}
      </span>
      <span
        v-if="displayTags.length > 3"
        class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-bg-tertiary text-text-secondary"
      >
        +{{ displayTags.length - 3 }}
      </span>
    </div>

    <div class="flex items-center justify-between pt-3 border-t-2 border-border-primary mt-auto">
      <span class="text-xs text-text-tertiary truncate mr-2 font-medium">
        {{ tool.author }}
      </span>
      <button
        class="flex items-center gap-1 text-sm font-semibold text-primary hover:text-primary group-hover:gap-1.5 transition-all flex-shrink-0"
        @click.stop="$emit('click')"
      >
        查看详情
        <ChevronRight class="w-4 h-4" />
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { ChevronRight } from 'lucide-vue-next'
import { getToolConfig, getToolIcon, getToolIconBg } from '@/config/tools'

interface Tool {
  code: string
  name: string
  description: string
  category: string
  icon: string
  tags: string[]
  version: string
  author: string
}

interface Props {
  tool: Tool
}

const props = defineProps<Props>()

defineEmits<{
  click: []
}>()

const toolConfig = computed(() => {
  return getToolConfig(props.tool.code)
})

const toolIcon = computed(() => {
  return getToolIcon(props.tool.code)
})

const toolIconBg = computed(() => {
  return getToolIconBg(props.tool.code)
})

const isLucideIcon = computed(() => {
  return typeof toolIcon.value !== 'string'
})

const displayTags = computed(() => {
  if (toolConfig.value?.tags && toolConfig.value.tags.length > 0) {
    return toolConfig.value.tags
  }

  if (Array.isArray(props.tool.tags)) {
    return props.tool.tags
  }

  if (typeof props.tool.tags === 'string') {
    return props.tool.tags
      .split(',')
      .map((tag) => tag.trim())
      .filter(Boolean)
  }

  return []
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
