<template>
  <main class="flex h-screen bg-bg-hover">
    <div class="w-[480px] bg-bg-elevated border-r border-border-primary overflow-y-auto">
      <div class="p-6">
        <div class="mb-6">
          <h1 class="text-2xl font-bold text-text-primary">
            {{ isEditing ? '编辑任务' : '添加定时任务' }}
          </h1>
          <p class="text-sm text-text-secondary mt-1">
            {{ isEditing ? '修改任务配置' : '选择工具创建定时任务' }}
          </p>
        </div>

        <form @submit.prevent="handleSubmit">
          <div class="mb-4">
            <label class="block text-sm font-medium text-text-secondary mb-2">
              任务名称 <span class="text-red-500">*</span>
            </label>
            <BaseInput v-model="form.name" placeholder="例如：每日签到" required />
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium text-text-secondary mb-2">
              执行规则 <span class="text-red-500">*</span>
            </label>
            <BaseSelect v-model="form.schedule_type" :options="scheduleTypeOptions" required />
          </div>

          <div class="mb-4">
            <BaseInput
              v-model="form.schedule_value"
              :placeholder="getScheduleValuePlaceholder()"
              required
            />
            <p class="text-xs text-text-tertiary mt-1">
              {{ getScheduleValueHint() }}
            </p>
          </div>

          <div class="mb-4">
            <label class="block text-sm font-medium text-text-secondary mb-2">
              选择工具 <span class="text-red-500">*</span>
            </label>
            <BaseSelect
              v-model="form.tool_code"
              :options="toolOptions"
              placeholder="请选择工具"
              required
              @change="handleToolChange"
            />
          </div>

          <div v-if="form.tool_code" class="mb-4">
            <label class="block text-sm font-medium text-text-secondary mb-2">
              工具配置 <span class="text-red-500">*</span>
            </label>
            <BaseButton
              type="button"
              variant="secondary"
              @click="showConfigDialog = true"
              class="w-full"
            >
              {{ isConfigured ? '✓ 已配置 - 点击修改' : '配置工具参数' }}
            </BaseButton>
          </div>

          <div class="flex gap-2 mt-6">
            <BaseButton type="button" variant="secondary" @click="goBack" class="flex-1">
              取消
            </BaseButton>
            <BaseButton
              type="submit"
              variant="primary"
              :disabled="submitting || !isConfigured"
              class="flex-1"
            >
              {{ submitting ? '创建中...' : '创建任务' }}
            </BaseButton>
          </div>
        </form>
      </div>
    </div>

    <div class="flex-1 p-6 overflow-y-auto">
      <div class="bg-bg-elevated border-2 border-border-primary rounded-lg p-6">
        <h2 class="text-lg font-semibold text-text-primary mb-4">任务预览</h2>

        <div v-if="form.name" class="mb-4">
          <div class="text-sm text-text-secondary">任务名称</div>
          <div class="text-base font-medium text-text-primary mt-1">{{ form.name }}</div>
        </div>

        <div v-if="form.tool_code" class="mb-4">
          <div class="text-sm text-text-secondary">使用工具</div>
          <div class="text-base font-medium text-text-primary mt-1">
            {{ selectedTool?.name || form.tool_code }}
          </div>
        </div>

        <div v-if="form.schedule_type" class="mb-4">
          <div class="text-sm text-text-secondary">执行规则</div>
          <div class="text-base font-medium text-text-primary mt-1">
            {{ getScheduleTypeName(form.schedule_type) }}：{{ form.schedule_value || '-' }}
          </div>
        </div>

        <div v-if="isConfigured" class="mb-4">
          <div class="text-sm text-text-secondary">配置状态</div>
          <div class="text-base font-medium text-emerald-600 mt-1">✓ 已完成配置</div>
        </div>
      </div>
    </div>

    <Dialog
      v-model="showConfigDialog"
      title="配置工具参数"
      width="600px"
      @confirm="saveConfig"
      confirm-text="保存配置"
    >
      <div v-if="form.tool_code === 'http_tool'" class="space-y-4">
        <div>
          <label class="block text-sm font-medium text-text-secondary mb-2">
            请求方式 <span class="text-red-500">*</span>
          </label>
          <BaseSelect v-model="toolConfig.method" :options="methodOptions" required />
        </div>

        <div>
          <label class="block text-sm font-medium text-text-secondary mb-2">
            接口地址 <span class="text-red-500">*</span>
          </label>
          <BaseInput
            v-model="toolConfig.url"
            placeholder="https://api.example.com/checkin"
            required
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-text-secondary mb-2"> 请求头（可选） </label>
          <textarea
            v-model="toolConfig.headers"
            class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-green-500 font-mono text-sm"
            rows="4"
            placeholder='{"Content-Type": "application/json"}'
          />
        </div>

        <div>
          <label class="block text-sm font-medium text-text-secondary mb-2">
            请求参数（可选）
          </label>
          <textarea
            v-model="toolConfig.body"
            class="w-full px-3 py-2 border-2 border-border-primary rounded-lg focus:outline-none focus:border-green-500 font-mono text-sm"
            rows="4"
            placeholder='{"key": "value"}'
          />
        </div>
      </div>

      <div v-else class="text-center py-8 text-text-tertiary">该工具暂无需配置参数</div>
    </Dialog>
  </main>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { message } from '@/utils/message'
import { createTask, updateTask, getTask } from '@/api/task'
import { getToolList, type Tool } from '@/api/tool'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'
import Dialog from '@/components/Dialog'

const router = useRouter()
const route = useRoute()
const isEditing = computed(() => !!route.params.id)
const submitting = ref(false)
const tools = ref<Tool[]>([])
const showConfigDialog = ref(false)

const form = ref({
  name: '',
  description: '',
  tool_code: '',
  schedule_type: 'daily',
  schedule_value: '09:00:00',
})

const toolConfig = ref<Record<string, any>>({
  url: '',
  method: 'GET',
  headers: '{}',
  body: '{}',
})

const isConfigured = ref(false)

const methodOptions = [
  { label: 'GET', value: 'GET' },
  { label: 'POST', value: 'POST' },
  { label: 'PUT', value: 'PUT' },
  { label: 'DELETE', value: 'DELETE' },
]

const scheduleTypeOptions = [
  { label: '每天', value: 'daily' },
  { label: '每周', value: 'weekly' },
  { label: '每月', value: 'monthly' },
  { label: '每小时', value: 'hourly' },
  { label: '间隔执行', value: 'interval' },
  { label: 'Cron表达式', value: 'cron' },
]

const toolOptions = computed(() =>
  tools.value.map((tool) => ({ label: tool.name, value: tool.code }))
)
const selectedTool = computed(() => tools.value.find((t) => t.code === form.value.tool_code))

onMounted(async () => {
  await loadTools()
  if (isEditing.value) {
    await loadTask()
  }
})

const loadTools = async () => {
  try {
    tools.value = await getToolList()
  } catch (error: any) {
    message.error(error.message || '加载工具列表失败')
  }
}

const loadTask = async () => {
  try {
    const task = await getTask(route.params.id as string)
    form.value = {
      name: task.name,
      description: task.description,
      tool_code: task.tool_code,
      schedule_type: task.schedule_type,
      schedule_value: task.schedule_value,
    }

    // 解析已有配置
    const config = JSON.parse(task.config)
    toolConfig.value = {
      url: config.url || '',
      method: config.method || 'GET',
      headers: JSON.stringify(config.headers || {}, null, 2),
      body: JSON.stringify(config.body || {}, null, 2),
    }
    isConfigured.value = true
  } catch (error: any) {
    message.error(error.message || '加载任务失败')
  }
}

const handleToolChange = () => {
  toolConfig.value = { url: '', method: 'GET', headers: '{}', body: '{}' }
  isConfigured.value = false
}

const saveConfig = () => {
  if (form.value.tool_code === 'http_tool') {
    if (!toolConfig.value.url) {
      message.error('请输入请求URL')
      return
    }
  }
  isConfigured.value = true
  showConfigDialog.value = false
}

const getScheduleTypeName = (type: string) => {
  const option = scheduleTypeOptions.find((o) => o.value === type)
  return option?.label || type
}

const getScheduleValuePlaceholder = () => {
  const map: Record<string, string> = {
    cron: '0 0 * * * *',
    daily: '09:00:00',
    weekly: '1:09:00:00',
    monthly: '1:09:00:00',
    hourly: '30:00',
    interval: '3600',
  }
  return map[form.value.schedule_type] || ''
}

const getScheduleValueHint = () => {
  const map: Record<string, string> = {
    daily: '每天在指定时间执行',
    weekly: '每周指定日期执行（1-7代表周一到周日）',
    monthly: '每月指定日期执行',
    hourly: '每小时的指定分秒执行',
    interval: '间隔秒数执行',
    cron: 'Cron表达式',
  }
  return map[form.value.schedule_type] || ''
}

const handleSubmit = async () => {
  if (!isConfigured.value) {
    message.error('请先配置工具参数')
    return
  }

  let config: Record<string, any> = {}
  if (form.value.tool_code === 'http_tool') {
    try {
      config = {
        url: toolConfig.value.url,
        method: toolConfig.value.method,
        headers: JSON.parse(toolConfig.value.headers || '{}'),
        body: JSON.parse(toolConfig.value.body || '{}'),
      }
    } catch {
      message.error('Headers 或 Body 不是有效的 JSON 格式')
      return
    }
  }

  submitting.value = true
  try {
    const data = {
      name: form.value.name,
      description: form.value.description,
      tool_code: form.value.tool_code,
      config,
      schedule_type: form.value.schedule_type,
      schedule_value: form.value.schedule_value,
    }

    if (isEditing.value) {
      await updateTask(route.params.id as string, data)
      message.success('更新成功')
    } else {
      await createTask(data)
      message.success('创建成功')
    }
    router.push('/')
  } catch (error: any) {
    message.error(error.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

const goBack = () => router.push('/')
</script>
