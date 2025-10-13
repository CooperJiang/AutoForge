<template>
  <Dialog
    v-model="isOpen"
    :title="tool?.name || '工具详情'"
    max-width="max-w-3xl"
    hide-footer
    @cancel="close"
  >
    <div v-if="tool" class="space-y-6">
      <!-- 工具头部 -->
      <div class="flex items-start gap-4">
        <div
          :class="[
            'flex-shrink-0 w-16 h-16 rounded-xl flex items-center justify-center text-white shadow-lg',
            iconBgClass
          ]"
        >
          <component :is="iconComponent" class="w-8 h-8" />
        </div>
        <div class="flex-1">
          <p class="text-text-secondary mb-2">
            {{ tool.description }}
          </p>
          <div class="flex flex-wrap items-center gap-3 text-sm">
            <span class="flex items-center gap-1 text-text-secondary">
              <Package class="w-4 h-4" />
              版本 {{ tool.version }}
            </span>
            <span class="flex items-center gap-1 text-text-secondary">
              <User class="w-4 h-4" />
              {{ tool.author }}
            </span>
            <span class="flex items-center gap-1 text-text-secondary">
              <Tag class="w-4 h-4" />
              {{ tool.category }}
            </span>
          </div>
        </div>
      </div>

      <!-- 工具标签 -->
      <div>
        <h3 class="text-sm font-semibold text-text-secondary mb-2">标签</h3>
        <div class="flex flex-wrap gap-2">
          <span
            v-for="tag in tool.tags"
            :key="tag"
            class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium bg-primary-light text-primary border border-primary"
          >
            {{ tag }}
          </span>
        </div>
      </div>

      <!-- 使用说明 -->
      <div>
        <h3 class="text-sm font-semibold text-text-primary mb-2 flex items-center gap-2">
          <BookOpen class="w-4 h-4" />
          使用说明
        </h3>
        <div class="bg-bg-hover rounded-lg p-4 space-y-3">
          <div v-if="tool.code === 'http_request'">
            <h4 class="font-semibold text-text-primary mb-2 text-sm">📡 HTTP 请求工具</h4>
            <p class="text-xs text-text-secondary mb-2">
              发送 HTTP 请求到指定的 URL，支持所有常见的 HTTP 方法。
            </p>
            <ul class="space-y-1.5 text-xs text-text-secondary">
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>支持 GET、POST、PUT、DELETE、PATCH 等方法</span>
              </li>
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>自定义请求头（Headers）、参数（Params）、请求体（Body）</span>
              </li>
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>支持粘贴 cURL 命令自动解析配置</span>
              </li>
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>适用场景：API 调用、数据抓取、Webhook 触发等</span>
              </li>
            </ul>
          </div>

          <div v-else-if="tool.code === 'email_sender'">
            <h4 class="font-semibold text-text-primary mb-2 text-sm">📧 邮件发送工具</h4>
            <p class="text-xs text-text-secondary mb-2">
              通过 SMTP 协议发送邮件通知，支持多收件人和 HTML 格式。
            </p>
            <ul class="space-y-1.5 text-xs text-text-secondary">
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>系统自动使用配置的 SMTP 服务器，无需用户提供</span>
              </li>
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>支持多个收件人、抄送（CC）</span>
              </li>
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>支持纯文本和 HTML 格式</span>
              </li>
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>适用场景：告警通知、报表发送、验证码邮件等</span>
              </li>
            </ul>
          </div>

          <div v-else-if="tool.code === 'health_checker'">
            <h4 class="font-semibold text-text-primary mb-2 text-sm">🏥 健康检查工具</h4>
            <p class="text-xs text-text-secondary mb-2">
              监控网站或 API 的可用性，检查 SSL 证书有效期。
            </p>
            <ul class="space-y-1.5 text-xs text-text-secondary">
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>检查网站状态码和响应时间</span>
              </li>
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>SSL 证书到期检查和告警（可设置警戒天数）</span>
              </li>
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>支持正则表达式匹配响应内容</span>
              </li>
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>支持复杂鉴权（自定义 Headers 和 Body）</span>
              </li>
              <li class="flex items-start gap-2">
                <Check class="w-3.5 h-3.5 text-green-600 mt-0.5 flex-shrink-0" />
                <span>适用场景：网站监控、API 健康检查、SSL 证书管理等</span>
              </li>
            </ul>
          </div>

          <div v-else>
            <p class="text-xs text-text-secondary">
              暂无详细说明，请直接使用该工具。
            </p>
          </div>
        </div>
      </div>

      <!-- 使用按钮 -->
      <div class="flex items-center justify-end gap-2 pt-2">
        <BaseButton variant="ghost" size="md" @click="close">
          取消
        </BaseButton>
        <BaseButton variant="primary" size="md" @click="handleUseTool" class="whitespace-nowrap">
          <Rocket class="w-4 h-4 mr-1.5 inline-block flex-shrink-0" />
          <span>立即使用</span>
        </BaseButton>
      </div>
    </div>
  </Dialog>
</template>

<script setup lang="ts">
import { computed, watch, ref } from 'vue'
import {
  Package, User, Tag, BookOpen, Check, Rocket,
  Globe, Mail, HeartPulse
} from 'lucide-vue-next'
import Dialog from '@/components/Dialog'
import BaseButton from '@/components/BaseButton'

interface Tool {
  code: string
  name: string
  description: string
  category: string
  icon: string
  tags: string[]
  version: string
  author: string
}

interface Props {
  modelValue: boolean
  tool: Tool | null
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'use-tool': [toolCode: string]
}>()

const isOpen = ref(props.modelValue)

watch(() => props.modelValue, (val) => {
  isOpen.value = val
})

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

// 根据工具代码获取图标
const iconComponent = computed(() => {
  if (!props.tool) return Globe
  const iconMap: Record<string, any> = {
    'http_request': Globe,
    'email_sender': Mail,
    'health_checker': HeartPulse
  }
  return iconMap[props.tool.code] || Globe
})

// 根据工具代码获取图标背景色
const iconBgClass = computed(() => {
  if (!props.tool) return 'bg-gradient-to-br from-primary to-accent'
  const colorMap: Record<string, string> = {
    'http_request': 'bg-gradient-to-br from-primary to-accent',
    'email_sender': 'bg-gradient-to-br from-purple-500 to-pink-600',
    'health_checker': 'bg-gradient-to-br from-primary to-accent'
  }
  return colorMap[props.tool.code] || 'bg-gradient-to-br from-primary to-accent'
})

const close = () => {
  isOpen.value = false
}

const handleUseTool = () => {
  if (props.tool) {
    emit('use-tool', props.tool.code)
    close()
  }
}
</script>
