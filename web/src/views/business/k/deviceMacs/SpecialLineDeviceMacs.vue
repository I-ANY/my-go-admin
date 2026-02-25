<template>
  <DeviceMacsBase :bizType="BizType.specialLine" />
</template>

<script setup lang="ts">
  import DeviceMacsBase from './DeviceMacsBase.vue';
  import { BizType } from '@/views/business/k/data';
</script>
