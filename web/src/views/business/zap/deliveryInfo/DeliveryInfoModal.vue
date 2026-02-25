<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="getTitle"
    @ok="handleSubmit"
    :style="{ top: '0px' }"
  >
    <BasicForm @register="registerForm" :actionColOptions="{ span: 24 }" :labelWidth="200" />
  </BasicModal>
</template>

<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import {
    deliveryInfoSchemas,
    deviceLoginList,
    submitTypeMap,
    deviceLoginTypeMap,
    accountTypeMap,
  } from './data';
  import { computed, onMounted, ref, unref } from 'vue';
  import { GetZapAreaList, UpdateDeliveryDevice } from '@/api/business/zap';
  import { message } from 'ant-design-vue';

  defineOptions({ name: 'DeliveryInfoModal' });
  const emit = defineEmits(['success', 'register']);
  const isUpdate = ref(true);
  const getTitle = computed(() => (!unref(isUpdate) ? '新增信息' : '编辑信息'));

  type Option = { label: string; value: string; children?: any[] };
  const provincesOptions = ref<Option[]>([]);
  onMounted(async () => {
    provincesOptions.value = await GetZapAreaList();
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    setModalProps({
      confirmLoading: false,
      destroyOnClose: true,
      width: 900,
      height: 700,
    });
    if (data?.record) {
      // 创建表单数据的副本（避免直接修改原始数据）
      const formData = { ...data.record };
      // 对部分字段进行转义处理
      // 处理提交类型字段
      if (formData.submit_type !== undefined) {
        formData.submit_type = submitTypeMap[formData.submit_type].dictLabel;
      }
      // 处理登录方式字段
      if (formData.ip_domain_option !== undefined) {
        formData.ip_domain_option = deviceLoginTypeMap[formData.ip_domain_option].dictLabel;
      }
      // 处理交付状态
      if (formData.order_status !== undefined) {
        formData.order_status = String(formData.order_status);
      }
      // 处理账号类型字段
      if (formData.account_type !== undefined && accountTypeMap[formData.account_type]) {
        formData.account_type = accountTypeMap[formData.account_type].dictLabel;
      }
      // 设置表单值
      await setFieldsValue(formData);
    }
  });

  const [registerForm, { setFieldsValue, validate }] = useForm({
    labelWidth: 120,
    schemas: deliveryInfoSchemas,
    showActionButtonGroup: false,
    actionColOptions: {
      span: 12,
    },
  });

  // 提交数据
  async function handleSubmit() {
    try {
      const values = await validate();
      setModalProps({ confirmLoading: true });
      values.ip_domain_option = Number(
        deviceLoginList.find((item) => item.label === values.ip_domain_option)?.value || 1,
      );
      values.id = Number(values.id);
      values.single_line_bandwidth = Number(values.single_line_bandwidth);
      values.total_bandwidth = Number(values.total_bandwidth);
      values.domain_login_port = Number(values.domain_login_port);
      values.submit_type = Number(
        deviceLoginList.find((item) => item.label === values.submit_type)?.value || 1,
      );
      values.order_status = Number(values.order_status);
      // 获取到对应的省份城市信息的中文信息
      provincesOptions.value.find((item) => {
        if (item.value === values.province_code) {
          values.province = item.label;
          values.city = item.children?.find((i) => i.value === values.city_code)?.label;
        }
      });
      await UpdateDeliveryDevice(values);
      message.success('编辑成功');
      closeModal();
      emit('success');
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
</script>

<style scoped lang="less"></style>
