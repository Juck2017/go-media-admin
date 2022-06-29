import request from '@/utils/request'

// 获取通道列表
export function getChannelList(data) {
  return request({
    url: '/rest/channel/list',
    method: 'post',
    data
  })
}

// 获取pusher客户端列表
export function getPusherList(data) {
  return request({
    url: '/rest/pusher/list',
    method: 'post',
    data
  })
}
