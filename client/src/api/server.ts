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
  releaseServer: function (data: any) {
    return HttpReq({
      url: '/release/server',
      data
    }) as unknown as BasicResp<any>
  },
  getServerList: function () {
    return HttpReq({
      url: '/main/queryServantGroup',
      method: 'get'
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
    })
  },
  getLogFileList: function (params: any) {
    return HttpReq({
      url: '/statlog/getLogFileList',
      method: 'get',
      params
    })
  },
  getLog: function (data: any) {
    return HttpReq({
      url: '/statlog/getLog',
      method: 'post',
      data
    })
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
  getGroup() {
    return HttpReq({
      url: '/devops/getGroups',
      method: 'get'
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
  getServants() {
    return HttpReq({
      url: '/devops/getServants',
      method: 'get'
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
  }
}

export default api
