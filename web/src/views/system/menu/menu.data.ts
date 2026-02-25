import { BasicColumn, FormSchema } from '@/components/Table';
import { h } from 'vue';
import { Tag } from 'ant-design-vue';
import Icon from '@/components/Icon/Icon.vue';
import { useI18n } from '@/hooks/web/useI18n';
import { getApiList } from '@/api/sys/api';

export const columns: BasicColumn[] = [
  {
    title: '菜单标题',
    dataIndex: 'title',
    width: 200,
    align: 'left',
    customRender({ record }) {
      const { t } = useI18n();
      return t(record.meta.title);
    },
  },
  {
    title: '菜单标识',
    dataIndex: 'name',
    width: 250,
    resizable: true,
    align: 'left',
  },
  {
    title: '图标',
    dataIndex: 'icon',
    width: 50,
    customRender: ({ record }) => {
      if (record.meta.icon) {
        return h(Icon, { icon: record.meta.icon });
      } else {
        return '';
      }
    },
  },
  {
    title: '权限标识',
    dataIndex: 'permission',
    width: 180,
  },
  {
    title: '组件',
    dataIndex: 'component',
  },
  {
    title: '排序',
    dataIndex: 'orderNo',
    width: 60,
    customRender({ record }) {
      return record.meta.orderNo;
    },
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 80,
    customRender: ({ record }) => {
      const status = record.status;
      const enable = ~~status === 0;
      const color = enable ? 'green' : 'red';
      const text = enable ? '启用' : '禁用';
      return h(Tag, { color: color }, () => text);
    },
  },
  {
    title: '类型',
    dataIndex: 'type',
    width: 180,
  },
];

const isDir = (type: number) => type === 0;
const isMenu = (type: number) => type === 1;
const isButton = (type: number) => type === 2;

export const searchFormSchema: FormSchema[] = [
  {
    field: 'menuName',
    label: '菜单名称',
    component: 'Input',
    colProps: { span: 8 },
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    componentProps: {
      options: [
        { label: '启用', value: 0 },
        { label: '禁用', value: 1 },
      ],
    },
    colProps: { span: 8 },
  },
];

export const formSchema: FormSchema[] = [
  {
    field: 'type',
    label: '菜单类型',
    component: 'RadioButtonGroup',
    defaultValue: 0,
    componentProps: {
      options: [
        { label: '目录', value: 0 },
        { label: '菜单', value: 1 },
        { label: '按钮', value: 2 },
      ],
    },
    colProps: { lg: 24, md: 24 },
  },
  {
    field: 'title',
    label: '菜单标题',
    component: 'Input',
    required: true,
  },
  {
    field: 'name',
    label: '菜单标识',
    component: 'Input',
    required: true,
  },
  {
    field: 'parentMenu',
    label: '上级菜单',
    component: 'TreeSelect',
    componentProps: {
      fieldNames: {
        label: 'menuTitle',
        value: 'id',
      },
      getPopupContainer: () => document.body,
    },
    ifShow: ({ values }) => !isDir(values.type),
  },
  {
    field: 'orderNo',
    label: '排序',
    component: 'InputNumber',
    required: true,
  },
  {
    field: 'redirect',
    label: '重定向',
    component: 'Input',
    ifShow: ({ values }) => !isButton(values.type),
  },
  {
    field: 'icon',
    label: '图标',
    component: 'IconPicker',
    ifShow: ({ values }) => !isButton(values.type),
  },
  {
    field: 'path',
    label: '路由地址',
    component: 'Input',
    // required: true,
    ifShow: ({ values }) => !isButton(values.type),
  },
  {
    field: 'component',
    label: '组件路径',
    component: 'Input',
    ifShow: ({ values }) => !isButton(values.type),
  },
  {
    field: 'permission',
    label: '权限标识',
    component: 'Input',
    // required: true,
    // ifShow: ({ values }) => !isDir(values.type),
  },
  {
    field: 'status',
    label: '状态',
    component: 'RadioButtonGroup',
    defaultValue: 0,
    componentProps: {
      options: [
        { label: '启用', value: 0 },
        { label: '禁用', value: 1 },
      ],
    },
  },
  {
    field: 'ignoreKeepAlive',
    label: '忽略缓存缓存',
    component: 'RadioButtonGroup',
    defaultValue: true,
    componentProps: {
      options: [
        { label: '否', value: false },
        { label: '是', value: true },
      ],
    },
    ifShow: ({ values }) => isMenu(values.type),
  },
  {
    field: 'hideMenu',
    label: '是否隐藏',
    component: 'RadioButtonGroup',
    defaultValue: false,
    componentProps: {
      options: [
        { label: '是', value: true },
        { label: '否', value: false },
      ],
    },
    ifShow: ({ values }) => !isButton(values.type),
  },
  {
    field: 'hideChildrenInMenu',
    label: '隐藏子菜单',
    component: 'RadioButtonGroup',
    defaultValue: false,
    componentProps: {
      options: [
        { label: '否', value: false },
        { label: '是', value: true },
      ],
    },
    ifShow: ({ values }) => !isButton(values.type),
  },
  {
    field: 'apiIds',
    label: 'API接口',
    component: 'ApiTransfer',
    colProps: {
      span: 24,
      lg: 24,
      md: 24,
    },
    componentProps: {
      api: getApiList,
      listStyle: {
        width: '50%',
        height: '300px',
      },
      titles: ['未选择', '已选择'],
      showSearch: true,
      labelField: 'handler',
      valueField: 'id',
      resultField: 'items',
      immediate: false,
      params: { pageSize: 10000, pageIndex: 1 },
      filterOption: (input: string, option: any) => {
        return option.title.toLowerCase().indexOf(input.toLowerCase()) >= 0;
      },
    },
  },
];
