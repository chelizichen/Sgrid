import { createRouter, createWebHashHistory } from 'vue-router'
import ServerView from '@/views/Server.vue'
import Login from '@/views/Login.vue'
import LogPage from '@/views/LogPage.vue'
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
    }
  ]
})

export default router
