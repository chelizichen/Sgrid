<script setup lang="ts">
import { RouterView } from "vue-router";
</script>
<template>
  <el-container>
    <el-header>
      <div style="position: relative; height: 100%">
        <div style="left: 3%" class="title title-pos" @click="$router.push('/server')">
          <img src="/icon.png" style="width: 50px;height:50px"/>
          Sgrid
        </div>

        <div
          style="left: 50%; bottom: 25%; font-size: 24px; color: rgb(207, 15, 124)"
          class="title title-pos"
          @click="$router.push('/devops')"
        >
          devops
        </div>
        <div style="right: 0%" class="title title-pos"></div>
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
