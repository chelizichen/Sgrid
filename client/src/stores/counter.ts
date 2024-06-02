import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { reduceMenu, reduceMenuToRouter } from '@/utils/obj'
import type { RouteRecordRaw } from 'vue-router'
import router from '@/router'
type userVo = {
  id: number
  password: string
  token: string
  userName: string
}

export const useUserStore = defineStore('user', () => {
  const userInfo = ref<userVo>({
    id: 0,
    password: '',
    userName: '',
    token: ''
  })
  const menus = ref<RouteRecordRaw[]>([])
  function setUserInfo(userInfoDto: userVo) {
    userInfo.value = userInfoDto
  }
  async function setMenu(obj: any) {
    const toRouters = reduceMenuToRouter(obj)
    menus.value = await toRouters
    menus.value.forEach((fatherRoute) => {
      router.addRoute('devops', fatherRoute)
      if (fatherRoute.children) {
        fatherRoute.children.forEach((childRoute) => {
          router.addRoute(fatherRoute.name!, childRoute)
        })
      }
    })
    console.log('getRoutes', router.getRoutes())
  }
  return { userInfo, setUserInfo, menus, setMenu }
})
