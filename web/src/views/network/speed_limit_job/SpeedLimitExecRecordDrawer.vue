<template>
  <BasicDrawer v-bind="$attrs" @register="registerModal" @close="onDrawerCancel">
    <Divider orientation="left" style="margin: 0">任务配置</Divider>
    <div class="pl-18px pr-18px pt-10px pb-10px w-full">
      <Descriptions :column="4" size="small" bordered>
        <DescriptionsItem label="任务名称">{{ data.record?.name || '-' }}</DescriptionsItem>
        <DescriptionsItem label="机房名称">{{ data.record?.roomName || '-' }}</DescriptionsItem>
        <DescriptionsItem label="交换机名称">
          {{ data.record?.switchDesc || '-' }}
        </DescriptionsItem>
        <DescriptionsItem label="状态">
          <Tag
            v-if="speedLimitJobStatusMap[data.record?.status]"
            :color="speedLimitJobStatusMap[data.record?.status]?.color || 'default'"
          >
            {{ speedLimitJobStatusMap[data.record?.status]?.dictLabel }}
          </Tag>
          <span v-else>{{ data.record?.status || '-' }}</span>
        </DescriptionsItem>
        <DescriptionsItem label="任务类型">
          <Tag
            v-if="speedLimitJobTypeMap[data.record?.jobType]"
            :color="speedLimitJobTypeMap[data.record?.jobType]?.color || 'default'"
          >
            {{ speedLimitJobTypeMap[data.record?.jobType]?.dictLabel }}
          </Tag>
          <span v-else>{{ data.record?.jobType || '-' }}</span>
        </DescriptionsItem>
        <DescriptionsItem label="重试次数">
          {{ data.record?.retryCount ?? '-' }}
        </DescriptionsItem>
        <DescriptionsItem label="最大延迟">
          {{
            data.record?.maxExecuteDelayMinutes ? `${data.record.maxExecuteDelayMinutes} 分钟` : '-'
          }}
        </DescriptionsItem>
        <DescriptionsItem label="操作对象">
          <Tag
            v-if="speedLimitTargetTypeMap[data.record?.strategies?.limitTarget?.operateType]"
            :color="
              speedLimitTargetTypeMap[data.record?.strategies?.limitTarget?.operateType]?.color ||
              'default'
            "
          >
            {{
              speedLimitTargetTypeMap[data.record?.strategies?.limitTarget?.operateType]?.dictLabel
            }}
          </Tag>
          <span v-else>{{ data.record?.strategies?.limitTarget?.operateType || '-' }}</span>
        </DescriptionsItem>
        <!-- 交换机端口类型时显示端口信息 -->
        <DescriptionsItem
          label="交换机端口"
          v-if="
            data.record?.strategies?.limitTarget?.operateType ==
            Network_SpeedLimitJob_OperateTargetType.SWITCH_PORT
          "
        >
          {{ data.record?.strategies?.limitTarget?.switchPort || '-' }}
        </DescriptionsItem>
        <DescriptionsItem
          label="端口范围"
          v-if="
            data.record?.strategies?.limitTarget?.operateType ==
            Network_SpeedLimitJob_OperateTargetType.SWITCH_PORT
          "
        >
          {{ data.record?.strategies?.limitTarget?.switchPortRange || '-' }}
        </DescriptionsItem>
        <!-- 业务标签类型时显示标签信息 -->
        <DescriptionsItem
          label="业务标签"
          :span="2"
          v-if="
            data.record?.strategies?.limitTarget?.operateType ==
            Network_SpeedLimitJob_OperateTargetType.BUSINESS_TAG
          "
        >
          <Tag
            v-for="tag in data.record?.strategies?.limitTarget?.businessTag"
            :key="tag"
            color="blue"
          >
            {{ tag }}
          </Tag>
          <span v-if="!data.record?.strategies?.limitTarget?.businessTag?.length">-</span>
        </DescriptionsItem>
      </Descriptions>
    </div>
    <Divider orientation="left" style="margin: 0">执行记录</Divider>
    <BasicTable @register="registerTable" style="padding-top: 2px; padding-bottom: 2px">
      <template #bodyCell="{ column, record }">
        <!-- 使用枚举值tag展示 -->
        <template v-for="(columnName, i) in Object.keys(data.showTagFields)" :key="i">
          <template v-if="column.key == columnName">
            <Tag
              style="font-weight: bold"
              v-if="data.showTagFields[columnName][record[columnName]]"
              :color="data.showTagFields[columnName][record[columnName]].color || 'default'"
              >{{ data.showTagFields[columnName][record[columnName]].dictLabel }}</Tag
            >
            <span v-else>{{ record[columnName] }}</span>
          </template>
        </template>
        <template v-if="column.key === 'consoleLog'">
          <Tooltip placement="top" :overlayStyle="getOverlayStyle(0)">
            <a>查看</a>
            <template #title>
              <p
                :style="data.pStyle"
                v-for="(item, index) in splitByLineAndTrim(record.consoleLog)"
                :key="index"
                >{{ item }}</p
              >
            </template>
          </Tooltip>
        </template>
        <template v-if="column.key === 'message'">
          <Tooltip v-if="record.message" :title="record.message" placement="topLeft">
            <span>{{ record.message }}</span>
          </Tooltip>
        </template>
      </template>
    </BasicTable>
  </BasicDrawer>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { BasicDrawer, useDrawerInner } from '@/components/Drawer';
  import { GetSpeedLimitJobExecRecordList } from '@/api/network/speed_limit_job';
  import {
    getSpeedLimitExecRecordColumns,
    getSpeedLimitExecRecordSearchForm,
    speedLimitJobExecStatusMap,
    speedLimitJobLimitTypeMap,
    speedLimitJobStatusMap,
    speedLimitJobTypeMap,
    speedLimitTargetTypeMap,
  } from './data';
  import { Tooltip, Tag, Divider, Descriptions, DescriptionsItem } from 'ant-design-vue';
  import { nextTick, reactive } from 'vue';
  import dayjs from 'dayjs';
  import { splitByLineAndTrim } from '@/utils/util';
  import { Network_SpeedLimitJob_OperateTargetType } from '@/enums/dictValueEnum';

  defineOptions({ name: 'SpeedLimitExecRecordDrawer' });
  const emit = defineEmits(['register', 'success']);
  const data = reactive({
    record: null as any,
    showTagFields: {
      executeStatus: speedLimitJobExecStatusMap,
      limitType: speedLimitJobLimitTypeMap,
    },
    pStyle: { display: 'block', margin: '0', padding: '0' },
  });
  const [registerModal, { setDrawerProps, closeDrawer }] = useDrawerInner(async (d) => {
    data.record = d.record;
    setDrawerProps({
      title: `执行记录-${data?.record?.name}(${data?.record?.roomName})`,
      width: 1400,
      // height: 630,
      destroyOnClose: true,
      showCancelBtn: true,
      showOkBtn: false,
      showFooter: false,
    });
    await resetOperTime();
    reload();
  });

  const [registerTable, { reload, getForm, getDataSource }] = useTable({
    // title: '机房限速记录',
    columns: getSpeedLimitExecRecordColumns(safeGetDataSource),
    api: async (params) => {
      const res = await GetSpeedLimitJobExecRecordList(data?.record.id, params);
      return res;
    },
    formConfig: {
      labelWidth: 120,
      schemas: getSpeedLimitExecRecordSearchForm(onTimePikerOpen),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      compact: true,
      // alwaysShowLines: 2,
      resetFunc() {
        nextTick(async () => {
          await resetOperTime();
        });
        return Promise.resolve();
      },
    },
    useSearchForm: true,
    showTableSetting: false,
    immediate: false,
    size: 'small',
    canResize: true,
    // scroll: { y: 'max-content' },
    bordered: true,
    showIndexColumn: false,
    rowKey: 'id',
    pagination: {
      // pageSizeOptions: ['1', '2', '5'],
    },
  });
  function safeGetDataSource() {
    try {
      return getDataSource() || [];
    } catch (error) {
      console.warn('getDataSource error:', error);
      return [];
    }
  }
  function resetOperTime(): Promise<void> {
    return getForm().setFieldsValue({
      executeTimeBegin: dayjs(
        dayjs().add(-15, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      executeTimeEnd: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    });
  }
  function onTimePikerOpen() {
    getForm().updateSchema({
      field: '[executeTimeBegin, executeTimeEnd]',
      componentProps: {
        presets: [
          {
            label: '今天',
            value: [
              dayjs(dayjs().format('YYYY-MM-DD 00:00:00'), 'YYYY-MM-DD HH:mm:ss'),
              dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
            ],
          },
          {
            label: '昨天',
            value: [
              dayjs(
                dayjs().subtract(1, 'day').format('YYYY-MM-DD 00:00:00'),
                'YYYY-MM-DD HH:mm:ss',
              ),
              dayjs(
                dayjs().subtract(1, 'day').format('YYYY-MM-DD 23:59:59'),
                'YYYY-MM-DD HH:mm:ss',
              ),
            ],
          },
          {
            label: '本周',
            value: [
              dayjs(
                dayjs().startOf('week').add(1, 'day').format('YYYY-MM-DD 00:00:00'),
                'YYYY-MM-DD HH:mm:ss',
              ),
              dayjs(dayjs().endOf('week').format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
            ],
          },
          { label: '最近1小时', value: [dayjs().add(-1, 'h'), dayjs()] },
          { label: '最近3小时', value: [dayjs().add(-3, 'h'), dayjs()] },
          { label: '最近6小时', value: [dayjs().add(-6, 'h'), dayjs()] },
          { label: '最近8小时', value: [dayjs().add(-8, 'h'), dayjs()] },
          { label: '最近12小时', value: [dayjs().add(-12, 'h'), dayjs()] },
          { label: '最近1天', value: [dayjs().add(-1, 'd'), dayjs()] },
          { label: '最近3天', value: [dayjs().add(-3, 'd'), dayjs()] },
          { label: '最近7天', value: [dayjs().add(-7, 'd'), dayjs()] },
        ],
      },
    });
  }
  function getOverlayStyle(len: number | undefined) {
    let style: any = {
      maxWidth: '1200px',
      maxHeight: '700px',
    };
    if (len == undefined || len < 30) {
      return style;
    }
    style.overflow = 'auto';
    return style;
  }
  function onDrawerCancel() {
    closeDrawer();
    emit('success');
  }
</script>
<style></style>
