<template>
  <DeviceMacsBase :bizType="BizType.normal" />
</template>

<script setup lang="ts">
  import DeviceMacsBase from './DeviceMacsBase.vue';
  import { BizType } from '@/views/business/k/data';
</script>
