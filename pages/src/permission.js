import router from './router'
import store from './store'
import { Message } from 'element-ui'
import NProgress from 'nprogress' // progress bar
import 'nprogress/nprogress.css' // progress bar style
import { getToken } from '@/utils/auth' // 从localStorage中获取token

NProgress.configure({ showSpinner: false }) // Nprogress页面跳转出现在浏览器顶部的进度条

const whiteList = ['/login', '/auth-redirect'] // 没有重定向白名单

router.beforeEach(async(to, from, next) => { // 路由前置守卫
  // 开启进度条
  NProgress.start()

  // 检查用户是否登录
  const hasToken = getToken()

  if (hasToken) {
    if (to.path === '/login') {
      // 如果用户已经登录,则重定向到主页面
      next({ path: '/' })
      NProgress.done() // hack: https://github.com/PanJiaChen/vue-element-admin/pull/2939
    } else {
      // determine whether the user has obtained his permission roles through getInfo
      const hasRoles = store.getters.roles && store.getters.roles.length > 0
      if (hasRoles) {
        next()
      } else {
        try {
          // 获取字典信息
          await store.dispatch('settings/setDictMap')
          // 获取用户权限信息
          const { privileges } = await store.dispatch('user/getUserInfo')

          // 基于权限编码生成动态路由
          const accessRoutes = await store.dispatch('permission/generateRoutes', privileges)

          // 动态添加路由
          router.addRoutes(accessRoutes)

          // hack method to ensure that addRoutes is complete
          // set the replace: true, so the navigation will not leave a history record
          next({ ...to, replace: true })
        } catch (error) {
          // remove token and go to login page to re-login
          await store.dispatch('user/resetToken')
          Message.error(error || 'Has Error')
          next(`/login?redirect=${to.path}`)
          NProgress.done()
        }
      }
    }
  } else {
    /* 没有token*/
    if (whiteList.indexOf(to.path) !== -1) {
      // 在登录白名单中,则放行
      next()
    } else {
      // 没有页面访问权限,并且重定向到登录页面
      next(`/login?redirect=${to.path}`)
      NProgress.done()
    }
  }
})

router.afterEach(() => {
  // finish progress bar
  NProgress.done()
})
