<template>
  <PageWrapper
    :title="'#' + taskId + ' 执行详情'"
    :contentFullHeight="false"
    contentBackground
    @back="goBack"
  >
    <BasicTable @register="registerTable">
      <template #headerCell="{ column }">
        <template v-if="column.key == 'hostname'">
          主机名
          <Tooltip title="主机名复制">
            <CopyOutlined class="ml-2" @click="handleCopyHostnames" />
          </Tooltip>
        </template>
        <template v-if="column.key == 'result'">
          结果详情
          <Tooltip title="结果复制">
            <CopyOutlined class="ml-2" @click="handleCopyResult" />
          </Tooltip>
        </template>
      </template>
      <template #toolbar>
        <a-button type="primary" size="small" @click="handleExportData" style="margin-right: 20px">
          {{ exportState.title }}
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'result'">
          <Popover
            trigger="hover"
            placement="top"
            :mouseLeaveDelay="0.2"
            :overlayStyle="{ maxWidth: '1200px' }"
            :overlayInnerStyle="{ padding: '0px' }"
          >
            <template #content>
              <div style="position: relative">
                <div class="code-copy-btn" @click="handleCopy(record.result)">
                  <CopyOutlined />
                  <span class="copy-text">复制</span>
                </div>
                <div class="result-code-box" v-html="formatResult(record.result)"></div>
              </div>
            </template>
            <span>{{ cutStringLength(record.result, 100) }}</span>
          </Popover>
          <span v-if="record.result && record.result.length > 100">{{
            cutStringLength(record.result, 500).slice(100)
          }}</span>
        </template>
      </template>
    </BasicTable>
  </PageWrapper>
</template>
<script setup lang="ts">
  import { BasicTable, useTable } from '@/components/Table';
  import { executeDetailColumns, executeDetailSearchFormSchema } from './data';
  import {
    cutStringLength,
    splitByLineAndTrim,
    highlightKeywords,
    handleCopyTableRowData,
  } from '@/utils/util';
  import { PageWrapper } from '@/components/Page';
  import { useRoute } from 'vue-router';
  import { useGo } from '@/hooks/web/usePage';
  import { notification, Popover, message } from 'ant-design-vue';
  import { GetExecuteRecordDetail } from '@/api/ops/execute';
  import { reactive, nextTick } from 'vue';
  import { defHttp } from '@/utils/http/axios';
  import { CopyOutlined } from '@ant-design/icons-vue';
  import Tooltip from 'ant-design-vue/lib/tooltip';

  const route = useRoute();
  let taskId: string = route.params.taskId as string;
  const go = useGo();
  const exportState = reactive({
    exporting: false,
    title: '导出',
  });

  const [registerTable, { getDataSource }] = useTable({
    title: '查询执行详情',
    api: GetExecuteRecordDetail,
    beforeFetch: async (params) => {
      params.taskId = taskId;
      params.hostnames = splitByLineAndTrim(params.hostnames);
    },
    columns: executeDetailColumns,
    formConfig: {
      labelWidth: 10,
      schemas: executeDetailSearchFormSchema,
      autoSubmitOnEnter: true,
      submitOnChange: true,
      baseRowStyle: {
        marginTop: '0',
        paddingTop: '0',
      },
      compact: true,
      actionColOptions: {
        span: 6,
      },
    },
    useSearchForm: true,
    showTableSetting: false,
    bordered: true,
    showIndexColumn: false,
    showSelectionBar: false, // 显示多选状态栏
    pagination: {
      // pageSizeOptions: ['10', '30', '50'],
    },
  });

  // 数据导出
  function handleExportData() {
    (async function () {
      exportState.exporting = true;
      exportState.title = '导出中...';
      try {
        await DoExportExeResultDetail({ taskId: Number(taskId) });
      } catch (error) {
        notification.error({
          message: '导出失败',
          description: error.message,
        });
        exportState.exporting = false;
        exportState.title = '导出';
      }
    })();
  }

  async function DoExportExeResultDetail(value: Recordable) {
    const res = await defHttp.post(
      {
        url: '/v1/ops/execute/tasks/result/export',
        responseType: 'blob',
        data: value,
        timeout: 10 * 60 * 1000,
      },
      { isReturnNativeResponse: true },
    );
    try {
      if (!res.headers['content-type'].includes('application/octet-stream')) {
        // 将 Blob 转换为 JSON
        const reader = new FileReader();
        reader.onload = () => {
          const jsonResponse = JSON.parse(reader.result as any);
          notification.error({
            message: '导出失败',
            description: jsonResponse.msg || '未知错误',
            duration: null,
          });
        };
        reader.readAsText(res.data);
        return;
      }
      const blob = new Blob([res.data], { type: res.headers['content-type'] });
      // 创建新的URL并指向File对象或者Blob对象的地址
      const blobURL = window.URL.createObjectURL(blob);
      // 创建a标签，用于跳转至下载链接
      const tempLink = document.createElement('a');
      tempLink.style.display = 'none';
      tempLink.href = blobURL;
      const contentDisposition =
        res.headers['content-disposition'] || `attachment;filename=hdd_device_info.csv`;
      const filename = contentDisposition.split(';')[1].split('=')[1].split("''")[1];
      tempLink.setAttribute('download', filename);
      // 兼容：某些浏览器不支持HTML5的download属性
      if (typeof tempLink.download === 'undefined') {
        tempLink.setAttribute('target', '_blank');
      }
      // 挂载a标签
      document.body.appendChild(tempLink);
      tempLink.click();
      document.body.removeChild(tempLink);
      // 释放blob URL地址
      window.URL.revokeObjectURL(blobURL);
      notification.success({
        message: '导出成功',
        description: '文件名：' + filename,
        duration: null,
      });
    } finally {
      nextTick(() => {
        exportState.exporting = false;
        exportState.title = '导出';
      });
    }
  }

  function goBack() {
    go('/ops/execute/history');
  }

  // 复制文本内容
  function handleCopy(text: string) {
    if (!text) {
      message.warning('内容为空');
      return;
    }
    try {
      if (navigator.clipboard) {
        navigator.clipboard.writeText(text);
      } else {
        // 兼容旧浏览器
        const textarea = document.createElement('textarea');
        textarea.value = text;
        document.body.appendChild(textarea);
        textarea.select();
        document.execCommand('copy');
        document.body.removeChild(textarea);
      }
      message.success('复制成功');
    } catch (error) {
      message.error('复制失败');
    }
  }

  // 格式化结果文本，针对特定格式进行换行处理
  function formatResult(text: string): string {
    let html = highlightKeywords(text);
    // 针对特定格式 ", result:" 后面添加换行
    // 匹配 pattern: , result: 后面允许有空格
    html = html.replace(/(, result:)\s*/gi, '$1\n');
    return html;
  }

  // 复制主机名到剪贴板
  async function handleCopyHostnames() {
    const dataSource = getDataSource();
    await handleCopyTableRowData(dataSource, (item) => item.hostname, '主机名');
  }

  // 复制结果到剪贴板
  async function handleCopyResult() {
    const dataSource = getDataSource();
    await handleCopyTableRowData(dataSource, (item) => item.result, '结果');
  }
</script>
<style scoped lang="less">
  .code-copy-btn {
    display: flex;
    position: absolute;
    z-index: 10;
    top: 5px;
    right: 10px;
    align-items: center;
    padding: 4px 8px;
    transition: all 0.3s;
    border-radius: 4px;
    opacity: 0.6;
    background-color: rgb(255 255 255 / 10%);
    color: #abb2bf;
    cursor: pointer;

    &:hover {
      opacity: 1;
      background-color: rgb(255 255 255 / 20%);
      color: #fff;
    }

    .anticon {
      font-size: 14px;
    }

    .copy-text {
      margin-left: 4px;
      font-size: 12px;
    }
  }

  .result-code-box {
    max-width: 1150px;
    max-height: 490px;
    padding: 20px;
    padding-top: 32px; // 增加顶部内边距，防止第一行文字被复制按钮遮挡
    overflow: auto;
    border-radius: 6px;
    background-color: #0b0b0b; // 深色背景，像 IDE
    color: #c6c1c1; // 浅色文字
    font-family: Consolas, Monaco, 'Andale Mono', 'Ubuntu Mono', monospace;
    font-size: 13px;
    line-height: 1.5;
    word-break: break-all;
    white-space: pre-wrap;
    user-select: text;
  }
</style>
