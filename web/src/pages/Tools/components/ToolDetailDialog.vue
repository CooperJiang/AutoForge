<template>
  <Dialog
    v-model="isOpen"
    :title="tool?.name || '工具详情'"
    max-width="max-w-3xl"
    hide-footer
    @cancel="close"
  >
    <div v-if="tool" class="space-y-6">
      <div class="flex items-start gap-4">
        <div
          :class="[
            'flex-shrink-0 w-16 h-16 rounded-xl flex items-center justify-center text-white shadow-lg',
            toolIconBg,
          ]"
        >
          <component v-if="isLucideIcon" :is="toolIcon" class="w-8 h-8" />
          <img v-else :src="toolIcon" alt="Tool Icon" class="w-8 h-8 object-contain" />
        </div>
        <div class="flex-1">
          <p class="text-text-secondary mb-2">
            {{ tool.description }}
          </p>
          <div class="flex flex-wrap items-center gap-3 text-sm">
            <span class="flex items-center gap-1 text-text-secondary">
              <Package class="w-4 h-4" />
              版本 {{ tool.version }}
            </span>
            <span class="flex items-center gap-1 text-text-secondary">
              <User class="w-4 h-4" />
              {{ tool.author }}
            </span>
            <span class="flex items-center gap-1 text-text-secondary">
              <Tag class="w-4 h-4" />
              {{ tool.category }}
            </span>
          </div>
        </div>
      </div>

      <!-- 工具标签 -->
      <div v-if="displayTags.length > 0">
        <h3 class="text-sm font-semibold text-text-secondary mb-2">标签</h3>
        <div class="flex flex-wrap gap-2">
          <span
            v-for="tag in displayTags"
            :key="tag"
            class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-primary-light text-primary border border-primary"
          >
            {{ tag }}
          </span>
        </div>
      </div>

      <!-- 使用说明 -->
      <div v-if="toolConfig">
        <h3 class="text-sm font-semibold text-text-primary mb-2 flex items-center gap-2">
          <BookOpen class="w-4 h-4" />
          使用说明
        </h3>
        <div class="bg-bg-hover rounded-lg p-4 space-y-3">
          <div>
            <h4 class="font-semibold text-text-primary mb-2 text-sm">{{ toolConfig.title }}</h4>
            <p v-if="toolConfig.usageDescription" class="text-xs text-text-secondary mb-2">
              {{ toolConfig.usageDescription }}
            </p>
            <ul
              v-if="toolConfig.usageItems.length > 0"
              class="space-y-1.5 text-xs text-text-secondary"
            >
              <li
                v-for="(item, index) in toolConfig.usageItems"
                :key="index"
                class="flex items-start gap-2"
              >
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>{{ item.text }}</span>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <div v-else>
        <h3 class="text-sm font-semibold text-text-primary mb-2 flex items-center gap-2">
          <BookOpen class="w-4 h-4" />
          使用说明
        </h3>
        <div class="bg-bg-hover rounded-lg p-4">
          <p class="text-xs text-text-secondary">暂无详细说明，请直接使用该工具。</p>
        </div>
      </div>

      <!-- 使用按钮 -->
      <div class="flex items-center justify-end gap-2 pt-2">
        <BaseButton variant="ghost" size="md" @click="close"> 取消 </BaseButton>
        <BaseButton variant="primary" size="md" @click="handleUseTool" class="whitespace-nowrap">
          <Rocket class="w-4 h-4 mr-1.5 inline-block flex-shrink-0" />
          <span>立即使用</span>
        </BaseButton>
      </div>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { computed, watch, ref } from 'vue'
import { Package, User, Tag, BookOpen, Check, Rocket } from 'lucide-vue-next'
import Dialog from '@/components/Dialog'
import BaseButton from '@/components/BaseButton'
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
  modelValue: boolean
  tool: Tool | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'use-tool': [toolCode: string]
}>()

const isOpen = ref(props.modelValue)

watch(
  () => props.modelValue,
  (val) => {
    isOpen.value = val
  }
)

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

const toolConfig = computed(() => {
  if (!props.tool) return null
  return getToolConfig(props.tool.code)
})

const toolIcon = computed(() => {
  if (!props.tool) return null
  return getToolIcon(props.tool.code)
})

const toolIconBg = computed(() => {
  if (!props.tool) return 'bg-gradient-to-br from-gray-500 to-gray-600'
  return getToolIconBg(props.tool.code)
})

const isLucideIcon = computed(() => {
  return typeof toolIcon.value !== 'string'
})

const displayTags = computed(() => {
  if (!props.tool) return []

  // Prefer tags from config, fallback to backend tags
  if (toolConfig.value?.tags && toolConfig.value.tags.length > 0) {
    return toolConfig.value.tags
  }

  // Parse backend tags (assuming it's a comma-separated string or array)
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

const close = () => {
  isOpen.value = false
}

const handleUseTool = () => {
  if (props.tool) {
    emit('use-tool', props.tool.code)
  }
  close()
}
</script>
