import request from '@/utils/request'
import type { Task } from './task'

export interface AdminLoginRequest {
  password: string
}

export interface AdminLoginResponse {
  token: string
  expires_in: number
}

export interface TaskQueryParams {
  page?: number
  page_size?: number
  user_id?: string
  status?: string
  keyword?: string
}

export interface TaskListResponse {
  total: number
  tasks: Task[]
}

export interface UserActivity {
  user_id: string
  task_count: number
  last_active: string
}

export interface StatsResponse {
  total_users: number
  total_tasks: number
  active_tasks: number
  today_executions: number
  success_rate: number
  total_workflows: number
  total_templates: number
  recent_users: UserActivity[]
}

export const login = (data: AdminLoginRequest) => {
  return request.post<AdminLoginResponse>('/api/v1/admin/login', data)
}

export const logout = () => {
  return request.post('/api/v1/admin/logout')
}

export const getTasks = (params: TaskQueryParams) => {
  return request.get<TaskListResponse>('/api/v1/admin/tasks', { params })
}

export const updateTaskStatus = (id: string, enabled: boolean) => {
  return request.put(`/api/v1/admin/tasks/${id}/status`, { enabled })
}

export const deleteTask = (id: string) => {
  return request.delete(`/api/v1/admin/tasks/${id}`)
}

export const executeTask = (id: string) => {
  return request.post(`/api/v1/admin/tasks/${id}/execute`)
}

export const getStats = () => {
  return request.get<StatsResponse>('/api/v1/admin/stats')
}

export interface ExecutionQueryParams {
  page?: number
  page_size?: number
  user_id?: string
  task_id?: string
  status?: string
}

export interface ExecutionListResponse {
  total: number
  executions: any[]
}

export const getExecutions = (params: ExecutionQueryParams) => {
  return request.get<ExecutionListResponse>('/api/v1/admin/executions', { params })
}

export const deleteExecution = (id: string) => {
  return request.delete(`/api/v1/admin/executions/${id}`)
}

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

export interface UserQueryParams {
  page?: number
  page_size?: number
  keyword?: string
  status?: number
}

export interface UserListResponse {
  total: number
  users: User[]
}

export const getUsers = (params: UserQueryParams) => {
  return request.get<UserListResponse>('/api/v1/admin/users', { params })
}

export const updateUserStatus = (id: string, status: number) => {
  return request.put(`/api/v1/admin/users/${id}/status`, { status })
}
