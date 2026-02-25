import type { BasicColumn, FormSchema } from '@/components/Table';
import { CommonEnum, RoomEnum } from '@/enums/dictTypeCode';
import { getDictDataMapFromDict, getSelectOptionsFromDict } from '@/utils/dict';
import { calTableCellRowSpan, commonCustomHeaderCell, customCeilDivide } from '@/utils/util';
import { getNumberResult } from '../room/peak/data';
import { GetEcdnRoomList } from '@/api/network/speed_limit_job';
import {
  Network_SpeedLimitJob_JobType,
  Network_SpeedLimitJob_LimitType,
  Network_SpeedLimitJob_OperateTargetType,
  Network_SpeedLimitJob_Status,
} from '@/enums/dictValueEnum';

export const speedLimitJobTypeOptions = getSelectOptionsFromDict(
  RoomEnum.NETWORK_SPEED_LIMIT_JOB_TYPE,
);
export const speedLimitJobTypeMap = getDictDataMapFromDict(RoomEnum.NETWORK_SPEED_LIMIT_JOB_TYPE);

export const speedLimitJobStatusOptions = getSelectOptionsFromDict(
  RoomEnum.NETWORK_SPEED_LIMIT_JOB_STATUS,
);
export const speedLimitJobStatusMap = getDictDataMapFromDict(
  RoomEnum.NETWORK_SPEED_LIMIT_JOB_STATUS,
);

export const speedLimitJobExecStatusOptions = getSelectOptionsFromDict(
  RoomEnum.NETWORK_SPEED_LIMIT_JOB_EXEC_STATUS,
);
export const speedLimitJobExecStatusMap = getDictDataMapFromDict(
  RoomEnum.NETWORK_SPEED_LIMIT_JOB_EXEC_STATUS,
);

export const speedLimitTargetTypeOptions = getSelectOptionsFromDict(
  RoomEnum.NETWORK_SPEED_LIMIT_TARGET_TYPE,
);
export const speedLimitTargetTypeMap = getDictDataMapFromDict(
  RoomEnum.NETWORK_SPEED_LIMIT_TARGET_TYPE,
);

export const speedLimitJobLimitTypeOptions = getSelectOptionsFromDict(
  RoomEnum.NETWORK_SPEED_LIMIT_LIMIT_TYPE,
);
export const speedLimitJobLimitTypeMap = getDictDataMapFromDict(
  RoomEnum.NETWORK_SPEED_LIMIT_LIMIT_TYPE,
);

export const ispOptions = getSelectOptionsFromDict(CommonEnum.ISP);
export const ispMap = getDictDataMapFromDict(CommonEnum.ISP);
export const roomChargeModeOptions = getSelectOptionsFromDict(RoomEnum.CHARGE_MODE);
export const roomChargeModeMap = getDictDataMapFromDict(RoomEnum.CHARGE_MODE);
export const roomTypeOptions = getSelectOptionsFromDict(RoomEnum.TYPE);
export const roomTypeMap = getDictDataMapFromDict(RoomEnum.TYPE);
export const operateTypeOptions = [
  { label: '开启限速', value: 1 },
  { label: '解除限速', value: 2 },
];
export const speedLimitDefaultValue = 15 * 1000; //默认限速值15Mbps

export const swithcPortValueValidator = function (getFormHandleFn: any) {
  return (_rule, _value) => {
    const { getFieldsValue } = getFormHandleFn();
    const values = getFieldsValue();
    const switchPort = values.switchPort || '';
    const switchPortRange = values.switchPortRange || '';
    if (
      !switchPortRange &&
      !switchPort &&
      values.operateType == Network_SpeedLimitJob_OperateTargetType.SWITCH_PORT
    ) {
      return Promise.reject('“交换机端口”和“交换机端口范围”不能同时为空');
    }
    return Promise.resolve();
  };
};

export const getSpeedLimitJobColumns = function (): BasicColumn[] {
  return [
    {
      title: '任务名称',
      dataIndex: 'name',
      width: 250,
      resizable: true,
    },
    {
      title: '群号|节点|机房',
      dataIndex: 'roomName',
      width: 160,
      resizable: true,
    },
    {
      title: '地区',
      dataIndex: 'location',
      width: 100,
      resizable: true,
    },
    {
      title: '运营商',
      dataIndex: 'isp',
      width: 80,
      resizable: true,
    },
    {
      title: '机房带宽\nGbps',
      dataIndex: 'bandwidth',
      width: 100,
      // sorter: true,
      resizable: true,
      customHeaderCell: commonCustomHeaderCell(),
      customRender: ({ record }) => {
        return getNumberResult(customCeilDivide(record.bandwidth, 1000 * 1000 * 1000, 2) as any);
      },
    },
    // {
    //   title: '交换机',
    //   dataIndex: 'switchDesc',
    //   width: 100,
    //   resizable: true,
    //   customRender: ({ record }) => {
    //     return record.strategies.limitTarget || '';
    //   },
    // },
    {
      title: '机房类型',
      dataIndex: 'roomType',
      width: 90,
    },
    {
      title: '计费方式',
      dataIndex: 'chargeMode',
      width: 90,
      resizable: true,
    },
    {
      title: '任务类型',
      dataIndex: 'jobType',
      width: 130,
      resizable: true,
    },
    {
      title: '状态',
      dataIndex: 'status',
      width: 90,
      resizable: true,
    },
    {
      title: '最后执行时间',
      dataIndex: 'lastExecuteTime',
      width: 160,
      resizable: true,
    },
    {
      title: '最后执行状态',
      dataIndex: 'lastExecuteStatus',
      width: 100,
      resizable: true,
    },
    {
      title: '更新人',
      dataIndex: 'updateUser',
      width: 100,
      resizable: true,
      customRender: ({ record }) => {
        return record?.updateUser?.nickName || '';
      },
    },
    {
      title: '更新时间',
      dataIndex: 'userUpdatedAt',
      width: 160,
      resizable: true,
    },
    // {
    //   title: '备注',
    //   dataIndex: 'remark',
    // },
  ];
};

export const getSpeedLimitJobSearchFormSchema = function (): FormSchema[] {
  return [
    {
      field: 'name',
      label: '任务名称',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'roomName',
      label: '群号|节点|机房',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'location',
      label: '地区',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'isps',
      label: '运营商',
      component: 'Select',
      colProps: { span: 6 },
      componentProps: {
        options: ispOptions,
        mode: 'multiple',
        allowClear: true,
      },
    },
    {
      field: 'roomTypes',
      label: '机房类型',
      component: 'Select',
      colProps: { span: 6 },
      componentProps: {
        options: roomTypeOptions,
        mode: 'multiple',
      },
    },
    {
      field: 'chargeModes',
      label: '计费方式',
      component: 'Select',
      colProps: { span: 6 },
      componentProps: {
        options: roomChargeModeOptions,
        mode: 'multiple',
      },
    },
    {
      field: 'jobTypes',
      label: '任务类型',
      component: 'Select',
      colProps: { span: 6 },
      componentProps: {
        options: speedLimitJobTypeOptions,
        mode: 'multiple',
      },
    },
    {
      field: 'status',
      label: '状态',
      component: 'Select',
      colProps: { span: 6 },
      componentProps: {
        options: speedLimitJobStatusOptions,
      },
    },
    {
      field: 'lastExecuteStatus',
      label: '最后执行状态',
      component: 'Select',
      colProps: { span: 6 },
      componentProps: {
        options: speedLimitJobExecStatusOptions,
        mode: 'multiple',
      },
    },
  ];
};
export const getSwitchPortValueValidator = function (getFormHandleFn: any) {
  return [{ validator: swithcPortValueValidator(getFormHandleFn), trigger: 'change' }];
};

export const getSpeedLimitJobFormSchema = function (
  onJobTypeChange: any,
  onRoomChange: any,
  onSwitchesChange: any,
  onOperateTypeChange: any,
  getFormHandleFn: any,
): FormSchema[] {
  return [
    {
      field: 'name',
      label: '任务名称',
      component: 'Input',
      colProps: { span: 24 },
      required: true,
    },
    {
      field: 'status',
      label: '状态',
      component: 'Switch',
      colProps: { span: 6 },
      componentProps: {
        checkedChildren: '启用',
        unCheckedChildren: '禁用',
        checkedValue: Network_SpeedLimitJob_Status.ENABLE,
        unCheckedValue: Network_SpeedLimitJob_Status.DISABLE,
      },
      defaultValue: Network_SpeedLimitJob_Status.ENABLE,
      required: true,
    },
    {
      field: 'retryCount',
      label: '重试次数',
      component: 'InputNumber',
      colProps: { span: 9 },
      componentProps: {
        min: 0,
        max: 10,
        precision: 0,
        step: 1,
      },
      defaultValue: 0,
      required: true,
      helpMessage: '当任务执行失败时，将进行重试，0为不重试，最大重试10次',
    },
    {
      field: 'maxExecuteDelayMinutes',
      label: '最大延迟时间',
      component: 'InputNumber',
      colProps: { span: 9 },
      componentProps: {
        min: 1,
        max: 45,
        precision: 0,
        step: 1,
        addonAfter: '分钟',
      },
      defaultValue: 5,
      required: true,
      helpMessage:
        '当任务因为重试或者其他未知原因导致执行延迟时，最大延迟的时间，如果任务超过该时间，将不再重试',
    },
    {
      field: 'ecdnRoomId',
      label: '群号|节点|机房',
      component: 'ApiSelect',
      componentProps: {
        api: async () => {
          const { items } = await GetEcdnRoomList({
            pageIndex: 1,
            pageSize: 99999,
            roomStatus: 1,
            roomType: [1, 2],
            origin: 1,
          });
          if (!items || items.length === 0) {
            return [];
          }
          return items.map((item) => ({
            label: `${item.name}-${item.location}-${item.isp}-${getNumberResult(customCeilDivide(item.bandwidth, 1000 * 1000 * 1000, 2) as any)}Gbps`,
            value: item.id,
          }));
        },
        filterOption: (input, option: any) => {
          return option.label.toLowerCase().indexOf(input.toLowerCase()) >= 0;
        },
        showSearch: true,
        allowClear: false,
        onChange: onRoomChange,
      },
      required: true,
      colProps: { span: 24 },
    },
    {
      field: 'switchId',
      label: '交换机',
      component: 'Select',
      componentProps: {
        // mode: 'multiple',
        options: [],
        showSearch: true,
        allowClear: false,
        // maxTagCount: 3,
        filterOption: (inputValue, option) => {
          return option?.label.toUpperCase().indexOf(inputValue.toUpperCase()) !== -1;
        },
        onChange: onSwitchesChange,
        loading: false,
      },
      required: true,
      colProps: { span: 24 },
    },
    {
      field: 'operateType',
      label: '操作对象',
      component: 'RadioButtonGroup',
      colProps: { span: 24 },
      componentProps: {
        options: speedLimitTargetTypeOptions as any,
        buttonStyle: 'solid',
        optionType: 'button',
        onChange: onOperateTypeChange,
        size: 'default',
      },
      required: true,
      defaultValue: String(Network_SpeedLimitJob_OperateTargetType.SWITCH_PORT),
    },
    {
      field: 'switchPort',
      label: '交换机端口',
      component: 'InputTextArea',
      componentProps: {
        placeholder: '请输入交换机端口，使用空格或者换行分隔，如：10GE1/0/2 10GE1/0/3 10GE1/0/4',
        rows: 3,
        onChange: async function () {
          try {
            await getFormHandleFn().validateFields(['switchPortRange', 'switchPort']);
          } catch (e) {
            // 忽略 outOfDate 错误
          }
        },
        // 当输入框失去焦点时，将输入框中的内容按行分割，并去除空行
        onBlur: function () {
          const { switchPort } = getFormHandleFn().getFieldsValue();
          if (!switchPort || switchPort.trim() === '') {
            return;
          }
          const data: string[] = [];
          switchPort.split('\n').forEach((item: string) => {
            if (item.trim()) {
              data.push(item.trim());
            }
          });
          // 去重
          const uniqueData = [...new Set(data)];
          // uniqueData.sort((a, b) => (a < b ? -1 : 1));
          getFormHandleFn().setFieldsValue({
            switchPort: uniqueData.join(' '),
          });
        },
      },
      rules: getSwitchPortValueValidator(getFormHandleFn) as any,
      colProps: { span: 24 },
      required: false,
    },
    {
      field: 'switchPortRange',
      label: '交换机端口范围',
      component: 'Input',
      colProps: { span: 24 },
      required: false,
      componentProps: {
        placeholder: '请输入端口范围，如：10GE1/0/2 to 10GE1/0/11',
        onChange: async function () {
          try {
            await getFormHandleFn().validateFields(['switchPortRange', 'switchPort']);
          } catch (e) {
            // 忽略 outOfDate 错误
          }
        },
      },
      rules: getSwitchPortValueValidator(getFormHandleFn) as any,
    },
    {
      field: 'businessTag',
      label: '业务标签',
      component: 'Select',
      componentProps: {
        options: [],
        mode: 'multiple',
        allowClear: true,
      },
      colProps: { span: 24 },
      required: true,
    },
    {
      field: 'jobType',
      label: '任务类型',
      component: 'RadioButtonGroup',
      colProps: { span: 24 },
      componentProps: {
        options: speedLimitJobTypeOptions as any,
        buttonStyle: 'solid',
        optionType: 'button',
        onChange: onJobTypeChange,
        size: 'default',
      },
      required: true,
      defaultValue: String(Network_SpeedLimitJob_JobType.IMMEDIATE_TASK),
    },
    {
      field: 'timeSettingGroup',
      label: '执行配置',
      slot: 'timeSettingGroup',
      colProps: { span: 24 },
    },
  ];
};

export const getSpeedLimitExecRecordSearchForm = function (onTimePikerOpen): FormSchema[] {
  return [
    {
      field: '[executeTimeBegin, executeTimeEnd]',
      label: '执行时间',
      component: 'RangePicker',
      componentProps: {
        allowClear: false,
        format: 'YYYY-MM-DD HH:mm:ss',
        showTime: { format: 'HH:mm:ss' },
        placeholder: ['开始时间', '结束时间'],
        style: {
          width: '100%',
        },
        onOpenChange: onTimePikerOpen,
      },
      colProps: { span: 8 },
      required: true,
    },
    {
      field: 'switchDesc',
      label: '交换机名称',
      component: 'Input',
      colProps: { span: 8 },
    },
    {
      label: '操作类型',
      field: 'limitType',
      component: 'Select',
      colProps: { span: 8 },
      componentProps: {
        options: speedLimitJobLimitTypeOptions,
      },
    },
    {
      label: '执行状态',
      field: 'executeStatus',
      component: 'Select',
      colProps: { span: 8 },
      componentProps: {
        options: speedLimitJobExecStatusOptions,
      },
    },
  ];
};
export function getSpeedLimitExecRecordColumns(safeGetDataSource: any): BasicColumn[] {
  return [
    {
      title: '任务时间',
      dataIndex: 'jobTime',
      width: 145,
      resizable: true,
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        // 计算单元格合并的行数
        const rowSpan = calTableCellRowSpan(dataSource, index, (current, next) => {
          return current.jobTime === next.jobTime;
        });
        return {
          rowSpan: rowSpan,
        };
      },
    },
    {
      title: '交换机名称',
      dataIndex: 'switchDesc',
      width: 300,
      resizable: true,
    },
    {
      title: '端口名称',
      dataIndex: 'portRange',
      width: 380,
      resizable: true,
    },
    {
      title: '操作类型',
      dataIndex: 'limitType',
      width: 90,
    },
    {
      title: '限速值(Mbps)',
      dataIndex: 'limitValueMbps',
      width: 110,
      customRender: ({ record }) => {
        if (record.limitType == Network_SpeedLimitJob_LimitType.LIMIT) {
          return record.limitValueMbps;
        } else {
          return '-';
        }
      },
    },
    {
      title: '执行状态',
      dataIndex: 'executeStatus',
      width: 90,
    },
    {
      title: '执行时间',
      dataIndex: 'executeTime',
      width: 145,
    },
    {
      title: '终端日志',
      dataIndex: 'consoleLog',
      width: 80,
      resizable: true,
    },
    {
      title: '更多信息',
      dataIndex: 'message',
      width: 160,
    },
  ];
}
