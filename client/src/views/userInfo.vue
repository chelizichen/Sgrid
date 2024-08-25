<template>
  <div></div>
</template>

<script setup lang="ts">
import { onMounted, onBeforeUnmount } from "vue";
onMounted(() => {
  L2Dwidget.init({
    model: {
      jsonPath:
        "https://unpkg.com/live2d-widget-model-izumi@1.0.0/assets/izumi.model.json",
      scale: 2,
    },
    display: {
      superSample: 2,
      width: 300, // 宽度
      height: 400, // 高度
      //   position: "left", // 模型位置
      hOffset: 600, // 纵向偏移
      vOffset: 0, // 横向偏移
    },
    mobile: {
      show: true,
      scale: 0.4,
    },
    react: {
      opacityDefault: 0.7,
      opacityOnHover: 0.2,
    },
  });
});
onBeforeUnmount(() => {
  const l2d_dom = document.getElementById("live2d-widget");
  console.log("l2d_dom", l2d_dom);
  l2d_dom?.remove();
});
</script>

<style scoped></style>
