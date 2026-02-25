<template>
  <div>
    <BasicTable @register="registerTable" @fetch-success="onFetchSuccess">
      <template #toolbar>
        <a-button type="primary" v-auth="'system:dept:add'" @click="handleCreate">
          新增部门
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'material-symbols:add',
                onClick: handleCreate.bind(null, record),
                auth: 'system:dept:add',
              },
              {
                icon: 'clarity:note-edit-line',
                onClick: handleEdit.bind(null, record),
                auth: 'system:dept:update',
              },
              {
                icon: 'ant-design:delete-outlined',
                color: 'error',
                popConfirm: {
                  title: '是否确认删除',
                  placement: 'left',
                  confirm: handleDelete.bind(null, record),
                },
                auth: 'system:dept:delete',
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <DeptModal @register="registerModal" @success="handleSuccess" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { getDeptList, deleteDept } from '@/api/demo/system';

  import { useModal } from '@/components/Modal';
  import DeptModal from './DeptModal.vue';

  import { columns } from './dept.data';
  import { nextTick } from 'vue';
  import { message } from 'ant-design-vue';

  defineOptions({ name: 'DeptManagement' });

  const [registerModal, { openModal }] = useModal();
  const [registerTable, { reload, expandAll }] = useTable({
    title: '部门列表',
    api: getDeptList,
    columns,
    // formConfig: {
    //   labelWidth: 120,
    //   schemas: searchFormSchema,
    // },
    pagination: false,
    striped: false,
    // useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    isTreeTable: true,
    showIndexColumn: false,
    canResize: false,
    actionColumn: {
      width: 120,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: undefined,
    },
  });

  function handleCreate(record: Recordable) {
    openModal(true, {
      isUpdate: false,
      parentDept: record?.id,
    });
  }

  function handleEdit(record: Recordable) {
    openModal(true, {
      record,
      isUpdate: true,
    });
  }

  async function handleDelete(record: Recordable) {
    await deleteDept(record.id);
    message.success('删除成功');
    reload();
  }

  function handleSuccess() {
    reload();
  }
  function onFetchSuccess() {
    // 演示默认展开所有表项
    nextTick(expandAll);
  }
</script>
