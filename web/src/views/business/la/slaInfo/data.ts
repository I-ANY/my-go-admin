import { BasicColumn, FormSchema } from '@/components/Table';
import { getSelectOptionsFromDict } from '@/utils/dict';
import { CommonEnum } from '@/enums/dictTypeCode';
import { RangePickPresetsExact } from '@/utils/common';
import dayjs from 'dayjs';

export const mergeIspOptions = getSelectOptionsFromDict(CommonEnum.MERGE_ISP);

export const searchFormSchema = (onTimePikerOpen): FormSchema[] => [
  {
    field: '[timeBegin, timeEnd]',
    label: '开始时间',
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
      disabledDate: (date) => {
        const todayEnd = dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss');
        return date.isAfter(todayEnd);
      },
    },
    colProps: { span: 6 },
    required: true,
  },
  {
    field: 'provinces',
    label: '省份',
    component: 'Select',
    componentProps: {
      allowClear: true,
      mode: 'multiple',
      maxTagCount: 3,
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'isp',
    label: '运营商',
    component: 'Select',
    componentProps: {
      options: mergeIspOptions,
      allowClear: true,
      // mode: 'multiple',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'hostname',
    label: '主机名',
    component: 'InputTextArea',
    componentProps: {
      allowClear: true,
      placeholder: '输入主机名进行搜索',
      rows: 3,
    },
    colProps: { span: 6 },
  },
  {
    field: 'node',
    label: '机房',
    component: 'Input',
    componentProps: {
      allowClear: true,
      placeholder: '输入机房搜索',
    },
    colProps: { span: 6 },
  },
  {
    field: 'deviceId',
    label: '设备ID',
    component: 'Input',
    componentProps: {
      allowClear: true,
      placeholder: '输入设备ID搜索',
    },
    colProps: { span: 6 },
  },
];

export const columns: BasicColumn[] = [
  {
    title: '省份',
    dataIndex: 'province',
    width: 100,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'isp',
    width: 80,
    resizable: true,
  },
  {
    title: '设备总数',
    dataIndex: 'deviceCount',
    width: 100,
    resizable: true,
    sorter: true,
  },
  {
    title: '正常设备数',
    dataIndex: 'normalCount',
    width: 100,
    resizable: true,
    sorter: true,
  },
  {
    title: '异常设备数',
    dataIndex: 'abnormalCount',
    width: 100,
    resizable: true,
    sorter: true,
  },
  {
    title: '正常占比',
    dataIndex: 'normalRate',
    width: 100,
    resizable: true,
    sorter: true,
  },
  {
    title: '时间',
    dataIndex: 'time',
    width: 160,
    resizable: true,
    sorter: true,
  },
  // {
  //   title: '创建时间',
  //   dataIndex: 'createdAt',
  //   width: 160,
  //   resizable: true,
  // },
];
export const slaAbnormalHostColumns = function (): BasicColumn[] {
  return [
    {
      title: '主机名',
      dataIndex: 'hostname',
      width: 240,
      resizable: true,
    },
    {
      title: '机房',
      dataIndex: 'node',
      width: 200,
      resizable: true,
    },
    {
      title: '设备ID',
      dataIndex: 'deviceId',
      width: 240,
      resizable: true,
    },
  ];
};

export const slaAbnormalHostSearchFormSchema = function (): FormSchema[] {
  return [
    {
      field: 'hostname',
      label: '主机名',
      component: 'Input',
      componentProps: {
        allowClear: true,
        placeholder: '输入主机名搜索',
      },
      colProps: { span: 8 },
    },
    {
      field: 'node',
      label: '机房',
      component: 'Input',
      componentProps: {
        allowClear: true,
        placeholder: '输入机房搜索',
      },
      colProps: { span: 8 },
    },
    {
      field: 'deviceId',
      label: '设备ID',
      component: 'Input',
      componentProps: {
        allowClear: true,
        placeholder: '输入设备ID搜索',
      },
      colProps: { span: 8 },
    },
  ];
};
