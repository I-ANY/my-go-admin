<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="'business:k:DailyPeak:export'"
          @click="handleExportData"
          :loading="data.exporting"
          >{{ data.exportButTitle }}</a-button
        >
        <a-button
          type="success"
          v-auth="'business:k:TrafficBilling:export'"
          @click="handleExportBillingData"
          :loading="data.exportBilling"
        >
          {{ data.exportBillingButTitle }}
        </a-button>
      </template>
      <template #title> 日峰值列表 </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'status'">
          <span
            v-if="[0, 1].includes(record.status)"
            :style="{
              color: record.status === 0 ? '#00b42a' : '#f53f3f',
              fontWeight: 600,
              fontSize: '14px',
            }"
          >
            {{ record.status === 0 ? '正常' : '异常' }}
          </span>
        </template>
        <template v-if="column.key == 'provider_id'">
          <span v-if="providerTypeMap[record.provider_id]">{{
            providerTypeMap[record.provider_id].dictLabel
          }}</span>
        </template>
      </template>
      <template #form-ecdn_diff_rate_range="{ model }">
        <FormItem>
          <Row :gutter="8">
            <Col :span="12">
              <InputNumber
                v-model:value="model.ecnd_diff_rate_min"
                placeholder="最小值"
                :min="-100"
                :max="100"
                :addon-after="h('span', { style: { color: '#000' } }, '%')"
                :precision="2"
                style="width: 100%"
              />
            </Col>
            <Col :span="12">
              <InputNumber
                v-model:value="model.ecdn_diff_rate_max"
                placeholder="最大值"
                :min="-100"
                :max="100"
                :addon-after="h('span', { style: { color: '#000' } }, '%')"
                :precision="2"
                style="width: 100%"
              />
            </Col>
          </Row>
        </FormItem>
      </template>
      <template #form-autoops_diff_rate_range="{ model }">
        <FormItem>
          <Row :gutter="8">
            <Col :span="12">
              <InputNumber
                v-model:value="model.autoops_diff_rate_min"
                class="w-full site-input-left"
                :min="-100"
                :max="100"
                :addon-after="h('span', { style: { color: '#000' } }, '%')"
                :precision="2"
                style="width: 100%"
                placeholder="最小值"
              />
            </Col>
            <Col :span="12">
              <InputNumber
                v-model:value="model.autoops_diff_rate_max"
                class="w-full site-input-right"
                :min="-100"
                :max="100"
                :addon-after="h('span', { style: { color: '#000' } }, '%')"
                :precision="2"
                style="width: 100%"
                placeholder="最大值"
            /></Col>
          </Row>
        </FormItem>
      </template>
    </BasicTable>
  </div>
</template>

<script setup lang="ts">
  import { onMounted, h, reactive, nextTick, ref } from 'vue';
  import { Modal, FormItem, Row, Col, InputNumber, DatePicker, Select } from 'ant-design-vue';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { BasicTable, useTable } from '@/components/Table';
  import { dailyPeekSearchFormSchema, dailyPeekColumns, providerTypeMap } from './data';
  import { Api, GetDailyPeakList } from '@/api/business/k';
  import { useAreaSelect } from '@/utils/kAreaSelect';
  import { useMessage } from '@/hooks/web/useMessage';
  import { downloadFileByUrl } from '@/utils/download';
  import dayjs from 'dayjs';
  // 获取数据时更新树形数据
  onMounted(async () => {
    try {
      const form = await getForm();
      // 区域联动筛选
      const { initAreaData } = useAreaSelect({
        form,
        fields: {
          area: 'area',
          province: 'province',
        },
      });
      await initAreaData();
    } catch (error) {
      console.error('获取数据失败:', error);
    }
  });

  const [registerTable, { getForm }] = useTable({
    title: '日峰值列表',
    api: GetDailyPeakList,
    columns: dailyPeekColumns,
    // 使用搜索表单配置
    useSearchForm: true,
    formConfig: {
      labelWidth: 120,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
      alwaysShowLines: 2,
      schemas: dailyPeekSearchFormSchema,
    },
    beforeFetch: (params) => {
      return params;
    },
    canResize: true,
    bordered: true,
    showIndexColumn: false,
    // rowKey: 'id',
  });

  // 导出数据
  const data = reactive({
    exporting: false,
    exportButTitle: '日峰值导出',
    exportBilling: false,
    exportBillingButTitle: '月账单导出',
  });
  const { notification } = useMessage();
  const exportBillingMonth = ref(dayjs().subtract(1, 'month'));
  const exportBillingProvider = ref();
  const providerOptions = [
    { label: '明赋专线', value: 55 },
    { label: '明赋汇聚', value: 80 },
    { label: '泓宁专线', value: 103 },
  ];

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
        await getForm().validate();
        const value = await getForm().getFieldsValue();
        if (value.provider_id) {
          value.provider_id = parseInt(value.provider_id);
          // 将isp,dev_type 转换为数字再转换为数组
          if (value.isp) {
            const isp_items: number[] = [];
            for (const isp_id of value.isp) {
              isp_items.push(parseInt(isp_id));
            }
            value.isp = isp_items;
          }
          if (value.dev_type) {
            const dev_type_items: number[] = [];
            for (const dev_type_id of value.dev_type) {
              dev_type_items.push(parseInt(dev_type_id));
            }
            value.dev_type = dev_type_items;
          }
        }
        nextTick(() => {
          data.exporting = true;
          data.exportButTitle = '导出中';
        });
        try {
          await exportDailyPeak(value);
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

  async function exportDailyPeak(value: Recordable) {
    try {
      let filename = await downloadFileByUrl(Api.ExportDailyPeak, 'POST', 5 * 60, value, null);
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

  function handleExportBillingData() {
    Modal.confirm({
      title: '导出流量账单数据',
      icon: h(ExclamationCircleOutlined),
      content: () =>
        h('div', [
          h('div', { style: 'color:red;margin-bottom:8px;' }, '导出将花费一些时间，请耐心等待'),
          h('div', { style: 'margin: 12px 0 4px 0; fontWeight: 500;' }, '月份'),
          h(DatePicker, {
            picker: 'month',
            value: exportBillingMonth.value,
            'onUpdate:value': (val: any) => (exportBillingMonth.value = val),
            style: 'width: 100%; margin-bottom: 12px;',
            placeholder: '请选择月份',
          }),
          h('div', { style: 'margin: 0 0 4px 0; fontWeight: 500;' }, '厂商'),
          h(
            Select,
            {
              value: exportBillingProvider.value,
              'onUpdate:value': (val) => (exportBillingProvider.value = val),
              style: 'width: 100%',
              placeholder: '请选择厂商',
              options: providerOptions,
              allowClear: false,
            },
            [],
          ),
        ]),
      async onOk() {
        if (!exportBillingMonth.value) {
          notification.error({ message: '请选择月份' });
          throw new Error('请选择月份');
        }
        if (!exportBillingProvider.value) {
          notification.error({ message: '请选择厂商' });
          throw new Error('请选择厂商');
        }
        await getForm().validate();
        const value = await getForm().getFieldsValue();
        value.month = exportBillingMonth.value.format('YYYY-MM');
        value.provider_id = exportBillingProvider.value;
        if (value.isp) {
          value.isp = value.isp.map((id: any) => parseInt(id));
        }
        if (value.dev_type) {
          value.dev_type = value.dev_type.map((id: any) => parseInt(id));
        }
        nextTick(() => {
          data.exportBilling = true;
          data.exportBillingButTitle = '导出中';
        });
        try {
          await exportTrafficBilling(value);
        } catch (error) {
          notification.error({
            message: '导出失败',
            description: error.message || '未知错误',
          });
          nextTick(() => {
            data.exportBilling = false;
            data.exportBillingButTitle = '月账单导出';
          });
        }
      },
    });
  }

  async function exportTrafficBilling(value: Recordable) {
    try {
      let filename = await downloadFileByUrl(Api.ExportTrafficBill, 'POST', 5 * 60, value, null);
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
        data.exportBilling = false;
        data.exportBillingButTitle = '月账单导出';
      });
    }
  }
</script>
