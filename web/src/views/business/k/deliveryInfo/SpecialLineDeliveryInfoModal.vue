<template>
  <BaseDeliveryInfoModal
    v-bind="$attrs"
    @register="registerModal"
    @success="handleSuccess"
    :bizType="bizType"
    :formSchema="formSchema"
    :showDifIspButton="true"
    :needAreaInfo="true"
    ref="baseModalRef"
  />
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  import { useModalInner } from '@/components/Modal';
  import { speciaLineInfoFormSchema } from './data';
  import BaseDeliveryInfoModal from './BaseDeliveryInfoModal.vue';

  defineOptions({ name: 'SpecialLineDeliveryInfoModal' });

  const emit = defineEmits(['success', 'register']);
  defineProps({ bizType: { type: String, required: true } });

  const baseModalRef = ref();

  // 直接使用静态schema，让基础组件动态绑定事件
  const formSchema = speciaLineInfoFormSchema();

  const [registerModal] = useModalInner();

  function handleSuccess() {
    emit('success');
  }
</script>
