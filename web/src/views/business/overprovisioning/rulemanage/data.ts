import { BasicColumn, FormSchema } from '@/components/Table';
import {
  GetBusinessGroupOptions,
  GetCategoryList,
  GetSubcategoryListAll,
} from '@/api/business/biz';

import { getNoRuleBusiness } from '@/api/business/overprovisioning';
import { h } from 'vue';

// 规则名称 ruleName  内存总大小 mem  ssd总大小 ssdSize
// hdd总大小 hddSize 存储带宽比 storageBwRatio 业务归属 businessBelong
// 规则状态 ruleStatus 添加人 addWho 最后修改人 updateWho
//将这些字段 都放到搜索表单中
export const searchFormSchema = (): FormSchema[] => [
  {
    field: 'ruleName',
    label: '规则名称',
    component: 'Input',
    componentProps: {
      placeholder: '请输入规则名称',
    },
    colProps: { span: 6 },
  },
  {
    field: 'businessIds',
    label: '业务名称',
    component: 'ApiSelect',
    componentProps: {
      placeholder: '请输入业务名称',
      api: GetSubcategoryListAll,
      params: {
        pageSize: 10000,
        pageIndex: 1,
        status: 1,
      },
      labelField: 'name',
      valueField: 'id',
      mode: 'multiple',
      resultField: 'items',
      showSearch: true,
      optionFilterProp: 'label',
      filterOption: (input: string, option: any) => {
        return option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'ruleStatus',
    label: '规则状态',
    component: 'Select',
    componentProps: {
      placeholder: '请选择规则状态',
      options: [
        { label: '启用', value: 0 },
        { label: '禁用', value: 1 },
      ],
    },
    colProps: { span: 6 },
  },
  {
    field: 'ruleDescription',
    label: '规则描述',
    component: 'Input',
    componentProps: {
      placeholder: '请输入规则描述',
    },
    colProps: { span: 6 },
  },
];

export const columns: BasicColumn[] = [
  {
    title: '规则名称',
    dataIndex: 'ruleName',
    width: 180,
    resizable: true,
  },
  {
    title: '业务归属',
    dataIndex: 'businesses',
    width: 250,
    resizable: true,
    ellipsis: true,
    customRender: ({ record }) => {
      const text =
        record.businesses && Array.isArray(record.businesses)
          ? record.businesses.map((b) => b.businessName).join('; ')
          : record.businessBelong || '';

      if (!text) return '';

      return h(
        'span',
        {
          style: {
            display: 'block',
            maxWidth: '100%',
            whiteSpace: 'pre-line',
          },
        },
        text,
      );
    },
  },
  {
    title: '超配规则',
    dataIndex: 'planConf',
    width: 120,
    resizable: true,
    customRender: ({ record }) => {
      const planconfData = record.planConf || record.planconf || record.PlanConf;
      if (!planconfData) return '';
      try {
        if (Array.isArray(planconfData)) {
          return `共${planconfData.length}条`;
        }
        return typeof planconfData === 'object' ? JSON.stringify(planconfData) : planconfData;
      } catch (e) {
        return planconfData;
      }
    },
  },
  {
    title: '系统盘大小(GB)',
    dataIndex: 'sysDiskSize',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      return record.sysDiskSize || record.sysdisksize || record.SysDiskSize || 0;
    },
  },
  {
    title: '总内存大小(GB)',
    dataIndex: 'mem',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      return record.mem || 0;
    },
  },
  {
    title: 'SSD总大小(TB)',
    dataIndex: 'ssdSize',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      return record.ssdSize || 0;
    },
  },
  {
    title: 'HDD总大小(TB)',
    dataIndex: 'hddSize',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      return record.hddSize || 0;
    },
  },
  {
    title: '存储带宽比',
    dataIndex: 'storageBwRatio',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      return record.storageBwRatio || 0;
    },
  },
  {
    title: '规则状态',
    dataIndex: 'ruleStatus',
    width: 100,
    resizable: true,
  },
  {
    title: '规则描述',
    dataIndex: 'ruleDescription',
    width: 200,
    resizable: true,
  },
  {
    title: '添加人',
    dataIndex: 'createUser',
    width: 100,
    resizable: true,
  },
  {
    title: '更新人',
    dataIndex: 'updateUser',
    width: 100,
    resizable: true,
  },
  {
    title: '操作',
    key: 'action',
    width: 150,
    fixed: 'right',
    // slots: { customRender: 'action' },
  },
];

export const ruleFormSchema: FormSchema[] = [
  {
    field: 'ruleName',
    label: '规则名称',
    component: 'Input',
    required: true,
    componentProps: {
      placeholder: '请输入规则名称',
    },
    colProps: { span: 24 },
  },
  {
    field: 'sysDiskSize',
    label: '系统盘大小',
    component: 'InputNumber',
    required: true,
    componentProps: {
      placeholder: '请输入系统盘大小(GB)',
      min: 0,
      style: { width: '100%' },
    },
    colProps: { span: 8 },
  },
  {
    field: 'enableSysDiskCheck',
    // label: '系统盘检测',
    component: 'Switch',
    defaultValue: false,
    componentProps: {
      checkedValue: true,
      unCheckedValue: false,
      checkedChildren: '开',
      unCheckedChildren: '关',
      style: { marginLeft: '10px' },
    },
    colProps: { span: 4 },
  },
  {
    field: 'mem',
    label: '内存总大小',
    component: 'InputNumber',
    required: true,
    componentProps: {
      placeholder: '请输入内存总大小(GB)',
      min: 0,
      style: { width: '100%' },
    },
    colProps: { span: 8 },
  },
  {
    field: 'enableMemoryCheck',
    component: 'Switch',
    defaultValue: true,
    componentProps: {
      checkedValue: true,
      unCheckedValue: false,
      checkedChildren: '开',
      unCheckedChildren: '关',
      style: { marginLeft: '10px' },
    },
    colProps: { span: 4 },
  },
  {
    field: 'ssdSize',
    label: 'SSD总大小',
    component: 'InputNumber',
    required: true,
    componentProps: {
      placeholder: '请输入SSD总大小(TB)',
      min: 0,
      style: { width: '100%' },
    },
    colProps: { span: 8 },
  },
  {
    field: 'enableSsdCheck',
    component: 'Switch',
    defaultValue: true,
    componentProps: {
      checkedValue: true,
      unCheckedValue: false,
      checkedChildren: '开',
      unCheckedChildren: '关',
      style: { marginLeft: '10px' },
    },
    colProps: { span: 4 },
  },
  {
    field: 'hddSize',
    label: 'HDD总大小',
    component: 'InputNumber',
    required: true,
    componentProps: {
      placeholder: '请输入HDD总大小(TB)',
      min: 0,
      style: { width: '100%' },
    },
    colProps: { span: 8 },
  },
  {
    field: 'enableHddCheck',
    component: 'Switch',
    defaultValue: true,
    componentProps: {
      checkedValue: true,
      unCheckedValue: false,
      checkedChildren: '开',
      unCheckedChildren: '关',
      style: { marginLeft: '10px' },
    },
    colProps: { span: 4 },
  },
  {
    field: 'storageBwRatioGroup',
    label: '存储带宽比',
    required: false,
    // rules: [{ required: true, validator: () => Promise.resolve() }],
    slot: 'storageBwRatioGroup',
    colProps: { span: 8 },
  },
  {
    field: 'enableStorageBwRatioCheck',
    component: 'Switch',
    defaultValue: false,
    componentProps: {
      checkedValue: true,
      unCheckedValue: false,
      checkedChildren: '开',
      unCheckedChildren: '关',
      style: { marginLeft: '10px' },
    },
    colProps: { span: 4 },
  },
  {
    field: 'storageBwRatioDiskType',
    label: '存储带宽比类型',
    component: 'Select',
    defaultValue: 0,
    componentProps: {
      placeholder: '请选择存储类型',
      options: [
        { label: 'SSD', value: 0 },
        { label: 'HDD', value: 1 },
        { label: '业务盘(SSD+HDD)', value: 2 },
      ],
      // style: { width: '100%' },
    },
    colProps: { span: 8 },
  },

  {
    field: 'enableDataDiskCheck',
    label: '业务盘检测',
    component: 'Switch',
    defaultValue: false,
    componentProps: {
      checkedValue: true,
      unCheckedValue: false,
      checkedChildren: '开',
      unCheckedChildren: '关',
      // style: { marginLeft: '10px' },
    },
    colProps: { span: 12 },
  },
  {
    field: 'ruleStatus',
    label: '规则状态',
    component: 'Switch',
    required: true,
    defaultValue: 0,
    componentProps: {
      checkedValue: 0,
      unCheckedValue: 1,
      checkedChildren: '启用',
      unCheckedChildren: '禁用',
    },
    colProps: { span: 12 },
  },
  {
    field: 'isProvinceScheduling',
    label: '是否跨省',
    component: 'Switch',
    defaultValue: 0,
    componentProps: {
      checkedValue: 1,
      unCheckedValue: 0,
      checkedChildren: '是',
      unCheckedChildren: '否',
    },
    colProps: { span: 12 },
  },
  // {
  //   field: 'enableDataDiskCheck',
  //   label: '业务盘检测',
  //   component: 'Switch',
  //   defaultValue: true,
  //   componentProps: {
  //     checkedValue: true,
  //     unCheckedValue: false,
  //     checkedChildren: '开',
  //     unCheckedChildren: '关',
  //   },
  //   colProps: { span: 12 },
  // },
  {
    field: 'businessIds',
    label: '业务归属',
    component: 'ApiSelect',
    required: true,
    componentProps: {
      api: async (params) => {
        const res = await GetSubcategoryListAll(params);
        const items = Array.isArray(res?.items) ? res.items : [];
        const hasAll = items.some((item) => item.id === 0);
        return {
          ...(res || {}),
          items: hasAll ? items : [{ id: 0, name: '全部业务' }, ...items],
        };
      },
      params: {
        pageSize: 10000,
        pageIndex: 1,
        status: 1,
      },
      labelField: 'name',
      valueField: 'id',
      mode: 'multiple',
      resultField: 'items',
      showSearch: true,
      optionFilterProp: 'label',
      filterOption: (input: string, option: any) => {
        return option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 24 },
  },
  {
    field: 'planConfGroup',
    label: '超配规则配置',
    slot: 'planconfGroup',
    colProps: { span: 24 },
  },
  {
    field: 'ruleDescription',
    label: '规则描述',
    component: 'InputTextArea', // 多行文本框
    componentProps: {
      placeholder: '请输入规则描述',
      autoSize: {
        minRows: 3,
        maxRows: 6,
      },
    },
    colProps: { span: 24 },
  },
];

// 无规则业务列表列配置
export const noRuleBusinessColumns: BasicColumn[] = [
  {
    title: '业务组',
    dataIndex: 'businessGroup',
    width: 150,
    resizable: true,
    ellipsis: true,
  },
  {
    title: '业务大类',
    dataIndex: 'categoryName',
    width: 150,
    resizable: true,
    ellipsis: true,
  },
  {
    title: '业务名称',
    dataIndex: 'businessName',
    width: 200,
    resizable: true,
    ellipsis: true,
  },
  {
    title: '备注',
    dataIndex: 'remark',
    // width: 200,
    resizable: true,
    ellipsis: true,
    customRender: ({ record }) => {
      return record.remark || '-';
    },
    edit: true,
    editComponent: 'InputTextArea',
    editComponentProps: {
      autoSize: {
        minRows: 1,
        maxRows: 3,
      },
    },
  },
];

export const noRuleBusinessSearchFormSchema: FormSchema[] = [
  {
    field: 'businessGroup',
    label: '业务组',
    component: 'ApiSelect',
    componentProps: {
      placeholder: '请选择业务组',
      api: GetBusinessGroupOptions,
      labelField: 'name',
      valueField: 'name',
    },
  },
  {
    field: 'businessId',
    label: '业务名称',
    component: 'ApiSelect',
    componentProps: {
      placeholder: '请选择业务名称',
      api: getNoRuleBusiness,
      params: {
        pageSize: 1000,
        pageIndex: 1,
        status: 1,
      },
      labelField: 'businessName',
      valueField: 'businessId',
      resultField: 'items',
      showSearch: true,
      optionFilterProp: 'label',
      filterOption: (input: string, option: any) => {
        return option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'categoryId',
    label: '业务大类',
    component: 'ApiSelect',
    componentProps: {
      placeholder: '请选择业务大类',
      api: GetCategoryList,
      params: {
        pageSize: 1000,
        pageIndex: 1,
        status: 1,
      },
      labelField: 'name',
      valueField: 'id',
      resultField: 'items',
      showSearch: true,
      optionFilterProp: 'label',
      filterOption: (input: string, option: any) => {
        return option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'remark',
    label: '备注',
    component: 'Input',
    componentProps: {
      placeholder: '请输入备注',
    },
    colProps: { span: 6 },
  },
  {
    field: 'isNoRule',
    label: '是否有备注',
    component: 'Switch',
    defaultValue: false,
    componentProps: ({ formActionType }) => ({
      checkedValue: true,
      unCheckedValue: false,
      checkedChildren: '是',
      unCheckedChildren: '否',
      onChange: () => {
        formActionType.submit();
      },
    }),
    colProps: { span: 6 },
  },
];
