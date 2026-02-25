<template>
  <BasicDrawer
    v-bind="$attrs"
    @register="register"
    title="编辑机房信息"
    :width="600"
    @ok="handleSubmit"
    destroyOnClose
    :maskClosable="true"
    :showFooter="true"
    :showCancelBtn="true"
    :showOkBtn="true"
  >
    <BasicForm @register="registerForm" />
    <div class="plans-container">
      <div
        v-for="(value, i) in allData.plans"
        :key="`${value}_${i}_${Math.random() * 1000000}`"
        style="display: inline-block; margin-right: 3px"
      >
        <Tooltip :title="value" :overlayStyle="{ maxWidth: '550px' }">
          <Tag
            class="sortable-tag"
            :closable="true"
            color="blue"
            :style="{ fontSize: '11px' }"
            @close="handleTagClose(value, i)"
            >{{ value.length > 70 ? value.slice(0, 67) + '...' : value }}</Tag
          >
        </Tooltip>
      </div>
    </div>
    <div>
      <FormItem
        v-if="allData.planInputVisible"
        :validateStatus="allData.plainValidateStatus"
        :help="allData.plainValidateHelp"
      >
        <Input
          ref="planInputRef"
          style="width: 98%"
          placeholder="添加方案，如需填写打峰方案请以“打峰方案-”开头"
          v-model:value="allData.planInputValue"
          :allowClear="true"
          @blur="handlePlanConfirm"
          @press-enter="handlePlanConfirm"
          @change="handlePlanChange"
          :maxlength="150"
        />
      </FormItem>
      <Tag v-else style="border-style: dashed; background: #fff" @click="showPlanInput">
        <plus-outlined />
        添加方案
      </Tag>
    </div>
  </BasicDrawer>
</template>

<script lang="ts" setup>
  import { BasicDrawer, useDrawerInner } from '@/components/Drawer';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { roomUseModEnum, updateRoomPeakFormSchema, validatePeakPlanFormat } from './data';
  import { Input, message, Tag, Tooltip, FormItem } from 'ant-design-vue';
  import { updatePeakRoom } from '@/api/network/room_peak';
  import { customCeilDivide } from '@/utils/util';
  import { nextTick, reactive, ref } from 'vue';
  import { PlusOutlined } from '@ant-design/icons-vue';

  defineOptions({
    name: 'UpdateRoomPeakDrawer',
  });
  const allData = reactive({
    plans: [] as string[],
    planInputValue: '' as string,
    planInputVisible: false,
    plainValidateStatus: 'success',
    plainValidateHelp: '',
  });
  const planInputRef = ref<HTMLElement | null>(null);
  const emit = defineEmits(['success', 'register']);
  let id: number = 0;
  let maxBw = 0;
  const [registerForm, { validate, setFieldsValue, setProps, updateSchema, getFieldsValue }] =
    useForm({
      labelWidth: 120,
      actionColOptions: {
        span: 24,
      },
      showActionButtonGroup: false,
    });

  const [register, { setDrawerProps, closeDrawer }] = useDrawerInner((data) => {
    const record = {
      ...data.record,
    };
    allData.plans = record.plans || [];
    id = record.id;
    maxBw = record.bandwidth || 0;
    setDrawerProps({
      title: `编辑机房信息：${record.name}`,
    });
    if (record) {
      convertValue(record);
      maxBw = record.bandwidth;
      nextTick(() => {
        setProps({
          schemas: updateRoomPeakFormSchema(maxBw, getFormFn),
        });

        setFieldsValue({
          guaranteedBw: record.guaranteedBw,
          preset95Bw: record.preset95Bw,
          billingStartAt: record.billingStartAt,
          remark: record.remark,
          bandwidth: record.bandwidth,
          useMode: record.useMode,
          canPeakShavingBw: record.canPeakShavingBw,
          singlePortPreset95Bw: record.singlePortPreset95Bw,
        });
      });
    }
  });

  function getFormFn() {
    return {
      getFieldsValue,
      setFieldsValue,
      validate,
      updateSchema,
    };
  }
  async function handleSubmit() {
    try {
      setDrawerProps({
        loading: true,
      });
      if (!handlePlanChange(null)) {
        return;
      }
      const values = await validate();
      parseValue(values);
      let requestBody = {
        guaranteedBw: values.guaranteedBw,
        preset95Bw: values.preset95Bw,
        billingStartAt: values.billingStartAt,
        remark: values.remark,
        useMode: values.useMode,
        plans: allData.plans,
      };
      // 单端口95削峰才能编辑可削峰带宽字段，否则由后端计算
      if (values.useMode == roomUseModEnum.SINGLE_PORT_95_PEAK_SHAVING) {
        requestBody['canPeakShavingBw'] = values.canPeakShavingBw; // 可削峰带宽
        requestBody['singlePortPreset95Bw'] = values.singlePortPreset95Bw; // 单口预95带宽
      }
      await updatePeakRoom(id, requestBody);
      message.success('操作成功');
      emit('success');
      closeDrawer();
    } finally {
      setDrawerProps({
        loading: false,
      });
    }
  }

  // b -> Gb
  function parseValue(data: Recordable) {
    if (data.guaranteedBw == null || data.guaranteedBw == undefined) {
      data.guaranteedBw = null;
    } else {
      data.guaranteedBw = parseInt(data.guaranteedBw * 1000 * 1000 * 1000);
    }
    if (data.preset95Bw == null || data.preset95Bw == undefined) {
      data.preset95Bw = null;
    } else {
      data.preset95Bw = parseInt(data.preset95Bw * 1000 * 1000 * 1000);
    }
    // 单端口95削峰-可削峰带宽、单口预95带宽才能编辑，否则由后端计算
    if (data.useMode == roomUseModEnum.SINGLE_PORT_95_PEAK_SHAVING) {
      if (data.canPeakShavingBw == null || data.canPeakShavingBw == undefined) {
        data.canPeakShavingBw = null;
      } else {
        data.canPeakShavingBw = parseInt(data.canPeakShavingBw * 1000 * 1000 * 1000);
      }
      if (data.singlePortPreset95Bw == null || data.singlePortPreset95Bw == undefined) {
        data.singlePortPreset95Bw = null;
      } else {
        data.singlePortPreset95Bw = parseInt(data.singlePortPreset95Bw * 1000 * 1000 * 1000);
      }
    } else {
      data.canPeakShavingBw = null;
      data.singlePortPreset95Bw = null;
    }
  }

  // Gb -> b, 保留2位小数
  function convertValue(params: Recordable) {
    if (params.guaranteedBw == null || params.guaranteedBw == undefined) {
      params.guaranteedBw = null;
    } else {
      params.guaranteedBw = parseFloat(
        customCeilDivide(params.guaranteedBw, 1000 * 1000 * 1000, 2) as any,
      );
    }
    if (params.preset95Bw == null || params.preset95Bw == undefined) {
      params.preset95Bw = null;
    } else {
      params.preset95Bw = parseFloat(
        customCeilDivide(params.preset95Bw, 1000 * 1000 * 1000, 2) as any,
      );
    }
    if (params.bandwidth == null || params.bandwidth == undefined) {
      params.bandwidth = null;
    } else {
      params.bandwidth = parseFloat(
        customCeilDivide(params.bandwidth, 1000 * 1000 * 1000, 2) as any,
      );
    }
    if (params.canPeakShavingBw == null || params.canPeakShavingBw == undefined) {
      params.canPeakShavingBw = null;
    } else {
      params.canPeakShavingBw = parseFloat(
        customCeilDivide(params.canPeakShavingBw, 1000 * 1000 * 1000, 2) as any,
      );
    }
    // 单口可削峰带宽
    if (params.singlePortPreset95Bw == null || params.singlePortPreset95Bw == undefined) {
      params.singlePortPreset95Bw = null;
    } else {
      params.singlePortPreset95Bw = parseFloat(
        customCeilDivide(params.singlePortPreset95Bw, 1000 * 1000 * 1000, 2) as any,
      );
    }
  }

  // 删除方案
  function handleTagClose(t, i) {
    const data: any[] = [];
    allData.plans.forEach((tag, index) => {
      if (!(index === i && tag === t)) {
        data.push(tag);
      }
    });
    allData.plans = data;
  }
  function showPlanInput() {
    allData.planInputVisible = true;
    nextTick(() => {
      planInputRef?.value?.focus();
    });
  }
  // 确认添加方案
  function handlePlanConfirm() {
    if (!allData.planInputValue) {
      return;
    }
    if (!handlePlanChange(null)) {
      return;
    }
    const inputValue = allData.planInputValue;
    let plans = [...allData.plans];
    plans.push(inputValue);
    allData.plans = plans;
    allData.planInputVisible = false;
    allData.planInputValue = '';
  }

  // 方案输入框变化时验证方案格式
  function handlePlanChange(_e: any | null): boolean {
    if (allData.plans.indexOf(allData.planInputValue) !== -1) {
      allData.plainValidateStatus = 'error';
      allData.plainValidateHelp = '方案已存在';
      return false;
    }
    if (validatePeakPlanFormat(allData.planInputValue)) {
      allData.plainValidateStatus = 'success';
      allData.plainValidateHelp = '';
      return true;
    } else {
      allData.plainValidateStatus = 'error';
      allData.plainValidateHelp = '“打峰方案”不满足格式要求';
      return false;
    }
  }
</script>

<style scoped>
  .plans-container {
    min-height: 100px; /* 确保容器有最小高度 */
    margin-bottom: 8px;
    padding: 8px;
    border: 1px solid gainsboro;
    border-radius: 8px; /* 更明显且一致的圆角 */
  }

  .sortable-tag {
    margin-right: 8px;
    margin-bottom: 8px;
  }
</style>
