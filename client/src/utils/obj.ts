import type { MenuVo, Tree } from '@/dto/dto'
import { loadView } from '@/permission'
import _ from 'lodash'
import type { RouteRecordRaw } from 'vue-router'
export function isEmptyObj(target: any): boolean {
  if(target == null || target == undefined){
    return true;
  }
  if (_.isObject(target)) {
    if (Object.keys(target).length == 0) {
      return true
    }
    return false
  }
  return true
}

export function reduceMenu(list: Array<MenuVo>): Tree[] {
  const menus: Tree[] = []
  // 拿到根节点
  list.forEach((e) => {
    if (e.parentId == 0 || !e.parentId) {
      menus.push({
        label: e.title,
        id: e.id
      })
    }
  })
  // 最多支持双层
  list.forEach((e) => {
    if (e.parentId || e.parentId != 0) {
      const item = menus.find((v) => v.id == e.parentId)
      if (!item) {
        return
      }
      if (!item.children) {
        item.children = []
      }
      item.children.push({
        label: e.title,
        id: e.id
      })
    }
  })
  return menus
}
              
export async function reduceMenuToRouter(list: Array<MenuVo>) {
    if(!list){
        list = []
    }
    const menus: Array<RouteRecordRaw> = []
  // 拿到根节点
  list.forEach(async (e) => {
    if (e.parentId == 0 || !e.parentId) {
        let component = null
        if(e.component){
            component = await loadView(`./views/${e.component}.vue`)
        }
      menus.push({
        title: e.title,
        id: e.id,
        path: e.path,
        name: e.name,
        icon: 'Grid',
        component:component,
      })
    }
  })
  // 最多支持双层
  list.forEach(async (e) => {
    if (e.parentId || e.parentId != 0) {
      const item = menus.find((v) => v.id == e.parentId)
      if (!item) {
        return
      }
      if (!item.children) {
        item.children = []
      }
      let component = null
      if(e.component){
          component = await loadView(`./views/${e.component}.vue`)
      }
      item.children.push({
        title: e.title,
        id: e.id,
        path: e.path,
        name: e.name,
        icon: 'Grid',
        component
      })
    }
  })
  return menus
}
