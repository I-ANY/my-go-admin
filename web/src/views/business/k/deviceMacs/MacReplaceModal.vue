<template>
  <BasicModal v-bind="$attrs" @register="registerModal" :title="getTitle" @ok="handleSubmit">
    <BasicForm @register="registerForm" :actionColOptions="{ span: 24 }" />
  </BasicModal>
</template>
<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { macInfoSchemas } from './data';
  import { computed, ref, unref } from 'vue';
  import { message } from 'ant-design-vue';
  import { UpdateDeviceMac } from '@/api/business/k';
  import { useUserStore } from '@/store/modules/user';

  defineOptions({ name: 'MacReplaceModal' });

  const isUpdate = ref(true);
  const macAddr = ref('');
  const getTitle = computed(() => (!unref(isUpdate) ? '新增信息' : 'MAC替换'));

  const [registerForm, { validate, setFieldsValue }] = useForm({
    labelWidth: 100,
    schemas: macInfoSchemas,
    showActionButtonGroup: false,
    actionColOptions: { span: 12 },
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    setModalProps({
      confirmLoading: false,
      destroyOnClose: true,
      width: 800,
      minHeight: 100,
    });

    // // 获取当前调度控制状态,响应式
    // const isProvinceScheduling = ref(data.record.is_province_scheduling);

    macAddr.value = data.record.mac_addr;
    const originalArea = `${data.record.area_name}/${data.record.province_name}/${data.record.city_name}`;
    const originalIsp = data.record.isp_name;
    await setFieldsValue({
      mac_addr: data.record.mac_addr,
      originalArea: originalArea,
      originalIsp: originalIsp,
    });
  });

  async function handleSubmit() {
    try {
      const values = await validate();
      const submitData = { ...values };

      setModalProps({ confirmLoading: true });
      await UpdateDeviceMac({
        ...submitData,
        operator: useUserStore().getUserInfo?.nickname,
        phone: useUserStore().getUserInfo?.tel,
        user_id: useUserStore().getUserInfo?.id,
      });
      message.success('操作成功');
      closeModal();
    } catch (e) {
      setModalProps({ confirmLoading: false });
      message.error('操作失败');
    }
  }
</script>
