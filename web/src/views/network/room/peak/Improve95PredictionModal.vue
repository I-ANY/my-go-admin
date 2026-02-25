<template>
  <BasicModal v-bind="$attrs" @register="registerModal" @cancel="onModalCancel"
    ><BasicTable @register="registerTable" :dataSource="data.items"
  /></BasicModal>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { Improve95PredictionColumns } from './data';
  import { reactive } from 'vue';

  defineOptions({ name: 'Improve95PredictionModal' });
  const emit = defineEmits(['register', 'success']);
  let data = reactive({
    items: [] as any[],
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (d) => {
    const { title, improve95Prediction } = d;
    setModalProps({
      title: title,
      width: 1000,
      height: 450,
      destroyOnClose: true,
      showCancelBtn: false,
      showOkBtn: false,
      canFullscreen: false,
    });
    data.items = improve95Prediction || [];
  });

  const [registerTable] = useTable({
    // title: '95值拉高预测',
    columns: Improve95PredictionColumns(),
    size: 'small',
    canResize: true,
    scroll: { y: 400 },
    useSearchForm: false,
    showTableSetting: false,
    bordered: true,
    showIndexColumn: false,
    rowKey: 'id',
    pagination: false,
  });
  function onModalCancel() {
    closeModal();
    emit('success');
  }
</script>
