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
  getServerList: function (id:number) {
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
  LoginByCache:function(data:any){
    return HttpReq({
      url:'/loginByCache',
      method:'post',
      data
    }) as unknown as Promise<BasicResp<any>>
  },
  getUserMenusByUserId:function(id:number){
    return HttpReq({
      url:'/getUserMenusByUserId',
      method:'get',
      params:{id}
    }) as unknown as Promise<BasicResp<any>>
  },
  getGroup(id:number) {
    return HttpReq({
      url: '/devops/getGroups',
      method: 'get',
      params:{id}
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
  getServants(id:number) {
    return HttpReq({
      url: '/devops/getServants',
      method: 'get',
      params:{id}
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
  updatePackageVersion(params:any){
    return HttpReq({
      url: '/release/updatePackageVersion',
      method: 'get',
      params
    }) as unknown as BasicResp<any>
  },
  getPropertys(data:any){
    return HttpReq({
      url: '/devops/getPropertys',
      method: 'post',
      data
    }) as unknown as BasicResp<any>
  },
  setProperty(data:any){
    return HttpReq({
      url: '/devops/setProperty',
      method: 'post',
      data
    }) as unknown as BasicResp<any>
  },
  delProperty(id:number){
    return HttpReq({
      url: '/devops/delProperty',
      method: 'get',
      params:{id}
    }) as unknown as BasicResp<any>
  },
  delServant(id:number){
    return HttpReq({
      url: '/devops/delServant',
      method: 'post',
      params:{id}
    }) as unknown as BasicResp<any>
  }
}

export default api
