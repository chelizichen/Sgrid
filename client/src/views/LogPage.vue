<template>
  <div style="display: flex">
    <div style="width: 30%">
      <el-form label-width="88px">
        <el-form-item>
          <div class="title">
            <el-icon style="color: rgb(207, 90, 124); font-size: 36px"><Help /></el-icon>
            SgridLog
          </div>
        </el-form-item>

        <el-form-item label="logFile">
          <el-input v-model="state.logFile" :disabled="true"></el-input>
        </el-form-item>
        <el-form-item label="pattern">
          <el-input v-model="state.pattern"></el-input
        ></el-form-item>
        <el-form-item label="rows">
          <el-select v-model="state.rows">
            <el-option
              v-for="item in rowSelect"
              :label="item"
              :value="item"
              :key="item"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="logList" v-for="(item, index) in logFileList" :key="index">
          <el-button
            :key="index"
            type="text"
            style="display: block"
            @click="state.logFile = item"
            >{{ item }}
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button @click="getLog">Search</el-button>
        </el-form-item>
      </el-form>
    </div>
    <div
      style="
        background-color: black;
        height: 100vh;
        padding: 5px 10px;
        width: 100%;
        overflow: scroll;
      "
    >
      <div style="color: aliceblue; margin: 2px" v-for="item in logger" :key="item">
        {{ item }}
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import API from "@/api/server";
import { computed, onMounted, reactive, ref } from "vue";
import { useRoute } from "vue-router";
const route = useRoute();
const query = computed(() => route.query);
const logFileList = ref([]);
const rowSelect = [10, 50, 100, 500, 1000];
const state = reactive({
  logFile: "",
  pattern: "",
  rows: 10,
});
const body = computed(() => {
  return {
    ...query.value,
    ...state,
  };
});
const logger = ref([]);
async function init() {
  const data = await API.getLogFileList({
    host: query.value.host,
    serverName: query.value.serverName,
  });
  logFileList.value = data.data;
}
async function getLog() {
  const res = await API.getLog(body.value);
  logger.value = res.data.split("\n");
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
