<template>
  <div>
    <Tabs v-model:activeKey="activeTabKey" type="card">
      <TabPane tab="汇聚" key="normal" v-auth="'business:k:demand:normal'">
        <BasicTable @register="tableConfigs.normal.register">
          <template #toolbar>
            <a-button
              type="primary"
              v-auth="'business:k:demand:normal:occupy'"
              @click="handleBatchOccupy('normal')"
              :disabled="!getSelectedCount('normal')"
              >批量占用</a-button
            >
            <a-button
              type="primary"
              v-auth="'business:k:demand:normal:export'"
              @click="handleExportData"
              :loading="data.exporting"
              >{{ data.exportButTitle }}</a-button
            >
          </template>
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'biz_type'">
              <span v-if="bizTypesMap[record.biz_type]">
                {{ bizTypesMap[record.biz_type].dictLabel || record.biz_type }}
              </span>
            </template>
            <template v-if="column.key === 'action'">
              <TableAction
                :actions="[
                  {
                    tooltip: '需求占用',
                    label: '占用',
                    onClick: handleOccupy.bind(null, record, 'normal'),
                    auth: 'business:k:demand:normal:occupy',
                    disabled:
                      dayjs(record.end_time).isBefore(dayjs().startOf('day')) ||
                      record.gap_bw - record.locked_bw <= 0,
                  },
                  {
                    label: record.is_locked ? '解锁' : '锁定',
                    tooltip: record.is_locked ? '解锁需求' : '锁定需求',
                    onClick: () => handleLockDemand(record, 'normal'),
                    auth: 'business:k:demand:normal:lock',
                    ifShow: true,
                    color: record.is_locked ? 'success' : 'error',
                  },
                ]"
              />
            </template>
          </template>
        </BasicTable>
      </TabPane>
      <TabPane tab="专线" key="specialLine" v-auth="'business:k:demand:specialline'">
        <BasicTable @register="tableConfigs.specialLine.register">
          <template #toolbar>
            <a-button
              type="primary"
              v-auth="'business:k:demand:specialline:occupy'"
              @click="handleBatchOccupy('specialLine')"
              :disabled="!getSelectedCount('specialLine')"
              >批量占用</a-button
            >
            <a-button
              type="primary"
              v-auth="'business:k:demand:specialline:export'"
              @click="handleExportData"
              :loading="data.exporting"
              >{{ data.exportButTitle }}</a-button
            >
          </template>
          <template #bodyCell="{ column, record }">
            <template v-if="column.key === 'biz_type'">
              <span v-if="bizTypesMap[record.biz_type]">
                {{ bizTypesMap[record.biz_type].dictLabel || record.biz_type }}
              </span>
            </template>
            <template v-if="column.key === 'action'">
              <TableAction
                :actions="[
                  {
                    tooltip: '需求占用',
                    label: '占用',
                    onClick: handleOccupy.bind(null, record, 'specialLine'),
                    auth: 'business:k:demand:specialline:occupy',
                    disabled:
                      dayjs(record.end_time).isBefore(dayjs().startOf('day')) ||
                      record.gap_bw - record.locked_bw <= 0,
                  },
                  {
                    label: record.is_locked ? '解锁' : '锁定',
                    tooltip: record.is_locked ? '解锁需求' : '锁定需求',
                    onClick: () => handleLockDemand(record, 'specialLine'),
                    auth: 'business:k:demand:specialline:lock',
                    ifShow: true,
                    color: record.is_locked ? 'success' : 'error',
                  },
                ]"
              />
            </template>
          </template>
        </BasicTable>
      </TabPane>
    </Tabs>
    <DemandOccupyModal @register="registerModal" @success="handleSuccess" />
    <LockDemandModal @register="registerLockModal" @success="handleSuccess" />
    <BatchOccupyModal @register="registerBatchModal" @success="handleBatchSuccess" />
  </div>
</template>

<script lang="ts" setup>
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { Tabs, TabPane, Modal } from 'ant-design-vue';
  import { demandColumns, searchFormSchema, bizTypesMap } from './data';
  import { Api, getDemandList, GetDevTypeList, lockDemand } from '@/api/business/k';
  import { ref, onMounted, watch, nextTick, h, reactive } from 'vue';
  import dayjs from 'dayjs';
  import DemandOccupyModal from './DemandOccupyModal.vue';
  import LockDemandModal from './LockDemandModal.vue';
  import BatchOccupyModal from './BatchOccupyModal.vue';
  import { useModal } from '@/components/Modal';
  import { useAreaSelect } from '@/utils/kAreaSelect';
  import { useMessage } from '@/hooks/web/useMessage';
  import { usePermission } from '@/hooks/web/usePermission';
  import { downloadFileStream } from '@/utils/download';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';

  // 权限检查
  const { hasPermission } = usePermission();

  // 当前激活的标签页
  const activeTabKey = ref<'normal' | 'specialLine'>('normal');

  // 检查权限并设置默认激活的 tab
  const checkPermissionsAndSetDefaultTab = () => {
    const hasNormalPermission = hasPermission('business:k:demand:normal');
    const hasSpecialLinePermission = hasPermission('business:k:demand:specialline');

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

  const devTypeOptions = ref<any[]>([]);

  // 获取设备类型数据
  const loadDevTypeData = async () => {
    try {
      const devTypeData = await GetDevTypeList();
      devTypeOptions.value = devTypeData.map((item) => ({
        label: item,
        value: item,
      }));
    } catch (error) {
      console.error('获取设备类型数据失败:', error);
    }
  };

  const [registerModal, { openModal }] = useModal();
  const [registerLockModal, { openModal: openLockModal }] = useModal();
  const [registerBatchModal, { openModal: openBatchModal }] = useModal();

  // 选中行状态
  const selectedRows = reactive<{
    normal: any[];
    specialLine: any[];
  }>({
    normal: [],
    specialLine: [],
  });

  // 记录当前已选中的设备类型（用于限制批量勾选同类型设备）
  const selectedDevType = reactive<{
    normal: string | null;
    specialLine: string | null;
  }>({
    normal: null,
    specialLine: null,
  });

  const { notification } = useMessage();

  // 记录是否已经手动操作过 end_time（用于区分首次加载和用户清空）
  const endTimeManuallySet = {
    normal: ref(false),
    specialLine: ref(false),
  };

  // 表格配置对象
  const tableConfigs = {
    normal: {
      register: null as any,
      getForm: null as any,
      reload: null as any,
      clearSelectedRowKeys: null as any,
    },
    specialLine: {
      register: null as any,
      getForm: null as any,
      reload: null as any,
      clearSelectedRowKeys: null as any,
    },
  };

  // 创建表格配置的函数
  const createTableConfig = (bizType: 'normal' | 'specialLine') => {
    const title = bizType === 'normal' ? '汇聚需求列表' : '专线需求列表';

    const [register, { getForm, reload, clearSelectedRowKeys }] = useTable({
      title,
      api: async (params) => {
        const result = await getDemandList(params);
        return result;
      },
      columns: demandColumns(bizType),
      formConfig: {
        labelWidth: 100,
        schemas: [...searchFormSchema(bizType)],
        autoSubmitOnEnter: true,
        showAdvancedButton: true,
        autoAdvancedLine: 4,
      },
      beforeFetch(params) {
        params.bizType = bizType;

        // 从 searchFormSchema 获取默认值
        const formSchema = searchFormSchema(bizType);
        const defaultProvider =
          formSchema.find((item) => item.field === 'provider')?.componentProps?.defaultValue ||
          'mf';
        const defaultOnlyCanSubmit =
          formSchema.find((item) => item.field === 'only_can_submit')?.componentProps
            ?.defaultValue || 1;

        if (!params.provider) params.provider = defaultProvider;

        // 仅显示可提交的
        if (typeof params.only_can_submit === 'undefined') {
          params.only_can_submit = defaultOnlyCanSubmit;
        }

        // 处理 demand_time 数据
        if (params.demand_time) {
          if (dayjs.isDayjs(params.demand_time)) {
            params.demand_time = params.demand_time.format('YYYY-MM');
          } else if (typeof params.demand_time === 'string') {
            params.demand_time = params.demand_time.substring(0, 7);
          }
        } else {
          // 如果是首次加载（未手动操作过），默认使用当月
          if (!endTimeManuallySet[bizType].value) {
            params.demand_time = dayjs().format('YYYY-MM');
          } else {
            // 用户主动清空了日期，删除该参数
            delete params.demand_time;
          }
        }

        // 处理区域和省份参数
        if (params.area) params.area_name = params.area;
        if (params.province) params.province_name = params.province;
        delete params.area;
        delete params.province;
        return params;
      },
      useSearchForm: true,
      showTableSetting: true,
      bordered: true,
      showIndexColumn: false,
      pagination: {},
      rowSelection: {
        type: 'checkbox',
        onChange: (_selectedRowKeys, rows) => {
          selectedRows[bizType] = rows;
          // 更新当前选中的设备类型：取第一条记录的设备类型，如果没有选中则重置为 null
          selectedDevType[bizType] = rows.length > 0 ? rows[0].dev_name : null;
        },
        getCheckboxProps: (record) => {
          // 根据业务类型判断禁用条件
          const isExpired = dayjs(record.end_time).isBefore(dayjs().startOf('day'));
          const noAvailableBw = record.gap_bw - record.locked_bw <= 0;
          // 汇聚额外检查 is_cover_diff_isp
          const isCoverDiffIsp = bizType === 'normal' && record.is_cover_diff_isp === 1;
          // 设备类型不匹配时禁用（如果已有选中的设备类型，则只允许勾选相同类型）
          const devTypeMismatch =
            selectedDevType[bizType] !== null && record.dev_name !== selectedDevType[bizType];

          return {
            disabled: isExpired || noAvailableBw || isCoverDiffIsp || devTypeMismatch,
          };
        },
      },
      rowKey: 'demand_id',
    });

    // 保存表格配置
    tableConfigs[bizType] = {
      register: register as any,
      getForm: getForm as any,
      reload: reload as any,
      clearSelectedRowKeys: clearSelectedRowKeys as any,
    };
  };

  // 创建两个表格的配置
  createTableConfig('normal');
  createTableConfig('specialLine');

  // 记录每个业务类型是否已初始化区域联动
  const areaSelectInited = {
    normal: ref(false),
    specialLine: ref(false),
  } as const;

  // 按需初始化指定业务类型的区域联动（仅初始化一次）
  async function initAreaFor(bizType: 'normal' | 'specialLine') {
    if (areaSelectInited[bizType].value) return;
    // 等待当前 Tab 对应的表格渲染完成，确保表单已就绪
    await nextTick();
    const form = await tableConfigs[bizType].getForm();

    const areaSelect = useAreaSelect({
      form,
      hasCity: false,
      fields: {
        area: 'area_name',
        province: 'province_name',
      },
    });
    await areaSelect.initAreaData();

    await form.updateSchema([
      {
        field: 'dev_type_names',
        componentProps: {
          options: devTypeOptions,
        },
      },
    ]);

    // 设置 demand_time 默认值为当月
    await form.setFieldsValue({
      demand_time: dayjs(),
    });

    // 标记已手动设置过，后续清空时可以删除参数
    endTimeManuallySet[bizType].value = true;

    areaSelectInited[bizType].value = true;
  }

  // 获取数据时更新树形数据
  onMounted(async () => {
    // 检查权限并设置默认激活的 tab
    checkPermissionsAndSetDefaultTab();

    await loadDevTypeData();

    try {
      // 仅初始化当前激活的 Tab，其他 Tab 在切换时再初始化
      await initAreaFor(activeTabKey.value);
    } catch (error) {
      console.error('获取数据失败:', error);
    }
  });

  // 监听 Tab 切换，懒加载初始化对应表格的区域联动
  watch(
    () => activeTabKey.value,
    async (key) => {
      if (key === 'normal' || key === 'specialLine') {
        await initAreaFor(key);
      }
    },
  );

  defineOptions({ name: 'KDemand' });

  // 需求占用处理函数
  const handleOccupy = (record, bizType) => {
    // 检查占用权限
    const occupyPermission =
      bizType === 'normal'
        ? 'business:k:demand:normal:occupy'
        : 'business:k:demand:specialline:occupy';

    if (!hasPermission(occupyPermission)) {
      notification.error({
        message: '权限不足',
        description: '您没有权限进行此操作',
      });
      return;
    }

    openModal(true, {
      record,
      isUpdate: true,
    });
  };

  // 需求锁定/解锁核心函数
  async function handleLock(record, bizType: 'normal' | 'specialLine') {
    try {
      await lockDemand({
        demand_id: record.demand_id,
        is_locked: !record.is_locked,
        locked_bw: record.is_locked ? undefined : record.locked_bw,
      });
      notification.success({
        message: record.is_locked ? '解锁成功' : '锁定成功',
      });
      tableConfigs[bizType].reload();
    } catch (e) {
      notification.error({
        message: record.is_locked ? '解锁失败' : '锁定失败',
        description: e.message,
      });
    }
  }

  // 需求锁定/解锁入口函数
  function handleLockDemand(record, bizType: 'normal' | 'specialLine') {
    // 检查锁定权限
    const lockPermission =
      bizType === 'normal' ? 'business:k:demand:normal:lock' : 'business:k:demand:specialline:lock';

    if (!hasPermission(lockPermission)) {
      notification.error({
        message: '权限不足',
        description: '您没有权限进行此操作',
      });
      return;
    }

    if (record.is_locked) {
      handleLock(record, bizType);
    } else {
      openLockModal(true, {
        record,
        is_locked: true,
      });
    }
  }

  // 成功处理后的刷新逻辑
  const handleSuccess = () => {
    // 刷新当前激活的表格
    tableConfigs[activeTabKey.value].reload();
  };

  // 批量占用成功后的刷新逻辑
  const handleBatchSuccess = () => {
    const bizType = activeTabKey.value;
    // 刷新当前激活的表格
    tableConfigs[bizType].reload();
    // 清空选中状态
    tableConfigs[bizType].clearSelectedRowKeys();
    selectedRows[bizType] = [];
    // 重置设备类型限制
    selectedDevType[bizType] = null;
  };

  // 获取选中行数量
  const getSelectedCount = (bizType: 'normal' | 'specialLine') => {
    return selectedRows[bizType]?.length || 0;
  };

  // 批量占用处理函数
  const handleBatchOccupy = (bizType: 'normal' | 'specialLine') => {
    const occupyPermission =
      bizType === 'normal'
        ? 'business:k:demand:normal:occupy'
        : 'business:k:demand:specialline:occupy';

    if (!hasPermission(occupyPermission)) {
      notification.error({
        message: '权限不足',
        description: '您没有权限进行此操作',
      });
      return;
    }

    const rows = selectedRows[bizType];
    if (!rows || rows.length === 0) {
      notification.warning({
        message: '提示',
        description: '请至少选择一条需求',
      });
      return;
    }

    // 根据当前选中的设备类型过滤行，只保留相同设备类型的需求
    const currentDevType = selectedDevType[bizType];
    const filteredRows = currentDevType
      ? rows.filter((row) => row.dev_name === currentDevType)
      : rows;

    if (filteredRows.length === 0) {
      notification.warning({
        message: '提示',
        description: '没有符合当前设备类型的需求',
      });
      return;
    }

    // 提取过滤后的 demand_id
    const demandIds = filteredRows.map((row) => row.demand_id);

    openBatchModal(true, {
      demandIds,
      bizType,
      provider: filteredRows[0].provider,
      devName: currentDevType, // 传递设备类型
    });
  };

  // 导出
  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
  });

  function handleExportData() {
    const currentBizType = activeTabKey.value;

    // 检查当前 tab 的导出权限
    const exportPermission =
      currentBizType === 'normal'
        ? 'business:k:demand:normal:export'
        : 'business:k:demand:specialline:export';

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

        // 处理 demand_time 数据
        if (value.demand_time) {
          if (dayjs.isDayjs(value.demand_time)) {
            value.demand_time = value.demand_time.format('YYYY-MM');
          } else if (typeof value.demand_time === 'string') {
            value.demand_time = value.demand_time.substring(0, 7);
          }
        } else {
          value.demand_time = '';
        }

        value.biz_type = currentBizType;

        // 从 searchFormSchema 中获取默认 provider 值
        const defaultProvider =
          searchFormSchema(currentBizType).find((item) => item.field === 'provider')?.componentProps
            ?.defaultValue || 'mf';
        if (!value.provider) value.provider = defaultProvider;

        nextTick(() => {
          data.exporting = true;
          data.exportButTitle = '导出中';
        });
        try {
          await ExportDemandList(value);
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

  async function ExportDemandList(value: Recordable) {
    const { success, fileName, error } = await downloadFileStream(Api.ExportDemandList, value);

    if (success) {
      notification.success({
        message: '导出成功',
        description: `文件名：${fileName}`,
        duration: null,
      });
      nextTick(() => {
        data.exporting = false;
        data.exportButTitle = '导出数据';
      });
    } else {
      notification.error({
        message: '导出失败',
        description: error || '未知错误',
      });
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
