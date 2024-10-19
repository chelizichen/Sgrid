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
      <div class="flex flex-row justify-between">
        <div @click="$router.push('/server')" class="flex items-center">
          <img src="/icon.png" class="w-8 h-8 mr-2" />
          <div
            class="text-xl font-serif tracking-wider text-indigo-600 font-bold transition-all hover:text-indigo-800"
          >
            Sgrid
          </div>
        </div>
        <div class="right">
          <el-menu class="border-none" mode="horizontal" :ellipsis="false">
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
