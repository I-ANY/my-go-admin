<template>
  <div>
    <BasicModal
      v-bind="$attrs"
      @register="registerModal"
      :title="title"
      width="1200px"
      :footer="null"
    >
      <BasicTable @register="registerTable">
        <template #bodyCell="{ record, column }">
          <template v-if="column.key == 'taskName'">
            <Tooltip placement="topLeft" :overlayStyle="{ maxWidth: '1200px', maxHeight: '500px' }">
              <template #title>
                {{ record.taskName }}
              </template>
              <span>{{ record.taskName }}</span>
            </Tooltip>
          </template>
          <template v-if="column.key == 'result'">
            <Tooltip placement="topLeft" :overlayStyle="{ maxWidth: '1200px', maxHeight: '500px' }">
              <template #title>
                <div class="tooltip-content" style="max-width: 1000px; max-height: 450px">
                  {{ record.result }}
                </div>
              </template>
              <span>{{ record.result }}</span>
            </Tooltip>
          </template>
        </template>
      </BasicTable>
    </BasicModal>
  </div>
</template>
<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { BasicTable, useTable } from '@/components/Table';
  import { deviceExecuteHistoryResultColumns } from '@/views/ops/batchExecute/data';
  import { GetExecuteRecordDetail } from '@/api/ops/execute';
  import Tooltip from 'ant-design-vue/lib/tooltip';
  import { ref } from 'vue';

  let title = ref('');
  let hostname = '';
  const [registerModal] = useModalInner(async (data) => {
    hostname = data.hostname;
    title.value = '执行历史：' + hostname;
    await reload();
  });

  const [registerTable, { reload }] = useTable({
    title: '',
    api: async (params) => {
      params.hostname = [hostname];
      return GetExecuteRecordDetail(params);
    },
    columns: deviceExecuteHistoryResultColumns,
    useSearchForm: false,
    showTableSetting: false,
    bordered: true,
    showIndexColumn: false,
    showSelectionBar: false, // 显示多选状态栏
    striped: true,
    pagination: {
      // pageSizeOptions: ['10', '30', '50'],
    },
  });
</script>

<style scoped lang="less">
  .tooltip-content {
    max-height: 400px;
    overflow: auto;
    word-break: break-all;
    white-space: pre-wrap;
  }
</style>
