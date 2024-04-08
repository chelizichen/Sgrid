<script lang="ts">
export default {
  name: "aside-component",
};
</script>
<script lang="ts" setup>
import type { Item } from "@/dto/dto";
const props = defineProps<{
  serverList: any[];
}>();
const emits = defineEmits(["handleOpen"]);
function handleOpen(item: Partial<Item>) {
  emits("handleOpen", item);
}

function toGit() {
  window.open("https://github.com/chelizichen/Simp");
}
</script>

<template>
  <div>
    <div class="app-bigger-size title" @click="toGit()">
      <el-icon style="color: rgb(207, 90, 124); font-size: 36px"><Help /></el-icon>
      Sgrid
    </div>
    <el-menu
      class="el-menu-vertical-demo"
      active-text-color="rgb(207, 15, 124)"
      style="border: none"
      :default-openeds="['2']"
    >
      <el-sub-menu index="-1">
        <template #title>
          <el-icon><Setting /></el-icon>
          <span>Sgrids</span>
        </template>
        <el-menu-item
          class="app-text-center"
          index="-1"
          key="-1"
          @click="
            handleOpen({
              id: 0,
              serverName: 'dashboard',
            })
          "
        >
          <el-icon class="app-not-show">
            <TrendCharts />
          </el-icon>
          <template #title>{{ "DashBoard" }}</template>
        </el-menu-item>
      </el-sub-menu>
      <template v-for="parent in props.serverList" :key="parent.id">
        <el-sub-menu :index="parent.id">
          <template #title>
            <template v-if="parent.tagEnglishName === 'Controller'">
              <el-icon><Setting /></el-icon>
            </template>
            <template v-else>
              <el-icon><Grid /></el-icon>
            </template>
            <span>{{ parent.tagEnglishName }}</span>
          </template>
          <el-menu-item
            v-for="(item, index) in parent.servants"
            class="app-text-center"
            :index="item"
            :key="index"
            @click="handleOpen(item)"
          >
            <template v-if="parent.tagEnglishName === 'Controller'">
              <el-icon class="app-not-show">
                <TrendCharts />
              </el-icon>
            </template>
            <template v-else>
              <el-icon class="app-not-show">
                <Menu />
              </el-icon>
            </template>
            <template #title>{{ item.serverName }}</template>
          </el-menu-item>
        </el-sub-menu>
      </template>
    </el-menu>
  </div>
</template>

<style>
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
</style>
