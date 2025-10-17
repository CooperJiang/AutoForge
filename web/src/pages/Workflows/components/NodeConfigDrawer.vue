<template>
  <Drawer v-model="isOpen" title="节点配置" size="xl" @close="handleClose">
    <div v-if="node" class="space-y-4">
      <div class="bg-bg-hover rounded-lg p-4">
        <h3 class="text-sm font-semibold text-text-primary mb-3">基本信息</h3>
        <div class="space-y-3">
          <div>
            <label class="block text-xs font-medium text-text-secondary mb-1">节点名称</label>
            <BaseInput v-model="localNode.name" placeholder="输入节点名称" />
          </div>
        </div>
      </div>

      <div v-if="node.type === 'tool' && node.toolCode" class="border-t border-border-primary pt-4">
        <h3 class="text-sm font-semibold text-text-primary mb-3">工具配置</h3>

        <div class="bg-bg-hover rounded-lg p-3 border border-border-primary mb-4">
          <div class="flex items-center justify-between mb-2">
            <h4 class="text-sm font-medium text-text-primary">变量助手</h4>
            <button
              type="button"
              @click="showVariableHelper = !showVariableHelper"
              class="text-xs text-primary hover:underline"
            >
              {{ showVariableHelper ? '隐藏' : '显示' }}
            </button>
          </div>

          <VariableHelper
            v-if="showVariableHelper"
            :show="true"
            :previous-nodes="props.previousNodes"
            :env-vars="props.envVars"
          />

          <p v-if="!showVariableHelper" class="text-xs text-text-tertiary">
            点击"显示"按钮查看可用的变量，点击变量即可复制到剪贴板
          </p>
        </div>

        <div v-if="node.toolCode === 'http_request'" class="space-y-4">
          <div class="bg-primary-light border border-primary rounded-lg p-3 text-xs text-primary">
            💡 小提示：按
            <kbd class="px-1.5 py-0.5 bg-bg-elevated border border-primary rounded">{{
              isMac ? 'Cmd' : 'Ctrl'
            }}</kbd>
            +
            <kbd class="px-1.5 py-0.5 bg-bg-elevated border border-primary rounded">V</kbd>
            可直接粘贴 cURL 命令自动解析
          </div>

          <div>
            <label class="block text-sm font-medium text-text-secondary mb-2">
              请求方式 <span class="text-red-500">*</span>
            </label>
            <BaseSelect v-model="localNode.config.method" :options="methodOptions" />
          </div>

          <div>
            <label class="block text-sm font-medium text-text-secondary mb-2">
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

          <div>
            <label class="block text-sm font-medium text-text-secondary mb-2">
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
                class="w-full py-2 text-sm text-text-secondary border-2 border-dashed border-slate-300 rounded-lg hover:border-slate-400 hover:text-text-secondary transition-colors"
              >
                + 添加请求头
              </button>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-text-secondary mb-2">
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
                class="w-full py-2 text-sm text-text-secondary border-2 border-dashed border-slate-300 rounded-lg hover:border-slate-400 hover:text-text-secondary transition-colors"
              >
                + 添加参数
              </button>
            </div>
          </div>

          <div>
            <button
              type="button"
              @click="bodyExpanded = !bodyExpanded"
              class="flex items-center justify-between w-full mb-2 text-left"
            >
              <label class="block text-sm font-medium text-text-secondary cursor-pointer">
                {{ bodyExpanded ? '▼' : '▶' }} 请求体 (Body)
                <span class="text-xs text-text-tertiary">(POST/PUT/PATCH)</span>
              </label>
            </button>
            <div v-show="bodyExpanded" class="space-y-1">
              <textarea
                v-model="localNode.config.body"
                class="w-full px-3 py-1.5 text-sm text-text-primary bg-bg-primary border-2 border-border-primary rounded-md transition-all duration-200 focus:border-border-focus focus:ring-2 focus:ring-primary-light focus:outline-none hover:border-border-secondary placeholder:text-text-placeholder font-mono"
                rows="8"
                placeholder='{"key": "value"}'
              />
              <div class="text-xs text-text-tertiary">支持 JSON、文本等格式</div>
            </div>
          </div>
        </div>

        <EmailToolConfig
          v-else-if="node.toolCode === 'email_sender'"
          v-model:config="localNode.config"
          :previous-nodes="props.previousNodes"
          :env-vars="props.envVars"
        />

        <HealthCheckerConfig
          v-else-if="node.toolCode === 'health_checker'"
          v-model:config="localNode.config"
        />

        <FeishuBotConfig
          v-else-if="node.toolCode === 'feishu_bot'"
          v-model:config="localNode.config"
          :previous-nodes="props.previousNodes"
          :env-vars="props.envVars"
        />

        <OpenAIConfig
          v-else-if="node.toolCode === 'openai_chatgpt'"
          v-model:config="localNode.config"
          :previous-nodes="props.previousNodes"
          :env-vars="props.envVars"
        />

        <OpenAIImageConfig
          v-else-if="node.toolCode === 'openai_image'"
          v-model:config="localNode.config"
          :previous-nodes="props.previousNodes"
          :env-vars="props.envVars"
        />

        <JsonTransformConfig
          v-else-if="node.toolCode === 'json_transform'"
          v-model:config="localNode.config"
          :previous-nodes="props.previousNodes"
          :env-vars="props.envVars"
        />

        <OutputFormatterConfig
          v-else-if="node.toolCode === 'output_formatter'"
          v-model:config="localNode.config"
          :previous-nodes="props.previousNodes"
          :env-vars="props.envVars"
        />

        <HtmlRenderConfig
          v-else-if="node.toolCode === 'html_render'"
          v-model:config="localNode.config"
          :previous-nodes="props.previousNodes"
          :env-vars="props.envVars"
        />

        <RedisContextConfig
          v-else-if="node.toolCode === 'redis_context'"
          v-model:config="localNode.config"
        />

        <ContextManagerConfig
          v-else-if="node.toolCode === 'context_manager'"
          v-model:config="localNode.config"
          :previous-nodes="props.previousNodes"
          :env-vars="props.envVars"
        />

        <PixelPunkUploadConfig
          v-else-if="node.toolCode === 'pixelpunk_upload'"
          v-model:config="localNode.config"
        />
        <FileDownloaderConfig
          v-else-if="node.toolCode === 'file_downloader'"
          v-model:config="localNode.config"
          :previous-nodes="props.previousNodes"
          :env-vars="props.envVars"
        />
        <AliyunOSSConfig
          v-else-if="node.toolCode === 'aliyun_oss'"
          v-model:config="localNode.config"
        />
        <TencentCOSConfig
          v-else-if="node.toolCode === 'tencent_cos'"
          v-model:config="localNode.config"
        />
        <QRCodeConfig
          v-else-if="node.toolCode === 'qrcode_generator'"
          v-model:config="localNode.config"
        />
        <GeminiConfig
          v-else-if="node.toolCode === 'gemini_chat'"
          v-model:config="localNode.config"
        />
      </div>

      <div v-if="node.type === 'external_trigger'" class="border-t border-border-primary pt-4">
        <h3 class="text-sm font-semibold text-text-primary mb-3">外部 API 触发配置</h3>
        <ExternalTriggerConfig v-model:config="localNode.config" @update="handleConfigUpdate" />
      </div>

      <div v-if="node.type === 'trigger'" class="border-t border-border-primary pt-4">
        <h3 class="text-sm font-semibold text-text-primary mb-3">触发配置</h3>
        <TriggerConfig v-model:config="localNode.config" />
      </div>

      <div v-if="node.type === 'condition'" class="border-t border-border-primary pt-4">
        <h3 class="text-sm font-semibold text-text-primary mb-3">条件配置</h3>

        <div class="bg-bg-hover rounded-lg p-3 border border-border-primary mb-4">
          <div class="flex items-center justify-between mb-2">
            <h4 class="text-sm font-medium text-text-primary">变量助手</h4>
            <button
              type="button"
              @click="showVariableHelper = !showVariableHelper"
              class="text-xs text-primary hover:underline"
            >
              {{ showVariableHelper ? '隐藏' : '显示' }}
            </button>
          </div>

          <VariableHelper
            v-if="showVariableHelper"
            :show="true"
            :previous-nodes="props.previousNodes"
            :env-vars="props.envVars"
          />

          <p v-if="!showVariableHelper" class="text-xs text-text-tertiary">
            点击"显示"按钮查看可用的变量，点击变量即可复制到剪贴板
          </p>
        </div>

        <ConditionConfig v-model:config="localNode.config" :previous-nodes="props.previousNodes" />
      </div>

      <div v-if="node.type === 'delay'" class="border-t border-border-primary pt-4">
        <h3 class="text-sm font-semibold text-text-primary mb-3">延迟配置</h3>

        <div class="bg-bg-hover rounded-lg p-3 border border-border-primary mb-4">
          <div class="flex items-center justify-between mb-2">
            <h4 class="text-sm font-medium text-text-primary">变量助手</h4>
            <button
              type="button"
              @click="showVariableHelper = !showVariableHelper"
              class="text-xs text-primary hover:underline"
            >
              {{ showVariableHelper ? '隐藏' : '显示' }}
            </button>
          </div>

          <VariableHelper
            v-if="showVariableHelper"
            :show="true"
            :previous-nodes="props.previousNodes"
            :env-vars="props.envVars"
          />

          <p v-if="!showVariableHelper" class="text-xs text-text-tertiary">
            点击"显示"按钮查看可用的变量，点击变量即可复制到剪贴板
          </p>
        </div>

        <DelayConfig
          v-model:config="localNode.config"
          :previous-nodes="props.previousNodes"
          :env-vars="props.envVars"
        />
      </div>

      <div v-if="node.type === 'switch'" class="border-t border-border-primary pt-4">
        <h3 class="text-sm font-semibold text-text-primary mb-3">开关配置</h3>

        <div class="bg-bg-hover rounded-lg p-3 border border-border-primary mb-4">
          <div class="flex items-center justify-between mb-2">
            <h4 class="text-sm font-medium text-text-primary">变量助手</h4>
            <button
              type="button"
              @click="showVariableHelper = !showVariableHelper"
              class="text-xs text-primary hover:underline"
            >
              {{ showVariableHelper ? '隐藏' : '显示' }}
            </button>
          </div>

          <VariableHelper
            v-if="showVariableHelper"
            :show="true"
            :previous-nodes="props.previousNodes"
            :env-vars="props.envVars"
          />

          <p v-if="!showVariableHelper" class="text-xs text-text-tertiary">
            点击"显示"按钮查看可用的变量，点击变量即可复制到剪贴板
          </p>
        </div>

        <SwitchConfig
          v-model:config="localNode.config"
          :previous-nodes="props.previousNodes"
          :env-vars="props.envVars"
        />
      </div>

      <div v-if="node.type === 'tool'" class="border-t border-border-primary pt-4">
        <h3 class="text-sm font-semibold text-text-primary mb-3">错误重试</h3>
        <RetryConfig
          :config="localNode.retry || defaultRetryConfig"
          @update:config="updateRetryConfig"
        />
      </div>

      <div v-if="testResult" class="border-t border-border-primary pt-4">
        <div class="bg-bg-hover rounded-lg p-4">
          <div class="flex items-center justify-between mb-3">
            <h3 class="text-sm font-semibold text-text-primary">测试结果</h3>
            <span
              :class="[
                'px-2 py-1 rounded text-xs font-medium',
                testResult.success ? 'bg-green-100 text-green-700' : 'bg-red-100 text-red-700',
              ]"
            >
              {{ testResult.success ? '✓ 成功' : '✗ 失败' }}
            </span>
          </div>

          <div
            v-if="testResult.error"
            class="mb-3 p-3 bg-red-50 border-l-4 border-red-400 text-sm text-red-700"
          >
            <div class="font-semibold mb-1">错误信息：</div>
            <div class="font-mono text-xs">{{ testResult.error }}</div>
          </div>

          <div v-if="testResult.output" class="space-y-2">
            <div class="text-xs font-semibold text-text-secondary">输出数据结构：</div>
            <div class="bg-slate-900 text-slate-100 rounded p-3 font-mono text-xs overflow-x-auto">
              <pre>{{ JSON.stringify(testResult.output, null, 2) }}</pre>
            </div>
            <div class="text-xs text-text-secondary">
              💡 可以在后续节点中通过
              <code class="px-1 py-0.5 bg-bg-tertiary rounded"
                >&#123;&#123;{{ node.id }}.fieldName&#125;&#125;</code
              >
              引用这些字段
            </div>
          </div>
        </div>
      </div>

      <div class="flex items-center justify-between pt-4 border-t border-border-primary">
        <BaseButton size="sm" variant="danger" @click="handleDelete">
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
          <BaseButton size="sm" variant="ghost" @click="handleClose"> 取消 </BaseButton>
          <BaseButton size="sm" @click="handleSave"> 保存配置 </BaseButton>
        </div>
      </div>
    </div>
  </Drawer>

  <Dialog
    v-model="showDeleteConfirm"
    title="删除节点"
    max-width="max-w-md"
    @confirm="confirmDelete"
    @cancel="showDeleteConfirm = false"
  >
    <div class="py-2">
      <p class="text-text-primary">
        确定要删除节点 <span class="font-semibold text-primary">{{ node?.name }}</span> 吗？
      </p>
      <p class="text-text-secondary text-sm mt-1">此操作无法撤销。</p>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted, computed } from 'vue'
import { Trash2, Play } from 'lucide-vue-next'
import Drawer from '@/components/Drawer'
import Dialog from '@/components/Dialog'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import ParamInput from '@/components/ParamInput'
import VariableSelector from '@/components/VariableSelector'
import EmailToolConfig from '@/components/tools/EmailToolConfig/index.vue'
import HealthCheckerConfig from '@/components/tools/HealthCheckerConfig/index.vue'
import FeishuBotConfig from '@/components/tools/FeishuBotConfig/index.vue'
import OpenAIConfig from '@/components/tools/OpenAIConfig/index.vue'
import OpenAIImageConfig from '@/components/tools/OpenAIImageConfig/index.vue'
import JsonTransformConfig from '@/components/tools/JsonTransformConfig/index.vue'
import OutputFormatterConfig from '@/components/tools/OutputFormatterConfig/index.vue'
import HtmlRenderConfig from '@/components/tools/HtmlRenderConfig/index.vue'
import RedisContextConfig from '@/components/tools/RedisContextConfig/index.vue'
import ContextManagerConfig from '@/components/tools/ContextManagerConfig/index.vue'
import PixelPunkUploadConfig from '@/components/tools/PixelPunkUploadConfig/index.vue'
import FileDownloaderConfig from '@/components/tools/FileDownloaderConfig/index.vue'
import TriggerConfig from '@/components/tools/TriggerConfig/index.vue'
import AliyunOSSConfig from '@/components/tools/AliyunOSSConfig/index.vue'
import TencentCOSConfig from '@/components/tools/TencentCOSConfig/index.vue'
import QRCodeConfig from '@/components/tools/QRCodeConfig/index.vue'
import GeminiConfig from '@/components/tools/GeminiConfig/index.vue'
import ConditionConfig from '@/components/tools/ConditionConfig/index.vue'
import DelayConfig from '@/components/tools/DelayConfig/index.vue'
import SwitchConfig from '@/components/tools/SwitchConfig/index.vue'
import ExternalTriggerConfig from './ExternalTriggerConfig.vue'
import RetryConfig from '@/components/RetryConfig'
import VariableHelper from '@/components/VariableHelper'
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
  envVars: () => [],
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
const showVariableHelper = ref(false)

const localNode = ref<WorkflowNode>({
  id: '',
  type: 'tool',
  name: '',
  config: {},
  position: { x: 0, y: 0 },
})

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
  { label: 'PATCH', value: 'PATCH' },
]

// 默认重试配置
const defaultRetryConfig: NodeRetryConfig = {
  enabled: false,
  maxRetries: 3,
  retryInterval: 5,
  exponentialBackoff: false,
}

// 格式化环境变量供 VariableSelector 使用
const formattedEnvVars = computed(() => {
  return props.envVars.map((v) => ({
    key: v.key,
    description: v.description || v.key,
  }))
})

// 更新重试配置
const updateRetryConfig = (config: NodeRetryConfig) => {
  localNode.value.retry = config
}

watch(
  () => props.modelValue,
  (val) => {
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
        // 只设置 enabled，其他由 TriggerConfig 组件管理
        localNode.value.config.enabled =
          localNode.value.config.enabled !== undefined ? localNode.value.config.enabled : true

        // 删除旧的 time 字段（如果存在）
        if ('time' in localNode.value.config) {
          delete localNode.value.config.time
        }
      }
    }
  }
)

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
      body: formattedBody,
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

const handleConfigUpdate = (config: any) => {
  localNode.value.config = config
}

const handleSave = () => {
  if (props.node) {
    const updateData = {
      name: localNode.value.name,
      config: localNode.value.config,
      retry: localNode.value.retry,
    }

    emit('update', props.node.id, updateData)
    handleClose()
  }
}

// 删除确认对话框
const showDeleteConfirm = ref(false)

const handleDelete = () => {
  if (props.node) {
    showDeleteConfirm.value = true
  }
}

const confirmDelete = () => {
  if (props.node) {
    emit('delete', props.node.id)
    showDeleteConfirm.value = false
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
          error: '请填写接口地址',
        }
        return
      }
    }

    message.info('正在测试运行节点...')

    // TODO: 调用API执行节点测试
    // const response = await nodeApi.testRun(localNode.value)

    // Mock 测试结果
    await new Promise((resolve) => setTimeout(resolve, 1500))

    if (localNode.value.toolCode === 'http_request') {
      testResult.value = {
        success: true,
        output: {
          status: 200,
          statusText: 'OK',
          headers: {
            'content-type': 'application/json',
          },
          body: {
            success: true,
            data: {
              id: 12345,
              message: '请求成功',
            },
          },
          responseTime: 234,
        },
      }
    } else if (localNode.value.toolCode === 'health_checker') {
      testResult.value = {
        success: true,
        output: {
          healthy: true,
          status: 200,
          responseTime: 156,
        },
      }
    } else if (localNode.value.toolCode === 'email') {
      testResult.value = {
        success: true,
        output: {
          sent: true,
          messageId: 'msg-12345',
        },
      }
    } else {
      testResult.value = {
        success: true,
        output: {
          result: '执行成功',
        },
      }
    }

    message.success('节点测试运行成功')
  } catch (error: any) {
    testResult.value = {
      success: false,
      error: error.message || '测试运行失败',
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
