import { getUserInfo, login } from '@/api/user'
import { getToken, setToken, removeToken } from '@/utils/auth'
import router, { resetRouter } from '@/router'
import { Message } from 'element-ui'

const state = {
  token: getToken(),
  org: {},
  privileges: [],
  user: {},
  roles: [],
  buttonPrivileges: []
}

const mutations = {
  SET_TOKEN: (state, token) => {
    state.token = token
  },
  SET_PRIVILEGES: (state, privileges) => {
    state.privileges = privileges
  },
  SET_BUTTON_PRIVILEGES: (state, buttonPrivileges) => {
    state.buttonPrivileges = buttonPrivileges
  },
  SET_ORG: (state, org) => {
    state.org = org
  },
  SET_USER: (state, user) => {
    state.user = user
  },
  SET_ROLES: (state, roles) => {
    state.roles = roles
  }
}

// actions中写调用接口和commit数据到mutations中逻辑,接收一个store实例具有相同方法和属性的context对象
const actions = {
  // 用户登录
  login({ commit }, userInfo) {
    const { username, password } = userInfo
    return new Promise((resolve, reject) => {
      login({ username: username.trim(), password: password }).then(response => {
        const { code, message, result } = response
        if (code === 200) {
          commit('SET_TOKEN', result.token) // 提交到vuex mutations中
          setToken(result.token) // 设置本地localStorage
          resolve(result)
          Message.success(message)
        } else {
          reject()
          Message.error(message)
        }
      }).catch(error => {
        reject(error)
      })
    })
  },

  // 获取用户信息
  getUserInfo({ commit, state }) {
    return new Promise((resolve, reject) => {
      getUserInfo(state.token).then(response => {
        const { result } = response
        if (!result) {
          reject('用户认证失败,请重新登录')
        }
        const { org, privileges, roles, user } = result

        // roles must be a non-empty array
        if (!roles || roles.length <= 0) {
          reject('getInfo: 该用户角色为空!')
        }
        const privilegesArr = []
        const privilegesButtonArr = []
        privileges.forEach(function(item) {
          if (item.type === 'menu') {
            privilegesArr.push(item.code)
          } else {
            privilegesButtonArr.push(item.code)
          }
        })
        if (privilegesArr.length === 0) {
          privilegesArr.push('p_zhuye')
        }
        // console.log('privilegesArr:', privilegesArr)
        // console.log('privilegesButtonArr:', privilegesButtonArr)
        commit('SET_ROLES', roles)
        commit('SET_USER', user)
        commit('SET_ORG', org)
        commit('SET_PRIVILEGES', privilegesArr)
        commit('SET_BUTTON_PRIVILEGES', privilegesButtonArr)
        resolve(result)
      }).catch(error => {
        reject(error)
      })
    })
  },

  // 用户登出
  logout({ commit, state, dispatch }) {
    return new Promise((resolve, reject) => {
      commit('SET_TOKEN', '')
      commit('SET_ROLES', [])
      removeToken() // 从localStorage中移除token
      resetRouter() // 重置路由

      // reset visited views and cached views
      // to fixed https://github.com/PanJiaChen/vue-element-admin/issues/2485
      dispatch('tagsView/delAllViews', null, { root: true })

      resolve()
    })
  },

  // 移除token
  resetToken({ commit }) {
    return new Promise(resolve => {
      commit('SET_TOKEN', '')
      commit('SET_ROLES', [])
      removeToken()
      resolve()
    })
  },

  // 动态修改权限
  async changeRoles({ commit, dispatch }, role) {
    const token = role + '-token'

    commit('SET_TOKEN', token)
    setToken(token)

    const { roles } = await dispatch('getInfo')

    resetRouter()

    // 基于权限生成可用的路由
    const accessRoutes = await dispatch('permission/generateRoutes', roles, { root: true })
    // 动态添加可用的路由
    router.addRoutes(accessRoutes)

    // 重置访客视图和缓存视图
    dispatch('tagsView/delAllViews', null, { root: true })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
