import request from '@/utils/request'
import type { Task } from './task'

// 管理员登录请求
export interface AdminLoginRequest {
  password: string
}

// 管理员登录响应
export interface AdminLoginResponse {
  token: string
  expires_in: number
}

// 任务查询参数
export interface TaskQueryParams {
  page?: number
  page_size?: number
  user_id?: string
  status?: string
  keyword?: string
}

// 任务列表响应
export interface TaskListResponse {
  total: number
  tasks: Task[]
}

// 用户活动
export interface UserActivity {
  user_id: string
  task_count: number
  last_active: string
}

// 统计数据响应
export interface StatsResponse {
  total_users: number
  total_tasks: number
  active_tasks: number
  today_executions: number
  success_rate: number
  recent_users: UserActivity[]
}

// 管理员登录
export const login = (data: AdminLoginRequest) => {
  return request.post<AdminLoginResponse>('/api/v1/admin/login', data)
}

// 管理员登出
export const logout = () => {
  return request.post('/api/v1/admin/logout')
}

// 获取所有任务
export const getTasks = (params: TaskQueryParams) => {
  return request.get<TaskListResponse>('/api/v1/admin/tasks', { params })
}

// 更新任务状态
export const updateTaskStatus = (id: string, enabled: boolean) => {
  return request.put(`/api/v1/admin/tasks/${id}/status`, { enabled })
}

// 删除任务
export const deleteTask = (id: string) => {
  return request.delete(`/api/v1/admin/tasks/${id}`)
}

// 立即执行任务
export const executeTask = (id: string) => {
  return request.post(`/api/v1/admin/tasks/${id}/execute`)
}

// 获取统计数据
export const getStats = () => {
  return request.get<StatsResponse>('/api/v1/admin/stats')
}

// 执行记录查询参数
export interface ExecutionQueryParams {
  page?: number
  page_size?: number
  user_id?: string
  task_id?: string
  status?: string
}

// 执行记录列表响应
export interface ExecutionListResponse {
  total: number
  executions: any[]
}

// 获取执行记录
export const getExecutions = (params: ExecutionQueryParams) => {
  return request.get<ExecutionListResponse>('/api/v1/admin/executions', { params })
}

// 用户接口
export interface User {
  id: string
  username: string
  email: string
  status: number
  role: number
  created_at: string
  updated_at: string
  total_tasks: number
  enabled_tasks: number
}

// 用户查询参数
export interface UserQueryParams {
  page?: number
  page_size?: number
  keyword?: string
  status?: number
}

// 用户列表响应
export interface UserListResponse {
  total: number
  users: User[]
}

// 获取用户列表
export const getUsers = (params: UserQueryParams) => {
  return request.get<UserListResponse>('/api/v1/admin/users', { params })
}

// 更新用户状态
export const updateUserStatus = (id: string, status: number) => {
  return request.put(`/api/v1/admin/users/${id}/status`, { status })
}
