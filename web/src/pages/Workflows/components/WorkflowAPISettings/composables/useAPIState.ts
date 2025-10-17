import { ref, computed, watch } from 'vue'
import type { Tab, CodeExample } from '../types'

/**
 * API 设置状态管理
 */
export function useAPIState(props: any) {
  const isOpen = ref(false)
  const activeTab = ref('overview')
  const apiEnabled = ref(false)
  const apiKey = ref('')
  const showApiKey = ref(false)
  const timeout = ref(30)
  const webhookURL = ref('')
  const testMode = ref('sync')
  const testParams = ref('{}')
  const testing = ref(false)
  const testResult = ref('')
  const selectedLanguage = ref('curl')

  // Tab列表
  const tabList: Tab[] = [
    { key: 'overview', label: '概览' },
    { key: 'examples', label: '调用示例' },
    { key: 'test', label: '在线测试' },
    { key: 'settings', label: '高级设置' },
  ]

  // API 端点
  const apiEndpoint = computed(() => {
    if (!props.workflow?.id) return ''
    const baseURL = window.location.origin
    return `${baseURL}/api/v1/workflows/${props.workflow.id}/execute`
  })

  // 显示的 API Key
  const displayApiKey = computed(() => {
    return apiKey.value || '未生成'
  })

  // 检查是否有外部触发节点
  const hasExternalTrigger = computed(() => {
    if (!props.workflow?.nodes) return false
    return props.workflow.nodes.some((n: any) => n.type === 'external_trigger')
  })

  // 外部参数列表
  const externalParams = computed(() => {
    if (!hasExternalTrigger.value) return []
    const triggerNode = props.workflow.nodes.find((n: any) => n.type === 'external_trigger')
    return triggerNode?.config?.params || []
  })

  // 代码示例
  const codeExamples: CodeExample[] = [
    {
      label: 'cURL',
      value: 'curl',
      code: '',
    },
    {
      label: 'JavaScript',
      value: 'javascript',
      code: '',
    },
    {
      label: 'Python',
      value: 'python',
      code: '',
    },
    {
      label: 'Go',
      value: 'go',
      code: '',
    },
  ]

  // 当前代码示例
  const currentCodeExample = computed(() => {
    const example = codeExamples.find((e) => e.value === selectedLanguage.value)
    if (!example) return ''

    const endpoint = apiEndpoint.value
    const key = apiKey.value || 'YOUR_API_KEY'
    const params = hasExternalTrigger.value
      ? JSON.stringify(
          externalParams.value.reduce((acc: any, p: any) => {
            acc[p.name] = `示例${p.name}`
            return acc
          }, {}),
          null,
          2
        )
      : '{}'

    switch (selectedLanguage.value) {
      case 'curl':
        return `curl -X POST "${endpoint}" \\
  -H "Content-Type: application/json" \\
  -H "X-API-Key: ${key}" \\
  -d '${params.replace(/\n/g, '\n  ')}'`

      case 'javascript':
        return `const response = await fetch("${endpoint}", {
  method: "POST",
  headers: {
    "Content-Type": "application/json",
    "X-API-Key": "${key}"
  },
  body: JSON.stringify(${params.replace(/\n/g, '\n    ')})
});

const data = await response.json();
console.log(data);`

      case 'python':
        return `import requests

response = requests.post(
    "${endpoint}",
    headers={
        "Content-Type": "application/json",
        "X-API-Key": "${key}"
    },
    json=${params.replace(/\n/g, '\n        ')}
)

print(response.json())`

      case 'go':
        return `package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

func main() {
    params := ${params.replace(/\n/g, '\n        ')}
    jsonData, _ := json.Marshal(params)

    req, _ := http.NewRequest("POST", "${endpoint}", bytes.NewBuffer(jsonData))
    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", "${key}")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    var result map[string]interface{}
    json.NewDecoder(resp.Body).Decode(&result)
    fmt.Println(result)
}`

      default:
        return ''
    }
  })

  // 监听 modelValue 变化
  watch(
    () => props.modelValue,
    (val) => {
      isOpen.value = val
    }
  )

  // 监听 workflow 变化
  watch(
    () => props.workflow,
    (workflow) => {
      if (workflow) {
        apiEnabled.value = workflow.api_enabled || false
        apiKey.value = workflow.api_key || ''
        timeout.value = workflow.api_timeout || 30
        webhookURL.value = workflow.api_webhook_url || ''
      }
    },
    { deep: true, immediate: true }
  )

  return {
    isOpen,
    activeTab,
    apiEnabled,
    apiKey,
    showApiKey,
    timeout,
    webhookURL,
    testMode,
    testParams,
    testing,
    testResult,
    selectedLanguage,
    tabList,
    apiEndpoint,
    displayApiKey,
    hasExternalTrigger,
    externalParams,
    codeExamples,
    currentCodeExample,
  }
}

