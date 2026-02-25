import { BasicColumn, FormSchema } from '@/components/Table';
import { DefaultOptionType, SelectValue } from 'ant-design-vue/es/select';
import { h } from 'vue';


export const columns: BasicColumn[] = [
  {
    title: 'ID',
    dataIndex: 'task_id',
    resizable: true,
    ellipsis: true,
  },
  {
    title: '供应商',
    dataIndex: 'iso_sign',
    width: 100,
    resizable: true,
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
    title: '创建人',
    dataIndex: 'creator',
    width: 120,
    resizable: true,
  },
  {
    title: '构建时间',
    dataIndex: 'updated_at',
    width: 150,
    format: 'date|YYYY-MM-DD HH:mm:ss',
    resizable: true,
  },
  {
    title: 'iso地址',
    dataIndex: 'iso_url',
    minWidth: 100,
    ellipsis: true,
    resizable: true,
    customRender: ({ record }) => {
      const url = record?.iso_url;
      if (!url) return '-';
      return h(
        'a',
        {
          href: url,
          target: '_blank',
          rel: 'noopener noreferrer',
          download: true,
        },
        record.iso_url
      );
    },
  },
  {
    title: 'iso_MD5',
    dataIndex: 'md5',
    width: 120,
    resizable: true,
  },
  {
    title: '任务状态',
    dataIndex: 'current_step',
    width: 120,
    resizable: true,
  },
];

export const searchFormSchema = (basenameList: any[]): FormSchema[] => [
  {
    field: 'iso_sign',
    label: '供应商',
    component: 'Input',
    colProps: { span: 6 },
  },
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
      field: 'iso_sign',
      label: '供应商名称',
      required: true,
      labelWidth:200,
      component: 'Input',
    },
    {
      field: 'basename',
      label: '基础配置名称',
      required: true,
      labelWidth:200,
      component: 'Select',
      componentProps: {
        options: [],
        allowClear: true,
        showSearch: true,
        filterOption: (inputValue: string, option: any) => {
          return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
        },
      },
    },
    {
      field: 'os_type',
      label: '操作系统',
      labelWidth:200,
      required: true,
      component: 'Select',
      componentProps: {
        options: [],
        allowClear: true,
        placeholder: '请先选择基础配置名称',
        showSearch: true,
        filterOption: (inputValue: string, option: any) => {
          return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
        },
      },
    },
    {
      field: 'tar_version',
      label: '业务包版本号',
      labelWidth:200,
      required: true,
      component: 'Select',
      componentProps: {
        options: [],
        allowClear: true,
        placeholder: '请先选择操作系统',
        showSearch: true,
        filterOption: (inputValue: string, option: any) => {
          return option.label.toLowerCase().indexOf(inputValue.toLowerCase()) >= 0;
        },
      },
    },
    {
      field: 'sysdisk_minsize',
      label: '系统盘最小容量',
      component: 'InputNumber',
      labelWidth:200,
      componentProps: {
        defaultValue: 160,
      }
    },
    {
      field: 'pppoe_type',
      label: '拨号类型',
      component: 'Input',
      labelWidth:200,
      componentProps: {
        defaultValue: 'pppoe',
      }
    },
    {
      field:'create_pcdn_index_data',
      label:'是否创建PCDN索引数据分区',
      component: 'Select',
      labelWidth:200,
      componentProps: {
        options:[{label:'是',value:true},{label:'否',value:false}] as unknown as DefaultOptionType[],
        defaultValue: false as unknown as SelectValue,
        showSearch: true,
      }
    },
    {
      field:'auto_register',
      label:'是否自动注册到平台',
      component: 'Select',
      labelWidth:200,
      componentProps: {
        options:[{label:'是',value:true},{label:'否',value:false}] as unknown as DefaultOptionType[],
        defaultValue: false as unknown as SelectValue,
        showSearch: true,
      }
    },
    {
      field:'syunhost_account',
      label:'云主机自动注册厂商ID',
      component: 'Input',
      labelWidth:200,
    },
    {
      field:'yunhost_password',
      label:'YUNHOST_PASSWORD',
      component: 'Input',
      labelWidth:200,
    }
  ];
}
