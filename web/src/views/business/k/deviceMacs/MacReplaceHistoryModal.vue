<template>
  <BasicModal v-bind="$attrs" @register="registerModal" :title="`${currentMac} 替换历史`">
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
  </BasicModal>
</template>

<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicTable, useTable } from '@/components/Table';
  import { GetMacReplaceHistory } from '@/api/business/k';
  import { ref } from 'vue';
  import { Tag, Tooltip } from 'ant-design-vue';
  import { isCoverDiffIspMap, isProvinceSchedulingMap, macReplaceHistoryColumns } from './data';

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

  defineProps<{
    macAddr?: string;
  }>();

  // 注册弹窗
  const [registerModal, { setModalProps }] = useModalInner(async (macAddr) => {
    // 接收父组件传递的mac地址和其他参数
    currentMac.value = macAddr || '';
    setModalProps({
      width: 1200,
      height: 900,
      destroyOnClose: true,
      showCancelBtn: false,
      showOkBtn: false,
    });
  });

  // 表格实例
  const [registerTable] = useTable({
    title: 'MAC替换历史',
    api: GetMacReplaceHistory,
    columns: macReplaceHistoryColumns,
    beforeFetch: (params) => {
      // 传递当前MAC地址到接口
      params.mac_addr = currentMac.value;
      return params;
    },
    canResize: true,
    scroll: { y: 300 },
    useSearchForm: false, // 无搜索框
    bordered: true,
    showIndexColumn: true,
    rowKey: 'id',
  });

  // 当前MAC地址
  const currentMac = ref('');
</script>
