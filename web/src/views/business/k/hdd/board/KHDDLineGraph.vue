<template>
  <div ref="chartRef" class="h-full w-full"></div>
</template>

<script setup lang="ts">
  import { getHddSummaryDailyPeak } from '@/api/business/k';
  import { useECharts } from '@/hooks/web/useECharts';
  import { ref, Ref, onMounted, nextTick, reactive } from 'vue';

  defineOptions({
    name: 'KHDDLineGraph',
  });
  const props = defineProps({
    reportTimeBegin: {
      type: String as PropType<string>,
      required: true,
    },
    reportTimeEnd: {
      type: String as PropType<string>,
      required: true,
    },
    minimumSize: {
      type: Number as PropType<number>,
      required: false,
    },
    maximumSize: {
      type: Number as PropType<number>,
      required: false,
    },
    minNetworkSpeed: {
      type: Number as PropType<number>,
      required: false,
    },
    maxNetworkSpeed: {
      type: Number as PropType<number>,
      required: false,
    },
    processStatus: {
      type: String as PropType<string>,
      required: false,
    },
    businessMountStatus: {
      type: String as PropType<string>,
      required: false,
    },
  });
  interface PeakData {
    dailyPeakTime: string;
    date: string;
    value: number;
  }

  const data = reactive({
    avgData: 4230,
    peakData: [] as PeakData[],
    xAxisData: [] as string[],
    // peakDataMounted: [] as PeakData[],
    // avgDataMounted: 0,
  });

  const chartRef = ref<any>(null);
  const { setOptions, getInstance } = useECharts(chartRef as Ref<HTMLDivElement>);

  const getOption = function () {
    return {
      title: {
        text: `K汇聚HDD盘在线统计看板`,
        subtext: `日峰平均值：${data.avgData}`,
        // subtext: `日峰平均值-全部分区：${data.avgData}; 日峰平均值-已挂载分区：${data.avgDataMounted}`,
      },
      tooltip: {
        trigger: 'axis',
        formatter: function (params: any) {
          let str = '日期：' + params[0].name + '<br/>';
          // 全部分区 - 蓝色 #1890ff
          let point = data.peakData.find((item) => item.date === params[0].name);
          if (point) {
            str +=
              '<span style="display:inline-block;margin-right:4px;border-radius:10px;width:10px;height:10px;background-color:#1890ff;"></span>' +
              '<span style="color:#1890ff;font-weight:500;">峰值信息：</span>' +
              '<br/>' +
              '&nbsp;&nbsp;&nbsp;日峰值：' +
              (point.value || 0) +
              '<br/>' +
              '&nbsp;&nbsp;&nbsp;峰值时刻：' +
              (point.dailyPeakTime || '-');
          }
          // 已挂载分区 - 绿色 #52c41a
          // let pointMounted = data.peakDataMounted.find((item) => item.date === params[0].name);
          // if (pointMounted) {
          //   str +=
          //     '<br/>' +
          //     '<span style="display:inline-block;margin-right:4px;border-radius:10px;width:10px;height:10px;background-color:#52c41a;"></span>' +
          //     '<span style="color:#52c41a;font-weight:500;">已挂载分区：</span>' +
          //     '<br/>' +
          //     '&nbsp;&nbsp;&nbsp;日峰值：' +
          //     (pointMounted.value || 0) +
          //     '<br/>' +
          //     '&nbsp;&nbsp;&nbsp;峰值时刻：' +
          //     (pointMounted.dailyPeakTime || '-');
          // }
          return str;
        },
      },
      legend: {
        selectedMode: true,
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true,
      },
      toolbox: {
        feature: {
          saveAsImage: {},
        },
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: data.xAxisData,
      },
      yAxis: {
        type: 'value',
        name: '最大值 非重复统计主机名（并集）',
        nameLocation: 'middle',
        nameGap: 60,
        nameTextStyle: {
          fontSize: 14,
          color: '#333',
        },
        max: getYAxisMax(),
        min: getYAxisMin(),
      },
      series: [
        {
          name: '日峰值',
          type: 'line',
          data: getPeakData(),
          smooth: true, // 平滑曲线
          symbol: 'circle', // 圆形数据点
          symbolSize: 8, // 数据点大小
          itemStyle: {
            color: '#1890ff',
            borderColor: '#fff',
            borderWidth: 2,
          },
          lineStyle: {
            color: '#1890ff',
            width: 2,
            shadowColor: 'rgba(24, 144, 255, 0.3)',
            shadowBlur: 4,
            shadowOffsetY: 2,
          },
          areaStyle: {
            color: {
              type: 'linear',
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                { offset: 0, color: 'rgba(24, 144, 255, 0.25)' },
                { offset: 1, color: 'rgba(24, 144, 255, 0.02)' },
              ],
            },
          },
          label: {
            show: true,
            position: 'top',
            distance: 8,
            formatter: '{c}',
            fontSize: 11,
            fontWeight: 500,
            color: '#1890ff',
            backgroundColor: 'rgba(255, 255, 255, 0.85)',
            padding: [2, 4],
            borderRadius: 2,
          },
        },
        // {
        //   name: '日峰值-已挂载分区',
        //   type: 'line',
        //   data: getPeakDataMounted(),
        //   smooth: true,
        //   symbol: 'circle',
        //   symbolSize: 8,
        //   itemStyle: {
        //     color: '#52c41a',
        //     borderColor: '#fff',
        //     borderWidth: 2,
        //   },
        //   lineStyle: {
        //     color: '#52c41a',
        //     width: 2,
        //     shadowColor: 'rgba(82, 196, 26, 0.3)',
        //     shadowBlur: 4,
        //     shadowOffsetY: 2,
        //   },
        //   areaStyle: {
        //     color: {
        //       type: 'linear',
        //       x: 0,
        //       y: 0,
        //       x2: 0,
        //       y2: 1,
        //       colorStops: [
        //         { offset: 0, color: 'rgba(82, 196, 26, 0.25)' },
        //         { offset: 1, color: 'rgba(82, 196, 26, 0.02)' },
        //       ],
        //     },
        //   },
        //   label: {
        //     show: true,
        //     position: 'bottom',
        //     distance: 8,
        //     formatter: '{c}',
        //     fontSize: 11,
        //     fontWeight: 500,
        //     color: '#52c41a',
        //     backgroundColor: 'rgba(255, 255, 255, 0.85)',
        //     padding: [2, 4],
        //     borderRadius: 2,
        //   },
        // },
        {
          name: '日峰平均值：' + data.avgData,
          type: 'line',
          data: getAvgData(),
          symbol: 'none',
          itemStyle: {
            color: '#69c0ff',
          },
          lineStyle: {
            color: '#69c0ff',
            width: 2,
            type: [8, 4], // 自定义虚线样式：线段8px，间隔4px
          },
          label: {
            show: true,
            position: [0, -15],
            formatter: `日峰平均值：${data.avgData}`,
            fontSize: 11,
            color: '#69c0ff',
            fontWeight: 500,
            backgroundColor: 'rgba(255, 255, 255, 0.85)',
            padding: [2, 4],
            borderRadius: 2,
          },
        },
        // {
        //   name: '日峰平均值-已挂载分区：' + data.avgDataMounted,
        //   type: 'line',
        //   data: getAvgDataMounted(),
        //   symbol: 'none',
        //   itemStyle: {
        //     color: '#95de64',
        //   },
        //   lineStyle: {
        //     color: '#95de64',
        //     width: 2,
        //     type: [8, 4], // 自定义虚线样式
        //   },
        //   label: {
        //     show: true,
        //     position: [0, 5], // 放到线条下方
        //     formatter: `日峰平均值-已挂载分区：${data.avgDataMounted}`,
        //     fontSize: 11,
        //     color: '#95de64',
        //     fontWeight: 500,
        //     backgroundColor: 'rgba(255, 255, 255, 0.85)',
        //     padding: [2, 4],
        //     borderRadius: 2,
        //   },
        // },
      ],
    };
  };
  function getYAxisMax() {
    // let max = Math.max(...getPeakData(), ...getPeakDataMounted());
    let max = Math.max(...getPeakData());
    let value = ((max - (max % 50)) / 50) * 50 + 150;
    return value;
  }
  function getYAxisMin() {
    // let min = Math.min(...getPeakData(), ...getPeakDataMounted());
    let min = Math.min(...getPeakData());
    let value = ((min - (min % 50)) / 50) * 50 - 150;
    return value <= 0 ? 0 : value;
  }

  function getAvgData() {
    const avgData = data.avgData;
    const avgDataArray: number[] = [];
    for (let i = 0; i < data.xAxisData.length; i++) {
      avgDataArray.push(avgData);
    }
    return avgDataArray;
  }
  // function getAvgDataMounted() {
  //   const avgDataMounted = data.avgDataMounted;
  //   const avgDataMountedArray: number[] = [];
  //   for (let i = 0; i < data.xAxisData.length; i++) {
  //     avgDataMountedArray.push(avgDataMounted);
  //   }
  //   return avgDataMountedArray;
  // }
  async function loadPeakData() {
    let res = await getHddSummaryDailyPeak({
      reportTimeBegin: props.reportTimeBegin,
      reportTimeEnd: props.reportTimeEnd,
      minimumSize: props.minimumSize,
      maximumSize: props.maximumSize,
      minNetworkSpeed: props.minNetworkSpeed,
      maxNetworkSpeed: props.maxNetworkSpeed,
      processStatus: props.processStatus,
      businessMountStatus: props.businessMountStatus,
    });
    let xAxisData: string[] = [];
    if (res && res.x?.length > 0) {
      xAxisData = res.x;
    }
    data.xAxisData = xAxisData;

    // 获取日峰值数据
    let peakDatas: PeakData[] = [];
    if (res && res.data?.length > 0) {
      peakDatas = res.data;
    }
    data.peakData = peakDatas;
    data.avgData = res.avg || 0;

    // 获取已挂载的日峰值数据
    // res = await getHddSummaryDailyPeak({
    //   reportTimeBegin: props.reportTimeBegin,
    //   reportTimeEnd: props.reportTimeEnd,
    //   businessMountStatus: 1,
    // });

    // let peakDatasMounted: PeakData[] = [];
    // if (res && res.data?.length > 0) {
    //   peakDatasMounted = res.data;
    // }
    // data.peakDataMounted = peakDatasMounted;
    // data.avgDataMounted = res.avg || 0;
  }
  function getPeakData() {
    const peakData = data.peakData;
    const peakDataArray: number[] = [];
    for (let i = 0; i < data.xAxisData.length; i++) {
      peakDataArray.push(peakData[i].value);
    }
    return peakDataArray;
  }
  // function getPeakDataMounted() {
  //   const peakDataMounted = data.peakDataMounted;
  //   const peakDataMountedArray: number[] = [];
  //   for (let i = 0; i < data.xAxisData.length; i++) {
  //     peakDataMountedArray.push(peakDataMounted[i].value);
  //   }
  //   return peakDataMountedArray;
  // }
  async function rebuildGraphic() {
    await loadPeakData();
    setOptions(getOption() as any);
    // 强制重新计算大小
    nextTick(() => {
      getInstance()?.resize();
    });
  }

  // 监听容器大小变化
  onMounted(() => {
    // 监听窗口大小变化
    const handleResize = () => {
      getInstance()?.resize();
    };

    window.addEventListener('resize', handleResize);

    // 使用ResizeObserver监听容器大小变化
    if (chartRef.value) {
      const resizeObserver = new ResizeObserver(() => {
        getInstance()?.resize();
      });
      resizeObserver.observe(chartRef.value);
    }
  });

  defineExpose({ rebuildGraphic });
</script>

<style scoped lang="scss">
  // 确保容器占满父元素
  div {
    min-height: 300px; // 设置最小高度
  }
</style>
