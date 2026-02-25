<template>
  <div class="reports-container">
    <Tabs v-model:activeKey="activeKey" @change="handleTabChange" destroyInactiveTabPane>
      <template v-for="tab in tabs" :key="tab.key">
        <Tabs.TabPane :tab="tab.label" :key="tab.key" v-if="hasPermission(tab.permission)">
          <template v-if="tab.key === 'overview'">
            <UtilizationRateOverview />
          </template>
          <!-- <template v-if="tab.key === 'assessment'">
            <UtilizationRateAssessment />
          </template> -->
          <template v-if="tab.key === 'biz'">
            <bizUtilizationRate />
          </template>
        </Tabs.TabPane>
      </template>
    </Tabs>
  </div>
</template>

<script lang="ts" setup>
  // import UtilizationRateAssessment from './components/utilizationRateAssessment.vue';
  import UtilizationRateOverview from './components/utilizationRateOverview.vue';
  import bizUtilizationRate from './components/bizUtilizationRate.vue';
  import { Tabs } from 'ant-design-vue';
  import { onMounted, ref } from 'vue';
  import { usePermission } from '@/hooks/web/usePermission';
  import { useRoute } from 'vue-router';

  const { hasPermission } = usePermission();
  const route = useRoute();
  defineOptions({ name: 'CompositeReports' });
  const tabs: any[] = [
    {
      label: '利用率概览',
      key: 'overview',
      permission: 'business:A:utilization:overview:view',
    },
    // {
    //   label: '利用率考核',
    //   key: 'assessment',
    //   permission: 'business:A:utilization:rate:assessment',
    // },
    {
      label: '业务利用率',
      key: 'biz',
      permission: 'business:A:utilization:rate:biz:view',
    },
  ];
  const defaultTab = getDefaultTab();
  const activeKey = ref(defaultTab);

  onMounted(() => {
    setUrl(activeKey.value);
  });

  const handleTabChange = async (key: any) => {
    activeKey.value = key;
    setUrl(key);
  };
  function setUrl(key: any) {
    const url = new URL(window.location.href);
    url.searchParams.set('tab', key);
    window.history.replaceState({}, '', url.toString());
  }
  function getDefaultTab(): string {
    const queryTab = route.query.tab as string;
    for (const tab of tabs) {
      if (queryTab && queryTab == tab.key && hasPermission(tab.permission)) {
        return tab.key;
      }
    }
    for (const tab of tabs) {
      if (hasPermission(tab.permission)) {
        return tab.key;
      }
    }
    return tabs[0].key;
  }
</script>
<style lang="less" scoped>
  :deep(.ant-tabs-nav-list) {
    margin-left: 20px !important;
  }

  .vben-basic-table-form-container {
    padding-top: 0 !important;
  }

  :deep(.ant-tabs-top > .ant-tabs-nav) {
    margin: 0 0 8px !important;
  }

  :deep(.ant-card .ant-card-body) {
    padding: 6px !important;
  }

  :deep(.vben-basic-table-form-container .ant-form) {
    margin-bottom: 8px !important;
  }

  :deep(.ant-picker) {
    width: 100% !important;
  }
</style>
