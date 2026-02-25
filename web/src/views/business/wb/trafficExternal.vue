<template>
  <div class="zp-traffic-external">
    <BasicTable @register="registerTable" />
  </div>
</template>

<script setup lang="ts">
  import { BasicTable, useTable } from '@/components/Table';
  import {
    trafficExternalColumns,
    trafficExternalSearchFormSchema,
  } from '@/views/business/wb/data';
  import { GetZPExtneralTrafficDetail } from '@/api/business/wb';
  import { splitByLineAndTrim } from '@/utils/util';
  import { onMounted } from 'vue';
  import dayjs from 'dayjs';

  const [registerTable, tableMethods] = useTable({
    title: 'ZP外部流量明细',
    api: GetZPExtneralTrafficDetail,
    columns: trafficExternalColumns,
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    immediate: false,
    formConfig: {
      labelWidth: 120,
      showAdvancedButton: true,
      autoAdvancedLine: 3,
      schemas: trafficExternalSearchFormSchema(),
    },
    pagination: {
      pageSizeOptions: ['10', '20', '50', '100', '200'],
    },
    beforeFetch: (params: Record<string, any>) => {
      const range = params['[metricsTimeBegin, metricsTimeEnd]'];
      if (Array.isArray(range) && range.length === 2) {
        params.metricsTimeBegin = range[0]?.format
          ? range[0].format('YYYY-MM-DD HH:mm:ss')
          : range[0];
        params.metricsTimeEnd = range[1]?.format
          ? range[1].format('YYYY-MM-DD HH:mm:ss')
          : range[1];
      }
      delete params['[metricsTimeBegin, metricsTimeEnd]'];

      const parsedHostnames = splitByLineAndTrim(params.hostnames);
      if (parsedHostnames && parsedHostnames.length > 0) {
        params.hostnames = parsedHostnames;
      } else {
        delete params.hostnames;
      }

      const parsedBusinesses = splitByLineAndTrim(params.businesses);
      if (parsedBusinesses && parsedBusinesses.length > 0) {
        params.businesses = parsedBusinesses;
      } else {
        delete params.businesses;
      }

      if (!params.sn) {
        delete params.sn;
      }
      if (!params.isp) {
        delete params.isp;
      }
      return params;
    },
  });

  const { getForm, reload } = tableMethods;

  onMounted(async () => {
    const start = dayjs().subtract(1, 'day').startOf('day');
    const end = dayjs().subtract(1, 'day').endOf('day');
    await getForm().setFieldsValue({
      '[metricsTimeBegin, metricsTimeEnd]': [start, end],
    });
    await reload({ page: 1 });
  });
</script>

<style scoped lang="less">
  .zp-traffic-external {
    padding: 16px;
  }
</style>
