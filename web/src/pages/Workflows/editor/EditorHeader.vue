<template>
  <div
    class="bg-bg-elevated border-b border-border-primary px-6 py-3 flex items-center justify-between flex-shrink-0"
  >
    <!-- Left Section -->
    <div class="flex items-center gap-4">
      <BaseButton size="sm" variant="ghost" @click="$emit('back')" class="shrink-0">
        <ArrowLeft class="w-4 h-4" />
      </BaseButton>
      <div
        class="input-wrapper flex items-center gap-2 px-2.5 py-1 rounded-md bg-bg-hover border border-border-primary hover:border-slate-300 transition-all duration-200"
      >
        <Workflow class="w-3.5 h-3.5 text-text-placeholder shrink-0" />
        <input
          :value="workflowName"
          type="text"
          placeholder="工作流名称"
          class="w-32 bg-transparent text-xs font-medium text-text-primary placeholder:text-text-placeholder"
          style="border: none; outline: none; box-shadow: none"
          @input="$emit('update:workflow-name', ($event.target as HTMLInputElement).value)"
          @focus="($event.target as HTMLInputElement).parentElement?.classList.add('input-focused')"
          @blur="($event.target as HTMLInputElement).parentElement?.classList.remove('input-focused')"
        />
      </div>
    </div>

    <!-- Right Section -->
    <div class="flex items-center gap-3">
      <!-- Enable/Disable Toggle -->
      <div
        :class="[
          'px-3 py-1.5 rounded-full text-xs font-medium flex items-center gap-1.5 transition-colors border',
          !hasWorkflowId
            ? 'bg-slate-500/10 text-slate-600 dark:text-slate-400 border-slate-500/20 opacity-50 cursor-not-allowed'
            : enabled
              ? 'bg-green-500/10 text-green-600 dark:text-green-400 border-green-500/20 hover:bg-green-500/20 cursor-pointer'
              : 'bg-slate-500/10 text-slate-600 dark:text-slate-400 border-slate-500/20 hover:bg-slate-500/20 cursor-pointer',
        ]"
        @click="hasWorkflowId && $emit('toggle-enabled')"
        :title="!hasWorkflowId ? '请先保存工作流后才能启用/禁用' : ''"
      >
        <Power :class="['w-3.5 h-3.5', enabled && 'animate-pulse']" />
        {{ enabled ? '已启用' : '已禁用' }}
      </div>

      <div class="h-6 w-px bg-bg-tertiary"></div>

      <!-- Toolbar Buttons -->
      <Tooltip text="API 设置" position="bottom">
        <BaseButton size="sm" variant="ghost" @click="$emit('open-api-settings')">
          <Globe class="w-4 h-4" />
        </BaseButton>
      </Tooltip>
      <Tooltip text="环境变量配置" position="bottom">
        <BaseButton size="sm" variant="ghost" @click="$emit('open-env-manager')">
          <Settings class="w-4 h-4" />
        </BaseButton>
      </Tooltip>
      <Tooltip text="导入工作流" position="bottom">
        <BaseButton size="sm" variant="ghost" @click="$emit('open-import')">
          <Upload class="w-4 h-4" />
        </BaseButton>
      </Tooltip>
      <Tooltip text="导出工作流" position="bottom">
        <BaseButton size="sm" variant="ghost" @click="$emit('export')">
          <Download class="w-4 h-4" />
        </BaseButton>
      </Tooltip>
      <Tooltip text="清除本地草稿" position="bottom">
        <BaseButton size="sm" variant="ghost" @click="$emit('clear-draft')">
          <Trash2 class="w-4 h-4" />
        </BaseButton>
      </Tooltip>
      <Tooltip v-if="isAdmin && hasWorkflowId" text="发布为模板" position="bottom">
        <BaseButton size="sm" variant="ghost" @click="$emit('open-publish')">
          <Package class="w-4 h-4" />
        </BaseButton>
      </Tooltip>

      <div class="h-6 w-px bg-bg-tertiary"></div>

      <!-- Action Buttons -->
      <BaseButton
        size="sm"
        variant="secondary"
        @click="$emit('execute')"
        :disabled="!canExecute"
      >
        <Play class="w-4 h-4 mr-1.5" />
        执行
      </BaseButton>
      <BaseButton size="sm" @click="$emit('save')">
        <Save class="w-4 h-4 mr-1.5" />
        保存
      </BaseButton>
    </div>
  </div>
</template>

<script setup lang="ts">
import {
  ArrowLeft,
  Workflow,
  Power,
  Globe,
  Settings,
  Upload,
  Download,
  Trash2,
  Package,
  Play,
  Save,
} from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import Tooltip from '@/components/Tooltip'

interface Props {
  workflowName: string
  enabled: boolean
  hasWorkflowId: boolean
  isAdmin: boolean
  canExecute: boolean
}

defineProps<Props>()

defineEmits<{
  'update:workflow-name': [value: string]
  back: []
  'toggle-enabled': []
  'open-api-settings': []
  'open-env-manager': []
  'open-import': []
  export: []
  'clear-draft': []
  'open-publish': []
  execute: []
  save: []
}>()
</script>

<style scoped>
.input-wrapper.input-focused {
  @apply border-primary;
}
</style>

