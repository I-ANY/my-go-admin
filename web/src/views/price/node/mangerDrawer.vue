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
  import { updateNode } from '@/api/price/price';

  const formSchema: FormSchema[] = [
    {
      field: 'price',
      label: '采购单价(元/Gbps)',
      required: true,
      component: 'InputNumber',
      componentProps: {
        precision: 2,
      },
      colProps: { span: 24 },
    },
    {
      field: 'priceType',
      label: '计费方式',
      component: 'Select',
      required: true,
      componentProps: {
        placeholder: '请选择计费方式',
        options: [
          { label: '未知', value: 0 },
          { label: '日95(集群日95)', value: 1 },
          { label: '单机日95', value: 2 },
          { label: '买断', value: 3 },
          { label: '月95', value: 4 },
          { label: '单口月95', value: 5 },
        ],
      },
      colProps: { span: 24 },
    },
  ];

  const emit = defineEmits(['success', 'register']);
  const isUpdate = ref(true);
  const recordId = ref(0);
  let checked = ref(false);
  const [registerForm, { resetFields, clearValidate, setFieldsValue, validate }] = useForm({
    labelWidth: 180,
    // baseColProps: { span: 24 },
    baseColProps: { lg: 12, md: 24 },
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

  const getTitle = computed(() => (!unref(isUpdate) ? '新增节点单价' : '编辑节点单价'));

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
        await updateNode(newValues);
        emit('success', false);
      }
      closeDrawer();
    } finally {
      setDrawerProps({ confirmLoading: false });
      resetValidate();
    }
  }

  async function resetValidate() {
    clearValidate();
  }
</script>
