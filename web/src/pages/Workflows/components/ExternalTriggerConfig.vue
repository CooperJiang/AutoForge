<template>
  <div class="space-y-4">
    <!-- 节点说明 -->
    <div class="bg-blue-50 border border-blue-200 rounded-lg p-3">
      <div class="flex items-start gap-2">
        <Globe class="w-4 h-4 text-blue-600 mt-0.5 flex-shrink-0" />
        <div class="text-xs text-blue-800">
          <p class="font-medium mb-1">外部 API 触发节点</p>
          <p class="text-blue-700">
            此节点定义工作流可接收的外部参数。启用工作流 API 后，外部系统可通过 API
            调用并传入这些参数。
          </p>
        </div>
      </div>
    </div>

    <!-- 参数列表 -->
    <div>
      <div class="flex items-center justify-between mb-2">
        <label class="text-xs font-medium text-text-primary">参数定义</label>
        <BaseButton size="xs" variant="ghost" @click="addParameter">
          <Plus class="w-3 h-3 mr-1" />
          添加参数
        </BaseButton>
      </div>

      <div
        v-if="params.length === 0"
        class="text-center py-8 border border-dashed border-border-primary rounded-lg"
      >
        <div class="text-text-tertiary text-xs">
          <Settings class="w-8 h-8 mx-auto mb-2 opacity-50" />
          <p>暂无参数</p>
          <p class="mt-1">点击上方按钮添加参数</p>
        </div>
      </div>

      <div v-else class="space-y-3">
        <div
          v-for="(param, index) in params"
          :key="index"
          class="border border-border-primary rounded-lg p-3 space-y-2 hover:border-primary transition-colors"
        >
          <!-- 参数名称和操作 -->
          <div class="flex items-center justify-between gap-2">
            <BaseInput
              v-model="param.key"
              placeholder="参数名称 (如: prompt)"
              size="sm"
              class="flex-1"
              @input="emitUpdate"
            />
            <BaseButton
              size="xs"
              variant="ghost"
              @click="removeParameter(index)"
              class="text-red-500 hover:text-red-600"
            >
              <Trash2 class="w-4 h-4" />
            </BaseButton>
          </div>

          <!-- 参数类型 -->
          <div>
            <label class="text-xs text-text-secondary mb-1 block">类型</label>
            <BaseSelect
              v-model="param.type"
              :options="typeOptions"
              size="sm"
              @change="emitUpdate"
            />
          </div>

          <!-- 必填选项 -->
          <div class="flex items-center gap-2">
            <input
              :id="`required-${index}`"
              v-model="param.required"
              type="checkbox"
              class="w-4 h-4 rounded border-border-primary text-primary focus:ring-primary cursor-pointer"
              @change="emitUpdate"
            />
            <label :for="`required-${index}`" class="text-xs text-text-secondary cursor-pointer">
              必填参数
            </label>
          </div>

          <!-- 默认值 -->
          <div>
            <label class="text-xs text-text-secondary mb-1 block">默认值</label>
            <BaseInput
              v-model="param.defaultValue"
              placeholder="默认值（可选）"
              size="sm"
              @input="emitUpdate"
            />
          </div>

          <!-- 描述 -->
          <div>
            <label class="text-xs text-text-secondary mb-1 block">描述</label>
            <textarea
              v-model="param.description"
              placeholder="参数说明（可选）"
              rows="2"
              class="w-full px-2 py-1.5 text-xs border border-border-primary rounded bg-bg-elevated text-text-primary resize-none focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent transition-all"
              @input="emitUpdate"
            />
          </div>

          <!-- 示例值 -->
          <div>
            <label class="text-xs text-text-secondary mb-1 block">示例值</label>
            <BaseInput
              v-model="param.example"
              placeholder="示例值（可选）"
              size="sm"
              @input="emitUpdate"
            />
          </div>
        </div>
      </div>
    </div>

    <!-- Webhook 配置 -->
    <div>
      <label class="text-xs font-medium text-text-primary mb-2 block"
        >Webhook 回调地址（可选）</label
      >
      <BaseInput
        v-model="webhookURL"
        type="url"
        placeholder="https://example.com/webhook"
        @input="emitUpdate"
      />
      <p class="text-xs text-text-tertiary mt-1">异步执行完成后，会将结果 POST 到此 URL</p>
    </div>

    <!-- 变量引用提示 -->
    <div class="bg-bg-tertiary rounded-lg p-3">
      <div class="flex items-start gap-2">
        <Info class="w-4 h-4 text-primary mt-0.5 flex-shrink-0" />
        <div class="text-xs text-text-secondary">
          <p class="font-medium text-text-primary mb-1">后续节点中引用这些参数</p>
          <p class="mb-1">在其他节点的配置中，使用以下变量格式：</p>
          <code
            class="text-xs bg-bg-elevated px-1.5 py-0.5 rounded text-primary"
            v-text="'{{external.参数名}}'"
          ></code>
          <div v-if="params.length > 0" class="mt-2 space-y-1">
            <p class="font-medium text-text-primary">当前可用变量：</p>
            <div v-for="param in params" :key="param.key" class="flex items-center gap-1">
              <code
                class="text-xs bg-bg-elevated px-1.5 py-0.5 rounded text-primary"
                v-text="`{{external.${param.key}}}`"
              ></code>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Globe, Plus, Trash2, Settings, Info } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'

interface Parameter {
  key: string
  type: 'string' | 'number' | 'boolean' | 'object' | 'array'
  required: boolean
  defaultValue: string
  description: string
  example: string
}

const props = defineProps<{
  config: any
}>()

const emit = defineEmits<{
  update: [config: any]
}>()

const params = ref<Parameter[]>(props.config?.params || [])
const webhookURL = ref<string>(props.config?.webhookURL || '')

// 类型选项
const typeOptions = [
  { label: '字符串 (string)', value: 'string' },
  { label: '数字 (number)', value: 'number' },
  { label: '布尔值 (boolean)', value: 'boolean' },
  { label: '对象 (object)', value: 'object' },
  { label: '数组 (array)', value: 'array' },
]

// 监听 props 变化
watch(
  () => props.config,
  (newConfig) => {
    if (newConfig) {
      params.value = newConfig.params || []
      webhookURL.value = newConfig.webhookURL || ''
    }
  },
  { deep: true }
)

// 添加参数
const addParameter = () => {
  params.value.push({
    key: '',
    type: 'string',
    required: false,
    defaultValue: '',
    description: '',
    example: '',
  })
  emitUpdate()
}

// 删除参数
const removeParameter = (index: number) => {
  params.value.splice(index, 1)
  emitUpdate()
}

// 触发更新
const emitUpdate = () => {
  emit('update', {
    params: params.value,
    webhookURL: webhookURL.value,
  })
}
</script>
