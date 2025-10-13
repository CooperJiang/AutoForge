<template>
  <div class="min-h-screen bg-gradient-to-br from-bg-secondary via-bg-secondary to-bg-secondary flex items-center justify-center p-4 overflow-hidden relative">
    <!-- 背景装饰 -->
    <div class="absolute inset-0 overflow-hidden pointer-events-none">
      <div class="absolute top-20 left-10 w-72 h-72 bg-green-200 rounded-full mix-blend-multiply filter blur-xl opacity-30 animate-blob"></div>
      <div class="absolute top-40 right-10 w-72 h-72 bg-cyan-200 rounded-full mix-blend-multiply filter blur-xl opacity-30 animate-blob animation-delay-2000"></div>
      <div class="absolute -bottom-8 left-1/2 w-72 h-72 bg-teal-200 rounded-full mix-blend-multiply filter blur-xl opacity-30 animate-blob animation-delay-4000"></div>
    </div>

    <div class="max-w-6xl w-full grid lg:grid-cols-2 gap-8 items-center relative z-10">
      <!-- 左侧介绍区域 -->
      <div class="hidden lg:block space-y-8 animate-fade-in-left">
        <div class="space-y-4">
          <div class="inline-flex items-center gap-2 px-3 py-1.5 bg-bg-elevated/50 backdrop-blur-sm border border-green-200 rounded-full">
            <div class="w-2 h-2 bg-green-500 rounded-full animate-pulse"></div>
            <span class="text-sm font-medium text-text-secondary">智能定时任务管理</span>
          </div>
          <h1 class="text-5xl font-bold text-text-primary leading-tight">
            让任务<br/>
            <span class="bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent">自动执行</span>
          </h1>
          <p class="text-lg text-text-secondary leading-relaxed max-w-md">
            简单易用的定时任务系统，支持多种调度规则，轻松管理您的 HTTP 定时任务
          </p>
        </div>

        <!-- 功能亮点 -->
        <div class="space-y-4">
          <div v-for="(feature, index) in features" :key="index"
               class="flex items-start gap-3 p-4 bg-bg-elevated/60 backdrop-blur-sm border border-border-primary rounded-xl hover:shadow-md transition-all duration-300 animate-fade-in-up"
               :style="{ animationDelay: `${index * 100}ms` }">
            <div class="flex-shrink-0 w-10 h-10 bg-gradient-to-br from-primary to-accent rounded-lg flex items-center justify-center">
              <component :is="feature.icon" class="w-5 h-5 text-white" />
            </div>
            <div>
              <h3 class="font-semibold text-text-primary text-sm">{{ feature.title }}</h3>
              <p class="text-xs text-text-secondary mt-0.5">{{ feature.description }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- 右侧登录/注册表单 -->
      <div class="flex justify-center lg:justify-end">
        <div class="bg-bg-elevated/80 backdrop-blur-xl border-2 border-bg-elevated shadow-2xl rounded-2xl p-8 max-w-md w-full animate-fade-in-up">
          <!-- 图标和标题 -->
          <div class="text-center mb-6">
            <div class="inline-flex items-center justify-center w-20 h-20 mb-4 animate-bounce-subtle">
              <img src="/logo.png" alt="Logo" class="w-full h-full object-contain" />
            </div>
            <h2 class="text-2xl font-bold text-text-primary mb-2">欢迎使用</h2>
            <p class="text-sm text-text-secondary">定时任务管理系统</p>
          </div>

          <!-- Tab 切换 -->
          <div class="flex bg-bg-tertiary rounded-lg p-1 mb-6">
            <button
              @click="activeTab = 'login'"
              :class="[
                'flex-1 py-2 px-4 text-sm font-medium rounded-md transition-all duration-200',
                activeTab === 'login'
                  ? 'bg-bg-elevated text-text-primary shadow-sm'
                  : 'text-text-secondary hover:text-text-primary'
              ]"
            >
              登录
            </button>
            <button
              @click="activeTab = 'register'"
              :class="[
                'flex-1 py-2 px-4 text-sm font-medium rounded-md transition-all duration-200',
                activeTab === 'register'
                  ? 'bg-bg-elevated text-text-primary shadow-sm'
                  : 'text-text-secondary hover:text-text-primary'
              ]"
            >
              注册
            </button>
          </div>

          <!-- 登录表单 -->
          <form v-if="activeTab === 'login'" @submit.prevent="handleLogin" class="space-y-4">
            <BaseInput
              v-model="loginForm.account"
              label="用户名或邮箱"
              placeholder="请输入用户名或邮箱"
              required
            />

            <BaseInput
              v-model="loginForm.password"
              type="password"
              label="密码"
              placeholder="请输入密码"
              required
            />

            <BaseButton
              type="submit"
              variant="primary"
              :loading="loading"
              class="w-full"
            >
              登录
            </BaseButton>

            <!-- OAuth2 登录（根据配置显示） -->
            <template v-if="oauth2Config.linuxdo.enabled">
              <!-- 分隔线 -->
              <div class="relative flex items-center justify-center my-4">
                <div class="border-t border-border-primary w-full absolute"></div>
                <span class="bg-bg-elevated px-3 text-sm text-text-tertiary relative z-10">或</span>
              </div>

              <!-- Linux.do OAuth2 登录 -->
              <a
                href="/api/v1/auth/linuxdo"
                class="w-full inline-flex items-center justify-center gap-2 px-4 py-2.5 bg-bg-elevated border-2 border-border-primary text-text-secondary font-medium rounded-lg hover:bg-bg-hover hover:border-slate-300 transition-all duration-200 shadow-sm"
              >
                <svg class="w-5 h-5" viewBox="0 0 24 24" fill="none">
                  <path d="M12 2L2 7L12 12L22 7L12 2Z" fill="currentColor" opacity="0.6"/>
                  <path d="M2 17L12 22L22 17" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  <path d="M2 12L12 17L22 12" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
                使用 Linux.do 登录
              </a>
            </template>
          </form>

          <!-- 注册表单 -->
          <form v-else @submit.prevent="handleRegister" class="space-y-4">
            <BaseInput
              v-model="registerForm.username"
              label="用户名"
              placeholder="请输入用户名"
              required
            />

            <BaseInput
              v-model="registerForm.email"
              type="email"
              label="邮箱"
              placeholder="请输入邮箱地址"
              required
            />

            <div class="flex gap-2">
              <BaseInput
                v-model="registerForm.code"
                label="验证码"
                placeholder="请输入验证码"
                required
                class="flex-1"
              />
              <BaseButton
                type="button"
                variant="secondary"
                :disabled="countdown > 0"
                @click="sendVerificationCode"
                class="mt-6"
              >
                {{ countdown > 0 ? `${countdown}s` : '获取验证码' }}
              </BaseButton>
            </div>

            <BaseInput
              v-model="registerForm.password"
              type="password"
              label="密码"
              placeholder="请输入密码（至少6位）"
              required
            />

            <BaseButton
              type="submit"
              variant="primary"
              :loading="loading"
              class="w-full"
            >
              注册
            </BaseButton>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { Clock, Zap, Shield, Gauge } from 'lucide-vue-next'
import { message } from '@/utils/message'
import * as userApi from '@/api/user'
import * as configApi from '@/api/config'
import BaseInput from '@/components/BaseInput'
import BaseButton from '@/components/BaseButton'
import SecureStorage, { STORAGE_KEYS } from '@/utils/storage'

const router = useRouter()
const route = useRoute()
const activeTab = ref<'login' | 'register'>('login')
const loading = ref(false)
const countdown = ref(0)

// OAuth2 配置
const oauth2Config = ref({
  linuxdo: {
    enabled: false
  }
})

// 功能亮点
const features = [
  {
    icon: Zap,
    title: '灵活调度',
    description: '支持每天、每周、每月、间隔、Cron表达式等多种调度方式'
  },
  {
    icon: Shield,
    title: '可靠执行',
    description: '自动执行HTTP任务，完整记录执行日志和响应结果'
  },
  {
    icon: Gauge,
    title: '实时监控',
    description: '查看任务执行状态、响应时间和详细的执行记录'
  }
]

// 登录表单
const loginForm = ref({
  account: '',
  password: ''
})

// 注册表单
const registerForm = ref({
  username: '',
  email: '',
  code: '',
  password: ''
})

// 发送验证码
const sendVerificationCode = async () => {
  if (!registerForm.value.email) {
    message.error('请先输入邮箱地址')
    return
  }

  try {
    await userApi.sendRegistrationCode(registerForm.value.email)
    message.success('验证码已发送到您的邮箱')

    // 开始倒计时
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error: any) {
    message.error(error.response?.data?.message || '发送验证码失败')
  }
}

// 处理登录
const handleLogin = async () => {
  loading.value = true
  try {
    const res = await userApi.login({
      account: loginForm.value.account,
      password: loginForm.value.password
    })

    // 保存token和用户信息到SecureStorage
    SecureStorage.setItem(STORAGE_KEYS.AUTH_TOKEN, res.data.token)
    SecureStorage.setItem(STORAGE_KEYS.AUTH_USER, res.data.user)

    message.success('登录成功')
    await router.push('/')
  } catch (error: any) {
    message.error(error.message || '登录失败')
  } finally {
    loading.value = false
  }
}

// 处理注册
const handleRegister = async () => {
  if (registerForm.value.password.length < 6) {
    message.error('密码至少需要6位')
    return
  }

  loading.value = true
  try {
    await userApi.register({
      username: registerForm.value.username,
      email: registerForm.value.email,
      password: registerForm.value.password,
      code: registerForm.value.code
    })

    message.success('注册成功，请登录')

    // 切换到登录页面并填充用户名
    activeTab.value = 'login'
    loginForm.value.account = registerForm.value.username
    loginForm.value.password = ''

    // 清空注册表单
    registerForm.value = {
      username: '',
      email: '',
      code: '',
      password: ''
    }
  } catch (error: any) {
    message.error(error.response?.data?.message || '注册失败')
  } finally {
    loading.value = false
  }
}

// 处理OAuth2回调
const handleOAuth2Callback = async () => {
  const code = route.query.code as string
  const state = route.query.state as string
  const error = route.query.error

  // 处理OAuth2错误
  if (error) {
    message.error('授权失败，请重试')
    router.replace('/')
    return
  }

  // 处理OAuth2成功回调
  if (code) {
    loading.value = true
    try {
      // 调用封装好的API
      const res = await userApi.linuxdoCallback({
        code,
        state: state || ''
      })

      // 保存token和用户信息到SecureStorage（与普通登录相同）
      SecureStorage.setItem(STORAGE_KEYS.AUTH_TOKEN, res.data.token)
      SecureStorage.setItem(STORAGE_KEYS.AUTH_USER, res.data.user)

      message.success('登录成功')
      await router.replace('/')
    } catch (error: any) {
      message.error(error.message || 'OAuth2登录失败')
      router.replace('/')
    } finally {
      loading.value = false
    }
  }
}

// 获取公开配置
const loadPublicConfig = async () => {
  try {
    const res = await configApi.getPublicConfig()
    oauth2Config.value = res.data.oauth2
  } catch (error) {
    console.error('获取配置失败:', error)
    // 配置获取失败不影响基本登录功能
  }
}

// 组件挂载时处理OAuth2回调和加载配置
onMounted(() => {
  loadPublicConfig()
  handleOAuth2Callback()
})
</script>

<style scoped>
@keyframes blob {
  0% {
    transform: translate(0px, 0px) scale(1);
  }
  33% {
    transform: translate(30px, -50px) scale(1.1);
  }
  66% {
    transform: translate(-20px, 20px) scale(0.9);
  }
  100% {
    transform: translate(0px, 0px) scale(1);
  }
}

@keyframes fade-in-up {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fade-in-left {
  from {
    opacity: 0;
    transform: translateX(-20px);
  }
  to {
    opacity: 1;
    transform: translateX(0);
  }
}

@keyframes bounce-subtle {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.animate-blob {
  animation: blob 7s infinite;
}

.animation-delay-2000 {
  animation-delay: 2s;
}

.animation-delay-4000 {
  animation-delay: 4s;
}

.animate-fade-in-up {
  animation: fade-in-up 0.6s ease-out forwards;
  opacity: 0;
}

.animate-fade-in-left {
  animation: fade-in-left 0.8s ease-out forwards;
}

.animate-bounce-subtle {
  animation: bounce-subtle 3s ease-in-out infinite;
}
</style>
