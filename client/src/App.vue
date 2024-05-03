<script setup lang="ts">
import { RouterView } from "vue-router";
</script>
<template>
  <el-container>
    <el-header>
      <div style="position: relative; height: 100%">
        <div style="left: 3%" class="title title-pos" @click="$router.push('/server')">
          <el-icon>
            <svg
              xmlns="http://www.w3.org/2000/svg"
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="rgb(207, 15, 124)"
              stroke-width="2"
              stroke-linecap="round"
              stroke-linejoin="round"
              class="lucide lucide-slack"
            >
              <rect width="3" height="8" x="13" y="2" rx="1.5" />
              <path d="M19 8.5V10h1.5A1.5 1.5 0 1 0 19 8.5" />
              <rect width="3" height="8" x="8" y="14" rx="1.5" />
              <path d="M5 15.5V14H3.5A1.5 1.5 0 1 0 5 15.5" />
              <rect width="8" height="3" x="14" y="13" rx="1.5" />
              <path d="M15.5 19H14v1.5a1.5 1.5 0 1 0 1.5-1.5" />
              <rect width="8" height="3" x="2" y="8" rx="1.5" />
              <path d="M8.5 5H10V3.5A1.5 1.5 0 1 0 8.5 5" />
            </svg>
          </el-icon>
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
