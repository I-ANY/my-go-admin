<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <Tooltip placement="topLeft" :overlayStyle="{ maxWidth: '800px' }">
          <template #title>从ECDN同步数据，目前只同步带宽相关数据</template>
          <a-button
            type="primary"
            v-show="canBatchDelivery"
            @click="syncEcdnData"
            :loading="syncEcdnButton"
            >同步ECDN数据
          </a-button>
        </Tooltip>
        <a-button
          type="primary"
          v-show="canBatchDelivery"
          @click="submittingDevicesRefresh"
          :loading="refreshButton"
          >状态更新
        </a-button>
        <a-button type="primary" v-show="canBatchDelivery" @click="batchImportDeliveryDevices"
          >批量录入
        </a-button>
        <a-button
          type="primary"
          v-show="canBatchDelivery"
          @click="batchCommitDeliveryDevices"
          :loading="commitButton"
        >
          批量提交
        </a-button>
        <Dropdown>
          <template #overlay>
            <Menu>
              <MenuItem
                key="2"
                v-show="canBatchDelivery"
                @click="batchCloseDeliveryDevices"
                :loading="batchCloseButton"
              >
                批量关闭
              </MenuItem>
              <MenuItem key="3" v-show="canBatchDelivery" @click="delZapPlatformToken">
                刷新Token
              </MenuItem>
            </Menu>
          </template>
          <a-button>
            更多
            <DownOutlined />
          </a-button>
        </Dropdown>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'outer_ips'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '200px' }">
            <template #title>{{ record.outer_ips }}</template>
            <span>{{ record.outer_ips }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key == 'remark'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '300px' }">
            <template #title>{{ record.remark }}</template>
            <span>{{ record.remark }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key == 'order_status'">
          <Tooltip>
            <Tag
              v-if="deliveryStatusMap[record.order_status]"
              :color="deliveryStatusMap[record.order_status].color || 'default'"
              >{{ deliveryStatusMap[record.order_status].dictLabel }}
            </Tag>
          </Tooltip>
        </template>
        <template v-if="column.key == 'submit_type'">
          <span>{{ submitTypeMap[record.submit_type].dictLabel }}</span>
        </template>
        <template v-if="column.key == 'result'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '500px' }">
            <template #title>{{ record.result }}</template>
            <span>{{ record.result }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                label: '编辑',
                disabled: !canEdit(record),
                onClick: handleEditInfo.bind(null, record),
                auth: 'business:z:delivery',
              },
              {
                // color: 'error',
                auth: 'business:z:delivery',
                disabled: !canEnter(record),
                label: '录入',
                tooltip: '录入设备到客户平台',
                onClick: importDeliveryDevice.bind(null, record),
              },
              {
                icon: 'el:ok',
                // color: 'error',
                auth: 'business:z:delivery',
                disabled: !canCommit(record),
                label: '提交',
                tooltip: '交付设备',
                onClick: commitDeliveryDevice.bind(null, record),
                loading: commitTag[record.id],
              },
            ]"
            :dropDownActions="[
              {
                label: '关闭',
                disabled: !canClose(record),
                auth: 'business:z:delivery',
                loading: closeTag[record.id],
                popConfirm: {
                  title: '是否确认关闭交付',
                  placement: 'left',
                  confirm: handleCloseDeliveryDevice.bind(null, record),
                },
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <DeliveryInfoModal @register="registerDeliveryModal" @success="handleSuccess" />
    <ImportModal @register="registerImportModal" @success="handleSuccess" />
  </div>
</template>

<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import {
    DeliveryDeviceCommit,
    DeliveryDeviceSubmittingRefresh,
    GetDeliveryDevice,
    UpdateDeliveryDeviceStatus,
    DelZapPlatformToken,
    SyncEcdnData,
  } from '@/api/business/zap';
  import { message, Tag, Tooltip, MenuItem, Dropdown, Menu, Modal } from 'ant-design-vue';
  import { defineOptions, onMounted, reactive, ref } from 'vue';
  import { usePermission } from '@/hooks/web/usePermission';
  import { columns, searchFormSchema, deliveryStatusMap, submitTypeMap } from './data';
  import { useModal } from '@/components/Modal';
  import DeliveryInfoModal from './DeliveryInfoModal.vue';
  import ImportModal from './ImportModal.vue';
  import { splitByLineAndTrim } from '@/utils/util';
  import { DownOutlined } from '@ant-design/icons-vue';

  const { hasPermission } = usePermission();
  const [registerDeliveryModal, { openModal: openDeliveryInfoModal }] = useModal();
  const [registerImportModal, { openModal: openImportModal }] = useModal();
  defineOptions({ name: 'ZapDevice' });
  let commitButton = ref(false);
  let refreshButton = ref(false);
  let syncEcdnButton = ref(false);
  let batchCloseButton = ref(false);
  let commitTag = reactive({});
  let closeTag = reactive({});

  onMounted(async () => {});

  const canBatchDelivery = (): boolean => {
    return hasPermission('business:z:delivery');
  };

  const [registerTable, { reload, getSelectRows, setSelectedRows }] = useTable({
    title: '交付设备列表',
    api: GetDeliveryDevice,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    rowSelection: {
      type: 'checkbox',
    },
    pagination: {
      // pageSizeOptions: ['10', '30', '50'],
    },
    showSelectionBar: true, // 显示多选状态栏
    clickToRowSelect: false,
    actionColumn: {
      width: 200,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: 'right',
      ifShow: function (): boolean {
        return hasPermission('business:z:delivery');
      },
    },
    beforeFetch: async (params) => {
      params.hostnames = splitByLineAndTrim(params.hostnames) || null;
    },
  });

  enum DeliveryStatus {
    ToEnter = 0, //待录入
    ToCommit = 1, // 待提交
    Delivering = 2, // 交付中
    ICSInit = 3, // ics初始化
    DeployService = 4, // 服务部署
    DeliverySuccess = 5, // 交付完成
    DeliveryFailed = 6, // 交付失败
    DeliveryClose = 8, // 交付关闭
  }

  function canEdit(record: Recordable): boolean {
    return (
      record.order_status == DeliveryStatus.ToEnter ||
      record.order_status == DeliveryStatus.ToCommit ||
      record.order_status == DeliveryStatus.DeliveryFailed ||
      record.order_status == DeliveryStatus.DeliveryClose ||
      record.order_status == DeliveryStatus.DeliverySuccess
    );
  }

  function canEnter(record: Recordable): boolean {
    return (
      record.order_status == DeliveryStatus.ToEnter ||
      record.order_status == DeliveryStatus.DeliveryFailed
    );
  }

  function canCommit(record: Recordable): boolean {
    return (
      record.order_status == DeliveryStatus.ToCommit ||
      record.order_status == DeliveryStatus.DeliveryFailed
    );
  }

  function canClose(record: Recordable): boolean {
    return (
      record.order_status == DeliveryStatus.ToEnter ||
      record.order_status == DeliveryStatus.ToCommit ||
      record.order_status == DeliveryStatus.DeliveryFailed
    );
  }

  // 编辑数据
  function handleEditInfo(record: Recordable) {
    openDeliveryInfoModal(true, {
      record,
      isUpdate: true,
    });
  }

  // 单设备录入
  function importDeliveryDevice(record: Recordable) {
    openImportModal(true, {
      ids: [record.id],
    });
  }

  // 批量录入设备到业务平台
  function batchImportDeliveryDevices() {
    let rows = getSelectRows();
    const ids = rows.map((item) => item.id);
    if (ids.length == 0) {
      message.error('未选择设备!');
      return;
    }
    openImportModal(true, {
      ids: ids,
    });
  }

  // 单设备提交（交付）
  async function commitDeliveryDevice(record: Recordable) {
    commitTag[record.id] = true;
    let params = {
      ids: [record.id],
    };
    try {
      await DeliveryDeviceCommit(params);
      message.success('操作成功!');
      await reload();
    } finally {
      commitTag[record.id] = false;
    }
  }

  // 批量设备提交
  async function batchCommitDeliveryDevices() {
    commitButton.value = true;
    let rows = getSelectRows();
    const ids = rows.map((item) => item.id);
    if (ids.length == 0) {
      commitButton.value = false;
      message.error('未选择设备!');
      return;
    }
    Modal.confirm({
      title: '确认操作',
      content: `确定要录入选中的 ${ids.length} 台设备吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await DeliveryDeviceCommit({ ids: ids });
          message.success('操作成功!');
          setSelectedRows([]);
          await reload();
        } finally {
          commitButton.value = false;
        }
      },
      onCancel: () => {
        commitButton.value = false;
      },
    });
  }

  async function submittingDevicesRefresh() {
    refreshButton.value = true;
    const ids = getSelectRows().map((item) => item.id);
    try {
      await DeliveryDeviceSubmittingRefresh({ ids: ids });
      await reload();
    } finally {
      refreshButton.value = false;
    }
  }

  // 单设备关闭交付
  async function handleCloseDeliveryDevice(record: Recordable) {
    closeTag[record.id] = true;
    try {
      await UpdateDeliveryDeviceStatus({ ids: [record.id] });
      message.success('操作成功!');
      await reload();
    } finally {
      closeTag[record.id] = false;
    }
  }

  // 批量设备关闭交付
  async function batchCloseDeliveryDevices() {
    batchCloseButton.value = true;
    const ids = getSelectRows().map((item) => item.id);
    try {
      await UpdateDeliveryDeviceStatus({ ids: ids });
      batchCloseButton.value = false;
      message.success('操作成功!');
      await reload();
    } finally {
      batchCloseButton.value = false;
    }
  }

  // 重新获取token
  async function delZapPlatformToken() {
    await DelZapPlatformToken();
    message.success('操作成功!');
  }

  function handleSuccess() {
    reload();
  }

  async function syncEcdnData() {
    const ids = getSelectRows().map((item) => item.id);
    if (ids.length == 0) {
      message.error('未选择设备!');
      return;
    }
    syncEcdnButton.value = true;
    try {
      await SyncEcdnData({ ids: ids });
      message.success('操作成功!');
    } finally {
      syncEcdnButton.value = false;
    }
  }
</script>
