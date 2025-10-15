<template>
  <div class="space-y-4">
    
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        Webhook 地址 <span class="text-red-500">*</span>
      </label>
      <VariableSelector
        v-model="localConfig.webhook_url"
        placeholder="https://open.feishu.cn/open-apis/bot/v2/hook/..."
        :previous-nodes="previousNodes"
        :env-vars="formattedEnvVars"
      />
      <p class="mt-1 text-xs text-text-tertiary">在飞书群中添加自定义机器人后获取</p>
    </div>

    
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2"> 签名密钥（可选） </label>
      <BaseInput
        v-model="localConfig.sign_secret"
        type="password"
        placeholder="留空表示不使用签名验证"
      />
      <p class="mt-1 text-xs text-text-tertiary">启用签名验证可以提高安全性</p>
    </div>

    
    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        消息类型 <span class="text-red-500">*</span>
      </label>
      <BaseSelect
        v-model="localConfig.msg_type"
        :options="msgTypeOptions"
        @update:model-value="handleMsgTypeChange"
      />
      <p class="mt-1 text-xs text-text-tertiary">💡 切换消息类型后,界面只显示当前类型的配置项</p>
    </div>

    
    <template v-if="localConfig.msg_type === 'text'">
      <div>
        <label
          class="block text-sm font-medium text-text-secondary mb-2 flex items-center justify-between"
        >
          <span>消息内容 <span class="text-red-500">*</span></span>
          <button
            type="button"
            @click="showVariableHelper = !showVariableHelper"
            class="text-xs text-primary hover:text-primary"
          >
            {{ showVariableHelper ? '隐藏' : '显示' }}变量助手
          </button>
        </label>

        
        <VariableHelper
          :show="showVariableHelper"
          :previous-nodes="previousNodes"
          :env-vars="formattedEnvVars"
          @insert-field="
            (nodeId, fieldName) => insertFieldVariable(nodeId, fieldName, contentTextareaRef)
          "
          @insert-node="(nodeId) => insertNodeVariable(nodeId, contentTextareaRef)"
          @insert-env="(key) => insertEnvVariable(key, contentTextareaRef)"
        />

        <textarea
          ref="contentTextareaRef"
          v-model="localConfig.content"
          rows="4"
          placeholder="输入要发送的文本内容...&#10;&#10;示例：执行结果：{{nodes.node_xxx.message}}"
          class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
        />
      </div>
    </template>

    
    <template v-if="localConfig.msg_type === 'post'">
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2"> 标题 </label>
        <VariableSelector
          v-model="localConfig.title"
          placeholder="消息标题"
          :previous-nodes="previousNodes"
          :env-vars="formattedEnvVars"
        />
      </div>

      <div>
        <label
          class="block text-sm font-medium text-text-secondary mb-2 flex items-center justify-between"
        >
          <span>富文本内容 <span class="text-red-500">*</span></span>
          <button
            type="button"
            @click="showVariableHelper = !showVariableHelper"
            class="text-xs text-primary hover:text-primary"
          >
            {{ showVariableHelper ? '隐藏' : '显示' }}变量助手
          </button>
        </label>

        
        <VariableHelper
          :show="showVariableHelper"
          :previous-nodes="previousNodes"
          :env-vars="formattedEnvVars"
          @insert-field="
            (nodeId, fieldName) => insertFieldVariable(nodeId, fieldName, postContentTextareaRef)
          "
          @insert-node="(nodeId) => insertNodeVariable(nodeId, postContentTextareaRef)"
          @insert-env="(key) => insertEnvVariable(key, postContentTextareaRef)"
        />

        <textarea
          ref="postContentTextareaRef"
          v-model="localConfig.post_content"
          rows="6"
          placeholder="支持 Markdown 格式：&#10;**粗体** *斜体* [链接](url)&#10;- 列表项&#10;&#10;示例：{{nodes.node_xxx.message}}"
          class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
        />
        <p class="mt-1 text-xs text-text-tertiary">
          支持 Markdown 语法：**粗体** *斜体* [链接文字](URL)
        </p>
      </div>
    </template>

    
    <template v-if="localConfig.msg_type === 'image'">
      
      <div
        class="bg-amber-50 dark:bg-amber-900/20 border border-amber-200 dark:border-amber-800 rounded-lg p-3 mb-4"
      >
        <p class="text-xs text-amber-800 dark:text-amber-200 mb-2">
          💡 <strong>要直接显示图片需要填写以下信息：</strong>
        </p>
        <ul class="text-xs text-amber-700 dark:text-amber-300 ml-4 space-y-1">
          <li>• 在飞书开放平台创建应用并获取 App ID 和 App Secret</li>
          <li>• 为应用添加"获取与上传图片"权限 (im:image)</li>
          <li>• 如不填写，将显示图片链接(需点击查看)</li>
        </ul>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2"> App ID (可选) </label>
        <BaseInput v-model="localConfig.app_id" placeholder="cli_xxxxxxxxxx" />
        <p class="mt-1 text-xs text-text-tertiary">飞书应用的 App ID，在"凭证与基础信息"中查看</p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          App Secret (可选)
        </label>
        <BaseInput v-model="localConfig.app_secret" type="password" placeholder="输入应用密钥" />
        <p class="mt-1 text-xs text-text-tertiary">
          飞书应用的 App Secret，在"凭证与基础信息"中查看
        </p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2"> 标题 </label>
        <VariableSelector
          v-model="localConfig.title"
          placeholder="图片消息标题(可选)"
          :previous-nodes="previousNodes"
          :env-vars="formattedEnvVars"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          图片 URL <span class="text-red-500">*</span>
        </label>
        <BaseInput
          v-model="localConfig.image_url"
          placeholder="https://example.com/image.png"
          @update:model-value="(val) => (localConfig.image_url = val)"
        />
        <p class="mt-1 text-xs text-text-tertiary">图片必须是公网可访问的 URL，大小 &lt; 10MB</p>
        <p class="mt-1 text-xs text-amber-600" v-if="!localConfig.image_url">⚠️ 请输入图片 URL</p>
      </div>
    </template>

    
    <template v-if="localConfig.msg_type === 'interactive'">
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2"> 卡片模板 </label>
        <BaseSelect v-model="localConfig.card_template" :options="cardTemplateOptions" />
      </div>

      <template v-if="localConfig.card_template !== 'custom'">
        <div>
          <label class="block text-sm font-medium text-text-secondary mb-2">
            标题 <span class="text-red-500">*</span>
          </label>
          <VariableSelector
            v-model="localConfig.title"
            placeholder="卡片标题"
            :previous-nodes="previousNodes"
            :env-vars="formattedEnvVars"
          />
        </div>

        <div>
          <label
            class="block text-sm font-medium text-text-secondary mb-2 flex items-center justify-between"
          >
            <span>内容</span>
            <button
              type="button"
              @click="showVariableHelper = !showVariableHelper"
              class="text-xs text-primary hover:text-primary"
            >
              {{ showVariableHelper ? '隐藏' : '显示' }}变量助手
            </button>
          </label>

          
          <VariableHelper
            :show="showVariableHelper"
            :previous-nodes="previousNodes"
            :env-vars="formattedEnvVars"
            @insert-field="
              (nodeId, fieldName) => insertFieldVariable(nodeId, fieldName, cardContentTextareaRef)
            "
            @insert-node="(nodeId) => insertNodeVariable(nodeId, cardContentTextareaRef)"
            @insert-env="(key) => insertEnvVariable(key, cardContentTextareaRef)"
          />

          <textarea
            ref="cardContentTextareaRef"
            v-model="localConfig.card_content"
            rows="3"
            placeholder="卡片主要内容描述...&#10;&#10;示例：{{nodes.node_xxx.message}}"
            class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-text-secondary mb-2"> 状态 </label>
          <BaseSelect v-model="localConfig.card_status" :options="cardStatusOptions" />
        </div>

        <div>
          <label class="block text-sm font-medium text-text-secondary mb-2">
            字段列表（JSON 格式）
          </label>
          <textarea
            v-model="localConfig.card_fields"
            rows="3"
            placeholder='[{"key":"任务名称","value":"数据同步"},{"key":"执行时间","value":"14:30:00"}]'
            class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
          />
          <p class="mt-1 text-xs text-text-tertiary">
            JSON 数组格式，每个对象包含 key 和 value 字段
          </p>
        </div>

        <div>
          <label class="block text-sm font-medium text-text-secondary mb-2">
            按钮列表（JSON 格式）
          </label>
          <textarea
            v-model="localConfig.card_buttons"
            rows="2"
            placeholder='[{"text":"查看详情","url":"https://example.com"}]'
            class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
          />
          <p class="mt-1 text-xs text-text-tertiary">
            JSON 数组格式，每个对象包含 text 和 url 字段
          </p>
        </div>
      </template>

      <template v-else>
        <div>
          <label class="block text-sm font-medium text-text-secondary mb-2">
            自定义卡片 JSON <span class="text-red-500">*</span>
          </label>
          <textarea
            v-model="localConfig.card_custom_json"
            rows="10"
            placeholder='{"config":{},"header":{},"elements":[]}'
            class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-primary bg-bg-elevated text-text-primary font-mono text-sm"
          />
          <p class="mt-1 text-xs text-text-tertiary">
            完整的飞书卡片 JSON 格式，参考
            <a
              href="https://open.feishu.cn/document/common-capabilities/message-card/message-cards-content/using-markdown-tags"
              target="_blank"
              class="text-primary hover:underline"
            >
              飞书开放平台文档
            </a>
          </p>
        </div>
      </template>
    </template>

    
    <div class="mt-4 p-3 bg-bg-secondary rounded-lg border border-border-primary">
      <div class="flex items-start gap-2">
        <span class="text-primary text-lg">💡</span>
        <div class="flex-1 text-sm text-text-secondary">
          <p class="font-medium mb-1">快速开始：</p>
          <ol class="list-decimal list-inside space-y-1 text-xs">
            <li>在飞书群中添加"自定义机器人"</li>
            <li>复制 Webhook URL 并粘贴到上方</li>
            <li>选择消息类型并配置内容</li>
            <li>保存工作流即可开始使用</li>
          </ol>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import VariableSelector from '@/components/VariableSelector'
import VariableHelper from '@/components/VariableHelper'

interface Props {
  config: Record<string, any>
  previousNodes?: Array<{ id: string; name: string; type: string; toolCode?: string }>
  envVars?: Array<{ key: string; value: string; description?: string }>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:config', config: Record<string, any>): void
}>()

// 本地配置状态
const localConfig = ref({
  webhook_url: props.config.webhook_url || '',
  sign_secret: props.config.sign_secret || '',
  app_id: props.config.app_id || '',
  app_secret: props.config.app_secret || '',
  msg_type: props.config.msg_type || 'text',
  content: props.config.content || '',
  title: props.config.title || '',
  post_content: props.config.post_content || '',
  image_url: props.config.image_url || '',
  card_template: props.config.card_template || 'notification',
  card_content: props.config.card_content || '',
  card_status: props.config.card_status || 'info',
  card_fields: props.config.card_fields || '',
  card_buttons: props.config.card_buttons || '',
  card_custom_json: props.config.card_custom_json || '',
})

// 消息类型选项
const msgTypeOptions = [
  { label: '文本消息', value: 'text' },
  { label: '富文本消息', value: 'post' },
  { label: '图片消息', value: 'image' },
  { label: '卡片消息', value: 'interactive' },
]

// 卡片模板选项
const cardTemplateOptions = [
  { label: '通知卡片', value: 'notification' },
  { label: '告警卡片', value: 'alert' },
  { label: '报告卡片', value: 'report' },
  { label: '自定义 JSON', value: 'custom' },
]

// 卡片状态选项
const cardStatusOptions = [
  { label: '成功 ✅', value: 'success' },
  { label: '警告 ⚠️', value: 'warning' },
  { label: '错误 ❌', value: 'error' },
  { label: '信息 ℹ️', value: 'info' },
]

// 变量助手状态
const showVariableHelper = ref(false)
const contentTextareaRef = ref<HTMLTextAreaElement>()
const postContentTextareaRef = ref<HTMLTextAreaElement>()
const cardContentTextareaRef = ref<HTMLTextAreaElement>()
const cardFieldsTextareaRef = ref<HTMLTextAreaElement>()
const cardButtonsTextareaRef = ref<HTMLTextAreaElement>()
const cardCustomJsonTextareaRef = ref<HTMLTextAreaElement>()

// 格式化环境变量
const formattedEnvVars = computed(() => {
  return props.envVars || []
})

// 上一次的消息类型
const previousMsgType = ref(localConfig.value.msg_type)

// 处理消息类型变化
const handleMsgTypeChange = (newType: string) => {
  if (previousMsgType.value !== newType) {
    // 根据新类型,只清空不需要的字段
    switch (newType) {
      case 'text':
        // 切换到文本,清空其他类型字段
        localConfig.value.title = ''
        localConfig.value.post_content = ''
        localConfig.value.image_url = ''
        localConfig.value.card_template = 'notification'
        localConfig.value.card_content = ''
        localConfig.value.card_status = 'info'
        localConfig.value.card_fields = ''
        localConfig.value.card_buttons = ''
        localConfig.value.card_custom_json = ''
        break
      case 'post':
        // 切换到富文本,清空其他类型字段
        localConfig.value.content = ''
        localConfig.value.image_url = ''
        localConfig.value.card_template = 'notification'
        localConfig.value.card_content = ''
        localConfig.value.card_status = 'info'
        localConfig.value.card_fields = ''
        localConfig.value.card_buttons = ''
        localConfig.value.card_custom_json = ''
        break
      case 'image':
        // 切换到图片,清空其他类型字段
        localConfig.value.content = ''
        localConfig.value.post_content = ''
        localConfig.value.card_template = 'notification'
        localConfig.value.card_content = ''
        localConfig.value.card_status = 'info'
        localConfig.value.card_fields = ''
        localConfig.value.card_buttons = ''
        localConfig.value.card_custom_json = ''
        break
      case 'interactive':
        // 切换到卡片,清空其他类型字段
        localConfig.value.content = ''
        localConfig.value.post_content = ''
        localConfig.value.image_url = ''
        break
    }

    previousMsgType.value = newType
  }
}

// 插入字段变量
const insertFieldVariable = (
  nodeId: string,
  fieldName: string,
  targetRef?: { value?: HTMLTextAreaElement }
) => {
  insertToTextarea(`{{nodes.${nodeId}.${fieldName}}}`, targetRef)
}

// 插入节点变量
const insertNodeVariable = (nodeId: string, targetRef?: { value?: HTMLTextAreaElement }) => {
  insertToTextarea(`{{nodes.${nodeId}.`, targetRef)
}

// 插入环境变量
const insertEnvVariable = (key: string, targetRef?: { value?: HTMLTextAreaElement }) => {
  insertToTextarea(`{{env.${key}}}`, targetRef)
}

// 插入变量到 textarea
const insertToTextarea = (text: string, targetRef?: { value?: HTMLTextAreaElement }) => {
  const textarea = targetRef?.value
  if (!textarea) return

  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const currentValue = textarea.value || ''

  // 根据 textarea 更新对应的 config 字段
  if (textarea === contentTextareaRef.value) {
    localConfig.value.content =
      currentValue.substring(0, start) + text + currentValue.substring(end)
  } else if (textarea === postContentTextareaRef.value) {
    localConfig.value.post_content =
      currentValue.substring(0, start) + text + currentValue.substring(end)
  } else if (textarea === cardContentTextareaRef.value) {
    localConfig.value.card_content =
      currentValue.substring(0, start) + text + currentValue.substring(end)
  } else if (textarea === cardFieldsTextareaRef.value) {
    localConfig.value.card_fields =
      currentValue.substring(0, start) + text + currentValue.substring(end)
  } else if (textarea === cardButtonsTextareaRef.value) {
    localConfig.value.card_buttons =
      currentValue.substring(0, start) + text + currentValue.substring(end)
  } else if (textarea === cardCustomJsonTextareaRef.value) {
    localConfig.value.card_custom_json =
      currentValue.substring(0, start) + text + currentValue.substring(end)
  }

  // 恢复光标位置
  setTimeout(() => {
    textarea.focus()
    const newPos = start + text.length
    textarea.setSelectionRange(newPos, newPos)
  }, 0)
}

// 监听外部配置变化（只在初始化和外部变化时更新）
watch(
  () => props.config,
  (newConfig) => {
    // 只在配置真正改变时更新，避免循环
    const hasChanged = Object.keys(localConfig.value).some((key) => {
      return localConfig.value[key] !== (newConfig[key] || '')
    })

    if (hasChanged) {
      localConfig.value = {
        webhook_url: newConfig.webhook_url || '',
        sign_secret: newConfig.sign_secret || '',
        app_id: newConfig.app_id || '',
        app_secret: newConfig.app_secret || '',
        msg_type: newConfig.msg_type || 'text',
        content: newConfig.content || '',
        title: newConfig.title || '',
        post_content: newConfig.post_content || '',
        image_url: newConfig.image_url || '',
        card_template: newConfig.card_template || 'notification',
        card_content: newConfig.card_content || '',
        card_status: newConfig.card_status || 'info',
        card_fields: newConfig.card_fields || '',
        card_buttons: newConfig.card_buttons || '',
        card_custom_json: newConfig.card_custom_json || '',
      }
    }
  },
  { immediate: true }
)

// 监听本地配置变化，同步到外部（使用节流）
let updateTimeout: number | null = null
watch(
  localConfig,
  (newConfig) => {
    if (updateTimeout) {
      clearTimeout(updateTimeout)
    }
    updateTimeout = window.setTimeout(() => {
      emit('update:config', { ...newConfig })
    }, 100)
  },
  { deep: true }
)
</script>
