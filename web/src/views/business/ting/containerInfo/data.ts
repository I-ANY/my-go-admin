import { BasicColumn, FormSchema } from '@/components/Table';
import { TingEnum } from '@/enums/dictTypeCode';
import { RangePickPresetsExact } from '@/utils/common';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';

export const registerStatusMap = getDictDataMapFromDict(TingEnum.REGISTER_STATUS);
export const containerStatusMap = getDictDataMapFromDict(TingEnum.CONTAINER_STATUS);
export const currentStatusMap = getDictDataMapFromDict(TingEnum.CURRENT_STATUS);

export const searchFormSchema = (onTimePikerOpen): FormSchema[] => [
  {
    field: 'uuid',
    label: 'UUID',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'ident',
    label: '主机名',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'currentStatus',
    label: '实时状态',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(TingEnum.CURRENT_STATUS),
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'status',
    label: '容器状态',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(TingEnum.CONTAINER_STATUS),
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'register',
    label: '注册状态',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(TingEnum.REGISTER_STATUS),
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'owner',
    label: '供应商',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: '[metricsTimeBegin, metricsTimeEnd]',
    label: '更新时间',
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

export const columns: BasicColumn[] = [
  {
    title: '容器名称',
    dataIndex: 'container',
    width: 200,
    resizable: true,
  },
  {
    title: 'UUID',
    dataIndex: 'uuid',
    width: 200,
    resizable: true,
  },
  {
    title: '主机名',
    dataIndex: 'ident',
    width: 200,
    resizable: true,
  },
  {
    title: '实时状态',
    dataIndex: 'currentStatus',
    width: 120,
    resizable: true,
  },
  {
    title: '容器状态',
    dataIndex: 'status',
    width: 120,
    resizable: true,
  },
  {
    title: '注册状态',
    dataIndex: 'register',
    width: 120,
    resizable: true,
  },
  {
    title: '更新时间',
    dataIndex: 'metricsTime',
    width: 200,
    resizable: true,
  },
  {
    title: '地区',
    dataIndex: 'location',
    width: 200,
    resizable: true,
  },
  {
    title: '供应商',
    dataIndex: 'owner',
    width: 200,
    resizable: true,
  },
];
