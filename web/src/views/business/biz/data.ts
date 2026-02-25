import {
  GetCategoryListAll,
  GetSubcategoryFilterList,
  GetSubcategoryListAll,
} from '@/api/business/biz';

import { FormSchema } from '@/components/Table';
import dayjs from 'dayjs';

export enum isVirtualEnum {
  YES = 1,
  NO = 0,
}

export const categorySchemas: FormSchema[] = [
  {
    field: 'name',
    label: '业务组名称',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'code',
    label: '业务组编码',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    colProps: { span: 6 },
    componentProps: {
      options: [
        {
          label: '启用',
          value: '1',
        },
        {
          label: '禁用',
          value: '0',
        },
      ],
    },
  },
  {
    field: 'isVirtual',
    label: '虚拟业务组',
    component: 'Select',
    colProps: { span: 6 },
    componentProps: {
      options: [
        {
          label: '是',
          value: isVirtualEnum.YES,
        },
        {
          label: '否',
          value: isVirtualEnum.NO,
        },
      ],
    },
  },
];

export const categoryColumns = [
  {
    title: '业务组',
    dataIndex: 'name',
    width: 150,
  },
  {
    title: '编码',
    dataIndex: 'code',
    width: 150,
    customRender: ({ record }) => {
      // 将子分类名称用逗号连接
      return record.code || '\\';
    },
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 80,
  },
  {
    title: '虚拟业务组',
    dataIndex: 'isVirtual',
    width: 120,
    customRender: ({ record }) => {
      return record.isVirtual == isVirtualEnum.YES ? '是' : '否';
    },
  },
  {
    title: '子业务详情',
    dataIndex: 'subcategories',
  },

  {
    title: '备注',
    dataIndex: 'describe',
    width: 80,
  },
];

export const subcategorySchemas: FormSchema[] = [
  {
    field: 'name',
    label: '',
    component: 'Input',
    colProps: { span: 24 },
    componentProps: {
      placeholder: '请输入子业务名称',
    },
  },
  {
    field: 'status',
    label: '',
    component: 'Select',
    colProps: { span: 8 },
    componentProps: {
      allowClear: true,
      placeholder: '请选择状态',
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
];

export const subcategoryColumns = [
  {
    title: '业务组',
    dataIndex: 'category',
    width: 100,
    customRender: ({ record }) => {
      return record.category.name;
    },
  },
  {
    title: '业务名称',
    dataIndex: 'name',
    width: 150,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 60,
  },
];

export function drawerFormSchema(onIsVirtualChange): FormSchema[] {
  return [
    {
      field: 'name',
      label: '业务组',
      required: true,
      component: 'Input',
      colProps: { span: 12 },
    },
    {
      field: 'code',
      label: '编码',
      component: 'Input',
      required: true,
      componentProps: {
        disabled: true,
      },
      colProps: { span: 12 },
    },
    {
      field: 'status',
      label: '状态',
      component: 'Switch',
      colProps: { span: 12 },
      componentProps: {
        checkedChildren: '启用',
        unCheckedChildren: '禁用',
        checkedValue: 1,
        unCheckedValue: 0,
      },
      defaultValue: 1,
      required: true,
    },
    {
      field: 'isVirtual',
      label: '虚拟业务组',
      component: 'Switch',
      colProps: { span: 12 },
      componentProps: {
        checkedChildren: '是',
        unCheckedChildren: '否',
        checkedValue: isVirtualEnum.YES,
        unCheckedValue: isVirtualEnum.NO,
        onChange: onIsVirtualChange,
      },
      defaultValue: isVirtualEnum.NO,
      required: false,
    },
    {
      label: '备注',
      field: 'describe',
      component: 'InputTextArea',
      colProps: { span: 24 },
    },
    {
      field: 'subcategories',
      label: '子业务',
      component: 'ApiTransfer',
      colProps: {
        span: 24,
        lg: 24,
        md: 24,
      },
      componentProps: {
        api: GetSubcategoryFilterList,
        listStyle: {
          width: '50%',
          height: '450px',
        },
        titles: ['未选择', '已选择'],
        showSearch: true,
        showSelectAll: true,
        labelField: 'name',
        valueField: 'id',
        resultField: 'items',
        immediate: false,
        filterOption: (input: string, option: any) => {
          return option.title.toLowerCase().indexOf(input.toLowerCase()) >= 0;
        },
      },
    },
    {
      field: 'virtualSubcategories',
      label: '虚拟子业务',
      component: 'ApiTransfer',
      colProps: {
        span: 24,
        lg: 24,
        md: 24,
      },
      componentProps: {
        api: GetSubcategoryListAll,
        params: { pageSize: 10000, pageIndex: 1 },
        listStyle: {
          width: '50%',
          height: '450px',
        },
        showSelectAll: true,
        titles: ['未选择', '已选择'],
        showSearch: true,
        labelField: 'name',
        valueField: 'id',
        resultField: 'items',
        immediate: false,
        filterOption: (input: string, option: any) => {
          return option.title.toLowerCase().indexOf(input.toLowerCase()) >= 0;
        },
      },
    },
  ];
}

// 业务组相关
export const groupSchemas: FormSchema[] = [
  {
    field: 'name',
    label: '业务组名称',
    component: 'Input',
    colProps: { span: 6 },
    componentProps: {
      placeholder: '请输入业务组名称',
    },
  },
];

export const groupColumns = [
  {
    title: '业务组名称',
    dataIndex: 'name',
    width: 200,
  },
  {
    title: '关联业务大类',
    dataIndex: 'categoryIds',
    customRender: ({ record }) => {
      const names = record.categoryNames || [];
      if (names.length === 0) return '-';
      return names.join(', ');
    },
  },
  {
    title: '创建时间',
    dataIndex: 'createdAt',
    width: 180,
    customRender: ({ record }) => {
      if (!record.createdAt) return '-';
      return dayjs(record.createdAt).format('YYYY-MM-DD HH:mm:ss');
    },
  },
  {
    title: '更新时间',
    dataIndex: 'updatedAt',
    width: 180,
    customRender: ({ record }) => {
      if (!record.updatedAt) return '-';
      return dayjs(record.updatedAt).format('YYYY-MM-DD HH:mm:ss');
    },
  },
];

export const groupModalFormSchema: FormSchema[] = [
  {
    field: 'name',
    label: '业务组名称',
    component: 'Input',
    required: true,
    componentProps: {
      placeholder: '请输入业务组名称',
    },
  },
  {
    field: 'categoryIds',
    label: '关联业务大类',
    component: 'ApiSelect',
    required: false,
    colProps: {
      span: 24,
    },
    componentProps: {
      options: [],
      showSearch: true,
      mode: 'multiple',
      placeholder: '请选择关联的业务大类',
      api: async () => {
        const data = await GetCategoryListAll({});
        return data.items.map((item: any) => ({
          label: item.name,
          value: item.id,
        }));
      },
    },
  },
];
