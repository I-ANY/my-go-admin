import { getResourceList } from '@/api/sys/resource';
import { BasicColumn, FormSchema } from '@/components/Table';
import { Rule } from 'ant-design-vue/es/form';
import { nextTick } from 'vue';

export enum AuthAllResource {
  YES = 1,
  NO = 0,
}
export enum ResourceIdentify {
  SUBCATEGORY = 'SUBCATEGORY',
}
export enum BusinessPermissionCode {
  GENERAL_PERMISSION = 'GENERAL_PERMISSION',
}

export const columns: BasicColumn[] = [
  {
    title: '角色名称',
    dataIndex: 'name',
    width: 300,
  },
  {
    title: '角色标识',
    dataIndex: 'identify',
    width: 300,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 120,
  },
  {
    title: '备注',
    dataIndex: 'remark',
  },
];

export const searchFormSchema: FormSchema[] = [
  {
    field: 'search',
    label: '搜索',
    component: 'Input',
    colProps: { span: 8 },
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
];

const checkIdentify = async (_rule: Rule, value: string) => {
  // 检查字符串是否只包含字母、数字和下划线
  const regex = /^[a-zA-Z0-9_]+$/;
  // 检查字符串是否为纯数字或纯下划线
  const isPureNumber = /^\d+$/.test(value);
  const isPureUnderscore = /^_+$/.test(value);
  if (value.length < 5) {
    return Promise.reject('长度不能小于5');
  }
  // 校验条件
  if (regex.test(value)) {
    if (isPureNumber || isPureUnderscore) {
      return Promise.reject('不能是纯数字或者下划线');
    } else {
      return Promise.resolve();
    }
  } else {
    return Promise.reject('只能包含字母、数字和下划线');
  }
};

export const formSchema: FormSchema[] = [
  {
    field: 'name',
    label: '角色名称',
    required: true,
    component: 'Input',
  },
  {
    field: 'identify',
    label: '角色标识',
    required: true,
    component: 'Input',
    rules: [{ validator: checkIdentify, trigger: 'change' }],
  },
  {
    field: 'status',
    label: '状态',
    component: 'RadioButtonGroup',
    required: true,
    defaultValue: 1,
    componentProps: {
      options: [
        { label: '启用', value: 1 },
        { label: '禁用', value: 2 },
      ],
    },
  },
  {
    label: '备注',
    field: 'remark',
    component: 'InputTextArea',
  },
  {
    label: ' ',
    field: 'menuIds',
    slot: 'menu',
  },
];

export function getResourcePermissionFormSchema(onResourceTypeChange, getFuncs): FormSchema[] {
  return [
    {
      field: 'resourceTypeId',
      label: '资源类型',
      component: 'ApiSelect',
      componentProps: {
        placeholder: '请选择资源类型',
        api: async () => {
          const { items } = await getResourceList({
            pageSize: 99999,
            pageIndex: 1,
          });
          if (items?.length > 0) {
            const { setFieldsValue } = getFuncs();
            // 选中第一个资源类型
            nextTick(() => {
              setFieldsValue({
                resourceTypeId: items[0].id,
              });
            });
            return items.map((item) => {
              return {
                label: item.name + '（' + item.identify + '）',
                value: item.id,
              };
            });
          } else {
            return [];
          }
        },
        // labelField: 'name',
        // valueField: 'id',
        showSearch: true,
        filterOption: (input: string, option: any) => {
          return option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0;
        },
        allowClear: false,
        onChange: onResourceTypeChange,
      },
      required: true,
      colProps: { span: 8 },
    },
  ];
}

export function getBusinessPermissionFormSchema(onAuthAllBusinessChange): FormSchema[] {
  return [
    {
      field: 'authAllBusiness',
      label: '授权所有业务',
      component: 'Switch',
      defaultValue: null,
      componentProps: {
        checkedValue: AuthAllResource.YES,
        unCheckedValue: AuthAllResource.NO,
        checkedChildren: '是',
        unCheckedChildren: '否',
        onChange: onAuthAllBusinessChange,
      },
      colProps: { span: 24 },
      required: true,
    },
    {
      label: ' ',
      field: 'businessIds',
      slot: 'businessPermission',
    },
  ];
}
