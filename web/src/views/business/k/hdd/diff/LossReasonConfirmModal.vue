<template>
  <BasicModal
    v-bind="$attrs"
    @register="register"
    title="确认流失原因"
    :width="600"
    @ok="handleSubmit"
    destroyOnClose
  >
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>

<script lang="ts" setup>
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form/index';
  import { getLossReasonFormSchema } from './data';
  import { message } from 'ant-design-vue';
  import { ConfirmHddDiff } from '@/api/business/k';

  defineOptions({
    name: 'LossReasonConfirmModal',
  });
  const emit = defineEmits(['success', 'register']);
  let ids: number[] = [];

  const [registerForm, { validate, updateSchema, getFieldsValue, clearValidate, setFieldsValue }] =
    useForm({
      labelWidth: 80,
      schemas: getLossReasonFormSchema(onLossReasonChange),
      actionColOptions: {
        span: 24,
      },
      showActionButtonGroup: false,
    });

  const [register, { setModalProps, closeModal }] = useModalInner((data) => {
    const { isBatch, ids: idsData, record } = data;
    ids = idsData;
    if (isBatch) {
      setModalProps({
        title: '批量确认流失原因',
      });
    } else {
      setModalProps({
        title: `确认流失原因：${record.hostname} - ${record.partitionName}`,
      });
    }
    if (record) {
      setFieldsValue({
        lossReason: record.lossReason == null ? null : record.lossReason + '',
        remark: record.remark,
      });
    }
  });
  async function handleSubmit() {
    try {
      setModalProps({
        loading: true,
      });
      const values = await validate();
      console.log(values);
      await ConfirmHddDiff({
        ids,
        lossReason: values.lossReason,
        remark: values.remark,
      });
      message.success('操作成功');
      emit('success');
      closeModal();
    } finally {
      setModalProps({
        loading: false,
      });
    }
  }
  function onLossReasonChange() {
    const values = getFieldsValue();
    if (values && values.lossReason === '7') {
      updateSchema({
        field: 'remark',
        required: true,
      });
    } else {
      updateSchema({
        field: 'remark',
        required: false,
      });
      clearValidate();
    }
  }
</script>
