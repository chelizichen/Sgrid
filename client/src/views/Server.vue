<script lang="ts">
export default {
  name: "server-component",
};
</script>

<template>
  <div>
    <el-container>
      <el-aside width="200px" class="py-2.5 h-full">
        <aside-component :server-list="state.serverList" @handle-open="handleOpen"></aside-component>
      </el-aside>
      <el-main>
        <el-tabs v-model="editableTabsValue" type="card" class="demo-tabs" closable @tab-remove="removeTab">
          <el-tab-pane v-for="item in editableTabs" :key="item.serverName" :label="item.serverName"
            :name="item.serverName">
            <controller v-if="editableTabsValue == HOME" @handle-open="handleOpen"></controller>
            <template v-else>
              <gridsComponent v-if="editableTabsValue == item.serverName && editableTabsValue != HOME"
                :grids-list="state.gridsList" :server-name="state.serverName" :servant-id="state.servantId"
                :servantLanguage="state.servantLanguage" @check-status="handleOpen(currentItem)"
                :server-version="state.serverVersion"></gridsComponent>
            </template>
          </el-tab-pane>
        </el-tabs>
      </el-main>
      <router-view></router-view>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import asideComponent from "@/components/aside.vue";
import gridsComponent from "@/components/grids.vue";
import controller from "@/components/controller.vue";
import { onMounted, reactive, ref, watch } from "vue";
import api from "../api/server";
import type { Item, T_Grid, T_ServerList } from "@/dto/dto";
import { useUserStore } from "@/stores/counter";
import _ from "lodash";

const userStore = useUserStore();
const state = reactive({
  serverName: "",
  serverList: <T_ServerList[]>[],
  gridsList: <T_Grid[]>[],
  servantId: -1,
  servantLanguage: "",
  serverVersion: 0,
  isHandlingOpen: false, // 添加一个新的属性来标记是否正在处理 handleOpen
});
const currentItem = ref();

async function handleOpen(item: Partial<Item>) {
  if (state.isHandlingOpen) return; // 如果正在处理中，直接返回
  state.isHandlingOpen = true; // 设置为正在处理
  console.log('debug.handleOpen', item);
  currentItem.value = item;
  state.servantId = item.id!;
  if (item.serverName !== HOME) {
    const grids = await api.queryGrid({ id: item.id });
    state.gridsList = grids.data;
    const serverVersion = await api.getPropertyByKey(`server.version.${item.id}`);
    if (_.isArray(serverVersion.data)) {
      state.serverVersion = Number(serverVersion.data[0].value) || 0;
    } else {
      state.serverVersion = 0;
    }
  }
  state.serverName = item.serverName!;
  state.servantLanguage = item.language!;
  editableTabsValue.value = item.serverName!
  if (editableTabs.value.find(tab => tab.serverName === item.serverName)) return
  editableTabs.value.push(item)
  state.isHandlingOpen = false; // 处理完成后重置标志
}

async function fetchServerList() {
  const resp = await api.getServerList(userStore.userInfo.id);
  state.serverList = (resp.data || []).sort((a, b) => a.id - b.id);
}

onMounted(() => {
  fetchServerList();
});

const HOME = "Home"
const editableTabsValue = ref(HOME)
const editableTabs = ref<Partial<Item>[]>([{
  serverName: HOME,
  id: -1,
  language: "zh-cn",
}])
const removeTab = (targetName: string) => {
  console.log('targetName', targetName);
  if (targetName === HOME) {
    return
  }
  editableTabs.value = editableTabs.value.filter((tab) => tab.serverName !== targetName)
  editableTabsValue.value = HOME
}

watch(editableTabsValue, function (val, oldVal) {
  if (val === oldVal) {
    return
  }
  const item = editableTabs.value.find((item) => item.serverName === val);
  if (item) {
    handleOpen(item)
  }
})
</script>