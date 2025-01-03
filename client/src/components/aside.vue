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
const emits = defineEmits(["handleOpen"]);
function handleOpen(item: Partial<Item>) {
  emits("handleOpen", item);
}
function toHome() {
  emits("handleOpen", {
    id: -1,
    serverName: "Home",
    language: "zh",
    tagEnglishName: "Home",
  });

}
</script>

<template>
  <div class="app">
    <el-menu active-text-color="var(--sgrid-primay-hover-color)" class="border-none">
      <el-menu-item index="1" class="h-8" @click="toHome">
        <template #title>
          <span>Home</span>
        </template>
      </el-menu-item>
      <template v-if="props.serverList && props.serverList.length">
        <template v-for="(parent, index) in props.serverList" :key="parent.id">
          <el-sub-menu :index="parent.id + '-' + index">
            <template #title>
              <span>{{ parent.tagEnglishName }}</span>
            </template>
            <el-menu-item v-for="(item, index) in parent.servants" :index="item" :key="index" @click="handleOpen(item)"
              class="h-8">
              <template #title>{{ item.serverName }}</template>
            </el-menu-item>
          </el-sub-menu>
        </template>
      </template>
      <template v-else>
        <el-empty />
      </template>
    </el-menu>
  </div>
</template>

<style lang="css" scoped>
.app >>> .el-sub-menu__title {
  height: 44px;
}
</style>
