import request from '@/utils/request'

export function deleteUser(query) {
  return request({
    url: '/rest/user/delete',
    method: 'post',
    data: query
  })
}

export function saveUser(query) {
  return request({
    url: '/rest/user/save',
    method: 'post',
    data: query
  })
}

// 保存用户
// export function saveUser(query) {
//   return request({
//     url: '/rest/sysUserDiy/save',
//     method: 'post',
//     data: query
//   })
// }

export function getUserList(query) {
  return request({
    url: '/rest/user/list',
    method: 'post',
    data: query
  })
}

export function privilegesList(query) {
  return request({
    url: '/rest/user/privileges',
    method: 'post',
    data: query
  })
}

export function configUserPrivileges(query) {
  return request({
    url: '/rest/user/configPrivileges',
    method: 'post',
    data: query
  })
}

export function enabledUser(query) {
  return request({
    url: '/rest/user/enabled',
    method: 'post',
    data: query
  })
}

export function resetPassword(query) {
  return request({
    url: '/rest/user/resetPassword',
    method: 'post',
    data: query
  })
}

export function getRoleList(query) {
  return request({
    url: '/rest/user/roles',
    method: 'post',
    data: query
  })
}

export function configRoles(query) {
  return request({
    url: '/rest/user/configRoles',
    method: 'post',
    data: query
  })
}

export function updatePassword(query) {
  return request({
    url: '/rest/user/updatePassword',
    method: 'post',
    data: query
  })
}
// 根绝orgId查询人员
export function getOrgUser(query) {
  return request({
    url: '/rest/user/orgUsers',
    method: 'post',
    data: query
  })
}
// 下拉框组织结构选择人员
export function userTreeData(query) {
  return request({
    url: '/rest/user/userTreeData',
    method: 'post',
    data: query
  })
}

