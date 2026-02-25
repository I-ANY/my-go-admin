<template>
  <div>
    <!-- 单独显示某个业务类型时，不显示 Tabs -->
    <template v-if="props.bizType !== 'all'">
      <!-- 汇聚表格 -->
      <BasicTable v-if="props.bizType === 'normal'" @register="tableConfigs.normal.register">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'result'">
            <a-tooltip v-if="record.result" :title="record.result">
              <span class="truncate-text">{{ record.result }}</span>
            </a-tooltip>
            <span v-else>-</span>
          </template>
        </template>
      </BasicTable>
      <!-- 专线表格 -->
      <BasicTable
        v-if="props.bizType === 'specialLine'"
        @register="tableConfigs.specialLine.register"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key === 'result'">
            <a-tooltip v-if="record.result" :title="record.result">
              <span class="truncate-text">{{ record.result }}</span>
            </a-tooltip>
            <span v-else>-</span>
          </template>
        </template>
      </BasicTable>
    </template>

    <!-- 显示全部时使用 Tabs -->
    <Tabs v-else v-model:activeKey="activeTabKey" type="card">
      <TabPane tab="汇聚" key="normal" v-auth="'business:k:occupytask:normal'">
        <BasicTable @register="tableConfigs.normal.register">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'result'">
              <a-tooltip v-if="record.result" :title="record.result">
                <span class="truncate-text">{{ record.result }}</span>
              </a-tooltip>
              <span v-else>-</span>
            </template>
          </template>
        </BasicTable>
      </TabPane>
      <TabPane tab="专线" key="specialLine" v-auth="'business:k:occupytask:specialline'">
        <BasicTable @register="tableConfigs.specialLine.register">
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'result'">
              <a-tooltip v-if="record.result" :title="record.result">
                <span class="truncate-text">{{ record.result }}</span>
              </a-tooltip>
              <span v-else>-</span>
            </template>
          </template>
        </BasicTable>
      </TabPane>
    </Tabs>
  </div>
</template>

<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { Tabs, TabPane } from 'ant-design-vue';
  import { OccupyTaskColumns } from './data';
  import { DemandOccupyTaskList } from '@/api/business/k';
  import type { FormSchema } from '@/components/Form';
  import { ref, onMounted } from 'vue';
  import { usePermission } from '@/hooks/web/usePermission';

  defineOptions({ name: 'OccupyTaskList' });

  // 定义 props
  const props = withDefaults(
    defineProps<{
      bizType?: 'normal' | 'specialLine' | 'all';
    }>(),
    {
      bizType: 'all',
    },
  );

  // 权限检查
  const { hasPermission } = usePermission();

  // 当前激活的标签页
  const activeTabKey = ref<'normal' | 'specialLine'>('normal');

  // 检查权限并设置默认激活的 tab
  const checkPermissionsAndSetDefaultTab = () => {
    // 如果指定了具体的业务类型，直接设置
    if (props.bizType !== 'all') {
      activeTabKey.value = props.bizType;
      return;
    }

    const hasNormalPermission = hasPermission('business:k:occupytask:normal');
    const hasSpecialLinePermission = hasPermission('business:k:occupytask:specialline');

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

  // 状态选项
  const statusOptions = [
    { label: '待执行', value: 0 },
    { label: '执行中', value: 1 },
    { label: '成功', value: 2 },
    { label: '失败', value: 3 },
  ];

  // 搜索表单配置
  const searchFormSchema: FormSchema[] = [
    {
      field: 'demand_id',
      label: '需求ID',
      component: 'Input',
      componentProps: {
        placeholder: '请输入需求ID',
        allowClear: true,
      },
      colProps: { span: 6 },
    },
    {
      field: 'status',
      label: '状态',
      component: 'Select',
      componentProps: {
        placeholder: '请选择状态',
        options: statusOptions,
        allowClear: true,
      },
      colProps: { span: 6 },
    },
  ];

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
    const title = bizType === 'normal' ? '汇聚占用任务列表' : '专线占用任务列表';

    const [register, { getForm, reload }] = useTable({
      title,
      api: async (params) => {
        const result = await DemandOccupyTaskList(params);
        return result;
      },
      columns: OccupyTaskColumns,
      formConfig: {
        labelWidth: 80,
        schemas: searchFormSchema,
        autoSubmitOnEnter: true,
        showAdvancedButton: true,
        autoAdvancedLine: 4,
      },
      beforeFetch(params) {
        params.biz_type = bizType;
        return params;
      },
      useSearchForm: true,
      showTableSetting: true,
      bordered: true,
      showIndexColumn: false,
      pagination: {},
      rowKey: 'id',
    });

    // 保存表格配置
    tableConfigs[bizType] = {
      register,
      getForm,
      reload,
    };
  };

  // 创建两个表格的配置
  createTableConfig('normal');
  createTableConfig('specialLine');

  // 获取数据时更新树形数据
  onMounted(async () => {
    // 检查权限并设置默认激活的 tab
    checkPermissionsAndSetDefaultTab();
  });
</script>

<style scoped>
  .truncate-text {
    display: inline-block;
    max-width: 200px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
</style>
