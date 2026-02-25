import type { BasicColumn, FormSchema } from '@/components/Table';
import { getBillingTypes, getOwners, getRoomTypes } from '@/api/business/a';

import { Tag } from 'ant-design-vue';
import { h } from 'vue';
import search from 'ant-design-vue/es/transfer/search';

export const columns: BasicColumn[] = [
  {
    title: '节点编号',
    dataIndex: 'owner',
    width: 150,
    sorter: true,
    fixed: 'left',
  },
  {
    title: '机房类型',
    dataIndex: 'roomType',
    width: 100,
    customRender: ({ record }) => {
      const typeMap = { 1: 'IDC', 2: 'ACDN', 3: 'MCDN' };
      return typeMap[record.roomType] || '-';
    },
  },
  {
    title: '机房归属',
    dataIndex: 'origin',
    width: 100,
    customRender: ({ record }) => {
      const originMap = { 1: '自建', 2: '招募' };
      return originMap[record.origin] || '-';
    },
  },
  {
    title: '运营商',
    dataIndex: 'isp',
    width: 120,
  },
  {
    title: '所在地',
    dataIndex: 'location',
    width: 120,
  },
  {
    title: '主线业务',
    dataIndex: 'mainBiz',
    width: 120,
  },
  {
    title: '统计类型',
    dataIndex: 'reportType',
    width: 100,
    customRender: ({ record }) => {
      const reportType = { 1: '机房总览', 2: '保底业务', 3: '削峰业务' };
      return reportType[record.reportType] || '-';
    },
  },
  {
    title: '是否特批',
    dataIndex: 'isSpecialApproval',
    width: 100,
    customRender: ({ record }) => {
      // 是红色，否绿色
      if (record.isSpecialApproval) {
        return h(Tag, { color: 'red' }, () => '是');
      } else {
        return h(Tag, { color: 'green' }, () => '否');
      }
    },
  },
  {
    title: '考核标准',
    dataIndex: 'assessmentStandard',
    width: 200,
  },
  {
    title: '创建时间',
    dataIndex: 'createdAt',
    width: 180,
    sorter: true,
    format: 'date|YYYY-MM-DD HH:mm:ss',
  },
  {
    title: '更新时间',
    dataIndex: 'updatedAt',
    width: 180,
    sorter: true,
    format: 'date|YYYY-MM-DD HH:mm:ss',
  },
];

export const searchFormSchema: FormSchema[] = [
  {
    field: 'owner',
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
      onChange: (v) => {
        search['owner'] = v;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'roomType',
    label: '机房类型',
    component: 'ApiSelect',
    componentProps: {
      placeholder: '请选择机房类型',
      api: getRoomTypes,
      resultField: '',
      labelField: 'label',
      valueField: 'value',
      immediate: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'billingType',
    label: '计费方式',
    component: 'ApiSelect',
    componentProps: {
      placeholder: '请选择计费方式',
      api: getBillingTypes,
      resultField: '',
      labelField: 'label',
      valueField: 'value',
      immediate: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'isp',
    label: '运营商',
    component: 'Select',
    componentProps: {
      placeholder: '请选择运营商',
      options: [
        { label: '移动', value: '移动' },
        { label: '联通', value: '联通' },
        { label: '电信', value: '电信' },
      ],
    },
    colProps: { span: 6 },
  },
  {
    field: 'location',
    label: '所在地',
    component: 'Input',
    componentProps: {
      placeholder: '请输入所在地',
    },
    colProps: { span: 6 },
  },
  {
    field: 'mainBiz',
    label: '主线业务',
    component: 'Input',
    componentProps: {
      placeholder: '请输入主要业务',
    },
    colProps: { span: 6 },
  },
  {
    field: 'reportType',
    label: '统计类型',
    component: 'Select',
    componentProps: {
      placeholder: '请选择统计类型',
      options: [
        { label: '机房总览', value: 1 },
        { label: '保底业务', value: 2 },
        { label: '削峰业务', value: 3 },
      ],
    },
    colProps: { span: 6 },
  },
  {
    field: 'isSpecialApproval',
    label: '是否特批',
    component: 'RadioButtonGroup',
    componentProps: {
      options: [
        { label: '是', value: true },
        { label: '否', value: false },
      ],
    },
    colProps: { span: 6 },
  },
];

export const formSchema: FormSchema[] = [
  {
    field: 'id',
    label: 'ID',
    component: 'Input',
    show: false,
  },
  {
    field: 'owner',
    label: '节点编号',
    component: 'Input',
    componentProps: {
      placeholder: '节点编号',
      disabled: true,
    },
  },
  {
    field: 'roomType',
    label: '机房类型',
    component: 'Select',
    componentProps: {
      placeholder: '请选择机房类型',
      options: [
        { label: 'IDC', value: 1 },
        { label: 'ACDN', value: 2 },
        { label: 'MCDN', value: 3 },
      ],
      disabled: true,
    },
  },
  {
    field: 'billingType',
    label: '计费方式',
    component: 'Select',
    componentProps: {
      placeholder: '请选择计费方式',
      options: [
        { label: '日95', value: '日95' },
        { label: '月95', value: '月95' },
        { label: '买断', value: '买断' },
      ],
      disabled: true,
    },
  },
  {
    field: 'utilizationRateThreshold',
    label: '日95&晚高峰95利用率(%)',
    component: 'InputNumber',
    required: true,
    componentProps: {
      placeholder: '请输入利用率阈值',
      min: 0,
      max: 100,
      precision: 1,
    },
  },
  {
    field: 'nightPeakPointsThreshold',
    label: '晚高峰峰值点数',
    component: 'InputNumber',
    required: true,
    componentProps: {
      placeholder: '请输入晚高峰峰值点数',
      min: 0,
      precision: 0,
    },
  },
];
