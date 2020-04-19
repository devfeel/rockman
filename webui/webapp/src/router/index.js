import Vue from 'vue'
import Router from 'vue-router'

Vue.use(Router)

export default new Router({
  mode: 'history',
  routes: [
    {
      path: '/static',
      redirect: { name: 'login' }
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
          path: '/static/settings',
          component: resolve => require(['../views/settings/home/index.vue'], resolve),
          children: [
            {
              path: '',
              name: 'settings',
              component: resolve => require(['../views/settings/nodes/index.vue'], resolve)
            },
            {
              path: 'nodes',
              name: 'nodes',
              component: resolve => require(['../views/settings/nodes/index.vue'], resolve)
            },
            {
              path: 'users',
              name: 'users',
              component: resolve => require(['../views/settings/users/index.vue'], resolve)
            }
          ]
        },
        {
          path: '/static/runtimes',
          component: resolve => require(['../views/runtimes/home/index.vue'], resolve),
          children: [
            {
              path: '',
              name: 'runtimes',
              component: resolve => require(['../views/runtimes/tasks/index.vue'], resolve)
            },
            {
              path: 'tasks',
              name: 'tasks',
              component: resolve => require(['../views/runtimes/tasks/index.vue'], resolve)
            },
            {
              path: 'taskdetail',
              name: 'taskdetail',
              component: resolve => require(['../views/runtimes/tasks/detail.vue'], resolve)
            },
            {
              path: 'logs',
              name: 'logs',
              component: resolve => require(['../views/runtimes/logs/index.vue'], resolve)
            }
          ]
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
