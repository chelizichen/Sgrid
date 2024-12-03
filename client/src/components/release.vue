<script lang="ts">
export default {
  name: "release-component",
};
</script>
<template>
  <el-dialog
    :model-value="$attrs.releaseVisible"
    :title="$attrs.serverName"
    @close="$emit('CLOSE_RELEASE_DIALOG')"
  >
    <div v-for="item in $attrs.selectionGrid" :key="item.id" style="margin: 3px">
      {{ item.gridNode.ip }}
    </div>
    <el-select
      v-model="selectValue"
      placeholder="Select"
      style="width: 100%; margin: 20px 0 0 0"
    >
      <el-option
        v-for="item in $attrs.releaseList"
        :key="item.id"
        :label="item.filePath"
        :value="item.id"
        :disabled="item.status == -1"
      >
        <span style="float: left">{{ item.filePath }}</span>
        <span
          style="
            float: right;
            color: red;
            font-size: 13px;
            display: block;
            margin-left: 20px;
          "
          @click="handleRemovePackage(item)"
          >删除</span
        >
        <span style="float: right; color: var(--el-text-color-secondary); font-size: 13px"
          >ID: {{ item.id }} @ Tag : {{ item.version || "--" }}</span
        >
      </el-option>
    </el-select>

    <div style="margin: 10px 0">
      <el-input v-model="selectVersion">
        <template #prepend>
          <span>Tag</span>
        </template>
      </el-input>
    </div>

    <div style="margin: 10px 0">
      <el-input
        type="textarea"
        v-model="selectContent"
        :disabled="true"
        :rows="5"
      ></el-input>
    </div>

    <span slot="footer">
      <div style="display: flex; align-items: center; justify-content: center">
        <el-button type="info" @click="addTag()">Add Tag</el-button>
        <el-button type="primary" @click="$emit('CLOSE_RELEASE_DIALOG')">Close</el-button>
        <el-button type="danger" @click="$emit('RELEASE_SERVER_BY_ID', selectValue)"
          >Release</el-button
        >
      </div>
    </span>
  </el-dialog>
</template>

<script lang="ts" setup>
import api from "@/api/server";
import { ElMessage, ElMessageBox, ElNotification } from "element-plus";
import { defineEmits, ref, useAttrs, watch } from "vue";
import _ from "lodash";

const attrs = useAttrs();
const selectValue = ref();
const selectContent = ref("");
const selectVersion = ref("");
watch(selectValue, function (newVal) {
  console.log("newVal", newVal);
  console.log("attrs", attrs);
  selectContent.value = attrs.releaseList.find((v) => v.id == newVal).content;
  selectVersion.value = attrs.releaseList.find((v) => v.id == newVal).version;
});
defineEmits(["CLOSE_RELEASE_DIALOG", "RELEASE_SERVER_BY_ID"]);

async function addTag() {
  const id = Number(selectValue.value);
  const version = selectVersion.value;
  const body = {
    id,
    version,
  };
  await api.updatePackageVersion(body);
  ElNotification.success(" Update Tag Success");
}

async function handleRemovePackage(item) {
  console.log("item", item);
  if (item.status == "-1") {
    return false;
  }
  item.serverName = attrs.serverName;
  const body = _.pick(item, ["id", "serverName"]);
  ElMessageBox.prompt("确认删除该服务包？删除后不可恢复!", "确认", {
    confirmButtonText: "确认",
    cancelButtonText: "取消",
    inputPlaceholder: "输入密码",
  }).then(async ({ value: password }) => {
    console.log("password", password);
    if (password != item.serverName) {
      return ElMessage.error(`密码错误`);
    }
    console.log("body", body);
    const data = await api.delPackage(body);
    if (data.code) {
      return ElMessage.error(`删除失败:${data.message}`);
    }
    return ElMessage.success("删除成功");
  });
}
</script>
