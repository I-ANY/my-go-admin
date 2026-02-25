import { BasicColumn, FormSchema } from '@/components/Table';

import { GetSubcategoryListAll } from '@/api/business/biz';
import dayjs from 'dayjs';

const machineTypeOptions = [
  { label: '物理机', value: 1 },
  { label: 'KVM', value: 2 },
  { label: '宿主机', value: 3 },
];

export const searchFormSchema: FormSchema[] = [
  {
    field: 'node',
    label: '节点',
    component: 'InputTextArea',
    colProps: { span: 6 },
    componentProps: {
      placeholder: '请输入节点，多个用换行分隔',
      rows: 3,
    },
  },
  {
    field: 'hostname',
    label: '主机名',
    component: 'InputTextArea',
    colProps: { span: 6 },
    componentProps: {
      placeholder: '请输入主机名，多个用换行分隔',
      rows: 3,
    },
  },
  {
    field: 'sn',
    label: 'SN',
    component: 'InputTextArea',
    colProps: { span: 6 },
    componentProps: {
      placeholder: '请输入SN，多个用换行分隔',
      rows: 3,
    },
  },
  {
    field: 'business',
    label: '业务名称',
    component: 'ApiSelect',
    colProps: { span: 6 },
    componentProps: {
      api: GetSubcategoryListAll,
      params: {
        pageSize: 10000,
        pageIndex: 1,
        status: 1,
      },
      labelField: 'name',
      valueField: 'name',
      mode: 'multiple',
      resultField: 'items',
      showSearch: true,
      optionFilterProp: 'label',
      filterOption: (input: string, option: any) => {
        return option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0;
      },
    },
  },
  {
    field: 'machineType',
    label: '机型',
    component: 'Select',
    componentProps: {
      options: machineTypeOptions,
    },
    colProps: { span: 6 },
  },
  {
    field: 'date',
    label: '检测日期',
    component: 'DatePicker',
    defaultValue: dayjs().format('YYYY-MM-DD'),
    componentProps: {
      format: 'YYYY-MM-DD',
      valueFormat: 'YYYY-MM-DD',
      placeholder: '请选择检测日期',
    },
    colProps: { span: 6 },
  },
];

export const columns: BasicColumn[] = [
  {
    title: '检测时间',
    dataIndex: 'collectDate',
    width: 120,
    sorter: true,
  },
  {
    title: '匹配规则',
    dataIndex: 'ruleName',
    width: 120,
  },
  {
    title: '业务组',
    dataIndex: 'businessGroup',
    width: 120,
  },
  {
    title: '业务大类',
    dataIndex: 'businessCategory',
    width: 120,
  },
  {
    title: '业务名称',
    dataIndex: 'businessName',
    width: 120,
  },
  {
    title: '节点',
    dataIndex: 'node',
    width: 180,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
  },
  {
    title: 'SN',
    dataIndex: 'sn',
    width: 180,
  },
  {
    title: '机型',
    dataIndex: 'machineType',
    width: 100,
  },
  {
    title: '检测结果',
    dataIndex: 'result',
    width: 100,
    customRender: ({ record }) => {
      return record.result;
    },
  },
  {
    title: '超配项',
    dataIndex: 'overProvisioningItem',
    width: 120,
  },
  {
    title: '标准配置',
    dataIndex: 'standard',
    width: 200,
    ellipsis: true,
  },
  {
    title: '当前配置',
    dataIndex: 'currentConfiguration',
    width: 200,
    ellipsis: true,
  },
];
