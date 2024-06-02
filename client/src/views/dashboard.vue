<script lang="ts">
export default {
  name: "dashboard",
};
</script>

<template>
  <div style="padding: 20px">
    <el-card shadow="always" style="margin: 20px 0">
      <div style="display: flex">
        <div style="flex: 1">
          <div>服务总数</div>
          <div>{{ statisticsObj.nums }}项</div>
        </div>
        <div style="flex: 1">
          <div>服务类型</div>
          <div>{{ statisticsObj.types }}种</div>
        </div>
        <div style="flex: 1">
          <div>最新创建服务</div>
          <div>{{ statisticsObj.lastCreateServer }}</div>
        </div>
        <div style="flex: 1">
          <div>内存使用率</div>
          <div>{{ Number(memObj.usedPercent).toFixed(2) }}%</div>
        </div>
      </div>
    </el-card>
    <el-card shadow="always">
      <template #header>服务统计</template>
      <div style="display: flex">
        <div id="server_bar" style="width: 70%; height: 500px"></div>
        <div id="server_pie" style="width: 30%; height: 500px"></div>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import api from "@/api/server";
import { onMounted, ref } from "vue";
import * as echarts from "echarts";
import { useDashboardServantBar, useDashboardServantPie } from "@/charts/dashboard";
import { getCpuInfo, getMemoryInfo } from "@/api/system";

const [chart, setChart] = useDashboardServantPie();
const [barChart, setBarChart] = useDashboardServantBar();
const statisticsObj = ref({
  nums: 0,
  types: 0,
  lastCreateServer: "",
});

async function initSeravntAbout() {
  const data = await api.getServants();
  statisticsObj.value.nums = data.data.length;
  statisticsObj.value.lastCreateServer = data.data[data.data.length - 1].serverName;
  if (true) {
    const seriesData = data.data.reduce((pre, curr) => {
      const item = pre.find((v) => v.name == curr.language);
      if (item) {
        item.value += 1;
      } else {
        pre.push({
          name: curr.language,
          value: 1,
        });
      }
      return pre;
    }, []);
    statisticsObj.value.types = seriesData.length;

    setChart(seriesData);
    //   setChart()
    echarts.init(document.getElementById("server_pie")).setOption(chart.value);
  }
  if (true) {
    setBarChart(data.data);
    console.log("data.data", barChart.value);
    echarts.init(document.getElementById("server_bar")).setOption(barChart.value);
  }
}

const memObj = ref({
  usedPercent: "",
});
async function statisticsGet() {
  const data = await getCpuInfo();
  console.log("data", data);

  const mem = await getMemoryInfo();
  memObj.value = JSON.parse(mem.data);
  console.log("mem", mem);
}
onMounted(() => {
  initSeravntAbout();
  statisticsGet();
});
</script>

<style scoped></style>
