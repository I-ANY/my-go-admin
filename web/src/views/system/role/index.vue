<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" v-auth="'system:role:add'" @click="handleCreate">
          新增角色
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'clarity:note-edit-line',
                onClick: handleEdit.bind(null, record),
                auth: 'system:role:edit',
                disabled: ['ADMIN'].includes(record.identify),
                tooltip: ['ADMIN'].includes(record.identify) ? '不可编辑' : '编辑角色',
              },
              {
                icon: 'clarity:lock-line',
                // label: '资源权限配置',
                onClick: handleOpenResourcePermissionModal.bind(null, record),
                auth: SystemPermissionCodeEnum.ROLE_RESOURCE_PERMISSION_CONFIG,
                tooltip: ['ADMIN'].includes(record.identify) ? '不可配置' : '资源权限配置',
                disabled: ['ADMIN'].includes(record.identify),
              },
              // {
              //   // label: '业务权限配置',
              //   icon: 'clarity:lock-line',
              //   onClick: handleEditBusinessPermission.bind(null, record),
              //   auth: SystemPermissionCodeEnum.ROLE_BUSINESS_PERMISSION_CONFIG,
              //   tooltip: ['ADMIN'].includes(record.identify) ? '不可配置' : '业务权限配置',
              //   disabled: ['ADMIN'].includes(record.identify),
              // },
              {
                icon: 'ant-design:delete-outlined',
                color: 'error',
                auth: 'system:role:delete',
                disabled: ['ADMIN', 'BASE_ROLE'].includes(record.identify),
                tooltip: ['ADMIN', 'BASE_ROLE'].includes(record.identify) ? '不可删除' : '删除角色',
                popConfirm: {
                  title: '是否确认删除',
                  placement: 'left',
                  confirm: handleDelete.bind(null, record),
                },
              },
            ]"
            :dropDownActions="getDropDownActions(record)"
          />
        </template>
        <template v-if="column.key == 'status'">
          <Tag :color="record.status == 1 ? 'green' : 'error'">{{
            record.status == 1 ? '启用' : '禁用'
          }}</Tag></template
        >
      </template>
    </BasicTable>
    <RoleDrawer @register="registerDrawer" @success="handleSuccess" />
    <ResourcePermission
      @register="registerResourcePermission"
      @success="handleSuccess"
      @reload="handleReload"
    />
    <BusinessPermissionDrawer
      @register="registerBusinessPermissionDrawer"
      @success="handleSuccess"
      @reload="handleReload"
    />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction, ActionItem } from '@/components/Table';
  import { getRoleListByPage, deleteRole } from '@/api/demo/system';
  import { Tag, message } from 'ant-design-vue';
  import { useDrawer } from '@/components/Drawer';
  import RoleDrawer from './RoleDrawer.vue';

  import { columns, searchFormSchema } from './role.data';
  import { usePermission } from '@/hooks/web/usePermission';
  import ResourcePermission from './ResourcePermission.vue';
  import { SystemPermissionCodeEnum } from '@/enums/permissionCodeEnum';
  import BusinessPermissionDrawer from './BusinessPermissionDrawer.vue';

  const { hasPermission } = usePermission();

  defineOptions({ name: 'RoleManagement' });

  const [registerDrawer, { openDrawer }] = useDrawer();
  const [registerResourcePermission, { openDrawer: openResourcePermissionDrawer }] = useDrawer();
  const [registerBusinessPermissionDrawer, { openDrawer: _openBusinessPermissionDrawer }] =
    useDrawer();
  const [registerTable, { reload }] = useTable({
    title: '角色列表',
    api: getRoleListByPage,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: true,
    pagination: {
      // pageSizeOptions: ['1'],
    },
    actionColumn: {
      width: 160,
      title: '操作',
      dataIndex: 'action',
      // slots: { customRender: 'action' },
      fixed: undefined,
      ifShow: function (): boolean {
        return hasPermission('system:role:edit') || hasPermission('system:role:delete');
      },
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
    await deleteRole(record.id);
    message.success('删除成功');
    reload();
  }

  function handleSuccess() {
    message.success('操作成功');
    reload();
  }
  function handleReload() {
    reload();
  }
  function handleOpenResourcePermissionModal(record: Recordable) {
    openResourcePermissionDrawer(true, {
      record,
    });
  }
  function getDropDownActions(_record: Recordable): ActionItem[] {
    let actions: ActionItem[] = [];
    // if (hasPermission(SystemPermissionCodeEnum.ROLE_RESOURCE_PERMISSION_CONFIG)) {
    //   actions.push({
    //     // icon: 'clarity:lock-line',
    //     label: '资源权限配置',
    //     onClick: handleOpenResourcePermissionModal.bind(null, record),
    //     auth: SystemPermissionCodeEnum.ROLE_RESOURCE_PERMISSION_CONFIG,
    //     tooltip: ['ADMIN'].includes(record.identify) ? '不可配置' : '资源权限配置',
    //     disabled: ['ADMIN'].includes(record.identify),
    //   });
    // }
    return actions;
  }
  // function handleEditBusinessPermission(record: Recordable) {
  //   openBusinessPermissionDrawer(true, {
  //     record,
  //   });
  // }
</script>
