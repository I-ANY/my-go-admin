<template>
  <div class="biz-container">
    <Tabs v-model:activeKey="activeKey" @change="handleTabChange" destroyInactiveTabPane>
      <template v-for="tab in tabs" :key="tab.key">
        <Tabs.TabPane :tab="tab.label" :key="tab.key" v-if="hasPermission(tab.permission)">
          <template v-if="tab.key === 'Rulemanage'">
            <Rulemanage />
          </template>
          <template v-if="tab.key === 'NoRuleBusiness'">
            <NoRuleBusiness />
          </template>
        </Tabs.TabPane>
      </template>
    </Tabs>
  </div>
</template>

<script lang="ts" setup>
  import Rulemanage from './rulemanage.vue';
  import NoRuleBusiness from './NoRuleBusiness.vue';
  import { Tabs } from 'ant-design-vue';
  import { onMounted, ref } from 'vue';
  import { usePermission } from '@/hooks/web/usePermission';
  import { useRoute } from 'vue-router';

  const { hasPermission } = usePermission();
  const route = useRoute();
  defineOptions({ name: 'CompositeReports' });
  const tabs: any[] = [
    {
      label: '规则管理',
      key: 'Rulemanage',
      permission: 'business:overprovisioning:rule:view',
    },
    {
      label: '无超配规则业务',
      key: 'NoRuleBusiness',
      permission: 'business:overprovisioning:no-rule:view',
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
