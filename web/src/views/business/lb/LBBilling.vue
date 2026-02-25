<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="'business:lb:bill:export'"
          @click="handleExportData"
          :loading="data.exporting"
          >{{ data.exportButTitle }}</a-button
        >
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'isp'">
          <span v-if="ispMap[record.isp]">{{ ispMap[record.isp].dictLabel }}</span>
        </template>
        <template v-if="column.key == 'billType'">
          <span v-if="billTypeMap[record.billType]">{{
            billTypeMap[record.billType].dictLabel
          }}</span>
        </template>
        <template v-if="column.key == 'status'">
          <Tag
            v-if="statusMap[record.status]"
            :color="statusMap[record.status].color || 'default'"
            >{{ statusMap[record.status].dictLabel }}</Tag
          >
          <span v-else>{{ record.status }}</span>
        </template>
      </template>
    </BasicTable>
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { columns, searchFormSchema, ispMap, billTypeMap, statusMap } from './data';
  import { Api, GetBillList } from '@/api/business/lb';
  import { Modal, Tag } from 'ant-design-vue';
  import { h, nextTick, onMounted, reactive } from 'vue';
  import dayjs from 'dayjs';
  import { useMessage } from '@/hooks/web/useMessage';
  import { downloadFileByUrl } from '@/utils/download';
  import { RangePickPresetsExact } from '@/utils/common';
  import { splitByLineAndTrim } from '@/utils/util';

  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
  });
  const { notification } = useMessage();

  defineOptions({ name: 'LBBilling' });
  const [registerTable, { getForm }] = useTable({
    title: '客户计费数据',
    api: GetBillList,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema(onTimePikerOpen),
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
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    rowKey: 'id',
    beforeFetch: (params) => {
      parseParams(params);
      return params;
    },
    // actionColumn: {
    //   width: 120,
    //   title: '操作',
    //   dataIndex: 'action',
    //   // slots: { customRender: 'action' },
    //   fixed: 'right',
    // },
  });

  function handleExportData() {
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
        parseParams(value);
        nextTick(() => {
          data.exporting = true;
          data.exportButTitle = '导出中';
        });
        await exportLbBill(value);
      },
    });
  }
  onMounted(async () => {
    await resetReportTime();
  });
  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      timeBegin: dayjs(dayjs().add(-7, 'day').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      timeEnd: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    });
  }

  async function exportLbBill(value: Recordable) {
    console.log('value', value);

    try {
      let filename = await downloadFileByUrl(Api.ExportBill, 'POST', 5 * 60, value, null);
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
  function onTimePikerOpen() {
    getForm().updateSchema({
      field: '[timeBegin, timeEnd]',
      componentProps: {
        presets: RangePickPresetsExact(),
      },
    });
  }
  function parseParams(params: Recordable) {
    params.hostnames = splitByLineAndTrim(params.hostnames);
    params.uids = splitByLineAndTrim(params.uids);
    params.deviceIds = splitByLineAndTrim(params.deviceIds);
  }
</script>
