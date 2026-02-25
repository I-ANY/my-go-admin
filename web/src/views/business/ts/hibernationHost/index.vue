<template>
  <div>
    <BasicTable @register="registerTable">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'hostname'">
          <span>{{ record.hostname }}</span>
        </template>
        <template v-if="column.key == 'hibernationTime'">
          <span>{{ record.hibernationTime }}</span>
        </template>
      </template>
    </BasicTable>
  </div>
</template>

<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { columns, searchFormSchema } from './data';
  import { RangeDatafor1dayPickPresetsExact } from '@/utils/common';
  import dayjs from 'dayjs';
  import { nextTick, onMounted } from 'vue';
  import { GetTsHibernationHost } from '@/api/business/ts';
  import { splitByLineAndTrim } from '@/utils/util';

  defineOptions({ name: 'HibernationHost' });
  const [registerTable, { getForm }] = useTable({
    title: '休眠主机',
    api: GetTsHibernationHost,
    columns,
    formConfig: {
      labelWidth: 150,
      schemas: searchFormSchema(onTimePikerOpen, disabledDate),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 3,
      submitOnReset: false,
      alwaysShowLines: 1,
      resetFunc() {
        nextTick(() => {
          resetReportTime();
        });
        return Promise.resolve();
      },
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    // rowKey: 'version',
    // actionColumn: {
    //   width: 120,
    //   title: '操作',
    //   dataIndex: 'action',
    //   // slots: { customRender: 'action' },
    //   fixed: 'right',
    // },

    beforeFetch: (params) => {
      // parseParams(params);
      // return params;
      parseValue(params);
      return params;
    },
  });

  function parseValue(params: Recordable) {
    params.hostnames = splitByLineAndTrim(params.hostnames);
    // if (params.reportDayBegin) {
    //   params.reportDayBegin = dayjs(params.reportDayBegin).format('YYYY-MM-DD');
    // }
    // if (params.reportDayEnd) {
    //   params.reportDayEnd = dayjs(params.reportDayEnd).format('YYYY-MM-DD');
    // }
  }

  function onTimePikerOpen() {
    getForm().updateSchema({
      field: '[metricsTimeBegin, metricsTimeEnd]',
      componentProps: {
        presets: RangeDatafor1dayPickPresetsExact(),
        // onChange: onTimeRangeChange,
        // disabledDate: disabledDate,
      },
    });
  }

  function disabledDate(current: any, dateStrings: any) {
    // 如果没有选择开始时间，不禁用任何日期
    if (!dateStrings || !dateStrings[0]) {
      return false;
    }

    const startDate = dayjs(dateStrings[0]);
    const currentDate = dayjs(current);

    // 如果当前日期在开始日期之前，不禁用
    if (currentDate.isBefore(startDate)) {
      return false;
    }

    // 如果当前日期超过开始日期1天，则禁用
    const diffDays = currentDate.diff(startDate, 'day', true);
    return diffDays > 1;
  }

  onMounted(async () => {
    await resetReportTime();
  });
  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      metricsTimeBegin: dayjs(
        dayjs().subtract(1, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      metricsTimeEnd: dayjs(
        dayjs().subtract(1, 'day').format('YYYY-MM-DD 23:59:59'),
        'YYYY-MM-DD HH:mm:ss',
      ),
    });
  }
</script>
