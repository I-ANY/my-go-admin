import { BasicColumn, FormSchema } from '@/components/Table';
import { RangePickPresetsExact } from '@/utils/common';
import { TingEnum } from '@/enums/dictTypeCode';
import { getDictDataMapFromDict } from '@/utils/dict';

export const currentStatusMap = getDictDataMapFromDict(TingEnum.CURRENT_STATUS);
export const registerStatusMap = getDictDataMapFromDict(TingEnum.REGISTER_STATUS);
export const containerStatusMap = getDictDataMapFromDict(TingEnum.CONTAINER_STATUS);
export const searchFormSchema = (onTimePikerOpen): FormSchema[] => [
  {
    field: '[startTimeBegin, startTimeEnd]',
    label: '开始时间',
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
  {
    field: '[stopTimeBegin, stopTimeEnd]',
    label: '结束时间',
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
    title: 'UUID',
    dataIndex: 'uuid',
    width: 200,
    resizable: true,
  },
  {
    title: '开始时间',
    dataIndex: 'startTime',
    width: 200,
    resizable: true,
  },
  {
    title: '结束时间',
    dataIndex: 'stopTime',
    width: 200,
    resizable: true,
  },
];

export const allDearchForm = (onTimePikerOpen): FormSchema[] => [
  {
    field: 'uuid',
    label: 'UUID',
    component: 'Input',
    colProps: { span: 8 },
  },
  {
    field: '[startTimeBegin, startTimeEnd]',
    label: '开始时间',
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
  {
    field: '[stopTimeBegin, stopTimeEnd]',
    label: '结束时间',
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
