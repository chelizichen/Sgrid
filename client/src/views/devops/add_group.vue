<template>
  <div>
    <el-form :model="formData" :rules="rules" label-width="100px">
      <el-form-item label="TagName" prop="tagName">
        <el-input v-model="formData.tagName"></el-input>
      </el-form-item>
      <el-form-item label="TagEnglishName" prop="tagEnglishName">
        <el-input v-model="formData.tagEnglishName"></el-input>
      </el-form-item>
      <el-form-item label="Operate">
        <el-button @click="resetForm">Reset</el-button>
        <el-button type="primary" @click="devopsAddGroup">Submit</el-button></el-form-item
      >
    </el-form>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import api from "@/api/server";
import { ElNotification } from "element-plus";
import { useUserStore } from "@/stores/counter";

const formData = ref({
  tagName: "",
  tagEnglishName: "",
  userId: 0,
});
const rules = ref({
  tagName: [{ required: true, message: "Please Input TagName", trigger: "blur" }],
  tagEnglishName: [
    { required: true, message: "Please Input TagEnglishName", trigger: "blur" },
  ],
});
const userStore = useUserStore();
async function devopsAddGroup() {
  formData.value.userId = userStore.userInfo.id;
  const data = await api.saveGroup(formData.value);
  if (data.code) {
    return ElNotification.error(data.message);
  }
  return ElNotification.success("Create Success");
}
</script>
