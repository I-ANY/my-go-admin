<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="NetworkPermissionCodeEnum.SPEED_LIMIT_JOB_ADD"
          @click="handleAddSpeedLimitJob"
          >添加任务</a-button
        >
      </template>
      <template #bodyCell="{ column, record }">
        <!-- 使用枚举值展示 -->
        <template v-for="(columnName, i) in Object.keys(data.showEnumFields)" :key="i">
          <template v-if="column.key == columnName">
            <span v-if="data.showEnumFields[columnName][record[columnName]]">
              {{ data.showEnumFields[columnName][record[columnName]].dictLabel }}
            </span>
            <span v-else>{{ record[columnName] }}</span>
          </template>
        </template>

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
        <template v-if="column.key === 'remark'">
          <Tooltip v-if="record.remark" :title="record.remark" placement="topLeft">
            {{ record.remark }}
          </Tooltip>
        </template>

        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                tooltip: '编辑任务',
                onClick: handleEditSpeedLimitJob.bind(null, record),
                label: '编辑',
                auth: NetworkPermissionCodeEnum.SPEED_LIMIT_JOB_EDIT,
              },
              {
                tooltip: '查看执行记录',
                onClick: handleViewExecRecord.bind(null, record),
                label: '执行记录',
                auth: NetworkPermissionCodeEnum.SPEED_LIMIT_JOB_VIEW_EXEC_RECORD,
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <SpeedLimitJobDrawer
      @register="registerSpeedLimitJobDrawer"
      @success="
        () => {
          reload();
        }
      "
    />
    <SpeedLimitExecRecordDrawer
      @register="registerSpeedLimitExecRecordDrawer"
      @success="
        () => {
          reload();
        }
      "
    />
  </div>
</template>
<script setup lang="ts">
  import { reactive } from 'vue';
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { GetSpeedLimitJobList } from '@/api/network/speed_limit_job';
  import { NetworkPermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import { Tag, Tooltip } from 'ant-design-vue';
  import {
    getSpeedLimitJobColumns,
    getSpeedLimitJobSearchFormSchema,
    speedLimitJobTypeMap,
    roomChargeModeMap,
    speedLimitJobStatusMap,
    speedLimitJobExecStatusMap,
    roomTypeMap,
    ispMap,
  } from './data';
  import SpeedLimitJobDrawer from './SpeedLimitJobDrawer.vue';
  import { useDrawer } from '@/components/Drawer';
  import SpeedLimitExecRecordDrawer from './SpeedLimitExecRecordDrawer.vue';

  defineOptions({ name: 'NetworkSpeedLimitJob' });
  const [registerSpeedLimitJobDrawer, { openDrawer: openSpeedLimitJobDrawer }] = useDrawer();
  const [registerSpeedLimitExecRecordDrawer, { openDrawer: openSpeedLimitExecRecordDrawer }] =
    useDrawer();
  const data = reactive({
    showEnumFields: {
      isp: ispMap,
    },
    showTagFields: {
      roomType: roomTypeMap,
      chargeMode: roomChargeModeMap,
      jobType: speedLimitJobTypeMap,
      status: speedLimitJobStatusMap,
      lastExecuteStatus: speedLimitJobExecStatusMap,
    },
  });
  const [registerTable, { reload }] = useTable({
    title: '定时任务列表',
    api: GetSpeedLimitJobList,
    columns: getSpeedLimitJobColumns(),
    formConfig: {
      labelWidth: 120,
      schemas: getSpeedLimitJobSearchFormSchema(),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      submitOnReset: false,
      alwaysShowLines: 1,
    },
    beforeFetch: (params) => {
      params = formatParams(params);
      return params;
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    immediate: true,
    rowKey: 'id',
    size: 'small',
    actionColumn: {
      width: 160,
      title: '操作',
      dataIndex: 'action',
      fixed: 'right',
    },
    // clickToRowSelect: false,
    // showSelectionBar: true,
  });

  function formatParams(params: Recordable) {
    return params;
  }

  async function handleAddSpeedLimitJob() {
    openSpeedLimitJobDrawer(true, {
      isUpdate: false,
    });
  }
  function handleEditSpeedLimitJob(record: Recordable) {
    openSpeedLimitJobDrawer(true, {
      record: record,
      isUpdate: true,
    });
  }
  function handleViewExecRecord(record: Recordable) {
    openSpeedLimitExecRecordDrawer(true, {
      record: record,
    });
  }
</script>
