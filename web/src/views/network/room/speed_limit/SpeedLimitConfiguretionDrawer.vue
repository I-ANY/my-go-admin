<template>
  <BasicDrawer
    v-bind="$attrs"
    @register="register"
    title="编辑机房限速配置"
    :width="650"
    @ok="handleSubmit"
    destroyOnClose
    :maskClosable="true"
    :showFooter="true"
    :showCancelBtn="true"
    :showOkBtn="true"
    @close="onDrawerClose"
  >
    <Spin :spinning="allData.loading">
      <BasicForm @register="registerForm" />
      <div ref="businessSortContainerRef" class="businessSortContainer">
        <div
          v-for="(value, i) in allData.businessSorts"
          :key="`${value}_${i}`"
          style="display: inline-block; margin-right: 3px"
        >
          <Tag class="sortable-tag" :closable="true" color="blue" @close="handleTagClose(value)">{{
            value
          }}</Tag>
        </div>
      </div>
      <div>
        <Select
          v-if="allData.selectVisible"
          ref="businessSelectRef"
          style="width: 140px"
          @blur="handleBusinessConfirm"
          placeholder="添加低优业务"
          v-model:value="allData.inputValue"
          :allowClear="true"
          :showSearch="true"
        >
          <SelectOption v-for="(value, i) in getCanSelectOptions()" :key="i" :value="value">{{
            value
          }}</SelectOption
          >"
        </Select>
        <Tag v-else style="border-style: dashed; background: #fff" @click="showBusinessSelect">
          <plus-outlined />
          添加低优业务
        </Tag>
      </div>
    </Spin>
  </BasicDrawer>
</template>

<script lang="ts" setup>
  import { BasicDrawer, useDrawerInner } from '@/components/Drawer';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { limitStatusEnum, speedLimitConfiguretionFormSchema, switchIsMainEnum } from './data';
  import { message, Spin, Tag, Select, SelectOption, Modal } from 'ant-design-vue';
  import {
    getPeakRoomTag,
    updatePeakRoomSpeedLimitConfig,
    getRoomSwitch,
  } from '@/api/network/room_peak';
  import { customCeilDivide } from '@/utils/util';
  import { nextTick, reactive, ref, createVNode, h } from 'vue';
  import type { Rule } from 'ant-design-vue/es/form';
  import Sortable from 'sortablejs';
  import { ExclamationCircleOutlined, PlusOutlined } from '@ant-design/icons-vue';

  defineOptions({
    name: 'SpeedLimitConfiguretionDrawer',
  });
  const businessSortContainerRef = ref<HTMLElement | null>(null);
  const businessSelectRef = ref<HTMLElement | null>(null);

  const emit = defineEmits(['success', 'register']);
  let allData = reactive({
    loading: false,
    record: {} as any,
    businessSorts: [] as string[],
    options: [] as string[],
    selectVisible: false,
    inputValue: null as any,
    popconfirmVisible: {} as any,
    allOnlineSwitches: [] as any,
  });
  const [
    registerForm,
    { validate, setFieldsValue, updateSchema, getFieldsValue, clearValidate, validateFields },
  ] = useForm({
    labelWidth: 130,
    schemas: speedLimitConfiguretionFormSchema(
      onStatusChange,
      unlimitThresholdChange,
      onValueChange,
      onExcludeSwitchesChange,
      getFormHandleFn,
    ),
    actionColOptions: {
      span: 24,
    },
    showActionButtonGroup: false,
  });

  const [register, { setDrawerProps, closeDrawer }] = useDrawerInner(async (data) => {
    try {
      allData.loading = true;
      allData.record = {
        ...data.record,
      };
      setDrawerProps({
        title: `编辑机房限速配置：${data.record.name}`,
      });
      const speedLimitConfig = data.record.speedLimitConfig || {};
      allData.businessSorts = [...(speedLimitConfig.businessSorts || [])];
      // 获排除交换机
      const excludeSwitchIds: number[] = [];
      speedLimitConfig.excludeSwitches?.forEach((item) => {
        excludeSwitchIds.push(item.id);
      });
      // 获取总流量交换机
      const totalTrafficSwitchIds: number[] = [];
      speedLimitConfig.totalTrafficSwitches?.forEach((item) => {
        totalTrafficSwitchIds.push(item.id);
      });
      // 获取扣减流量交换机
      const deductTrafficSwitchIds: number[] = [];
      speedLimitConfig.deductTrafficSwitches?.forEach((item) => {
        deductTrafficSwitchIds.push(item.id);
      });
      // 获取所有交换机信息
      const { items: switches } = await getRoomSwitch(data.record.id, { status: 3 });
      allData.allOnlineSwitches = switches;
      await updateSchema({
        field: 'excludeSwitches',
        componentProps: {
          options: allData.allOnlineSwitches
            .filter((s) => s.isMain == switchIsMainEnum.NO)
            .map((s) => {
              return {
                label: s.description,
                value: s.id,
              };
            }),
        },
      });

      // 机房可所有的tag
      const res = await getPeakRoomTag(allData.record.id);
      if (res.tags && res.tags.length > 0) {
        allData.options = res.tags;
      } else {
        allData.options = [];
      }
      convertValue(speedLimitConfig);
      nextTick(async () => {
        await setFieldsValue({
          status: speedLimitConfig.status == 1 ? 1 : 0,
          limitValue: speedLimitConfig.limitValue,
          limitThreshold: speedLimitConfig.limitThreshold,
          unlimitThreshold: speedLimitConfig.unlimitThreshold,
          excludeSwitches: excludeSwitchIds,
          totalTrafficSwitches: totalTrafficSwitchIds,
          deductTrafficSwitches: deductTrafficSwitchIds,
        });
        await resetOptions();
      });
      initSorter();
    } finally {
      onStatusChange();
      allData.loading = false;
    }
  });
  function onStatusChange() {
    const value = getFieldsValue();
    if (value?.status == limitStatusEnum.LIMIT) {
      updateSchema([
        {
          field: 'unlimitThreshold',
          required: true,
        },
        {
          field: 'limitThreshold',
          required: true,
        },
        {
          field: 'limitValue',
          required: true,
        },
      ]);
    } else {
      updateSchema([
        {
          field: 'unlimitThreshold',
          required: false,
        },
        {
          field: 'limitThreshold',
          required: false,
        },
        {
          field: 'limitValue',
          required: false,
        },
      ]);
      nextTick(() => {
        onValueChange();
      });
    }
  }
  function unlimitThresholdChange(_rule: Rule, _value: number) {
    const values = getFieldsValue();
    if (!values.unlimitThreshold && !values.limitThreshold && !values.limitValue) {
      return Promise.resolve();
    }
    if (!values.unlimitThreshold || !values.limitThreshold || !values.limitValue) {
      Promise.reject('不满足：解除限速阈值<=限速值<=开启限速阈值');
    }
    if (
      !(values.unlimitThreshold <= values.limitValue && values.limitValue <= values.limitThreshold)
    ) {
      return Promise.reject('不满足：解除限速阈值<=限速值<=开启限速阈值');
    }
    return Promise.resolve();
  }
  function onValueChange() {
    clearValidate();
    validate();
  }
  function initSorter() {
    if (businessSortContainerRef.value) {
      Sortable.create(businessSortContainerRef.value, {
        onEnd: (evt) => {
          const { newIndex, oldIndex } = evt;
          if (newIndex === oldIndex) {
            return;
          }
          if (newIndex == undefined || oldIndex == undefined) {
            return;
          }
          // 创建一个新的数组副本以避免直接修改响应式数据
          const data = [...allData.businessSorts];

          // 保存被移动的元素
          const movedItem = data.splice(oldIndex, 1)[0];

          // 将元素插入到新位置
          data.splice(newIndex, 0, movedItem);

          // 更新原始数组
          allData.businessSorts = data;
        },
      });
    } else {
      console.log('xxxxxxxxxxxxxxx');
    }
  }
  function handleBusinessConfirm() {
    const inputValue = allData.inputValue;
    let tags = [...allData.businessSorts];
    if (inputValue && tags.indexOf(inputValue) === -1) {
      tags.push(inputValue);
      allData.businessSorts = tags;
    }
    allData.selectVisible = false;
    allData.inputValue = null;
  }
  function showBusinessSelect() {
    allData.selectVisible = true;
    nextTick(() => {
      businessSelectRef?.value?.focus();
    });
  }

  function getCanSelectOptions() {
    const data = allData.options.filter((tag) => !allData.businessSorts.includes(tag));
    return data;
  }
  function handleTagClose(t) {
    const data = allData.businessSorts.filter((tag) => tag !== t);
    allData.businessSorts = data;

    // handlePopconfirmCancel(t);
  }
  // function onTagCloseClick(value) {
  //   allData.popconfirmVisible[value] = true;
  // }
  // function handlePopconfirmCancel(value) {
  //   allData.popconfirmVisible[value] = false;
  // }
  // b -> Gb
  function parseValue(data: Recordable) {
    if (data.limitValue == null || data.limitValue == undefined) {
      data.limitValue = null;
    } else {
      data.limitValue = parseInt(data.limitValue * 1000 * 1000 * 1000);
    }
    if (data.limitThreshold == null || data.limitThreshold == undefined) {
      data.limitThreshold = null;
    } else {
      data.limitThreshold = parseInt(data.limitThreshold * 1000 * 1000 * 1000);
    }
    if (data.unlimitThreshold == null || data.unlimitThreshold == undefined) {
      data.unlimitThreshold = null;
    } else {
      data.unlimitThreshold = parseInt(data.unlimitThreshold * 1000 * 1000 * 1000);
    }
  }

  // Gb -> b, 保留2位小数
  function convertValue(params: Recordable) {
    if (params.limitValue == null || params.limitValue == undefined) {
      params.limitValue = null;
    } else {
      params.limitValue = parseFloat(
        customCeilDivide(params.limitValue, 1000 * 1000 * 1000, 2) as any,
      );
    }
    if (params.limitThreshold == null || params.limitThreshold == undefined) {
      params.limitThreshold = null;
    } else {
      params.limitThreshold = parseFloat(
        customCeilDivide(params.limitThreshold, 1000 * 1000 * 1000, 2) as any,
      );
    }
    if (params.unlimitThreshold == null || params.unlimitThreshold == undefined) {
      params.unlimitThreshold = null;
    } else {
      params.unlimitThreshold = parseFloat(
        customCeilDivide(params.unlimitThreshold, 1000 * 1000 * 1000, 2) as any,
      );
    }
  }
  function onDrawerClose() {
    emit('success');
  }
  async function handleSubmit() {
    const sourceValue = await validate();
    const values = { ...sourceValue };
    values.businessSorts = allData.businessSorts;
    parseValue(values);

    // 获取已选择的排除交换机名字
    const switchNames: string[] = [];
    allData.allOnlineSwitches.forEach((item) => {
      if (values?.excludeSwitches?.includes(item.id)) {
        switchNames.push(item.description);
      }
    });

    // 获取已选择的总流量交换机名字
    const totalTrafficSwitchNames: string[] = [];
    allData.allOnlineSwitches.forEach((item) => {
      if (values?.totalTrafficSwitches?.includes(item.id)) {
        totalTrafficSwitchNames.push(item.description);
      }
    });
    // 获取已选择的流量扣减交换机名字
    const deductTrafficSwitchNames: string[] = [];
    allData.allOnlineSwitches.forEach((item) => {
      if (values?.deductTrafficSwitches?.includes(item.id)) {
        deductTrafficSwitchNames.push(item.description);
      }
    });

    if (
      values.status == limitStatusEnum.LIMIT &&
      (!values.businessSorts || values.businessSorts.length == 0)
    ) {
      message.error('请选择需要限速的业务');
      return;
    }
    if (values.status == limitStatusEnum.LIMIT) {
      Modal.confirm({
        title: '是否确认开启机房限速?',
        icon: createVNode(ExclamationCircleOutlined),
        maskClosable: true,
        width: 650,
        content: h('div', [
          h(
            'p',
            { style: 'color:red;margin:2px;font-weight: bold;' },
            '限速值（Gbps）：' + sourceValue.limitValue,
          ),
          h(
            'p',
            { style: 'color:red;margin:2px;font-weight: bold;' },
            '开启限速阈值（Gbps）：' + sourceValue.limitThreshold,
          ),
          h(
            'p',
            { style: 'color:red;margin:2px;font-weight: bold;' },
            '解除限速阈值（Gbps）：' + sourceValue.unlimitThreshold,
          ),
          h(
            'p',
            { style: 'color:red;margin:2px;font-weight: bold;' },
            '排除交换机：' + switchNames.join('；'),
          ),
          h(
            'p',
            { style: 'color:red;margin:2px;font-weight: bold;' },
            '总流量交换机：' + totalTrafficSwitchNames.join('；'),
          ),
          h(
            'p',
            { style: 'color:red;margin:2px;font-weight: bold;' },
            '流量扣减交换机：' + deductTrafficSwitchNames.join('；'),
          ),
          h(
            'p',
            { style: 'color:red;margin:2px;font-weight: bold;' },
            '限速顺序：' + values.businessSorts.join(' -> '),
          ),
        ]),
        async onOk() {
          await submit(values);
        },
      });
    } else {
      await submit(values);
    }
  }
  async function submit(values) {
    try {
      setDrawerProps({
        loading: true,
      });
      await updatePeakRoomSpeedLimitConfig(allData.record.id, {
        status: values.status,
        limitValue: values.limitValue,
        limitThreshold: values.limitThreshold,
        unlimitThreshold: values.unlimitThreshold,
        businessSorts: values.businessSorts.length > 0 ? values.businessSorts : null,
        excludeSwitches: values.excludeSwitches,
        totalTrafficSwitches: values.totalTrafficSwitches,
        deductTrafficSwitches: values.deductTrafficSwitches,
      });
      message.success('操作成功');
      emit('success');
      closeDrawer();
    } finally {
      setDrawerProps({
        loading: false,
      });
    }
  }
  function onExcludeSwitchesChange() {
    // 重新设置总流量交换机options和值
    resetOptions();
  }
  async function resetOptions() {
    // 重新设置总流量交换机options
    const values = await getFieldsValue();
    const excludeSwitches: number[] = values.excludeSwitches || [];
    const switchOptions = allData.allOnlineSwitches
      .filter((s) => !excludeSwitches.includes(s.id)) // 过滤掉排除的交换机，排除后的交换机不能再被选择
      .map((s) => {
        return {
          label: s.description,
          value: s.id,
        };
      });

    await updateSchema([
      {
        field: 'totalTrafficSwitches',
        componentProps: {
          options: switchOptions,
        },
      },
      {
        field: 'deductTrafficSwitches',
        componentProps: {
          options: switchOptions,
        },
      },
    ]);
    // 重新设置总流量交换机的值，过滤掉排除的交换机
    await setFieldsValue({
      totalTrafficSwitches: (values.totalTrafficSwitches || []).filter(
        (s) => !excludeSwitches.includes(s), // 过滤掉排除的交换机，排除后的交换机不能再被选择
      ),
      deductTrafficSwitches: (values.deductTrafficSwitches || []).filter(
        (s) => !excludeSwitches.includes(s), // 过滤掉排除的交换机，排除后的交换机不能再被选择
      ),
    });
  }
  function getFormHandleFn() {
    return {
      getFieldsValue,
      updateSchema,
      setFieldsValue,
      validateFields,
    };
  }
</script>
<style scoped>
  .business-sort-container {
    cursor: pointer; /* 整个容器显示拖拽光标 */
  }

  .sortable-tag {
    margin-right: 8px;
    margin-bottom: 8px;
    cursor: pointer; /* 每个标签显示拖拽光标 */
    user-select: none; /* 防止拖拽时选中文本 */
  }

  .sortable-tag:hover {
    cursor: pointer; /* 悬停时确保显示手型光标 */
  }

  .businessSortContainer {
    min-height: 70px; /* 确保容器有最小高度 */
    margin-bottom: 10px;
    padding: 8px;
    border: 1px solid gainsboro;
    border-radius: 8px; /* 更明显且一致的圆角 */
  }
</style>
