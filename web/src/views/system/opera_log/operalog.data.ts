import { BasicColumn, FormSchema } from '@/components/Table';
import { SysEnum } from '@/enums/dictTypeCode';
import { RangePickPresetsExact } from '@/utils/common';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';

export const requestMethodMap = getDictDataMapFromDict(SysEnum.REQUEST_METHOD);
export const handleSourceMap = getDictDataMapFromDict(SysEnum.HANDLE_SOURCE);

export const handleSourceOptions = getSelectOptionsFromDict(SysEnum.HANDLE_SOURCE);
export const getOperaLogSearchFormSchema = function (
  onSearchUser: any,
  onTimePikerOpen,
): FormSchema[] {
  return [
    {
      field: '[requestTimeRangeStart, requestTimeRangeEnd]',
      // field: 'time',
      label: '操作时间',
      component: 'RangePicker',
      componentProps: {
        format: 'YYYY-MM-DD HH:mm:ss',
        showTime: { format: 'HH:mm:ss' },
        placeholder: ['开始时间', '结束时间'],
        presets: RangePickPresetsExact(),
        allowClear: false,
        onOpenChange: onTimePikerOpen,
      },
      required: true,
      colProps: { span: 6 },
    },
    {
      field: 'handleSource',
      label: '操作来源',
      component: 'Select',
      colProps: { span: 6 },
      componentProps: {
        options: handleSourceOptions,
        allowClear: false,
      },
      required: true,
      defaultValue: handleSourceOptions[0].value as any,
    },
    {
      field: 'uri',
      label: 'URI',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'userId',
      label: '操作人',
      component: 'Select',
      componentProps: {
        placeholder: '输入用户信息进行搜索',
        allowClear: true,
        showSearch: true,
        loading: false,
        onSearch: onSearchUser,
        notFoundContent: '未找到用户信息',
        filterOption: (input: string, option: any) => {
          return (
            option.nickName?.toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
            option.username?.toLowerCase().indexOf(input.toLowerCase()) >= 0
          );
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'httpCode',
      label: 'HTTP状态码',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'bizCode',
      label: '业务状态码',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'requestMethod',
      label: '请求方法',
      component: 'Select',
      colProps: { span: 6 },
      componentProps: {
        options: getSelectOptionsFromDict(SysEnum.REQUEST_METHOD),
        allowClear: true,
      },
    },
    {
      field: 'handler',
      label: 'handler',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'clientIp',
      label: '客户端IP',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'requestId',
      label: 'RequestId',
      component: 'Input',
      colProps: { span: 6 },
    },
  ];
};

export const operaLogColumns: BasicColumn[] = [
  {
    title: 'URI',
    dataIndex: 'uri',
    ellipsis: true,
    width: 280,
    resizable: true,
  },
  {
    title: '请求方法',
    dataIndex: 'requestMethod',
    width: 80,
    resizable: true,
  },
  {
    title: 'HTTP状态码',
    dataIndex: 'httpCode',
    width: 90,
    resizable: true,
  },
  {
    title: '业务状态码',
    dataIndex: 'bizCode',
    width: 90,
    resizable: true,
  },
  {
    title: '请求耗时(ms)',
    dataIndex: 'latencyTime',
    width: 100,
    resizable: true,
  },
  {
    title: '客户端IP',
    dataIndex: 'clientIp',
    width: 150,
    resizable: true,
  },
  {
    title: '操作时间',
    dataIndex: 'requestTime',
    width: 160,
    resizable: true,
    // customRender: ({ record }) => {
    //   return formatToDateTime(record.createdAt);
    // },
  },
  {
    title: '操作来源',
    dataIndex: 'handleSource',
    width: 120,
  },
  {
    title: '操作人',
    dataIndex: 'user',
    width: 120,
    customRender: ({ record }) => {
      return record.userInfo?.nickName;
    },
  },
  {
    title: 'handler',
    dataIndex: 'handler',
    width: 200,
    resizable: true,
  },
  {
    title: '接口',
    dataIndex: 'api',
    ellipsis: true,
    width: 160,
    resizable: true,
  },
  {
    title: 'UA',
    dataIndex: 'userAgent',
    width: 120,
    resizable: true,
  },
  {
    title: 'RequestId',
    dataIndex: 'requestId',
    width: 200,
    resizable: true,
  },
];
