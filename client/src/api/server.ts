import type { BasicResp } from '@/dto/dto'
import HttpReq from '@/utils/request'

const API = {
  uploadSgridServer: function (data:any) {
    return HttpReq({
      url:'/upload/uploadServer',
      data,
      timeout: 10 * 60 * 1000
    }) as unknown as BasicResp<any>
  },
  releaseServer: function (data:any) {
    return HttpReq({
      url:'/release/server',
      data,
    }) as unknown as BasicResp<any>
  },
  GetServerList: function () {
    return HttpReq({
      url:'/main/queryServantGroup',
      method:'get'
    }) as unknown as BasicResp<any>
  },
  getUploadList: function (params:any) {
    return HttpReq({
      url:'/upload/getList',
      params,
      method:'get',
    }) as unknown as BasicResp<any>
  },
  CheckServer: function (data:any) {
    return HttpReq({
      url:'/checkServer',
      data
    }) as unknown as BasicResp<any>
  },
  CreateServer: function (data:any) {
    return HttpReq({
      url:'/createServer',
      data
    }) as unknown as BasicResp<any>
  },
  CheckConfig: function (data:any) {
    return HttpReq({
      url:'/checkConfig',
      data
    }) as unknown as BasicResp<any>
  },
  DeletePackage: function (data:any) {
    return HttpReq({
      url:'/deletePackage',
      data
    }) as unknown as BasicResp<any>
  },
  queryGrid:function(params:any){
    return HttpReq({
      url:'/main/queryGrid',
      params,
      method:'get',
    })
  },
  shutdownServer:function(data:any){
    return HttpReq({
      url:'/release/shutdown',
      data
    })
  },
  getStatLogList:function(params:any){
    return HttpReq({
      url:'/statlog/getlist',
      method:'get',
      params
    })
  },
  getLogFileList:function(params:any){
    return HttpReq({
      url:'/statlog/getLogFileList',
      method:'get',
      params
    })
  },
  getLog:function(data:any){
    return HttpReq({
      url:'/statlog/getLog',
      method:'post',
      data
    })
  }
}

export default API
