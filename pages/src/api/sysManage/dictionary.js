import request from '@/utils/request'

export function getDicList() {
  return request({
    url: '/rest/dic/list',
    method: 'post'
  })
}

export function saveDic(formData) {
  return request({
    url: '/rest/dic/save',
    method: 'post',
    data: formData
  })
}

export function deleteDic(id) {
  return request({
    url: '/rest/dic/delete',
    method: 'post',
    data: id
  })
}
