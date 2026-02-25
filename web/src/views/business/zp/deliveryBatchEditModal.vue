<template>
  <BasicModal v-bind="$attrs" @register="registerModal" :title="getTitle" @ok="handleSubmit">
    <BasicForm @register="registerForm" :actionColOptions="{ span: 12 }" :labelWidth="100" />
  </BasicModal>
</template>

<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { deliveryBatchEditSchemas } from './data';
  import { computed, ref, unref } from 'vue';
  import { DeliveryDeviceBatchEdit } from '@/api/business/zp';
  import { message } from 'ant-design-vue';

  defineOptions({ name: 'DeliveryBatchEditModal' });
  const emit = defineEmits(['success', 'register']);
  const isUpdate = ref(true);
  const getTitle = computed(() => (!unref(isUpdate) ? '新增信息' : '批量编辑'));
  let ids = ref([]);

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    setModalProps({
      confirmLoading: false,
      destroyOnClose: true,
      width: 500,
      minHeight: 100,
      // height: 200,
    });
    if (data?.ids) {
      ids.value = data.ids;
    }
  });

  const [registerForm, { validate }] = useForm({
    labelWidth: 120,
    schemas: deliveryBatchEditSchemas,
    showActionButtonGroup: false,
    actionColOptions: {
      span: 12,
    },
  });

  // 提交数据
  async function handleSubmit() {
    try {
      setModalProps({ confirmLoading: true });
      const values = await validate();
      values.single_line_bw = Number(values.single_line_bw);
      values.status = Number(values.status);
      values.ids = ids.value;
      await DeliveryDeviceBatchEdit(values);
      message.success('操作成功');
      closeModal();
      emit('success');
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
</script>

<style scoped lang="less"></style>
