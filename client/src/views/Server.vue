<script lang="ts">
export default {
  name: "server-component",
};
</script>

<template>
  <div>
    <el-container>
      <el-aside width="200px" class="py-2.5 h-full">
        <aside-component
          :server-list="state.serverList"
          @handle-open="handleOpen"
        ></aside-component>
      </el-aside>
      <el-main>
        <div>
          <gridsComponent
            :grids-list="state.gridsList"
            :server-name="state.serverName"
            :servant-id="state.servantId"
            :servantLanguage="state.servantLanguage"
            @check-status="handleOpen(currentItem)"
            :server-version="state.serverVersion"
          ></gridsComponent>
        </div>
      </el-main>
      <router-view></router-view>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import asideComponent from "@/components/aside.vue";
import gridsComponent from "@/components/grids.vue";
import { onMounted, reactive, ref } from "vue";
import api from "../api/server";
import type { Item } from "@/dto/dto";
import { useUserStore } from "@/stores/counter";
import _ from "lodash";

const userStore = useUserStore();
const state = reactive({
  serverName: "",
  serverList: [],
  gridsList: [],
  servantId: 0,
  servantLanguage: "",
  serverVersion: 0,
});
const currentItem = ref();
async function handleOpen(item: Item) {
  currentItem.value = item;
  console.log("item", item);
  const grids = await api.queryGrid({ id: item.id });
  state.gridsList = grids.data;
  state.servantId = item.id;
  console.log("state.grids", state.gridsList);
  state.serverName = item.serverName;
  state.servantLanguage = item.language;
  const serverVersion = await api.getPropertyByKey(`server.version.${item.id}`);
  if (_.isArray(serverVersion.data)) {
    state.serverVersion = Number(serverVersion.data[0].value) || 0;
  } else {
    state.serverVersion = 0;
  }
  console.log("serverVersion", serverVersion);
}

async function fetchServerList() {
  const resp = await api.getServerList(userStore.userInfo.id);
  state.serverList = (resp.data || []).sort((a, b) => a.id - b.id);
}

onMounted(() => {
  fetchServerList();
});
</script>
