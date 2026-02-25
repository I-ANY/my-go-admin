<template>
    <!-- <div class="reports-container"> -->
        <Tabs v-model:activeKey="activeKey" @change="handleTabChange" destroyInactiveTabPane>
            <template v-for="tab in availableTabs" :key="tab.key">
                <Tabs.TabPane :tab="tab.label">
                    <template v-if="tab.key === 'businessPackpageConfig'">
                        <BusinessPackpageConfig />
                    </template>
                    <template v-if="tab.key === 'taskBuildCenter'">
                        <TaskBuildCenter />
                    </template>
                </Tabs.TabPane>
            </template>
        </Tabs>
    <!-- </div> -->
</template>

<script lang="ts" setup>
import BusinessPackpageConfig from './businessPackpageConfig/index.vue';
import TaskBuildCenter from './taskBuildCenter/index.vue';
import { Tabs } from 'ant-design-vue';
import { ref, computed, watch } from 'vue';
import { usePermission } from '@/hooks/web/usePermission';

const { hasPermission } = usePermission();
defineOptions({ name: 'CompositeReports' });
const tabs = computed(() => [
    {
        label: '业务包配置',
        key: 'businessPackpageConfig',
        permission: 'business:encryption:businessPackpageConfig:view',
    },
    {
        label: '定制系统构建',
        key: 'taskBuildCenter',
        permission: 'business:encryption:taskBuildCenter:view',
    },
]);

// 过滤出有权限的tabs
const availableTabs = computed(() => {
    return tabs.value.filter(tab => hasPermission(tab.permission));
});

// 获取第一个有权限的tab的key（如果只有一个tab有权限，就选中它；如果有多个，选中第一个）
const firstAvailableTabKey = computed(() => {
    return availableTabs.value[0]?.key || '';
});

const activeKey = ref(firstAvailableTabKey.value);

// 监听有权限的tabs变化，确保始终选中第一个有权限的tab
watch(firstAvailableTabKey, (newKey) => {
    if (newKey && newKey !== activeKey.value) {
        activeKey.value = newKey;
    }
}, { immediate: true });

const handleTabChange = async (key: any) => {
    activeKey.value = key;
};

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