import { createRouter, createWebHashHistory } from 'vue-router'
import ServerView from '@/views/server.vue'
import Login from '@/views/login.vue'
import LogPage from '@/views/logger.vue'
import Devops from '@/views/devops.vue'
const router = createRouter({
  history: createWebHashHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/login'
    },
    {
      path: '/server',
      name: 'server',
      component: ServerView
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    },
    {
      path:'/logpage',
      name:'logpage',
      component: LogPage
    },
    {
      path:'/devops',
      name:'devops',
      component: Devops
    },
  ]
})

export default router
