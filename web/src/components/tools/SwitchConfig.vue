<template>
  <div class="space-y-4">
    <div class="bg-accent/10 border-l-4 border-accent p-3">
      <p class="text-sm text-accent">
        <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
          <path
            fill-rule="evenodd"
            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
            clip-rule="evenodd"
          />
        </svg>
        根据字段值匹配不同的分支，类似编程语言中的 switch 语句
      </p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        检查字段 <span class="text-red-500">*</span>
      </label>
      <BaseInput
        v-model="localConfig.field"
        placeholder="例如: status, type, level"
        @update:model-value="emitUpdate"
      />
      <p class="text-xs text-text-tertiary mt-1">从上一个节点的输出结果中读取该字段</p>
    </div>

    <!-- 分支配置 -->
    <div class="border-t border-border-primary pt-4">
      <div class="flex items-center justify-between mb-3">
        <label class="block text-sm font-medium text-text-secondary">
          分支条件 <span class="text-red-500">*</span>
        </label>
        <BaseButton size="xs" @click="addCase"> + 添加分支 </BaseButton>
      </div>

      <div class="space-y-3">
        <div
          v-for="(caseItem, index) in localConfig.cases"
          :key="index"
          class="bg-bg-hover rounded-lg p-3 border border-border-primary"
        >
          <div class="flex items-start justify-between mb-2">
            <div class="flex items-center gap-2">
              <div
                class="w-6 h-6 rounded flex items-center justify-center text-xs font-bold text-white"
                :style="{ backgroundColor: getBranchColor(index) }"
              >
                {{ index + 1 }}
              </div>
              <span class="text-xs font-medium text-text-secondary">分支 {{ index + 1 }}</span>
            </div>
            <button type="button" @click="removeCase(index)" class="text-error hover:text-rose-700">
              <X class="w-4 h-4" />
            </button>
          </div>

          <div class="space-y-2">
            <div>
              <BaseInput
                v-model="caseItem.label"
                placeholder="分支标签（显示在连接点上）"
                size="sm"
                @update:model-value="emitUpdate"
              />
            </div>
            <div>
              <BaseInput
                v-model="caseItem.value"
                placeholder="匹配值"
                size="sm"
                @update:model-value="emitUpdate"
              />
              <p class="text-xs text-text-tertiary mt-1">当字段值等于此值时，执行该分支</p>
            </div>
          </div>
        </div>

        <!-- Default 分支说明 -->
        <div class="bg-bg-tertiary rounded-lg p-3 border border-slate-300">
          <div class="flex items-center gap-2 mb-1">
            <div
              class="w-6 h-6 rounded flex items-center justify-center text-xs font-bold bg-bg-hover0 text-white"
            >
              D
            </div>
            <span class="text-xs font-medium text-text-secondary">默认分支 (Default)</span>
          </div>
          <p class="text-xs text-text-secondary ml-8">当所有条件都不匹配时，执行此分支</p>
        </div>
      </div>
    </div>

    <!-- 示例 -->
    <div class="bg-bg-hover rounded-lg p-3">
      <div class="text-xs font-semibold text-text-secondary mb-2">示例：根据HTTP状态码分支</div>
      <div class="text-xs text-text-secondary space-y-1 font-mono">
        <div>• 检查字段: <span class="text-indigo-600">status</span></div>
        <div>• 分支1: 值=<span class="text-indigo-600">200</span> → 成功处理</div>
        <div>• 分支2: 值=<span class="text-indigo-600">404</span> → 资源不存在</div>
        <div>• 分支3: 值=<span class="text-indigo-600">500</span> → 服务器错误</div>
        <div>• 默认分支 → 其他情况处理</div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { X } from 'lucide-vue-next'
import BaseInput from '@/components/BaseInput'
import BaseButton from '@/components/BaseButton'

interface Props {
  config: Record<string, any>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:config': [config: Record<string, any>]
}>()

const localConfig = ref({
  field: '',
  cases: [] as Array<{ label: string; value: string }>,
  ...props.config,
})

// 确保cases是数组
if (!Array.isArray(localConfig.value.cases)) {
  localConfig.value.cases = []
}

watch(
  () => props.config,
  (newVal) => {
    localConfig.value = {
      ...localConfig.value,
      ...newVal,
      cases: Array.isArray(newVal.cases) ? newVal.cases : [],
    }
  },
  { deep: true }
)

const addCase = () => {
  localConfig.value.cases.push({
    label: `Case ${localConfig.value.cases.length + 1}`,
    value: '',
  })
  emitUpdate()
}

const removeCase = (index: number) => {
  localConfig.value.cases.splice(index, 1)
  emitUpdate()
}

const getBranchColor = (index: number) => {
  const colors = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981']
  return colors[index % colors.length]
}

const emitUpdate = () => {
  emit('update:config', localConfig.value)
}
</script>
