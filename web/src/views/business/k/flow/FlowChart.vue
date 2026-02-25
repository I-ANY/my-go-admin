<template>
  <div class="p-4">
    <div class="bg-white rounded-md shadow-sm p-4">
      <BasicForm :formSchema="searchFormSchema" @register="registerForm" @submit="fetchData" />
      <div ref="chartRef" style="height: 600px; margin-top: 20px"></div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  import { BasicForm, useForm } from '@/components/Form';
  import { GetFlowStats } from '@/api/business/k';
  import { searchFormSchema, getChartOptions } from './data';
  import { ref, Ref, onMounted } from 'vue';
  import { useECharts } from '@/hooks/web/useECharts';
  import dayjs from 'dayjs';
  import { useAreaSelect } from '@/utils/kAreaSelect';

  // 折线图配置
  const chartRef = ref<HTMLDivElement>();
  const { setOptions } = useECharts(chartRef as Ref<HTMLDivElement>);

  const [registerForm, { updateSchema, getFieldsValue, setFieldsValue }] = useForm({
    labelWidth: 100,
    schemas: searchFormSchema,
    autoSubmitOnEnter: true,
    showAdvancedButton: true,
    autoAdvancedLine: 4,
    actionColOptions: {
      span: 24,
    },
  });

  onMounted(async () => {
    // 获取表单默认值
    const defaultParams = await getFieldsValue();
    // 触发数据加载
    fetchData(defaultParams);
    try {
      // 区域省份
      const { initAreaData } = useAreaSelect({
        form: {
          updateSchema,
          setFieldsValue,
          getFieldsValue,
        },
        fields: {
          area: 'area',
          province: 'province',
        },
      });
      await initAreaData();
    } catch (error) {
      console.error('获取数据失败:', error);
    }
  });

  const fetchData = async (params: any) => {
    if (!params.provider_id) {
      params.provider_id = 55;
    }
    if (!params.date) {
      params.date = dayjs().subtract(1, 'day').format('YYYY-MM-DD');
    }

    const { items = [] } = await GetFlowStats(params);
    setOptions(getChartOptions(items));
  };
</script>
