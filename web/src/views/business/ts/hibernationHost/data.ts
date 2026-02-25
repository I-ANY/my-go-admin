import { BasicColumn, FormSchema } from '@/components/Table';
import { RangePickPresetsExact } from '@/utils/common';

// const uts_status_map = [
//   {
//     label: '正常',
//     value: 'false',
//   },
//   {
//     label: '异常',
//     value: 'true',
//   },
// ];

export const searchFormSchema = (
  onTimePikerOpen,
  // onTimeRangeChange,
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
      // onChange: onTimeRangeChange,
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
];

export const columns: BasicColumn[] = [
  {
    title: '休眠时间',
    dataIndex: 'hibernationTime',
    width: 200,
    resizable: true,
  },
  {
    title: '休眠主机名',
    dataIndex: 'hostname',
    width: 200,
    resizable: true,
  },
];
