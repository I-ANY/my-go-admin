<template>
  <BasicModal v-bind="$attrs" @register="registerModal" @ok="handleSubmit">
    <div class="info-grid mb-4 p-4 rounded bg-gray-50">
      <div class="mb-2">
        <label class="text-gray-600">已选中需求数量：</label>
        <span class="font-medium text-primary">{{ modalData.demandIds?.length || 0 }} 条</span>
      </div>
      <div class="mb-2" v-if="modalData.devName">
        <label class="text-gray-600">设备类型：</label>
        <span class="font-medium text-primary">{{ modalData.devName }}</span>
      </div>
      <div class="text-gray-500 text-sm">
        需求ID：{{ (modalData.demandIds || []).join(', ') }}
      </div>
    </div>
    <BasicForm v-if="modalData.bizType === 'specialLine'" @register="registerForm" />
    <div v-else class="text-center text-gray-500 py-4"> 汇聚类型将直接提交选中的需求进行占用 </div>
  </BasicModal>
</template>

<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { notification } from 'ant-design-vue';
  import { DemandOccupyTask } from '@/api/business/k';
  import { ref } from 'vue';
  import type { FormSchema } from '@/components/Form';

  const emit = defineEmits(['success', 'register']);

  const hostnameSchema: FormSchema[] = [
    {
      field: 'hostnames',
      label: '主机名列表',
      component: 'InputTextArea',
      componentProps: {
        placeholder: '请输入需要占用的主机名，多个用换行分隔',
        rows: 6,
        maxlength: 5000,
        allowClear: true,
      },
      required: true,
      colProps: { span: 24 },
    },
  ];

  const [registerForm, { validate, resetFields }] = useForm({
    labelWidth: 100,
    schemas: hostnameSchema,
    showActionButtonGroup: false,
  });

  const modalData = ref<{
    demandIds: string[];
    bizType: 'normal' | 'specialLine';
    provider: string;
    devName: string | null;
  }>({
    demandIds: [],
    bizType: 'normal',
    provider: '',
    devName: null,
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    modalData.value = data;
    const title = data.bizType === 'normal' ? '批量占用 - 汇聚' : '批量占用 - 专线';
    setModalProps({
      title,
      confirmLoading: false,
      destroyOnClose: true,
      width: 600,
    });
    // 重置表单
    if (data.bizType === 'specialLine') {
      await resetFields();
    }
  });

  async function handleSubmit() {
    const demandIds = modalData.value.demandIds;

    if (!demandIds || demandIds.length === 0) {
      notification.warning({
        message: '提示',
        description: '请至少选择一条需求',
      });
      return;
    }

    let hostnames: string[] = [];

    // 专线需要验证并获取 hostnames
    if (modalData.value.bizType === 'specialLine') {
      const values = await validate();
      // 将换行分隔的字符串转换为数组
      hostnames = (values.hostnames || '')
        .split('\n')
        .map((h: string) => h.trim())
        .filter((h: string) => h);

      if (hostnames.length === 0) {
        notification.warning({
          message: '提示',
          description: '请输入至少一个主机名',
        });
        return;
      }
    }

    setModalProps({ confirmLoading: true });

    try {
      const params: {
        biz_type: 'normal' | 'specialLine';
        provider: string;
        demand_ids: string[];
        dev_name?: string;
        hostnames?: string[];
      } = {
        biz_type: modalData.value.bizType,
        provider: modalData.value.provider as string,
        demand_ids: demandIds,
      };

      // 传递设备类型
      if (modalData.value.devName) {
        params.dev_name = modalData.value.devName;
      }

      // 专线才传 hostnames
      if (modalData.value.bizType === 'specialLine') {
        params.hostnames = hostnames;
      }

      await DemandOccupyTask(params);

      notification.success({
        message: '操作成功',
        description: '批量占用任务已提交',
        duration: 3,
      });
      closeModal();
      emit('success');
    } catch (error: any) {
      notification.error({
        message: '操作失败',
        description: error.message || '未知错误',
        duration: 5,
      });
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
</script>
