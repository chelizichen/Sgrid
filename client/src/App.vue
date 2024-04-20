<script setup lang="ts">
import { RouterView } from "vue-router";
</script>
<template>
  <el-container>
    <el-header>
      <div style="position: relative">
        <div style="left: 0%" class="title title-pos" @click="$router.push('/server')">
          Sgrid
        </div>
        <div style="left: 50%" class="title title-pos" @click="$router.push('/server')">
          <el-icon><Tools /></el-icon>
          DevOps
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
  width: 200px;
  justify-content: center;
  cursor: pointer;
}
.title-pos {
  position: absolute;
  top: 10px;
}
</style>
