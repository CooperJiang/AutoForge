import request from '@/utils/request'

// 注册请求参���
export interface RegisterRequest {
  username: string
  email: string
  password: string
  code: string
}

// 登录请求参数
export interface LoginRequest {
  account: string // 用户名或邮箱
  password: string
}

// 登录响应
export interface LoginResponse {
  user: {
    id: string
    username: string
    email: string
    role: number
  }
  token: string
  expires_at: string
}

// 发送注册验证码
export const sendRegistrationCode = (email: string) => {
  return request.post('/api/v1/user/send-registration-code', { email })
}

// 注册
export const register = (data: RegisterRequest) => {
  return request.post('/api/v1/user/register', data)
}

// 登录
export const login = (data: LoginRequest) => {
  return request.post<LoginResponse>('/api/v1/user/login', data)
}

// 获取用户信息
export const getUserInfo = () => {
  return request.get('/api/v1/user/info')
}

// 更新用户资料请求参数
export interface UpdateProfileRequest {
  username?: string
  email?: string
  avatar?: string
  code?: string
}

// 更新用户资料
export const updateProfile = (data: UpdateProfileRequest) => {
  return request.put('/api/v1/user/profile', data)
}

// 修改密码请求参数
export interface ChangePasswordRequest {
  oldPassword: string
  newPassword: string
}

// 修改密码
export const changePassword = (data: ChangePasswordRequest) => {
  return request.post('/api/v1/user/change-password', data)
}

// 发送修改邮箱验证码
export const sendChangeEmailCode = (email: string) => {
  return request.post('/api/v1/user/send-change-email-code', { email })
}

// OAuth2 回调请求参数
export interface OAuth2CallbackRequest {
  code: string
  state?: string
}

// Linux.do OAuth2 回调
export const linuxdoCallback = (data: OAuth2CallbackRequest) => {
  return request.post<LoginResponse>('/api/v1/auth/linuxdo/callback', data)
}
