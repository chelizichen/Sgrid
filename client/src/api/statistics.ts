import type { BasicResp } from '@/dto/dto'
import request from '@/utils/request'

export function getStatisticsByType(TYPE: string) {
  return request({
    url: '/server/statistics/getByType',
    method: 'get',
    params: { TYPE }
  }) as unknown as BasicResp<any>
}
