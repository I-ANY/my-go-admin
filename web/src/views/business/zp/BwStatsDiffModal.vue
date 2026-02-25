<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    title="差异数据详情"
    :width="1200"
    :destroyOnClose="true"
    :footer="null"
  >
    <BasicTable @register="registerTable">
      <template #tableTitle>
        <span class="text-gray-500">时间点：{{ datetime }}</span>
      </template>
    </BasicTable>
  </BasicModal>
</template>

<script lang="ts" setup>
  import { h, ref } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicTable, useTable } from '@/components/Table';
  import { GetBwStatsDiffDetail } from '@/api/business/zp';
  import type { FormSchema } from '@/components/Form';
  import type { BasicColumn } from '@/components/Table';
  import { splitByLineAndTrim } from '@/utils/util';
  import { Tag } from 'ant-design-vue';

  defineOptions({ name: 'BwStatsDiffModal' });
  defineEmits(['register']);

  const datetime = ref('');
  const isps = ref<string[]>([]);

  // 表格列定义
  const columns: BasicColumn[] = [
    // {
    //   title: '时间',
    //   dataIndex: 'datetime',
    //   width: 150,
    // },
    {
      title: '机房节点',
      dataIndex: 'owner',
      width: 200,
      // 排序
      sorter: (a, b) => {
        return a.owner.localeCompare(b.owner);
      },
    },
    {
      title: '设备ID',
      dataIndex: 'device_id',
      width: 300,
    },
    {
      title: '主机名',
      dataIndex: 'provider_device_id',
      width: 200,
    },
    {
      title: '变更类型',
      dataIndex: 'change_type',
      width: 80,
      customRender: ({ record }) => {
        const color = record.change_type === 'added' ? 'green' : 'red';
        const text = record.change_type === 'added' ? '新增' : '删除';
        return h(Tag, { color: color }, () => text);
      },
    },
  ];

  // 搜索表单
  const searchFormSchema: FormSchema[] = [
    {
      field: 'hostnames',
      label: ' ',
      component: 'InputTextArea',
      colProps: { span: 6 },
      componentProps: {
        rows: 2,
        placeholder: '多个主机名，换行输入',
      },
    },
    {
      field: 'device_ids',
      label: ' ',
      component: 'InputTextArea',
      colProps: { span: 8 },
      componentProps: {
        rows: 2,
        placeholder: '多个设备ID，换行输入',
      },
    },
    {
      field: 'change_type',
      label: ' ',
      component: 'Select',
      colProps: { span: 4 },
      componentProps: {
        options: [
          {
            label: '新增',
            value: 'added',
          },
          {
            label: '删除',
            value: 'removed',
          },
        ],
      },
    },
  ];

  const [registerTable, { reload, setTableData }] = useTable({
    api: GetBwStatsDiffDetail,
    columns,
    formConfig: {
      labelWidth: 20,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
    },
    useSearchForm: true,
    showTableSetting: false,
    bordered: true,
    showIndexColumn: false,
    immediate: false,
    pagination: {
      pageSizeOptions: ['10', '30', '50', '100'],
    },
    beforeFetch: (params) => {
      // 处理多行输入
      if (params.hostnames) {
        params.hostnames = splitByLineAndTrim(params.hostnames);
      }
      if (params.device_ids) {
        params.device_ids = splitByLineAndTrim(params.device_ids);
      }
      params.datetime = datetime.value;
      return params;
    },
  });

  const [registerModal] = useModalInner(async (data) => {
    datetime.value = data.datetime || '';
    isps.value = data.isps || [];
    // 清空之前的数据
    setTableData([]);
    // 加载数据
    await reload();
  });
</script>
