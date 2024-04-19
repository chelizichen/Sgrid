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
    <el-divider content-position="left">
      <el-button type="text" @click="$emit('checkStatus')"
        ><el-icon><Loading /></el-icon> Grids
      </el-button></el-divider
    >
    <el-table
      :data="props.gridsList"
      style="width: 100%"
      border
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column label="Grid">
        <template #default="scoped">
          <el-button type="text" @click="toLog(scoped.row.gridNode.ip)"
            >{{ scoped.row.gridNode.ip }}:{{ scoped.row.port }}</el-button
          >
        </template>
      </el-table-column>
      <el-table-column label="Status">
        <template #default="scoped">
          <div :class="gridStatus[scoped.row.status] || 'offline'">
            {{ gridStatus[scoped.row.status] || "offline" }}
          </div>
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
      <el-button type="text" @click="getLogList($props.gridsList)"
        ><el-icon><Loading /></el-icon> StatLog
      </el-button>
    </el-divider>
    <el-table :data="statLogList" style="width: 100%; margin-top: 20px" border>
      <el-table-column prop="id" label="id" width="180" />
      <el-table-column prop="name" label="name" width="180" />
      <el-table-column prop="threads" label="threads" width="180" />
      <el-table-column prop="isRunning" label="isRunning" width="180" />
      <el-table-column prop="createTime" label="createTime" />
      <el-table-column label="Grid">
        <template #default="scoped">
          <div>{{ scoped.row.gridInfo.gridNode.ip }}:{{ scoped.row.gridInfo.port }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="stat" label="stat" />

      <el-table-column prop="pid" label="pid"> </el-table-column>
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
import { useRouter } from "vue-router";

const props = defineProps<{
  gridsList: any[];
  serverName: string;
  servantId: number;
}>();
const emits = defineEmits(["checkStatus"]);
async function getLogList(gridList) {
  const resp = await Promise.all(
    gridList.map(async (v) => {
      const list = await api.getStatLogList({
        id: v.id,
      });
      const ret = list.data.list.map((item) => {
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
    await checkStat(newVal);
    await getLogList(newVal);
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
  releaseList.value = data.data;
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

async function batchShutdown() {
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

  const data = await api.shutdownServer(body);
  if (data.code) {
    ElMessage.error(data.message);
  }
  ElMessage.success("关闭成功");
}

async function shutDown(v) {
  const body = {
    req: [
      {
        pid: v.pid,
        gridId: v.id,
        host: v.gridNode.ip,
        port: v.port,
      },
    ],
  };
  const data = await api.shutdownServer(body);
  if (data.code) {
    ElMessage.error(data.message);
  }
  ElMessage.success("关闭成功");
}

const router = useRouter();
function toLog(host: string) {
  const text = router.resolve({
    path: "/logpage",
    query: {
      host: host,
      serverName: props.serverName,
    },
  });
  window.open(text.href, "_blank");
}

function checkStat(list) {
  const body = list.map((v) => {
    return {
      pid: v.pid,
      host: v.gridNode.ip,
      gridId: v.id,
    };
  });
  api.checkStat({
    hostPids: body,
  });
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
.online {
  color: #55bd55;
}
.offline {
  color: red;
}
</style>
