<template>
  <BasicModal v-bind="$attrs" @register="registerModal" :title="getTitle" @ok="handleSubmit">
    <BasicForm @register="registerForm" name="userModal" />
  </BasicModal>
</template>
<script lang="ts" setup>
  import { ref, computed, unref } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { accountFormSchema } from './account.data';
  import { getDeptList, addUser, updateUser } from '@/api/demo/system';
  import { message } from 'ant-design-vue';

  defineOptions({ name: 'AccountModal' });

  const emit = defineEmits(['success', 'register']);

  const isUpdate = ref(true);
  const rowId = ref('');

  const [registerForm, { setFieldsValue, updateSchema, resetFields, validate }] = useForm({
    labelWidth: 100,
    baseColProps: { span: 24 },
    schemas: accountFormSchema,
    showActionButtonGroup: false,
    actionColOptions: {
      span: 23,
    },
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    await resetFields();
    setModalProps({ confirmLoading: false, destroyOnClose: true });
    isUpdate.value = !!data?.isUpdate;

    if (unref(isUpdate)) {
      rowId.value = data.record.id;
      await updateSchema([
        {
          field: 'password',
          required: false,
          componentProps: {
            placeholder: '如果不修改用户密码则不填写',
          },
        },
      ]);
      await setFieldsValue({
        ...data.record,
        password: undefined,
      });
    } else {
      setFieldsValue({
        password: '123456',
      });
      updateSchema([
        {
          field: 'password',
          required: true,
          componentProps: {
            placeholder: '请输入用户初始密码',
          },
        },
      ]);
    }

    let res = await getDeptList();
    updateSchema([
      {
        field: 'deptId',
        componentProps: { treeData: res.items },
      },
    ]);
  });

  const getTitle = computed(() => (!unref(isUpdate) ? '新增用户' : '编辑用户'));

  async function handleSubmit() {
    try {
      const values = await validate();
      setModalProps({ confirmLoading: true });
      if (!unref(isUpdate)) {
        await addUser(values);
        message.success('添加成功');
      } else {
        await updateUser(rowId.value, values);
        message.success('编辑成功');
      }
      closeModal();
      emit('success');
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
</script>
