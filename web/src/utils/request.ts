import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import type { BaseResponse } from '@/types'
import { useMessage } from '@/composables/useMessage'
import SecureStorage, { STORAGE_KEYS } from './storage'

const request: AxiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
})

const { error: showError } = useMessage()

request.interceptors.request.use(
  (config) => {
    const token = SecureStorage.getItem<string>(STORAGE_KEYS.AUTH_TOKEN)
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }

    if (config.data instanceof FormData) {
      delete config.headers['Content-Type']
    }

    if (config.method === 'get') {
      config.params = {
        ...config.params,
        _t: Date.now(),
      }
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

request.interceptors.response.use(
  (response: AxiosResponse<BaseResponse>) => {
    const { data } = response

    if (data.code !== 200) {
      const errorMessage = (data as any).detail || data.message || '请求失败'
      const error: any = new Error(errorMessage)

      error.response = response
      error.code = data.code
      error.message = errorMessage
      error.detail = (data as any).detail

      return Promise.reject(error)
    }

    return {
      ...response,
      data: data.data,
    } as AxiosResponse
  },
  (error) => {
    let message = '网络错误'
    let shouldShowError = true

    if (error.response) {
      const { status, data } = error.response

      if (status === 401) {
        message = '登录已过期，请重新登录'
        shouldShowError = false

        SecureStorage.removeItem(STORAGE_KEYS.AUTH_TOKEN)
        SecureStorage.removeItem(STORAGE_KEYS.AUTH_USER)

        setTimeout(() => {
          if (window.location.pathname !== '/auth') {
            window.location.replace('/auth')
          }
        }, 100)
        const customError = new Error(message)
        return Promise.reject(customError)
      }

      if (data?.detail) {
        message = data.detail
      } else if (data?.message) {
        message = data.message
      } else {
        switch (status) {
          case 403:
            message = '拒绝访问'
            break
          case 404:
            message = '请求的资源不存在'
            break
          case 500:
            message = '服务器内部错误'
            break
          default:
            message = `请求失败 (${status})`
        }
      }
    } else if (error.request) {
      message = '网络连接失败，请检查网络连接'
    }

    if (shouldShowError) {
      showError(message)
    }

    const customError = new Error(message)
    return Promise.reject(customError)
  }
)

export const ApiClient = {
  async get<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await request.get(url, config)
    return response.data
  },

  async post<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await request.post(url, data, config)
    return response.data
  },

  async put<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await request.put(url, data, config)
    return response.data
  },

  async delete<T = unknown>(url: string, config?: AxiosRequestConfig): Promise<T> {
    const response = await request.delete(url, config)
    return response.data
  },

  async patch<T = unknown>(url: string, data?: unknown, config?: AxiosRequestConfig): Promise<T> {
    const response = await request.patch(url, data, config)
    return response.data
  },
}

export default request
