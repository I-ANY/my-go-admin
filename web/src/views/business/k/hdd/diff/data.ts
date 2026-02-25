import { GetBusinessList } from '@/api/business/k';
import { BasicColumn, FormSchema } from '@/components/Table';
import { KEnum } from '@/enums/dictTypeCode';
import { RangeDataPickPresetsExact } from '@/utils/common';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { commonCustomHeaderCell } from '@/utils/util';

export const hddLoseReasonMap = getDictDataMapFromDict(KEnum.HDD_LOSE_REASON);
export const hddLoseReasonOptions = getSelectOptionsFromDict(KEnum.HDD_LOSE_REASON);

export const isProvinceSchedulingMap = getDictDataMapFromDict(KEnum.SCHEDULE_TYPE);
export const isProvinceSchedulingOptions = getSelectOptionsFromDict(KEnum.SCHEDULE_TYPE);
export const isCoverDiffIspMap = getDictDataMapFromDict(KEnum.IS_COVER_DIFF_ISP);
export const isCoverDiffIspOptions = getSelectOptionsFromDict(KEnum.IS_COVER_DIFF_ISP);

export const ispOptions = getSelectOptionsFromDict(KEnum.ISP)?.filter((item) => {
  return item.value != '-1';
});
export const ispMap = getDictDataMapFromDict(KEnum.ISP);

export const devTypeOptions = getSelectOptionsFromDict(KEnum.DEVICE_TYPE);
export const devTypeMap = getDictDataMapFromDict(KEnum.DEVICE_TYPE);

// 包装函数：将字符串数组转换为 ApiSelect 需要的格式
const GetBusinessListOptions = async (params: Recordable) => {
  const data = await GetBusinessList(params);
  // 如果返回的是字符串数组，转换为对象数组
  if (Array.isArray(data)) {
    return data.map((item) => ({
      label: item,
      value: item,
    }));
  }

  return [];
};

export const searchFormSchema = (onOpenChange): FormSchema[] => [
  {
    field: '[reportDayBegin, reportDayEnd]',
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
      onOpenChange: onOpenChange,
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
    field: 'business',
    label: '业务名称',
    component: 'ApiSelect',
    componentProps: {
      mode: 'multiple',
      api: GetBusinessListOptions,
      allowClear: true,
      params: {
        page: 1,
        pageSize: 10000,
        bizType: 'normal',
      },
      maxTagCount: 3,
    },
    colProps: { span: 6 },
  },
  {
    field: 'lossReasons',
    label: '流失原因',
    component: 'Select',
    componentProps: {
      options: hddLoseReasonOptions,
      allowClear: true,
      mode: 'multiple',
      maxTagCount: 3,
    },
    colProps: { span: 6 },
  },
  {
    field: 'devTypes',
    label: '设备小类',
    component: 'Select',
    componentProps: {
      options: devTypeOptions,
      allowClear: true,
      mode: 'multiple',
      maxTagCount: 3,
    },
    colProps: { span: 6 },
  },
  {
    field: 'owner',
    label: '节点',
    component: 'Input',
    colProps: { span: 6 },
  },
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
    title: '分区名',
    dataIndex: 'partitionName',
    width: 120,
    resizable: true,
  },
  {
    title: '业务名称',
    dataIndex: 'business',
    width: 140,
    resizable: true,
  },
  {
    title: '节点',
    dataIndex: 'owner',
    width: 100,
    resizable: true,
  },
  {
    title: '设备小类',
    dataIndex: 'devType',
    width: 80,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'ispId',
    width: 80,
    resizable: true,
  },
  {
    title: '是否异网',
    dataIndex: 'isCoverDiffIsp',
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
    title: '日期',
    dataIndex: 'reportDay',
    width: 100,
    resizable: true,
  },
  {
    title: '晚高峰时刻',
    dataIndex: 'peakMoment',
    width: 160,
    resizable: true,
  },
  {
    title: '前一日\n晚高峰时刻',
    dataIndex: 'lastDayPeakMoment',
    width: 160,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '确认人',
    dataIndex: 'confirmUser',
    width: 100,
    resizable: true,
  },
  {
    title: '流失原因',
    dataIndex: 'lossReason',
    width: 120,
    resizable: true,
  },
  {
    title: '备注',
    dataIndex: 'remark',
    // width: 100,
    // resizable: true,
  },
];

export const getLossReasonFormSchema = function (onLossReasonChange): FormSchema[] {
  return [
    {
      field: 'lossReason',
      label: '流失原因',
      component: 'Select',
      componentProps: {
        options: hddLoseReasonOptions,
        allowClear: false,
        onChange: onLossReasonChange,
      },
      colProps: { span: 24 },
      required: true,
    },
    {
      field: 'remark',
      label: '备注',
      component: 'InputTextArea',
      componentProps: {
        allowClear: true,
        placeholder: '输入备注',
        rows: 4,
      },
      colProps: { span: 24 },
      required: false,
    },
  ];
};
