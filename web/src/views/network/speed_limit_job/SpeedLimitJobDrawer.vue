<template>
  <BasicDrawer v-bind="$attrs" @register="registerModal" @close="onDrawerCancel" @ok="onDrawerOk">
    <BasicForm @register="registerForm">
      <template #timeSettingGroup>
        <div class="flex flex-col gap-2 p-3 border border-gray-200 rounded bg-gray-50/40">
          <!-- 立即执行 -->
          <div v-if="data.jobType == Network_SpeedLimitJob_JobType.IMMEDIATE_TASK">
            <!-- 标题 -->
            <div class="flex items-center justify-between mb-2">
              <span class="font-medium">立即执行</span>
            </div>
            <!-- 操作类型 -->
            <div class="flex gap-2 items-center mb-2">
              <span style="width: 80px"><span class="text-red-500">*</span>操作类型：</span>
              <FormItemRest>
                <Select
                  v-model:value="data.immediateTask.type"
                  placeholder="请选择操作类型"
                  style="flex: 1"
                  :allow-clear="false"
                  required
                  :options="operateTypeOptions"
                />
              </FormItemRest>
            </div>
            <!-- 限速值 -->
            <div
              class="flex gap-2 items-center mb-2"
              v-if="data.immediateTask.type == Network_SpeedLimitJob_LimitType.LIMIT"
            >
              <span style="width: 80px"><span class="text-red-500">*</span>限速值：</span>
              <FormItemRest>
                <InputNumber
                  v-model:value="data.immediateTask.limitValue"
                  placeholder="请输入限速值"
                  :min="1"
                  :max="999999999"
                  :precision="0"
                  :step="1"
                  required
                  style="flex: 1"
                  addonAfter="Kbps"
                />
              </FormItemRest>
            </div>
          </div>
          <!-- 定时任务，一次性 -->
          <div v-else-if="data.jobType == Network_SpeedLimitJob_JobType.ONCE_TASK">
            <!-- 标题 -->
            <div class="flex items-center justify-between mb-2">
              <span class="font-medium">定时执行（一次性）</span>
            </div>
            <!-- 操作类型 -->
            <div class="flex gap-2 items-center mb-2">
              <span style="width: 80px"><span class="text-red-500">*</span>操作时间：</span>
              <FormItemRest>
                <DatePicker
                  v-model:value="data.onceTask.executeTime"
                  style="flex: 1"
                  show-time
                  :allow-clear="false"
                  required
                  format="YYYY-MM-DD HH:mm"
                  placeholder="请选择操作时间"
                />
              </FormItemRest>
            </div>
            <div class="flex gap-2 items-center mb-2">
              <span style="width: 80px"><span class="text-red-500">*</span>操作类型：</span>
              <FormItemRest>
                <Select
                  v-model:value="data.onceTask.type"
                  placeholder="请选择操作类型"
                  style="flex: 1"
                  required
                  :allow-clear="false"
                  :options="operateTypeOptions"
                />
              </FormItemRest>
            </div>
            <!-- 限速值 -->
            <div
              class="flex gap-2 items-center mb-2"
              v-if="data.onceTask.type == Network_SpeedLimitJob_LimitType.LIMIT"
            >
              <span style="width: 80px"><span class="text-red-500">*</span>限速值：</span>
              <FormItemRest>
                <InputNumber
                  v-model:value="data.onceTask.limitValue"
                  placeholder="请输入限速值"
                  :min="1"
                  :max="999999999"
                  :precision="0"
                  :step="1"
                  required
                  style="flex: 1"
                  addonAfter="Kbps"
                />
              </FormItemRest>
            </div>
          </div>
          <!-- 周期任务 -->
          <div v-else-if="data.jobType == Network_SpeedLimitJob_JobType.CYCLE_TASK">
            <div class="flex items-center justify-between mb-2">
              <span class="font-medium">周期任务</span>
              <Button
                type="text"
                size="small"
                class="text-blue-500 hover:text-blue-600 hover:bg-blue-50"
                @click="importPeakSchema"
              >
                <ImportOutlined />
                导入打峰方案
              </Button>
            </div>
            <!-- 循环执行周期 -->
            <div
              v-for="(_cycleTask, index) in data.cycleTask"
              :key="index"
              class="border border-gray-300 border-l-4 border-l-blue-400 rounded bg-white mt-2 mb-3 p-3 shadow-sm transition-all duration-200 cursor-pointer hover:shadow-lg hover:border-blue-400 hover:bg-blue-50/50 hover:border-l-blue-500"
            >
              <!-- 执行周期标题 -->
              <div class="flex items-center justify-between mb-2">
                <span class="font-medium text-blue-600">执行周期{{ index + 1 }}</span>
                <Button
                  type="text"
                  size="small"
                  class="text-blue-500 hover:text-blue-600 hover:bg-blue-50"
                  @click="removeCycleTaskItem(index)"
                  :disabled="data.cycleTask.length === 1"
                >
                  <DeleteOutlined />
                  删除执行周期
                </Button>
              </div>
              <!-- 执行周期选择 -->
              <div class="flex gap-2 items-center mb-2">
                <span style="width: 80px"><span class="text-red-500">*</span>日期：</span>
                <FormItemRest>
                  <RangePicker
                    v-model:value="data.cycleTask[index].dateRange"
                    style="flex: 1"
                    format="YYYY-MM-DD"
                    :placeholder="['开始日期', '结束日期']"
                  />
                </FormItemRest>
              </div>
              <!-- 操作列表 -->
              <div
                v-for="(_value, i) in data.cycleTask[index].timeRange"
                :key="i"
                class="border border-gray-300 border-l-4 border-l-amber-400 rounded flex flex-col gap-2 p-3 mt-2 mb-2 bg-white shadow-sm transition-all duration-200 cursor-pointer hover:shadow-lg hover:border-amber-400 hover:bg-amber-50/50 hover:border-l-amber-500"
              >
                <div class="flex items-center justify-between mb-2">
                  <span class="font-medium text-amber-600">操作{{ i + 1 }}</span>
                  <Button
                    type="text"
                    size="small"
                    class="text-amber-500 hover:text-amber-600 hover:bg-amber-50"
                    @click="removeCycleTaskTimeRangeItem(index, i)"
                    :disabled="data.cycleTask[index].timeRange.length == 1"
                  >
                    <DeleteOutlined />
                    删除操作
                  </Button>
                </div>
                <div class="flex gap-2 items-center mb-2">
                  <span style="width: 80px"><span class="text-red-500">*</span>操作时间：</span>
                  <FormItemRest>
                    <TimePicker
                      v-model:value="data.cycleTask[index].timeRange[i].time"
                      style="flex: 1"
                      format="HH:mm"
                      :allow-clear="false"
                      required
                      placeholder="请选择操作时间"
                    />
                  </FormItemRest>
                </div>
                <!-- 操作类型 -->
                <div class="flex gap-2 items-center mb-2">
                  <span style="width: 80px"><span class="text-red-500">*</span>操作类型：</span>
                  <FormItemRest>
                    <Select
                      v-model:value="data.cycleTask[index].timeRange[i].type"
                      placeholder="请选择操作类型"
                      style="flex: 1"
                      :allow-clear="false"
                      required
                      :options="operateTypeOptions"
                    />
                  </FormItemRest>
                </div>
                <!-- 限速值 -->
                <div
                  class="flex gap-2 items-center mb-2"
                  v-if="
                    data.cycleTask[index].timeRange[i].type == Network_SpeedLimitJob_LimitType.LIMIT
                  "
                >
                  <span style="width: 80px"><span class="text-red-500">*</span>限速值：</span>
                  <FormItemRest>
                    <InputNumber
                      v-model:value="data.cycleTask[index].timeRange[i].limitValue"
                      placeholder="请输入限速值"
                      :min="1"
                      :max="999999999"
                      :precision="0"
                      :step="1"
                      style="flex: 1"
                      addonAfter="Kbps"
                      required
                    />
                  </FormItemRest>
                </div>
              </div>
              <Button type="dashed" class="mt-2" block @click="addCycleTaskTimeRangeItem(index)">
                <PlusOutlined />
                添加操作
              </Button>
            </div>
            <Button type="dashed" class="mt-2" block @click="addCycleTaskItem">
              <PlusOutlined />
              添加执行周期
            </Button>
          </div>
        </div>
      </template>
    </BasicForm>
  </BasicDrawer>
</template>
<script setup lang="ts">
  import { BasicDrawer, useDrawerInner } from '@/components/Drawer';
  import { createVNode, nextTick, reactive } from 'vue';
  import { BasicForm, useForm } from '@/components/Form';
  import {
    getSpeedLimitJobFormSchema,
    getSwitchPortValueValidator,
    operateTypeOptions,
    speedLimitDefaultValue,
  } from './data';
  import {
    CreateSpeedLimitJob,
    GetRoomSwitches,
    GetSwitchBusinessTags,
    ImportPeakSchema,
    UpdateSpeedLimitJob,
  } from '@/api/network/speed_limit_job';
  import {
    Network_SpeedLimitJob_ExecuteNow,
    Network_SpeedLimitJob_JobType,
    Network_SpeedLimitJob_LimitType,
    Network_SpeedLimitJob_OperateTargetType,
    Network_SpeedLimitJob_Status,
  } from '@/enums/dictValueEnum';
  import {
    InputNumber,
    Select,
    Form,
    DatePicker,
    Button,
    RangePicker,
    TimePicker,
    Modal,
    message,
  } from 'ant-design-vue';
  import {
    PlusOutlined,
    DeleteOutlined,
    ImportOutlined,
    ExclamationCircleOutlined,
  } from '@ant-design/icons-vue';
  import dayjs from 'dayjs';

  const FormItemRest = Form.ItemRest;
  defineOptions({ name: 'SpeedLimitJobDrawer' });
  const emit = defineEmits(['register', 'success']);
  interface ImmediateTask {
    type: number | undefined;
    limitValue: number | undefined;
  }
  interface OnceTaskType {
    executeTime: string | undefined;
    type: number | undefined;
    limitValue: number | undefined;
  }
  interface CycleTaskType {
    dateRange: [any, any];
    timeRange: TimeRange[];
  }
  interface TimeRange {
    time: string | undefined;
    type: number | undefined;
    limitValue: number | undefined;
  }

  const data = reactive({
    isRewriting: false, // 新增标志位
    record: null as any,
    title: '',
    isUpdate: false,
    jobType: null,
    immediateTask: {
      type: undefined,
      limitValue: speedLimitDefaultValue,
    } as ImmediateTask,
    onceTask: {
      executeTime: undefined,
      type: undefined,
      limitValue: speedLimitDefaultValue,
    } as OnceTaskType,
    cycleTask: [
      {
        dateRange: [undefined, undefined],
        timeRange: [
          {
            time: undefined,
            type: undefined,
            limitValue: speedLimitDefaultValue,
          },
        ] as TimeRange[],
      },
    ] as CycleTaskType[],
  });
  const [
    registerForm,
    { setFieldsValue, getFieldsValue, updateSchema, validateFields, clearValidate },
  ] = useForm({
    labelWidth: 130,
    baseColProps: { span: 24 },
    schemas: getSpeedLimitJobFormSchema(
      onJobTypeChange,
      onRoomChange,
      onSwitchesChange,
      onOperateTypeChange,
      getFormHandleFn,
    ),
    showActionButtonGroup: false,
  });
  const [registerModal, { setDrawerProps, closeDrawer }] = useDrawerInner(async (d) => {
    data.record = d.record;
    data.isUpdate = !!d.isUpdate;
    initData();
    data.title = data.isUpdate ? '编辑定时任务' : '新增定时任务';
    setDrawerProps({
      title: data.title,
      width: 800,
      // height: 630,
      destroyOnClose: true,
      showCancelBtn: true,
      showOkBtn: true,
      showFooter: true,
    });
    // 新增
    if (data.isUpdate) {
      rewriteFormData();
      // 更新
    } else {
      // 设置默认值
      setFieldsValue({
        jobType: String(Network_SpeedLimitJob_JobType.IMMEDIATE_TASK),
        operateType: String(Network_SpeedLimitJob_OperateTargetType.SWITCH_PORT),
      });
    }
  });
  async function rewriteFormData() {
    if (!data.record) {
      return;
    }
    data.isRewriting = true;
    try {
      await setFieldsValue({
        name: data.record.name,
        status: data.record.status,
        ecdnRoomId: data.record.ecdnRoomId,
        switchId: data.record?.strategies?.limitTarget?.switchIds?.[0] || null,
        operateType: String(data.record?.strategies?.limitTarget?.operateType),
        switchPort: data.record?.strategies?.limitTarget?.switchPort || undefined,
        switchPortRange: data.record?.strategies?.limitTarget?.switchPortRange || undefined,
        businessTag: data.record?.strategies?.limitTarget?.businessTag || undefined,
        jobType: String(data.record?.jobType),
        retryCount: Number(data.record?.retryCount) || 0,
        maxExecuteDelayMinutes: Number(data.record?.maxExecuteDelayMinutes) || null,
      });
      await updateSchema({
        field: 'ecdnRoomId',
        componentProps: {
          readonly: true,
          disabled: true,
        },
      });
      // 立即执行-写表单数据
      if (data.record.jobType == Network_SpeedLimitJob_JobType.IMMEDIATE_TASK) {
        data.immediateTask = {
          type: data.record?.strategies?.immediateTaskStrategy?.limitType?.type,
          limitValue:
            data.record?.strategies?.immediateTaskStrategy?.limitType?.limitValue ||
            speedLimitDefaultValue,
        };
      } else if (data.record.jobType == Network_SpeedLimitJob_JobType.ONCE_TASK) {
        let executeTime = undefined as any;
        if (data.record?.strategies?.onceTaskStrategy?.executeTime) {
          executeTime = dayjs(
            data.record?.strategies?.onceTaskStrategy?.executeTime,
            'YYYY-MM-DD HH:mm:ss',
          );
        }
        // 定时任务(一次性)-写表单数据
        data.onceTask = {
          executeTime: executeTime,
          type: data.record?.strategies?.onceTaskStrategy?.limitType?.type,
          limitValue:
            data.record?.strategies?.onceTaskStrategy?.limitType?.limitValue ||
            speedLimitDefaultValue,
        };
      } else if (data.record.jobType == Network_SpeedLimitJob_JobType.CYCLE_TASK) {
        // 周期任务-写表单数据
        data.cycleTask =
          data.record?.strategies?.cycleTaskStrategy?.dates?.map((item: any) => {
            // 开始日期
            let startDate = item.startDate ? dayjs(item.startDate, 'YYYY-MM-DD') : undefined;
            // 结束日期
            let endDate = item.endDate ? dayjs(item.endDate, 'YYYY-MM-DD') : undefined;
            // 操作时间
            let timeRange = item?.times?.map((t: any) => {
              let time = t.time ? dayjs(t.time, 'HH:mm:ss') : undefined;
              return {
                time: time,
                type: t?.limitType?.type,
                limitValue: t?.limitType?.limitValue || speedLimitDefaultValue,
              };
            });
            return {
              dateRange: [startDate, endDate],
              timeRange: timeRange,
            };
          }) || [];
      }
    } finally {
      nextTick(() => {
        data.isRewriting = false;
      });
    }
  }
  function getFormHandleFn() {
    return {
      getFieldsValue,
      updateSchema,
      setFieldsValue,
      validateFields,
    };
  }
  function onDrawerCancel() {
    closeDrawer();
    emit('success');
  }
  function onJobTypeChange(value: any) {
    data.jobType = value;
  }
  function removeCycleTaskItem(index: number) {
    if (data.cycleTask.length > 1) {
      data.cycleTask.splice(index, 1);
    }
  }
  function addCycleTaskItem() {
    data.cycleTask.push({
      dateRange: [undefined, undefined],
      timeRange: [
        {
          time: undefined,
          type: undefined,
          limitValue: speedLimitDefaultValue,
        },
      ],
    });
  }
  function initData() {
    data.jobType = null;
    data.immediateTask = {
      type: undefined,
      limitValue: speedLimitDefaultValue,
    };
    data.onceTask = {
      executeTime: undefined,
      type: undefined,
      limitValue: speedLimitDefaultValue,
    };
    data.cycleTask = [
      {
        dateRange: [undefined, undefined],
        timeRange: [
          {
            time: undefined,
            type: undefined,
            limitValue: speedLimitDefaultValue,
          },
        ],
      },
    ];
  }
  function addCycleTaskTimeRangeItem(index: number) {
    data.cycleTask[index].timeRange.push({
      time: undefined,
      type: undefined,
      limitValue: speedLimitDefaultValue,
    });
  }
  function removeCycleTaskTimeRangeItem(index: number, i: number) {
    if (data.cycleTask[index].timeRange.length > 1) {
      data.cycleTask[index].timeRange.splice(i, 1);
    }
  }
  async function importPeakSchema() {
    let { ecdnRoomId } = await getFieldsValue();
    if (!ecdnRoomId) {
      message.error('请选择 群号|节点|机房 信息');
      return;
    }
    Modal.confirm({
      title: '是否确认从打峰方案中导入执行配置?',
      icon: createVNode(ExclamationCircleOutlined),
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        let res = await ImportPeakSchema({ ecdnRoomId: ecdnRoomId });
        if (!res || !res.schemas || res.schemas.length === 0) {
          message.warn('节点未配置打峰方案');
          return;
        }
        // 写入打峰方案到数据中
        let newCycleTask: CycleTaskType[] = [];
        res.schemas.forEach((item: any) => {
          newCycleTask.push({
            dateRange: [
              (item.startDate ? dayjs(item.startDate) : undefined) as any,
              (item.endDate ? dayjs(item.endDate) : undefined) as any,
            ],
            timeRange:
              item.timeRange?.map((t: any) => ({
                time: t.time ? dayjs(t.time, 'HH:mm:00') : undefined,
                type: t.type,
                limitValue: t.limitValue || speedLimitDefaultValue,
              })) || [],
          });
        });
        data.cycleTask = newCycleTask;
        message.success('导入成功');
      },
    });
  }
  async function onRoomChange(_value: any) {
    // 编辑回写数据，不清空设置的值
    if (!data.isRewriting) {
      await setFieldsValue({
        switchId: null,
      });
    }
    // 重新获取交换机信息
    let switchOptions: any[] = [];
    let values = await getFieldsValue();
    if (values.ecdnRoomId) {
      let { switches } = await GetRoomSwitches({ ecdnRoomId: values.ecdnRoomId, status: 3 });
      if (!switches || switches.length === 0) {
        switchOptions = [];
      } else {
        switchOptions = switches.map((item: any) => ({
          label: item.description,
          value: item.id,
        }));
      }
    }
    // 重新设置交换机信息
    await updateSchema({
      field: 'switchId',
      componentProps: { options: switchOptions },
    });
  }
  function onSwitchesChange(_value: any) {
    resetBusinessTagOptions();
  }
  async function onOperateTypeChange(_value: any) {
    let value = await getFieldsValue();
    const operateType = value.operateType;
    // 业务端口
    if (operateType == Network_SpeedLimitJob_OperateTargetType.BUSINESS_TAG) {
      // 显示业务标签
      updateSchema([
        {
          field: 'businessTag',
          required: true,
          show: true,
        },
        {
          field: 'switchPort',
          show: false,
          rules: [],
        },
        {
          field: 'switchPortRange',
          show: false,
          rules: [],
        },
      ]);
      // 清空交换机端口
      // await setFieldsValue({
      //   switchPort: null,
      //   switchPortRange: null,
      // });
    } else {
      // 显示交换机端口
      updateSchema([
        {
          field: 'businessTag',
          show: false,
          required: false,
        },
        {
          field: 'switchPort',
          show: true,
          rules: getSwitchPortValueValidator(getFormHandleFn) as any,
        },
        {
          field: 'switchPortRange',
          show: true,
          rules: getSwitchPortValueValidator(getFormHandleFn) as any,
        },
      ]);
      // 清空业务标签
      // await setFieldsValue({
      //   businessTag: undefined,
      // });
    }
    clearValidate();
  }
  async function onDrawerOk() {
    try {
      setDrawerProps({
        confirmLoading: true,
      });
      await validateFields();
      let payload = (await GetFormData()) as any;
      let errInfo = ValidatePayload(payload);
      if (errInfo) {
        message.error(errInfo);
        return;
      }

      if (
        payload.jobType == Network_SpeedLimitJob_JobType.IMMEDIATE_TASK &&
        payload.status == Network_SpeedLimitJob_Status.ENABLE
      ) {
        Modal.confirm({
          title: '是否确认执行任务?',
          content: '当前任务为立即执行任务，是否立即执行',
          icon: createVNode(ExclamationCircleOutlined),
          okText: '立即执行',
          cancelText: '暂不执行',
          onCancel: async () => {
            payload.executeNow = Network_SpeedLimitJob_ExecuteNow.NOT_EXECUTE; //不立即执行
            //立即不执行
            if (data.isUpdate) {
              await UpdateSpeedLimitJob(data.record.id as number, payload);
            } else {
              await CreateSpeedLimitJob(payload);
            }
            message.success('操作成功');
            emit('success');
            closeDrawer();
          },
          onOk: async () => {
            payload.executeNow = Network_SpeedLimitJob_ExecuteNow.EXECUTE; // 立即执行
            if (data.isUpdate) {
              await UpdateSpeedLimitJob(data.record.id as number, payload);
            } else {
              await CreateSpeedLimitJob(payload);
            }
            message.success('操作成功');
            emit('success');
            closeDrawer();
          },
        });
      } else {
        if (data.isUpdate) {
          await UpdateSpeedLimitJob(data.record.id as number, payload);
        } else {
          await CreateSpeedLimitJob(payload);
        }
        message.success('操作成功');
        emit('success');
        closeDrawer();
      }
    } finally {
      setDrawerProps({
        confirmLoading: false,
      });
    }
  }
  async function resetBusinessTagOptions() {
    // 编辑回写数据，不清空设置的值
    if (!data.isRewriting) {
      await setFieldsValue({
        businessTag: undefined,
      });
    }
    updateSchema({
      field: 'businessTag',
      componentProps: {
        loading: true,
        notFoundContent: '数据加载中...',
        options: [],
      },
    });
    try {
      let options: any[] = [];
      let { switchId } = await getFieldsValue();
      if (switchId) {
        let res = await GetSwitchBusinessTags({ switchIds: [switchId] });
        if (res && res.tags && res.tags.length > 0) {
          options = res.tags.map((item: any) => ({
            label: item,
            value: item,
          }));
        }
      }
      // let emptyContent = h(Empty, { image: Empty.PRESENTED_IMAGE_SIMPLE });
      updateSchema([
        {
          field: 'businessTag',
          componentProps: {
            options,
            notFoundContent: '暂无数据',
          },
        },
      ]);
    } finally {
      updateSchema({
        field: 'businessTag',
        componentProps: {
          loading: false,
        },
      });
    }
  }
  // 计算限速值
  function CalLimitValue(type: number | undefined, limitValue: number | undefined): number | null {
    // 解除限速，返回null
    if (type == Network_SpeedLimitJob_LimitType.UNLIMIT) {
      return null;
    }
    // 限速，返回限速值
    if (type == Network_SpeedLimitJob_LimitType.LIMIT) {
      if (!limitValue) {
        // return speedLimitDefaultValue;
        return null;
      }
      return Number(limitValue);
    }
    return null;
  }
  function GetFormData() {
    let value = getFieldsValue();
    let payload = {
      name: value.name,
      status: value.status,
      ecdnRoomId: value.ecdnRoomId,
      jobType: Number(value.jobType),
      retryCount: Number(value.retryCount),
      maxExecuteDelayMinutes: Number(value.maxExecuteDelayMinutes),
      strategies: {
        immediateTaskStrategy: null as any,
        onceTaskStrategy: null as any,
        cycleTaskStrategy: null as any,
        limitTarget: {
          switchIds: value.switchId ? [value.switchId] : [],
          operateType: Number(value.operateType),
          switchPort: undefined,
          switchPortRange: undefined,
          businessTag: undefined,
        },
      },
    };
    // 业务标签类型
    if (value.operateType == Network_SpeedLimitJob_OperateTargetType.BUSINESS_TAG) {
      payload.strategies.limitTarget.businessTag = value.businessTag;
    } else {
      payload.strategies.limitTarget.switchPort = value.switchPort;
      payload.strategies.limitTarget.switchPortRange = value.switchPortRange;
    }
    // 判断Job类型
    if (value.jobType == Network_SpeedLimitJob_JobType.IMMEDIATE_TASK) {
      // 立即执行
      payload.strategies.immediateTaskStrategy = {
        limitType: {
          type: data.immediateTask.type,
          limitValue: CalLimitValue(data.immediateTask.type, data.immediateTask.limitValue),
        },
      };
    } else if (value.jobType == Network_SpeedLimitJob_JobType.ONCE_TASK) {
      // 定时任务(一次性)
      payload.strategies.onceTaskStrategy = {
        executeTime: data.onceTask.executeTime
          ? dayjs(data.onceTask.executeTime).format('YYYY-MM-DD HH:mm:00')
          : undefined,
        limitType: {
          type: data.onceTask.type,
          limitValue: CalLimitValue(data.onceTask.type, data.onceTask.limitValue),
        },
      };
    } else if (value.jobType == Network_SpeedLimitJob_JobType.CYCLE_TASK) {
      // 周期任务
      payload.strategies.cycleTaskStrategy = {
        dates: [],
      };
      data.cycleTask.forEach((item: any) => {
        const date = {
          startDate: item.dateRange[0] ? dayjs(item.dateRange[0]).format('YYYY-MM-DD') : undefined,
          endDate: item.dateRange[1] ? dayjs(item.dateRange[1]).format('YYYY-MM-DD') : undefined,
          times: item.timeRange.map((t: any) => ({
            time: t.time ? dayjs(t.time, 'HH:mm:00').format('HH:mm:00') : undefined,
            limitType: {
              type: t.type,
              limitValue: CalLimitValue(t.type, t.limitValue),
            },
          })),
        };
        payload.strategies.cycleTaskStrategy.dates.push(date);
      });
    }
    return payload;
  }
  function ValidatePayload(payload: any): string | null {
    if (!payload) {
      return '参数不能为空';
    }
    if (!payload.name) {
      return '请输入任务名称';
    }
    if (!payload.status && payload.status != Network_SpeedLimitJob_Status.DISABLE) {
      return '请选择任务状态';
    }
    if (!payload.ecdnRoomId) {
      return '请选择群号|节点|机房 信息';
    }
    if (!payload.strategies) {
      return '请选择任务类型';
    }
    if (
      !payload.strategies.limitTarget.switchIds ||
      payload.strategies.limitTarget.switchIds.length === 0
    ) {
      return '请选择需要限速的交换机';
    }
    if (!payload.strategies.limitTarget || !payload.strategies.limitTarget.operateType) {
      return '请选择操作类型';
    }
    if (
      payload.strategies.limitTarget.operateType ==
        Network_SpeedLimitJob_OperateTargetType.BUSINESS_TAG &&
      (!payload.strategies.limitTarget.businessTag ||
        payload.strategies.limitTarget.businessTag.length === 0)
    ) {
      return '请选择业务标签';
    }
    if (
      payload.strategies.limitTarget.operateType ==
        Network_SpeedLimitJob_OperateTargetType.SWITCH_PORT &&
      !payload.strategies.limitTarget.switchPort &&
      !payload.strategies.limitTarget.switchPortRange
    ) {
      return '请填写交换机端口或交换机端口范围';
    }

    if (payload.jobType == Network_SpeedLimitJob_JobType.IMMEDIATE_TASK) {
      if (!payload.strategies.immediateTaskStrategy) {
        return '请填写立即执行配置';
      }
      if (
        !payload.strategies.immediateTaskStrategy.limitType ||
        !payload.strategies.immediateTaskStrategy.limitType.type
      ) {
        return '请选操作类型';
      }
      if (
        payload.strategies.immediateTaskStrategy.limitType.type ==
          Network_SpeedLimitJob_LimitType.LIMIT &&
        !payload.strategies.immediateTaskStrategy.limitType.limitValue
      ) {
        return '请填写限速值';
      }
    } else if (payload.jobType == Network_SpeedLimitJob_JobType.ONCE_TASK) {
      if (!payload.strategies.onceTaskStrategy) {
        return '请填写定时执行配置';
      }
      if (
        !payload.strategies.onceTaskStrategy.limitType ||
        !payload.strategies.onceTaskStrategy.limitType.type
      ) {
        return '请选操作类型';
      }
      if (!payload.strategies.onceTaskStrategy.executeTime) {
        return '请填写操作时间';
      }
      if (
        payload.strategies.onceTaskStrategy.limitType.type ==
          Network_SpeedLimitJob_LimitType.LIMIT &&
        !payload.strategies.onceTaskStrategy.limitType.limitValue
      ) {
        return '请填写限速值';
      }
    }
    if (payload.jobType == Network_SpeedLimitJob_JobType.CYCLE_TASK) {
      if (!payload.strategies.cycleTaskStrategy) {
        return '请填写周期任务配置';
      }
      if (
        !payload.strategies.cycleTaskStrategy.dates ||
        payload.strategies.cycleTaskStrategy.dates.length === 0
      ) {
        return '请填写周期任务配置';
      }
      for (let i = 0; i < payload.strategies.cycleTaskStrategy.dates.length; i++) {
        let date = payload.strategies.cycleTaskStrategy.dates[i];
        if (!date.startDate || !date.endDate) {
          return `请填写执行周期${i + 1}的开始日期和结束日期`;
        }
        for (let j = 0; j < date.times.length; j++) {
          let time = date.times[j];
          if (!time.time) {
            return `请填写执行周期${i + 1}-操作${j + 1}的操作时间`;
          }
          if (!time.limitType || !time.limitType.type) {
            return `请选执行周期${i + 1}-操作${j + 1}的操作类型`;
          }
          if (
            time.limitType.type == Network_SpeedLimitJob_LimitType.LIMIT &&
            !time.limitType.limitValue
          ) {
            return `请填写执行周期${i + 1}-操作${j + 1}的限速值`;
          }
        }
      }
    }
    return null;
  }
</script>

<style scoped></style>
