import { BasicColumn, FormSchema } from '@/components/Table';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { LAEnum, CommonEnum } from '@/enums/dictTypeCode';
import { RangePickPresetsExact } from '@/utils/common';
import { commonCustomHeaderCell, customCeilDivide } from '@/utils/util';

export const isDiffIspMap = getDictDataMapFromDict(LAEnum.IS_DIFF_ISP);
export const isDiffIspOptions = getSelectOptionsFromDict(LAEnum.IS_DIFF_ISP);
export const businessStatusMap = getDictDataMapFromDict(LAEnum.BUSINESS_STATUS);
export const businessStatusOptions = getSelectOptionsFromDict(LAEnum.BUSINESS_STATUS);

export const ispOptions = getSelectOptionsFromDict(CommonEnum.ISP);

export const searchFormSchema = (onTimePikerOpen, onAreasChange): FormSchema[] => [
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
  {
    field: 'deviceIds',
    label: '设备ID',
    component: 'InputTextArea',
    componentProps: {
      allowClear: true,
      placeholder: '输入设备ID进行搜索',
      rows: 3,
    },
    colProps: { span: 6 },
  },
  {
    field: 'business',
    label: '业务名称',
    component: 'Select',
    componentProps: {
      allowClear: true,
      mode: 'multiple',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'areas',
    label: '大区',
    component: 'Select',
    componentProps: {
      allowClear: true,
      mode: 'multiple',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
      onChange: onAreasChange,
    },
    colProps: { span: 6 },
  },
  {
    field: 'provinces',
    label: '省份',
    component: 'Select',
    componentProps: {
      allowClear: true,
      mode: 'multiple',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'isps',
    label: '运营商',
    component: 'Select',
    componentProps: {
      options: ispOptions,
      allowClear: true,
      mode: 'multiple',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'diffIsps',
    label: '异网运营商',
    component: 'Select',
    componentProps: {
      options: ispOptions,
      allowClear: true,
      mode: 'multiple',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'businessStatus',
    label: '设备状态',
    component: 'Select',
    componentProps: {
      options: businessStatusOptions,
      allowClear: true,
      // mode: 'multiple',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'isDiffIsp',
    label: '是否异网',
    component: 'Select',
    componentProps: {
      options: isDiffIspOptions,
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  // {
  //   field: 'frankID',
  //   label: 'frankID',
  //   component: 'Input',
  //   componentProps: {
  //     allowClear: true,
  //     placeholder: '请输入frankID进行搜索',
  //   },
  //   colProps: { span: 6 },
  // },
];

export const columns: BasicColumn[] = [
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
    fixed: 'left',
  },
  {
    title: '设备ID',
    dataIndex: 'deviceId',
    width: 180,
    resizable: true,
  },
  {
    title: '业务名称',
    dataIndex: 'business',
    width: 120,
    resizable: true,
  },
  {
    title: '业务状态',
    dataIndex: 'businessStatus',
    width: 80,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'isp',
    width: 80,
    resizable: true,
  },
  {
    title: '是否异网',
    dataIndex: 'isDiffIsp',
    width: 80,
    resizable: true,
  },
  {
    title: '异网运营商',
    dataIndex: 'diffIsp',
    width: 100,
    resizable: true,
  },
  {
    title: '大区',
    dataIndex: 'area',
    width: 80,
    resizable: true,
  },
  {
    title: '省份',
    dataIndex: 'province',
    width: 100,
    resizable: true,
  },
  {
    title: '配置带宽\nMbps',
    dataIndex: 'bandwidth',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      return customCeilDivide(record.bandwidth);
    },
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '日志速率*1000\nMbps',
    dataIndex: 'logNetworkSpeed1000',
    width: 120,
    resizable: true,
    // helpMessage: '从CSS日志获取并*1000',
    customRender: ({ record }) => {
      return customCeilDivide(record.logNetworkSpeed1000);
    },
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '日志速率*1024\nMbps',
    dataIndex: 'logNetworkSpeed1024',
    width: 120,
    resizable: true,
    // helpMessage: '从CSS日志获取并*1024',
    customRender: ({ record }) => {
      return customCeilDivide(record.logNetworkSpeed1024);
    },
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: 'API速率\nMbps',
    dataIndex: 'apiNetworkSpeed',
    // helpMessage: '从客户接口获取',
    width: 120,
    resizable: true,
    customRender: ({ record }) => {
      return customCeilDivide(record.apiNetworkSpeed);
    },
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: 'ECDN速率\nMbps',
    dataIndex: 'ecdnNetworkSpeed',
    width: 120,
    resizable: true,
    // helpMessage: '从ECDN获取',
    customRender: ({ record }) => {
      return customCeilDivide(record.ecdnNetworkSpeed);
    },
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '时间',
    dataIndex: 'time',
    width: 160,
    resizable: true,
    // minWidth: 160,
  },
];
