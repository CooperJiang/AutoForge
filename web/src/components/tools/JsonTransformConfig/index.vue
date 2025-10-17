<template>
  <div class="space-y-4">
    <div class="bg-primary-light border-l-4 border-border-focus p-3">
      <p class="text-sm text-primary">
        选择数据来源（支持变量，例如
        <code v-pre class="px-1 py-0.5 bg-bg-tertiary rounded">{{ nodes.node_xxx.output }}</code
        >），然后编写 JS 表达式处理数据。 表达式中可直接使用
        <code class="px-1 py-0.5 bg-bg-tertiary rounded">data</code> 与
        <code class="px-1 py-0.5 bg-bg-tertiary rounded">ctx</code>。
      </p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        数据来源 <span class="text-red-500">*</span>
      </label>
      <VariableSelector
        v-model="localConfig.data_source"
        placeholder="{{nodes.node_...}}"
        :previous-nodes="previousNodes"
        :env-vars="helperEnvVars"
      />
      <p class="text-xs text-text-tertiary mt-1">
        可直接引用前置节点输出或环境变量，无需手动拼接 JSON。
      </p>
    </div>

    <div>
      <div class="flex items-center justify-between mb-2">
        <label class="block text-sm font-medium text-text-secondary">
          JS 表达式 <span class="text-red-500">*</span>
        </label>
        <button
          type="button"
          class="text-xs text-primary hover:underline"
          @click="toggleCheatsheet"
        >
          {{ showCheatsheet ? '隐藏示例' : '查看示例' }}
        </button>
      </div>
      <textarea
        v-model="localConfig.expression"
        rows="6"
        class="w-full px-3 py-2 text-sm font-mono border-2 border-border-primary rounded-lg focus:outline-none focus:border-border-focus focus:ring-2 focus:ring-primary-light bg-bg-primary text-text-primary placeholder:text-text-placeholder"
        placeholder="data.map(item => item.url)"
      ></textarea>
      <div v-if="showCheatsheet" class="mt-2 text-xs text-text-tertiary space-y-1">
        <div>• 获取数组字段：<code>data.map(item => item.url)</code></div>
        <div>
          • 结合上下文：<code v-pre
            >ctx.nodes.node_xxx.output.filter(it => it.status === 'ok')</code
          >
        </div>
        <div>• 构造对象：<code>({ count: data.length, first: data[0] })</code></div>
      </div>
    </div>

    <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-1">输出字段名称</label>
        <BaseInput v-model="localConfig.output_name" placeholder="result" />
      </div>
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-1">执行超时 (ms)</label>
        <BaseInput
          v-model.number="localConfig.timeout_ms"
          type="number"
          min="100"
          placeholder="1500"
        />
      </div>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"
        >样例 JSON（可选，仅用于预览）</label
      >
      <textarea
        v-model="localConfig.sample_json"
        rows="4"
        class="w-full px-3 py-2 text-sm font-mono border-2 border-border-primary rounded-lg focus:outline-none focus:border-border-focus focus:ring-2 focus:ring-primary-light bg-bg-primary text-text-primary placeholder:text-text-placeholder"
        placeholder='[{"url":"https://example.com"}]'
      ></textarea>
      <p class="text-xs text-text-tertiary mt-1">提供样例数据后，可在前端预览表达式输出。</p>
    </div>

    <div class="space-y-2">
      <div class="flex items-center gap-2">
        <BaseButton size="sm" @click="runPreview" :loading="previewLoading"> 预览输出 </BaseButton>
        <span class="text-xs text-text-tertiary">仅基于上方样例 JSON 进行演算</span>
      </div>

      <div
        v-if="previewError"
        class="bg-red-50 border border-red-200 text-xs text-red-700 rounded-lg p-2"
      >
        {{ previewError }}
      </div>

      <div
        v-else-if="previewResult"
        class="border border-border-primary rounded-lg overflow-hidden"
      >
        <JsonViewer :content="previewResult" />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import VariableSelector from '@/components/VariableSelector'
import BaseInput from '@/components/BaseInput'
import BaseButton from '@/components/BaseButton'
import JsonViewer from '@/components/JsonViewer'
import { testTask } from '@/api/task'
import { message } from '@/utils/message'

interface Props {
  config: Record<string, any>
  previousNodes?: Array<{ id: string; name: string; type: string; toolCode?: string }>
  envVars?: Array<{ key: string; value: string; description?: string }>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:config': [config: Record<string, any>]
}>()

const defaultConfig = () => ({
  data_source: '',
  expression: '',
  output_name: 'result',
  timeout_ms: 1500,
  sample_json: '',
})

const localConfig = ref({ ...defaultConfig(), ...props.config })
const showCheatsheet = ref(false)

const previewLoading = ref(false)
const previewResult = ref('')
const previewError = ref('')

const helperEnvVars = computed(() => props.envVars || [])

const updatingFromProps = ref(false)

watch(
  () => props.config,
  (newVal) => {
    updatingFromProps.value = true
    localConfig.value = { ...defaultConfig(), ...(newVal || {}) }
    updatingFromProps.value = false
    previewResult.value = ''
    previewError.value = ''
  },
  { deep: true, immediate: true }
)

watch(
  localConfig,
  (val) => {
    if (updatingFromProps.value) return
    emit('update:config', { ...val })
  },
  { deep: true }
)

watch(
  () => [localConfig.value.expression, localConfig.value.sample_json],
  () => {
    previewResult.value = ''
    previewError.value = ''
  }
)

const runPreview = async () => {
  const sample = localConfig.value.sample_json?.trim()
  if (!sample) {
    message.warning('请先填写样例 JSON 后再预览')
    return
  }

  if (!localConfig.value.expression?.trim()) {
    message.warning('请先填写 JS 表达式')
    return
  }

  previewLoading.value = true
  previewResult.value = ''
  previewError.value = ''

  try {
    const response = await testTask({
      tool_code: 'json_transform',
      config: {
        data_source: sample,
        expression: localConfig.value.expression,
        output_name: localConfig.value.output_name,
        timeout_ms: localConfig.value.timeout_ms,
      },
    })

    if (response.success) {
      const outputData = response.output?.[localConfig.value.output_name] ?? response.output ?? null
      if (outputData !== null) {
        previewResult.value = JSON.stringify(outputData, null, 2)
      } else {
        previewResult.value = JSON.stringify(response, null, 2)
      }
      if (response.message) {
        previewError.value = ''
      }
    } else {
      previewError.value = response.error_message || response.message || '预览失败'
    }
  } catch (error: any) {
    previewError.value = error.response?.data?.message || error.message || '预览失败'
  } finally {
    previewLoading.value = false
  }
}

const toggleCheatsheet = () => {
  showCheatsheet.value = !showCheatsheet.value
}
</script>
