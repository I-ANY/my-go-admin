import { BasicColumn, FormSchema } from '@/components/Table';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { ZEnum } from '@/enums/dictTypeCode';
import { h } from 'vue';
import { Tag } from 'ant-design-vue';

export const deviceStatusMap = getDictDataMapFromDict(ZEnum.Z_DEVICE_STATUS);

interface common {
  label: string;
  value: any;
}

export const ispList: common[] = [
  { label: '电信', value: '电信' },
  { label: '联通', value: '联通' },
  { label: '移动', value: '移动' },
];
export const statusList: common[] = [
  { label: 'Active', value: 1 },
  { label: 'Offline', value: 2 },
  { label: 'Deactive', value: 4 },
];

export const periodList: common[] = [
  { label: '建设中', value: 0 },
  { label: '建设完成', value: 1 },
  { label: '节点下架', value: 2 },
];

export const statesList: common[] = [
  { label: '正常', value: 0 },
  { label: '降线', value: 1 },
  { label: '挂起', value: 2 },
];

export const searchFormSchema: FormSchema[] = [
  {
    field: 'labels',
    label: ' ',
    component: 'InputTextArea',
    colProps: { span: 5 },
    componentProps: {
      allowClear: true,
      placeholder: '输入节点label进行搜索，一行一个',
      rows: 2,
    },
  },
  {
    field: 'hostnames',
    label: ' ',
    component: 'InputTextArea',
    colProps: { span: 5 },
    componentProps: {
      allowClear: true,
      placeholder: '输入主机名进行搜索，一行一个',
      rows: 2,
    },
  },
  {
    field: 'owner',
    label: ' ',
    component: 'Input',
    colProps: { span: 4 },
    componentProps: {
      allowClear: true,
      placeholder: '输入用户进行搜索',
    },
  },
  {
    field: 'name',
    label: ' ',
    component: 'Input',
    colProps: { span: 4 },
    componentProps: {
      allowClear: true,
      placeholder: '输入节点名称进行搜索',
    },
  },
  {
    field: 'ip',
    label: ' ',
    component: 'Input',
    colProps: { span: 3 },
    componentProps: {
      allowClear: true,
      placeholder: '输入IP进行搜索',
    },
  },
  {
    field: 'province',
    label: ' ',
    component: 'Input',
    colProps: { span: 3 },
    componentProps: {
      allowClear: true,
      placeholder: '输入省份进行搜索',
    },
  },
  {
    field: 'city',
    label: ' ',
    component: 'Input',
    colProps: { span: 3 },
    componentProps: {
      allowClear: true,
      placeholder: '输入城市进行搜索',
    },
  },
  {
    field: 'account_type',
    label: ' ',
    component: 'Select',
    colProps: { span: 3 },
    componentProps: {
      options: getSelectOptionsFromDict(ZEnum.Z_ACCOUNT_TYPE),
      allowClear: true,
      placeholder: '选择账号类型',
    },
  },
  {
    field: 'resource_properties',
    label: ' ',
    component: 'Select',
    colProps: { span: 3 },
    componentProps: {
      options: [
        {
          label: '汇聚',
          value: '汇聚',
        },
        {
          label: '专线',
          value: '专线',
        },
      ],
      allowClear: true,
      placeholder: '选择资源属性',
    },
  },
  {
    field: 'status',
    label: ' ',
    component: 'Select',
    componentProps: {
      options: statusList,
      allowClear: true,
      mode: 'multiple', // 启用多选模式
      placeholder: '选择节点状态',
      maxTagCount: 3, // 最多显示3个标签，超出显示+X
    },
    colProps: { span: 3 },
  },
  {
    field: 'isp',
    label: ' ',
    component: 'Select',
    componentProps: {
      options: ispList,
      allowClear: true,
      placeholder: '选择运营商',
    },
    colProps: { span: 3 },
  },
  {
    field: 'isp_status',
    label: ' ',
    component: 'Select',
    componentProps: {
      options: [
        {
          label: '是',
          value: 1,
        },
        {
          label: '否',
          value: 0,
        },
      ],
      allowClear: true,
      placeholder: '选择运营商是否一致',
    },
    colProps: { span: 3 },
  },
];
export const columns: BasicColumn[] = [
  {
    title: '用户',
    dataIndex: 'owner',
    width: 150,
    resizable: true,
  },
  {
    title: '节点名称',
    dataIndex: 'name',
    width: 200,
    resizable: true,
  },
  {
    title: '账号(日|月95)',
    dataIndex: 'account_type',
    width: 120,
    resizable: true,
    customRender: ({ record }) => {
      if (record.account_type == 'mful') {
        return '月95';
      } else {
        return '日95';
      }
    },
  },
  {
    title: '资源属性',
    dataIndex: 'resource_properties',
    width: 150,
    resizable: true,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
  },
  {
    title: '节点label',
    dataIndex: 'label',
    width: 180,
    resizable: true,
  },
  {
    title: '节点状态',
    dataIndex: 'status',
    width: 100,
    resizable: true,
  },
  {
    title: '生命周期',
    dataIndex: 'period',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      let statusText = '';
      let color = '';
      if (record.period === 0) {
        statusText = '建设中';
        color = '#f3872f';
      } else if (record.period === 1) {
        statusText = '建设完成';
        color = '#298ded';
      } else if (record.period === 2) {
        statusText = '节点下架';
        color = '#f30000';
      } else if (record.period === 3) {
        statusText = '取消建设';
        color = '#FF1100';
      } else {
        statusText = '';
        color = '#cccccc'; // 默认颜色
      }
      return h(Tag, { color: color }, () => statusText);
    },
  },
  {
    title: '字节运营商',
    dataIndex: 'isp',
    width: 100,
    resizable: true,
  },
  {
    title: 'ECDN运营商',
    dataIndex: 'isp_ecdn',
    width: 100,
    resizable: true,
  },
  {
    title: '运营商是否一致',
    dataIndex: 'isp_status',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      let statusText = '';
      let color = '';
      if (record.isp_status === 1) {
        statusText = '是';
        color = '#2eb438';
      } else if (record.isp_status === 0) {
        statusText = '否';
        color = '#FF1100';
      } else {
        statusText = '';
        color = '#cccccc'; // 默认颜色
      }
      return h(Tag, { color: color }, () => statusText);
    },
  },
  {
    title: '容量',
    dataIndex: 'capacity',
    width: 80,
    resizable: true,
  },
  {
    title: '带宽(G)',
    dataIndex: 'isp_bandwidth_max',
    width: 80,
    resizable: true,
  },
  {
    title: '保底(G)',
    dataIndex: 'guarantee',
    width: 80,
    resizable: true,
  },
  {
    title: '是否覆盖本省',
    dataIndex: 'is_only_cover',
    width: 120,
    resizable: true,
  },
  {
    title: '是否内网资源',
    dataIndex: 'is_intranet_resource',
    width: 120,
    resizable: true,
  },
  {
    title: '公网ipv4',
    dataIndex: 'ipv4',
    width: 120,
    resizable: true,
  },
  {
    title: '公网ipv6',
    dataIndex: 'ipv6',
    width: 120,
    resizable: true,
  },
  {
    title: '省份',
    dataIndex: 'Province',
    width: 120,
    resizable: true,
  },
  {
    title: '城市',
    dataIndex: 'city',
    width: 120,
    resizable: true,
  },
  {
    title: '是否是直接签约',
    dataIndex: 'is_direct_sign',
    width: 120,
    resizable: true,
  },
  {
    title: '是否是规避节点',
    dataIndex: 'is_avoidance_node',
    width: 120,
    resizable: true,
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    width: 150,
    resizable: true,
  },
  {
    title: '上线时间',
    dataIndex: 'online_time',
    width: 150,
    resizable: true,
  },
  {
    title: '更新时间',
    dataIndex: 'updated_at',
    width: 150,
    resizable: true,
  },
  {
    title: '临时下线时间',
    dataIndex: 'temp_offline_time',
    width: 150,
    resizable: true,
  },
  {
    title: '永久下线时间',
    dataIndex: 'offline_time',
    width: 150,
    resizable: true,
  },
  {
    title: '执行结果',
    dataIndex: 'exec_result',
    width: 250,
    resizable: true,
  },
  {
    title: '操作账户',
    dataIndex: 'operation_account',
    width: 120,
    resizable: true,
  },
];

export const realTimeNodeSchema: FormSchema[] = [
  {
    field: 'hostname',
    label: '主机名',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'label',
    label: '节点label',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    label: '类型',
    field: 'account_type',
    component: 'Select',
    defaultValue: 'day95',
    required: true,
    colProps: {
      span: 4,
    },
    componentProps: {
      options: [
        {
          label: '日95',
          value: 'day95',
        },
        {
          label: '月95',
          value: 'month95',
        },
      ],
    },
  },
];

export const realTimeStatusColumns: BasicColumn[] = [
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
  },
  {
    title: '节点label',
    dataIndex: 'label',
    width: 180,
    resizable: true,
  },
  {
    title: '节点名称',
    dataIndex: 'node_cn_name',
    width: 200,
    resizable: true,
  },
  {
    title: '节点状态',
    dataIndex: 'status',
    width: 100,
    resizable: true,
  },
  {
    title: '跑量低原因',
    dataIndex: 'conclusion',
    width: 200,
    resizable: true,
  },
  {
    title: '质量调度状态',
    dataIndex: 'state',
    width: 120,
    resizable: true,
  },
  {
    title: '当前带宽',
    dataIndex: 'nginx_bw',
    width: 150,
    resizable: true,
  },
  {
    title: '可用空间(G)',
    dataIndex: 'available_space',
    width: 120,
    resizable: true,
  },
  {
    title: '质量指标',
    dataIndex: 'metric',
    width: 500,
    resizable: true,
  },
  {
    title: '压量系数',
    dataIndex: 'coef',
    width: 100,
    resizable: true,
  },
  {
    title: '质量反复异常次数',
    dataIndex: 'discrete_reduce_count',
    width: 180,
    resizable: true,
  },
  {
    title: '异常质量指标',
    dataIndex: 'is_metric_abnormal',
    width: 150,
    resizable: true,
  },
  {
    title: '上次异常压量系数',
    dataIndex: 'last_abnormal_rate',
    width: 150,
    resizable: true,
  },
  {
    title: '探测异常点数(1/min)',
    dataIndex: 'abnormal_count',
    width: 200,
    resizable: true,
  },
  {
    title: '质量降线buffer(单位:M)',
    dataIndex: 'no80_node_buffer',
    width: 200,
    resizable: true,
  },

  {
    title: '节点规划上限(单位:G)',
    dataIndex: 'node_limit_bandwidth',
    width: 200,
    resizable: true,
  },
  {
    title: '节点安全水位(单位:G)',
    dataIndex: 'node_safe_bandwidth',
    width: 200,
    resizable: true,
  },
  {
    title: '汇聚资源可用线路',
    dataIndex: 'nat_eth_info',
    width: 150,
    resizable: true,
  },
  {
    title: '可用ip',
    dataIndex: 'enable_ips_info',
    width: 180,
    resizable: true,
  },
  {
    title: '不可用ip',
    dataIndex: 'disable_ips_info',
    width: 180,
    resizable: true,
  },
  {
    title: '是否为NAT资源',
    dataIndex: 'is_intranet_resource',
    width: 150,
    resizable: true,
  },
];
