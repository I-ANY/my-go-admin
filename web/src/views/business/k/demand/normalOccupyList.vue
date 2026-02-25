<template>
  <DemandOccupyList :bizType="BizType.normal" />
</template>

<script setup lang="ts">
  import DemandOccupyList from './DemandOccupyList.vue';
  import { BizType } from '@/views/business/k/data';
</script>
