<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="exportTraffic" :loading="loading">导出</a-button>
      </template>
      <template #bodyCell="{ record, column }">
        <template v-if="column.key == 'is_publish'">
          <Switch
            :checked="record.is_publish"
            :checked-children="'已发布'"
            :un-checked-children="'未发布'"
            :loading="switchTag[record.id] || false"
            @change="(checked) => handleTrafficStatus(record, checked)"
          />
        </template>
      </template>
    </BasicTable>
  </div>
</template>
<script setup lang="ts">
  import { BasicTable, useTable } from '@/components/Table';
  import { GetZPTraffic95List, UpdateZPTraffic95 } from '@/api/business/zp';
  import { traffic95Columns, traffic95SearchFormSchema } from '@/views/business/zp/data';
  import { splitByLineAndTrim } from '@/utils/util';
  import { RangePickPresetsExact } from '@/utils/common';
  import { onMounted, reactive, ref } from 'vue';
  import dayjs from 'dayjs';
  import { message, Switch } from 'ant-design-vue';

  let loading = ref(false);
  let switchTag: Record<string, boolean> = reactive({});

  const [registerTable, { getForm, reload }] = useTable({
    title: '95值统计',
    api: GetZPTraffic95List,
    columns: traffic95Columns,
    formConfig: {
      labelWidth: 120,
      schemas: traffic95SearchFormSchema(onTimePikerOpen),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      pageSizeOptions: ['10', '30', '50', '100', '500', '1000'],
    },
    beforeFetch: (params) => {
      params.hostnames = splitByLineAndTrim(params.hostnames) || null;
      params.device_ids = splitByLineAndTrim(params.device_ids) || null;
    },
  });

  function onTimePikerOpen() {
    getForm().updateSchema({
      field: '[startTimestmp, endTimestmp]',
      componentProps: {
        presets: RangePickPresetsExact(),
      },
    });
  }

  onMounted(async () => {
    await resetReportTime();
  });

  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      startTimestmp: dayjs(dayjs().add(-1, 'day').format('YYYY-MM-DD'), 'YYYY-MM-DD'),
      endTimestmp: dayjs(dayjs().format('YYYY-MM-DD'), 'YYYY-MM-DD'),
    });
  }

  function exportTraffic() {
    console.log('exportTraffic');
  }

  async function handleTrafficStatus(record, checked) {
    console.log(checked);
    switchTag[record.id] = true;
    record.is_publish = checked;
    try {
      await UpdateZPTraffic95({ ...record });
      message.success('操作成功');
      switchTag[record.id] = false;
      await reload();
    } catch (e) {
      switchTag[record.id] = false;
      await reload();
    }
  }
</script>

<style scoped lang="less"></style>
