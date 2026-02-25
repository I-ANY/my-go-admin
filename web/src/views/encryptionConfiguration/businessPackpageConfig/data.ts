import { BasicColumn, FormSchema } from '@/components/Table';


export const columns: BasicColumn[] = [
  {
    title: '序号',
    dataIndex: 'index',
    width: 50,
    customRender: ({ record, index }) => {
      return index + 1;
    },
  },
  {
    title: '基础配置名称',
    dataIndex: 'basename',
    width: 200,
    resizable: true,
  },
  {
    title: '操作系统',
    dataIndex: 'os_type',
    width: 100,
    resizable: true,
  },
  {
    title: '业务包版本号',
    dataIndex: 'tar_version',
    width: 100,
    resizable: true,
  },
  {
    title: 'tar包地址',
    dataIndex: 'tar_url',
    width: 200,
    resizable: true,
  },
  {
    title: 'tar包名称',
    dataIndex: 'tar_name',
    width: 200,
    resizable: true,
  },

  {
    title: '操作人',
    dataIndex: 'operator',
    width: 120,
    resizable: true,
  },
  {
    title: '更新时间',
    dataIndex: 'updated_at',
    width: 150,
    format: 'date|YYYY-MM-DD HH:mm:ss',
    resizable: true,
  },
];

export const searchFormSchema = (basenameList: any[]): FormSchema[] => [
  {
    field: 'basename',
    label: '基础配置名称',
    component: 'Select',
    colProps: { span: 6 },
    componentProps: {
      options: basenameList,
      allowClear: true,
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
  },
  {
    field: 'os_type',
    label: '操作系统',
    component: 'Select',
    colProps: { span: 6 },
    componentProps: {
      options: [],
      allowClear: true,
      placeholder: '请先选择基础配置名称',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
  },
  {
    field: 'tar_version',
    label: '版本号',
    component: 'Select',
    colProps: { span: 6 },
    componentProps: {
      options: [],
      allowClear: true,
      placeholder: '请先选择操作系统和基础配置名称',
      filterOption: (inputValue, option: any) => {
        return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
      },
    },
  },
];

export function getFormSchema(): FormSchema[] {
  return [
    {
      field: 'basename',
      label: '基础配置名称',
      required: true,
      component: 'Input',
    },
    {
      field: 'os_type',
      label: '操作系统',
      required: true,
      component: 'Input',
    },
    {
      field: 'tar_version',
      label: '业务包版本号',
      required: true,
      component: 'Input',
    },
    {
      field: 'tar_url',
      label: 'tar包地址',
      component: 'Input',
      required: true,
    },
    {
      field: 'tar_name',
      label: 'tar包名称',
      required: true,
      component: 'Input',
    }
  ];
}
