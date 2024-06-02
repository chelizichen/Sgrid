<script lang="ts">
export default {
  name: "aside-component",
};
</script>
<script lang="ts" setup>
import type { Item } from "@/dto/dto";
import { useRouter } from "vue-router";
const props = defineProps<{
  serverList: any[];
}>();
const router = useRouter();
function toDashboard() {
  router.push({
    path: "/server/dashboard",
  });
}
const emits = defineEmits(["handleOpen"]);
function handleOpen(item: Partial<Item>) {
  emits("handleOpen", item);
}
</script>

<template>
  <div>
    <el-menu
      class="el-menu-vertical-demo"
      active-text-color="rgb(207, 15, 124)"
      style="border: none"
    >
      <template v-for="(parent, index) in props.serverList" :key="parent.id">
        <el-sub-menu :index="parent.id + '-' + index">
          <template #title>
            <el-icon><Grid /></el-icon>
            <span>{{ parent.tagEnglishName }}</span>
          </template>
          <el-menu-item
            v-for="(item, index) in parent.servants"
            class="app-text-center"
            :index="item"
            :key="index"
            @click="handleOpen(item)"
          >
            <el-icon>
              <TrendCharts />
            </el-icon>
            <template #title>{{ item.serverName }}</template>
          </el-menu-item>
        </el-sub-menu>
      </template>
    </el-menu>
  </div>
</template>

<style scoped></style>
