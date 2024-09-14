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
          <div @click="$router.push('/dashboard')">
            <SvgDashboard class="icon" />
          </div>
          <div @click="$router.push('/devops')">
            <SvgMenu class="icon" />
          </div>
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
