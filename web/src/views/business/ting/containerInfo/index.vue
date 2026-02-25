<template>
  <div>
    <BasicTable @register="registerTable">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'uuid'">
          <Tooltip
            placement="top"
            title="查看容器运行详情"
            v-if="hasPermission(['business:ting:containerDetail:view'])"
          >
            <a @click="handleClick(record)">{{ record.uuid }}</a>
          </Tooltip>
        </template>
        <template v-if="column.key == 'register'">
          <Tag
            v-if="registerStatusMap[record.register]"
            :color="registerStatusMap[record.register].color || 'default'"
            >{{ registerStatusMap[record.register].dictLabel || record.register }}</Tag
          >
        </template>
        <template v-if="column.key == 'status'">
          <Tag
            v-if="containerStatusMap[record.status]"
            :color="containerStatusMap[record.status].color || 'default'"
            >{{ containerStatusMap[record.status].dictLabel || record.status }}</Tag
          >
        </template>
        <template v-if="column.key == 'currentStatus'">
          <Tag
            v-if="currentStatusMap[record.currentStatus]"
            :color="currentStatusMap[record.currentStatus].color || 'default'"
            >{{ currentStatusMap[record.currentStatus].dictLabel || record.currentStatus }}</Tag
          >
        </template>
      </template>
    </BasicTable>
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';

  import {
    columns,
    searchFormSchema,
    registerStatusMap,
    containerStatusMap,
    currentStatusMap,
  } from './data';
  import { GetContainerInfoList } from '@/api/business/ting';
  import { Tag, Tooltip } from 'ant-design-vue';
  import { useGo } from '@/hooks/web/usePage';
  import { usePermission } from '@/hooks/web/usePermission';
  import { RangePickPresetsExact } from '@/utils/common';

  const { hasPermission } = usePermission();
  const go = useGo();
  defineOptions({ name: 'TingContainerInfo' });

  const [registerTable, { setProps, getForm }] = useTable({
    title: '容器信息',
    api: GetContainerInfoList,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema(onTimePikerOpen),
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
    afterFetch: () => {
      getAvailableContainerTotal();
    },
  });
  function handleClick(record: any) {
    go(`/business/ting/container/${record.id}/detail`);
  }
  async function getAvailableContainerTotal() {
    const data = await GetContainerInfoList({ pageSize: 1, pageIndex: 1, currentStatus: 1 });
    setProps({
      title: `容器信息，当前可用总数：${data.total}`,
    });
  }
  function onTimePikerOpen() {
    getForm().updateSchema([
      {
        field: '[metricsTimeBegin, metricsTimeEnd]',
        componentProps: {
          presets: RangePickPresetsExact(),
        },
      },
    ]);
  }
</script>
