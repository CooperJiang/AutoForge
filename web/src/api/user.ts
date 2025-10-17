import request from '@/utils/request'

export interface RegisterRequest {
  username: string
  email: string
  password: string
  code: string
}

export interface LoginRequest {
  account: string
  password: string
}

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

export const sendRegistrationCode = (email: string) => {
  return request.post('/api/v1/user/send-registration-code', { email })
}

export const register = (data: RegisterRequest) => {
  return request.post('/api/v1/user/register', data)
}

export const login = (data: LoginRequest) => {
  return request.post<LoginResponse>('/api/v1/user/login', data)
}

export const getUserInfo = () => {
  return request.get('/api/v1/user/info')
}

export interface UpdateProfileRequest {
  username?: string
  email?: string
  avatar?: string
  code?: string
}

export const updateProfile = (data: UpdateProfileRequest) => {
  return request.put('/api/v1/user/profile', data)
}

export interface ChangePasswordRequest {
  oldPassword: string
  newPassword: string
}

export const changePassword = (data: ChangePasswordRequest) => {
  return request.post('/api/v1/user/change-password', data)
}

export const sendChangeEmailCode = (email: string) => {
  return request.post('/api/v1/user/send-change-email-code', { email })
}

export interface OAuth2CallbackRequest {
  code: string
  state?: string
}

export const linuxdoCallback = (data: OAuth2CallbackRequest) => {
  return request.post<LoginResponse>('/api/v1/auth/linuxdo/callback', data)
}

export const githubCallback = (data: OAuth2CallbackRequest) => {
  return request.post<LoginResponse>('/api/v1/auth/github/callback', data)
}
