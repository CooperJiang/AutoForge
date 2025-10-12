<template>
  <Drawer v-model="isOpen" title="节点配置" size="xl" @close="handleClose">
    <div v-if="node" class="space-y-4">
      <!-- 节点基本信息 -->
      <div class="bg-slate-50 rounded-lg p-4">
        <h3 class="text-sm font-semibold text-slate-900 mb-3">基本信息</h3>
        <div class="space-y-3">
          <div>
            <label class="block text-xs font-medium text-slate-700 mb-1">节点名称</label>
            <BaseInput
              v-model="localNode.name"
              placeholder="输入节点名称"
            />
          </div>
        </div>
      </div>

      <!-- 工具配置 - 使用Tasks的配置组件 -->
      <div v-if="node.type === 'tool' && node.toolCode" class="border-t border-slate-200 pt-4">
        <h3 class="text-sm font-semibold text-slate-900 mb-3">工具配置</h3>

        <!-- HTTP请求 -->
        <div v-if="node.toolCode === 'http_request'" class="space-y-4">
          <!-- Curl 粘贴提示 -->
          <div class="bg-blue-50 border border-blue-200 rounded-lg p-3 text-xs text-blue-800">
            💡 小提示：按 <kbd class="px-1.5 py-0.5 bg-white border border-blue-300 rounded">{{ isMac ? 'Cmd' : 'Ctrl' }}</kbd> + <kbd class="px-1.5 py-0.5 bg-white border border-blue-300 rounded">V</kbd> 可直接粘贴 cURL 命令自动解析
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">
              请求方式 <span class="text-red-500">*</span>
            </label>
            <BaseSelect
              v-model="localNode.config.method"
              :options="methodOptions"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">
              接口地址 <span class="text-red-500">*</span>
            </label>
            <VariableSelector
              v-model="localNode.config.url"
              placeholder="https://api.example.com/checkin 或使用 {{变量}}"
              :previous-nodes="props.previousNodes"
              :env-vars="formattedEnvVars"
              :show-trigger-data="true"
            />
          </div>

          <!-- Headers -->
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">
              请求头（可选）
            </label>
            <div class="space-y-2">
              <ParamInput
                v-for="(header, index) in localNode.config.headers"
                :key="index"
                :param="header"
                key-placeholder="Header名称"
                value-placeholder="Header值"
                @update:param="updateHeader(index, $event)"
                @remove="removeHeader(index)"
              />
              <button
                type="button"
                @click="addHeader"
                class="w-full py-2 text-sm text-slate-600 border-2 border-dashed border-slate-300 rounded-lg hover:border-slate-400 hover:text-slate-700 transition-colors"
              >
                + 添加请求头
              </button>
            </div>
          </div>

          <!-- Params -->
          <div>
            <label class="block text-sm font-medium text-slate-700 mb-2">
              请求参数（可选）
            </label>
            <div class="space-y-2">
              <ParamInput
                v-for="(param, index) in localNode.config.params"
                :key="index"
                :param="param"
                key-placeholder="参数名"
                value-placeholder="参数值"
                @update:param="updateParam(index, $event)"
                @remove="removeParam(index)"
              />
              <button
                type="button"
                @click="addParam"
                class="w-full py-2 text-sm text-slate-600 border-2 border-dashed border-slate-300 rounded-lg hover:border-slate-400 hover:text-slate-700 transition-colors"
              >
                + 添加参数
              </button>
            </div>
          </div>

          <!-- Body -->
          <div>
            <button
              type="button"
              @click="bodyExpanded = !bodyExpanded"
              class="flex items-center justify-between w-full mb-2 text-left"
            >
              <label class="block text-sm font-medium text-slate-700 cursor-pointer">
                {{ bodyExpanded ? '▼' : '▶' }} 请求体 (Body) <span class="text-xs text-slate-500">(POST/PUT/PATCH)</span>
              </label>
            </button>
            <div v-show="bodyExpanded" class="space-y-1">
              <textarea
                v-model="localNode.config.body"
                class="w-full px-3 py-2 border-2 border-slate-200 rounded-lg focus:outline-none focus:border-green-500 font-mono text-sm"
                rows="8"
                placeholder='{"key": "value"}'
              />
              <div class="text-xs text-slate-500">支持 JSON、文本等格式</div>
            </div>
          </div>
        </div>

        <!-- 邮件发送 -->
        <EmailToolConfig
          v-else-if="node.toolCode === 'email_sender'"
          v-model:config="localNode.config"
        />

        <!-- 健康检查 -->
        <HealthCheckerConfig
          v-else-if="node.toolCode === 'health_checker'"
          v-model:config="localNode.config"
        />
      </div>

      <!-- 触发器配置 -->
      <div v-if="node.type === 'trigger'" class="border-t border-slate-200 pt-4">
        <h3 class="text-sm font-semibold text-slate-900 mb-3">触发配置</h3>
        <TriggerConfig v-model:config="localNode.config" />
      </div>

      <!-- 条件配置 -->
      <div v-if="node.type === 'condition'" class="border-t border-slate-200 pt-4">
        <h3 class="text-sm font-semibold text-slate-900 mb-3">条件配置</h3>
        <ConditionConfig v-model:config="localNode.config" />
      </div>

      <!-- 延迟配置 -->
      <div v-if="node.type === 'delay'" class="border-t border-slate-200 pt-4">
        <h3 class="text-sm font-semibold text-slate-900 mb-3">延迟配置</h3>
        <DelayConfig v-model:config="localNode.config" />
      </div>

      <!-- 开关配置 -->
      <div v-if="node.type === 'switch'" class="border-t border-slate-200 pt-4">
        <h3 class="text-sm font-semibold text-slate-900 mb-3">开关配置</h3>
        <SwitchConfig v-model:config="localNode.config" />
      </div>

      <!-- 错误重试配置 -->
      <div v-if="node.type === 'tool'" class="border-t border-slate-200 pt-4">
        <h3 class="text-sm font-semibold text-slate-900 mb-3">错误重试</h3>
        <RetryConfig
          :config="localNode.retry || defaultRetryConfig"
          @update:config="updateRetryConfig"
        />
      </div>

      <!-- 测试运行结果 -->
      <div v-if="testResult" class="border-t border-slate-200 pt-4">
        <div class="bg-slate-50 rounded-lg p-4">
          <div class="flex items-center justify-between mb-3">
            <h3 class="text-sm font-semibold text-slate-900">测试结果</h3>
            <span
              :class="[
                'px-2 py-1 rounded text-xs font-medium',
                testResult.success ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700'
              ]"
            >
              {{ testResult.success ? '✓ 成功' : '✗ 失败' }}
            </span>
          </div>

          <div v-if="testResult.error" class="mb-3 p-3 bg-red-50 border-l-4 border-red-400 text-sm text-red-700">
            <div class="font-semibold mb-1">错误信息：</div>
            <div class="font-mono text-xs">{{ testResult.error }}</div>
          </div>

          <div v-if="testResult.output" class="space-y-2">
            <div class="text-xs font-semibold text-slate-700">输出数据结构：</div>
            <div class="bg-slate-900 text-slate-100 rounded p-3 font-mono text-xs overflow-x-auto">
              <pre>{{ JSON.stringify(testResult.output, null, 2) }}</pre>
            </div>
            <div class="text-xs text-slate-600">
              💡 可以在后续节点中通过 <code class="px-1 py-0.5 bg-slate-200 rounded">&#123;&#123;{{ node.id }}.fieldName&#125;&#125;</code> 引用这些字段
            </div>
          </div>
        </div>
      </div>

      <!-- 操作按钮 -->
      <div class="flex items-center justify-between pt-4 border-t border-slate-200">
        <BaseButton
          size="sm"
          variant="danger"
          @click="handleDelete"
        >
          <Trash2 class="w-4 h-4 mr-1" />
          删除节点
        </BaseButton>
        <div class="flex gap-2">
          <BaseButton
            v-if="node.type === 'tool'"
            size="sm"
            variant="secondary"
            @click="handleTestRun"
            :disabled="testRunning"
          >
            <Play class="w-4 h-4 mr-1" />
            {{ testRunning ? '测试中...' : '测试运行' }}
          </BaseButton>
          <BaseButton
            size="sm"
            variant="ghost"
            @click="handleClose"
          >
            取消
          </BaseButton>
          <BaseButton
            size="sm"
            @click="handleSave"
          >
            保存配置
          </BaseButton>
        </div>
      </div>
    </div>
  </Drawer>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
import { Trash2, GitBranch, Play } from 'lucide-vue-next'
import Drawer from '@/components/Drawer'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import ParamInput from '@/components/ParamInput'
import VariableSelector from '@/components/VariableSelector'
import EmailToolConfig from '@/components/tools/EmailToolConfig.vue'
import HealthCheckerConfig from '@/components/tools/HealthCheckerConfig.vue'
import TriggerConfig from '@/components/tools/TriggerConfig.vue'
import ConditionConfig from '@/components/tools/ConditionConfig.vue'
import DelayConfig from '@/components/tools/DelayConfig.vue'
import SwitchConfig from '@/components/tools/SwitchConfig.vue'
import RetryConfig from '@/components/RetryConfig'
import type { WorkflowNode, WorkflowEnvVar, NodeRetryConfig } from '@/types/workflow'
import { message } from '@/utils/message'
import { parseCurl } from '@/utils/curlParser'

interface Props {
  modelValue: boolean
  node: WorkflowNode | null
  previousNodes?: WorkflowNode[]
  envVars?: WorkflowEnvVar[]
}

const props = withDefaults(defineProps<Props>(), {
  previousNodes: () => [],
  envVars: () => []
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  update: [nodeId: string, updates: Partial<WorkflowNode>]
  delete: [nodeId: string]
}>()

const isOpen = ref(props.modelValue)
const isMac = ref(/Mac/.test(navigator.userAgent))
const bodyExpanded = ref(false)
const testRunning = ref(false)
const testResult = ref<{ success: boolean; output?: any; error?: string } | null>(null)

const localNode = ref<WorkflowNode>({
  id: '',
  type: 'tool',
  name: '',
  config: {},
  position: { x: 0, y: 0 }
})

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' }
]

// 默认重试配置
const defaultRetryConfig: NodeRetryConfig = {
  enabled: false,
  maxRetries: 3,
  retryInterval: 5,
  exponentialBackoff: false
}

// 格式化环境变量供 VariableSelector 使用
const formattedEnvVars = computed(() => {
  return props.envVars.map(v => ({
    key: v.key,
    description: v.description || v.key
  }))
})

// 更新重试配置
const updateRetryConfig = (config: NodeRetryConfig) => {
  localNode.value.retry = config
}

watch(() => props.modelValue, (val) => {
  isOpen.value = val
  if (val && props.node) {
    localNode.value = JSON.parse(JSON.stringify(props.node))

    // 初始化HTTP配置
    if (localNode.value.type === 'tool' && localNode.value.toolCode === 'http_request') {
      if (!localNode.value.config.method) localNode.value.config.method = 'GET'
      if (!localNode.value.config.url) localNode.value.config.url = ''
      if (!localNode.value.config.headers) localNode.value.config.headers = []
      if (!localNode.value.config.params) localNode.value.config.params = []
      if (!localNode.value.config.body) localNode.value.config.body = ''
    }
    // 其他工具默认配置
    else if (localNode.value.type === 'tool') {
      if (localNode.value.toolCode === 'health_checker') {
        localNode.value.config.timeout = localNode.value.config.timeout || 30
        localNode.value.config.expectedStatus = localNode.value.config.expectedStatus || 200
      }
    }
    // 触发器默认配置
    else if (localNode.value.type === 'trigger') {
      localNode.value.config.scheduleType = localNode.value.config.scheduleType || 'daily'
      localNode.value.config.time = localNode.value.config.time || '09:00'
      localNode.value.config.enabled = localNode.value.config.enabled !== undefined ? localNode.value.config.enabled : true
    }
  }
})

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

// Headers 操作
const addHeader = () => {
  if (!Array.isArray(localNode.value.config.headers)) {
    localNode.value.config.headers = []
  }
  localNode.value.config.headers.push({ key: '', value: '' })
}

const updateHeader = (index: number, param: any) => {
  localNode.value.config.headers[index] = param
}

const removeHeader = (index: number) => {
  localNode.value.config.headers.splice(index, 1)
}

// Params 操作
const addParam = () => {
  if (!Array.isArray(localNode.value.config.params)) {
    localNode.value.config.params = []
  }
  localNode.value.config.params.push({ key: '', value: '' })
}

const updateParam = (index: number, param: any) => {
  localNode.value.config.params[index] = param
}

const removeParam = (index: number) => {
  localNode.value.config.params.splice(index, 1)
}

// cURL 粘贴解析
const handlePaste = (e: ClipboardEvent) => {
  if (!isOpen.value || !props.node) return
  if (localNode.value.type !== 'tool' || localNode.value.toolCode !== 'http_request') return

  const text = e.clipboardData?.getData('text')
  if (!text || !text.trim().startsWith('curl')) return

  e.preventDefault()

  const parsed = parseCurl(text)
  if (parsed) {
    // 格式化 body
    let formattedBody = parsed.body || ''
    if (formattedBody) {
      try {
        const bodyObj = JSON.parse(formattedBody)
        formattedBody = JSON.stringify(bodyObj, null, 2)
      } catch {
        // 如果不是 JSON，保持原样
      }
    }

    localNode.value.config = {
      url: parsed.url,
      method: parsed.method,
      headers: parsed.headers,
      params: parsed.params,
      body: formattedBody
    }

    // 如果有 body，自动展开
    if (formattedBody) {
      bodyExpanded.value = true
    }

    message.success('cURL 命令解析成功')
  } else {
    message.error('cURL 命令解析失败')
  }
}

const handleClose = () => {
  isOpen.value = false
}

const handleSave = () => {
  if (props.node) {
    emit('update', props.node.id, {
      name: localNode.value.name,
      config: localNode.value.config,
      retry: localNode.value.retry
    })
    handleClose()
  }
}

const handleDelete = () => {
  if (props.node && confirm('确定要删除此节点吗？')) {
    emit('delete', props.node.id)
  }
}

// 测试运行节点
const handleTestRun = async () => {
  if (!props.node) return

  testRunning.value = true
  testResult.value = null

  try {
    // 验证节点配置
    if (localNode.value.toolCode === 'http_request') {
      if (!localNode.value.config.url) {
        testResult.value = {
          success: false,
          error: '请填写接口地址'
        }
        return
      }
    }

    message.info('正在测试运行节点...')

    // TODO: 调用API执行节点测试
    // const response = await nodeApi.testRun(localNode.value)

    // Mock 测试结果
    await new Promise(resolve => setTimeout(resolve, 1500))

    if (localNode.value.toolCode === 'http_request') {
      testResult.value = {
        success: true,
        output: {
          status: 200,
          statusText: 'OK',
          headers: {
            'content-type': 'application/json'
          },
          body: {
            success: true,
            data: {
              id: 12345,
              message: '请求成功'
            }
          },
          responseTime: 234
        }
      }
    } else if (localNode.value.toolCode === 'health_checker') {
      testResult.value = {
        success: true,
        output: {
          healthy: true,
          status: 200,
          responseTime: 156
        }
      }
    } else if (localNode.value.toolCode === 'email') {
      testResult.value = {
        success: true,
        output: {
          sent: true,
          messageId: 'msg-12345'
        }
      }
    } else {
      testResult.value = {
        success: true,
        output: {
          result: '执行成功'
        }
      }
    }

    message.success('节点测试运行成功')
  } catch (error: any) {
    testResult.value = {
      success: false,
      error: error.message || '测试运行失败'
    }
    message.error('节点测试运行失败')
  } finally {
    testRunning.value = false
  }
}

// 监听粘贴事件
onMounted(() => {
  window.addEventListener('paste', handlePaste)
})

onUnmounted(() => {
  window.removeEventListener('paste', handlePaste)
})
</script>
