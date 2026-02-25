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
  import {
    updateRecord,
    createRecord,
    getCategory,
    getRegion,
    getZone,
    getMode,
  } from '@/api/price/price';

  const ispOption = [
    {
      label: '电信',
      value: '电信',
    },
    {
      label: '移动',
      value: '移动',
    },
    {
      label: '联通',
      value: '联通',
    },
  ];

  const diffIspOptionData = {
    电信: [
      {
        label: '移动',
        value: '移动',
      },
      {
        label: '联通',
        value: '联通',
      },
    ],
    联通: [
      {
        label: '电信',
        value: '电信',
      },
      {
        label: '移动',
        value: '移动',
      },
    ],
    移动: [
      {
        label: '电信',
        value: '电信',
      },
      {
        label: '联通',
        value: '联通',
      },
    ],
  };

  const formSchema: FormSchema[] = [
    {
      field: 'localIsp',
      label: '本网运营商',
      required: true,
      component: 'Select',
      componentProps: ({ formModel, formActionType }) => {
        return {
          options: ispOption,
          onChange: (e: any) => {
            let diffIspOptions = diffIspOptionData[e];
            formModel.bizIsp = undefined;
            const { updateSchema } = formActionType;
            updateSchema({
              field: 'bizIsp',
              componentProps: {
                options: diffIspOptions,
              },
            });
          },
        };
      },
    },
    {
      label: '业务组',
      field: 'categoryID',
      component: 'ApiSelect',
      componentProps: {
        api: getCategory,
        params: { pageSize: 50000, pageIndex: 1 },
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
      required: true,
    },
    {
      label: '大区',
      field: 'regionID',
      component: 'ApiSelect',
      componentProps: {
        api: getRegion,
        params: { pageSize: 5000, pageIndex: 1 },
        resultField: 'items',
        labelField: 'name',
        valueField: 'id',
        showSearch: true,
        placeholder: '大区与省份不能同时填写',
        filterOption: (input: string, option: any) => {
          return (
            option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
            option.identify?.toLowerCase().indexOf(input.toLowerCase()) >= 0
          );
        },
      },
      dynamicDisabled: ({ values }) => {
        return !!values.zoneID;
      },
    },
    {
      label: '省份',
      field: 'zoneID',
      component: 'ApiSelect',
      componentProps: {
        api: getZone,
        params: { pageSize: 5000, pageIndex: 1 },
        resultField: 'items',
        labelField: 'name',
        valueField: 'id',
        showSearch: true,
        placeholder: '大区与省份不能同时填写',
        filterOption: (input: string, option: any) => {
          return (
            option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
            option.identify?.toLowerCase().indexOf(input.toLowerCase()) >= 0
          );
        },
      },
      dynamicDisabled: ({ values }) => {
        return !!values.regionID;
      },
    },
    {
      field: 'bizIsp',
      label: '异网运营商',
      component: 'Select',
      componentProps: {
        options: [],
        placeholder: '异网业务需填写',
        onChange: (value) => {
          if (!value) {
            resetValidate();
          }
        },
      },
      dynamicRules: ({ values }) => {
        return values.bizIspMode ? [{ required: true, message: '异网运营商必填' }] : [];
      },
    },
    {
      field: 'bizIspMode',
      label: '跨网计费方式',
      component: 'Select',
      componentProps: {
        options: [
          { label: '本网计费', value: 1 },
          { label: '异网计费', value: 2 },
        ],
        placeholder: '异网业务需填写',
        onChange: (value) => {
          if (!value) {
            resetValidate();
          }
        },
      },
      dynamicRules: ({ values }) => {
        return values.bizIsp ? [{ required: true, message: '跨网计费方式必填' }] : [];
      },
    },
    {
      label: '计费方式',
      field: 'modeID',
      component: 'ApiSelect',
      componentProps: {
        api: getMode,
        params: { pageSize: 50000, pageIndex: 1 },
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
      required: true,
    },
    {
      field: 'price',
      label: '单价(元)',
      required: true,
      component: 'InputNumber',
    },
    {
      field: 'low',
      label: '溜缝业务',
      component: 'Select',
      componentProps: {
        options: [
          { label: '是', value: 1 },
          { label: '否', value: 0 },
        ],
      },
      defaultValue: 0,
      required: true,
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
  const [registerForm, { resetFields, clearValidate, setFieldsValue, validate }] = useForm({
    labelWidth: 90,
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

  const getTitle = computed(() => (!unref(isUpdate) ? '新增业务单价' : '编辑业务单价'));

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
        await updateRecord(newValues);
        emit('success', false);
      } else {
        await createRecord(values);
        emit('success', true);
      }
      closeDrawer();
    } finally {
      setDrawerProps({ confirmLoading: false });
    }
  }

  async function resetValidate() {
    clearValidate();
  }
</script>
