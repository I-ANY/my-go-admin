<template>
  <DemandOccupyList :bizType="BizType.specialLine" />
</template>

<script setup lang="ts">
  import DemandOccupyList from './DemandOccupyList.vue';
  import { BizType } from '@/views/business/k/data';
</script>
