<template>
  <BasicTable @register="registerTable">
    <template #toolbar>
      <a-button
        type="primary"
        v-auth="'business:b:traffic:export'"
        @click="handleExportTrafficInfo"
        :loading="data.exporting"
        >{{ data.exportButTitle }}</a-button
      >
    </template>
    <template #bodyCell="{ column, record }">
      <template v-if="column.key == 'deviceType'">
        <span v-if="deviceTypeMap[record.deviceType]"
          >{{ deviceTypeMap[record.deviceType].dictLabel }}
        </span>
      </template>
      <template v-if="column.key == 'ipType'">
        <span v-if="ipTypeMap[record.ipType]">{{ ipTypeMap[record.ipType].dictLabel }} </span>
      </template>
      <template v-if="column.key == 'mfyIpType'">
        <span v-if="ipTypeMap[record.mfyIpType]">{{ ipTypeMap[record.mfyIpType].dictLabel }} </span>
      </template>
      <template v-if="column.key == 'agentStatus'">
        <Tag
          :color="agentStatusMap[record.agentStatus].color"
          v-if="agentStatusMap[record.agentStatus]"
          >{{ agentStatusMap[record.agentStatus].dictLabel }}
        </Tag>
      </template>
      <template v-if="column.key == 'businessStatus'">
        <Tag
          :color="businessStatusMap[record.businessStatus].color"
          v-if="businessStatusMap[record.businessStatus]"
          >{{ businessStatusMap[record.businessStatus].dictLabel }}
        </Tag>
      </template>
      <template v-if="column.key == 'localCoverageOnly'">
        <span v-if="localCoverageOnlyMap[record.localCoverageOnly]"
          >{{ localCoverageOnlyMap[record.localCoverageOnly].dictLabel }}
        </span>
        <span v-else></span>
      </template>
      <template v-if="column.key == 'networkSpeed'">
        <span>{{ customCeilDivide(record.networkSpeed) }}</span>
      </template>
      <template v-if="column.key == 'mfyNetworkSpeed'">
        <span>{{ customCeilDivide(record.mfyNetworkSpeed) }}</span>
      </template>
    </template>
  </BasicTable>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import {
    deviceTypeMap,
    ipTypeMap,
    agentStatusMap,
    businessStatusMap,
    localCoverageOnlyMap,
    trafficFormSchema,
    trafficColumns,
  } from './data';
  import { Api, GetTrafficList, GetTrafficOptions } from '@/api/business/b';
  import { h, nextTick, onMounted, reactive } from 'vue';
  import dayjs from 'dayjs';
  import { Modal, notification, Tag } from 'ant-design-vue';
  import { customCeilDivide } from '@/utils/util';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { downloadFileByUrl } from '@/utils/download';
  import { RangePickPresetsExact } from '@/utils/common';

  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
  });
  defineOptions({ name: 'BTrafficInfo' });
  const [registerTable, { getForm }] = useTable({
    title: 'xagent数据',
    api: GetTrafficList,
    columns: trafficColumns,
    formConfig: {
      labelWidth: 120,
      schemas: trafficFormSchema(onTimeChange, onTimePikerOpen),
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
    // async beforeFetch(params: Recordable) {
    //   parseValue(params);
    //   return params;
    // },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    rowKey: 'id',
    // actionColumn: {
    //   width: 120,
    //   title: '操作',
    //   dataIndex: 'action',
    //   // slots: { customRender: 'action' },
    //   fixed: 'right',
    // },
  });

  onMounted(async () => {
    await resetReportTime();
  });
  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      reportTimeBegin: dayjs(
        dayjs().add(-7, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      reportTimeEnd: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    });
  }
  function handleExportTrafficInfo() {
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

        await exportHddDeviceStatus(value);
      },
    });
  }
  async function exportHddDeviceStatus(value) {
    // TODO 修改 参数
    nextTick(() => {
      data.exporting = true;
      data.exportButTitle = '导出中';
    });
    try {
      let filename = await downloadFileByUrl(Api.ExportTraffic, 'POST', 5 * 60, value, null);
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
  async function onTimeChange() {
    await restBusinessOptions();
  }
  async function restBusinessOptions() {
    const data = await getForm().getFieldsValue();
    const res = await GetTrafficOptions({
      reportTimeBegin: data.reportTimeBegin,
      reportTimeEnd: data.reportTimeEnd,
      type: 'business',
    });
    getForm().updateSchema({
      field: 'business',
      componentProps: {
        options: res,
      },
    });
  }
  function onTimePikerOpen() {
    getForm().updateSchema({
      field: '[reportTimeBegin, reportTimeEnd]',
      componentProps: {
        presets: RangePickPresetsExact(),
      },
    });
  }
</script>
