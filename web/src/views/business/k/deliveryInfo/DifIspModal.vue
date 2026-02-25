<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="getTitle"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>
<script lang="ts" setup>
  import { ref, computed, unref, h } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { getDifIspForm } from './data';
  import { message, Modal } from 'ant-design-vue';
  import { GetAreaList, DeliveryDifIsp } from '@/api/business/k';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';

  defineOptions({ name: 'DeliveryInfoModal' });

  const emit = defineEmits(['success', 'register']);
  const prop = defineProps({ bizType: { type: String, required: true } });

  const isUpdate = ref(true);
  let record: any = {};
  let areaInfo: Record<string, string[]> = {};
  let provinceOptions: any[] = [];

  const [
    registerForm,
    { setFieldsValue, resetFields, validate, updateSchema, getFieldsValue, clearValidate },
  ] = useForm({
    labelWidth: 120,
    baseColProps: { span: 24 },
    schemas: getDifIspForm(onIsProvinceSchedulingChange, onProvinceChange),
    showActionButtonGroup: false,
    actionColOptions: {
      span: 24,
    },
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    await resetFields();
    provinceOptions = [];
    areaInfo = {};
    record = data.record;
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
    setModalProps({
      confirmLoading: false,
      destroyOnClose: true,
      width: 800,
      height: 100,
      title: `${record.hostname} 异网跨省下发`,
    });
    setFieldsValue({
      isProvinceScheduling: record.isProvinceScheduling,
    });
  });

  const getTitle = computed(() => (!unref(isUpdate) ? '新增信息' : '异网跨省下发'));

  async function handleSubmit() {
    try {
      const values = await validate();
      values.id = record.id;
      setModalProps({ confirmLoading: true });
      values.bizType = prop.bizType;
      const m = Modal.confirm({
        title: '是否确认异网跨省下发?',
        icon: h(ExclamationCircleOutlined),
        maskClosable: true,
        async onOk() {
          try {
            await DeliveryDifIsp(values);
            message.success('下发成功');
            closeModal();
            emit('success');
          } finally {
            m.destroy();
            //  ignore
          }
        },
      });
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }

  function buildAreaInfo(areaInfo: Array<any>): Record<string, string[]> {
    let data: Record<string, any> = {};
    for (let i = 0; i < areaInfo?.length; i++) {
      let item = areaInfo[i];
      let province = item.province;
      if (!data[province]) {
        data[province] = [];
      }
      data[province].push({ label: item.city, value: item.city });
    }
    return data;
  }

  async function onIsProvinceSchedulingChange(value) {
    // 仅本省
    if (value == 1) {
      await updateSchema([
        {
          field: 'province',
          componentProps: {
            disabled: true,
            placeholder: '请选择省份',
          },
          required: false,
        },
        {
          field: 'city',
          componentProps: {
            disabled: true,
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
            placeholder: '不填为“全国”',
          },
          required: false,
        },
        {
          field: 'city',
          componentProps: {
            disabled: false,
          },
          required: false,
        },
      ]);
    }
    setFieldsValue({
      province: null,
      city: null,
    });
    clearValidate(['city', 'province']);
  }
  async function onProvinceChange(value) {
    let data = await getFieldsValue();
    // 不限制
    if (data.isProvinceScheduling == 0) {
      await setFieldsValue({
        city: null,
      });
      if (data.province) {
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
  function handleCancel() {
    emit('success');
  }
</script>
