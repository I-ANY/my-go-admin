<template>
  <PageWrapper :title="hostname + ' 巡检记录'" contentFullHeight contentBackground @back="goBack">
    <div>
      <BasicTable @register="registerTable" @change="handleTableChange">
        <template #toolbar>
          <a-button
            type="primary"
            v-auth="PermissionCodeEnum.BUSINESS_SERVER_INSPECT_RECORD_EXPORT"
            @click="handleDataExport"
            :loading="data.exporting"
            >{{ data.exportButTitle }}</a-button
          ></template
        >
        <template #bodyCell="{ column, record }">
          <template v-if="column.key == 'taskId'">
            <Tooltip
              v-if="hasPermission(PermissionCodeEnum.BUSINESS_INSPECT_TASK_RESULT)"
              placement="top"
              title="查看详情"
            >
              <a @click="onTaskIdClick(record)">#{{ record.taskId }}</a>
            </Tooltip>
            <span v-else>#{{ record.id }}</span>
          </template>
          <template v-if="column.key == 'status'">
            <Tooltip :overlayStyle="{ maxWidth: '600px' }" :title="record.remark">
              <Tag
                v-if="execStatusMap[record.status]"
                :color="execStatusMap[record.status].color || 'default'"
                style="font-weight: bold"
                >{{ execStatusMap[record.status].dictLabel }}</Tag
              >
              <span v-else>{{ record.status }}</span>
            </Tooltip>
          </template>
          <template v-if="column.key == 'inspectResult'">
            <Tag
              v-if="inspectResultMap[record.inspectResult]"
              :color="inspectResultMap[record.inspectResult].color || 'default'"
              style="font-weight: bold"
              >{{ inspectResultMap[record.inspectResult].dictLabel }}</Tag
            >
            <span v-else>{{ record.inspectResult }}</span>
          </template>

          <!-- 使用枚举值tag展示 -->
          <template v-for="(c, i) in Object.keys(allData.enumTagColumn)" :key="i">
            <template v-if="column.key == c">
              <Tag
                v-if="allData.enumTagColumn[c][record[c]]"
                :color="allData.enumTagColumn[c][record[c]].color || 'default'"
                style="font-weight: bold"
                >{{ allData.enumTagColumn[c][record[c]].dictLabel }}</Tag
              >
              <span v-else>{{ record[c] }}</span>
            </template>
          </template>

          <!-- 使用枚举值展示 -->
          <template v-for="(c, i) in Object.keys(allData.enumNormalColumn)" :key="i">
            <template v-if="column.key == c">
              <span v-if="allData.enumTagColumn[c][record[c]]">{{
                allData.enumNormalColumn[c][record[c]].dictLabel
              }}</span>
              <span v-else>{{ record[c] }}</span>
            </template>
          </template>

          <!-- 显示detail -->
          <template v-for="(c, i) in allData.showDetailColumn" :key="i">
            <template v-if="column.key == c">
              <span v-if="record[c]" :style="{ color: record[c].value == 1 ? 'green' : 'red' }">{{
                record[c].detail
              }}</span>
              <span v-else></span>
            </template>
          </template>

          <!-- 显示正常值和异常值 -->
          <template v-for="(c, i) in allData.normalAndAnomalyColumn" :key="i">
            <template v-if="column.key == c">
              <Tooltip
                v-if="record[c]"
                :overlayStyle="getOverlayStyle(splitByLineAndTrim(record[c].detail)?.length)"
              >
                <span :style="{ color: record[c].value == 1 ? 'green' : 'red' }">{{
                  record[c].value ? '正常' : '异常'
                }}</span>
                <template #title>
                  <p
                    :style="allData.pStyle"
                    v-for="(item, index) in splitByLineAndTrim(record[c].detail)"
                    :key="index"
                    >{{ item }}</p
                  >
                </template>
              </Tooltip>
              <span v-else></span>
            </template>
          </template>

          <!-- 基础输出型字段 -->
          <template v-for="(c, i) in allData.outPutColumn" :key="i">
            <template v-if="column.key == c">
              <Tooltip :title="'查看详情'" v-if="record[c]">
                <a @click="handShowLog(record[c])">查看</a>
              </Tooltip>
              <span v-else></span>
            </template>
          </template>

          <!-- 扩展检查字段 -->
          <template v-if="column.key == 'checkFields'">
            <div class="extend-field">
              <div v-for="(field, iii) in record.checkFields" :key="iii">
                <Tooltip :overlayStyle="getOverlayStyle(splitByLineAndTrim(field.detail)?.length)">
                  <span :style="{ color: field.value == 1 ? 'green' : 'red' }">{{
                    field.name
                  }}</span>
                  <template #title>
                    <p
                      :style="allData.pStyle"
                      v-for="(item, index) in splitByLineAndTrim(field.detail)"
                      :key="index"
                      >{{ item }}</p
                    >
                  </template>
                </Tooltip>
                <Divider
                  v-if="iii < record.checkFields?.length - 1"
                  type="vertical"
                  style="margin: 0 5px; background-color: grey"
                />
              </div>
            </div>
          </template>

          <!-- 扩展输出字段 -->
          <template v-if="column.key == 'outputFields'">
            <div class="extend-field">
              <div v-for="(field, ii) in record.outputFields" :key="ii">
                <Tooltip :overlayStyle="getOverlayStyle(0)" title="查看详情">
                  <a :color="field.value == 1 ? 'green' : 'red'" @click="handShowLog(field)">{{
                    field.name
                  }}</a>
                </Tooltip>
                <Divider
                  v-if="ii < record.outputFields?.length - 1"
                  type="vertical"
                  style="margin: 0 5px; background-color: grey"
                />
              </div>
            </div>
          </template>
        </template>
      </BasicTable>
      <LogModal @register="registerLogModal" @success="handleSuccess" />
    </div>
  </PageWrapper>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import {
    inspectRecordColumns,
    inspectRecordSearchFormSchema,
    execStatusMap,
    inspectResultMap,
  } from './data';
  import { Api, GetResultList, GetResultFieldName } from '@/api/business/inspect';
  import { Tag, Tooltip, Divider, notification, Modal } from 'ant-design-vue';
  import { useRoute } from 'vue-router';
  import { h, nextTick, onMounted, reactive } from 'vue';
  import { splitByLineAndTrim } from '@/utils/util';
  import LogModal from '../result/LogModal.vue';
  import { useModal } from '@/components/Modal';
  import { useGo } from '@/hooks/web/usePage';
  import { usePermission } from '@/hooks/web/usePermission';
  import { PermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import { PageWrapper } from '@/components/Page';
  import { RangePickPresetsExact } from '@/utils/common';
  import dayjs from 'dayjs';
  import { useCommonStore } from '@/store/modules/common';
  import { downloadFileByUrl } from '@/utils/download';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';

  const { hasPermission } = usePermission();
  const go = useGo();
  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
    filters: {} as any,
  });

  const [registerLogModal, { openModal: openLogModal }] = useModal();
  const route = useRoute();
  let hostname: string = route.query.hostname as string;
  const allData = reactive({
    pStyle: { display: 'block', margin: '0', padding: '0' },
    normalAndAnomalyColumn: [
      'rateCompare',
      'dialStatus',
      'pingPacketLoss',
      'businessStatus',
      'businessConfiguration',
      'diskMountStatus',
      'ispConsistency',
      'businessTrafficStatus',
      'monitorpyStatus',
      'natType',
    ], // 显示正常、异常的列
    showDetailColumn: ['realTimeRate', 'systemLoad', 'rootPartitionUsage'], // 显示详情的列
    enumTagColumn: {}, // 使用枚举显示Tag的列
    enumNormalColumn: {}, // 使用枚举显示的列
    outPutColumn: ['lbCpu', 'nicSpeed', 'kernelLogs', 'businessLog'], // 输出的列，点击查看日志
  });

  defineOptions({ name: 'InspectRecord' });
  const [registerTable, { getForm, reload }] = useTable({
    title: '巡检记录',
    api: GetResultList,
    beforeFetch: (params) => {
      params.hostnames = [hostname];
      return params;
    },
    columns: inspectRecordColumns,
    formConfig: {
      labelWidth: 120,
      schemas: inspectRecordSearchFormSchema(onTimePikerOpen, checkFieldsChange),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 2,
      submitOnReset: false,
      alwaysShowLines: 1,
      showResetButton: false,
      resetFunc() {
        nextTick(async () => {
          await resetReportTime();
          resetCheckFields();
        });
        return Promise.resolve();
      },
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      pageSizeOptions: ['10', '50', '100'],
    },
    rowKey: 'id',
    // actionColumn: {
    //   width: 120,
    //   title: '操作',
    //   dataIndex: 'action',
    //   // slots: { customRender: 'action' },
    //   fixed: 'right',
    // },
  });
  function handShowLog(field: Recordable) {
    openLogModal(true, field);
  }
  function handleSuccess() {
    reload();
  }
  onMounted(async () => {
    await resetReportTime();
    resetCheckFields();
  });
  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      startTimeBegin: dayjs(
        dayjs().add(-30, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      startTimeEnd: dayjs(dayjs().format('YYYY-MM-DD 23:59:59'), 'YYYY-MM-DD HH:mm:ss'),
    });
  }
  function onTimePikerOpen() {
    getForm().updateSchema({
      field: '[startTimeBegin, startTimeEnd]',
      componentProps: {
        presets: RangePickPresetsExact(),
      },
    });
  }
  function onTaskIdClick(record: Recordable) {
    const commonStore = useCommonStore();
    commonStore.setInspectHostname(hostname);
    go('/business/inspect/task/' + record.taskId + '/result');
  }
  function goBack() {
    go('/business/inspect/server');
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
  function handleTableChange(pagination, filters) {
    data.filters = filters;
  }

  async function handleDataExport() {
    Modal.confirm({
      title: '是否确认导出筛选的数据?',
      icon: h(ExclamationCircleOutlined),
      content: h(
        'div',
        { style: 'color:red;' },
        '导出将花费一些时间，请耐心等待！如果数据过大，可能会导出失败，请缩小导出条件后重试',
      ),
      async onOk() {
        await getForm().validate();
        let value = await getForm().getFieldsValue();
        value.hostnames = [hostname];
        value = { ...value, ...data.filters };
        nextTick(() => {
          data.exporting = true;
          data.exportButTitle = '导出中';
        });
        await exportInspectResult(value);
      },
    });
  }
  async function exportInspectResult(value: Recordable) {
    console.log('value', value);

    try {
      let filename = await downloadFileByUrl(Api.ExportInspectResult, 'POST', 5 * 60, value, null);
      notification.success({
        message: '导出成功',
        description: '文件名：' + filename,
        duration: null,
      });
    } catch (error) {
      notification.error({
        message: '导出失败',
        description: error.message,
        duration: null,
      });
    } finally {
      nextTick(() => {
        data.exporting = false;
        data.exportButTitle = '导出数据';
      });
    }
  }
  async function resetCheckFields() {
    // 获取字段名
    let res = await GetResultFieldName({ hostnames: [hostname], type: 1 });
    let checkFieldsOps: Array<any> = [];
    if (res?.names?.length > 0) {
      res.names.forEach((item) => {
        checkFieldsOps.push({
          label: item,
          value: item,
        });
      });
    }
    getForm().updateSchema({
      field: 'checkFields',
      componentProps: {
        options: checkFieldsOps,
      },
    });
  }
  async function checkFieldsChange() {
    let values = await getForm().getFieldsValue();
    let checkFields = values?.checkFields;
    if (checkFields && checkFields.length > 0) {
      await getForm().updateSchema({
        field: 'checkFieldStatus',
        componentProps: {
          disabled: false,
        },
        required: true,
      });
    } else {
      await getForm().updateSchema({
        field: 'checkFieldStatus',
        componentProps: {
          disabled: true,
        },
        required: false,
      });
      getForm().setFieldsValue({
        checkFieldStatus: null,
      });
    }
    getForm().clearValidate();
  }
</script>
<style scoped>
  ::v-deep(form.ant-form-default) {
    margin-bottom: 0;
    padding-top: 5px;
    padding-bottom: 5px;
  }

  ::v-deep(div.vben-basic-table.vben-basic-table-form-container) {
    padding-top: 8px;
    padding-bottom: 8px;
  }

  ::v-deep(div.ant-page-header) {
    padding-top: 8px;
    padding-bottom: 8px;
  }

  .extend-field {
    display: flex;
    flex-wrap: wrap;
    justify-content: center; /* 水平居中 */
    text-align: center;
    gap: 0;
  }
</style>
