<template>
  <div>
    <!-- 单独显示某个业务类型时，不显示 Tabs -->
    <template v-if="props.bizType !== 'all'">
      <!-- 汇聚表格 -->
      <BasicTable v-if="props.bizType === 'normal'" @register="tableConfigs.normal.register">
        <template #toolbar>
          <a-button
            type="primary"
            v-auth="'business:k:demandoccupy:normal:export'"
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
          <template v-if="column.key == 'status'">
            <Tooltip>
              <Tag
                v-if="OccupyStatusMap[record.status]"
                :color="OccupyStatusMap[record.status].color || 'default'"
                >{{ OccupyStatusMap[record.status].dictLabel }}
              </Tag>
            </Tooltip>
          </template>
        </template>
      </BasicTable>
      <!-- 专线表格 -->
      <BasicTable
        v-if="props.bizType === 'specialLine'"
        @register="tableConfigs.specialLine.register"
      >
        <template #toolbar>
          <a-button
            type="primary"
            v-auth="'business:k:demandoccupy:specialline:export'"
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
          <template v-if="column.key == 'status'">
            <Tooltip>
              <Tag
                v-if="OccupyStatusMap[record.status]"
                :color="OccupyStatusMap[record.status].color || 'default'"
                >{{ OccupyStatusMap[record.status].dictLabel }}
              </Tag>
            </Tooltip>
          </template>
        </template>
      </BasicTable>
    </template>

    <!-- 显示全部时使用 Tabs -->
    <Tabs v-else v-model:activeKey="activeTabKey" type="card">
      <TabPane tab="汇聚" key="normal" v-auth="'business:k:demandoccupy:normal'">
        <BasicTable @register="tableConfigs.normal.register">
          <template #toolbar>
            <a-button
              type="primary"
              v-auth="'business:k:demandoccupy:normal:export'"
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
            <template v-if="column.key == 'status'">
              <Tooltip>
                <Tag
                  v-if="OccupyStatusMap[record.status]"
                  :color="OccupyStatusMap[record.status].color || 'default'"
                  >{{ OccupyStatusMap[record.status].dictLabel }}
                </Tag>
              </Tooltip>
            </template>
          </template>
        </BasicTable>
      </TabPane>
      <TabPane tab="专线" key="specialLine" v-auth="'business:k:demandoccupy:specialline'">
        <BasicTable @register="tableConfigs.specialLine.register">
          <template #toolbar>
            <a-button
              type="primary"
              v-auth="'business:k:demandoccupy:specialline:export'"
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
            <template v-if="column.key == 'status'">
              <Tooltip>
                <Tag
                  v-if="OccupyStatusMap[record.status]"
                  :color="OccupyStatusMap[record.status].color || 'default'"
                  >{{ OccupyStatusMap[record.status].dictLabel }}
                </Tag>
              </Tooltip>
            </template>
          </template>
        </BasicTable>
      </TabPane>
    </Tabs>
  </div>
</template>

<script lang="ts" setup>
  import { BasicTable, useTable } from '@/components/Table';
  import { Tabs, TabPane, Modal, Tag, Tooltip } from 'ant-design-vue';
  import {
    DemandOccupyListColumns,
    searchOccupyListSchema,
    bizTypesMap,
    OccupyStatusMap,
  } from './data';
  import { Api, DemandOccupyList, GetDevTypeList } from '@/api/business/k';
  import { h, ref, nextTick, onMounted, reactive, watch } from 'vue';
  import dayjs from 'dayjs';
  import { useMessage } from '@/hooks/web/useMessage';
  import { ExclamationCircleOutlined } from '@ant-design/icons-vue';
  import { downloadFileStream } from '@/utils/download';
  import { useAreaSelect } from '@/utils/kAreaSelect';
  import { usePermission } from '@/hooks/web/usePermission';
  import { BizType } from '@/views/business/k/data';

  // 定义 props
  const props = withDefaults(
    defineProps<{
      bizType?: BizType | 'all';
    }>(),
    {
      bizType: 'all',
    },
  );

  // 权限检查
  const { hasPermission } = usePermission();

  // 当前激活的标签页
  const activeTabKey = ref<'normal' | 'specialLine'>('normal');

  // 检查权限并设置默认激活的 tab
  const checkPermissionsAndSetDefaultTab = () => {
    // 如果指定了具体的业务类型，直接设置
    if (props.bizType !== 'all') {
      activeTabKey.value = props.bizType;
      return;
    }

    const hasNormalPermission = hasPermission('business:k:demandoccupy:normal');
    const hasSpecialLinePermission = hasPermission('business:k:demandoccupy:specialline');

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

  const { notification } = useMessage();

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
    const title = bizType === 'normal' ? '汇聚需求占用列表' : '专线需求占用列表';

    // 从 searchFormSchema 中获取默认 provider 值
    const defaultProvider =
      searchOccupyListSchema(bizType).find((item) => item.field === 'provider')?.componentProps
        ?.defaultValue || 'mf';

    // 定义搜索表单 schema
    const searchSchema = [...searchOccupyListSchema(bizType)];

    const [register, { getForm, reload }] = useTable({
      title,
      api: async (params) => {
        const result = await DemandOccupyList(params);
        return result;
      },
      columns: DemandOccupyListColumns,
      formConfig: {
        labelWidth: 120,
        schemas: searchSchema,
        autoSubmitOnEnter: true,
        showAdvancedButton: true,
        autoAdvancedLine: 4,
      },
      beforeFetch(params) {
        params.biz_type = bizType;
        if (!params.provider) params.provider = defaultProvider;

        if (params.occupy_time_range) {
          params.start_time = dayjs(params.occupy_time_range[0]).format('YYYY-MM-DD HH:mm:ss');
          params.end_time = dayjs(params.occupy_time_range[1]).format('YYYY-MM-DD HH:mm:ss');
        }
        delete params.occupy_time_range;
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
      register,
      getForm,
      reload,
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

    try {
      const form = await tableConfigs[bizType].getForm();
      if (!form) return;

      // 区域联动筛选
      const { initAreaData } = useAreaSelect({
        form,
        fields: {
          area: 'area',
          province: 'province',
        },
      });
      await initAreaData();

      // 更新设备类型选项
      await form.updateSchema([
        {
          field: 'dev_name',
          componentProps: {
            options: devTypeOptions,
          },
        },
      ]);

      areaSelectInited[bizType].value = true;
    } catch (error) {
      console.error(`初始化${bizType}表格数据失败:`, error);
    }
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

  defineOptions({ name: 'KDemandOccupy' });

  // 导出数据
  const data = reactive({
    exporting: false,
    exportButTitle: '导出数据',
  });

  function handleExportData() {
    const currentBizType = props.bizType !== 'all' ? props.bizType : activeTabKey.value;

    // 检查当前 tab 的导出权限
    const exportPermission =
      currentBizType === 'normal'
        ? 'business:k:demandoccupy:normal:export'
        : 'business:k:demandoccupy:specialline:export';

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

        // 处理 occupy_time 数据
        if (value.occupy_time_range) {
          value.start_time = dayjs(value.occupy_time_range[0]).format('YYYY-MM-DD HH:mm:ss');
          value.end_time = dayjs(value.occupy_time_range[1]).format('YYYY-MM-DD HH:mm:ss');
        }
        delete value.occupy_time_range;
        value.biz_type = currentBizType;

        // 从 searchFormSchema 中获取默认 provider 值
        const defaultProvider =
          searchOccupyListSchema(currentBizType).find((item) => item.field === 'provider')
            ?.componentProps?.defaultValue || 'mf';
        if (!value.provider) value.provider = defaultProvider;

        nextTick(() => {
          data.exporting = true;
          data.exportButTitle = '导出中';
        });
        try {
          await ExportDemandOccupy(value);
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
  async function ExportDemandOccupy(value: Recordable) {
    const { success, fileName, error } = await downloadFileStream(Api.ExportDemandOccupy, value);

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
