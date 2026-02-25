import type { BasicColumn, FormSchema } from '@/components/Table';
import { getBillingTypes, getOrigins, getOwners, getRoomTypes } from '@/api/business/a';

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
      return typeMap[record.roomType] || record.roomType;
    },
    filters: [
      { text: 'IDC', value: '1' },
      { text: 'ACDN', value: '2' },
      { text: 'MCDN', value: '3' },
    ],
  },
  {
    title: '机房归属',
    dataIndex: 'origin',
    width: 100,
    customRender: ({ record }) => {
      const originMap = { 1: '自建', 2: '招募' };
      return originMap[record.origin] || record.origin;
    },
    filters: [
      { text: '自建', value: '1' },
      { text: '招募', value: '2' },
    ],
  },
  {
    title: '运营商',
    dataIndex: 'Isp',
    width: 120,
  },
  {
    title: '所在地',
    dataIndex: 'location',
    width: 120,
  },
  {
    title: '计费方式',
    dataIndex: 'billingType',
    width: 100,
  },
  {
    title: '总出口带宽',
    dataIndex: 'planBw',
    width: 140,
    sorter: true,
    // 自定义渲染
    customRender: ({ record }) => {
      // 默认bps，转换为 kbps、Mbps、Gbps
      if (record.planBw >= 1000000 * 1000) {
        return `${(record.planBw / (1000000 * 1000)).toFixed(2)} Gbps`;
      } else if (record.planBw > 1000000) {
        return `${(record.planBw / 1000000).toFixed(2)} Mbps`;
      } else if (record.planBw > 1000) {
        return `${(record.planBw / 1000).toFixed(2)} kbps`;
      } else if (record.planBw > 0) {
        return `${record.planBw} bps`;
      } else if (record.planBw === 0) {
        return '0';
      }
      return '-';
    },
  },
  {
    title: '保底带宽',
    dataIndex: 'minBw',
    width: 140,
    sorter: true,
    customRender: ({ record }) => {
      if (record.billingType !== '买断') {
        if (record.minBw >= 1000000 * 1000) {
          return `${(record.minBw / (1000000 * 1000)).toFixed(2)} Gbps`;
        } else if (record.minBw > 1000000) {
          return `${(record.minBw / 1000000).toFixed(2)} Mbps`;
        } else if (record.minBw > 1000) {
          return `${(record.minBw / 1000).toFixed(2)} kbps`;
        } else if (record.minBw > 0) {
          return `${record.minBw} bps`;
        } else if (record.minBw === 0) {
          return '0';
        }
      }
      return '-';
    },
  },
  {
    title: '是否仅异网节点',
    dataIndex: 'isExternalOnly',
    width: 120,
    customRender: ({ record }) => {
      if (record.isExternalOnly === null) return '-';
      return record.isExternalOnly ? '是' : '否';
    },
  },
  {
    title: '是否考核',
    dataIndex: 'isAssessment',
    width: 100,
    customRender: ({ record }) => {
      if (record.isAssessment === null) return '-';
      return record.isAssessment ? '是' : '否';
    },
  },
  {
    title: '跨省占比',
    dataIndex: 'trZoneRatio',
    width: 150,
    customRender: ({ record }) => {
      if (record.trZoneRatio === '-1') return '-';
      return record.trZoneRatio;
    },
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
      mode: 'multiple',
    },
    colProps: { span: 6 },
  },
  {
    field: 'origin',
    label: '机房归属',
    component: 'ApiSelect',
    componentProps: {
      placeholder: '请选择机房归属',
      api: getOrigins,
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
    field: 'isExternal',
    label: '是否仅异网节点',
    component: 'Select',
    componentProps: {
      placeholder: '请选择是否仅异网节点',
      options: [
        { label: '是', value: 'true' },
        { label: '否', value: 'false' },
      ],
    },
    colProps: { span: 6 },
  },
  {
    field: 'isAssessment',
    label: '是否考核',
    component: 'Select',
    componentProps: {
      placeholder: '请选择是否考核',
      options: [
        { label: '是', value: 'true' },
        { label: '否', value: 'false' },
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
    required: true,
    componentProps: {
      placeholder: '请输入节点编号',
      disabled: true,
    },
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
  },
  {
    field: 'isExternalOnly',
    label: '是否仅异网节点',
    component: 'RadioButtonGroup',
    defaultValue: false,
    componentProps: {
      options: [
        { label: '是', value: true },
        { label: '否', value: false },
      ],
    },
  },
  {
    field: 'isAssessment',
    label: '是否考核',
    component: 'RadioButtonGroup',
    componentProps: {
      options: [
        { label: '是', value: true },
        { label: '否', value: false },
      ],
    },
  },
];

export const importFormSchema: FormSchema[] = [
  {
    field: 'file',
    label: '选择文件',
    component: 'Upload',
    componentProps: {
      accept: ['.xlsx', '.xls'],
      helpText: '支持Excel格式文件(.xlsx, .xls)，包含节点编号、计费方式等字段',
    },
  },
];
