<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button @click="guidCheck" type="primary" :loading="loading">状态更新</a-button>
      </template>
      <!--      <template #headerCell="{ column }">-->
      <!--        <template v-if="column.key == 'hostname'"> 主机名</template>-->
      <!--      </template>-->
    </BasicTable>
  </div>
</template>

<script setup lang="ts">
  import { BasicTable, useTable } from '@/components/Table';
  import { GetGuidList, GuidCheck } from '@/api/business/sa';
  import { guidTableColumns, guidSearchSchema } from '@/views/business/sa/guid/data';
  import { splitByLineAndTrim } from '@/utils/util';
  import { message } from 'ant-design-vue';
  import { ref } from 'vue';

  const [registerTable, { reload }] = useTable({
    title: 'GUID列表',
    api: GetGuidList,
    columns: guidTableColumns,
    formConfig: {
      labelWidth: 120,
      schemas: guidSearchSchema,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      // pageSizeOptions: ['10', '30', '50', '100', '500', '1000'],
    },
    beforeFetch: (params) => {
      params.hostnames = splitByLineAndTrim(params.hostnames) || null;
      params.guids = splitByLineAndTrim(params.guids) || null;
    },
  });

  const loading = ref(false);

  // GUID检测
  async function guidCheck() {
    loading.value = true;
    await GuidCheck();
    message.success('状态更新完成');
    loading.value = false;
    reload();
  }
</script>

<style scoped lang="less"></style>
