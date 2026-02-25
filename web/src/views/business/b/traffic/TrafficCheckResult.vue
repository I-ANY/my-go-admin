<template>
  <BasicTable @register="registerTable">
    <template #bodyCell="{ column, record }">
      <template v-if="column.key == 'errorType'">
        <Tag :color="errorTypeMap[record.errorType].color" v-if="errorTypeMap[record.errorType]"
          >{{ errorTypeMap[record.errorType].dictLabel }}
        </Tag>
        <span v-else></span>
      </template>
      <template v-if="column.key == 'agentStatus'">
        <Tag
          :color="agentStatusMap[record.agentStatus].color"
          v-if="agentStatusMap[record.agentStatus]"
          >{{ agentStatusMap[record.agentStatus].dictLabel }}
        </Tag>
      </template>
      <template v-if="column.key == 'businessStatus'">
        <Tag
          :color="businessStatusMap[record.businessStatus].color"
          v-if="businessStatusMap[record.businessStatus]"
          >{{ businessStatusMap[record.businessStatus].dictLabel }}
        </Tag>
      </template>
      <template v-if="column.key == 'networkSpeed'">
        <span>{{ customCeilDivide(record.networkSpeed) }}</span>
      </template>
      <template v-if="column.key == 'mfNetworkSpeed'">
        <span>{{ customCeilDivide(record.mfNetworkSpeed) }}</span>
      </template>
    </template>
  </BasicTable>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import {
    errorTypeMap,
    checkResultFormSchema,
    checkResultColumns,
    businessStatusMap,
    agentStatusMap,
  } from './data';
  import { GetTrafficCheckResultList, GetTrafficOptions } from '@/api/business/b';
  import { nextTick, onMounted } from 'vue';
  import dayjs from 'dayjs';
  import { Tag } from 'ant-design-vue';
  import { customCeilDivide } from '@/utils/util';
  import { RangePickPresetsExact } from '@/utils/common';

  defineOptions({ name: 'BTrafficCheckResult' });
  const [registerTable, { getForm }] = useTable({
    title: '流量差异数据',
    api: GetTrafficCheckResultList,
    columns: checkResultColumns,
    formConfig: {
      labelWidth: 120,
      schemas: checkResultFormSchema(onTimeChange, onTimePikerOpen),
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
    // async beforeFetch(params: Recordable) {
    //   parseValue(params);
    //   return params;
    // },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    rowKey: 'id',
    // actionColumn: {
    //   width: 120,
    //   title: '操作',
    //   dataIndex: 'action',
    //   // slots: { customRender: 'action' },
    //   fixed: 'right',
    // },
  });

  onMounted(async () => {
    await resetReportTime();
  });
  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      reportTimeBegin: dayjs(
        dayjs().add(-7, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      reportTimeEnd: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    });
  }
  async function onTimeChange() {
    await restBusinessOptions();
  }
  async function restBusinessOptions() {
    const data = await getForm().getFieldsValue();
    const res = await GetTrafficOptions({
      reportTimeBegin: data.reportTimeBegin,
      reportTimeEnd: data.reportTimeEnd,
      type: 'business',
    });
    getForm().updateSchema({
      field: 'business',
      componentProps: {
        options: res,
      },
    });
  }
  function onTimePikerOpen() {
    getForm().updateSchema({
      field: '[reportTimeBegin, reportTimeEnd]',
      componentProps: {
        presets: RangePickPresetsExact(),
      },
    });
  }
</script>
