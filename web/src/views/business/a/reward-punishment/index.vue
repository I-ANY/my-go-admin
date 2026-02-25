<template>
  <div class="reward-punishment-container">
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" :loading="exportState.loading" @click="handleExport">
          {{ exportState.text }}
        </a-button>
        <!-- <a-button style="margin-left: 8px" @click="handleReset">重置</a-button> -->
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column && (column.key === 'action' || column.dataIndex === 'action')">
          <!-- 只在业务组第一行显示操作按钮，其他行通过 rowSpan 合并 -->
          <TableAction
            v-if="record && record.businessGroupRowSpan && record.businessGroupRowSpan > 0"
            :actions="[
              {
                label: '历史方案',
                icon: 'ant-design:history-outlined',
                tooltip: '查看历史方案',
                onClick: () => handleViewPlanHistory(record),
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <!-- 历史方案内容查看模态框 -->
    <HistoryPlanModal @register="registerPlanModal" />
  </div>
</template>

<script lang="ts" setup>
  import { reactive, nextTick, h } from 'vue';
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { Modal } from 'ant-design-vue';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { columns, searchFormSchema } from './data';
  import { getRewardPunishmentDetails, Api } from '@/api/business/a';
  import { downloadFileByUrl } from '@/utils/download';
  import { useMessage } from '@/hooks/web/useMessage';
  import HistoryPlanModal from './HistoryPlanModal.vue';
  import dayjs from 'dayjs';
  import { useModal } from '@/components/Modal';

  defineOptions({ name: 'RewardPunishmentList' });

  // 历史方案模态框
  const [registerPlanModal, { openModal: openPlanModal }] = useModal();

  const { notification } = useMessage();

  const exportState = reactive({
    loading: false,
    text: '导出',
  });

  // 处理合并单元格的数据
  // 由于已禁用分页，可以基于所有数据进行合并单元格处理
  const processMergeData = (dataSource: any[]) => {
    if (!dataSource || !Array.isArray(dataSource) || dataSource.length === 0) {
      return [];
    }

    // 过滤掉 null/undefined 的项，并创建新对象避免修改原数据
    const processedData = dataSource.filter((item) => item != null).map((item) => ({ ...item }));

    // 按业务组、考核业务、节点编号、日期排序，确保相同内容的数据相邻
    processedData.sort((a, b) => {
      if (a.businessGroup !== b.businessGroup) {
        return (a.businessGroup || '').localeCompare(b.businessGroup || '');
      }
      if (a.assessmentBiz !== b.assessmentBiz) {
        return (a.assessmentBiz || '').localeCompare(b.assessmentBiz || '');
      }
      if (a.roomNo !== b.roomNo) {
        return (a.roomNo || '').localeCompare(b.roomNo || '');
      }
      if (a.date !== b.date) {
        return (a.date || '').localeCompare(b.date || '');
      }
      return 0;
    });

    // 初始化所有行的rowSpan为1
    processedData.forEach((item) => {
      item.businessGroupRowSpan = 1;
      item.businessRowSpan = 1;
    });

    // 计算业务组的合并行数
    for (let i = 0; i < processedData.length; i++) {
      if (processedData[i].businessGroupRowSpan > 0) {
        let businessGroupCount = 1;
        for (let j = i + 1; j < processedData.length; j++) {
          if (processedData[j].businessGroup === processedData[i].businessGroup) {
            businessGroupCount++;
            processedData[j].businessGroupRowSpan = 0;
          } else {
            break;
          }
        }
        processedData[i].businessGroupRowSpan = businessGroupCount;
      }
    }

    // 在每个业务组内，计算考核业务的合并行数
    let groupStart = 0;
    while (groupStart < processedData.length) {
      // 找到业务组的开始和结束
      let groupEnd = groupStart;
      while (
        groupEnd + 1 < processedData.length &&
        processedData[groupEnd + 1].businessGroup === processedData[groupStart].businessGroup
      ) {
        groupEnd++;
      }

      // 在业务组内处理考核业务的合并
      for (let i = groupStart; i <= groupEnd; i++) {
        if (processedData[i].businessRowSpan > 0) {
          let businessCount = 1;
          for (let j = i + 1; j <= groupEnd; j++) {
            if (processedData[j].assessmentBiz === processedData[i].assessmentBiz) {
              businessCount++;
              processedData[j].businessRowSpan = 0;
            } else {
              break;
            }
          }
          processedData[i].businessRowSpan = businessCount;
        }
      }

      groupStart = groupEnd + 1;
    }

    return processedData;
  };

  // 转换日期范围为 startDate 和 endDate
  const transformDateRange = (params: any) => {
    const transformedParams = { ...params };

    // 如果存在 dateRange，转换为 startDate 和 endDate
    if (transformedParams.dateRange && Array.isArray(transformedParams.dateRange)) {
      transformedParams.startDate = transformedParams.dateRange[0];
      transformedParams.endDate = transformedParams.dateRange[1];
      // 删除 dateRange，避免传递给后端
      delete transformedParams.dateRange;
    }

    return transformedParams;
  };

  // 表格配置
  const [registerTable, { getForm }] = useTable({
    title: '奖惩明细表',
    api: getRewardPunishmentDetails,
    columns,
    formConfig: {
      labelWidth: 100,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
    },
    useSearchForm: true,
    bordered: true,
    showIndexColumn: false, // 关闭索引列，避免与合并单元格冲突
    showTableSetting: true,
    canResize: true,
    rowKey: 'id',
    // 禁用分页，一次性加载所有数据
    pagination: false,
    searchInfo: {
      dateRange: [dayjs().subtract(7, 'day').format('YYYY-MM-DD'), dayjs().format('YYYY-MM-DD')],
    },
    actionColumn: {
      width: 120,
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      fixed: 'right',
      // 操作列按业务组合并单元格
      customCell: (record: any) => {
        if (!record) return {};
        const rowSpan = record.businessGroupRowSpan ?? 1;
        if (rowSpan === 0) {
          return { rowSpan: 0 };
        }
        return { rowSpan };
      },
    },
    // 请求前处理参数，将 dateRange 转换为 startDate 和 endDate
    beforeFetch: (params) => {
      return transformDateRange(params);
    },
    // 处理返回的数据
    afterFetch: (dataSource) => {
      if (!dataSource || !Array.isArray(dataSource)) {
        return [];
      }
      const processedData = processMergeData(dataSource);
      return processedData || [];
    },
  });

  // 查看历史方案内容
  const handleViewPlanHistory = (record: any) => {
    // 从搜索表单获取日期范围
    const form = getForm();
    const searchValues = form.getFieldsValue() || {};
    const dateRange = searchValues.dateRange || [
      dayjs().subtract(7, 'day').format('YYYY-MM-DD'),
      dayjs().format('YYYY-MM-DD'),
    ];

    openPlanModal(true, {
      record,
      startDate: dateRange[0],
      endDate: dateRange[1],
      businessGroup: record.businessGroup,
    });
  };

  // 导出功能
  function handleExport() {
    Modal.confirm({
      title: '确认导出',
      icon: h(ExclamationCircleOutlined),
      content: '是否导出当前筛选条件下的奖惩明细？',
      okText: '导出',
      cancelText: '取消',
      async onOk() {
        const form = await getForm();
        await form.validate();
        const values = await form.getFieldsValue();

        // 如果没有日期范围，使用默认值
        if (!values.dateRange) {
          values.dateRange = [
            dayjs().subtract(7, 'day').format('YYYY-MM-DD'),
            dayjs().format('YYYY-MM-DD'),
          ];
        }

        // 转换日期范围
        const payload = transformDateRange({ ...values });

        // 移除空值和未定义的参数（但保留 startDate 和 endDate）
        Object.keys(payload).forEach((key) => {
          if (key === 'startDate' || key === 'endDate') {
            return; // 保留日期参数
          }
          if (payload[key] === undefined || payload[key] === null || payload[key] === '') {
            delete payload[key];
          }
          // 如果是空数组，也删除
          if (Array.isArray(payload[key]) && payload[key].length === 0) {
            delete payload[key];
          }
        });

        // 确保至少要有日期参数
        if (!payload.startDate || !payload.endDate) {
          notification.error({
            message: '导出失败',
            description: '请选择日期范围',
          });
          return;
        }

        console.log('导出参数:', payload);

        nextTick(() => {
          exportState.loading = true;
          exportState.text = '导出中...';
        });

        try {
          let filename = await downloadFileByUrl(
            Api.RewardPunishmentDetailsExport,
            'POST',
            5 * 60,
            payload,
            null,
          );
          notification.success({
            message: '导出成功',
            description: '文件名：' + filename,
            duration: null,
          });
        } catch (error) {
          notification.error({
            message: '导出失败',
            description: error.message,
            duration: null,
          });
        } finally {
          nextTick(() => {
            exportState.loading = false;
            exportState.text = '导出';
          });
        }
      },
    });
  }
</script>

<style lang="scss" scoped>
  .reward-punishment-container {
    padding: 16px;

    :deep(.ant-picker) {
      width: 100% !important;
    }
  }
</style>
