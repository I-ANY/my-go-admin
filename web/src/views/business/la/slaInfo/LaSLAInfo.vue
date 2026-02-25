<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="LAPermissionCodeEnum.SLA_DETAIL_EXPORT"
          @click="handleExportData"
          :loading="data.exporting"
          >{{ data.exportButTitle }}</a-button
        >
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'abnormalCount'">
          <Tooltip title="查看异常设备信息">
            <a @click="handleViewAbnormalHost(record)" style="text-decoration: underline">
              {{ record.abnormalCount }}
            </a>
          </Tooltip>
        </template>
        <template v-if="column.key == 'normalRate'">
          <template v-if="record.normalRate || record.normalRate === 0">
            <template v-if="record.normalRate * 100 == 100">
              <span style="color: green">{{ (record.normalRate * 100).toFixed(2) }}% </span>
            </template>
            <template v-else-if="record.normalRate * 100 < 50">
              <span style="color: red">{{ (record.normalRate * 100).toFixed(2) }}%</span>
            </template>
            <template v-else>
              <span style="color: #ff9800">{{ (record.normalRate * 100).toFixed(2) }}%</span>
            </template>
          </template>
          <span v-else></span>
        </template>
      </template>
    </BasicTable>
    <LaSLAAbnormalHostModal @register="registerAbnormalHostModal" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { columns, searchFormSchema } from './data';
  import { Modal, notification, Tooltip } from 'ant-design-vue';
  import { h, nextTick, onMounted, reactive } from 'vue';
  import { Api, GetSLAInfoList } from '@/api/business/la';
  import { getAreaData, getProvince } from '@/utils/kAreaSelect';
  import LaSLAAbnormalHostModal from './LaSLAAbnormalHostModal.vue';
  import { useModal } from '@/components/Modal';
  import dayjs from 'dayjs';
  import { RangePickPresetsExact } from '@/utils/common';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { downloadFileByUrl } from '@/utils/download';
  import { LAPermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import { splitByLineAndTrim } from '@/utils/util';

  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
  });

  const [registerAbnormalHostModal, { openModal }] = useModal();
  defineOptions({ name: 'LaSLAInfo' });
  const [registerTable, { getForm }] = useTable({
    title: '设备SLA信息',
    api: GetSLAInfoList,
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
        nextTick(async () => {
          await resetReportTime();
        });
        return Promise.resolve();
      },
      // actionColOptions: {
      //   span: 5,
      // },
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    rowKey: 'id',
    beforeFetch: (params) => {
      parseParams(params);
      return params;
    },
    pagination: {},
  });

  onMounted(async () => {
    await resetReportTime();
    await resetProvinceOptions();
  });
  async function resetProvinceOptions() {
    const areaData = await getAreaData(getProvince);
    const allProvinceOptions = areaData.allProvinceOptions;
    getForm().updateSchema({
      field: 'provinces',
      componentProps: {
        options: allProvinceOptions,
      },
    });
  }
  function handleViewAbnormalHost(record: Recordable) {
    openModal(true, {
      record,
    });
  }
  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      timeBegin: dayjs(dayjs().format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      timeEnd: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    });
  }
  function onTimePikerOpen() {
    getForm().updateSchema({
      field: '[timeBegin, timeEnd]',
      componentProps: {
        presets: RangePickPresetsExact(),
      },
    });
  }
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
        await exportTraffic(value);
      },
    });
  }
  async function exportTraffic(value: Recordable) {
    try {
      let filename = await downloadFileByUrl(Api.ExportSLAInfo(), 'POST', 5 * 60, value, null);
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
  function parseParams(params: Recordable) {
    params.hostname = splitByLineAndTrim(params.hostname);
  }
</script>
