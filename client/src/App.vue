<script setup lang="ts">
import { RouterView } from "vue-router";
</script>
<template>
  <el-container>
    <el-header>
      <div style="position: relative; height: 100%">
        <div
          style="left: 3%; text-decoration: underline"
          class="title title-pos"
          @click="$router.push('/server')"
        >
          <el-icon><Sunny /></el-icon>
          Sgrid
        </div>

        <div
          style="
            left: 50%;
            bottom: 25%;
            font-size: 24px;
            color: rgb(207, 15, 124);
            text-decoration: underline;
          "
          class="title title-pos"
          @click="$router.push('/devops')"
        >
          <el-icon><Lightning /></el-icon> devops
        </div>
        <div style="right: 0%" class="title title-pos">
          <el-icon><Tools /></el-icon>
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
</style>
