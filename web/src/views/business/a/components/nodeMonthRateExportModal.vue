<template>
  <div>
    <BasicModal v-bind="$attrs" @register="registerModal" @ok="handleOk">
      <BasicForm @register="registerForm" />
    </BasicModal>
  </div>
</template>
<script lang="ts" setup>
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, FormSchema, useForm } from '@/components/Form';
  import { Api, getOwners } from '@/api/business/a';
  import { downloadFileByUrl } from '@/utils/download';
  import { message, notification } from 'ant-design-vue';
  import dayjs from 'dayjs';

  defineEmits(['register']);

  const [registerForm, { validate, resetFields }] = useForm({
    schemas: getFormSchemas(),
    showActionButtonGroup: false,
    labelWidth: 80,
  });
  const [registerModal, { closeModal, setModalProps }] = useModalInner(async (_d) => {
    resetFields();
    setModalProps({
      confirmLoading: false,
      title: '节点利用率导出',
      width: 350,
      minHeight: 50,
      height: 180,
      canFullscreen: false,
      maskClosable: false,
    });
  });

  const handleOk = async () => {
    let values = await validate();
    values.month = dayjs(values.month).format('YYYY-MM');
    if (values?.roomNo?.length > 10) {
      message.error('最多一次性导出10个节点信息');
      return;
    }
    try {
      setModalProps({
        confirmLoading: true,
      });
      await exportData(values);
    } catch (error) {
      notification.error({
        message: '导出失败',
        description: error.message,
        duration: null,
      });
    } finally {
      setModalProps({
        confirmLoading: false,
      });
      closeModal();
    }
  };

  async function exportData(value: Recordable) {
    let filename = await downloadFileByUrl(Api.ExportNodeMonthRate, 'POST', 5 * 60, value, null);
    notification.success({
      message: '导出成功',
      description: '文件名：' + filename,
      duration: null,
    });
  }

  function getFormSchemas(): FormSchema[] {
    return [
      {
        field: 'roomNo',
        component: 'ApiSelect',
        label: '节点编号',
        componentProps: {
          options: [],
          mode: 'multiple',
          maxTagCount: 2,
          showSearch: true,
          placeholder: '请选择节点编号',
          api: async () => {
            const data = await getOwners();
            return data.map((item: any) => ({
              label: item.name,
              value: item.id,
            }));
          },
        },
        required: true,
        colProps: { span: 24 },
      },
      {
        field: 'month',
        label: '月份',
        component: 'DatePicker',
        componentProps: {
          allowClear: false,
          format: 'YYYY-MM',
          style: {
            width: '100%',
          },
          picker: 'month',
        },
        required: true,
        colProps: { span: 24 },
      },
    ];
  }
</script>
