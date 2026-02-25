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
  import { updateCategory, createCategory, getNotRelBiz } from '@/api/price/price';

  const formSchema: FormSchema[] = [
    {
      field: 'name',
      label: '业务组',
      required: true,
      component: 'Input',
    },
    {
      field: 'outName',
      label: '对外业务名',
      component: 'Input',
    },
    {
      label: '备注',
      field: 'describe',
      component: 'InputTextArea',
    },
    {
      field: 'bizsId',
      label: '业务详情',
      required: true,
      component: 'ApiTransfer',
      colProps: {
        span: 24,
        lg: 24,
        md: 24,
      },
      componentProps: {
        api: getNotRelBiz,
        listStyle: {
          width: '50%',
          height: '300px',
        },
        titles: ['未选择', '已选择'],
        showSearch: true,
        labelField: 'name',
        valueField: 'id',
        resultField: 'items',
        immediate: false,
        filterOption: (input: string, option: any) => {
          return option.title.toLowerCase().indexOf(input.toLowerCase()) >= 0;
        },
      },
    },
  ];

  const emit = defineEmits(['success', 'register']);
  const isUpdate = ref(true);
  const recordId = ref(0);
  let checked = ref(false);
  const [registerForm, { resetFields, setFieldsValue, validate, updateSchema }] = useForm({
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

    // 临时禁用搜索框（用于清空搜索状态）
    updateSchema({
      field: 'bizsId',
      componentProps: {
        immediate: true,
        params: { id: data?.record?.id || 0 },
        showSearch: false, // 先禁用搜索框
      },
    });

    setTimeout(() => {
      // 重新启用搜索框
      updateSchema({
        field: 'bizsId',
        componentProps: {
          immediate: true,
          params: { id: data?.record?.id || 0 },
          showSearch: true, // 重新启用搜索框，确保搜索可用
        },
      });
    }, 100);

    if (unref(isUpdate)) {
      setFieldsValue({
        ...data.record,
      });
      recordId.value = data.record.id;
    }
  });

  const getTitle = computed(() => (!unref(isUpdate) ? '新增业务大类' : '编辑业务大类'));

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
        await updateCategory(newValues);
        emit('success', false);
      } else {
        await createCategory(values);
        emit('success', true);
      }
      closeDrawer();
    } finally {
      setDrawerProps({ confirmLoading: false });
    }
  }
</script>
