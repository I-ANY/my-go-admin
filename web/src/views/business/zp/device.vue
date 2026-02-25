<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="syncZPDevice" :loading="loading">同步设备</a-button>
        <a-button type="primary" @click="updateCookie">Cookie信息</a-button>
        <a-button
          type="primary"
          v-auth="'business:z:DeviceInfo:export'"
          @click="handleExportData"
          :loading="data.exporting"
          >{{ data.exportButTitle }}
        </a-button>
      </template>
    </BasicTable>

    <CookieModal @register="RegisterCookieModal" />
  </div>
</template>
<script setup lang="ts">
  import { BasicTable, useTable } from '@/components/Table';
  import { GetDeviceList, GetZPCookie, SyncZPDevice, Api } from '@/api/business/zp';
  import { deviceColumns, deviceSearchFormSchema } from '@/views/business/zp/data';
  import { splitByLineAndTrim } from '@/utils/util';
  import { message, notification } from 'ant-design-vue';
  import CookieModal from '@/views/business/zp/cookieModal.vue';
  import { useModal } from '@/components/Modal';
  import { nextTick, reactive, ref } from 'vue';
  import { defHttp } from '@/utils/http/axios';

  const loading = ref(false);
  const [RegisterCookieModal, { openModal: openCookieModal }] = useModal();

  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
  });

  const [registerTable, { reload, getForm }] = useTable({
    title: '设备列表',
    api: GetDeviceList,
    columns: deviceColumns,
    formConfig: {
      labelWidth: 120,
      schemas: deviceSearchFormSchema,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      // alwaysShowLines: 2,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {
      pageSizeOptions: ['10', '30', '50', '100', '500', '1000'],
    },
    beforeFetch: (params) => {
      params.hostnames = splitByLineAndTrim(params.hostnames) || null;
      params.device_ids = splitByLineAndTrim(params.device_ids) || null;
      // params.acceptance_status = splitByLineAndTrim(params.acceptance_status) || null;
      // params.real_time_status = splitByLineAndTrim(params.real_time_status) || null;
    },
  });

  async function syncZPDevice() {
    loading.value = true;
    try {
      await SyncZPDevice();
      message.success('同步成功');
      loading.value = false;
      await reload();
    } catch (e) {
      loading.value = false;
      message.error(e.message);
    }
  }

  async function updateCookie() {
    const record = await GetZPCookie();
    openCookieModal(true, {
      record,
      isUpdate: true,
    });
  }

  function handleExportData() {
    const formValue = getForm().getFieldsValue();
    formValue.hostnames = splitByLineAndTrim(formValue.hostnames) || null;
    formValue.labels = splitByLineAndTrim(formValue.labels) || null;
    (async function () {
      await getForm().validate();
      data.exporting = true;
      data.exportButTitle = '导出中...';
      try {
        await ExportDeviceList(formValue);
      } catch (error) {
        notification.error({
          message: '导出失败',
          description: error.message,
        });
        data.exporting = false;
        data.exportButTitle = '导出数据';
      }
    })();
  }

  async function ExportDeviceList(value: Recordable) {
    const res = await defHttp.post(
      {
        url: Api.ExportDeviceList,
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
        res.headers['content-disposition'] || `attachment;filename=zp_device_info.csv`;
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
        data.exporting = false;
        data.exportButTitle = '导出数据';
      });
    }
  }
</script>

<style scoped lang="less"></style>
