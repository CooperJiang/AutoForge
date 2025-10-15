interface ParsedCurl {
  url: string
  method: string
  headers: { key: string; value: string }[]
  params: { key: string; value: string }[]
  body?: string
}

export function parseCurl(curlCommand: string): ParsedCurl | null {
  try {
    // 清理命令，移除换行和多余空格
    const cleaned = curlCommand
      .replace(/\\\n/g, ' ')
      .replace(/\n/g, ' ')
      .replace(/\s+/g, ' ')
      .trim()

    // 检查是否是 curl 命令
    if (!cleaned.startsWith('curl')) {
      return null
    }

    // 提取 URL
    const urlMatch = cleaned.match(/curl\s+'([^']+)'|curl\s+"([^"]+)"|curl\s+([^\s-]+)/)
    if (!urlMatch) {
      return null
    }
    const url = urlMatch[1] || urlMatch[2] || urlMatch[3]

    // 解析URL，分离查询参数
    const urlObj = new URL(url)
    const params: { key: string; value: string }[] = []
    urlObj.searchParams.forEach((value, key) => {
      params.push({ key, value })
    })

    // 提取 HTTP 方法
    const methodMatch = cleaned.match(/-X\s+([A-Z]+)/)
    const method = methodMatch ? methodMatch[1] : 'GET'

    // 提取 Headers
    const headers: { key: string; value: string }[] = []
    const headerRegex = /-H\s+'([^:]+):\s*([^']+)'|-H\s+"([^:]+):\s*([^"]+)"/g
    let headerMatch

    while ((headerMatch = headerRegex.exec(cleaned)) !== null) {
      const key = headerMatch[1] || headerMatch[3]
      const value = headerMatch[2] || headerMatch[4]
      headers.push({ key: key.trim(), value: value.trim() })
    }

    // 提取 body 数据
    let body: string | undefined

    // 支持多种 body 数据格式：--data, -d, --data-raw, --data-binary, --data-urlencode, --json
    const bodyPatterns = [
      /--data-raw\s+'([^']*)'|--data-raw\s+"([^"]*)"/,
      /--data-binary\s+'([^']*)'|--data-binary\s+"([^"]*)"/,
      /--data-urlencode\s+'([^']*)'|--data-urlencode\s+"([^"]*)"/,
      /--data\s+'([^']*)'|--data\s+"([^"]*)"/,
      /-d\s+'([^']*)'|-d\s+"([^"]*)"/,
      /--json\s+'([^']*)'|--json\s+"([^"]*)"/,
    ]

    for (const pattern of bodyPatterns) {
      const dataMatch = cleaned.match(pattern)
      if (dataMatch) {
        body = dataMatch[1] || dataMatch[2]
        break
      }
    }

    // 如果没有匹配到带引号的，尝试匹配不带引号的格式
    if (!body) {
      const unquotedMatch = cleaned.match(
        /(?:--data-raw|--data-binary|--data|-d|--json)\s+([^\s-][^\s]*?)(?:\s+--|$|\s+-[a-zA-Z])/
      )
      if (unquotedMatch) {
        body = unquotedMatch[1]
      }
    }

    // 如果有 body 且方法是 GET，改为 POST
    if (body && method === 'GET') {
      return {
        url: urlObj.origin + urlObj.pathname,
        method: 'POST',
        headers,
        params,
        body,
      }
    }

    return {
      url: urlObj.origin + urlObj.pathname,
      method,
      headers,
      params,
      body,
    }
  } catch (error) {
    console.error('解析 curl 失败:', error)
    return null
  }
}
