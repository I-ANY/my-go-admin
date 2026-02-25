<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="LAPermissionCodeEnum.TRAFFIC_DATA_EXPORT"
          @click="handleExportData"
          :loading="data.exporting"
          >{{ data.exportButTitle }}</a-button
        >

        <Dropdown v-if="showDropdown()">
          <template #overlay>
            <Menu>
              <MenuItem
                key="1"
                @click="handleExportDeviceInfo"
                v-if="hasPermission(LAPermissionCodeEnum.TRAFFIC_DEVICE_INFO_EXPORT)"
              >
                导出设备信息
              </MenuItem>
            </Menu>
          </template>
          <a-button>
            更多
            <DownOutlined />
          </a-button>
        </Dropdown>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'isDiffIsp'">
          <span v-if="isDiffIspMap[record.isDiffIsp]">{{
            isDiffIspMap[record.isDiffIsp].dictLabel
          }}</span>
          <span v-else></span>
        </template>
        <template v-if="column.key == 'businessStatus'">
          <Tag
            v-if="businessStatusMap[record.businessStatus]"
            style="font-weight: bold"
            :color="businessStatusMap[record.businessStatus].color || 'default'"
            >{{ businessStatusMap[record.businessStatus].dictLabel }}</Tag
          >
          <span v-else></span>
        </template>
      </template>
    </BasicTable>
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { ExclamationCircleOutlined, DownOutlined } from '@ant-design/icons-vue';
  import { columns, searchFormSchema, isDiffIspMap, businessStatusMap } from './data';
  import { Modal, Tag, Dropdown, Menu, MenuItem } from 'ant-design-vue';
  import { h, nextTick, onMounted, reactive } from 'vue';
  import dayjs from 'dayjs';
  import { useMessage } from '@/hooks/web/useMessage';
  import { downloadFileByUrl } from '@/utils/download';
  import { RangePickPresetsExact } from '@/utils/common';
  import { splitByLineAndTrim } from '@/utils/util';
  import { LAPermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import { Api, GetBusiness, GetTrafficList } from '@/api/business/la';
  import { getAreaData, getProvince } from '@/utils/kAreaSelect';
  import { DefaultOptionType } from 'ant-design-vue/es/select';
  import { usePermission } from '@/hooks/web/usePermission';

  const { hasPermission } = usePermission();
  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
    areaProvinceTree: {} as Record<string, DefaultOptionType[]>,
    allProvinceOptions: [] as DefaultOptionType[],
    areaOptions: [] as DefaultOptionType[],
  });
  const { notification } = useMessage();

  defineOptions({ name: 'LaTrafficInfo' });
  const [registerTable, { getForm }] = useTable({
    title: '流量数据',
    api: GetTrafficList,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema(onTimePikerOpen, resetProvinces),
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
    params.deviceIds = splitByLineAndTrim(params.deviceIds);
  }
  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      timeBegin: dayjs(dayjs().add(-2, 'day').format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
      timeEnd: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    });
  }
  onMounted(async () => {
    resetBusinessOptions();
    await resetReportTime();
    await buildAreaInfo();
    await resetAreaNameOptions();
    await resetProvinces();
  });
  async function buildAreaInfo() {
    const areaData = await getAreaData(getProvince);
    data.areaProvinceTree = areaData.areaProvinceTree;
    data.allProvinceOptions = areaData.allProvinceOptions;
    data.areaOptions = areaData.areaOptions;
  }
  async function resetAreaNameOptions() {
    await getForm().updateSchema({
      field: 'areas',
      componentProps: {
        options: data.areaOptions,
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

  function handleExportDeviceInfo() {
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
        await exportDeviceInfo(value);
      },
    });
  }
  async function exportDeviceInfo(value: Recordable) {
    try {
      let filename = await downloadFileByUrl(Api.ExportDeviceInfo(), 'POST', 5 * 60, value, null);
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
    }
  }
  function showDropdown() {
    return hasPermission(LAPermissionCodeEnum.TRAFFIC_DEVICE_INFO_EXPORT);
  }
  async function resetProvinces() {
    // 清空当前选中的值
    await getForm().setFieldsValue({
      provinces: [],
    });
    //省份的options
    let proviceOptions: DefaultOptionType[] = [];

    // 根据所选的大区获取大区的省份
    let values = await getForm().getFieldsValue();
    if (values.areas?.length > 0) {
      values.areas.forEach((areaName: string) => {
        if (data.areaProvinceTree[areaName]?.length > 0) {
          proviceOptions = proviceOptions.concat(data.areaProvinceTree[areaName] || []);
        }
      });
    } else {
      proviceOptions = data.allProvinceOptions;
    }

    // 设置省份的options
    await getForm().updateSchema({
      field: 'provinces',
      componentProps: {
        options: proviceOptions,
      },
    });
  }
  async function resetBusinessOptions() {
    let options = [] as Array<any>;
    let business = await GetBusiness({});
    business.forEach((item) => {
      options.push({
        label: item,
        value: item,
      });
    });
    await getForm().updateSchema({
      field: 'business',
      componentProps: {
        options: options,
      },
    });
  }
</script>
