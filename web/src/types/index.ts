export interface BaseResponse<T = unknown> {
  code: number
  message: string
  data: T
  request_id?: string
  timestamp?: string
}

export interface PaginationParams {
  page?: number
  size?: number
}

export interface PaginationResponse<T> {
  items: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

export interface BaseModel {
  id: string
  created_at: string
  updated_at: string
}

export interface User {
  id?: string
  username: string
  email: string
  status: number
  role?: number
  avatar?: string
  bio?: string
  created_at?: string
  updated_at?: string
}

export interface LoginRequest {
  account: string
  password: string
}

export interface LoginResponse {
  token: string
  user: User
}

export interface RegisterRequest {
  username: string
  email: string
  password: string
  code: string
}

export interface RegisterResponse {
  message: string
  user: User
}

export interface SendCodeRequest {
  email: string
}

export interface SendCodeResponse {
  message: string
}

export interface ResetPasswordRequest {
  email: string
  code: string
  newPassword: string
}

export interface ChangePasswordRequest {
  oldPassword: string
  newPassword: string
}

export interface UpdateProfileRequest {
  username?: string
  email?: string
  avatar?: string
  code?: string
}

export enum Status {
  INACTIVE = 0,
  ACTIVE = 1,
}

export enum UserRole {
  SUPER_ADMIN = 1,
  ADMIN = 2,
  USER = 3,
}

export interface FormField {
  name: string
  label: string
  type: 'text' | 'email' | 'password' | 'number' | 'select' | 'textarea'
  required?: boolean
  placeholder?: string
  options?: { label: string; value: string | number }[]
  rules?: Array<(value: unknown) => boolean | string>
}

export interface MenuItem {
  id: string
  title: string
  icon?: string
  path?: string
  children?: MenuItem[]
}

export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
  success: boolean
}

export interface ApiError {
  code: number
  message: string
  details?: Record<string, unknown>
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  limit: number
}

export interface PageQuery {
  page?: number
  limit?: number
  search?: string
}

export interface PageResponse<T> {
  list: T[]
  total: number
  page: number
  limit: number
}

export interface RequestConfig {
  timeout?: number
  headers?: Record<string, string>
  params?: Record<string, unknown>
  data?: unknown
}

export interface HttpResponse<T = unknown> {
  data: T
  status: number
  statusText: string
  headers: Record<string, string>
}
