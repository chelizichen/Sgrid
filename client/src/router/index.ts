import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import ServerView from '@/views/server.vue'
import Login from '@/views/login.vue'
import LogPage from '@/views/logger.vue'
import Devops from '@/views/devops/aside.vue'
import { localGet, constants } from '@/constant'

export const adminMenu: RouteRecordRaw[] = [
  {
    title: '系统管理',
    icon: 'Grid',
    path: 'system',
    name: 'system',
    children: [
      {
        title: '用户管理',
        icon: 'Grid',
        path: 'user',
        name: 'user',
        component: ()=>import('@/views/devops/system/user_admin.vue')
      },
      {
        title: '角色管理',
        icon: 'Grid',
        path: 'role',
        name: 'role',
        component: ()=>import('@/views/devops/system/role_admin.vue')
      },
      {
        title: '菜单管理',
        icon: 'Grid',
        path: 'menu',
        name: 'menu',
        component: ()=>import('@/views/devops/system/menu_admin.vue')
      }
    ]
  },
  {
    title: '服务管理',
    icon: 'Grid',
    path: 'servant',
    name: 'servant',
    children: [
      {
        title: '添加服务组',
        icon: 'Grid',
        path: 'add_group',
        name: 'add_group',
        component: ()=>import('@/views/devops/add_group.vue')
      },
      {
        title: '添加服务',
        icon: 'Grid',
        path: 'add_servant',
        name: 'add_servant',
        component: ()=>import('@/views/devops/add_servant.vue')
      },
      {
        title: '服务列表',
        icon: 'Grid',
        path: 'servant_admin',
        name: 'servant_admin',
        component: ()=>import('@/views/devops/servant_list.vue')
      }
    ]
  },
  {
    title: '节点管理',
    icon: 'Grid',
    path: 'grid',
    name: 'grid',
    children: [
      {
        title: '添加节点',
        icon: 'Grid',
        path: 'add_node',
        name: 'add_node',
        component: ()=>import('@/views/devops/add_node.vue')
      }
    ]
  },
  {
    title: '网关管理',
    icon: 'Grid',
    path: 'gateway',
    name: 'gateway',
    children: [
      {
        title: '网关配置',
        icon: 'Grid',
        path: 'admin',
        name: 'admin',
        component: ()=>import('@/views/devops/gateway_conf.vue')
      }
    ]
  },
  {
    title: '属性管理',
    icon: 'Grid',
    path: 'property',
    name: 'property',
    children: [
      {
        title: '属性设置',
        icon: 'Grid',
        path: 'set',
        name: 'set',
        component: ()=>import('@/views/devops/property_admin.vue')
      }
    ]
  }
]

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
      component: Devops,
      children: adminMenu
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
