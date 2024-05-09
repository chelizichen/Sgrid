import { createRouter, createWebHashHistory } from 'vue-router'
import ServerView from '@/views/server.vue'
import Login from '@/views/login.vue'
import LogPage from '@/views/logger.vue'
import Devops from '@/views/devops.vue'
import { localGet, constants } from '@/constant'
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
      path: '/logpage',
      name: 'logpage',
      component: LogPage
    },
    {
      path: '/devops',
      name: 'devops',
      component: Devops
    }
  ]
})

const whileList = ['/login']

router.beforeEach(async (to, from, next) => {
  const tkn = localGet(constants.TOKEN)
  if (whileList.includes(to.path)) {
    next()
  } else if (tkn) {
    next()
  } else {
    next({ path: '/login' })
  }
})

export default router
