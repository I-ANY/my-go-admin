import { GetSubcategoryListAll } from '@/api/business/biz';
import type { BasicColumn, FormSchema } from '@/components/Table';
import { CommonEnum } from '@/enums/dictTypeCode';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';

function customRenderDscps(values: number[] | null): string {
  if (!values || values.length == 0) {
    return '无';
  }
  return values.join(' / ');
}
export enum ModifyResultStatus {
  SUCCESS = 1,
  FAILED = 2,
}

// 修改DSCP值结果
export type ModifyDscpResult = {
  hostname: string;
  status: ModifyResultStatus;
  message?: string | null;
};

export const dscpAgentStatusMap = getDictDataMapFromDict(CommonEnum.ONLINE_STATUS);
export const dscpAgentStatusOptions = getSelectOptionsFromDict(CommonEnum.ONLINE_STATUS);

export const serverOnlineStatusMap = getDictDataMapFromDict(CommonEnum.ONLINE_STATUS);
export const serverOnlineStatusOptions = getSelectOptionsFromDict(CommonEnum.ONLINE_STATUS);

export const bizPriorityColumns = function (): BasicColumn[] {
  return [
    {
      title: '供应商',
      dataIndex: 'owner',
      width: 180,
      resizable: true,
    },
    {
      title: '业务类型',
      dataIndex: 'business',
      // width: 180,
      // resizable: true,
    },
    {
      title: '全局优先级',
      dataIndex: 'globalPriority',
      width: 150,
      sorter: true,
      customRender: ({ record }) => {
        return customRenderDscps(record.globalPriority);
      },
      resizable: true,
    },
    {
      title: '全局DSCP值',
      dataIndex: 'default',
      width: 150,
      sorter: true,
      customRender: ({ record }) => {
        return customRenderDscps(record.default);
      },
      resizable: true,
    },
    {
      title: '同网省内DSCP值',
      dataIndex: 'sameInner',
      width: 150,
      sorter: true,
      customRender: ({ record }) => {
        return customRenderDscps(record.sameInner);
      },
      resizable: true,
    },
    {
      title: '同网省外DSCP值',
      dataIndex: 'sameOuter',
      width: 150,
      sorter: true,
      customRender: ({ record }) => {
        return customRenderDscps(record.sameOuter);
      },
      resizable: true,
    },
    {
      title: '异网省外DSCP值',
      dataIndex: 'diffInner',
      width: 150,
      sorter: true,
      customRender: ({ record }) => {
        return customRenderDscps(record.diffInner);
      },
    },
    {
      title: '异网省外DSCP值',
      dataIndex: 'diffOuter',
      width: 150,
      sorter: true,
      customRender: ({ record }) => {
        return customRenderDscps(record.diffOuter);
      },
      resizable: true,
    },
  ];
};

export const bizPrioritySearchFormSchema = function (): FormSchema[] {
  return [
    {
      field: 'owner',
      label: '供应商',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'business',
      label: '业务类型',
      component: 'ApiSelect',
      componentProps: {
        mode: 'multiple',
        maxTagCount: 2,
        allowClear: true,
        autoClearSearchValue: false,
        api: GetSubcategoryListAll,
        params: {
          pageSize: 99999,
          pageIndex: 1,
        },
        resultField: 'items',
        valueField: 'name',
        labelField: 'name',
        filterOption: (inputValue, option: any) => {
          return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
        },
      },
      colProps: { span: 6, xl: 6 },
    },
    {
      field: 'hostname',
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
      field: 'ecdnIp',
      label: 'IP',
      component: 'InputTextArea',
      componentProps: {
        allowClear: true,
        placeholder: '输入IP进行搜索',
        rows: 3,
      },
      colProps: { span: 6 },
    },
    {
      field: 'globalPriority',
      label: '全局优先级',
      component: 'InputNumber',
      componentProps: {
        min: 0,
        max: 7,
        step: 1,
        style: {
          width: '100%',
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'dscp',
      label: 'DSCP值',
      slot: 'dscp',
      colProps: { span: 6 },
    },
    {
      field: 'dscpType',
      label: 'DSCP类型',
      component: 'Select',
      show: false,
    },
  ];
};

export const serverPriorityColumns = function (): BasicColumn[] {
  return [
    {
      title: '主机名',
      dataIndex: 'hostname',
      width: 170,
      resizable: true,
    },
    {
      title: 'IP',
      dataIndex: 'ecdnIp',
      width: 120,
      resizable: true,
      customRender: ({ record }) => {
        return record.ecdnIp ? record.ecdnIp : '无';
      },
    },
    {
      title: '全局优先级',
      dataIndex: 'globalPriority',
      width: 130,
      sorter: true,
      customRender: ({ record }) => {
        const ipv4 = record.globalPriority !== null ? record.globalPriority : '无';
        const ipv6 = record.globalPriorityV6 !== null ? record.globalPriorityV6 : '无';
        return `IPV4: ${ipv4} IPV6: ${ipv6}`;
      },
      resizable: true,
    },
    {
      title: '全局DSCP值',
      dataIndex: 'default',
      width: 130,
      sorter: true,
      customRender: ({ record }) => {
        const ipv4 = record.default !== null ? record.default : '无';
        const ipv6 = record.defaultV6 !== null ? record.defaultV6 : '无';
        return `IPV4: ${ipv4} IPV6: ${ipv6}`;
      },
      resizable: true,
    },
    {
      title: '同网省内DSCP值',
      dataIndex: 'sameInner',
      width: 140,
      sorter: true,
      customRender: ({ record }) => {
        const ipv4 = record.sameInner !== null ? record.sameInner : '无';
        const ipv6 = record.sameInnerV6 !== null ? record.sameInnerV6 : '无';
        return `IPV4: ${ipv4} IPV6: ${ipv6}`;
      },
      resizable: true,
    },
    {
      title: '同网省外DSCP值',
      dataIndex: 'sameOuter',
      width: 140,
      sorter: true,
      customRender: ({ record }) => {
        const ipv4 = record.sameOuter !== null ? record.sameOuter : '无';
        const ipv6 = record.sameOuterV6 !== null ? record.sameOuterV6 : '无';
        return `IPV4: ${ipv4} IPV6: ${ipv6}`;
      },
      resizable: true,
    },
    {
      title: '异网省外DSCP值',
      dataIndex: 'diffInner',
      width: 140,
      sorter: true,
      customRender: ({ record }) => {
        const ipv4 = record.diffInner !== null ? record.diffInner : '无';
        const ipv6 = record.diffInnerV6 !== null ? record.diffInnerV6 : '无';
        return `IPV4: ${ipv4} IPV6: ${ipv6}`;
      },
    },
    {
      title: '异网省外DSCP值',
      dataIndex: 'diffOuter',
      width: 140,
      sorter: true,
      resizable: true,
      customRender: ({ record }) => {
        const ipv4 = record.diffOuter !== null ? record.diffOuter : '无';
        const ipv6 = record.diffOuterV6 !== null ? record.diffOuterV6 : '无';
        return `IPV4: ${ipv4} IPV6: ${ipv6}`;
      },
    },
    {
      title: '设备状态',
      dataIndex: 'online',
      width: 80,
    },
    {
      title: 'DSCP Agent状态',
      dataIndex: 'status',
      width: 120,
    },
    {
      title: '上报时间',
      dataIndex: 'reportTime',
      width: 160,
    },
  ];
};

export enum modifyType {
  SUMMARY = 'summary',
  SERVER = 'serverPriority',
}

export const serverPrioritySearchFormSchema = function (): FormSchema[] {
  return [
    {
      field: 'hostname',
      label: '主机名',
      component: 'InputTextArea',
      componentProps: {
        allowClear: true,
        placeholder: '输入主机名进行搜索',
        rows: 2,
      },
      colProps: { span: 6 },
    },
    {
      field: 'ecdnIp',
      label: 'IP',
      component: 'InputTextArea',
      componentProps: {
        allowClear: true,
        placeholder: '输入IP进行搜索',
        rows: 2,
      },
      colProps: { span: 6 },
    },
    {
      field: 'globalPriority',
      label: '全局优先级',
      component: 'InputNumber',
      componentProps: {
        min: 0,
        max: 7,
        step: 1,
        style: {
          width: '100%',
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'defaultDscp',
      label: '全局优DSCP值',
      component: 'InputNumber',
      componentProps: {
        min: 0,
        max: 63,
        step: 1,
        style: {
          width: '100%',
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'status',
      label: 'DSCP Agent状态',
      component: 'Select',
      componentProps: {
        options: dscpAgentStatusOptions,
      },
      colProps: { span: 6 },
    },
    {
      field: 'online',
      label: '设备状态',
      component: 'Select',
      componentProps: {
        options: serverOnlineStatusOptions,
      },
      colProps: { span: 6 },
    },
  ];
};

export const modifyDscpSchemas = function (getFormHandleFn: any): FormSchema[] {
  return [
    {
      field: 'scope',
      label: '生效范围',
      component: 'CheckboxGroup',
      componentProps: {
        options: [
          { label: 'IPV4', value: 1 },
          { label: 'IPV6', value: 2 },
        ],
      },
      required: true,
    },
    {
      field: 'default',
      label: '全局DSCP值',
      component: 'InputNumber',
      rules: [{ validator: dscpValueValidator(getFormHandleFn), trigger: 'change' }],
      componentProps: {
        min: 0,
        max: 63,
        step: 1,
        style: {
          width: '100%',
        },
        onChange: function () {
          getFormHandleFn().validateFields([
            'default',
            'sameInner',
            'sameOuter',
            'diffInner',
            'diffOuter',
          ]);
        },
      },
    },
    {
      field: 'sameInner',
      label: '同网省内DSCP值',
      component: 'InputNumber',
      rules: [{ validator: dscpValueValidator(getFormHandleFn), trigger: 'change' }],
      componentProps: {
        min: 0,
        max: 63,
        step: 1,
        style: {
          width: '100%',
        },
        onChange: function () {
          getFormHandleFn().validateFields([
            'default',
            'sameInner',
            'sameOuter',
            'diffInner',
            'diffOuter',
          ]);
        },
      },
    },
    {
      field: 'sameOuter',
      label: '同网省外DSCP值',
      component: 'InputNumber',
      rules: [{ validator: dscpValueValidator(getFormHandleFn), trigger: 'change' }],
      componentProps: {
        min: 0,
        max: 63,
        step: 1,
        style: {
          width: '100%',
        },
        onChange: function () {
          getFormHandleFn().validateFields([
            'default',
            'sameInner',
            'sameOuter',
            'diffInner',
            'diffOuter',
          ]);
        },
      },
    },
    {
      field: 'diffInner',
      label: '异网省内DSCP值',
      component: 'InputNumber',
      rules: [{ validator: dscpValueValidator(getFormHandleFn), trigger: 'change' }],
      componentProps: {
        min: 0,
        max: 63,
        step: 1,
        style: {
          width: '100%',
        },
        onChange: function () {
          getFormHandleFn().validateFields([
            'default',
            'sameInner',
            'sameOuter',
            'diffInner',
            'diffOuter',
          ]);
        },
      },
    },
    {
      field: 'diffOuter',
      label: '异网省外DSCP值',
      component: 'InputNumber',
      rules: [{ validator: dscpValueValidator(getFormHandleFn), trigger: 'change' }],
      componentProps: {
        min: 0,
        max: 63,
        step: 1,
        style: {
          width: '100%',
        },
        onChange: function () {
          getFormHandleFn().validateFields([
            'default',
            'sameInner',
            'sameOuter',
            'diffInner',
            'diffOuter',
          ]);
        },
      },
    },
  ];
};

export const dscpValueValidator = function (getFormHandleFn: any) {
  return (_rule, _value) => {
    const { getFieldsValue } = getFormHandleFn();
    const values = getFieldsValue();
    if (
      (values.default === null || values.default === undefined) &&
      (values.sameInner === null || values.sameInner === undefined) &&
      (values.sameOuter === null || values.sameOuter === undefined) &&
      (values.diffInner === null || values.diffInner === undefined) &&
      (values.diffOuter === null || values.diffOuter === undefined)
    ) {
      return Promise.reject('至少填写一个DSCP值');
    }
    return Promise.resolve();
  };
};
