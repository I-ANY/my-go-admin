<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    title="历史利用率考核"
    @ok="handleOk"
    width="80%"
    :show-ok-btn="false"
    :show-cancel-btn="true"
    :destroy-on-close="true"
    cancel-text="关闭"
  >
    <div class="history-assessment-modal">
      <div class="header-info">
        <div class="info-item">
          <span class="label">节点编号:</span>
          <span class="value">{{ currentRecord?.roomNo || '-' }}</span>
        </div>
        <div class="info-item">
          <span class="label">运营商:</span>
          <span class="value">{{ currentRecord?.localIsp || '-' }}</span>
        </div>
        <div class="info-item">
          <span class="label">所在地:</span>
          <span class="value">{{ currentRecord?.location || '-' }}</span>
        </div>
        <div class="info-item">
          <span class="label">统计类型:</span>
          <span class="value">{{ reportTypeMap[currentRecord?.reportType] || '-' }}</span>
        </div>
        <template v-if="currentRecord?.roomType === 3">
          <div class="info-item">
            <span class="label">主线业务:</span>
            <span class="value">{{ currentRecord?.biz || '-' }}</span>
          </div>
        </template>
      </div>

      <!-- 搜索表单 -->
      <BasicForm @register="registerSearchForm" />

      <BasicTable @register="registerTable" />
    </div>
  </BasicModal>
</template>

<script lang="ts" setup>
  import { ref, h } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicTable, useTable, BasicColumn, FormSchema } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form';
  import { getCompositeHistory } from '@/api/business/a';
  import { Tag } from 'ant-design-vue';
  import dayjs from 'dayjs';

  defineOptions({ name: 'HistoryAssessmentModal' });

  const currentRecord = ref<any>(null);

  // 搜索表单配置
  const searchFormSchema: FormSchema[] = [
    {
      field: 'dateRange',
      label: '日期范围',
      component: 'RangePicker',
      componentProps: {
        format: 'YYYY-MM-DD',
        valueFormat: 'YYYY-MM-DD',
        placeholder: ['开始日期', '结束日期'],
      },
      colProps: { span: 12 },
    },
  ];

  // 搜索处理函数
  const handleSearch = async () => {
    await loadHistoryData();
  };

  const [registerSearchForm, { getFieldsValue, setFieldsValue }] = useForm({
    labelWidth: 80,
    schemas: searchFormSchema,
    showActionButtonGroup: true,
    actionColOptions: { span: 12 },
    autoSubmitOnEnter: true,
    submitFunc: handleSearch,
  });

  // 格式化百分比函数，保留两位小数
  const formatPercentage = (value: number | null | undefined): string => {
    if (value == null || value === undefined) return '-';
    return `${(value * 100).toFixed(2)}%`;
  };

  const stateMap = {
    1: '是',
    2: '否',
    3: '暂无要求',
  };

  const stateManualMap = {
    1: '是',
    2: '否',
    3: '暂无要求',
    4: '已裁撤',
    5: '不计费',
    6: '其他',
    7: '待确认',
  };

  const planTypeMap = {
    1: '不扣分',
    2: '有方案',
    3: '无方案',
    4: '临时特批',
    5: '其他',
  };

  const reportTypeMap = {
    1: '节点',
    2: '保底',
    3: '削峰',
  };

  const columns: BasicColumn[] = [
    {
      title: '日期',
      dataIndex: 'date',
      width: 120,
      customRender: ({ record }) => {
        return dayjs(record.date).format('YYYY-MM-DD');
      },
    },
    {
      title: '考核业务',
      dataIndex: 'assessmentBusiness',
      width: 150,
      customRender: ({ record }) => {
        const bizDisplay = record.assessmentBusiness ? record.assessmentBusiness.join(', ') : '-';
        return bizDisplay;
      },
    },
    {
      title: '日95利用率',
      dataIndex: 'bwUsageRateDay',
      width: 120,
      customRender: ({ record }) => {
        return formatPercentage(record.bwUsageRateDay);
      },
    },
    {
      title: '晚高峰利用率',
      dataIndex: 'bwUsageRateNight',
      width: 120,
      customRender: ({ record }) => {
        return formatPercentage(record.bwUsageRateNight);
      },
    },
    {
      title: '达标点数',
      dataIndex: 'nightPointNum',
      width: 100,
    },
    {
      title: '是否达标',
      dataIndex: 'state',
      width: 100,
      customRender: ({ record }) => {
        const color = record.state === 1 ? 'green' : record.state === 2 ? 'red' : 'default';
        return h(Tag, { color }, () => stateMap[record.state] || '-');
      },
    },
    {
      title: '达标确认',
      dataIndex: 'stateManual',
      width: 100,
      customRender: ({ record }) => {
        const color =
          record.stateManual === 1 ? 'green' : record.stateManual === 2 ? 'red' : 'default';
        return h(Tag, { color }, () => stateManualMap[record.stateManual] || '-');
      },
    },
    {
      title: '评分',
      dataIndex: 'score',
      width: 80,
      customRender: ({ record }) => {
        return record.score || '0';
      },
    },
    {
      title: '方案类型',
      dataIndex: 'planType',
      width: 120,
      customRender: ({ record }) => {
        return planTypeMap[record.planType] || '-';
      },
    },
    {
      title: '方案内容',
      dataIndex: 'planContent',
      width: 200,
      customRender: ({ record }) => {
        return record.planContent || '-';
      },
    },
    {
      title: '备注',
      dataIndex: 'describe',
      width: 200,
      customRender: ({ record }) => {
        return record.describe || '-';
      },
    },
  ];

  const [registerTable, { setTableData, setLoading }] = useTable({
    columns,
    bordered: true,
    showIndexColumn: true,
    pagination: false,
    rowKey: (record) => record.id || record.date,
    immediate: false,
  });

  const [registerModal, { closeModal }] = useModalInner(async (data) => {
    currentRecord.value = data.record;

    // 设置默认日期范围：currentRecord.date的前三天到currentRecord.date
    if (currentRecord.value?.date) {
      const currentDate = dayjs(currentRecord.value.date);
      const startDate = currentDate.subtract(3, 'day').format('YYYY-MM-DD');
      const endDate = currentDate.format('YYYY-MM-DD');

      setFieldsValue({
        dateRange: [startDate, endDate],
      });
    }

    await loadHistoryData();
  });

  const loadHistoryData = async () => {
    if (!currentRecord.value) return;

    try {
      setLoading(true);
      const formValues = getFieldsValue();
      console.log('formValues', formValues);
      const params: any = {
        roomNo: currentRecord.value.roomNo,
        localIsp: currentRecord.value.localIsp,
        location: currentRecord.value.location,
        reportType: currentRecord.value.reportType,
      };
      if (currentRecord.value.roomType === 3) {
        params.biz = currentRecord.value.biz;
      }

      // 添加日期范围参数
      if (formValues.dateRange && formValues.dateRange.length === 2) {
        params.startDate = formValues.dateRange[0];
        params.endDate = formValues.dateRange[1];
      }
      const result = await getCompositeHistory(params);
      const historyData = result?.items || [];
      setTableData(historyData);
    } catch (error) {
      console.error('获取历史数据失败:', error);
      setTableData([]);
    } finally {
      setLoading(false);
    }
  };

  const handleOk = () => {
    closeModal();
  };
</script>

<style lang="less" scoped>
  .history-assessment-modal {
    .header-info {
      display: flex;
      flex-wrap: wrap;
      margin-bottom: 16px;
      padding: 16px;
      border-radius: 6px;
      background-color: #f5f5f5;
      gap: 16px;

      .info-item {
        display: flex;
        align-items: center;
        min-width: 200px;

        .label {
          margin-right: 8px;
          color: #666;
          font-weight: 500;
        }

        .value {
          color: #333;
          font-weight: 600;
        }
      }
    }

    :deep(.ant-form) {
      margin-bottom: 16px;
      padding: 16px;
      border-radius: 6px;
      background-color: #fafafa;
    }
  }
</style>
