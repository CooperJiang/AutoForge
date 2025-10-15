import { computed } from 'vue'
import { useRoute } from 'vue-router'

export interface BreadcrumbItem {
  label: string
  to?: string
}

export function useBreadcrumb() {
  const route = useRoute()

  const breadcrumbItems = computed<BreadcrumbItem[]>(() => {
    const path = route.path
    const items: BreadcrumbItem[] = []

    // 首页/任务管理
    if (path === '/' || path === '/tasks') {
      items.push({ label: '任务管理' })
    }
    // 创建任务
    else if (path === '/tasks/create') {
      items.push({ label: '任务管理', to: '/' })
      items.push({ label: '创建任务' })
    }
    // 编辑任务
    else if (path.startsWith('/tasks/') && path.endsWith('/edit')) {
      items.push({ label: '任务管理', to: '/' })
      items.push({ label: '编辑任务' })
    }
    // 工作流
    else if (path === '/workflows') {
      items.push({ label: '工作流' })
    }
    // 创建工作流
    else if (path === '/workflows/create') {
      items.push({ label: '工作流', to: '/workflows' })
      items.push({ label: '创建工作流' })
    }
    // 编辑工作流
    else if (path.startsWith('/workflows/') && path.endsWith('/edit')) {
      items.push({ label: '工作流', to: '/workflows' })
      items.push({ label: '编辑工作流' })
    }
    // 执行历史
    else if (path.match(/^\/workflows\/[^/]+\/executions$/)) {
      items.push({ label: '工作流', to: '/workflows' })
      items.push({ label: '执行历史' })
    }
    // 执行详情
    else if (path.match(/^\/workflows\/[^/]+\/executions\/[^/]+$/)) {
      const workflowId = route.params.id as string
      items.push({ label: '工作流', to: '/workflows' })
      items.push({ label: '执行历史', to: `/workflows/${workflowId}/executions` })
      items.push({ label: '执行详情' })
    }
    // 工具箱
    else if (path === '/tools') {
      items.push({ label: '工具箱' })
    }
    // 个人中心
    else if (path === '/profile') {
      items.push({ label: '个人中心' })
    }
    // 管理后台
    else if (path.startsWith('/admin')) {
      items.push({ label: '管理后台' })
    }

    return items
  })

  return {
    breadcrumbItems,
  }
}
