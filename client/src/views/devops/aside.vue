<script lang="ts">
// devops component
// 1.选择服务组 ｜ 创建
// 2.创建服务
// 3.选择节点
// 4.添加至服务网格
// 5.选择端口
export default {
  name: "devops-component",
};
</script>
<template>
  <div>
    <el-container>
      <el-aside style="width: 200px; height: 100vh">
        <el-menu class="el-menu-vertical-demo" active-text-color="rgb(207, 15, 124)">
          <el-sub-menu
            v-for="(item, index) in adminMenu"
            :index="String(index)"
            :key="index"
          >
            <template #title>
              <el-icon><grid /></el-icon>
              <span>{{ item.title }}</span>
            </template>
            <el-menu-item
              v-for="(s_item, s_index) in item.children"
              :key="s_index"
              :index="index + '_' + s_index"
              @click="push(item, s_item)"
              >{{ s_item.title }}</el-menu-item
            >
          </el-sub-menu>
        </el-menu>
      </el-aside>
      <el-main>
        <router-view></router-view>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { adminMenu } from "@/router";
import { useRouter } from "vue-router";
const router = useRouter();
const base_path = "/devops/";
function push(item, s_item) {
  const path = base_path + item.path + "/" + s_item.path;
  router.push(path);
  console.log(item, s_item);
}
</script>

<style scoped></style>
