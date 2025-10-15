import request from '@/utils/request'


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


export const getPublicConfig = () => {
  return request.get<PublicConfigResponse>('/api/v1/config')
}
