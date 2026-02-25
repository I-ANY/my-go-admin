<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" v-auth="'business:k:batchDelivery'" @click="handleBatchDelivery"
          >批量交付</a-button
        >
        <a-button
          type="primary"
          v-auth="'business:k:specialLine:batchBindDemand'"
          @click="handlebatchBindDemand"
          >批量绑定需求单</a-button
        >
        <Dropdown v-if="showDropdown()">
          <template #overlay>
            <Menu>
              <MenuItem
                key="1"
                @click="handleBatchCloseDelivery"
                v-if="hasPermission('business:k:specialLine:delivery:batchClose')"
              >
                批量关闭交付
              </MenuItem>
              <MenuItem
                key="2"
                @click="handleBatchEdit"
                v-if="hasPermission('business:k:specialLine:delivery:batchEdit')"
              >
                批量编辑
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
        <template v-if="column.key == 'deliveryStatus'">
          <Tooltip>
            <Tag
              v-if="deliveryInfoMap[record.deliveryStatus]"
              :color="deliveryInfoMap[record.deliveryStatus].color || 'default'"
              >{{ deliveryInfoMap[record.deliveryStatus].dictLabel }}</Tag
            >
            <template #title v-if="record.remark">
              <span>{{ record.remark }}</span>
            </template>
          </Tooltip>
        </template>
        <template v-if="column.key == 'deliveryType'">
          <Tooltip>
            <Tag
              v-if="deliveryTypeMap[record.deliveryType]"
              :color="deliveryTypeMap[record.deliveryType].color || 'default'"
              >{{ deliveryTypeMap[record.deliveryType].dictLabel }}</Tag
            >
            <template #title v-if="record.remark">
              <span>{{ record.remark }}</span>
            </template>
          </Tooltip>
        </template>

        <template v-if="column.key == 'isCoverDiffIsp'">
          <span v-if="isCoverDiffIspMap[record.isCoverDiffIsp]">{{
            isCoverDiffIspMap[record.isCoverDiffIsp].dictLabel
          }}</span>
        </template>
        <template v-if="column.key == 'isProvinceScheduling'">
          <span v-if="isProvinceSchedulingMap[record.isProvinceScheduling]">{{
            isProvinceSchedulingMap[record.isProvinceScheduling].dictLabel
          }}</span>
        </template>
        <template v-if="column.key == 'bizType'">
          <span v-if="bizTypesMap[record.bizType]">{{
            bizTypesMap[record.bizType].dictLabel
          }}</span>
        </template>
        <template v-if="column.key == 'provider'">
          <span v-if="providersMap[record.provider]">{{
            providersMap[record.provider].dictLabel
          }}</span>
        </template>
        <template v-if="column.key == 'deviceType'">
          <span v-if="deviceTypeMap[record.deviceType]">{{
            deviceTypeMap[record.deviceType].dictLabel
          }}</span>
          <span v-else></span>
        </template>
        <template v-if="column.key == 'remark'">
          <Tooltip>
            {{ record.remark }}
            <template #title v-if="record.remark">
              <span>{{ record.remark }}</span>
            </template>
          </Tooltip>
        </template>
        <template v-if="column.key == 'bwCount'">
          <Tooltip>
            <a style="display: block" @click="handleClickBwCount(record)">{{ record.bwCount }}</a>
            <template #title>
              <!-- <div v-for="mac in record.macs" :key="mac.id">{{ mac.mac }}</div> -->
              <span>查看详情</span>
            </template>
          </Tooltip>
        </template>
        <template v-if="column.key == 'demandId'">
          <Tooltip>
            <Badge color="red" v-if="record.bindError">
              <span
                :style="{
                  display: 'block',
                  marginRight: '6px',
                  color: record.bindError == true ? 'red' : 'green',
                }"
                >{{ record.demandId }}
              </span>
              <template #count>
                <ExclamationCircleFilled style="color: red" />
              </template>
            </Badge>
            <div v-else
              ><span>{{ record.demandId }}</span></div
            >
            <template #title v-if="record.bindError">
              <span>{{ record.errorMsg }}</span>
            </template>
          </Tooltip>
        </template>
        <!-- <template v-if="column.key == 'deliveryBw'">
          {{ record.deliveryBw / 1000 || '' }}
        </template> -->
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'el:ok',
                // color: 'error',
                auth: 'business:k:delivery',
                disabled: !canDelivery(record),
                label: '交付',
                tooltip: '交付设备',
                onClick: handleDelivery.bind(null, record),
              },
              {
                // icon: 'clarity:note-edit-line',
                tooltip: '绑定需求单号',
                disabled: !canEdit(record),
                label: '绑定需求单',
                onClick: handleEditDemand.bind(null, record),
                auth: 'business:k:demand:edit',
              },
            ]"
            :dropDownActions="[
              {
                label: '编辑',
                disabled: !canEdit(record),
                onClick: handleEditInfo.bind(null, record),
                auth: 'business:k:info:edit',
              },
              {
                label: '异网跨省下发',
                disabled: !canDifIsp(record),
                onClick: handleDifIsp.bind(null, record),
                auth: 'business:k:specialLine:difIsp',
              },
              {
                label: '业务接入',
                disabled: canBusinessJoin(record),
                auth: 'business:k:delivery:businessJoin',
                popConfirm: {
                  title: '是否确认重新接入',
                  placement: 'left',
                  confirm: handleBusinessJoin.bind(null, record),
                },
              },
              {
                label: '关闭交付',
                disabled: !canEdit(record),
                auth: 'business:k:specialLine:delivery:close',
                popConfirm: {
                  title: '是否确认关闭交付',
                  placement: 'left',
                  confirm: handleCloseDeliveryInfo.bind(null, record),
                },
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <SpecialLineDemandSelectModal
      @register="registerModal"
      @success="handleSuccess"
      :biz-type="bizType"
    />
    <SpecialLineDeliveryInfoModal
      @register="registerDeliveryModal"
      @success="handleSuccess"
      :biz-type="bizType"
    />
    <DeliveryInfoMacModal @register="registerDeliveryInfoMacModal" />
    <SpecialLineBatchDemandModal
      @register="registerSpecialLineBatchDemandModal"
      @success="handleSuccess"
      :biz-type="bizType"
    />
    <DifIspModal :biz-type="bizType" @register="registerDifIspModal" @success="handleSuccess" />
    <SpecialLineBatchEditModal
      @register="registerBatchEditModal"
      @success="handleSuccess"
      :biz-type="bizType"
    />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';

  import {
    specialLineColumns,
    specialLineFormSchema,
    deliveryInfoMap,
    deliveryTypeMap,
    DiveryStatus,
    bizTypesMap,
    isProvinceSchedulingMap,
    isCoverDiffIspMap,
    providersMap,
    deviceTypeMap,
  } from './data';
  import {
    BatchBindDemandCheck,
    BatchUpdateDeliveryStatus,
    getDeliveryInfoList,
    deliveryInfoDelivery,
    UpdateDeliveryInfoStatus,
    BusinessJoin,
  } from '@/api/business/k';
  import { usePermission } from '@/hooks/web/usePermission';
  import { Tag, Modal, message, Tooltip, Badge, Dropdown, MenuItem, Menu } from 'ant-design-vue';
  import { h, nextTick } from 'vue';
  import {
    ExclamationCircleOutlined,
    ExclamationCircleFilled,
    DownOutlined,
  } from '@ant-design/icons-vue';
  import SpecialLineDemandSelectModal from './SpecialLineDemandSelectModal.vue';
  import SpecialLineBatchDemandModal from './SpecialLineBatchDemandModal.vue';
  import SpecialLineDeliveryInfoModal from './SpecialLineDeliveryInfoModal.vue';
  import DifIspModal from './DifIspModal.vue';
  import DeliveryInfoMacModal from './DeliveryInfoMacModal.vue';
  import SpecialLineBatchEditModal from './SpecialLineBatchEditModal.vue';
  import { useModal } from '@/components/Modal';
  import { BizType } from '../data';
  import { splitByLineAndTrim } from '@/utils/util';
  import { useMessage } from '@/hooks/web/useMessage';

  const [registerModal, { openModal }] = useModal();
  const [registerDeliveryModal, { openModal: openDeliverModal }] = useModal();
  const [registerSpecialLineBatchDemandModal, { openModal: openSpecialLineBatchDemandModal }] =
    useModal();
  const [registerDeliveryInfoMacModal, { openModal: openDeliveryInfoMacModal }] = useModal();
  const [registerDifIspModal, { openModal: openDifIspModal }] = useModal();
  const [registerBatchEditModal, { openModal: openBatchEditModal }] = useModal();
  defineOptions({ name: 'KSpecialLineDelivery' });
  const { hasPermission } = usePermission();
  const { notification } = useMessage();
  const bizType = BizType.specialLine;
  const [registerTable, { getSelectRowKeys, reload, setSelectedRowKeys, getSelectRows }] = useTable(
    {
      title: '设备交付',
      api: getDeliveryInfoList,
      columns: specialLineColumns,
      formConfig: {
        labelWidth: 120,
        schemas: specialLineFormSchema(),
        autoSubmitOnEnter: true,
        showAdvancedButton: true,
        autoAdvancedLine: 4,
        // alwaysShowLines: 2,
        alwaysShowLines: 1,
      },
      useSearchForm: true,
      showTableSetting: true,
      bordered: true,
      showIndexColumn: false,
      pagination: {
        // pageSizeOptions: ['1', '2', '5'],
      },
      beforeFetch: (params) => {
        params.bizType = bizType;
        params.hostnames = splitByLineAndTrim(params.hostnames);
        params.macAddrs = splitByLineAndTrim(params.macAddrs);
        params.owners = splitByLineAndTrim(params.owners);
        return params;
      },
      rowKey: 'id',
      rowSelection: {
        type: 'checkbox',
        getCheckboxProps: (record) => {
          return {
            disabled: !canChecked(record),
          };
        },
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
          return (
            hasPermission('business:k:delivery') ||
            hasPermission('business:k:demand:edit') ||
            hasPermission('business:k:info:edit') ||
            hasPermission('business:k:specialLine:difIsp') ||
            hasPermission('business:k:specialLine:delivery:close') ||
            hasPermission('business:k:specialLine:delivery:businessJoin')
          );
        },
      },
    },
  );
  function handleEditDemand(record: Recordable) {
    openModal(true, {
      record,
      isUpdate: true,
    });
  }

  async function handleDelivery(record: Recordable) {
    const res = await getDeliveryInfoList({ ids: [record.id], bizType: bizType });
    if (!res.items || res.items.length === 0) {
      message.error({
        content: '未找到交付记录',
      });
      return;
    }
    const deliveryInfo = res.items[0];

    //  交付总带宽
    let deliveryBw = deliveryInfo.totalBw;
    // 单线带宽
    let singleBw = deliveryInfo.uploadBw;
    if (deliveryInfo.deliveryBw && deliveryInfo.deliveryBw > 0) {
      deliveryBw = deliveryInfo.deliveryBw;
      if (deliveryInfo.bwCount > 0) {
        singleBw = (deliveryInfo.deliveryBw / deliveryInfo.bwCount).toFixed(0);
      } else {
        singleBw = '';
      }
    }

    Modal.confirm({
      title: '是否确认交付该设备?',
      icon: h(ExclamationCircleOutlined),
      content: h('div', [
        h(
          'p',
          { style: 'color:black;margin:2px;font-weight: bold;' },
          '主机名：' + deliveryInfo.hostname,
        ),
        h(
          'p',
          { style: 'color:red;margin:2px;font-weight: bold;' },
          '交付总带宽：' + deliveryBw + 'Mbps',
        ),
        h(
          'p',
          { style: 'color:red;margin:2px;font-weight: bold;' },
          '线路数：' + deliveryInfo.bwCount,
        ),
        h(
          'p',
          { style: 'color:red;margin:2px;font-weight: bold;' },
          '交付单线单宽：' + singleBw + 'Mbps',
        ),
      ]),
      maskClosable: true,
      async onOk() {
        const data = {
          ids: [record.id],
          bizType: bizType,
        };
        await deliveryInfoDelivery(data);
        message.success({ content: '提交成功' });
        reload();
        setSelectedRowKeys([]);
      },
      class: 'isDelivery',
    });
  }
  // async function handleDelivery(record: Recordable) {
  //   const data = {
  //     ids: [record.id],
  //     bizType: bizType,
  //   };
  //   await deliveryInfoDelivery(data);
  //   message.success({ content: '提交成功' });
  //   reload();
  //   setSelectedRowKeys([]);
  // }
  function handleBatchDelivery() {
    const selectKeys = getSelectRowKeys();
    if (selectKeys.length === 0) {
      message.warning({ content: '请选择要交付的设备', key: 'batchDelivery:noSelect' });
      return;
    }

    // 获取选中行的数据
    const selectedRows = getSelectRows();
    const hostnames = selectedRows.map((row) => row.hostname).join('\n');

    Modal.confirm({
      title: '是否确认批量交付设备?',
      icon: h(ExclamationCircleOutlined),
      content: h('div', [
        h('p', { style: 'color:red;margin-bottom: 8px;' }, '共' + selectKeys.length + '个设备'),
        h('div', { style: 'color:black;overflow-y: auto;white-space: pre-line;' }, hostnames),
      ]),
      async onOk() {
        const data = {
          ids: selectKeys,
          bizType: bizType,
        };
        await deliveryInfoDelivery(data);
        message.success({ content: '提交成功' });
        reload();
        setSelectedRowKeys([]);
      },
      class: 'isDelivery',
    });
  }
  function handleSuccess() {
    reload();
  }
  function handleEditInfo(record: Recordable) {
    openDeliverModal(true, {
      record,
      isUpdate: true,
    });
  }
  function handleClickBwCount(record: Recordable) {
    openDeliveryInfoMacModal(true, { record });
  }
  function canDelivery(record: Recordable): boolean {
    return (
      (record.deliveryStatus == DiveryStatus.DeliveryFailed ||
        record.deliveryStatus == DiveryStatus.DeploySuccess ||
        record.deliveryStatus == DiveryStatus.DifIsping ||
        record.deliveryStatus == DiveryStatus.DifIspSuccess ||
        record.deliveryStatus == DiveryStatus.DifIspFailed ||
        record.deliveryStatus == DiveryStatus.JoinSuccess ||
        record.deliveryStatus == DiveryStatus.JoinFailed) &&
      record.demandId &&
      !record.bindError &&
      record.macs &&
      record.macs.length > 0 &&
      record.bwCount > 0 &&
      record.taskId == null // 交付数据有ID 则不能再提交交付
    );
  }
  function canEdit(record: Recordable): boolean {
    return (
      record.deliveryStatus == DiveryStatus.DeliveryFailed ||
      record.deliveryStatus == DiveryStatus.DeploySuccess ||
      record.deliveryStatus == DiveryStatus.DifIspFailed ||
      record.deliveryStatus == DiveryStatus.DifIspSuccess ||
      record.deliveryStatus == DiveryStatus.DifIsping ||
      record.deliveryStatus == DiveryStatus.JoinSuccess ||
      record.deliveryStatus == DiveryStatus.JoinFailed
    );
  }
  async function handleCloseDeliveryInfo(record: Recordable) {
    const data = {
      bizType: bizType,
      deliveryStatus: DiveryStatus.DeliveryClosed,
    };
    await UpdateDeliveryInfoStatus(record.id, data);
    message.success({ content: '操作成功' });
    reload();
  }
  function canChecked(record: Recordable): boolean {
    return canDelivery(record) || canEdit(record);
  }
  function handleBatchCloseDelivery() {
    const selectKeys = getSelectRowKeys();
    if (selectKeys.length === 0) {
      message.warning({ content: '请选择要关闭交付的设备' });
      return;
    }
    Modal.confirm({
      title: '是否确认批量关闭交付?',
      icon: h(ExclamationCircleOutlined),
      content: h('div', { style: 'color:red;' }, '共' + selectKeys.length + '个设备'),
      async onOk() {
        const data = {
          ids: selectKeys,
          bizType: bizType,
          deliveryStatus: DiveryStatus.DeliveryClosed,
        };
        const res = await BatchUpdateDeliveryStatus(data);
        if (res?.errorMessage) {
          notification.error({
            message: '批量关闭交付失败',
            description: res.errorMessage,
            duration: null,
          });
          nextTick(() => {
            setSelectedRowKeys([]);
          });
          return;
        }
        message.success({ content: '操作成功' });
        reload();
        nextTick(() => {
          setSelectedRowKeys([]);
        });
      },
    });
  }
  function showDropdown() {
    return hasPermission('business:k:specialLine:delivery:batchClose');
  }
  async function handlebatchBindDemand() {
    const selectKeys = getSelectRowKeys();
    if (selectKeys.length === 0) {
      message.warning({ content: '请选择要绑定需求单的设备' });
      return;
    }
    const res = await BatchBindDemandCheck({ ids: selectKeys, bizType: bizType });
    if (!res?.success) {
      message.error({ content: res?.message });
      return;
    } else if (res.message) {
      message.warning({ content: res?.message, duration: 10, key: 'batchBindDemand:warning' });
    }
    const data = { ids: selectKeys, bizType: bizType };
    openSpecialLineBatchDemandModal(true, data);
  }

  function canDifIsp(record: Recordable) {
    return (
      record.deliveryStatus == DiveryStatus.DeploySuccess ||
      record.deliveryStatus == DiveryStatus.DifIspFailed ||
      record.deliveryStatus == DiveryStatus.DeliveryFailed ||
      record.deliveryStatus == DiveryStatus.DifIspSuccess ||
      record.deliveryStatus == DiveryStatus.DifIsping
    );
  }
  function handleDifIsp(record: Recordable) {
    openDifIspModal(true, {
      record,
      bizType: bizType,
    });
  }

  // 批量编辑功能
  async function handleBatchEdit() {
    const selectKeys = getSelectRowKeys();
    if (selectKeys.length === 0) {
      message.warning({ content: '请选择要编辑的设备' });
      return;
    }

    const selectedRows = getSelectRows();

    // 检查选中的设备是否可以编辑
    const canEditDevices = selectedRows.filter((record) => canEdit(record));

    if (canEditDevices.length === 0) {
      message.warning({ content: '选中的设备都不允许编辑' });
      return;
    }

    // 如果有不可编辑的设备，给出提示
    if (canEditDevices.length < selectedRows.length) {
      message.warning({
        content: `只有${canEditDevices.length}个设备可以编辑，${selectedRows.length - canEditDevices.length}个设备状态不允许编辑`,
      });
    }

    openBatchEditModal(true, {
      selectedDevices: canEditDevices,
    });
  }

  function canBusinessJoin(record: Recordable) {
    return (
      record.deliveryStatus == DiveryStatus.JoinSuccess ||
      record.deliveryStatus == DiveryStatus.DeliveryClosed
    );
  }

  async function handleBusinessJoin(record: Recordable) {
    try {
      const params = {
        biz_type: bizType,
        provider: record.provider,
        id: record.id,
        hostname: record.hostname,
      };
      const response = await BusinessJoin(params);

      // 检查response和response.data是否存在
      if (!response || !response.data) {
        notification.error({
          message: '操作失败',
          description: '服务器响应异常',
          duration: 5,
          placement: 'top',
        });
        return;
      }

      const { code, msg } = response.data;

      if (code === 200) {
        notification.success({
          message: '操作成功',
          description: msg,
          duration: 5,
          placement: 'top',
        });
      } else {
        notification.error({
          message: '操作失败',
          description: msg,
          duration: 5,
          placement: 'top',
        });
      }
      reload();
    } catch (error) {
      console.error('handleBusinessJoin error:', error);
      notification.error({
        message: '操作失败',
        description: '请求异常，请稍后重试',
        duration: 5,
        placement: 'top',
      });
    }
  }
</script>
