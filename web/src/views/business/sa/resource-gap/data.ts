import { FormSchema } from '@/components/Table';
import { h } from 'vue';
import { Tag } from 'ant-design-vue';
import { SAEnum } from '@/enums/dictTypeCode';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { formatToDateTime } from '@/utils/dateUtil';

export const saResourceGapMap = getDictDataMapFromDict(SAEnum.SA_RESOURCE_GAP_STATUS);

export const resourceGapSearchSchema: FormSchema[] = [
  {
    field: 'gap_type',
    label: '资源缺口类型',
    component: 'Select',
    colProps: { span: 8 },
    componentProps: {
      allowClear: true,
      options: [],
    },
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      allowClear: true,
      options: getSelectOptionsFromDict(SAEnum.SA_RESOURCE_GAP_STATUS),
    },
  },
  {
    field: 'province',
    label: '省份',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      showSearch: true,
      allowClear: true,
    },
  },
  {
    field: 'isp',
    label: '运行商',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      allowClear: true,
      options: [
        {
          label: '电信',
          value: '电信',
        },
        {
          label: '移动',
          value: '移动',
        },
        {
          label: '联通',
          value: '联通',
        },
        {
          label: '联通转电信',
          value: '联通转电信',
        },
        {
          label: '联通转移动',
          value: '联通转移动',
        },
        {
          label: '移动转电信',
          value: '移动转电信',
        },
        {
          label: '移动转联通',
          value: '移动转联通',
        },
        {
          label: '电信转移动',
          value: '电信转移动',
        },
        {
          label: '电信转联通',
          value: '电信转联通',
        },
      ],
    },
  },
];

export const resourceGapTableColumns = [
  {
    title: '资源缺口类型',
    dataIndex: 'gap_type',
    width: 250,
  },
  {
    title: '省份',
    dataIndex: 'province',
    width: 120,
  },
  {
    title: '运行商',
    dataIndex: 'isp',
    width: 80,
  },
  {
    title: '带宽(最新值)',
    dataIndex: 'bandwidth_gap',
    width: 80,
  },
  {
    title: '带宽(上个值)',
    dataIndex: 'previous_bandwidth_gap',
    width: 80,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 80,
    customRender: ({ record }) => {
      const status = record.status;
      const color = saResourceGapMap[status].color || 'default';
      const text = saResourceGapMap[status].dictLabel;
      return h(Tag, { color }, text);
    },
  },
  {
    title: '文件时间',
    dataIndex: 'file_time',
    width: 100,
  },
  {
    title: '文件时间(上次修改)',
    dataIndex: 'previous_file_time',
    width: 100,
  },
  {
    title: '更新时间',
    dataIndex: 'updatedAt',
    width: 120,
    customRender: ({ record }) => {
      return formatToDateTime(record.updatedAt);
    },
  },
];
