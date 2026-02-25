<template>
  <BasicModal v-bind="$attrs" @register="registerModal" :title="getTitle" @ok="handleSubmit">
    <BasicForm @register="registerForm" :actionColOptions="{ span: 24 }" />
  </BasicModal>
</template>
<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { computed } from 'vue';
  import { message } from 'ant-design-vue';
  import { EditDeviceMac } from '@/api/business/k';

  defineOptions({ name: 'MacEditModal' });

  let currentRecord: Recordable | undefined;

  const getTitle = computed(() => '编辑备注');

  const [registerForm, { validate, setFieldsValue }] = useForm({
    labelWidth: 50,
    schemas: [
      {
        label: '备注',
        field: 'remark',
        component: 'InputTextArea',
        componentProps: {
          placeholder: '请输入备注',
          maxlength: 500,
          showCount: true,
        },
        colProps: { span: 24 },
        required: false,
      },
    ],
    showActionButtonGroup: false,
    actionColOptions: { span: 12 },
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    currentRecord = data.record;
    setModalProps({
      confirmLoading: false,
      destroyOnClose: true,
      width: 400,
      minHeight: 100,
    });
    await setFieldsValue({
      remark: data.record?.remark || '',
    });
  });

  const emit = defineEmits(['success']);

  async function handleSubmit() {
    try {
      const values = await validate();
      setModalProps({ confirmLoading: true });
      await EditDeviceMac({
        mac_addr: currentRecord?.mac_addr,
        remark: values.remark,
      });
      message.success('备注修改成功');
      emit('success');
      closeModal();
    } catch (e) {
      setModalProps({ confirmLoading: false });
      message.error('备注修改失败');
    }
  }
</script>
