import { createRouter, createWebHashHistory } from 'vue-router'
import type { RouteRecordRaw } from 'vue-router'
import ServerView from '@/views/server.vue'
import Login from '@/views/login.vue'
import LogPage from '@/views/logger.vue'
import Dashboard from '@/views/dashboard.vue'
import Devops from '@/views/devops/aside.vue'
import { localGet, constants } from '@/constant'
import { useUserStore } from '@/stores/counter'
import api from '@/api/server'
import { ElNotification } from 'element-plus'
import NProgress from 'nprogress'; // Progress 进度条
import 'nprogress/nprogress.css';// Progress 进度条样式

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
      path:"/dashboard",
      name:"dashboard",
      component:Dashboard
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
      // children: adminMenu
    },
  ]
})

const whileList = ['/login']

router.beforeEach(async (to, from, next) => {
  NProgress.start();
  const tkn = localGet(constants.TOKEN)
  const userStore = useUserStore()
  // 没登陆，但是在白名单内，直接跳转
  if (whileList.includes(to.path)) {
    next()
  }
  // 登陆了，但是store里面没有用户信息，需要重新获取用户信息
  else if (tkn &&  (!userStore.userInfo.userName && !userStore.userInfo.password && !userStore.userInfo.token)) {
    const data = await api.LoginByCache({
      name:'',
      password:''
    })
    if(!data.data){
      ElNotification.error("登陆过期，请重新登陆")
      return next({path:"/login"})
    }
    userStore.setUserInfo(data.data)
    const menus = await api.getUserMenusByUserId(data.data.id)
    await userStore.setMenu(menus.data)
    next({
      replace:true,
      ...to
    })
  } 
  // 登陆了，Store里也有对应的数据
  else if(tkn){
      next()
  }
  // 走登陆
  else{
    next({ path: '/login' })
  }
})

router.afterEach(() => {
  NProgress.done(); // 结束Progress
});

export default router
