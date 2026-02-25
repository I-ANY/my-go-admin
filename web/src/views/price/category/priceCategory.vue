<template>
  <div class="h-full flex p-4">
    <div class="flex flex-col pr-4 w-1/3">
      <BasicTable @register="registerBizTable">
        <template #bodyCell="{ column, record }">
          <template v-if="column.key == 'hide'">
            {{ record.hide ? '隐藏' : '显示' }}
          </template>
          <template v-if="column.key == 'enableNotify'">
            <Switch
              v-model:checked="record.enableNotify"
              @change="(checked) => handleNotifyStatus(record, checked)"
            />
          </template>
          <template v-if="column.key === 'action'">
            <TableAction
              :actions="[
                {
                  icon: 'clarity:note-edit-line',
                  onClick: (event) => {
                    event.stopPropagation();
                    handleEditBiz(record);
                  },
                  tooltip: '编辑',
                },
              ]"
            />
          </template>
        </template>
      </BasicTable>
      <BizDrawer @register="registerDrawerBiz" @success="handleSuccess" />
    </div>
    <div class="flex-1 flex flex-col h-full">
      <div class="h-1/3 mb-4">
        <div class="flex-1 h-0">
          <BasicTable @register="registerTable">
            <template #toolbar>
              <Excel :excelType="excelType" @success-upload="uploadSuccess" />
              <a-button type="primary" @click="handleCreate"> 新增 </a-button>
            </template>
            <template #bodyCell="{ column, record }">
              <template v-if="column.key == 'bizs'">
                {{ (record.bizs || []).map((biz) => biz.name).join(', ') }}
              </template>
              <template v-if="column.key == 'outName'">
                {{ record.outName ? record.outName : '/' }}
              </template>
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
          <CategoryDrawer @register="registerDrawer" @success="handleSuccess" />
        </div>
      </div>
    </div>
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction, BasicColumn, FormSchema } from '@/components/Table';
  import { getCategory, deleteCategory, getNotRelBiz, updateNotify } from '@/api/price/price';
  import { message, Switch } from 'ant-design-vue';
  import { useDrawer } from '@/components/Drawer';
  import CategoryDrawer from './categoryDrawer.vue';
  import BizDrawer from './bizDrawer.vue';
  import Excel from '../excel.vue';

  defineOptions({ name: 'PriceCategory' });

  const excelType = 'category';

  const bizColumns: BasicColumn[] = [
    {
      title: '业务名',
      dataIndex: 'name',
    },
    {
      title: '状态',
      dataIndex: 'hide',
      width: 80,
    },
    {
      title: '通知',
      dataIndex: 'enableNotify',
      width: 100,
    },
  ];
  const [registerDrawerBiz, { openDrawer: openDrawerBiz }] = useDrawer();
  const [registerBizTable, { reload: reloadBiz }] = useTable({
    title: '未分配业务列表',
    api: getNotRelBiz,
    columns: bizColumns,
    useSearchForm: false,
    bordered: true,
    showIndexColumn: false,
    showTableSetting: true,
    pagination: false,
    actionColumn: {
      width: 80,
      title: '操作',
      dataIndex: 'action',
      fixed: 'right',
    },
  });

  const columns: BasicColumn[] = [
    {
      title: '业务组',
      dataIndex: 'name',
      width: 200,
    },
    {
      title: '对外业务名',
      dataIndex: 'outName',
      width: 150,
    },
    {
      title: '业务详情',
      dataIndex: 'bizs',
    },
    {
      title: '备注',
      dataIndex: 'describe',
      width: 100,
    },
  ];
  const searchFormSchema: FormSchema[] = [
    {
      field: 'name',
      label: '业务组',
      component: 'Input',
      colProps: { span: 6 },
    },
    {
      field: 'outName',
      label: '对外业务名',
      component: 'Input',
      colProps: { span: 6 },
    },
  ];

  const [registerDrawer, { openDrawer }] = useDrawer();
  const [registerTable, { reload }] = useTable({
    title: '业务详情列表',
    api: getCategory,
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
      fixed: 'right',
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

  function handleEditBiz(record: Recordable) {
    openDrawerBiz(true, {
      record,
      isUpdate: true,
    });
  }

  async function handleDelete(record: Recordable) {
    await deleteCategory(record.id);
    message.success('删除成功');
    reload();
    reloadBiz();
  }

  function handleSuccess(isAdd: boolean) {
    reload();
    reloadBiz();
    let msg = '新增成功';
    if (!isAdd) {
      msg = '编辑成功';
    }
    message.success(msg);
  }

  function uploadSuccess(msg: string) {
    reload();
    reloadBiz();
    message.success(msg, 10);
  }

  function handleNotifyStatus(record, checked) {
    updateNotify({ id: record.id, enableNotify: checked }).finally(() => {
      reloadBiz();
    });
  }
</script>
