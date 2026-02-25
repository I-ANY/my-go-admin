<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="getTitle"
    :destroyOnClose="true"
    @ok="handleSubmit"
  >
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>
<script setup lang="ts">
  import { computed, defineOptions, unref, ref } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { optionGroupFormSchema } from './data';
  import { AddOperation, EditOperation } from '@/api/ops/execute';
  import { message } from 'ant-design-vue';

  defineOptions({
    name: 'OptionsGroupModal',
  });
  const emit = defineEmits(['success', 'register']);

  const isUpdate = ref(false);
  let businessId = '';
  let id = 0;
  const getTitle = computed(() => (!unref(isUpdate) ? '新增操作组' : '编辑操作组'));

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    setModalProps({ width: 500, minHeight: 60 });
    isUpdate.value = data.isUpdate;
    businessId = data.businessId;
    if (data.isUpdate) {
      id = data.optionGroupId;
      await setFieldsValue({
        name: data.title,
      });
    }
  });

  const [registerForm, { getFieldsValue, setFieldsValue, validate }] = useForm({
    labelWidth: 80,
    baseColProps: { span: 24 },
    schemas: optionGroupFormSchema,
    showActionButtonGroup: false,
  });

  async function handleSubmit() {
    validate().then(async () => {
      const params = getFieldsValue();
      params.businessId = businessId;
      if (isUpdate.value) {
        try {
          await EditOperation(id, params);
          message.success('操作成功');
          emit('success');
        } finally {
          closeModal();
        }
      } else {
        try {
          await AddOperation(params);
          message.success('操作成功');
          emit('success');
        } finally {
          closeModal();
        }
      }
    });
  }
</script>

<style scoped lang="less"></style>
