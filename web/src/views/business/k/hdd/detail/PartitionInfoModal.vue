<template>
  <BasicModal v-bind="$attrs" @register="registerModal"
    ><BasicTable @register="registerTable">
      <template #bodyCell="{ column, record: r }">
        <template v-if="column.dataIndex === 'capacity'">
          {{ byte2GB(r.capacity) }}
        </template>
        <template v-if="column.dataIndex === 'businessMountStatus'">
          <Tag
            v-if="businessMountStatusMap[r.businessMountStatus]"
            :color="businessMountStatusMap[r.businessMountStatus].color || 'default'"
          >
            {{ businessMountStatusMap[r.businessMountStatus].dictLabel || r.businessMountStatus }}
          </Tag>
        </template>
      </template>
    </BasicTable></BasicModal
  >
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import {
    businessMountStatusMap,
    getPartitionInfoColumns,
    getPartitionInfoSearchForm,
  } from './data';
  import { GetHddPartitionList } from '@/api/business/k';
  import { Tag } from 'ant-design-vue';

  defineOptions({ name: 'KHDDPartitionInfoModal' });
  defineEmits(['register']);
  let record: any = undefined;
  let formValue: any = undefined;

  const [registerModal, { setModalProps }] = useModalInner(async (data) => {
    record = data.record;
    formValue = data.formValue;
    setModalProps({
      title: `${record.hostname}分区信息`,
      width: 1200,
      height: 500,
      destroyOnClose: true,
      showCancelBtn: false,
      showOkBtn: false,
    });
  });

  const [registerTable] = useTable({
    title: '分区信息',
    api: GetHddPartitionList,
    columns: getPartitionInfoColumns(),
    beforeFetch: (params) => {
      params.batchId = record.batchId;
      if (formValue.minimumSize) {
        params.minimumSize = parseInt(formValue.minimumSize) * 1000 * 1000 * 1000;
      }
      if (formValue.maximumSize) {
        params.maximumSize = parseInt(formValue.maximumSize) * 1000 * 1000 * 1000;
      }
      if (formValue.minUsage) {
        params.minUsage = formValue.minUsage;
      }
      if (formValue.maxUsage) {
        params.maxUsage = formValue.maxUsage;
      }
      // 判断businessMountStatus存在并且大于等于0
      if (formValue.businessMountStatus != null && formValue.businessMountStatus >= 0) {
        params.businessMountStatus = formValue.businessMountStatus;
      }
      return params;
    },
    size: 'small',
    canResize: true,
    scroll: { y: 250 },
    formConfig: {
      labelWidth: 120,
      schemas: getPartitionInfoSearchForm(),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    rowKey: 'id',
    pagination: {
      // pageSizeOptions: ['1', '2', '5'],
    },
    clickToRowSelect: false,
    showSelectionBar: false,
  });
  function byte2GB(byte) {
    if (!byte) return 0;
    return (byte / 1000 / 1000 / 1000).toFixed(2);
  }
</script>
