import { BasicColumn, FormSchema } from '@/components/Table';

import { Tag } from 'ant-design-vue';
import { formatToDateTime } from '@/utils/dateUtil';
import { h } from 'vue';

export const taskConfigSearchSchema: FormSchema[] = [
  {
    field: 'scriptName',
    label: '',
    component: 'Input',
    colProps: { span: 6 },
    componentProps: {
      placeholder: '脚本名称',
    },
  },
];

export const taskConfigColumns: BasicColumn[] = [
  {
    title: 'ID',
    dataIndex: 'id',
    width: 80,
    ifShow: false,
  },
  {
    title: '操作组ID',
    dataIndex: 'parentId',
    width: 80,
    ifShow: false,
  },
  {
    title: '名称',
    dataIndex: 'name',
    width: 200,
  },
  {
    title: '脚本目录',
    dataIndex: 'scriptPath',
    width: 100,
  },
  {
    title: '脚本名称',
    dataIndex: 'scriptName',
    width: 180,
  },
  {
    title: '参数配置详情',
    dataIndex: 'params',
    width: 200,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 80,
    customRender: ({ record }) => {
      return h(
        Tag,
        { color: record.status === 1 ? 'green' : 'red' },
        record.status === 1 ? '启用' : '禁用',
      );
    },
  },
  {
    title: '台数限制',
    dataIndex: 'maxHosts',
    width: 80,
  },
  {
    title: '超时时间',
    dataIndex: 'timeout',
    width: 80,
  },
  {
    title: '操作人',
    dataIndex: 'operator',
    width: 120,
  },
  {
    title: '创建时间',
    dataIndex: 'createdAt',
    width: 150,
    customRender: ({ record }) => {
      return formatToDateTime(record.createdAt, 'YYYY-MM-DD HH:mm:ss');
    },
  },
  {
    title: '备注',
    dataIndex: 'description',
    width: 180,
  },
];
// 任务配置表单
export const configFormSchema: FormSchema[] = [
  {
    field: 'id',
    label: 'ID',
    component: 'InputNumber',
    show: false,
  },
  {
    field: 'businessId',
    label: '业务大类',
    component: 'Select',
    required: true,
    colProps: { span: 12 },
    componentProps: {
      placeholder: '请选择所属业务',
    },
  },
  {
    field: 'parentId',
    label: '操作组',
    component: 'Select',
    required: true,
    colProps: { span: 12 },
    componentProps: {
      placeholder: '请选择所操作组',
    },
  },
  {
    field: 'name',
    label: '名称',
    component: 'Input',
    required: true,
    colProps: { span: 24 },
    componentProps: {
      placeholder: '请输入任务名称',
    },
  },
  {
    field: 'scriptPath',
    label: '脚本目录',
    component: 'Input',
    required: true,
    colProps: { span: 12 },
    componentProps: {
      placeholder: '只需要输入对应业务的子目录名称即可，如公共脚本，填写general',
    },
  },
  {
    field: 'scriptName',
    label: '脚本名称',
    component: 'Input',
    required: true,
    colProps: { span: 12 },
    componentProps: {
      placeholder: '请输入任务名称',
    },
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    required: true,
    colProps: { span: 8 },
    componentProps: {
      placeholder: '选择状态',
      defaultValue: 1,
      options: [
        {
          label: '启用',
          value: 1,
        },
        {
          label: '禁用',
          value: 0,
        },
      ],
    },
  },
  {
    field: 'maxHosts',
    label: '台数限制',
    component: 'InputNumber',
    colProps: { span: 8 },
    componentProps: {
      placeholder: '最大可执行台数，默认0，不限制',
    },
  },
  {
    field: 'timeout',
    label: '超时时间',
    component: 'InputNumber',
    colProps: { span: 8 },
    componentProps: {
      placeholder: '执行超时时间，默认300s',
    },
  },
  {
    field: 'params',
    label: '操作详情',
    // component: 'Input',
    slot: 'params',
    colProps: { span: 24 },
  },
  {
    field: 'description',
    label: '备注',
    component: 'InputTextArea',
    colProps: { span: 24 },
    componentProps: {
      placeholder: '备注信息',
      rows: 2,
    },
  },
];

export const optionGroupFormSchema: FormSchema[] = [
  // {
  //   field: 'businessID',
  //   label: '所属业务',
  //   component: 'ApiSelect',
  //   required: true,
  //   colProps: { span: 20 },
  //   componentProps: {
  //     api: GetAuthedBiz,
  //     resultField: 'categories',
  //     labelField: 'name',
  //     valueField: 'id',
  //     placeholder: '请选择所属业务',
  //   },
  // },
  {
    field: 'name',
    label: '组名',
    component: 'Input',
    required: true,
    colProps: { span: 20 },
    componentProps: {
      placeholder: '请输入组名称',
    },
  },
  {
    field: 'sort',
    label: '排序',
    component: 'InputNumber',
    required: false,
    colProps: { span: 20 },
    helpMessage: '数字越大越靠前',
    defaultValue: 0,
    componentProps: {
      placeholder: '请输入组名称',
    },
  },
  {
    field: 'type',
    label: '类型',
    component: 'Input',
    defaultValue: 'group',
    show: false,
  },
  {
    field: 'parentId',
    label: '父级ID',
    component: 'Input',
    show: false,
    defaultValue: 0,
  },
];

// 拷贝配置表单
export const copyFormSchema: FormSchema[] = [
  {
    field: 'targetBizId',
    label: '目标业务',
    component: 'Select',
    required: true,
    colProps: { span: 24 },
    componentProps: {
      placeholder: '请选择目标业务',
    },
  },
  {
    field: 'targetParentId',
    label: '目标操作组',
    component: 'Select',
    required: true,
    colProps: { span: 24 },
    componentProps: {
      placeholder: '请选择目标操作组',
    },
  },
  {
    field: 'scriptPath',
    label: '脚本目录',
    component: 'Input',
    required: false,
    colProps: { span: 24 },
    componentProps: {
      placeholder: '可选，留空则保留原来目录名',
    },
  },
];
