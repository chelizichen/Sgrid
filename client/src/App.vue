<script setup lang="ts">
import { RouterView } from 'vue-router'
import { useServersStore } from './stores/counter'
import { onMounted } from 'vue'
import SvgDashboard from './svg/dashboard.vue'
import SvgMenu from './svg/menu.vue'
const serverStore = useServersStore()

onMounted(() => {
  serverStore.refreshServants()
})
</script>
<template>
  <el-container>
    <el-header>
      <div class="head">
        <div @click="$router.push('/server')" class="left">
          <img src="/icon_v2.png" class="logo" />
        </div>
        <div class="right">
          <el-menu class="el-menu-demo" mode="horizontal" :ellipsis="false">
            <el-menu-item index="0" @click="$router.push('/server')">Control</el-menu-item>
            <el-menu-item index="1" @click="$router.push('/dashboard')">Monitor</el-menu-item>
            <el-menu-item index="2" @click="$router.push('/devops')">Devops</el-menu-item>
          </el-menu>
        </div>
      </div>
    </el-header>
    <RouterView />
  </el-container>
</template>

<style scoped lang="less">
.head {
  position: relative;
  height: 100%;
  border-bottom: 1px gainsboro;
  border-bottom-style: dotted;
  display: flex;
  justify-content: space-between;
  .left {
    display: flex;
    cursor: pointer;
    .logo {
      width: 120px;
      height: 120px;
      margin-top: -30px;
    }
  }
  .right {
    display: flex;
    align-items: center;
    .icon {
      width: 32px;
      height: 32px;
      margin-right: 10px;
      cursor: pointer;
    }
  }
}
</style>
