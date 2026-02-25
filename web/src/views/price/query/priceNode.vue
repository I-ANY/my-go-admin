<template>
  <div>
    <BasicTable @register="registerTable">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'priceType'">
          <span>{{ PriceTypeMap[record.priceType] }}</span>
        </template>
        <template v-if="column.key == 'bizs'">
          <a-button type="primary" @click="handleQueryBiz(record.bizs)">业务详情</a-button>
        </template>
      </template>
    </BasicTable>
    <Modal
      :bodyStyle="{ height: '300px', overflowY: 'auto' }"
      v-model:open="open"
      title="业务详情"
      @ok="open = false"
    >
      <template v-for="item in bizs" :key="item">
        <div
          style="
            display: flex;
            align-items: center;
            padding: 8px 16px;
            border-bottom: 1px solid #f0f0f0;
          "
        >
          <span style="flex: 1; text-align: left">{{ item }}</span>
        </div>
      </template>
    </Modal>
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, BasicColumn, FormSchema } from '@/components/Table';
  import { getNodesByCategory, getCategory } from '@/api/price/price';
  import { ref } from 'vue';
  import { Modal } from 'ant-design-vue';

  defineOptions({ name: 'PriceNodeCategoryQuery' });

  const bizs = ref([]);
  const open = ref(false);

  const PriceTypeMap = {
    0: '未知',
    1: '日95(集群日95)',
    2: '单机日95',
    3: '买断',
    4: '月95',
    5: '单口月95',
  };

  const cellContent = () => ({
    colSpan: 1,
  });

  function getMergeHeaderColumns(): BasicColumn[] {
    return [
      {
        title: '级别',
        dataIndex: 'level',
        width: 60,
        customCell: (value: any, index: number | undefined) => {
          if (index === undefined) return { rowSpan: 0 };
          const dataSource = getDataSource();
          if (!dataSource) return { rowSpan: 0 };
          if (index > 0 && dataSource[index - 1].level === value.level) {
            return { rowSpan: 0 };
          }
          let rowSpan = 1;
          for (let i = index + 1; i < dataSource.length; i++) {
            if (dataSource[i].level === value.level) {
              rowSpan++;
            } else {
              break;
            }
          }
          return { rowSpan };
        },
      },
      {
        title: '节点',
        dataIndex: 'owner',
        minWidth: 100,
        maxWidth: 300,
        customCell: cellContent,
      },
      {
        title: '运营商',
        dataIndex: 'localIsp',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '所在地',
        dataIndex: 'location',
        minWidth: 100,
        maxWidth: 300,
        customCell: cellContent,
      },
      {
        title: '业务大类',
        dataIndex: 'categoryName',
        width: 150,
        customCell: cellContent,
      },
      {
        title: '业务详情',
        dataIndex: 'bizs',
        customCell: cellContent,
      },
      {
        title: '计费方式',
        dataIndex: 'priceType',
        customCell: cellContent,
      },
    ];
  }
  const searchFormSchema: FormSchema[] = [
    {
      label: '业务组',
      field: 'name',
      component: 'ApiSelect',
      componentProps: {
        api: getCategory,
        params: { pageSize: 50000, pageIndex: 1 },
        resultField: 'items',
        labelField: 'name',
        valueField: 'name',
        showSearch: true,
        filterOption: (input: string, option: any) => {
          return (
            option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
            option.identify?.toLowerCase().indexOf(input.toLowerCase()) >= 0
          );
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'sort',
      label: '排序',
      component: 'Select',
      componentProps: {
        options: [
          { label: '升序', value: 1 },
          { label: '降序', value: 0 },
        ],
      },
      defaultValue: 0,
      colProps: { span: 4 },
    },
  ];

  const [registerTable, { getDataSource }] = useTable({
    title: '业务查询列表',
    api: getNodesByCategory,
    columns: getMergeHeaderColumns(),
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
    },
    useSearchForm: true,
    bordered: true,
    pagination: false,
    showIndexColumn: false,
    showTableSetting: true,
  });

  function handleQueryBiz(record) {
    bizs.value = record;
    open.value = true;
  }
</script>
