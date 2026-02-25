<template>
  <Card :body-style="{ padding: '10px' }">
    <div ref="chartRef" :style="{ width, height }"></div>
  </Card>
</template>
<script setup lang="ts">
  import { onMounted, Ref, ref } from 'vue';
  // import { defineOptions, defineEmits, defineProps } from 'vue';
  import { Card } from 'ant-design-vue';
  import { useECharts } from '@/hooks/web/useECharts';

  const emits = defineEmits(['onFieldSelect']);

  defineOptions({ name: 'AnomalyServerPie' });
  const props = defineProps({
    taskId: {
      type: [Number, String],
      required: true,
    },
    data: {
      type: Array as PropType<any[]>,
      default: () => [],
    },
    title: {
      type: [Number, String],
      required: true,
    },
    number: {
      type: [Number, String],
      required: true,
    },

    width: {
      type: String as PropType<string>,
      default: '100%',
    },
    height: {
      type: String as PropType<string>,
      default: '150px',
    },
    legendLeft: {
      type: String as PropType<string>,
      default: '20%',
    },
    seriesLeft: {
      type: String as PropType<string>,
      default: '10%',
    },
  });
  const chartRef = ref<any>(null);
  const { setOptions, getInstance } = useECharts(chartRef as Ref<HTMLDivElement>);
  const getOption = function () {
    return {
      color: [
        '#003366',
        '#006699',
        '#4cabce',
        '#e5323e',
        '#FF9900',
        '#109618',
        '#990099',
        '#0099C6',
        '#DD4477',
        '#66AA00',
        '#B82E2E',
        '#316395',
        '#F6B26B',
        '#93C47D',
        '#8E7CC3',
        '#FFD966',
        '#E06666',
        '#6FA8DC',
        '#9FC5E8',
        '#A64D79',
        '#5B8E7D',
        '#BC8034',
        '#6A5ACD',
        '#20B2AA',
        '#FF6347',
        '#4682B4',
        '#32CD32',
        '#DA70D6',
        '#FF4500',
      ],
      tooltip: {
        // trigger: 'item',
        show: false,
      },
      legend: {
        // height: '30%', // 限制图例高度
        // right: 20, // 右侧位置
        top: 'center', // 垂直居中
        left: props.legendLeft,
        itemGap: 10,
        itemWidth: 20, // 标记宽度
        itemHeight: 14, // 标记高度
        formatter: function (name) {
          if (props.data.length == 0) {
            return props.title;
          }
          const dataItem = props.data.find((item) => item.name === name);
          if (!dataItem) return name;
          // 设置最大显示长度（字符数）
          const maxLength = 24;
          let displayName = name;

          // 截断过长的名称
          if (name.length > maxLength) {
            displayName = name.substring(0, maxLength - 1) + '...';
          }

          return `${displayName} - ${dataItem.value}`;
        },
        orient: 'vertical',
        textStyle: {
          fontSize: 13,
          fontWeight: 'bold',
          color: '#333',
          itemStyle: {
            borderWidth: 0, // 去除边框
          },
          textStyle: {
            fontSize: 12,
            padding: [0, 0, 0, 5], // 文字左边距
          },
        },
        selectedMode: true,
      },
      series: [
        {
          legendHoverLink: true,
          center: [props.seriesLeft, '50%'], // 调整饼图位置为左侧
          name: props.title,
          type: 'pie',
          radius: ['70%', '90%'],
          avoidLabelOverlap: false,
          label: {
            silent: true,
            show: true,
            position: 'center',
            color: '#4c4a4a',
            formatter: '{active|' + props.title + '}' + '\n\r' + '{total|' + props.number + '}',
            rich: {
              active: {
                // fontFamily: '微软雅黑',
                fontSize: 12,
                color: '#6c7a89',
              },
              total: {
                fontSize: 30,
                // fontFamily: '微软雅黑',
                color: '#333',
                fontBold: true,
                lineHeight: 40,
              },
            },
          },
          emphasis: {
            disabled: false,
            label: {
              show: true,
            },
          },
          labelLine: {
            show: true,
          },
          itemStyle: {
            borderRadius: 1,
            borderColor: '#fff',
            borderWidth: 1,
          },
          data: props.data,
        },
      ],
      graphic: [],
    };
  };

  function rebuildGraphic() {
    setOptions(getOption() as any);
  }
  defineExpose({ rebuildGraphic });
  // watch(
  //   () => props.data,
  //   () => {
  //     setOptions(getOption() as any);
  //   },
  // );

  // watch(
  //   () => props.number,
  //   () => {
  //     setOptions(getOption() as any);
  //   },
  // );

  // watch(
  //   () => props.title,
  //   () => {
  //     setOptions(getOption() as any);
  //   },
  // );
  onMounted(() => {
    getInstance()?.on('legendselectchanged', function (params: any) {
      // 立即恢复所有图例的选中状态
      const allSelected = {};
      props.data.forEach((item) => {
        allSelected[item.name] = true;
      });
      getInstance()?.setOption({
        legend: {
          selected: allSelected,
        },
      });
      // if (params?.name.indexOf(props.title) != -1) {
      //   return;
      // }
      emits('onFieldSelect', params?.name);
      return false;
    });
  });
</script>
<style lang="css" scoped></style>
