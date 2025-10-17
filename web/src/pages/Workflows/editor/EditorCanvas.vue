<template>
  <div class="flex-1 relative" @drop="onDrop" @dragover.prevent @dragenter.prevent>
    <VueFlow
      v-if="isReady"
      v-model:nodes="nodes"
      v-model:edges="edges"
      :default-zoom="1"
      :min-zoom="0.2"
      :max-zoom="4"
      @node-click="$emit('node-click', $event)"
      @edge-click="$emit('edge-click', $event)"
      @connect="$emit('connect', $event)"
    >
      <Background variant="dots" pattern-color="#94a3b8" :gap="16" :size="1" />
      <Controls />

      <!-- Node Templates -->
      <template #node-tool="{ data }">
        <ToolNode :data="data" @delete="$emit('node-delete', data)" />
      </template>

      <template #node-external_trigger="{ data }">
        <ExternalTriggerNode :data="data" @delete="$emit('node-delete', data)" />
      </template>

      <template #node-trigger="{ data }">
        <TriggerNode :data="data" @delete="$emit('node-delete', data)" />
      </template>

      <template #node-condition="{ data }">
        <ConditionNode :data="data" @delete="$emit('node-delete', data)" />
      </template>

      <template #node-delay="{ data }">
        <DelayNode :data="data" @delete="$emit('node-delete', data)" />
      </template>

      <template #node-switch="{ data }">
        <SwitchNode :data="data" @delete="$emit('node-delete', data)" />
      </template>
    </VueFlow>

    <!-- Empty State -->
    <div
      v-if="showEmptyState"
      class="absolute inset-0 flex items-center justify-center pointer-events-none"
    >
      <div class="text-center text-text-placeholder">
        <Workflow class="w-16 h-16 mx-auto mb-4 opacity-50" />
        <p class="text-lg font-medium mb-2">从左侧拖拽工具开始构建工作流</p>
        <p class="text-sm">或点击工具添加到画布</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { VueFlow } from '@vue-flow/core'
import { Background } from '@vue-flow/background'
import { Controls } from '@vue-flow/controls'
import { Workflow } from 'lucide-vue-next'
import ToolNode from '../components/ToolNode.vue'
import ExternalTriggerNode from '../components/ExternalTriggerNode.vue'
import TriggerNode from '../components/TriggerNode.vue'
import ConditionNode from '../components/ConditionNode.vue'
import DelayNode from '../components/DelayNode.vue'
import SwitchNode from '../components/SwitchNode.vue'

import '@vue-flow/core/dist/style.css'
import '@vue-flow/core/dist/theme-default.css'
import '@vue-flow/controls/dist/style.css'

interface Props {
  nodes: any[]
  edges: any[]
  isReady: boolean
  showEmptyState: boolean
}

defineProps<Props>()

defineEmits<{
  'update:nodes': [nodes: any[]]
  'update:edges': [edges: any[]]
  'node-click': [event: any]
  'edge-click': [event: any]
  connect: [event: any]
  'node-delete': [node: any]
  drop: [event: DragEvent]
}>()

const onDrop = (event: DragEvent) => {
  // 父组件处理 drop 事件
  event.preventDefault()
  // Emit 事件让父组件处理
}
</script>

