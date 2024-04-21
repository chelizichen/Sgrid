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
</script>

<template>
  <div>
    <el-menu
      class="el-menu-vertical-demo"
      active-text-color="rgb(207, 15, 124)"
      style="border: none"
      :default-openeds="['-1']"
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
          <el-icon>
            <TrendCharts />
          </el-icon>
          <template #title>{{ "DashBoard" }}</template>
        </el-menu-item>
      </el-sub-menu>
      <template v-for="(parent, index) in props.serverList" :key="parent.id">
        <el-sub-menu :index="parent.id + '-' + index">
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
              <el-icon>
                <TrendCharts />
              </el-icon>
            </template>
            <template v-else>
              <el-icon>
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

<style scoped></style>
