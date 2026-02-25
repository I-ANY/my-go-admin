<template>
  <BasicTable @register="registerTable">
    <!-- <template #bodyCell="{ column, record }">
        
      </template> -->
  </BasicTable>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { columns, allDearchForm } from './data';
  import { GetContainerDetailList } from '@/api/business/ting';
  import { RangePickPresetsExact } from '@/utils/common';

  defineOptions({ name: 'TingAllContainerDetail' });

  const [registerTable, { getForm }] = useTable({
    title: '容器运行详情',
    api: GetContainerDetailList,
    immediate: true,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: allDearchForm(onTimePikerOpen),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 1,
      // alwaysShowLines: 2,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: true,
    pagination: {
      // pageSizeOptions: ['1', '2', '5'],
    },
  });
  function onTimePikerOpen() {
    console.log('onTimePikerOpen');
    getForm().updateSchema([
      {
        field: '[startTimeBegin, startTimeEnd]',
        componentProps: {
          presets: RangePickPresetsExact(),
        },
      },
      {
        field: '[stopTimeBegin, stopTimeEnd]',
        componentProps: {
          presets: RangePickPresetsExact(),
        },
      },
    ]);
  }
</script>
