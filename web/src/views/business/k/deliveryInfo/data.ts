import { BasicColumn, FormSchema } from '@/components/Table';
import { CommonEnum, KEnum } from '@/enums/dictTypeCode';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';

import { BizType } from '../data';
import { commonCustomHeaderCell } from '@/utils/util';
import { getDeliveryOptions } from '@/api/business/k';
import { h } from 'vue';

export const roomTypeOptions = getSelectOptionsFromDict(KEnum.NETWORK_TYPE);
export const deliveryInfoMap = getDictDataMapFromDict(KEnum.DELIVERY_STATUS);
export const deliveryTypeMap = getDictDataMapFromDict(KEnum.DELIVERY_TYPE);
export const bizTypesMap = getDictDataMapFromDict(KEnum.BIZ_TYPE);
export const providersMap = getDictDataMapFromDict(KEnum.PRODIVER);
export const isProvinceSchedulingMap = getDictDataMapFromDict(KEnum.SCHEDULE_TYPE);
export const isCoverDiffIspMap = getDictDataMapFromDict(KEnum.IS_COVER_DIFF_ISP);
export const deviceTypeMap = getDictDataMapFromDict(KEnum.DEVICE_TYPE);
export const sameAreaReplaceMap = getDictDataMapFromDict(KEnum.SAME_AREA_REPLACE);
export const hddTransformMap = getDictDataMapFromDict(KEnum.HDD_TRANSFORM);
export const isIndependentDeployMap = getDictDataMapFromDict(CommonEnum.YES_NO);

export const ispOptions = [
  {
    label: '本网',
    value: '0',
  },
  {
    label: '电信',
    value: '1',
  },
  {
    label: '联通',
    value: '2',
  },
  {
    label: '移动',
    value: '3',
  },
];
export enum DiveryStatus {
  DeploySuccess = 1,
  Delivering = 2,
  DeliverySuccess = 3,
  DeliveryFailed = 4,
  DeliveryClosed = 5,
  DifIsping = 6,
  DifIspSuccess = 7,
  DifIspFailed = 8,
  JoinSuccess = 9,
  JoinFailed = 10,
}

export const specialLineFormSchema = function (): FormSchema[] {
  const providerOptions = getSelectOptionsFromDict(KEnum.PRODIVER);
  return [
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
      field: 'owners',
      label: '节点/机房',
      component: 'InputTextArea',
      componentProps: {
        allowClear: true,
        placeholder: '输入主机名进行搜索',
        rows: 3,
      },
      colProps: { span: 6 },
    },
    {
      field: 'demandId',
      label: '需求单号',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'provider',
      label: '供应商',
      component: 'Select',
      componentProps: {
        options: providerOptions,
        allowClear: true,
      },
      colProps: { span: 4 },
      defaultValue: providerOptions?.length > 0 ? providerOptions[0].value : null,
    },
    {
      field: 'deliveryType',
      label: '交付类型',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        options: getSelectOptionsFromDict(KEnum.DELIVERY_TYPE),
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'deliveryStatus',
      label: '交付状态',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        options: getSelectOptionsFromDict(KEnum.DELIVERY_STATUS),
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'isProvinceScheduling',
      label: '调度控制',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.SCHEDULE_TYPE),
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'isCoverDiffIsp',
      label: '是否异网设备',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.IS_COVER_DIFF_ISP),
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'business',
      label: '业务',
      colProps: { span: 4 },
      component: 'ApiSelect',
      componentProps: {
        allowClear: true,
        showSearch: true,
        loading: false,
        api: getDeliveryOptions,
        valueField: 'value',
        labelField: 'label',
        resultField: 'data',
        params: {
          bizType: BizType.specialLine,
          type: 'business',
        },
        filterOption: (input: string, option: any) => {
          return option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0;
        },
      },
    },
  ];
};
export const normalFormSchema = function (): FormSchema[] {
  return [
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
      field: 'owners',
      label: '节点/机房',
      component: 'InputTextArea',
      componentProps: {
        allowClear: true,
        placeholder: '输入主机名进行搜索',
        rows: 3,
      },
      colProps: { span: 6 },
    },
    {
      field: 'demandId',
      label: '需求单号',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'deliveryType',
      label: '交付类型',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        options: getSelectOptionsFromDict(KEnum.DELIVERY_TYPE),
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'deliveryStatus',
      label: '交付状态',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        options: getSelectOptionsFromDict(KEnum.DELIVERY_STATUS),
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'isProvinceScheduling',
      label: '调度控制',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.SCHEDULE_TYPE),
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'isCoverDiffIsp',
      label: '是否异网设备',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.IS_COVER_DIFF_ISP),
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    // {
    //   field: 'provider',
    //   label: '供应商',
    //   component: 'Select',
    //   componentProps: {
    //     options: getSelectOptionsFromDict(KEnum.PRODIVER),
    //     allowClear: true,
    //   },
    //   colProps: { span: 6 },
    // },
    {
      field: 'business',
      label: '业务',
      colProps: { span: 6 },
      component: 'ApiSelect',
      componentProps: {
        allowClear: true,
        showSearch: true,
        loading: false,
        api: getDeliveryOptions,
        valueField: 'value',
        labelField: 'label',
        resultField: 'data',
        params: {
          bizType: BizType.normal,
          type: 'business',
        },
        filterOption: (input: string, option: any) => {
          return option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0;
        },
      },
    },
  ];
};
export const specialLineColumns: BasicColumn[] = [
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
  },
  {
    title: '需求单号',
    dataIndex: 'demandId',
    width: 215,
    resizable: true,
  },
  {
    title: '节点/机房',
    dataIndex: 'owner',
    width: 160,
    resizable: true,
  },
  {
    title: '业务',
    dataIndex: 'business',
    width: 120,
    resizable: true,
  },
  {
    title: '供应商',
    dataIndex: 'provider',
    width: 80,
    resizable: true,
  },
  {
    title: '交付类型',
    dataIndex: 'deliveryType',
    width: 85,
    resizable: true,
  },
  {
    title: `单线带宽\n(Mbps)`,
    dataIndex: 'uploadBw',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  // {
  //   title: '单线下行带宽(Mbps)',
  //   dataIndex: 'downloadBw',
  //   width: 120,
  //   resizable: true,
  // },
  // {
  //   title: '单线下行带宽(Mbps)',
  //   dataIndex: 'downloadBw',
  //   width: 140,
  //   resizable: true,
  // },
  {
    title: '线路数',
    dataIndex: 'bwCount',
    width: 55,
    resizable: true,
  },
  {
    title: '总带宽\n(Mbps)',
    dataIndex: 'totalBw',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '总交付带宽\n(Mbps)',
    dataIndex: 'deliveryBw',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '区域',
    dataIndex: 'areaName',
    width: 60,
    resizable: true,
  },
  {
    title: '省份',
    dataIndex: 'province',
    width: 100,
    resizable: true,
  },
  {
    title: '城市',
    dataIndex: 'city',
    width: 80,
    resizable: true,
  },
  {
    title: '调度控制',
    dataIndex: 'isProvinceScheduling',
    width: 90,
    resizable: true,
  },
  {
    title: '是否异网',
    dataIndex: 'isCoverDiffIsp',
    width: 70,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'isp',
    width: 60,
    resizable: true,
  },
  {
    title: '异网覆盖\n供应商',
    dataIndex: 'coverDiffIsp',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '设备类型',
    dataIndex: 'deviceType',
    width: 80,
    resizable: true,
  },
  {
    title: '网络接入类型',
    dataIndex: 'networkType',
    width: 100,
    resizable: true,
    // customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '外网IP',
    dataIndex: 'ip',
    width: 120,
    resizable: true,
  },
  {
    title: '任务ID',
    dataIndex: 'taskId',
    width: 100,
    resizable: true,
  },
  {
    title: '状态',
    dataIndex: 'deliveryStatus',
    width: 100,
    resizable: true,
    fixed: 'right',
  },
  {
    title: '创建时间',
    dataIndex: 'createdAt',
    width: 200,
    resizable: true,
  },
  {
    title: '备注',
    dataIndex: 'remark',
    // width: 100,
    // resizable: true,
  },
  {
    title: '操作人',
    dataIndex: 'updateName',
    width: 100,
    resizable: true,
  },
];
export const normalColumns: BasicColumn[] = [
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
  },
  {
    title: '需求单号',
    dataIndex: 'demandId',
    width: 215,
    resizable: true,
  },
  {
    title: '节点/机房',
    dataIndex: 'owner',
    width: 160,
    resizable: true,
  },
  {
    title: '业务',
    dataIndex: 'business',
    width: 120,
    resizable: true,
  },
  {
    title: '交付类型',
    dataIndex: 'deliveryType',
    width: 85,
    resizable: true,
  },
  {
    title: `单线带宽\n(Mbps)`,
    dataIndex: 'uploadBw',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  // {
  //   title: '单线下行带宽(Mbps)',
  //   dataIndex: 'downloadBw',
  //   width: 120,
  //   resizable: true,
  // },
  // {
  //   title: '单线下行带宽(Mbps)',
  //   dataIndex: 'downloadBw',
  //   width: 140,
  //   resizable: true,
  // },
  {
    title: '线路数',
    dataIndex: 'bwCount',
    width: 55,
    resizable: true,
  },
  {
    title: '总带宽\n(Mbps)',
    dataIndex: 'totalBw',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '总交付带宽\n(Mbps)',
    dataIndex: 'deliveryBw',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '区域',
    dataIndex: 'areaName',
    width: 60,
    resizable: true,
  },
  {
    title: '省份',
    dataIndex: 'province',
    width: 100,
    resizable: true,
  },
  {
    title: '城市',
    dataIndex: 'city',
    width: 80,
    resizable: true,
  },
  {
    title: '调度控制',
    dataIndex: 'isProvinceScheduling',
    width: 90,
    resizable: true,
  },
  {
    title: '是否异网',
    dataIndex: 'isCoverDiffIsp',
    width: 70,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'isp',
    width: 60,
    resizable: true,
  },
  {
    title: '异网覆盖\n供应商',
    dataIndex: 'coverDiffIsp',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '设备类型',
    dataIndex: 'deviceType',
    width: 80,
    resizable: true,
  },
  // {
  //   title: '供应商',
  //   dataIndex: 'provider',
  //   width: 80,
  //   resizable: true,
  // },
  {
    title: '外网IP',
    dataIndex: 'ip',
    width: 120,
    resizable: true,
  },
  {
    title: '网络接入类型',
    dataIndex: 'networkType',
    width: 100,
    resizable: true,
  },
  {
    title: '创建时间',
    dataIndex: 'createdAt',
    width: 160,
    resizable: true,
  },
  {
    title: '任务ID',
    dataIndex: 'taskId',
    width: 100,
    resizable: true,
  },
  {
    title: '状态',
    dataIndex: 'deliveryStatus',
    width: 100,
    resizable: true,
    fixed: 'right',
  },
  {
    title: '备注',
    dataIndex: 'remark',
    // width: 100,
    // resizable: true,
  },
  {
    title: '操作人',
    dataIndex: 'updateName',
    width: 100,
    resizable: true,
  },
];
export const specialLineDemandColumns = function (): BasicColumn[] {
  return [
    {
      title: '需求单号',
      dataIndex: 'demandId',
      width: 215,
      resizable: true,
      fixed: 'left',
    },
    // {
    //   title: '业务类型',
    //   dataIndex: 'bizType',
    //   width: 80,
    //   resizable: true,
    // },
    {
      title: '剩余可绑定\n带宽(Mbps)',
      dataIndex: 'canBoundedBw',
      width: 90,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '缺口带宽\n(Mbps)',
      dataIndex: 'gapBw',
      width: 90,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '是否独立\n部署',
      dataIndex: 'isIndependentDeploy',
      width: 80,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '需求带宽\n(Mbps)',
      dataIndex: 'demandBw',
      width: 90,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '已交付带宽\n(Mbps)',
      dataIndex: 'deliveredBw',
      width: 90,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '供应商',
      dataIndex: 'provider',
      width: 70,
      resizable: true,
    },
    {
      title: '设备类型',
      dataIndex: 'devId',
      width: 80,
      resizable: true,
    },
    {
      title: '区域名称',
      dataIndex: 'areaName',
      width: 80,
      resizable: true,
    },
    {
      title: '省份',
      dataIndex: 'provinceName',
      width: 100,
      resizable: true,
    },
    {
      title: '运营商',
      dataIndex: 'ispName',
      width: 80,
      resizable: true,
    },
    {
      title: '是否跨省',
      dataIndex: 'isProvinceScheduling',
      width: 80,
      customRender: ({ text }) => {
        return text === 1 ? '仅本省调度' : '不限制';
      },
    },
    {
      title: '是否异网覆盖',
      dataIndex: 'isCoverDiffIsp',
      width: 80,
      resizable: true,
      customRender: ({ text }) => {
        return text === 1 ? '是' : '否';
      },
    },
    {
      title: '验收中任\n务数',
      dataIndex: 'checkingTaskNum',
      width: 80,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '需求类型\n名称',
      dataIndex: 'demandTypeName',
      width: 80,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '允许同大\n区替换',
      dataIndex: 'isSameAreaReplace',
      width: 80,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    // {
    //   title: '提交验收标识',
    //   dataIndex: 'uploadFlag',
    //   width: 120,
    //   resizable: true,
    // },
    {
      title: 'HDD改造',
      dataIndex: 'hddTransform',
      width: 80,
      resizable: true,
    },
    {
      title: 'IPv4网段',
      dataIndex: 'networkSegment',
      width: 120,
      resizable: true,
    },
    {
      title: 'IPv6网段',
      dataIndex: 'ipv6NetworkSegment',
      width: 120,
      resizable: true,
    },
    {
      title: '机房名称',
      dataIndex: 'idcName',
      width: 120,
      resizable: true,
    },
    {
      title: '开始时间',
      dataIndex: 'startTime',
      // width: 120,
      resizable: true,
    },
    {
      title: '结束时间',
      dataIndex: 'endTime',
      // width: 120,
      // resizable: true,
    },
  ];
};
export const normalDemandColumns = function (): BasicColumn[] {
  return [
    {
      title: '需求单号',
      dataIndex: 'demandId',
      width: 215,
      resizable: true,
      fixed: 'left',
    },
    {
      title: '剩余可绑定\n带宽(Mbps)',
      dataIndex: 'canBoundedBw',
      width: 90,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '缺口带宽\n(Mbps)',
      dataIndex: 'gapBw',
      width: 90,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '是否独立\n部署',
      dataIndex: 'isIndependentDeploy',
      width: 80,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '需求带宽\n(Mbps)',
      dataIndex: 'demandBw',
      width: 90,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '已交付带宽\n(Mbps)',
      dataIndex: 'deliveredBw',
      width: 90,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '区域名称',
      dataIndex: 'areaName',
      width: 80,
      resizable: true,
    },
    {
      title: '省份',
      dataIndex: 'provinceName',
      width: 100,
      resizable: true,
    },
    {
      title: '设备类型',
      dataIndex: 'devId',
      width: 80,
      resizable: true,
    },
    {
      title: '运营商',
      dataIndex: 'ispName',
      width: 80,
      resizable: true,
    },
    {
      title: '调度控制',
      dataIndex: 'isProvinceScheduling',
      width: 80,
      customRender: ({ text }) => {
        if (text === 1) {
          return h('span', { style: { color: 'red' } }, '仅本省');
        } else if (text === 0) {
          return h('span', { style: { color: 'green' } }, '不限制');
        } else {
          return '未知';
        }
      },
    },
    {
      title: '是否异网',
      dataIndex: 'isCoverDiffIsp',
      width: 80,
      resizable: true,
      customRender: ({ text }) => {
        if (text === 1) {
          return h('span', { style: { color: 'green' } }, '是');
        } else if (text === 0) {
          return h('span', { style: { color: 'red' } }, '否');
        } else {
          return '未知';
        }
      },
    },
    {
      title: '验收中任\n务数',
      dataIndex: 'checkingTaskNum',
      width: 80,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '需求类型\n名称',
      dataIndex: 'demandTypeName',
      width: 80,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '允许同大\n区替换',
      dataIndex: 'isSameAreaReplace',
      width: 80,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    // {
    //   title: '提交验收标识',
    //   dataIndex: 'uploadFlag',
    //   width: 120,
    //   resizable: true,
    // },
    {
      title: 'HDD改造',
      dataIndex: 'hddTransform',
      width: 80,
      resizable: true,
    },
    {
      title: 'IPv4网段',
      dataIndex: 'networkSegment',
      width: 120,
      resizable: true,
    },
    {
      title: 'IPv6网段',
      dataIndex: 'ipv6NetworkSegment',
      width: 120,
      resizable: true,
    },
    {
      title: '机房名称',
      dataIndex: 'idcName',
      width: 120,
      resizable: true,
    },
    {
      title: '开始时间',
      dataIndex: 'startTime',
      // width: 100,
      // resizable: true,
    },
    {
      title: '结束时间',
      dataIndex: 'endTime',
      // width: 120,
      // resizable: true,
    },
  ];
};
export const searchDemandFormSchema = function (bizType: string): FormSchema[] {
  return [
    {
      field: 'demandId',
      label: '需求单号',
      component: 'Input',
      colProps: { span: 8 },
    },
    {
      field: 'area',
      label: '区域',
      component: 'Select',
      colProps: { span: 6 },
      ifShow: () => bizType === 'specialLine',
    },
    {
      field: 'province',
      label: '省份',
      component: 'Select',
      colProps: { span: 6 },
      ifShow: () => bizType === 'specialLine',
    },
    {
      field: 'isp',
      label: '运营商',
      component: 'Select',
      colProps: { span: 6 },
      ifShow: () => bizType === 'specialLine',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.ISP),
        allowClear: true,
      },
    },
  ];
};

export const normalInfoFormSchema = function (): FormSchema[] {
  return [
    {
      field: 'singleDiverybw',
      label: '单线交付带宽(Mbps)',
      component: 'InputNumber',
      componentProps: {
        min: 0,
      },
      helpMessage: '期望交付的单线路带宽',
    },
    {
      field: 'bwCount',
      label: '线路数',
      component: 'InputNumber',
      componentProps: {
        min: 0,
        disabled: true,
        readonly: true,
      },
    },
    {
      field: 'deliveryBw',
      label: '总交付带宽(Mbps)',
      component: 'InputNumber',
      componentProps: {
        min: 0,
        disabled: true,
        readonly: true,
      },
      helpMessage: '总交付带宽=单线交付带宽 * 线路数',
    },
    {
      field: 'isProvinceScheduling',
      label: '调度控制',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.SCHEDULE_TYPE),
        showSearch: true,
        allowClear: false,
      },
      required: true,
      colProps: { span: 8 },
    },
    {
      field: 'province',
      label: '省份',
      component: 'Select',
      componentProps: {
        showSearch: true,
        allowClear: true,
        placeholder: '请选择省份',
      },
      colProps: { span: 8 },
    },
    {
      field: 'city',
      label: '城市',
      component: 'Select',
      componentProps: {
        showSearch: true,
        allowClear: true,
        placeholder: '请选择城市',
      },
      colProps: { span: 8 },
    },
    {
      field: 'isCoverDiffIsp',
      label: '是否异网',
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: '否', value: 0 },
          { label: '是', value: 1 },
        ],
      },
      colProps: { span: 8 },
    },
    {
      field: 'difIsp',
      label: '异网运营商',
      component: 'Select',
      componentProps: {
        options: ispOptions,
        allowClear: false,
        placeholder: '请选择异网运营商',
      },
      colProps: { span: 8 },
    },
  ];
};

export const speciaLineInfoFormSchema = function (): FormSchema[] {
  return [
    {
      field: 'singleDiverybw',
      label: '单线交付带宽(Mbps)',
      component: 'InputNumber',
      componentProps: {
        min: 0,
      },
      helpMessage: '期望交付的单线路带宽',
    },
    {
      field: 'bwCount',
      label: '线路数',
      component: 'InputNumber',
      componentProps: {
        min: 0,
        disabled: true,
        readonly: true,
      },
    },
    {
      field: 'deliveryBw',
      label: '总交付带宽(Mbps)',
      component: 'InputNumber',
      componentProps: {
        min: 0,
        disabled: true,
        readonly: true,
      },
      helpMessage: '总交付带宽=单线交付带宽 * 线路数',
    },
    {
      field: 'networkType',
      label: '网络接入类型',
      component: 'Select',
      componentProps: {
        options: roomTypeOptions,
        allowClear: false,
      },
      required: true,
      colProps: { span: 8 },
    },
    {
      field: 'isProvinceScheduling',
      label: '调度控制',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.SCHEDULE_TYPE),
        showSearch: true,
        allowClear: false,
      },
      required: true,
      colProps: { span: 8 },
    },
    {
      field: 'province',
      label: '省份',
      component: 'Select',
      componentProps: {
        showSearch: true,
        allowClear: true,
        placeholder: '请选择省份',
      },
      colProps: { span: 8 },
    },
    {
      field: 'city',
      label: '城市',
      component: 'Select',
      componentProps: {
        showSearch: true,
        allowClear: true,
        placeholder: '请选择城市',
      },
      colProps: { span: 8 },
    },
    {
      field: 'isCoverDiffIsp',
      label: '是否异网',
      component: 'RadioGroup',
      componentProps: {
        options: [
          { label: '否', value: 0 },
          { label: '是', value: 1 },
        ],
      },
      colProps: { span: 8 },
    },
    {
      field: 'difIsp',
      label: '异网运营商',
      component: 'Select',
      componentProps: {
        options: ispOptions,
        allowClear: false,
        placeholder: '请选择异网运营商',
      },
      colProps: { span: 8 },
    },
  ];
};

export const getDeliveryInfoMacColumns = function (): BasicColumn[] {
  return [
    {
      title: 'Mac地址',
      dataIndex: 'mac',
      width: 200,
      resizable: true,
    },
  ];
};

export const getDeliveryInfoMacSearchForm = function (): FormSchema[] {
  return [
    {
      field: 'mac',
      label: 'Mac地址',
      component: 'Input',
      colProps: { span: 16 },
    },
  ];
};

export const getDifIspForm = function (
  onIsProvinceSchedulingChange,
  onProvinceChange,
): FormSchema[] {
  return [
    {
      field: 'isProvinceScheduling',
      label: '调度控制',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.SCHEDULE_TYPE),
        onChange: onIsProvinceSchedulingChange,
        showSearch: true,
        allowClear: false,
      },
      required: true,
      colProps: { span: 8 },
    },
    {
      field: 'province',
      label: '省份',
      component: 'Select',
      componentProps: {
        onChange: onProvinceChange,
        showSearch: true,
        allowClear: true,
        placeholder: '请选择省份',
      },
      colProps: { span: 8 },
    },
    {
      field: 'city',
      label: '城市',
      component: 'Select',
      componentProps: {
        showSearch: true,
        allowClear: true,
        placeholder: '请选择城市',
      },
      colProps: { span: 8 },
    },
    {
      field: 'difIsp',
      label: '异网运营商',
      component: 'Select',
      componentProps: {
        options: ispOptions,
        allowClear: false,
      },
      required: true,
      colProps: { span: 8 },
    },
  ];
};
