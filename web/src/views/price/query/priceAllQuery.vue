<template>
  <div>
    <BasicTable @register="registerTable">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'category'">
          <span>{{ record.category ? record.category.name : '/' }}</span>
        </template>
        <template v-if="column.key == 'bizs'">
          <a-button type="primary" @click="handleQueryBiz(record.categoryID)">详情</a-button>
        </template>
        <template v-if="column.key == 'region'">
          {{ record.region ? record.region.name : '/' }}
        </template>
        <template v-if="column.key == 'zone'">
          {{ record.zone ? record.zone.name : '/' }}
        </template>
        <template v-if="column.key == 'isNotlocalIsp'">
          {{ record.bizIspMode || record.bizIsp ? '是' : '否' }}
        </template>
        <template v-if="column.key == 'bizIsp'">
          {{ record.bizIsp ? record.bizIsp : '/' }}
        </template>
        <template v-if="column.key == 'bizIspMode'">
          {{ record.bizIspMode == 2 ? '异网计费' : record.bizIspMode == 1 ? '本网计费' : '/' }}
        </template>
        <template v-if="column.key == 'mode'">
          {{ record.mode ? record.mode.name : '/' }}
        </template>
      </template>
    </BasicTable>
    <Modal
      :bodyStyle="{ height: '300px', overflowY: 'auto' }"
      v-model:open="open"
      title="业务详情"
      @ok="open = false"
    >
      <template v-for="item in bizData.items" :key="item.id">
        <div
          style="
            display: flex;
            align-items: center;
            padding: 8px 16px;
            border-bottom: 1px solid #f0f0f0;
          "
        >
          <span style="flex: 3; text-align: left">{{ item.name || '无业务' }}</span>
        </div>
      </template>
    </Modal>
  </div>
</template>
<script lang="ts" setup>
  import { Modal, message } from 'ant-design-vue';
  import { ref, reactive } from 'vue';
  import { BasicTable, useTable, BasicColumn, FormSchema } from '@/components/Table';
  import { getAllRecord, getBiz, getZone, getRegion, getCategory } from '@/api/price/price';

  defineOptions({ name: 'PriceAllQuery' });

  const open = ref<boolean>(false);
  const bizData = reactive({ items: [] as { id: number; name: string }[] });

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
        title: '本网运营商',
        dataIndex: 'localIsp',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '业务组',
        dataIndex: 'category',
        width: 150,
        customCell: cellContent,
      },
      {
        title: '业务详情',
        dataIndex: 'bizs',
        width: 100,
        customCell: cellContent,
      },
      // {
      //   title: '溜缝业务',
      //   dataIndex: 'low',
      //   width: 100,
      //   customRender: ({ text }) => (text ? '是' : '否'),
      // },
      {
        title: '大区',
        dataIndex: 'region',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '省份',
        dataIndex: 'zone',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '是否异网',
        dataIndex: 'isNotlocalIsp',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '异网运营商',
        dataIndex: 'bizIsp',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '跨网计费方式',
        dataIndex: 'bizIspMode',
        width: 100,
        customCell: cellContent,
      },
      {
        title: '计费方式',
        dataIndex: 'mode',
        width: 300,
        customCell: cellContent,
      },
      {
        title: '备注',
        dataIndex: 'describe',
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
        mode: 'multiple',
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
      field: 'localIsp',
      label: '本网运营商',
      component: 'Select',
      componentProps: {
        options: [
          { label: '电信', value: '电信' },
          { label: '移动', value: '移动' },
          { label: '联通', value: '联通' },
        ],
      },
      colProps: { span: 6 },
    },
    {
      field: 'diffNet',
      label: '是否异网',
      component: 'Select',
      componentProps: {
        options: [
          { label: '否', value: 0 },
          { label: '是', value: 1 },
        ],
      },
      colProps: { span: 6 },
    },
    {
      label: '大区',
      field: 'regionID',
      component: 'ApiSelect',
      componentProps: {
        api: getRegion,
        params: { pageSize: 5000, pageIndex: 1 },
        resultField: 'items',
        labelField: 'name',
        valueField: 'id',
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
      label: '省份',
      field: 'zoneID',
      component: 'ApiSelect',
      componentProps: {
        api: getZone,
        params: { pageSize: 5000, pageIndex: 1 },
        resultField: 'items',
        labelField: 'name',
        valueField: 'id',
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
      colProps: { span: 6 },
    },
  ];

  const [registerTable, { getDataSource }] = useTable({
    title: '业务单价列表',
    api: getAllRecord,
    columns: getMergeHeaderColumns(),
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
    },
    useSearchForm: true,
    bordered: true,
    showTableSetting: true,
    showIndexColumn: false,
  });

  function handleQueryBiz(id: Recordable) {
    bizData.items = [];
    getBiz({ id: id })
      .then(function (resp) {
        for (var i = 0; i < resp.length; i++) {
          bizData.items.push(resp[i]);
        }
        open.value = true;
      })
      .catch(function (err) {
        message.error(err);
      });
  }
</script>
