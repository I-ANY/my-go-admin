import { BasicColumn, FormSchema } from '@/components/Table';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { LAEnum, CommonEnum } from '@/enums/dictTypeCode';
import { commonCustomHeaderCell, customCeilDivide } from '@/utils/util';

export const isDiffIspMap = getDictDataMapFromDict(LAEnum.IS_DIFF_ISP);
export const isDiffIspOptions = getSelectOptionsFromDict(LAEnum.IS_DIFF_ISP);
export const businessStatusMap = getDictDataMapFromDict(LAEnum.BUSINESS_STATUS);
export const businessStatusOptions = getSelectOptionsFromDict(LAEnum.BUSINESS_STATUS);

export const ispOptions = getSelectOptionsFromDict(CommonEnum.ISP);

export const searchFormSchema = (onAreasChange): FormSchema[] => [
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
    field: 'isDiffIsp',
    label: '是否异网',
    component: 'Select',
    componentProps: {
      options: isDiffIspOptions,
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'node',
    label: '机房',
    component: 'Input',
    componentProps: {
      allowClear: true,
      placeholder: '输入机房进行搜索',
    },
    colProps: { span: 6 },
  },
];

export const columns: BasicColumn[] = [
  {
    title: '设备ID',
    dataIndex: 'deviceId',
    width: 180,
    resizable: true,
    fixed: 'left',
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
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
    title: '机房',
    dataIndex: 'node',
    width: 100,
    resizable: true,
  },
  {
    title: '创建时间',
    dataIndex: 'deviceCreatedTime',
    width: 160,
    resizable: true,
  },
  {
    title: '更新时间',
    dataIndex: 'updatedAt',
    width: 160,
    resizable: true,
  },
];
export const hostHistoryColumns = function (): BasicColumn[] {
  return [
    {
      title: '主机名',
      dataIndex: 'hostname',
      width: 180,
      resizable: true,
    },
    {
      title: '时间',
      dataIndex: 'time',
      width: 160,
      resizable: true,
    },
  ];
};

export const hostHistorySearchFormSchema = function (): FormSchema[] {
  return [
    {
      field: 'hostname',
      label: '主机名',
      component: 'Input',
      componentProps: {
        allowClear: true,
        placeholder: '输入主机名进行搜索',
      },
      colProps: { span: 12 },
    },
  ];
};
