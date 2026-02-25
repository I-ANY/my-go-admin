<template>
  <div>
    <BasicModal v-bind="$attrs" @register="registerModal" @ok="handleSubmit">
      <BasicTable @register="registerTable" @selection-change="handleSelectionChange">
        <template #tableTitle
          ><span style="color: red; font-size: 18px"
            >当前选中设备（{{ ids.length }}个）的真实交付带宽总和：{{
              allData.realDeliveryBw
            }}Mbps</span
          ></template
        >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key == 'bizType'">
            <span v-if="bizTypesMap[record.bizType]">{{
              bizTypesMap[record.bizType].dictLabel
            }}</span>
          </template>
          <template v-if="column.key == 'provider'">
            <span v-if="providersMap[record.provider]">{{
              providersMap[record.provider].dictLabel
            }}</span>
          </template>
          <template v-if="column.key == 'devId'">
            <span v-if="deviceTypeMap[record.devId]">{{
              deviceTypeMap[record.devId].dictLabel
            }}</span>
            <span v-else></span>
          </template>
          <template v-if="column.key == 'demandId'">
            <Tooltip v-if="isFutureMonth(record.startTime)">
              <Badge color="orange" :offset="[10, 10]">
                <span
                  :style="{
                    display: 'block',
                    color: 'orange',
                    marginRight: '6px',
                  }"
                >
                  {{ record.demandId }}
                </span>
                <template #count>
                  <ExclamationCircleFilled style="color: orange" />
                </template>
              </Badge>
              <template #title>次月需求不能当月提交</template>
            </Tooltip>
            <span v-else>{{ record.demandId }}</span>
          </template>
          <template v-if="column.key == 'canBoundedBw'">
            <Tooltip v-if="record.canBoundedBw < allData.realDeliveryBw">
              <Badge color="red">
                <span
                  :style="{
                    display: 'block',
                    marginRight: '6px',
                    color: 'red',
                  }"
                  >{{ record.canBoundedBw }}
                </span>
                <template #count>
                  <ExclamationCircleFilled style="color: red" />
                </template>
              </Badge>
              <template #title>剩余可绑定带宽不足{{ allData.realDeliveryBw }}Mbps</template>
            </Tooltip>
            <span
              v-if="record.canBoundedBw >= allData.realDeliveryBw"
              :style="{
                color: 'green',
              }"
              >{{ record.canBoundedBw }}</span
            >
          </template>
          <template v-if="column.key == 'isSameAreaReplace'">
            <span v-if="sameAreaReplaceMap[record.isSameAreaReplace]">{{
              sameAreaReplaceMap[record.isSameAreaReplace].dictLabel
            }}</span>
            <span v-else>{{ record.isSameAreaReplace }}</span>
          </template>
          <template v-if="column.key == 'hddTransform'">
            <span v-if="hddTransformMap[record.hddTransform]">{{
              hddTransformMap[record.hddTransform].dictLabel
            }}</span>
            <span v-else>{{ record.hddTransform }}</span>
          </template>
          <template v-if="column.key == 'isIndependentDeploy'">
            <span v-if="isIndependentDeployMap[record.isIndependentDeploy]">{{
              isIndependentDeployMap[record.isIndependentDeploy].dictLabel
            }}</span>
            <span v-else>{{ record.isIndependentDeploy }}</span>
          </template>
        </template></BasicTable
      ></BasicModal
    >
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import {
    specialLineDemandColumns,
    searchDemandFormSchema,
    providersMap,
    bizTypesMap,
    deviceTypeMap,
    sameAreaReplaceMap,
    hddTransformMap,
    isIndependentDeployMap,
  } from './data';
  import { getRealDeliveryBw, GetBatchBindDemandList, BatchBindDemand } from '@/api/business/k';
  import { ExclamationCircleFilled } from '@ant-design/icons-vue';
  import { message, Tooltip, Badge } from 'ant-design-vue';
  import { reactive } from 'vue';
  import { useAreaSelect } from '@/utils/kAreaSelect';
  import dayjs from 'dayjs';

  defineOptions({ name: 'SpecialLineBatchDemandModal' });
  const emit = defineEmits(['register', 'success']);
  let ids: Array<number> = [];
  const prop = defineProps({ bizType: { type: String, required: true } });
  const allData = reactive({
    realDeliveryBw: 0 as number,
  });

  // 判断是否为次月需求（开始时间大于当前月份）
  const isFutureMonth = (startTime: string) => {
    if (!startTime) return false;
    return dayjs(startTime).isAfter(dayjs().endOf('month'));
  };

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    ids = data.ids;
    setModalProps({
      title: `批量绑定需求单号`,
      width: 1300,
      height: 520,
      destroyOnClose: true,
    });
    setProps({
      title: `需求单列表`,
    });

    setSelectedRowKeys([]);
    const res = await getRealDeliveryBw({ ids: ids });
    allData.realDeliveryBw = res.realDeliveryBw || 0;

    // modal 打开时初始化区域数据
    setTimeout(() => {
      initAreaSelect();
    }, 200);
  });

  const [registerTable, { setSelectedRowKeys, getSelectRowKeys, setProps, getForm }] = useTable({
    title: '需求单列表',
    api: GetBatchBindDemandList,
    columns: specialLineDemandColumns(),
    beforeFetch: (params) => {
      params.deliveryInfoIds = ids;
      params.bizType = prop.bizType;
      return params;
    },
    // size: 'small',
    canResize: true,
    scroll: { y: 260 },
    formConfig: {
      labelWidth: 120,
      schemas: searchDemandFormSchema(prop.bizType),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
      baseRowStyle: {
        marginBottom: 0,
      },
    },
    // afterFetch: (data) => {
    //   data.forEach((item) => {
    //     // 用户没有手动选择
    //     if (item.demandId == deliveryRecord.demandId && !userSelected) {
    //       nextTick(() => {
    //         setSelectedRowKeys([item.demandId]);
    //       });
    //     }
    //   });
    //   return data;
    // },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      // pageSizeOptions: ['1', '2', '5'],
    },
    rowKey: 'demandId',
    rowSelection: {
      type: 'radio',
      getCheckboxProps: (record) => {
        // 带宽不足禁用
        const isBwInsufficient = record.canBoundedBw < allData.realDeliveryBw;
        // 需求开始时间大于当前月份禁用（次月需求不能当月提交）
        const isFutureMonthDisabled = record.startTime
          ? dayjs(record.startTime).isAfter(dayjs().endOf('month'))
          : false;

        return {
          disabled: isBwInsufficient || isFutureMonthDisabled,
        };
      },
    },
    clickToRowSelect: false,
    showSelectionBar: true,
  });
  async function handleSubmit() {
    try {
      setModalProps({ confirmLoading: true });
      const selectKeys = getSelectRowKeys();
      if (selectKeys.length == 0) {
        message.warn('请选择需求单号');
        return;
      }
      const data = {
        BizType: prop.bizType,
        ids: ids,
        demandId: selectKeys.length == 1 ? selectKeys[0] : '',
      };
      await BatchBindDemand(data);
      message.success('编辑成功');
      emit('success');
      closeModal();
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
  function handleSelectionChange(data) {
    // 因为此处是单选框，所以需要把原来的值给删除，否则第二次点击的时候单选框无法选中
    if (typeof data.keys == 'object' && data.keys.length > 1) {
      setSelectedRowKeys([data.keys[1]]);
    }
  }
  // 区域联动
  const initAreaSelect = async () => {
    try {
      const form = await getForm();
      if (form) {
        // 区域联动筛选
        const { initAreaData } = useAreaSelect({
          form,
          fields: {
            area: 'area',
            province: 'province',
          },
        });
        await initAreaData();
      }
    } catch (error) {
      console.error('初始化区域数据失败:', error);
    }
  };
</script>
