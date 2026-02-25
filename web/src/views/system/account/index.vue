<template>
  <PageWrapper dense contentFullHeight contentClass="flex">
    <DeptTree class="w-1/6 xl:w-1/6" @select="handleSelect" />
    <BasicTable @register="registerTable" class="w-3/4 xl:w-4/5" :searchInfo="searchInfo">
      <template #toolbar>
        <a-button type="primary" v-auth="'system:user:add'" @click="handleCreate"
          >新增用户</a-button
        >
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              // {
              //   icon: 'clarity:info-standard-line',
              //   tooltip: '查看用户详情',
              //   onClick: handleView.bind(null, record),
              // },
              {
                icon: 'clarity:note-edit-line',
                tooltip: '编辑用户',
                onClick: handleEdit.bind(null, record),
                auth: 'system:user:edit',
              },
              {
                icon: 'ant-design:delete-outlined',
                color: 'error',
                auth: 'system:user:delete',
                tooltip: '删除用户',
                popConfirm: {
                  title: '是否确认删除',
                  placement: 'left',
                  confirm: handleDelete.bind(null, record),
                },
              },
            ]"
          />
        </template>
        <template v-if="column.key === 'status'">
          <Tag :color="record.status == 1 ? 'green' : 'red'">{{
            record.status == 1 ? '启用' : '禁用'
          }}</Tag>
        </template>
        <template v-if="column.key === 'dept'">
          <Tag color="green" v-if="record.dept">{{ record.dept?.name }}</Tag>
        </template>
        <template v-if="column.key === 'roles'">
          <Tag color="green" v-for="role in record.roles" :key="role.id">{{ role.name }}</Tag>
        </template>
        <template v-if="column.key === 'source'">
          <Tag
            v-if="userFromMap[record.source]"
            :color="userFromMap[record.source].color || 'default'"
            >{{ userFromMap[record.source].dictLabel }}</Tag
          >
        </template>
      </template>
    </BasicTable>
    <AccountModal @register="registerModal" @success="handleSuccess" />
  </PageWrapper>
</template>
<script lang="ts" setup>
  import { reactive } from 'vue';

  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { deleteUser, getAccountList } from '@/api/demo/system';
  import { PageWrapper } from '@/components/Page';
  import DeptTree from './DeptTree.vue';

  import { useModal } from '@/components/Modal';
  import AccountModal from './AccountModal.vue';

  import { columns, searchFormSchema, userFromMap } from './account.data';
  // import { useGo } from '@/hooks/web/usePage';
  import { Tag, message } from 'ant-design-vue';
  import { usePermission } from '@/hooks/web/usePermission';

  const { hasPermission } = usePermission();
  defineOptions({ name: 'AccountManagement' });

  // const go = useGo();
  const [registerModal, { openModal }] = useModal();
  const searchInfo = reactive<Recordable>({});
  const [registerTable, { reload }] = useTable({
    title: '用户列表',
    api: getAccountList,
    rowKey: 'id',
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
      autoAdvancedLine: 1,
    },
    // pagination: {
    //   pageSizeOptions: ['1'],
    // },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    handleSearchInfoFn(info) {
      return info;
    },
    actionColumn: {
      width: 120,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      ifShow: function (): boolean {
        return hasPermission('system:user:edit') || hasPermission('system:user:delete');
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
    await deleteUser(record.id);
    message.success('删除成功');
    reload();
  }

  function handleSuccess() {
    // if (isUpdate) {
    //   // 演示不刷新表格直接更新内部数据。
    //   // 注意：updateTableDataRecord要求表格的rowKey属性为string并且存在于每一行的record的keys中
    //   const result = updateTableDataRecord(values.id, values);
    //   console.log(result);
    // } else {
    reload();
    // }
  }

  function handleSelect(deptId = '') {
    searchInfo.deptId = deptId;
    reload();
  }

  // function handleView(record: Recordable) {
  //   go('/system/account_detail/' + record.id);
  // }
</script>
