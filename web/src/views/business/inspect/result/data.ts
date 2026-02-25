import { GetCategoryList } from '@/api/business/biz';
import { BasicColumn, FormSchema } from '@/components/Table';
import { InspectEnum } from '@/enums/dictTypeCode';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { commonCustomHeaderCell, customCeilDivide } from '@/utils/util';

export const inspectResultOptions = getSelectOptionsFromDict(InspectEnum.INSPECT_RESULT);
export const inspectResultMap = getDictDataMapFromDict(InspectEnum.INSPECT_RESULT);

export const idcTypeOptions = getSelectOptionsFromDict(InspectEnum.IDC_TYPE);
export const idcTypeResulteMap = getDictDataMapFromDict(InspectEnum.IDC_TYPE);

export const execStatusOptions = getSelectOptionsFromDict(InspectEnum.EXEC_STATUS);
export const execStatusMap = getDictDataMapFromDict(InspectEnum.EXEC_STATUS);

const statusFilterOptions = [
  { text: '正常', value: '1' },
  { text: '异常', value: '0' },
];
const checkFieldStatusOptions = [
  { label: '正常', value: '1' },
  { label: '异常', value: '0' },
];
export const searchFormSchema = (
  onIdcTypeChange,
  onCategoryIdsChange,
  checkFieldsChange,
): FormSchema[] => [
  {
    field: 'categoryIds',
    label: '业务大类',
    component: 'ApiSelect',
    componentProps: {
      api: GetCategoryList,
      resultField: 'items',
      valueField: 'id',
      labelField: 'name',
      mode: 'multiple',
      params: {
        pageSize: 99999,
        pageIndex: 1,
      },
      onChange: onCategoryIdsChange,
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'idcType',
    label: '机房类型',
    component: 'Select',
    componentProps: {
      options: idcTypeOptions,
      onChange: onIdcTypeChange,
    },
    colProps: { span: 6 },
    helpMessage: '选择“自建”则仅显示_MF结尾的业务名称，“招募”反之',
  },
  {
    field: 'business',
    label: '业务名称',
    component: 'Select',
    componentProps: {
      options: [],
      mode: 'multiple',
      maxTagCount: 3,
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
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
    field: 'owner',
    label: '节点',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'inspectResult',
    label: '巡检结果',
    component: 'Select',
    componentProps: {
      options: inspectResultOptions,
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'checkFields',
    label: '检查指标',
    component: 'Select',
    componentProps: {
      mode: 'tags',
      allowClear: true,
      maxTagCount: 3,
      onChange: checkFieldsChange,
    },
    helpMessage: '该筛选条件需要配合“检查指标状态”才能生效，多个字段为“或”的关系',
    colProps: { span: 6 },
  },
  {
    field: 'checkFieldStatus',
    label: '指标状态',
    component: 'Select',
    componentProps: {
      options: checkFieldStatusOptions,
      allowClear: true,
      disabled: true,
    },
    helpMessage: '筛选“检查指标”的状态',
    colProps: { span: 6 },
  },
  {
    field: 'status',
    label: '执行状态',
    component: 'Select',
    componentProps: {
      options: execStatusOptions,
      allowClear: true,
    },
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
    title: '执行状态',
    dataIndex: 'status',
    width: 100,
    resizable: true,
    fixed: 'left',
  },
  {
    title: '巡检结果',
    dataIndex: 'inspectResult',
    width: 80,
    resizable: true,
    fixed: 'left',
  },
  {
    title: '业务名称',
    dataIndex: 'business',
    width: 120,
    resizable: true,
  },
  {
    title: '节点',
    dataIndex: 'owner',
    width: 120,
    resizable: true,
  },
  {
    title: '晚高峰95\n利用率',
    dataIndex: 'peakUtilization95',
    width: 90,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return record.peakUtilization95 ? (record.peakUtilization95 * 100).toFixed(2) + '%' : '';
    },
    sorter: true,
  },
  // {
  //   title: '网络速率\nMbps',
  //   dataIndex: 'networkSpeed',
  //   width: 100,
  //   resizable: true,
  //   customHeaderCell: commonCustomHeaderCell(),
  //   customRender: ({ record }) => {
  //     if (
  //       record.plannedBandwidth > 0 &&
  //       record.networkSpeed > 0 &&
  //       record.networkSpeed / record.plannedBandwidth < 0.3
  //     ) {
  //       return h(
  //         'span',
  //         { style: { color: 'red' } },
  //         customCeilDivide(record.networkSpeed, 1000 * 1000, 2),
  //       );
  //     } else {
  //       return customCeilDivide(record.networkSpeed, 1000 * 1000, 2);
  //     }
  //   },
  //   sorter: true,
  // },
  {
    title: '设备带宽\nMbps',
    dataIndex: 'bandwidth',
    width: 100,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return customCeilDivide(record.bandwidth, 1000 * 1000, 1);
    },
    sorter: true,
  },
  {
    title: '规划带宽\nMbps',
    dataIndex: 'plannedBandwidth',
    width: 100,
    resizable: true,
    customHeaderCell: commonCustomHeaderCell(),
    customRender: ({ record }) => {
      return customCeilDivide(record.plannedBandwidth, 1000 * 1000, 1);
    },
    sorter: true,
  },
  {
    title: '检查型指标',
    dataIndex: 'check',
    children: [
      {
        title: '网络速率\nMbps',
        dataIndex: 'realTimeRate',
        width: 100,
        resizable: true,
        customHeaderCell: commonCustomHeaderCell(),
        filters: statusFilterOptions,
      },
      {
        title: '流量环比/\n同比检查',
        dataIndex: 'rateCompare',
        width: 100,
        resizable: true,
        customHeaderCell: commonCustomHeaderCell(),
        filters: statusFilterOptions,
      },
      {
        title: '系统负载',
        dataIndex: 'systemLoad',
        width: 95,
        resizable: true,
        filters: statusFilterOptions,
      },
      {
        title: '根分区\n使用率',
        dataIndex: 'rootPartitionUsage',
        width: 85,
        resizable: true,
        customHeaderCell: commonCustomHeaderCell(),
        filters: statusFilterOptions,
      },
      {
        title: '拨号状态',
        dataIndex: 'dialStatus',
        width: 95,
        resizable: true,
        filters: statusFilterOptions,
      },
      {
        title: 'ping丢包',
        dataIndex: 'pingPacketLoss',
        width: 95,
        resizable: true,
        filters: statusFilterOptions,
      },
      {
        title: '业务状态',
        dataIndex: 'businessStatus',
        width: 95,
        resizable: true,
        filters: statusFilterOptions,
      },
      {
        title: '业务配置',
        dataIndex: 'businessConfiguration',
        width: 95,
        resizable: true,
        filters: statusFilterOptions,
      },
      {
        title: '磁盘挂载\n状态',
        dataIndex: 'diskMountStatus',
        width: 95,
        resizable: true,
        customHeaderCell: commonCustomHeaderCell(),
        filters: statusFilterOptions,
      },
      {
        title: 'nat类型',
        dataIndex: 'natType',
        width: 95,
        resizable: true,
        filters: statusFilterOptions,
      },
      {
        title: '运营商\n一致性',
        dataIndex: 'ispConsistency',
        width: 85,
        resizable: true,
        customHeaderCell: commonCustomHeaderCell(),
        filters: statusFilterOptions,
      },
      {
        title: 'business_traffic\n状态',
        dataIndex: 'businessTrafficStatus',
        width: 130,
        resizable: true,
        customHeaderCell: commonCustomHeaderCell(),
        filters: statusFilterOptions,
      },
      {
        title: 'monitor-py\n状态',
        dataIndex: 'monitorpyStatus',
        width: 100,
        resizable: true,
        customHeaderCell: commonCustomHeaderCell(),
        filters: statusFilterOptions,
      },
      {
        title: '扩展字段',
        dataIndex: 'checkFields',
        customHeaderCell: commonCustomHeaderCell(),
        minWidth: 150,
        width: 420,
        resizable: true,
      },
    ],
  },
  {
    title: '输出型指标',
    dataIndex: 'output',
    children: [
      {
        title: '软中断\n数据包',
        dataIndex: 'lbCpu',
        width: 75,
        resizable: true,
        customHeaderCell: commonCustomHeaderCell(),
      },
      {
        title: '网卡速率',
        dataIndex: 'nicSpeed',
        width: 100,
        resizable: true,
      },
      {
        title: '内核日志',
        dataIndex: 'kernelLogs',
        width: 75,
        resizable: true,
      },
      {
        title: '业务日志',
        dataIndex: 'businessLog',
        width: 75,
        resizable: true,
      },
      {
        title: '扩展字段',
        dataIndex: 'outputFields',
        customHeaderCell: commonCustomHeaderCell(),
        minWidth: 150,
        width: 420,
        resizable: true,
      },
    ],
  },
  {
    title: '开始时间',
    dataIndex: 'startTime',
    width: 160,
    resizable: true,
  },
  {
    title: '完成时间',
    dataIndex: 'finishTime',
    width: 160,
    resizable: true,
  },
];
