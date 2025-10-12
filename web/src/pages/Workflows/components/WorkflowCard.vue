<template>
  <div class="group bg-white rounded-lg shadow-sm hover:shadow-xl transition-all duration-300 p-4 border border-slate-200 hover:border-blue-400">
    <!-- 头部 -->
    <div class="flex items-start justify-between mb-3">
      <div class="flex-1 min-w-0">
        <h3 class="text-sm font-semibold text-slate-900 truncate mb-1">
          {{ workflow.name }}
        </h3>
        <p class="text-xs text-slate-600 line-clamp-2">
          {{ workflow.description || '暂无描述' }}
        </p>
      </div>
      <div
        :class="[
          'flex-shrink-0 px-2 py-0.5 rounded text-xs font-medium ml-2',
          workflow.enabled ? 'bg-green-50 text-green-700' : 'bg-slate-100 text-slate-600'
        ]"
      >
        {{ workflow.enabled ? '已启用' : '已禁用' }}
      </div>
    </div>

    <!-- 统计信息 -->
    <div class="flex items-center gap-3 text-xs text-slate-500 mb-3">
      <span class="flex items-center gap-1">
        <Box class="w-3.5 h-3.5" />
        {{ workflow.nodes.length }} 个节点
      </span>
      <span class="flex items-center gap-1">
        <GitBranch class="w-3.5 h-3.5" />
        {{ workflow.edges.length }} 个连接
      </span>
    </div>

    <!-- 操作按钮 -->
    <div class="space-y-2 pt-3 border-t border-slate-100">
      <!-- 第一行：主要操作 -->
      <div class="flex items-center gap-2">
        <BaseButton size="sm" variant="ghost" @click="$emit('edit', workflow)" class="flex-1">
          <Edit3 class="w-3.5 h-3.5 mr-1" />
          编辑
        </BaseButton>
        <BaseButton size="sm" variant="ghost" @click="$emit('executions', workflow)" class="flex-1">
          <History class="w-3.5 h-3.5 mr-1" />
          历史
        </BaseButton>
        <BaseButton
          size="sm"
          variant="primary"
          @click="$emit('execute', workflow)"
          class="flex-1"
          :disabled="workflow.nodes.length === 0"
        >
          <Play class="w-3.5 h-3.5 mr-1" />
          执行
        </BaseButton>
      </div>

      <!-- 第二行：状态和删除 -->
      <div class="flex items-center gap-2">
        <BaseButton
          size="sm"
          :variant="workflow.enabled ? 'success' : 'secondary'"
          @click="$emit('toggle', workflow)"
          class="flex-1"
        >
          <Power class="w-3.5 h-3.5 mr-1" />
          {{ workflow.enabled ? '已启用' : '已禁用' }}
        </BaseButton>
        <BaseButton size="sm" variant="danger" @click="handleDelete">
          <Trash2 class="w-3.5 h-3.5" />
        </BaseButton>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Box, GitBranch, Edit3, Trash2, History, Play, Power } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import type { Workflow } from '@/types/workflow'

interface Props {
  workflow: Workflow
}

const props = defineProps<Props>()

const emit = defineEmits<{
  edit: [workflow: Workflow]
  executions: [workflow: Workflow]
  execute: [workflow: Workflow]
  delete: [workflow: Workflow]
  toggle: [workflow: Workflow]
}>()

const handleDelete = () => {
  if (confirm('确定要删除这个工作流吗？')) {
    emit('delete', props.workflow)
  }
}
</script>
