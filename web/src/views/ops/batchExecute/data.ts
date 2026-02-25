import { BasicColumn } from '@/components/Table';
import { FormSchema } from '@/components/Form';
import { Tag } from 'ant-design-vue';
import { formatToDateTime } from '@/utils/dateUtil';
import { h } from 'vue';
import { useDebounceFn } from '@vueuse/core';

export const executeSearchFormSchema = (onHostnamesChange): FormSchema[] => [
  {
    field: 'hostnameSelect',
    label: '主机选择',
    component: 'Input',
    required: false,
    colSlot: 'colSlot_selectDevice',
    colProps: {
      span: 24,
    },
  },
  {
    field: 'hostnames',
    label: '主机名',
    component: 'InputTextArea',
    required: true,
    colProps: {
      span: 24,
    },
    componentProps: {
      allowClear: true,
      rows: 5,
      placeholder: '输入主机，一行一个',
      // onBlur: onHostnamesChange,
      onChange: useDebounceFn(onHostnamesChange, 500),
    },
  },
  {
    field: 'functionOptions',
    label: '选择功能',
    component: 'Cascader',
    required: true,
    colProps: {
      span: 24,
    },
    componentProps: {
      options: [],
      allowClear: true,
      showSearch: true, // 启用搜索功能
      placeholder: '选择|搜索 功能菜单',
      dropdownClassName: 'function-options-cascader-popup',
    },
  },
  {
    field: 'functionParams',
    label: '参数',
    component: 'Input',
    colProps: {
      span: 24,
    },
    show: false,
    componentProps: {
      disabled: true,
      allowClear: true,
      placeholder: '',
    },
  },
];

export const deviceSearchFormSchema = (onCategoryIdsChange): FormSchema[] => [
  {
    field: 'hostnames',
    label: '  ',
    component: 'InputTextArea',
    colProps: { span: 6 },
    componentProps: {
      allowClear: true,
      placeholder: '输入主机名进行搜索',
    },
  },
  {
    field: 'categoryId',
    label: ' ',
    component: 'Select',
    slot: 'categoryId',
    componentProps: {
      placeholder: '选择业务大类',
      onChange: onCategoryIdsChange,
      allowClear: true,
      showSearch: true, // 启用搜索功能
      optionLabelProp: 'text', // 回显使用 text 字段，避免显示 VNode
      filterOption: (inputValue, option: any) => {
        const text = option.text || option.label;
        return (
          (typeof text === 'string' ? text : '').toLowerCase().indexOf(inputValue.toLowerCase()) >=
          0
        );
      },
    },
    colProps: { span: 4 },
  },
  {
    field: 'businessIds',
    label: '  ',
    component: 'Select',
    componentProps: {
      options: [],
      allowClear: true,
      mode: 'multiple', // 启用多选模式
      placeholder: '选择业务小类',
      maxTagCount: 10, // 最多显示3个标签，超出显示+X
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
    colProps: { span: 14 },
  },
  {
    field: 'owner',
    label: '  ',
    component: 'Input',
    colProps: { span: 6 },
    componentProps: {
      allowClear: true,
      placeholder: '输入机房进行搜索',
    },
  },
  {
    field: 'status',
    label: '  ',
    component: 'Select',
    componentProps: {
      options: [
        { label: '待审核', value: 0 },
        { label: '已上线', value: 1 },
        { label: '已审核', value: 4 },
        { label: '未通过', value: 5 },
        { label: '休眠', value: 6 },
      ],
      allowClear: true,
      mode: 'multiple', // 启用多选模式
      placeholder: '选择审批状态',
      maxTagCount: 3, // 最多显示3个标签，超出显示+X
    },
    defaultValue: [1, 4],
    colProps: { span: 6 },
  },
  {
    field: 'online',
    label: '  ',
    component: 'Select',
    defaultValue: 1,
    componentProps: {
      options: [
        { label: '离线', value: 0 },
        { label: '在线', value: 1 },
      ],
      allowClear: true,
      placeholder: '选择在线状态',
      maxTagCount: 3, // 最多显示3个标签，超出显示+X
    },
    colProps: { span: 3 },
  },
  {
    field: 'isp',
    label: '  ',
    component: 'Select',
    componentProps: {
      options: [
        { label: '移动', value: '移动' },
        { label: '联通', value: '联通' },
        { label: '电信', value: '电信' },
      ],
      allowClear: true,
      mode: 'multiple', // 启用多选模式
      placeholder: '选择运营商',
      maxTagCount: 3, // 最多显示3个标签，超出显示+X
    },
    colProps: { span: 3 },
  },
];

export const deviceColumns: BasicColumn[] = [
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 100,
    resizable: true,
  },
  {
    title: '机房',
    dataIndex: 'owner',
    width: 80,
    resizable: true,
  },
  {
    title: '业务名',
    dataIndex: 'business',
    width: 80,
    resizable: true,
  },
  {
    title: '运营商',
    dataIndex: 'carrier',
    width: 80,
    resizable: true,
  },
  {
    title: '审核状态',
    dataIndex: 'status',
    width: 80,
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
    },
  },
  {
    title: '设备状态',
    dataIndex: 'online',
    width: 80,
    resizable: true,
    customRender: ({ record }) => {
      const online = record.online;
      if (online === 0) {
        return h(Tag, { color: 'red' }, () => '离线');
      } else {
        return h(Tag, { color: 'green' }, () => '在线');
      }
    },
  },
];

export const deviceExecuteHistoryResultColumns: BasicColumn[] = [
  {
    title: '任务',
    dataIndex: 'taskName',
    width: 250,
    resizable: true,
  },
  {
    title: '返回码',
    dataIndex: 'code',
    width: 50,
    resizable: true,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 60,
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
    title: '结果',
    dataIndex: 'result',
    width: 250,
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
    width: 110,
    resizable: true,
    customRender: ({ record }) => {
      return formatToDateTime(record.createdAt, 'YYYY-MM-DD HH:mm:ss');
    },
  },
];
