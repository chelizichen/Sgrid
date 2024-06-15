import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import { reduceMenu, reduceMenuToRouter } from '@/utils/obj'
import type { RouteRecordRaw } from 'vue-router'
import router from '@/router'
import _, {uniqWith,isEqual} from 'lodash'
import api from '@/api/server'
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
    if(obj == null || obj == undefined){
      obj = []
    }
    obj = uniqWith(obj, isEqual)
    // obj = Array.from(new Set())
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

export const useServersStore = defineStore('servers',()=>{
  const user = useUserStore()
  const servers = ref([])
  function refreshServants(){
    api.getServants(user.userInfo.id).then(res=>{
      servers.value = res.data
    })
  }
  function getServerNameById(id:number){
    return _.get(servers.value.find(v=>v.id === id),'serverName','--')
  }
  return {
    servers,
    refreshServants,
    getServerNameById
  }
})