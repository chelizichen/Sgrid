<template>
  <div style="display: flex">
    <div style="width: 30%">
      <el-card>
        <el-input v-model="state.logFile" :disabled="true"></el-input>
        <el-input v-model="state.pattern"></el-input>
        <el-select v-model="state.rows">
          <el-option
            v-for="item in rowSelect"
            :label="item"
            :value="item"
            :key="item"
          ></el-option>
        </el-select>
        <div v-for="(item, index) in logFileList" :key="index">
          <el-button
            :key="index"
            type="text"
            style="display: block"
            @click="state.logFile = item"
            >{{ item }}
          </el-button>
        </div>
        <el-button @click="getLog">Search</el-button>
      </el-card>
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
  logFileList.value = data.Data;
}
async function getLog() {
  const res = await API.getLog(body.value);
  logger.value = res.Data.split("\n");
}
onMounted(() => {
  init();
});
</script>

<style scoped></style>
