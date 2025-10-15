import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { authApi } from '../../api/auth'
import { useAuthStorage } from '../common'
import { STORAGE_KEYS } from '../../utils/storage'
import type { LoginRequest, RegisterRequest, ResetPasswordRequest, User } from '../../types'

export function useAuth() {
  const router = useRouter()
  const loading = ref(false)
  const error = ref('')


  const [token, setToken, removeToken] = useAuthStorage<string | null>(
    STORAGE_KEYS.AUTH_TOKEN,
    null
  )
  const [user, setUser, removeUser] = useAuthStorage<User | null>(STORAGE_KEYS.AUTH_USER, null)


  const isAuthenticated = computed(() => !!token.value && !!user.value)


  const login = async (credentials: LoginRequest) => {
    try {
      loading.value = true
      error.value = ''

      const response = await authApi.login(credentials)


      const responseData = (response as any).data || response

      setToken(responseData.token)
      setUser(responseData.user)

      return response
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : '登录失败'
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }


  const register = async (userData: RegisterRequest) => {
    try {
      loading.value = true
      error.value = ''

      const response = await authApi.register(userData)
      return response
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : '注册失败'
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }


  const sendRegistrationCode = async (email: string) => {
    try {
      loading.value = true
      error.value = ''

      const response = await authApi.sendRegistrationCode({ email })
      return response
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : '发送验证码失败'
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }


  const sendResetPasswordCode = async (email: string) => {
    try {
      loading.value = true
      error.value = ''

      const response = await authApi.sendResetPasswordCode({ email })
      return response
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : '发送验证码失败'
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }


  const resetPassword = async (data: ResetPasswordRequest) => {
    try {
      loading.value = true
      error.value = ''

      const response = await authApi.resetPassword(data)
      return response
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : '重置密码失败'
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }


  const logout = (shouldRedirect: boolean = true) => {
    removeToken()
    removeUser()
    if (shouldRedirect) {
      router.push('/login')
    }
  }


  const getUserInfo = async () => {
    try {
      loading.value = true
      const response = await authApi.getUserInfo()

      const userInfo = (response as any).data || response
      setUser(userInfo)
      return userInfo
    } catch (err: unknown) {
      const errorMessage = err instanceof Error ? err.message : '获取用户信息失败'
      error.value = errorMessage
      throw err
    } finally {
      loading.value = false
    }
  }

  return {

    loading,
    error,
    user,
    token,
    isAuthenticated,


    login,
    register,
    sendRegistrationCode,
    sendResetPasswordCode,
    resetPassword,
    logout,
    getUserInfo,
  }
}
