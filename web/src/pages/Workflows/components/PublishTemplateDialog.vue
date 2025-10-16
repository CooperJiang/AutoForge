<template>
  <Dialog
    :modelValue="visible"
    title="发布为模板"
    max-width="max-w-2xl"
    @update:modelValue="handleClose"
  >
    <div class="space-y-4">
      <div>
        <label class="block text-sm font-medium text-text-primary mb-2">模板名称 *</label>
        <BaseInput v-model="form.name" placeholder="输入模板名称" />
      </div>

      <div>
        <label class="block text-sm font-medium text-text-primary mb-2">描述</label>
        <textarea
          v-model="form.description"
          placeholder="输入模板描述"
          rows="3"
          class="w-full px-3 py-2 rounded-lg bg-surface-secondary border border-border-primary text-text-primary placeholder:text-text-placeholder focus:outline-none focus:border-primary resize-none"
        ></textarea>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-primary mb-2">分类 *</label>
        <div class="flex gap-2">
          <BaseSelect
            v-model="form.category"
            :options="categoryOptions"
            placeholder="选择分类"
            class="flex-1"
          />
          <BaseInput
            v-if="form.category === 'custom'"
            v-model="customCategory"
            placeholder="输入自定义分类"
            class="flex-1"
          />
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-primary mb-2">图标类型</label>
        <div class="flex gap-2 mb-2">
          <BaseButton
            size="sm"
            :variant="iconType === 'icon' ? 'primary' : 'outline'"
            @click="iconType = 'icon'"
          >
            图标库
          </BaseButton>
          <BaseButton
            size="sm"
            :variant="iconType === 'url' ? 'primary' : 'outline'"
            @click="iconType = 'url'"
          >
            图片 URL
          </BaseButton>
        </div>

        <!-- 图标库选择器 -->
        <div v-if="iconType === 'icon'" class="space-y-2">
          <!-- 搜索框 -->
          <BaseInput
            v-model="iconSearchKeyword"
            placeholder="搜索图标..."
            @input="handleIconSearch"
          />

          <!-- 图标网格 -->
          <div
            class="grid grid-cols-8 gap-2 p-3 bg-surface-secondary rounded-lg border border-border-primary max-h-60 overflow-y-auto"
          >
            <button
              v-for="iconName in filteredIcons"
              :key="iconName"
              @click="form.icon = iconName"
              :class="[
                'p-3 rounded hover:bg-surface-tertiary transition-colors flex items-center justify-center',
                form.icon === iconName && 'bg-primary/20 ring-2 ring-primary',
              ]"
              type="button"
              :title="iconName"
            >
              <component :is="getIconComponent(iconName)" class="w-5 h-5" />
            </button>
          </div>

          <!-- 当前选择 -->
          <div class="flex items-center gap-2 p-2 bg-surface-tertiary rounded">
            <component v-if="form.icon" :is="getIconComponent(form.icon)" class="w-5 h-5" />
            <span class="text-sm text-text-secondary">{{ form.icon || '未选择' }}</span>
          </div>
        </div>

        <!-- URL 输入 -->
        <div v-else>
          <BaseInput
            v-model="form.cover_image"
            placeholder="输入图片 URL，例如: https://example.com/icon.png"
          />
          <p class="text-xs text-text-tertiary mt-1">建议尺寸: 512x512 或 1:1 比例</p>
        </div>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-primary mb-2">使用指南</label>
        <textarea
          v-model="form.usage_guide"
          placeholder="输入使用指南和注意事项"
          rows="4"
          class="w-full px-3 py-2 rounded-lg bg-surface-secondary border border-border-primary text-text-primary placeholder:text-text-placeholder focus:outline-none focus:border-primary resize-none"
        ></textarea>
      </div>

      <div class="flex items-center gap-2">
        <input
          id="is_featured"
          v-model="form.is_featured"
          type="checkbox"
          class="w-4 h-4 rounded border-border-primary text-primary focus:ring-primary"
        />
        <label for="is_featured" class="text-sm text-text-primary cursor-pointer">
          设为精选模板
        </label>
      </div>
    </div>

    <template #footer>
      <BaseButton variant="secondary" @click="handleClose">取消</BaseButton>
      <BaseButton @click="handlePublish" :disabled="!isValid || publishing">
        <Package class="w-4 h-4 mr-1" />
        {{ publishing ? '发布中...' : '发布模板' }}
      </BaseButton>
    </template>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import * as LucideIcons from 'lucide-vue-next'
import Dialog from '@/components/Dialog'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import { templateApi } from '@/api/template'
import { message } from '@/utils/message'

const { Package } = LucideIcons

const props = defineProps<{
  visible: boolean
  workflowId: string
}>()

const emit = defineEmits<{
  close: []
  success: []
}>()

const publishing = ref(false)
const iconType = ref<'icon' | 'url'>('icon')
const customCategory = ref('')
const iconSearchKeyword = ref('')

const form = ref({
  name: '',
  description: '',
  category: '',
  icon: 'Package',
  cover_image: '',
  usage_guide: '',
  is_featured: false,
})

const categoryOptions = [
  { label: '请选择', value: '' },
  { label: '自动化', value: 'automation' },
  { label: '通知', value: 'notification' },
  { label: '数据处理', value: 'data' },
  { label: '集成', value: 'integration' },
  { label: '监控', value: 'monitoring' },
  { label: '定时任务', value: 'scheduled' },
  { label: '自定义分类', value: 'custom' },
]

// 常用图标列表（从 lucide-vue-next）
const commonIcons = [
  'Package',
  'Rocket',
  'Zap',
  'Wrench',
  'Settings',
  'Bell',
  'Mail',
  'MessageSquare',
  'BarChart',
  'TrendingUp',
  'Database',
  'HardDrive',
  'Search',
  'Clock',
  'Calendar',
  'CheckCircle',
  'XCircle',
  'AlertTriangle',
  'Target',
  'Palette',
  'Film',
  'FileText',
  'Link',
  'Globe',
  'Lock',
  'Key',
  'Lightbulb',
  'Gift',
  'Star',
  'Sparkles',
  'Flame',
  'Droplet',
  'Cloud',
  'Sun',
  'Moon',
  'Smartphone',
  'Laptop',
  'Code',
  'GitBranch',
  'Users',
  'User',
  'Heart',
  'ThumbsUp',
  'Share2',
  'Download',
  'Upload',
  'RefreshCw',
  'PlayCircle',
  'PauseCircle',
  'StopCircle',
  'Trash2',
  'Edit3',
  'Copy',
  'Folder',
  'FolderOpen',
  'File',
  'Image',
  'Music',
  'Video',
  'Archive',
  'Bookmark',
  'Tag',
  'Filter',
  'Sliders',
  'Cpu',
  'HardDrive',
  'Server',
  'Shield',
  'Eye',
  'Layers',
  'Grid',
  'List',
  'Layout',
  'Maximize2',
  'Minimize2',
  'Move',
  'Navigation',
]

// 过滤后的图标列表
const filteredIcons = computed(() => {
  if (!iconSearchKeyword.value) {
    return commonIcons
  }
  const keyword = iconSearchKeyword.value.toLowerCase()
  return commonIcons.filter((icon) => icon.toLowerCase().includes(keyword))
})

// 获取图标组件
const getIconComponent = (iconName: string) => {
  return (LucideIcons as any)[iconName] || Package
}

// 处理图标搜索
const handleIconSearch = () => {
  // 搜索逻辑已在 computed 中处理
}

const isValid = computed(() => {
  const hasName = form.value.name.trim() !== ''
  const hasCategory = form.value.category !== ''
  const hasValidCategory = form.value.category !== 'custom' || customCategory.value.trim() !== ''
  return hasName && hasCategory && hasValidCategory
})

const handleClose = () => {
  emit('close')
}

const handlePublish = async () => {
  if (!isValid.value) {
    message.error('请填写必填项')
    return
  }

  publishing.value = true
  try {
    // 处理自定义分类
    const finalCategory =
      form.value.category === 'custom' ? customCategory.value.trim() : form.value.category

    // 构建请求数据
    const requestData: any = {
      name: form.value.name,
      description: form.value.description,
      category: finalCategory,
      workflow_id: props.workflowId,
      usage_guide: form.value.usage_guide,
      is_featured: form.value.is_featured,
    }

    // 根据图标类型设置图标或封面图
    if (iconType.value === 'emoji') {
      requestData.icon = form.value.icon
    } else {
      requestData.cover_image = form.value.cover_image
    }

    await templateApi.create(requestData)

    message.success('模板发布成功')
    emit('success')
    emit('close')

    // Reset form
    form.value = {
      name: '',
      description: '',
      category: '',
      icon: 'Package',
      cover_image: '',
      usage_guide: '',
      is_featured: false,
    }
    customCategory.value = ''
    iconType.value = 'icon'
    iconSearchKeyword.value = ''
  } catch (error: any) {
    console.error('Failed to publish template:', error)
    message.error(error.response?.data?.message || '发布失败')
  } finally {
    publishing.value = false
  }
}
</script>
