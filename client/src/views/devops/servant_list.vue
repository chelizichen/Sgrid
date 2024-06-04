<template>
  <div>
    <el-table border :data="servantList">
      <el-table-column type="index" label="序号" width="180"></el-table-column>
      <el-table-column prop="serverName" label="serverName"></el-table-column>
      <el-table-column prop="stat" label="服务状态">
        <template #default="scoped">
          <div>
            <span v-if="!scoped.row.stat">运行中</span>
            <span v-if="scoped.row.stat == -1">已停用</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column prop="servantGroupId" label="servantGroupId"></el-table-column>
      <el-table-column prop="execPath" label="execPath"></el-table-column>
      <el-table-column prop="language" label="language"></el-table-column>
      <el-table-column prop="protocol" label="protocol"></el-table-column>
      <el-table-column label="操作">
        <template #default="scoped">
          <el-button @click="updateServant(scoped.row)">修改</el-button>
          <el-button
            v-if="scoped.row.stat != -1"
            @click="deleteServant(scoped.row, -1)"
            type="danger"
            >停用</el-button
          >
          <el-button v-else @click="deleteServant(scoped.row, 0)" type="success"
            >启用</el-button
          >
        </template>
      </el-table-column>
    </el-table>
  </div>
  <el-dialog v-model="editDialogVisible" title="编辑服务信息" width="50%">
    <el-form :model="servant" label-width="100px">
      <el-form-item label="ID">
        <el-input v-model="servant.id" disabled></el-input>
      </el-form-item>
      <el-form-item label="服务器名称">
        <el-input v-model="servant.serverName"></el-input>
      </el-form-item>
      <el-form-item label="语言">
        <el-input v-model="servant.language"></el-input>
      </el-form-item>
      <el-form-item label="协议">
        <el-input v-model="servant.protocol"></el-input>
      </el-form-item>
      <el-form-item label="执行路径">
        <el-input v-model="servant.execPath"></el-input>
      </el-form-item>
      <el-form-item label="服务组ID">
        <el-skeleton style="width: 100%" :loading="editLoading" animated>
          <el-select v-model="servant.servantGroupId">
            <el-option
              v-for="item in groupList"
              :label="item.tagEnglishName"
              :key="item.id"
              :value="Number(item.id)"
            ></el-option>
          </el-select>
        </el-skeleton>
      </el-form-item>
      <el-button type="primary" @click="confirmUpdateServant">更新</el-button>
      <el-button @click="editDialogVisible = false">取消</el-button>
    </el-form>
  </el-dialog>
</template>

<script setup lang="ts">
import api from "@/api/server";
import { useUserStore } from "@/stores/counter";
import { ElMessage, ElMessageBox } from "element-plus";
import _ from "lodash";
import { onMounted, ref } from "vue";

const userStore = useUserStore();
const servantDemo = {
  id: 2,
  serverName: "ShellServer",
  language: "node",
  protocol: "grpc",
  execPath: "service_go",
  servantGroupId: 1,
};
type T_Servant = typeof servantDemo;
const groupList = ref<Array<{ tagEnglishName: string; tagName: string; id: number }>>([]);
const servant = ref<Partial<T_Servant>>({});
const servantList = ref<Array<T_Servant>>([]);
const editDialogVisible = ref(false);
const editLoading = ref(true);
async function getServantList() {
  const servantsResp = await api.getServants(userStore.userInfo.id);
  servantList.value = servantsResp.data;
  console.log("servantResp", servantsResp);
}
onMounted(async () => {
  await getServantList();
});

function updateServant(row: T_Servant) {
  editDialogVisible.value = true;
  servant.value = _.cloneDeep(row);
  api.getGroup(userStore.userInfo.id).then((getGroup) => {
    console.log("getGroup", getGroup);
    editLoading.value = false;
    groupList.value = getGroup.data;
  });
  console.log("row", row);
}

function confirmUpdateServant() {}

async function deleteServant(row: T_Servant, stat: number) {
  const text = stat == -1 ? "停用" : " 启用";
  ElMessageBox.confirm(`确认${text}?`, {
    confirmButtonText: "确认",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(async () => {
      const resp = await api.delServant(row.id, stat);
      if (resp.code) {
        return ElMessage.error({
          type: "error",
          message: resp.message,
        });
      }
      await getServantList();
      ElMessage({
        type: "success",
        message: text + "成功",
      });
    })
    .catch(() => {
      ElMessage({
        type: "info",
        message: "取消" + text,
      });
    });
}
</script>

<style scoped></style>
