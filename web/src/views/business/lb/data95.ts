import { BasicColumn, FormSchema } from '@/components/Table';
import { CommonEnum, LBEnum } from '@/enums/dictTypeCode';
import { RangePickPresetsExact } from '@/utils/common';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { customCeilDivide } from '@/utils/util';

export const ispOptions = getSelectOptionsFromDict(CommonEnum.ISP);
export const ispMap = getDictDataMapFromDict(CommonEnum.ISP);

export const billTypeOptions = getSelectOptionsFromDict(LBEnum.BILL_TYPE);
export const billTypeMap = getDictDataMapFromDict(LBEnum.BILL_TYPE);

export const statusOptions = getSelectOptionsFromDict(LBEnum.DEVICE_STATUS);
export const statusMap = getDictDataMapFromDict(LBEnum.DEVICE_STATUS);

export const searchFormSchema = (onTimePikerOpen): FormSchema[] => [
  {
    field: '[timeBegin, timeEnd]',
    label: '时间',
    component: 'RangePicker',
    componentProps: {
      allowClear: false,
      format: 'YYYY-MM-DD HH:mm:ss',
      showTime: { format: 'HH:mm:ss' },
      placeholder: ['开始时间', '结束时间'],
      style: {
        width: '100%',
      },
      presets: RangePickPresetsExact(),
      onOpenChange: onTimePikerOpen,
    },
    required: true,
    colProps: { span: 8 },
  },
  {
    field: 'isp',
    label: '运营商',
    component: 'Select',
    componentProps: {
      options: ispOptions,
      allowClear: true,
      mode: 'multiple',
    },
    colProps: { span: 8 },
  },
  {
    field: 'billType',
    label: '账单类型',
    component: 'Select',
    componentProps: {
      options: billTypeOptions,
      allowClear: true,
    },
    colProps: { span: 8 },
  },
];

export const columns: BasicColumn[] = [
  {
    title: '95时间',
    dataIndex: 'time',
    width: 200,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'isp',
    width: 120,
    resizable: true,
  },
  {
    title: '账单类型',
    dataIndex: 'billType',
    width: 120,
    resizable: true,
  },
  {
    title: '日期',
    dataIndex: 'version',
    width: 120,
    resizable: true,
  },
  {
    title: '网络速率（Mbps）',
    dataIndex: 'bandwidth',
    width: 220,
    resizable: true,
    customRender: ({ record }) => {
      return customCeilDivide(record.bandwidth);
    },
  },
];
