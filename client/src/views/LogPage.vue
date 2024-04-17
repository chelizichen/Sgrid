<template>
  <div></div>
</template>

<script setup lang="ts">
import API from "@/api/server";
import { computed, onMounted } from "vue";
import { useRoute } from "vue-router";
const route = useRoute();
const query = computed(() => route.query);
async function init() {
  console.log("query", query);

  const data = await API.getLogFileList({
    host: query.value.host,
    serverName: query.value.serverName,
  });
  console.log("data", data);
}
onMounted(() => {
  init();
});
</script>

<style scoped></style>
