import { getRoleListByPage } from '@/api/demo/system';
import { BasicColumn, FormSchema } from '@/components/Table';
import { SysEnum } from '@/enums/dictTypeCode';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
/**
 * transform mock data
 * {
 *  0: '华东分部',
 * '0-0': '华东分部-研发部'
 * '0-1': '华东分部-市场部',
 *  ...
 * }
 */
export const deptMap = (() => {
  // const pDept = ['华东分部', '华南分部', '西北分部'];
  // const cDept = ['研发部', '市场部', '商务部', '财务部'];
  // return pDept.reduce((map, p, pIdx) => {
  //   map[pIdx] = p;
  //   cDept.forEach((c, cIndex) => (map[`${pIdx}-${cIndex}`] = `${p}-${c}`));
  //   return map;
  // }, {});
})();

export const userFromMap = getDictDataMapFromDict(SysEnum.USER_FROM);

export const columns: BasicColumn[] = [
  {
    title: '用户名',
    dataIndex: 'username',
    width: 120,
  },
  {
    title: '姓名',
    dataIndex: 'nickName',
    width: 120,
  },
  {
    title: '邮箱',
    dataIndex: 'email',
    width: 180,
  },
  {
    title: '手机号',
    dataIndex: 'tel',
    width: 160,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 120,
  },
  {
    title: '角色',
    dataIndex: 'roles',
    width: 260,
    resizable: true,
  },
  {
    title: '用户来源',
    dataIndex: 'source',
    width: 160,
  },
  {
    title: '所属部门',
    dataIndex: 'dept',
    width: 160,
  },
  // {
  //   title: '备注',
  //   dataIndex: 'remark',
  // },
];

export const searchFormSchema: FormSchema[] = [
  // {
  //   field: 'search',
  //   label: '搜索',
  //   component: 'Input',
  //   colProps: { span: 8 },
  // },
  {
    field: 'username',
    label: '用户名',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'nickName',
    label: '姓名',
    component: 'Input',
    colProps: { span: 6 },
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
      allowClear: true,
    },
    colProps: { span: 6 },
  },

  {
    field: 'email',
    label: '邮箱',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'tel',
    label: '手机号',
    component: 'Input',
    colProps: { span: 6 },
  },
  {
    field: 'source',
    label: '用户来源',
    component: 'Select',
    componentProps: {
      options: getSelectOptionsFromDict(SysEnum.USER_FROM),
      allowClear: true,
    },
    colProps: { span: 6 },
  },
  {
    field: 'roleId',
    label: '角色',
    component: 'ApiSelect',
    componentProps: {
      api: getRoleListByPage,
      params: { pageSize: 500, pageIndex: 1 },
      resultField: 'items',
      labelField: 'name',
      valueField: 'id',
      showSearch: true,
      // mode: 'multiple',
      allowClear: true,
      filterOption: (input: string, option: any) => {
        return (
          option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
          option.identify?.toLowerCase().indexOf(input.toLowerCase()) >= 0
        );
      },
    },
    colProps: { span: 6 },
  },
];

export const accountFormSchema: FormSchema[] = [
  {
    field: 'username',
    label: '用户名',
    component: 'Input',
    required: true,
    rules: [{ min: 4, message: '长度必须大于4', trigger: 'blur' }],
    componentProps: {
      autocomplete: 'new-user',
    },

    // helpMessage: ['本字段演示异步验证', '不能输入带有admin的用户名'],
    // rules: [
    //   {
    //     required: true,
    //     message: '请输入用户名',
    //   },
    //   {
    //     trigger: 'blur',
    //     validator(_, value) {
    //       return new Promise((resolve, reject) => {
    //         if (!value) return resolve();
    //         isAccountExist(value)
    //           .then(resolve)
    //           .catch((err) => {
    //             reject(err.message || '验证失败');
    //           });
    //       });
    //     },
    //   },
    // ],
  },
  {
    field: 'nickName',
    label: '姓名',
    rules: [{ min: 2, message: '长度必须大于2', trigger: 'change' }],
    component: 'Input',
    required: true,
    componentProps: {
      autocomplete: 'new-user',
    },
  },
  {
    label: '角色',
    field: 'roleIds',
    component: 'ApiSelect',
    componentProps: {
      api: getRoleListByPage,
      params: { pageSize: 500, pageIndex: 1 },
      resultField: 'items',
      labelField: 'name',
      valueField: 'id',
      mode: 'multiple',
      filterOption: (input: string, option: any) => {
        return (
          option.label?.toLowerCase().indexOf(input.toLowerCase()) >= 0 ||
          option.identify?.toLowerCase().indexOf(input.toLowerCase()) >= 0
        );
      },
    },
    required: true,
  },
  {
    field: 'deptId',
    label: '所属部门',
    component: 'TreeSelect',
    componentProps: {
      fieldNames: {
        label: 'name',
        value: 'id',
      },
      getPopupContainer: () => document.body,
      treeDefaultExpandAll: true,
    },
    required: true,
  },
  {
    label: '邮箱',
    field: 'email',
    component: 'Input',
    required: true,
    rules: [{ type: 'email', trigger: 'change' }],
  },
  {
    label: '手机号',
    field: 'tel',
    component: 'Input',
    rules: [{ trigger: 'change', min: 11, max: 14, message: '请输入合法的手机号' }],
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
    required: true,
  },
  {
    field: 'password',
    label: '密码',
    component: 'InputPassword',
    ifShow: true,
    componentProps: {
      autocomplete: 'new-password',
    },
  },
];
