import { FormSchema } from '@/components/Table';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { SAEnum } from '@/enums/dictTypeCode';
import { formatToDateTime } from '@/utils/dateUtil';
import { RangePickPresetsExact } from '@/utils/common';
import { h } from 'vue';
import { Tag } from 'ant-design-vue';

export const utilizationVersionMap = getDictDataMapFromDict(SAEnum.SA_UTILIZATION_VERSION);
export const utilizationVersionOptions = getSelectOptionsFromDict(SAEnum.SA_UTILIZATION_VERSION);
export const utilizationSearchSchema = (onTimePikerOpen): FormSchema[] => [
  {
    field: '[start_statistical_time, end_statistical_time]',
    label: '时间',
    component: 'RangePicker',
    colProps: { span: 7 },
    componentProps: {
      allowClear: false,
      format: 'YYYY-MM-DD ',
      showTime: false,
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
    field: 'owners',
    label: '机房',
    component: 'InputTextArea',
    colProps: { span: 5 },
    componentProps: {
      allowClear: true,
      placeholder: '输入机房进行搜索',
      rows: 2,
    },
  },
  {
    field: 'hostnames',
    label: '主机名',
    component: 'InputTextArea',
    colProps: { span: 5 },
    componentProps: {
      allowClear: true,
      placeholder: '输入主机名进行搜索',
      rows: 2,
    },
  },
  {
    field: 'guids',
    label: 'GUID',
    component: 'InputTextArea',
    colProps: { span: 7 },
    componentProps: {
      allowClear: true,
      placeholder: '输入GUID进行搜索',
      rows: 2,
    },
  },
  {
    field: 'server_ids',
    label: '服务器ID',
    component: 'InputTextArea',
    colProps: { span: 7 },
    componentProps: {
      allowClear: true,
      placeholder: '输入服务器ID进行搜索',
      rows: 2,
    },
  },
  {
    field: 'versions',
    label: '业务',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      allowClear: true,
      options: utilizationVersionOptions,
    },
  },
  {
    field: 'across_province_type',
    label: '跨省类型',
    component: 'Select',
    colProps: { span: 3 },
    componentProps: {
      allowClear: true,
      options: [
        { label: '本省', value: '本省' },
        { label: '跨省', value: '跨省' },
      ],
    },
  },
  {
    field: 'across_net_type',
    label: '跨网类型',
    component: 'Select',
    colProps: { span: 3 },
    componentProps: {
      allowClear: true,
      options: [
        { label: '同网', value: '同网' },
        { label: '跨网', value: '跨网' },
      ],
    },
  },
  {
    field: 'is_limit',
    label: '是否被拉黑',
    component: 'Select',
    colProps: { span: 3 },
    componentProps: {
      allowClear: true,
      options: [
        { label: '是', value: 'true' },
        { label: '否', value: 'false' },
      ],
    },
  },
  {
    field: 'is_different',
    label: '是否异常',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      allowClear: true,
      options: [
        { label: '是', value: 0 },
        { label: '否', value: 1 },
      ],
    },
  },
];

export const utilizationTableColumns = [
  {
    title: '日期时间',
    dataIndex: 'statistical_time',
    width: 180,
    sorter: true,
    customRender: ({ record }) => {
      return formatToDateTime(record.statistical_time);
    },
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    align: 'left',
    sorter: true,
  },
  {
    title: '机房',
    dataIndex: 'owner',
    width: 120,
    sorter: true,
  },
  {
    title: '服务器ID',
    dataIndex: 'server_id',
    width: 260,
    align: 'left',
    sorter: true,
  },
  {
    title: 'GUID',
    dataIndex: 'guid',
    width: 350,
    align: 'left',
    sorter: true,
  },
  {
    title: '业务',
    dataIndex: 'version',
    width: 160,
    customRender: ({ record }) => {
      return utilizationVersionMap[record.version]?.dictLabel || record.version;
    },
  },
  {
    title: '调度运营商',
    dataIndex: 'scheduled_isp',
    width: 80,
  },
  {
    title: 'ECDN运营商',
    dataIndex: 'ecdn_isp',
    width: 100,
  },
  {
    title: '节点运营商',
    dataIndex: 'isp',
    width: 80,
  },
  {
    title: '跨省类型',
    dataIndex: 'across_province_type',
    width: 80,
  },
  {
    title: 'ECDN跨省类型',
    dataIndex: 'ecdn_across_province_type',
    width: 110,
  },
  {
    title: '跨网类型',
    dataIndex: 'across_net_type',
    width: 80,
  },
  {
    title: 'ECDN跨网类型',
    dataIndex: 'ecdn_across_net_type',
    width: 110,
  },
  {
    title: '是否异常',
    dataIndex: 'is_different',
    width: 100,
    customRender: ({ record }) => {
      if (
        record.isp_status === 0 ||
        record.across_province_type_status === 0 ||
        record.across_net_type_status === 0
      ) {
        return h(Tag, { color: 'red' }, () => '异常');
      }
      return h(Tag, { color: 'green' }, () => '正常');
    },
  },
  {
    title: '上报带宽(Gbps)',
    dataIndex: 'bw_report',
    width: 130,
    sorter: true,
  },
  {
    title: '拉黑后带宽(Gbps)',
    dataIndex: 'bw_limit',
    width: 145,
    sorter: true,
  },
  {
    title: '跑量带宽(Gbps)',
    dataIndex: 'bw_actual',
    width: 130,
    sorter: true,
  },
  {
    title: '利用率',
    dataIndex: 'utilization',
    width: 80,
    helpMessage: '利用率—拉黑后口径',
    sorter: true,
  },
  {
    title: '上报带宽利用率',
    dataIndex: 'bw_report_utilization',
    width: 130,
    helpMessage: '跑量带宽/上报带宽',
    sorter: true,
  },
  {
    title: '拉黑原因',
    dataIndex: 'reason',
    width: 500,
    align: 'left',
  },
  {
    title: '更新时间',
    dataIndex: 'updated_at',
    width: 180,
    sorter: true,
    customRender: ({ record }) => {
      return formatToDateTime(record.updated_at);
    },
  },
];

export const LimitListColumns = [
  {
    title: '日期时间',
    dataIndex: 'date',
    width: 80,
  },
  {
    title: '网卡',
    dataIndex: 'net_name',
    width: 110,
  },
  {
    title: '上报带宽(Gbps)',
    dataIndex: 'bw_report',
    width: 110,
  },
  {
    title: '拉黑带宽(Gbps)',
    dataIndex: 'bw_limit',
    width: 110,
  },
  {
    title: '跑量带宽(Gbps)',
    dataIndex: 'bw_actual',
    width: 110,
  },
  {
    title: '拉黑说明',
    dataIndex: 'description',
    width: 250,
  },
  {
    title: '拉黑原因',
    dataIndex: 'reason',
    width: 250,
  },
];
