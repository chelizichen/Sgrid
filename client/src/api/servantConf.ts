import type { BasicResp } from '@/dto/dto'
import HttpReq from '@/utils/request'

export function getConfig(params: { id: number }) {
  return HttpReq({
    method: 'get',
    params: params,
    url: '/devops/getConfig'
  })
}
export function updateConfig(data: any) {
  return HttpReq({
    method: 'post',
    data: data,
    url: '/devops/updateConfig'
  })
}
