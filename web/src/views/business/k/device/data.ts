import { BasicColumn, FormSchema } from '@/components/Table';
import { KEnum } from '@/enums/dictTypeCode';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { commonCustomHeaderCell } from '@/utils/util';

// const bizTypes: enumItem[] = [
//   {
//     name: '专线',
//     value: 'specialLine',
//     color: 'green',
//   },
//   {
//     name: '汇聚',
//     value: 'normal',
//     color: 'orange',
//   },
// ];
// const providers: enumItem[] = [
//   {
//     name: '明赋',
//     value: 'mf',
//   },
//   {
//     name: '泓宁',
//     value: 'hn',
//   },
// ];
// const isProvinceScheduling: enumItem[] = [
//   {
//     name: '不限制',
//     value: 0,
//   },
//   {
//     name: '仅本省调度',
//     value: 1,
//   },
// ];

// const isCoverDiffIsp: enumItem[] = [
//   {
//     name: '否',
//     value: 0,
//   },
//   {
//     name: '是',
//     value: 1,
//   },
// ];

export const bizTypesMap = getDictDataMapFromDict(KEnum.BIZ_TYPE);
export const providersMap = getDictDataMapFromDict(KEnum.PRODIVER);
export const isProvinceSchedulingMap = getDictDataMapFromDict(KEnum.SCHEDULE_TYPE);
export const isCoverDiffIspMap = getDictDataMapFromDict(KEnum.IS_COVER_DIFF_ISP);
export const coverDiffIspMap = getDictDataMapFromDict(KEnum.ISP);

export const normalSearchFormSchema: FormSchema[] = [
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
  // {
  //   field: 'biz_type',
  //   label: '业务类型',
  //   component: 'Select',
  //   componentProps: {
  //     options: getSelectOptionsFromDict(KEnum.BIZ_TYPE),
  //     allowClear: true,
  //   },
  //   colProps: { span: 6 },
  // },
  {
    field: 'wan_ip',
    label: '外网IP',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'is_province_scheduling',
    label: '调度控制',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(KEnum.SCHEDULE_TYPE),
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'is_cover_diff_isp',
    label: '是否异网设备',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(KEnum.IS_COVER_DIFF_ISP),
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'op_status_name',
    label: '运营状态',
    component: 'Select',
    colProps: { span: 6 },
    componentProps: {
      options: getSelectOptionsFromDict(KEnum.OP_STATUS),
      allowClear: true,
    },
    defaultValue:
      getSelectOptionsFromDict(KEnum.OP_STATUS).length > 1
        ? getSelectOptionsFromDict(KEnum.OP_STATUS)[0].value
        : null,
  },
  // {
  //   field: 'isp_name',
  //   label: '供应商',
  //   component: 'Input',
  //   colProps: { span: 6 },
  // },
];

export const specialLineSearchFormSchema: FormSchema[] = [
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
  // {
  //   field: 'biz_type',
  //   label: '业务类型',
  //   component: 'Select',
  //   componentProps: {
  //     options: getSelectOptionsFromDict(KEnum.BIZ_TYPE),
  //     allowClear: true,
  //   },
  //   colProps: { span: 6 },
  // },
  {
    field: 'provider',
    label: '供应商',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(KEnum.PRODIVER),
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'wan_ip',
    label: '外网IP',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'is_province_scheduling',
    label: '调度控制',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(KEnum.SCHEDULE_TYPE),
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'is_cover_diff_isp',
    label: '是否异网设备',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(KEnum.IS_COVER_DIFF_ISP),
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'op_status_name',
    label: '运营状态',
    component: 'Select',
    colProps: { span: 6 },
    componentProps: {
      options: getSelectOptionsFromDict(KEnum.OP_STATUS),
      allowClear: true,
    },
    defaultValue:
      getSelectOptionsFromDict(KEnum.OP_STATUS).length > 1
        ? getSelectOptionsFromDict(KEnum.OP_STATUS)[0].value
        : null,
  },
  // {
  //   field: 'province_name',
  //   label: '省份名称',
  //   component: 'Input',
  //   colProps: { span: 6 },
  // },
  // {
  //   field: 'isp_name',
  //   label: '供应商',
  //   component: 'Input',
  //   colProps: { span: 6 },
  // },
];

export const normalColumns: BasicColumn[] = [
  {
    title: 'mac地址',
    dataIndex: 'mac_addr',
    width: 160,
    resizable: true,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
  },
  {
    title: '设备类型',
    dataIndex: 'device_type_name',
    width: 120,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'isp_name',
    width: 120,
    resizable: true,
  },
  {
    title: '省份名称',
    dataIndex: 'province_name',
    width: 120,
    resizable: true,
  },
  {
    title: '入库时间',
    dataIndex: 'pass_time',
    width: 160,
    resizable: true,
  },
  {
    title: '运营状态',
    dataIndex: 'op_status_name',
    width: 120,
    resizable: true,
  },
  {
    title: '上传带宽\n(Mbps)',
    dataIndex: 'link_upload_speed',
    width: 120,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '调度控制',
    dataIndex: 'is_province_scheduling',
    width: 120,
    resizable: true,
  },
  {
    title: '是否异网',
    dataIndex: 'is_cover_diff_isp',
    width: 120,
    resizable: true,
  },
  {
    title: '异网覆盖\n供应商',
    dataIndex: 'cover_diff_isp_id',
    width: 120,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '外网IP',
    dataIndex: 'wan_ip',
    width: 120,
    resizable: true,
  },
];

export const specialLineColumns: BasicColumn[] = [
  {
    title: 'mac地址',
    dataIndex: 'mac_addr',
    width: 160,
    resizable: true,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
  },
  {
    title: '供应商',
    dataIndex: 'provider',
    width: 120,
    resizable: true,
  },
  {
    title: '设备类型',
    dataIndex: 'device_type_name',
    width: 100,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'isp_name',
    width: 80,
    resizable: true,
  },
  {
    title: '省份名称',
    dataIndex: 'province_name',
    width: 120,
    resizable: true,
  },
  {
    title: '入库时间',
    dataIndex: 'pass_time',
    width: 160,
    resizable: true,
  },
  {
    title: '运营状态',
    dataIndex: 'op_status_name',
    width: 80,
    resizable: true,
  },
  {
    title: '上传带宽\n(Mbps)',
    dataIndex: 'link_upload_speed',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '调度控制',
    dataIndex: 'is_province_scheduling',
    width: 120,
    resizable: true,
  },
  {
    title: '是否异网',
    dataIndex: 'is_cover_diff_isp',
    width: 80,
    resizable: true,
  },
  {
    title: '异网覆盖\n供应商',
    dataIndex: 'cover_diff_isp_id',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '外网IP',
    dataIndex: 'wan_ip',
    width: 120,
    resizable: true,
  },
];
