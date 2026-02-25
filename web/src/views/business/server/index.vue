<template>
  <div>
    <BasicTable @register="registerTable">
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'notes'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '300px' }">
            <template #title>{{ record.notes }}</template>
            <span>{{ record.notes }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key === 'cactiNotes'">
          <Tooltip placement="topLeft" :overlayInnerStyle="{ width: '300px' }">
            <template #title>{{ record.cactiNotes }}</template>
            <span>{{ record.cactiNotes }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                label: '业务追溯',
                onClick: handleShowHistory.bind(null, record),
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <HistoryModal @register="registerHistoryModal" />
  </div>
</template>

<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { Tooltip } from 'ant-design-vue';
  import { useModal } from '@/components/Modal';
  import { columns, searchFormSchema } from './data';
  import { GetAuthedBiz, GetAuthedServer } from '@/api/business/biz';
  import { splitByLineAndTrim } from '@/utils/util';
  import HistoryModal from './HistoryModal.vue';
  import { onMounted } from 'vue';

  defineOptions({ name: 'ServerList' });

  const [registerHistoryModal, { openModal: openHistoryModal }] = useModal();

  function handleShowHistory(record: Recordable) {
    openHistoryModal(true, {
      hostname: record.hostname,
      sn: record.sn,
      location: record.location,
      owner: record.owner,
    });
  }

  let subCategoryOptions: any[] = [];

  const [registerTable, { getForm }] = useTable({
    title: '服务器列表',
    api: GetAuthedServer,
    columns,
    formConfig: {
      labelWidth: 10,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
    },
    useSearchForm: true,
    showSelectionBar: false,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      pageSizeOptions: ['10', '30', '50', '100', '200', '500'],
    },
    clickToRowSelect: false,
    scroll: { x: 3000 },
    actionColumn: {
      width: 100,
      title: '操作',
      dataIndex: 'action',
      fixed: 'right',
    },
    beforeFetch: async (params) => {
      // 处理 hostnames 字段，支持多行输入
      if (params.hostnames) {
        params.hostnames = splitByLineAndTrim(params.hostnames);
      }
      if (params.query) {
        params.query = splitByLineAndTrim(params.query);
      }
      params.hiddenVirtualBizServer = 1;
      return params;
    },
  });

  onMounted(async () => {
    const res = await GetAuthedBiz({});
    if (res.categories) {
      res.categories.forEach((item) => {
        if (item.subcategories) {
          item.subcategories.forEach((subItem) => {
            subCategoryOptions.push({
              value: subItem.name,
              label: subItem.name,
            });
          });
        }
      });
    }
    await getForm().updateSchema([
      {
        field: 'business',
        componentProps: {
          options: subCategoryOptions,
        },
      },
    ]);
  });
</script>
