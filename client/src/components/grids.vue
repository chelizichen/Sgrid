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
          <el-button
            class="text"
            type="text"
            @click="releaseServer"
            :disabled="selectionGrid.length == 0"
            >Release</el-button
          >
        </div>
        <div class="flex-item">
          <div class="text">Restart</div>
        </div>
        <div class="flex-item">
          <el-button
            class="text"
            type="text"
            @click="batchShutdown"
            :disabled="selectionGrid.length == 0"
            >BatchShutDown</el-button
          >
        </div>
      </div>
    </el-card>
    <el-divider content-position="left">GridNodes</el-divider>
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
          <div>{{ gridStatus[scoped.row.status] || "offline" }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="pid" label="PID"></el-table-column>
      <el-table-column label="Type">
        <template #default="scoped">
          <div>{{ scoped.row.gridServant.language }}</div>
        </template>
      </el-table-column>
      <el-table-column label="Protocol">
        <template #default="scoped">
          <div>{{ scoped.row.gridServant.protocol }}</div>
        </template>
      </el-table-column>
      <el-table-column label="Shutdown">
        <template #default="scoped">
          <div>
            <div class="text" @click="shutDown(scoped.row)">Shutdown</div>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <el-divider content-position="left">
      <el-button type="text" @click="getLogList($props.gridsList)">StatLog</el-button>
    </el-divider>
    <el-table :data="statLogList" style="width: 100%; margin-top: 20px" border>
      <el-table-column prop="id" label="id" width="180" />
      <el-table-column prop="createTime" label="createTime" />
      <el-table-column label="Grid">
        <template #default="scoped">
          <el-button type="text"
            >{{ scoped.row.gridInfo.gridNode.ip }}:{{
              scoped.row.gridInfo.port
            }}</el-button
          >
        </template>
      </el-table-column>
      <el-table-column prop="stat" label="stat" />
    </el-table>
    <uploadComponent
      :upload-visible="state.uploadVisible"
      :serverName="$props.serverName"
      :servantId="$props.servantId"
      @CLOSE_UPLOAD_DIALOG="() => (state.uploadVisible = false)"
    ></uploadComponent>
    <releaseComponent
      :releaseVisible="state.releaseVisible"
      :serverName="$props.serverName"
      :releaseList="releaseList"
      :selectionGrid="selectionGrid"
      @CLOSE_RELEASE_DIALOG="() => (state.releaseVisible = false)"
      @RELEASE_SERVER_BY_ID="(val) => handleRelease(val)"
    >
    </releaseComponent>
  </div>
</template>
<script lang="ts" setup>
import { reactive, ref, watch } from "vue";
import uploadComponent from "./upload.vue";
import releaseComponent from "./release.vue";
import moment from "moment";
import api from "@/api/server";
import { ElMessage } from "element-plus";

const props = defineProps<{
  gridsList: any[];
  serverName: string;
  servantId: number;
}>();

async function getLogList(gridList) {
  const resp = await Promise.all(
    gridList.map(async (v) => {
      const list = await api.getStatLogList({
        id: v.id,
      });
      const ret = list.Data.list.map((item) => {
        item.gridInfo = v;
        return item;
      });
      return ret;
    })
  );
  const newArr: any[] = [];
  resp.forEach((v) => {
    newArr.push(...v);
  });

  statLogList.value = newArr.map((v) => {
    v.createTime = moment(v.createTime).format("YYYY-MM-DD HH:mm:ss");
    return v;
  });
}

const statLogList = ref([]);
watch(
  () => props.gridsList,
  async function (newVal) {
    getLogList(newVal);
  }
);

const state = reactive({
  uploadVisible: false,
  releaseVisible: false,
});
const releaseList = ref([]);
async function releaseServer() {
  const data = await api.getUploadList({
    id: props.servantId,
  });
  state.releaseVisible = true;
  releaseList.value = data.Data;
}
async function handleRelease(id) {
  const releaseItem = releaseList.value.find((v) => v.id == id);
  const servantBaseInfo = props.gridsList[0];
  console.log("selectionGrid", selectionGrid);

  const body = {
    serverName: props.serverName,
    filePath: releaseItem.filePath,
    serverLanguage: servantBaseInfo.gridServant.language,
    serverProtocol: servantBaseInfo.gridServant.protocol,
    execPath: servantBaseInfo.gridServant.execPath,
    servantGrids: selectionGrid.value.map((v) => ({
      ip: v.gridNode.ip,
      port: v.port,
      gridId: v.id,
    })),
  };

  console.log("body", body);

  const data = await api.releaseServer(body);
  ElMessage.success("success!");
  state.releaseVisible = false;
}
const selectionGrid = ref([]);
function handleSelectionChange(value) {
  selectionGrid.value = value;
}
const gridStatus = {
  "1": "online",
  "0": "offline",
};

function batchShutdown() {
  console.log("releaseList", selectionGrid.value);
  const body = {
    req: selectionGrid.value
      .filter((v) => v.status != 0 && v.status)
      .map((v) => ({
        pid: v.pid,
        gridId: v.id,
        host: v.gridNode.ip,
        port: v.port,
      })),
  };

  api.shutdownServer(body);
  console.log("body", body);
}

function shutDown(v) {
  const body = {
    req: {
      pid: v.pid,
      gridId: v.id,
      host: v.gridNode.ip,
      port: v.port,
    },
  };
  api.shutdownServer(body);
  console.log("body", body);
}
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
