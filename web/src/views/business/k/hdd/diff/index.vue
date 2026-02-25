<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="KPermissionCodeEnum.BUSINESS_K_HDD_DIFF_BATCH_CONFIRM"
          @click="batchDiffConfirm"
          ><template #icon> <CheckOutlined /> </template>批量确认</a-button
        >
        <a-button
          type="primary"
          v-auth="KPermissionCodeEnum.BUSINESS_K_HDD_DIFF_EXPORT"
          @click="onExportHddDiffData"
          :loading="data.exporting"
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
              v-if="data.showTagFields[columnName][record[columnName]]"
              :color="data.showTagFields[columnName][record[columnName]].color || 'default'"
              >{{ data.showTagFields[columnName][record[columnName]].dictLabel }}</Tag
            >
            <span v-else>{{ record[columnName] }}</span>
          </template>
        </template>

        <!-- 确认人 -->
        <template v-if="column.key === 'confirmUser'">
          <span v-if="record.confirmUser">{{ record?.confirmUser?.nickName }}</span>
          <span v-else></span>
        </template>
        <!-- 备注 -->
        <template v-if="column.key === 'remark'">
          <Tooltip v-if="record.remark" :title="record.remark">
            <span>{{ record.remark }}</span>
          </Tooltip>
        </template>

        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'el:ok',
                onClick: diffConfirm.bind(null, record),
                tooltip: '确认流失原因',
                // label: '确认',
                auth: KPermissionCodeEnum.BUSINESS_K_HDD_DIFF_CONFIRM,
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <LossReasonConfirmModal @register="registerModal" @success="onConfirmSuccess" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import {
    columns,
    devTypeMap,
    hddLoseReasonMap,
    isCoverDiffIspMap,
    ispMap,
    isProvinceSchedulingMap,
    searchFormSchema,
  } from './data';
  import { Api, GetHddDiffList } from '@/api/business/k';
  import { useModal } from '@/components/Modal';
  import { h, nextTick, onMounted, reactive } from 'vue';
  import { splitByLineAndTrim } from '@/utils/util';
  import dayjs from 'dayjs';
  import { KPermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import { message, Modal, notification, Tag, Tooltip } from 'ant-design-vue';
  import { CheckOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import LossReasonConfirmModal from './LossReasonConfirmModal.vue';
  import { downloadFileByUrl } from '@/utils/download';
  import { RangeDataPickPresetsExact } from '@/utils/common';

  const [registerModal, { openModal }] = useModal();

  defineOptions({ name: 'KHddDiff' });

  const data = reactive({
    showEnumFields: {
      isProvinceScheduling: isProvinceSchedulingMap,
      isCoverDiffIsp: isCoverDiffIspMap,
      devType: devTypeMap,
      ispId: ispMap,
    },
    showTagFields: {
      lossReason: hddLoseReasonMap,
    },
    exporting: false,
    exportButTitle: '导出数据',
  });

  const [registerTable, { getForm, reload, getSelectRowKeys, setSelectedRowKeys }] = useTable({
    title: '计费差值列表',
    api: GetHddDiffList,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema(onOpenChange),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 3,
      submitOnReset: false,
      alwaysShowLines: 1,
      resetFunc() {
        nextTick(() => {
          resetReportTime();
        });
        return Promise.resolve();
      },
    },
    beforeFetch: async (params: Recordable) => {
      parseValue(params);
      return params;
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    rowKey: 'id',
    clickToRowSelect: false,
    showSelectionBar: true,
    rowSelection: {
      type: 'checkbox',
      // getCheckboxProps: () => {
      //   return {
      //     disabled: false,
      //   };
      // },
    },
    immediate: true,
    actionColumn: {
      width: 120,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: 'right',
    },
  });

  onMounted(async () => {
    await resetReportTime();
    reload();
  });
  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      reportDayBegin: dayjs(
        dayjs().add(-1, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      reportDayEnd: dayjs(
        dayjs().add(-1, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
    });
  }
  function parseValue(params: Recordable) {
    params.hostnames = splitByLineAndTrim(params.hostnames);
    if (params.reportDayBegin) {
      params.reportDayBegin = dayjs(params.reportDayBegin).format('YYYY-MM-DD');
    }
    if (params.reportDayEnd) {
      params.reportDayEnd = dayjs(params.reportDayEnd).format('YYYY-MM-DD');
    }
  }

  async function batchDiffConfirm() {
    const ids = getSelectRowKeys();
    if (ids.length === 0) {
      message.warning('请选择要确认的记录');
      return;
    }
    openModal(true, {
      isBatch: true,
      ids,
      record: null,
    });
  }
  async function diffConfirm(record: Recordable) {
    openModal(true, {
      isBatch: false,
      ids: [record.id],
      record,
    });
  }
  async function onConfirmSuccess() {
    reload();
    setSelectedRowKeys([]);
  }
  async function onExportHddDiffData() {
    const formValue = getForm().getFieldsValue();
    formValue.hostnames = splitByLineAndTrim(formValue.hostnames) || null;
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
        const value = await getForm().getFieldsValue();
        parseValue(value);
        nextTick(() => {
          data.exporting = true;
          data.exportButTitle = '导出中';
        });
        await exportHddDiffData(value);
      },
    });
  }
  async function exportHddDiffData(value: Recordable) {
    try {
      let filename = await downloadFileByUrl(Api.ExportHddDiff, 'POST', 5 * 60, value, null);
      notification.success({
        message: '导出成功',
        description: '文件名：' + filename,
        duration: null,
      });
    } catch (error) {
      notification.error({
        message: '导出失败',
        description: error.message,
      });
    } finally {
      nextTick(() => {
        data.exporting = false;
        data.exportButTitle = '导出数据';
      });
    }
  }
  async function onOpenChange() {
    getForm().updateSchema({
      field: '[reportDayBegin, reportDayEnd]',
      componentProps: {
        presets: RangeDataPickPresetsExact(),
      },
    });
  }
</script>
