import type { LucideIcon } from 'lucide-vue-next'
import {
  Globe,
  Mail,
  Activity,
  Shuffle,
  Zap,
  MessageSquare,
  Palette,
  FileJson,
  MessageCircle,
  Sparkles,
  Image as ImageIcon,
  Database,
  Upload,
  Cloud,
} from 'lucide-vue-next'

export interface ToolUsageItem {
  text: string
}

export interface ToolConfig {
  code: string
  title: string
  description: string
  icon: LucideIcon | string
  iconBg: string
  usageTitle?: string
  usageDescription?: string
  usageItems: ToolUsageItem[]
  tags?: string[]
}

export const TOOL_CONFIGS: Record<string, ToolConfig> = {
  http_request: {
    code: 'http_request',
    title: '📡 HTTP 请求工具',
    description: 'Send HTTP requests to any URL with full control over methods, headers, and body',
    icon: Globe,
    iconBg: 'bg-gradient-to-br from-blue-500 to-blue-600',
    usageTitle: 'HTTP Request Tool',
    usageDescription: '发送 HTTP 请求到指定的 URL，支持所有常见的 HTTP 方法。',
    usageItems: [
      { text: '支持 GET、POST、PUT、DELETE、PATCH 等方法' },
      { text: '自定义请求头（Headers）、参数（Params）、请求体（Body）' },
      { text: '支持粘贴 cURL 命令自动解析配置' },
      { text: '适用场景：API 调用、数据抓取、Webhook 触发等' },
    ],
    tags: ['HTTP', 'API', 'Request', 'Web'],
  },

  email_sender: {
    code: 'email_sender',
    title: '📧 邮件发送工具',
    description: 'Send emails with SMTP protocol, supports multiple recipients and HTML format',
    icon: Mail,
    iconBg: 'bg-gradient-to-br from-red-500 to-pink-600',
    usageTitle: 'Email Sender',
    usageDescription: '通过 SMTP 协议发送邮件通知，支持多收件人和 HTML 格式。',
    usageItems: [
      { text: '系统自动使用配置的 SMTP 服务器，无需用户提供' },
      { text: '支持多个收件人、抄送（CC）' },
      { text: '支持纯文本和 HTML 格式邮件' },
      { text: '适用场景：通知提醒、报告发送、告警邮件等' },
    ],
    tags: ['Email', 'SMTP', 'Notification'],
  },

  health_checker: {
    code: 'health_checker',
    title: '🏥 健康检查工具',
    description: 'Monitor website/API availability, response time, and SSL certificate validity',
    icon: Activity,
    iconBg: 'bg-gradient-to-br from-green-500 to-emerald-600',
    usageTitle: 'Health Checker',
    usageDescription: '检查网站/API的可用性、响应时间、状态码和内容匹配。',
    usageItems: [
      { text: '支持 HTTP/HTTPS 健康检查' },
      { text: '监控响应时间和状态码' },
      { text: '检查 SSL 证书有效期' },
      { text: '支持内容匹配（正则表达式）' },
      { text: '适用场景：服务监控、可用性检测、告警通知等' },
    ],
    tags: ['Health', 'Monitor', 'SSL', 'Uptime'],
  },

  json_transform: {
    code: 'json_transform',
    title: '🔄 JSON 转换工具',
    description: 'Transform JSON data using JavaScript expressions with variable support',
    icon: Shuffle,
    iconBg: 'bg-gradient-to-br from-purple-500 to-indigo-600',
    usageTitle: 'JSON Transform',
    usageDescription: '基于变量系统和 JS 表达式动态转换数据。',
    usageItems: [
      { text: '支持引用前置节点的输出数据作为输入' },
      { text: '使用 JavaScript 表达式进行数据转换' },
      { text: '内置 data 和 ctx 变量，无需声明' },
      { text: '支持数组映射、过滤、聚合等操作' },
      { text: '示例：data.map(item => item.url)' },
      { text: '适用场景：数据清洗、格式转换、字段提取等' },
    ],
    tags: ['JSON', 'Transform', 'JavaScript', 'Data Processing'],
  },

  redis_context: {
    code: 'redis_context',
    title: '💾 Redis 上下文存储',
    description:
      'Store and retrieve workflow context data using Redis for cross-execution state management',
    icon: Zap,
    iconBg: 'bg-gradient-to-br from-red-600 to-orange-600',
    usageTitle: 'Redis Context Storage',
    usageDescription: '使用 Redis 存储和获取工作流上下文数据，实现跨执行的状态管理。',
    usageItems: [
      { text: '支持 SET（存储）和 GET（获取）操作' },
      { text: '可设置过期时间（TTL），自动清理过期数据' },
      { text: '支持存储任意 JSON 数据' },
      { text: '适用场景：状态记录、去重判断、计数器、会话管理等' },
    ],
    tags: ['Redis', 'Storage', 'Context', 'State'],
  },

  feishu_bot: {
    code: 'feishu_bot',
    title: '📱 飞书机器人',
    description: 'Send messages to Feishu via webhook, supports text, rich text, images and cards',
    icon: MessageSquare,
    iconBg: 'bg-gradient-to-br from-blue-600 to-cyan-600',
    usageTitle: 'Feishu Bot',
    usageDescription: '通过飞书机器人 Webhook 发送消息通知，支持文本、富文本、图片和卡片消息。',
    usageItems: [
      { text: '支持多种消息类型：文本、富文本、图片、卡片' },
      { text: '卡片消息支持预设模板：成功、警告、错误、信息' },
      { text: '支持自定义卡片字段和按钮' },
      { text: '可配置签名密钥进行安全验证' },
      { text: '适用场景：团队通知、告警提醒、工作报告等' },
    ],
    tags: ['Feishu', 'Bot', 'Notification', 'IM'],
  },

  output_formatter: {
    code: 'output_formatter',
    title: '🎨 输出格式化',
    description: 'Format output data as image, video, HTML, Markdown, or JSON for rich display',
    icon: Palette,
    iconBg: 'bg-gradient-to-br from-pink-500 to-rose-600',
    usageTitle: 'Output Formatter',
    usageDescription: '将数据格式化为指定的展示类型（图片、视频、HTML、Markdown 等）。',
    usageItems: [
      { text: '支持多种输出类型：image、video、html、markdown、json 等' },
      { text: '可引用前置节点的输出作为内容源' },
      { text: '支持图片画廊（gallery）展示多张图片' },
      { text: 'HTML 支持两种方式：直接内容或 URL 地址' },
      { text: '适用场景：数据可视化、报告生成、多媒体展示等' },
    ],
    tags: ['Format', 'Display', 'Output', 'Visualization'],
  },

  html_render: {
    code: 'html_render',
    title: '🌐 HTML 内容保存',
    description: 'Save HTML content as a web page file and generate shareable preview URL',
    icon: FileJson,
    iconBg: 'bg-gradient-to-br from-orange-500 to-amber-600',
    usageTitle: 'HTML Render',
    usageDescription: '将 HTML 内容保存为网页文件并生成可访问的 URL。',
    usageItems: [
      { text: '将 HTML 内容保存到服务器并生成预览 URL' },
      { text: '支持设置网页标题' },
      { text: '可配置链接过期时间（0 表示永不过期）' },
      { text: '生成的 URL 可直接分享或嵌入' },
      { text: '适用场景：报告分享、网页归档、内容预览等' },
    ],
    tags: ['HTML', 'Render', 'Web', 'Preview'],
  },

  openai_chatgpt: {
    code: 'openai_chatgpt',
    title: '🤖 OpenAI 对话',
    description: 'Use OpenAI Chat API for conversations, supports GPT-3.5, GPT-4, GPT-4o models',
    icon: MessageCircle,
    iconBg: 'bg-gradient-to-br from-emerald-500 to-teal-600',
    usageTitle: 'OpenAI Chat',
    usageDescription: '使用 OpenAI Chat API 进行对话，支持 GPT-3.5、GPT-4、GPT-4o 等模型。',
    usageItems: [
      { text: '支持多种 GPT 模型：gpt-3.5-turbo、gpt-4、gpt-4o 等' },
      { text: '可自定义系统角色（System Message）设定 AI 行为' },
      { text: '可调节 temperature 参数控制回复随机性' },
      { text: '支持设置最大 token 数量' },
      { text: '适用场景：智能问答、内容生成、文本分析等' },
    ],
    tags: ['OpenAI', 'GPT', 'AI', 'Chat'],
  },

  openai_image: {
    code: 'openai_image',
    title: '🎨 OpenAI 图片生成',
    description: 'Generate images using OpenAI DALL-E models from text descriptions',
    icon: ImageIcon,
    iconBg: 'bg-gradient-to-br from-violet-500 to-purple-600',
    usageTitle: 'OpenAI Image Generation',
    usageDescription: '使用 OpenAI DALL-E 模型生成图片。',
    usageItems: [
      { text: '支持 DALL-E 2、DALL-E 3、gpt-image-1等模型' },
      { text: '可自定义图片尺寸和数量' },
      { text: '支持标准和高清质量' },
      { text: '输出格式：URL 或 Base64 编码' },
      { text: '适用场景：图片创作、设计辅助、内容配图等' },
    ],
    tags: ['OpenAI', 'DALL-E', 'Image', 'AI'],
  },

  context_manager: {
    code: 'context_manager',
    title: '💬 对话上下文管理器',
    description:
      'Manage conversation context with prepare and persist modes for multi-turn dialogues',
    icon: Database,
    iconBg: 'bg-gradient-to-br from-indigo-500 to-purple-600',
    usageTitle: 'Context Manager',
    usageDescription:
      '管理多轮对话的上下文历史，支持准备消息（Prepare）和保存消息（Persist）两种模式。',
    usageItems: [
      { text: 'Prepare 模式：读取历史 + 拼接当前消息 → 输出 messages_json' },
      { text: 'Persist 模式：保存 AI 回复到 Redis 历史记录' },
      { text: '支持会话隔离（通过 session_key 区分不同用户）' },
      { text: '自动裁剪到窗口大小，防止上下文过长' },
      { text: '通用设计，可与任意 LLM（OpenAI、Gemini、Claude 等）配合使用' },
      { text: '适用场景：智能客服、AI 助手、教育辅导、知识问答等' },
    ],
    tags: ['Context', 'Memory', 'Chat', 'Conversation', 'Session'],
  },

  pixelpunk_upload: {
    code: 'pixelpunk_upload',
    title: '📸 PixelPunk 图床上传',
    description: '上传图片到 PixelPunk 图床，返回 CDN URL',
    icon: Upload,
    iconBg: 'bg-gradient-to-br from-cyan-500 to-blue-600',
    usageTitle: 'PixelPunk Image Upload',
    usageDescription: '将图片上传到 PixelPunk 图床，获取永久可访问的 CDN URL。',
    usageItems: [
      { text: '支持多种访问级别：public（公开）、private（私有）、protected（受保护）' },
      { text: '可选图片优化压缩，减少文件大小' },
      { text: '支持虚拟路径管理和文件夹分类' },
      { text: '返回原图 URL、缩略图 URL、图片尺寸等完整信息' },
      { text: '适用场景：内容发布、图片存储、CDN 加速等' },
    ],
    tags: ['Image', 'Upload', 'CDN', 'Storage', 'PixelPunk'],
  },

  aliyun_oss: {
    code: 'aliyun_oss',
    title: '☁️ 阿里云 OSS 上传',
    description: '上传文件到阿里云对象存储服务',
    icon: Cloud,
    iconBg: 'bg-gradient-to-br from-orange-500 to-red-600',
    usageTitle: 'Aliyun OSS Upload',
    usageDescription: '将文件上传到阿里云 OSS，获取永久可访问的文件 URL。',
    usageItems: [
      { text: '支持任意类型文件上传' },
      { text: '自动识别文件类型（Content-Type）' },
      { text: '可自定义 OSS 存储路径' },
      { text: '配置统一管理在后端，安全可靠' },
      { text: '适用场景：文件存储、图片管理、视频上传、文档分发等' },
    ],
    tags: ['Storage', 'Upload', 'Aliyun', 'OSS', 'Cloud'],
  },

  tencent_cos: {
    code: 'tencent_cos',
    title: '☁️ 腾讯云 COS 上传',
    description: '上传文件到腾讯云对象存储服务',
    icon: Cloud,
    iconBg: 'bg-gradient-to-br from-blue-500 to-cyan-600',
    usageTitle: 'Tencent COS Upload',
    usageDescription: '将文件上传到腾讯云 COS，获取永久可访问的文件 URL。',
    usageItems: [
      { text: '支持任意类型文件上传' },
      { text: '自动识别文件类型（Content-Type）' },
      { text: '可自定义 COS 存储路径' },
      { text: '配置统一管理在后端，安全可靠' },
      { text: '适用场景：文件存储、图片管理、视频上传、文档分发等' },
    ],
    tags: ['Storage', 'Upload', 'Tencent', 'COS', 'Cloud'],
  },
}

/**
 * Get tool configuration by code
 */
export function getToolConfig(code: string): ToolConfig | undefined {
  return TOOL_CONFIGS[code]
}

/**
 * Get all tool configurations
 */
export function getAllToolConfigs(): ToolConfig[] {
  return Object.values(TOOL_CONFIGS)
}

/**
 * Get tool icon component or path
 */
export function getToolIcon(code: string): LucideIcon | string {
  const config = getToolConfig(code)
  return config?.icon || Sparkles
}

/**
 * Get tool icon background class
 */
export function getToolIconBg(code: string): string {
  const config = getToolConfig(code)
  return config?.iconBg || 'bg-gradient-to-br from-gray-500 to-gray-600'
}
