<template>
  <Drawer v-model="isOpen" title="环境变量管理" size="md" @close="handleClose">
    <div class="space-y-4">
      <div class="bg-primary-light border border-primary rounded-lg p-3 text-sm">
        <div class="font-semibold text-text-primary mb-1">💡 什么是环境变量？</div>
        <div class="text-primary text-xs space-y-1">
          <p>环境变量可以在工作流的任何节点中引用，用于存储 API 密钥、配置参数等敏感信息。</p>
          <p>
            使用方式：<code class="px-1 py-0.5 bg-primary-light rounded"
              >&#123;&#123;env.VARIABLE_NAME&#125;&#125;</code
            >
          </p>
        </div>
      </div>

      <div class="space-y-2">
        <div
          v-for="(envVar, index) in localEnvVars"
          :key="index"
          class="border-2 border-border-primary rounded-lg p-3 space-y-2"
        >
          <div class="grid grid-cols-2 gap-2">
            <div>
              <label class="block text-xs font-medium text-text-secondary mb-1">变量名</label>
              <BaseInput v-model="envVar.key" placeholder="API_KEY" class="font-mono text-sm" />
            </div>
            <div>
              <label class="block text-xs font-medium text-text-secondary mb-1">描述（可选）</label>
              <BaseInput v-model="envVar.description" placeholder="API密钥" />
            </div>
          </div>
          <div>
            <label class="block text-xs font-medium text-text-secondary mb-1">值</label>
            <div class="relative">
              <BaseInput
                v-model="envVar.value"
                :type="!envVar.encrypted || showPassword[index] ? 'text' : 'password'"
                placeholder="变量值"
                input-class="pr-10 font-mono"
              />
              <button
                v-if="envVar.encrypted"
                type="button"
                @click="togglePasswordVisibility(index)"
                class="absolute right-2 top-1/2 -translate-y-1/2 text-text-secondary hover:text-text-primary"
                :title="showPassword[index] ? '隐藏' : '显示'"
              >
                <Eye v-if="showPassword[index]" class="w-4 h-4" />
                <EyeOff v-else class="w-4 h-4" />
              </button>
            </div>
          </div>
          <div class="flex items-center justify-between pt-2">
            <div class="flex items-center gap-2">
              <BaseCheckbox v-model="envVar.encrypted" label="敏感信息（加密存储）" />
            </div>
            <button
              type="button"
              @click="removeEnvVar(index)"
              class="text-xs text-red-600 hover:text-red-700 font-medium"
            >
              删除
            </button>
          </div>
        </div>

        <button
          type="button"
          @click="addEnvVar"
          class="w-full py-3 text-sm text-text-secondary border-2 border-dashed border-slate-300 rounded-lg hover:border-green-500 hover:text-green-600 transition-colors"
        >
          + 添加环境变量
        </button>
      </div>

      <div class="flex gap-2 pt-4 border-t border-border-primary">
        <BaseButton size="sm" variant="ghost" class="flex-1" @click="handleClose">
          取消
        </BaseButton>
        <BaseButton size="sm" class="flex-1" @click="handleSave"> 保存 </BaseButton>
      </div>
    </div>
  </Drawer>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { Eye, EyeOff } from 'lucide-vue-next'
import Drawer from '@/components/Drawer'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import BaseCheckbox from '@/components/BaseCheckbox/index.vue'
import type { WorkflowEnvVar } from '@/types/workflow'
import { message } from '@/utils/message'

interface Props {
  modelValue: boolean
  envVars: WorkflowEnvVar[]
}

const props = defineProps<Props>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  'update:envVars': [envVars: WorkflowEnvVar[]]
}>()

const isOpen = ref(props.modelValue)
const localEnvVars = ref<WorkflowEnvVar[]>([])
const showPassword = ref<Record<number, boolean>>({})

watch(
  () => props.modelValue,
  (val) => {
    isOpen.value = val
    if (val) {
      // 深拷贝环境变量
      localEnvVars.value = JSON.parse(JSON.stringify(props.envVars))
      // 重置密码显示状态
      showPassword.value = {}
    }
  }
)

watch(isOpen, (val) => {
  emit('update:modelValue', val)
})

const togglePasswordVisibility = (index: number) => {
  showPassword.value[index] = !showPassword.value[index]
}

const addEnvVar = () => {
  localEnvVars.value.push({
    key: '',
    value: '',
    description: '',
    encrypted: false,
  })
}

const removeEnvVar = (index: number) => {
  localEnvVars.value.splice(index, 1)
  // 清除对应的密码显示状态
  delete showPassword.value[index]
}

const handleClose = () => {
  isOpen.value = false
}

const handleSave = () => {
  // 验证
  for (const envVar of localEnvVars.value) {
    if (!envVar.key.trim()) {
      message.error('请填写变量名')
      return
    }
    if (!envVar.value.trim()) {
      message.error(`请填写变量 ${envVar.key} 的值`)
      return
    }
    // 验证变量名格式（字母、数字、下划线）
    if (!/^[A-Z_][A-Z0-9_]*$/i.test(envVar.key)) {
      message.error(`变量名 ${envVar.key} 格式不正确，只能包含字母、数字和下划线，且不能以数字开头`)
      return
    }
  }

  // 检查重复
  const keys = localEnvVars.value.map((v) => v.key)
  const duplicates = keys.filter((key, index) => keys.indexOf(key) !== index)
  if (duplicates.length > 0) {
    message.error(`变量名重复：${duplicates[0]}`)
    return
  }

  emit('update:envVars', localEnvVars.value)
  message.success('环境变量已保存')
  handleClose()
}
</script>
