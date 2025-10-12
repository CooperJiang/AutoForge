import { createRouter, createWebHistory } from 'vue-router'
import SecureStorage, { STORAGE_KEYS } from '@/utils/storage'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/auth',
      name: 'auth',
      component: () => import('@/pages/Identity/index.vue'),
      meta: {
        title: '登录注册',
      },
    },
    {
      path: '/',
      component: () => import('@/layouts/DefaultLayout.vue'),
      children: [
        {
          path: '',
          name: 'home',
          component: () => import('@/pages/Tasks/index.vue'),
          meta: {
            title: '任务管理',
          },
        },
        {
          path: '/tasks',
          redirect: '/',
        },
        {
          path: '/tasks/create',
          name: 'tasks-create',
          component: () => import('@/pages/Tasks/create.vue'),
          meta: {
            title: '创建任务',
          },
        },
        {
          path: '/tasks/:id/edit',
          name: 'tasks-edit',
          component: () => import('@/pages/Tasks/create.vue'),
          meta: {
            title: '编辑任务',
          },
        },
        {
          path: '/profile',
          name: 'profile',
          component: () => import('@/pages/Profile.vue'),
          meta: {
            title: '个人中心',
          },
        },
        {
          path: '/tools',
          name: 'tools',
          component: () => import('@/pages/Tools/index.vue'),
          meta: {
            title: '工具箱',
          },
        },
        {
          path: '/workflows',
          name: 'workflows',
          component: () => import('@/pages/Workflows/index.vue'),
          meta: {
            title: '工作流',
          },
        },
        {
          path: '/workflows/create',
          name: 'workflows-create',
          component: () => import('@/pages/Workflows/editor.vue'),
          meta: {
            title: '创建工作流',
          },
        },
        {
          path: '/workflows/:id/edit',
          name: 'workflows-edit',
          component: () => import('@/pages/Workflows/editor.vue'),
          meta: {
            title: '编辑工作流',
          },
        },
        {
          path: '/workflows/:id/executions',
          name: 'workflows-executions',
          component: () => import('@/pages/Workflows/executions.vue'),
          meta: {
            title: '执行历史',
          },
        },
        {
          path: '/workflows/:id/executions/:executionId',
          name: 'workflows-execution-detail',
          component: () => import('@/pages/Workflows/execution-detail.vue'),
          meta: {
            title: '执行详情',
          },
        },
      ],
    },
    // 管理员路由
    {
      path: '/admin/login',
      name: 'admin-login',
      component: () => import('@/pages/Admin/Login.vue'),
      meta: {
        title: '管理员登录',
      },
    },
    {
      path: '/admin',
      component: () => import('@/layouts/AdminLayout.vue'),
      meta: {
        requiresAdminAuth: true,
      },
      children: [
        {
          path: '',
          redirect: '/admin/dashboard',
        },
        {
          path: 'dashboard',
          name: 'admin-dashboard',
          component: () => import('@/pages/Admin/Dashboard.vue'),
          meta: {
            title: '管理后台',
            requiresAdminAuth: true,
          },
        },
      ],
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found',
      component: () => import('@/pages/NotFound/index.vue'),
      meta: {
        title: '页面未找到',
      },
    },
  ],
})

// 路由守卫 - 检查用户登录状态
router.beforeEach((to, from, next) => {
  // 从SecureStorage获取token
  const token = SecureStorage.getItem<string>(STORAGE_KEYS.AUTH_TOKEN)

  // 白名单：不需要登录的页面
  const publicPaths = ['/auth', '/admin/login']
  const isPublicPath = publicPaths.includes(to.path)

  // 如果访问需要登录的页面但没有token，跳转到登录页
  if (!isPublicPath && !token) {
    next({ name: 'auth' })
    return
  }

  // 如果已登录且访问登录页，跳转到首页
  if (to.path === '/auth' && token) {
    next({ name: 'home' })
    return
  }

  next()
})

// 路由后置守卫 - 设置页面标题
router.afterEach((to) => {
  const defaultTitle = '定时任务系统'
  const pageTitle = to.meta.title as string

  if (pageTitle) {
    document.title = `${pageTitle} - ${defaultTitle}`
  } else {
    document.title = defaultTitle
  }
})

export default router
