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
        <label class="block text-sm font-medium text-text-primary mb-2">封面图片 URL</label>
        <BaseInput
          v-model="form.cover_image"
          placeholder="输入图片 URL（可选），例如: https://example.com/cover.png"
        />
        <p class="text-xs text-text-tertiary mt-1">建议尺寸: 512x512 或 1:1 比例</p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-primary mb-2">图标 Emoji</label>
        <BaseInput
          v-model="form.icon"
          placeholder="输入 Emoji 图标（可选），例如: 📦"
          maxlength="10"
        />
        <p class="text-xs text-text-tertiary mt-1">用于在列表中显示，留空则使用默认图标</p>
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
import { ref, computed, watch } from 'vue'
import Dialog from '@/components/Dialog'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import { templateApi } from '@/api/template'
import type { TemplateCategory } from '@/api/template'
import { message } from '@/utils/message'

const props = defineProps<{
  visible: boolean
  workflowId: string
}>()

const emit = defineEmits<{
  close: []
  success: []
}>()

const publishing = ref(false)
const customCategory = ref('')
const categories = ref<TemplateCategory[]>([])

const form = ref({
  name: '',
  description: '',
  category: '',
  icon: '',
  cover_image: '',
  usage_guide: '',
  is_featured: false,
})

// 加载分类列表
const loadCategories = async () => {
  try {
    const res = await templateApi.listCategories({ page_size: 100, is_active: true })
    categories.value = res.items
  } catch (error: any) {
    console.error('加载分类失败:', error)
  }
}

// 监听对话框打开时加载分类
watch(
  () => props.visible,
  (newVal) => {
    if (newVal) {
      loadCategories()
    }
  }
)

// 动态构建分类选项
const categoryOptions = computed(() => {
  const options = [{ label: '请选择', value: '' }]
  categories.value.forEach((cat) => {
    options.push({
      label: cat.name,
      value: cat.name,
    })
  })
  options.push({ label: '+ 自定义分类', value: 'custom' })
  return options
})

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
      icon: form.value.icon,
      cover_image: form.value.cover_image,
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
      icon: '',
      cover_image: '',
      usage_guide: '',
      is_featured: false,
    }
    customCategory.value = ''
  } catch (error: any) {
    console.error('Failed to publish template:', error)
    message.error(error.response?.data?.message || '发布失败')
  } finally {
    publishing.value = false
  }
}
</script>
