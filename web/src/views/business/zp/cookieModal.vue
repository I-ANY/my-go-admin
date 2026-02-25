<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="getTitle"
    :width="1000"
    @ok="handleSubmit"
    :okText="'提交'"
  >
    <BasicForm @register="registerForm" :model="modelRef" />
  </BasicModal>
</template>

<script setup lang="ts">
  import { BasicForm, FormSchema, useForm } from '@/components/Form';
  import { ref } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { formatToDateTime } from '@/utils/dateUtil';
  import { message } from 'ant-design-vue';
  import { UpdateZPCookie } from '@/api/business/zp';

  defineOptions({ name: 'CookieModal' });
  const getTitle = 'Cookie信息';
  // const emit = defineEmits(['success', 'register']);
  // const prop = defineProps();

  const schemas: FormSchema[] = [
    {
      field: 'updated_at',
      label: '更新时间',
      component: 'Input',
      colProps: {
        span: 12,
      },
      componentProps: {
        readonly: true,
      },
    },
    {
      field: 'cookie',
      label: 'cookie',
      component: 'InputTextArea',
      colProps: {
        span: 24,
      },
      componentProps: {
        readonly: true,
        rows: 8,
      },
    },
    {
      field: 'newCookie',
      component: 'InputTextArea',
      label: '新cookie',
      colProps: {
        span: 24,
      },
      required: true,
      componentProps: {
        placeholder: '可输入新cookie进行替换',
        rows: 8,
      },
    },
  ];

  const [registerForm, { validateFields }] = useForm({
    labelWidth: 80,
    schemas,
    showActionButtonGroup: false,
    actionColOptions: {
      span: 24,
    },
  });

  const [registerModal, { closeModal }] = useModalInner((data) => {
    data && onDataReceive(data);
  });
  const modelRef = ref({});

  function onDataReceive(data) {
    // 赋值到表单
    modelRef.value = {
      updated_at: formatToDateTime(data.record.updated_at),
      cookie: data.record.cookie,
      newCookie: '',
    };
  }

  async function handleSubmit() {
    const formData = await validateFields();
    if (formData.newCookie) {
      if (formData.newCookie.trim() === formData.cookie.trim()) {
        message.warning('Cookie信息一致，无需更新！！！');
        closeModal();
      } else {
        await UpdateZPCookie({
          cookie: formData.newCookie.trim(),
        });
        closeModal();
        message.success('更新成功!!!');
      }
    } else {
      closeModal();
    }
    // if modelRef.value.new
  }
</script>

<style scoped lang="less"></style>
