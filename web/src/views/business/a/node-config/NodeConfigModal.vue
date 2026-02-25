<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="getTitle"
    @ok="handleSubmit"
    width="50%"
  >
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>

<script lang="ts" setup>
  import { ref, computed, unref } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { formSchema } from './data';
  import { createNodeConfig, updateNodeConfig } from '@/api/business/a';
  import { message } from 'ant-design-vue';

  defineOptions({ name: 'NodeConfigModal' });

  const emit = defineEmits(['success', 'register']);

  const isUpdate = ref(true);
  const rowId = ref('');

  const [registerForm, { setFieldsValue, clearValidate, validate, resetFields }] = useForm({
    labelWidth: 120,
    schemas: formSchema,
    showActionButtonGroup: false,
    baseColProps: { lg: 24, md: 24 },
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    resetFields();
    clearValidate();
    setModalProps({ confirmLoading: false });
    isUpdate.value = !!data?.isUpdate;

    if (unref(isUpdate)) {
      rowId.value = data.record.id;

      let formData = { ...data.record };
      // 编辑时转换保底带宽格式：bps -> Gbps
      if (
        formData.minBw !== undefined &&
        formData.minBw !== null &&
        data.record.billingType !== '买断'
      ) {
        formData.minBw = formData.minBw / 1000 / 1000 / 1000;
      }
      setFieldsValue(formData);
    } else {
      // 新增时设置默认值
      setFieldsValue({
        billingType: '日95',
      });
    }
  });

  const getTitle = computed(() => (!unref(isUpdate) ? '新增节点配置' : '编辑节点配置'));

  async function handleSubmit() {
    try {
      const values = await validate();
      setModalProps({ confirmLoading: true });

      const submitData = {
        ...values,
      };
      submitData.id = Number(rowId.value);
      submitData.minBw = submitData.minBw * 1000 * 1000 * 1000; // 转换为bps

      if (unref(isUpdate)) {
        await updateNodeConfig(Number(rowId.value), submitData);
        message.success('编辑成功');
      } else {
        await createNodeConfig(submitData);
        message.success('新增成功');
      }

      closeModal();
      emit('success');
    } catch (error) {
      console.error('操作失败:', error);
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
</script>
