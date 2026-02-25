<template>
  <div>
    <Tabs v-model:activeKey="activeTabKey" type="card">
      <TabPane tab="汇聚" key="normal" v-auth="'business:k:devicemacs:normal'">
        <BasicTable @register="tableConfigs.normal.register">
          <template #toolbar>
            <a-button
              type="primary"
              v-auth="'business:k:devicemacs:normal:export'"
              @click="handleExportData"
              :loading="data.exporting"
              >{{ data.exportButTitle }}</a-button
            >
          </template>
          <template #bodyCell="{ column, record }">
            <template v-if="column.key == 'biz_type'">
              <span v-if="bizTypesMap[record.biz_type]">{{
                bizTypesMap[record.biz_type].dictLabel || record.biz_type
              }}</span>
            </template>
            <template v-if="column.key == 'provider'">
              <span v-if="providersMap[record.provider]">{{
                providersMap[record.provider].dictLabel || record.provider
              }}</span>
            </template>
            <template v-if="column.key == 'is_cover_diff_isp'">
              <span v-if="isCoverDiffIspMap[record.is_cover_diff_isp]">{{
                isCoverDiffIspMap[record.is_cover_diff_isp].dictLabel || record.is_cover_diff_isp
              }}</span>
            </template>
            <template v-if="column.key == 'op_status'">
              <Tooltip>
                <Tag
                  v-if="opStatusMap[record.op_status]"
                  :color="opStatusMap[record.op_status].color || 'default'"
                  >{{ opStatusMap[record.op_status].dictLabel }}
                </Tag>
              </Tooltip>
            </template>
            <template v-if="column.key == 'business_status'">
              <Tooltip>
                <Tag
                  v-if="BusinessStatusMap[record.business_status]"
                  :color="BusinessStatusMap[record.business_status].color || 'default'"
                  >{{ BusinessStatusMap[record.business_status].dictLabel }}
                </Tag>
              </Tooltip>
            </template>
            <template v-if="column.key == 'flow_upload_tx_status'">
              <Tag
                v-if="[0, 1, 2].includes(record.flow_upload_tx_status)"
                :color="record.flow_upload_tx_status === 0 ? 'error' : 'success'"
                style="font-weight: bold"
              >
                {{ record.flow_upload_tx_status === 0 ? '异常' : '正常' }}
              </Tag>
            </template>
            <template v-if="column.key == 'flow_upload_autoops_status'">
              <Tag
                v-if="[0, 1, 2].includes(record.flow_upload_autoops_status)"
                :color="record.flow_upload_autoops_status === 0 ? 'error' : 'success'"
                style="font-weight: bold"
              >
                {{ record.flow_upload_autoops_status === 0 ? '异常' : '正常' }}
              </Tag>
            </template>
            <template v-if="column.key == 'is_province_scheduling'">
              <span v-if="isProvinceSchedulingMap[record.is_province_scheduling]">{{
                isProvinceSchedulingMap[record.is_province_scheduling].dictLabel ||
                record.is_province_scheduling
              }}</span>
            </template>
            <template v-if="column.key == 'cover_diff_isp_id'">
              <span v-if="coverDiffIspMap[record.cover_diff_isp_id]">{{
                coverDiffIspMap[record.cover_diff_isp_id].dictLabel
              }}</span>
            </template>
            <template v-if="column.key == 'is_disabled'">
              <Tag
                v-if="[0, 1].includes(record.is_disabled)"
                :color="record.is_disabled === 0 ? 'green' : 'red'"
              >
                {{ record.is_disabled === 1 ? '是' : '否' }}
              </Tag>
            </template>
            <template v-if="column.key == 'is_pass'">
              <Tag :color="record.is_pass === 1 ? 'green' : 'red'">
                {{ record.is_pass === 1 ? '是' : '否' }}
              </Tag>
            </template>
            <template v-if="column.key == 'is_first_mac'">
              <Tag :color="record.is_first_mac === true ? 'green' : 'red'">
                {{ record.is_first_mac === true ? '是' : '否' }}
              </Tag>
            </template>
            <template v-if="column.key == 'is_continuous_disabled'">
              <span>{{ record.is_continuous_disabled === 1 ? '是' : '否' }}</span>
            </template>
            <template v-if="column.key == 'is_only_ipv6'">
              <span>{{ record.is_only_ipv6 === 1 ? '是' : '否' }}</span>
            </template>
            <template v-if="column.key == 'is_independent_deploy'">
              <span>{{ record.is_independent_deploy === 1 ? '是' : '否' }}</span>
            </template>
            <template v-if="column.key == 'disabled_reason'">
              <Tooltip :title="record.disabled_reason">
                <span
                  style="
                    display: inline-block;
                    max-width: 120px;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    vertical-align: middle;
                    white-space: nowrap;
                  "
                >
                  {{ record.disabled_reason }}
                </span>
              </Tooltip>
            </template>
            <template v-if="column.key === 'action'">
              <TableAction
                :actions="[
                  {
                    label: 'MAC替换',
                    disabled: !canOperate(record),
                    onClick: handleMacReplace.bind(null, record),
                    auth: 'business:k:mac:normal:replace',
                  },
                  {
                    label: '备注',
                    onClick: handleEditRemark.bind(null, record),
                    auth: 'business:k:mac:normal:edit',
                  },
                ]"
              />
            </template>
          </template>
          <template #form-offlineTime="{ model }">
            <FormItem>
              <Row>
                <Col :span="11">
                  <InputNumber
                    v-model:value="model['min_offline_time']"
                    class="w-full site-input-left"
                    :min="0"
                    style="text-align: center"
                    placeholder="最小时间(天)"
                  />
                </Col>
                <Col :span="2">
                  <Input
                    class="w-full site-input-split"
                    placeholder="~"
                    disabled
                    style="text-align: center"
                /></Col>
                <Col :span="11">
                  <InputNumber
                    v-model:value="model['max_offline_time']"
                    class="w-full site-input-right"
                    :min="1"
                    style="text-align: center"
                    placeholder="最大时间(天)"
                  />
                </Col>
              </Row>
            </FormItem>
          </template>
        </BasicTable>
      </TabPane>
      <TabPane tab="专线" key="specialLine" v-auth="'business:k:devicemacs:specialline'">
        <BasicTable @register="tableConfigs.specialLine.register">
          <template #toolbar>
            <a-button
              type="primary"
              v-auth="'business:k:devicemacs:specialline:export'"
              @click="handleExportData"
              :loading="data.exporting"
              >{{ data.exportButTitle }}</a-button
            >
          </template>
          <template #bodyCell="{ column, record }">
            <template v-if="column.key == 'biz_type'">
              <span v-if="bizTypesMap[record.biz_type]">{{
                bizTypesMap[record.biz_type].dictLabel || record.biz_type
              }}</span>
            </template>
            <template v-if="column.key == 'provider'">
              <span v-if="providersMap[record.provider]">{{
                providersMap[record.provider].dictLabel || record.provider
              }}</span>
            </template>
            <template v-if="column.key == 'is_cover_diff_isp'">
              <span v-if="isCoverDiffIspMap[record.is_cover_diff_isp]">{{
                isCoverDiffIspMap[record.is_cover_diff_isp].dictLabel || record.is_cover_diff_isp
              }}</span>
            </template>
            <template v-if="column.key == 'op_status'">
              <Tooltip>
                <Tag
                  v-if="opStatusMap[record.op_status]"
                  :color="opStatusMap[record.op_status].color || 'default'"
                  >{{ opStatusMap[record.op_status].dictLabel }}
                </Tag>
              </Tooltip>
            </template>
            <template v-if="column.key == 'business_status'">
              <Tooltip>
                <Tag
                  v-if="BusinessStatusMap[record.business_status]"
                  :color="BusinessStatusMap[record.business_status].color || 'default'"
                  >{{ BusinessStatusMap[record.business_status].dictLabel }}
                </Tag>
              </Tooltip>
            </template>
            <template v-if="column.key == 'flow_upload_tx_status'">
              <Tag
                v-if="[0, 1, 2].includes(record.flow_upload_tx_status)"
                :color="record.flow_upload_tx_status === 0 ? 'error' : 'success'"
                style="font-weight: bold"
              >
                {{ record.flow_upload_tx_status === 0 ? '异常' : '正常' }}
              </Tag>
            </template>
            <template v-if="column.key == 'flow_upload_autoops_status'">
              <Tag
                v-if="[0, 1, 2].includes(record.flow_upload_autoops_status)"
                :color="record.flow_upload_autoops_status === 0 ? 'error' : 'success'"
                style="font-weight: bold"
              >
                {{ record.flow_upload_autoops_status === 0 ? '异常' : '正常' }}
              </Tag>
            </template>
            <template v-if="column.key == 'is_province_scheduling'">
              <span v-if="isProvinceSchedulingMap[record.is_province_scheduling]">{{
                isProvinceSchedulingMap[record.is_province_scheduling].dictLabel ||
                record.is_province_scheduling
              }}</span>
            </template>
            <template v-if="column.key == 'cover_diff_isp_id'">
              <span v-if="coverDiffIspMap[record.cover_diff_isp_id]">{{
                coverDiffIspMap[record.cover_diff_isp_id].dictLabel
              }}</span>
            </template>
            <template v-if="column.key == 'is_disabled'">
              <Tag
                v-if="[0, 1].includes(record.is_disabled)"
                :color="record.is_disabled === 0 ? 'green' : 'red'"
              >
                {{ record.is_disabled === 1 ? '是' : '否' }}
              </Tag>
            </template>
            <template v-if="column.key == 'is_pass'">
              <Tag :color="record.is_pass === 1 ? 'green' : 'red'">
                {{ record.is_pass === 1 ? '是' : '否' }}
              </Tag>
            </template>
            <template v-if="column.key == 'is_first_mac'">
              <Tag :color="record.is_first_mac === true ? 'green' : 'red'">
                {{ record.is_first_mac === true ? '是' : '否' }}
              </Tag>
            </template>
            <template v-if="column.key == 'is_continuous_disabled'">
              <span>{{ record.is_continuous_disabled === 1 ? '是' : '否' }}</span>
            </template>
            <template v-if="column.key == 'is_only_ipv6'">
              <span>{{ record.is_only_ipv6 === 1 ? '是' : '否' }}</span>
            </template>
            <template v-if="column.key == 'is_independent_deploy'">
              <span>{{ record.is_independent_deploy === 1 ? '是' : '否' }}</span>
            </template>
            <template v-if="column.key == 'disabled_reason'">
              <Tooltip :title="record.disabled_reason">
                <span
                  style="
                    display: inline-block;
                    max-width: 120px;
                    overflow: hidden;
                    text-overflow: ellipsis;
                    vertical-align: middle;
                    white-space: nowrap;
                  "
                >
                  {{ record.disabled_reason }}
                </span>
              </Tooltip>
            </template>
            <template v-if="column.key === 'action'">
              <TableAction
                :actions="[
                  {
                    label: 'MAC替换',
                    disabled: !canOperate(record),
                    onClick: handleMacReplace.bind(null, record),
                    auth: 'business:k:mac:specialline:replace',
                  },
                  {
                    label: '备注',
                    onClick: handleEditRemark.bind(null, record),
                    auth: 'business:k:mac:specialline:edit',
                  },
                ]"
              />
            </template>
          </template>
          <template #form-offlineTime="{ model }">
            <FormItem>
              <Row>
                <Col :span="11">
                  <InputNumber
                    v-model:value="model['min_offline_time']"
                    class="w-full site-input-left"
                    :min="0"
                    style="text-align: center"
                    placeholder="最小时间(天)"
                  />
                </Col>
                <Col :span="2">
                  <Input
                    class="w-full site-input-split"
                    placeholder="~"
                    disabled
                    style="text-align: center"
                /></Col>
                <Col :span="11">
                  <InputNumber
                    v-model:value="model['max_offline_time']"
                    class="w-full site-input-right"
                    :min="1"
                    style="text-align: center"
                    placeholder="最大时间(天)"
                  />
                </Col>
              </Row>
            </FormItem>
          </template>
        </BasicTable>
      </TabPane>
    </Tabs>
    <MacEditModal @register="registerMacEditModal" @success="handleSuccess" />
    <MacReplaceModal @register="registerMacReplaceModal" @success="handleSuccess" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, TableAction, useTable } from '@/components/Table';
  import {
    Tabs,
    TabPane,
    Col,
    FormItem,
    Input,
    InputNumber,
    Row,
    Tag,
    Tooltip,
    Modal,
  } from 'ant-design-vue';
  // import MacReplaceHistoryModal from './MacReplaceHistoryModal.vue';

  import {
    normalColumns,
    MacSearchFormSchema,
    coverDiffIspMap,
    providersMap,
    bizTypesMap,
    isCoverDiffIspMap,
    isProvinceSchedulingMap,
    opStatusMap,
    BusinessStatusMap,
  } from './data';
  import {
    Api,
    GetBusinessList,
    GetDeviceMacsList,
    GetAreaList,
    GetDevTypeList,
  } from '@/api/business/k';
  import { splitByLineAndTrim } from '@/utils/util';
  import { useModal } from '@/components/Modal';
  import MacReplaceModal from './MacReplaceModal.vue';
  import MacEditModal from './MacEditModal.vue';
  import { usePermission } from '@/hooks/web/usePermission';
  import { h, nextTick, onMounted, reactive, ref } from 'vue';
  import { useMessage } from '@/hooks/web/useMessage';
  import { defHttp } from '@/utils/http/axios';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';

  const [registerMacEditModal, { openModal: openMacEditModal }] = useModal();
  const [registerMacReplaceModal, { openModal: openMacReplaceModal }] = useModal();
  const { hasPermission } = usePermission();

  // 当前激活的标签页
  const activeTabKey = ref<'normal' | 'specialLine'>('normal');

  // 检查权限并设置默认激活的 tab
  const checkPermissionsAndSetDefaultTab = () => {
    const hasNormalPermission = hasPermission('business:k:devicemacs:normal');
    const hasSpecialLinePermission = hasPermission('business:k:devicemacs:specialline');

    if (hasNormalPermission && !hasSpecialLinePermission) {
      activeTabKey.value = 'normal';
    } else if (!hasNormalPermission && hasSpecialLinePermission) {
      activeTabKey.value = 'specialLine';
    } else if (hasNormalPermission && hasSpecialLinePermission) {
      // 两个都有权限，保持默认值
      activeTabKey.value = 'normal';
    } else {
      // 都没有权限，设置为 normal（这种情况应该不会发生，因为页面本身应该有权限控制）
      activeTabKey.value = 'normal';
    }
  };

  // 表格配置对象
  const tableConfigs = {
    normal: {
      register: null as any,
      getForm: null as any,
      reload: null as any,
    },
    specialLine: {
      register: null as any,
      getForm: null as any,
      reload: null as any,
    },
  };

  // 创建表格配置的函数
  const createTableConfig = (bizType: 'normal' | 'specialLine') => {
    const title = bizType === 'normal' ? '汇聚设备列表' : '专线设备列表';
    const MacsColumns = normalColumns(bizType);
    const SearchFormSchema = MacSearchFormSchema(bizType);

    const [register, { getForm, reload }] = useTable({
      title,
      api: GetDeviceMacsList,
      columns: MacsColumns,
      formConfig: {
        labelWidth: 120,
        schemas: SearchFormSchema,
        autoSubmitOnEnter: true,
        showAdvancedButton: true,
        autoAdvancedLine: 4,
      },
      actionColumn: {
        width: 150,
        title: '操作',
        dataIndex: 'action',
        fixed: 'right',
        ifShow: function (): boolean {
          const editPermission =
            bizType === 'normal' ? 'business:k:mac:normal:edit' : 'business:k:mac:specialline:edit';
          return hasPermission(editPermission);
        },
      },
      beforeFetch(params) {
        params.biz_type = bizType;
        params.hostnames = splitByLineAndTrim(params.hostnames);
        params.macAddrs = splitByLineAndTrim(params.macAddrs);
        params.nodes = splitByLineAndTrim(params.nodes);
        if (params.min_offline_time) {
          params.min_offline_time = parseInt(params.min_offline_time) * 24;
        }
        if (params.max_offline_time) {
          params.max_offline_time = parseInt(params.max_offline_time) * 24;
        }

        return params;
      },
      useSearchForm: true,
      showTableSetting: true,
      bordered: true,
      showIndexColumn: false,
      pagination: {},
    });

    // 保存表格配置
    tableConfigs[bizType] = {
      register: register as any,
      getForm: getForm as any,
      reload: reload as any,
    };
  };

  // 创建两个表格的配置
  createTableConfig('normal');
  createTableConfig('specialLine');

  function canOperate(record: Recordable): boolean {
    return record.business_status === 0 && record.is_pass === 1;
  }

  function handleMacReplace(record: Recordable) {
    // 检查编辑权限
    const editPermission =
      record.biz_type === 'normal'
        ? 'business:k:mac:normal:replace'
        : 'business:k:mac:specialline:replace';

    if (!hasPermission(editPermission)) {
      notification.error({
        message: '权限不足',
        description: '您没有权限进行此操作',
      });
      return;
    }

    openMacReplaceModal(true, {
      record,
      isUpdate: true,
    });
  }

  function handleEditRemark(record: Recordable) {
    // 检查编辑权限
    const editPermission =
      record.biz_type === 'normal'
        ? 'business:k:mac:normal:edit'
        : 'business:k:mac:specialline:edit';

    if (!hasPermission(editPermission)) {
      notification.error({
        message: '权限不足',
        description: '您没有权限进行此操作',
      });
      return;
    }

    openMacEditModal(true, {
      record: record,
    });
  }

  function handleSuccess() {
    // 刷新当前激活的表格
    tableConfigs[activeTabKey.value].reload();
  }

  async function fetchOptions() {
    // 设备类型
    let devTypeData = await GetDevTypeList();
    let devTypeOptions: Array<any> = [];
    devTypeData.forEach((e) => {
      const option = {
        label: e,
        value: e,
      };
      devTypeOptions.push(option);
    });

    // 省份列表
    let params = {};
    let data = await GetAreaList(params);
    let proOptions: Array<any> = [];
    let pro = '';
    data.items.forEach((e) => {
      let option = {
        label: e.province_name,
        value: e.province_name,
      };
      if (pro != e.province) {
        proOptions.push(option);
      }
      pro = e.province;
    });

    // 为两个表格更新选项
    for (const bizType of ['normal', 'specialLine'] as const) {
      // 业务名称列表
      let businessData = await GetBusinessList({ bizType });
      let businessOptions: Array<any> = [];
      businessData.forEach((e) => {
        const option = {
          label: e,
          value: e,
        };
        businessOptions.push(option);
      });

      const form = await tableConfigs[bizType].getForm();
      await form.updateSchema([
        {
          field: 'business',
          componentProps: {
            options: businessOptions,
          },
        },
        {
          field: 'dev_type_names',
          componentProps: {
            options: devTypeOptions,
          },
        },
        {
          field: 'province_names',
          componentProps: {
            options: proOptions,
          },
        },
      ]);
    }
  }

  onMounted(async () => {
    // 检查权限并设置默认激活的 tab
    checkPermissionsAndSetDefaultTab();
    await fetchOptions();
  });

  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
  });
  const { notification } = useMessage();

  function handleExportData() {
    const currentBizType = activeTabKey.value;

    // 检查当前 tab 的导出权限
    const exportPermission =
      currentBizType === 'normal'
        ? 'business:k:devicemacs:normal:export'
        : 'business:k:devicemacs:specialline:export';

    if (!hasPermission(exportPermission)) {
      notification.error({
        message: '权限不足',
        description: '您没有权限导出当前类型的数据',
      });
      return;
    }

    Modal.confirm({
      title: '是否确认导出筛选的数据?',
      icon: h(ExclamationCircleOutlined),
      content: h(
        'div',
        { style: 'color:red;' },
        '导出将花费一些时间，请耐心等待！如果数据过大，可能会导出失败，请缩小导出条件后重试',
      ),
      async onOk() {
        const form = await tableConfigs[currentBizType].getForm();
        await form.validate();
        const value = await form.getFieldsValue();
        value.biz_type = currentBizType;
        nextTick(() => {
          data.exporting = true;
          data.exportButTitle = '导出中';
        });
        try {
          await ExportDeviceMacsStatus(value);
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
  async function ExportDeviceMacsStatus(value: Recordable) {
    const res = await defHttp.post(
      {
        url: Api.ExportDeviceMacs,
        responseType: 'json',
        data: value,
        timeout: 10 * 60 * 1000,
      },
      { isReturnNativeResponse: true },
    );
    try {
      const data = res.data;
      if (data.code === 200) {
        const { file_name, file_type, file_content } = data.data;
        // 将 Base64 字符串转换为 Uint8Array
        const byteCharacters = atob(file_content);
        const byteNumbers = new Array(byteCharacters.length);
        for (let i = 0; i < byteCharacters.length; i++) {
          byteNumbers[i] = byteCharacters.charCodeAt(i);
        }
        const byteArray = new Uint8Array(byteNumbers);
        // 创建 Blob 对象
        const blob = new Blob([byteArray], { type: file_type });
        // 创建下载链接
        const downloadUrl = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = downloadUrl;
        a.download = file_name;
        document.body.appendChild(a);
        a.click();
        URL.revokeObjectURL(downloadUrl);
        document.body.removeChild(a);
        notification.success({
          message: '导出成功',
          description: '文件名：' + file_name,
          duration: null,
        });
      } else {
        notification.error({
          message: '导出失败',
          description: data.msg || '未知错误',
        });
      }
    } finally {
      nextTick(() => {
        data.exporting = false;
        data.exportButTitle = '导出数据';
      });
    }
  }

  // // MAC替换历史
  // const [registerReplaceHistoryModal, { openModal: openReplaceHistoryModal }] = useModal();

  // async function handleReplaceHistory(macAddr: string) {
  //   openReplaceHistoryModal(true, macAddr);
  // }
</script>

<style scoped>
  :deep(.site-input-split) {
    border-right: 0;
    border-left: 0;
    border-radius: 0%;
    background-color: #fff;
  }

  :deep(.site-input-right) {
    border-left-width: 0;
    border-top-left-radius: 0%;
    border-bottom-left-radius: 0%;
  }

  :deep(.site-input-left) {
    border-right-width: 0;
    border-top-right-radius: 0%;
    border-bottom-right-radius: 0%;
  }

  :deep(.site-input-right:hover),
  :deep(.site-input-right:focus) {
    border-left-width: 1px;
  }

  :deep(.site-input-left:hover),
  :deep(.site-input-left:focus) {
    border-right-width: 1px;
  }

  [data-theme='dark'] :deep(.site-input-split) {
    background-color: transparent;
  }
</style>
