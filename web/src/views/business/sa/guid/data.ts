import { FormSchema } from '@/components/Table';
import { h } from 'vue';
import { Tag } from 'ant-design-vue';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { EcdnEnum } from '@/enums/dictTypeCode';

export const EcdnServerStatusMap = getDictDataMapFromDict(EcdnEnum.SERVER_STATUS);

export const guidSearchSchema: FormSchema[] = [
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
    field: 'guids',
    label: 'GUID',
    component: 'InputTextArea',
    colProps: { span: 8 },
    componentProps: {
      allowClear: true,
      placeholder: '输入GUID进行搜索',
      rows: 2,
    },
  },
  {
    field: 'guid_status',
    label: 'GUID状态',
    component: 'Select',
    colProps: { span: 5 },
    componentProps: {
      allowClear: true,
      options: [
        { label: '正常', value: 1 },
        { label: '异常', value: 0 },
      ],
    },
    helpMessage: 'hostname是否与GUID中的一致',
  },
  {
    field: 'status',
    label: '状态',
    component: 'Select',
    colProps: { span: 5 },
    componentProps: {
      allowClear: true,
      options: [
        { label: '在线', value: 1 },
        { label: '离线', value: 0 },
      ],
    },
    helpMessage: '机器不正常上报GUID信息，则为离线',
  },
  {
    field: 'status_ecdn',
    label: '状态(ECDN)',
    component: 'Select',
    colProps: { span: 6 },
    componentProps: {
      allowClear: true,
      options: getSelectOptionsFromDict(EcdnEnum.SERVER_STATUS),
    },
  },
];

export const guidTableColumns = [
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 120,
    align: 'left',
  },
  {
    title: 'GUID',
    dataIndex: 'guid',
    width: 250,
    align: 'left',
  },
  {
    title: 'GUID是否异常',
    dataIndex: 'guid_status',
    width: 100,
    helpMessage:
      'hostname是否与GUID中的是否一致，如不一致且有在ECDN打了标签(GUID与主机名不一致)，则也置为正常',
    customRender: ({ record }) => {
      const status = record.guid_status;
      const enable = 1;
      const color = status === enable ? 'green' : 'red';
      const text = status === enable ? '正常' : '异常';
      return h(Tag, { color: color }, () => text);
    },
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 80,
    helpMessage: '机器不正常上报GUID信息，则为离线',
    customRender: ({ record }) => {
      const status = record.status;
      const enable = 1;
      const color = status === enable ? 'green' : 'red';
      const text = status === enable ? '在线' : '离线';
      return h(Tag, { color: color }, () => text);
    },
  },
  {
    title: '状态(ECDN)',
    dataIndex: 'status_ecdn',
    width: 90,
    helpMessage: 'ECDN平台上主机的状态',
    customRender: ({ record }) => {
      const status = record.status_ecdn;
      const color = EcdnServerStatusMap[status].color || 'default';
      const text = EcdnServerStatusMap[status].dictLabel;
      return h(Tag, { color: color }, () => text);
    },
  },
  {
    title: '业务服务',
    dataIndex: 'service_name',
    width: 100,
  },
  {
    title: '服务状态',
    dataIndex: 'service_status',
    width: 100,
    customRender: ({ record }) => {
      const status = record.service_status;
      if (status != '') {
        const color = status == 'active' ? 'green' : 'red';
        return h(Tag, { color: color }, () => status);
      }
    },
  },
  {
    title: '上行带宽(Mbps)',
    dataIndex: 'multi_line_speed',
    width: 110,
    customRender: ({ record }) => {
      return record.multi_line_speed * 8;
    },
  },
  {
    title: '创建时间',
    dataIndex: 'created_at',
    width: 120,
  },
  {
    title: '更新时间',
    dataIndex: 'updated_at',
    width: 120,
  },
];
