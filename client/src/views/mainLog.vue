<template>
  <div class="w-full p-6 border-solid border-gray-100 border-2">
    <el-form :inline="true" :model="pagination">
      <el-form-item>
        <el-input v-model="pagination.keyword" clearable placeholder="search" />
      </el-form-item>
      <el-form-item>
        <el-button @click="getList">Select</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="tableData" border>
      <el-table-column prop="id" label="id" width="100" />
      <el-table-column prop="type" label="type" />
      <el-table-column prop="info" label="info" />
      <el-table-column prop="createTime" label="createTime" width="180" />
    </el-table>
    <el-pagination
      background
      layout="prev, pager, next"
      :total="total"
      class="pagination"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script lang="ts" setup>
import api from "@/api/server";
import moment from "moment";
import { onMounted, ref } from "vue";

const pagination = ref({
  keyword: "",
  offset: 0,
});
const tableData = ref([]);
const total = ref(0);

async function getList() {
  const data = await api.getMainLogger(pagination.value);
  tableData.value = data.data.map((v) => {
    v.createTime = moment(v.createTime).format("YYYY-MM-DD HH:mm:ss");
    return v;
  });
  total.value = data.total;
}

onMounted(() => {
  getList();
});

async function handleCurrentChange(curr: number) {
  pagination.value.offset = (curr - 1) * 20;
  await getList();
}
</script>

<style scoped lang="less">
.container {
  padding: 20px;
  width: 100%;
  box-sizing: border-box;
}
.pagination {
  float: right;
  margin-top: 20px;
}
</style>
