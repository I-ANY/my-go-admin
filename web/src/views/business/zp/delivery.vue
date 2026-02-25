<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="batchHandleEdit"> 批量编辑</a-button>
        <a-button type="primary" :loading="deviceSyncEcdnBwLoading" @click="handleSyncEcdnBw">
          带宽同步
        </a-button>
        <a-button type="primary" :loading="batchSubmitLoading" @click="batchHandleSubmit">
          批量提交
        </a-button>
        <a-button type="primary" :loading="deviceStatusRefreshLoading" @click="deviceStatusRefresh">
          状态更新
        </a-button>
      </template>
      <template #bodyCell="{ record, column }">
        <template v-if="column.key === 'network'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '300px', height: '300px' }">
            <template #title>
              <div class="network-content">{{ formatToJson(record.network) }}</div>
            </template>
            <span>{{ record.network }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key === 'storage'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '300px', height: '300px' }">
            <template #title>
              <div class="network-content">{{ formatToJson(record.storage) }}</div>
            </template>
            <span>{{ record.storage }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'clarity:note-edit-line',
                onClick: handleEdit.bind(null, record),
                tooltip: '编辑',
              },
              {
                icon: '',
                label: '提交',
                popConfirm: {
                  title: '是否提交？',
                  confirm: handleSubmit.bind(null, record),
                },
                loading: submitLoading[record.id],
                disabled: !canSubmit(record),
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <DeliveryInfoModal @register="registerDeliveryModal" @success="handleSuccess" />
    <DeliveryBatchEditModal @register="registerDeliveryBatchEditModal" @success="handleSuccess" />
  </div>
</template>
<script setup lang="ts">
  import { BasicTable, TableAction, useTable } from '@/components/Table';
  import {
    DeliveryDeviceList,
    DeliveryDeviceStatusRefresh,
    DeliveryDeviceSubmit,
    DeliveryDeviceSyncEcdnBw,
  } from '@/api/business/zp';
  import { deliveryDeviceColumns, deliveryDeviceSearchFormSchema } from '@/views/business/zp/data';
  import { splitByLineAndTrim } from '@/utils/util';
  import { reactive, ref } from 'vue';
  import DeliveryInfoModal from '@/views/business/zp/deliveryModal.vue';
  import DeliveryBatchEditModal from '@/views/business/zp/deliveryBatchEditModal.vue';
  import { useModal } from '@/components/Modal';
  import { Tooltip, message, Modal } from 'ant-design-vue';

  const [registerDeliveryModal, { openModal: openDeliveryInfoModal }] = useModal();
  const [registerDeliveryBatchEditModal, { openModal: openDeliveryBatchEditModal }] = useModal();
  let batchSubmitLoading = ref(false);
  let deviceStatusRefreshLoading = ref(false);
  let deviceSyncEcdnBwLoading = ref(false);
  let submitLoading = reactive({});

  // const data = reactive({
  //   exporting: false,
  //   exportButTitle: '导出数据',
  // });

  const [registerTable, { reload, getSelectRows, setSelectedRows }] = useTable({
    title: '设备列表',
    api: DeliveryDeviceList,
    columns: deliveryDeviceColumns,
    formConfig: {
      labelWidth: 120,
      schemas: deliveryDeviceSearchFormSchema(),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
      actionColOptions: {
        span: 5,
      },
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    clickToRowSelect: false,
    rowSelection: {
      type: 'checkbox',
    },
    showSelectionBar: true,
    pagination: {
      pageSizeOptions: ['10', '30', '50', '100', '500', '1000'],
    },
    actionColumn: {
      width: 120,
      title: '操作',
      dataIndex: 'action',
      fixed: 'right',
    },
    beforeFetch: (params) => {
      params.provider_device_ids = splitByLineAndTrim(params.provider_device_ids) || null;
      params.device_ids = splitByLineAndTrim(params.device_ids) || null;
      setSelectedRows([]);
    },
  });

  async function handleEdit(record: Recordable) {
    openDeliveryInfoModal(true, {
      record,
      isUpdate: true,
    });
  }

  function canSubmit(record: Recordable): boolean {
    return record.status === 0 || record.status === 3;
  }

  async function handleSubmit(record) {
    submitLoading[record.id] = true;
    try {
      await DeliveryDeviceSubmit({ ids: [record.id] });
      message.success('操作成功');
      await reload();
    } finally {
      submitLoading[record.id] = false;
    }
  }

  async function batchHandleSubmit() {
    batchSubmitLoading.value = true;
    let rows = getSelectRows();
    if (!rows.length) {
      message.warning('请选择设备');
      batchSubmitLoading.value = false;
      return;
    }
    const ids = rows.map((item) => item.id);
    Modal.confirm({
      title: '确认操作',
      content: `确定要提交选中的 ${ids.length} 台设备吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await DeliveryDeviceSubmit({ ids: ids });
          message.success('操作成功');
          await reload();
        } finally {
          batchSubmitLoading.value = false;
        }
      },
      onCancel: () => {
        batchSubmitLoading.value = false;
      },
    });
  }

  async function deviceStatusRefresh() {
    deviceStatusRefreshLoading.value = true;
    try {
      await DeliveryDeviceStatusRefresh({});
      message.success('状态更新成功');
    } finally {
      deviceStatusRefreshLoading.value = false;
      await reload();
    }
  }

  async function batchHandleEdit() {
    let rows = getSelectRows();
    if (!rows.length) {
      message.warning('请选择设备');
      return;
    }
    const ids = rows.map((item) => item.id);
    openDeliveryBatchEditModal(true, {
      ids: ids,
      isUpdate: true,
    });
  }

  async function handleSyncEcdnBw() {
    deviceSyncEcdnBwLoading.value = true;
    let rows = getSelectRows();
    try {
      if (!rows.length) {
        message.warning('请选择设备');
        return;
      }
      const ids = rows.map((item) => item.id);
      await DeliveryDeviceSyncEcdnBw({ ids: ids });
      message.success('操作成功');
      await reload();
    } finally {
      deviceSyncEcdnBwLoading.value = false;
    }
  }

  async function handleSuccess() {
    await reload();
  }

  function formatToJson(Data) {
    try {
      // 其他情况返回原始数据
      return JSON.parse(Data);
    } catch (e) {
      // 解析失败时返回原始数据
      return Data;
    }
  }
</script>
<style lang="less" scoped>
  .network-content {
    max-height: 300px;
    overflow: auto;
    word-break: break-all;
    white-space: pre-wrap;
  }
</style>
