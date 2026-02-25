<template>
  <MacDisable :bizType="BizType.specialLine" />
</template>

<script setup lang="ts">
  import MacDisable from './MacDisable.vue';
  import { BizType } from '@/views/business/k/data';
</script>
