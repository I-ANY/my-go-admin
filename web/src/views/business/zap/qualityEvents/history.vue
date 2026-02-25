<template>
  <BasicTable @register="registerTable">
    <template #headerCell="{ column }">
      <template v-if="column.key == 'hostname'">
        主机名
        <CopyOutlined class="ml-2" @click="copy_data('hostname')" />
      </template>
      <template v-if="column.key == 'node'">
        节点label
        <CopyOutlined class="ml-2" @click="copy_data('node')" />
      </template>
    </template>
    <template #bodyCell="{ column, record }">
      <template v-if="column.key == 'fields'">
        <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '800px' }">
          <template #title>{{ record.fields }}</template>
          <span>{{ truncatedFields(record.fields) }}</span>
        </Tooltip>
      </template>
    </template>
  </BasicTable>
</template>

<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { columns, searchFormSchema, truncatedFields } from './data';
  import { GetQualityEvents } from '@/api/business/zap';
  import { message, Tooltip } from 'ant-design-vue';
  import { splitByLineAndTrim } from '@/utils/util';
  import { CopyOutlined } from '@ant-design/icons-vue';

  defineOptions({ name: 'ZapDevice' });
  const [registerTable, { getDataSource }] = useTable({
    title: '质量事件列表(历史)',
    api: GetQualityEvents,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 2,
      actionColOptions: { span: 12 },
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    scroll: {
      y: 500,
    },
    pagination: {
      // pageSizeOptions: ['10', '30', '50'],
    },
    beforeFetch: (params) => {
      params.nodes = splitByLineAndTrim(params.nodes);
      params.hostnames = splitByLineAndTrim(params.hostnames);
    },
  });

  function copy_data(column: string) {
    const data = getDataSource()
      .map((item) => item[column])
      .join('\n');
    navigator.clipboard.writeText(data);
    message.success('已复制到剪切板');
  }

  // // 定义一个方法来截取前五十个字符
  // function truncatedFields(fields: string) {
  //   const jdata = JSON.stringify(fields);
  //   return jdata.length > 120 ? jdata.slice(0, 120) + '...' : jdata;
  // }
</script>
