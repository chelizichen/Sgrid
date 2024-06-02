<template>
  <div>
    <el-form :inline="true">
      <el-form-item>
        <el-button type="primary" @click="editUservisible = true">创建角色</el-button>
      </el-form-item>
    </el-form>

    <el-table
      :data="userList"
      border
      highlight-current-row
      @current-change="handleCurrentChange"
    >
      <el-table-column type="index" label="序号" width="90"></el-table-column>
      <el-table-column prop="name" label="角色名"></el-table-column>
      <el-table-column prop="description" label="角色描述"></el-table-column>
      <el-table-column prop="createTime" label="创建时间"></el-table-column>
      <el-table-column label="操作">
        <template #default="scoped">
          <div style="display: flex; flex-wrap: wrap">
            <el-button @click="handleEdit(scoped.row)">修改</el-button>
            <el-button @click="handleDel(scoped.row)" type="danger">删除</el-button>
          </div>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="editUservisible" title="角色管理">
      <el-form label-width="100px">
        <el-form-item label="角色名">
          <el-input v-model="editRoleObj.name"></el-input>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="editRoleObj.description" type="textarea"></el-input>
        </el-form-item>
        <el-form-item label="操作">
          <el-button @click="reset">Reset</el-button>
          <el-button type="primary" @click="submitEdit">Submit</el-button></el-form-item
        >
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { delRole, getRole, saveRole } from "@/api/system";
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
import _ from "lodash";
import { onMounted, ref } from "vue";

type RoleVo = {
  id: number;
  name: string;
  createTime: string;
  description: string;
};

const userList = ref();
async function getRoleList() {
  const servantsResp = await getRole(undefined);
  userList.value = servantsResp.data;
  console.log("servantResp", servantsResp);
}

const editUservisible = ref(false);

const editRoleObj = ref<Partial<RoleVo>>({
  createTime: "",
  id: 0,
  name: "",
  description: "",
});
function handleEdit(row: RoleVo) {
  editRoleObj.value = _.cloneDeep(row);
  editUservisible.value = true;
}
async function submitEdit() {
  if (!editRoleObj.value.id) {
    editRoleObj.value.createTime = undefined;
  }
  const data = await saveRole(editRoleObj.value);
  if (data.code) {
    return ElNotification.error(data.message);
  }
  editUservisible.value = false;
  getRoleList();
  return ElNotification.success("success");
}

function reset() {
  editRoleObj.value.createTime = "";
  editRoleObj.value.id = 0;
  editRoleObj.value.name = "";
}
onMounted(async () => {
  await getRoleList();
});
const emits = defineEmits(["recvRole"]);
function handleCurrentChange(row: RoleVo) {
  if (row) {
    emits("recvRole", row.id);
  }
}

async function handleDel(row: RoleVo) {
  ElMessageBox.confirm("确认删除?", {
    confirmButtonText: "确认",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(async () => {
      await delRole(row.id);
      await getRoleList();
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
</script>

<style scoped></style>
