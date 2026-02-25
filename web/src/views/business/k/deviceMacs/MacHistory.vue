<template>
  <div>
    <BasicTable @register="registerTable">
      <template #title> MAC历史记录 </template>
    </BasicTable>
  </div>
</template>
<script setup lang="ts">
  import { BasicTable, useTable } from '@/components/Table';
  import { MacHistoryColumns, MacHistorySearchFormSchema } from './data';
  import { GetMacHistory } from '@/api/business/k';

  const [registerTable] = useTable({
    title: 'MAC历史记录',
    api: GetMacHistory,
    columns: MacHistoryColumns,
    // 使用搜索表单配置
    useSearchForm: true,
    formConfig: {
      labelWidth: 120,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      alwaysShowLines: 2,
      schemas: MacHistorySearchFormSchema(),
    },
    beforeFetch: (params) => {
      return params;
    },
    canResize: true,
    bordered: true,
    showIndexColumn: true,
    rowKey: 'id',
  });
</script>
