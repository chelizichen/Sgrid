import type { BasicResp } from '@/dto/dto'
import request from '@/utils/request'

export function getUser(data: any) {
  return request({
    url: '/system/user/get',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function getRole(data: any) {
  return request({
    url: '/system/role/get',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function getMenu(data: any) {
  return request({
    url: '/system/menu/get',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function saveUser(data: any) {
  return request({
    url: '/system/user/save',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function saveRole(data: any) {
  return request({
    url: '/system/role/save',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function saveMenu(data: any) {
  return request({
    url: '/system/menu/save',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function delMenu(id: number) {
  return request({
    url: '/system/menu/del',
    method: 'get',
    params: { id }
  }) as unknown as BasicResp<any>
}

export function delRole(id: number) {
  return request({
    url: '/system/role/del',
    method: 'get',
    params: { id }
  }) as unknown as BasicResp<any>
}

export function setUserToRole(data: any) {
  return request({
    url: '/system/setUserToRole',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function setRoleToMenu(data: any) {
  return request({
    url: '/system/setRoleToMenu',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function getUserToRoleRelation(id: number) {
  return request({
    url: '/system/getUserToRoleRelation',
    method: 'get',
    params: { id }
  }) as unknown as BasicResp<any>
}

export function getMenuListByRoleId(id: number) {
  return request({
    url: '/system/getMenuListByRoleId',
    method: 'get',
    params: { id }
  }) as unknown as BasicResp<any>
}
