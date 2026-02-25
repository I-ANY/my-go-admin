import { BasicColumn, FormSchema } from '@/components/Table';
import { h } from 'vue';
import { Tag } from 'ant-design-vue';
import { formatToDateTime } from '@/utils/dateUtil';

export const executeResultSearchFormSchema: FormSchema[] = [
  {
    field: 'hostname',
    label: '主机名',
    component: 'InputTextArea',
    colProps: { span: 6 },
    componentProps: {
      allowClear: true,
      placeholder: '输入主机名进行搜索',
    },
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    colProps: { span: 6 },
    componentProps: {
      allowClear: true,
      placeholder: '输入状态进行搜索',
      options: [
        {
          label: '执行中',
          value: 0,
        },
        {
          label: '成功',
          value: 1,
        },
        {
          label: '失败',
          value: 2,
        },
      ],
    },
  },
  {
    field: 'operator',
    label: '操作人',
    component: 'Input',
    colProps: { span: 6 },
    componentProps: {
      allowClear: true,
      placeholder: '输入操作人(中文)进行搜索',
    },
  },
];

export const executeResultColumns: BasicColumn[] = [
  // {
  //   title: '任务',
  //   dataIndex: 'taskName',
  //   width: 150,
  //   resizable: true,
  // },
  {
    title: '机房',
    dataIndex: 'owner',
    width: 100,
    resizable: true,
  },
  {
    title: '业务名',
    dataIndex: 'business',
    width: 100,
    resizable: true,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 120,
    resizable: true,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 60,
    resizable: true,
  },
  {
    title: '结果',
    dataIndex: 'result',
    width: 350,
    resizable: true,
  },
  {
    title: '操作人',
    dataIndex: 'operator',
    width: 80,
    resizable: true,
  },
  {
    title: '操作时间',
    dataIndex: 'createdAt',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      return formatToDateTime(record.createdAt, 'YYYY-MM-DD HH:mm:ss');
    },
  },
];

export const executeTaskColumns: BasicColumn[] = [
  {
    title: '业务名',
    dataIndex: 'business',
    width: 80,
    resizable: true,
  },
  {
    title: '操作',
    dataIndex: 'name',
    width: 150,
    resizable: true,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 60,
    resizable: true,
    customRender: ({ record }) => {
      const status = record.status;
      if (status === 0) {
        return h(Tag, { color: 'orange' }, () => '执行中');
      } else if (status === 1) {
        return h(Tag, { color: 'green' }, () => '成功');
      } else {
        return h(Tag, { color: 'red' }, () => '失败');
      }
    },
  },
  {
    title: '结果',
    dataIndex: 'result',
    width: 120,
    resizable: true,
    customRender: ({ record }) => {
      const successText = '成功: ' + record.successCount;
      const failedText = '失败: ' + record.failedCount;
      const runningText =
        '执行中: ' + (record.totalCount - record.successCount - record.failedCount);
      return h(
        'div',
        { style: 'display: flex; gap: 8px;justify-content: center; align-items: center' },
        [
          h(Tag, { color: 'green' }, () => successText),
          h(Tag, { color: 'red' }, () => failedText),
          h(Tag, { color: 'orange' }, () => runningText),
        ],
      );
    },
  },
  {
    title: '详情',
    dataIndex: 'detail',
    width: 80,
    resizable: true,
  },
  {
    title: '操作人',
    dataIndex: 'operator',
    width: 80,
    resizable: true,
  },
  {
    title: '操作时间',
    dataIndex: 'createdAt',
    width: 120,
    resizable: true,
    customRender: ({ record }) => {
      return formatToDateTime(record.createdAt, 'YYYY-MM-DD HH:mm:ss');
    },
  },
];

export const executeDetailSearchFormSchema: FormSchema[] = [
  {
    field: 'hostname',
    label: ' ',
    component: 'InputTextArea',
    colProps: { span: 6 },
    componentProps: {
      allowClear: true,
      placeholder: '输入主机名进行搜索',
    },
  },
  {
    field: 'owner',
    label: ' ',
    component: 'InputTextArea',
    colProps: { span: 4 },
    componentProps: {
      allowClear: true,
      rows: 2,
      placeholder: '输入机房进行搜索',
    },
  },
  {
    field: 'business',
    label: ' ',
    component: 'InputTextArea',
    colProps: { span: 4 },
    componentProps: {
      allowClear: true,
      rows: 2,
      placeholder: '输入业务名进行搜索',
    },
  },
  {
    field: 'status',
    label: ' ',
    component: 'Select',
    colProps: { span: 4 },
    componentProps: {
      options: [
        {
          label: '执行中',
          value: 0,
        },
        {
          label: '成功',
          value: 1,
        },
        {
          label: '失败',
          value: 2,
        },
      ],
      allowClear: true,
      placeholder: '选择状态',
    },
  },
];

export const executeDetailColumns: BasicColumn[] = [
  {
    title: '机房',
    dataIndex: 'owner',
    width: 100,
    resizable: true,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 120,
    resizable: true,
  },
  {
    title: '业务名',
    dataIndex: 'business',
    width: 100,
    resizable: true,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 40,
    resizable: true,
    customRender: ({ record }) => {
      let text = '';
      let color = '';
      if (record.status === 0) {
        text = '执行中';
        color = 'orange';
      }
      if (record.status === 1) {
        text = '成功';
        color = 'green';
      }
      if (record.status === 2) {
        text = '失败';
        color = 'red';
      }
      return h(Tag, { color: color }, () => text);
    },
  },
  {
    title: '结果详情',
    dataIndex: 'result',
    width: 400,
    resizable: true,
  },
  {
    title: '操作人',
    dataIndex: 'operator',
    width: 80,
    resizable: true,
  },
  {
    title: '操作时间',
    dataIndex: 'createdAt',
    width: 100,
    resizable: true,
    customRender: ({ record }) => {
      return formatToDateTime(record.createdAt, 'YYYY-MM-DD HH:mm:ss');
    },
  },
];
