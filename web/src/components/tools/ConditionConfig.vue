<template>
  <div class="space-y-4">
    <div class="bg-amber-50 border-l-4 border-amber-400 p-3">
      <p class="text-sm text-amber-700">
        <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
        </svg>
        根据上一个节点的执行结果进行条件判断，决定工作流的分支走向
      </p>
    </div>

    <div>
      <label class="block text-sm font-medium text-slate-700 mb-2">
        条件类型 <span class="text-red-500">*</span>
      </label>
      <BaseSelect
        v-model="localConfig.conditionType"
        :options="conditionTypeOptions"
        @update:model-value="emitUpdate"
      />
    </div>

    <!-- 简单条件 -->
    <div v-if="localConfig.conditionType === 'simple'" class="space-y-3 border-t border-slate-200 pt-4">
      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">
          检查字段 <span class="text-red-500">*</span>
        </label>
        <BaseInput
          v-model="localConfig.field"
          placeholder="例如: status, code, success"
          @update:model-value="emitUpdate"
        />
        <p class="text-xs text-slate-500 mt-1">从上一个节点的输出结果中读取该字段</p>
      </div>

      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">
          判断条件 <span class="text-red-500">*</span>
        </label>
        <BaseSelect
          v-model="localConfig.operator"
          :options="operatorOptions"
          @update:model-value="emitUpdate"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">
          期望值 <span class="text-red-500">*</span>
        </label>
        <BaseInput
          v-model="localConfig.value"
          placeholder="例如: 200, success, true"
          @update:model-value="emitUpdate"
        />
        <p class="text-xs text-slate-500 mt-1">与检查字段进行比较的值</p>
      </div>

      <!-- 示例 -->
      <div class="bg-slate-50 rounded-lg p-3">
        <div class="text-xs font-semibold text-slate-700 mb-2">条件示例：</div>
        <div class="text-xs text-slate-600 space-y-1 font-mono">
          <div>• 检查HTTP状态码：<span class="text-blue-600">status == 200</span></div>
          <div>• 检查健康状态：<span class="text-blue-600">healthy == true</span></div>
          <div>• 检查错误信息：<span class="text-blue-600">error contains "timeout"</span></div>
        </div>
      </div>
    </div>

    <!-- 表达式条件 -->
    <div v-if="localConfig.conditionType === 'expression'" class="space-y-3 border-t border-slate-200 pt-4">
      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">
          条件表达式 <span class="text-red-500">*</span>
        </label>
        <textarea
          v-model="localConfig.expression"
          @input="emitUpdate"
          class="w-full px-3 py-2 border-2 border-slate-200 rounded-lg focus:outline-none focus:border-amber-500 font-mono text-sm"
          rows="4"
          placeholder="status == 200 && success == true"
        />
        <p class="text-xs text-slate-500 mt-1">支持 &&（与）、||（或）、!（非）等逻辑运算符</p>
      </div>

      <div class="bg-slate-50 rounded-lg p-3">
        <div class="text-xs font-semibold text-slate-700 mb-2">表达式示例：</div>
        <div class="text-xs text-slate-600 space-y-1 font-mono">
          <div>• <span class="text-blue-600">status == 200 && success == true</span></div>
          <div>• <span class="text-blue-600">code >= 200 && code < 300</span></div>
          <div>• <span class="text-blue-600">error == null || error == ""</span></div>
        </div>
      </div>
    </div>

    <!-- JavaScript脚本 -->
    <div v-if="localConfig.conditionType === 'script'" class="space-y-3 border-t border-slate-200 pt-4">
      <div>
        <label class="block text-sm font-medium text-slate-700 mb-2">
          JavaScript 脚本 <span class="text-red-500">*</span>
        </label>
        <textarea
          v-model="localConfig.script"
          @input="emitUpdate"
          class="w-full px-3 py-2 border-2 border-slate-200 rounded-lg focus:outline-none focus:border-amber-500 font-mono text-sm"
          rows="8"
          placeholder="// 可使用 input 变量访问上一节点的输出&#10;// 返回 true 或 false&#10;&#10;if (input.status === 200) {&#10;  return true;&#10;}&#10;return false;"
        />
        <p class="text-xs text-slate-500 mt-1">通过 <code class="bg-slate-100 px-1 rounded">input</code> 变量访问上一节点的输出，必须返回布尔值</p>
      </div>
    </div>

    <!-- 分支说明 -->
    <div class="border-t border-slate-200 pt-4">
      <div class="text-sm font-semibold text-slate-700 mb-3">分支说明</div>
      <div class="grid grid-cols-2 gap-3">
        <div class="bg-emerald-50 border border-emerald-200 rounded-lg p-3">
          <div class="flex items-center gap-2 mb-1">
            <div class="w-3 h-3 rounded-full bg-emerald-500"></div>
            <span class="text-sm font-medium text-emerald-700">True 分支</span>
          </div>
          <p class="text-xs text-emerald-600">条件为真时执行此分支</p>
        </div>
        <div class="bg-rose-50 border border-rose-200 rounded-lg p-3">
          <div class="flex items-center gap-2 mb-1">
            <div class="w-3 h-3 rounded-full bg-rose-500"></div>
            <span class="text-sm font-medium text-rose-700">False 分支</span>
          </div>
          <p class="text-xs text-rose-600">条件为假时执行此分支</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import BaseInput from '@/components/BaseInput'
import BaseSelect from '@/components/BaseSelect'

interface Props {
  config: Record<string, any>
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:config': [config: Record<string, any>]
}>()

const localConfig = ref({
  conditionType: 'simple',
  field: '',
  operator: 'equals',
  value: '',
  expression: '',
  script: '',
  ...props.config
})

const conditionTypeOptions = [
  { label: '简单条件', value: 'simple' },
  { label: '表达式', value: 'expression' },
  { label: 'JavaScript 脚本', value: 'script' }
]

const operatorOptions = [
  { label: '等于 (==)', value: 'equals' },
  { label: '不等于 (!=)', value: 'not_equals' },
  { label: '大于 (>)', value: 'greater_than' },
  { label: '小于 (<)', value: 'less_than' },
  { label: '大于等于 (>=)', value: 'greater_or_equal' },
  { label: '小于等于 (<=)', value: 'less_or_equal' },
  { label: '包含', value: 'contains' },
  { label: '不包含', value: 'not_contains' },
  { label: '以...开始', value: 'starts_with' },
  { label: '以...结束', value: 'ends_with' },
  { label: '为空', value: 'is_empty' },
  { label: '不为空', value: 'is_not_empty' }
]

watch(() => props.config, (newVal) => {
  localConfig.value = { ...localConfig.value, ...newVal }
}, { deep: true })

const emitUpdate = () => {
  emit('update:config', localConfig.value)
}
</script>
