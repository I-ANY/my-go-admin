import { BasicColumn, FormSchema } from '@/components/Table';
import { InspectEnum } from '@/enums/dictTypeCode';
import { RangePickPresetsExact } from '@/utils/common';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';

export const operatorTypeOptions = getSelectOptionsFromDict(InspectEnum.OPERATOR_TYPE);
export const operatorTypeMap = getDictDataMapFromDict(InspectEnum.OPERATOR_TYPE);

export const execStatusOptions = getSelectOptionsFromDict(InspectEnum.EXEC_STATUS);
export const execStatusMap = getDictDataMapFromDict(InspectEnum.EXEC_STATUS);

export const searchFormSchema = (onTimePikerOpen): FormSchema[] => [
  {
    field: '[startTimeBegin, startTimeEnd]',
    label: '操作时间',
    component: 'RangePicker',
    componentProps: {
      allowClear: false,
      format: 'YYYY-MM-DD HH:mm:ss',
      showTime: { format: 'HH:mm:ss' },
      placeholder: ['开始时间', '结束时间'],
      style: {
        width: '100%',
      },
      presets: RangePickPresetsExact(),
      onOpenChange: onTimePikerOpen,
    },
    required: true,
    colProps: { span: 6 },
  },
  {
    field: 'taskName',
    label: '任务名称',
    component: 'Select',
    componentProps: {
      options: [],
      allowClear: true,
      showSearch: true,
      filterOption: (input, option) => {
        return (option?.label ?? '').toLowerCase().includes(input.toLowerCase());
      },
    },
    colProps: { span: 6 },
  },
  {
    field: 'operatorType',
    label: '操作类型',
    component: 'Select',
    componentProps: {
      options: operatorTypeOptions,
      allowClear: true,
      // mode: 'multiple',
    },
    colProps: { span: 6 },
  },
  {
    field: 'status',
    label: '执行状态',
    component: 'Select',
    componentProps: {
      options: execStatusOptions,
      allowClear: true,
    },
    colProps: { span: 6 },
  },
];

export const columns: BasicColumn[] = [
  {
    title: '任务ID',
    dataIndex: 'id',
    width: 200,
    resizable: true,
  },
  {
    title: '任务名称',
    dataIndex: 'taskName',
    width: 350,
    resizable: true,
  },
  {
    title: '操作时间',
    dataIndex: 'startTime',
    width: 180,
    resizable: true,
  },
  {
    title: '完成时间',
    dataIndex: 'finishTime',
    width: 180,
    resizable: true,
  },
  {
    title: '操作类型',
    dataIndex: 'operatorType',
    width: 140,
    resizable: true,
  },
  {
    title: '执行状态',
    dataIndex: 'status',
    width: 140,
    resizable: true,
  },
  {
    title: '设备总数',
    dataIndex: 'serverCount',
    width: 100,
    resizable: true,
  },
  {
    title: '正常/异常数量',
    dataIndex: 'count',
    width: 120,
    resizable: true,
  },
  // {
  //   title: '正常数量',
  //   dataIndex: 'normalCount',
  //   width: 100,
  //   resizable: true,
  // },
  // {
  //   title: '异常数量',
  //   dataIndex: 'abnormalCount',
  //   width: 100,
  //   resizable: true,
  // },
  {
    title: '操作人',
    dataIndex: 'nickName',
    // width: 200,
    // resizable: true,
    customRender: ({ record }) => {
      return record?.operatorUser?.nickName || '';
    },
  },
];
