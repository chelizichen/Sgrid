<script lang="ts">
// devops component
// 1.选择服务组 ｜ 创建
// 2.创建服务
// 3.选择节点
// 4.添加至服务网格
// 5.选择端口
export default {
  name: "devops-component",
};
</script>

<template>
  <div>
    <el-container>
      <el-aside style="width: 150px">
        <el-menu
          class="el-menu-vertical-demo"
          active-text-color="rgb(207, 15, 124)"
          style="border: none"
          :default-openeds="['2']"
        >
          <el-menu-item index="1" key="1" @click="switchShow">
            <template #title>添加服务组</template>
          </el-menu-item>
          <el-menu-item index="2" key="2" @click="switchShow">
            <template #title>添加服务</template>
          </el-menu-item>
          <el-menu-item index="3" key="3" @click="switchShow">
            <template #title>选择节点</template>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <div v-if="modelIndex === '1'">
          <el-form :model="formData" :rules="rules" label-width="100px">
            <el-form-item label="服务标签" prop="tagName">
              <el-input v-model="formData.tagName"></el-input>
            </el-form-item>
            <el-form-item label="英文标签" prop="tagEnglishName">
              <el-input v-model="formData.tagEnglishName"></el-input>
            </el-form-item>
            <el-form-item label="操作">
              <el-button @click="resetForm">重置</el-button>
              <el-button type="primary" @click="devopsAddGroup"
                >确定</el-button
              ></el-form-item
            >
          </el-form>
        </div>
        <div v-if="modelIndex === '2'">
          <el-form :model="formData" :rules="rules" label-width="100px">
            <el-form-item label="选择项">
              <div>{{ selectOpt() }}</div>
              <!-- <el-input :disabled="true" v-model="selectOpt"></el-input> -->
            </el-form-item>
            <el-form-item label="选择服务组">
              <el-select v-model="groupId">
                <el-option
                  v-for="item in groups"
                  :label="item.tagName"
                  :key="item.id"
                  :value="Number(item.id)"
                ></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="服务名称">
              <el-input v-model="servantItem.serverName"></el-input>
            </el-form-item>
            <el-form-item label="语言">
              <el-input v-model="servantItem.language"></el-input>
            </el-form-item>
            <el-form-item label="协议">
              <el-input v-model="servantItem.protocol"></el-input>
            </el-form-item>
            <el-form-item label="可执行路径">
              <el-input v-model="servantItem.execPath"></el-input>
            </el-form-item>
            <el-form-item label="操作">
              <el-button @click="resetServant">重置</el-button>
              <el-button type="primary" @click="devopsAddServant"
                >确定</el-button
              ></el-form-item
            >
            <!-- <el-form-item label="端口">
              <el-input v-model="servantItem.Protocol"></el-input>
            </el-form-item> -->
          </el-form>
        </div>
        <div v-if="modelIndex === '3'">
          <el-card style="margin-bottom: 10px">
            <div style="display: flex; justify-content: space-around">
              <div style="color: rgb(207, 15, 124); cursor: pointer" @click="addGrid">
                服务部署
              </div>
              <div style="color: rgb(207, 15, 124); cursor: pointer" @click="addNode">
                添加节点
              </div>
            </div>
          </el-card>

          <el-table :data="nodes" border @selection-change="handleSelectionChange">
            <el-table-column type="selection" width="55" />
            <el-table-column prop="id" label="id"></el-table-column>
            <el-table-column prop="ip" label="ip"></el-table-column>
            <el-table-column prop="main" label="main"></el-table-column>
            <el-table-column prop="nodeStatus" label="nodeStatus"></el-table-column>
            <el-table-column prop="platform" label="platform"></el-table-column>
          </el-table>

          <el-dialog v-model="addGridVisible" title="服务部署">
            <el-form label-width="100px">
              <template v-for="(item, index) in selectionNodes" :key="index">
                <el-form-item label="分配节点">
                  <el-input v-model.number="addGridForm.port[index]">
                    <template #prepend>{{ item.ip }}</template>
                  </el-input>
                </el-form-item>
              </template>
              <el-form-item label="选择服务">
                <el-select v-model="addGridForm.servantId">
                  <el-option
                    v-for="item in servants"
                    :label="item.serverName"
                    :key="item.id"
                    :value="Number(item.id)"
                  ></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="操作">
                <el-button @click="resetServant">重置</el-button>
                <el-button type="primary" @click="devopsAddGrid"
                  >确定</el-button
                ></el-form-item
              >
            </el-form>
          </el-dialog>

          <el-dialog v-model="addNodeVisible" title="节点部署">
            <el-form label-width="100px">
              <el-form-item label="主机地址">
                <el-input v-model="addNodeForm.ip"></el-input>
              </el-form-item>
              <el-form-item label="操作系统">
                <el-input v-model="addNodeForm.platForm"></el-input>
              </el-form-item>
              <el-form-item label="主备">
                <el-select v-model="addNodeForm.main">
                  <el-option label="主" value="1"></el-option>
                  <el-option label="备" value="0"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="节点状态">
                <el-select v-model="addNodeForm.nodeStatus">
                  <el-option label="启用" :value="1"></el-option>
                  <el-option label="停止" :value="2"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="操作">
                <el-button @click="resetNode">重置</el-button>
                <el-button type="primary" @click="devopsAddNode"
                  >确定</el-button
                ></el-form-item
              >
            </el-form>
          </el-dialog>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import API from "@/api/server";
import { ElMessage } from "element-plus";
import { ref, watch } from "vue";

const selectOpt = () => {
  return `服务组ID :  ${groupId.value}
  | 服务名称： ${servantItem.value.serverName}
  | 语言：${servantItem.value.language}
  | 协议：${servantItem.value.protocol}
  | 可执行路径(golang 默认 sgrid_app) : ${servantItem.value.execPath}`;
};
const modelIndex = ref("");

// 服务组model
const formData = ref({
  tagName: "",
  tagEnglishName: "",
});
const rules = ref({
  tagName: [{ required: true, message: "请输入服务标签", trigger: "blur" }],
  tagEnglishName: [{ required: true, message: "请输入英文标签", trigger: "blur" }],
});

async function devopsAddGroup() {
  const data = await API.saveGroup(formData.value);
  if (data.code) {
    return ElMessage.error(data.message);
  }
  resetForm();
  return ElMessage.success("创建成功");
}

const resetForm = () => {
  formData.value.tagName = "";
  formData.value.tagEnglishName = "";
};

const resetServant = () => {
  servantItem.value.serverName = "";
  servantItem.value.language = "";
  servantItem.value.protocol = "";
  servantItem.value.execPath = "";
};

const resetNode = () => {
  addNodeForm.value.ip = "";
  addNodeForm.value.main = "";
  addNodeForm.value.nodeStatus = "";
  addNodeForm.value.platForm = "";
};

function switchShow(value: any) {
  modelIndex.value = value.index;
}

const groups = ref([]);
const groupId = ref(1);

const servantItem = ref({
  serverName: "",
  language: "",
  protocol: "",
  execPath: "sgrid_app",
});

async function devopsAddServant() {
  const body = {
    serverName: servantItem.value.serverName,
    language: servantItem.value.language,
    protocol: servantItem.value.protocol,
    execPath: servantItem.value.execPath,
    servantGroupId: groupId.value,
  };

  const data = await API.saveServant(body);
  if (data.code) {
    return ElMessage.error(data.message);
  }
  resetServant();
  return ElMessage.success("创建成功");
}

const nodes = ref([]);
const servants = ref([]);
const selectionNodes = ref([]);
const addGridVisible = ref(false);
const addGridForm = ref({
  port: [],
  servantId: 1,
  selectionNodes: [],
});
function addGrid() {
  addGridVisible.value = true;
}

function handleSelectionChange(value: never[]) {
  selectionNodes.value = value;
  addGridForm.value.selectionNodes = value;
}

async function devopsAddGrid() {
  const body = addGridForm.value.selectionNodes.map((item, index) => {
    return {
      nodeId: item.id,
      port: addGridForm.value.port[index],
      servantId: addGridForm.value.servantId,
    };
  });
  const ret = await Promise.all(body.map((v) => API.saveGrid(v)));
  if (ret.every((item) => item.code == 0)) {
    ElMessage.success("部署成功");
    addGridVisible.value = false;
  }
}

const addNodeVisible = ref(false);
const addNodeForm = ref({ ip: "", platForm: "", main: "0", nodeStatus: 0 });
function addNode() {
  addNodeVisible.value = true;
}

async function devopsAddNode() {
  const data = await API.saveNode(addNodeForm.value);
  if (data.code) {
    return ElMessage.error(data.message);
  }
  resetNode();
  ElMessage.success("部署成功");
  return (addNodeVisible.value = false);
}

watch(
  () => modelIndex.value,
  async function (newVal) {
    if (!newVal) {
      return;
    }
    if (newVal == "1") {
      resetForm();
    }
    if (newVal == "2") {
      const data = await API.getGroup();
      groups.value = data.data;
    }
    if (newVal == "3") {
      const data = await API.queryNodes();
      nodes.value = data.data;
      const servantsResp = await API.getServants();
      servants.value = servantsResp.data.sort((a, b) => b.id - a.id);
      addGridForm.value.servantId = servants.value[0].id;
    }
  }
);
</script>

<style scoped></style>
