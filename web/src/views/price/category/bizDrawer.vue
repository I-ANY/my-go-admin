<template>
  <BasicDrawer
    v-bind="$attrs"
    @register="registerDrawer"
    showFooter
    :title="getTitle"
    width="40%"
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
  import { getCategory, updateBiz } from '@/api/price/price';

  const formSchema: FormSchema[] = [
    {
      field: 'categoryID',
      label: '业务组',
      required: true,
      component: 'ApiSelect',
      componentProps: {
        api: getCategory,
        params: { pageSize: 5000, pageIndex: 1 },
        resultField: 'items',
        labelField: 'name',
        valueField: 'id',
        showSearch: true,
        filterOption: (input: string, option: any) => {
          return (
            option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
            option.identify?.toLowerCase().indexOf(input.toLowerCase()) >= 0
          );
        },
      },
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
        categoryID: null,
      });
      recordId.value = data.record.id;
    }
  });

  const getTitle = computed(() => (!unref(isUpdate) ? '新增业务详情' : '编辑业务详情'));

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
        await updateBiz(newValues);
        emit('success', false);
      } else {
        await updateBiz(values);
        emit('success', true);
      }
      closeDrawer();
    } finally {
      setDrawerProps({ confirmLoading: false });
    }
  }
</script>
