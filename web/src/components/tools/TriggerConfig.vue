<template>
  <div class="space-y-4">
    <!-- 触发器类型选择 -->
    <div>
      <label class="block text-sm font-medium text-slate-700 mb-2">
        触发器类型 <span class="text-red-500">*</span>
      </label>
      <BaseSelect
        v-model="localConfig.triggerType"
        :options="triggerTypeOptions"
        @update:model-value="handleTriggerTypeChange"
      />
    </div>

    <!-- 定时触发配置 -->
    <div v-if="localConfig.triggerType === 'schedule'">
      <div class="bg-blue-50 border-l-4 border-blue-400 p-3 mb-3">
        <p class="text-sm text-blue-700">
          <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
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

    <!-- Webhook 触发配置 -->
    <div v-if="localConfig.triggerType === 'webhook'" class="space-y-3">
      <div class="bg-purple-50 border-l-4 border-purple-400 p-3">
        <p class="text-sm text-purple-700">
          <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
          </svg>
          通过 HTTP 请求触发工作流执行，适用于接收外部系统的回调通知
        </p>
      </div>

      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">
          请求方法 <span class="text-red-500">*</span>
        </label>
        <BaseSelect
          v-model="localConfig.webhookMethod"
          :options="methodOptions"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">
          Webhook 路径 <span class="text-red-500">*</span>
        </label>
        <div class="flex items-center gap-2">
          <span class="text-sm text-slate-500 whitespace-nowrap">/api/webhook/</span>
          <BaseInput
            v-model="localConfig.webhookPath"
            placeholder="my-workflow"
            @update:model-value="emitUpdate"
          />
        </div>
        <p class="text-xs text-slate-500 mt-1">
          完整URL: <code class="px-1 py-0.5 bg-slate-100 rounded">https://your-domain.com/api/webhook/{{ localConfig.webhookPath || 'my-workflow' }}</code>
        </p>
      </div>

      <!-- 示例说明 -->
      <div class="bg-slate-50 rounded-lg p-3">
        <div class="text-xs font-semibold text-slate-700 mb-2">使用示例：</div>
        <div class="bg-slate-900 text-slate-100 rounded p-2 font-mono text-xs overflow-x-auto">
          <div>curl -X {{ localConfig.webhookMethod || 'POST' }} \</div>
          <div class="ml-2">https://your-domain.com/api/webhook/{{ localConfig.webhookPath || 'my-workflow' }} \</div>
          <div class="ml-2">-H "Content-Type: application/json" \</div>
          <div class="ml-2">-d '{"key": "value"}'</div>
        </div>
        <p class="text-xs text-slate-600 mt-2">
          请求体数据可通过 <code class="px-1 py-0.5 bg-slate-200 rounded">&#123;&#123;trigger.data&#125;&#125;</code> 在后续节点中访问
        </p>
      </div>
    </div>

    <!-- 手动触发配置 -->
    <div v-if="localConfig.triggerType === 'manual'">
      <div class="bg-green-50 border-l-4 border-green-400 p-3">
        <p class="text-sm text-green-700">
          <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
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

const localConfig = ref({
  triggerType: 'schedule',
  scheduleType: 'daily',
  scheduleValue: '09:00',
  webhookMethod: 'POST',
  webhookPath: '',
  ...props.config
})

const triggerTypeOptions = [
  { label: '定时触发', value: 'schedule' },
  { label: 'Webhook 触发', value: 'webhook' },
  { label: '手动触发', value: 'manual' }
]

const methodOptions = [
  { label: 'POST', value: 'POST' },
  { label: 'GET', value: 'GET' },
  { label: 'PUT', value: 'PUT' }
]

watch(() => props.config, (newVal) => {
  localConfig.value = { ...localConfig.value, ...newVal }
}, { deep: true })

const handleTriggerTypeChange = () => {
  // 重置相关配置
  if (localConfig.value.triggerType === 'schedule') {
    localConfig.value.scheduleType = 'daily'
    localConfig.value.scheduleValue = '09:00'
  } else if (localConfig.value.triggerType === 'webhook') {
    localConfig.value.webhookMethod = 'POST'
    localConfig.value.webhookPath = ''
  }
  emitUpdate()
}

const updateScheduleType = (type: string) => {
  localConfig.value.scheduleType = type
  // 根据类型设置默认值
  if (type === 'daily') {
    localConfig.value.scheduleValue = '09:00'
  } else if (type === 'hourly') {
    localConfig.value.scheduleValue = '00:00'
  } else if (type === 'interval') {
    localConfig.value.scheduleValue = '300'
  } else if (type === 'weekly') {
    localConfig.value.scheduleValue = '1'
  } else if (type === 'monthly') {
    localConfig.value.scheduleValue = '1'
  } else if (type === 'cron') {
    localConfig.value.scheduleValue = '0 0 * * * *'
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
