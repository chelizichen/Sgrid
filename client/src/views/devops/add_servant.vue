<template>
  <div>
    <el-form label-width="100px">
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
        <el-select v-model="servantItem.language">
          <el-option
            v-for="item in languages"
            :label="item"
            :key="item"
            :value="item"
          ></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="Protocol">
        <el-select v-model="servantItem.protocol">
          <el-option
            v-for="item in protocols"
            :label="item"
            :key="item"
            :value="item"
          ></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="ExecPath">
        <el-input v-model="servantItem.execPath"></el-input>
      </el-form-item>
      <el-form-item label="Operate">
        <el-button @click="resetServant">Reset</el-button>
        <el-button type="primary" @click="devopsAddServant"
          >Submit</el-button
        ></el-form-item
      >
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from "vue";
import api from "@/api/server";
import { ElNotification } from "element-plus";
import { useUserStore } from "@/stores/counter";

const languages = ["node", "java", "go"];
const protocols = ["http", "grpc"];
const selectOpt = () => {
  return `Group:ID :  ${groupId.value}
  | ServantName： ${servantItem.value.serverName}
  | Language：${servantItem.value.language}
  | Protocol：${servantItem.value.protocol}
  | Exec Path (golang :: default ::sgrid_app) : ${servantItem.value.execPath}`;
};

const groupId = ref(0);
const groups = ref<Array<{ tagEnglishName: string; tagName: string; id: number }>>([]);
const servantItem = ref({
  serverName: "",
  language: "node",
  protocol: "http",
  execPath: "sgrid_app",
});
async function devopsAddServant() {
  const body = {
    serverName: servantItem.value.serverName,
    language: servantItem.value.language,
    protocol: servantItem.value.protocol,
    execPath: servantItem.value.execPath,
    servantGroupId: groupId.value,
    userId: userStore.userInfo.id,
  };

  const data = await api.saveServant(body);
  if (data.code) {
    return ElNotification.error(data.message);
  }
  resetServant();
  return ElNotification.success("Create Success");
}

const userStore = useUserStore();

onMounted(async () => {
  const data = await api.getGroup(userStore.userInfo.id);
  groups.value = data.data;
  groupId.value = data.data[data.data.length - 1].id;
});
const resetServant = () => {
  servantItem.value.serverName = "";
  servantItem.value.language = "";
  servantItem.value.protocol = "";
  servantItem.value.execPath = "";
};
</script>
