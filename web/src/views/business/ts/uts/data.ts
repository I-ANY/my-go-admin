import { BasicColumn, FormSchema } from '@/components/Table';
import { RangePickPresetsExact } from '@/utils/common';

const uts_status_map = [
  {
    label: '正常',
    value: 'false',
  },
  {
    label: '异常',
    value: 'true',
  },
];

export const searchFormSchema = (
  onTimePikerOpen,
  onTimeRangeChange,
  disabledDate,
): FormSchema[] => [
  {
    field: '[metricsTimeBegin, metricsTimeEnd]',
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
      onChange: onTimeRangeChange,
      disabledDate: disabledDate,
    },
    required: true,
    colProps: { span: 6 },
  },
  {
    field: 'hostnames',
    label: '主机名',
    component: 'InputTextArea',
    componentProps: {
      allowClear: true,
      placeholder: '输入主机名进行搜索',
      rows: 3,
    },
    colProps: { span: 6 },
  },
  // {
  //     field: 'uts_vm_hostnames',
  //     label: 'kvm主机名称',
  //     component: 'InputTextArea',
  //     componentProps: {
  //       allowClear: true,
  //       placeholder: '输入主机名进行搜索',
  //       rows: 3,
  //     },
  //     colProps: { span: 8 },
  // },
  {
    field: 'exceeded95',
    label: '流量超出95线',
    component: 'Select',
    componentProps: {
      options: uts_status_map,
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'exceeded_plan_bw',
    label: '流量超出规划带宽',
    component: 'Select',
    componentProps: {
      options: uts_status_map,
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'exceeded_time',
    label: '超过八点存在流量',
    component: 'Select',
    componentProps: {
      options: uts_status_map,
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'businessStatus',
    label: '宿主机业务状态',
    component: 'Select',
    componentProps: {
      options: uts_status_map,
      allowClear: true,
    },
    colProps: { span: 6 },
  },
];

export const columns: BasicColumn[] = [
  {
    title: '时间',
    dataIndex: 'stime',
    width: 200,
    resizable: true,
  },
  {
    title: '主机名称',
    dataIndex: 'uts_vm_hostname',
    width: 200,
    resizable: true,
  },
  {
    title: '宿主机名称',
    dataIndex: 'hostname',
    width: 120,
    resizable: true,
  },
  {
    title: '流量超出95线',
    dataIndex: 'exceeded95',
    width: 120,
    resizable: true,
  },
  {
    title: '流量超出规划带宽',
    dataIndex: 'exceeded_plan_bw',
    width: 120,
    resizable: true,
  },
  {
    title: '超过八点存在流量',
    dataIndex: 'exceeded_time',
    width: 120,
    resizable: true,
  },
  {
    title: '宿主机业务状态',
    dataIndex: 'businessStatus',
    width: 120,
    resizable: true,
  },
];
