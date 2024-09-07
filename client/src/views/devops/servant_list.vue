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
      <el-table-column prop="preview" label="preview"></el-table-column>
      <el-table-column label="操作" align="center">
        <template #default="scoped">
          <el-button type="primary" @click="setCronTask(scoped.row)">
            定时设置
          </el-button>
        </template>
      </el-table-column>

      <el-table-column label="操作" align="center">
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
      <el-form-item label="服务名">
        <el-input v-model="servant.serverName"></el-input>
      </el-form-item>
      <el-form-item label="语言">
        <el-select v-model="servant.language">
          <el-option
            v-for="item in languages"
            :label="item"
            :key="item"
            :value="item"
          ></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="协议">
        <el-select v-model="servant.protocol">
          <el-option
            v-for="item in protocols"
            :label="item"
            :key="item"
            :value="item"
          ></el-option>
        </el-select>
      </el-form-item>
      <el-form-item label="执行路径">
        <el-input v-model="servant.execPath"></el-input>
      </el-form-item>
      <el-form-item label="预览地址">
        <el-input v-model="servant.preview"></el-input>
      </el-form-item>
      <el-form-item label="服务组">
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

  <el-dialog v-model="cronTaskDialogVisible" title="定时设置" width="80%">
    <el-form :model="form" label-width="auto">
      <el-form-item label="Server Grids">
        <el-table :data="cronTaskNodes" border @selection-change="handleSelectionChange">
          <el-table-column type="selection" width="55" />
          <el-table-column label="Grid">
            <template #default="scoped">
              <el-button type="text" @click="toLog(scoped.row)"
                >{{ scoped.row.gridNode.ip }}:{{ scoped.row.port }}</el-button
              >
            </template>
          </el-table-column>
          <el-table-column label="Status">
            <template #default="scoped">
              <div
                :class="gridStatus[scoped.row.status] || 'offline'"
                @click="$emit('checkStatus')"
              >
                {{ gridStatus[scoped.row.status] || "offline" }}
              </div>
            </template>
          </el-table-column>
          <el-table-column prop="pid" label="PID"></el-table-column>
          <el-table-column label="Type">
            <template #default="scoped">
              <div>{{ scoped.row.gridServant.language }}</div>
            </template>
          </el-table-column>
          <el-table-column label="Protocol">
            <template #default="scoped">
              <div>{{ scoped.row.gridServant.protocol }}</div>
            </template>
          </el-table-column>
        </el-table>
      </el-form-item>
      <el-form-item label="Server Name">
        <el-input v-model="selectRow.serverName" disabled />
      </el-form-item>
      <el-form-item label="Activity Time">
        <el-date-picker
          v-model="form.activeTime"
          type="datetime"
          placeholder="上线时间"
        />
      </el-form-item>
      <el-form-item label="ExpireTime Time">
        <el-date-picker
          v-model="form.expireTime"
          type="datetime"
          placeholder="下线时间"
        />
      </el-form-item>
      <el-form-item label="Set Mark">
        <el-input v-model="form.mark" type="textarea" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">设置</el-button>
        <el-button>取消</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script setup lang="ts">
import api from "@/api/server";
import assetsApi from "@/api/assets";
import { useUserStore } from "@/stores/counter";
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
import _ from "lodash";
import moment from "moment";
import { computed, onMounted, reactive, ref } from "vue";
import { useRouter } from "vue-router";
const languages = ["node", "java", "java(jar)", "go", "exe", "custom command"];
const protocols = ["http", "grpc"];

const userStore = useUserStore();
const servantDemo = {
  id: 2,
  serverName: "ShellServer",
  language: "node",
  protocol: "grpc",
  execPath: "service_go",
  servantGroupId: 1,
  preview: "",
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

function confirmUpdateServant() {
  ElMessageBox.confirm(`确认修改?`, {
    confirmButtonText: "确认",
    cancelButtonText: "取消",
    type: "warning",
  })
    .then(async () => {
      const resp = await api.updateServant(servant.value);
      if (resp.code) {
        return ElMessage.error({
          type: "error",
          message: resp.message,
        });
      }
      await getServantList();
      ElMessage({
        type: "success",
        message: "修改成功",
      });
      editDialogVisible.value = false;
    })
    .catch(() => {
      ElMessage({
        type: "info",
        message: "取消修改",
      });
    });
}

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

const cronTaskDialogVisible = ref(false);
const cronTaskNodes = ref([]);
const selectRow = ref<T_Servant>({});
function setCronTask(row: T_Servant) {
  console.log("row", row);
  cronTaskDialogVisible.value = true;
  api.queryGrid({ id: row.id }).then((res) => {
    cronTaskNodes.value = res.data;
    selectRow.value = row;
  });
}

const router = useRouter();
function toLog(row: any) {
  const text = router.resolve({
    path: "/logpage",
    query: {
      host: row.gridNode.ip,
      serverName: selectRow.value.serverName,
      gridId: row.id,
    },
  });
  window.open(text.href, "_blank");
}

const selectionGrid = ref([]);
function handleSelectionChange(value: any) {
  selectionGrid.value = value;
}

const gridStatus: any = {
  "1": "online",
  "0": "offline",
};

const form = reactive({
  mark: "",
  activeTime: "",
  expireTime: "",
});

const submitBody = computed(() => {
  return {
    mark: form.mark,
    activeTime: form.activeTime,
    expireTime: form.expireTime,
    gridIds: selectionGrid.value.map((item) => item.id),
    operateUserId: userStore.userInfo.id,
    servantId: selectRow.value.id,
  };
});

async function onSubmit() {
  const now = moment();
  if (!submitBody.value.activeTime && !submitBody.value.expireTime) {
    return ElMessage.error({
      type: "error",
      message: "请至少选择一个上线时间或过期时间",
    });
  }
  if (submitBody.value.activeTime && now.isAfter(submitBody.value.activeTime)) {
    return ElMessage.error({
      type: "error",
      message: "上线时间不能小于当前时间",
    });
  }
  if (submitBody.value.expireTime && now.isAfter(submitBody.value.expireTime)) {
    return ElMessage.error({
      type: "error",
      message: "过期时间不能小于当前时间",
    });
  }
  if (submitBody.value.gridIds.length == 0) {
    return ElMessage.error({
      type: "error",
      message: "请至少选择一个节点",
    });
  }

  const bodys = submitBody.value.gridIds.map((v) => {
    return {
      gridId: v,
      mark: submitBody.value.mark,
      activeTime: moment(submitBody.value.activeTime).toISOString(),
      expireTime: moment(submitBody.value.expireTime).toISOString(),
      operateUserId: submitBody.value.operateUserId,
      ServantId: submitBody.value.servantId,
    };
  });
  await Promise.all(
    bodys.map(async (body) => {
      await assetsApi.upsertAsset(body);
    })
  );
  console.log("bodys", bodys);
  ElNotification.success("设置成功");
  cronTaskDialogVisible.value = false;
}
</script>

<style scoped></style>
