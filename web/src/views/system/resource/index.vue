<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="SystemPermissionCodeEnum.RESOURCE_ADD"
          @click="handleCreate"
        >
          新增资源
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                label: '查看资源',
                // icon: 'ant-design:eye-filled',
                icon: 'iconamoon:eye-light',
                tooltip: '查看资源数据',
                auth: SystemPermissionCodeEnum.RESOURCE_VIEW,
                onClick: handleView.bind(null, record),
              },
              {
                label: '编辑',
                icon: 'clarity:note-edit-line',
                tooltip: '编辑资源',
                auth: SystemPermissionCodeEnum.RESOURCE_EDIT,
                onClick: handleEdit.bind(null, record),
              },
              {
                label: '删除',
                icon: 'ant-design:delete',
                tooltip: '删除资源',
                color: 'error',
                auth: SystemPermissionCodeEnum.RESOURCE_DELETE,
                popConfirm: {
                  title: '确定删除该资源吗？',
                  confirm: handleDelete.bind(null, record),
                },
              },
            ]"
            :dropDownActions="getActions(record)"
          />
        </template>
        <template v-else-if="column.dataIndex === 'permissionTypes'">
          <div class="warp-tag">
            <Tag color="blue" v-for="(p, i) in record.permissionTypes" :key="i">{{ p.name }}</Tag>
          </div>
        </template>
        <template v-else-if="column.dataIndex === 'filter'">
          <Tooltip v-if="record.filter" placement="topLeft"
            >{{ record.filter }}
            <template #title>
              <span>{{ record.filter }}</span>
            </template></Tooltip
          >
        </template>
      </template>
    </BasicTable>
    <ResourceDrawer
      @register="registerResourceEditDrawer"
      @success="handleSuccess"
      @reload="handleReload"
    />
    <ResourceViewModal
      @register="registerResourceViewModal"
      @success="handleSuccess"
      @reload="handleReload"
    />
  </div>
</template>
<script lang="ts" setup>
  import { message, Tag, Tooltip } from 'ant-design-vue';
  import { BasicTable, useTable, TableAction, ActionItem } from '@/components/Table';
  import { getResourceList, deleteResource } from '@/api/sys/resource';
  import { getResourceColumns, getSearchResourceColumnsFormSchema } from './data';
  import { usePermission } from '@/hooks/web/usePermission';
  import { SystemPermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import ResourceDrawer from './ResourceDrawer.vue';
  import { useDrawer } from '@/components/Drawer';
  import { useModal } from '@/components/Modal';
  import ResourceViewModal from './ResourceViewModal.vue';

  const { hasPermission } = usePermission();

  defineOptions({ name: 'ResourceManagement' });
  const [registerResourceEditDrawer, { openDrawer: openResourceEditDrawer }] = useDrawer();
  const [registerResourceViewModal, { openModal: openResourceViewModal }] = useModal();
  const [registerTable, { reload }] = useTable({
    title: '资源列表',
    api: getResourceList,
    columns: getResourceColumns(),
    formConfig: {
      labelWidth: 120,
      schemas: getSearchResourceColumnsFormSchema(),
      autoSubmitOnEnter: true,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      // pageSizeOptions: ['1'],
    },
    actionColumn: {
      width: 250,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: 'right',
      ifShow: function (): boolean {
        return (
          hasPermission(SystemPermissionCodeEnum.RESOURCE_EDIT) ||
          hasPermission(SystemPermissionCodeEnum.RESOURCE_VIEW)
        );
      },
    },
  });

  function handleCreate() {
    openResourceEditDrawer(true, {
      isUpdate: false,
    });
  }

  function handleEdit(record: Recordable) {
    openResourceEditDrawer(true, {
      isUpdate: true,
      record,
    });
  }

  function handleSuccess() {
    message.success('操作成功');
    reload();
  }
  function handleView(record: Recordable) {
    openResourceViewModal(true, {
      record,
    });
  }
  function getActions(_record: Recordable): ActionItem[] {
    let actions: ActionItem[] = [];
    return actions;
  }
  function handleReload() {
    reload();
  }
  // function showActionColumn(_record: Recordable): boolean {
  //   return (
  //     hasPermission(SystemPermissionCodeEnum.RESOURCE_EDIT) ||
  //     hasPermission(SystemPermissionCodeEnum.RESOURCE_VIEW)
  //   );
  // }
  async function handleDelete(record: Recordable) {
    await deleteResource(record.id);
    message.success('删除成功');
    reload();
  }
</script>
<style lang="less" scoped>
  .warp-tag {
    display: flex;
    flex-wrap: wrap;
    gap: 3px;
    justify-content: center; /* 水平居中 */
  }
</style>
