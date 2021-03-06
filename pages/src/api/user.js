import request from '@/utils/request'

export function login(data) {
  return request({
    url: '/login',
    method: 'post',
    data
  })
}

export function getUserInfo(token) {
  return request({
    url: '/rest/auth/getUserInfo',
    method: 'post'
  })
}

// export function logout() {
//   return request({
//     url: '/logout',
//     method: 'post'
//   })
// }
