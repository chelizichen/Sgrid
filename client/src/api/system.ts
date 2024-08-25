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

export function getGroup(data: any) {
  return request({
    url: '/system/group/get',
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

export function saveUserGroup(data:any){
  return request({
    url: '/system/group/save',
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

export function delUserGroup(id: number) {
  return request({
    url: '/system/group/del',
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

export function setUserToUserGroup(data: any) {
  return request({
    url: '/system/setUserToUserGroup',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function setUserGroupToServantGroup(data: any) {
  return request({
    url: '/system/setUserGroupToServantGroup',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function getServantGroupsByUserGroupId(data) {
  return request({
    url: '/system/spec/getServantGroupsByUserGroupId',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}

export function getUsersByUserGroup(data) {
  return request({
    url: '/system/spec/getUsersByUserGroup',
    method: 'post',
    data
  }) as unknown as BasicResp<any>
}