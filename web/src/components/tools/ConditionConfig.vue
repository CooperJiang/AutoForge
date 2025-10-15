<template>
  <div class="space-y-4">
    <div class="bg-warning-light border-l-4 border-warning rounded-lg p-3">
      <p class="text-sm text-warning-text">
        <svg class="inline-block w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
          <path
            fill-rule="evenodd"
            d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z"
            clip-rule="evenodd"
          />
        </svg>
        根据上一个节点的执行结果进行条件判断，决定工作流的分支走向
      </p>
    </div>

    <div>
      <label class="block text-sm font-medium text-text-secondary mb-2">
        条件类型 <span class="text-red-500">*</span>
      </label>
      <BaseSelect
        v-model="localConfig.conditionType"
        :options="conditionTypeOptions"
        @update:model-value="emitUpdate"
      />
    </div>

    
    <div
      v-if="localConfig.conditionType === 'simple'"
      class="space-y-3 border-t border-border-primary pt-4"
    >
      <div>
        <label
          class="block text-sm font-medium text-text-secondary mb-2 flex items-center justify-between"
        >
          <span>检查字段 <span class="text-red-500">*</span></span>
          <button
            v-if="previousNodes && previousNodes.length > 0"
            type="button"
            @click="showFieldHelper = !showFieldHelper"
            class="text-xs text-primary hover:text-primary"
          >
            {{ showFieldHelper ? '隐藏' : '显示' }}字段助手
          </button>
        </label>

        
        <div
          v-if="showFieldHelper && previousNodes && previousNodes.length > 0"
          class="mb-2 p-3 bg-bg-hover rounded-lg border border-border-primary"
        >
          <div class="text-xs font-semibold text-text-secondary mb-2">可用字段（点击插入）：</div>
          <div class="space-y-2">
            <div
              v-for="node in previousNodes"
              :key="node.id"
              class="border border-border-primary rounded p-2"
            >
              <div class="font-semibold text-text-primary mb-1">{{ node.name }}</div>
              <div class="text-text-tertiary text-[10px] mb-2">ID: {{ node.id }}</div>

              
              <div class="flex flex-wrap gap-1">
                <button
                  v-for="field in getCommonFields(node.type, node.toolCode)"
                  :key="field.name"
                  type="button"
                  @click="insertField(field.name)"
                  class="px-1.5 py-0.5 bg-primary-light hover:bg-primary-light text-primary rounded text-[10px] font-mono"
                  :title="field.description"
                >
                  {{ field.name }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <BaseInput
          v-model="localConfig.field"
          placeholder="例如: status, code, success 或嵌套字段: data.user.name"
          @update:model-value="emitUpdate"
        />
        <p class="text-xs text-text-tertiary mt-1">
          从上一个节点的输出结果中读取该字段，支持嵌套字段（用 . 分隔）
        </p>
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          判断条件 <span class="text-red-500">*</span>
        </label>
        <BaseSelect
          v-model="localConfig.operator"
          :options="operatorOptions"
          @update:model-value="emitUpdate"
        />
      </div>

      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          期望值 <span class="text-red-500">*</span>
        </label>
        <BaseInput
          v-model="localConfig.value"
          placeholder="例如: 200, success, true"
          @update:model-value="emitUpdate"
        />
        <p class="text-xs text-text-tertiary mt-1">与检查字段进行比较的值</p>
      </div>

      
      <div class="bg-bg-hover rounded-lg p-3">
        <div class="text-xs font-semibold text-text-secondary mb-2">条件示例：</div>
        <div class="text-xs text-text-secondary space-y-1 font-mono">
          <div>• 检查HTTP状态码：<span class="text-primary">status == 200</span></div>
          <div>• 检查健康状态：<span class="text-primary">healthy == true</span></div>
          <div>• 检查错误信息：<span class="text-primary">error contains "timeout"</span></div>
        </div>
      </div>
    </div>

    
    <div
      v-if="localConfig.conditionType === 'expression'"
      class="space-y-3 border-t border-border-primary pt-4"
    >
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          条件表达式 <span class="text-red-500">*</span>
        </label>
        <textarea
          v-model="localConfig.expression"
          @input="emitUpdate"
          class="w-full px-3 py-1.5 text-sm text-text-primary bg-bg-primary border-2 border-border-primary rounded-md transition-all duration-200 focus:border-border-focus focus:ring-2 focus:ring-primary-light focus:outline-none hover:border-border-secondary placeholder:text-text-placeholder font-mono"
          rows="4"
          placeholder="status == 200 && success == true"
        />
        <p class="text-xs text-text-tertiary mt-1">支持 &&（与）、||（或）、!（非）等逻辑运算符</p>
      </div>

      <div class="bg-bg-hover rounded-lg p-3">
        <div class="text-xs font-semibold text-text-secondary mb-2">表达式示例：</div>
        <div class="text-xs text-text-secondary space-y-1 font-mono">
          <div>• <span class="text-primary">status == 200 && success == true</span></div>
          <div>• <span class="text-primary">code >= 200 && code &lt; 300</span></div>
          <div>• <span class="text-primary">error == null || error == ""</span></div>
        </div>
      </div>
    </div>

    
    <div
      v-if="localConfig.conditionType === 'script'"
      class="space-y-3 border-t border-border-primary pt-4"
    >
      <div>
        <label class="block text-sm font-medium text-text-secondary mb-2">
          JavaScript 脚本 <span class="text-red-500">*</span>
        </label>
        <textarea
          v-model="localConfig.script"
          @input="emitUpdate"
          class="w-full px-3 py-1.5 text-sm text-text-primary bg-bg-primary border-2 border-border-primary rounded-md transition-all duration-200 focus:border-border-focus focus:ring-2 focus:ring-primary-light focus:outline-none hover:border-border-secondary placeholder:text-text-placeholder font-mono"
          rows="8"
          placeholder="// 可使用 input 变量访问上一节点的输出&#10;// 返回 true 或 false&#10;&#10;if (input.status === 200) {&#10;  return true;&#10;}&#10;return false;"
        />
        <p class="text-xs text-text-tertiary mt-1">
          通过
          <code class="bg-bg-tertiary px-1 rounded">input</code>
          变量访问上一节点的输出，必须返回布尔值
        </p>
      </div>
    </div>

    
    <div class="border-t border-border-primary pt-4">
      <div class="text-sm font-semibold text-text-secondary mb-3">分支说明</div>
      <div class="grid grid-cols-2 gap-3">
        <div class="bg-success-light border border-success rounded-lg p-3">
          <div class="flex items-center gap-2 mb-1">
            <div class="w-3 h-3 rounded-full bg-success"></div>
            <span class="text-sm font-medium text-success-text">True 分支</span>
          </div>
          <p class="text-xs text-success-text">条件为真时执行此分支</p>
        </div>
        <div class="bg-error-light border border-error rounded-lg p-3">
          <div class="flex items-center gap-2 mb-1">
            <div class="w-3 h-3 rounded-full bg-error"></div>
            <span class="text-sm font-medium text-error-text">False 分支</span>
          </div>
          <p class="text-xs text-error-text">条件为假时执行此分支</p>
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
  previousNodes?: Array<{ id: string; name: string; type: string; toolCode?: string }>
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
  ...props.config,
})

const conditionTypeOptions = [
  { label: '简单条件', value: 'simple' },
  { label: '表达式', value: 'expression' },
  { label: 'JavaScript 脚本', value: 'script' },
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
  { label: '不为空', value: 'is_not_empty' },
]

const showFieldHelper = ref(false)

// 获取节点的常见字段
const getCommonFields = (nodeType: string, toolCode?: string) => {
  // HTTP 请求的常见返回字段
  if (toolCode === 'http_request') {
    return [
      { name: 'status', description: 'HTTP 状态码' },
      { name: 'message', description: '返回消息' },
      { name: 'success', description: '成功标识' },
      { name: 'data', description: '返回数据' },
      { name: 'code', description: '业务代码' },
      { name: 'error', description: '错误信息' },
    ]
  }

  // 邮件发送的常见返回字段
  if (toolCode === 'email_sender') {
    return [
      { name: 'success', description: '发送成功' },
      { name: 'message_id', description: '消息ID' },
      { name: 'error', description: '错误信息' },
    ]
  }

  // 健康检查的常见返回字段
  if (toolCode === 'health_checker') {
    return [
      { name: 'healthy', description: '健康状态' },
      { name: 'status', description: '检查状态' },
      { name: 'message', description: '状态消息' },
      { name: 'latency', description: '响应延迟' },
    ]
  }

  // 条件节点的返回字段
  if (nodeType === 'condition') {
    return [
      { name: 'result', description: '条件判断结果 (true/false)' },
      { name: 'message', description: '判断说明' },
    ]
  }

  // 通用字段
  return [
    { name: 'success', description: '成功标识' },
    { name: 'message', description: '消息' },
    { name: 'data', description: '数据' },
    { name: 'error', description: '错误' },
  ]
}

// 插入字段
const insertField = (fieldName: string) => {
  localConfig.value.field = fieldName
  emitUpdate()
}

watch(
  () => props.config,
  (newVal) => {
    localConfig.value = { ...localConfig.value, ...newVal }
  },
  { deep: true }
)

const emitUpdate = () => {
  emit('update:config', localConfig.value)
}
</script>
