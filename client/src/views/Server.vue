<script lang="ts">
export default {
  name: "server-component",
};
</script>

<template>
  <div>
    <el-container>
      <el-aside width="200px" class="aside">
        <aside-component
          :server-list="state.serverList"
          @handle-open="handleOpen"
        ></aside-component>
      </el-aside>
      <el-main>
        <el-card>
          <gridsComponent
            :grids-list="state.gridsList"
            :server-name="state.serverName"
            :servant-id="state.servantId"
            @check-status="handleOpen(currentItem)"
          ></gridsComponent>
        </el-card>
      </el-main>
      <router-view></router-view>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import asideComponent from "@/components/aside.vue";
import gridsComponent from "@/components/grids.vue";
import { onMounted, reactive, ref } from "vue";
import API from "../api/server";
import type { Item } from "@/dto/dto";

const state = reactive({
  serverName: "",
  serverList: [],
  gridsList: [],
  servantId: 0,
});
const currentItem = ref();
async function handleOpen(item: Item) {
  currentItem.value = item;
  console.log("item", item);
  const grids = await API.queryGrid({ id: item.id });
  state.gridsList = grids.data;
  state.servantId = item.id;
  console.log("state.grids", state.gridsList);
  state.serverName = item.serverName;
}

async function fetchServerList() {
  const resp = await API.getServerList(0);
  state.serverList = (resp.data || []).sort((a, b) => a.id - b.id);
}

onMounted(() => {
  fetchServerList();
});
</script>

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
