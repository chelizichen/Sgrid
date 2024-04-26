import type { BasicResp } from '@/dto/dto'
import HttpReq from '@/utils/exp_req'

export function getProxyList() {
  return HttpReq({
    url: '/nginx/getProxyList',
    method: 'get'
  }) as unknown as BasicResp<any>
}

export function merge(data: any) {
  return HttpReq({
    url: '/nginx/merge',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function reload() {
  return HttpReq({
    url: '/nginx/reload',
    method: 'post'
  }) as unknown as BasicResp<any>
}

export function getBackupList() {
  return HttpReq({
    url: '/nginx/getBackupList',
    method: 'get'
  }) as unknown as BasicResp<any>
}

export function getBackupFile(params: any) {
  return HttpReq({
    url: '/nginx/getBackupFile',
    method: 'get',
    params
  }) as unknown as BasicResp<any>
}

export function test() {
  return HttpReq({
    url: '/nginx/test',
    method: 'get'
  }) as unknown as BasicResp<any>
}
