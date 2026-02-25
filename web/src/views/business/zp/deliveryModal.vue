<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :style="{ top: '0px' }"
    :title="getTitle"
    @ok="handleSubmit"
  >
    <BasicForm @register="registerForm" :actionColOptions="{ span: 24 }" :labelWidth="200" />
  </BasicModal>
</template>

<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { deliveryInfoSchemas, zpDeliveryStatusMap } from './data';
  import { computed, ref, unref } from 'vue';
  import { DeliveryDeviceEdit } from '@/api/business/zp';
  import { message } from 'ant-design-vue';

  defineOptions({ name: 'DeliveryInfoModal' });
  const emit = defineEmits(['success', 'register']);
  const isUpdate = ref(true);
  const getTitle = computed(() => (!unref(isUpdate) ? '新增信息' : '编辑信息'));

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
      if (formData.network && typeof formData.network === 'string') {
        try {
          let networkInfo = JSON.parse(formData.network);
          formData.lines = networkInfo.length;
          formData.single_line_bw = networkInfo[0].bandwidth ? networkInfo[0].bandwidth : 0;

          formData.network = parseFields(formData.network);
        } catch (e) {
          console.error('JSON解析失败:', e);
        }
      }
      // 处理 network 字段，将数组转换为换行显示的字符串
      formData.storage = parseFields(formData.storage);

      // 处理资源类型字段
      if (formData.status !== undefined) {
        formData.status = zpDeliveryStatusMap[formData.status].dictLabel || formData.status;
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
      values.id = +values.id;
      values.status = +values.status;
      values.single_line_bw = Number(values.single_line_bw);
      await DeliveryDeviceEdit(values);
      message.success('编辑成功');
      closeModal();
      emit('success');
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }

  function parseFields(str: string) {
    if (str) {
      const arrayData = JSON.parse(str);
      if (Array.isArray(arrayData)) {
        return JSON.stringify(arrayData, null, 2);
      }
    }
    return str;
  }
</script>

<style scoped lang="less"></style>
