<template>
  <BasicDrawer
    v-bind="$attrs"
    @register="registerDrawer"
    showFooter
    :title="getTitle"
    width="500px"
    @ok="handleSubmit"
  >
    <BasicForm @register="registerForm" />
  </BasicDrawer>
</template>
<script lang="ts" setup>
  import { ref, computed, unref } from 'vue';
  import { FormSchema } from '@/components/Table';
  import { BasicForm, useForm } from '@/components/Form';
  import { BasicDrawer, useDrawerInner } from '@/components/Drawer';
  import { updateMode, createMode } from '@/api/price/price';

  const formSchema: FormSchema[] = [
    {
      field: 'name',
      label: '计费方式',
      required: true,
      component: 'Input',
    },
    {
      label: '备注',
      field: 'describe',
      component: 'InputTextArea',
    },
  ];

  const emit = defineEmits(['success', 'register']);
  const isUpdate = ref(true);
  const recordId = ref(0);
  let checked = ref(false);
  const [registerForm, { resetFields, setFieldsValue, validate }] = useForm({
    labelWidth: 90,
    baseColProps: { span: 24 },
    schemas: formSchema,
    showActionButtonGroup: false,
  });

  const [registerDrawer, { setDrawerProps, closeDrawer }] = useDrawerInner(async (data) => {
    checked.value = false;
    resetFields();
    setDrawerProps({ confirmLoading: false, destroyOnClose: false });
    isUpdate.value = !!data?.isUpdate;
    if (unref(isUpdate)) {
      setFieldsValue({
        ...data.record,
      });
      recordId.value = data.record.id;
    }
  });

  const getTitle = computed(() => (!unref(isUpdate) ? '新增计费方式' : '编辑计费方式'));

  async function handleSubmit() {
    try {
      const values = await validate();
      if (values.menuIds?.checked) {
        values.menuIds = values.menuIds.checked;
      }
      setDrawerProps({ confirmLoading: true, destroyOnClose: true });
      if (unref(isUpdate)) {
        //附加ID
        var newValues = {
          ...values,
          id: unref(recordId),
        };
        await updateMode(newValues);
        emit('success', false);
      } else {
        await createMode(values);
        emit('success', true);
      }
      closeDrawer();
    } finally {
      setDrawerProps({ confirmLoading: false });
    }
  }
</script>
