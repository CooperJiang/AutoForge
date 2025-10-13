<template>
  <div
    class="group bg-gradient-to-br from-bg-elevated to-bg-secondary rounded-xl shadow-md hover:shadow-2xl transition-all duration-300 p-4 cursor-pointer border-2 border-border-primary hover:border-primary hover:-translate-y-1 flex flex-col min-h-[180px]"
    @click="$emit('click')"
  >
    <!-- 工具图标和标题 -->
    <div class="flex items-start gap-3 mb-3">
      <div
        :class="[
          'flex-shrink-0 w-12 h-12 rounded-xl flex items-center justify-center text-white',
          'shadow-lg group-hover:scale-110 transition-transform duration-300',
          iconBgClass
        ]"
      >
        <component :is="iconComponent" class="w-6 h-6" />
      </div>
      <div class="flex-1 min-w-0">
        <h3 class="text-base font-bold text-text-primary mb-1 group-hover:text-primary transition-colors truncate">
          {{ tool.name }}
        </h3>
        <p class="text-xs text-text-tertiary truncate">
          {{ tool.version }} · {{ tool.category }}
        </p>
      </div>
    </div>

    <!-- 工具描述 -->
    <p class="text-sm text-text-secondary mb-3 line-clamp-2 flex-shrink-0 leading-relaxed">
      {{ tool.description }}
    </p>

    <!-- 工具标签 -->
    <div class="flex flex-wrap gap-1.5 mb-3 flex-shrink-0">
      <span
        v-for="tag in tool.tags.slice(0, 3)"
        :key="tag"
        class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-primary-light text-primary"
      >
        {{ tag }}
      </span>
      <span
        v-if="tool.tags.length > 3"
        class="inline-flex items-center px-2 py-1 rounded-md text-xs font-medium bg-bg-tertiary text-text-secondary"
      >
        +{{ tool.tags.length - 3 }}
      </span>
    </div>

    <!-- 底部信息 -->
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
import { Globe, Mail, HeartPulse, ChevronRight } from 'lucide-vue-next'

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

// 根据工具代码获取图标
const iconComponent = computed(() => {
  const iconMap: Record<string, any> = {
    'http_request': Globe,
    'email_sender': Mail,
    'health_checker': HeartPulse
  }
  return iconMap[props.tool.code] || Globe
})

// 根据工具代码获取图标背景色
const iconBgClass = computed(() => {
  const colorMap: Record<string, string> = {
    'http_request': 'bg-gradient-to-br from-primary to-accent',
    'email_sender': 'bg-gradient-to-br from-purple-500 to-pink-600',
    'health_checker': 'bg-gradient-to-br from-primary to-accent'
  }
  return colorMap[props.tool.code] || 'bg-gradient-to-br from-primary to-accent'
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
