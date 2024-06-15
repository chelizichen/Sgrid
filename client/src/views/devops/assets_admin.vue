<template>
  <div>
    <el-table border :data="assetsList">
      <el-table-column type="index" label="序号" width="180"></el-table-column>
      <el-table-column prop="servantId" label="服务ID"></el-table-column>
      <el-table-column label="服务名称">
        <template #default="scoped">
          {{ serverStore.getServerNameById(scoped.row.servantId) }}
        </template>
      </el-table-column>
      <el-table-column prop="gridId" label="网格ID"></el-table-column>
      <el-table-column prop="activeTime" label="生效时间"></el-table-column>
      <el-table-column prop="expireTime" label="过期时间"></el-table-column>
      <el-table-column prop="mark" label="备注"></el-table-column>
      <el-table-column prop="operateUserId" label="操作人"></el-table-column>
      <el-table-column prop="updateTime" label="更新时间"></el-table-column>
      <el-table-column prop="createTime" label="创建时间"></el-table-column>
      <el-table-column label="操作" align="center">
        <template #default="scoped">
          <el-button @click="delAssetsById(scoped.row.gridId)" type="danger"
            >删除</el-button
          >
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import api from "@/api/assets";
import { ref, onMounted } from "vue";
import { ElMessageBox, ElMessage } from "element-plus";
import { useServersStore } from "@/stores/counter";
import moment from "moment";
const serverStore = useServersStore();
const queryParams = ref({
  offset: 0,
  size: 10,
});
const F_M_T = "YYYY-MM-DD HH:mm:ss";
const assetsList = ref<Array<any>>([]);
async function getAssetsList() {
  const servantsResp = await api.getList(queryParams.value);
  assetsList.value = servantsResp.data.map((v) => {
    v.activeTime = moment(v.activeTime).format(F_M_T);
    v.createTime = moment(v.createTime).format(F_M_T);
    v.expireTime = moment(v.expireTime).format(F_M_T);
    v.updateTime = moment(v.updateTime).format(F_M_T);
    return v;
  });
  console.log("servantResp", servantsResp);
}

async function delAssetsById(id: number) {
  ElMessageBox.confirm("确认删除?", {
    confirmButtonText: "确认",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(async () => {
      const resp = await api.delAssert({ id });
      await getAssetsList();
      if (resp.code) {
        return ElMessage.error({
          type: "error",
          message: resp.message,
        });
      }
      ElMessage({
        type: "success",
        message: "删除成功",
      });
    })
    .catch(() => {
      ElMessage({
        type: "info",
        message: "取消删除",
      });
    });
}
onMounted(async () => {
  await getAssetsList();
});
</script>

<style scoped></style>
