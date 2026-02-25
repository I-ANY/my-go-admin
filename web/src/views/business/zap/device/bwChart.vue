<template>
  <div class="bg-white p-4 flex flex-col" style="height: calc(100vh - 120px)">
    <div class="mb-2 text-lg font-bold flex-shrink-0">设备带宽汇总数据展示(单位:G)</div>
    <BasicForm @register="register">
      <template #chartTypeSlot>
        <a-button type="primary" @click="toggleChartType">
          <template #icon>
            <LineChartOutlined v-if="chartType === 'bar'" />
            <BarChartOutlined v-else />
          </template>
          {{ chartType === 'line' ? '切换柱状图' : '切换折线图' }}
        </a-button>
      </template>
    </BasicForm>
    <div v-loading="loading" class="flex-1 flex flex-col min-h-0 w-full overflow-hidden">
      <div ref="chartRef1" class="flex-1 min-h-0 w-full"></div>
      <div ref="chartRef2" class="flex-1 min-h-0 w-full mt-8"></div>
    </div>
  </div>
</template>
<script lang="ts" setup>
  import { onMounted, ref, Ref } from 'vue';
  import { useECharts } from '@/hooks/web/useECharts';
  import { GetBwChartData } from '@/api/business/zap';
  import { RangePickPresetsExact } from '@/utils/common';
  import { BasicForm, useForm, FormSchema } from '@/components/Form';
  import { LineChartOutlined, BarChartOutlined } from '@ant-design/icons-vue';
  import dayjs from 'dayjs';

  const chartRef1 = ref<HTMLDivElement | null>(null);
  const chartRef2 = ref<HTMLDivElement | null>(null);
  const loading = ref(false);
  const chartType = ref<'line' | 'bar'>('line'); // 默认折线图
  const { setOptions: setOptions1 } = useECharts(chartRef1 as Ref<HTMLDivElement>);
  const { setOptions: setOptions2 } = useECharts(chartRef2 as Ref<HTMLDivElement>);

  // 缓存数据用于切换图表类型时重新渲染
  const cachedData = ref<{
    month95: { days: string[]; mobile: number[]; dianlian: number[] } | null;
    day95: { days: string[]; mobile: number[]; dianlian: number[] } | null;
  }>({ month95: null, day95: null });

  // 切换图表类型
  function toggleChartType() {
    chartType.value = chartType.value === 'line' ? 'bar' : 'line';
    // 使用缓存数据重新渲染图表
    if (cachedData.value.month95) {
      const { days, mobile, dianlian } = cachedData.value.month95;
      setOptions1(getBaseOptions('月95', mobile, dianlian, days));
    }
    if (cachedData.value.day95) {
      const { days, mobile, dianlian } = cachedData.value.day95;
      setOptions2(getBaseOptions('日95', mobile, dianlian, days));
    }
  }

  function onTimePikerOpen() {
    // console.log('onTimePikerOpen');
  }

  const utilizationSearchSchema = (onTimePikerOpen, onDateChange): FormSchema[] => [
    {
      field: '[start_date, end_date]',
      label: '时间范围',
      component: 'RangePicker',
      colProps: { span: 7 },
      componentProps: {
        allowClear: false,
        format: 'YYYY-MM-DD',
        valueFormat: 'YYYY-MM-DD',
        showTime: false,
        placeholder: ['开始时间', '结束时间'],
        style: {
          width: '100%',
        },
        presets: RangePickPresetsExact(),
        onOpenChange: onTimePikerOpen,
        onChange: onDateChange,
      },
      required: true,
    },
    {
      field: 'chartType',
      label: ' ',
      component: 'Input',
      slot: 'chartTypeSlot',
      colProps: { span: 4 },
    },
  ];

  const [register, { setFieldsValue, getFieldsValue }] = useForm({
    labelWidth: 100,
    schemas: utilizationSearchSchema(onTimePikerOpen, fetchData),
    actionColOptions: { span: 24 },
    showActionButtonGroup: false,
    autoSubmitOnEnter: true,
  });

  // 解析数据
  const parseData = (items: any[]) => {
    const days: string[] = [];
    const mobile: number[] = [];
    const dianlian: number[] = [];
    items.forEach((item) => {
      days.push(item.date.substring(5));
      mobile.push(item.cmcc);
      dianlian.push(Number((item.cucc + item.ctcc).toFixed(2)));
    });
    return { days, mobile, dianlian };
  };

  const getBaseOptions = (title: string, dataA: number[], dataB: number[], days: string[]): any => {
    const isLine = chartType.value === 'line';
    const seriesType = chartType.value;
    return {
      title: {
        text: title,
        left: 'center',
      },
      tooltip: {
        trigger: 'axis' as const,
        axisPointer: {
          type: isLine ? ('line' as const) : ('shadow' as const),
        },
      },
      legend: {
        data: ['移动', '电联'],
        bottom: 0, // 放到图表底部
        top: 'auto', // 取消顶部定位
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '10%', // 留出空间给底部的图例
        containLabel: true,
        top: 60, // 为标题和图例留出空间
      },
      xAxis: [
        {
          type: 'category' as const,
          data: days,
          axisTick: {
            alignWithLabel: true,
          },
          boundaryGap: !isLine, // 折线图不留边距
        },
      ],
      yAxis: [
        {
          type: 'value' as const,
        },
      ],
      series: [
        {
          name: '移动',
          type: seriesType,
          barWidth: '30%',
          label: {
            show: true,
            position: 'top',
          },
          itemStyle: {
            color: '#fa9751', // 橙色
          },
          smooth: isLine, // 折线图平滑
          symbol: isLine ? 'circle' : 'none',
          symbolSize: isLine ? 8 : 0,
          data: dataA,
        },
        {
          name: '电联',
          type: seriesType,
          barWidth: '30%',
          barGap: '0%',
          label: {
            show: true,
            position: 'top',
          },
          itemStyle: {
            color: '#3aaede', // 青色
          },
          smooth: isLine, // 折线图平滑
          symbol: isLine ? 'circle' : 'none',
          symbolSize: isLine ? 8 : 0,
          data: dataB,
        },
      ],
    };
  };

  async function fetchData() {
    const values = getFieldsValue();
    loading.value = true;
    try {
      const [res_month95, res_day95] = await Promise.all([
        GetBwChartData({ account_type: 'mful', ...values }),
        GetBwChartData({ account_type: 'mfulone', ...values }),
      ]);

      if (res_month95 && res_month95.items) {
        const { days, mobile, dianlian } = parseData(res_month95.items);
        cachedData.value.month95 = { days, mobile, dianlian }; // 缓存数据
        setOptions1(getBaseOptions('月95', mobile, dianlian, days));
      }

      if (res_day95 && res_day95.items) {
        const { days, mobile, dianlian } = parseData(res_day95.items);
        cachedData.value.day95 = { days, mobile, dianlian }; // 缓存数据
        setOptions2(getBaseOptions('日95', mobile, dianlian, days));
      }
    } catch (error) {
      console.error(error);
    } finally {
      loading.value = false;
    }
  }

  function resetReportTime(): Promise<void> {
    return setFieldsValue({
      start_date: dayjs().add(-30, 'day').format('YYYY-MM-DD'),
      end_date: dayjs().format('YYYY-MM-DD'),
    });
  }

  onMounted(async () => {
    await resetReportTime();
    // await fetchData();
  });
</script>
