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
          <el-button class="text" type="text" @click="showConfiguration"
            >Configuration</el-button
          >
        </div>
        <div class="flex-item">
          <el-button
            class="text"
            type="text"
            @click="restartServer"
            :disabled="selectionGrid.length == 0"
            >Restart</el-button
          >
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
    <el-table
      :data="props.gridsList"
      style="width: 100%"
      border
      @selection-change="handleSelectionChange"
    >
      <el-table-column type="selection" width="55" />
      <el-table-column label="Grid">
        <template #default="scoped">
          <el-button type="text" @click="toLog(scoped.row)"
            >{{ scoped.row.gridNode.ip }}:{{ scoped.row.port }}</el-button
          >
        </template>
      </el-table-column>
      <el-table-column prop="gridServant.serverName" label="serverName">
        <template #default="scoped">
          <template v-if="scoped.row.gridServant.preview">
            <el-button
              type="text"
              @click="toPreview(scoped.row.gridServant.preview)"
              style="color: #4a4aca"
              >{{ scoped.row.gridServant.serverName }}</el-button
            >
          </template>
          <template v-else>
            <span>{{ scoped.row.gridServant.serverName }}</span>
          </template>
        </template>
      </el-table-column>
      <el-table-column label="Status">
        <template #default="scoped">
          <div
            :class="gridStatus[scoped.row.status] || 'offline'"
            @click="$emit('checkStatus')"
          >
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
      <el-table-column label="Delete">
        <template #default="scoped">
          <div>
            <div class="danger" @click="deleteGridById(scoped.row)">Delete</div>
          </div>
        </template>
      </el-table-column>
    </el-table>
    <el-divider content-position="left">
      <el-button type="text" @click="getLogList($props.gridsList)"
        ><el-icon style="font-size: large; font-weight: 900; color: gray"
          ><Loading
        /></el-icon>
      </el-button>
    </el-divider>
    <el-table :data="statLogList" style="width: 100%; margin-top: 20px" border>
      <el-table-column prop="id" label="id" width="180" />
      <el-table-column prop="name" label="name" width="180">
        <template #default="scoped">
          <div>{{ scoped.row.name || "--" }}</div>
        </template>
      </el-table-column>
      <el-table-column label="Grid">
        <template #default="scoped">
          <div>{{ scoped.row.gridInfo.gridNode.ip }}:{{ scoped.row.gridInfo.port }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="pid" label="pid"> </el-table-column>
      <el-table-column prop="threads" label="threads" width="180">
        <template #default="scoped">
          <div>{{ scoped.row.threads || "--" }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="isRunning" label="isRunning" width="180">
        <template #default="scoped">
          <div>{{ scoped.row.isRunning || "--" }}</div>
        </template>
      </el-table-column>
      <el-table-column prop="createTime" label="createTime" />
      <el-table-column prop="stat" label="behavior" />
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
    <servantConf
      :servantId="showConfigurationId"
      :dialogVisible="showConfigurationVisible"
      @CLOSE_RELEASE_DIALOG="() => (showConfigurationVisible = false)"
    ></servantConf>
  </div>
</template>
<script lang="ts" setup>
import { reactive, ref, watch } from "vue";
import uploadComponent from "./upload.vue";
import releaseComponent from "./release.vue";
import servantConf from "./servantConf.vue";
import moment from "moment";
import api from "@/api/server";
import { ElNotification } from "element-plus";
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

  statLogList.value = newArr
    .map((v) => {
      v.createTime = moment(v.createTime).format("YYYY-MM-DD HH:mm:ss");
      return v;
    })
    .sort((a, b) => b.id - a.id);
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
async function handleRelease(id: number) {
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
    servantId: Number(props.servantId),
  };

  console.log("body", body);

  const data = await api.releaseServer(body);
  if (data.code) {
    return ElNotification.error(data.message);
  }
  ElNotification.success("success!");
  state.releaseVisible = false;
}

async function restartServer() {
  const servantBaseInfo = props.gridsList[0];
  const body = {
    serverName: props.serverName,
    filePath: "",
    serverLanguage: servantBaseInfo.gridServant.language,
    serverProtocol: servantBaseInfo.gridServant.protocol,
    execPath: servantBaseInfo.gridServant.execPath,
    servantGrids: selectionGrid.value.map((v) => ({
      ip: v.gridNode.ip,
      port: v.port,
      gridId: v.id,
    })),
    servantId: Number(props.servantId),
  };
  const data = await api.restartServer(body);
  if (data.code) {
    return ElNotification.error(data.message);
  }
  ElNotification.success("success!");
}

const selectionGrid = ref([]);
function handleSelectionChange(value) {
  selectionGrid.value = value;
}
const gridStatus: any = {
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
    ElNotification.error(data.message);
  }
  ElNotification.success("关闭成功");
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
    ElNotification.error(data.message);
  }
  ElNotification.success("关闭成功");
}

const router = useRouter();
function toLog(row) {
  const text = router.resolve({
    path: "/logpage",
    query: {
      host: row.gridNode.ip,
      serverName: props.serverName,
      gridId: row.id,
    },
  });
  window.open(text.href, "_blank");
}

function toPreview(path: string) {
  window.open(path, "_blank");
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

async function deleteGridById(row) {
  if (row.status) {
    return ElNotification.error("error/client :: this grid still alive");
  }
  const data = await api.deleteGrid({
    id: row.id,
  });
  if (data.code) {
    return ElNotification.error(data.message);
  }
  ElNotification.success("delete success");
  emits("checkStatus");
}

const showConfigurationVisible = ref(false);
const showConfigurationId = ref(0);
async function showConfiguration() {
  showConfigurationVisible.value = true;
  showConfigurationId.value = props.servantId;
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
.danger {
  color: red;
  cursor: pointer;
}
.flex-item {
  text-align: center;
  width: 15%;
  padding: 10px;
}
.online {
  color: #55bd55;
  cursor: pointer;
}
.offline {
  color: red;
  cursor: pointer;
}
</style>
