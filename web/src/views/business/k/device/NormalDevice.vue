<template>
  <div>
    <BasicTable @register="registerTable">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'biz_type'">
          <span v-if="bizTypesMap[record.biz_type]">{{
            bizTypesMap[record.biz_type].dictLabel || record.biz_type
          }}</span>
        </template>
        <template v-if="column.key == 'provider'">
          <span v-if="providersMap[record.provider]">{{
            providersMap[record.provider].dictLabel || record.provider
          }}</span>
        </template>

        <template v-if="column.key == 'is_cover_diff_isp'">
          <span v-if="isCoverDiffIspMap[record.is_cover_diff_isp]">{{
            isCoverDiffIspMap[record.is_cover_diff_isp].dictLabel || record.is_cover_diff_isp
          }}</span>
        </template>
        <template v-if="column.key == 'is_province_scheduling'">
          <span v-if="isProvinceSchedulingMap[record.is_province_scheduling]">{{
            isProvinceSchedulingMap[record.is_province_scheduling].dictLabel ||
            record.is_province_scheduling
          }}</span>
        </template>
        <template v-if="column.key == 'cover_diff_isp_id'">
          <span v-if="coverDiffIspMap[record.cover_diff_isp_id]">{{
            coverDiffIspMap[record.cover_diff_isp_id].dictLabel
          }}</span>
        </template>
      </template>
    </BasicTable>
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';

  import {
    normalColumns,
    normalSearchFormSchema,
    coverDiffIspMap,
    providersMap,
    bizTypesMap,
    isCoverDiffIspMap,
    isProvinceSchedulingMap,
  } from './data';
  import { GetDeviceList } from '@/api/business/k';
  import { BizType } from '../data';
  import { splitByLineAndTrim } from '@/utils/util';

  defineOptions({ name: 'KNormalDevice' });
  const bizType = BizType.normal;
  const [registerTable] = useTable({
    title: '设备列表',
    api: GetDeviceList,
    columns: normalColumns,
    formConfig: {
      labelWidth: 120,
      schemas: normalSearchFormSchema,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
    },
    beforeFetch(params) {
      params.biz_type = bizType;
      params.hostnames = splitByLineAndTrim(params.hostnames);
      params.macAddrs = splitByLineAndTrim(params.macAddrs);
      return params;
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      // pageSizeOptions: ['1', '2', '5'],
    },
  });
</script>
