import { BasicColumn } from '@/components/Table';
import { FormSchema } from '@/components/Form';
import { Tag } from 'ant-design-vue';
import { formatToDateTime } from '@/utils/dateUtil';
import { h } from 'vue';

// 状态选项
export const statusOptions = [
  { label: '待审核', value: 0 },
  { label: '已上线', value: 1 },
  { label: '已审核', value: 4 },
  { label: '未通过', value: 5 },
  { label: '休眠', value: 6 },
];

// 在线状态选项
export const onlineOptions = [
  { label: '离线', value: 0 },
  { label: '在线', value: 1 },
];

// 运营商选项
export const ispOptions = [
  { label: '移动', value: '移动' },
  { label: '联通', value: '联通' },
  { label: '电信', value: '电信' },
];

// 机房类型选项
export const roomTypeOptions = [
  { label: 'IDC', value: 1 },
  { label: 'ACDN', value: 2 },
  { label: 'MCDN', value: 3 },
];

// 可跨省选项
export const interprovincialOptions = [
  { label: '未知', value: 0 },
  { label: '是', value: 1 },
  { label: '否', value: 2 },
];

// 机房归属选项
export const originOptions = [
  { label: '自建', value: 1 },
  { label: '招募', value: 2 },
];

// 计费类型选项
export const chargeModeOptions = [
  { label: '买断', value: 1 },
  { label: '95', value: 2 },
  { label: '单机95', value: 3 },
];

// 搜索表单配置
export const searchFormSchema: FormSchema[] = [
  {
    field: 'query',
    label: ' ',
    component: 'InputTextArea',
    colProps: { span: 8 },
    componentProps: {
      allowClear: true,
      placeholder: '输入FrankID或serverID或sn，多个换行，一行一个',
    },
  },
  {
    field: 'hostnames',
    label: ' ',
    component: 'InputTextArea',
    colProps: { span: 6 },
    componentProps: {
      allowClear: true,
      placeholder: '输入主机名，多个换行，一行一个',
      rows: 2,
    },
  },
  {
    field: 'owner',
    label: ' ',
    component: 'Input',
    colProps: { span: 3 },
    componentProps: {
      allowClear: true,
      placeholder: '输入机房',
    },
  },
  {
    field: 'status',
    label: ' ',
    component: 'Select',
    colProps: { span: 5 },
    componentProps: {
      options: statusOptions,
      allowClear: true,
      mode: 'multiple',
      placeholder: '审批状态',
      maxTagCount: 2,
    },
  },
  {
    field: 'online',
    label: ' ',
    component: 'Select',
    colProps: { span: 2 },
    componentProps: {
      options: onlineOptions,
      allowClear: true,
      placeholder: '设备状态',
    },
  },
  {
    field: 'business',
    label: ' ',
    component: 'Select',
    colProps: { span: 8 },
    componentProps: {
      options: [],
      allowClear: true,
      mode: 'multiple',
      placeholder: '输入业务名称',
      maxTagCount: 5,
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
  },
  {
    field: 'isp',
    label: ' ',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      options: ispOptions,
      allowClear: true,
      mode: 'multiple',
      placeholder: '选择运营商',
      maxTagCount: 2,
    },
  },
  {
    field: 'province',
    label: ' ',
    component: 'Input',
    colProps: { span: 4 },
    componentProps: {
      allowClear: true,
      placeholder: '输入省份',
    },
  },
];

// 表格列配置
export const columns: BasicColumn[] = [
  {
    title: '机房',
    dataIndex: 'owner',
    width: 100,
    resizable: true,
  },
  {
    title: '机房类型',
    dataIndex: 'roomType',
    width: 90,
    resizable: true,
    customRender: ({ record }) => {
      const roomType = record.roomType;
      if (roomType === 1) {
        return h(Tag, { color: 'blue' }, () => 'IDC');
      } else if (roomType === 2) {
        return h(Tag, { color: 'green' }, () => 'ACDN');
      } else if (roomType === 3) {
        return h(Tag, { color: 'purple' }, () => 'MCDN');
      }
      return roomType;
    },
  },
  {
    title: '服务器ID',
    dataIndex: 'serverID',
    width: 280,
    resizable: true,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
  },
  {
    title: '宿主机',
    dataIndex: 'parent',
    width: 180,
    resizable: true,
    customRender: ({ record }) => {
      return record.parent || '-';
    },
  },
  {
    title: 'FrankID',
    dataIndex: 'frankID',
    width: 400,
    resizable: true,
    align: 'left',
    customRender: ({ record }) => {
      return record.frankID || '-';
    },
  },
  {
    title: '审批状态',
    dataIndex: 'status',
    width: 90,
    resizable: true,
    customRender: ({ record }) => {
      const status = record.status;
      if (status === 0) {
        return h(Tag, { color: 'default' }, () => '待审核');
      } else if (status === 1) {
        return h(Tag, { color: 'green' }, () => '已上线');
      } else if (status === 4) {
        return h(Tag, { color: 'orange' }, () => '已审核');
      } else if (status === 5) {
        return h(Tag, { color: 'red' }, () => '未通过');
      } else if (status === 6) {
        return h(Tag, { color: 'red' }, () => '休眠');
      }
      return status;
    },
  },
  {
    title: '在线状态',
    dataIndex: 'online',
    width: 90,
    resizable: true,
    customRender: ({ record }) => {
      const online = record.online;
      if (online === 1) {
        return h(Tag, { color: 'green' }, () => '在线');
      }
      return h(Tag, { color: 'red' }, () => '离线');
    },
  },
  {
    title: '业务名称',
    dataIndex: 'business',
    width: 120,
    resizable: true,
  },
  {
    title: 'IP地址',
    dataIndex: 'ip',
    width: 140,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'carrier',
    width: 80,
    resizable: true,
  },
  {
    title: '总带宽',
    dataIndex: 'bwTotal',
    width: 90,
    resizable: true,
    customRender: ({ record }) => {
      return record.bwTotal / 1000000000;
    },
  },
  {
    title: '规划带宽',
    dataIndex: 'bwPlan',
    width: 90,
    resizable: true,
    customRender: ({ record }) => {
      return record.bwPlan / 1000000000;
    },
  },
  {
    title: '线路数',
    dataIndex: 'bwCount',
    width: 80,
    resizable: true,
  },
  {
    title: '单线路带宽',
    dataIndex: 'bwSingle',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      return record.bwSingle / 1000000000;
    },
  },
  //   {
  //     title: '当前速率(bps)',
  //     dataIndex: 'speedNow',
  //     width: 120,
  //     resizable: true,
  //   },
  //   {
  //     title: '计费类型',
  //     dataIndex: 'chargeMode',
  //     width: 90,
  //     resizable: true,
  //     customRender: ({ record }) => {
  //       const mode = record.chargeMode;
  //       if (mode === 1) {
  //         return h(Tag, { color: 'green' }, () => '买断');
  //       } else if (mode === 2) {
  //         return h(Tag, { color: 'blue' }, () => '95');
  //       } else if (mode === 3) {
  //         return h(Tag, { color: 'purple' }, () => '单机95');
  //       }
  //       return '-';
  //     },
  //   },
  //   {
  //     title: '业务ID',
  //     dataIndex: 'businessId',
  //     width: 80,
  //     resizable: true,
  //   },
  {
    title: '地区',
    dataIndex: 'location',
    width: 100,
    resizable: true,
  },
  {
    title: '省份',
    dataIndex: 'province',
    width: 80,
    resizable: true,
  },
  {
    title: '跨省标签',
    dataIndex: 'provStatus',
    width: 100,
    resizable: true,
  },
  {
    title: '是否可跨省',
    dataIndex: 'isInterprovincial',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      const val = record.isInterprovincial;
      if (val === 1) {
        return h(Tag, { color: 'green' }, () => '是');
      } else if (val === 2) {
        return h(Tag, { color: 'red' }, () => '否');
      }
      return h(Tag, { color: 'default' }, () => '未知');
    },
  },
  {
    title: '异网标签',
    dataIndex: 'reteMirabileStatus',
    width: 100,
    resizable: true,
  },
  {
    title: '日95利用率',
    dataIndex: 'allDay95Utilization',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      const val = record.allDay95Utilization;
      if (val !== null && val !== undefined) {
        return `${(val * 100).toFixed(2)}%`;
      }
      return '-';
    },
  },
  {
    title: '晚高峰95利用率',
    dataIndex: 'eveningPeak95Utilization',
    width: 120,
    resizable: true,
    customRender: ({ record }) => {
      const val = record.eveningPeak95Utilization;
      if (val !== null && val !== undefined) {
        return `${(val * 100).toFixed(2)}%`;
      }
      return '-';
    },
  },
  {
    title: '业务流状态',
    dataIndex: 'deployStatus',
    width: 100,
    resizable: true,
  },
  {
    title: '审核时间',
    dataIndex: 'timepassed',
    width: 160,
    resizable: true,
    customRender: ({ record }) => {
      return record.timepassed ? formatToDateTime(record.timepassed) : '-';
    },
  },
  {
    title: '机房归属',
    dataIndex: 'origin',
    width: 90,
    resizable: true,
    customRender: ({ record }) => {
      const origin = record.origin;
      if (origin === 1) {
        return h(Tag, { color: 'blue' }, () => '自建');
      } else if (origin === 2) {
        return h(Tag, { color: 'orange' }, () => '招募');
      }
      return '-';
    },
  },
  {
    title: '监控标签',
    dataIndex: 'cactiNotes',
    width: 200,
    resizable: true,
    align: 'left',
  },
  {
    title: 'SN',
    dataIndex: 'sn',
    width: 180,
    resizable: true,
    align: 'left',
  },
  {
    title: '备注',
    dataIndex: 'notes',
    width: 120,
    resizable: true,
    ellipsis: true,
    align: 'left',
  },
];
