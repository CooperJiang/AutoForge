<template>
  <div class="variable-tree-selector">
    <div class="flex flex-wrap gap-2">
      <button
        v-for="node in previousNodes"
        :key="node.id"
        type="button"
        @click="toggleNodeTree(node.id)"
        class="inline-flex items-center gap-1.5 px-2.5 py-1.5 text-xs font-medium bg-primary-light hover:bg-primary text-primary hover:text-white border border-primary rounded-md transition-all"
        :title="`选择 ${node.name} 的变量`"
      >
        <Folder class="w-3.5 h-3.5" />
        <span>{{ node.name }}</span>
      </button>
    </div>

    <Teleport to="body">
      <div
        v-if="activeNodeId"
        class="fixed inset-0 z-50 flex items-start justify-center pt-20"
        @click.self="closeTree"
      >
        <div
          class="bg-bg-elevated rounded-lg shadow-2xl border border-border-primary max-w-2xl w-full mx-4 max-h-[70vh] flex flex-col"
        >
          <div class="px-4 py-3 border-b border-border-primary flex items-center justify-between">
            <div>
              <h3 class="text-sm font-semibold text-text-primary">选择变量</h3>
              <p class="text-xs text-text-tertiary mt-0.5">
                {{ activeNode?.name }} - 点击字段复制变量语法
              </p>
            </div>
            <button @click="closeTree" class="p-1 hover:bg-bg-hover rounded transition-colors">
              <X class="w-5 h-5 text-text-secondary" />
            </button>
          </div>

          <div
            class="px-4 py-2 bg-amber-50 dark:bg-amber-900/20 border-b border-amber-200 dark:border-amber-800"
          >
            <p class="text-xs text-amber-800 dark:text-amber-200">
              💡 以下是系统已知字段。部分来自外部接口的动态字段未列出，您可以手动输入完整路径访问。
            </p>
          </div>

          <div class="flex-1 overflow-y-auto p-4">
            <div class="space-y-1">
              <VariableTreeNode
                v-for="(value, key) in getNodeStructure(activeNodeId)"
                :key="key"
                :name="String(key)"
                :value="value"
                :path="getVariablePath(activeNodeId, String(key))"
                :level="0"
                @copy="handleCopy"
              />
            </div>
          </div>

          <div class="px-4 py-3 border-t border-border-primary bg-bg-hover">
            <p class="text-xs text-text-secondary">
              <strong>提示：</strong>点击任意字段复制变量，格式为
              <code
                class="px-1 py-0.5 bg-bg-elevated rounded text-primary font-mono"
                v-text="'{{nodes.xxx.field}}'"
              ></code>
            </p>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { X, Folder } from 'lucide-vue-next'
import VariableTreeNode from '@/components/VariableTreeNode'
import { message } from '@/utils/message'
import { getToolList, describeToolOutput, type OutputFieldDef } from '@/api/tool'

interface NodeData {
  id: string
  name: string
  type: string
  toolCode?: string
  config?: Record<string, any>
}

interface Props {
  previousNodes?: NodeData[]
}

const props = defineProps<Props>()

const toolOutputSchemas = ref<Record<string, Record<string, OutputFieldDef>>>({})
const activeNodeId = ref<string | null>(null)
const dynamicSchemas = ref<Record<string, Record<string, OutputFieldDef>>>({})

const activeNode = computed(() => {
  return props.previousNodes?.find((n) => n.id === activeNodeId.value)
})

const toggleNodeTree = (nodeId: string) => {
  if (activeNodeId.value === nodeId) {
    closeTree()
  } else {
    activeNodeId.value = nodeId
    // 尝试加载动态输出结构（如工具支持）
    const node = props.previousNodes?.find((n) => n.id === nodeId)
    if (node && node.toolCode) {
      // 仅在未缓存或配置变化时请求；这里简单按节点ID缓存一次
      if (!dynamicSchemas.value[nodeId]) {
        const cfg = node.config || {}
        describeToolOutput(node.toolCode, cfg)
          .then((schema) => {
            if (schema && Object.keys(schema).length > 0) {
              dynamicSchemas.value[nodeId] = schema
            }
          })
          .catch(() => {
            // 动态描述失败时静默回退
          })
      }
    }
  }
}

const closeTree = () => {
  activeNodeId.value = null
}

const getVariablePath = (nodeId: string | null, key: string) => {
  if (!nodeId) return key

  const node = props.previousNodes?.find((n) => n.id === nodeId)
  if (!node) return `nodes.${nodeId}.${key}`

  if (node.type === 'external_trigger') {
    return `external.${key}`
  }

  return `nodes.${nodeId}.${key}`
}

onMounted(async () => {
  try {
    const tools = await getToolList()
    const schemas: Record<string, Record<string, OutputFieldDef>> = {}

    tools.forEach((tool) => {
      if (tool.output_fields_schema) {
        schemas[tool.code] = tool.output_fields_schema
      }
    })

    toolOutputSchemas.value = schemas
  } catch {
    console.error('Failed to load tool schemas:', error)
  }
})

const convertFieldDefToStructure = (fieldDef: OutputFieldDef): string | Record<string, any> => {
  if (fieldDef.type === 'string' || fieldDef.type === 'number' || fieldDef.type === 'boolean') {
    return fieldDef.label
  }

  if ((fieldDef.type === 'object' || fieldDef.type === 'array') && fieldDef.children) {
    const result: Record<string, any> = {}

    if (fieldDef.label) {
      result._description = fieldDef.label
    }

    Object.entries(fieldDef.children).forEach(([key, childDef]) => {
      result[key] = convertFieldDefToStructure(childDef)
    })

    return result
  }

  return fieldDef.label
}

const getNodeStructure = (nodeId: string | null): Record<string, any> => {
  if (!nodeId) return {}

  const node = props.previousNodes?.find((n) => n.id === nodeId)
  if (!node) return {}

  // 优先使用动态输出结构（如果后端提供）
  if (dynamicSchemas.value[node.id]) {
    const structure: Record<string, any> = {}
    Object.entries(dynamicSchemas.value[node.id]).forEach(([key, fieldDef]) => {
      structure[key] = convertFieldDefToStructure(fieldDef)
    })
    return structure
  }

  if (node.toolCode && toolOutputSchemas.value[node.toolCode]) {
    const schema = toolOutputSchemas.value[node.toolCode]
    const structure: Record<string, any> = {}

    Object.entries(schema).forEach(([key, fieldDef]) => {
      structure[key] = convertFieldDefToStructure(fieldDef)
    })

    return structure
  }

  if (node.type === 'external_trigger') {
    return getExternalTriggerStructure(node)
  }

  return getGenericNodeStructure(node)
}

const getExternalTriggerStructure = (node: NodeData): Record<string, string> => {
  const params = node.config?.params || []

  if (params.length === 0) {
    return {
      _empty: 'No parameters configured',
    }
  }

  const structure: Record<string, string> = {}
  const typeLabels: Record<string, string> = {
    string: 'String',
    number: 'Number',
    boolean: 'Boolean',
    object: 'Object',
    array: 'Array',
  }

  params.forEach((param: any) => {
    if (param.key) {
      const typeLabel = typeLabels[param.type] || param.type
      const description = param.description || `${typeLabel} parameter`
      const required = param.required ? ' [Required]' : ' [Optional]'
      const example = param.example ? ` (e.g., ${param.example})` : ''

      structure[param.key] = `${description}${required} - ${typeLabel}${example}`
    }
  })

  return structure
}

const getGenericNodeStructure = (node: NodeData): Record<string, string> => {
  if (node.type === 'condition') {
    return {
      result: 'Condition result (true/false)',
      message: 'Condition message',
    }
  }

  if (node.type === 'trigger') {
    return {
      triggered: 'Trigger status',
      timestamp: 'Trigger timestamp',
    }
  }

  return {
    success: 'Success flag',
    message: 'Execution message',
    data: 'Output data',
    error: 'Error message',
  }
}

const handleCopy = async (path: string) => {
  const variableText = `{{${path}}}`

  try {
    await navigator.clipboard.writeText(variableText)
    message.success(`Copied: ${variableText}`)
  } catch (err) {
    console.error('Copy failed:', err)
    message.error('Failed to copy')
  }
}
</script>

<style scoped>
code {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
}
</style>
