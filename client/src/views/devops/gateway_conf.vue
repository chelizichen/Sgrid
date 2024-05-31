<template>
  <div>
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
</template>

<script setup lang="ts">
import { ref, watch } from "vue";
import { ElNotification } from "element-plus";
import { getBackupList, getBackupFile, merge, test, reload } from "@/api/nginx";

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
    return ElNotification.error("Merge Error|" + data.message);
  }
  ElNotification.success("Merge Success");
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
    return ElNotification.error("Test Error|" + data.message);
  }
  ElNotification.success("Test Success");
  expansionForm.value.chooseFileContent = data.data;
}

async function nginxReload() {
  const data = await reload();
  if (data.code) {
    return ElNotification.error("Reload Error|" + data.message);
  }
  ElNotification.success("Reload Success");
  expansionForm.value.chooseFileContent = data.data;
}
</script>
