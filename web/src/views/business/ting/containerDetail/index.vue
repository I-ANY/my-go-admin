<template>
  <PageWrapper :title="data.title" contentFullHeight contentBackground @back="goBack">
    <template #headerContent>
      <span class="m-20px">容器名称：{{ data.containerInfo.container }}</span>
      <span>
        实时状态：
        <Tag
          v-if="currentStatusMap[data.containerInfo.currentStatus]"
          :color="currentStatusMap[data.containerInfo.currentStatus].color || 'default'"
          >{{
            currentStatusMap[data.containerInfo.currentStatus].dictLabel ||
            data.containerInfo.currentStatus
          }}</Tag
        >
      </span>
      <span>
        容器状态：
        <Tag
          v-if="containerStatusMap[data.containerInfo.status]"
          :color="containerStatusMap[data.containerInfo.status].color || 'default'"
          >{{
            containerStatusMap[data.containerInfo.status].dictLabel || data.containerInfo.status
          }}</Tag
        >
      </span>
      <span>
        注册状态：
        <Tag
          v-if="registerStatusMap[data.containerInfo.register]"
          :color="registerStatusMap[data.containerInfo.register].color || 'default'"
          >{{
            registerStatusMap[data.containerInfo.register].dictLabel || data.containerInfo.register
          }}</Tag
        >
      </span>
    </template>
    <BasicTable @register="registerTable">
      <!-- <template #bodyCell="{ column, record }">
        
      </template> -->
    </BasicTable>
  </PageWrapper>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { PageWrapper } from '@/components/Page';
  import {
    columns,
    searchFormSchema,
    currentStatusMap,
    registerStatusMap,
    containerStatusMap,
  } from './data';
  import { GetContainerDetailList, GetContainerInfo } from '@/api/business/ting';
  import { useRoute } from 'vue-router';
  import { useGo } from '@/hooks/web/usePage';
  import { onMounted, reactive } from 'vue';
  import { Tag } from 'ant-design-vue';
  import { RangePickPresetsExact } from '@/utils/common';

  defineOptions({ name: 'TingContainerDetail' });

  const route = useRoute();
  const go = useGo();
  const data = reactive({
    title: '',
    id: 0,
    containerInfo: {} as Recordable,
  });
  data.id = route.params.id as any;

  const [registerTable, { reload, getForm }] = useTable({
    title: '容器运行详情',
    api: GetContainerDetailList,
    beforeFetch(params) {
      params.uuid = data.containerInfo.uuid;
      return params;
    },
    immediate: false,
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
  });
  function goBack() {
    go('/business/ting/container/info');
  }
  onMounted(async () => {
    const res = await GetContainerInfo(data.id);
    data.containerInfo = res;
    reload();
    data.title = `${data.containerInfo.uuid} 运行详情`;
  });

  function onTimePikerOpen() {
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
