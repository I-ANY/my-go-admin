<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="'business:overprovisioning:records:export'"
          :loading="exportState.loading"
          @click="handleExportRecords"
        >
          {{ exportState.text }}
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'result'">
          <Tag :color="record.result === '正常' ? 'green' : 'red'">
            {{ record.result }}
          </Tag>
        </template>
        <template v-if="column.key == 'standard'">
          <Tooltip :title="record.standard">
            <span>{{ record.standard }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key == 'currentConfiguration'">
          <Tooltip :title="record.currentConfiguration">
            <span>{{ record.currentConfiguration }}</span>
          </Tooltip>
        </template>
      </template>
    </BasicTable>
  </div>
</template>

<script lang="ts" setup>
  import { reactive, nextTick, h } from 'vue';
  import { Tag, Tooltip, Modal } from 'ant-design-vue';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { BasicTable, useTable } from '@/components/Table';
  import { useMessage } from '@/hooks/web/useMessage';
  import { getRecordList, Api } from '@/api/business/overprovisioning';
  import { downloadFileByUrl } from '@/utils/download';
  import { columns, searchFormSchema } from './data';

  defineOptions({ name: 'OverProvisioningRecords' });

  const { notification } = useMessage();

  const exportState = reactive({
    loading: false,
    text: '导出数据',
  });

  const [registerTable, { getForm }] = useTable({
    title: '检测记录',
    api: getRecordList,
    columns,
    formConfig: {
      labelWidth: 80,
      schemas: searchFormSchema,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    rowKey: 'id',
    beforeFetch: (params) => {
      return normalizeFilters(params);
    },
  });

  function normalizeFilters(source: Record<string, any>) {
    const params = { ...source };

    const normalizeMultilineField = (field: string) => {
      if (params[field] !== undefined && params[field] !== null) {
        if (typeof params[field] === 'string') {
          const list = params[field]
            .split('\n')
            .map((item: string) => item.trim())
            .filter((item: string) => item.length > 0);
          if (list.length > 0) {
            params[field] = list;
          } else {
            delete params[field];
          }
        } else if (Array.isArray(params[field]) && params[field].length === 0) {
          delete params[field];
        }
      }
    };

    normalizeMultilineField('node');
    normalizeMultilineField('hostname');
    normalizeMultilineField('sn');

    if (Array.isArray(params.business) && params.business.length === 0) {
      delete params.business;
    }

    if (params.dateRange) {
      params.fromDate = params.dateRange[0];
      params.toDate = params.dateRange[1];
      delete params.dateRange;
    }
    return params;
  }

  function handleExportRecords() {
    Modal.confirm({
      title: '确认导出',
      icon: h(ExclamationCircleOutlined),
      content: '是否导出当前筛选条件下的检测记录？',
      okText: '导出',
      cancelText: '取消',
      async onOk() {
        const form = await getForm();
        await form.validate();
        const values = await form.getFieldsValue();
        const payload = normalizeFilters({ ...values });

        nextTick(() => {
          exportState.loading = true;
          exportState.text = '导出中...';
        });

        try {
          await downloadFileByUrl(Api.RecordListExport, 'POST', 5 * 60, payload, null);
          notification.success({
            message: '导出成功',
            description: '文件已开始下载，请稍候查收。',
          });
        } catch (error: any) {
          notification.error({
            message: '导出失败',
            description: error?.message || '未知错误',
          });
        } finally {
          nextTick(() => {
            exportState.loading = false;
            exportState.text = '导出数据';
          });
        }
      },
    });
  }
</script>
