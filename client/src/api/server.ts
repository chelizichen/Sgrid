import type { BasicResp } from '@/dto/dto'
import HttpReq from '@/utils/request'

const api = {
  uploadSgridServer: function (data: any) {
    return HttpReq({
      url: '/upload/uploadServer',
      data,
      timeout: 10 * 60 * 1000
    }) as unknown as BasicResp<any>
  },
  releaseServer: function (data: any, params: { releaseId: string }) {
    return HttpReq({
      url: '/release/server',
      data,
      params
    }) as unknown as BasicResp<any>
  },
  restartServer: function (data: any) {
    return HttpReq({
      url: '/restart/server',
      data
    }) as unknown as BasicResp<any>
  },
  getServerList: function (id: number) {
    return HttpReq({
      url: '/main/queryServantGroup',
      method: 'post',
      data: {
        id
      }
    }) as unknown as BasicResp<any>
  },
  getUploadList: function (params: any) {
    return HttpReq({
      url: '/upload/getList',
      params,
      method: 'get'
    }) as unknown as BasicResp<any>
  },
  queryGrid: function (params: any) {
    return HttpReq({
      url: '/main/queryGrid',
      params,
      method: 'get'
    })
  },
  shutdownServer: function (data: any) {
    return HttpReq({
      url: '/release/shutdown',
      data
    })
  },
  getStatLogList: function (params: any) {
    return HttpReq({
      url: '/statlog/getlist',
      method: 'get',
      params
    }) as unknown as BasicResp<any>
  },
  getLogFileList: function (params: any) {
    return HttpReq({
      url: '/statlog/getLogFileList',
      method: 'get',
      params
    }) as unknown as BasicResp<any>
  },
  getLog: function (data: any) {
    return HttpReq({
      url: '/statlog/getLog',
      method: 'post',
      data
    }) as unknown as BasicResp<any>
  },
  checkStat: function (data: any) {
    return HttpReq({
      url: '/statlog/check',
      method: 'post',
      data
    })
  },
  Login: function (data: any) {
    return HttpReq({
      url: '/login',
      data
    }) as unknown as BasicResp<any>
  },
  LoginByCache: function (data: any) {
    return HttpReq({
      url: '/loginByCache',
      method: 'post',
      data
    }) as unknown as Promise<BasicResp<any>>
  },
  getUserMenusByUserId: function (id: number) {
    return HttpReq({
      url: '/getUserMenusByUserId',
      method: 'get',
      params: { id }
    }) as unknown as Promise<BasicResp<any>>
  },
  getGroup(id: number) {
    return HttpReq({
      url: '/devops/getGroups',
      method: 'get',
      params: { id }
    })
  },
  saveServant(data: any) {
    return HttpReq({
      url: '/devops/saveServant',
      method: 'post',
      data: data
    }) as unknown as BasicResp<any>
  },
  queryNodes() {
    return HttpReq({
      url: '/devops/queryNodes',
      method: 'get'
    })
  },
  getServants(id?: number) {
    return HttpReq({
      url: '/devops/getServants',
      method: 'get',
      params: { id }
    })
  },
  saveGrid(data: any) {
    return HttpReq({
      url: '/devops/saveGrid',
      method: 'post',
      data
    }) as unknown as BasicResp<any>
  },
  deleteGrid(data: any) {
    return HttpReq({
      url: '/devops/deleteGrid',
      method: 'post',
      data
    }) as unknown as BasicResp<any>
  },
  saveGroup(data: any) {
    return HttpReq({
      url: '/devops/saveGroup',
      method: 'post',
      data
    }) as unknown as BasicResp<any>
  },
  saveNode(data: any) {
    return HttpReq({
      url: '/devops/saveNode',
      method: 'post',
      data
    }) as unknown as BasicResp<any>
  },
  updatePackageVersion(params: any) {
    return HttpReq({
      url: '/release/updatePackageVersion',
      method: 'get',
      params
    }) as unknown as BasicResp<any>
  },
  getPropertys(data: any) {
    return HttpReq({
      url: '/devops/getPropertys',
      method: 'post',
      data
    }) as unknown as BasicResp<any>
  },
  getPropertyByKey(key: string) {
    return HttpReq({
      url: '/devops/getPropertyByKey',
      method: 'post',
      params: {
        key
      }
    }) as unknown as BasicResp<any>
  },
  setProperty(data: any) {
    return HttpReq({
      url: '/devops/setProperty',
      method: 'post',
      data
    }) as unknown as BasicResp<any>
  },
  delProperty(id: number) {
    return HttpReq({
      url: '/devops/delProperty',
      method: 'get',
      params: { id }
    }) as unknown as BasicResp<any>
  },
  delServant(id: number, stat: number) {
    return HttpReq({
      url: '/devops/delServant',
      method: 'post',
      params: { id, stat }
    }) as unknown as BasicResp<any>
  },
  deleteGroup(id: number) {
    return HttpReq({
      url: '/devops/delGroup',
      method: 'post',
      params: { id }
    }) as unknown as BasicResp<any>
  },
  getRandomPort() {
    return HttpReq({
      url: '/main/port/random',
      method: 'get'
    }) as unknown as BasicResp<any>
  },
  updateServant(data: any) {
    return HttpReq({
      url: '/devops/updateServant',
      method: 'post',
      data
    }) as unknown as BasicResp<any>
  },
  delPackage(params: any) {
    return HttpReq({
      url: '/upload/removePackage',
      method: 'get',
      params
    }) as unknown as BasicResp<any>
  },
  downLoadPackage(serverName: string, fileName: string) {
    return HttpReq({
      url: '/download/serverPackage',
      method: 'get',
      params: { serverName, fileName }
    }) as unknown as BasicResp<any>
  },
  getMainLogger(data: any) {
    return HttpReq({
      url: '/main/logger/get',
      method: 'post',
      data
    }) as unknown as BasicResp<any>
  },
  invokeWithCmd(data: any) {
    return HttpReq({
      url: '/main/invokeWithCmd',
      method: 'post',
      data
    }) as unknown as BasicResp<any>
  }
}

export default api
