<template>
  <div class="min-h-screen flex items-center justify-center bg-gradient-to-br from-slate-100 via-green-50 to-cyan-50">
    <div class="w-full max-w-md px-6">
      <!-- Logo 和标题 -->
      <div class="text-center mb-8">
        <div class="inline-flex items-center justify-center w-20 h-20 mb-4">
          <img src="/logo.png" alt="Logo" class="w-full h-full object-contain" />
        </div>
        <h1 class="text-3xl font-bold text-slate-900 mb-2">管理后台</h1>
        <p class="text-slate-600">请输入管理员密码以继续</p>
      </div>

      <!-- 登录卡片 -->
      <div class="bg-white rounded-2xl shadow-xl border-2 border-slate-200 p-8">
        <form @submit.prevent="handleLogin">
          <div class="mb-6">
            <label class="block text-sm font-medium text-slate-700 mb-2">
              管理员密码
            </label>
            <input
              v-model="password"
              type="password"
              placeholder="请输入密码"
              class="w-full px-4 py-3 border-2 border-slate-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-green-500 focus:border-transparent transition-all"
              :disabled="loading"
              autofocus
            />
          </div>

          <button
            type="submit"
            :disabled="loading || !password"
            class="w-full bg-gradient-to-r from-green-500 to-cyan-500 text-white font-semibold py-3 rounded-lg hover:from-green-600 hover:to-cyan-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all shadow-lg hover:shadow-xl"
          >
            <span v-if="!loading">登录</span>
            <span v-else>登录中...</span>
          </button>
        </form>

        <!-- 提示信息 -->
        <div class="mt-6 text-center text-sm text-slate-500">
          <p>默认密码请查看配置文件</p>
        </div>
      </div>

      <!-- 返回首页 -->
      <div class="text-center mt-6">
        <router-link
          to="/"
          class="text-slate-600 hover:text-slate-900 text-sm font-medium transition-colors"
        >
          ← 返回首页
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import * as adminApi from '@/api/admin'
import { message } from '@/utils/message'

const router = useRouter()
const password = ref('')
const loading = ref(false)

const handleLogin = async () => {
  if (!password.value) {
    message.warning('请输入密码')
    return
  }

  loading.value = true

  try {
    const res = await adminApi.login({ password: password.value })

    // 存储 token
    localStorage.setItem('admin_token', res.data.token)

    message.success('登录成功')

    // 跳转到管理面板
    router.push('/admin/dashboard')
  } catch (error: any) {
    // 错误已在 request.ts 拦截器中显示
    console.error('登录失败:', error)
  } finally {
    loading.value = false
  }
}
</script>
