<template>
  <div>
    <el-form :inline="true">
      <el-form-item>
        <el-input v-model="parmas.keyword" placeholder="请输入查询的用户名"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button @click="getServantList(true)">查询</el-button>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="createUser()">创建用户</el-button>
      </el-form-item>
    </el-form>

    <el-table border :data="userList">
      <el-table-column type="index" label="序号" width="180"></el-table-column>
      <el-table-column prop="userName" label="用户名"></el-table-column>
      <el-table-column prop="turthName" label="真实姓名"></el-table-column>
      <el-table-column prop="createTime" label="创建时间"></el-table-column>
      <el-table-column prop="lastLoginTime" label="上次登陆时间"></el-table-column>
      <el-table-column label="操作">
        <template #default="scoped">
          <el-button @click="handleEdit(scoped.row)">修改</el-button>
          <el-button @click="handleDel(scoped.row)" type="danger">删除</el-button>
          <el-button @click="handleSetRole(scoped.row)" type="info">设置权限</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="editUservisible" title="编辑用户">
      <el-form label-width="100px">
        <el-form-item label="userName">
          <el-input v-model="editUserObj.userName"></el-input>
        </el-form-item>
        <el-form-item label="turthName">
          <el-input v-model="editUserObj.turthName"></el-input>
        </el-form-item>
        <el-form-item label="password">
          <el-input v-model="editUserObj.password" :disabled="true"></el-input>
        </el-form-item>
        <el-form-item label="createTime">
          <el-input v-model="editUserObj.createTime" :disabled="true"></el-input>
        </el-form-item>
        <el-form-item label="lastLoginTime">
          <el-input v-model="editUserObj.lastLoginTime" :disabled="true"></el-input>
        </el-form-item>
        <el-form-item label="Operate">
          <el-button @click="reset">Reset</el-button>
          <el-button type="primary" @click="submitEdit">Submit</el-button></el-form-item
        >
      </el-form>
    </el-dialog>
    <el-dialog v-model="editSetRoleVisible" title="用户管理">
      <el-form label-width="100px">
        <el-form-item label="角色选择">
          <el-checkbox-group v-model="relList" :min="1">
            <el-checkbox
              style="display: block"
              v-for="item in roleList"
              :key="item.id"
              :value="item.id"
              :label="item.id"
            >
              {{ item.name }}
            </el-checkbox>
          </el-checkbox-group>
        </el-form-item>
        <el-form-item label="操作">
          <el-button type="primary" @click="submitSetRole"
            >Submit</el-button
          ></el-form-item
        >
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import {
  getRole,
  getUser,
  getUserToRoleRelation,
  saveUser,
  setUserToRole,
} from "@/api/system";
import { isEmptyObj } from "@/utils/obj";
import { ElNotification } from "element-plus";
import _ from "lodash";
import { onMounted, ref } from "vue";

type UserVo = {
  id: number;
  password: string;
  userName: string;
  createTime: string;
  lastLoginTime: string;
  turthName: string;
};
const parmas = ref({
  offset: 0,
  size: 10,
  keyword: "",
});
const userList = ref();
async function getServantList(init: boolean) {
  if (init) {
    parmas.value.offset = 0;
    parmas.value.size = 10;
  }
  const servantsResp = await getUser(parmas.value);
  userList.value = servantsResp.data;
  console.log("servantResp", servantsResp);
}

const editUservisible = ref(false);

const editUserObj = ref<UserVo>({
  password: "",
  userName: "",
  createTime: "",
  lastLoginTime: "",
  id: 0,
  turthName: "",
});
function handleEdit(row: UserVo) {
  editUserObj.value = _.cloneDeep(row);
  editUservisible.value = true;
}
function createUser() {
  editUservisible.value = true;
  reset();
}
async function submitEdit() {
  const submitBody = _.omit(editUserObj.value, ["createTime", "lastLoginTime"]);
  const data = await saveUser(submitBody);
  if (data.code) {
    return ElNotification.error(data.message);
  }
  editUservisible.value = false;
  getServantList(true);
  return ElNotification.success("success");
}
const editSetRoleVisible = ref(false);
const relList = ref([]);
const roleList = ref([]);
const userId = ref(0);
async function handleSetRole(row: UserVo) {
  editSetRoleVisible.value = true;
  userId.value = row.id;
  const rels = await getUserToRoleRelation(row.id);
  const roles = await getRole(undefined);
  roleList.value = roles.data;
  relList.value = rels.data.filter((v: any) => !isEmptyObj(v)).map((v: any) => v.id);
}
async function submitSetRole() {
  const body = {
    userId: userId.value,
    roleIds: relList.value.filter((v) => v),
  };
  const ret = setUserToRole(body);
  if (ret.code) {
    return ElNotification.error(ret.message);
  }
  editSetRoleVisible.value = false;
  return ElNotification.success("success");
}

function reset() {
  editUserObj.value.createTime = "";
  editUserObj.value.userName = "";
  editUserObj.value.password = "";
  editUserObj.value.lastLoginTime = "";
  editUserObj.value.turthName = "";
  editUserObj.value.id = 0;
}
onMounted(async () => {
  await getServantList(true);
});

function handleDel(row: UserVo) {
  console.log("row", row);
}
</script>

<style scoped></style>
