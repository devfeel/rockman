// The Vue build version to load with the `import` command
// (runtime-only or standalone) has been set in webpack.base.conf with an alias.
import Vue from 'vue'
import App from './App'
import 'element-ui/lib/theme-chalk/index.css'
import router from './router'
import ViewUI from 'view-design';
import ElementUI from 'element-ui'
// 引入字体图标库
import 'view-design/dist/styles/iview.css';
import './assets/iconfont/iconfont.css'
import GLOBAL from './common/global.js'
import store from './store/store.js'
import * as Utils from './common/utils.js'

Vue.use(ElementUI);
Vue.use(ViewUI)

// 挂载到Vue实例上面
Vue.prototype.GLOBAL = GLOBAL
Vue.prototype.Utils = Utils

if (process.env.NODE_ENV === 'development') {
  require('./api/mock');
}

// 页面刷新时，重新赋值
if (window.sessionStorage.getItem('Token')) {
  // let data = JSON.parse(window.sessionStorage.getItem('Token'));
  // let token = window.sessionStorage.getItem('Token');
  // store.commit('SET_TOKEN', token)
  // store.dispatch('ChangeTheme', data.theme)
}

router.beforeEach(({meta, path}, from, next) => {
  // var {auth = true} = meta
  // true用户已登录， false用户未登录
  // var isLogin = Boolean(store.state.user.token)
  // if (auth && !isLogin && path !== '/static/login') {
  //    return next({ path: '/static/login' })
  // }
  // router.push({ path: '/static/home' });
  next()
})

Vue.config.productionTip = false

let vm = new Vue({
  el: '#app',
  store,
  router,
  components: {App},
  template: '<App/>'
})

Vue.use({
  vm
})
