<template>
  <BasicTable @register="registerTable">
    <template #bodyCell="{ column, record }">
      <template v-if="column.key == 'id'">
        <Tooltip
          v-if="hasPermission(PermissionCodeEnum.BUSINESS_INSPECT_TASK_RESULT)"
          placement="top"
          title="查看详情"
        >
          <a @click="onIdClick(record)">#{{ record.id }}</a>
        </Tooltip>
        <span v-else>#{{ record.id }}</span>
      </template>
      <template v-if="column.key == 'operatorType'">
        <Tag
          v-if="operatorTypeMap[record.operatorType]"
          :color="operatorTypeMap[record.operatorType].color || 'default'"
          style="font-weight: bold"
          >{{ operatorTypeMap[record.operatorType].dictLabel }}</Tag
        >
        <span v-else>{{ record.status }}</span>
      </template>
      <template v-if="column.key == 'status'">
        <Tag
          v-if="execStatusMap[record.status]"
          :color="execStatusMap[record.status].color || 'default'"
          style="font-weight: bold"
          >{{ execStatusMap[record.status].dictLabel }}</Tag
        >
        <span v-else>{{ record.status }}</span>
      </template>
      <template v-if="column.key == 'count'">
        <span style="color: green"> {{ record.normalCount }} </span>/
        <span style="color: red">
          {{ record.abnormalCount }}
        </span>
      </template>
      <!-- <template v-if="column.key == 'normalCount'">
        <span style="color: green">
          {{ record.normalCount }}
        </span>
      </template>
      <template v-if="column.key == 'abnormalCount'">
        <span style="color: red">
          {{ record.abnormalCount }}
        </span>
      </template> -->
    </template>
  </BasicTable>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { columns, searchFormSchema, execStatusMap, operatorTypeMap } from './data';
  import { GetTaskList, GetTaskNames } from '@/api/business/inspect';
  import { Tag, Tooltip } from 'ant-design-vue';
  import { nextTick, onMounted } from 'vue';
  import dayjs from 'dayjs';
  import { RangePickPresetsExact } from '@/utils/common';
  import { useGo } from '@/hooks/web/usePage';
  import { usePermission } from '@/hooks/web/usePermission';
  import { PermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import { useCommonStore } from '@/store/modules/common';
  import { DefaultOptionType } from 'ant-design-vue/es/select';

  const { hasPermission } = usePermission();
  const go = useGo();
  defineOptions({ name: 'InspectTask' });
  const [registerTable, { getForm }] = useTable({
    title: '巡检任务',
    api: GetTaskList,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema(onTimePikerOpen),
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 3,
      submitOnReset: false,
      alwaysShowLines: 1,
      resetFunc() {
        nextTick(() => {
          resetReportTime();
          loadTaskNames();
        });
        return Promise.resolve();
      },
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    rowKey: 'id',
    // actionColumn: {
    //   width: 120,
    //   title: '操作',
    //   dataIndex: 'action',
    //   // slots: { customRender: 'action' },
    //   fixed: 'right',
    // },
  });

  onMounted(async () => {
    loadTaskNames();
    await resetReportTime();
  });
  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      startTimeBegin: dayjs(
        dayjs().add(-1, 'day').format('YYYY-MM-DD 00:00:00'),
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
  function onIdClick(record: Recordable) {
    const commonStore = useCommonStore();
    commonStore.clearInspectHostname();
    go('/business/inspect/task/' + record.id + '/result');
  }
  async function loadTaskNames() {
    const res = await GetTaskNames();
    let options: DefaultOptionType[] = [];
    res?.forEach((item) => {
      options.push({
        label: item,
        value: item,
      });
    });
    getForm().updateSchema({
      field: 'taskName',
      componentProps: {
        options: options,
      },
    });
  }
</script>
