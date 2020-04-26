import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/static',
      redirect: { name: 'home' }
    },
    {
      path: '/static/index',
      component: resolve => require(['../views/main/index.vue'], resolve),
      children: [
        {
          path: '/static/home',
          name: 'home',
          component: resolve => require(['../views/home/index.vue'], resolve)
        },
        {
          path: '/static/node',
          name: 'node',
          component: resolve => require(['../views/node/index.vue'], resolve)
        },
        {
          path: '/static/task',
          name: 'task',
          component: resolve => require(['../views/task/index.vue'], resolve)
        },
        {
          path: '/static/task/detail',
          name: 'detail',
          component: resolve => require(['../views/task/detail.vue'], resolve)
        },
        {
          path: '/static/task/logdetail',
          name: 'logdetail',
          component: resolve => require(['../views/task/logdetail.vue'], resolve)
        }
      ]
    },
    {
      path: '/static/login',
      name: 'login',
      component: resolve => require(['../views/login.vue'], resolve)
    }
  ]
})
