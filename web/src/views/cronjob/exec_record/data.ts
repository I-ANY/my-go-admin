import { getJobList } from '@/api/cronjob/job';
import { BasicColumn, FormSchema } from '@/components/Table';
import { RangePickPresetsExact } from '@/utils/common';

export const columns: BasicColumn[] = [
  {
    title: '开始时间',
    dataIndex: 'startTime',
    width: 200,
  },
  {
    title: '结束时间',
    dataIndex: 'endTime',
    width: 200,
  },
  {
    title: '执行状态',
    dataIndex: 'runStatus',
    width: 120,
  },
  {
    title: '触发类型',
    dataIndex: 'triggerType',
    width: 120,
  },
  {
    title: '执行耗时',
    dataIndex: 'latencyTime',
    width: 120,
  },
];

export const searchFormSchema = (onTimePikerOpen): FormSchema[] => [
  {
    field: 'runStatus',
    label: '执行状态',
    component: 'Select',
    componentProps: {
      options: [
        { label: '运行中', value: 1 },
        { label: '执行成功', value: 2 },
        { label: '执行失败', value: 3 },
      ],
    },
    colProps: { span: 6 },
  },
  {
    field: 'triggerType',
    label: '触发类型',
    component: 'Select',
    componentProps: {
      options: [
        { label: '定时触发', value: 1 },
        { label: '手动触发', value: 2 },
      ],
    },
    colProps: { span: 6 },
  },
  {
    field: '[startTimeBegin, startTimeEnd]',
    label: '执行时间',
    component: 'RangePicker',
    componentProps: {
      allowClear: true,
      format: 'YYYY-MM-DD HH:mm:ss',
      showTime: { format: 'HH:mm:ss' },
      placeholder: ['开始时间', '结束时间'],
      presets: RangePickPresetsExact(),
      onOpenChange: onTimePikerOpen,
    },
    colProps: { span: 8 },
  },
];

export const allRecordColumns: BasicColumn[] = [
  {
    title: '任务名称',
    dataIndex: 'jobName',
    width: 250,
  },
  {
    title: '开始时间',
    dataIndex: 'startTime',
    width: 200,
  },
  {
    title: '结束时间',
    dataIndex: 'endTime',
    width: 200,
  },
  {
    title: '执行状态',
    dataIndex: 'runStatus',
    width: 120,
  },
  {
    title: '触发类型',
    dataIndex: 'triggerType',
    width: 120,
  },
  {
    title: '执行耗时',
    dataIndex: 'latencyTime',
    width: 120,
  },
];

export const allRecordsearchFormSchema = (onTimePikerOpen): FormSchema[] => [
  {
    field: 'jobId',
    label: '任务名称',
    component: 'ApiSelect',
    componentProps: {
      allowClear: true,
      api: getJobList,
      params: { pageSize: 1000, pageIndex: 1 },
      resultField: 'items',
      labelField: 'jobName',
      valueField: 'id',
      showSearch: true,
      filterOption: (input: string, option: any) => {
        return option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'runStatus',
    label: '执行状态',
    component: 'Select',
    componentProps: {
      showSearch: true,
      options: [
        { label: '运行中', value: 1 },
        { label: '执行成功', value: 2 },
        { label: '执行失败', value: 3 },
      ],
    },
    colProps: { span: 6 },
  },
  {
    field: 'startTime',
    label: '执行时间',
    component: 'RangePicker',
    componentProps: {
      allowClear: true,
      format: 'YYYY-MM-DD HH:mm:ss',
      showTime: { format: 'HH:mm:ss' },
      placeholder: ['开始时间', '结束时间'],
      presets: RangePickPresetsExact(),
      onOpenChange: onTimePikerOpen,
    },
    colProps: { span: 6 },
  },
  {
    field: 'triggerType',
    label: '触发类型',
    component: 'Select',
    componentProps: {
      showSearch: true,
      options: [
        { label: '定时触发', value: 1 },
        { label: '手动触发', value: 2 },
      ],
    },
    colProps: { span: 6 },
  },
];
