import { ref, computed, type Ref } from 'vue'
import type { WorkflowExecution, StatusType, NodeStatusType } from '../types'
import { Clock, Play, User, Globe } from 'lucide-vue-next'

/**
 * 执行详情状态管理
 */
export function useExecutionState() {
  const execution = ref<WorkflowExecution | null>(null)
  const loading = ref(true)
  const polling = ref(false)
  const openSections = ref<Record<string, Set<string>>>({})
  const expandedErrors = ref<Record<string, boolean>>({})

  // ============ 状态样式函数 ============

  const getStatusClass = (status: StatusType) => {
    const classes = {
      running: 'bg-primary-light text-primary',
      success: 'bg-green-100 text-green-700',
      failed: 'bg-red-100 text-red-700',
      cancelled: 'bg-bg-tertiary text-text-secondary',
    }
    return classes[status] || classes.cancelled
  }

  const getStatusDotClass = (status: StatusType) => {
    const classes = {
      running: 'bg-[var(--color-primary)] animate-pulse',
      success: 'bg-green-500',
      failed: 'bg-red-500',
      cancelled: 'bg-bg-hover',
    }
    return classes[status] || classes.cancelled
  }

  const getStatusText = (status: StatusType) => {
    const texts = {
      running: '运行中',
      success: '执行成功',
      failed: '执行失败',
      cancelled: '已取消',
    }
    return texts[status] || '未知'
  }

  // ============ 节点状态样式函数 ============

  const getNodeStatusClass = (status: NodeStatusType) => {
    const classes = {
      pending: 'bg-bg-tertiary text-text-secondary',
      running: 'bg-primary-light text-primary',
      success: 'bg-green-100 text-green-700',
      failed: 'bg-red-100 text-red-700',
      skipped: 'bg-amber-100 text-amber-700',
    }
    return classes[status] || classes.pending
  }

  const getNodeStatusBgClass = (status: NodeStatusType) => {
    const classes = {
      pending: 'bg-bg-tertiary text-text-secondary',
      running: 'bg-[var(--color-primary)] text-white',
      success: 'bg-green-500 text-white',
      failed: 'bg-red-500 text-white',
      skipped: 'bg-amber-500 text-white',
    }
    return classes[status] || classes.pending
  }

  const getNodeStatusText = (status: NodeStatusType) => {
    const texts = {
      pending: '等待中',
      running: '运行中',
      success: '成功',
      failed: '失败',
      skipped: '已跳过',
    }
    return texts[status] || '未知'
  }

  // ============ 触发方式图标和文本 ============

  const getTriggerIcon = (triggerType: string) => {
    const icons: Record<string, any> = {
      manual: Play,
      schedule: Clock,
      api: Globe,
      user: User,
    }
    return icons[triggerType] || Clock
  }

  const getTriggerText = (triggerType: string) => {
    const texts: Record<string, string> = {
      manual: '手动触发',
      schedule: '定时触发',
      api: 'API 触发',
      user: '用户触发',
    }
    return texts[triggerType] || '未知'
  }

  // ============ 展开/收起逻辑 ============

  const toggleSection = (nodeId: string, section: 'input' | 'output') => {
    if (!openSections.value[nodeId]) {
      openSections.value[nodeId] = new Set()
    }
    if (openSections.value[nodeId].has(section)) {
      openSections.value[nodeId].delete(section)
    } else {
      openSections.value[nodeId].add(section)
    }
  }

  const toggleErrorExpand = (nodeId: string) => {
    expandedErrors.value[nodeId] = !expandedErrors.value[nodeId]
  }

  const isSectionOpen = (nodeId: string, section: 'input' | 'output') => {
    return openSections.value[nodeId]?.has(section) || false
  }

  return {
    execution,
    loading,
    polling,
    openSections,
    expandedErrors,

    // 样式函数
    getStatusClass,
    getStatusDotClass,
    getStatusText,
    getNodeStatusClass,
    getNodeStatusBgClass,
    getNodeStatusText,

    // 触发方式
    getTriggerIcon,
    getTriggerText,

    // 展开/收起
    toggleSection,
    toggleErrorExpand,
    isSectionOpen,
  }
}

