<template>
  <div>
    <BasicTable @register="register">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="'business:k:demanddetect:export'"
          @click="handleExportData"
          :loading="data.exporting"
          >{{ data.exportButTitle }}</a-button
        >
      </template>
    </BasicTable>
  </div>
</template>

<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { Modal } from 'ant-design-vue';
  import { DemandDetectResultColumns, searchDetectResultSchema } from './data';
  import { GetDemandDetectList, Api } from '@/api/business/k';
  import { reactive, nextTick, h } from 'vue';
  import dayjs from 'dayjs';
  import { useMessage } from '@/hooks/web/useMessage';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { downloadFileByUrl } from '@/utils/download';

  const { notification } = useMessage();

  const [register, { getForm }] = useTable({
    title: 'K专线定制需求探测结果',
    api: async (params) => {
      const result = await GetDemandDetectList(params);
      return result;
    },
    columns: DemandDetectResultColumns,
    formConfig: {
      labelWidth: 120,
      schemas: [...searchDetectResultSchema()],
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
    },
    beforeFetch(params) {
      if (params.detect_status !== undefined && params.detect_status !== null) {
        if (params.detect_status === '') {
          delete params.detect_status;
        } else {
          params.detect_status = Number(params.detect_status);
        }
      }
      return params;
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
  });

  nextTick(async () => {
    try {
      const form = await getForm();
      await form.setFieldsValue({
        detect_date: dayjs(),
      });
    } catch (error) {
      console.error('初始化探测日期失败:', error);
    }
  });

  defineOptions({ name: 'DemandDetectResult' });

  // 导出数据
  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
  });

  function handleExportData() {
    Modal.confirm({
      title: '是否确认导出筛选的数据?',
      icon: h(ExclamationCircleOutlined),
      content: h(
        'div',
        { style: 'color:red;' },
        '导出将花费一些时间，请耐心等待！如果数据过大，可能会导出失败，请缩小导出条件后重试',
      ),
      async onOk() {
        const form = await getForm();
        await form.validate();
        const value = await form.getFieldsValue();

        // 处理 detect_date 数据
        if (value.detect_date) {
          if (dayjs.isDayjs(value.detect_date)) {
            value.detect_date = value.detect_date.format('YYYY-MM-DD');
          } else if (typeof value.detect_date === 'string') {
            value.detect_date = value.detect_date.substring(0, 10);
          }
        } else {
          value.detect_date = '';
        }

        nextTick(() => {
          data.exporting = true;
          data.exportButTitle = '导出中';
        });
        try {
          await ExportDemandDetectData(value);
        } catch (error) {
          notification.error({
            message: '导出失败',
            description: error.message || '未知错误',
          });
          nextTick(() => {
            data.exporting = false;
            data.exportButTitle = '导出数据';
          });
        }
      },
    });
  }

  async function ExportDemandDetectData(value: Recordable) {
    try {
      let filename = await downloadFileByUrl(Api.ExportDemandDetect, 'POST', 5 * 60, value, null);
      notification.success({
        message: '导出成功',
        description: '文件名：' + filename,
        duration: null,
      });
    } catch (error) {
      notification.error({
        message: '导出失败',
        description: error.message,
      });
    } finally {
      nextTick(() => {
        data.exporting = false;
        data.exportButTitle = '导出数据';
      });
    }
  }
</script>
<style scoped>
  :deep(.expired-row) {
    background-color: #f5f5f5 !important;
  }
</style>
