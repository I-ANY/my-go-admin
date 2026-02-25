<template>
  <BasicDrawer
    v-bind="$attrs"
    @register="registerDrawer"
    showFooter
    :title="getTitle"
    width="40%"
    @ok="handleSubmit()"
  >
    <BasicForm @register="registerForm" />
  </BasicDrawer>
</template>
<script lang="ts" setup>
  import { ref, computed, unref } from 'vue';
  import { BasicForm, useForm } from '@/components/Form';
  import { BasicDrawer, useDrawerInner } from '@/components/Drawer';
  import { UpdateCategory, CreateCategory } from '@/api/business/biz';
  import { drawerFormSchema, isVirtualEnum } from './data';

  const getTitle = computed(() => (!isUpdate.value ? '新增业务大类' : '编辑业务大类'));
  const emit = defineEmits(['success', 'register']);
  const isUpdate = ref(true);
  const recordId = ref(0);
  let checked = ref(false);
  const [registerForm, { resetFields, setFieldsValue, validate, updateSchema }] = useForm({
    labelWidth: 140,
    baseColProps: { span: 24 },
    schemas: drawerFormSchema(onIsVirtualChange),
    showActionButtonGroup: false,
  });

  const [registerDrawer, { setDrawerProps, closeDrawer }] = useDrawerInner(async (data) => {
    checked.value = false;
    await resetFields();
    setDrawerProps({ confirmLoading: false, destroyOnClose: false });
    isUpdate.value = !!data?.isUpdate;
    const ids = data?.record?.subcategories?.map((item: any) => item.id) || [];
    const virtualIds = data?.record?.virtualSubcategories?.map((item: any) => item.id) || [];
    if (isUpdate.value) {
      recordId.value = data.record.id;
      await setFieldsValue({
        ...data.record,
        subcategories: ids || [],
        virtualSubcategories: virtualIds || [],
      });
    } else {
      await setFieldsValue({
        isVirtual: isVirtualEnum.NO,
      });
    }
  });

  async function handleSubmit() {
    try {
      const values = await validate();

      if (values.isVirtual == isVirtualEnum.YES) {
        values.subcategories = null;
      } else {
        values.virtualSubcategories = null;
      }

      setDrawerProps({ confirmLoading: true, destroyOnClose: true });
      if (unref(isUpdate)) {
        //附加ID
        let newValues = {
          ...values,
          id: unref(recordId),
        } as any;
        await UpdateCategory(unref(recordId), newValues);
        emit('success', false);
      } else {
        await CreateCategory(unref(values));
        emit('success', true);
      }
      closeDrawer();
    } finally {
      setDrawerProps({ confirmLoading: false });
    }
  }
  function onIsVirtualChange(value: number) {
    // 虚拟子业务
    if (value === isVirtualEnum.YES) {
      updateSchema([
        {
          field: 'virtualSubcategories',
          show: true,
        },
        {
          field: 'subcategories',
          show: false,
          componentProps: {
            params: { 'categoryIds[]': unref(recordId) <= 0 ? null : unref(recordId) },
          },
        },
        {
          field: 'code',
          componentProps: {
            disabled: unref(isUpdate),
          },
        },
      ]);
    } else {
      updateSchema([
        {
          field: 'virtualSubcategories',
          show: false,
        },
        {
          field: 'subcategories',
          show: true,
          componentProps: {
            params: { 'categoryIds[]': unref(recordId) <= 0 ? null : unref(recordId) },
          },
        },
        {
          field: 'code',
          componentProps: {
            disabled: unref(isUpdate),
          },
        },
      ]);
    }
  }
</script>
