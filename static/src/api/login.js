import request from '@/utils/request'

export function login(username, password) {
  return request({
    url: '/login',
    method: 'post',
    data: {
      username,
      password
    }
  })
}

export function getInfo(token) {
  return request({
    url: 'api/v1/user/info',
    method: 'get',
    params: { token }
  })
}
export function getUsers(token) {
  return request({
    url: 'stock/users',
    method: 'get',
    params: { token }
  })
}
export function logout() {
  return request({
    url: 'api/v1/user/logout',
    method: 'post'
  })
}
