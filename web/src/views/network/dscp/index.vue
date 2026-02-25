<template>
  <div class="network-dscp-summary-container">
    <BasicTable @register="registerTable" @change="handleTableChange">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="NetworkPermissionCodeEnum.BIZ_DSCP_MODIFY"
          @click="handleModifyDscp"
          >修改DSCP值</a-button
        >
        <a-button
          type="primary"
          v-auth="NetworkPermissionCodeEnum.BIZ_DSCP_EXPORT"
          :loading="data.exporting"
          @click="handleExportBizDscp"
          >{{ data.exportButTitle }}</a-button
        >
      </template>
      <template #form-dscp="{ model, field }">
        <InputNumber
          v-model:value="model[field]"
          placeholder="DSCP值"
          :precision="0"
          :min="0"
          :max="63"
          :step="1"
        >
          <!-- 全局、同网省内、同网省外、异网省内、异网省 -->
          <template #addonAfter>
            <Select
              v-model:value="model.dscpType"
              allowClear
              :options="[
                { label: '全局DSCP ', value: 1 },
                { label: '同网省内DSCP', value: 2 },
                { label: '同网省外DSCP', value: 3 },
                { label: '异网省内DSCP', value: 4 },
                { label: '异网省外DSCP', value: 5 },
              ]"
              placeholder="选择DSCP类型"
            />
          </template>
        </InputNumber>
      </template>
      <template #tableTop>
        <div style="margin-bottom: 8px; padding: 8px; background-color: white">
          <div style="display: flex; align-items: center; gap: 4px">
            <span style="color: #4ebef1; white-space: nowrap">当前节点业务全局优先级：</span
            ><span>{{ data.businessPriority || '-' }}</span></div
          >
        </div>
      </template>
      <template #bodyCell="{ column, record }">
        <!-- 使用枚举值展示 -->
        <template v-for="(columnName, i) in Object.keys(data.showEnumFields)" :key="i">
          <template v-if="column.key == columnName">
            <span v-if="data.showEnumFields[columnName][record[columnName]]">
              {{ data.showEnumFields[columnName][record[columnName]].dictLabel }}
            </span>
            <span v-else>{{ record[columnName] }}</span>
          </template>
        </template>

        <!-- 使用枚举值tag展示 -->
        <template v-for="(columnName, i) in Object.keys(data.showTagFields)" :key="i">
          <template v-if="column.key == columnName">
            <Tag
              style="font-weight: bold"
              v-if="data.showTagFields[columnName][record[columnName]]"
              :color="data.showTagFields[columnName][record[columnName]].color || 'default'"
              >{{ data.showTagFields[columnName][record[columnName]].dictLabel }}</Tag
            >
            <span v-else>{{ record[columnName] }}</span>
          </template>
        </template>
        <template v-if="column.key === 'business'">
          <Tooltip
            title="查看详情"
            v-if="hasPermission(NetworkPermissionCodeEnum.SERVER_DSCP_DETAIL)"
          >
            <a @click="handleViewDscpDetail(record)">{{ record.business }}</a>
          </Tooltip>
          <span v-else>{{ record.business }}</span>
        </template>

        <!-- <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                tooltip: '查看详情',
                onClick: handleViewDscpDetail.bind(null, record),
                label: '详情',
              },
            ]"
          />
        </template> -->
      </template>
    </BasicTable>
    <!-- 设备详情 -->
    <ServerPriorityModal
      @register="registerServerPriorityModal"
      @success="
        () => {
          reload();
        }
      "
    />
    <!-- 修改DSCP值 -->
    <ModifyDscpModal
      @register="registerModifyDscpModal"
      @modify-dscp-success="handleModifyDscpSuccess"
      @success="
        () => {
          reload();
        }
      "
    />
    <!-- 修改DSCP值结果 -->
    <ModifyDscpResultModal
      @register="registerModifyDscpResultModal"
      @success="
        () => {
          reload();
        }
      "
    />
  </div>
</template>
<script setup lang="ts">
  import { h, nextTick, reactive } from 'vue';
  import { Select, InputNumber, Tag, message, Tooltip, notification, Modal } from 'ant-design-vue';
  import { BasicTable, PaginationProps, SorterResult, useTable } from '@/components/Table';
  import { GetBizOwnerPriority, GetBizPriorityList, Api } from '@/api/network/dscp';
  import { NetworkPermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import {
    bizPriorityColumns,
    bizPrioritySearchFormSchema,
    ModifyDscpResult,
    modifyType,
  } from './data';
  import { splitByLineAndTrim } from '@/utils/util';
  import { useModal } from '@/components/Modal';
  import ServerPriorityModal from './ServerPriorityModal.vue';
  import ModifyDscpModal from './ModifyDscpModal.vue';
  import ModifyDscpResultModal from './ModifyDscpResultModal.vue';
  import { usePermission } from '@/hooks/web/usePermission';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { downloadFileByUrl } from '@/utils/download';

  const { hasPermission } = usePermission();
  const [registerModifyDscpResultModal, { openModal: openModifyDscpResultModal }] = useModal();

  const [registerServerPriorityModal, { openModal: openServerPriorityModal }] = useModal();
  const [registerModifyDscpModal, { openModal: openModifyDscpModal }] = useModal();
  defineOptions({ name: 'NetworkBusinessPriority' });
  const data = reactive({
    showEnumFields: {},
    showTagFields: {},
    businessPriority: '',
    selectOwner: '',
    exporting: false,
    exportButTitle: '导出数据',
    sorter: { field: null as string | null, order: null as string | null },
  });
  const [registerTable, { getForm, getSelectRowKeys, getSelectRows, reload }] = useTable({
    title: '业务优先级',
    api: GetBizPriorityList,
    columns: bizPriorityColumns(),
    formConfig: {
      labelWidth: 120,
      schemas: bizPrioritySearchFormSchema(),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      submitOnReset: false,
      alwaysShowLines: 1,
    },
    beforeFetch: (params) => {
      formatBusinessPriority();
      params = formatParams(params);
      return params;
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    immediate: true,
    rowKey: 'id',
    size: 'small',
    // actionColumn: {
    //   width: 100,
    //   title: '操作',
    //   dataIndex: 'action',
    //   fixed: 'right',
    // },
    clickToRowSelect: false,
    showSelectionBar: true,
    rowSelection: {
      type: 'checkbox',
      hideSelectAll: true,
      getCheckboxProps: (record) => {
        return {
          disabled: data.selectOwner != '' && data.selectOwner != record.owner,
        };
      },
      onChange: (selectedRowKeys, selectedRows) => {
        if (selectedRows.length > 0) {
          data.selectOwner = selectedRows[0]?.owner || '';
        } else {
          data.selectOwner = '';
        }
      },
    },
  });
  function handleViewDscpDetail(record: Recordable) {
    openServerPriorityModal(true, { record });
  }
  function formatParams(params: Recordable) {
    params.hostname = splitByLineAndTrim(params.hostname);
    params.ecdnIp = splitByLineAndTrim(params.ecdnIp);
    return params;
  }
  async function formatBusinessPriority() {
    let values = (await getForm().getFieldsValue()) || {};
    let params = { ...values };
    params = formatParams(params);
    const { priority } = await GetBizOwnerPriority(params);
    let businessPriority = '';
    for (let i = 0; i < priority?.length || 0; i++) {
      let operator = '';
      let lastBusinessPriority = priority[i - 1];
      let currentBusinessPriority = priority[i];
      if (i > 0) {
        if ((currentBusinessPriority.priority || 0) == (lastBusinessPriority.priority || 0)) {
          operator = ' = ';
        } else {
          operator = ' > ';
        }
      }
      businessPriority += operator + currentBusinessPriority.business;
    }
    data.businessPriority = businessPriority;
  }
  async function handleModifyDscp() {
    const selectedRowKeys = getSelectRowKeys();
    if (!selectedRowKeys || selectedRowKeys.length == 0) {
      message.error({ content: '请选择要修改的业务' });
      return;
    }
    openModifyDscpModal(true, {
      type: modifyType.SUMMARY,
      ids: selectedRowKeys as number[],
      owner: getSelectRows()[0]?.owner,
    });
  }
  function handleModifyDscpSuccess(results: Array<ModifyDscpResult>) {
    openModifyDscpResultModal(true, { results });
  }

  // 导出数据
  function handleExportBizDscp() {
    Modal.confirm({
      title: '是否确认导出筛选的数据?',
      icon: h(ExclamationCircleOutlined),
      content: h(
        'div',
        { style: 'color:red;' },
        '导出将花费一些时间，请耐心等待！如果数据过大，可能会导出失败，请缩小导出条件后重试',
      ),
      async onOk() {
        await getForm().validate();
        let value = await getForm().getFieldsValue();
        // 追加过滤条件和排序条件
        value = { ...value };
        value = formatParams(value);
        if (data?.sorter) {
          value.field = data.sorter.field || null;
          value.order = data.sorter.order || null;
        }

        nextTick(() => {
          data.exporting = true;
          data.exportButTitle = '导出中';
        });
        await exportData(value);
      },
    });
  }
  async function exportData(value: Recordable) {
    try {
      let filename = await downloadFileByUrl(Api.ExportBizPriority(), 'POST', 5 * 60, value, null);
      notification.success({
        message: '导出成功',
        description: '文件名：' + filename,
        duration: null,
      });
    } catch (error) {
      notification.error({
        message: '导出失败',
        description: error.message,
        duration: null,
      });
    } finally {
      nextTick(() => {
        data.exporting = false;
        data.exportButTitle = '导出数据';
      });
    }
  }
  function handleTableChange(
    _pagination: PaginationProps,
    _filters: Recordable,
    sorter: SorterResult,
  ) {
    if (sorter && sorter.field && sorter.order) {
      data.sorter = {
        field: sorter.field,
        order: sorter.order,
      };
    } else {
      data.sorter = {
        field: null,
        order: null,
      };
    }
  }
</script>
<style>
  .network-dscp-summary-container .vben-basic-table.vben-basic-table-form-container form {
    margin-bottom: 8px !important;
  }
</style>
