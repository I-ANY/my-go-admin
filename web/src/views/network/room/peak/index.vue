<template>
  <div>
    <BasicTable @register="registerTable">
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
        <template v-if="column.key === 'status'">
          <Tag
            style="font-weight: bold"
            v-if="record.speedLimitConfig?.status == limitStatusEnum.LIMIT + ''"
            :color="speedLimitStausMap[limitStatusEnum.LIMIT + ''].color || 'default'"
            >{{ speedLimitStausMap[limitStatusEnum.LIMIT + ''].dictLabel }}</Tag
          >
          <Tag
            style="font-weight: bold"
            v-else
            :color="speedLimitStausMap[limitStatusEnum.UNLIMITED + ''].color || 'default'"
            >{{ speedLimitStausMap[limitStatusEnum.UNLIMITED + ''].dictLabel }}</Tag
          >
        </template>

        <!-- 确认人 -->
        <template v-if="column.key === 'updateUser'">
          <span v-if="record.updateUser">{{ record?.updateUser?.nickName }}</span>
          <span v-else></span>
        </template>
        <!-- <template v-if="column.key === 'peakBusinessNames'">
          <Tooltip
            v-if="record.peakBusinessNames?.length >= 1"
            style="display: 'inline-block'; width: 100%"
            placement="topLeft"
          >
            <template #title>
              <p
                style="margin: 1px"
                v-for="(item, index) in record.peakBusinessNames"
                :key="index"
                >{{ item }}</p
              >
            </template>
            <span v-for="item in record.peakBusinessNames" :key="item">
              <Tag :key="item" color="blue"> {{ item }} </Tag>&nbsp;
            </span>
          </Tooltip>
        </template> -->
        <template v-if="column.key == 'peakBusinessNames'">
          <template v-if="record.peakBusinessNames && record.peakBusinessNames.length > 0">
            <div class="cell-tags">
              <Tag v-for="(item, index) in record.peakBusinessNames" :key="index" :color="'blue'">
                {{ item }}
              </Tag>
            </div>
          </template>
        </template>
        <!-- 备注 -->
        <template v-if="column.key === 'remark'">
          <Tooltip v-if="record.remark" :title="record.remark" placement="topLeft">
            {{ record.remark }}
            <!-- <div
              :style="{
                textAlign: 'left',
                wordWrap: 'break-word',
                whiteSpace: 'pre-wrap',
                overflowWrap: 'break-word',
              }"
              >{{
                record.remark?.length > 38 ? record.remark?.substring(0, 35) + '...' : record.remark
              }}</div
            > -->
          </Tooltip>
        </template>
        <!-- 计算时间 -->
        <template v-if="column.key === 'calculateTime'">
          <Tooltip title="超过60分钟未重新计算95值，请检查" v-if="calTimeIsOver60Minutes(record)">
            <span v-if="record.calculateTime" :style="getCalTimeStyle(record)">{{
              record.calculateTime
            }}</span>
          </Tooltip>
          <span v-else>{{ record.calculateTime }}</span>
        </template>
        <!-- 方案 -->
        <template v-if="column.key === 'plans'">
          <template v-if="record.plans && record.plans.length > 0">
            <Tooltip :overlayStyle="{ maxWidth: '900px' }">
              <div class="cell-tags">
                <Tag
                  title="方案"
                  v-for="(item, i) in record.plans"
                  :key="i"
                  color="blue"
                  :style="{ fontSize: '11px' }"
                  >{{ item }}</Tag
                >
              </div>
              <template #title>
                <p
                  v-for="(item, i) in record.plans"
                  :key="i"
                  :style="{ display: 'block', margin: '0', padding: '0' }"
                  >{{ item }}</p
                >
              </template>
            </Tooltip>
          </template>
        </template>
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'material-symbols:monitoring',
                tooltip: '95值拉升预测',
                onClick: handleImprove95Prediction.bind(null, record),
                label: '预测',
                ifShow: () => {
                  return record.useMode != roomUseModEnum.SINGLE_PORT_95_PEAK_SHAVING;
                },
              },
              {
                icon: 'mdi:chart-box-outline',
                tooltip: '查看单端口95详情',
                onClick: handleSinglePort95.bind(null, record),
                label: '端口95',
                auth: NetworkPermissionCodeEnum.SINGLE_PORT_95_PEAK_SHAVING_PORT,
                ifShow: () => {
                  return record.useMode == roomUseModEnum.SINGLE_PORT_95_PEAK_SHAVING;
                },
              },
              {
                icon: 'clarity:note-edit-line',
                tooltip: '编辑',
                onClick: handleEdit.bind(null, record),
                label: '编辑',
                auth: NetworkPermissionCodeEnum.NETWORK_ROOM_PEAK_EDIT,
              },
            ]"
            :dropDownActions="getDropDownActions(record)"
          />
        </template>
      </template>
    </BasicTable>
    <UpdateRoomPeakDrawer @register="registerDrawer" @success="onUpdateSuccess" />
    <Improve95PredictionModal
      @register="registerImprove95PredictionModal"
      @success="onImprove95PredictionSuccess"
    />
    <SpeedLimitConfiguretionDrawer
      @register="registerSpeedLimitConfiguretionDrawer"
      @success="onSpeedLimitConfiguretionSuccess"
    />
    <SpeedLimitRecordDrawer
      @register="registerSpeedLimitRecordDrawer"
      @success="onSpeedLimitRecordSuccess"
    />
    <SinglePort95Modal @register="registerSinglePort95Modal" @success="onSinglePort95Success" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import {
    columns,
    searchFormSchema,
    roomStatusMap,
    roomTypeMap,
    roomChargeModeMap,
    ispMap,
    roomPeakEnableStatusMap,
    roomUseModeMap,
    roomUseModEnum,
    getNumberResult,
  } from './data';
  import { speedLimitStausMap, limitStatusEnum } from '../speed_limit/data';
  import UpdateRoomPeakDrawer from './UpdateRoomPeakDrawer.vue';
  import { useDrawer } from '@/components/Drawer';
  import { nextTick, onMounted, reactive } from 'vue';
  import { Tooltip, Tag } from 'ant-design-vue';
  import dayjs from 'dayjs';
  import { getPeakBusiness, getPeakRoomList } from '@/api/network/room_peak';
  import { NetworkPermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import Improve95PredictionModal from './Improve95PredictionModal.vue';
  import { useModal } from '@/components/Modal';
  import SpeedLimitConfiguretionDrawer from '../speed_limit/SpeedLimitConfiguretionDrawer.vue';
  import SpeedLimitRecordDrawer from '../speed_limit/SpeedLimitRecordDrawer.vue';
  import SinglePort95Modal from './SinglePort95Modal.vue';
  import { usePermission } from '@/hooks/web/usePermission';
  import { customCeilDivide } from '@/utils/util';

  const { hasPermission } = usePermission();
  const [registerDrawer, { openDrawer }] = useDrawer();
  const [registerImprove95PredictionModal, { openModal: openImprove95PredictionModal }] =
    useModal();
  const [registerSinglePort95Modal, { openModal: openSinglePort95Modal }] = useModal();
  const [registerSpeedLimitRecordDrawer, { openDrawer: openSpeedLimitRecordDrawer }] = useDrawer();
  const [registerSpeedLimitConfiguretionDrawer, { openDrawer: openSpeedLimitConfiguretionDrawer }] =
    useDrawer();
  const data = reactive({
    showEnumFields: {
      isp: ispMap,
      roomType: roomTypeMap,
    },
    showTagFields: {
      roomStatus: roomStatusMap,
      chargeMode: roomChargeModeMap,
      peakEnabled: roomPeakEnableStatusMap,
      useMode: roomUseModeMap,
      // status: speedLimitStausMap,
    },
  });

  defineOptions({ name: 'NetworkRoomPeak' });
  const [registerTable, { getForm, reload }] = useTable({
    title: '机房列表',
    api: getPeakRoomList,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema(),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      submitOnReset: false,
      alwaysShowLines: 1,
      resetFunc() {
        nextTick(() => {
          resetReportTime();
        });
        return Promise.resolve();
      },
    },
    async beforeFetch(params: Recordable) {
      parseValue(params);
      return params;
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    immediate: false,
    rowKey: 'id',
    actionColumn: {
      width: 180,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: 'right',
    },
  });

  onMounted(async () => {
    await resetReportTime();
    reload();
    reloadPeakBusiness();
  });
  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      month: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    });
  }

  function parseValue(params: Recordable) {
    params.month = dayjs(params.month).format('YYYY-MM');
  }
  function handleEdit(record: Recordable) {
    openDrawer(true, {
      record,
    });
  }
  async function reloadPeakBusiness() {
    let params = getForm().getFieldsValue();
    const month = dayjs(params.month).format('YYYY-MM');
    const { business } = await getPeakBusiness({ recentMonth: 10, month });
    let businessOptions: any[] = [];
    business?.forEach((item) => {
      businessOptions.push({
        label: item,
        value: item,
      });
    });
    await getForm().updateSchema({
      field: 'peakBusiness',
      componentProps: {
        options: businessOptions,
      },
    });
  }
  function handleSpeedLimitConfiguretion(record: Recordable) {
    openSpeedLimitConfiguretionDrawer(true, {
      record,
    });
  }
  function handleOpenSpeedLimitRecord(record: Recordable) {
    openSpeedLimitRecordDrawer(true, {
      record,
    });
  }
  function handleSinglePort95(record: Recordable) {
    openSinglePort95Modal(true, {
      record,
    });
  }
  function onUpdateSuccess() {
    reload();
  }
  function onImprove95PredictionSuccess() {
    reload();
  }
  function onSpeedLimitConfiguretionSuccess() {
    reload();
  }
  function onSpeedLimitRecordSuccess() {
    reload();
  }
  function onSinglePort95Success() {
    reload();
  }
  function handleImprove95Prediction(record: Recordable) {
    const preset95Bw =
      record.preset95Bw == null || record.preset95Bw == undefined
        ? record.guaranteedBw
        : record.preset95Bw;
    let preset95Bwbps = getNumberResult(customCeilDivide(preset95Bw, 1000 * 1000 * 1000, 2) as any);
    const title = `95值拉升预测-${record.name}-预95值：${preset95Bwbps || ' - '}Gbps`;
    const improve95Prediction = record.improve95Prediction || [];
    openImprove95PredictionModal(true, {
      title,
      improve95Prediction,
    });
  }
  function getCalTimeStyle(record: Recordable) {
    if (calTimeIsOver60Minutes(record)) {
      return { color: 'red' };
    }
    return {};
  }
  function calTimeIsOver60Minutes(record: Recordable) {
    if (!record.calculateTime) {
      return false;
    }
    const diff = dayjs().unix() - dayjs(record.calculateTime).unix();
    return diff > 60 * 60;
  }
  function getDropDownActions(record: Recordable) {
    let actions: any[] = [];
    if (hasPermission(NetworkPermissionCodeEnum.NETWORK_ROOM_PEAK_SPEED_LIMIT_CONFIGURATION)) {
      actions.push({
        label: '限速配置',
        onClick: handleSpeedLimitConfiguretion.bind(null, record),
        auth: NetworkPermissionCodeEnum.NETWORK_ROOM_PEAK_SPEED_LIMIT_CONFIGURATION,
      });
    }
    if (hasPermission(NetworkPermissionCodeEnum.NETWORK_ROOM_PEAK_SPEED_LIMIT_RECORD)) {
      actions.push({
        label: '限速记录',
        onClick: handleOpenSpeedLimitRecord.bind(null, record),
        auth: NetworkPermissionCodeEnum.NETWORK_ROOM_PEAK_SPEED_LIMIT_RECORD,
      });
    }
    return actions;
  }
</script>

<style scoped>
  .cell-tags {
    display: flex;
    flex-wrap: wrap;
    gap: 2px;
    justify-content: center; /* 水平居中 */
  }
</style>
