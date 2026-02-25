import { BasicColumn, FormSchema } from '@/components/Table';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { ZEnum } from '@/enums/dictTypeCode';
import { ref } from 'vue';
import { GetZapAreaList } from '@/api/business/zap';
import { formatToDateTime } from '@/utils/dateUtil';

export const submitTypeMap = getDictDataMapFromDict(ZEnum.Z_SUBMIT_TYPE);
export const submitTypeList = getSelectOptionsFromDict(ZEnum.Z_SUBMIT_TYPE);
export const ispMap = getDictDataMapFromDict(ZEnum.Z_ISP_NAME);
export const deliveryStatusMap = getDictDataMapFromDict(ZEnum.Z_DELIVERY_STATUS);
export const deviceLoginTypeMap = getDictDataMapFromDict(ZEnum.Z_DEVICE_LOGIN_TYPE);
export const deliveryResourceTypeMap = getDictDataMapFromDict(ZEnum.Z_DELIVERY_RESOURCE_TYPE);
export const deliveryResourceTypeList = getSelectOptionsFromDict(ZEnum.Z_DELIVERY_RESOURCE_TYPE);
export const deviceLoginList = getSelectOptionsFromDict(ZEnum.Z_DEVICE_LOGIN_TYPE);
export const deviceSystemList = getSelectOptionsFromDict(ZEnum.Z_MACHINE_TYPE);
export const accountTypeMap = getDictDataMapFromDict(ZEnum.Z_ACCOUNT_TYPE);

type Option = { label: string; value: string; children?: any[] };
const provincesOptions = ref<Option[]>([]);

// 获取省份城市数据
export const loadProvincesData = async () => {
  const res = await GetZapAreaList();
  provincesOptions.value = res || [];
};

// 在模块加载时自动获取数据（可选）
loadProvincesData();

export const searchFormSchema: FormSchema[] = [
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
    field: 'ip',
    label: 'ip地址',
    component: 'Input',
    colProps: { span: 5 },
  },
  {
    field: 'order_status',
    label: '交付状态',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(ZEnum.Z_DELIVERY_STATUS),
      allowClear: true,
    },
    colProps: { span: 5 },
  },
  {
    field: 'business',
    label: '业务',
    component: 'Select',
    colProps: { span: 5 },
  },
  {
    field: 'resource_isp',
    label: '运行商',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(ZEnum.Z_ISP_NAME),
      allowClear: true,
    },
    colProps: { span: 4 },
  },
  {
    field: 'province_code',
    label: '省份',
    component: 'Select',
    colProps: { span: 5 },
    componentProps: ({ formModel, formActionType }) => {
      return {
        options: provincesOptions.value,
        onChange: (v: any) => {
          formModel.city_code = undefined; //  reset city value
          const cityOptions =
            provincesOptions.value.find((item) => item.value === v)?.children || [];
          const { updateSchema } = formActionType;
          updateSchema({
            field: 'city_code',
            componentProps: {
              options: cityOptions,
            },
          });
        },
      };
    },
  },
  {
    field: 'city_code',
    label: '城市',
    component: 'Select',
    colProps: { span: 5 },
  },
];

export const columns: BasicColumn[] = [
  {
    title: '节点ID',
    dataIndex: 'task_id',
    width: 120,
    resizable: true,
    helpMessage: '录入业务平台之后返回的ID',
    customRender: ({ record }) => {
      return record.task_id != 0 ? record.task_id : '';
    },
  },
  {
    title: '账号(月|日95)',
    dataIndex: 'account_type',
    width: 110,
    resizable: true,
    customRender: ({ record }) => {
      if (record.account_type == 'mful') {
        return '月95';
      }
      if (record.account_type == 'mfulone') {
        return '日95';
      }
    },
  },
  {
    title: '节点label',
    dataIndex: 'label',
    width: 180,
    resizable: true,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
  },
  {
    title: '业务',
    dataIndex: 'business',
    width: 180,
    resizable: true,
  },
  {
    title: '资源类型',
    dataIndex: 'resource_type',
    width: 80,
    resizable: true,
    customRender: ({ record }) => {
      return deliveryResourceTypeList.find((item) => item.value === record.resource_type)?.label;
    },
  },
  {
    title: '机器类型',
    dataIndex: 'machine_type',
    width: 100,
    resizable: true,
  },
  {
    title: '提交类型',
    dataIndex: 'submit_type',
    width: 100,
    resizable: true,
  },
  {
    title: '运行商',
    dataIndex: 'resource_isp',
    width: 80,
    resizable: true,
    customRender: ({ record }) => {
      return ispMap[record.resource_isp].dictLabel;
    },
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
    title: '内网ip',
    dataIndex: 'inner_ip',
    width: 150,
    resizable: true,
  },
  {
    title: '外网ip',
    dataIndex: 'outer_ips',
    width: 150,
    resizable: true,
  },
  {
    title: '是否只覆盖本省',
    dataIndex: 'is_only_cover',
    width: 120,
    resizable: true,
    customRender: ({ record }) => {
      const text = record.is_only_cover === true ? '是' : '否';
      return text;
    },
  },
  {
    title: '是否内网NAT1资源',
    dataIndex: 'is_intranet_resource',
    width: 140,
    resizable: true,
    customRender: ({ record }) => {
      const text = record.is_intranet_resource === true ? '是' : '否';
      return text;
    },
  },
  {
    title: '是否仅晚高峰',
    dataIndex: 'only_evening_peak',
    width: 140,
    resizable: true,
    customRender: ({ record }) => {
      const text = record.only_evening_peak === true ? '是' : '否';
      return text;
    },
  },
  {
    title: '总带宽(Mbps)',
    dataIndex: 'total_bandwidth',
    width: 120,
    resizable: true,
  },
  {
    title: '单线路带宽(Mbps)',
    dataIndex: 'single_line_bandwidth',
    width: 130,
    resizable: true,
  },
  {
    title: '线路总数',
    dataIndex: 'bwcount',
    width: 120,
    resizable: true,
  },
  {
    title: '登录方式',
    dataIndex: 'ip_domain_option',
    width: 120,
    resizable: true,
    customRender: ({ record }) => {
      const text = record.ip_domain_option === 1 ? '域名' : 'IP';
      return text;
    },
  },
  {
    title: '登录域名',
    dataIndex: 'domain',
    width: 120,
    resizable: true,
  },
  {
    title: '登录账号',
    dataIndex: 'domain_account_name',
    width: 120,
    resizable: true,
  },
  {
    title: '登录密码',
    dataIndex: 'domain_password',
    width: 120,
    resizable: true,
  },
  {
    title: '登录端口',
    dataIndex: 'domain_login_port',
    width: 120,
    resizable: true,
  },
  {
    title: '交付状态',
    dataIndex: 'order_status',
    width: 120,
    resizable: true,
    fixed: 'right',
  },
  {
    title: '供应商名称',
    dataIndex: 'supplier_name',
    width: 200,
    resizable: true,
  },
  {
    title: '结果',
    dataIndex: 'result',
    width: 350,
    resizable: true,
  },
  {
    title: '操作账户',
    dataIndex: 'operation_account',
    width: 100,
    resizable: true,
  },
  {
    title: '备注',
    dataIndex: 'remark',
    width: 150,
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
    dataIndex: 'update_at',
    width: 150,
    resizable: true,
    customRender: ({ record }) => {
      return formatToDateTime(record.update_at);
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
    label: 'frankId',
    field: 'frank_id',
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
    field: 'hostname',
    component: 'Input',
    colProps: {
      span: 22,
    },
    required: true,
    componentProps: {
      disabled: true,
      readonly: true,
    },
  },
  {
    label: '资源类型',
    field: 'resource_type',
    component: 'Select',
    colProps: {
      span: 22,
    },
    required: true,
    componentProps: {
      options: getSelectOptionsFromDict(ZEnum.Z_DELIVERY_RESOURCE_TYPE),
    },
  },
  {
    label: '机器类型',
    field: 'machine_type',
    component: 'Select',
    colProps: {
      span: 22,
    },
    required: true,
    componentProps: {
      options: deviceSystemList,
    },
  },
  {
    label: '提交类型',
    field: 'submit_type',
    component: 'Select',
    colProps: {
      span: 22,
    },
    required: true,
    componentProps: {
      options: getSelectOptionsFromDict(ZEnum.Z_SUBMIT_TYPE),
    },
  },
  {
    label: '资源归属运行商',
    field: 'resource_isp',
    component: 'Select',
    colProps: {
      span: 22,
    },
    required: true,
    componentProps: {
      options: getSelectOptionsFromDict(ZEnum.Z_ISP_NAME),
    },
  },
  {
    label: '省份',
    field: 'province_code',
    component: 'Select',
    colProps: {
      span: 6,
    },
    required: true,
    componentProps: ({ formModel, formActionType }) => {
      return {
        options: provincesOptions.value,
        onChange: (v: any) => {
          formModel.city_code = undefined; //  reset city value
          const cityOptions =
            provincesOptions.value.find((item) => item.value === v)?.children || [];
          const { updateSchema } = formActionType;
          updateSchema({
            field: 'city_code',
            componentProps: {
              options: cityOptions,
            },
          });
        },
      };
    },
  },
  {
    label: '城市',
    field: 'city_code',
    component: 'Select',
    colProps: {
      span: 8,
    },
    required: true,
    componentProps: {
      options: [],
    },
  },
  {
    label: '内网IP',
    field: 'inner_ip',
    component: 'Input',
    colProps: {
      span: 22,
    },
    required: true,
  },
  {
    label: '外网IP',
    field: 'outer_ips',
    component: 'Input',
    colProps: {
      span: 22,
    },
    required: true,
  },
  {
    label: '是否只覆盖本省',
    field: 'is_only_cover',
    component: 'RadioGroup',
    colProps: {
      span: 8,
    },
    required: true,
    componentProps: {
      options: [
        {
          label: '是',
          value: true,
        },
        {
          label: '否',
          value: false,
        },
      ],
    },
  },
  {
    label: '内网NAT1资源',
    field: 'is_intranet_resource',
    component: 'RadioGroup',
    colProps: {
      span: 8,
    },
    required: true,
    componentProps: {
      options: [
        {
          label: '是',
          value: true,
        },
        {
          label: '否',
          value: false,
        },
      ],
    },
  },
  {
    label: '节点总带宽',
    field: 'total_bandwidth',
    component: 'Input',
    colProps: {
      span: 22,
    },
    required: true,
  },
  {
    label: '单线路带宽',
    field: 'single_line_bandwidth',
    component: 'Input',
    colProps: {
      span: 22,
    },
    required: true,
  },
  {
    label: '登录方式',
    field: 'ip_domain_option',
    component: 'Select',
    colProps: {
      span: 22,
    },
    required: true,
    componentProps: {
      options: deviceLoginList,
    },
  },
  {
    label: '登录域名',
    field: 'domain',
    component: 'Input',
    colProps: {
      span: 22,
    },
    required: true,
  },
  {
    label: '登录端口',
    field: 'domain_login_port',
    component: 'Input',
    colProps: {
      span: 22,
    },
    required: true,
  },
  {
    label: '登录密码',
    field: 'domain_password',
    component: 'Input',
    colProps: {
      span: 22,
    },
  },
  {
    label: '交付账号',
    field: 'account_type',
    component: 'Select',
    colProps: {
      span: 22,
    },
    componentProps: {
      options: getSelectOptionsFromDict(ZEnum.Z_ACCOUNT_TYPE),
    },
  },
  {
    label: '交付状态',
    field: 'order_status',
    component: 'Select',
    colProps: {
      span: 22,
    },
    componentProps: {
      options: getSelectOptionsFromDict(ZEnum.Z_DELIVERY_STATUS),
    },
  },
];
