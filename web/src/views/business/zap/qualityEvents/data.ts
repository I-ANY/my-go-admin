import { BasicColumn, FormSchema } from '@/components/Table';

export const searchFormSchema: FormSchema[] = [
  {
    field: 'hostnames',
    label: '主机名',
    component: 'InputTextArea',
    colProps: { span: 6 },
    componentProps: {
      allowClear: true,
      placeholder: '输入主机名进行搜索',
      rows: 2,
    },
  },
  {
    field: 'nodes',
    label: '节点label',
    component: 'InputTextArea',
    colProps: { span: 6 },
    componentProps: {
      allowClear: true,
      placeholder: '输入节点label进行搜索',
      rows: 2,
    },
  },
];
export const columns: BasicColumn[] = [
  {
    title: '告警次数',
    dataIndex: 'count',
    width: 70,
    resizable: true,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 180,
    resizable: true,
  },
  {
    title: '节点label',
    dataIndex: 'node',
    width: 180,
    resizable: true,
  },
  {
    title: '详情',
    dataIndex: 'fields',
    width: 600,
    resizable: true,
  },
  {
    title: '开始告警时间',
    dataIndex: 'first_alert_time',
    width: 150,
    resizable: true,
  },
  {
    title: '结束告警时间',
    dataIndex: 'last_alert_time',
    width: 150,
    resizable: true,
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    width: 120,
    resizable: true,
  },
];

// 定义一个方法来截取前五十个字符
export function truncatedFields(fields: string) {
  const jdata = JSON.stringify(fields);
  return jdata.length > 120 ? jdata.slice(0, 120) + '...' : jdata;
}
