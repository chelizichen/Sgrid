<template>
  <div class="flex w-full">
    <div class="w-1/4 h-screen bg-slate-50">
      <el-form label-width="88px">
        <el-form-item label="logFile">
          <el-input v-model="state.LogFileName" :disabled="true"></el-input>
        </el-form-item>
        <el-form-item label="pattern">
          <el-input v-model="state.keyword"></el-input></el-form-item>
        <el-form-item label="rows">
          <el-select v-model="state.len">
            <el-option v-for="item in rowSelect" :label="item" :value="item" :key="item"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="logList">
          <el-select v-model="state.LogFileName">
            <el-option v-for="(item, index) in logFileList" :label="getFileName(item)" :value="getFileName(item)"
              :key="index"></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="Operate">
          <el-button @click="getLog">Search</el-button>
          <el-button @click="deleteLog" type="primary">Delete</el-button>
        </el-form-item>
        <!-- <el-form-item label="Pagination">
          <el-pagination layout="prev, pager, next" :total="total" size="small" :page-size="state.size"
            @current-change="handleCurrentChange" />
        </el-form-item> -->
        <!-- <el-form-item label="Shell">
          <el-button @click="showShell">Shell</el-button>
        </el-form-item> -->
      </el-form>
    </div>
    <div class="w-3/4 h-screen">
      <div class="bg-black px-4 pt-2 h-full w-full overflow-scroll">
        <div v-for="item in logger" :key="item" class="break-all text-wrap text-white w-full">
          <div v-html="item"></div>
        </div>
      </div>
    </div>
    <el-dialog v-model="shellVisible" class="w-3/4" style="height: 800px" title="SgridShell">
      <div class="w-full h-full">
        <iframe :src="SHELL_PATH" class="w-full" style="height: 700px"></iframe>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts" name="logger-page">
import API from "@/api/server";
import { ElNotification, ElMessageBox } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute } from "vue-router";
const shellVisible = ref(false);
// const showShell = () => (shellVisible.value = !shellVisible.value);
const route = useRoute();
const query = computed(() => route.query);
const logFileList = ref([]);
const rowSelect = [10, 50, 100, 500, 1000];
const state = reactive({
  LogFileName: "",
  keyword: "",
  len: 10,
  offset: 0,
});

function getFileName(path: string) {
  const arr = path.split("/");
  return arr[arr.length - 1];
}

const SHELL_PATH = ref("");

const body = computed(() => {
  return Object.assign(
    {},
    {
      ...query.value,
      ...state,
    },
    { gridId: Number(query.value.gridId) }
  );
});
const logger = ref([]);
const total = ref(0);
async function init() {
  const data = await API.getLogFileList({
    host: query.value.host,
    serverName: query.value.serverName,
    gridId: query.value.gridId,
  });
  logFileList.value = data.data;
  const webConsole = await API.getPropertyByKey(`WebConsole@${query.value.host}`);
  if (!webConsole.code) {
    SHELL_PATH.value = webConsole.data[0].value;
  }
  console.log("webConsole", webConsole);
}
async function getLog() {
  const res = await API.getLog(body.value);
  if (res.code) {
    return ElNotification.error(`error:${res.message}`);
  }
  if (!res.data) {
    logger.value = []
    return
  }
  logger.value = res.data.map((v: string) => {
    if (state.keyword) {
      v = v.replaceAll(
        state.keyword,
        `<span class="text-red-500">${state.keyword}</span>`
      );
    }
    return v;
  });
  total.value = res.total;
}
async function deleteLog() {
  ElMessageBox.confirm("确认删除该节点？删除后不可恢复!", "Confirm", {
    confirmButtonText: "OK",
    cancelButtonText: "Cancel",
  }).then(async () => {
    const res = await API.deleteByLogType(body.value);
    if (res.code) {
      return ElNotification.error(`删除失败:${res.message}`);
    }
    ElNotification.success("删除成功");
    init();
  });
}
// async function handleCurrentChange(curr: number) {
//   console.log("curr", curr);
//   state.offset = (curr - 1) * state.len;
//   await getLog();
// }
onMounted(() => {
  init();
});
</script>

<style scoped>
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
