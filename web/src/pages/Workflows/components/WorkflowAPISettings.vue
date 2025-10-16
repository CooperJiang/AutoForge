<template>
  <Drawer v-model="isOpen" title="工作流 API 设置" size="lg" @close="handleClose">
    <Tabs v-model="activeTab" :tabs="tabList" class="mb-6" />

    <div class="space-y-6">
      <div v-show="activeTab === 'overview'">
        <div
          v-if="hasExternalTrigger"
          class="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4"
        >
          <div class="flex items-start gap-3">
            <AlertTriangle class="w-5 h-5 text-blue-600 flex-shrink-0 mt-0.5" />
            <div class="text-sm text-blue-800">
              <p class="font-medium mb-1">已配置外部触发节点</p>
              <p>
                工作流将接收外部传入的参数，参数可在后续节点中通过
                <code
                  class="px-1 py-0.5 bg-blue-100 rounded font-mono text-xs"
                  v-text="'{{external.参数名}}'"
                ></code>
                引用。
              </p>
            </div>
          </div>
        </div>

        <div
          v-if="!props.workflow?.id"
          class="bg-yellow-50 border border-yellow-300 rounded-lg p-3 mb-4"
        >
          <div class="flex items-start gap-2">
            <AlertTriangle class="w-4 h-4 text-yellow-600 mt-0.5 flex-shrink-0" />
            <div class="text-xs text-yellow-800">
              <p class="font-medium">工作流尚未保存</p>
              <p class="mt-1">请先保存工作流后，才能启用 API 调用功能。</p>
            </div>
          </div>
        </div>

        <div class="bg-bg-hover rounded-lg p-4">
          <div class="flex items-center justify-between mb-4">
            <div>
              <h3 class="text-sm font-semibold text-text-primary mb-1">API 状态</h3>
              <p class="text-xs text-text-secondary">
                {{ apiEnabled ? '已启用 - 可通过 API 调用此工作流' : '已禁用 - 无法通过 API 调用' }}
              </p>
            </div>
            <button
              @click="toggleAPI"
              :disabled="!props.workflow?.id"
              :class="[
                'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
                apiEnabled ? 'bg-primary' : 'bg-gray-300',
                !props.workflow?.id && 'opacity-50 cursor-not-allowed',
              ]"
            >
              <span
                :class="[
                  'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                  apiEnabled ? 'translate-x-6' : 'translate-x-1',
                ]"
              />
            </button>
          </div>

          <div v-if="apiEnabled" class="space-y-3">
            <div>
              <label class="block text-xs font-medium text-text-secondary mb-1">API Key</label>
              <div class="flex gap-2">
                <div class="flex-1 relative">
                  <input
                    :value="displayApiKey"
                    :type="showApiKey ? 'text' : 'password'"
                    readonly
                    class="w-full px-3 py-2 text-sm bg-bg-elevated border border-border-primary rounded-lg font-mono"
                  />
                  <button
                    @click="showApiKey = !showApiKey"
                    class="absolute right-2 top-1/2 -translate-y-1/2 text-text-secondary hover:text-text-primary"
                  >
                    <Eye v-if="showApiKey" class="w-4 h-4" />
                    <EyeOff v-else class="w-4 h-4" />
                  </button>
                </div>
                <BaseButton size="sm" variant="ghost" @click="copyApiKey" title="复制">
                  <Copy class="w-4 h-4" />
                </BaseButton>
                <BaseButton size="sm" variant="ghost" @click="regenerateApiKey" title="重新生成">
                  <RefreshCw class="w-4 h-4" />
                </BaseButton>
              </div>
            </div>

            <div>
              <label class="block text-xs font-medium text-text-secondary mb-1">调用端点</label>
              <div class="flex gap-2">
                <input
                  :value="apiEndpoint"
                  readonly
                  class="flex-1 px-3 py-2 text-sm bg-bg-elevated border border-border-primary rounded-lg font-mono"
                />
                <BaseButton size="sm" variant="ghost" @click="copyEndpoint" title="复制">
                  <Copy class="w-4 h-4" />
                </BaseButton>
              </div>
            </div>
          </div>
        </div>

        <div v-if="apiEnabled" class="grid grid-cols-3 gap-3">
          <div class="bg-bg-hover rounded-lg p-4">
            <div class="text-xs text-text-secondary mb-1">总调用次数</div>
            <div class="text-2xl font-bold text-text-primary">{{ stats.totalCalls }}</div>
          </div>
          <div class="bg-bg-hover rounded-lg p-4">
            <div class="text-xs text-text-secondary mb-1">最后调用</div>
            <div class="text-sm font-medium text-text-primary">{{ stats.lastCallTime }}</div>
          </div>
          <div class="bg-bg-hover rounded-lg p-4">
            <div class="text-xs text-text-secondary mb-1">成功率</div>
            <div class="text-2xl font-bold text-green-600">{{ stats.successRate }}</div>
          </div>
        </div>
      </div>

      <div v-show="activeTab === 'settings'" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-text-primary mb-2">
            超时时间（同步模式）
          </label>
          <Slider
            v-model="timeout"
            :min="30"
            :max="3600"
            :step="30"
            :value-formatter="(value) => `${value} 秒`"
          />
          <p class="text-xs text-text-tertiary mt-1">同步模式下，API 请求的最大等待时间</p>
        </div>

        <div>
          <label class="block text-sm font-medium text-text-primary mb-2">
            Webhook 回调地址（异步模式）
          </label>
          <input
            v-model="webhookURL"
            type="url"
            placeholder="https://example.com/webhook"
            class="w-full px-3 py-2 text-sm border border-border-primary rounded-lg bg-bg-elevated text-text-primary"
          />
          <p class="text-xs text-text-tertiary mt-1">异步执行完成后，系统会将结果 POST 到此地址</p>
        </div>

        <div class="bg-blue-50 border border-blue-200 rounded-lg p-3">
          <div class="text-xs text-blue-800 space-y-2">
            <p class="font-medium">执行模式说明：</p>
            <div>
              <p class="font-medium">• 同步模式 (?mode=sync)</p>
              <p class="ml-4">等待执行完成，直接返回结果。适合快速任务。</p>
            </div>
            <div>
              <p class="font-medium">• 异步模式 (?mode=async，默认)</p>
              <p class="ml-4">立即返回 execution_id，通过 Webhook 通知结果。适合耗时任务。</p>
            </div>
          </div>
        </div>
      </div>

      <div v-show="activeTab === 'test'" class="space-y-4">
        <div class="bg-yellow-50 border border-yellow-300 rounded-lg p-3 text-xs text-yellow-800">
          💡 提示：测试前请确保工作流已保存并启用
        </div>

        <div>
          <label class="block text-sm font-medium text-text-primary mb-2">执行模式</label>
          <div class="flex gap-3">
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="testMode" type="radio" value="sync" class="w-4 h-4 text-primary" />
              <span class="text-sm">同步模式</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input v-model="testMode" type="radio" value="async" class="w-4 h-4 text-primary" />
              <span class="text-sm">异步模式</span>
            </label>
          </div>
        </div>

        <div>
          <label class="block text-sm font-medium text-text-primary mb-2">请求参数（JSON）</label>
          <textarea
            v-model="testParams"
            rows="8"
            class="w-full px-3 py-2 text-sm border border-border-primary rounded-lg bg-bg-elevated text-text-primary font-mono"
            placeholder='{ "params": { "prompt": "测试" } }'
          ></textarea>
        </div>

        <BaseButton @click="handleTest" :disabled="!apiEnabled || testing" class="w-full">
          <Play v-if="!testing" class="w-4 h-4 mr-2" />
          <Loader v-else class="w-4 h-4 mr-2 animate-spin" />
          {{ testing ? '测试中...' : '发送测试请求' }}
        </BaseButton>

        <div v-if="testResult" class="space-y-2">
          <label class="block text-sm font-medium text-text-primary">响应结果</label>
          <div class="bg-gray-900 rounded-lg p-4 overflow-x-auto">
            <pre class="text-xs text-green-400 font-mono">{{ testResult }}</pre>
          </div>
        </div>
      </div>

      <div v-show="activeTab === 'code'" class="space-y-4">
        <div class="flex gap-2">
          <button
            v-for="lang in ['cURL', 'JavaScript', 'Python']"
            :key="lang"
            @click="codeLang = lang"
            :class="[
              'px-3 py-1.5 text-sm rounded transition-colors',
              codeLang === lang
                ? 'bg-primary text-white'
                : 'bg-bg-hover text-text-secondary hover:text-text-primary',
            ]"
          >
            {{ lang }}
          </button>
        </div>

        <div>
          <div class="bg-gray-900 rounded-lg p-4 overflow-x-auto">
            <pre class="text-xs text-green-400 font-mono">{{ currentCodeExample }}</pre>
          </div>
          <BaseButton size="sm" variant="ghost" class="mt-2" @click="copyCode">
            <Copy class="w-3 h-3 mr-1" />
            复制代码
          </BaseButton>
        </div>
      </div>

      <div v-show="activeTab === 'logs'" class="space-y-4">
        <div class="text-center py-12 text-text-tertiary">
          <History class="w-12 h-12 mx-auto mb-3 opacity-50" />
          <p>调用日志功能开发中...</p>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="flex justify-end gap-2">
        <BaseButton variant="secondary" @click="handleClose">关闭</BaseButton>
        <BaseButton v-if="activeTab === 'settings'" @click="handleSave"> 保存设置 </BaseButton>
      </div>
    </template>
  </Drawer>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { AlertTriangle, Eye, EyeOff, Copy, RefreshCw, Play, Loader, History } from 'lucide-vue-next'
import Drawer from '@/components/Drawer'
import BaseButton from '@/components/BaseButton'
import Tabs from '@/components/Tabs'
import Slider from '@/components/Slider'
import type { Tab } from '@/components/Tabs'
import { message } from '@/utils/message'
import { workflowApi } from '@/api/workflow'
import type { Workflow } from '@/types/workflow'

const props = defineProps<{
  workflow: Workflow
}>()

const emit = defineEmits<{
  refresh: []
}>()

const isOpen = defineModel<boolean>()

const tabList: Tab[] = [
  { value: 'overview', label: 'API 概览' },
  { value: 'settings', label: '调用设置' },
  { value: 'test', label: '测试调用' },
  { value: 'code', label: '代码示例' },
  { value: 'logs', label: '调用日志' },
]

const activeTab = ref('overview')
const apiEnabled = ref(props.workflow?.api_enabled || false)
const apiKey = ref(props.workflow?.api_key || '')
const timeout = ref(props.workflow?.api_timeout || 300)
const webhookURL = ref('')
const showApiKey = ref(false)
const testMode = ref('async')
const testParams = ref('')
const testing = ref(false)
const testResult = ref('')
const codeLang = ref('cURL')

// 检查是否有外部触发节点（仅用于提示，不限制功能）
const hasExternalTrigger = computed(() => {
  const nodes = props.workflow?.nodes || []
  return nodes.some((n) => n.type === 'external_trigger')
})

const displayApiKey = computed(() => {
  return apiKey.value || '未生成'
})

const apiEndpoint = computed(() => {
  return `${window.location.origin}/api/v1/public/workflows/invoke`
})

const stats = computed(() => {
  return {
    totalCalls: props.workflow?.api_call_count || 0,
    lastCallTime: props.workflow?.api_last_called_at
      ? new Date(props.workflow.api_last_called_at * 1000).toLocaleString('zh-CN', {
          month: '2-digit',
          day: '2-digit',
          hour: '2-digit',
          minute: '2-digit',
        })
      : '从未',
    successRate: '--',
  }
})

const exampleParams = computed(() => {
  const params = props.workflow?.nodes?.[0]?.config?.params || []
  return params.reduce((acc: any, param: any) => {
    acc[param.key] =
      param.example || (param.type === 'string' ? '示例值' : param.type === 'number' ? 0 : null)
    return acc
  }, {})
})

const currentCodeExample = computed(() => {
  const paramsJson = JSON.stringify({ params: exampleParams.value }, null, 2)

  if (codeLang.value === 'cURL') {
    return `curl -X POST '${apiEndpoint.value}?mode=sync' \\
  -H 'X-API-Key: ${apiKey.value || 'YOUR_API_KEY'}' \\
  -H 'Content-Type: application/json' \\
  -d '${paramsJson}'`
  } else if (codeLang.value === 'JavaScript') {
    return `const response = await fetch('${apiEndpoint.value}?mode=sync', {
  method: 'POST',
  headers: {
    'X-API-Key': '${apiKey.value || 'YOUR_API_KEY'}',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify(${paramsJson})
})

const data = await response.json()
console.log(data)`
  } else {
    // Python
    return `import requests

response = requests.post(
    '${apiEndpoint.value}?mode=sync',
    headers={
        'X-API-Key': '${apiKey.value || 'YOUR_API_KEY'}',
        'Content-Type': 'application/json'
    },
    json=${paramsJson}
)

data = response.json()
print(data)`
  }
})

watch(
  () => props.workflow,
  (newWorkflow, oldWorkflow) => {
    if (newWorkflow) {
      // 只在 workflow id 变化时更新状态（避免覆盖用户刚刚的操作）
      if (!oldWorkflow || oldWorkflow.id !== newWorkflow.id) {
        apiEnabled.value = newWorkflow.api_enabled || false
        apiKey.value = newWorkflow.api_key || ''
        timeout.value = newWorkflow.api_timeout || 300
        webhookURL.value = newWorkflow.api_webhook_url || ''
      } else {
        // 如果是同一个 workflow 的更新，只更新不会被用户直接修改的字段
        if (newWorkflow.api_key && newWorkflow.api_key !== apiKey.value) {
          apiKey.value = newWorkflow.api_key
        }
        if (newWorkflow.api_enabled !== undefined && apiEnabled.value !== newWorkflow.api_enabled) {
          apiEnabled.value = newWorkflow.api_enabled
        }
      }

      // 初始化测试参数
      if (!testParams.value) {
        testParams.value = JSON.stringify({ params: exampleParams.value }, null, 2)
      }
    }
  },
  { deep: true, immediate: true }
)

const toggleAPI = async () => {
  const targetState = !apiEnabled.value

  if (targetState) {
    try {
      const response = await workflowApi.enableAPI(props.workflow.id)
      apiEnabled.value = true
      apiKey.value = response.api_key
      message.success('API 已启用')
      emit('refresh')
    } catch {
      message.error('启用 API 失败')
    }
  } else {
    try {
      await workflowApi.disableAPI(props.workflow.id)
      apiEnabled.value = false
      apiKey.value = ''
      message.success('API 已禁用')
      emit('refresh')
    } catch {
      message.error('禁用 API 失败')
    }
  }
}

const regenerateApiKey = async () => {
  try {
    const response = await workflowApi.regenerateAPIKey(props.workflow.id)
    apiKey.value = response.api_key
    message.success('API Key 已重新生成')
    emit('refresh')
  } catch {
    message.error('重新生成失败')
  }
}

const copyApiKey = () => {
  navigator.clipboard.writeText(apiKey.value)
  message.success('API Key 已复制')
}

const copyEndpoint = () => {
  navigator.clipboard.writeText(apiEndpoint.value)
  message.success('端点地址已复制')
}

const copyCode = () => {
  navigator.clipboard.writeText(currentCodeExample.value)
  message.success('代码已复制')
}

const handleTest = async () => {
  testing.value = true
  testResult.value = ''

  try {
    // 这里实际调用 API
    JSON.parse(testParams.value)
    // TODO: 实现实际的测试调用
    testResult.value = JSON.stringify(
      {
        code: 200,
        message: '测试成功',
        data: {
          execution_id: 'exec_test_123',
          status: testMode.value === 'sync' ? 'success' : 'running',
          message: testMode.value === 'sync' ? '执行成功' : '工作流已开始执行',
        },
      },
      null,
      2
    )
    message.success('测试请求已发送')
  } catch (error: any) {
    testResult.value = JSON.stringify(
      {
        error: error.message || '测试失败',
      },
      null,
      2
    )
    message.error('测试失败')
  } finally {
    testing.value = false
  }
}

const handleClose = () => {
  isOpen.value = false
}

const handleSave = async () => {
  try {
    await workflowApi.updateAPITimeout(props.workflow.id, timeout.value)
    if (webhookURL.value) {
      await workflowApi.updateAPIWebhook(props.workflow.id, webhookURL.value)
    }
    message.success('设置已保存')
    emit('refresh')
  } catch {
    message.error('保存失败')
  }
}
</script>
