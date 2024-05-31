<script setup lang="ts">
import { RouterView } from "vue-router";
</script>
<template>
  <el-container>
    <el-header>
      <div style="position: relative; height: 100%">
        <div style="left: 1%" class="title title-pos" @click="$router.push('/server')">
          <img src="/icon.png" style="width: 50px; height: 50px" />
          Sgrid
        </div>
        <div class="title title-pos devops" @click="$router.push('/devops')">
          <el-icon><Setting /></el-icon>
        </div>
      </div>
    </el-header>
    <RouterView />
  </el-container>
</template>

<style scoped>
.el-container {
  min-height: 100vh;
}
.aside >>> .el-card__body {
  padding: 10px 0;
  height: 100vh;
}
.title {
  color: rgb(207, 15, 124);
  text-align: center;
  display: flex;
  align-items: center;
  font-size: 30px;
  justify-content: center;
  cursor: pointer;
}
.title-pos {
  position: absolute;
  top: 10px;
}
.devops {
  right: 0;
  bottom: 0;
  font-size: 24px;
}
</style>
