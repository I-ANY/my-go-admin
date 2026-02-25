<template>
  <div>
    <BasicTable @register="registerTable">
      <!-- <template #toolbar>
        <a-button
          type="primary"
          v-auth="LAPermissionCodeEnum.TRAFFIC_DATA_EXPORT"
          @click="handleExportData"
          :loading="data.exporting"
          >{{ data.exportButTitle }}</a-button
        >
      </template> -->
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
        <template v-if="column.key == 'hostname'">
          <Tooltip title="查看设备历史主机名">
            <a @click="handleViewHistoryHostname(record)">
              {{ record.hostname }}
            </a>
          </Tooltip>
        </template>
      </template>
    </BasicTable>
    <LaDeviceHostStoryModal @register="registerHostStoryModal" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { columns, searchFormSchema, isDiffIspMap, businessStatusMap } from './data';
  import { Tag, Tooltip } from 'ant-design-vue';
  import { onMounted, reactive } from 'vue';
  import { splitByLineAndTrim } from '@/utils/util';
  import { GetBusiness, GetDeviceList } from '@/api/business/la';
  import { getAreaData, getProvince } from '@/utils/kAreaSelect';
  import { DefaultOptionType } from 'ant-design-vue/es/select';
  import LaDeviceHostStoryModal from './LaDeviceHostStoryModal.vue';
  import { useModal } from '@/components/Modal';

  const data = reactive({
    exporting: false,
    // exportButTitle: '导出数据',
    areaProvinceTree: {} as Record<string, DefaultOptionType[]>,
    allProvinceOptions: [] as DefaultOptionType[],
    areaOptions: [] as DefaultOptionType[],
  });
  const [registerHostStoryModal, { openModal }] = useModal();
  defineOptions({ name: 'LaDeviceInfo' });
  const [registerTable, { getForm }] = useTable({
    title: '设备列表',
    api: GetDeviceList,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema(resetProvinces),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 3,
      submitOnReset: false,
      alwaysShowLines: 1,
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
  /*
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
  }*/
  function parseParams(params: Recordable) {
    params.hostnames = splitByLineAndTrim(params.hostnames);
    params.deviceIds = splitByLineAndTrim(params.deviceIds);
  }
  onMounted(async () => {
    resetBusinessOptions();
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
  function handleViewHistoryHostname(record: Recordable) {
    openModal(true, {
      record,
    });
  }
</script>
