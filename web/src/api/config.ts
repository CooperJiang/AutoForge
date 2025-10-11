import request from '@/utils/request'

// 公开配置响应
export interface PublicConfigResponse {
  oauth2: {
    linuxdo: {
      enabled: boolean
    }
  }
  app: {
    name: string
  }
}

// 获取公开配置
export const getPublicConfig = () => {
  return request.get<PublicConfigResponse>('/api/v1/config')
}
