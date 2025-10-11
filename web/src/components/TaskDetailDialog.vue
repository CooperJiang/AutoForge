<template>
  <div v-if="modelValue" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
    <div class="bg-white rounded-xl shadow-2xl max-w-3xl w-full max-h-[90vh] overflow-hidden flex flex-col">
      <!-- Header -->
      <div class="px-6 py-4 border-b-2 border-slate-200 flex justify-between items-center flex-shrink-0">
        <h3 class="text-xl font-bold text-slate-900">任务详情</h3>
        <button
          @click="$emit('update:modelValue', false)"
          class="text-slate-400 hover:text-slate-600 transition-colors"
        >
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Content -->
      <div class="flex-1 overflow-y-auto p-6">
        <div v-if="task" class="space-y-6">
          <!-- 基本信息 -->
          <div>
            <h4 class="text-sm font-semibold text-slate-700 mb-3 flex items-center gap-2">
              <span class="w-1 h-4 bg-blue-500 rounded"></span>
              基本信息
            </h4>
            <div class="grid grid-cols-2 gap-4 bg-slate-50 rounded-lg p-4">
              <div>
                <p class="text-xs text-slate-500 mb-1">任务ID</p>
                <p class="text-sm text-slate-900 font-mono">{{ task.id }}</p>
              </div>
              <div>
                <p class="text-xs text-slate-500 mb-1">任务名称</p>
                <p class="text-sm text-slate-900 font-medium">{{ task.name }}</p>
              </div>
              <div>
                <p class="text-xs text-slate-500 mb-1">用户ID</p>
                <p class="text-sm text-slate-900 font-mono">{{ task.user_id }}</p>
              </div>
              <div>
                <p class="text-xs text-slate-500 mb-1">状态</p>
                <span :class="[
                  'inline-flex px-2 py-1 text-xs font-medium rounded-full',
                  task.enabled ? 'bg-green-100 text-green-700' : 'bg-gray-100 text-gray-700'
                ]">
                  {{ task.enabled ? '已启用' : '已禁用' }}
                </span>
              </div>
            </div>
          </div>

          <!-- 请求配置 -->
          <div>
            <h4 class="text-sm font-semibold text-slate-700 mb-3 flex items-center gap-2">
              <span class="w-1 h-4 bg-green-500 rounded"></span>
              请求配置
            </h4>
            <div class="bg-slate-50 rounded-lg p-4 space-y-3">
              <div>
                <p class="text-xs text-slate-500 mb-1">请求方法</p>
                <span class="inline-flex px-2 py-1 text-xs font-semibold rounded bg-blue-100 text-blue-700">
                  {{ getConfigValue('method', 'GET') }}
                </span>
              </div>
              <div>
                <p class="text-xs text-slate-500 mb-1">请求URL</p>
                <p class="text-sm text-slate-900 font-mono break-all">{{ getConfigValue('url', '-') }}</p>
              </div>
              <div v-if="getConfigValue('headers')">
                <p class="text-xs text-slate-500 mb-1">请求头</p>
                <pre class="text-xs text-slate-900 bg-white rounded p-2 border border-slate-200 overflow-x-auto">{{ formatJSON(getConfigValue('headers')) }}</pre>
              </div>
              <div v-if="getConfigValue('body')">
                <p class="text-xs text-slate-500 mb-1">请求体</p>
                <pre class="text-xs text-slate-900 bg-white rounded p-2 border border-slate-200 overflow-x-auto">{{ formatJSON(getConfigValue('body')) }}</pre>
              </div>
            </div>
          </div>

          <!-- 调度配置 -->
          <div>
            <h4 class="text-sm font-semibold text-slate-700 mb-3 flex items-center gap-2">
              <span class="w-1 h-4 bg-purple-500 rounded"></span>
              调度配置
            </h4>
            <div class="bg-slate-50 rounded-lg p-4 space-y-3">
              <div class="grid grid-cols-2 gap-4">
                <div>
                  <p class="text-xs text-slate-500 mb-1">调度类型</p>
                  <p class="text-sm text-slate-900">{{ getScheduleType(task.schedule_type) }}</p>
                </div>
                <div>
                  <p class="text-xs text-slate-500 mb-1">调度值</p>
                  <p class="text-sm text-slate-900 font-mono">{{ task.schedule_value }}</p>
                </div>
              </div>
              <div>
                <p class="text-xs text-slate-500 mb-1">下次执行时间</p>
                <p class="text-sm text-slate-900">{{ task.next_run_time || '未安排' }}</p>
              </div>
            </div>
          </div>

          <!-- 时间信息 -->
          <div>
            <h4 class="text-sm font-semibold text-slate-700 mb-3 flex items-center gap-2">
              <span class="w-1 h-4 bg-orange-500 rounded"></span>
              时间信息
            </h4>
            <div class="grid grid-cols-2 gap-4 bg-slate-50 rounded-lg p-4">
              <div>
                <p class="text-xs text-slate-500 mb-1">创建时间</p>
                <p class="text-sm text-slate-900">{{ formatTime(task.created_at) }}</p>
              </div>
              <div>
                <p class="text-xs text-slate-500 mb-1">更新时间</p>
                <p class="text-sm text-slate-900">{{ formatTime(task.updated_at) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="px-6 py-4 border-t-2 border-slate-200 flex justify-end gap-3 flex-shrink-0">
        <button
          @click="$emit('update:modelValue', false)"
          class="px-4 py-2 bg-slate-100 hover:bg-slate-200 text-slate-700 font-medium rounded-lg transition-colors"
        >
          关闭
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits } from 'vue'
import type { Task } from '@/api/task'

const props = defineProps<{
  modelValue: boolean
  task: Task | null
}>()

defineEmits(['update:modelValue'])

const getConfigValue = (key: string, defaultValue: any = null) => {
  if (!props.task?.config) return defaultValue

  try {
    const config = typeof props.task.config === 'string'
      ? JSON.parse(props.task.config)
      : props.task.config
    return config[key] ?? defaultValue
  } catch {
    return defaultValue
  }
}

const formatJSON = (value: any) => {
  if (!value) return ''
  if (typeof value === 'string') return value
  try {
    return JSON.stringify(value, null, 2)
  } catch {
    return String(value)
  }
}

const getScheduleType = (type: string) => {
  const typeMap: Record<string, string> = {
    daily: '每天',
    weekly: '每周',
    monthly: '每月',
    hourly: '每小时',
    interval: '间隔',
    cron: 'Cron表达式'
  }
  return typeMap[type] || type
}

const formatTime = (time: any) => {
  if (!time) return '--'
  return new Date(time).toLocaleString('zh-CN')
}
</script>
