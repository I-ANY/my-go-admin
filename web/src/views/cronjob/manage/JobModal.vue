<template>
  <BasicModal v-bind="$attrs" @register="register" @ok="handleSubmit">
    <BasicForm @register="registerForm"
  /></BasicModal>
</template>
<script lang="ts" setup>
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { getFormSchema } from './data';
  import { reactive } from 'vue';
  import { createJob, updateJob } from '@/api/cronjob/job';

  const emits = defineEmits(['success', 'register']);
  const allData = reactive({
    isUpdate: false,
    record: {} as any,
  });
  const [registerForm, { resetFields, setFieldsValue, validate, updateSchema }] = useForm({
    labelWidth: 120,
    baseColProps: { span: 24 },
    schemas: getFormSchema(onJobTypeSelect),
    showActionButtonGroup: false,
  });
  const [register, { closeModal, setModalProps }] = useModalInner(async (data) => {
    resetFields();
    allData.isUpdate = !!data.isUpdate;
    setModalProps({
      confirmLoading: false,
      destroyOnClose: true,
      title: allData.isUpdate ? '编辑任务' : '新增任务',
      width: 800,
      height: 500,
    });
    if (allData.isUpdate) {
      onJobTypeSelect(data.record.jobType);
      await updateSchema({
        field: 'scheduleNode',
        componentProps: {
          disabled: true,
        },
      });
      allData.record = data.record;
      setFieldsValue({
        ...data.record,
      });
    } else {
      await updateSchema({
        field: 'scheduleNode',
        componentProps: {
          disabled: false,
        },
      });
    }
  });
  async function handleSubmit() {
    let data = await validate();
    if (allData.isUpdate) {
      await updateJob(allData.record.id, data);
    } else {
      await createJob(data);
    }
    emits('success', !allData.isUpdate);

    closeModal();
  }

  function onJobTypeSelect(value: string | number) {
    const tipInfo = `{
    "uri":"/ping",
    "method":"GET/POST/DELETE/PUT",
    "header":{ "k1":"v1" },
    "query":{ "k2":"v2" },
    "body":{ "k3":"v3" },
    "timeout": 15
  }`;
    // 类型为flask-api
    if (value == 1) {
      updateSchema([
        {
          field: 'invokeTarget',
          required: false,
          show: false,
        },
        {
          field: 'args',
          required: true,
          componentProps: {
            placeholder: tipInfo,
          },
        },
      ]);
    } else {
      // exec
      updateSchema([
        {
          field: 'invokeTarget',
          required: true,
          show: true,
        },
        {
          field: 'args',
          required: false,
          componentProps: {
            placeholder: '请输入args',
          },
        },
      ]);
    }
  }
</script>
