import { BasicColumn, FormSchema } from '@/components/Table';

import { KEnum } from '@/enums/dictTypeCode';
import dayjs from 'dayjs';
import { getSelectOptionsFromDict } from '@/utils/dict';

// 质量异常搜索表单配置
export function qualityAbnormalSearchFormSchema(bizType: string): FormSchema[] {
  return [
    {
      label: '厂商',
      field: 'provider',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.PRODIVER),
        allowClear: true,
        defaultValue: 'mf',
      },
      colProps: { span: 4 },
      ifShow: () => bizType === 'specialLine',
    },
    {
      label: 'MAC地址',
      field: 'mac',
      component: 'Input',
      colProps: { span: 6 },
      componentProps: {
        placeholder: '请输入MAC地址',
        allowClear: true,
      },
    },
    {
      label: '严重程度',
      field: 'severity',
      component: 'Select',
      colProps: { span: 4 },
      componentProps: {
        options: [
          { label: '轻度', value: '轻度' },
          { label: '中度', value: '中度' },
          { label: '严重', value: '严重' },
        ],
        allowClear: true,
        placeholder: '请选择严重程度',
      },
    },
    {
      label: '时间范围',
      field: 'time_range',
      component: 'RangePicker',
      colProps: { span: 10 },
      componentProps: {
        format: 'YYYY-MM-DD HH:mm:ss',
        showTime: true,
        defaultValue: [dayjs().subtract(24, 'hour'), dayjs()],
        placeholder: ['开始时间', '结束时间'],
      },
    },
  ];
}

// 质量异常表格列配置
export const qualityAbnormalColumns: BasicColumn[] = [
  {
    title: 'MAC地址',
    dataIndex: 'mac',
    width: 180,
    resizable: true,
    fixed: 'left',
  },
  {
    title: '运营商',
    dataIndex: 'provider',
    width: 100,
    resizable: true,
    customRender: ({ text }) => {
      const providerMap = {
        mf: '明赋',
        hn: '泓宁',
      };
      return providerMap[text] || text;
    },
  },
  {
    title: '事件类型',
    dataIndex: 'incident_type',
    key: 'incident_type',
    width: 250,
    resizable: true,
    ellipsis: true,
  },
  {
    title: '严重程度',
    dataIndex: 'severity',
    width: 100,
    resizable: true,
    customRender: ({ text }) => text,
    customCell: (record) => {
      const colorMap = {
        轻度: '#52c41a',
        中度: '#faad14',
        严重: '#ff4d4f',
      };
      return {
        style: {
          color: colorMap[record.severity] || '#000',
          fontWeight: 'bold',
        },
      };
    },
  },
  {
    title: '状态描述',
    dataIndex: 'status',
    key: 'status',
    width: 300,
    resizable: true,
    ellipsis: true,
  },
  {
    title: '开始时间',
    dataIndex: 'gnow',
    width: 180,
    resizable: true,
  },
  {
    title: '结束时间',
    dataIndex: 'gend',
    width: 180,
    resizable: true,
  },
];
