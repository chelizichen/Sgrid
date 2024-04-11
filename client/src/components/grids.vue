<script lang="ts">
export default {
  name: "grids-component",
};
</script>

<template>
  <div>
    <el-card shadow="hover" v-if="props.serverName" style="margin-bottom: 10px">
      <div class="card">
        <div class="flex-item">
          <div style="font-weight: 700">ServerName</div>
          <div class="text" @click="state.uploadVisible = true">
            {{ props.serverName }}
          </div>
        </div>

        <div class="flex-item">
          <div class="text">Shutdown</div>
        </div>
        <div class="flex-item">
          <div class="text">Release</div>
        </div>
        <div class="flex-item">
          <div class="text">Restart</div>
        </div>
      </div>
    </el-card>
    <el-table
      :data="props.gridsList"
      style="width: 100%"
      border
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column label="Grid">
        <template #default="scoped">
          <el-button type="text"
            >{{ scoped.row.gridNode.ip }}:{{ scoped.row.port }}</el-button
          >
        </template>
      </el-table-column>
      <el-table-column label="Status">
        <template #default="scoped">
          <div>{{ gridStatus[scoped.row.status] }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="pid" label="PID"></el-table-column>
      <el-table-column label="Type">
        <template #default="scoped">
          <div>{{ scoped.row.gridServant.language }}</div>
        </template>
      </el-table-column>
    </el-table>
    <uploadComponent
      :upload-visible="state.uploadVisible"
      :serverName="$props.serverName"
      :servantId="$props.servantId"
      @CLOSE_UPLOAD_DIALOG="() => (state.uploadVisible = false)"
    ></uploadComponent>
  </div>
</template>
<script lang="ts" setup>
import type { Item } from "@/dto/dto";
import { reactive } from "vue";
import uploadComponent from "./upload.vue";
const props = defineProps<{
  gridsList: any[];
  serverName: string;
  servantId: number;
}>();

const emits = defineEmits(["handleOpen"]);

const state = reactive({
  uploadVisible: false,
});
function handleOpen(item: Item) {
  emits("handleOpen", item);
}

function handleSelectionChange(value) {
  console.log("value", value);
}
const gridStatus = {
  "1": "online",
  "0": "offline",
};
</script>
<style scoped>
.card {
  height: 70px;
  display: flex;
  align-items: center;
  justify-content: space-around;
}
.text {
  color: rgb(207, 15, 124);
  cursor: pointer;
}
.flex-item {
  text-align: center;
  width: 15%;
  padding: 10px;
}
</style>
