import request from '@/utils/request'

export function getRoleList(param) {
  return request({
    url: 'rest/role/list',
    method: 'post',
    data: param
  })
}

export function deleteRole(param) {
  return request({
    url: 'rest/role/delete',
    method: 'post',
    data: param
  })
}

export function saveRole(param) {
  return request({
    url: 'rest/role/save',
    method: 'post',
    data: param
  })
}

export function configRolePrivileges(param) {
  return request({
    url: 'rest/role/configPrivileges',
    method: 'post',
    data: param
  })
}

export function privilegesRole(param) {
  return request({
    url: 'rest/role/privileges',
    method: 'post',
    data: param
  })
}
