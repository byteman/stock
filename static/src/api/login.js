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
export function updateUser(uid, user) {
  return request({
    url: 'stock/user/' + uid,
    method: 'put',
    data: user
  })
}
export function deleteUser(uid) {
  return request({
    url: 'stock/user/' + uid,
    method: 'delete'
  })
}
export function logout() {
  return request({
    url: 'api/v1/user/logout',
    method: 'post'
  })
}
