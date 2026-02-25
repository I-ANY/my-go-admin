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
  import { updateWebhook, createWebhook, getUsers } from '@/api/price/price';

  const urlTypeOptions = [
    { label: '企业微信机器人', value: 1 },
    { label: '钉钉机器人', value: 2 },
  ];

  const notifyOptions = [
    { label: '未关联业务组', value: 1 },
    { label: '业务组更新', value: 2 },
    { label: '采购单价未配置', value: 3 },
  ];

  const formSchema: FormSchema[] = [
    {
      field: 'urlType',
      label: '类型',
      required: true,
      component: 'Select',
      componentProps: {
        options: urlTypeOptions,
      },
    },
    {
      label: '通知用户',
      field: 'notifyUser',
      required: true,
      component: 'ApiSelect',
      componentProps: {
        api: getUsers,
        params: { pageSize: 50000, pageIndex: 1 },
        resultField: 'items',
        labelField: 'nickName',
        valueField: 'id',
        showSearch: true,
        mode: 'multiple',
        filterOption: (input: string, option: any) => {
          return (
            option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
            option.identify?.toLowerCase().indexOf(input.toLowerCase()) >= 0
          );
        },
      },
    },
    {
      field: 'url',
      label: 'webhook',
      component: 'Input',
      required: true,
      rules: [
        {
          required: true,
          message: '请输入 URL 地址',
        },
        {
          pattern: /^(https?:\/\/)?([\w-]+(\.[\w-]+)+)([\w.,@?^=%&:/~+#-]*[\w@?^=%&/~+#-])?$/,
          message: '请输入有效的 URL 地址',
        },
      ],
    },
    {
      label: '通知事件',
      field: 'notifyType',
      component: 'Select',
      required: true,
      componentProps: {
        mode: 'multiple',
        options: notifyOptions,
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
      });
      recordId.value = data.record.id;
    }
  });

  const getTitle = computed(() => (!unref(isUpdate) ? '新增webhook' : '编辑webhook'));

  async function handleSubmit() {
    try {
      const values = await validate();
      if (values.menuIds?.checked) {
        values.menuIds = values.menuIds.checked;
      }
      setDrawerProps({ confirmLoading: true, destroyOnClose: true });
      var newValues = {
        ...values,
      };
      if (unref(isUpdate)) {
        //附加ID
        newValues['id'] = unref(recordId);
        await updateWebhook(newValues);
        emit('success', false);
      } else {
        await createWebhook(newValues);
        emit('success', true);
      }
      closeDrawer();
    } finally {
      setDrawerProps({ confirmLoading: false });
    }
  }
</script>
