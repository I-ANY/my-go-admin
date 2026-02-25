<template>
  <BasicModal v-bind="$attrs" @register="register">
    <BasicTable @register="registerTable" />
  </BasicModal>
</template>
<script lang="ts" setup>
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicTable, useTable } from '@/components/Table';
  import { slaAbnormalHostColumns, slaAbnormalHostSearchFormSchema } from './data';
  import { GetAbnormalHostList } from '@/api/business/la';
  import { reactive } from 'vue';
  import dayjs from 'dayjs';

  defineOptions({ name: 'LaSLAAbnormalHostModal' });

  const data = reactive({
    slaInfoId: '',
  });

  const [register, { setModalProps }] = useModalInner(({ record }) => {
    data.slaInfoId = record.id;
    setModalProps({
      title: `${record.province}-${record.isp}-${dayjs(record.time)?.format('YYYY-MM-DD HH:mm:ss')} 异常设备信息`,
      width: 1000,
      // height: 600,
      destroyOnClose: true,
      showCancelBtn: false,
      showOkBtn: false,
    });
  });
  const [registerTable] = useTable({
    title: '异常设备',
    api: GetAbnormalHostList,
    beforeFetch: (params) => {
      params.slaInfoId = data.slaInfoId;
      return params;
    },
    size: 'small',
    canResize: true,
    scroll: { y: 300 },
    columns: slaAbnormalHostColumns(),
    formConfig: {
      labelWidth: 120,
      schemas: slaAbnormalHostSearchFormSchema(),
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
