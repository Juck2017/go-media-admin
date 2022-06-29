import request from '@/utils/request'

export function getPrivilegeList() {
  return request({
    url: '/rest/privilege/list',
    method: 'post'
  })
}

export function savePrivilege(data) {
  return request({
    url: '/rest/privilege/save',
    method: 'post',
    data
  })
}

export function deletePrivilege(id) {
  return request({
    url: '/rest/privilege/delete',
    method: 'post',
    data: id
  })
}
