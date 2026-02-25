<template>
  <BasicModal v-bind="$attrs" @register="registerModal" :title="getTitle" @ok="handleSubmit">
    <!-- 异网调度信息修改提示 -->
    <div class="diff-isp-notice">
      <Alert
        message="修改异网调度信息会自动触发异网跨省下发，如需手动下发请前往异网跨省下发页面操作"
        type="warning"
        show-icon
        :closable="false"
        style="margin-bottom: 16px"
      />
    </div>
    <BasicForm @register="registerForm" class="delivery-info-form" />
  </BasicModal>
</template>

<script lang="ts" setup>
  import { ref, computed, unref } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { message, Alert } from 'ant-design-vue';
  import { getDeliveryInfoList, updateDeliveryInfo, GetAreaList } from '@/api/business/k';

  defineOptions({ name: 'BaseDeliveryInfoModal' });

  const emit = defineEmits(['success', 'register']);

  const props = defineProps({
    bizType: { type: String, required: true },
    formSchema: { type: Array, required: true },
    showDifIspButton: { type: Boolean, default: false },
    needAreaInfo: { type: Boolean, default: false },
    modalWidth: { type: Number, default: 800 },
    modalHeight: { type: Number, default: 600 },
  });

  const isUpdate = ref(true);
  const rowId = ref(0);
  let record: any = {};
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

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
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

    isUpdate.value = !!data?.isUpdate;
    record = data.record || {};

    if (unref(isUpdate)) {
      rowId.value = data.record.id;
      let renderData = {};
      let res = await getDeliveryInfoList({ ids: [data.record.id], bizType: props.bizType });
      if (res.items?.length > 0) {
        record = res.items[0];
        let singleDiverybw: number | string | null = null;
        if (
          typeof record.deliveryBw === 'number' &&
          typeof record.bwCount === 'number' &&
          record.deliveryBw >= 0 &&
          record.bwCount > 0
        ) {
          singleDiverybw = (record.deliveryBw / record.bwCount).toFixed(2);
        }
        renderData = {
          ...record,
          singleDiverybw,
        };
      }
      await setFieldsValue(renderData);

      // 根据调度控制状态初始化省份城市字段
      if (props.needAreaInfo && record.isProvinceScheduling !== undefined) {
        if (record.isProvinceScheduling == 1) {
          // 仅本省：禁用省份城市但保持原值
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
                  record.province && areaInfo[record.province] ? areaInfo[record.province] : [],
              },
              required: false,
            },
          ]);
        } else {
          // 不限制：启用省份城市，并设置城市选项
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
                  record.province && areaInfo[record.province] ? areaInfo[record.province] : [],
              },
              required: record.province ? true : false,
            },
          ]);
        }
      }
    }
  });

  const getTitle = computed(() => {
    return !unref(isUpdate) ? '新增信息' : '编辑信息';
  });

  // 基础提交功能
  async function handleSubmit() {
    try {
      const values = await validate();
      setModalProps({ confirmLoading: true });
      values.bizType = props.bizType;

      // 处理省份为空的情况，设置为'全国'
      if (props.needAreaInfo && !values.province) {
        values.province = '全国';
        values.city = null; // 省份为全国时，城市清空
      }

      await updateDeliveryInfo(rowId.value, values);
      message.success('编辑成功');
      closeModal();
      emit('success');
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }

  // 单条线路带宽变化计算
  async function onSingleDiverybwChange(value) {
    if (record.bwCount && record.bwCount > 0) {
      let d: any = {};
      if (value != null) {
        d.deliveryBw = value * record.bwCount;
      } else {
        d.deliveryBw = null;
      }
      await setFieldsValue(d);
    }
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

  // 省份调度控制变化
  async function onIsProvinceSchedulingChange(value) {
    if (props.needAreaInfo) {
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
              options:
                record.province && areaInfo[record.province] ? areaInfo[record.province] : [],
            },
            required: false,
          },
        ]);
        // 切换到仅本省时恢复原始省份城市数据（如果有）
        if (record.province || record.city) {
          setFieldsValue({
            province: record.province || null,
            city: record.city || null,
          });
        }
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
              options:
                record.province && areaInfo[record.province] ? areaInfo[record.province] : [],
            },
            required: record.province ? true : false,
          },
        ]);
        // 切换到不限制时，如果有原始数据则恢复，否则保持为空
        if (record.province || record.city) {
          setFieldsValue({
            province: record.province || null,
            city: record.city || null,
          });
        }
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

          // 检查原始城市是否属于当前省份
          const isCityInProvince =
            record.city &&
            record.province === value &&
            areaInfo[value]?.some((city) => city.value === record.city);

          // 如果用户没有改变省份，且原始城市属于该省份，则恢复城市值
          if (data.province === value && isCityInProvince) {
            await setFieldsValue({
              city: record.city,
            });
          }
          // 如果用户没有改变省份，但原始城市不属于该省份，则清空城市
          else if (data.province === value && record.city) {
            await setFieldsValue({
              city: null,
            });
          }
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

  // 将运营商名称转换为运营商ID
  function getIspIdByName(ispName: string): string | null {
    const ispMap: Record<string, string> = {
      本网: '0',
      电信: '1',
      联通: '2',
      移动: '3',
    };
    return ispMap[ispName] || null;
  }

  // 是否异网变化处理
  async function onIsCoverDiffIspChange(value) {
    if (value === 1) {
      // 选择"是"时，异网运营商字段变为必填，并自动填充当前isp字段
      await updateSchema([
        {
          field: 'difIsp',
          required: true,
        },
      ]);

      // 自动填充当前isp字段，将运营商名称转换为ID
      if (record.isp && !getFieldsValue().difIsp) {
        const ispId = getIspIdByName(record.isp);
        if (ispId) {
          await setFieldsValue({
            difIsp: ispId,
          });
        }
      }
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

  // 暴露方法给子组件使用
  defineExpose({
    onSingleDiverybwChange,
    onIsProvinceSchedulingChange,
    onProvinceChange,
    onIsCoverDiffIspChange,
    buildAreaInfo,
    getIspIdByName,
  });
</script>

<style lang="less" scoped>
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
