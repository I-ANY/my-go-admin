import { FormSchema } from '@/components/Form';
import { KEnum } from '@/enums/dictTypeCode';
import { RangeDataPickPresetsExact } from '@/utils/common';
import { getSelectOptionsFromDict } from '@/utils/dict';

export const hddBoardFormSchema = function (onTimeChange): FormSchema[] {
  return [
    {
      field: '[reportTimeBegin, reportTimeEnd]',
      label: '日期',
      component: 'RangePicker',
      componentProps: {
        allowClear: false,
        format: 'YYYY-MM-DD',
        showTime: false,
        placeholder: ['开始日期', '结束日期'],
        style: {
          width: '100%',
        },
        presets: RangeDataPickPresetsExact(),
        onChange: onTimeChange,
      },
      required: true,
      colProps: { span: 6 },
    },
    {
      field: 'capacity',
      slot: 'capacity',
      label: '容量范围(G)',
      colProps: { span: 6 },
    },
    {
      field: 'networkSpeed',
      slot: 'networkSpeed',
      label: '网速范围(Mbps)',
      colProps: { span: 6 },
    },
    {
      field: 'minimumSize',
      label: '最小容量',
      component: 'InputNumber',
      show: false,
    },
    {
      field: 'maximumSize',
      label: '最大容量',
      component: 'InputNumber',
      show: false,
    },
    {
      field: 'minNetworkSpeed',
      label: '最小网速',
      component: 'InputNumber',
      show: false,
    },
    {
      field: 'maxNetworkSpeed',
      label: '最大网速',
      component: 'InputNumber',
      show: false,
    },
    {
      field: 'processStatus',
      label: '进程状态',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.DEVICE_STATUS),
        allowClear: true,
      },
      colProps: { span: 6 },
      helpMessage: '该字段来源服务器上报',
    },
    {
      field: 'businessMountStatus',
      label: '挂载状态',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.PARTITION_MOUNT_STATUS),
        allowClear: true,
      },
      colProps: { span: 6 },
    },
  ];
};
