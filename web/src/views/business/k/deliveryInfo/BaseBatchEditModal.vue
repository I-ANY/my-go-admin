<template>
  <BasicModal v-bind="$attrs" @register="registerModal" :title="getTitle" @ok="handleSubmit">
    <!-- 异网调度信息修改提示 -->
    <div class="diff-isp-notice">
      <Alert
        message="修改异网调度信息会自动触发异网跨省下发，如需手动下发请前往异网跨省下发页面操作"
        type="warning"
        show-icon
        :closable="false"
      />
    </div>
    <!-- 批量编辑模式：显示选中设备数量 -->
    <!-- <div v-if="isBatchMode" class="batch-edit-header">
      <h3>已选择 {{ selectedDeviceCount }} 个设备</h3>
    </div> -->
    <BasicForm @register="registerForm" class="delivery-info-form" />
  </BasicModal>
</template>

<script lang="ts" setup>
  import { computed } from 'vue';
  import type { PropType } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { message, Alert } from 'ant-design-vue';
  import { batchUpdateDeliveryInfo, GetAreaList } from '@/api/business/k';

  defineOptions({ name: 'BaseBatchEditModal' });

  const emit = defineEmits(['success', 'register']);

  interface DeviceInfo {
    id: number;
    deliveryBw?: number;
    bwCount?: number;
    singleDiverybw?: number;
    province?: string;
    city?: string;
    isProvinceScheduling?: number;
    isCoverDiffIsp?: number;
    remark?: string;
    [key: string]: any;
  }

  const props = defineProps({
    bizType: { type: String, required: true },
    formSchema: { type: Array, required: true },
    showDifIspButton: { type: Boolean, default: false },
    needAreaInfo: { type: Boolean, default: false },
    modalWidth: { type: Number, default: 800 },
    modalHeight: { type: Number, default: 600 },
    isBatchMode: { type: Boolean, default: true },
    selectedDevices: { type: Array as PropType<DeviceInfo[]>, default: () => [] },
  });

  const selectedDeviceCount = computed(() => props.selectedDevices.length);
  let areaInfo: Record<string, { label: string; value: string }[]> = {};
  let provinceOptions: any[] = [];

  const [
    registerForm,
    { setFieldsValue, resetFields, validate, updateSchema, getFieldsValue, clearValidate },
  ] = useForm({
    labelWidth: 120,
    baseColProps: { lg: 24, md: 24 },
    schemas: getDynamicFormSchema(),
    showActionButtonGroup: false,
    actionColOptions: {
      span: 24,
    },
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async () => {
    await resetFields();

    // 获取地区信息（如果需要）
    if (props.needAreaInfo) {
      provinceOptions = [];
      areaInfo = {};
      let res = await GetAreaList({});
      areaInfo = buildAreaInfo(res.items);

      Object.keys(areaInfo).forEach((key) => {
        let option = {
          label: key,
          value: key,
        };
        provinceOptions.push(option);
      });

      await updateSchema([
        {
          field: 'province',
          componentProps: {
            options: provinceOptions,
          },
        },
        {
          field: 'city',
          componentProps: {
            options: [],
          },
        },
      ]);
    }

    setModalProps({
      confirmLoading: false,
      destroyOnClose: true,
      width: props.modalWidth,
      height: props.modalHeight,
    });

    // 批量编辑模式预填充表单数据
    if (props.isBatchMode) {
      await prefillBatchFormData();
    }
  });

  const getTitle = computed(() => {
    if (props.isBatchMode) {
      return `批量编辑 (${selectedDeviceCount.value}个设备)`;
    }
    return '批量编辑';
  });

  // 批量表单数据预填充
  async function prefillBatchFormData() {
    if (!props.selectedDevices || props.selectedDevices.length === 0) {
      return;
    }

    // 检查设备字段是否一致，不一致则显示空值
    // const firstDevice = props.selectedDevices[0];
    const fieldsToCheck = [
      'deliveryBw',
      'bwCount',
      'singleDiverybw',
      'province',
      'city',
      'isProvinceScheduling',
      'isCoverDiffIsp',
      'remark',
    ];

    const prefillData: any = {};

    for (const field of fieldsToCheck) {
      const values = props.selectedDevices.map((device: DeviceInfo) => device[field]);
      const allSame = values.every((value) => value === values[0]);

      if (allSame && values[0] !== undefined && values[0] !== null) {
        prefillData[field] = values[0];
      } else {
        prefillData[field] = null; // 不一致时显示空值
      }
    }

    await setFieldsValue(prefillData);

    // 根据调度控制状态初始化省份城市字段
    if (props.needAreaInfo && prefillData.isProvinceScheduling !== undefined) {
      if (prefillData.isProvinceScheduling == 1) {
        // 仅本省：禁用省份城市
        await updateSchema([
          {
            field: 'province',
            componentProps: {
              disabled: true,
              placeholder: '请选择省份',
              options: provinceOptions,
            },
            required: false,
          },
          {
            field: 'city',
            componentProps: {
              disabled: true,
              options:
                prefillData.province && areaInfo[prefillData.province]
                  ? areaInfo[prefillData.province]
                  : [],
            },
            required: false,
          },
        ]);
      } else {
        // 不限制：启用省份城市
        await updateSchema([
          {
            field: 'province',
            componentProps: {
              disabled: false,
              placeholder: '不填为"全国"',
              options: provinceOptions,
            },
            required: false,
          },
          {
            field: 'city',
            componentProps: {
              disabled: false,
              options:
                prefillData.province && areaInfo[prefillData.province]
                  ? areaInfo[prefillData.province]
                  : [],
            },
            required: prefillData.province ? true : false,
          },
        ]);
      }
    }
  }

  // 批量编辑提交功能
  async function handleSubmit() {
    try {
      const values = await validate();
      setModalProps({ confirmLoading: true });

      // 处理省份为空的情况，设置为'全国'
      if (props.needAreaInfo && !values.province) {
        values.province = '全国';
        values.city = null; // 省份为全国时，城市清空
      }

      // 构建批量更新数据
      const { bizType, ...updateData } = values;
      console.log('bizType', bizType);

      // 类型转换：将数字字段转换为int类型
      if (
        updateData.isProvinceScheduling !== undefined &&
        updateData.isProvinceScheduling !== null
      ) {
        updateData.isProvinceScheduling = parseInt(updateData.isProvinceScheduling, 10);
      }
      if (updateData.isCoverDiffIsp !== undefined && updateData.isCoverDiffIsp !== null) {
        updateData.isCoverDiffIsp = parseInt(updateData.isCoverDiffIsp, 10);
      }
      if (updateData.difIsp !== undefined && updateData.difIsp !== null) {
        updateData.difIsp = parseInt(updateData.difIsp, 10);
      }

      // 如果有单线带宽变化，计算总交付带宽
      if (updateData.singleDiverybw !== undefined && updateData.singleDiverybw !== null) {
        updateData.singleDiverybw = parseFloat(updateData.singleDiverybw);
        // 如果线路数也存在且一致，则计算总交付带宽
        if (updateData.bwCount !== undefined && updateData.bwCount !== null) {
          updateData.bwCount = parseInt(updateData.bwCount, 10);
          updateData.deliveryBw = updateData.bwCount * updateData.singleDiverybw;
        }
      }

      if (updateData.bwCount !== undefined && updateData.bwCount !== null) {
        updateData.bwCount = parseInt(updateData.bwCount, 10);
      }

      const batchData = {
        ids: props.selectedDevices.map((device: DeviceInfo) => device.id),
        bizType: props.bizType,
        updateData: {
          ...updateData,
        },
      };

      // 调用批量更新接口
      const response = await batchUpdateDeliveryInfo(batchData);

      // 解析业务数据
      const { successCount, failedCount, failedIds, errorMsg } = response;
      console.log('批量更新结果:', response);

      if (failedCount > 0) {
        // 部分失败的情况
        message.warning(`批量更新完成：成功 ${successCount} 个，失败 ${failedCount} 个`);

        // 显示失败的设备信息
        if (failedIds && failedIds.length > 0) {
          console.error('批量更新失败的设备ID:', failedIds.join(', '));
        }
        if (errorMsg) {
          console.error('批量更新错误信息:', errorMsg);
        }
      } else {
        // 全部成功
        message.success(`批量更新成功，共更新 ${successCount} 个设备`);
      }

      closeModal();
      emit('success');
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }

  // 省份调度控制变化
  async function onIsProvinceSchedulingChange(value) {
    if (props.needAreaInfo) {
      const data = await getFieldsValue();
      // 仅本省
      if (value == 1) {
        await updateSchema([
          {
            field: 'province',
            componentProps: {
              disabled: true,
              placeholder: '请选择省份',
              options: provinceOptions,
            },
            required: false,
          },
          {
            field: 'city',
            componentProps: {
              disabled: true,
              options: data.province && areaInfo[data.province] ? areaInfo[data.province] : [],
            },
            required: false,
          },
        ]);
      } else {
        await updateSchema([
          {
            field: 'province',
            componentProps: {
              disabled: false,
              placeholder: '不填为"全国"',
              options: provinceOptions,
            },
            required: false,
          },
          {
            field: 'city',
            componentProps: {
              disabled: false,
              options: data.province && areaInfo[data.province] ? areaInfo[data.province] : [],
            },
            required: data.province ? true : false,
          },
        ]);
      }
      clearValidate(['city', 'province']);
    }
  }

  // 省份变化处理
  async function onProvinceChange(value) {
    if (props.needAreaInfo) {
      let data = await getFieldsValue();
      // 不限制
      if (data.isProvinceScheduling == 0) {
        // 用户主动改变省份时清空城市
        if (data.province !== value) {
          await setFieldsValue({
            city: null,
          });
        }

        if (value) {
          await updateSchema([
            {
              field: 'city',
              componentProps: {
                disabled: false,
                options: areaInfo[value],
              },
              required: true,
            },
          ]);
        } else {
          await updateSchema([
            {
              field: 'city',
              componentProps: {
                disabled: false,
                options: [],
              },
              required: false,
            },
          ]);
        }
      }
      clearValidate(['city', 'province']);
    }
  }

  // 是否异网变化处理
  async function onIsCoverDiffIspChange(value) {
    if (value === 1) {
      // 选择"是"时，异网运营商字段变为必填
      await updateSchema([
        {
          field: 'difIsp',
          required: true,
        },
      ]);
    } else {
      // 选择"否"时，异网运营商字段变为非必填并清空
      await updateSchema([
        {
          field: 'difIsp',
          required: false,
        },
      ]);
      await setFieldsValue({
        difIsp: null,
      });
    }
    clearValidate(['difIsp']);
  }

  // 构建地区信息
  function buildAreaInfo(
    rawAreaInfo: Array<any>,
  ): Record<string, { label: string; value: string }[]> {
    let data: Record<string, { label: string; value: string }[]> = {};
    for (let i = 0; i < rawAreaInfo?.length; i++) {
      let item = rawAreaInfo[i];
      let province = item.province_name;
      if (!data[province]) {
        data[province] = [];
      }
      data[province].push({ label: item.city_name, value: item.city_name });
    }
    return data;
  }

  // 动态生成表单schema
  function getDynamicFormSchema() {
    return props.formSchema.map((schema: any) => {
      const newSchema = { ...schema };

      // 为需要的字段绑定事件处理函数
      if (schema.field === 'singleDiverybw') {
        newSchema.componentProps = {
          ...schema.componentProps,
          onChange: onSingleDiverybwChange,
        };
      }

      if (schema.field === 'isProvinceScheduling') {
        newSchema.componentProps = {
          ...schema.componentProps,
          onChange: onIsProvinceSchedulingChange,
        };
      }

      if (schema.field === 'province') {
        newSchema.componentProps = {
          ...schema.componentProps,
          onChange: onProvinceChange,
        };
      }

      if (schema.field === 'isCoverDiffIsp') {
        newSchema.componentProps = {
          ...schema.componentProps,
          onChange: onIsCoverDiffIspChange,
        };
      }

      return newSchema;
    });
  }

  // 单条线路带宽变化计算
  async function onSingleDiverybwChange(value) {
    const data = await getFieldsValue();
    if (data.bwCount && data.bwCount > 0) {
      let d: any = {};
      if (value != null) {
        d.deliveryBw = value * data.bwCount;
      } else {
        d.deliveryBw = null;
      }
      await setFieldsValue(d);
    }
  }

  // 暴露方法给子组件使用
  defineExpose({
    onSingleDiverybwChange,
    onIsProvinceSchedulingChange,
    onProvinceChange,
    onIsCoverDiffIspChange,
    buildAreaInfo,
  });
</script>

<style lang="less" scoped>
  .batch-edit-header {
    margin-bottom: 16px;
    padding: 12px;
    border-radius: 6px;
    background-color: #f0f2f5;
    text-align: center;

    h3 {
      margin: 0;
      color: #1890ff;
      font-size: 16px;
      font-weight: 500;
    }
  }

  .diff-isp-notice {
    margin-bottom: 16px;
  }

  .delivery-info-form {
    /deep/ .ant-form-item-label {
      align-items: flex-start !important;
      min-width: 120px !important;
      padding-top: 1px !important;
      line-height: 1 !important;
      white-space: normal !important;
    }
  }
</style>
