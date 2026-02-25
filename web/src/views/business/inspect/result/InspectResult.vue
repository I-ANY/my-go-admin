<template>
  <PageWrapper
    :title="'#' + taskId + ' 巡检详情'"
    :contentFullHeight="false"
    contentBackground
    @back="goBack"
  >
    <div>
      <Row>
        <Col :span="11">
          <AnomalyServerPie
            :taskId="taskId"
            :data="data.pie.status.data"
            :legendLeft="'55%'"
            :seriesLeft="'40%'"
            :title="data.pie.status.title"
            :number="data.pie.status.total"
            ref="statusPieRef"
            @on-field-select="onStatusPieSelectField"
          />
        </Col>
        <Col :span="13">
          <AnomalyServerPie
            :taskId="taskId"
            :legendLeft="'30%'"
            :seriesLeft="'18%'"
            :data="data.pie.field.data"
            :title="data.pie.field.title"
            :number="data.pie.field.total"
            ref="fieldPieRef"
            @on-field-select="onFieldPieSelectField"
          />
        </Col>
      </Row>
    </div>
    <!-- <template #headerContent>
      <span class="m-20px">巡检时间：{{ allData.taskInfo?.startTime }}</span>
      <span>
        操作类型
        <Tag
          v-if="operatorTypeMap[allData.taskInfo?.operatorType]"
          :color="operatorTypeMap[allData.taskInfo?.operatorType].color || 'default'"
          >{{ operatorTypeMap[allData.taskInfo?.operatorType].dictLabel }}</Tag
        >
        <span v-else>{{ allData.taskInfo.operatorType }}</span>
      </span>
      <span>
        状态：
        <Tag
          v-if="execStatusMap[allData.taskInfo?.status]"
          :color="execStatusMap[allData.taskInfo?.status].color || 'default'"
          >{{ execStatusMap[allData.taskInfo?.status].dictLabel }}</Tag
        >
        <span v-else>{{ allData.taskInfo.status }}</span>
      </span>
      <span> 操作人：{{ allData.taskInfo?.operatorUser?.nickName }} </span>
    </template> -->
    <div>
      <BasicTable @register="registerTable" @change="handleTableChange">
        <template #toolbar>
          <a-button
            type="primary"
            v-auth="PermissionCodeEnum.BUSINESS_INSPECT_RESULT_EXPORT"
            @click="handleDataExport"
            :loading="data.exporting"
            >{{ data.exportButTitle }}</a-button
          ></template
        >
        <template #bodyCell="{ column, record }">
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
    </div>
    <LogModal @register="registerLogModal" @success="handleSuccess" />
  </PageWrapper>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { columns, searchFormSchema, execStatusMap, inspectResultMap } from './data';
  import {
    Api,
    GetResultList,
    GetTaskList,
    GetResultFieldName,
    GetResultSummaryStatus,
    GetResultSummaryField,
  } from '@/api/business/inspect';
  import { Tag, Tooltip, Divider, notification, Modal, Row, Col } from 'ant-design-vue';
  import { GetSubcategoryList } from '@/api/business/biz';
  import { useRoute } from 'vue-router';
  import { h, nextTick, onMounted, reactive, ref } from 'vue';
  import { splitByLineAndTrim } from '@/utils/util';
  import LogModal from './LogModal.vue';
  import { useModal } from '@/components/Modal';
  import { PageWrapper } from '@/components/Page';
  import { useGo } from '@/hooks/web/usePage';
  import { useCommonStore } from '@/store/modules/common';
  import { PermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import { downloadFileByUrl } from '@/utils/download';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import AnomalyServerPie from './AnomalyServerPie.vue';

  const go = useGo();
  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
    filters: {} as any,
    pie: {
      status: {
        data: [] as any[],
        total: 0,
        title: '设备总数(台)',
      },
      field: {
        data: [] as any[],
        total: 0,
        title: '异常设备(台)',
      },
    },
  });
  let statusPieRef = ref<any>();
  let fieldPieRef = ref<any>();
  const [registerLogModal, { openModal: openLogModal }] = useModal();
  const route = useRoute();
  let taskId: string = route.params.taskId as string;
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
    taskInfo: {} as any,
  });

  defineOptions({ name: 'InspectResult' });
  const [registerTable, { getForm, reload }] = useTable({
    title: '巡检结果',
    api: GetResultList,
    beforeFetch: (params) => {
      params.taskId = taskId;
      params.hostnames = splitByLineAndTrim(params.hostnames);
      return params;
    },
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema(onIdcTypeChange, onCategoryIdsChange, checkFieldsChange),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
      submitOnReset: false,
      showResetButton: false,
      resetFunc() {
        nextTick(async () => {
          reloadBusinessOptions(false);
          resetCheckFields();
          // checkFieldsChange();
        });
        return Promise.resolve();
      },
      compact: true,
    },
    scroll: {
      y: 'max-content',
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      pageSizeOptions: ['10', '50', '100', '500'],
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

  onMounted(() => {
    reloadBusinessOptions();
    resetCheckFields();
    loadPieStatus();
    loadPieField();
  });
  function onIdcTypeChange() {
    reloadBusinessOptions();
  }

  function onCategoryIdsChange() {
    reloadBusinessOptions();
  }

  async function reloadBusinessOptions(includeParmas = true) {
    let selectedBusiness: any[] = [];
    let params = {
      idcType: null,
      categoryIds: null,
    };

    if (includeParmas) {
      const values = await getForm().getFieldsValue();
      params.idcType = values.idcType;
      params.categoryIds = values.categoryIds;
      selectedBusiness = values.business;
    }

    let subcategory = await GetSubcategoryList(params);
    let options: any[] = [];
    subcategory?.items?.forEach((item) => {
      options.push({
        label: item.name,
        value: item.name,
      });
    });
    await getForm().updateSchema([
      {
        field: 'business',
        componentProps: {
          options: options,
        },
      },
    ]);
    // 获取已选择的业务，过滤掉Options中不存在的选项
    let keepBusiness: any[] = [];
    selectedBusiness?.forEach((item) => {
      options.forEach((option) => {
        if (option.value === item) {
          keepBusiness.push(item);
        }
      });
    });
    getForm().setFieldsValue({
      business: keepBusiness,
    });
  }
  function handShowLog(field: Recordable) {
    openLogModal(true, field);
  }
  function handleSuccess() {
    reload();
  }
  function goBack() {
    const commonStore = useCommonStore();
    const inspectHostname = commonStore.getInspectHostname;
    // 返回后需要清除这个值
    commonStore.clearInspectHostname();
    if (inspectHostname) {
      go('/business/inspect/record?hostname=' + inspectHostname);
    } else {
      go('/business/inspect/task');
    }
  }
  onMounted(async () => {
    let data = await GetTaskList({ id: taskId });
    if (data?.items?.length > 0) {
      allData.taskInfo = data.items[0];
    }
  });
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
        value.hostnames = splitByLineAndTrim(value.hostnames);
        value = { ...value, ...data.filters, taskId: taskId };
        nextTick(() => {
          data.exporting = true;
          data.exportButTitle = '导出中';
        });
        await exportInspectResult(value);
      },
    });
  }
  async function exportInspectResult(value: Recordable) {
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
    let res = await GetResultFieldName({ taskId: taskId, type: 1 });
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

  async function loadPieStatus() {
    let res = await GetResultSummaryStatus({ taskId: taskId });
    if (res.items) {
      res.items.forEach((item) => {
        switch (item.name) {
          case '正常设备':
            item.itemStyle = {
              color: '#5fb842',
            };
            break;
          case '异常设备':
            item.itemStyle = {
              color: '#e80919',
            };
            break;
          case '执行失败':
            item.itemStyle = {
              color: '#e062ae',
            };
            break;
          case '执行中':
            item.itemStyle = {
              color: '#fba31d',
            };
            break;
        }
      });
      data.pie.status.data = res.items;
    } else {
      data.pie.status.data = [];
    }
    if (res.total) {
      data.pie.status.total = res.total;
    } else {
      data.pie.status.total = 0;
    }
    nextTick(() => {
      statusPieRef.value.rebuildGraphic();
    });
  }
  async function loadPieField() {
    let res = await GetResultSummaryField({ taskId: taskId });
    if (res.items) {
      data.pie.field.data = res.items;
    } else {
      data.pie.field.data = [];
    }
    if (res.total) {
      data.pie.field.total = res.total;
    } else {
      data.pie.field.total = 0;
    }
    nextTick(() => {
      fieldPieRef.value.rebuildGraphic();
    });
  }
  async function onStatusPieSelectField(name: string) {
    if (name.indexOf(data.pie.status.title) != -1) {
      return;
    }
    await getForm().resetFields();
    let values: any = {};
    switch (name) {
      case '执行中':
        values.status = '2';
        break;
      case '执行失败':
        values.status = '3';
        break;
      case '正常设备':
        values.inspectResult = '1';
        break;
      case '异常设备':
        values.inspectResult = '0';
        break;
    }
    await getForm().setFieldsValue(values);
    reload();
  }
  async function onFieldPieSelectField(name: string) {
    if (name.indexOf(data.pie.field.title) != -1) {
      return;
    }

    await getForm().resetFields();
    let values: any = {
      checkFields: [name],
      checkFieldStatus: '0',
    };
    await getForm().setFieldsValue(values);
    reload();
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
