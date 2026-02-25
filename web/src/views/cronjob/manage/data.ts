import { BasicColumn, FormSchema } from '@/components/Table';
import { Rule } from 'ant-design-vue/es/form';
import { CjEnum } from '@/enums/dictTypeCode';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';

export const nodeNameMap = getDictDataMapFromDict(CjEnum.CRON_NODE);
const checkCron = async (_rule: Rule, value: string) => {
  if (!value || value.length < 11) {
    return Promise.reject('非法的cron表达式');
  }
  const strs = value.split(' ');
  const validStrs: string[] = [];
  strs.forEach((e) => {
    validStrs.push(e.trim());
  });

  if (validStrs.length == 6) {
    // 校验条件
    return Promise.resolve();
  } else {
    return Promise.reject('非法的cron表达式');
  }
};

export const columns: BasicColumn[] = [
  {
    title: '名称',
    dataIndex: 'jobName',
    width: 300,
    resizable: true,
  },
  {
    title: '类型',
    dataIndex: 'jobType',
    width: 100,
  },
  {
    title: 'cron表达式',
    dataIndex: 'cronExpression',
    width: 200,
  },
  {
    title: '调度任务',
    dataIndex: 'invokeTarget',
    width: 200,
  },
  {
    title: '调度节点',
    dataIndex: 'scheduleNode',
    width: 200,
  },

  {
    title: '状态',
    dataIndex: 'status',
    width: 120,
  },
  {
    title: '上次执行状态',
    dataIndex: 'runStatus',
    width: 120,
  },
  {
    title: '运行参数',
    dataIndex: 'args',
    width: 200,
  },
];

export const searchFormSchema: FormSchema[] = [
  {
    field: 'jobName',
    label: '名称',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'jobType',
    label: '类型',
    component: 'Select',
    componentProps: {
      options: [
        { label: 'EXEC', value: 2 },
        { label: 'FLASK-API', value: 1 },
      ],
    },
    colProps: { span: 6 },
  },
  {
    field: 'invokeTarget',
    label: '调度任务',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'scheduleNode',
    label: '调度节点',
    component: 'Select',
    colProps: { span: 6 },
    componentProps: {
      options: getSelectOptionsFromDict(CjEnum.CRON_NODE),
      allowClear: true,
    },
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    componentProps: {
      options: [
        { label: '启用', value: 1 },
        { label: '禁用', value: 2 },
      ],
    },
    colProps: { span: 6 },
  },
  {
    field: 'runStatus',
    label: '上次执行状态',
    component: 'Select',
    componentProps: {
      options: [
        { label: '运行中', value: 1 },
        { label: '执行成功', value: 2 },
        { label: '执行失败', value: 3 },
      ],
    },
    colProps: { span: 6 },
  },
];

export function getFormSchema(onJobTypeSelect: any): FormSchema[] {
  return [
    {
      field: 'jobName',
      label: '名称',
      required: true,
      component: 'Input',
    },
    {
      field: 'jobType',
      label: '类型',
      required: true,
      component: 'Select',
      componentProps: {
        options: [
          { label: 'FLASK-API', value: 1 },
          { label: 'EXEC', value: 2 },
        ],
        onSelect: onJobTypeSelect,
      },
    },
    {
      field: 'invokeTarget',
      label: '调度任务',
      required: true,
      component: 'Input',
    },
    {
      field: 'scheduleNode',
      label: '调度节点',
      component: 'Select',
      required: true,
      componentProps: {
        options: getSelectOptionsFromDict(CjEnum.CRON_NODE),
      },
      helpMessage: '创建后该字段值不可修改',
    },
    {
      field: 'cronExpression',
      label: 'cron表达式',
      required: true,
      component: 'Input',
      componentProps: {
        placeholder: '* * * * * * 表示每秒执行一次',
      },
      helpMessage: '* * * * * 表示每秒执行一次',
      rules: [{ validator: checkCron, trigger: 'change' }],
    },
    {
      field: 'status',
      label: '状态',
      component: 'RadioButtonGroup',
      defaultValue: 1,
      componentProps: {
        options: [
          { label: '启用', value: 1 },
          { label: '禁用', value: 2 },
        ],
      },
    },
    {
      field: 'args',
      label: '运行参数',
      component: 'InputTextArea',
      componentProps: {
        rows: 8,
        placeholder: '请输入args',
      },
    },
  ];
}
