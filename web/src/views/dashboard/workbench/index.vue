<template>
  <PageWrapper>
    <template #headerContent> <WorkbenchHeader /> </template>

    <div class="lg:flex"> </div>
  </PageWrapper>
</template>

<script lang="ts" setup>
  import { ref } from 'vue';
  // import { Card } from 'ant-design-vue';
  import { PageWrapper } from '@/components/Page';
  import WorkbenchHeader from './components/WorkbenchHeader.vue';

  const loading = ref(true);

  setTimeout(() => {
    loading.value = false;
  }, 1500);
</script>
