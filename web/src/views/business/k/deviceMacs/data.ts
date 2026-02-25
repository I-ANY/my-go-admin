import { BasicColumn, FormSchema } from '@/components/Table';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';

import { KEnum } from '@/enums/dictTypeCode';
import { commonCustomHeaderCell } from '@/utils/util';
import dayjs from 'dayjs';

export const bizTypesMap = getDictDataMapFromDict(KEnum.BIZ_TYPE);
export const providersMap = getDictDataMapFromDict(KEnum.PRODIVER);
export const isProvinceSchedulingMap = getDictDataMapFromDict(KEnum.SCHEDULE_TYPE);
export const isCoverDiffIspMap = getDictDataMapFromDict(KEnum.IS_COVER_DIFF_ISP);
export const coverDiffIspMap = getDictDataMapFromDict(KEnum.ISP);
export const isDisabledMap = getDictDataMapFromDict(KEnum.IS_DISABLED);
export const opStatusMap = getDictDataMapFromDict(KEnum.OP_STATUS);
export const BusinessStatusMap = getDictDataMapFromDict(KEnum.BUSINESS_STATUS);

function handleOpStatusOptions(data: any): [] {
  for (let i = data.length - 1; i >= 0; i--) {
    const item = data[i];
    if (item.value && isNaN(Number(item.value))) {
      data.splice(i, 1);
    }
  }
  return data;
}

export const areaSelectOptions = [
  { label: '全国', value: '全国' },
  { label: '东北', value: '东北' },
  { label: '华北', value: '华北' },
  { label: '华南', value: '华南' },
  { label: '华东', value: '华东' },
  { label: '西北', value: '西北' },
  { label: '西南', value: '西南' },
];

export function MacSearchFormSchema(bizType: string): FormSchema[] {
  return [
    {
      field: 'nodes',
      label: '节点名',
      component: 'InputTextArea',
      componentProps: {
        allowClear: true,
        placeholder: '输入节点名进行搜索',
        rows: 2,
      },
      colProps: { span: 4 },
    },
    {
      field: 'hostnames',
      label: '主机名',
      component: 'InputTextArea',
      componentProps: {
        allowClear: true,
        placeholder: '输入主机名进行搜索',
        rows: 2,
      },
      colProps: { span: 4 },
    },
    {
      field: 'macAddrs',
      label: 'mac地址',
      component: 'InputTextArea',
      colProps: { span: 4 },
      componentProps: {
        allowClear: true,
        placeholder: '输入MAC进行搜索',
        rows: 2,
      },
    },
    {
      field: 'business',
      label: '业务名',
      component: 'ApiSelect',
      componentProps: {
        options: [],
        mode: 'multiple',
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'dev_type_names',
      label: '设备类型',
      component: 'ApiSelect',
      componentProps: {
        options: [],
        mode: 'multiple',
        allowClear: true,
      },
      colProps: { span: 4 },
    },

    {
      field: 'wan_ip',
      label: '外网IP',
      component: 'Input',
      colProps: { span: 4 },
    },
    {
      field: 'is_province_scheduling',
      label: '调度控制',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.SCHEDULE_TYPE),
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'is_cover_diff_isp',
      label: '是否异网设备',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.IS_COVER_DIFF_ISP),
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'is_independent_deploy',
      label: '是否云平台',
      component: 'Select',
      componentProps: {
        options: [
          { label: '否', value: '0' },
          { label: '是', value: '1' },
        ],
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'is_only_ipv6',
      label: '是否IPV6设备',
      component: 'Select',
      componentProps: {
        options: [
          { label: '否', value: '0' },
          { label: '是', value: '1' },
        ],
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'is_pass',
      label: '是否已入库',
      component: 'Select',
      componentProps: {
        options: [
          { label: '否', value: '0' },
          { label: '是', value: '1' },
        ],
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      label: '是否禁用',
      field: 'is_disabled',
      component: 'Select',
      colProps: { span: 4 },
      componentProps: {
        options: [
          { label: '否', value: '0' },
          { label: '是', value: '1' },
        ],
        allowClear: true,
      },
    },
    {
      field: 'op_status',
      label: '运营状态',
      component: 'Select',
      colProps: { span: 4 },
      componentProps: {
        options: handleOpStatusOptions(getSelectOptionsFromDict(KEnum.OP_STATUS)),
        allowClear: true,
      },
    },
    {
      field: 'business_status',
      label: '业务状态',
      component: 'Select',
      colProps: { span: 4 },
      componentProps: {
        options: handleOpStatusOptions(getSelectOptionsFromDict(KEnum.BUSINESS_STATUS)),
        allowClear: true,
      },
    },
    {
      field: 'area_name',
      label: '区域',
      component: 'ApiSelect',
      colProps: {
        span: 4,
      },
      componentProps: {
        options: areaSelectOptions,
        mode: 'multiple',
        allowClear: true,
      },
    },
    {
      field: 'province_names',
      label: '省份',
      component: 'ApiSelect',
      colProps: {
        span: 4,
      },
      componentProps: {
        options: [],
        mode: 'multiple',
        allowClear: true,
      },
    },
    {
      field: 'isp_name',
      label: '运营商',
      component: 'Select',
      colProps: { span: 4 },
      componentProps: {
        options: [
          { label: '电信', value: '电信' },
          { label: '联通', value: '联通' },
          { label: '移动', value: '移动' },
        ],
      },
    },
    {
      field: 'provider',
      label: '厂商',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.PRODIVER),
        allowClear: true,
      },
      colProps: { span: 4 },
      ifShow: () => bizType === 'specialLine',
    },
    {
      field: 'flow_upload_tx_status',
      label: '上报腾讯状态',
      component: 'Select',
      colProps: { span: 4 },
      componentProps: {
        options: [
          { label: '异常', value: '0' },
          { label: '正常', value: '1' },
        ],
        allowClear: true,
      },
    },
    {
      field: 'flow_upload_autoops_status',
      label: '上报autoops状态',
      component: 'Select',
      colProps: { span: 4 },
      componentProps: {
        options: [
          { label: '异常', value: '0' },
          { label: '正常', value: '1' },
        ],
        allowClear: true,
      },
    },
    {
      field: 'offline_time',
      label: '业务离线时长',
      slot: 'offlineTime',
      colProps: {
        span: 8,
      },
    },
    {
      field: 'is_first_mac',
      label: '是否首Mac',
      component: 'Select',
      componentProps: {
        options: [
          { label: '否', value: 0 },
          { label: '是', value: 1 },
        ],
        allowClear: true,
      },
      colProps: {
        span: 4,
      },
    },
  ];
}

export function normalColumns(bizType: string): BasicColumn[] {
  return [
    {
      title: 'mac地址',
      dataIndex: 'mac_addr',
      width: 160,
      resizable: true,
      fixed: 'left',
    },
    {
      title: '首mac',
      dataIndex: 'is_first_mac',
      width: 80,
      resizable: true,
    },
    {
      title: '厂商',
      dataIndex: 'provider',
      width: 80,
      resizable: true,
      ifShow: () => bizType === 'specialLine',
      customRender: ({ text }) => {
        const providerMap = {
          mf: '明赋',
          hn: '泓宁',
          hnk: '泓宁',
        };
        return providerMap[text] || text;
      },
    },
    {
      title: '节点/机房',
      dataIndex: 'node',
      width: 100,
      resizable: true,
    },
    {
      title: '主机名',
      dataIndex: 'hostname',
      width: 180,
      resizable: true,
    },
    {
      title: '业务名',
      dataIndex: 'business',
      width: 120,
      resizable: true,
    },
    {
      title: '设备类型',
      dataIndex: 'dev_type_name',
      width: 80,
      resizable: true,
    },
    {
      title: '运营商',
      dataIndex: 'isp_name',
      width: 80,
      resizable: true,
    },
    {
      title: '区域',
      dataIndex: 'area_name',
      width: 80,
      resizable: true,
    },
    {
      title: '省份',
      dataIndex: 'province_name',
      width: 100,
      resizable: true,
    },
    {
      title: '城市',
      dataIndex: 'city_name',
      width: 100,
      resizable: true,
    },
    {
      title: '是否已入库',
      dataIndex: 'is_pass',
      width: 100,
      resizable: true,
    },
    {
      title: '运营状态',
      dataIndex: 'op_status',
      width: 80,
      resizable: true,
    },
    {
      title: '业务状态',
      dataIndex: 'business_status',
      width: 80,
      resizable: true,
    },
    {
      title: '流量上报\n腾讯状态',
      dataIndex: 'flow_upload_tx_status',
      width: 80,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '流量上报\nautoops状态',
      dataIndex: 'flow_upload_autoops_status',
      width: 100,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '带宽(Mbps)',
      dataIndex: 'link_upload_speed',
      width: 100,
      resizable: true,
    },
    {
      title: '外网IP',
      dataIndex: 'wan_ip',
      width: 120,
      resizable: true,
    },
    {
      title: '历史使用过的IP',
      dataIndex: 'ip_history',
      width: 200,
      resizable: true,
    },
    {
      title: '调度控制',
      dataIndex: 'is_province_scheduling',
      width: 120,
      resizable: true,
    },
    {
      title: '是否云平台',
      dataIndex: 'is_independent_deploy',
      width: 100,
      resizable: true,
    },
    {
      title: '是否IPV6设备',
      dataIndex: 'is_only_ipv6',
      width: 100,
      resizable: true,
    },
    {
      title: '是否异网',
      dataIndex: 'is_cover_diff_isp',
      width: 100,
      resizable: true,
    },
    {
      title: '异网覆盖\n供应商',
      dataIndex: 'cover_diff_isp_id',
      width: 100,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '入库时间',
      dataIndex: 'pass_time',
      width: 160,
      resizable: true,
    },
    {
      title: '最后在线时间',
      dataIndex: 'last_online_time',
      width: 160,
      resizable: true,
    },
    {
      title: '离线时长(天)',
      dataIndex: 'offline_time',
      width: 100,
      resizable: true,
      customRender: ({ record }) => {
        const days = record.offline_time / 24;
        const displayValue = days.toFixed(2);
        return displayValue;
      },
    },
    {
      title: '是否禁用',
      dataIndex: 'is_disabled',
      width: 80,
      resizable: true,
    },
    {
      title: '禁用原因',
      dataIndex: 'disabled_reason',
      width: 120,
      resizable: true,
    },
    {
      title: '禁用时间',
      dataIndex: 'disabled_time',
      width: 160,
      resizable: true,
    },
    {
      title: '是否为持续\n禁用设备',
      dataIndex: 'is_continuous_disabled',
      width: 120,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '备注',
      dataIndex: 'remark',
      width: 120,
      resizable: true,
    },
  ];
}

export const macInfoSchemas: FormSchema[] = [
  {
    label: 'mac地址',
    field: 'mac_addr',
    component: 'Input',
    colProps: {
      span: 24,
    },
    componentProps: {
      readonly: true,
    },
    helpMessage: '当前主机mac地址',
  },
  {
    label: '原区域',
    field: 'originalArea',
    component: 'Input',
    colProps: { span: 24 },
    componentProps: {
      readonly: true,
      bordered: false,
    },
  },
  {
    label: '运营商',
    field: 'originalIsp',
    component: 'Input',
    colProps: { span: 24 },
    componentProps: {
      readonly: true,
      bordered: false,
    },
  },
  {
    label: '替换信息',
    field: 'ReplaceInfo',
    component: 'Divider', // 添加分隔线提升可视化
    colProps: { span: 24 },
  },
  {
    label: '目标host',
    field: 'hostname',
    component: 'Input',
    colProps: {
      span: 24,
    },
    componentProps: {
      allowClear: true,
      placeholder: '请输入目标主机名(hostname)',
    },
    required: true,
  },
  {
    field: 'modify_ip',
    label: '是否修改IP',
    component: 'RadioGroup',
    colProps: { span: 12 },
    componentProps: {
      options: [
        { label: '是', value: 1 },
        { label: '否', value: 0 },
      ],
    },
    required: true,
  },
];

// 新增：MAC替换历史表格列配置
export const macReplaceHistoryColumns: BasicColumn[] = [
  {
    title: 'MAC地址',
    dataIndex: 'mac_addr',
    width: 160,
    resizable: true,
    fixed: 'left',
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 120,
    resizable: true,
  },
  {
    title: '节点/机房',
    dataIndex: 'node',
    width: 120,
    resizable: true,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
  },
  {
    title: 'IP地址',
    dataIndex: 'ip',
    width: 150,
    resizable: true,
  },
  {
    title: '区域',
    dataIndex: 'area',
    width: 80,
    resizable: true,
  },
  {
    title: '省份',
    dataIndex: 'province',
    width: 80,
    resizable: true,
  },
  {
    title: '城市',
    dataIndex: 'city',
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
    title: '调度控制',
    dataIndex: 'is_province_scheduling',
    width: 100,
    resizable: true,
  },
  {
    title: '是否异网',
    dataIndex: 'is_cover_diff_isp',
    width: 80,
    resizable: true,
  },
  // {
  //   title: '异网运营商',
  //   dataIndex: 'cover_diff_isp',
  //   width: 80,
  //   resizable: true,
  // },
  // {
  //   title: '调度区域',
  //   dataIndex: 'schedule_area',
  //   width: 80,
  //   resizable: true,
  // },
  // {
  //   title: '调度省份',
  //   dataIndex: 'schedule_province',
  //   width: 80,
  //   resizable: true,
  // },
  // {
  //   title: '调度城市',
  //   dataIndex: 'schedule_city',
  //   width: 80,
  //   resizable: true,
  // },
  {
    title: '替换时间',
    dataIndex: 'created_at',
    width: 160,
    resizable: true,
  },
  {
    title: '操作人',
    dataIndex: 'operator',
    width: 100,
    resizable: true,
  },
];

//MAC禁用查询
export function macDisableSearchFormSchema(bizType: string): FormSchema[] {
  return [
    {
      label: '厂商',
      field: 'provider',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.PRODIVER),
        allowClear: true,
        defaultValue: 'mf',
      },
      colProps: { span: 4 },
      ifShow: () => bizType === 'specialLine',
    },
    {
      label: 'MAC地址',
      field: 'mac',
      component: 'Input',
      colProps: { span: 6 },
      componentProps: {
        placeholder: '请输入MAC地址',
      },
    },
    {
      label: '禁用时间范围',
      field: 'disable_time_range',
      component: 'RangePicker',
      colProps: { span: 12 },
      componentProps: {
        format: 'YYYY-MM-DD HH:mm:ss',
        showTime: true,
        defaultValue: [dayjs().subtract(1, 'hour'), dayjs()],
      },
    },
  ];
}

// MAC禁用流水表列配置
export const macDisableHistoryColumns: BasicColumn[] = [
  {
    title: 'MAC地址',
    dataIndex: 'mac',
    width: 180,
    resizable: true,
    fixed: 'left',
  },
  {
    title: '开始时间',
    dataIndex: 'disable_time',
    width: 220,
    resizable: true,
    // 格式化时间 +8 小时
    // format: (text: string) => {
    //   return dayjs(text, 'ddd, DD MMM YYYY HH:mm:ss [GMT]')
    //     .add(8, 'hour')
    //     .format('YYYY-MM-DD HH:mm:ss');
    // },
  },
  {
    title: '原因',
    dataIndex: 'status',
  },
];

// MAC历史记录
export const MacHistoryColumns: BasicColumn[] = [
  {
    title: 'MAC地址',
    dataIndex: 'mac_addr',
    width: 200,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 200,
  },
  {
    title: '业务名',
    dataIndex: 'business',
    width: 200,
  },
  {
    title: '设备ID',
    dataIndex: 'frank_id',
  },
  {
    title: '使用时间',
    dataIndex: 'created_at',
    width: 200,
  },
];

export function MacHistorySearchFormSchema(): FormSchema[] {
  return [
    {
      label: 'MAC地址',
      field: 'mac_addr',
      component: 'Input',
      colProps: { span: 8 },
      componentProps: {
        placeholder: '请输入MAC地址',
      },
    },
    {
      label: '主机名',
      field: 'hostname',
      component: 'Input',
      colProps: { span: 8 },
      componentProps: {
        placeholder: '请输入主机名地址',
      },
    },
  ];
}

// MAC替换记录搜索
export function MacReplaceSearchFormSchema(): FormSchema[] {
  return [
    {
      label: 'MAC地址',
      field: 'mac_addr',
      component: 'Input',
      colProps: { span: 8 },
    },
    {
      label: '主机名',
      field: 'hostname',
      component: 'Input',
      colProps: { span: 8 },
    },
  ];
}
