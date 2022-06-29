import request from '@/utils/request'

export function getOrgList(param) {
  return request({
    url: 'rest/org/list',
    method: 'post',
    data: param
  })
}

export function deleteOrg(param) {
  return request({
    url: 'rest/org/delete',
    method: 'post',
    data: param
  })
}

export function saveOrg(param) {
  return request({
    url: 'rest/org/save',
    method: 'post',
    data: param
  })
}

export function configPrivileges(param) {
  return request({
    url: 'rest/org/configPrivileges',
    method: 'post',
    data: param
  })
}

export function getOrgPrivileges(param) {
  return request({
    url: 'rest/org/privileges',
    method: 'post',
    data: param
  })
}

export function getRootList(param) {
  return request({
    url: 'rest/org/rootList',
    method: 'post',
    data: param
  })
}

export function getOrgRole(param) {
  return request({
    url: 'rest/org/roles',
    method: 'post',
    data: param
  })
}
