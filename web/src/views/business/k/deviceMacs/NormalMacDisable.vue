<template>
  <MacDisable :bizType="BizType.normal" />
</template>

<script setup lang="ts">
  import MacDisable from './MacDisable.vue';
  import { BizType } from '@/views/business/k/data';
</script>
