<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="'business:z:Device:action'"
          @click="batchHandleOnline"
          :loading="batchOnlineTag"
          >批量上线
        </a-button>
        <a-button
          type="primary"
          v-auth="'business:z:Device:action'"
          @click="batchHandleOffline"
          :loading="batchOfflineTag"
          >批量临时下线
        </a-button>
        <a-button
          type="primary"
          v-auth="'business:z:DeviceInfo:export'"
          x
          @click="handleExportData"
          :loading="data.exporting"
          >{{ data.exportButTitle }}
        </a-button>
        <Dropdown>
          <template #overlay>
            <Menu>
              <MenuItem key="1" @click="copyHostnameLabel"> 复制主机名&label</MenuItem>
              <Menu>
                <MenuItem
                  key="2"
                  @click="batchHandleOfflineForever"
                  :loading="batchForeverOfflineTag"
                  v-show="canForeverOffline"
                >
                  批量永久下线
                </MenuItem>
              </Menu>
            </Menu>
          </template>
          <a-button>
            更多
            <DownOutlined />
          </a-button>
        </Dropdown>
      </template>
      <template #headerCell="{ column }">
        <template v-if="column.key == 'hostname'">
          主机名
          <CopyOutlined class="ml-2" @click="copy_data('hostname')" />
        </template>
        <template v-if="column.key == 'label'">
          节点label
          <CopyOutlined class="ml-2" @click="copy_data('label')" />
        </template>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'ipv4'">
          <Tooltip placement="topLeft" :overlayStyle="{ maxWidth: '220px' }">
            <template #title>{{ formatIpList(record.ipv4) }}</template>
            <span>{{ record.ipv4 }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key == 'ipv6'">
          <Tooltip placement="topLeft" :overlayStyle="{ maxWidth: '220px' }">
            <template #title>{{ formatIpList(record.ipv6) }}</template>
            <span>{{ record.ipv6 }}</span>
          </Tooltip>
        </template>
        <template v-if="column.key == 'status'">
          <Tooltip>
            <Tag
              v-if="deviceStatusMap[record.status]"
              :color="deviceStatusMap[record.status].color || 'default'"
              >{{ deviceStatusMap[record.status].dictLabel }}
            </Tag>
          </Tooltip>
        </template>
        <template v-if="column.key == 'period'">
          <span>{{ periodList.filter((item) => item.value == record.period)[0].label }}</span>
        </template>
        <template v-if="column.key == 'is_only_cover'">
          <span>{{ record.is_only_cover === true ? '是' : '否' }}</span>
        </template>
        <template v-if="column.key == 'is_intranet_resource'">
          <span>{{ record.is_intranet_resource === true ? '是' : '否' }}</span>
        </template>
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                label: '临时下线',
                // onClick: handleOffline.bind(null, record),
                disabled: !canTempOffline(record),
                color: 'error',
                auth: 'business:z:device:action',
                popConfirm: {
                  title: '是否临时下线？',
                  confirm: handleTempOffline.bind(null, record),
                },
              },
              {
                label: '上线',
                auth: 'business:z:device:action',
                disabled: !canOnline(record),
                loading: onlineTag[record.id],
                onClick: handleOnline.bind(null, record),
              },
            ]"
            :dropDownActions="[
              {
                label: '永久关闭',
                disabled: !canTempOffline(record),
                auth: 'business:z:device:action',
                color: 'error',
                popConfirm: {
                  title: '是否确认永久下线',
                  placement: 'left',
                  confirm: handleOffline.bind(null, record),
                },
              },
              {
                label: '修改带宽',
                auth: 'business:z:device:action',
                onClick: handleModifyBandwidth.bind(null, record),
              },
              {
                label: '修改覆盖规则',
                auth: 'business:z:device:action',
                onClick: handleModifyCoverRule.bind(null, record),
              },
              {
                label: '修改进晚高峰规则',
                auth: 'business:z:device:action',
                onClick: handleModifyLatePeakRule.bind(null, record),
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <ModifyRuleModal @register="registerModifyModal" @success="reload" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, TableAction, useTable } from '@/components/Table';
  import { columns, searchFormSchema, periodList, deviceStatusMap } from './data';
  import {
    GetDeviceList,
    Api,
    DeviceTempOffline,
    DeviceOnline,
    DeviceForeverOffline,
    DeviceUpdate,
  } from '@/api/business/zap';
  import {
    Tooltip,
    notification,
    Tag,
    message,
    MenuItem,
    Menu,
    Dropdown,
    Modal,
  } from 'ant-design-vue';
  import { splitByLineAndTrim } from '@/utils/util';
  import { nextTick, reactive, ref } from 'vue';
  import { defHttp } from '@/utils/http/axios';
  import { CopyOutlined, DownOutlined } from '@ant-design/icons-vue';
  import { usePermission } from '@/hooks/web/usePermission';
  import { useModal } from '@/components/Modal';
  import ModifyRuleModal from './ModifyRuleModal.vue';

  defineOptions({ name: 'ZapDevice' });

  const [registerModifyModal, { openModal: openModifyModal }] = useModal();

  const batchOfflineTag = ref(false);
  const batchForeverOfflineTag = ref(false);
  const batchOnlineTag = ref(false);
  const onlineTag = reactive({});
  const offlineTag = reactive({});
  const tempOfflineTag = reactive({});

  const { hasPermission } = usePermission();

  const canForeverOffline = (): boolean => {
    return hasPermission('business:z:device:action');
  };

  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
  });

  const [registerTable, { getForm, getDataSource, getSelectRows, setSelectedRows, reload }] =
    useTable({
      title: '设备列表',
      api: GetDeviceList,
      columns,
      formConfig: {
        labelWidth: 10,
        schemas: searchFormSchema,
        autoSubmitOnEnter: true,
        showAdvancedButton: true,
        autoAdvancedLine: 1,
        alwaysShowLines: 2,
      },
      useSearchForm: true,
      showTableSetting: true,
      actionColumn: {
        width: '150px',
        title: '操作',
        dataIndex: 'action',
        fixed: 'right',
        ifShow: function (): boolean {
          return hasPermission('business:z:device:action');
        },
      },
      showSelectionBar: true, // 显示多选状态栏
      rowSelection: {
        type: 'checkbox',
      },
      bordered: true,
      showIndexColumn: false,
      pagination: {
        pageSizeOptions: ['10', '30', '50', '100', '500', '1000'],
      },
      beforeFetch: (params) => {
        params.hostnames = splitByLineAndTrim(params.hostnames) || null;
        params.labels = splitByLineAndTrim(params.labels) || null;
      },
    });

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
        data.exporting = false;
        data.exportButTitle = '导出数据';
      });
    }
  }

  async function copy_data(column: string) {
    try {
      const data = getDataSource()
        .map((item) => item[column])
        .join('\n');
      // 优先使用 Clipboard API
      if (navigator.clipboard) {
        await navigator.clipboard.writeText(data);
      } else {
        // 兼容旧浏览器
        const textarea = document.createElement('textarea');
        textarea.value = data;
        document.body.appendChild(textarea);
        textarea.select();
        document.execCommand('copy');
        document.body.removeChild(textarea);
      }
      message.success('已复制到剪切板');
    } catch (err) {
      console.error('复制失败:', err);
    }
  }

  async function copyHostnameLabel() {
    try {
      const data = getDataSource()
        .map((item) => item['hostname'] + '  ' + item['label'])
        .join('\n');
      // 优先使用 Clipboard API
      if (navigator.clipboard) {
        await navigator.clipboard.writeText(data);
      } else {
        // 兼容旧浏览器
        const textarea = document.createElement('textarea');
        textarea.value = data;
        document.body.appendChild(textarea);
        textarea.select();
        document.execCommand('copy');
        document.body.removeChild(textarea);
      }
      message.success('已复制到剪切板');
    } catch (err) {
      console.error('复制失败:', err);
    }
  }

  function canTempOffline(record: Recordable): boolean {
    return record.status == 1 && record.period != 2;
  }

  function canOnline(record: Recordable): boolean {
    return record.status == 2 || record.status == 4;
  }

  // 临时下线
  async function handleTempOffline(record: Recordable) {
    tempOfflineTag[record.id] = true;
    let param = {
      ids: [record.id],
    };
    try {
      await DeviceTempOffline(param);
      message.success('操作成功!');
      await reload();
    } finally {
      tempOfflineTag[record.id] = false;
    }
  }

  // 永久下线
  async function handleOffline(record: Recordable) {
    offlineTag[record.id] = true;
    let param = {
      ids: [record.id],
    };
    try {
      await DeviceForeverOffline(param);
      message.success('操作成功!');
      await reload();
    } finally {
      offlineTag[record.id] = false;
    }
  }

  // 批量临时下线
  async function batchHandleOffline() {
    batchOfflineTag.value = true;
    const selectRows = getSelectRows();
    const ids = selectRows.filter((item) => item.status === 1).map((item) => item.id);
    if (ids.length === 0) {
      message.error('未选择设备或所选设备均非Active状态!');
      batchOfflineTag.value = false;
      return;
    }
    if (selectRows.length !== ids.length) {
      message.error('有选择非Active设备，请重新选择!');
      batchOfflineTag.value = false;
      return;
    }
    Modal.confirm({
      title: '确认操作',
      content: `确定要临时下线选中的 ${ids.length} 台设备吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await DeviceTempOffline({ ids: ids });
          message.success('操作成功!');
          await reload();
          setSelectedRows([]);
        } finally {
          batchOfflineTag.value = false;
        }
      },
      onCancel: () => {
        batchOfflineTag.value = false;
      },
    });
  }

  // 批量永久下线
  async function batchHandleOfflineForever() {
    batchForeverOfflineTag.value = true;
    const selectRows = getSelectRows();
    const ids = selectRows.filter((item) => item.status === 1).map((item) => item.id);
    if (ids.length === 0) {
      message.error('未选择设备或所选设备均非Active状态!');
      batchForeverOfflineTag.value = false;
      return;
    }
    if (selectRows.length !== ids.length) {
      message.error('有选择非Active设备，请重新选择!');
      batchForeverOfflineTag.value = false;
      return;
    }
    if (ids.length > 20) {
      message.error('最多只能选择20台设备!');
      batchForeverOfflineTag.value = false;
      return;
    }
    // 添加确认弹窗
    Modal.confirm({
      title: '确认操作',
      content: `确定要永久下线选中的 ${ids.length} 台设备吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await DeviceForeverOffline({ ids: ids });
          message.success('操作成功!');
        } finally {
          setSelectedRows([]);
          batchForeverOfflineTag.value = false;
          await reload();
        }
      },
    });
  }

  // 上线
  async function handleOnline(record: Recordable) {
    onlineTag[record.id] = true;
    let param = {
      ids: [record.id],
    };
    try {
      await DeviceOnline(param);
      message.success('操作成功!');
      await reload();
    } finally {
      offlineTag[record.id] = false;
    }
  }

  // 批量上线
  async function batchHandleOnline() {
    batchOnlineTag.value = true;
    const selectRows = getSelectRows();
    const ids = selectRows.filter((item) => item.status != 1).map((item) => item.id);
    if (ids.length === 0) {
      message.error('未选择设备或所选设备中有Active状态设备!');
      batchOnlineTag.value = false;
      return;
    }

    if (selectRows.length !== ids.length) {
      message.error('有选择Active设备，请重新选择!');
      batchOnlineTag.value = false;
      return;
    }
    // 添加确认弹窗
    Modal.confirm({
      title: '确认操作',
      content: `确定要上线选中的 ${ids.length} 台设备吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await DeviceOnline({ ids: ids });
          message.success('操作成功!');
          await reload();
          setSelectedRows([]);
        } finally {
          batchOnlineTag.value = false;
        }
      },
      onCancel: () => {
        batchOnlineTag.value = false;
      },
    });
  }

  // 模拟 API，请替换为真实 API
  const updateBandwidthApi = async (params: any) => {
    await DeviceUpdate(params);
    message.success('操作成功!');
  };

  const updateCoverRuleApi = async (params: any) => {
    await DeviceUpdate(params);
    message.success('操作成功!');
  };

  const updateLatePeakRuleApi = async (params: any) => {
    console.log('Update Late Peak Rule:', params);
    message.warning('功能开发中，敬请期待！');
  };

  function handleModifyBandwidth(record: Recordable) {
    openModifyModal(true, {
      title: '修改带宽',
      recordId: record.id,
      schemas: [
        {
          field: 'total_bandwidth',
          label: '节点带宽容量(Gbps)',
          component: 'InputNumber',
          required: true,
          colProps: { span: 20 },
        },
        {
          field: 'single_line_bandwidth',
          label: '单线带宽(Mbps)',
          component: 'InputNumber',
          required: true,
          colProps: { span: 20 },
        },
        {
          field: 'bwcount',
          label: '线路数',
          component: 'InputNumber',
          required: true,
          colProps: { span: 20 },
        },
      ],
      values: {
        total_bandwidth: record.capacity,
        single_line_bandwidth: getSingleLineBandwidth(record),
        bwcount: countIpv4(record.ipv4),
      },
      submitFunc: updateBandwidthApi,
    });
  }

  function handleModifyCoverRule(record: Recordable) {
    openModifyModal(true, {
      title: '修改覆盖规则',
      recordId: record.id,
      schemas: [
        {
          field: 'is_only_cover',
          label: '是否只覆盖本省',
          component: 'Select',
          required: true,
          colProps: { span: 20 },
          componentProps: {
            options: [
              {
                label: '是',
                value: true,
              },
              {
                label: '否',
                value: false,
              },
            ],
          },
        },
      ],
      values: {
        is_only_cover: record.is_only_cover,
      },
      submitFunc: updateCoverRuleApi,
    });
  }

  function handleModifyLatePeakRule(record: Recordable) {
    openModifyModal(true, {
      title: '修改仅晚高峰规则',
      recordId: record.id,
      schemas: [
        {
          field: 'only_evening_peak',
          label: '是否只覆盖晚高峰',
          component: 'Select',
          required: true,
          colProps: { span: 20 },
          componentProps: {
            options: [
              {
                label: '是',
                value: true,
              },
              {
                label: '否',
                value: false,
              },
            ],
          },
        },
      ],
      values: {
        only_evening_peak: record.only_evening_peak,
      },
      submitFunc: updateLatePeakRuleApi,
    });
  }

  // 添加格式化IP地址的函数
  function formatIpList(ipStr: string): string {
    if (!ipStr) return '';
    return ipStr.split(',').join('\n');
  }

  // 计算节点的IPV4的个数
  function countIpv4(ipStr: string): number {
    if (!ipStr) return 0; // 处理空字符串的情况
    return ipStr.split(',').length;
  }

  // 计算单线带宽
  function getSingleLineBandwidth(record: Recordable) {
    if (!record.ipv4) return 0;
    return (record.capacity / countIpv4(record.ipv4)) * 1000;
  }
</script>
