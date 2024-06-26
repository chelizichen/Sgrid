import { constants } from '@/constant'
import axios from 'axios'

const HttpReq = axios.create({
  // baseURL, // api的base_url
  timeout: 15000, // 请求超时时间,
  method: 'post',
  baseURL: '/sgirdcloud/'
})
HttpReq.interceptors.response.use((resp) => {
  if (resp.data.code != 0) {
    console.log('resp.data.Message |', resp.data.Message)
    // this.$message.error(resp.data.Message)
  }
  return resp.data
})

HttpReq.interceptors.request.use((config) => {
  const tkn = localStorage.getItem(constants.TOKEN)
  config.headers['token'] = tkn
  return config
})

export default HttpReq
