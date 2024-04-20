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
      <el-aside>
        <el-menu
          class="el-menu-vertical-demo"
          active-text-color="rgb(207, 15, 124)"
          style="border: none"
          :default-openeds="['2']"
        >
          <el-menu-item index="0" key="0" @click="switchShow('showAddGroup')">
            <el-icon>
              <TrendCharts />
            </el-icon>
            <template #title>{{ "Add Group" }}</template>
          </el-menu-item>
          <el-sub-menu index="1">
            <template #title>
              <el-icon><Setting /></el-icon>
              <span>Add Server</span>
            </template>
            <el-menu-item index="1" key="1" @click="switchShow('showSwitchGroup')">
              <el-icon>
                <TrendCharts />
              </el-icon>
              <template #title>{{ "SwitchGroup" }}</template>
            </el-menu-item>
          </el-sub-menu>
        </el-menu>
      </el-aside>
      <el-main>
        <div v-if="state.showAddGroup">
          <el-form
            ref="form"
            :model="formData"
            :rules="rules"
            label-width="100px"
            class="add-servant-group-form"
          >
            <el-form-item label="服务标签" prop="tagName">
              <el-input v-model="formData.tagName"></el-input>
            </el-form-item>
            <el-form-item label="英文标签" prop="tagEnglishName">
              <el-input v-model="formData.tagEnglishName"></el-input>
            </el-form-item>
            <el-form-item label="确定">
              <el-button @click="resetForm">重置</el-button>
              <el-button type="primary">确定</el-button></el-form-item
            >
          </el-form>
        </div>

        <div v-if="state.showSwitchGroup"></div>
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import API from "@/api/server";
import { ref, watch } from "vue";

const state = ref<Record<string, boolean>>({
  showAddGroup: false,
  showSwitchGroup: false,
});
const resetState = () => {
  state.value.showAddGroup = false;
  state.value.showSwitchGroup = false;
};

// 服务组model
const formData = ref({
  tagName: "",
  tagEnglishName: "",
});
const rules = ref({
  tagName: [{ required: true, message: "请输入服务标签", trigger: "blur" }],
  tagEnglishName: [{ required: true, message: "请输入英文标签", trigger: "blur" }],
});
const resetForm = () => {
  formData.value.tagName = "";
  formData.value.tagEnglishName = "";
};

function switchShow(key: string) {
  resetState();
  state.value[key] = true;
}

watch(
  () => state.value.showSwitchGroup,
  async function (newVal) {
    if (!newVal) {
      return;
    }
  }
);
</script>

<style scoped></style>
