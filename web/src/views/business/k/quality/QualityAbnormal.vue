<template>
  <div>
    <Tabs v-model:activeKey="activeTabKey" type="card">
      <TabPane tab="汇聚" key="normal" v-auth="'business:k:qualityabnormal:normal'">
        <BasicTable @register="tableConfigs.normal.register">
          <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'incident_type'">
              <Tooltip :title="record.incident_type">
                <div class="text-ellipsis">{{ record.incident_type }}</div>
              </Tooltip>
            </template>
            <template v-if="column.dataIndex === 'status'">
              <Tooltip :title="record.status">
                <div class="text-ellipsis">{{ record.status }}</div>
              </Tooltip>
            </template>
          </template>
        </BasicTable>
      </TabPane>
      <TabPane tab="专线" key="specialLine" v-auth="'business:k:qualityabnormal:specialline'">
        <BasicTable @register="tableConfigs.specialLine.register">
          <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'incident_type'">
              <Tooltip :title="record.incident_type">
                <div class="text-ellipsis">{{ record.incident_type }}</div>
              </Tooltip>
            </template>
            <template v-if="column.dataIndex === 'status'">
              <Tooltip :title="record.status">
                <div class="text-ellipsis">{{ record.status }}</div>
              </Tooltip>
            </template>
          </template>
        </BasicTable>
      </TabPane>
    </Tabs>
  </div>
</template>
<script setup lang="ts">
  import { BasicTable, useTable } from '@/components/Table';
  import { Tabs, TabPane, Tooltip } from 'ant-design-vue';
  import { qualityAbnormalColumns, qualityAbnormalSearchFormSchema } from './data';
  import { GetQualityAbnormal } from '@/api/business/k';
  import { usePermission } from '@/hooks/web/usePermission';
  import { ref, onMounted } from 'vue';
  import dayjs from 'dayjs';

  // 权限检查
  const { hasPermission } = usePermission();

  // 当前激活的标签页
  const activeTabKey = ref<'normal' | 'specialLine'>('normal');

  // 检查权限并设置默认激活的 tab
  const checkPermissionsAndSetDefaultTab = () => {
    const hasNormalPermission = hasPermission('business:k:qualityabnormal:normal');
    const hasSpecialLinePermission = hasPermission('business:k:qualityabnormal:specialline');

    if (hasNormalPermission && !hasSpecialLinePermission) {
      activeTabKey.value = 'normal';
    } else if (!hasNormalPermission && hasSpecialLinePermission) {
      activeTabKey.value = 'specialLine';
    } else if (hasNormalPermission && hasSpecialLinePermission) {
      // 两个都有权限，保持默认值
      activeTabKey.value = 'normal';
    } else {
      // 都没有权限，设置为 normal（这种情况应该不会发生，因为页面本身应该有权限控制）
      activeTabKey.value = 'normal';
    }
  };

  // 表格配置对象
  const tableConfigs = {
    normal: {
      register: null as any,
      getForm: null as any,
      reload: null as any,
    },
    specialLine: {
      register: null as any,
      getForm: null as any,
      reload: null as any,
    },
  };

  // 创建表格配置的函数
  const createTableConfig = (bizType: 'normal' | 'specialLine') => {
    const title = bizType === 'normal' ? '汇聚质量异常' : '专线质量异常';

    // 从 searchFormSchema 中获取默认 provider 值
    const defaultProvider =
      qualityAbnormalSearchFormSchema(bizType).find((item) => item.field === 'provider')
        ?.componentProps?.defaultValue || 'mf';

    const [register, { getForm, reload }] = useTable({
      title,
      api: GetQualityAbnormal,
      columns: qualityAbnormalColumns,
      // 使用搜索表单配置
      useSearchForm: true,
      formConfig: {
        labelWidth: 120,
        autoSubmitOnEnter: true,
        showAdvancedButton: true,
        autoAdvancedLine: 3,
        alwaysShowLines: 2,
        schemas: qualityAbnormalSearchFormSchema(bizType),
      },
      beforeFetch: (params) => {
        params.biz_type = bizType;
        if (!params.provider) params.provider = defaultProvider;

        // 设置默认时间范围为最近24小时
        if (!params.time_range) {
          params.time_range = [dayjs().subtract(24, 'hour'), dayjs()];
        }
        if (params.time_range) {
          params.begin_time = dayjs(params.time_range[0]).format('YYYY-MM-DD HH:mm:ss');
          params.end_time = dayjs(params.time_range[1]).format('YYYY-MM-DD HH:mm:ss');
        }
        delete params.time_range;
        return params;
      },
      canResize: true,
      bordered: true,
      showIndexColumn: true,
      rowKey: 'id',
    });

    // 保存表格配置
    tableConfigs[bizType] = {
      register: register as any,
      getForm: getForm as any,
      reload: reload as any,
    };
  };

  // 创建两个表格的配置
  createTableConfig('normal');
  createTableConfig('specialLine');

  onMounted(() => {
    // 检查权限并设置默认激活的 tab
    checkPermissionsAndSetDefaultTab();
  });
</script>

<style scoped>
  .text-ellipsis {
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
