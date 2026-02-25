<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="title"
    @ok="handleSubmit"
    :minHeight="60"
  >
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>
<script lang="ts" setup>
  import { ref } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';

  const emit = defineEmits(['success', 'register']);
  const title = ref('');
  const recordId = ref('');

  // 保存提交函数，由父组件传递
  let submitFunc: Function | null = null;

  const [registerForm, { validate, resetFields, setProps, setFieldsValue }] = useForm({
    labelWidth: 150,
    schemas: [],
    showActionButtonGroup: false,
    actionColOptions: {
      span: 23,
    },
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    await resetFields();
    setModalProps({ confirmLoading: false });

    title.value = data.title;
    recordId.value = data.recordId;
    submitFunc = data.submitFunc;

    // 动态设置表单 Schema
    if (data.schemas) {
      // 必须使用 setProps 动态更新 schemas
      // 注意：setProps 会合并配置，对于 schemas 需要确保能够正确替换
      // Vben Admin 的 BasicForm 动态修改 schemas 推荐用 updateSchema 或 setProps({ schemas })
      // 这里如果 schemas 变动大，可能需要重置。
      // 最稳妥是用 appendSchema 或 updateSchema，但这里是全新的。
      // 这里的 setProps({ schemas }) 应该是可行的。
      await setProps({ schemas: data.schemas });
    }

    if (data.values) {
      await setFieldsValue(data.values);
    }
  });

  async function handleSubmit() {
    try {
      const values = await validate();
      setModalProps({ confirmLoading: true });
      console.log('values', values);
      if (submitFunc) {
        // 将 ID 和表单值一起传给回调
        await submitFunc({ ids: [recordId.value], ...values });
      }

      closeModal();
      emit('success');
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
</script>
