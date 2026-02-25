<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    title="设备录入"
    @ok="handleSubmit"
    :minHeight="40"
  >
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>
<script lang="ts" setup>
  import { ref } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { DeliveryDeviceImport } from '@/api/business/zap';
  import { notification } from 'ant-design-vue';
  import { getSelectOptionsFromDict } from '@/utils/dict';
  import { ZEnum } from '@/enums/dictTypeCode';

  const emit = defineEmits(['success', 'register']);
  const ids = ref<string[]>([]);

  const [registerForm, { validate, resetFields }] = useForm({
    labelWidth: 100,
    schemas: [
      {
        field: 'account_type',
        label: '账号类型',
        component: 'Select',
        required: true,
        defaultValue: 'mful',
        colProps: {
          span: 20,
        },
        componentProps: {
          options: getSelectOptionsFromDict(ZEnum.Z_ACCOUNT_TYPE),
        },
      },
    ],
    showActionButtonGroup: false,
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    await resetFields();
    setModalProps({ confirmLoading: false });
    ids.value = data.ids;
  });

  async function handleSubmit() {
    try {
      const values = await validate();
      setModalProps({ confirmLoading: true });

      const params = {
        ids: ids.value,
        account_type: values.account_type,
      };

      const result = await DeliveryDeviceImport(params);

      if (result) {
        result.forEach((item: any) => {
          if (item.success) {
            notification.success({ message: item.hostname, description: item.msg });
          } else {
            notification.error({ message: item.hostname, description: item.msg });
          }
        });
      }

      closeModal();
      emit('success');
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
</script>
