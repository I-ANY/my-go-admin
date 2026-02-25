<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="handleExportData" :loading="data.exporting">
          {{ data.exportButTitle }}
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'reason'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '500px' }">
            <template #title>{{ record.reason }}</template>
            <span>{{ truncatedFields(record.reason) }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                label: '拉黑详情',
                onClick: selectLimitInfo.bind(null, record),
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <LimitListModal @register="registerLimitList" />
  </div>
</template>

<script setup lang="ts">
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { GetUtilizationList, Api } from '@/api/business/sa';
  import {
    utilizationSearchSchema,
    utilizationTableColumns,
  } from '@/views/business/sa/utilization/data';
  import { splitByLineAndTrim } from '@/utils/util';
  import { truncatedFields } from '@/views/business/zap/qualityEvents/data';
  import { notification, Tooltip } from 'ant-design-vue';
  import { useModal } from '@/components/Modal';
  import { RangePickPresetsExact } from '@/utils/common';
  import { nextTick, onMounted, reactive } from 'vue';
  import dayjs from 'dayjs';
  import LimitListModal from './limitListModal.vue';
  import { formatToDateTime } from '@/utils/dateUtil';
  import { defHttp } from '@/utils/http/axios';

  const data = reactive({
    exporting: false,
    exportButTitle: '导出',
  });

  const [registerLimitList, { openModal: openLimitListModal }] = useModal();
  const [registerTable, { getForm }] = useTable({
    title: '点播利用率信息列表',
    api: GetUtilizationList,
    columns: utilizationTableColumns,
    formConfig: {
      labelWidth: 100,
      schemas: utilizationSearchSchema(onTimePikerOpen),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      pageSizeOptions: ['5', '10', '50'],
    },
    actionColumn: {
      width: 100,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: 'right',
    },
    beforeFetch: (params) => {
      params.hostnames = splitByLineAndTrim(params.hostnames) || null;
      params.guids = splitByLineAndTrim(params.guids) || null;
      params.owners = splitByLineAndTrim(params.owners) || null;
      params.server_ids = splitByLineAndTrim(params.server_ids) || null;
      // 处理排序参数
      if (params.order && params.field) {
        params.sort_order = params.order === 'ascend' ? 'asc' : 'desc';
        params.sort_field = params.field;
        delete params.order;
        delete params.field;
      }
      params.end_statistical_time = params.end_statistical_time.split(' ')[0] += ' 23:59:59';
    },
  });

  function onTimePikerOpen() {
    getForm().updateSchema({
      field: '[start_statistical_time, end_statistical_time]',
      componentProps: {
        presets: RangePickPresetsExact(),
      },
    });
  }

  onMounted(async () => {
    await resetReportTime();
  });

  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      start_statistical_time: dayjs(
        dayjs().add(-1, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      end_statistical_time: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    });
  }

  async function selectLimitInfo(record) {
    openLimitListModal(true, {
      date: formatToDateTime(record.statistical_time, 'YYYYMMDD'),
      guid: record.guid,
    });
  }

  function handleExportData() {
    const formValue = getForm().getFieldsValue();
    formValue.end_statistical_time = formValue.end_statistical_time.split(' ')[0] += ' 23:59:59';
    formValue.hostnames = splitByLineAndTrim(formValue.hostnames) || null;
    formValue.server_ids = splitByLineAndTrim(formValue.server_ids) || null;
    formValue.owners = splitByLineAndTrim(formValue.owners) || null;
    formValue.guids = splitByLineAndTrim(formValue.guids) || null;
    formValue.is_limit =
      formValue.is_limit === 'true' ? true : formValue.is_limit === 'false' ? false : null;

    (async function () {
      await getForm().validate();
      data.exporting = true;
      data.exportButTitle = '导出中...';
      try {
        await ExportUtilizationInfo(formValue);
      } catch (error) {
        notification.error({
          message: '导出失败',
          description: error.message,
        });
        data.exporting = false;
        data.exportButTitle = '导出';
      }
    })();
  }

  async function ExportUtilizationInfo(value: Recordable) {
    const res = await defHttp.post(
      {
        url: Api.ExportUtilizationList,
        responseType: 'blob',
        data: value,
        timeout: 10 * 60 * 1000,
      },
      { isReturnNativeResponse: true },
    );
    try {
      if (!res.headers['content-type'].includes('application/octet-stream')) {
        // 将 Blob 转换为 JSON
        const reader = new FileReader();
        reader.onload = () => {
          const jsonResponse = JSON.parse(reader.result as any);
          notification.error({
            message: '导出失败',
            description: jsonResponse.msg || '未知错误',
            duration: null,
          });
        };
        reader.readAsText(res.data);
        return;
      }
      const blob = new Blob([res.data], { type: res.headers['content-type'] });
      // 创建新的URL并指向File对象或者Blob对象的地址
      const blobURL = window.URL.createObjectURL(blob);
      // 创建a标签，用于跳转至下载链接
      const tempLink = document.createElement('a');
      tempLink.style.display = 'none';
      tempLink.href = blobURL;
      const contentDisposition =
        res.headers['content-disposition'] || `attachment;filename=hdd_device_info.csv`;
      const filename = contentDisposition.split(';')[1].split('=')[1].split("''")[1];
      tempLink.setAttribute('download', filename);
      // 兼容：某些浏览器不支持HTML5的download属性
      if (typeof tempLink.download === 'undefined') {
        tempLink.setAttribute('target', '_blank');
      }
      // 挂载a标签
      document.body.appendChild(tempLink);
      tempLink.click();
      document.body.removeChild(tempLink);
      // 释放blob URL地址
      window.URL.revokeObjectURL(blobURL);
      notification.success({
        message: '导出成功',
        description: '文件名：' + filename,
        duration: null,
      });
    } finally {
      nextTick(() => {
        data.exporting = false;
        data.exportButTitle = '导出';
      });
    }
  }
</script>

<script setup lang="ts"></script>
<style scoped lang="less"></style>
