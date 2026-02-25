<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="'business:k:hddDeviceInfo:export'"
          @click="handleExportData"
          :loading="data.exporting"
          >{{ data.exportButTitle }}</a-button
        >
      </template>
      <template #bodyCell="{ column, record }">
        <!-- 使用同一个枚举值，如果有变更，需要替换 -->
        <template v-if="column.key == 'ispName'">
          <span v-if="ispNameMap[record.ispName]">{{
            ispNameMap[record.ispName].dictLabel || record.ispName
          }}</span>
        </template>
        <template v-if="column.key == 'deviceIsp'">
          <span v-if="ispNameMap[record.deviceIsp]">{{
            ispNameMap[record.deviceIsp].dictLabel || record.deviceIsp
          }}</span>
        </template>

        <!-- 使用同一个枚举值，如果有变更，需要替换 -->
        <template v-if="column.key == 'status'">
          <span v-if="deviceStatusMap[record.status]">{{
            deviceStatusMap[record.status].dictLabel || record.status
          }}</span>
        </template>
        <template v-if="column.key == 'isFirstMac'">
          <span v-if="isFirstMacMap[record.isFirstMac]">{{
            isFirstMacMap[record.isFirstMac].dictLabel
          }}</span>
          <span v-else> </span>
        </template>
        <template v-if="column.key == 'processStatus'">
          <span v-if="deviceStatusMap[record.processStatus]">{{
            deviceStatusMap[record.processStatus].dictLabel || record.processStatus
          }}</span>
        </template>
        <template v-if="column.key == 'ReportFlowServiceStatus'">
          <span v-if="deviceStatusMap[record.ReportFlowServiceStatus]">{{
            deviceStatusMap[record.ReportFlowServiceStatus].dictLabel ||
            record.ReportFlowServiceStatus
          }}</span>
        </template>

        <!-- 挂载详情 -->
        <template v-if="column.key == 'businessMountDetail'">
          <Tooltip
            :overlayStyle="getOverlayStyle(splitByLineAndTrim(record.businessMountDetail)?.length)"
            placement="topLeft"
            v-if="record.businessMountDetail"
          >
            <span>{{ record.businessMountDetail }}</span>
            <template #title>
              <p
                v-for="(item, index) in splitByLineAndTrim(record.businessMountDetail)"
                style=" margin: 0;padding: 0"
                :key="index"
                >{{ item }}</p
              >
            </template>
          </Tooltip>
          <span v-else>- </span>
        </template>
        <!-- 进程信息 -->
        <template v-if="column.key == 'businessProcessInfo'">
          <Tooltip
            :overlayStyle="getOverlayStyle(splitByLineAndTrim(record.businessProcessInfo)?.length)"
            placement="topLeft"
            v-if="record.businessProcessInfo"
          >
            <span>{{ record.businessProcessInfo }}</span>
            <template #title>
              <p
                v-for="(item, index) in splitByLineAndTrim(record.businessProcessInfo)"
                style=" margin: 0;padding: 0"
                :key="index"
                >{{ item }}</p
              >
            </template>
          </Tooltip>
          <span v-else>- </span>
        </template>
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                // color: 'error',
                onClick: handleViewPartitionInfo.bind(null, record),
                tooltip: '查看分区信息',
                label: '查看',
              },
            ]"
          />
        </template>
      </template>
      <template #form-capacity="{ model }">
        <FormItem>
          <Row>
            <Col :span="11">
              <InputNumber
                v-model:value="model['minimumSize']"
                class="w-full site-input-left"
                :min="0"
                style="text-align: center"
                placeholder="最小容量"
              />
            </Col>
            <Col :span="2">
              <Input
                class="w-full site-input-split pl-0 pr-0"
                placeholder="~"
                disabled
                style="text-align: center"
            /></Col>
            <Col :span="11">
              <InputNumber
                v-model:value="model['maximumSize']"
                class="w-full site-input-right"
                :min="1"
                style="text-align: center"
                placeholder="最大容量"
            /></Col>
          </Row>
        </FormItem>
      </template>
      <template #form-networkSpeed="{ model }">
        <FormItem>
          <Row>
            <Col :span="11">
              <InputNumber
                v-model:value="model['minNetworkSpeed']"
                class="w-full site-input-left"
                :min="0"
                style="text-align: center"
                placeholder="最小网速"
              />
            </Col>
            <Col :span="2">
              <Input
                class="w-full site-input-split pl-0 pr-0"
                placeholder="~"
                disabled
                style="text-align: center"
            /></Col>
            <Col :span="11">
              <InputNumber
                v-model:value="model['maxNetworkSpeed']"
                class="w-full site-input-right"
                :min="1"
                style="text-align: center"
                placeholder="最大网速"
            /></Col>
          </Row>
        </FormItem>
      </template>
      <template #form-usage="{ model }">
        <FormItem>
          <Row>
            <Col :span="11">
              <InputNumber
                v-model:value="model['minUsage']"
                class="w-full site-input-left"
                :min="0"
                style="text-align: center"
                placeholder="最小使用率"
              />
            </Col>
            <Col :span="2">
              <Input
                class="w-full site-input-split"
                placeholder="~"
                disabled
                style="text-align: center"
            /></Col>
            <Col :span="11">
              <InputNumber
                v-model:value="model['maxUsage']"
                class="w-full site-input-right"
                :min="1"
                style="text-align: center"
                placeholder="最大使用率"
            /></Col>
          </Row>
        </FormItem>
      </template>
    </BasicTable>
    <PartitionInfoModal @register="registerModal" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { columns, searchFormSchema, ispNameMap, deviceStatusMap, isFirstMacMap } from './data';
  import { Api, GetHddDeviceStatusList } from '@/api/business/k';
  import PartitionInfoModal from './PartitionInfoModal.vue';
  import { useModal } from '@/components/Modal';
  import { InputNumber, Input, Row, Col, FormItem, Modal, Tooltip } from 'ant-design-vue';
  import { h, nextTick, onMounted, reactive } from 'vue';
  import { splitByLineAndTrim } from '@/utils/util';
  import dayjs from 'dayjs';
  import { useMessage } from '@/hooks/web/useMessage';
  import { downloadFileByUrl } from '@/utils/download';
  import { RangePickPresetsExact } from '@/utils/common';

  const [registerModal, { openModal }] = useModal();
  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
  });
  const { notification } = useMessage();

  defineOptions({ name: 'KHddDeviceStatus' });
  const [registerTable, { getForm }] = useTable({
    title: '设备列表',
    api: GetHddDeviceStatusList,
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
    async beforeFetch(params: Recordable) {
      parseValue(params);
      return params;
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    rowKey: 'id',
    actionColumn: {
      width: 120,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: 'right',
    },
  });
  function handleViewPartitionInfo(record: Recordable) {
    const formValue = getForm().getFieldsValue();
    openModal(true, { record, formValue });
  }

  function handleExportData() {
    const formValue = getForm().getFieldsValue();
    formValue.hostnames = splitByLineAndTrim(formValue.hostnames) || null;
    formValue.macAddrs = splitByLineAndTrim(formValue.macAddrs) || null;
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
        await exportHddDeviceStatus(value);
      },
    });
  }
  onMounted(async () => {
    await resetReportTime();
  });
  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      reportTimeBegin: dayjs(
        dayjs().add(-4, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      reportTimeEnd: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    });
  }
  function parseValue(params: Recordable) {
    params.hostnames = splitByLineAndTrim(params.hostnames);
    params.macAddrs = splitByLineAndTrim(params.macAddrs);
    // params.includePartition = true;

    if (params.minimumSize) {
      params.minimumSize = parseInt(params.minimumSize) * 1000 * 1000 * 1000;
    }
    if (params.maximumSize) {
      params.maximumSize = parseInt(params.maximumSize) * 1000 * 1000 * 1000;
    }

    if (params.minNetworkSpeed) {
      params.minNetworkSpeed = parseInt(params.minNetworkSpeed) * 1000 * 1000;
    }
    if (params.maxNetworkSpeed) {
      params.maxNetworkSpeed = parseInt(params.maxNetworkSpeed) * 1000 * 1000;
    }
  }
  function getOverlayStyle(len: number | undefined) {
    let style: any = {
      maxWidth: '1200px',
      maxHeight: '700px',
    };
    if (len == undefined || len < 30) {
      return style;
    }
    style.overflow = 'auto';
    return style;
  }
  async function exportHddDeviceStatus(value: Recordable) {
    try {
      let filename = await downloadFileByUrl(
        Api.ExportHddDeviceStatus,
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
      field: '[reportTimeBegin, reportTimeEnd]',
      componentProps: {
        presets: RangePickPresetsExact(),
      },
    });
  }
</script>
<style scoped>
  :deep(.site-input-split) {
    border-right: 0;
    border-left: 0;
    border-radius: 0%;
    background-color: #fff;
  }

  :deep(.site-input-right) {
    border-left-width: 0;
    border-top-left-radius: 0%;
    border-bottom-left-radius: 0%;
  }

  :deep(.site-input-left) {
    border-right-width: 0;
    border-top-right-radius: 0%;
    border-bottom-right-radius: 0%;
  }

  :deep(.site-input-right:hover),
  :deep(.site-input-right:focus) {
    border-left-width: 1px;
  }

  :deep(.site-input-left:hover),
  :deep(.site-input-left:focus) {
    border-right-width: 1px;
  }

  [data-theme='dark'] :deep(.site-input-split) {
    background-color: transparent;
  }
</style>
