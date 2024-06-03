<script setup lang="ts">
import { RouterView } from "vue-router";
</script>
<template>
  <el-container>
    <el-header>
      <div class="head">
        <div style="left: 1%" class="title title-pos" @click="$router.push('/server')">
          <img src="@/assets/icon.png" style="width: 50px; height: 50px" />
          Sgrid
        </div>
        <div
          style="left: 88%; bottom: 1px"
          class="title title-pos"
          @click="$router.push('/dashboard')"
        >
          <img src="@/assets/dashboard.png" style="width: 38px; height: 38px" />
        </div>
        <div class="title title-pos" style="right: 5%; bottom:0"  @click="$router.push('/devops')">
          <img src="@/assets/admin.png" style="width: 100px; height: 100px" />
        </div>
        <div
          style="right: 2%; bottom: 8px"
          class="title title-pos"
          @click="$router.push('/userInfo')"
        >
          <img src="@/assets/user.png" style="width: 70px; height: 70px" />
        </div>
        
      </div>
    </el-header>
    <RouterView />
  </el-container>
</template>

<style scoped>
.head {
  position: relative;
  height: 100%;
  border-bottom: 1px gainsboro;
  border-bottom-style: dotted;
}
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
