<template>
  <div class="p-4">
    <Card title="质量事件实时查询" :bordered="false">
      <BasicForm
        autoFocusFirstItem
        :labelWidth="100"
        :schemas="schemas"
        :actionColOptions="{ span: 4 }"
        @submit="handleSubmit"
      />
    </Card>
    <BasicTable :loading="loading" @register="registerRefundTable" v-show="show">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'fields'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '800px' }">
            <template #title>{{ record.fields }}</template>
            <span>{{ truncatedFields(record.fields) }}</span>
          </Tooltip>
        </template>
      </template>
    </BasicTable>
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, FormSchema, useTable } from '@/components/Table';
  import { BasicForm } from '@/components/Form';
  import { columns, truncatedFields } from './data';
  import { GetQualityEventsRealtimeBytedance } from '@/api/business/zap';
  import { ref, reactive, onMounted } from 'vue';
  import { Card, Tooltip } from 'ant-design-vue';
  import { useMessage } from '@/hooks/web/useMessage';
  import { useLoading } from '@/components/Loading';

  const show = ref(false);
  const loading = ref(false);
  const schemas: FormSchema[] = [
    {
      field: 'hostname',
      component: 'Input',
      label: '主机名',
      defaultValue: '',
      valueField: 'hostname',
      colProps: {
        span: 6,
      },
    },
    {
      field: 'label',
      component: 'Input',
      label: '节点label',
      defaultValue: '',
      valueField: 'label',
      colProps: {
        span: 6,
      },
    },
    {
      label: '类型',
      field: 'account_type',
      component: 'Select',
      defaultValue: 'day95',
      required: true,
      colProps: {
        span: 4,
      },
      componentProps: {
        options: [
          {
            label: '日95',
            value: 'day95',
          },
          {
            label: '月95',
            value: 'month95',
          },
        ],
      },
    },
    {
      field: 'sync_scope',
      component: 'Select',
      defaultValue: 0.5,
      label: '时间范围',
      colProps: {
        span: 4,
      },
      componentProps: {
        options: [
          {
            label: '十分钟内',
            value: 0.16,
            key: '1',
          },
          {
            label: '半小时内',
            value: 0.5,
            key: '2',
          },
          {
            label: '一小时内',
            value: 1,
            key: '3',
          },
          {
            label: '两小时内',
            value: 2,
            key: '4',
          },
          {
            label: '一天内',
            value: 24,
            key: '5',
          },
          {
            label: '两天内',
            value: 48,
            key: '6',
          },
          {
            label: '三天内',
            value: 72,
            key: '7',
          },
          {
            label: '一周内',
            value: 166,
            key: '8',
          },
        ],
      },
    },
  ];

  const tableData = reactive<any>([]);
  const [registerRefundTable, { setTableData }] = useTable({
    title: '',
    // dataSource: tableData.data,
    columns: columns,
    pagination: false,
    showIndexColumn: false,
    scroll: { y: 150 },
    showSummary: true,
  });

  const { createMessage } = useMessage();
  const { info, error } = createMessage;
  const [openFullLoading, closeFullLoading] = useLoading({
    tip: '加载中...',
  });

  async function handleSubmit(formData: any) {
    if (formData.label == '' && formData.hostname == '') {
      error('请输入节点label或主机名！');
      return;
    }
    if (formData.label && formData.hostname) {
      error('不能同时输入节点label和主机名！');
      return;
    }
    openFullLoading();
    const params = {
      label: formData.label,
      hostname: formData.hostname,
      sync_scope: formData.sync_scope,
      account_type: formData.account_type,
    };
    await GetQualityEventsRealtimeBytedance(params)
      .then((res) => {
        if (res.length > 0) {
          closeFullLoading();
          show.value = true;
          setTableData(res);
        } else {
          show.value = false;
          closeFullLoading();
          info('查询成功，暂无数据！');
        }
      })
      .catch(() => {
        closeFullLoading();
      });
  }

  onMounted(() => {
    tableData.length = 0;
  });
</script>
