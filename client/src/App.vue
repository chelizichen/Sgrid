<script setup lang="ts">
import { RouterView } from "vue-router";
import { useServersStore } from "./stores/counter";
import { onMounted } from "vue";
const serverStore = useServersStore();

onMounted(() => {
  serverStore.refreshServants();
});
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
            <el-sub-menu index="1">
              <template #title>
                <el-icon><SwitchFilled /></el-icon> Workspace</template
              >
              <el-menu-item index="1-0" @click="$router.push('/server')">
                <el-icon><Connection /></el-icon>
                Control</el-menu-item
              >
              <el-menu-item index="1-1" @click="$router.push('/dashboard')">
                <el-icon><DataAnalysis /></el-icon>
                Monitor</el-menu-item
              >
              <el-menu-item index="1-2" @click="$router.push('/devops')">
                <el-icon><Operation /></el-icon>
                Devops</el-menu-item
              >
              <el-menu-item index="1-3" @click="$router.push('/mainLog')">
                <el-icon><ChatLineSquare /></el-icon>
                MainLog</el-menu-item
              >
            </el-sub-menu>
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
