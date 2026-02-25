<template>
  <BasicModal v-bind="$attrs" @register="registerModal" @ok="handleSubmit">
    <div class="info-grid mb-4 p-4 rounded">
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="text-gray-600">缺口带宽：</label>
          <span class="font-medium">{{ modalData?.record?.gap_bw }} Mbps</span>
        </div>
        <div>
          <label class="text-gray-600">区域：</label>
          <span class="font-medium">{{ modalData?.record?.area_name }}</span>
        </div>
        <div>
          <label class="text-gray-600">省份：</label>
          <span class="font-medium">{{ modalData?.record?.province_name }}</span>
        </div>
        <div>
          <label class="text-gray-600">运营商：</label>
          <span class="font-medium">{{ modalData?.record?.isp_name }}</span>
        </div>
      </div>
    </div>
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>

<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { notification } from 'ant-design-vue';
  import { lockDemand } from '@/api/business/k';
  import { ref } from 'vue';

  const [registerForm, { validate, setFieldsValue }] = useForm({
    labelWidth: 80,
    schemas: [
      {
        field: 'locked_bw',
        label: '锁定带宽',
        component: 'InputNumber',
        componentProps: {
          min: 1,
          style: { width: '100%' },
          placeholder: '请输入锁定带宽, 默认全部锁定',
        },
        colProps: { span: 22 },
      },
    ],
    showActionButtonGroup: false,
  });

  const modalData = ref<any>({});

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    setModalProps({
      title: `锁定${data.record.demand_id}需求`,
      width: 500,
      destroyOnClose: true,
    });
    modalData.value = data;
    setFieldsValue(data);
  });

  async function handleSubmit() {
    try {
      const values = await validate();
      setModalProps({ confirmLoading: true });

      const response = await lockDemand({
        demand_id: modalData.value.record.demand_id,
        is_locked: true,
        locked_bw: values.locked_bw,
      });

      const { code, msg } = response.data;

      if (code === 200) {
        notification.success({
          message: '操作成功',
          description: msg,
          duration: 5,
          placement: 'top',
        });
        closeModal();
        emit('success'); // 触发父组件刷新
      } else {
        notification.error({
          message: '操作失败',
          description: msg,
          duration: 5,
          placement: 'top',
        });
        closeModal();
      }
      emit('success');
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }

  const emit = defineEmits(['success']);
</script>
