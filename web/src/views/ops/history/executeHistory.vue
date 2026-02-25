<template>
  <BasicTable @register="registerTable">
    <template #bodyCell="{ column, record }">
      <template v-if="column.key == 'detail'">
        <Tooltip placement="top" title="查看详情">
          <a @click="onIdClick(record)">查看详情</a>
        </Tooltip>
      </template>
    </template>
  </BasicTable>
</template>
<script setup lang="ts">
  import { BasicTable, useTable } from '@/components/Table';
  import { GetExecuteTaskRecord } from '@/api/ops/execute';
  import { executeResultSearchFormSchema, executeTaskColumns } from './data';
  import { Tooltip } from 'ant-design-vue';
  import { useGo } from '@/hooks/web/usePage';
  import { useCommonStore } from '@/store/modules/common';
  import { splitByLineAndTrim } from '@/utils/util';

  const go = useGo();

  const [registerTable] = useTable({
    title: '执行记录',
    api: GetExecuteTaskRecord,
    columns: executeTaskColumns,
    formConfig: {
      schemas: executeResultSearchFormSchema,
      labelWidth: 80,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
    },
    useSearchForm: true,
    showTableSetting: false,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      // pageSizeOptions: ['10', '30', '50'],
    },
    beforeFetch: (params) => {
      params.hostname = splitByLineAndTrim(params.hostname);
    },
  });

  function onIdClick(record: Recordable) {
    const commonStore = useCommonStore();
    commonStore.clearInspectHostname();
    go('/ops/batchExecute/result/detail/' + record.id);
  }
</script>

<style scoped lang="less"></style>
