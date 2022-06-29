import { asyncRoutes, constantRoutes } from '@/router'

/**
 * Use meta.role to determine if the current user has permission
 * @param roles
 * @param route
 */
function hasPermission(privilege, route) {
  if (route.meta && route.meta.privilege === privilege.code) {
    return true
    // return roles.some(role => route.meta.roles.includes(role))
  } else {
    return false
  }
}

/**
 * 通过权限编码动态过滤路由
 * @param routes asyncRoutes
 * @param roles
 */
export function filterAsyncRoutes(routes, privileges) {
  const res = []
  routes.forEach(route => {
    const tmp = { ...route }
    privileges.forEach(privilege => {
      if (hasPermission(privilege, tmp)) {
        if (tmp.children) {
          tmp.children = filterAsyncRoutes(tmp.children, privileges)
        }
        res.push(tmp)
      }
    })
  })
  return res
}

const state = {
  routes: [],
  addRoutes: []
}

const mutations = {
  SET_ROUTES: (state, routes) => {
    state.addRoutes = routes
    state.routes = constantRoutes.concat(routes)
  }
}

const actions = {
  // 基于权限编码生成动态路由
  generateRoutes({ commit }, privileges) {
    return new Promise(resolve => {
      const accessedRoutes = filterAsyncRoutes(asyncRoutes, privileges)
      commit('SET_ROUTES', accessedRoutes)
      resolve(accessedRoutes)
    })
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
