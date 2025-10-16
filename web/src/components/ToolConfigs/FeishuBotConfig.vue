<template>
  <div class="space-y-4">
    <div>
      <label class="block text-sm font-medium text-text-primary mb-1">
        Webhook 地址 <span class="text-error">*</span>
      </label>
      <input
        v-model="localConfig.webhook_url"
        type="text"
        placeholder="https://open.feishu.cn/open-apis/bot/v2/hook/..."
        class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary"
        @input="emitUpdate"
      />
      <p class="mt-1 text-xs text-text-tertiary">在飞书群中添加自定义机器人后获取</p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-primary mb-1"> 签名密钥（可选） </label>
      <input
        v-model="localConfig.sign_secret"
        type="password"
        placeholder="留空表示不使用签名验证"
        class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary"
        @input="emitUpdate"
      />
      <p class="mt-1 text-xs text-text-tertiary">启用签名验证可以提高安全性</p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-primary mb-1">
        消息类型 <span class="text-error">*</span>
      </label>
      <select
        v-model="localConfig.msg_type"
        class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary"
        @change="emitUpdate"
      >
        <option value="text">文本消息</option>
        <option value="post">富文本消息</option>
        <option value="image">图片消息</option>
        <option value="interactive">卡片消息</option>
      </select>
    </div>

    <template v-if="localConfig.msg_type === 'text'">
      <div>
        <label class="block text-sm font-medium text-text-primary mb-1">
          消息内容 <span class="text-error">*</span>
        </label>
        <textarea
          v-model="localConfig.content"
          rows="4"
          placeholder="输入要发送的文本内容..."
          class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary"
          @input="emitUpdate"
        />
      </div>
    </template>

    <template v-if="localConfig.msg_type === 'post'">
      <div>
        <label class="block text-sm font-medium text-text-primary mb-1"> 标题 </label>
        <input
          v-model="localConfig.title"
          type="text"
          placeholder="消息标题"
          class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary"
          @input="emitUpdate"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-text-primary mb-1">
          富文本内容 <span class="text-error">*</span>
        </label>
        <textarea
          v-model="localConfig.post_content"
          rows="6"
          placeholder="支持 Markdown 格式：&#10;**粗体** *斜体* [链接](url)&#10;- 列表项"
          class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary font-mono text-sm"
          @input="emitUpdate"
        />
        <p class="mt-1 text-xs text-text-tertiary">
          支持 Markdown 语法：**粗体** *斜体* [链接文字](URL)
        </p>
      </div>
    </template>

    <template v-if="localConfig.msg_type === 'image'">
      <div>
        <label class="block text-sm font-medium text-text-primary mb-1">
          图片 URL <span class="text-error">*</span>
        </label>
        <input
          v-model="localConfig.image_url"
          type="text"
          placeholder="https://example.com/image.png"
          class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary"
          @input="emitUpdate"
        />
        <p class="mt-1 text-xs text-text-tertiary">图片必须是公网可访问的 URL，大小 &lt; 10MB</p>
      </div>
    </template>

    <template v-if="localConfig.msg_type === 'interactive'">
      <div>
        <label class="block text-sm font-medium text-text-primary mb-1"> 卡片模板 </label>
        <select
          v-model="localConfig.card_template"
          class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary"
          @change="emitUpdate"
        >
          <option value="notification">通知卡片</option>
          <option value="alert">告警卡片</option>
          <option value="report">报告卡片</option>
          <option value="custom">自定义 JSON</option>
        </select>
      </div>

      <template v-if="localConfig.card_template !== 'custom'">
        <div>
          <label class="block text-sm font-medium text-text-primary mb-1">
            标题 <span class="text-error">*</span>
          </label>
          <input
            v-model="localConfig.title"
            type="text"
            placeholder="卡片标题"
            class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary"
            @input="emitUpdate"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-text-primary mb-1"> 内容 </label>
          <textarea
            v-model="localConfig.card_content"
            rows="3"
            placeholder="卡片主要内容描述..."
            class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary"
            @input="emitUpdate"
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-text-primary mb-1"> 状态 </label>
          <select
            v-model="localConfig.card_status"
            class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary"
            @change="emitUpdate"
          >
            <option value="success">成功 ✅</option>
            <option value="warning">警告 ⚠️</option>
            <option value="error">错误 ❌</option>
            <option value="info">信息 ℹ️</option>
          </select>
        </div>

        <div>
          <label class="block text-sm font-medium text-text-primary mb-1">
            字段列表（JSON 格式）
          </label>
          <textarea
            v-model="localConfig.card_fields"
            rows="3"
            placeholder='[{"key":"任务名称","value":"数据同步"},{"key":"执行时间","value":"14:30:00"}]'
            class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary font-mono text-sm"
            @input="emitUpdate"
          />
          <p class="mt-1 text-xs text-text-tertiary">
            JSON 数组格式，每个对象包含 key 和 value 字段
          </p>
        </div>

        <div>
          <label class="block text-sm font-medium text-text-primary mb-1">
            按钮列表（JSON 格式）
          </label>
          <textarea
            v-model="localConfig.card_buttons"
            rows="2"
            placeholder='[{"text":"查看详情","url":"https://example.com"}]'
            class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary font-mono text-sm"
            @input="emitUpdate"
          />
          <p class="mt-1 text-xs text-text-tertiary">
            JSON 数组格式，每个对象包含 text 和 url 字段
          </p>
        </div>
      </template>

      <template v-else>
        <div>
          <label class="block text-sm font-medium text-text-primary mb-1">
            自定义卡片 JSON <span class="text-error">*</span>
          </label>
          <textarea
            v-model="localConfig.card_custom_json"
            rows="10"
            placeholder='{"config":{},"header":{},"elements":[]}'
            class="w-full px-3 py-2 border border-border-primary rounded-lg bg-bg-primary text-text-primary focus:outline-none focus:ring-2 focus:ring-primary font-mono text-sm"
            @input="emitUpdate"
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
import { ref, watch } from 'vue'

interface Props {
  config: Record<string, any>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  (e: 'update:config', config: Record<string, any>): void
}>()

// 本地配置状态
const localConfig = ref({
  webhook_url: props.config.webhook_url || '',
  sign_secret: props.config.sign_secret || '',
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

// 监听外部配置变化
watch(
  () => props.config,
  (newConfig) => {
    localConfig.value = {
      webhook_url: newConfig.webhook_url || '',
      sign_secret: newConfig.sign_secret || '',
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
  },
  { deep: true }
)

// 发送更新事件
const emitUpdate = () => {
  emit('update:config', { ...localConfig.value })
}
</script>
