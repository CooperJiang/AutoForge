<template>
  <div class="max-w-6xl mx-auto px-4 py-6">
    <div class="grid grid-cols-1 lg:grid-cols-4 gap-6">
      <!-- 左侧菜单 -->
      <div class="lg:col-span-1">
        <div class="bg-white rounded-xl shadow-lg border-2 border-slate-200 p-6 sticky top-6">
          <!-- 用户信息 -->
          <div class="text-center mb-6">
            <div class="w-24 h-24 mx-auto rounded-full bg-gradient-to-br from-emerald-400 to-cyan-500 flex items-center justify-center text-white text-4xl font-bold shadow-lg mb-4">
              {{ userInitial }}
            </div>
            <h2 class="text-xl font-bold text-slate-900 mb-1">{{ userName }}</h2>
            <p class="text-sm text-slate-600 mb-3">{{ userEmail }}</p>
            <span :class="['inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold', roleBadgeClass]">
              {{ roleText }}
            </span>
          </div>

          <!-- 菜单列表 -->
          <nav class="space-y-2">
            <button
              v-for="item in menuItems"
              :key="item.id"
              @click="activeTab = item.id"
              :class="[
                'w-full flex items-center gap-3 px-4 py-3 rounded-lg text-sm font-medium transition-colors',
                activeTab === item.id
                  ? 'bg-emerald-50 text-emerald-700 border-2 border-emerald-200'
                  : 'text-slate-700 hover:bg-slate-50 border-2 border-transparent'
              ]"
            >
              <component :is="item.icon" class="w-5 h-5" />
              <span>{{ item.label }}</span>
            </button>
          </nav>
        </div>
      </div>

      <!-- 右侧内容区 -->
      <div class="lg:col-span-3">
        <!-- 个人资料 -->
        <div v-show="activeTab === 'profile'" class="bg-white rounded-xl shadow-lg border-2 border-slate-200 p-8">
          <div class="mb-6">
            <h3 class="text-2xl font-bold text-slate-900 mb-2">个人资料</h3>
            <p class="text-sm text-slate-600">管理您的用户名和邮箱信息</p>
          </div>

          <!-- 修改用户名 -->
          <div class="mb-8 pb-8 border-b-2 border-slate-200">
            <h4 class="text-lg font-semibold text-slate-900 mb-4">修改用户名</h4>
            <form @submit.prevent="handleUpdateUsername" class="space-y-4">
              <div>
                <label class="block text-sm font-semibold text-slate-700 mb-2">当前用户名</label>
                <div class="px-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-lg text-slate-900">
                  {{ userName }}
                </div>
              </div>
              <BaseInput
                v-model="usernameForm.username"
                label="新用户名"
                placeholder="输入新用户名（2-20个字符）"
              />
              <div class="flex justify-end gap-3">
                <BaseButton type="button" variant="secondary" @click="usernameForm.username = ''" size="sm">
                  重置
                </BaseButton>
                <BaseButton type="submit" variant="primary" :disabled="updatingUsername" size="sm">
                  {{ updatingUsername ? '保存中...' : '保存修改' }}
                </BaseButton>
              </div>
            </form>
          </div>

          <!-- 修改邮箱 -->
          <div>
            <h4 class="text-lg font-semibold text-slate-900 mb-4">修改邮箱</h4>
            <form @submit.prevent="handleUpdateEmail" class="space-y-4">
              <div>
                <label class="block text-sm font-semibold text-slate-700 mb-2">当前邮箱</label>
                <div class="px-4 py-3 bg-slate-50 border-2 border-slate-200 rounded-lg text-slate-900">
                  {{ userEmail }}
                </div>
              </div>
              <BaseInput
                v-model="emailForm.email"
                label="新邮箱地址"
                type="email"
                placeholder="输入新的邮箱地址"
              />
              <div>
                <label class="block text-sm font-semibold text-slate-700 mb-2">邮箱验证码</label>
                <div class="flex gap-3">
                  <BaseInput
                    v-model="emailForm.code"
                    placeholder="输入6位验证码"
                    class="flex-1"
                  />
                  <BaseButton
                    type="button"
                    variant="secondary"
                    @click="handleSendEmailCode"
                    :disabled="sendingCode || countdown > 0"
                    class="whitespace-nowrap"
                    size="sm"
                  >
                    {{ countdown > 0 ? `${countdown}秒` : sendingCode ? '发送中...' : '发送验证码' }}
                  </BaseButton>
                </div>
              </div>
              <div class="flex justify-end gap-3">
                <BaseButton type="button" variant="secondary" @click="emailForm = { email: '', code: '' }" size="sm">
                  重置
                </BaseButton>
                <BaseButton type="submit" variant="primary" :disabled="updatingEmail" size="sm">
                  {{ updatingEmail ? '保存中...' : '保存修改' }}
                </BaseButton>
              </div>
            </form>
          </div>
        </div>

        <!-- 修改密码 -->
        <div v-show="activeTab === 'password'" class="bg-white rounded-xl shadow-lg border-2 border-slate-200 p-8">
          <div class="mb-6">
            <h3 class="text-2xl font-bold text-slate-900 mb-2">修改密码</h3>
            <p class="text-sm text-slate-600">设置一个安全性更高的密码来保护您的账户</p>
          </div>
          <form @submit.prevent="handleChangePassword" class="space-y-6">
            <BaseInput
              v-model="passwordForm.oldPassword"
              label="当前密码"
              type="password"
              placeholder="输入当前密码"
              required
            />
            <BaseInput
              v-model="passwordForm.newPassword"
              label="新密码"
              type="password"
              placeholder="输入新密码（6-20位）"
              required
            />
            <BaseInput
              v-model="passwordForm.confirmPassword"
              label="确认新密码"
              type="password"
              placeholder="再次输入新密码"
              required
            />
            <div class="bg-amber-50 border-2 border-amber-200 rounded-lg p-4">
              <p class="text-sm text-amber-800">
                <strong>提示：</strong>修改密码成功后，您需要重新登录。
              </p>
            </div>
            <div class="flex justify-end gap-3">
              <BaseButton type="button" variant="secondary" @click="passwordForm = { oldPassword: '', newPassword: '', confirmPassword: '' }">
                重置
              </BaseButton>
              <BaseButton type="submit" variant="primary" :disabled="changingPassword">
                {{ changingPassword ? '修改中...' : '修改密码' }}
              </BaseButton>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Mail, Lock } from 'lucide-vue-next'
import BaseButton from '@/components/BaseButton'
import BaseInput from '@/components/BaseInput'
import SecureStorage, { STORAGE_KEYS } from '@/utils/storage'
import { message } from '@/utils/message'
import * as userApi from '@/api/user'

const router = useRouter()

// 角色常量
const USER_ROLE = {
  SUPER_ADMIN: 1,
  ADMIN: 2,
  USER: 3
}

// 菜单项
const menuItems = [
  { id: 'profile', label: '个人资料', icon: User },
  { id: 'password', label: '修改密码', icon: Lock }
]

// 当前激活的标签
const activeTab = ref('profile')

// 从 SecureStorage 获取用户信息
const user = computed(() => SecureStorage.getItem<{
  username: string
  id: string
  email: string
  role: number
}>(STORAGE_KEYS.AUTH_USER))

// 用户名
const userName = computed(() => user.value?.username || '未知用户')

// 邮箱
const userEmail = computed(() => user.value?.email || '未设置')

// 用户名首字母（大写）
const userInitial = computed(() => {
  const name = userName.value
  return name ? name.charAt(0).toUpperCase() : '?'
})

// 角色文本
const roleText = computed(() => {
  const role = user.value?.role
  if (role === USER_ROLE.SUPER_ADMIN) return '超级管理员'
  if (role === USER_ROLE.ADMIN) return '管理员'
  return '普通用户'
})

// 角色徽章样式
const roleBadgeClass = computed(() => {
  const role = user.value?.role
  if (role === USER_ROLE.SUPER_ADMIN) return 'bg-purple-100 text-purple-700'
  if (role === USER_ROLE.ADMIN) return 'bg-blue-100 text-blue-700'
  return 'bg-emerald-100 text-emerald-700'
})

// 修改用户名表单
const usernameForm = ref({
  username: ''
})
const updatingUsername = ref(false)

// 修改邮箱表单
const emailForm = ref({
  email: '',
  code: ''
})
const updatingEmail = ref(false)
const sendingCode = ref(false)
const countdown = ref(0)

// 修改密码表单
const passwordForm = ref({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})
const changingPassword = ref(false)

// 更新用户名
const handleUpdateUsername = async () => {
  if (!usernameForm.value.username.trim()) {
    message.error('请输入用户名')
    return
  }

  if (usernameForm.value.username === userName.value) {
    message.warning('用户名未修改')
    return
  }

  updatingUsername.value = true
  try {
    await userApi.updateProfile({ username: usernameForm.value.username })

    // 更新本地存储
    const currentUser = user.value
    if (currentUser) {
      currentUser.username = usernameForm.value.username
      SecureStorage.setItem(STORAGE_KEYS.AUTH_USER, currentUser)
    }

    message.success('用户名修改成功')
    usernameForm.value.username = ''
  } catch (error: any) {
    message.error(error.response?.data?.message || '修改失败')
  } finally {
    updatingUsername.value = false
  }
}

// 发送邮箱验证码
const handleSendEmailCode = async () => {
  if (!emailForm.value.email.trim()) {
    message.error('请输入邮箱地址')
    return
  }

  if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(emailForm.value.email)) {
    message.error('请输入有效的邮箱地址')
    return
  }

  sendingCode.value = true
  try {
    await userApi.sendChangeEmailCode(emailForm.value.email)
    message.success('验证码已发送，请查收邮箱')

    // 启动倒计时
    countdown.value = 60
    const timer = setInterval(() => {
      countdown.value--
      if (countdown.value <= 0) {
        clearInterval(timer)
      }
    }, 1000)
  } catch (error: any) {
    message.error(error.response?.data?.message || '发送失败')
  } finally {
    sendingCode.value = false
  }
}

// 更新邮箱
const handleUpdateEmail = async () => {
  if (!emailForm.value.email.trim() || !emailForm.value.code.trim()) {
    message.error('请填写完整信息')
    return
  }

  updatingEmail.value = true
  try {
    await userApi.updateProfile({
      email: emailForm.value.email,
      code: emailForm.value.code
    })

    // 更新本地存储
    const currentUser = user.value
    if (currentUser) {
      currentUser.email = emailForm.value.email
      SecureStorage.setItem(STORAGE_KEYS.AUTH_USER, currentUser)
    }

    message.success('邮箱修改成功')
    emailForm.value.email = ''
    emailForm.value.code = ''
  } catch (error: any) {
    message.error(error.response?.data?.message || '修改失败')
  } finally {
    updatingEmail.value = false
  }
}

// 修改密码
const handleChangePassword = async () => {
  if (!passwordForm.value.oldPassword || !passwordForm.value.newPassword || !passwordForm.value.confirmPassword) {
    message.error('请填写完整信息')
    return
  }

  if (passwordForm.value.newPassword.length < 6 || passwordForm.value.newPassword.length > 20) {
    message.error('新密码长度必须在6-20位之间')
    return
  }

  if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
    message.error('两次输入的新密码不一致')
    return
  }

  changingPassword.value = true
  try {
    await userApi.changePassword({
      oldPassword: passwordForm.value.oldPassword,
      newPassword: passwordForm.value.newPassword
    })

    message.success('密码修改成功，请重新登录')

    // 清空表单
    passwordForm.value = {
      oldPassword: '',
      newPassword: '',
      confirmPassword: ''
    }

    // 延迟后退出登录
    setTimeout(() => {
      SecureStorage.removeItem(STORAGE_KEYS.AUTH_TOKEN)
      SecureStorage.removeItem(STORAGE_KEYS.AUTH_USER)
      router.push('/auth')
    }, 1500)
  } catch (error: any) {
    message.error(error.response?.data?.message || '修改失败')
  } finally {
    changingPassword.value = false
  }
}

// 检查是否已登录
onMounted(() => {
  if (!user.value) {
    message.warning('请先登录')
    router.push('/auth')
  }
})
</script>
