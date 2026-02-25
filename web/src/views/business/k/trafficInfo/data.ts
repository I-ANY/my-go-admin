import { BasicColumn, FormSchema } from '@/components/Table';
import { KEnum, CommonEnum } from '@/enums/dictTypeCode';
import { RangePickPresetsExact } from '@/utils/common';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { commonCustomHeaderCell, customCeilDivide } from '@/utils/util';

export const ispOptions = getSelectOptionsFromDict(KEnum.ISP)?.filter((item) => {
  return item.value != '-1';
});
export const ispMap = getDictDataMapFromDict(KEnum.ISP);

export const devTypeOptions = getSelectOptionsFromDict(KEnum.DEVICE_TYPE);
export const devTypeMap = getDictDataMapFromDict(KEnum.DEVICE_TYPE);

export const deviceTypeOptions = getSelectOptionsFromDict(KEnum.DEVICE_CATEGORY);
export const deviceTypeMap = getDictDataMapFromDict(KEnum.DEVICE_CATEGORY);

export const isProvinceSchedulingMap = getDictDataMapFromDict(KEnum.SCHEDULE_TYPE);
export const isProvinceSchedulingOptions = getSelectOptionsFromDict(KEnum.SCHEDULE_TYPE);
export const isCoverDiffIspMap = getDictDataMapFromDict(KEnum.IS_COVER_DIFF_ISP);
export const isCoverDiffIspOptions = getSelectOptionsFromDict(KEnum.IS_COVER_DIFF_ISP);

export const providerTypeMap = getDictDataMapFromDict(KEnum.PROVIDER_TYPE);
export const providerTypeOptions = getSelectOptionsFromDict(KEnum.PROVIDER_TYPE);

export const ispConsistencyMap = getDictDataMapFromDict(CommonEnum.YES_NO);
export const ispConsistencyOptions = getSelectOptionsFromDict(CommonEnum.YES_NO);

export const searchFormSchema = (onTimePikerOpen, onDemandAreasChange): FormSchema[] => [
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
    field: 'macAddrs',
    label: '上报Mac地址',
    component: 'InputTextArea',
    componentProps: {
      allowClear: true,
      placeholder: '输入Mac进行搜索',
      rows: 3,
    },
    colProps: { span: 6 },
  },
  {
    field: 'providerIds',
    label: '厂商',
    component: 'Select',
    componentProps: {
      options: providerTypeOptions,
      allowClear: true,
      mode: 'multiple',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'ispIds',
    label: '上报运营商',
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
    field: 'demandIspIds',
    label: '需求运营商',
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
    field: 'ispConsistency',
    label: '运营商一致',
    component: 'Select',
    componentProps: {
      options: ispConsistencyOptions,
      allowClear: true,
      // mode: 'multiple',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'deviceTypes',
    label: '上报设备大类',
    component: 'Select',
    componentProps: {
      options: deviceTypeOptions,
      allowClear: true,
      mode: 'multiple',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'devTypes',
    label: '上报设备小类',
    component: 'Select',
    componentProps: {
      options: devTypeOptions,
      allowClear: true,
      mode: 'multiple',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'demandAreas',
    label: '需求大区',
    component: 'Select',
    componentProps: {
      allowClear: true,
      mode: 'multiple',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
      onChange: onDemandAreasChange,
    },
    colProps: { span: 6 },
  },
  {
    field: 'demandProvinces',
    label: '需求省份',
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
    field: 'isProvinceScheduling',
    label: '调度控制',
    component: 'Select',
    componentProps: {
      options: isProvinceSchedulingOptions,
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'isCoverDiffIsp',
    label: '是否异网',
    component: 'Select',
    componentProps: {
      options: isCoverDiffIspOptions,
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'owner',
    label: '机房',
    component: 'Input',
    componentProps: {
      allowClear: true,
      placeholder: '请输入机房进行搜索',
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
  },
  {
    title: '上报Mac地址',
    dataIndex: 'macAddr',
    width: 160,
    resizable: true,
  },
  {
    title: '上报\n运营商',
    dataIndex: 'ispId',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '需求\n运营商',
    dataIndex: 'demandIspId',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '运营商一致',
    dataIndex: 'ispConsistency',
    width: 120,
    resizable: true,
    helpMessage: '当上报运营商与需求运营商一致时，为是，否则为否',
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '机房',
    dataIndex: 'owner',
    width: 120,
    resizable: true,
  },
  {
    title: '上报设备\n大类',
    dataIndex: 'deviceType',
    width: 80,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '上报设备\n小类',
    dataIndex: 'devType',
    width: 100,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '厂商',
    dataIndex: 'providerId',
    width: 100,
    resizable: true,
  },
  {
    title: '需求省份',
    dataIndex: 'demandProvince',
    width: 100,
    resizable: true,
  },
  {
    title: '需求大区',
    dataIndex: 'demandArea',
    width: 80,
    resizable: true,
  },
  {
    title: '调度控制',
    dataIndex: 'isProvinceScheduling',
    width: 100,
    resizable: true,
  },
  {
    title: '是否异网',
    dataIndex: 'isCoverDiffIsp',
    width: 80,
    resizable: true,
  },
  {
    title: '时间',
    dataIndex: 'time',
    width: 160,
    resizable: true,
  },
  // {
  //   title: '已入库流量\nbyte',
  //   dataIndex: 'registerFlow',
  //   width: 120,
  //   resizable: true,
  //   customHeaderCell: commonCustomHeaderCell(),
  // },
  // {
  //   title: '未入库流量\nbyte',
  //   dataIndex: 'unRegisterFlow',
  //   width: 120,
  //   resizable: true,
  //   customHeaderCell: commonCustomHeaderCell(),
  // },
  // {
  //   title: '总流量\nbyte',
  //   dataIndex: 'allFlow',
  //   width: 120,
  //   resizable: true,
  //   customHeaderCell: commonCustomHeaderCell(),
  // },
  {
    title: '已入库速率\nMbps',
    dataIndex: 'registerSpeed',
    width: 120,
    resizable: true,
    customRender: ({ record }) => {
      return customCeilDivide(record.registerSpeed);
    },
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '未入库速率\nMbps',
    dataIndex: 'unRegisterSpeed',
    width: 120,
    resizable: true,
    customRender: ({ record }) => {
      return customCeilDivide(record.unRegisterSpeed);
    },
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '总速率\nMbps',
    dataIndex: 'allSpeed',
    // width: 120,
    // resizable: true,
    customRender: ({ record }) => {
      return customCeilDivide(record.allSpeed);
    },
    customHeaderCell: commonCustomHeaderCell(),
  },
  // {
  //   title: 'frankID',
  //   dataIndex: 'frankID',
  // },
];
