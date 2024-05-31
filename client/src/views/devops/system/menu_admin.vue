<template>
  <div class="menu-to-role">
    <div class="item">
      <el-card>
        <el-form :inline="true">
          <el-form-item>
            <el-button type="primary" @click="createParentMenu()">创建父级菜单</el-button>
          </el-form-item>
          <el-form-item>
            <el-button @click="submitRoleMenu()">保存</el-button>
          </el-form-item>
        </el-form>
        <el-tree
          :data="dataSource"
          show-checkbox
          node-key="id"
          default-expand-all
          :expand-on-click-node="false"
          ref="treeMenuRef"
        >
          <template #default="{ node, data }">
            <div class="custom-tree-node">
              <div>{{ node.label }}</div>
              <div>
                <el-button @click="appendMenu(data)" type="text" v-if="node.level == 1">
                  添加
                </el-button>
                <el-button @click="changeMenu(data)" type="text">修改</el-button>
                <el-button
                  style="margin-left: 8px"
                  type="text"
                  @click="removeMenu(node, data)"
                  v-if="node.level == 2 || !data.children || data.children.length == 0"
                >
                  删除
                </el-button>
              </div>
            </div>
          </template>
        </el-tree>
      </el-card>

      <el-dialog v-model="editMenuVisible" title="角色管理">
        <el-form label-width="100px">
          <el-form-item label="父级ID">
            <el-input v-model="editMenuObj.parentId" :disabled="true"></el-input>
          </el-form-item>
          <el-form-item label="菜单标题">
            <el-input v-model="editMenuObj.title"></el-input>
          </el-form-item>
          <el-form-item label="组建路径">
            <el-input v-model="editMenuObj.component"></el-input>
          </el-form-item>
          <el-form-item label="名称">
            <el-input v-model="editMenuObj.name"></el-input>
          </el-form-item>
          <el-form-item label="路由">
            <el-input v-model="editMenuObj.path"></el-input>
          </el-form-item>
          <el-form-item label="操作">
            <el-button @click="reset">Reset</el-button>
            <el-button type="primary" @click="submitEdit">Submit</el-button></el-form-item
          >
        </el-form>
      </el-dialog>
    </div>
    <div class="item">
      <el-card>
        <RoleAdmin @recvRole="handleRecvRole"></RoleAdmin>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { getMenu, getMenuListByRoleId, saveMenu, setRoleToMenu } from "@/api/system";
import { ElNotification } from "element-plus";
import _ from "lodash";
import { computed, onMounted, ref } from "vue";
import RoleAdmin from "./role_admin.vue";
const treeMenuRef = ref();
interface Tree {
  id: number;
  label: string;
  children?: Tree[];
}
type MenuVo = {
  id: number;
  name: string;
  title: string;
  path: string;
  component: string;
  parentId: number;
};

const menuList = ref<Array<MenuVo>>([]);
async function getMenuList() {
  const servantsResp = await getMenu(undefined);
  menuList.value = servantsResp.data;
  console.log("servantResp", servantsResp);
}

const editMenuVisible = ref(false);

const editMenuObj = ref<Partial<MenuVo>>({
  title: "",
  path: "",
  name: "",
  component: "",
  parentId: 0,
});
async function submitEdit() {
  const data = await saveMenu(editMenuObj.value);
  if (data.code) {
    return ElNotification.error(data.message);
  }
  editMenuVisible.value = false;
  getMenuList();
  return ElNotification.success("success");
}

function reset() {
  editMenuObj.value.id = 0;
  editMenuObj.value.name = "";
  editMenuObj.value.title = "";
  editMenuObj.value.component = "";
  editMenuObj.value.path = "";
  editMenuObj.value.parentId = 0;
}

function reduceMenu(list: Array<MenuVo>): Tree[] {
  const menus: Tree[] = [];
  // 拿到根节点
  list.forEach((e) => {
    if (e.parentId == 0 || !e.parentId) {
      menus.push({
        label: e.title,
        id: e.id,
      });
    }
  });
  // 最多支持双层
  list.forEach((e) => {
    if (e.parentId || e.parentId != 0) {
      const item = menus.find((v) => v.id == e.parentId);
      if (!item) {
        return;
      }
      if (!item.children) {
        item.children = [];
      }
      item.children.push({
        label: e.title,
        id: e.id,
      });
    }
  });
  return menus;
}

const dataSource = computed(() => {
  const item = reduceMenu(menuList.value);
  console.log("item", item);
  return item;
});

onMounted(async () => {
  await getMenuList();
});
function createParentMenu() {
  reset();
  editMenuVisible.value = true;
}
async function appendMenu(data: Tree) {
  reset();
  editMenuObj.value.parentId = data.id;
  editMenuVisible.value = true;
}
async function changeMenu(data: Tree) {
  const item = menuList.value.find((v) => v.id == data.id);
  editMenuObj.value = _.cloneDeep(item);
  editMenuVisible.value = true;
}

async function removeMenu(node, data) {}
const recvRoleId = ref<number>(0);
async function handleRecvRole(id: number) {
  console.log("id", id);
  recvRoleId.value = id;
  const data = await getMenuListByRoleId(id);
  if (data.data instanceof Array) {
    treeMenuRef.value.setCheckedKeys(data.data.map((v) => v.menuId));
    console.log("data", data);
  }
}
async function submitRoleMenu() {
  const ids = treeMenuRef.value.getCheckedKeys();
  const body = {
    roleId: recvRoleId.value,
    menuIds: ids,
  };
  const data = await setRoleToMenu(body);
  if (data.code) {
    return ElNotification.error(data.message);
  }
  ElNotification.success("success");
  console.log("ref", treeMenuRef.value.getCheckedKeys());
}
</script>

<style scoped>
.custom-tree-node {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: 14px;
  padding-right: 8px;
}
.menu-to-role {
  display: flex;
  justify-content: space-between;
}
.item {
  width: 49%;
}
</style>
