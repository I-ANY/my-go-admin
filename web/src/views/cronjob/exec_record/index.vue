<template>
  <PageWrapper :title="data.title" contentFullHeight contentBackground @back="goBack">
    <BasicTable @register="registerTable">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'icon-park-outline:view-list',
                // color: 'error',
                tooltip: '查看日志',
                auth: 'cronjob:log:view',
                onClick: handleViewLog.bind(null, record),
              },
            ]"
          />
        </template>
        <template v-if="column.key == 'runStatus'">
          <Tag :bordered="true" v-if="record.runStatus == 1" :color="'orange'">运行中</Tag>
          <Tag :bordered="true" v-else-if="record.runStatus == 2" :color="'green'">执行成功</Tag>
          <Tag :bordered="true" v-else-if="record.runStatus == 3" :color="'error'">执行失败</Tag>
          <span v-else></span>
        </template>
        <template v-if="column.key == 'triggerType'">
          <Tag :bordered="true" v-if="record.triggerType == 1" :color="'blue'">定时触发</Tag>
          <Tag :bordered="true" v-else-if="record.triggerType == 2" :color="'purple'">手动触发</Tag>
          <span v-else></span>
        </template>
      </template>
    </BasicTable>
    <ExecRecordLogModal @register="registerModal" @success="handleSuccess" />
  </PageWrapper>
</template>
<script lang="ts" setup>
  import { useRoute } from 'vue-router';
  import { PageWrapper } from '@/components/Page';
  import { useGo } from '@/hooks/web/usePage';
  import { onMounted, reactive } from 'vue';
  import { useTable, BasicTable, TableAction } from '@/components/Table';
  import { usePermission } from '@/hooks/web/usePermission';
  import { columns, searchFormSchema } from './data';
  import { getJob, getJobExecRecordListApi } from '@/api/cronjob/job';
  import { Tag } from 'ant-design-vue';
  import ExecRecordLogModal from './ExecRecordLogModal.vue';
  import { useModal } from '@/components/Modal';
  import { RangePickPresetsExact } from '@/utils/common';

  defineOptions({ name: 'JobExecRecord' });
  const { hasPermission } = usePermission();
  const [registerModal, { openModal }] = useModal();

  const route = useRoute();
  const go = useGo();
  const data = reactive({
    title: '',
    jobId: 0,
  });
  data.jobId = route.params.id as any;
  const [registerTable, { reload, getForm }] = useTable({
    title: '执行记录',
    api: getJobExecRecordListApi,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema(onTimePikerOpen),
      autoSubmitOnEnter: true,
      autoAdvancedLine: 1,
      actionColOptions: {
        span: 4,
      },
    },
    size: 'small',
    beforeFetch: (params) => {
      params.jobId = data.jobId;
      return params;
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: true,
    actionColumn: {
      width: 80,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: 'right',
      ifShow: function (): boolean {
        return hasPermission('cronjob:log:view');
      },
    },
  });

  onMounted(async () => {
    const job = await getJob(data.jobId);
    data.title = job?.jobName;
  });

  // 页面左侧点击返回链接时的操作
  function goBack() {
    go('/cronjob/job/list');
  }
  function handleSuccess() {
    reload();
  }
  function handleViewLog(record) {
    openModal(true, record);
  }
  function onTimePikerOpen() {
    getForm().updateSchema([
      {
        field: '[startTimeBegin, startTimeEnd]',
        componentProps: {
          presets: RangePickPresetsExact(),
        },
      },
    ]);
  }
</script>
