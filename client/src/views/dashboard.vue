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
      <div>智慧监控系统</div>
    </dv-decoration-11>
    <!-- 右边雷达图 -->
    <dv-decoration-12 class="right-radar" />

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
      <div id="chart2-inner" style="width: 100%; height: 100%"></div>
    </dv-border-box12>

    <dv-border-box12 class="box4">
      <dv-scroll-board :config="configStatisticsGetStatus" class="table4" />
    </dv-border-box12>

    <dv-border-box12 class="box5">
      <dv-scroll-board :config="configStatisticsGetServants" class="table5" />
    </dv-border-box12>
    <dv-border-box12 class="chart1-out">
      <div id="chart1-inner" style="width: 100%; height: 100%"></div>
    </dv-border-box12>

    <dv-border-box12 class="chart3-out">
      <div id="chart3-inner" style="width: 100%; height: 100%"></div>
    </dv-border-box12>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref, reactive } from "vue";
import { getStatisticsByType, getNodesInfo } from "@/api/statistics";
import * as echarts from "echarts";
import {
  useDashboardChart1,
  useDashboardChart2,
  useDashboardChart3,
} from "@/charts/dashboard";
import _ from "lodash";
import api from "@/api/server";
// 系统统计
const [opt1, setOpt1] = useDashboardChart1();
// 服务类型
const [opt2, setOpt2] = useDashboardChart2();
// 发布统计
const [opt3, setOpt3] = useDashboardChart3();

enum Req_TYPE {
  StatisticsGetServerPackage = "1",
  StatisticsGetLatestLog = "2",
  StatisticsGetStatus = "3",
  StatisticsGetServerType = "4",
  StatisticsGetServants = "5",
  StatisticsgetNodes = "6",
}
const configStatisticsGetServerPackage = reactive({
  header: ["序号", "服务名", "发布总次数"],
  data: [],
  columnWidth: [50, 200, 100],
});

const configStatisticsGetLatestLog = reactive({
  header: ["序号", "服务名", "日志"],
  data: [],
  columnWidth: [100, 200, 1200],
});
const configStatisticsGetStatus = reactive({
  header: ["序号", "服务(PID)", "状态"],
  data: [],
  columnWidth: [50, 150, 50],
});
const configStatisticsGetServerType = reactive({
  header: ["服务类型", "数量"],
  data: [],
  columnWidth: [150, 50],
});

const configStatisticsGetServants = reactive({
  header: ["序号", "服务名", "创建时间"],
  data: [],
  columnWidth: [70, 150, 250],
});

async function init() {
  const data1 = await getStatisticsByType(Req_TYPE.StatisticsGetServerPackage);
  configStatisticsGetServerPackage.data = data1.data
    .sort((a, b) => b.value - a.value)
    .map((v) => {
      return [v.id, v.label, v.value + " (次)"];
    });
  const data2 = await getStatisticsByType(Req_TYPE.StatisticsGetLatestLog);
  configStatisticsGetLatestLog.data = data2.data.map((v) => {
    return [v.id, v.label, v.value];
  });
  const data3 = await getStatisticsByType(Req_TYPE.StatisticsGetStatus);
  configStatisticsGetStatus.data = data3.data.map((v) => {
    return [v.id, v.label, v.value];
  });
  const data4 = await getStatisticsByType(Req_TYPE.StatisticsGetServerType);
  configStatisticsGetServerType.data = data4.data.map((v) => {
    return { name: v.label, value: v.value };
  });
  const data5 = await getStatisticsByType(Req_TYPE.StatisticsGetServants);
  configStatisticsGetServants.data = data5.data.map((v) => {
    return [v.id, v.label, v.value];
  });
}

const systemInfo = ref({});
const nodesInfo = ref([]);
async function systemStatistics() {
  const nodes = await getStatisticsByType(Req_TYPE.StatisticsgetNodes);
  nodesInfo.value = nodes.data;

  const infos = await getNodesInfo();
  systemInfo.value = infos.data;
  console.log("infos", infos);
}

async function totalSetChart() {
  {
    const systemInfosToMap = _.keyBy(systemInfo.value, "host");
    const chartvalue = nodesInfo.value
      .map((v) => {
        return Object.assign({}, v, systemInfosToMap[v.id]);
      })
      .map((v) => {
        return {
          name: v.label,
          value: [v.cpuLength, v.cpuPercent, v.memorySize, v.memoryPercent, v.value],
        };
      });
    console.log("chartvalue", chartvalue);

    setOpt1(chartvalue);
    echarts.init(document.getElementById("chart1-inner")).setOption(opt1.value);
  }
  {
    const chartvalue = configStatisticsGetServerType.data;
    setOpt2(chartvalue);
    echarts.init(document.getElementById("chart2-inner")).setOption(opt2.value);
  }
  {
    const data = await api.getServants();
    setOpt3(data.data);
    echarts.init(document.getElementById("chart3-inner")).setOption(opt3.value);
  }
}

onMounted(async () => {
  await init();
  await systemStatistics();
  totalSetChart();
});
</script>

<style scoped lang="less">
.bg {
  background-image: url(@/assets/big-screen-bg.png);
  background-size: 100% 100%;
  padding: 20px;
  min-height: 100vh;
  // background: #222222;
  position: relative;
  width: 100vw;
  height: 100vh;
  box-sizing: border-box;
  .chart1-out {
    position: absolute;
    top: 41vh;
    height: 36vh;
    width: 40vw;
    padding: 20px;
    box-sizing: border-box;
  }
  .chart3-out {
    position: absolute;
    top: 41vh;
    height: 36vh;
    width: 59vw;
    left: 40vw;
    padding: 20px;
    box-sizing: border-box;
  }
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
    top: 15vh;
    width: 30vw;
    height: 25vh;
    left: 22vw;
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
    top: 78vh;
    width: 98vw;
    height: 22vh;
    box-sizing: border-box;
    .table2 {
      width: 98%;
      height: 98%;
      position: absolute;
      top: 1%;
      left: 1%;
    }
  }
  .box3 {
    position: absolute;
    top: 15vh;
    width: 18vw;
    height: 25vh;
    right: 30.5vw;
    box-sizing: border-box;
    padding: 20px;
  }
  .box4 {
    position: absolute;
    top: 15vh;
    width: 21vw;
    height: 25vh;
    left: 1vw;
    box-sizing: border-box;
    .table4 {
      width: 96%;
      height: 92%;
      position: absolute;
      top: 4%;
      left: 2%;
    }
  }
  .box5 {
    position: absolute;
    top: 15vh;
    width: 30vw;
    height: 25vh;
    right: 1vw;
    box-sizing: border-box;
    .table5 {
      width: 98%;
      height: 92%;
      position: absolute;
      top: 4%;
      left: 1%;
    }
  }
}
</style>
