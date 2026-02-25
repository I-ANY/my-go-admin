<template>
  <div>
    <BasicTable @register="registerTable" @expanded-rows-change="onExpandedRowsChange">
      <template #toolbar>
        <Button @click="handleExpand">{{ data.expand ? '收起' : '展开' }}所有</Button>
        <Button type="primary" @click="handleCreate" v-auth="'system:menu:add'"> 新增菜单 </Button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'material-symbols:add',
                onClick: handleCreate.bind(null, record),
                auth: 'system:menu:add',
                tooltip: '添加子菜单',
              },
              {
                icon: 'clarity:note-edit-line',
                onClick: handleEdit.bind(null, record),
                auth: 'system:menu:edit',
                tooltip: '编辑菜单',
              },
              {
                icon: 'ant-design:delete-outlined',
                color: 'error',
                auth: 'system:menu:delete',
                tooltip: '删除菜单',
                popConfirm: {
                  title: '是否确认删除',
                  placement: 'left',
                  confirm: handleDelete.bind(null, record),
                },
              },
            ]"
          />
        </template>
        <template v-if="column.dataIndex === 'type'">
          <Tag v-if="record.type == 0" color="#2db7f5">目录</Tag>
          <Tag v-else-if="record.type == 1" color="#0960bd">菜单</Tag>
          <Tag v-else-if="record.type == 2" color="default">按钮</Tag>
        </template>
      </template>
    </BasicTable>
    <MenuDrawer @register="registerDrawer" @success="handleSuccess" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';

  import { useDrawer } from '@/components/Drawer';
  import MenuDrawer from './MenuDrawer.vue';
  import { Tag, message, Button } from 'ant-design-vue';
  import { columns } from './menu.data';
  import { getAllMenus } from '@/api/demo/system';
  import { deleteMenu } from '@/api/sys/menu';
  import { nextTick, reactive } from 'vue';
  import { usePermission } from '@/hooks/web/usePermission';

  const { hasPermission } = usePermission();

  defineOptions({ name: 'MenuManagement' });
  const data = reactive({
    expand: false, //是否展开
    expandedRows: [] as any[],
  });
  const [registerDrawer, { openDrawer }] = useDrawer();
  const [registerTable, { reload, expandAll, collapseAll, expandRows }] = useTable({
    title: '菜单列表',
    api: getAllMenus,
    columns,
    // formConfig: {
    //   labelWidth: 120,
    //   schemas: searchFormSchema,
    // },
    isTreeTable: true,
    pagination: false,
    striped: false,
    // useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    canResize: false,
    rowKey: 'id',
    actionColumn: {
      width: 120,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: undefined,
      ifShow: function (): boolean {
        return (
          hasPermission('system:menu:add') ||
          hasPermission('system:menu:edit') ||
          hasPermission('system:menu:delete')
        );
      },
    },
    afterFetch: () => {
      nextTick(() => {
        if (data.expandedRows.length > 0) {
          expandRows(data.expandedRows);
        }
      });
    },
  });

  function handleCreate(record: Recordable) {
    if (record) {
      openDrawer(true, {
        isUpdate: false,
        parentMenu: record.id,
      });
    } else {
      openDrawer(true, {
        isUpdate: false,
      });
    }
  }

  function handleEdit(record: Recordable) {
    openDrawer(true, {
      record,
      isUpdate: true,
    });
  }

  async function handleDelete(record: Recordable) {
    await deleteMenu(record.id);
    reload();
    data.expand = false;
    message.success('删除成功');
  }

  function handleSuccess() {
    reload();
    data.expand = false;
  }

  function handleExpand() {
    if (data.expand) {
      nextTick(collapseAll);
    } else {
      nextTick(expandAll);
    }
    data.expand = !data.expand;
  }
  function onExpandedRowsChange(expandedRows: any) {
    data.expandedRows = expandedRows;
  }
</script>
