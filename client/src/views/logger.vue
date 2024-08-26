<template>
  <div style="display: flex; height: 95vh">
    <div style="width: 30%; padding: 0 20px 0 0">
      <el-form label-width="88px">
        <el-form-item label="logFile">
          <el-input v-model="state.logFile" :disabled="true"></el-input>
        </el-form-item>
        <el-form-item label="pattern">
          <el-input v-model="state.pattern"></el-input
        ></el-form-item>
        <el-form-item label="rows">
          <el-select v-model="state.size">
            <el-option
              v-for="item in rowSelect"
              :label="item"
              :value="item"
              :key="item"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="logList" :key="index">
          <el-select v-model="state.logFile">
            <el-option
              v-for="(item, index) in logFileList"
              :label="item.logType + '_' + item.dateTime"
              :value="index"
              :key="item"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="Operate">
          <el-button @click="getLog">Search</el-button>
        </el-form-item>
        <el-form-item label="Pagination">
          <el-pagination
            layout="prev, pager, next"
            :total="total"
            size="small"
            :page-size="state.size"
            @current-change="handleCurrentChange"
          />
        </el-form-item>
      </el-form>
    </div>
    <div style="width: 100%">
      <div
        style="
          background-color: black;
          height: 60vh;
          padding: 5px 10px;
          width: 100%;
          overflow: scroll;
        "
      >
        <div style="color: aliceblue; margin: 2px" v-for="item in logger" :key="item">
          {{ item }}
        </div>
      </div>

      <div
        style="border-top: 5px solid #ddd; height: 37vh; width: 100%; overflow: scroll"
      >
        <iframe :src="SHELL_PATH" style="width: inherit; height: inherit"></iframe>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import API from "@/api/server";
import { ElNotification } from "element-plus";
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute } from "vue-router";
const route = useRoute();
const query = computed(() => route.query);
const logFileList = ref([]);
const rowSelect = [10, 50, 100, 500, 1000];
const state = reactive({
  logFile: 0,
  pattern: "",
  size: 10,
  offset: 0,
});

const selectFile = computed(() => {
  return logFileList.value[state.logFile];
});

const SHELL_PATH = ref("");

const body = computed(() => {
  return Object.assign(
    {},
    {
      ...query.value,
      ...selectFile.value,
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
  logger.value = res.data;
  total.value = res.total;
}

async function handleCurrentChange(curr: number) {
  console.log("curr", curr);
  state.offset = (curr - 1) * state.size;
  await getLog();
}
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
