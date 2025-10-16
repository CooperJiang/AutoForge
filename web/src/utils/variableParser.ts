/**
 * 变量解析工具
 * 支持格式：{{node_id.path.to.value}} 或 {{env.VARIABLE_NAME}}
 */

export interface VariableContext {
  nodes: Record<string, any>
  env: Record<string, any>
  trigger?: any
}

/**
 * 解析字符串中的变量引用
 * 例如: "{{node_1.response.data}}" -> 实际值
 */
export function parseVariables(text: string, context: VariableContext): any {
  if (typeof text !== 'string') return text

  const variableRegex = /\{\{([^}]+)\}\}/g

  return text.replace(variableRegex, (match, variable) => {
    const value = resolveVariable(variable.trim(), context)
    return value !== undefined ? value : match
  })
}

/**
 * 解析对象中的所有变量引用
 */
export function parseObjectVariables(obj: any, context: VariableContext): any {
  if (obj === null || obj === undefined) return obj

  if (typeof obj === 'string') {
    return parseVariables(obj, context)
  }

  if (Array.isArray(obj)) {
    return obj.map((item) => parseObjectVariables(item, context))
  }

  if (typeof obj === 'object') {
    const result: any = {}
    for (const key in obj) {
      result[key] = parseObjectVariables(obj[key], context)
    }
    return result
  }

  return obj
}

/**
 * 解析变量路径，获取实际值
 * 例如: "node_1.response.data" -> context.nodes.node_1.response.data
 */
function resolveVariable(path: string, context: VariableContext): any {
  const parts = path.split('.')
  const root = parts[0]

  let value: any

  if (root === 'env') {
    value = context.env
  } else if (root === 'trigger') {
    value = context.trigger
  } else if (context.nodes[root]) {
    value = context.nodes[root]
  } else {
    return undefined
  }

  for (let i = 1; i < parts.length; i++) {
    if (value === undefined || value === null) return undefined
    value = value[parts[i]]
  }

  return value
}

/**
 * 提取字符串中的所有变量引用
 * 返回格式: ['node_1.response.data', 'env.API_KEY']
 */
export function extractVariables(text: string): string[] {
  if (typeof text !== 'string') return []

  const variableRegex = /\{\{([^}]+)\}\}/g
  const variables: string[] = []
  let match

  while ((match = variableRegex.exec(text)) !== null) {
    variables.push(match[1].trim())
  }

  return variables
}

/**
 * 提取对象中的所有变量引用
 */
export function extractObjectVariables(obj: any): string[] {
  if (obj === null || obj === undefined) return []

  if (typeof obj === 'string') {
    return extractVariables(obj)
  }

  if (Array.isArray(obj)) {
    return obj.flatMap((item) => extractObjectVariables(item))
  }

  if (typeof obj === 'object') {
    const variables: string[] = []
    for (const key in obj) {
      variables.push(...extractObjectVariables(obj[key]))
    }
    return variables
  }

  return []
}

/**
 * 验证变量引用是否有效
 */
export function validateVariable(
  variable: string,
  availableNodes: string[]
): {
  valid: boolean
  message?: string
} {
  const parts = variable.split('.')
  const root = parts[0]

  if (root === 'env' || root === 'trigger') {
    return { valid: true }
  }

  if (!availableNodes.includes(root)) {
    return {
      valid: false,
      message: `节点 "${root}" 不存在`,
    }
  }

  return { valid: true }
}

/**
 * 获取节点的可用输出字段（用于UI提示）
 */
export function getNodeOutputSchema(nodeType: string, toolCode?: string): Record<string, string> {
  if (toolCode === 'http_request') {
    return {
      'response.status': '响应状态码',
      'response.data': '响应数据',
      'response.headers': '响应头',
      error: '错误信息（如果失败）',
    }
  }

  if (toolCode === 'email_sender') {
    return {
      success: '是否发送成功',
      messageId: '邮件ID',
      error: '错误信息（如果失败）',
    }
  }

  if (toolCode === 'health_checker') {
    return {
      healthy: '是否健康',
      status: '状态码',
      responseTime: '响应时间（毫秒）',
      error: '错误信息（如果失败）',
    }
  }

  if (nodeType === 'condition') {
    return {
      result: '条件结果 (true/false)',
      branch: '执行的分支 (true/false)',
    }
  }

  if (nodeType === 'switch') {
    return {
      matchedCase: '匹配的分支',
      value: '实际值',
    }
  }

  if (nodeType === 'trigger') {
    return {
      timestamp: '触发时间',
      type: '触发类型',
    }
  }

  return {}
}
