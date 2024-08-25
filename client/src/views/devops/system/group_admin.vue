<template>
  <div>
    <el-form :inline="true">
      <el-form-item>
        <el-button type="primary" @click="editUservisible = true">创建小组</el-button>
      </el-form-item>
    </el-form>

    <el-table
      :data="pageGroupList"
      border
      highlight-current-row
      @current-change="handleCurrentChange"
    >
      <el-table-column type="index" label="序号" width="90"></el-table-column>
      <el-table-column prop="name" label="小组名称">
        <template #default="scoped">
          <el-button type="text" @click="editUserGroupToServantGroup(scoped.row)">{{
            scoped.row.name
          }}</el-button>
        </template>
      </el-table-column>
      <el-table-column prop="total" label="小组总人数">
        <template #default="scoped">
          <el-button type="text" @click="editUsersByUserGroup(scoped.row)"
            >{{ scoped.row.total }}（人）</el-button
          >
        </template>
      </el-table-column>
      <el-table-column prop="description" label="小组描述"></el-table-column>
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

    <el-dialog v-model="editUservisible" title="编辑小组">
      <el-form label-width="100px">
        <el-form-item label="小组名称">
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
    <el-dialog v-model="editUserGroupToServantGroupVisible" title="编辑服务组映射关系">
      <el-form label-width="100px">
        <el-form-item label="选择服务组">
          <el-select v-model="bindServantGroupList" multiple>
            <el-option
              v-for="item in groupList"
              :value="item.id"
              :label="`${item.tagName}  (${item.tagEnglishName})`"
              :key="item.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="操作">
          <el-button @click="confirmSetUserGroupToServantGroup"
            >保存</el-button
          ></el-form-item
        >
      </el-form>
    </el-dialog>
    <el-dialog v-model="editUsersByUserGroupVisible" title="编辑团队成员映射关系">
      <el-form label-width="100px">
        <el-form-item label="选择团队成员">
          <el-select v-model="bindUserList" multiple>
            <el-option
              v-for="item in userList"
              :value="item.id"
              :label="`${item.turthName}  (${item.userName})`"
              :key="item.id"
            ></el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="操作">
          <el-button @click="confirmSetUserToUserGroup">保存</el-button></el-form-item
        >
      </el-form>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import api from "@/api/server";
import {
  delUserGroup,
  getGroup,
  saveUserGroup,
  getUsersByUserGroup,
  getServantGroupsByUserGroupId,
  setUserGroupToServantGroup,
  getUser,
  setUserToUserGroup,
} from "@/api/system";
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
import _, { template } from "lodash";
import { onMounted, ref } from "vue";

type RoleVo = {
  id: number;
  name: string;
  createTime: string;
  description: string;
};

const pagination = ref({
  offset: 0,
  size: 10,
  keyword: "",
});

const pageGroupList = ref();
async function getPageList() {
  const servantsResp = await getGroup(pagination.value);
  pageGroupList.value = servantsResp.data;
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
  const data = await saveUserGroup(editRoleObj.value);
  if (data.code) {
    return ElNotification.error(data.message);
  }
  editUservisible.value = false;
  getPageList();
  return ElNotification.success("success");
}

function reset() {
  editRoleObj.value.createTime = "";
  editRoleObj.value.id = 0;
  editRoleObj.value.name = "";
}
onMounted(async () => {
  await getPageList();
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
      await delUserGroup(row.id);
      await getPageList();
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

const editUsersByUserGroupVisible = ref(false);
const bindUserList = ref([]);
const userList = ref([]);
async function editUsersByUserGroup(row: any) {
  bindGroupId.value = row.id;
  const dto = {
    id: row.id,
    offset: 0,
    size: 9999,
  };
  const _bindUserList = await getUsersByUserGroup(dto);
  if (_bindUserList.data) {
    bindUserList.value = _bindUserList.data.map((v) => v.userId);
  }

  const userDto = {
    offset: 0,
    size: 9999,
  };
  const _userList = await getUser(userDto);
  userList.value = _userList.data;
  editUsersByUserGroupVisible.value = true;
}

async function confirmSetUserToUserGroup() {
  const body = {
    groupId: bindGroupId.value,
    userIds: bindUserList.value,
  };
  await setUserToUserGroup(body);
  editUsersByUserGroupVisible.value = false;
  getPageList();
  return ElNotification.success("success");
}

const editUserGroupToServantGroupVisible = ref(false);
const groupList = ref([]);
const bindGroupId = ref(0);
const bindServantGroupList = ref([]);
async function editUserGroupToServantGroup(row: any) {
  bindGroupId.value = row.id;
  const dto = {
    id: row.id,
    offset: 0,
    size: 9999,
  };
  const _groupList = await api.getGroup(0);
  groupList.value = _groupList.data;
  console.log("groupList", groupList.value);
  const _bindServantGroupList = await getServantGroupsByUserGroupId(dto);
  if (_bindServantGroupList.data) {
    bindServantGroupList.value = _bindServantGroupList.data.map((v) => v.servantGroupId);
    console.log("bindServantGroupList", bindServantGroupList.value);
  }
  editUserGroupToServantGroupVisible.value = true;
}
async function confirmSetUserGroupToServantGroup() {
  const body = {
    groupId: bindGroupId.value,
    servantGroupIds: bindServantGroupList.value,
  };
  await setUserGroupToServantGroup(body);
  editUserGroupToServantGroupVisible.value = false;
  getPageList();
  return ElNotification.success("success");
}
</script>

<style scoped></style>
