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
      style="width: 100%; margin: 20px 0"
    >
      <el-option
        v-for="item in $attrs.releaseList"
        :key="item.id"
        :label="item.filePath"
        :value="item.id"
      >
        <span style="float: left">{{ item.filePath }}</span>
        <span style="float: right; color: var(--el-text-color-secondary); font-size: 13px"
          >ID: {{ item.id }}</span
        >
      </el-option>
    </el-select>
    <span slot="footer">
      <div style="display: flex; align-items: center; justify-content: center">
        <el-button type="primary" @click="$emit('CLOSE_RELEASE_DIALOG')">Close</el-button>
        <el-button type="danger" @click="$emit('RELEASE_SERVER_BY_ID', selectValue)"
          >Release</el-button
        >
      </div>
    </span>
  </el-dialog>
</template>

<script lang="ts" setup>
import { defineEmits, ref } from "vue";

const selectValue = ref();
defineEmits(["CLOSE_RELEASE_DIALOG", "RELEASE_SERVER_BY_ID"]);
</script>
