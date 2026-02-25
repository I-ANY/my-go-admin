<template>
  <BasicModal v-bind="$attrs" @register="registerModal" @ok="handleSubmit" @cancel="handleCancel">
    <Alert :message="data.alertMessage" type="info" show-icon style="margin-bottom: 16px" />
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>

<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { modifyDscpSchemas, modifyType } from './data';
  import { reactive } from 'vue';
  import { Alert } from 'ant-design-vue';
  import { UpdateBizPriority, UpdateServerPriority } from '@/api/network/dscp';

  defineOptions({ name: 'NetworkModifyDscpModal' });

  const data = reactive({
    ids: [] as number[],
    type: '' as string,
    alertMessage: '' as string,
    owner: '' as string,
    business: '' as string,
  });
  const emit = defineEmits(['modifyDscpSuccess', 'register', 'success']);
  const [registerForm, { validate, resetFields, getFieldsValue, validateFields }] = useForm({
    schemas: modifyDscpSchemas(getFormHandleFn),
    showActionButtonGroup: false,
    labelWidth: 120,
    baseColProps: { span: 24 },
  });
  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (d) => {
    resetFields();
    data.ids = d.ids;
    data.type = d.type;
    let title = '';
    data.owner = d.owner;
    data.business = d.business;
    if (data.type === modifyType.SUMMARY) {
      title = '修改业务DSCP值';
      data.alertMessage = `节点：${data.owner}，已选择${data.ids.length}个业务`;
    } else if (data.type === modifyType.SERVER) {
      title = '修改设备DSCP值';
      data.alertMessage = `节点：${data.owner}，业务：${data.business}，已选择${data.ids.length}个设备`;
    }
    setModalProps({
      confirmLoading: false,
      destroyOnClose: true,
      title: title,
    });
  });

  async function handleSubmit() {
    const values = await validate();
    values.ids = data.ids;
    // Modal.confirm({
    //   title: '是否确认修改?',
    //   onOk: async () => {
    setModalProps({ confirmLoading: true });
    try {
      let res: any;
      if (data.type === modifyType.SUMMARY) {
        res = await UpdateBizPriority(values);
      } else if (data.type === modifyType.SERVER) {
        res = await UpdateServerPriority(values);
      }
      const results = res.results || [];
      closeModal();
      emit('modifyDscpSuccess', results);
    } finally {
      setModalProps({ confirmLoading: false });
    }
    //   },
    // });
  }
  function handleCancel() {
    emit('success');
  }
  function getFormHandleFn() {
    return {
      getFieldsValue,
      validateFields,
    };
  }
</script>
