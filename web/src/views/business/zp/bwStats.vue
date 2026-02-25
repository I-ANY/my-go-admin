<template>
  <div class="bg-white p-4 flex flex-col" style="height: calc(100vh - 120px)">
    <div class="mb-2 text-lg font-bold flex-shrink-0">带宽统计折线图</div>
    <BasicForm @register="register" />
    <div v-loading="loading" class="flex-1 min-h-0 w-full overflow-hidden">
      <div ref="chartRef" class="h-full w-full"></div>
    </div>
    <!-- 运营商切换按钮（多选） -->
    <div class="flex justify-center gap-2 mt-3">
      <a-button
        v-for="item in ispButtons"
        :key="item.value"
        :type="selectedIsps.includes(item.value) ? 'primary' : 'default'"
        size="small"
        @click="toggleIsp(item.value)"
      >
        {{ item.label }}
      </a-button>
    </div>
    <!-- 差异数据弹框 -->
    <BwStatsDiffModal @register="registerDiffModal" />
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, ref, Ref } from 'vue';
  import { useECharts } from '@/hooks/web/useECharts';
  import { GetBwStats } from '@/api/business/zp';
  import { RangePickPresetsExact } from '@/utils/common';
  import { BasicForm, useForm, FormSchema } from '@/components/Form';
  import { useModal } from '@/components/Modal';
  import dayjs from 'dayjs';
  import BwStatsDiffModal from './BwStatsDiffModal.vue';

  defineOptions({ name: 'ZpBwStats' });

  const chartRef = ref<HTMLDivElement | null>(null);
  const loading = ref(false);
  const { setOptions, getInstance } = useECharts(chartRef as Ref<HTMLDivElement>);

  // 差异数据弹框
  const [registerDiffModal, { openModal: openDiffModal }] = useModal();

  // 当前选中的运营商（支持多选）
  const selectedIsps = ref<string[]>(['联通', '移动', '电信']);

  // 运营商按钮配置
  const ispButtons = [
    { label: '联通', value: '联通' },
    { label: '移动', value: '移动' },
    { label: '电信', value: '电信' },
    { label: '汇总', value: '汇总' },
  ];

  // 颜色配置 - 每个运营商对应两种颜色（真实带宽和平台带宽）
  const colorConfig: Record<string, Record<string, string>> = {
    联通: {
      real_bandwidth: '#1890ff', // 蓝色
      platform_bandwidth: '#69c0ff', // 浅蓝色
    },
    移动: {
      real_bandwidth: '#52c41a', // 绿色
      platform_bandwidth: '#95de64', // 浅绿色
    },
    电信: {
      real_bandwidth: '#fa8c16', // 橙色
      platform_bandwidth: '#ffc069', // 浅橙色
    },
    汇总: {
      real_bandwidth: '#722ed1', // 紫色
      platform_bandwidth: '#b37feb', // 浅紫色
    },
  };

  // 线条样式配置
  const lineStyleConfig: Record<string, string> = {
    real_bandwidth: 'solid',
    platform_bandwidth: 'dashed',
  };

  const searchFormSchema: FormSchema[] = [
    {
      field: '[start_time, end_time]',
      label: '时间范围',
      component: 'RangePicker',
      colProps: { span: 10 },
      componentProps: {
        allowClear: false,
        format: 'YYYY-MM-DD HH:mm',
        valueFormat: 'YYYY-MM-DD HH:mm:ss',
        showTime: { format: 'HH:mm' },
        placeholder: ['开始时间', '结束时间'],
        style: { width: '100%' },
        presets: RangePickPresetsExact(),
        onChange: fetchData,
      },
      required: true,
    },
  ];

  const [register, { setFieldsValue, getFieldsValue }] = useForm({
    labelWidth: 100,
    schemas: searchFormSchema,
    actionColOptions: { span: 24 },
    showActionButtonGroup: false,
    autoSubmitOnEnter: true,
  });

  // 缓存原始数据
  let cachedData: any[] = [];

  // 中文运营商对应后端数据的 key
  const ispKeyMap: Record<string, string> = {
    联通: 'cucc',
    移动: 'cmcc',
    电信: 'ctcc',
  };

  // 切换运营商（多选）
  function toggleIsp(isp: string) {
    const index = selectedIsps.value.indexOf(isp);
    if (index > -1) {
      // 至少保留一个选中
      if (selectedIsps.value.length > 1) {
        selectedIsps.value.splice(index, 1);
      }
    } else {
      selectedIsps.value.push(isp);
    }
    updateChart();
  }

  // 查看差异数据
  function handleViewDiff(datetime: string) {
    openDiffModal(true, { datetime, isps: selectedIsps.value });
  }

  // 绑定图表点击事件
  function bindChartClick() {
    const chart = getInstance();
    if (chart) {
      chart.off('click'); // 先移除之前的监听
      chart.on('click', (params: any) => {
        if (params.componentType === 'series') {
          const datetime = params.name || '';
          if (datetime) {
            handleViewDiff(datetime);
          }
        }
      });
    }
  }

  // 解析数据并生成图表配置
  function generateChartOptions(items: any[], isps: string[]) {
    if (!items || items.length === 0 || isps.length === 0) {
      return null;
    }

    // 提取时间轴
    const xAxisData = items.map((item) => item.datetime);

    // 生成系列数据和图例
    const series: any[] = [];
    const legendData: string[] = [];

    isps.forEach((isp) => {
      // 获取后端数据对应的 key
      const dataKey = ispKeyMap[isp];
      const isTotal = isp === '汇总';

      // 真实带宽线
      const realBandwidthData = items.map((item) => {
        const bandwidthData = item.real_bandwidth || {};
        if (isTotal) {
          // 汇总：计算三个运营商的总和
          return (bandwidthData.cucc || 0) + (bandwidthData.cmcc || 0) + (bandwidthData.ctcc || 0);
        }
        return bandwidthData[dataKey] || 0;
      });
      const realSeriesName = `${isp}-真实带宽`;
      legendData.push(realSeriesName);
      series.push({
        name: realSeriesName,
        type: 'line' as const,
        smooth: true,
        symbol: 'circle',
        symbolSize: isTotal ? 8 : 6,
        lineStyle: {
          width: isTotal ? 3 : 2,
          type: lineStyleConfig.real_bandwidth,
        },
        itemStyle: {
          color: colorConfig[isp]?.real_bandwidth || '#1890ff',
        },
        data: realBandwidthData,
      });

      // 平台带宽线
      const platformBandwidthData = items.map((item) => {
        const bandwidthData = item.platform_bandwidth || {};
        if (isTotal) {
          // 汇总：计算三个运营商的总和
          return (bandwidthData.cucc || 0) + (bandwidthData.cmcc || 0) + (bandwidthData.ctcc || 0);
        }
        return bandwidthData[dataKey] || 0;
      });
      const platformSeriesName = `${isp}-平台带宽`;
      legendData.push(platformSeriesName);
      series.push({
        name: platformSeriesName,
        type: 'line' as const,
        smooth: true,
        symbol: 'circle',
        symbolSize: isTotal ? 8 : 6,
        lineStyle: {
          width: isTotal ? 3 : 2,
          type: lineStyleConfig.platform_bandwidth,
        },
        itemStyle: {
          color: colorConfig[isp]?.platform_bandwidth || '#fa8c16',
        },
        data: platformBandwidthData,
      });
    });

    return {
      title: {
        text: '带宽趋势图',
        left: 'center',
      },
      tooltip: {
        trigger: 'axis' as const,
        axisPointer: {
          type: 'cross' as const,
        },
        formatter: (params: any) => {
          const paramList = Array.isArray(params) ? params : [params];
          const datetime = paramList[0]?.axisValue || '';
          let result = `<div style="font-weight: bold; margin-bottom: 8px;">${datetime}</div>`;
          paramList.forEach((param: any) => {
            const value =
              typeof param.value === 'number' ? param.value.toLocaleString() : param.value;
            result += `<div style="display: flex; justify-content: space-between; min-width: 200px;">
              <span>${param.marker} ${param.seriesName}</span>
              <span style="font-weight: bold; margin-left: 20px;">${value}</span>
            </div>`;
          });
          result += `<div style="margin-top: 8px; color: #999; font-size: 11px; text-align: center;">点击数据点查看差异数据</div>`;
          return result;
        },
      },
      legend: {
        data: legendData,
        bottom: 10,
        selectedMode: 'multiple' as const, // 支持多选点击切换
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '15%',
        top: 60,
        containLabel: true,
      },
      xAxis: {
        type: 'category' as const,
        boundaryGap: false,
        data: xAxisData,
        axisLabel: {
          show: true,
          interval: 'auto',
          formatter: (value: string) => {
            if (value && value.length >= 16) {
              return value.substring(5, 16); // 显示 MM-DD HH:mm
            }
            return value;
          },
          rotate: 45,
          fontSize: 12,
        },
        axisTick: {
          show: true,
          alignWithLabel: true,
        },
        axisLine: {
          show: true,
        },
      },
      yAxis: {
        type: 'value' as const,
        name: '带宽',
        axisLabel: {
          formatter: (value: number) => {
            if (value >= 1000) {
              return (value / 1000).toFixed(1) + 'G';
            }
            return value.toString();
          },
        },
      },
      series,
    };
  }

  // 更新图表
  async function updateChart() {
    const options = generateChartOptions(cachedData, selectedIsps.value);
    if (options) {
      await setOptions(options);
      // 绑定点击事件
      bindChartClick();
    }
  }

  // 获取数据
  async function fetchData() {
    const values = getFieldsValue();
    if (!values.start_time || !values.end_time) {
      return;
    }

    loading.value = true;
    try {
      const res = await GetBwStats({
        start_time: values.start_time,
        end_time: values.end_time,
      });

      if (res && res.items) {
        cachedData = res.items;
        await updateChart();
      }
    } catch (error) {
      console.error('获取带宽统计数据失败:', error);
    } finally {
      loading.value = false;
    }
  }

  // 初始化时间范围
  function resetReportTime(): Promise<void> {
    return setFieldsValue({
      start_time: dayjs().add(-1, 'day').format('YYYY-MM-DD 00:00:00'),
      end_time: dayjs().format('YYYY-MM-DD HH:mm:ss'),
    });
  }

  onMounted(async () => {
    await resetReportTime();
    await fetchData();
  });
</script>

<style scoped lang="less"></style>
