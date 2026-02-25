<template>
  <div>
    <BasicTable @register="registerTable">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'uts_vm_hostname'">
          <span>{{ record.uts_vm_hostname }}</span>
        </template>
        <template v-if="column.key == 'hostname'">
          <span>{{ record.hostname }}</span>
        </template>
        <template v-if="column.key == 'exceeded95'">
          <Tag :color="record.exceeded95 ? 'red' : 'green'">
            {{ record.exceeded95 ? '异常' : '正常' }}
          </Tag>
        </template>
        <template v-if="column.key == 'exceeded_plan_bw'">
          <Tag :color="record.exceeded_plan_bw ? 'red' : 'green'">
            {{ record.exceeded_plan_bw ? '异常' : '正常' }}
          </Tag>
        </template>
        <template v-if="column.key == 'exceeded_time'">
          <Tag :color="record.exceeded_time ? 'red' : 'green'">
            {{ record.exceeded_time ? '异常' : '正常' }}
          </Tag>
        </template>
        <template v-if="column.key == 'businessStatus'">
          <Tag :color="record.businessStatus ? 'red' : 'green'">
            {{ record.businessStatus ? '异常' : '正常' }}
          </Tag>
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
  import { GetTsDetail } from '@/api/business/ts';
  import { splitByLineAndTrim } from '@/utils/util';
  import { Tag, message } from 'ant-design-vue';

  defineOptions({ name: 'UTS带宽抢占告警' });
  const [registerTable, { getForm }] = useTable({
    title: 'UTS带宽抢占告警',
    api: GetTsDetail,
    columns,
    formConfig: {
      labelWidth: 150,
      schemas: searchFormSchema(onTimePikerOpen, onTimeRangeChange, disabledDate),
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
    beforeFetch: (params) => {
      // parseParams(params);
      // return params;
      parseValue(params);
      return params;
    },
    // actionColumn: {
    //   width: 120,
    //   title: '操作',
    //   dataIndex: 'action',
    //   // slots: { customRender: 'action' },
    //   fixed: 'right',
    // },
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
        onChange: onTimeRangeChange,
        disabledDate: disabledDate,
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

  function onTimeRangeChange(dates: any) {
    if (dates && dates.length === 2) {
      const [startDate, endDate] = dates;
      const start = dayjs(startDate);
      const end = dayjs(endDate);

      // 计算时间差（以天为单位）
      const diffDays = end.diff(start, 'day', true);

      // 如果选择的时间范围超过1天，则重置为1天
      if (diffDays > 1) {
        message.warning('只能选择一天的时间范围');
        // 重置为从开始时间算起的1天
        const newEndDate = start.add(1, 'day').subtract(1, 'second');
        getForm().setFieldsValue({
          metricsTimeBegin: start,
          metricsTimeEnd: newEndDate,
        });
      }
    }
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
