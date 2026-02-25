<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button type="primary" @click="gatherTrafficData" :loading="loading">95值计算</a-button>
        <a-button type="primary" @click="exportTraffic" :loading="loading">
          {{ data.exportButTitle }}
        </a-button>
      </template>
    </BasicTable>
    <modal95 @register="register" />
  </div>
</template>
<script setup lang="ts">
  import { BasicTable, useTable } from '@/components/Table';
  import { Api, GetZPTrafficDetailList } from '@/api/business/zp';
  import { trafficColumns, trafficSearchFormSchema } from '@/views/business/zp/data';
  import { splitByLineAndTrim } from '@/utils/util';
  import { RangePickPresetsExact } from '@/utils/common';
  import { nextTick, onMounted, reactive, ref, h } from 'vue';
  import dayjs from 'dayjs';
  import { useModal } from '@/components/Modal';
  import modal95 from './95modal.vue';
  import { notification, Modal } from 'ant-design-vue';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { defHttp } from '@/utils/http/axios';

  const [register, { openModal: openModal95 }] = useModal();
  let loading = ref(false);
  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
  });
  const [registerTable, { getForm }] = useTable({
    title: '流量计费数据',
    api: GetZPTrafficDetailList,
    columns: trafficColumns,
    formConfig: {
      labelWidth: 120,
      schemas: trafficSearchFormSchema(onTimePikerOpen),
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
      params.z_device_ids = splitByLineAndTrim(params.z_device_ids) || null;
    },
  });

  function onTimePikerOpen() {
    getForm().updateSchema({
      field: '[startTimestmp, endTimestmp]',
      componentProps: {
        presets: RangePickPresetsExact(),
      },
    });
  }

  onMounted(async () => {
    await resetReportTime();
  });

  function resetReportTime(): Promise<void> {
    return getForm().setFieldsValue({
      startTimestmp: dayjs(
        dayjs().add(-1, 'day').format('YYYY-MM-DD 00:00:00'),
        'YYYY-MM-DD HH:mm:ss',
      ),
      endTimestmp: dayjs(
        dayjs().add(-1, 'day').format('YYYY-MM-DD 23:59:59'),
        'YYYY-MM-DD HH:mm:ss',
      ),
    });
  }

  function gatherTrafficData() {
    // 获取表单的实际数据
    const formData = getForm().getFieldsValue();
    openModal95(true, {
      // formData: formData, // 传递表单数据
      props: {
        formData, // 传递表单数据
      },
    });
  }

  function exportTraffic() {
    const formValue = getForm().getFieldsValue();
    formValue.hostnames = splitByLineAndTrim(formValue.hostnames) || null;
    formValue.labels = splitByLineAndTrim(formValue.labels) || null;
    Modal.confirm({
      title: '是否确认导出筛选的数据?',
      icon: h(ExclamationCircleOutlined),
      content: h(
        'div',
        { style: 'color:red;' },
        '导出将花费一些时间，请耐心等待！如果数据过大，可能会导出失败，请缩小导出条件后重试',
      ),
      async onOk() {
        nextTick(() => {
          data.exporting = true;
          data.exportButTitle = '导出中';
        });
        await ExportTrafficDetailList(formValue);
      },
    });
  }

  async function ExportTrafficDetailList(value: Recordable) {
    const res = await defHttp.post(
      {
        url: Api.ExportTrafficDetailList,
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
        res.headers['content-disposition'] || `attachment;filename=ZP_device_traffic_detail.csv`;
      console.log(contentDisposition);
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
