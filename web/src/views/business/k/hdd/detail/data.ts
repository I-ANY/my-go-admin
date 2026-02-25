import { BasicColumn, FormSchema } from '@/components/Table';
import { KEnum, CommonEnum } from '@/enums/dictTypeCode';
import { RangePickPresetsExact } from '@/utils/common';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';

export const deviceStatusMap = getDictDataMapFromDict(KEnum.DEVICE_STATUS);
export const ispNameMap = getDictDataMapFromDict(KEnum.ISP_NAME);
export const isFirstMacMap = getDictDataMapFromDict(CommonEnum.YES_NO);
export const isFirstMacOptions = getSelectOptionsFromDict(CommonEnum.YES_NO);
export const businessMountStatusMap = getDictDataMapFromDict(KEnum.PARTITION_MOUNT_STATUS);

export const searchFormSchema = (onTimePikerOpen): FormSchema[] => [
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
  {
    field: 'macAddrs',
    label: 'mac地址',
    component: 'InputTextArea',
    colProps: { span: 6 },
    componentProps: {
      allowClear: true,
      placeholder: '输入MAC进行搜索',
      rows: 3,
    },
  },
  {
    field: '[reportTimeBegin, reportTimeEnd]',
    label: '上报时间',
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
    colProps: { span: 6 },
  },
  {
    field: 'ispName',
    label: '运营商',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(KEnum.ISP_NAME),
      allowClear: true,
    },
    colProps: { span: 6 },
    helpMessage: '该字段来源客户平台',
  },
  {
    field: 'deviceIsp',
    label: '设备运营商',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(KEnum.ISP_NAME),
      allowClear: true,
    },
    colProps: { span: 6 },
    helpMessage: '该字段来源服务器上报',
  },
  {
    field: 'isFirstMac',
    label: '首mac',
    component: 'Select',
    componentProps: {
      options: isFirstMacOptions,
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(KEnum.DEVICE_STATUS),
      allowClear: true,
    },
    colProps: { span: 6 },
    helpMessage: '该字段来源客户平台',
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
    field: 'reportFlowServiceStatus',
    label: '上报服务状态',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(KEnum.DEVICE_STATUS),
      allowClear: true,
    },
    colProps: { span: 6 },
    helpMessage: '流量上报服务状态，该字段来源服务器上报',
  },
  {
    field: 'capacity',
    slot: 'capacity',
    label: '容量范围(G)',
    colProps: { span: 6 },
  },
  {
    field: 'usage',
    slot: 'usage',
    label: '磁盘使用率（%）',
    colProps: { span: 6 },
  },
  {
    field: 'networkSpeed',
    slot: 'networkSpeed',
    label: '网速范围(Mbps)',
    colProps: { span: 6 },
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

export const columns: BasicColumn[] = [
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 200,
    resizable: true,
  },
  {
    title: 'mac地址',
    dataIndex: 'macAddr',
    width: 200,
    resizable: true,
  },
  {
    title: '首mac',
    dataIndex: 'isFirstMac',
    width: 120,
    resizable: true,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 120,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'ispName',
    width: 120,
    resizable: true,
  },
  {
    title: '进程状态',
    dataIndex: 'processStatus',
    width: 120,
    resizable: true,
  },
  {
    title: '流量上报服务状态',
    dataIndex: 'reportFlowServiceStatus',
    width: 120,
    resizable: true,
  },
  {
    title: '设备运营商',
    dataIndex: 'deviceIsp',
    width: 120,
    resizable: true,
  },
  {
    title: '上报时间',
    dataIndex: 'reportTime',
    width: 180,
    resizable: true,
  },
  {
    title: '网络速率（Mbps）',
    dataIndex: 'networkSpeed',
    width: 220,
    resizable: true,
    customRender: ({ record }) => {
      return record.networkSpeed != null ? (record.networkSpeed / 1000 / 1000).toFixed(2) : '';
    },
  },
  {
    title: '挂载详情',
    dataIndex: 'businessMountDetail',
    width: 220,
    resizable: true,
  },
  {
    title: '进程信息',
    dataIndex: 'businessProcessInfo',
    width: 220,
    resizable: true,
  },
];
export const getPartitionInfoColumns = function (): BasicColumn[] {
  return [
    {
      title: '分区名称',
      dataIndex: 'name',
      width: 200,
      resizable: true,
    },
    {
      title: '业务挂载状态',
      dataIndex: 'businessMountStatus',
      width: 200,
      resizable: true,
    },
    {
      title: '挂载点',
      dataIndex: 'mountPath',
      width: 200,
      resizable: true,
    },
    {
      title: '容量（GB）',
      dataIndex: 'capacity',
      width: 120,
      resizable: true,
    },
    {
      title: '磁盘使用率（%）',
      dataIndex: 'usage',
      width: 200,
      resizable: true,
      // customRender: ({ record }) => {
      //   return record.usage || '';
      // },
    },
  ];
};

export const getPartitionInfoSearchForm = function (): FormSchema[] {
  return [
    {
      field: 'name',
      label: '分区名称',
      component: 'Input',
      colProps: { span: 6 },
    },
    // {
    //   field: 'businessMountStatus',
    //   label: '挂载状态',
    //   component: 'Select',
    //   componentProps: {
    //     options: getSelectOptionsFromDict(KEnum.PARTITION_MOUNT_STATUS),
    //     allowClear: true,
    //   },
    //   colProps: { span: 6 },
    // },
    {
      field: 'mountPath',
      label: '挂载点',
      component: 'Input',
      colProps: { span: 6 },
    },
  ];
};
