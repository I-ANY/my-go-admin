import { BasicColumn, FormSchema } from '@/components/Table';
import { BEnum, CommonEnum } from '@/enums/dictTypeCode';
import { RangePickPresetsExact } from '@/utils/common';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';

export const deviceTypeMap = getDictDataMapFromDict(BEnum.DEVICE_TYPE);
export const ipTypeMap = getDictDataMapFromDict(BEnum.IP_TYPE);
export const agentStatusMap = getDictDataMapFromDict(BEnum.AGENT_STATUS);
export const businessStatusMap = getDictDataMapFromDict(BEnum.BUSINESS_STATUS);
export const errorTypeMap = getDictDataMapFromDict(BEnum.ERROR_TYPE);
export const localCoverageOnlyMap = getDictDataMapFromDict(BEnum.LOCAL_COVERAGE_ONLY);

export const deviceTypeOptions = getSelectOptionsFromDict(BEnum.DEVICE_TYPE);
export const ipTypeOptions = getSelectOptionsFromDict(BEnum.IP_TYPE);
export const agentStatusOptions = getSelectOptionsFromDict(BEnum.AGENT_STATUS);
export const businessStatusOptions = getSelectOptionsFromDict(BEnum.BUSINESS_STATUS);
export const errorTypeOptions = getSelectOptionsFromDict(BEnum.ERROR_TYPE);
export const ispOptions = getSelectOptionsFromDict(CommonEnum.ISP);
export const localCoverageOnlyOptions = getSelectOptionsFromDict(BEnum.LOCAL_COVERAGE_ONLY);

export const checkResultFormSchema = function (onTimeChange, onTimePikerOpen): FormSchema[] {
  return [
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
        onChange: onTimeChange,
        onOpenChange: onTimePikerOpen,
      },
      required: true,
      colProps: { span: 8 },
    },
    {
      field: 'hostname',
      label: '主机名',
      component: 'Input',
      colProps: { span: 8 },
    },
    {
      field: 'errorType',
      label: '错误类型',
      component: 'Select',
      componentProps: {
        options: errorTypeOptions,
        allowClear: true,
        mode: 'multiple',
      },
      colProps: { span: 8 },
    },
    {
      field: 'agentStatus',
      label: 'xagent状态',
      component: 'Select',
      componentProps: {
        options: agentStatusOptions,
        allowClear: true,
      },
      colProps: { span: 8 },
    },
    {
      field: 'businessStatus',
      label: '业务状态',
      component: 'Select',
      componentProps: {
        options: businessStatusOptions,
        allowClear: true,
      },
      colProps: { span: 8 },
    },
    {
      field: 'business',
      label: '业务名称',
      component: 'Select',
      componentProps: {
        options: [],
        allowClear: true,
        mode: 'multiple',
      },
      colProps: { span: 8 },
    },
    {
      field: 'frankID',
      label: 'frankID',
      component: 'Input',
      colProps: { span: 8 },
    },
  ];
};

export const checkResultColumns: BasicColumn[] = [
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
    fixed: 'left',
  },
  {
    title: '网速(Mbps)',
    dataIndex: 'networkSpeed',
    width: 110,
    resizable: true,
    helpMessage: '数据来源xagent',
  },
  {
    title: 'ECDN_网速(Mbps)',
    dataIndex: 'mfNetworkSpeed',
    width: 150,
    resizable: true,
    helpMessage: '数据来源ECDN平台',
  },
  {
    title: '运营商',
    dataIndex: 'isp',
    width: 100,
    resizable: true,
    helpMessage: '数据来源xagent',
  },
  {
    title: 'ECDN_运营商',
    dataIndex: 'mfIsp',
    width: 130,
    resizable: true,
    helpMessage: '数据来源ECDN平台',
  },

  {
    title: 'xagent状态',
    dataIndex: 'agentStatus',
    width: 100,
  },
  {
    title: '业务状态',
    dataIndex: 'businessStatus',
    width: 100,
  },
  {
    title: '业务名称',
    dataIndex: 'business',
    width: 100,
  },
  {
    title: '容器数量',
    dataIndex: 'containerCount',
    width: 100,
    helpMessage: '以“xcdn-”开头的容器数量',
  },
  {
    title: '错误类型',
    dataIndex: 'errorType',
    width: 120,
    resizable: true,
  },
  {
    title: '上报时间',
    dataIndex: 'reportTime',
    width: 160,
    resizable: true,
  },
  {
    title: '备注',
    dataIndex: 'remark',
    width: 380,
    resizable: true,
  },
  {
    title: 'frankID',
    dataIndex: 'frankID',
    // width: 200,
    // resizable: true,
  },
];

export const trafficFormSchema = function (onTimeChange, onTimePikerOpen): FormSchema[] {
  return [
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
        onChange: onTimeChange,
        presets: RangePickPresetsExact(),
        onOpenChange: onTimePikerOpen,
      },
      required: true,
      colProps: { span: 8 },
    },
    {
      field: 'hostname',
      label: '主机名',
      component: 'Input',
      colProps: { span: 8 },
    },
    {
      field: 'isps',
      label: '运营商',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        options: ispOptions,
        allowClear: true,
      },
      colProps: { span: 8 },
    },
    {
      field: 'ecdnIsps',
      label: 'ECDN运营商',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        options: ispOptions,
        allowClear: true,
      },
      colProps: { span: 8 },
    },
    {
      field: 'mfyIsps',
      label: '明赋_运营商',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        options: ispOptions,
        allowClear: true,
      },
      colProps: { span: 8 },
    },
    {
      field: 'agentStatus',
      label: 'xagent状态',
      component: 'Select',
      componentProps: {
        options: agentStatusOptions,
        allowClear: true,
      },
      colProps: { span: 8 },
    },
    {
      field: 'businessStatus',
      label: '业务状态',
      component: 'Select',
      componentProps: {
        options: businessStatusOptions,
        allowClear: true,
      },
      colProps: { span: 8 },
    },
    {
      field: 'business',
      label: '业务名称',
      component: 'Select',
      componentProps: {
        options: [],
        allowClear: true,
        mode: 'multiple',
      },
      colProps: { span: 8 },
    },
    {
      field: 'deviceType',
      label: '设备类型',
      component: 'Select',
      componentProps: {
        options: deviceTypeOptions,
        allowClear: true,
      },
      colProps: { span: 8 },
    },
    {
      field: 'ipType',
      label: 'IP类型',
      component: 'Select',
      componentProps: {
        options: ipTypeOptions,
        allowClear: true,
      },
      colProps: { span: 8 },
    },
    {
      field: 'mfyIpType',
      label: '明赋_IP类型',
      component: 'Select',
      componentProps: {
        options: ipTypeOptions,
        allowClear: true,
      },
      colProps: { span: 8 },
    },
    {
      field: 'localCoverageOnly',
      label: '省内调度',
      component: 'Select',
      componentProps: {
        options: localCoverageOnlyOptions,
        allowClear: true,
      },
      colProps: { span: 8 },
    },
    {
      field: 'frankID',
      label: 'frankID',
      component: 'Input',
      colProps: { span: 8 },
    },
  ];
};

export const trafficColumns: BasicColumn[] = [
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
    fixed: 'left',
  },
  {
    title: '网速(Mbps)',
    dataIndex: 'networkSpeed',
    width: 110,
    resizable: true,
    helpMessage: '数据来源xagent接口',
  },
  {
    title: 'ECDN_网速(Mbps)',
    dataIndex: 'mfyNetworkSpeed',
    width: 150,
    resizable: true,
    helpMessage: '数据来源ECDN',
  },
  {
    title: '运营商',
    dataIndex: 'isp',
    width: 100,
    helpMessage: '数据来源xagent接口',
  },
  {
    title: 'ECDN运营商',
    dataIndex: 'ecdnIsp',
    width: 130,
    helpMessage: '数据来源ECDN',
  },
  {
    title: '明赋_运营商',
    dataIndex: 'mfyIsp',
    width: 120,
    helpMessage: '数据来源业务配置文件',
  },
  {
    title: 'xagent状态',
    dataIndex: 'agentStatus',
    width: 100,
  },
  {
    title: '业务状态',
    dataIndex: 'businessStatus',
    width: 100,
  },
  {
    title: '业务名称',
    dataIndex: 'business',
    width: 100,
  },
  {
    title: '容器数量',
    dataIndex: 'containerCount',
    width: 100,
    helpMessage: '以“xcdn-”开头的容器数量',
  },
  {
    title: '上报时间',
    dataIndex: 'reportTime',
    width: 160,
    resizable: true,
  },
  {
    title: '省内调度',
    dataIndex: 'localCoverageOnly',
    width: 100,
    resizable: true,
    helpMessage: '数据来源xagent',
  },
  {
    title: '省份',
    dataIndex: 'prov',
    width: 100,
    helpMessage: '数据来源xagent接口',
  },
  {
    title: '明赋_省份',
    dataIndex: 'mfyProv',
    width: 100,
    helpMessage: '数据来源业务配置文件',
  },
  {
    title: 'IP类型',
    dataIndex: 'ipType',
    width: 100,
    helpMessage: '数据来源xagent接口',
  },
  {
    title: '明赋_IP类型',
    dataIndex: 'mfyIpType',
    width: 120,
    helpMessage: '数据来源业务配置文件',
  },
  {
    title: '设备类型',
    dataIndex: 'deviceType',
    width: 100,
  },
  {
    title: 'frankID',
    dataIndex: 'frankID',
    // width: 200,
    // resizable: true,
  },
];
