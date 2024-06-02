import _ from 'lodash'
import moment from 'moment'
import { ref } from 'vue'

type t_servant_type = {
  name: string
  value: number
}

export function useDashboardServantPie() {
  const dashboard_servant_type_option = ref({
    title: {
        text: "服务类型",
        textStyle: {
            color: '#999999',
            fontSize: 12,
        }
    },
    tooltip: {
      trigger: 'item',
      formatter: '{a} <br/>{b}: {c} ({d}%)'
    },
    legend: {
      top: '30%',
      left: '65%',
      orient: 'vertical',
      icon: 'circle',
      itemHeight: '10',
      textStyle: {
        fontSize: '13'
      }
    },
    series: [
      {
        name: '服务类型',
        type: 'pie',
        radius: '55%',
        center: ['30%', '50%'],
        data: [],
        avoidLabelOverlap: false,
        label: {
          show: false
        },
        labelLine: {
          show: false
        },
        itemStyle: {
          emphasis: {
            shadowBlur: 10,
            shadowOffsetX: 0,
            shadowColor: 'rgba(0, 0, 0, 0.5)'
          },
          normal: {
            color: function (colors) {
              const colorList = ['#975fe4', '#41a0ff', '#3bcbcb', '#4dcb73', '#f9d337', '#f2637b']
              return colorList[colors.dataIndex]
            }
          }
        }
      }
    ]
  })

  function set_dashboard_servant_type_option(data: Array<t_servant_type>) {
    dashboard_servant_type_option.value.series[0].data = data
    dashboard_servant_type_option.value.legend.data = data.map((v) => v.name)
  }
  return [dashboard_servant_type_option, set_dashboard_servant_type_option]
}

const __servant_type = {
    "id": 2,
    "serverName": "ShellServer",
    "language": "node",
    "upStreamName": "up_shell_server",
    "location": "/shellserver/",
    "protocol": "grpc",
    "execPath": "service_go",
    "servantGroupId": 1,
    "createTime":""
}

type servant_type = typeof __servant_type

export function useDashboardServantBar(){
    function setOpt(servers:servant_type[]){
        option.value.xAxis.data = Array.from(new Set(servers.map(v=>moment(v.createTime).format("YYYY-MM-DD"))))
        option.value.series[0].data = servers.reduce((pre,curr)=>{
            const item = pre.find((v) => v.name == moment(curr.createTime).format("YYYY-MM-DD"));
            if (item) {
              item.value += 1;
            } else {
              pre.push({
                name: moment(curr.createTime).format("YYYY-MM-DD"),
                value: 1,
              });
            }
            return pre;
        },[])
    }

    const option = ref({
        xAxis: {
            type: 'category',
            name: "日期",
            data:[],
            boundaryGap: [0, 0.01],
        },
        yAxis: {
            type: 'value'
        },
        series: [
            {
                data: [],
                type: "bar",
                showBackground: true,
                barWidth: '40%',
                color: "#09a3f6",
                name: "创建服务数",
                label: {
                    show: true,
                    position: 'top',
                },
            }
        ],
        tooltip: {
            trigger: "axis",
            formatter: "{b0}<br/>{a0} : {c0}",
            axisPointer: {
              type: "cross",
              crossStyle: {
                color: "#999",
              },
            },
        },
        grid: {
            bottom: '10%'
        },
        title: {
            text: "服务流水",
            textStyle: {
                color: '#999999',
                fontSize: 12,
            }
        },
    })

    return [option,setOpt]
}