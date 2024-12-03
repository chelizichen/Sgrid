<template>
  <el-dialog
    :model-value="props.dialogVisible"
    title="配置中心"
    @close="$emit('CLOSE_RELEASE_DIALOG')"
  >
    <el-form :model="configForm" ref="configFormRef" label-width="120px">
      <el-form-item label="修改配置">
        <el-switch v-model="isAble" inline-prompt />
      </el-form-item>
      <el-form-item label="配置ID">
        <el-input v-model="configForm.id" disabled></el-input>
      </el-form-item>
      <el-form-item label="配置项">
        <el-input
          v-if="isAble"
          v-model="configForm.conf"
          placeholder="请输入配置项"
          type="textarea"
          rows="20"
        ></el-input>
        <el-card style="width: 100%" v-if="!isAble">
          <highlightjs language="yaml" :code="configForm.conf"></highlightjs>
        </el-card>
      </el-form-item>
      <el-form-item label="操作">
        <el-button type="primary" @click="submitConfig">保存</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script lang="ts" setup>
import { getConfig, updateConfig } from "@/api/servantConf";
import { ElNotification } from "element-plus";
import { ref, watch } from "vue";
const isAble = ref(false);
const configForm = ref({
  id: 0, // 假设我们有一个初始ID，或者你可以设置为null或undefined
  conf: "",
  servantId: 0,
});

const props = defineProps<{
  servantId: number;
  dialogVisible: boolean;
}>();

watch(
  () => props.dialogVisible,
  function (newVal) {
    if (!newVal) {
      return;
    }
    fetchConfig();
  }
);

const fetchConfig = async () => {
  try {
    console.log("props", props);

    const response = await getConfig({ id: props.servantId });
    if (response.data) {
      configForm.value = { ...configForm.value, ...response.data };
    }
  } catch (error) {
    console.error("Failed to fetch config:", error);
    // 处理错误
  }
};

const submitConfig = async () => {
  try {
    configForm.value.servantId = Number(props.servantId);
    await updateConfig(configForm.value);
    // 处理成功后的逻辑，比如提示用户或重新获取配置
    ElNotification.success("保存成功");
    // fetchConfig(); // 如果需要，可以重新获取配置
  } catch (error) {
    console.error("Failed to submit config:", error);
    // 处理错误
  }
};
</script>

<style scoped>
/* 你的样式代码 */
</style>
