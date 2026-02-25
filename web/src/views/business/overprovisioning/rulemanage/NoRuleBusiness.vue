<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" :loading="exportState.loading" @click="handleExport">
          {{ exportState.text }}
        </a-button>
      </template>
    </BasicTable>
  </div>
</template>

<script lang="ts" setup>
  import { reactive, nextTick, h } from 'vue';
  import { message, Modal } from 'ant-design-vue';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { BasicTable, useTable } from '@/components/Table';
  import { useMessage } from '@/hooks/web/useMessage';
  import {
    getNoRuleBusiness,
    updateNoRuleBusinessRemark,
    Api,
  } from '@/api/business/overprovisioning';
  import { downloadFileByUrl } from '@/utils/download';
  import { noRuleBusinessColumns, noRuleBusinessSearchFormSchema } from './data';

  defineOptions({ name: 'NoRuleBusiness' });

  const { notification } = useMessage();

  const exportState = reactive({
    loading: false,
    text: '导出数据',
  });

  // 备注编辑提交前的回调
  const beforeEditSubmit = async (data: Recordable) => {
    const { record, key, value } = data;
    if (key === 'remark') {
      try {
        // 尝试多种可能的字段名获取 businessId
        const businessId = record.businessId || record.BusinessId || record.id;
        if (!businessId) {
          console.error('无法获取业务ID，record:', record);
          message.error('无法获取业务ID，请刷新页面重试');
          return false;
        }
        const postData = {
          businessId: Number(businessId),
          remark: value || '',
        };
        await updateNoRuleBusinessRemark(postData);
        message.success('备注保存成功');
        reload();
        return true;
      } catch (error: any) {
        message.error(error?.message || '备注保存失败');
        return false;
      }
    }
    return true;
  };

  const [registerTable, { reload, getForm }] = useTable({
    title: '无超配规则业务',
    api: getNoRuleBusiness,
    columns: noRuleBusinessColumns,
    useSearchForm: true,
    formConfig: {
      labelWidth: 100,
      schemas: noRuleBusinessSearchFormSchema,
    },
    showTableSetting: true,
    bordered: true,
    showIndexColumn: true,
    rowKey: 'businessId',
    pagination: {
      pageSize: 10,
    },
    beforeEditSubmit,
  });

  // 导出功能
  function handleExport() {
    Modal.confirm({
      title: '确认导出',
      icon: h(ExclamationCircleOutlined),
      content: '是否导出当前筛选条件下的无规则业务列表？',
      okText: '导出',
      cancelText: '取消',
      async onOk() {
        const form = await getForm();
        await form.validate();
        const values = await form.getFieldsValue();

        nextTick(() => {
          exportState.loading = true;
          exportState.text = '导出中...';
        });

        try {
          // 使用 POST 请求导出
          await downloadFileByUrl(Api.NoRuleBusinessExport, 'POST', 5 * 60, values, null);
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
