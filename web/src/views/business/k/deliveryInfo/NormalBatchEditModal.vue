<template>
  <BaseBatchEditModal
    v-bind="$attrs"
    @register="registerModal"
    @success="handleSuccess"
    :bizType="bizType"
    :formSchema="formSchema"
    :showDifIspButton="true"
    :needAreaInfo="true"
    :selectedDevices="selectedDevices"
    ref="baseModalRef"
  />
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { useModalInner } from '@/components/Modal';
  import { normalInfoFormSchema } from './data';
  import BaseBatchEditModal from './BaseBatchEditModal.vue';

  // 定义设备信息接口
  interface DeviceInfo {
    id: number;
    deliveryBw?: number;
    bwCount?: number;
    singleDiverybw?: number;
    province?: string;
    city?: string;
    isProvinceScheduling?: number;
    isCoverDiffIsp?: number;
    remark?: string;
    [key: string]: any;
  }

  defineOptions({ name: 'NormalBatchEditModal' });

  const emit = defineEmits(['success', 'register']);
  defineProps({
    bizType: { type: String, required: true },
  });

  const baseModalRef = ref();
  const selectedDevices = ref<DeviceInfo[]>([]);

  // 直接使用静态schema，让基础组件动态绑定事件
  const formSchema = normalInfoFormSchema();

  const [registerModal] = useModalInner((data) => {
    // 接收父组件传递的selectedDevices数据
    if (data && data.selectedDevices) {
      selectedDevices.value = data.selectedDevices;
    }
  });

  function handleSuccess() {
    emit('success');
  }
</script>
