<template>
  <div class="json-viewer relative group">
    <button
      @click="copyToClipboard"
      class="absolute top-1.5 right-1.5 opacity-0 group-hover:opacity-100 transition-all bg-white/90 hover:bg-white text-slate-600 hover:text-slate-900 px-1.5 py-1 rounded shadow-sm border border-slate-200 text-xs flex items-center gap-1 z-10"
      :class="{ '!bg-emerald-50 !text-emerald-700 !border-emerald-200': copied }"
    >
      <Copy v-if="!copied" :size="12" />
      <Check v-else :size="12" />
    </button>
    <pre v-if="isJson" class="bg-slate-900 p-3 rounded font-mono text-xs overflow-auto whitespace-pre-wrap break-words" style="max-height: 400px" v-html="highlightedJson"></pre>
    <div v-else class="bg-slate-50 border-2 border-slate-200 p-3 rounded font-mono text-xs text-slate-700 overflow-auto whitespace-pre-wrap break-words" style="max-height: 400px">{{ content }}</div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { Copy, Check } from 'lucide-vue-next'
import { message } from '@/utils/message'

interface Props {
  content: string
}

const props = defineProps<Props>()

const copied = ref(false)

const isJson = computed(() => {
  try {
    JSON.parse(props.content)
    return true
  } catch {
    return false
  }
})

const formattedJson = computed(() => {
  if (!isJson.value) return props.content

  try {
    const obj = JSON.parse(props.content)
    return JSON.stringify(obj, null, 2)
  } catch {
    return props.content
  }
})

const highlightedJson = computed(() => {
  if (!isJson.value) return props.content

  const json = formattedJson.value
  return json
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"([^"]+)":/g, '<span class="text-blue-400">"$1"</span>:') // 键名
    .replace(/: "([^"]*)"/g, ': <span class="text-green-400">"$1"</span>') // 字符串值
    .replace(/: (\d+)/g, ': <span class="text-amber-400">$1</span>') // 数字
    .replace(/: (true|false)/g, ': <span class="text-purple-400">$1</span>') // 布尔值
    .replace(/: (null)/g, ': <span class="text-slate-500">$1</span>') // null
    .replace(/^(\s*)([\{\}\[\]])/gm, '$1<span class="text-slate-400">$2</span>') // 括号
})

const copyToClipboard = async () => {
  try {
    const textToCopy = isJson.value ? formattedJson.value : props.content
    await navigator.clipboard.writeText(textToCopy)
    copied.value = true
    message.success('已复制到剪贴板')
    setTimeout(() => {
      copied.value = false
    }, 2000)
  } catch (error) {
    message.error('复制失败')
  }
}
</script>

<style scoped>
.json-viewer :deep(pre) {
  margin: 0;
  padding: 0;
}

.json-viewer :deep(span) {
  white-space: pre-wrap;
  word-break: break-word;
}
</style>
