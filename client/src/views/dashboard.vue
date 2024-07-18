<script lang="ts">
export default {
  name: "dashboard-component",
};
</script>

<template>
  <div class="bg">
    <dv-decoration8 class="left-deco"> </dv-decoration8>
    <dv-decoration8 :reverse="true" class="right-deco"> </dv-decoration8>
    <dv-decoration-10 class="left-line" />
    <dv-decoration-10 :reverse="true" class="right-line" />
    <!-- 中间文字 -->
    <dv-decoration-11 class="center-text">
      <div>云智慧监控系统</div>
    </dv-decoration-11>
    <!-- 右边雷达图 -->
    <dv-decoration-12 class="right-radar" />
    <!-- 内存使用率 -->
    <dv-water-level-pond :config="configSystemMem" class="pond" />

    <!-- 服务包 -->
    <dv-border-box12 class="box1">
      <dv-scroll-board :config="configStatisticsGetServerPackage" class="table1" />
    </dv-border-box12>
    <!-- 实时日志 -->
    <dv-border-box12 class="box2">
      <dv-scroll-board :config="configStatisticsGetLatestLog" class="table2" />
    </dv-border-box12>
    <!-- 服务类型 -->
    <dv-border-box12 class="box3">
      <dv-scroll-ranking-board :config="configStatisticsGetServerType" class="table3" />
    </dv-border-box12>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from "vue";
import { getStatisticsByType } from "@/api/statistics";
import { getCpuInfo, getMemoryInfo } from "@/api/system";
enum Req_TYPE {
  StatisticsGetServerPackage = "1",
  StatisticsGetLatestLog = "2",
  StatisticsGetStatus = "3",
  StatisticsGetServerType = "4",
  StatisticsGetServants = "5",
}
const configStatisticsGetServerPackage = reactive({
  header: ["序号", "服务名", "发布总次数"],
  data: [],
  columnWidth: [50, 250, 100],
});

const configStatisticsGetLatestLog = reactive({
  header: ["序号", "服务名", "发布总次数"],
  data: [],
  columnWidth: [100, 200],
});
const configStatisticsGetServerType = reactive({
  data: [],
  unit: "个",
});
async function init() {
  const data1 = await getStatisticsByType(Req_TYPE.StatisticsGetServerPackage);
  configStatisticsGetServerPackage.data = data1.data.map((v) => {
    return [v.id, v.label, v.value + " (次)"];
  });
  const data2 = await getStatisticsByType(Req_TYPE.StatisticsGetLatestLog);
  configStatisticsGetLatestLog.data = data2.data.map((v) => {
    return [v.id, v.label, v.value];
  });
  const data3 = await getStatisticsByType(Req_TYPE.StatisticsGetStatus);

  const data4 = await getStatisticsByType(Req_TYPE.StatisticsGetServerType);
  configStatisticsGetServerType.data = data4.data.map((v) => {
    return {
      name: v.label,
      value: v.value,
    };
  });
  const data5 = await getStatisticsByType(Req_TYPE.StatisticsGetServants);
}

const configSystemMem = reactive({
  data: [],
});
async function systemStatistics() {
  const mem = await getMemoryInfo();
  configSystemMem.data = [JSON.parse(mem.data).usedPercent.toFixed(2)];
  console.log("configSystemMem.data", configSystemMem.data);
}
onMounted(() => {
  init();
  systemStatistics();
});
</script>

<style scoped lang="less">
.bg {
  padding: 20px;
  min-height: 100vh;
  background: #222222;
  position: relative;
  width: 100vw;
  height: 100vh;
  box-sizing: border-box;
  .left-deco {
    position: absolute;
    top: 10vh;
    left: 25%;
    width: 300px;
    height: 50px;
  }
  .right-deco {
    position: absolute;
    top: 10vh;
    right: 25%;
    width: 300px;
    height: 50px;
  }
  .left-line {
    position: absolute;
    width: 20%;
    height: 5px;
    top: 10vh;
    left: 5%;
  }
  .right-line {
    position: absolute;
    width: 20%;
    height: 5px;
    top: 10vh;
    right: 5%;
  }
  .center-text {
    position: absolute;
    z-index: 100%;
    width: 30vw;
    height: 60px;
    color: #9494fc;
    left: 35%;
    top: 6vh;
  }
  .right-radar {
    position: absolute;
    width: 60px;
    height: 60px;
    right: 10vw;
    top: 3vh;
  }
  .box1 {
    position: absolute;
    top: 20vh;
    width: 30vw;
    height: 30vh;
    box-sizing: border-box;
    .table1 {
      width: 94%;
      height: 92%;
      position: absolute;
      top: 4%;
      left: 3%;
    }
  }
  .box2 {
    position: absolute;
    top: 50vh;
    width: 100vw;
    height: 40vh;
    box-sizing: border-box;
    .table2 {
      width: 98%;
      height: 92%;
      position: absolute;
      top: 4%;
      left: 1%;
    }
  }
  .box3 {
    position: absolute;
    top: 20vh;
    width: 30vw;
    height: 30vh;
    right: 0;
    box-sizing: border-box;
    .table3 {
      width: 96%;
      height: 96%;
      position: absolute;
      top: 2%;
      left: 2%;
    }
  }
  .pond {
    position: absolute;
    width: 30vw;
    height: 10vh;
    top: 40vh;
    left: 35vw;
  }
}
</style>
