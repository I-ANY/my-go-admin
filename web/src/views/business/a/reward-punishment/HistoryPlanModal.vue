<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="getTitle"
    :width="1000"
    :showOkBtn="false"
    :showCancelBtn="false"
    :destroyOnClose="true"
  >
    <div class="history-plan-container">
      <div class="mt-4">
        <BasicTable @register="registerTable" />
      </div>
    </div>
  </BasicModal>
</template>

<script lang="ts" setup>
  import { ref, computed, nextTick } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicTable, useTable } from '@/components/Table';
  import dayjs from 'dayjs';
  import { getRewardPunishmentHistoryPlans } from '@/api/business/a';
  import { historyPlanColumns, historyPlanSearchFormSchema } from './data';

  defineOptions({ name: 'HistoryPlanModal' });

  const planData = ref<any>(null);
  const defaultDateRange = ref<[string, string] | undefined>(undefined);
  const businessGroup = ref<string>('');
  const defaultAssessmentBizs = ref<string[]>([]);

  // 将嵌套数据展开为扁平结构，每个 historyPlan 作为一行
  const processTableData = (sourceData: any[]) => {
    if (!sourceData || sourceData.length === 0) {
      return [];
    }

    const flattenedData: any[] = [];
    const data = [...sourceData];

    // 按扣分日期、考核业务、节点编号排序
    data.sort((a, b) => {
      if (a.date !== b.date) {
        return a.date.localeCompare(b.date);
      }
      if (a.assessmentBiz !== b.assessmentBiz) {
        return a.assessmentBiz.localeCompare(b.assessmentBiz);
      }
      if (a.roomNo !== b.roomNo) {
        return a.roomNo.localeCompare(b.roomNo);
      }
      return 0;
    });

    // 展开嵌套数据
    data.forEach((item) => {
      const historyPlans = item.historyPlans || [];
      if (historyPlans.length === 0) {
        // 如果没有历史方案，至少显示一行
        flattenedData.push({
          ...item,
          assessmentDate: '-',
          planContent: '-',
          rowSpan: 1,
        });
      } else {
        // 对历史方案按考核日期倒序排序
        const sortedPlans = [...historyPlans].sort((a: any, b: any) => {
          return b.assessmentDate.localeCompare(a.assessmentDate);
        });

        sortedPlans.forEach((plan: any, index: number) => {
          flattenedData.push({
            ...item,
            assessmentDate: plan.assessmentDate,
            planContent: plan.planContent,
            rowSpan: index === 0 ? sortedPlans.length : 0, // 第一行设置 rowSpan，其他行设为 0
          });
        });
      }
    });

    return flattenedData;
  };

  const [registerModal, { setModalProps }] = useModalInner(async (data) => {
    if (data) {
      setModalProps({ confirmLoading: false });

      // 设置默认日期范围和业务组
      const { startDate, endDate, businessGroup: bg } = data;
      const finalStartDate = startDate || dayjs().subtract(7, 'day').format('YYYY-MM-DD');
      const finalEndDate = endDate || dayjs().format('YYYY-MM-DD');
      defaultDateRange.value = [finalStartDate, finalEndDate];
      businessGroup.value = bg || '';

      planData.value = {
        startDate: finalStartDate,
        endDate: finalEndDate,
        businessGroup: bg,
      };

      // 设置搜索表单的默认值（使用 nextTick 确保表单已初始化）
      await nextTick();
      const form = getForm();
      if (form) {
        form.setFieldsValue({
          dateRange: defaultDateRange.value,
        });
        // 重新加载表格数据
        reload();
      }
    }
  });

  const getTitle = computed(() => `历史方案内容-${planData.value?.businessGroup}`);

  // 表格配置
  const [registerTable, { getForm, reload }] = useTable({
    title: '',
    api: getRewardPunishmentHistoryPlans,
    columns: historyPlanColumns,
    formConfig: {
      labelWidth: 80,
      wrapperCol: {
        span: 24,
      },
      schemas: historyPlanSearchFormSchema,
      autoSubmitOnEnter: true,
    },
    useSearchForm: true,
    bordered: true,
    showIndexColumn: false,
    showTableSetting: false,
    canResize: false, // 禁用自动调整大小，使用自定义高度
    pagination: false,
    scroll: { x: 'max-content', y: 500 }, // 自定义表格滚动区域大小
    size: 'small', // 表格尺寸
    // 设置默认搜索值
    searchInfo: computed(() => {
      if (defaultDateRange.value) {
        return {
          dateRange: defaultDateRange.value,
        };
      }
      return {};
    }),
    // 请求前处理参数
    beforeFetch: (params) => {
      // 将 dateRange 转换为 startDate 和 endDate
      if (params.dateRange && Array.isArray(params.dateRange) && params.dateRange.length === 2) {
        params.startDate = params.dateRange[0];
        params.endDate = params.dateRange[1];
        delete params.dateRange;
      }
      // 添加业务组参数
      if (businessGroup.value) {
        params.businessGroup = businessGroup.value;
      }
      return params;
    },
    // 处理返回的数据
    afterFetch: (dataSource) => {
      if (!dataSource || !Array.isArray(dataSource)) {
        return [];
      }

      // 从原始数据中提取所有唯一的考核业务
      const uniqueAssessmentBizs = new Set<string>();
      dataSource.forEach((item: any) => {
        if (item.assessmentBiz) {
          uniqueAssessmentBizs.add(item.assessmentBiz);
        }
      });
      defaultAssessmentBizs.value = Array.from(uniqueAssessmentBizs);

      // 设置表单的默认考核业务值
      nextTick(() => {
        const form = getForm();
        if (form && defaultAssessmentBizs.value.length > 0) {
          form.setFieldsValue({
            assessmentBizs: defaultAssessmentBizs.value,
          });
        }
      });

      return processTableData(dataSource);
    },
  });
</script>

<style lang="scss" scoped>
  .history-plan-container {
    .mt-4 {
      margin-top: 16px;
    }

    :deep(.ant-table-tbody > tr > td) {
      vertical-align: top;
    }
  }
</style>
