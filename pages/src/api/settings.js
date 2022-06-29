import request from '@/utils/request'

export function getDictMap() {
  return request({
    url: '/rest/dic/getDicMap',
    method: 'post'
  })
}
