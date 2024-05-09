<template>
  <el-dialog
    :model-value="props.dialogVisible"
    title="配置中心"
    @close="$emit('CLOSE_RELEASE_DIALOG')"
  >
    <el-form :model="configForm" ref="configFormRef" label-width="120px">
      <el-form-item label="配置ID">
        <el-input v-model="configForm.id" disabled></el-input>
      </el-form-item>
      <el-form-item label="配置项">
        <el-input
          v-model="configForm.conf"
          placeholder="请输入配置项"
          type="textarea"
          rows="20"
        ></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="submitConfig">保存</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>

<script lang="ts" setup>
import { getConfig, updateConfig } from "@/api/servantConf";
import { ElNotification } from "element-plus";
import { ref, watch } from "vue";

const configForm = ref({
  id: 0, // 假设我们有一个初始ID，或者你可以设置为null或undefined
  conf: "",
  servantId: 0,
});

const props = defineProps<{
  servantId: string;
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
