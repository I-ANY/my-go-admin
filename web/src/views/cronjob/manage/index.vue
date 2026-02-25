<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" v-auth="'cronjob:job:add'" @click="handleCreate">
          新增任务
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'lsicon:play-outline',
                // color: 'error',
                tooltip: '立即执行',
                auth: 'cronjob:job:exec',
                popConfirm: {
                  title: '是否确认执行',
                  placement: 'left',
                  confirm: handleExec.bind(null, record),
                },
              },
              {
                icon: 'clarity:note-edit-line',
                tooltip: '编辑任务',
                onClick: handleEdit.bind(null, record),
                auth: 'cronjob:job:edit',
              },
              {
                icon: 'ant-design:delete-outlined',
                color: 'error',
                tooltip: '删除任务',
                auth: 'cronjob:job:delete',
                popConfirm: {
                  title: '是否确认删除',
                  placement: 'left',
                  confirm: handleDelete.bind(null, record),
                },
              },
            ]"
          />
        </template>
        <template v-if="column.key == 'status'">
          <Tag :bordered="true" :color="record.status == 1 ? 'green' : 'error'">{{
            record.status == 1 ? '启用' : '禁用'
          }}</Tag></template
        >
        <template v-if="column.key == 'jobType'">
          <Tag :bordered="true" :color="record.jobType == 2 ? 'green' : 'purple'">{{
            record.jobType == 2 ? 'EXEC' : 'FLASK-API'
          }}</Tag></template
        >
        <template v-if="column.key == 'runStatus'">
          <Tag :bordered="true" v-if="record.runStatus == 1" :color="'orange'">运行中</Tag>
          <Tag :bordered="true" v-else-if="record.runStatus == 2" :color="'green'">执行成功</Tag>
          <Tag :bordered="true" v-else-if="record.runStatus == 3" :color="'error'">执行失败</Tag>
          <span v-else></span>
        </template>
        <template v-if="column.key == 'jobName'">
          <Tooltip
            placement="top"
            title="查看执行记录"
            v-if="hasPermission(['cronjob:job:exec-record'])"
          >
            <a @click="gotoRecord(record.id)">{{ record.jobName }}</a>
          </Tooltip>
          <span v-else>{{ record.jobName }}</span>
        </template>
      </template>
    </BasicTable>
    <JobMadal @register="registerMoal" @success="handleSuccess" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { getJobList, deleteJob, execJob } from '@/api/cronjob/job';
  import { Tag, message, Tooltip } from 'ant-design-vue';
  import { useModal } from '@/components/Modal';
  import JobMadal from './JobModal.vue';

  import { columns, searchFormSchema } from './data';
  import { usePermission } from '@/hooks/web/usePermission';
  import { useGo } from '@/hooks/web/usePage';

  defineOptions({ name: 'JobManagement' });
  const { hasPermission } = usePermission();
  const go = useGo();

  const [registerMoal, { openModal }] = useModal();
  const [registerTable, { reload }] = useTable({
    title: '任务列表',
    api: getJobList,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
      // autoAdvancedLine: 1,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    actionColumn: {
      width: 120,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: 'right',
      ifShow: function (): boolean {
        return (
          hasPermission('cronjob:job:exec') ||
          hasPermission('cronjob:job:edit') ||
          hasPermission('cronjob:job:delete')
        );
      },
    },
  });

  function handleCreate() {
    openModal(true, {
      isUpdate: false,
    });
  }

  function handleEdit(record: Recordable) {
    openModal(true, {
      record,
      isUpdate: true,
    });
  }

  async function handleDelete(record: Recordable) {
    await deleteJob(record.id);
    message.success('删除成功');
    reload();
  }
  async function handleExec(record: Recordable) {
    await execJob(record.id);
    message.success('触发执行成功');
    reload();
  }
  function handleSuccess(isAdd: boolean) {
    reload();
    let msg = '新增成功';
    if (!isAdd) {
      msg = '编辑成功';
    }
    message.success(msg);
  }
  function gotoRecord(id: number | string) {
    go('/cronjob/job/' + id + '/exec-record', true);
  }
</script>
