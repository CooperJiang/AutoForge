<template>
  <div
    class="group bg-white rounded-lg shadow-sm hover:shadow-xl transition-all duration-300 p-2.5 cursor-pointer border border-slate-200 hover:border-blue-400 hover:-translate-y-1 flex flex-col h-[155px]"
    @click="$emit('click')"
  >
    <!-- 工具图标和标题 -->
    <div class="flex items-start gap-2 mb-1.5">
      <div
        :class="[
          'flex-shrink-0 w-9 h-9 rounded-lg flex items-center justify-center text-white',
          'shadow-md group-hover:scale-110 transition-transform duration-300',
          iconBgClass
        ]"
      >
        <component :is="iconComponent" class="w-4.5 h-4.5" />
      </div>
      <div class="flex-1 min-w-0">
        <h3 class="text-sm font-semibold text-slate-900 mb-0.5 group-hover:text-blue-600 transition-colors truncate">
          {{ tool.name }}
        </h3>
        <p class="text-xs text-slate-500 truncate">
          {{ tool.version }} · {{ tool.category }}
        </p>
      </div>
    </div>

    <!-- 工具描述 -->
    <p class="text-xs text-slate-600 mb-1.5 line-clamp-2 flex-shrink-0">
      {{ tool.description }}
    </p>

    <!-- 工具标签 -->
    <div class="flex flex-wrap gap-1 mb-1.5 flex-shrink-0">
      <span
        v-for="tag in tool.tags.slice(0, 2)"
        :key="tag"
        class="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-slate-100 text-slate-700"
      >
        {{ tag }}
      </span>
      <span
        v-if="tool.tags.length > 2"
        class="inline-flex items-center px-1.5 py-0.5 rounded text-xs font-medium bg-slate-100 text-slate-500"
      >
        +{{ tool.tags.length - 2 }}
      </span>
    </div>

    <!-- 底部信息 -->
    <div class="flex items-center justify-between pt-1.5 border-t border-slate-100 mt-auto">
      <span class="text-xs text-slate-500 truncate mr-2">
        {{ tool.author }}
      </span>
      <button
        class="flex items-center gap-0.5 text-xs font-medium text-blue-600 hover:text-blue-700 group-hover:gap-1 transition-all flex-shrink-0"
        @click.stop="$emit('click')"
      >
        查看
        <ChevronRight class="w-3 h-3" />
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
    'http_request': 'bg-gradient-to-br from-blue-500 to-purple-600',
    'email_sender': 'bg-gradient-to-br from-purple-500 to-pink-600',
    'health_checker': 'bg-gradient-to-br from-indigo-500 to-blue-600'
  }
  return colorMap[props.tool.code] || 'bg-gradient-to-br from-blue-500 to-purple-600'
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
