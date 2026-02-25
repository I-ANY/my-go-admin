<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="getTitle"
    @ok="handleSubmit"
    width="600px"
  >
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>

<script setup lang="ts">
  import { ref, computed, unref } from 'vue';
  import { BasicForm, useForm } from '@/components/Form';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { CreateBusinessGroup, UpdateBusinessGroup } from '@/api/business/biz';
  import { groupModalFormSchema } from './data';
  import { message } from 'ant-design-vue';

  const emit = defineEmits(['success', 'register']);
  const isUpdate = ref(true);
  const recordId = ref(0);

  const getTitle = computed(() => (!unref(isUpdate) ? '新增业务组' : '编辑业务组'));

  const [registerForm, { resetFields, setFieldsValue, validate }] = useForm({
    labelWidth: 100,
    baseColProps: { span: 24 },
    schemas: groupModalFormSchema,
    showActionButtonGroup: false,
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    await resetFields();
    setModalProps({ confirmLoading: false });
    isUpdate.value = !!data?.isUpdate;

    if (unref(isUpdate)) {
      await setFieldsValue({
        name: data.record.name,
        categoryIds: data.record.categoryIds || [],
      });
      recordId.value = data.record.id;
    }
  });

  async function handleSubmit() {
    try {
      const values = await validate();
      setModalProps({ confirmLoading: true });

      if (unref(isUpdate)) {
        await UpdateBusinessGroup(unref(recordId), values);
        message.success('修改成功');
      } else {
        await CreateBusinessGroup(values);
        message.success('新增成功');
      }

      closeModal();
      emit('success');
    } catch (error) {
      message.error('操作失败');
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
</script>
