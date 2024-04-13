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
  ShutDownServer: function (data, url = '/shutdownServer') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetLogger: function (data, url = '/getServerLog') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetApiJson: function (data, url = '/getApiJson') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetDoc: function (data, url = '/getDoc') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetLogList: function (data, url = '/getLogList') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  CoverConfig: function (data, url = '/coverConfig') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  DeleteServer: function (data, url = '/deleteServer') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetMainLogList: function (data, url = '/main/getLogList') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  GetMainLogger: function (data, url = '/main/getServerLog') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  Login: function (data, url = '/login') {
    return HttpReq({
      url,
      data
    }) as unknown as BasicResp<any>
  },
  getChildStats: function (data, url = '/getChildStats') {
    return HttpReq({
      url,
      data
    })
  },
  queryGrid:function(params:any){
    return HttpReq({
      url:'/main/queryGrid',
      params,
      method:'get',
    })
  }
}

export default API
