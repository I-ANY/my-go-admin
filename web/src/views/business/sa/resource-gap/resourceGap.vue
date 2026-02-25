<template>
  <div>
    <BasicTable @register="registerTable">
      <template #bodyCell="{ column, record }">
        <template
          v-if="column.key == 'bandwidth_gap' && (record.status == 0 || record.status == 3)"
        >
          <span style="color: red"> {{ record.bandwidth_gap }} </span>
        </template>
        <template
          v-if="
            column.key == 'previous_bandwidth_gap' && (record.status == 0 || record.status == 3)
          "
        >
          <span style="color: red"> {{ record.previous_bandwidth_gap }} </span>
        </template>
      </template>
    </BasicTable>
  </div>
</template>

<script setup lang="ts">
  import { BasicTable, useTable } from '@/components/Table';
  import {
    GetResourceGapList,
    GetResourceGapTypeList,
    GetResourceGapProvinceList,
  } from '@/api/business/sa';
  import {
    resourceGapTableColumns,
    resourceGapSearchSchema,
  } from '@/views/business/sa/resource-gap/data';
  import { onMounted } from 'vue';

  const [registerTable, { getForm }] = useTable({
    title: '带宽资源缺口详情',
    api: GetResourceGapList,
    columns: resourceGapTableColumns,
    formConfig: {
      labelWidth: 100,
      schemas: resourceGapSearchSchema,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // 监听输入变化，防抖提交
      submitOnChange: true,
      // alwaysShowLines: 2,
      actionColOptions: {
        span: 4,
      },
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      // pageSizeOptions: ['10', '30', '50', '100', '500', '1000'],
    },
  });

  // 获取到资源缺口类型列表
  async function getResourceGapTypeOptions() {
    try {
      const res = await GetResourceGapTypeList();
      getForm().updateSchema({
        field: 'gap_type',
        componentProps: {
          options: res,
        },
      });
    } catch (e) {
      console.log(e);
    }
  }

  // 获取到省份信息
  async function getResourceGapProvinceList() {
    try {
      const res = await GetResourceGapProvinceList();
      getForm().updateSchema({
        field: 'province',
        componentProps: {
          options: res,
          filterOption: (input: string, option: any) => {
            return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
          },
        },
      });
    } catch (e) {
      console.log(e);
    }
  }

  onMounted(async () => {
    await getResourceGapTypeOptions();
    await getResourceGapProvinceList();
  });
</script>
<style scoped lang="less"></style>
