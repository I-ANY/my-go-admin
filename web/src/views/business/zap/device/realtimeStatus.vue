<template>
  <div>
    <BasicTable @register="registerTable1">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'metric'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '500px' }">
            <template #title>
              <pre style="margin: 0; white-space: pre-wrap">{{ formatJson(record.metric) }}</pre>
            </template>
            <span>{{ record.metric }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key == 'conclusion'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '600px' }">
            <template #title>{{ record.conclusion }}</template>
            <span
              style="display: block; text-align: left; word-break: break-all; white-space: pre-wrap"
              >{{ record.conclusion ? record.conclusion : '无' }}</span
            >
          </Tooltip>
        </template>
        <template v-if="column.key == 'enable_ips_info'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '300px' }">
            <template #title>{{ record.enable_ips_info }}</template>
            <span
              style="display: block; text-align: left; word-break: break-all; white-space: pre-wrap"
              >{{ record.enable_ips_info ? record.enable_ips_info : '无' }}</span
            >
          </Tooltip>
        </template>
        <template v-if="column.key == 'nat_eth_info'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '500px' }">
            <template #title>{{ record.nat_eth_info }}</template>
            <span>{{ record.nat_eth_info }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key == 'is_metric_abnormal'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '500px' }">
            <template #title>{{ record.is_metric_abnormal }}</template>
            <span>{{ record.is_metric_abnormal }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key == 'disable_ips_info'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '300px' }">
            <template #title>{{ record.disable_ips_info }}</template>
            <span
              style="display: block; text-align: left; word-break: break-all; white-space: pre-wrap"
              >{{ record.disable_ips_info ? record.disable_ips_info : '无' }}</span
            >
          </Tooltip>
        </template>
        <template v-if="column.key == 'state'">
          <span>{{ statesList.filter((item) => item.value == record.state)[0].label }}</span>
        </template>
        <template v-if="column.key == 'is_intranet_resource'">
          <span>{{ record.is_intranet_resource }}</span>
        </template>
      </template>
    </BasicTable>
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { realTimeStatusColumns, realTimeNodeSchema, statesList } from './data';
  import { GetNodeStatusRealtimeBytedance } from '@/api/business/zap';
  import { Tooltip } from 'ant-design-vue';

  defineOptions({ name: 'ZapDeviceStatus' });

  const [registerTable1] = useTable({
    title: '详情',
    api: GetNodeStatusRealtimeBytedance,
    columns: realTimeStatusColumns,
    formConfig: {
      labelWidth: 80,
      schemas: realTimeNodeSchema,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
    },
    immediate: false,
    pagination: true,
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
  });
  function formatJson(value?: string) {
    if (!value) return '无';
    try {
      return JSON.stringify(JSON.parse(value), null, 2);
    } catch (error) {
      return value;
    }
  }
</script>
