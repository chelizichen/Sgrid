<script lang="ts">
export default {
  name: "devops-component",
};
</script>
<template>
  <div>
    <el-container>
      <el-aside class="app h-screen w-48">
        <el-menu class="border-none" active-text-color="var(--sgrid-primay-hover-color)">
          <el-sub-menu v-for="(item, index) in menus" :index="String(index)" :key="index">
            <template #title>
              <span>{{ item.title }}</span>
            </template>
            <el-menu-item
              v-for="(s_item, s_index) in item.children"
              :key="s_index"
              :index="index + '_' + s_index"
              @click="push(item, s_item)"
              class="h-8"
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
import { useUserStore } from "@/stores/counter";
import { computed } from "vue";
import { useRouter } from "vue-router";
const router = useRouter();
const userStore = useUserStore();
const menus = computed(() => {
  console.log("userStore.menus", userStore.menus);
  return userStore.menus;
});
const base_path = "/devops/";
function push(item: any, s_item: any) {
  const path = base_path + item.path + "/" + s_item.path;
  router.push(path);
  console.log(item, s_item);
}
</script>

<style scoped>
.app >>> .el-sub-menu__title {
  height: 44px;
}
</style>
