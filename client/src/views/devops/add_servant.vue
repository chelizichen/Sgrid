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
      <el-form-item label="ExecPath">
        <el-input v-model="servantItem.execPath"></el-input>
      </el-form-item>
      <el-form-item label="Operate">
        <el-button @click="resetServant">Reset</el-button>
        <el-button type="primary" @click="devopsAddServant">Submit</el-button>
        <el-button type="danger" @click="tips">Tips</el-button>
      </el-form-item>
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from "vue";
import api from "@/api/server";
import { ElMessageBox, ElNotification } from "element-plus";
import { useUserStore } from "@/stores/counter";

const languages = ["node", "java", "java(jar)", "go", "exe", "custom command", "python(tar)", "python(exe)"];
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

function tips() {
  ElMessageBox.confirm(
    `
        <div>
          Protocl 选项的选择暂时与业务方无任何太大的联系，选择哪一个都可以通过命令启动服务
        </div>

        <br/>
        <div>
          Language 选项一开始只是为了方便进行跨语言执行不同的命令，在随后 <b>20240907</b>
          日的更新中弱化了这个概念，添加了
          <b> custom command</b> 这个选项
        </div>
        <br/>
        <div>
          在 <b>Language</b> 选择
          <b>custom command</b>
          选项后，代表着会通过自定义命令启动服务，此时命令需要开发人员自行确认 ExecPath
          的值，值为 字符串数组
          <br />
          <code><b>样例 : ["node","server.js"]</b></code>
          <br />
          <br />
          我在golang后台自行做了反序列化和切片,保证命令是正确的
          <br />
          <code>
            <b> var parseExecArgs []string </b>
            <br />
            <b> err = json.Unmarshal([]byte(execFilePath), &parseExecArgs) </b>
            <br />
            <b> cmd = exec.Command(parseExecArgs[0], parseExecArgs[1:]...) </b>
          </code>
          <br />
        </div>
      `,
    "tips",
    {
      dangerouslyUseHTMLString: true,
      customStyle: {
        width: "800px",
      },
    }
  );
}
</script>
