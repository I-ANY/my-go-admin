import { FormSchema } from '@/components/Form';
import { BasicColumn } from '@/components/Table';
import { formatToDateTime } from '@/utils/dateUtil';
import { RangePickPresetsExact } from '@/utils/common';

export const trafficExternalSearchFormSchema = (): FormSchema[] => [
  {
    field: '[metricsTimeBegin, metricsTimeEnd]',
    label: '时间范围',
    component: 'RangePicker',
    colProps: { span: 8 },
    componentProps: {
      allowClear: false,
      format: 'YYYY-MM-DD HH:mm:ss',
      showTime: { format: 'HH:mm:ss' },
      placeholder: ['开始时间', '结束时间'],
      style: {
        width: '100%',
      },
      presets: RangePickPresetsExact(),
    },
    required: true,
  },
  {
    field: 'hostnames',
    label: '主机名',
    component: 'InputTextArea',
    colProps: { span: 8 },
    componentProps: {
      allowClear: true,
      placeholder: '每行一个主机名',
      rows: 2,
    },
  },
  {
    field: 'Businesses',
    label: '业务小类',
    component: 'InputTextArea',
    colProps: { span: 8 },
    componentProps: {
      allowClear: true,
      placeholder: '每行一个业务小类名',
      rows: 2,
    },
  },
  {
    field: 'sn',
    label: 'SN',
    component: 'Input',
    colProps: { span: 4 },
    componentProps: {
      allowClear: true,
      placeholder: '',
    },
  },
  {
    field: 'isp',
    label: '运营商',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      allowClear: true,
      options: [
        { label: '电信', value: '电信' },
        { label: '联通', value: '联通' },
        { label: '移动', value: '移动' },
        { label: '广电', value: '广电' },
      ],
      placeholder: '请选择运营商',
    },
  },
];

export const trafficExternalColumns: BasicColumn[] = [
  {
    title: '时间',
    dataIndex: 'standard_time',
    width: 180,
    resizable: true,
    align: 'left',
    customRender: ({ record }) => {
      if (!record?.standard_time) {
        return '--';
      }
      return formatToDateTime(record.standard_time);
    },
  },
  {
    title: '业务小类',
    dataIndex: 'business',
    width: 180,
    resizable: true,
    align: 'left',
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
    align: 'left',
  },
  {
    title: 'SN',
    dataIndex: 'sn',
    width: 200,
    resizable: true,
    align: 'left',
  },
  {
    title: '运营商',
    dataIndex: 'isp',
    width: 100,
    resizable: true,
  },
  {
    title: '上行流量(bps)',
    dataIndex: 'speed',
    width: 140,
    resizable: true,
    align: 'right',
  },
  {
    title: '下行流量(bps)',
    dataIndex: 'down',
    width: 140,
    resizable: true,
    align: 'right',
  },
];
