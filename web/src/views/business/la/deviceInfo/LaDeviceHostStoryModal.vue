<template>
  <BasicModal v-bind="$attrs" @register="register">
    <BasicTable @register="registerTable" />
  </BasicModal>
</template>
<script lang="ts" setup>
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicTable, useTable } from '@/components/Table';
  import { hostHistoryColumns, hostHistorySearchFormSchema } from './data';
  import { GetDeviceHostHistoryList } from '@/api/business/la';
  import { reactive } from 'vue';

  defineOptions({ name: 'LaDeviceHostStoryModal' });

  const data = reactive({
    deviceId: '',
  });

  const [register, { setModalProps }] = useModalInner(({ record }) => {
    data.deviceId = record.deviceId;
    setModalProps({
      title: `设备 ${record.deviceId} 历史主机名`,
      width: 1000,
      // height: 500,
      destroyOnClose: true,
      showCancelBtn: false,
      showOkBtn: false,
    });
  });
  const [registerTable] = useTable({
    title: '历史主机名',
    api: GetDeviceHostHistoryList,
    beforeFetch: (params) => {
      params.deviceId = data.deviceId;
      return params;
    },
    size: 'small',
    canResize: true,
    scroll: { y: 300 },
    columns: hostHistoryColumns(),
    formConfig: {
      labelWidth: 120,
      schemas: hostHistorySearchFormSchema(),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 3,
      submitOnReset: false,
      alwaysShowLines: 1,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    rowKey: 'id',
  });
</script>
