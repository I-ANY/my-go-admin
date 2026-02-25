import { BasicColumn, FormSchema } from '@/components/Table';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { RoomEnum } from '@/enums/dictTypeCode';
import { customCeilDivide } from '@/utils/util';
import { getNumberResult } from '../peak/data';

export const speedLimitStausOptions = getSelectOptionsFromDict(RoomEnum.SPEED_LIMIT_STATUS);
export const speedLimitStausMap = getDictDataMapFromDict(RoomEnum.SPEED_LIMIT_STATUS);
export const speedLimitOperStatusOptions = getSelectOptionsFromDict(
  RoomEnum.SPEED_LIMIT_OPER_STATUS,
);
export const speedLimitOperStatusMap = getDictDataMapFromDict(RoomEnum.SPEED_LIMIT_OPER_STATUS);
export const speedLimitOperTypeOptions = getSelectOptionsFromDict(RoomEnum.SPEED_LIMIT_OPER_TYPE);
export const speedLimitOperTypeMap = getDictDataMapFromDict(RoomEnum.SPEED_LIMIT_OPER_TYPE);

export enum limitStatusEnum {
  UNLIMITED = 0,
  LIMIT = 1,
}
export enum switchIsMainEnum {
  YES = 1,
  NO = 2,
}
export const swithcValueValidator = function (getFormHandleFn: any) {
  return (_rule, _value) => {
    const { getFieldsValue } = getFormHandleFn();
    const values = getFieldsValue();
    const totalTrafficSwitches = values.totalTrafficSwitches || [];
    const deductTrafficSwitches = values.deductTrafficSwitches || [];
    if (totalTrafficSwitches.length > 0 && deductTrafficSwitches.length > 0) {
      return Promise.reject('“总流量交换机”和“流量扣减交换机”不能同时填写');
    }
    return Promise.resolve();
  };
};

export const speedLimitConfiguretionFormSchema = (
  onStatusChange,
  limitValueValidator,
  onValueChange,
  onExcludeSwitchesChange,
  getFormHandleFn,
): FormSchema[] => {
  const res: FormSchema[] = [
    {
      field: 'status',
      label: '限速开关',
      component: 'Switch',
      // required: true,
      colProps: { span: 12 },
      componentProps: {
        checkedChildren: '开启',
        unCheckedChildren: '关闭',
        checkedValue: limitStatusEnum.LIMIT,
        unCheckedValue: limitStatusEnum.UNLIMITED,
        onChange: onStatusChange,
      },
      defaultValue: limitStatusEnum.UNLIMITED,
    },
    {
      field: 'limitValue',
      label: '限速值',
      component: 'InputNumber',
      helpMessage: '当机房实时带宽高于“开启限速阈值”值时将带宽限制为该值',
      componentProps: {
        min: 0.01,
        max: 9999999.99,
        addonAfter: 'Gbps',
        precision: 2,
        step: 0.01,
        style: {
          width: '100%',
        },
        onChange: onValueChange,
      },
      rules: [{ validator: limitValueValidator, trigger: 'change' }],
      colProps: { span: 12 },
      required: false,
    },
    {
      field: 'limitThreshold',
      label: '开启限速阈值',
      component: 'InputNumber',
      helpMessage: '当机房实时带宽高于该值时触发限速动作',
      componentProps: {
        min: 0.01,
        max: 9999999.99,
        addonAfter: 'Gbps',
        precision: 2,
        step: 0.01,
        style: {
          width: '100%',
        },
        onChange: onValueChange,
      },
      rules: [{ validator: limitValueValidator, trigger: 'change' }],
      colProps: { span: 12 },
      required: false,
    },
    {
      field: 'unlimitThreshold',
      label: '解除限速阈值',
      component: 'InputNumber',
      helpMessage: '当机房实时带宽低于该值时触发解除限速动作',
      componentProps: {
        min: 0.01,
        max: 9999999.99,
        addonAfter: 'Gbps',
        precision: 2,
        step: 0.01,
        style: {
          width: '100%',
        },
        onChange: onValueChange,
      },
      rules: [{ validator: limitValueValidator, trigger: 'change' }],
      colProps: { span: 12 },
      required: false,
    },
    {
      field: 'excludeSwitches',
      label: '排除交换机',
      component: 'Select',
      helpMessage: '限速时将不采集选中的交换机信息，也不会针对交换机做限速',
      componentProps: {
        mode: 'multiple',
        options: [],
        maxTagCount: 3,
        filterOption: (inputValue, option) => {
          return option?.label.toUpperCase().indexOf(inputValue.toUpperCase()) !== -1;
        },
        onChange: onExcludeSwitchesChange,
      },
      colProps: { span: 24 },
    },
    {
      field: 'totalTrafficSwitches',
      label: '总流量交换机',
      component: 'Select',
      helpMessage:
        '不填写时将使用机房内所有主交换机总流量作为机房总流量，填写后将使用指定交换机总流量作为机房总流量',
      componentProps: {
        mode: 'multiple',
        options: [],
        maxTagCount: 3,
        filterOption: (inputValue, option) => {
          return option?.label.toUpperCase().indexOf(inputValue.toUpperCase()) !== -1;
        },
        onChange: function () {
          getFormHandleFn().validateFields(['deductTrafficSwitches', 'totalTrafficSwitches']);
        },
      },
      rules: [{ validator: swithcValueValidator(getFormHandleFn), trigger: 'change' }],
      colProps: { span: 24 },
    },
    {
      field: 'deductTrafficSwitches',
      label: '流量扣减交换机',
      component: 'Select',
      helpMessage: '当填写时，在计算机房总流量时会扣减选中的交换机的流量',
      componentProps: {
        mode: 'multiple',
        options: [],
        maxTagCount: 3,
        filterOption: (inputValue, option) => {
          return option?.label.toUpperCase().indexOf(inputValue.toUpperCase()) !== -1;
        },
        onChange: function () {
          getFormHandleFn().validateFields(['totalTrafficSwitches', 'deductTrafficSwitches']);
        },
      },
      rules: [{ validator: swithcValueValidator(getFormHandleFn), trigger: 'change' }],
      colProps: { span: 24 },
    },
    {
      field: 'divider-selects',
      component: 'Divider',
      label: '低优业务限速顺序（拖动排序，越靠前越先限速）',
      helpMessage: [
        '低优业务限速顺序（越靠前越先限速）！！！',
        '高优业务不要填写！！！',
        '高优业务不要填写！！！',
        '高优业务不要填写！！！',
      ],
      componentProps: {
        orientation: 'center',
      },
      colProps: {
        span: 24,
      },
    },
  ];
  return res;
};

export const getSpeedLimitRecordSearchForm = function (onTimePikerOpen): FormSchema[] {
  return [
    {
      field: '[operatorTimeBegin, operatorTimeEnd]',
      label: '操作时间',
      component: 'RangePicker',
      componentProps: {
        allowClear: false,
        format: 'YYYY-MM-DD HH:mm:ss',
        showTime: { format: 'HH:mm:ss' },
        placeholder: ['开始时间', '结束时间'],
        style: {
          width: '100%',
        },
        onOpenChange: onTimePikerOpen,
      },
      colProps: { span: 8 },
      required: true,
    },
    {
      field: 'switchDesc',
      label: '交换机名称',
      component: 'Input',
      colProps: { span: 8 },
    },
    {
      label: '端口名称',
      field: 'portName',
      component: 'Input',
      colProps: { span: 8 },
    },
    {
      label: '操作类型',
      field: 'handleType',
      component: 'Select',
      colProps: { span: 8 },
      componentProps: {
        options: speedLimitOperTypeOptions,
      },
    },
    {
      label: '操作状态',
      field: 'status',
      component: 'Select',
      colProps: { span: 8 },
      componentProps: {
        options: speedLimitOperStatusOptions,
      },
    },
  ];
};
export const getSpeedLimitRecordColumns = function (): BasicColumn[] {
  return [
    {
      title: '交换机名称',
      dataIndex: 'switchDesc',
      width: 300,
      resizable: true,
    },
    {
      title: '机房速率(Mbps)',
      dataIndex: 'roomSpeed',
      width: 120,
      customRender: ({ record }) => {
        return getNumberResult(customCeilDivide(record.roomSpeed, 1000 * 1000, 2) as any);
      },
    },
    {
      title: '端口名称',
      dataIndex: 'portName',
      width: 100,
      resizable: true,
    },
    {
      title: '端口业务',
      dataIndex: 'business',
      width: 80,
    },
    {
      title: '端口速率(Mbps)',
      dataIndex: 'ifSpeed',
      width: 120,
      customRender: ({ record }) => {
        return getNumberResult(customCeilDivide(record.ifSpeed, 1000 * 1000, 2) as any);
      },
    },
    {
      title: '操作类型',
      dataIndex: 'handleType',
      width: 90,
    },
    {
      title: '操作状态',
      dataIndex: 'status',
      width: 90,
    },
    {
      title: '操作时间',
      dataIndex: 'operatorTime',
      width: 150,
    },
    {
      title: '终端日志',
      dataIndex: 'consoleLog',
      // width: 80,
    },
    {
      title: '备注',
      dataIndex: 'remark',
      width: 300,
    },
  ];
};
