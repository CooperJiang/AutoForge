<template>
  <div class="http-tool-config">
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <el-form-item label="URL" prop="url">
        <el-input v-model="form.url" placeholder="https://example.com/api/endpoint" />
      </el-form-item>

      <el-form-item label="请求方法" prop="method">
        <el-select v-model="form.method" placeholder="选择请求方法">
          <el-option label="GET" value="GET" />
          <el-option label="POST" value="POST" />
          <el-option label="PUT" value="PUT" />
          <el-option label="DELETE" value="DELETE" />
          <el-option label="PATCH" value="PATCH" />
        </el-select>
      </el-form-item>

      <el-form-item label="Headers">
        <div class="key-value-list">
          <div v-for="(header, index) in form.headers" :key="index" class="key-value-item">
            <el-input
              v-model="header.key"
              placeholder="Key"
              style="width: 40%; margin-right: 10px"
            />
            <el-input
              v-model="header.value"
              placeholder="Value"
              style="width: 40%; margin-right: 10px"
            />
            <el-button type="danger" size="small" @click="removeHeader(index)" icon="Delete">
            </el-button>
          </div>
          <el-button type="primary" size="small" @click="addHeader" icon="Plus">
            添加Header
          </el-button>
        </div>
      </el-form-item>

      <el-form-item label="Query Params">
        <div class="key-value-list">
          <div v-for="(param, index) in form.params" :key="index" class="key-value-item">
            <el-input
              v-model="param.key"
              placeholder="Key"
              style="width: 40%; margin-right: 10px"
            />
            <el-input
              v-model="param.value"
              placeholder="Value"
              style="width: 40%; margin-right: 10px"
            />
            <el-button type="danger" size="small" @click="removeParam(index)" icon="Delete">
            </el-button>
          </div>
          <el-button type="primary" size="small" @click="addParam" icon="Plus">
            添加参数
          </el-button>
        </div>
      </el-form-item>

      <el-form-item label="Body" v-if="['POST', 'PUT', 'PATCH'].includes(form.method)">
        <el-input
          v-model="form.body"
          type="textarea"
          :rows="6"
          placeholder='{"key": "value"}'
        />
      </el-form-item>

      <el-form-item label="超时时间(秒)" prop="timeout">
        <el-input-number v-model="form.timeout" :min="1" :max="300" />
      </el-form-item>

      <el-form-item label="重试次数" prop="retry_times">
        <el-input-number v-model="form.retry_times" :min="0" :max="5" />
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, watch } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'

interface KeyValue {
  key: string
  value: string
}

interface HttpConfig {
  url: string
  method: string
  headers: KeyValue[]
  params: KeyValue[]
  body: string
  timeout: number
  retry_times: number
}

interface Props {
  modelValue?: Record<string, any>
}

interface Emits {
  (e: 'update:modelValue', value: Record<string, any>): void
  (e: 'validate', valid: boolean): void
}

const props = defineProps<Props>()
const emit = defineEmits<Emits>()

const formRef = ref<FormInstance>()

const form = reactive<HttpConfig>({
  url: '',
  method: 'GET',
  headers: [],
  params: [],
  body: '',
  timeout: 30,
  retry_times: 0,
})

// 初始化表单
if (props.modelValue) {
  Object.assign(form, {
    url: props.modelValue.url || '',
    method: props.modelValue.method || 'GET',
    headers: props.modelValue.headers || [],
    params: props.modelValue.params || [],
    body: props.modelValue.body || '',
    timeout: props.modelValue.timeout || 30,
    retry_times: props.modelValue.retry_times || 0,
  })
}

const rules = reactive<FormRules>({
  url: [
    { required: true, message: '请输入URL', trigger: 'blur' },
    { type: 'url', message: '请输入有效的URL', trigger: 'blur' },
  ],
  method: [{ required: true, message: '请选择请求方法', trigger: 'change' }],
  timeout: [
    { required: true, message: '请输入超时时间', trigger: 'blur' },
    { type: 'number', min: 1, max: 300, message: '超时时间必须在1-300秒之间', trigger: 'blur' },
  ],
  retry_times: [
    { required: true, message: '请输入重试次数', trigger: 'blur' },
    { type: 'number', min: 0, max: 5, message: '重试次数必须在0-5之间', trigger: 'blur' },
  ],
})

const addHeader = () => {
  form.headers.push({ key: '', value: '' })
}

const removeHeader = (index: number) => {
  form.headers.splice(index, 1)
}

const addParam = () => {
  form.params.push({ key: '', value: '' })
}

const removeParam = (index: number) => {
  form.params.splice(index, 1)
}

// 监听表单变化并向父组件发送
watch(
  form,
  () => {
    emit('update:modelValue', { ...form })
    // 触发验证
    formRef.value?.validate((valid) => {
      emit('validate', valid)
    })
  },
  { deep: true }
)

// 暴露验证方法
defineExpose({
  validate: () => formRef.value?.validate(),
})
</script>

<style scoped>
.http-tool-config {
  padding: 10px;
}

.key-value-list {
  width: 100%;
}

.key-value-item {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}
</style>
