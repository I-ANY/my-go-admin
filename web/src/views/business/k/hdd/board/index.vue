<template>
  <div class="p-4 h-full">
    <div class="h-full w-full p-15px bg-white">
      <div class="pt-10px">
        <BasicForm @register="registerForm" @submit="handleSubmit">
          <template #advanceAfter>
            <Dropdown v-if="showDropdown()" :loading="data.exporting">
              <template #overlay>
                <Menu>
                  <MenuItem
                    key="1"
                    @click="handleExportDailyPeak"
                    v-if="hasPermission(KPermissionCodeEnum.BUSINESS_K_HDD_DAILYPEAK_EXPORT)"
                  >
                    导出日峰值数据
                  </MenuItem>
                  <MenuItem
                    key="2"
                    @click="handleExportDailyPeakDetail"
                    v-if="hasPermission(KPermissionCodeEnum.BUSINESS_K_HDD_DAILYPEAK_DETAIL_EXPORT)"
                  >
                    导出日峰值明细数据
                  </MenuItem>
                  <MenuItem
                    key="3"
                    @click="handleExportDaily"
                    v-if="hasPermission(KPermissionCodeEnum.BUSINESS_K_HDD_5MIN_SUMMARY_EXPORT)"
                  >
                    导出5min点数据
                  </MenuItem>
                </Menu>
              </template>
              <a-button>
                {{ data.exportButTitle }}
                <DownOutlined />
              </a-button>
            </Dropdown>
          </template>
          <template #capacity="{ model }">
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
          <template #networkSpeed="{ model }">
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
        </BasicForm>
      </div>
      <div class="h-10/11 pt-10px">
        <KHDDLineGraph
          :reportTimeBegin="data.reportTimeBegin"
          :reportTimeEnd="data.reportTimeEnd"
          :minimumSize="data.minimumSize"
          :maximumSize="data.maximumSize"
          :minNetworkSpeed="data.minNetworkSpeed"
          :maxNetworkSpeed="data.maxNetworkSpeed"
          :processStatus="data.processStatus"
          :businessMountStatus="data.businessMountStatus"
          ref="lineGraphRef"
          class="h-full w-full"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
  import { BasicForm, useForm } from '@/components/Form';
  import { hddBoardFormSchema } from './data';
  import dayjs from 'dayjs';
  import { h, nextTick, onMounted, reactive, ref } from 'vue';
  import KHDDLineGraph from './KHDDLineGraph.vue';
  import {
    Dropdown,
    MenuItem,
    Menu,
    notification,
    Modal,
    Input,
    InputNumber,
    FormItem,
    Row,
    Col,
  } from 'ant-design-vue';
  import { usePermission } from '@/hooks/web/usePermission';
  import { KPermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import { DownOutlined, ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { downloadFileByUrl } from '@/utils/download';
  import { Api } from '@/api/business/k';

  const { hasPermission } = usePermission();
  let data = reactive({
    reportTimeBegin: '',
    reportTimeEnd: '',
    minimumSize: undefined as number | undefined,
    maximumSize: undefined as number | undefined,
    minNetworkSpeed: undefined as number | undefined,
    maxNetworkSpeed: undefined as number | undefined,
    processStatus: undefined,
    businessMountStatus: undefined,
    avgData: 10,
    peakData: [],
    exporting: false,
    exportButTitle: '导出数据',
  });

  // 分区挂载状态选项
  // const mountStatusOptions = getSelectOptionsFromDict(KEnum.PARTITION_MOUNT_STATUS);
  // 导出时选择的挂载类型
  const exportMountStatus = ref<string | undefined>(undefined);
  const lineGraphRef = ref<InstanceType<typeof KHDDLineGraph>>();
  const [registerForm, { setFieldsValue, validate, getFieldsValue }] = useForm({
    labelWidth: 120,
    baseColProps: { span: 24 },
    // actionColOptions: {
    //   span: 6,
    // },
    schemas: hddBoardFormSchema(onTimeChange),
    showActionButtonGroup: true,
    compact: true,
    autoAdvancedLine: 1,
    resetFunc() {
      nextTick(async () => {
        await resetReportTime();
      });
      return Promise.resolve();
    },
  });

  onMounted(() => {
    resetReportTime();
  });

  async function onTimeChange() {
    await handleSubmit();
  }

  async function resetReportTime() {
    await setFieldsValue({
      reportTimeBegin: dayjs(
        dayjs().add(-15, 'd').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      reportTimeEnd: dayjs(
        dayjs().add(-1, 'd').format('YYYY-MM-DD 23:59:59'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      minimumSize: 3758,
      minNetworkSpeed: 1,
    });
  }

  async function handleSubmit() {
    await validate();
    const values = await getFormValues();
    data.reportTimeBegin = values.reportTimeBegin;
    data.reportTimeEnd = values.reportTimeEnd;
    data.minimumSize = values.minimumSize;
    data.maximumSize = values.maximumSize;
    data.minNetworkSpeed = values.minNetworkSpeed;
    data.maxNetworkSpeed = values.maxNetworkSpeed;
    data.processStatus = values.processStatus;
    data.businessMountStatus = values.businessMountStatus;

    nextTick(() => {
      lineGraphRef.value?.rebuildGraphic();
    });
  }

  async function getFormValues() {
    const values = await getFieldsValue();
    let result = {} as any;
    let reportTimeBegin = dayjs(values.reportTimeBegin).format('YYYY-MM-DD');
    let reportTimeEnd = dayjs(values.reportTimeEnd).format('YYYY-MM-DD');
    result.reportTimeBegin = reportTimeBegin;
    result.reportTimeEnd = reportTimeEnd;

    if (values.minimumSize || values.minimumSize == 0) {
      result.minimumSize = parseInt(values.minimumSize) * 1000 * 1000 * 1000;
    } else {
      result.minimumSize = undefined;
    }

    if (values.maximumSize || values.maximumSize == 0) {
      result.maximumSize = parseInt(values.maximumSize) * 1000 * 1000 * 1000;
    } else {
      result.maximumSize = undefined;
    }

    if (values.minNetworkSpeed || values.minNetworkSpeed == 0) {
      result.minNetworkSpeed = parseInt(values.minNetworkSpeed) * 1000 * 1000;
    } else {
      result.minNetworkSpeed = undefined;
    }

    if (values.maxNetworkSpeed || values.maxNetworkSpeed == 0) {
      result.maxNetworkSpeed = parseInt(values.maxNetworkSpeed) * 1000 * 1000;
    } else {
      result.maxNetworkSpeed = undefined;
    }
    result.processStatus = values.processStatus;
    result.businessMountStatus = values.businessMountStatus;
    return result;
  }

  function showDropdown() {
    return (
      hasPermission(KPermissionCodeEnum.BUSINESS_K_HDD_DAILYPEAK_EXPORT) ||
      hasPermission(KPermissionCodeEnum.BUSINESS_K_HDD_DAILYPEAK_DETAIL_EXPORT) ||
      hasPermission(KPermissionCodeEnum.BUSINESS_K_HDD_5MIN_SUMMARY_EXPORT)
    );
  }

  async function handleExportDailyPeak() {
    await validate();
    let formValues = await getFormValues();
    Modal.confirm({
      title: '是否确认导出日峰值数据?',
      icon: h(ExclamationCircleOutlined),
      content: () =>
        h('div', [
          h(
            'div',
            { style: 'color:red; margin-bottom: 12px;' },
            '导出将花费一些时间，请耐心等待！如果数据过大，可能会导出失败，请缩小导出条件后重试',
          ),
          // h('div', { style: 'margin: 0 0 4px 0; font-weight: 500;' }, '分区挂载状态'),
          // h(Select, {
          //   value: exportMountStatus.value,
          //   'onUpdate:value': (val: any) => (exportMountStatus.value = val),
          //   style: 'width: 100%',
          //   placeholder: '请选择分区挂载状态（不选则导出所有）',
          //   options: mountStatusOptions,
          //   allowClear: true,
          // }),
        ]),
      async onOk() {
        await exportDailyPeak({
          ...formValues,
        });
      },
    });
  }
  async function exportDailyPeak(value: Recordable) {
    try {
      nextTick(() => {
        data.exporting = true;
        data.exportButTitle = '导出中';
      });
      let filename = await downloadFileByUrl(Api.ExportHddDailyPeak, 'POST', 5 * 60, value, null);
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

  async function handleExportDailyPeakDetail() {
    await validate();
    let formValues = await getFormValues();
    Modal.confirm({
      title: '是否确认导出日峰值明细数据?',
      icon: h(ExclamationCircleOutlined),
      content: () =>
        h('div', [
          h(
            'div',
            { style: 'color:red; margin-bottom: 12px;' },
            '导出将花费一些时间，请耐心等待！如果数据过大，可能会导出失败，请缩小导出条件后重试',
          ),
          // h('div', { style: 'margin: 0 0 4px 0; font-weight: 500;' }, '分区挂载状态'),
          // h(Select, {
          //   value: exportMountStatus.value,
          //   'onUpdate:value': (val: any) => (exportMountStatus.value = val),
          //   style: 'width: 100%',
          //   placeholder: '请选择分区挂载状态（不选则导出所有）',
          //   options: mountStatusOptions,
          //   allowClear: true,
          // }),
        ]),
      async onOk() {
        await exportDailyPeakDetail({
          ...formValues,
        });
      },
    });
  }

  async function exportDailyPeakDetail(value: Recordable) {
    try {
      nextTick(() => {
        data.exporting = true;
        data.exportButTitle = '导出中';
      });
      let filename = await downloadFileByUrl(
        Api.ExportHddDailyPeakDetail,
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

  async function handleExportDaily() {
    await validate();
    let formValues = await getFormValues();
    exportMountStatus.value = undefined;
    Modal.confirm({
      title: '是否确认导出5min点数据?',
      icon: h(ExclamationCircleOutlined),
      content: () =>
        h('div', [
          h(
            'div',
            { style: 'color:red; margin-bottom: 12px;' },
            '导出将花费一些时间，请耐心等待！如果数据过大，可能会导出失败，请缩小导出条件后重试',
          ),
          // h('div', { style: 'margin: 0 0 4px 0; font-weight: 500;' }, '分区挂载状态'),
          // h(Select, {
          //   value: exportMountStatus.value,
          //   'onUpdate:value': (val: any) => (exportMountStatus.value = val),
          //   style: 'width: 100%',
          //   placeholder: '请选择分区挂载状态（不选则导出所有）',
          //   options: mountStatusOptions,
          //   allowClear: true,
          // }),
        ]),
      async onOk() {
        await export5minSummary({
          ...formValues,
        });
      },
    });
  }
  async function export5minSummary(value: Recordable) {
    try {
      nextTick(() => {
        data.exporting = true;
        data.exportButTitle = '导出中';
      });
      let filename = await downloadFileByUrl(Api.ExportHdd5minSummary, 'POST', 5 * 60, value, null);
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
</script>

<style scoped lang="scss">
  // 覆盖PageWrapper的默认margin
  :deep(.ant-page-header) {
    margin: 0;
  }

  :deep(.vben-page-wrapper-content) {
    height: 100%;
    margin: 0 !important;
  }

  // 确保flex容器占满高度
  .flex-col {
    height: 100%;
  }

  // 允许flex子元素收缩
  .min-h-0 {
    min-height: 0;
  }

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
