<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    class="network-server-priority-modal"
    :width="1450"
    :destroy-on-close="true"
    @cancel="handelCancel"
  >
    <BasicTable @register="registerTable" @change="handleTableChange">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="NetworkPermissionCodeEnum.SERVER_DSCP_MODIFY"
          @click="handleModifyDscp"
          >修改DSCP值</a-button
        >
        <a-button
          type="primary"
          v-auth="NetworkPermissionCodeEnum.SERVER_DSCP_EXPORT"
          :loading="data.exporting"
          @click="handleExportServerPriority"
          >{{ data.exportButTitle }}</a-button
        >
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
      </template></BasicTable
    >
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
  </BasicModal>
</template>

<script setup lang="ts">
  import { h, nextTick, reactive } from 'vue';
  import { BasicModal, useModalInner, useModal } from '@/components/Modal';
  import { BasicTable, PaginationProps, SorterResult, useTable } from '@/components/Table';
  import { Api, GetServerPriorityList } from '@/api/network/dscp';
  import { splitByLineAndTrim } from '@/utils/util';
  import { message, Modal, notification, Tag } from 'ant-design-vue';
  import {
    serverPriorityColumns,
    serverPrioritySearchFormSchema,
    dscpAgentStatusMap,
    serverOnlineStatusMap,
    modifyType,
    ModifyDscpResult,
  } from './data';
  import ModifyDscpModal from './ModifyDscpModal.vue';
  import ModifyDscpResultModal from './ModifyDscpResultModal.vue';
  import { NetworkPermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import { downloadFileByUrl } from '@/utils/download';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';

  const [registerModifyDscpModal, { openModal: openModifyDscpModal }] = useModal();
  const [registerModifyDscpResultModal, { openModal: openModifyDscpResultModal }] = useModal();
  defineOptions({ name: 'NetworkServerPriorityModal' });
  const emit = defineEmits(['register', 'success']);

  const data = reactive({
    showEnumFields: {},
    showTagFields: { status: dscpAgentStatusMap, online: serverOnlineStatusMap },
    record: {} as Recordable,
    exporting: false,
    exportButTitle: '导出数据',
    sorter: { field: null as string | null, order: null as string | null },
  });

  const [registerTable, { getSelectRowKeys, reload, getForm }] = useTable({
    title: '设备详情',
    api: GetServerPriorityList,
    columns: serverPriorityColumns(),
    formConfig: {
      labelWidth: 120,
      schemas: serverPrioritySearchFormSchema(),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      submitOnReset: false,
      alwaysShowLines: 1,
    },
    beforeFetch: (params) => {
      params = formatParams(params);
      params.summaryId = data.record.id;
      return params;
    },
    scroll: { y: '340px' },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    immediate: true,
    rowKey: 'id',
    size: 'small',
    clickToRowSelect: false,
    showSelectionBar: true,
    rowSelection: {
      type: 'checkbox',
    },
  });

  const [registerModal, { setModalProps }] = useModalInner(async (d) => {
    data.record = d.record || {};
    setModalProps({
      footer: null,
      confirmLoading: false,
      destroyOnClose: true,
      showCancelBtn: false,
      showOkBtn: false,
      title: `设备详情-${data.record.owner}-${data.record.business}`,
    });
  });
  async function handleModifyDscp() {
    const selectedRowKeys = getSelectRowKeys();
    if (!selectedRowKeys || selectedRowKeys.length == 0) {
      message.error({ content: '请选择要修改的设备' });
      return;
    }
    openModifyDscpModal(true, {
      type: modifyType.SERVER,
      ids: selectedRowKeys as number[],
      owner: data.record.owner,
      business: data.record.business,
    });
  }
  function handleModifyDscpSuccess(results: Array<ModifyDscpResult>) {
    openModifyDscpResultModal(true, { results });
  }
  function handelCancel() {
    emit('success');
  }

  // 导出数据
  function handleExportServerPriority() {
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
        value.summaryId = data.record.id;

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
      let filename = await downloadFileByUrl(
        Api.ExportServerPriority(),
        'POST',
        5 * 60,
        value,
        null,
      );
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
  function formatParams(params: Recordable) {
    params.hostname = splitByLineAndTrim(params.hostname);
    params.ecdnIp = splitByLineAndTrim(params.ecdnIp);
    return params;
  }
</script>

<style>
  .network-server-priority-modal .scrollbar.scroll-container {
    padding-top: 10px !important;
    padding-bottom: 16px !important;
  }

  .network-server-priority-modal .scrollbar__wrap.scrollbar__wrap--hidden-default {
    margin-bottom: 0 !important;
  }

  .network-server-priority-modal .vben-basic-table.vben-basic-table-form-container {
    padding: 0 !important;
  }

  .network-server-priority-modal .vben-basic-table.vben-basic-table-form-container form {
    margin-bottom: 5px !important;
    padding: 0 !important;
  }
</style>
