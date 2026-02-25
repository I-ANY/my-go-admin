<template>
  <BasicModal v-bind="$attrs" @register="registerModal"
    ><BasicTable @register="registerTable"
  /></BasicModal>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { GetDeliveryInfoMacList } from '@/api/business/k';
  import { getDeliveryInfoMacColumns, getDeliveryInfoMacSearchForm } from './data';

  defineOptions({ name: 'DeliveryInfoMacModal' });
  defineEmits(['register']);
  let record: any = {};

  const [registerModal, { setModalProps }] = useModalInner(async (data) => {
    record = data.record;
    setModalProps({
      title: `${record.hostname} Mac信息`,
      width: 800,
      height: 500,
      destroyOnClose: true,
      showCancelBtn: false,
      showOkBtn: false,
    });
  });

  const [registerTable] = useTable({
    title: 'Mac信息列表',
    api: GetDeliveryInfoMacList,
    columns: getDeliveryInfoMacColumns(),
    beforeFetch: (params) => {
      params.deliveryInfoId = record.id;
      return params;
    },
    canResize: true,
    scroll: { y: 250 },
    formConfig: {
      labelWidth: 120,
      schemas: getDeliveryInfoMacSearchForm(),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: true,
    rowKey: 'id',
    pagination: {
      // pageSizeOptions: ['1', '2', '5'],
    },
    clickToRowSelect: false,
    showSelectionBar: false,
  });
</script>
