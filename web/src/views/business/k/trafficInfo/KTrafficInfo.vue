<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="'business:k:traffic:export'"
          @click="handleExportData"
          :loading="data.exporting"
          >{{ data.exportButTitle }}</a-button
        >
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'ispId'">
          <span v-if="ispMap[record.ispId]">{{ ispMap[record.ispId].dictLabel }}</span>
        </template>
        <template v-if="column.key == 'devType'">
          <span v-if="devTypeMap[record.devType]">{{ devTypeMap[record.devType].dictLabel }}</span>
        </template>
        <template v-if="column.key == 'deviceType'">
          <span v-if="deviceTypeMap[record.deviceType]">{{
            deviceTypeMap[record.deviceType].dictLabel
          }}</span>
        </template>
        <template v-if="column.key == 'demandIspId'">
          <span v-if="ispMap[record.demandIspId]">{{ ispMap[record.demandIspId].dictLabel }}</span>
        </template>
        <template v-if="column.key == 'isProvinceScheduling'">
          <span v-if="isProvinceSchedulingMap[record.isProvinceScheduling]">{{
            isProvinceSchedulingMap[record.isProvinceScheduling].dictLabel
          }}</span>
        </template>
        <template v-if="column.key == 'isCoverDiffIsp'">
          <span v-if="isCoverDiffIspMap[record.isCoverDiffIsp]">{{
            isCoverDiffIspMap[record.isCoverDiffIsp].dictLabel
          }}</span>
        </template>
        <template v-if="column.key == 'providerId'">
          <span v-if="providerTypeMap[record.providerId]">{{
            providerTypeMap[record.providerId].dictLabel
          }}</span>
        </template>
        <template v-if="column.key == 'ispConsistency'">
          <Tag
            v-if="ispConsistencyMap[record.ispConsistency]"
            :color="ispConsistencyMap[record.ispConsistency].color || 'default'"
            >{{ ispConsistencyMap[record.ispConsistency].dictLabel }}</Tag
          >
          <span v-else></span>
        </template>
      </template>
    </BasicTable>
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import {
    columns,
    searchFormSchema,
    ispMap,
    devTypeMap,
    deviceTypeMap,
    isProvinceSchedulingMap,
    isCoverDiffIspMap,
    providerTypeMap,
    ispConsistencyMap,
  } from './data';
  import { Api, GetAreaList, GetTrafficList } from '@/api/business/k';
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
    areaInfo: [] as any[],
    allProvinceOptions: [] as string[],
  });
  const { notification } = useMessage();

  defineOptions({ name: 'KTrafficInfo' });
  const [registerTable, { getForm }] = useTable({
    title: '流量数据',
    api: GetTrafficList,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema(onTimePikerOpen, resetDemandProvinces),
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
        await exportKTraffic(value);
      },
    });
  }
  onMounted(async () => {
    await resetReportTime();
  });
  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      timeBegin: dayjs(dayjs().add(-1, 'day').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      timeEnd: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    });
  }

  async function exportKTraffic(value: Recordable) {
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
        duration: null,
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
    params.macAddrs = splitByLineAndTrim(params.macAddrs);
  }
  onMounted(async () => {
    await buildAreaInfo();
    await resetAreaNameOptions();
    await resetDemandProvinces();
  });
  async function buildAreaInfo() {
    let res = await GetAreaList({});
    let areaInfo = res.items;
    let areaInfos: Record<string, string[]> = {};

    for (let i = 0; i < areaInfo?.length; i++) {
      let item = areaInfo[i];
      let areaName = item.area_name;
      if (!areaInfos[areaName]) {
        areaInfos[areaName] = [];
      }
      // 不存在则添加
      if (areaInfos[areaName].indexOf(item.province_name) < 0) {
        areaInfos[areaName].push(item.province_name);
      }
    }
    let allProvinceOptions: any[] = [];
    Object.keys(areaInfos).forEach((areaName: string) => {
      let areaOptions: any[] = [];
      areaInfos[areaName]?.forEach((province: string) => {
        let option = {
          label: province,
          value: province,
        };
        areaOptions.push(option);
        allProvinceOptions.push(option);
      });
      data.areaInfo[areaName] = areaOptions;
    });
    data.allProvinceOptions = allProvinceOptions;
  }
  async function resetAreaNameOptions() {
    let options: any[] = [];
    Object.keys(data?.areaInfo).forEach((key) => {
      let option = {
        label: key,
        value: key,
      };
      options.push(option);
    });
    await getForm().updateSchema({
      field: 'demandAreas',
      componentProps: {
        options: options,
      },
    });
  }
  async function resetDemandProvinces() {
    // 清空当前选中的值
    await getForm().setFieldsValue({
      demandProvinces: [],
    });
    //省份的options
    let proviceOptions: any[] = [];

    // 根据所选的大区获取大区的省份
    let values = await getForm().getFieldsValue();
    if (values.demandAreas?.length > 0) {
      values.demandAreas.forEach((areaName: string) => {
        if (data.areaInfo[areaName]?.length > 0) {
          proviceOptions = proviceOptions.concat(data.areaInfo[areaName]);
        }
      });
    } else {
      proviceOptions = data.allProvinceOptions;
    }

    // 设置省份的options
    await getForm().updateSchema({
      field: 'demandProvinces',
      componentProps: {
        options: proviceOptions,
      },
    });
  }
</script>
