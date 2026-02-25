import { BasicColumn, FormSchema } from '@/components/Table';
import { CommonEnum, RoomEnum } from '@/enums/dictTypeCode';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { commonCustomHeaderCell, customCeilDivide } from '@/utils/util';
import { h } from 'vue';
import { speedLimitStausOptions } from '../speed_limit/data';
import { Tooltip } from 'ant-design-vue';

export const roomStatusOptions = getSelectOptionsFromDict(RoomEnum.STATUS);
export const roomStatusMap = getDictDataMapFromDict(RoomEnum.STATUS);

export const roomChargeModeOptions = getSelectOptionsFromDict(RoomEnum.CHARGE_MODE);
export const roomChargeModeMap = getDictDataMapFromDict(RoomEnum.CHARGE_MODE);

export const roomPeakEnableStatusOptions = getSelectOptionsFromDict(RoomEnum.PEAK_ENABLE_STATUS);
export const roomPeakEnableStatusMap = getDictDataMapFromDict(RoomEnum.PEAK_ENABLE_STATUS);
export const roomTypeOptions = getSelectOptionsFromDict(RoomEnum.TYPE);
export const roomTypeMap = getDictDataMapFromDict(RoomEnum.TYPE);
export const roomUseModeOptions = getSelectOptionsFromDict(RoomEnum.USE_MODE);
export const roomUseModeMap = getDictDataMapFromDict(RoomEnum.USE_MODE);

export const ispOptions = getSelectOptionsFromDict(CommonEnum.ISP);
export const ispMap = getDictDataMapFromDict(CommonEnum.ISP);

// 1-单端口95削峰，2-合并95削峰，3-不削峰
export enum roomUseModEnum {
  SINGLE_PORT_95_PEAK_SHAVING = 1,
  MERGE_95_PEAK_SHAVING = 2,
  NO_PEAK_SHAVING = 3,
}

export const getNumberResult = (
  value: number | undefined | null,
): number | null | undefined | string => {
  if (value === null || value === undefined) {
    return;
  }
  if (value === 0) {
    return 0;
  }
  if (value == 0.01) {
    return '≤0.01';
  }
  return value;
};

export const searchFormSchema = (): FormSchema[] => [
  {
    field: 'month',
    label: '月份',
    component: 'DatePicker',
    componentProps: {
      allowClear: false,
      format: 'YYYY-MM',
      style: {
        width: '100%',
      },
      picker: 'month',
    },
    required: true,
    colProps: { span: 8, xxl: 6, lg: 6 },
  },
  {
    field: 'name',
    label: '群号|节点|机房',
    component: 'Input',
    componentProps: {
      allowClear: true,
      placeholder: '输入群号|节点|机房进行搜索',
    },
    colProps: { span: 8, xxl: 6, lg: 6 },
  },
  {
    field: 'useModes',
    label: '使用模式',
    component: 'Select',
    colProps: { span: 8, xxl: 6, lg: 6 },
    componentProps: {
      allowClear: true,
      options: roomUseModeOptions,
      mode: 'multiple',
      maxTagCount: 1,
    },
  },
  {
    field: 'peakBusiness',
    label: '削峰业务',
    component: 'Select',
    componentProps: {
      options: [],
      allowClear: true,
      showSearch: true,
      // mode: 'multiple',
    },
    colProps: { span: 8, xxl: 6, lg: 6 },
  },
  {
    field: 'roomStatus',
    label: '机房状态',
    component: 'Select',
    colProps: { span: 8, xxl: 6, lg: 6 },
    componentProps: {
      allowClear: true,
      options: roomStatusOptions,
    },
  },
  {
    field: 'roomTypes',
    label: '机房类型',
    component: 'Select',
    colProps: { span: 8, xxl: 6, lg: 6 },
    componentProps: {
      allowClear: true,
      options: roomTypeOptions,
      mode: 'multiple',
      maxTagCount: 1,
    },
  },
  {
    field: 'chargeModes',
    label: '计费模式',
    component: 'Select',
    componentProps: {
      options: roomChargeModeOptions,
      allowClear: true,
      mode: 'multiple',
      maxTagCount: 1,
    },
    colProps: { span: 8, xxl: 6, lg: 6 },
  },
  {
    field: 'peakEnabled',
    label: '削峰配置状态',
    component: 'Select',
    componentProps: {
      options: roomPeakEnableStatusOptions,
      allowClear: true,
    },
    colProps: { span: 8, xxl: 6, lg: 6 },
  },
  {
    field: 'isps',
    label: '运营商',
    component: 'Select',
    componentProps: {
      options: ispOptions,
      allowClear: true,
      mode: 'multiple',
    },
    colProps: { span: 8, xxl: 6, lg: 6 },
  },
  {
    field: 'status',
    label: '限速开关',
    component: 'Select',
    componentProps: {
      options: speedLimitStausOptions,
      allowClear: true,
      // mode: 'multiple',
    },
    colProps: { span: 8, xxl: 6, lg: 6 },
  },
];
export const columns: BasicColumn[] = [
  {
    title: '群号|节点|机房',
    dataIndex: 'name',
    width: 140,
    resizable: true,
    fixed: 'left',
  },
  {
    title: '机房类型',
    dataIndex: 'roomType',
    width: 75,
    // resizable: true,
  },
  {
    title: '地区',
    dataIndex: 'location',
    width: 100,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'isp',
    width: 75,
    // resizable: true,
  },
  {
    title: '月份',
    dataIndex: 'month',
    width: 75,
    // resizable: true,
  },
  {
    title: '机房状态',
    dataIndex: 'roomStatus',
    width: 75,
    // resizable: true,
  },
  {
    title: '使用模式',
    dataIndex: 'useMode',
    width: 110,
    // resizable: true,
  },
  {
    title: '机房带宽\nGbps',
    dataIndex: 'bandwidth',
    width: 85,
    sorter: true,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return getNumberResult(customCeilDivide(record.bandwidth, 1000 * 1000 * 1000, 2) as any);
    },
  },
  {
    title: '保底带宽\nGbps',
    dataIndex: 'guaranteedBw',
    width: 85,
    sorter: true,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return getNumberResult(customCeilDivide(record.guaranteedBw, 1000 * 1000 * 1000, 2) as any);
    },
  },
  {
    title: '预95带宽\nGbps',
    dataIndex: 'preset95Bw',
    width: 85,
    sorter: true,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return getNumberResult(customCeilDivide(record.preset95Bw, 1000 * 1000 * 1000, 2) as any);
    },
  },
  {
    title: '可削峰带\n宽Gbps',
    dataIndex: 'canPeakShavingBw',
    width: 85,
    sorter: true,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return getNumberResult(
        customCeilDivide(record.canPeakShavingBw, 1000 * 1000 * 1000, 2) as any,
      );
    },
  },
  {
    title: '上月95值\nGbps',
    dataIndex: 'lastMonth95Bw',
    width: 85,
    sorter: true,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return getNumberResult(customCeilDivide(record.lastMonth95Bw, 1000 * 1000 * 1000, 2) as any);
    },
  },
  {
    title: '当前95值\nGbps',
    dataIndex: 'now95Bw',
    width: 85,
    sorter: true,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return getNumberResult(customCeilDivide(record.now95Bw, 1000 * 1000 * 1000, 2) as any);
    },
  },
  {
    title: '全月95值\nGbps',
    dataIndex: 'month95Bw',
    width: 85,
    sorter: true,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return getNumberResult(customCeilDivide(record.month95Bw, 1000 * 1000 * 1000, 2) as any);
    },
  },
  {
    title: '剩余/已用/总共\n削峰点数',
    dataIndex: 'rest95PointCount',
    customHeaderCell: commonCustomHeaderCell(),
    width: 140,
    resizable: true,
    sorter: true,
    customRender: ({ record }) => {
      const used95PointCount =
        record.used95PointCount === null || record.used95PointCount === undefined
          ? '-'
          : record.used95PointCount;
      const maxPeakShavingPointCount =
        record.maxPeakShavingPointCount === null || record.maxPeakShavingPointCount === undefined
          ? '-'
          : record.maxPeakShavingPointCount;
      const restCount =
        record.rest95PointCount === null || record.rest95PointCount === undefined
          ? '-'
          : record.rest95PointCount;

      return `${restCount}/${used95PointCount}/${maxPeakShavingPointCount}`;
    },
  },
  {
    title: '剩余削峰\n时长(h)',
    dataIndex: 'restPeakShavingMinute',
    width: 85,
    sorter: true,
    customHeaderCell: commonCustomHeaderCell(),
    resizable: true,
    customRender: ({ record }) => {
      if (record.restPeakShavingMinute == null || record.restPeakShavingMinute == undefined) {
        return '';
      }
      const result = getNumberResult(customCeilDivide(record.restPeakShavingMinute, 60, 2) as any);
      if (record.restPeakShavingMinute < 180 * (record.portCount || 1)) {
        if (record.useMode == roomUseModEnum.SINGLE_PORT_95_PEAK_SHAVING) {
          return h(
            Tooltip,
            {
              title: `剩余削峰时长小于 ${3} * ${record.portCount || 1} = ${3 * (record.portCount || 1)} 小时`,
            },
            {
              default: () => h('span', { style: { color: 'red' } }, result as any),
            },
          );
        } else {
          return h('span', { style: { color: 'red' } }, result as any);
        }
      }
      return h('span', {}, result as any);
    },
  },
  {
    title: '削峰业务',
    dataIndex: 'peakBusinessNames',
    width: 210,
    resizable: true,
  },
  {
    title: '削峰配\n置状态',
    dataIndex: 'peakEnabled',
    width: 75,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '计费方式',
    dataIndex: 'chargeMode',
    width: 75,
    resizable: true,
  },
  {
    title: '起始计费\n时间',
    dataIndex: 'billingStartAt',
    width: 145,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '限速开关',
    dataIndex: 'status',
    width: 75,
  },
  {
    title: '方案',
    dataIndex: 'plans',
    width: 480,
    resizable: true,
  },
  {
    title: '备注',
    dataIndex: 'remark',
    width: 130,
    resizable: true,
  },
  {
    title: '上次计算\n95值时间',
    dataIndex: 'calculateTime',
    width: 145,
    resizable: true,
    sorter: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '最后\n更新人',
    dataIndex: 'updateUser',
    minWidth: 100,
    width: 100,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
  },
  {
    title: '更新时间',
    dataIndex: 'updateByUserAt',
    width: 145,
    resizable: true,
    sorter: true,
  },
  {
    title: '创建时间',
    dataIndex: 'createdAt',
    // width: 150,
    // resizable: true,
    sorter: true,
  },
  // {
  //   title: '限速配置',
  //   dataIndex: 'limitSpeedConfig',
  //   children: [
  //     {
  //       title: '开启限速',
  //       dataIndex: 'status',
  //       width: 85,
  //     },
  //     {
  //       title: '限速值Gbps',
  //       dataIndex: 'limitValue',
  //       width: 85,
  //     },
  //     {
  //       title: '开启限速\n阈值Gbps',
  //       dataIndex: 'limitThreshold',
  //       width: 85,
  //       customHeaderCell: commonCustomHeaderCell(),
  //     },
  //     {
  //       title: '解除限速\n阈值Gbps',
  //       dataIndex: 'unlimitThreshold',
  //       width: 85,
  //       customHeaderCell: commonCustomHeaderCell(),
  //     },
  //     {
  //       title: '限速顺序',
  //       dataIndex: 'businessSorts',
  //       width: 120,
  //       resizable: true,
  //     },
  //   ],
  // },
];

export const updateRoomPeakFormSchema = function (
  maxBw: number,
  getFormFn: () => any,
): FormSchema[] {
  return [
    {
      field: 'useMode',
      label: '使用模式',
      component: 'Select',
      colProps: { span: 12 },
      required: true,
      componentProps: {
        allowClear: false,
        options: roomUseModeOptions,
        onChange: (_value: any) => {
          const { getFieldsValue, updateSchema } = getFormFn();
          const values = getFieldsValue();
          if (values.useMode == roomUseModEnum.SINGLE_PORT_95_PEAK_SHAVING) {
            updateSchema([
              {
                field: 'canPeakShavingBw',
                show: true,
                required: true,
              },
              {
                field: 'singlePortPreset95Bw',
                show: true,
                required: true,
              },
            ]);
          } else {
            updateSchema([
              {
                field: 'canPeakShavingBw',
                show: false,
                required: false,
              },
              {
                field: 'singlePortPreset95Bw',
                show: false,
                required: false,
              },
            ]);
          }
        },
      },
    },
    {
      field: 'bandwidth',
      label: '机房带宽',
      component: 'InputNumber',
      componentProps: {
        min: 0,
        max: 9999999.99,
        addonAfter: 'Gbps',
        precision: 2,
        step: 0.01,
        style: {
          width: '100%',
        },
        disabled: true,
        readonly: true,
      },
      colProps: { span: 12 },
    },
    {
      field: 'guaranteedBw',
      label: '保底带宽',
      component: 'InputNumber',
      componentProps: {
        min: 0,
        max: 9999999.99,
        addonAfter: 'Gbps',
        precision: 2,
        step: 0.01,
        style: {
          width: '100%',
        },
      },
      colProps: { span: 12 },
      required: true,
    },
    {
      field: 'preset95Bw',
      label: '预95带宽',
      component: 'InputNumber',
      componentProps: {
        min: 0,
        max: 9999999.99,
        addonAfter: 'Gbps',
        precision: 2,
        step: 0.01,
        style: {
          width: '100%',
        },
      },
      colProps: { span: 12 },
    },
    {
      field: 'canPeakShavingBw',
      label: '可削峰带宽',
      component: 'InputNumber',
      componentProps: {
        min: 0,
        max: 9999999.99,
        addonAfter: 'Gbps',
        precision: 2,
        step: 0.01,
        style: {
          width: '100%',
        },
      },
      show: false,
      colProps: { span: 12 },
    },
    {
      field: 'singlePortPreset95Bw',
      label: '单口预95带宽',
      component: 'InputNumber',
      helpMessage: '对于单口95削峰类型的机房，该值用于计算单端口的削峰时长',
      componentProps: {
        min: 0,
        max: 9999999.99,
        addonAfter: 'Gbps',
        precision: 2,
        step: 0.01,
        style: {
          width: '100%',
        },
      },
      show: false,
      colProps: { span: 12 },
      required: true,
    },
    {
      field: 'billingStartAt',
      label: '起始计费时间',
      component: 'DatePicker',
      componentProps: {
        allowClear: true,
        format: 'YYYY-MM-DD HH:mm:ss',
        style: {
          width: '100%',
        },
        showNow: true,
        showTime: true,
      },
      // required: true,
      colProps: { span: 24 },
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
    {
      field: 'plans',
      component: 'Divider',
      label: '方案',
      helpMessage:
        '打峰方案格式必须满足格式，例如：打峰方案-10.11-10.31日，1口打峰（18:00-23:59），峰值214.1G削峰55.7G，保底限速158.4G',
      componentProps: {
        orientation: 'center',
      },
      colProps: { span: 24 },
    },
  ];
};
export const Improve95PredictionColumns = (): BasicColumn[] => [
  {
    title: '拉升带宽\nGbps',
    dataIndex: 'improveBw',
    width: 100,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return getNumberResult(customCeilDivide(record.improveBw, 1000 * 1000 * 1000, 2) as any);
    },
  },
  {
    title: '拉升后带宽\nGbps',
    dataIndex: 'afterImproveBw',
    width: 100,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return getNumberResult(customCeilDivide(record.afterImproveBw, 1000 * 1000 * 1000, 2) as any);
    },
  },
  {
    title: '剩余/已用/总共\n削峰点数',
    dataIndex: 'rest95PointCount',
    customHeaderCell: commonCustomHeaderCell(),
    width: 140,
    resizable: true,
    customRender: ({ record }) => {
      const used95PointCount =
        record.used95PointCount === null || record.used95PointCount === undefined
          ? '-'
          : record.used95PointCount;
      const maxPeakShavingPointCount =
        record.maxPeakShavingPointCount === null || record.maxPeakShavingPointCount === undefined
          ? '-'
          : record.maxPeakShavingPointCount;
      const restCount =
        record.rest95PointCount === null || record.rest95PointCount === undefined
          ? '-'
          : record.rest95PointCount;

      return `${restCount}/${used95PointCount}/${maxPeakShavingPointCount}`;
    },
  },
  {
    title: '剩余时长(h)',
    dataIndex: 'restPeakShavingMinute',
    width: 85,
    resizable: true,
    customRender: ({ record }) => {
      if (record.restPeakShavingMinute != null && record.restPeakShavingMinute != undefined) {
        const result = getNumberResult(
          customCeilDivide(record.restPeakShavingMinute, 60, 2) as any,
        );
        if (record.restPeakShavingMinute < 180) {
          return h('span', { style: { color: 'red' } }, result as any);
        }
        return h('span', {}, result as any);
      }
      return '';
    },
  },
];

export const SinglePort95Columns = (getTabelFunctions): BasicColumn[] => [
  {
    title: '群号|节点|机房',
    dataIndex: 'roomName',
    width: 160,
    resizable: true,
    fixed: 'left',
    customCell: (_, index) => {
      const { getDataSource } = getTabelFunctions();
      const ds = getDataSource();
      if (ds.length > 0 && index === ds.length - 1) {
        return {
          colSpan: 3,
        };
      }
      return {
        colSpan: 1,
      };
    },
  },
  {
    title: '端口名称',
    dataIndex: 'name',
    width: 160,
    resizable: true,
    customCell: (_, index) => {
      const { getDataSource } = getTabelFunctions();
      const ds = getDataSource();
      if (ds.length > 0 && index === ds.length - 1) {
        return {
          colSpan: 0,
        };
      }
      return {
        colSpan: 1,
      };
    },
  },
  {
    title: '预95带宽\nGbps',
    dataIndex: 'preset95Bw',
    width: 85,
    sorter: true,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return getNumberResult(customCeilDivide(record.preset95Bw, 1000 * 1000 * 1000, 2) as any);
    },
    customCell: (_, index) => {
      const { getDataSource } = getTabelFunctions();
      const ds = getDataSource();
      if (ds.length > 0 && index === ds.length - 1) {
        return {
          colSpan: 0,
        };
      }
      return {
        colSpan: 1,
      };
    },
  },

  {
    title: '当前95值\nGbps',
    dataIndex: 'now95Bw',
    width: 85,
    sorter: true,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return getNumberResult(customCeilDivide(record.now95Bw, 1000 * 1000 * 1000, 2) as any);
    },
  },
  {
    title: '全月95值\nGbps',
    dataIndex: 'month95Bw',
    width: 85,
    sorter: true,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return getNumberResult(customCeilDivide(record.month95Bw, 1000 * 1000 * 1000, 2) as any);
    },
  },
  {
    title: '剩余/已用/总共\n削峰点数',
    dataIndex: 'rest95PointCount',
    customHeaderCell: commonCustomHeaderCell(),
    width: 140,
    resizable: true,
    sorter: true,
    customRender: ({ record }) => {
      const used95PointCount =
        record.used95PointCount === null || record.used95PointCount === undefined
          ? '-'
          : record.used95PointCount;
      const maxPeakShavingPointCount =
        record.maxPeakShavingPointCount === null || record.maxPeakShavingPointCount === undefined
          ? '-'
          : record.maxPeakShavingPointCount;
      const restCount =
        record.rest95PointCount === null || record.rest95PointCount === undefined
          ? '-'
          : record.rest95PointCount;

      return `${restCount}/${used95PointCount}/${maxPeakShavingPointCount}`;
    },
  },
  {
    title: '剩余削峰\n时长(h)',
    dataIndex: 'restPeakShavingMinute',
    width: 85,
    sorter: true,
    customHeaderCell: commonCustomHeaderCell(),
    resizable: true,
    customRender: ({ record }) => {
      if (record.restPeakShavingMinute != null && record.restPeakShavingMinute != undefined) {
        const result = getNumberResult(
          customCeilDivide(record.restPeakShavingMinute, 60, 2) as any,
        );
        if (record.restPeakShavingMinute < 180) {
          return h('span', { style: { color: 'red' } }, result as any);
        }
        return h('span', {}, result as any);
      }
      return '';
    },
  },
];

/**
 * 验证打峰方案格式的正则表达式校验函数
 * 格式：打峰方案-日期范围，描述（时间范围），峰值数据削峰数据，保底限速数据
 * 示例：打峰方案-10.11-10.31日，1口打峰（18:00-23:59），峰值214.1G削峰55.7G，保底限速158.4G
 */
export function validatePeakPlanFormat(text: string): boolean {
  if (!text.startsWith('打峰方案')) {
    return true;
  }

  // 正则表达式匹配打峰方案格式
  const peakPlanRegex =
    /^打峰方案-\d{1,2}\.\d{1,2}-\d{1,2}\.\d{1,2}日，\d+口打峰（\d{1,2}:\d{2}-\d{1,2}:\d{2}），峰值\d+(?:\.\d+)?G削峰\d+(?:\.\d+)?G，保底限速\d+(?:\.\d+)?G.*/;
  return peakPlanRegex.test(text);
}
