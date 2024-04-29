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
            <template #title>AddGroup</template>
          </el-menu-item>
          <el-menu-item index="2" key="2" @click="switchShow">
            <template #title>AddServant</template>
          </el-menu-item>
          <el-menu-item index="3" key="3" @click="switchShow">
            <template #title>AddNode</template>
          </el-menu-item>
          <el-menu-item index="4" key="4" @click="switchShow">
            <template #title>GatewayConf</template>
          </el-menu-item>
          <el-menu-item index="5" key="5" @click="switchShow">
            <template #title>ServantAdmin</template>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <div v-if="modelIndex === '1'">
          <el-form :model="formData" :rules="rules" label-width="100px">
            <el-form-item label="TagName" prop="tagName">
              <el-input v-model="formData.tagName"></el-input>
            </el-form-item>
            <el-form-item label="TagEnglishName" prop="tagEnglishName">
              <el-input v-model="formData.tagEnglishName"></el-input>
            </el-form-item>
            <el-form-item label="Operate">
              <el-button @click="resetForm">Reset</el-button>
              <el-button type="primary" @click="devopsAddGroup"
                >Submit</el-button
              ></el-form-item
            >
          </el-form>
        </div>
        <div v-if="modelIndex === '2'">
          <el-form :model="formData" :rules="rules" label-width="100px">
            <el-form-item label="Options">
              <div>{{ selectOpt() }}</div>
              <!-- <el-input :disabled="true" v-model="selectOpt"></el-input> -->
            </el-form-item>
            <el-form-item label="SelectGroup">
              <el-select v-model="groupId">
                <el-option
                  v-for="item in groups"
                  :label="item.tagEnglishName"
                  :key="item.id"
                  :value="Number(item.id)"
                ></el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="ServerName">
              <el-input v-model="servantItem.serverName"></el-input>
            </el-form-item>
            <el-form-item label="Language">
              <el-input v-model="servantItem.language"></el-input>
            </el-form-item>
            <el-form-item label="Protocol">
              <el-input v-model="servantItem.protocol"></el-input>
            </el-form-item>
            <el-form-item label="可执行路径">
              <el-input v-model="servantItem.execPath"></el-input>
            </el-form-item>
            <el-form-item label="Operate">
              <el-button @click="resetServant">Reset</el-button>
              <el-button type="primary" @click="devopsAddServant"
                >Submit</el-button
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
                Server Release
              </div>
              <div style="color: rgb(207, 15, 124); cursor: pointer" @click="addNode">
                Add Node
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

          <el-dialog v-model="addGridVisible" title="Server Release">
            <el-form label-width="100px">
              <template v-for="(item, index) in selectionNodes" :key="index">
                <el-form-item label="ChooseNode">
                  <el-input v-model.number="addGridForm.port[index]">
                    <template #prepend>{{ item.ip }}</template>
                  </el-input>
                </el-form-item>
              </template>
              <el-form-item label="ChooseServant">
                <el-select v-model="addGridForm.servantId">
                  <el-option
                    v-for="item in servants"
                    :label="item.serverName"
                    :key="item.id"
                    :value="Number(item.id)"
                  ></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="Operate">
                <el-button @click="resetServant">Reset</el-button>
                <el-button type="primary" @click="devopsAddGrid"
                  >Confirm</el-button
                ></el-form-item
              >
            </el-form>
          </el-dialog>

          <el-dialog v-model="addNodeVisible" title="AddNode">
            <el-form label-width="100px">
              <el-form-item label="Host">
                <el-input v-model="addNodeForm.ip"></el-input>
              </el-form-item>
              <el-form-item label="Os">
                <el-input v-model="addNodeForm.platForm"></el-input>
              </el-form-item>
              <el-form-item label="IsMaster">
                <el-select v-model="addNodeForm.main">
                  <el-option label="Master" value="1"></el-option>
                  <el-option label="Slave" value="0"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="NodeStat">
                <el-select v-model="addNodeForm.nodeStatus">
                  <el-option label="Use" :value="1"></el-option>
                  <el-option label="Stop" :value="2"></el-option>
                </el-select>
              </el-form-item>
              <el-form-item label="Operate">
                <el-button @click="resetNode">Reset</el-button>
                <el-button type="primary" @click="devopsAddNode"
                  >Confirm</el-button
                ></el-form-item
              >
            </el-form>
          </el-dialog>
        </div>
        <div v-if="modelIndex === '4'">
          <div style="display: flex">
            <el-card style="width: 20%">
              <el-form label-position="left" label-width="100px">
                <el-form-item label="File Manager">
                  <el-select v-model="expansionForm.chooseFile" @change="selectFile">
                    <el-option
                      v-for="item in expansionForm.list"
                      :label="item"
                      :value="item"
                      :key="item"
                    ></el-option>
                  </el-select>
                </el-form-item>
                <el-form-item label="Operate">
                  <el-button type="primary" @click="mergeContent">Merge</el-button>
                </el-form-item>
                <el-form-item label="Operate">
                  <el-button type="primary" @click="nginxTest">Test</el-button>
                </el-form-item>
                <el-form-item label="Operate">
                  <el-button type="primary" @click="nginxReload">Reload</el-button>
                </el-form-item>
              </el-form>
            </el-card>
            <el-card style="width: 80%">
              <template #header>
                <div class="card-header">
                  <el-switch
                    style="margin-left: 20px"
                    v-model="expansionForm.couldEdit"
                    inline-prompt
                  />
                </div>
              </template>
              <el-input
                type="textarea"
                rows="40"
                v-model="expansionForm.chooseFileContent"
                :disabled="!expansionForm.couldEdit"
              ></el-input>
            </el-card>
          </div>
        </div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import API from "@/api/server";
import { getBackupList, getBackupFile, merge, test, reload } from "@/api/nginx";
import { ElMessage } from "element-plus";
import { ref, watch } from "vue";

const selectOpt = () => {
  return `Group:ID :  ${groupId.value}
  | ServantName： ${servantItem.value.serverName}
  | Language：${servantItem.value.language}
  | Protocol：${servantItem.value.protocol}
  | Exec Path (golang :: default ::sgrid_app) : ${servantItem.value.execPath}`;
};
const modelIndex = ref("");

// 服务组model
const formData = ref({
  tagName: "",
  tagEnglishName: "",
});
const rules = ref({
  tagName: [{ required: true, message: "Please Input TagName", trigger: "blur" }],
  tagEnglishName: [
    { required: true, message: "Please Input TagEnglishName", trigger: "blur" },
  ],
});

async function devopsAddGroup() {
  const data = await API.saveGroup(formData.value);
  if (data.code) {
    return ElMessage.error(data.message);
  }
  resetForm();
  return ElMessage.success("Create Success");
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
  addNodeForm.value.nodeStatus = 0;
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
  return ElMessage.success("Create Success");
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
    ElMessage.success("Release Success");
    addGridVisible.value = false;
  } else {
    const item = ret.find((item) => item.code != 0);
    ElMessage.error(item!.message);
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
  ElMessage.success("Release Success");
  return (addNodeVisible.value = false);
}

const expansionForm = ref({
  list: [],
  chooseFile: "",
  chooseFileContent: "",
  couldEdit: false,
});

async function selectFile(f: string) {
  const file = await getBackupFile({
    fileName: f,
  });
  expansionForm.value.chooseFileContent = file.data;
}

async function mergeContent() {
  const data = await merge({ content: expansionForm.value.chooseFileContent });
  if (data.code) {
    return ElMessage.error("Merge Error|" + data.message);
  }
  ElMessage.success("Merge Success");
  expansionForm.value.couldEdit = false;
  console.log("data", data);
  if (!data.code) {
    const list = await getBackupList();
    list.data.unshift("origin");
    expansionForm.value.list = list.data;
  }
}

async function nginxTest() {
  const data = await test();
  if (data.code) {
    return ElMessage.error("Test Error|" + data.message);
  }
  ElMessage.success("Test Success");
  expansionForm.value.chooseFileContent = data.data;
}

async function nginxReload() {
  const data = await reload();
  if (data.code) {
    return ElMessage.error("Reload Error|" + data.message);
  }
  ElMessage.success("Reload Success");
  expansionForm.value.chooseFileContent = data.data;
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
    if (newVal == "4") {
      const data = await getBackupList();
      data.data.unshift("origin");
      expansionForm.value.list = data.data;
    }
  }
);
</script>

<style scoped></style>
