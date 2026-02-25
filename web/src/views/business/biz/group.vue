<template>
  <div class="group-page">
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="handleCreate"> 新增业务组</a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'clarity:note-edit-line',
                onClick: handleEdit.bind(null, record),
                tooltip: '编辑',
              },
              {
                icon: 'ant-design:delete-outlined',
                color: 'error',
                onClick: handleDelete.bind(null, record),
                tooltip: '删除',
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <GroupModal @register="registerModal" @success="handleSuccess" />
  </div>
</template>

<script setup lang="ts">
  import { BasicTable, TableAction, useTable } from '@/components/Table';
  import { GetBusinessGroupList, DeleteBusinessGroup } from '@/api/business/biz';
  import { groupSchemas, groupColumns } from './data';
  import { useModal } from '@/components/Modal';
  import GroupModal from './groupModal.vue';
  import { message } from 'ant-design-vue';

  const [registerTable, { reload }] = useTable({
    title: '业务组列表',
    api: GetBusinessGroupList,
    columns: groupColumns,
    formConfig: {
      labelWidth: 100,
      schemas: groupSchemas,
      autoSubmitOnEnter: true,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: true,
    actionColumn: {
      width: 120,
      title: '操作',
      dataIndex: 'action',
      fixed: 'right',
    },
  });

  const [registerModal, { openModal }] = useModal();

  function handleCreate() {
    openModal(true, {
      isUpdate: false,
      record: {},
    });
  }

  function handleEdit(record: Recordable) {
    openModal(true, {
      isUpdate: true,
      record,
    });
  }

  async function handleDelete(record: Recordable) {
    try {
      await DeleteBusinessGroup(record.id);
      message.success('删除成功');
      await reload();
    } catch (error) {
      message.error('删除失败');
    }
  }

  function handleSuccess() {
    reload();
  }
</script>

<style scoped lang="less">
  .group-page {
    padding: 16px;
  }
</style>
