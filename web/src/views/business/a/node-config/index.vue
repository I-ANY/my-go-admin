<template>
  <div class="node-config-container">
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="success"
          @click="handleImport"
          v-if="hasPermission('businessA:node-config:import')"
        >
          批量导入
        </a-button>
      </template>

      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'clarity:note-edit-line',
                tooltip: '编辑节点信息',
                onClick: handleEdit.bind(null, record),
                ifShow: () => hasPermission('businessA:node-config:edit'),
              },
            ]"
          />
        </template>
        <template v-if="column.key === 'billingType'">
          <Tag :color="getBillingTypeColor(record.billingType)">
            {{ record.billingType }}
          </Tag>
        </template>
      </template>
    </BasicTable>

    <!-- 新增/编辑模态框 -->
    <NodeConfigModal @register="registerModal" @success="handleSuccess" />

    <!-- 导入模态框 -->
    <ImportModal @register="registerImportModal" @success="handleSuccess" />
  </div>
</template>

<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { useModal } from '@/components/Modal';
  import { Tag } from 'ant-design-vue';
  import { columns, searchFormSchema } from './data';
  import NodeConfigModal from './NodeConfigModal.vue';
  import ImportModal from './ImportModal.vue';
  import { getNodeConfigs } from '@/api/business/a';
  import { usePermission } from '@/hooks/web/usePermission';

  defineOptions({ name: 'NodeConfigList' });

  const { hasPermission } = usePermission();

  // 模态框
  const [registerModal, { openModal }] = useModal();
  const [registerImportModal, { openModal: openImportModal }] = useModal();

  // 表格配置
  const [registerTable, { reload }] = useTable({
    title: '节点配置列表',
    api: getNodeConfigs,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    actionColumn: {
      width: 120,
      title: '操作',
      dataIndex: 'action',
      fixed: 'right',
    },
  });

  // 计费类型颜色映射
  const getBillingTypeColor = (type: string) => {
    const colorMap = {
      日95: 'blue',
      月95: 'green',
      买断: 'orange',
    };
    return colorMap[type] || 'default';
  };

  // 编辑
  function handleEdit(record: Recordable) {
    openModal(true, {
      record,
      isUpdate: true,
    });
  }

  // 批量导入
  function handleImport() {
    openImportModal(true);
  }

  // 操作成功回调
  function handleSuccess() {
    reload();
  }
</script>

<style lang="less" scoped>
  .node-config-container {
    padding: 16px;
  }
</style>
