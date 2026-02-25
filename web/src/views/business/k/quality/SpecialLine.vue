<template>
  <QualityAbnormal :bizType="BizType.specialLine" />
</template>

<script setup lang="ts">
  import QualityAbnormal from './QualityAbnormal.vue';
  import { BizType } from '@/views/business/k/data';
</script>
