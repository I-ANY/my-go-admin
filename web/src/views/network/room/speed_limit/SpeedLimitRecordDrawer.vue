<template>
  <BasicDrawer v-bind="$attrs" @register="registerModal" @cancel="onModalCancel">
    <Divider orientation="left" style="margin: 0">限速配置</Divider>
    <div style="margin: 20px 20px 35px">
      <Row :gutter="16" style="margin-bottom: 12px">
        <Col :span="5"
          ><div class="message-item">
            <div class="message-item-label">限速开关：</div>
            <div class="message-item-value">
              <Tag
                style="font-weight: bold"
                v-if="data?.record?.speedLimitConfig?.status == limitStatusEnum.LIMIT"
                :color="speedLimitStausMap[limitStatusEnum.LIMIT].color || 'default'"
                >{{ speedLimitStausMap[limitStatusEnum.LIMIT].dictLabel }}</Tag
              >
              <Tag
                v-else
                style="font-weight: bold"
                :color="speedLimitStausMap[limitStatusEnum.UNLIMITED].color"
                >{{ speedLimitStausMap[limitStatusEnum.UNLIMITED].dictLabel }}</Tag
              >
            </div>
          </div></Col
        >
        <Col :span="5"
          ><div class="message-item"
            ><div class="message-item-label">限速值(Gbps)：</div>
            <div class="message-item-value">{{
              getNumberResult(
                customCeilDivide(
                  data?.record?.speedLimitConfig?.limitValue,
                  1000 * 1000 * 1000,
                  2,
                ) as any,
              )
            }}</div></div
          ></Col
        >
        <Col :span="5"
          ><div class="message-item">
            <div class="message-item-label">开启限速阈值(Gbps)：</div>
            <div class="message-item-value">
              {{
                getNumberResult(
                  customCeilDivide(
                    data?.record?.speedLimitConfig?.limitThreshold,
                    1000 * 1000 * 1000,
                    2,
                  ) as any,
                )
              }}</div
            ></div
          ></Col
        >
        <Col :span="5"
          ><div class="message-item">
            <div class="message-item-label">解除限速阈值(Gbps)：</div>
            <div class="message-item-value">
              {{
                getNumberResult(
                  customCeilDivide(
                    data?.record?.speedLimitConfig?.unlimitThreshold,
                    1000 * 1000 * 1000,
                    2,
                  ) as any,
                )
              }}</div
            ></div
          ></Col
        >
      </Row>
      <Row style="margin-bottom: 12px">
        <Col :span="24"
          ><div class="message-item">
            <div class="message-item-label">排除交换机：</div>
            <!-- <div class="message-item-value">{{
              getSwitchNames(data?.record?.speedLimitConfig?.excludeSwitches)
            }}</div> -->
            <div class="message-item-value">
              <Tag v-for="item in data?.record?.speedLimitConfig?.excludeSwitches" :key="item.id">{{
                item.description
              }}</Tag>
            </div>
          </div></Col
        >
      </Row>
      <Row style="margin-bottom: 12px">
        <Col :span="24"
          ><div class="message-item">
            <div class="message-item-label">总流量交换机：</div>
            <!-- <div class="message-item-value">{{
              getSwitchNames(data?.record?.speedLimitConfig?.totalTrafficSwitches)
            }}</div> -->
            <div class="message-item-value">
              <Tag
                v-for="item in data?.record?.speedLimitConfig?.totalTrafficSwitches"
                :key="item.id"
                >{{ item.description }}</Tag
              >
            </div>
          </div></Col
        >
      </Row>
      <Row style="margin-bottom: 12px">
        <Col :span="24"
          ><div class="message-item">
            <div class="message-item-label">流量扣减交换机：</div>
            <div class="message-item-value">
              <Tag
                v-for="item in data?.record?.speedLimitConfig?.deductTrafficSwitches"
                :key="item.id"
                >{{ item.description }}</Tag
              >
            </div>
          </div></Col
        >
      </Row>
      <Row :gutter="16">
        <Col :span="24"
          ><div class="message-item">
            <div class="message-item-label">限速顺序：</div>
            <div class="message-item-value">{{
              data?.record?.speedLimitConfig?.businessSorts?.join(' -> ')
            }}</div>
            <!-- <div class="message-item-value">
              <Tag
                v-for="item in data?.record?.speedLimitConfig?.businessSorts"
                :key="item"
                :color="'blue'"
                >{{ item }}</Tag
              >
            </div> -->
          </div></Col
        >
      </Row>
    </div>
    <Divider orientation="left" style="margin-bottom: 0">限速记录</Divider>
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
        <template v-if="column.key === 'remark'">
          <Tooltip v-if="record.remark" :title="record.remark" placement="topLeft">
            <span>{{ record.remark }}</span>
          </Tooltip>
        </template>
      </template>
    </BasicTable></BasicDrawer
  >
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { BasicDrawer, useDrawerInner } from '@/components/Drawer';
  import {
    getSpeedLimitRecordColumns,
    getSpeedLimitRecordSearchForm,
    speedLimitOperStatusMap,
    speedLimitOperTypeMap,
    speedLimitStausMap,
    limitStatusEnum,
  } from './data';
  import { customCeilDivide, splitByLineAndTrim } from '@/utils/util';
  import { getNumberResult } from '../peak/data';
  import { getSpeedLimitRecordList } from '@/api/network/room_peak';
  import { nextTick, reactive } from 'vue';
  import { Tooltip, Tag, Divider, Row, Col } from 'ant-design-vue';
  import dayjs from 'dayjs';

  defineOptions({ name: 'SpeedLimitRecordModal' });
  const emit = defineEmits(['register', 'success']);
  const data = reactive({
    showTagFields: {
      status: speedLimitOperStatusMap,
      handleType: speedLimitOperTypeMap,
    },
    pStyle: { display: 'block', margin: '0', padding: '0' },
    record: null as any,
  });

  const [registerModal, { setDrawerProps, closeDrawer }] = useDrawerInner(async (d) => {
    data.record = d.record;
    setDrawerProps({
      title: `机房限速记录-${data?.record?.name}(${data?.record?.month})`,
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

  const [registerTable, { reload, getForm }] = useTable({
    // title: '机房限速记录',
    columns: getSpeedLimitRecordColumns(),
    api: getSpeedLimitRecordList,
    beforeFetch: (params) => {
      params.id = data?.record.id;
      return Promise.resolve(params);
    },

    formConfig: {
      labelWidth: 120,
      schemas: getSpeedLimitRecordSearchForm(onTimePikerOpen),
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
  function onModalCancel() {
    closeDrawer();
    emit('success');
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
  function resetOperTime(): Promise<void> {
    return getForm().setFieldsValue({
      operatorTimeBegin: dayjs(
        dayjs().add(-6, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      operatorTimeEnd: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    });
  }
  function onTimePikerOpen() {
    getForm().updateSchema({
      field: '[operatorTimeBegin, operatorTimeEnd]',
      componentProps: {
        disabledDate: (currentDate) => {
          return currentDate.format('YYYY-MM') !== data?.record.month;
        },
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
  // function getSwitchNames(excludeSwitches: any[] | null | undefined) {
  //   if (excludeSwitches == null || excludeSwitches.length == 0) {
  //     return '';
  //   }
  //   let names: string[] = [];
  //   excludeSwitches.forEach((excludeSwitch) => {
  //     names.push(excludeSwitch.description);
  //   });
  //   return names.join('；');
  // }
</script>
<style scoped>
  .message-item {
    display: flex;
    gap: 4px;
    align-items: flex-start;
    justify-content: left;
  }

  .message-item-label {
    width: 160px;
    font-weight: bold;
    text-align: right;
  }

  .message-item-value {
    flex: 1;
  }
</style>
