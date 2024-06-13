import type { BasicResp } from '@/dto/dto'
import HttpReq from '@/utils/request'

const api = {
    getList: function (data: any) {
        return HttpReq({
            url: '/assets/admin/getList',
            data,
        }) as unknown as BasicResp<any>
    },
    upsertAsset: function (data: any) {
        return HttpReq({
            url: '/assets/admin/upsertAsset',
            data,
        }) as unknown as BasicResp<any>
    },
    delAssert: function (params: any) {
        return HttpReq({
            url: '/assets/admin/delAssert',
            params,
            method: 'get',
        }) as unknown as BasicResp<any>
    },
}

export default api
