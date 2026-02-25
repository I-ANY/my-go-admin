import { FormSchema } from '@/components/Form';
import { BasicColumn } from '@/components/Table';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { ZEnum } from '@/enums/dictTypeCode';
import { h } from 'vue';
import { Tag } from 'ant-design-vue';
import { formatToDateTime } from '@/utils/dateUtil';
import { RangePickPresetsExact } from '@/utils/common';
import { customCeilDivide } from '@/utils/util';

export const deviceAcceptanceStatusMap = getDictDataMapFromDict(ZEnum.ZP_ACCEPTANCE_STATUS);
export const zpDeliveryStatusMap = getDictDataMapFromDict(ZEnum.ZP_DELIVERY_STATUS);
export const zpDeliveryTypeMap = getDictDataMapFromDict(ZEnum.ZP_DELIVERY_TYPE);

export const deviceSearchFormSchema: FormSchema[] = [
  {
    field: 'hostnames',
    label: '主机名',
    component: 'InputTextArea',
    colProps: { span: 6 },
    componentProps: {
      allowClear: true,
      placeholder: '输入主机名进行搜索',
      rows: 2,
    },
  },
  {
    field: 'device_ids',
    label: '字节设备ID',
    component: 'InputTextArea',
    colProps: { span: 8 },
    componentProps: {
      allowClear: true,
      placeholder: '输入字节设备ID进行搜索',
      rows: 2,
    },
  },
  {
    field: 'acceptance_status',
    label: '验收状态',
    component: 'Select',
    colProps: { span: 5 },
    componentProps: {
      mode: 'multiple',
      options: getSelectOptionsFromDict(ZEnum.ZP_ACCEPTANCE_STATUS),
    },
  },
  {
    field: 'real_time_status',
    label: '实时状态',
    component: 'Select',
    colProps: { span: 5 },
    componentProps: {
      mode: 'multiple',
      options: [
        { label: '正常', value: '正常' },
        { label: '异常', value: '异常' },
        { label: '离线', value: '离线' },
      ],
    },
  },
  {
    field: 'is_bw_equal',
    label: '带宽是否一致',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      options: [
        { label: '是', value: 'true' },
        { label: '否', value: 'false' },
      ],
    },
  },
  {
    field: 'isp',
    label: '运营商',
    component: 'Select',
    componentProps: {
      options: [
        { label: '移动', value: '移动' },
        { label: '联通', value: '联通' },
        { label: '电信', value: '电信' },
      ],
      allowClear: true,
    },
    colProps: { span: 4 },
  },
  {
    field: 'isp_status',
    label: '运营商是否一致',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      options: [
        { label: '是', value: 1 },
        { label: '否', value: 0 },
      ],
    },
  },
  {
    field: 'province',
    label: '省份',
    component: 'Input',
    colProps: { span: 6 },
  },
];
export const deviceColumns: BasicColumn[] = [
  {
    title: '字节设备ID',
    dataIndex: 'device_id',
    width: 300,
    resizable: true,
    align: 'left',
  },
  {
    title: '主机名',
    dataIndex: 'provider_device_id',
    width: 200,
    resizable: true,
    align: 'left',
  },
  {
    title: '设备类型',
    dataIndex: 'type',
    width: 80,
    resizable: true,
  },
  {
    title: '省份(字节)',
    dataIndex: 'province',
    width: 100,
    resizable: true,
  },
  {
    title: '运营商(字节)',
    dataIndex: 'isp',
    width: 100,
    resizable: true,
  },
  {
    title: '省份(ECDN)',
    dataIndex: 'ecdn_province',
    width: 100,
    resizable: true,
  },
  {
    title: '运营商(ECDN)',
    dataIndex: 'ecdn_isp',
    width: 100,
    resizable: true,
  },
  {
    title: '验收状态',
    dataIndex: 'acceptance_status',
    width: 120,
    resizable: true,
    customRender: ({ record }) => {
      const status = record.acceptance_status;
      const color = deviceAcceptanceStatusMap[status].color || 'default';
      return h(Tag, { color: color }, () => status);
    },
  },
  {
    title: '实时状态',
    dataIndex: 'real_time_status',
    width: 80,
    resizable: true,
    customRender: ({ record }) => {
      const status = record.real_time_status;
      const color = status === '正常' ? '#2eb438' : '#FF1100';
      return h(Tag, { color: color }, () => status);
    },
  },
  {
    title: '运营商是否一致',
    dataIndex: 'isp_status',
    width: 80,
    resizable: true,
    customRender: ({ record }) => {
      const status = record.isp_status === 1 ? '是' : '否';
      const color = status === '是' ? '#2eb438' : '#FF1100';
      return h(Tag, { color: color }, () => status);
    },
  },
  {
    title: '真实带宽(Mbps)',
    dataIndex: 'real_bw',
    width: 100,
    resizable: true,
    helpMessage: '设备真实上报的带宽，单位为Mbps',
    customRender: ({ text, record }) => {
      if (record.real_bw != null) {
        const isDiff = record.real_bw !== record.total_bandwidth;
        return h('span', { style: isDiff ? { color: '#FF1100' } : {} }, text);
      }
    },
  },
  {
    title: '总带宽(Mbps)',
    dataIndex: 'total_bandwidth',
    width: 100,
    resizable: true,
    customRender: ({ text, record }) => {
      if (record.real_bw != null) {
        const isDiff = record.real_bw !== record.total_bandwidth;
        return h('span', { style: isDiff ? { color: '#FF1100' } : {} }, text);
      }
    },
  },
  {
    title: 'CPU',
    dataIndex: 'cpu',
    width: 100,
    resizable: true,
  },
  {
    title: '内存',
    dataIndex: 'memory',
    width: 100,
    resizable: true,
  },
  {
    title: '总磁盘容量',
    dataIndex: 'total_disk_capacity',
    width: 100,
    resizable: true,
  },
  {
    title: '注册时间',
    dataIndex: 'created_at',
    width: 180,
    resizable: true,
    customRender: ({ text }) => {
      return formatToDateTime(text); // 根据你的工具函数调整参数
    },
  },
  {
    title: '更新时间',
    dataIndex: 'updated_at',
    width: 180,
    resizable: true,
    customRender: ({ text }) => {
      return formatToDateTime(text); // 根据你的工具函数调整参数
    },
  },
];

export const trafficSearchFormSchema = (onTimePikerOpen): FormSchema[] => [
  {
    field: '[startTimestmp, endTimestmp]',
    label: '时间',
    component: 'RangePicker',
    colProps: { span: 8 },
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
  },
  {
    field: 'hostnames',
    label: '主机名',
    component: 'InputTextArea',
    colProps: { span: 8 },
    componentProps: {
      allowClear: true,
      placeholder: '输入主机名进行搜索',
      rows: 2,
    },
  },
  {
    field: 'z_device_ids',
    label: '字节设备ID',
    component: 'InputTextArea',
    colProps: { span: 8 },
    componentProps: {
      allowClear: true,
      placeholder: '输入字节设备ID进行搜索',
      rows: 2,
    },
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    colProps: { span: 8 },
    componentProps: {
      options: [
        { label: '正常', value: 1 },
        { label: '异常', value: 0 },
        { label: '无数据', value: 2 },
      ],
    },
  },
  {
    field: 'is_cmcc',
    label: '是否电联',
    component: 'Select',
    colProps: { span: 8 },
    componentProps: {
      options: [
        { label: '是', value: 0 },
        { label: '否', value: 1 },
      ],
    },
  },
];

const statusMap = {
  1: {
    text: '正常',
    color: 'green',
  },
  0: {
    text: '异常',
    color: 'red',
  },
  2: {
    text: '无数据',
    color: 'yellow',
  },
};

export const traffic95SearchFormSchema = (): FormSchema[] => [
  {
    field: '[start_date, end_date]',
    label: '日期范围',
    component: 'RangePicker',
    colProps: { span: 8 },
    componentProps: {
      format: 'YYYY-MM-DD',
      placeholder: ['开始日期', '结束日期'],
    },
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    colProps: { span: 6 },
    componentProps: {
      options: [
        { label: '正常', value: 1 },
        { label: '异常', value: 0 },
      ],
    },
  },
];

// 流量信息字段
export const trafficColumns: BasicColumn[] = [
  {
    title: '时间',
    dataIndex: 'timestmp',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      return formatToDateTime(record.timestmp);
    },
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 100,
    resizable: true,
  },
  {
    title: '字节设备ID',
    dataIndex: 'z_device_id',
    width: 160,
    resizable: true,
  },
  {
    title: '省份(ECDN)',
    dataIndex: 'province',
    width: 50,
    resizable: true,
  },
  {
    title: '省份(字节)',
    dataIndex: 'z_province',
    width: 50,
    resizable: true,
  },
  {
    title: '运营商(ECDN)',
    dataIndex: 'isp',
    width: 60,
    resizable: true,
  },
  {
    title: '运营商(字节)',
    dataIndex: 'z_isp',
    width: 60,
    resizable: true,
  },
  {
    title: '是否电联',
    dataIndex: 'is_cmcc',
    width: 50,
    resizable: true,
    customRender: ({ record }) => {
      const status = record.is_cmcc ? '否' : '是';
      const color = status === '是' ? 'green' : 'red';
      return h(Tag, { color: color }, () => status);
    },
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 50,
    resizable: true,
    customRender: ({ record }) => {
      const status = record.status;
      const color = statusMap[status]['color'];
      const text = statusMap[status]['text'];
      return h(Tag, { color: color }, () => text);
    },
  },
  {
    title: '网络速率(Mbps)',
    dataIndex: 'netsentspeed',
    width: 80,
    resizable: true,
    customRender: ({ record }) => {
      return customCeilDivide(record.netsentspeed);
    },
  },
];

export const traffic95Columns: BasicColumn[] = [
  {
    title: '日期',
    dataIndex: 'date',
    width: 80,
    resizable: true,
  },
  {
    title: '原运营商',
    dataIndex: 'isp',
    width: 50,
    resizable: true,
  },
  {
    title: 'ecdn95时间',
    dataIndex: 'ecdn_timestmp',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      return formatToDateTime(record.ecdn_timestmp);
    },
  },
  {
    title: 'ecdn95值(Gbps)',
    dataIndex: 'ecdn_total',
    width: 80,
    resizable: true,
    customRender: ({ record }) => {
      return customCeilDivide(record.ecdn_total, 1000 * 1000 * 1000);
    },
  },
  {
    title: '字节95时间',
    dataIndex: 'z_timestmp',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      return formatToDateTime(record.z_timestmp);
    },
  },
  {
    title: '字节95值(Gbps)',
    dataIndex: 'z_total',
    width: 80,
    resizable: true,
    customRender: ({ record }) => {
      return customCeilDivide(record.z_total, 1000 * 1000 * 1000);
    },
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 80,
    resizable: true,
    customRender: ({ record }) => {
      const status = record.status;
      const color = statusMap[status]['color'];
      const text = statusMap[status]['text'];
      return h(Tag, { color: color }, () => text);
    },
  },
  {
    title: '发布状态',
    dataIndex: 'is_publish',
    width: 80,
    resizable: true,
  },
];

export const deliveryDeviceSearchFormSchema = (): FormSchema[] => [
  {
    field: 'provider_device_ids',
    label: '主机名',
    component: 'InputTextArea',
    colProps: { span: 5 },
    componentProps: {
      rows: 2,
      placeholder: '输入主机名进行搜索',
      allowClear: true,
    },
  },
  {
    field: 'device_ids',
    label: '字节设备ID',
    component: 'InputTextArea',
    colProps: { span: 6 },
    componentProps: {
      rows: 2,
      placeholder: '输入字节设备ID进行搜索',
      allowClear: true,
    },
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      placeholder: '请选择',
      allowClear: true,
      options: getSelectOptionsFromDict(ZEnum.ZP_DELIVERY_STATUS),
    },
  },
  {
    field: 'type',
    label: '类型',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      placeholder: '请选择',
      allowClear: true,
      options: getSelectOptionsFromDict(ZEnum.ZP_DELIVERY_TYPE),
    },
  },
];

export const deliveryDeviceColumns: BasicColumn[] = [
  {
    title: '主机名',
    dataIndex: 'provider_device_id',
    width: 170,
    resizable: true,
  },
  {
    title: '字节设备ID',
    dataIndex: 'device_id',
    width: 300,
    resizable: true,
  },
  {
    title: '类型',
    dataIndex: 'type',
    width: 80,
    resizable: true,
    customRender: ({ record }) => {
      const status = record.type;
      return zpDeliveryTypeMap[status].dictLabel;
    },
  },
  {
    title: '省份',
    dataIndex: 'province',
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
    title: '交付运营商',
    dataIndex: 'submit_isp',
    width: 80,
    resizable: true,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 80,
    resizable: true,
    fixed: 'right',
    customRender: ({ record }) => {
      const status = record.status;
      const color = zpDeliveryStatusMap[status]['color'] || 'default';
      const text = zpDeliveryStatusMap[status].dictLabel;
      return h(Tag, { color: color }, () => text);
    },
  },
  {
    title: '网卡信息',
    dataIndex: 'network',
    width: 350,
    resizable: true,
    align: 'left',
  },
  {
    title: '磁盘信息',
    dataIndex: 'storage',
    width: 350,
    resizable: true,
    align: 'left',
  },
  {
    title: '操作账号',
    dataIndex: 'operation_account',
    width: 80,
    resizable: true,
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    width: 150,
    resizable: true,
    customRender: ({ record }) => {
      return formatToDateTime(record.created_at);
    },
  },
  {
    title: '更新时间',
    dataIndex: 'updated_at',
    width: 150,
    resizable: true,
    customRender: ({ record }) => {
      return formatToDateTime(record.updated_at);
    },
  },
];

export const deliveryInfoSchemas: FormSchema[] = [
  {
    label: 'id',
    field: 'id',
    component: 'Input',
    colProps: {
      span: 22,
    },
    required: true,
    show: false,
    componentProps: {
      disabled: true,
      readonly: true,
    },
  },
  {
    label: '主机名',
    field: 'provider_device_id',
    component: 'Input',
    colProps: {
      span: 22,
    },
    required: true,
    show: true,
    componentProps: {
      readonly: true,
    },
  },
  {
    label: '字节设备ID',
    field: 'device_id',
    component: 'Input',
    colProps: {
      span: 22,
    },
    required: true,
    show: true,
    componentProps: {
      readonly: true,
    },
  },
  {
    label: '类型',
    field: 'type',
    component: 'Input',
    colProps: {
      span: 22,
    },
    required: true,
    show: true,
    componentProps: {
      readonly: true,
    },
  },
  {
    label: '省份',
    field: 'province',
    component: 'Input',
    colProps: {
      span: 22,
    },
    required: true,
    show: true,
    componentProps: {},
  },
  {
    label: '原运营商',
    field: 'isp',
    component: 'Input',
    colProps: {
      span: 8,
    },
    required: true,
    show: true,
    componentProps: {
      readonly: true,
    },
  },
  {
    label: '交付运营商',
    field: 'submit_isp',
    component: 'Select',
    colProps: {
      span: 8,
    },
    required: true,
    show: true,
    componentProps: {
      options: [
        {
          label: '移动',
          value: '移动',
        },
        {
          label: '联通',
          value: '联通',
        },
        {
          label: '电信',
          value: '电信',
        },
      ],
    },
  },
  {
    label: '磁盘信息',
    field: 'storage',
    component: 'InputTextArea',
    colProps: {
      span: 22,
    },
    required: true,
    show: true,
    helpMessage:
      '[\n' +
      '  {\n' +
      '    "name": "/dev/vdb",\n' +
      '    "capacity": 649,\n' +
      '    "type": "ssd"\n' +
      '  },\n...\n]\n',
    componentProps: {
      rows: 5,
    },
  },
  {
    label: '线路数',
    field: 'lines',
    component: 'Input',
    colProps: {
      span: 8,
    },
    required: true,
    show: true,
    componentProps: {},
  },
  {
    label: '单线路带宽',
    field: 'single_line_bw',
    component: 'Input',
    colProps: {
      span: 8,
    },
    required: true,
    show: true,
    componentProps: {},
  },
  {
    label: '当前网卡信息',
    field: 'network',
    component: 'InputTextArea',
    colProps: {
      span: 22,
    },
    required: true,
    show: true,
    componentProps: {
      rows: 5,
    },
  },
  {
    label: '状态',
    field: 'status',
    component: 'Select',
    colProps: {
      span: 22,
    },
    required: true,
    show: true,
    componentProps: {
      options: getSelectOptionsFromDict(ZEnum.ZP_DELIVERY_STATUS),
    },
  },
];

export const deliveryBatchEditSchemas: FormSchema[] = [
  {
    label: '交付运营商',
    field: 'submit_isp',
    component: 'Select',
    colProps: {
      span: 22,
    },
    show: true,
    componentProps: {
      options: [
        {
          label: '移动',
          value: '移动',
        },
        {
          label: '联通',
          value: '联通',
        },
        {
          label: '电信',
          value: '电信',
        },
      ],
    },
  },
  {
    label: '单线路带宽',
    field: 'single_line_bw',
    component: 'Input',
    colProps: {
      span: 22,
    },
    show: true,
    componentProps: {},
  },
  {
    label: '交付状态',
    field: 'status',
    component: 'Select',
    colProps: {
      span: 22,
    },
    show: true,
    componentProps: {
      options: getSelectOptionsFromDict(ZEnum.ZP_DELIVERY_STATUS),
    },
  },
];
