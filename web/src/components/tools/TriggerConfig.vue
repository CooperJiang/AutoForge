<template>
  <div class="space-y-4">
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        触发器类型 <span class="text-red-500">*</span>
      </label>
      <BaseSelect
        v-model="localConfig.triggerType"
        :options="triggerTypeOptions"
        @update:model-value="handleTriggerTypeChange"
      />
    </div>

    <div v-if="localConfig.triggerType === 'schedule'">
      <div class="bg-primary-light border-l-4 border-border-focus p-3 mb-3">
        <p class="text-sm text-primary">
          <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path
              fill-rule="evenodd"
              d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
              clip-rule="evenodd"
            />
          </svg>
          配置定时触发器的执行计划
        </p>
      </div>

      <ScheduleSelector
        :type="localConfig.scheduleType"
        :value="localConfig.scheduleValue"
        @update:type="updateScheduleType"
        @update:value="updateScheduleValue"
      />
    </div>

    <div v-if="localConfig.triggerType === 'webhook'" class="space-y-3">
      <div class="bg-info-light border-l-4 border-info rounded-lg p-3">
        <p class="text-sm text-info-text">
          <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path
              fill-rule="evenodd"
              d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
              clip-rule="evenodd"
            />
          </svg>
          通过 HTTP 请求触发工作流执行，适用于接收外部系统的回调通知
        </p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          请求方法 <span class="text-red-500">*</span>
        </label>
        <BaseSelect v-model="localConfig.webhookMethod" :options="methodOptions" />
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          Webhook 路径 <span class="text-red-500">*</span>
        </label>
        <div class="flex items-center gap-2">
          <span class="text-sm text-text-tertiary whitespace-nowrap">/api/webhook/</span>
          <BaseInput
            v-model="localConfig.webhookPath"
            placeholder="my-workflow"
            @update:model-value="emitUpdate"
          />
        </div>
        <p class="text-xs text-text-tertiary mt-1">
          完整URL:
          <code class="px-1 py-0.5 bg-bg-tertiary rounded"
            >https://your-domain.com/api/webhook/{{
              localConfig.webhookPath || 'my-workflow'
            }}</code
          >
        </p>
      </div>

      <div class="bg-bg-hover rounded-lg p-3">
        <div class="text-xs font-semibold text-text-secondary mb-2">使用示例：</div>
        <div class="bg-slate-900 text-slate-100 rounded p-2 font-mono text-xs overflow-x-auto">
          <div>curl -X {{ localConfig.webhookMethod || 'POST' }} \</div>
          <div class="ml-2">
            https://your-domain.com/api/webhook/{{ localConfig.webhookPath || 'my-workflow' }} \
          </div>
          <div class="ml-2">-H "Content-Type: application/json" \</div>
          <div class="ml-2">-d '{"key": "value"}'</div>
        </div>
        <p class="text-xs text-text-secondary mt-2">
          请求体数据可通过
          <code class="px-1 py-0.5 bg-bg-tertiary rounded"
            >&#123;&#123;trigger.data&#125;&#125;</code
          >
          在后续节点中访问
        </p>
      </div>
    </div>

    <div v-if="localConfig.triggerType === 'manual'">
      <div class="bg-success-light border-l-4 border-green-400 p-3">
        <p class="text-sm text-success">
          <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path
              fill-rule="evenodd"
              d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
              clip-rule="evenodd"
            />
          </svg>
          手动触发，需要在工作流列表页面手动点击执行按钮
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import ScheduleSelector from '@/pages/Tasks/components/ScheduleSelector.vue'
import BaseSelect from '@/components/BaseSelect'
import BaseInput from '@/components/BaseInput'

interface Props {
  config: Record<string, any>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:config': [config: Record<string, any>]
}>()

// 初始化配置，确保有默认值
const initConfig = () => {
  return {
    triggerType: props.config.triggerType || 'schedule',
    scheduleType: props.config.scheduleType || 'daily',
    scheduleValue: props.config.scheduleValue || '09:00:00',
    webhookMethod: props.config.webhookMethod || 'POST',
    webhookPath: props.config.webhookPath || '',
    enabled: props.config.enabled !== undefined ? props.config.enabled : true,
    ...props.config, // 最后用 props.config 覆盖所有字段
  }
}

const localConfig = ref(initConfig())

const triggerTypeOptions = [
  { label: '定时触发', value: 'schedule' },
  { label: 'Webhook 触发', value: 'webhook' },
  { label: '手动触发', value: 'manual' },
]

const methodOptions = [
  { label: 'POST', value: 'POST' },
  { label: 'GET', value: 'GET' },
  { label: 'PUT', value: 'PUT' },
]

// 监听外部配置变化，但保留 localConfig 中已有的值
watch(
  () => props.config,
  (newVal) => {
    // 只更新外部传入的字段，保留本地已有的其他字段
    Object.keys(newVal).forEach((key) => {
      if (newVal[key] !== undefined) {
        localConfig.value[key] = newVal[key]
      }
    })
  },
  { deep: true }
)

const handleTriggerTypeChange = () => {
  // 重置相关配置
  if (localConfig.value.triggerType === 'schedule') {
    localConfig.value.scheduleType = 'daily'
    localConfig.value.scheduleValue = '09:00:00'
  } else if (localConfig.value.triggerType === 'webhook') {
    localConfig.value.webhookMethod = 'POST'
    localConfig.value.webhookPath = ''
  }
  emitUpdate()
}

const updateScheduleType = (type: string) => {
  const oldType = localConfig.value.scheduleType
  localConfig.value.scheduleType = type

  // 只在切换到不同类型时设置默认值
  if (oldType !== type) {
    if (type === 'daily') {
      localConfig.value.scheduleValue = '09:00:00'
    } else if (type === 'hourly') {
      localConfig.value.scheduleValue = '00:00'
    } else if (type === 'interval') {
      localConfig.value.scheduleValue = '300'
    } else if (type === 'weekly') {
      localConfig.value.scheduleValue = '1:09:00:00'
    } else if (type === 'monthly') {
      localConfig.value.scheduleValue = '1:09:00:00'
    } else if (type === 'cron') {
      localConfig.value.scheduleValue = '0 0 * * * *'
    }
  }
  emitUpdate()
}

const updateScheduleValue = (value: string) => {
  localConfig.value.scheduleValue = value
  emitUpdate()
}

const emitUpdate = () => {
  emit('update:config', localConfig.value)
}
</script>
