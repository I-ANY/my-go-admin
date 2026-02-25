import { BasicColumn, FormSchema } from '@/components/Table';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';

// import { BizType } from '@/views/business/k/data';
import { KEnum } from '@/enums/dictTypeCode';
import { commonCustomHeaderCell } from '@/utils/util';
import dayjs from 'dayjs';
import { h } from 'vue';

export const IsCoverDiffIspMap: { [key: string]: string } = {
  '0': '否',
  '1': '是',
  '-1': '未知',
};

export const IsProvinceScheduling: { [key: string]: string } = {
  '0': '不限制',
  '1': '仅本省调度',
  '-1': '未知',
};

export const OccupyStatusMap = getDictDataMapFromDict(KEnum.OCCUPY_STATUS);
export const bizTypesMap = getDictDataMapFromDict(KEnum.BIZ_TYPE);
export const providersMap = getDictDataMapFromDict(KEnum.PRODIVER);

export function demandColumns(bizType: string): BasicColumn[] {
  return [
    {
      title: '需求ID',
      dataIndex: 'demand_id',
      width: 250,
    },
    {
      title: '需求类型',
      dataIndex: 'demand_type_name',
      width: 120,
      ifShow: () => bizType === 'specialLine',
    },
    {
      title: '设备类型',
      dataIndex: 'dev_name',
      width: 120,
    },
    {
      title: '需求带宽Mb',
      dataIndex: 'demand_bw',
      width: 100,
    },
    {
      title: '已交付带宽Mb',
      dataIndex: 'delivered_bw',
      width: 100,
    },
    {
      title: '缺口带宽Mb',
      dataIndex: 'gap_bw',
      width: 100,
    },
    {
      title: '区域',
      dataIndex: 'area_name',
      width: 80,
    },
    {
      title: '省份',
      dataIndex: 'province_name',
      width: 80,
    },
    {
      title: '运营商',
      dataIndex: 'isp_name',
      width: 80,
    },
    {
      title: '调度控制',
      dataIndex: 'is_province_scheduling',
      width: 90,
      customRender: ({ text }) => {
        if (text === null) return '-';
        return IsProvinceScheduling[text];
      },
    },
    {
      title: '异网覆盖',
      dataIndex: 'is_cover_diff_isp',
      width: 80,
      customRender: ({ text }) => {
        if (text === null) return '-';
        return IsCoverDiffIspMap[text];
      },
    },
    {
      title: '需求开始时间',
      dataIndex: 'start_time',
      width: 100,
    },
    {
      title: '需求结束时间',
      dataIndex: 'end_time',
      width: 100,
      customRender: ({ text }) => {
        const endDate = dayjs(text);
        const today = dayjs();
        const diffDays = endDate.diff(today, 'day');

        // 过期显示灰色
        if (diffDays < 0) {
          return h('span', { style: { color: 'gray' } }, text);
        }
        // 7天内显示红色
        if (diffDays <= 7) {
          return h('span', { style: { color: 'red', fontWeight: 500 } }, text);
        }
        return text;
      },
    },
    {
      title: '锁定状态',
      dataIndex: 'is_locked',
      width: 80,
      customRender: ({ text, record }) => {
        return text === true
          ? h('span', { style: { color: 'red' } }, `已锁定 ${record.locked_bw}M`)
          : h('span', { style: { color: 'green' } }, '未锁定');
      },
    },
    {
      title: '设备类型',
      dataIndex: 'dev_type_name',
      width: 80,
    },
    {
      title: '机房',
      dataIndex: 'idc_name',
      width: 80,
    },
    {
      title: '是否独\n立部署',
      dataIndex: 'is_independent_deploy',
      width: 80,
      customHeaderCell: commonCustomHeaderCell(),
      customRender: ({ text }) => {
        return text === true ? '是' : '否';
      },
    },
    {
      title: '验收中\n任务数',
      dataIndex: 'checking_task_num',
      width: 80,
      customHeaderCell: commonCustomHeaderCell(),
    },
    {
      title: '操作',
      dataIndex: 'action',
      width: 150,
      fixed: 'right',
    },
  ];
}

export function searchFormSchema(bizType: string): FormSchema[] {
  const currentMonth = dayjs();
  return [
    {
      field: 'provider',
      label: '供应商',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.PRODIVER),
        allowClear: true,
        defaultValue: 'mf',
      },
      colProps: { span: 4 },
      ifShow: () => bizType === 'specialLine',
    },
    {
      field: 'demand_id',
      label: '需求ID',
      component: 'Input',
      colProps: { span: 4 },
    },
    {
      field: 'area_name',
      label: '区域',
      component: 'Select',
      colProps: {
        span: 4,
      },
      componentProps: {
        placeholder: '请选择区域',
      },
    },
    {
      field: 'province_name',
      label: '省份',
      component: 'Select',
      colProps: {
        span: 4,
      },
      componentProps: {
        placeholder: '请选择省份',
      },
    },
    {
      field: 'isp_name',
      label: '运营商',
      component: 'Select',
      colProps: { span: 4 },
      componentProps: {
        options: [
          { label: '电信', value: '电信' },
          { label: '联通', value: '联通' },
          { label: '移动', value: '移动' },
        ],
      },
    },
    {
      field: 'dev_type_names',
      label: '设备类型',
      component: 'ApiSelect',
      componentProps: {
        options: [],
        mode: 'multiple',
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'is_province_scheduling',
      label: '调度控制',
      component: 'Select',
      componentProps: {
        options: Object.entries(IsProvinceScheduling).map(([value, label]) => ({
          value,
          label,
        })),
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'demand_time',
      label: '需求月份',
      component: 'MonthPicker',
      colProps: { span: 4 },
      componentProps: {
        allowClear: true,
        format: 'YYYY-MM',
        showTime: { format: 'YYYY-MM' },
        defaultValue: currentMonth,
      },
    },
    {
      field: 'end_time',
      label: '需求结束时间',
      component: 'DatePicker',
      colProps: { span: 4 },
      componentProps: {
        allowClear: true,
        format: 'YYYY-MM-DD',
        showTime: { format: 'YYYY-MM-DD' },
      },
    },
    {
      field: 'is_cover_diff_isp',
      label: '是否异网覆盖',
      component: 'Switch',
      componentProps: ({ formActionType }) => ({
        checkedValue: 1,
        unCheckedValue: 0,
        checkedChildren: '是',
        unCheckedChildren: '否',
        onChange: () => {
          formActionType.submit();
        },
      }),
      colProps: { span: 4 },
    },
    {
      field: 'is_independent_deploy',
      label: '是否独立部署',
      component: 'Switch',
      colProps: { span: 4 },
      componentProps: ({ formActionType }) => ({
        checkedValue: 1,
        unCheckedValue: 0,
        checkedChildren: '是',
        unCheckedChildren: '否',
        onChange: () => {
          formActionType.submit();
        },
      }),
    },
    {
      field: 'only_can_submit',
      label: '仅显示余量需求',
      component: 'Switch',
      colProps: { span: 4 },
      defaultValue: 1,
      componentProps: ({ formActionType }) => ({
        checkedValue: 1,
        unCheckedValue: 0,
        checkedChildren: '是',
        unCheckedChildren: '否',
        onChange: () => {
          formActionType.submit();
        },
      }),
    },
  ];
}

// 需求占用
export function DemandOccupySchema(is_cover_diff_isp: number, isp_name: string): FormSchema[] {
  return [
    {
      field: 'submit_bw',
      label: '提交带宽',
      component: 'InputNumber',
      colProps: { span: 11 },
      componentProps: {
        placeholder: '请输入提交带宽',
      },
    },
    {
      field: 'mac_count',
      label: 'MAC地址数量',
      component: 'InputNumber',
      colProps: { span: 11 },
      componentProps: {
        placeholder: '请输入MAC地址数量',
      },
    },
    {
      field: 'cover_diff_isp',
      label: '异网覆盖供应商',
      component: 'RadioGroup',
      colProps: { span: 22 },
      componentProps: {
        options: [
          { label: '电信', value: '电信', disabled: isp_name === '电信' },
          { label: '移动', value: '移动', disabled: isp_name === '移动' },
          { label: '联通', value: '联通', disabled: isp_name === '联通' },
        ],
      },
      ifShow: () => is_cover_diff_isp === 1,
      required: () => is_cover_diff_isp === 1,
    },
    {
      field: 'hostname',
      label: '主机名',
      component: 'InputTextArea',
      componentProps: {
        placeholder: '请输入需要占用的主机名，多个用换行分隔',
        rows: 5,
        maxlength: 1000,
        allowClear: true,
      },
      // required: true,
      colProps: { span: 22 },
    },
  ];
}

// 需求占用历史列表
export const DemandOccupyListColumns: BasicColumn[] = [
  {
    title: '需求ID',
    dataIndex: 'demand_id',
    width: 180,
    resizable: true,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 150,
    resizable: true,
  },
  {
    title: '设备类型',
    dataIndex: 'dev_name',
    width: 100,
  },
  {
    title: 'MAC地址',
    dataIndex: 'mac_addr',
    width: 120,
  },
  {
    title: '占用带宽(Mbps)',
    dataIndex: 'upload_bw',
    width: 100,
  },
  {
    title: '区域',
    dataIndex: 'area',
    width: 60,
  },
  {
    title: '省份',
    dataIndex: 'province',
    width: 80,
  },
  {
    title: '城市',
    dataIndex: 'city',
    width: 80,
  },
  {
    title: '运营商',
    dataIndex: 'isp',
    width: 60,
  },
  {
    title: '调度控制',
    dataIndex: 'is_province_scheduling',
    width: 60,
    customRender: ({ text }) => {
      if (text === 1) {
        return h('span', { style: { color: 'green' } }, '仅本省');
      } else if (text === 0) {
        return h('span', { style: { color: 'red' } }, '不限制');
      } else {
        return '未知';
      }
    },
  },
  {
    title: '是否异网覆盖',
    dataIndex: 'is_cover_diff_isp',
    width: 60,
    customRender: ({ text }) => {
      if (text === 1) {
        return h('span', { style: { color: 'green' } }, '是');
      } else if (text === 0) {
        return h('span', { style: { color: 'red' } }, '否');
      } else {
        return '未知';
      }
    },
  },
  {
    title: '任务ID',
    dataIndex: 'task_id',
    width: 60,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 60,
  },
  {
    title: '结果',
    dataIndex: 'result',
    width: 80,
  },
  {
    title: '操作时间',
    dataIndex: 'created_at',
    width: 120,
  },
];

export function searchOccupyListSchema(bizType: string): FormSchema[] {
  return [
    {
      field: 'demand_id',
      label: '需求ID',
      component: 'InputTextArea',
      colProps: { span: 4 },
      componentProps: {
        placeholder: '多个用换行分割',
      },
    },
    {
      field: 'hostname',
      label: '主机名',
      component: 'InputTextArea',
      colProps: { span: 4 },
      componentProps: {
        placeholder: '多个用换行分割',
      },
    },
    {
      field: 'mac',
      label: 'MAC地址',
      component: 'InputTextArea',
      colProps: { span: 4 },
      componentProps: {
        placeholder: '多个用换行分割',
      },
    },
    {
      field: 'isp',
      label: '运营商',
      component: 'Select',
      colProps: { span: 4 },
      componentProps: {
        options: [
          { label: '电信', value: '电信' },
          { label: '联通', value: '联通' },
          { label: '移动', value: '移动' },
        ],
      },
    },
    {
      field: 'dev_name',
      label: '设备类型',
      component: 'ApiSelect',
      componentProps: {
        options: [],
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      field: 'provider',
      label: '供应商',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.PRODIVER),
        allowClear: true,
        defaultValue: 'mf',
      },
      colProps: { span: 4 },
      ifShow: () => bizType === 'specialLine',
    },
    {
      field: 'area',
      label: '区域',
      component: 'ApiSelect',
      colProps: {
        span: 4,
      },
      componentProps: {
        options: [],
        allowClear: true,
      },
    },
    {
      field: 'province',
      label: '省份',
      component: 'ApiSelect',
      colProps: {
        span: 4,
      },
      componentProps: {
        options: [],
        allowClear: true,
      },
    },
    {
      field: 'status',
      label: '占用状态',
      component: 'Select',
      componentProps: {
        options: getSelectOptionsFromDict(KEnum.OCCUPY_STATUS),
        allowClear: true,
      },
      colProps: { span: 4 },
    },
    {
      label: '占用时间',
      field: 'occupy_time_range',
      component: 'RangePicker',
      colProps: { span: 8 },
      componentProps: {
        format: 'YYYY-MM-DD HH:mm:ss',
        showTime: true,
      },
    },
  ];
}

const DetectStatusMap: Record<number, { label: string; color: string }> = {
  0: { label: '探测中', color: '#1890ff' },
  1: { label: '通过', color: '#52c41a' },
  2: { label: '不通过', color: '#ff4d4f' },
};

// 需求探测结果列表
export const DemandDetectResultColumns: BasicColumn[] = [
  {
    title: '节点',
    dataIndex: 'owner',
    width: 120,
  },
  {
    title: '运营商',
    dataIndex: 'carrier',
    width: 100,
  },
  {
    title: '省份',
    dataIndex: 'province',
    width: 100,
  },
  {
    title: '探测ID',
    dataIndex: 'detect_id',
    width: 120,
  },
  {
    title: '探测设备',
    dataIndex: 'hostname',
    width: 150,
  },
  {
    title: '探测结果',
    dataIndex: 'detect_status',
    width: 100,
    customRender: ({ text }) => {
      const status = DetectStatusMap[Number(text)];
      if (!status) return '-';
      return h('span', { style: { color: status.color } }, status.label);
    },
  },
  {
    title: '探测时间',
    dataIndex: 'detect_time',
    width: 180,
  },
];

// 需求探测结果搜索表单
export function searchDetectResultSchema(): FormSchema[] {
  return [
    {
      field: 'owner',
      label: '节点',
      component: 'Input',
      colProps: { span: 6 },
      componentProps: {
        placeholder: '请输入节点',
        allowClear: true,
      },
    },
    {
      field: 'carrier',
      label: '运营商',
      component: 'Select',
      colProps: { span: 6 },
      componentProps: {
        options: [
          { label: '电信', value: '电信' },
          { label: '联通', value: '联通' },
          { label: '移动', value: '移动' },
        ],
        allowClear: true,
      },
    },
    {
      field: 'detect_status',
      label: '探测结果',
      component: 'Select',
      colProps: { span: 6 },
      componentProps: {
        placeholder: '请选择探测结果',
        options: Object.entries(DetectStatusMap).map(([value, item]) => ({
          value: Number(value),
          label: item.label,
        })),
        allowClear: true,
      },
    },
    {
      field: 'detect_date',
      label: '探测日期',
      component: 'DatePicker',
      defaultValue: dayjs().format('YYYY-MM-DD'),
      colProps: { span: 6 },
      componentProps: {
        format: 'YYYY-MM-DD',
        valueFormat: 'YYYY-MM-DD',
        allowClear: true,
        placeholder: '请选择探测日期',
      },
    },
  ];
}

export const OccupyTaskColumns: BasicColumn[] = [
  {
    title: '需求ID',
    dataIndex: 'demand_id',
    width: 180,
    resizable: true,
  },
  {
    title: '主机名',
    dataIndex: 'hostname',
    width: 120,
    resizable: true,
  },
  {
    title: '状态',
    dataIndex: 'status',
    width: 100,
    resizable: true,
    customRender: ({ text }) => {
      // 0: 待执行, 1: 执行中, 2: 成功, 3: 失败
      if (text === 0) {
        return h('span', { style: { color: 'blue' } }, '待执行');
      } else if (text === 1) {
        return h('span', { style: { color: 'orange' } }, '执行中');
      } else if (text === 2) {
        return h('span', { style: { color: 'green' } }, '成功');
      } else if (text === 3) {
        return h('span', { style: { color: 'red' } }, '失败');
      } else {
        return '-';
      }
    },
  },
  {
    title: '结果',
    dataIndex: 'result',
    resizable: true,
    width: 180,
    customRender: ({ text }) => {
      return text ? text : '-';
    },
  },
  {
    title: '操作时间',
    dataIndex: 'created_at',
    width: 100,
  },
  {
    title: '执行时间',
    dataIndex: 'execute_at',
    width: 100,
  },
];
