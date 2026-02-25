import type { BasicColumn, FormSchema } from '@/components/Table';
import { GetBusinessGroupOptions, GetCategoryList } from '@/api/business/biz';

import dayjs from 'dayjs';
import { getOwners } from '@/api/business/a';
import { h } from 'vue';

export const columns: BasicColumn[] = [
  {
    title: '业务组',
    dataIndex: 'businessGroup',
    key: 'businessGroup',
    width: 120,
    customCell: (record: any) => {
      if (!record) return {};
      const rowSpan = record.businessGroupRowSpan ?? 1;
      if (rowSpan === 0) {
        return { rowSpan: 0 };
      }
      return { rowSpan };
    },
  },
  {
    title: '业务组扣分',
    dataIndex: 'businessGroupScore',
    key: 'businessGroupScore',
    width: 120,
    sorter: true,
    customCell: (record: any) => {
      if (!record) return {};
      const rowSpan = record.businessGroupRowSpan ?? 1;
      if (rowSpan === 0) {
        return { rowSpan: 0 };
      }
      return { rowSpan };
    },
    customRender: ({ record }) => {
      const score = record?.businessGroupScore ? record.businessGroupScore.toFixed(1) : '-';
      if (score === '-') return score;
      return h('span', { style: { color: '#ff4d4f' } }, score);
    },
  },
  {
    title: '考核业务',
    dataIndex: 'assessmentBiz',
    key: 'assessmentBiz',
    width: 150,
    sorter: true,
    customCell: (record: any) => {
      if (!record) return {};
      const rowSpan = record.businessRowSpan ?? 1;
      if (rowSpan === 0) {
        return { rowSpan: 0 };
      }
      return { rowSpan };
    },
  },
  {
    title: '业务扣分',
    dataIndex: 'businessScore',
    key: 'businessScore',
    width: 100,
    sorter: true,
    customCell: (record: any) => {
      if (!record) return {};
      const rowSpan = record.businessRowSpan ?? 1;
      if (rowSpan === 0) {
        return { rowSpan: 0 };
      }
      return { rowSpan };
    },
    customRender: ({ record }) => {
      const score = record?.businessScore ? record.businessScore.toFixed(1) : '-';
      if (score === '-') return score;
      return h('span', { style: { color: '#fa8c16' } }, score);
    },
  },
  {
    title: '节点编号',
    dataIndex: 'roomNo',
    key: 'roomNo',
    width: 150,
    sorter: true,
  },
  {
    title: '节点扣分',
    dataIndex: 'nodeScore',
    key: 'nodeScore',
    width: 100,
    sorter: true,
    customRender: ({ record }) => {
      const score = record?.nodeScore ? record.nodeScore.toFixed(1) : '-';
      if (score === '-') return score;
      return h('span', { style: { color: '#52c41a' } }, score);
    },
  },
];

export const searchFormSchema: FormSchema[] = [
  {
    label: '日期范围',
    field: 'dateRange',
    component: 'RangePicker',
    colProps: { span: 8 },
    componentProps: {
      format: 'YYYY-MM-DD',
      valueFormat: 'YYYY-MM-DD',
      defaultValue: [dayjs().subtract(7, 'day').format('YYYY-MM-DD'), dayjs().format('YYYY-MM-DD')],
    },
  },
  {
    field: 'businessGroup',
    label: '业务组',
    component: 'ApiSelect',
    componentProps: {
      placeholder: '请选择业务组',
      options: [],
      showSearch: true,
      api: async () => {
        const data = await GetBusinessGroupOptions();
        return data.map((item: any) => ({
          label: item.name,
          value: item.name,
        }));
      },
    },
    colProps: { span: 8 },
  },
  {
    field: 'assessmentBizs',
    label: '考核业务',
    component: 'ApiSelect',
    componentProps: {
      options: [],
      showSearch: true,
      mode: 'multiple',
      placeholder: '请选择考核业务',
      api: async () => {
        const data = await GetCategoryList({ pageSize: 1000, pageIndex: 1 });
        return data.items.map((item: any) => ({
          label: item.name,
          value: item.name,
        }));
      },
    },
    colProps: { span: 8 },
  },
  {
    field: 'roomNo',
    label: '节点编号',
    component: 'ApiSelect',
    componentProps: {
      options: [],
      showSearch: true,
      placeholder: '请选择节点编号',
      api: async () => {
        const data = await getOwners();
        return data.map((item: any) => ({
          label: item.name,
          value: item.id,
        }));
      },
    },
    colProps: { span: 8 },
  },
];

// 历史方案表格列配置
export const historyPlanColumns: BasicColumn[] = [
  {
    title: '扣分日期',
    dataIndex: 'date',
    key: 'date',
    width: 120,
    customCell: (record: any) => {
      if (!record) return {};
      const rowSpan = record.rowSpan ?? 1;
      if (rowSpan === 0) {
        return { rowSpan: 0 };
      }
      return { rowSpan };
    },
  },
  {
    title: '考核业务',
    dataIndex: 'assessmentBiz',
    key: 'assessmentBiz',
    width: 120,
    customCell: (record: any) => {
      if (!record) return {};
      const rowSpan = record.rowSpan ?? 1;
      if (rowSpan === 0) {
        return { rowSpan: 0 };
      }
      return { rowSpan };
    },
  },
  {
    title: '节点编号',
    dataIndex: 'roomNo',
    key: 'roomNo',
    width: 150,
    customCell: (record: any) => {
      if (!record) return {};
      const rowSpan = record.rowSpan ?? 1;
      if (rowSpan === 0) {
        return { rowSpan: 0 };
      }
      return { rowSpan };
    },
  },
  {
    title: '考核日期',
    dataIndex: 'assessmentDate',
    key: 'assessmentDate',
    width: 150,
  },
  {
    title: '方案内容',
    dataIndex: 'planContent',
    key: 'planContent',
    ellipsis: true,
    align: 'left',
    customRender: ({ record }) => {
      if (!record.planContent || record.planContent === '-') return '-';
      // 如果包含换行符，显示为多行文本
      if (record.planContent.includes('\n')) {
        return h(
          'div',
          {
            style: { whiteSpace: 'pre-line' },
          },
          record.planContent,
        );
      }
      return record.planContent;
    },
  },
];

// 历史方案搜索表单配置
export const historyPlanSearchFormSchema: FormSchema[] = [
  {
    label: '日期范围',
    field: 'dateRange',
    component: 'RangePicker',
    colProps: { span: 8 },
    componentProps: {
      format: 'YYYY-MM-DD',
      valueFormat: 'YYYY-MM-DD',
      placeholder: ['开始日期', '结束日期'],
    },
  },
  {
    field: 'assessmentBizs',
    label: '考核业务',
    component: 'ApiSelect',
    componentProps: {
      options: [],
      showSearch: true,
      mode: 'multiple',
      placeholder: '请选择考核业务',
      api: async () => {
        const data = await GetCategoryList({ pageSize: 1000, pageIndex: 1 });
        return data.items.map((item: any) => ({
          label: item.name,
          value: item.name,
        }));
      },
    },
    colProps: { span: 8 },
  },
  {
    field: 'roomNo',
    label: '节点编号',
    component: 'Input',
    componentProps: {
      placeholder: '请输入节点编号',
      allowClear: true,
    },
    colProps: { span: 8 },
  },
];
