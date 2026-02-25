<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="handleCreate"> 新增 </a-button>
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
                tooltip: '删除',
                popConfirm: {
                  title: '确认删除?',
                  placement: 'left',
                  confirm: handleDelete.bind(null, record),
                },
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <ModeDrawer @register="registerDrawer" @success="handleSuccess" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction, BasicColumn, FormSchema } from '@/components/Table';
  import { getMode, deleteMode } from '@/api/price/price';
  import { message } from 'ant-design-vue';
  import { useDrawer } from '@/components/Drawer';
  import ModeDrawer from './ModeDrawer.vue';

  defineOptions({ name: 'PriceMode' });

  const columns: BasicColumn[] = [
    {
      title: '计费方式',
      dataIndex: 'name',
      width: 300,
    },
    {
      title: '备注',
      dataIndex: 'describe',
    },
  ];
  const searchFormSchema: FormSchema[] = [
    {
      field: 'name',
      label: '计费方式',
      component: 'Input',
      colProps: { span: 4 },
    },
  ];

  const [registerDrawer, { openDrawer }] = useDrawer();
  const [registerTable, { reload }] = useTable({
    title: '计费方式列表',
    api: getMode,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
    },
    useSearchForm: true,
    bordered: true,
    showIndexColumn: true,
    showTableSetting: true,
    actionColumn: {
      width: 80,
      title: '操作',
      dataIndex: 'action',
      fixed: undefined,
    },
  });

  function handleCreate() {
    openDrawer(true, {
      isUpdate: false,
    });
  }

  function handleEdit(record: Recordable) {
    openDrawer(true, {
      record,
      isUpdate: true,
    });
  }

  async function handleDelete(record: Recordable) {
    await deleteMode(record.id);
    message.success('删除成功');
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
</script>
