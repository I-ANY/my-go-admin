<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="handleCreate"> 新增 </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'notifyType'">
          <template v-for="id in record.notifyType" :key="id">
            <Tag style="margin: 1px">{{ notifyTypeMap[id] }}</Tag>
          </template>
        </template>
        <template v-if="column.key == 'users'">
          {{ (record.users || []).map((user) => user.nickName).join(', ') }}
        </template>
        <template v-if="column.key == 'urlType'">
          <span>{{ urlTypeOptions[record.urlType] }}</span>
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
    <WebhookDrawer @register="registerDrawer" @success="handleSuccess" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, TableAction, BasicColumn } from '@/components/Table';
  import { getWebhooks, deleteWebhook } from '@/api/price/price';
  import { message, Tag } from 'ant-design-vue';
  import { useDrawer } from '@/components/Drawer';
  import WebhookDrawer from './webhookDrawer.vue';

  defineOptions({ name: 'PriceWebhook' });

  const notifyTypeMap = {
    1: '未关联业务组',
    2: '业务组更新',
    3: '采购单价未配置',
  };
  const urlTypeOptions = {
    1: '企业微信机器人',
    2: '钉钉机器人',
  };

  const columns: BasicColumn[] = [
    {
      title: '类型',
      dataIndex: 'urlType',
      width: 200,
    },
    {
      title: '通知用户',
      dataIndex: 'users',
      width: 300,
    },
    {
      title: 'webhook',
      dataIndex: 'url',
    },
    {
      title: '通知事件',
      dataIndex: 'notifyType',
      width: 200,
    },
  ];

  const [registerDrawer, { openDrawer }] = useDrawer();
  const [registerTable, { reload }] = useTable({
    title: 'webhook列表',
    api: getWebhooks,
    columns,
    useSearchForm: false,
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
    await deleteWebhook(record.id);
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
