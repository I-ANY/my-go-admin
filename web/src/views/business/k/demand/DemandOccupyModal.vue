<template>
  <BasicModal v-bind="$attrs" @register="registerModal" @ok="handleSubmit">
    <div class="info-grid mb-4 p-4 rounded">
      <div class="grid grid-cols-2 gap-4" style="padding-left: 30px">
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
        <div class="col-span-2">
          <span>默认根据缺口带宽按10G全部占用，可调整提交带宽和MAC地址数量</span>
        </div>
      </div>
    </div>
    <BasicForm @register="registerForm" />
  </BasicModal>
</template>

<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicForm, useForm } from '@/components/Form';
  import { DemandOccupySchema } from './data';
  import { notification } from 'ant-design-vue';
  import { DemandOccupy } from '@/api/business/k';
  import { ref, defineEmits } from 'vue';

  const modalData = ref<any>({});

  const [registerForm, { validate, setFieldsValue, resetSchema }] = useForm({
    labelWidth: 120,
    schemas: DemandOccupySchema(0, ''),
    showActionButtonGroup: false,
  });

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    setModalProps({
      title: `${data.record.demand_id}需求占用`,
      confirmLoading: false,
      destroyOnClose: true,
      width: 800,
    });
    modalData.value = data;
    // 根据 is_cover_diff_isp 动态更新表单 schema
    await resetSchema(
      DemandOccupySchema(data.record?.is_cover_diff_isp ?? 0, data.record?.isp_name ?? ''),
    );
    setFieldsValue(data);
  });

  async function handleSubmit() {
    const values = await validate();

    // 验证 Mac 数量：gap_bw > submit_bw * mac_count 时不允许提交
    const gapBw = modalData.value.record?.gap_bw || 0;
    const submitBw = values.submit_bw || 0;
    const macCount = values.mac_count || 0;

    if (submitBw > 0 && macCount > 0 && gapBw < submitBw * macCount) {
      notification.warning({
        message: '验证失败',
        description: `缺口带宽(${gapBw}) 小于 提交带宽(${submitBw}) × MAC数量(${macCount}) = ${submitBw * macCount}，请调整参数`,
      });
      throw new Error('验证失败');
    }

    // 合并业务参数
    Object.assign(values, {
      biz_type: modalData.value.record.biz_type,
      provider: modalData.value.record.provider,
      demand_id: modalData.value.record.demand_id,
    });

    setModalProps({ confirmLoading: true });

    // 调用接口并解构响应
    const response = await DemandOccupy(values);
    const { code, msg } = response.data;
    console.log(code, msg);

    if (code === 200) {
      notification.success({
        message: '操作成功',
        description: msg,
        duration: 5,
        placement: 'top',
      });
      closeModal();
      defineEmits(['success']); // 触发父组件刷新
    } else {
      notification.error({
        message: '操作失败',
        description: msg,
        duration: 5,
        placement: 'top',
      });
      closeModal();
    }
  }
</script>
