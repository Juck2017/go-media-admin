import Vue from 'vue'
import Cookies from 'js-cookie'
import 'normalize.css/normalize.css' // a modern alternative to CSS resets
import Element from 'element-ui' // element ui组件
import './styles/element-variables.scss' // element ui样式
import '@/styles/index.scss' // 全局css
import App from './App'
import store from './store'
import router from './router'
import './icons' // 图标
import './permission' // 权限控制
import './utils/error-log' // 错误日志
import * as filters from './filters' // 全局过滤器

// 全局注册element ui组件
Vue.use(Element, {
  size: Cookies.get('size') || 'medium' // set element-ui default size
})

// 注册全局过滤器
Object.keys(filters).forEach(key => {
  Vue.filter(key, filters[key])
})

Vue.config.productionTip = false

new Vue({
  el: '#app',
  router,
  store,
  render: h => h(App)
})
