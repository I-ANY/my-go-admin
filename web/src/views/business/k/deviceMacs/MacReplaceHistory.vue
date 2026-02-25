<template>
  <div>
    <BasicTable @register="registerTable">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'is_province_scheduling'">
          <span v-if="isProvinceSchedulingMap[record.is_province_scheduling]">{{
            isProvinceSchedulingMap[record.is_province_scheduling].dictLabel ||
            record.is_province_scheduling
          }}</span>
        </template>
        <template v-if="column.key == 'is_cover_diff_isp'">
          <span v-if="isCoverDiffIspMap[record.is_cover_diff_isp]">{{
            isCoverDiffIspMap[record.is_cover_diff_isp].dictLabel || record.is_cover_diff_isp
          }}</span>
        </template>
        <template v-if="column.key == 'status'">
          <Tooltip>
            <Tag :color="statusColorMap[record.status]">
              {{ statusLabelMap[record.status] || record.status }}
            </Tag>
            <template #title v-if="record.remark">
              <span>{{ record.remark }}</span>
            </template>
          </Tooltip>
        </template>
      </template>
    </BasicTable>
  </div>
</template>
<script setup lang="ts">
  import { BasicTable, useTable } from '@/components/Table';
  import {
    macReplaceHistoryColumns,
    MacReplaceSearchFormSchema,
    isCoverDiffIspMap,
    isProvinceSchedulingMap,
  } from './data';
  import { GetMacReplaceHistory } from '@/api/business/k';
  import { Tag, Tooltip } from 'ant-design-vue';

  // 新增状态颜色映射
  const statusColorMap = {
    0: 'orange',
    1: 'success',
    2: 'error',
    3: 'orange',
    4: 'blue',
    5: 'error',
  };

  // 新增状态标签映射
  const statusLabelMap = {
    0: '替换中',
    1: '成功',
    2: '失败',
    3: '异网下发中',
    4: '异网下发成功',
    5: '异网下发失败',
  };

  const [registerTable] = useTable({
    title: 'MAC替换记录',
    api: GetMacReplaceHistory,
    columns: macReplaceHistoryColumns,
    // 使用搜索表单配置
    useSearchForm: true,
    formConfig: {
      labelWidth: 120,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      alwaysShowLines: 2,
      schemas: MacReplaceSearchFormSchema(),
    },
    beforeFetch: (params) => {
      return params;
    },
    canResize: true,
    bordered: true,
    showIndexColumn: true,
    rowKey: 'id',
  });
</script>
