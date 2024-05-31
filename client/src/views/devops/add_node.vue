<template>
  <div>
    <el-form :inline="true">
      <el-form-item>
        <el-button type="primary" @click="addGrid">配置服务节点</el-button>
      </el-form-item>
      <el-form-item>
        <el-button @click="addNode">添加节点</el-button>
      </el-form-item>
    </el-form>
    <el-table :data="nodes" border @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="id" label="id"></el-table-column>
      <el-table-column prop="ip" label="ip"></el-table-column>
      <el-table-column prop="main" label="main"></el-table-column>
      <el-table-column prop="nodeStatus" label="nodeStatus"></el-table-column>
      <el-table-column prop="platform" label="platform"></el-table-column>
    </el-table>

    <el-dialog v-model="addGridVisible" title="配置服务节点">
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
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import api from "@/api/server";
import { ElNotification } from "element-plus";

const nodes = ref([]);
const servants = ref([]);
const selectionNodes = ref([]);
const addGridVisible = ref(false);
const addGridForm = ref({
  port: [],
  servantId: "",
  selectionNodes: [],
});
function addGrid() {
  addGridVisible.value = true;
}
function addNode() {
  addNodeVisible.value = true;
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
  const ret = await Promise.all(body.map((v) => api.saveGrid(v)));
  if (ret.every((item) => item.code == 0)) {
    ElNotification.success("Release Success");
    addGridVisible.value = false;
  } else {
    const item = ret.find((item) => item.code != 0);
    ElNotification.error(item!.message);
  }
}

const resetNode = () => {
  addNodeForm.value.ip = "";
  addNodeForm.value.main = "";
  addNodeForm.value.nodeStatus = 0;
  addNodeForm.value.platForm = "";
};

const addNodeVisible = ref(false);
const addNodeForm = ref({ ip: "", platForm: "", main: "0", nodeStatus: 0 });

async function devopsAddNode() {
  const data = await api.saveNode(addNodeForm.value);
  if (data.code) {
    return ElNotification.error(data.message);
  }
  resetNode();
  ElNotification.success("Release Success");
  return (addNodeVisible.value = false);
}
const resetServant = () => {
  addGridForm.value.servantId = "";
};
onMounted(async () => {
  const data = await api.queryNodes();
  nodes.value = data.data;
  const servantsResp = await api.getServants();
  servants.value = servantsResp.data.sort((a, b) => b.id - a.id);
  addGridForm.value.servantId = servants.value[0].id;
});
</script>
