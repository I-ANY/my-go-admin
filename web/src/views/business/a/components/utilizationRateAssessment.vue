<template>
  <div class="utilization-rate-assessment-container">
    <BasicTable @register="registerTable">
      <template #toolbar>
        <div v-if="hasPermission(APermissionCodeEnum.BUSINESS_A_UTILIZATION_RATE_ASSESSMENT_SCORE)">
          <Spin size="small" :spinning="spinning">
            <a-button
              type="success"
              @click="handleAssessmentScore"
              :disabled="!canUseAssessmentScore()"
              :class="{ 'disabled-btn': !canUseAssessmentScore() }"
            >
              考核打分
            </a-button>
          </Spin>
        </div>
        <div>
          <Spin size="small" :spinning="spinning">
            <a-button type="primary" @click="handleExport">导出 </a-button>
          </Spin>
        </div>
        <a-button style="margin-left: 8px" @click="handleResetData">重置</a-button>
      </template>
      <template #bodyCell="{ column, record, index }">
        <template v-if="column.key == 'roomNo'">
          <template v-if="calculateRowSpanByNode(safeGetDataSource(), index, 'roomNo') > 0">
            {{ record.roomNo }}
          </template>
        </template>
        <template v-if="column.key == 'roomType'">
          <template v-if="calculateRowSpanByNode(safeGetDataSource(), index, 'roomType') > 0">
            <Tag style="margin: 1px">{{ roomTypeMap[record.roomType] }}</Tag>
          </template>
        </template>
        <template v-if="column.key == 'localIsp'">
          <template v-if="calculateRowSpanByNode(safeGetDataSource(), index, 'localIsp') > 0">
            {{ record.localIsp }}
          </template>
        </template>
        <template v-if="column.key == 'location'">
          <template v-if="calculateRowSpanByNode(safeGetDataSource(), index, 'location') > 0">
            {{ record.location }}
          </template>
        </template>
        <template v-if="column.key == 'state'">
          <Tag :color="record.state === 1 ? 'green' : 'red'" style="margin: 1px">
            {{ stateMap[record.state] }}
          </Tag>
        </template>
        <template v-if="column.key == 'reportType'">
          <Tag style="margin: 1px">{{ reportTypeMap[record.reportType] }}</Tag>
        </template>
        <template v-if="column.key == 'bwUsageRateDay'">
          {{ formatPercentage(record.bwUsageRateDay) }}
        </template>
        <template v-if="column.key == 'bwUsageRateNight'">
          {{ formatPercentage(record.bwUsageRateNight) }}
        </template>
        <template v-if="column.key == 'bizs'">
          {{ record.bizs ? record.bizs.join(', ') : '-' }}
        </template>
        <template v-if="column.key == 'biz'">
          <template v-if="calculateRowSpanByNode(safeGetDataSource(), index, 'biz') > 0">
            {{ record.biz || '-' }}
          </template>
        </template>
        <template v-if="column.key == 'stateManual'">
          <template v-if="record.state === 1">
            <!-- 达标时直接显示，不可点击 -->
            <Tag :color="record.stateManual === 1 ? 'green' : 'red'" style="margin: 1px">
              {{ stateManualMap[record.stateManual] }}
            </Tag>
          </template>
          <template v-else>
            <!-- 使用Dropdown组件实现 -->
            <Dropdown :disabled="!canEditRecord(record)">
              <template #overlay>
                <Menu>
                  <Menu.Item
                    v-for="item in stateManualOptions"
                    :key="item.value"
                    :disabled="!canEditRecord(record)"
                    @click="handleStateManualChange(item.value, record)"
                  >
                    {{ item.label }}
                  </Menu.Item>
                </Menu>
              </template>
              <a
                class="ant-dropdown-link"
                style="display: inline-block; min-width: 80px"
                @click.stop
              >
                {{ stateManualMap[record.stateManual] || '请选择' }}
                <DownOutlined />
              </a>
            </Dropdown>
          </template>
        </template>
        <template v-if="column.key == 'planType'">
          <template v-if="record.state === 1 || record.state === 3">
            <!-- 达标时或暂无要求时直接显示 "-"，不可点击 -->
            <span style="display: inline-block; min-width: 100px; color: #999">-</span>
          </template>
          <template v-else-if="record.stateManual === 2">
            <!-- 只有达标确认==否时才显示可编辑的 Dropdown -->
            <Dropdown :disabled="!canEditRecord(record)">
              <template #overlay>
                <Menu>
                  <Menu.Item
                    v-for="item in planTypeOptions"
                    :key="item.value"
                    :disabled="!canEditRecord(record)"
                    @click="handleplanTypeChange(item.value, record)"
                  >
                    {{ item.label }}
                  </Menu.Item>
                </Menu>
              </template>
              <a
                class="ant-dropdown-link"
                :class="{ 'ant-dropdown-link-disabled': !canEditRecord(record) }"
                style="display: inline-block; min-width: 100px"
                @click.stop
              >
                {{ planTypeMap[record.planType] || '请选择' }}
                <DownOutlined />
              </a>
            </Dropdown>
          </template>
          <template v-else>
            <!-- 达标确认!=否时显示 "-" -->
            <span style="display: inline-block; min-width: 100px; color: #999">-</span>
          </template>
        </template>
        <template v-if="column.key == 'nonComplianceReason'">
          <template v-if="record.state === 1 || record.state === 3 || record.stateManual !== 2">
            <!-- 达标/暂无要求或达标确认非“否”时不可编辑 -->
            <span style="display: inline-block; min-width: 140px; color: #999">-</span>
          </template>
          <template v-else>
            <Select
              mode="multiple"
              style="min-width: 120px"
              :value="normalizeNonComplianceReason(record.nonComplianceReason)"
              :options="nonComplianceReasonOptions"
              placeholder="请选择"
              allowClear
              :disabled="!canEditRecord(record)"
              @change="(vals) => handleNonComplianceReasonChange(vals, record)"
            />
          </template>
        </template>
        <template v-if="column.key == 'origin'">
          <Tag :color="record.origin === '自建' ? 'green' : 'blue'">
            {{ record.origin || '自建' }}
          </Tag>
        </template>
        <template v-if="column.key == 'billingType'">
          <Tag :color="getBillingTypeColor(record.billingType)">
            {{ record.billingType || '-' }}
          </Tag>
        </template>

        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                label: '查看历史',
                tooltip: '查看历史',
                type: 'primary',
                size: 'small',
                onClick: handleViewHistory.bind(null, record),
              },
            ]"
          />
        </template>
      </template>

      <template #tableTop>
        <div style="margin-bottom: 8px; padding: 8px; background-color: white">
          <div
            >1、<span style="color: #4ebef1">考核标准</span
            >：IDC&ACDN类型：日95&晚高峰95利用率均<span style="color: #4ebef1">>99%</span
            >且晚高峰的点数<span style="color: #4ebef1">>=36</span
            >;MCDN类型：日95&晚高峰95利用率均<span style="color: #4ebef1">>95%</span
            >且晚高峰的点数<span style="color: #4ebef1">>=36</span
            >;考核数据依据对应节点考核指标进行展示；</div
          >
          <div
            >2、<span style="color: #4ebef1">晚高峰达标点数</span
            >：晚高峰利用率大于考核标准的点数;</div
          >
        </div>
      </template>
    </BasicTable>

    <!-- 历史利用率考核模态框 -->
    <HistoryAssessmentModal @register="registerHistoryModal" />
  </div>
</template>
<script lang="ts" setup>
  import { BasicTable, useTable, BasicColumn, FormSchema, TableAction } from '@/components/Table';
  import {
    getOwners,
    getComposite,
    getExportData,
    editReports,
    getBizs,
    getAssessment,
    assessmentAutoScore,
  } from '@/api/business/a';
  import dayjs from 'dayjs';
  import { h, reactive, ref, onUnmounted, onMounted } from 'vue';
  import { jsonToSheetXlsx, ExportModalResult } from '@/components/Excel';
  import { Tag, Spin, message, Dropdown, Menu, Tooltip, Modal, Select } from 'ant-design-vue';
  import DownOutlined from '@ant-design/icons-vue/DownOutlined';
  import { useModal } from '@/components/Modal';
  import HistoryAssessmentModal from './HistoryAssessmentModal.vue';
  import { usePermission } from '@/hooks/web/usePermission';
  import { APermissionCodeEnum } from '@/enums/permissionCodeEnum';

  const serache = reactive({});

  // 历史考核模态框
  const [registerHistoryModal, { openModal: openHistoryModal }] = useModal();

  // 权限检查
  const { hasPermission } = usePermission();

  // 自动刷新定时器
  const autoRefreshTimer = ref<number | null>(null);
  const spinning = ref<boolean>(false);

  // 格式化百分比函数，保留两位小数
  const formatPercentage = (value: number | null | undefined): string => {
    if (value == null || value === undefined) return '-';
    return `${(value * 100).toFixed(2)}%`;
  };

  // 获取计费类型颜色
  const getBillingTypeColor = (type: string): string => {
    const colorMap = {
      日95: 'blue',
      月95: 'green',
      买断: 'orange',
    };
    return colorMap[type] || 'default';
  };

  // 获取前一天
  const yesterday = dayjs().subtract(1, 'day').format('YYYY-MM-DD');
  // 获取当天日期
  const today = dayjs().format('YYYY-MM-DD');

  // 检查是否可以编辑：当天数据所有人都可编辑，历史数据需要特定权限
  const canEditRecord = (record: any): boolean => {
    const recordDate = record.date || record.calculateDate;
    if (!recordDate) return false;
    // 如果是当天数据，所有人都可以编辑
    if (recordDate === today) return true;
    // 如果是历史数据，需要检查权限
    return hasPermission(APermissionCodeEnum.BUSINESS_A_UTILIZATION_RATE_ASSESSMENT_HISTORY_EDIT);
  };

  const normalizeSelectedDate = (value: string | dayjs.Dayjs | undefined): string | undefined => {
    if (!value) return undefined;
    if (typeof value === 'string') return value;
    if (dayjs.isDayjs(value)) {
      return value.format('YYYY-MM-DD');
    }
    return String(value);
  };

  // 检查考核打分按钮是否可用：昨天可用，其他日期需要历史编辑权限
  const canUseAssessmentScore = (): boolean => {
    const rawDate = serache['date'];
    const selectedDate = normalizeSelectedDate(rawDate) || yesterday;
    // 如果是昨天，所有人可用
    if (selectedDate === yesterday) return true;
    // 如果是其他日期，需要历史编辑权限
    return hasPermission(APermissionCodeEnum.BUSINESS_A_UTILIZATION_RATE_ASSESSMENT_HISTORY_EDIT);
  };

  const roomTypeMap = {
    1: 'IDC',
    2: 'ACDN',
    3: 'MCDN',
  };
  const stateMap = {
    0: '未知',
    1: '是',
    2: '否',
    3: '暂无要求',
  };
  const reportTypeMap = {
    1: '节点',
    2: '保底',
    3: '削峰',
  };
  const stateManualMap = {
    1: '是',
    2: '否',
    3: '暂无要求',
    4: '已裁撤',
    5: '不计费',
    6: '其他',
  };
  // 评分映射
  const planTypeMap = {
    1: '不扣分',
    2: '有方案',
    3: '无方案',
    4: '临时特批',
    5: '其他',
  };
  const stateOptions = [
    { label: '是', value: 1 },
    { label: '否', value: 2 },
    { label: '暂无要求', value: 3 },
    { label: '未知', value: 0 },
  ];
  // 达标确认选项
  const stateManualOptions = [
    { label: '是', value: 1 },
    { label: '否', value: 2 },
    { label: '暂无要求', value: 3 },
    { label: '已裁撤', value: 4 },
    { label: '不计费', value: 5 },
    { label: '其他', value: 6 },
  ];

  // 方案类型选项
  const planTypeOptions = [
    { label: '不扣分', value: 1 },
    { label: '有方案', value: 2 },
    { label: '无方案', value: 3 },
    { label: '临时特批', value: 4 },
    { label: '其他', value: 5 },
    { label: '未知', value: 0 },
  ];

  // 是否仅异网节点
  const isExternalOptions = [
    { label: '是', value: 'true' },
    { label: '否', value: 'false' },
  ];

  // 未达标原因选项
  const nonComplianceReasonOptions = [
    { label: '补发硬件', value: '补发硬件' },
    { label: '硬件问题', value: '硬件问题' },
    { label: '网络问题', value: '网络问题' },
    { label: '缓存问题', value: '缓存问题' },
    { label: '性能问题', value: '性能问题' },
    { label: '起量较晚', value: '起量较晚' },
    { label: '业务测试', value: '业务测试' },
    { label: '其他问题', value: '其他问题' },
  ];

  // 将后端存储值转换为多选数组
  const normalizeNonComplianceReason = (value: any): string[] => {
    if (!value) return [];
    if (Array.isArray(value)) return value.filter((v) => v && String(v).trim());
    if (typeof value === 'string') {
      return value
        .split(',')
        .map((v) => v.trim())
        .filter((v) => v);
    }
    return [];
  };

  // 处理达标确认变更
  const handleStateManualChange = async (value: number, record: any) => {
    // 检查是否可以编辑
    if (!canEditRecord(record)) {
      message.warning('只有当天数据或拥有历史数据编辑权限的用户才能修改');
      return;
    }
    // 保存原值，用于失败时恢复
    const originalValue = record.stateManual;
    console.log('originalValue: ', value);

    try {
      const params: Recordable = {
        id: record.id,
        stateManual: value,
      };

      // 如果达标确认不是"否"(2)，需要清除方案类型、方案内容、评分
      if (value !== 2) {
        params.planType = null;
        params.planContent = '';
        params.score = null;
      }

      await editReports(params);
      message.success('达标手动确认更新成功');

      // 更新本地数据，避免重新加载整个表格
      record.stateManual = value;
      if (value !== 2) {
        record.planType = null;
        record.planContent = '';
        record.score = null;
      }
    } catch (error) {
      message.error('达标手动确认更新失败，请重试');
      // 失败时恢复原值
      record.stateManual = originalValue;
      // 强制重新渲染表格
      // reload();
    }
  };

  // 处理方案类型变更
  const handleplanTypeChange = async (value: number, record: any) => {
    // 检查是否可以编辑
    if (!canEditRecord(record)) {
      message.warning('只有当天数据或拥有历史数据编辑权限的用户才能修改');
      return;
    }
    // 保存原值，用于失败时恢复
    const originalValue = record.planType;

    try {
      const params: Recordable = {
        id: record.id,
        planType: value,
      };
      await editReports(params);
      message.success('方案类型更新成功');
      // 更新本地数据，避免重新加载整个表格
      record.planType = value;
    } catch (error) {
      message.error('方案类型更新失败，请重试');
      // 失败时恢复原值
      record.planType = originalValue;
      // 强制重新渲染表格
      // reload();
    }
  };
  // 处理未达标原因变更（支持多选）
  const handleNonComplianceReasonChange = async (values: any, record: any) => {
    // 检查是否可以编辑
    if (!canEditRecord(record)) {
      message.warning('只有当天数据或拥有历史数据编辑权限的用户才能修改');
      return;
    }
    try {
      const value = Array.isArray(values)
        ? values.filter((v) => v && String(v).trim()).join(',')
        : '';
      const params: Recordable = {
        id: record.id,
        nonComplianceReason: value,
      };
      await editReports(params);
      message.success('未达标原因修改成功');
      // 更新本地数据，避免重新加载整个表格
      record.nonComplianceReason = value;
    } catch (error) {
      console.error('未达标原因修改失败:', error);
      message.error('未达标原因修改失败');
      // reload();
    }
  };

  // 处理考核业务变更（专线多业务节点）
  const handleAssessmentBusinessChange = async (values: string[], record: any) => {
    // 检查是否可以编辑
    if (!canEditRecord(record)) {
      message.warning('只有当天数据或拥有历史数据编辑权限的用户才能修改');
      // 回滚选择
      return;
    }
    try {
      const params: Recordable = {
        id: record.id,
        assessmentBusiness: values || [],
      };
      await editReports(params);
      message.success('考核业务更新成功');
      // 更新本地数据
      record.assessmentBusiness = values;
    } catch (error: any) {
      console.error('考核业务更新失败:', error);
      message.error(error?.message || '考核业务更新失败');
    }
  };

  // 查看历史利用率考核
  const handleViewHistory = (record: any) => {
    openHistoryModal(true, {
      record,
    });
  };

  // 存储当前的评分类型，用于切换API
  let currentAssessmentType = ref<number | null>(null);

  // 考核打分
  const handleAssessmentScore = () => {
    const date = serache['date'] || yesterday;
    Modal.confirm({
      title: '确认操作',
      content: `确定要对 ${date} 的数据进行考核自动打分吗？`,
      okText: '确认',
      cancelText: '取消',
      onOk: async () => {
        try {
          spinning.value = true;
          await assessmentAutoScore({ date });
          message.success('考核打分成功');
          reload();
        } catch (error: any) {
          message.error(error?.message || '考核打分失败');
        } finally {
          spinning.value = false;
        }
      },
    });
  };

  // 统一的API调用函数
  const unifiedApiCall = async (params: Recordable) => {
    if (currentAssessmentType.value) {
      // 如果设置了评分类型，调用评分接口
      const assessmentParams: Recordable = {
        ...params,
        date: params.date || yesterday,
        assessmentType: currentAssessmentType.value,
        // pageSize: 1000, // 确保获取所有评分数据
      };
      return await getAssessment(assessmentParams);
    } else {
      // 否则调用原始接口
      return await getComposite(params);
    }
  };

  // 重置数据到原始列表
  const handleResetData = () => {
    currentAssessmentType.value = null;
    reload();
    message.success('已重置为原始数据');
  };

  const handleExport = () => {
    const timestamp = dayjs().format('YYYYMMDDHHmm');
    defaultHeader({
      filename: `利用率考核导出数据_${timestamp}.xlsx`,
      bookType: 'xlsx',
    });
  };

  // 停止自动刷新
  const stopAutoRefresh = () => {
    if (autoRefreshTimer.value) {
      clearInterval(autoRefreshTimer.value);
      autoRefreshTimer.value = null;
    }
  };

  // 组件卸载时清除定时器
  onUnmounted(() => {
    stopAutoRefresh();
  });
  onMounted(() => {
    serache['date'] = yesterday;
  });
  function defaultHeader({ filename, bookType }: ExportModalResult) {
    spinning.value = true;
    getExportData({
      ...serache,
    })
      .then((resp) => {
        spinning.value = false;
        let data = resp;
        console.log(data);
        let result: { [key: string]: any }[] = [];
        for (let i = 0; i < data.length; i++) {
          let map = {};
          map['计算日期'] = data[i]['计算日期'];
          map['机房类型'] = data[i]['机房类型'];
          map['机房归属'] = data[i]['机房归属'];
          map['节点编号'] = data[i]['节点编号'];
          map['运营商'] = data[i]['运营商'];
          map['所在地'] = data[i]['所在地'];
          map['统计类型'] = data[i]['统计类型'];
          map['计费方式'] = data[i]['计费方式'];
          map['是否仅异网节点'] = data[i]['是否仅异网节点'];
          map['日95利用率'] = data[i]['日95利用率'];
          map['晚高峰利用率'] = data[i]['晚高峰利用率'];
          map['达标点数'] = data[i]['达标点数'];
          map['主线业务'] = data[i]['主线业务'];
          map['考核业务'] = data[i]['考核业务'];
          map['未达标原因'] = data[i]['未达标原因'];
          map['带宽占比'] = data[i]['带宽占比'];
          map['方案类型'] = data[i]['方案类型'];
          map['方案内容'] = data[i]['方案内容'];
          map['跨省占比'] = data[i]['跨省占比'];
          map['评分'] = data[i]['评分'];
          map['是否达标'] = data[i]['是否达标'];
          map['是否达标(手动确认)'] = data[i]['是否达标(手动确认)'];
          map['备注'] = data[i]['备注'];
          map['节点备注'] = data[i]['节点备注'];
          result.push(map);
        }
        jsonToSheetXlsx({
          data: result,
          filename: filename,
          write2excelOpts: {
            bookType,
          },
        });
      })
      .catch(() => {
        spinning.value = false;
      });
  }
  const columns: BasicColumn[] = [
    {
      title: '计算日期',
      dataIndex: 'date',
      width: 100,
      resizable: true,
      // fixed: 'left',
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'date');
        return {
          rowSpan: rowSpan,
        };
      },
      // defaultHidden: true,
    },
    {
      title: '机房类型',
      dataIndex: 'roomType',
      width: 80,
      resizable: true,
      fixed: 'left',
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'roomType');
        return {
          rowSpan: rowSpan,
        };
      },
    },
    {
      title: '节点编号',
      dataIndex: 'roomNo',
      width: 100,
      resizable: true,
      fixed: 'left',
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'roomNo');
        return {
          rowSpan: rowSpan,
        };
      },
    },
    {
      title: '计费方式',
      dataIndex: 'billingType',
      width: 100,
      resizable: true,
      fixed: 'left',
    },
    {
      title: '运营商',
      dataIndex: 'localIsp',
      width: 100,
      resizable: true,
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'localIsp');
        return {
          rowSpan: rowSpan,
        };
      },
      defaultHidden: true,
    },
    {
      title: '所在地',
      dataIndex: 'location',
      width: 100,
      resizable: true,
      defaultHidden: true,
    },
    {
      title: '是否仅异网节点',
      dataIndex: 'isExternalOnly',
      width: 80,
      customRender: ({ record }) => {
        if (record.isExternalOnly === null) return '-';
        return record.isExternalOnly === true ? '是' : '否';
      },
    },
    {
      title: '是否特批',
      dataIndex: 'isSpecialApproval',
      width: 80,
      customRender: ({ record }) => {
        if (record.isSpecialApproved === null) return '-';
        return record.isSpecialApproval === true ? '是' : '否';
      },
    },
    {
      title: '统计类型',
      dataIndex: 'reportType',
      width: 100,
      resizable: true,
      helpMessage: '统计类型：该字段按照利用率考核维度去统计对应节点考核数据的统计类型',
      customRender: ({ record }) => {
        if (record.reportType === 2) {
          return '保底';
        } else if (record.reportType === 3) {
          return '削峰';
        } else if (record.reportType === 1) {
          return '机房总览';
        }
      },
      defaultHidden: true,
    },
    {
      title: '日95利用率',
      dataIndex: 'bwUsageRateDay',
      width: 100,
      sorter: true,
      resizable: true,
    },
    {
      title: '晚高峰利用率',
      dataIndex: 'bwUsageRateNight',
      sorter: true,
      width: 100,
      resizable: true,
    },
    {
      title: '达标点数',
      dataIndex: 'nightPointNum',
      sorter: true,
      width: 100,
      resizable: true,
    },
    {
      title: '主线业务',
      dataIndex: 'biz',
      width: 100,
      sorter: true,
      resizable: true,
      customRender: ({ record }) => {
        const obj = {
          1: '机房总览',
          2: '保底',
          3: '削峰',
        };
        if (record.reportType === 1) {
          return record.biz;
        } else {
          return record.biz + '-' + record.localIsp + '-' + obj[record.reportType];
        }
      },
    },
    {
      title: '考核业务',
      dataIndex: 'assessmentBusiness',
      key: 'assessmentBusiness',
      width: 150,
      sorter: true,
      resizable: true,
      customRender: ({ record }) => {
        // 多业务节点：显示待选考核业务的多选下拉框
        if (record.candidateAssessmentBiz && record.candidateAssessmentBiz.length > 0) {
          console.log(record.candidateAssessmentBiz);
          const options = record.candidateAssessmentBiz.map((v: string) => ({
            label: v,
            value: v,
          }));
          return h(Select, {
            mode: 'multiple',
            value: record.assessmentBusiness || [],
            options: options,
            placeholder: '请选择',
            allowClear: true,
            style: { minWidth: '120px' },
            disabled: !canEditRecord(record),
            onChange: (vals: any) => handleAssessmentBusinessChange(vals, record),
            // maxTagCount: 2,
          });
        }
        // MCDN/单业务节点：显示已同步的考核业务
        const bizDisplay = record.assessmentBusiness ? record.assessmentBusiness.join(', ') : '-';
        return h(
          'span',
          {
            style: {
              color:
                record.assessmentBusiness && record.assessmentBusiness.length > 0
                  ? '#1890ff'
                  : '#999',
            },
          },
          bizDisplay,
        );
      },
    },
    {
      title: '未达标原因',
      dataIndex: 'nonComplianceReason',
      width: 150,
      sorter: true,
      resizable: true,
    },
    {
      title: '带宽占比',
      dataIndex: 'bwRatio',
      sorter: true,
      width: 100,
      resizable: true,
      customRender: ({ record }) => {
        if (record.bwRatio === -1) return '-';
        return `${(record.bwRatio * 100).toFixed(2)}%`;
      },
      defaultHidden: true,
    },
    {
      title: '跨省占比',
      dataIndex: 'trZoneRatio',
      sorter: true,
      width: 150,
      resizable: true,
      edit: true,
      editComponent: 'InputTextArea',
      editComponentProps: {
        autoSize: {
          minRows: 1,
          maxRows: 3,
        },
        submitOnEnter: false,
      },
      editRender: ({ record }) => {
        if (record.trZoneRatio === -1) return '-';
        if (!record.trZoneRatio) return '';

        // 将文本按换行符分割
        const lines = String(record.trZoneRatio).split('\n');
        const children: Array<any> = [];

        // 为每一行创建文本节点，行之间插入 <br> 元素
        lines.forEach((line, index) => {
          if (index > 0) {
            children.push(h('br'));
          }
          children.push(line);
        });

        return h('span', children);
      },
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'roomNo');
        return {
          rowSpan: rowSpan,
        };
      },
    },
    {
      title: '节点备注',
      dataIndex: 'nodeRemark',
      width: 200,
      resizable: true,
      edit: true,
      fixed: 'left',
      editComponent: 'InputTextArea',
      editComponentProps: {
        style: {
          textAlign: 'left',
          height: '100%',
        },
        autoSize: {
          minRows: 1,
          maxRows: 3,
        },
        placeholder: '请输入节点备注',
        submitOnEnter: false,
      },
      customRender: () => {},
      editRender: ({ record }) => {
        const text = record.nodeRemark || '';
        const isLongText = text.length > 60;
        const displayText = isLongText ? text.substring(0, 60) + '...' : text;
        return h(
          Tooltip,
          {
            title: text,
            placement: 'topLeft',
          },
          {
            default: () =>
              h(
                'div',
                {
                  style: {
                    wordWrap: 'break-word',
                    whiteSpace: 'pre-wrap',
                    overflowWrap: 'break-word',
                    textAlign: 'left',
                  },
                },
                displayText,
              ),
          },
        );
      },
      customCell: (_record: any, index: number | undefined) => {
        if (index === undefined) return {};
        const dataSource = safeGetDataSource();
        const rowSpan = calculateRowSpanByNode(dataSource, index, 'roomNo');
        return {
          rowSpan: rowSpan,
        };
      },
      editRule: async (text, record) => {
        if (text.length > 200) {
          return '不能超过200字符';
        }
        if (!canEditRecord(record)) {
          return '历史数据无权限编辑';
        }
        return '';
      },
    },
    {
      title: '机房归属',
      dataIndex: 'origin',
      width: 100,
      resizable: true,
      customRender: ({ record }) => {
        const obj = {
          1: '自建',
          2: '招募',
        };
        return obj[record.origin] || '-';
      },
      defaultHidden: true,
    },
    {
      title: '方案类型',
      dataIndex: 'planType',
      width: 120,
      resizable: true,
    },
    {
      title: '方案内容',
      dataIndex: 'planContent',
      width: 200,
      resizable: true,
      edit: true,
      editComponent: 'InputTextArea',
      editComponentProps: {
        autoSize: {
          minRows: 1,
          maxRows: 4,
        },
        submitOnEnter: false,
      },
      editRender: ({ record }) => {
        // 如果是否达标为"是"或暂无要求，统一显示 "-"
        if (record.state === 1 || record.state === 3) {
          return '-';
        }
        // 如果达标确认!=否，显示 "-"
        if (record.stateManual !== 2) {
          return '-';
        }
        // if (!record.planContent) return '';

        // 将文本按换行符分割
        const lines = String(record.planContent).split('\n');
        const children: Array<any> = [];

        // 为每一行创建文本节点，行之间插入 <br> 元素
        lines.forEach((line, index) => {
          if (index > 0) {
            children.push(h('br'));
          }
          children.push(line);
        });

        // return h('span', children);
        return h(
          Tooltip,
          {
            title: children,
            placement: 'topLeft',
          },
          {
            default: () =>
              h(
                'div',
                {
                  style: {
                    wordWrap: 'break-word',
                    whiteSpace: 'pre-wrap',
                    overflowWrap: 'break-word',
                    textAlign: 'left',
                  },
                },
                record.planContent,
              ),
          },
        );
      },
      editRule: async (text, record) => {
        // 如果是否达标为"是"或暂无要求，不允许编辑方案内容
        if (record && (record.state === 1 || record.state === 3)) {
          return '达标或暂无要求记录的方案内容不可编辑';
        }
        // 如果达标确认!=否，不允许编辑
        if (record && record.stateManual !== 2) {
          return '只有达标确认为"否"时才能编辑方案内容';
        }
        if (!canEditRecord(record)) {
          return '历史数据无权限编辑';
        }
        // if (text && text.length > 200) {
        //   return '不能超过200字符';
        // }
        return '';
      },
    },

    {
      title: '备注',
      dataIndex: 'describe',
      width: 150,
      resizable: true,
      edit: true,
      editComponent: 'InputTextArea',
      editComponentProps: {
        autoSize: {
          minRows: 1,
          maxRows: 3,
        },
        submitOnEnter: false,
      },
      editRule: async (text, record) => {
        if (text.length > 200) {
          return '不能超过200字符';
        }
        if (!canEditRecord(record)) {
          return '历史数据无权限编辑';
        }
        return '';
      },
    },
    {
      title: '操作人',
      dataIndex: 'updateUser',
      width: 100,
      resizable: true,
    },
    {
      title: '操作时间',
      dataIndex: 'updatedAt',
      width: 180,
      resizable: true,
      customRender: ({ record }) => {
        return dayjs(new Date(record.updatedAt)).format('YYYY-MM-DD HH:mm:ss');
      },
    },
    {
      dataIndex: 'state',
      title: '是否达标',
      width: 100,
      resizable: true,
      fixed: 'right',
      customRender: ({ record }) => {
        const color = record.state === 1 ? 'green' : record.state === 2 ? 'red' : 'default';
        return h(Tag, { color }, () => stateMap[record.state] || '未知');
      },
    },
    {
      title: '达标确认',
      dataIndex: 'stateManual',
      width: 100,
      resizable: true,
      fixed: 'right',
    },
    {
      title: '评分',
      dataIndex: 'score',
      width: 120,
      resizable: true,
      sorter: true,
      fixed: 'right',
    },
  ];
  const searchFormSchema: FormSchema[] = [
    {
      field: 'date',
      label: '计算日期',
      component: 'DatePicker',
      componentProps: {
        placeholder: '请选择日期',
        format: 'YYYY-MM-DD',
        valueFormat: 'YYYY-MM-DD',
        disabledDate: (current: any) => {
          // 不能选择未来的日期
          return current && current > new Date();
        },
        defaultValue: yesterday,
        onChange: (v) => {
          console.log(v);
          serache['date'] = v;
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'roomNo',
      label: '节点编号',
      component: 'ApiSelect',
      colProps: { span: 6 },
      componentProps: {
        options: [],
        mode: 'multiple',
        showSearch: true,
        placeholder: '请选择节点编号',
        api: async () => {
          const data = await getOwners();
          return data.map((item: any) => ({
            label: item.name,
            value: item.id,
          }));
        },
        onChange: (v) => {
          serache['roomNo'] = v;
        },
      },
    },
    {
      field: 'roomType',
      label: '机房类型',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        placeholder: '请选择机房类型',
        options: [
          { label: 'IDC', value: 1 },
          { label: 'ACDN', value: 2 },
          { label: 'MCDN', value: 3 },
        ],
        onChange: (v) => {
          serache['roomType'] = v;
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'origin',
      label: '机房归属',
      component: 'Select',
      componentProps: {
        placeholder: '请选择机房归属',
        options: [
          { label: '自建', value: 1 },
          { label: '招募', value: 2 },
          { label: '未知', value: 0 },
        ],
        onChange: (v) => {
          serache['origin'] = v;
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'billingType',
      label: '计费方式',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        placeholder: '请选择计费方式',
        options: [
          { label: '日95', value: '日95' },
          { label: '月95', value: '月95' },
          { label: '买断', value: '买断' },
        ],
        onChange: (v) => {
          serache['billingType'] = v;
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'localIsp',
      label: '运营商',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        placeholder: '请选择运营商',
        options: [
          { label: '电信', value: '电信' },
          { label: '移动', value: '移动' },
          { label: '联通', value: '联通' },
        ],
        onChange: (v) => {
          serache['localIsp'] = v;
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'reportType',
      label: '统计类型',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        placeholder: '请选择统计类型',
        options: [
          { label: '节点', value: 1 },
          { label: '保底', value: 2 },
          { label: '削峰', value: 3 },
        ],
        onChange: (v) => {
          serache['reportType'] = v;
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'biz',
      label: '业务名称',
      component: 'ApiSelect',
      componentProps: {
        options: [],
        mode: 'multiple',
        showSearch: true,
        placeholder: '请选择业务名称',
        api: async () => {
          // 获取当前选中的日期，如果没有则使用昨天作为默认值
          const selectedDate = serache['date'] || yesterday;
          const data = await getBizs({ date: selectedDate });
          console.log('data', data);
          return data.bizs.map((item: any) => ({
            label: item,
            value: item,
          }));
        },
        onChange: (v) => {
          serache['biz'] = v;
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'planType',
      label: '方案类型',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        placeholder: '请选择方案类型',
        options: planTypeOptions,
        onChange: (v) => {
          serache['planType'] = v;
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'state',
      label: '是否达标',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        placeholder: '是否达标',
        options: stateOptions,
        onChange: (v) => {
          serache['state'] = v;
        },
      },
      colProps: { span: 6 },
    },
    {
      field: 'stateManual',
      label: '达标确认',
      component: 'Select',
      componentProps: {
        mode: 'multiple',
        placeholder: '请选择达标确认方式',
        options: [...stateManualOptions, { label: '未知', value: 0 }],
        onChange: (v) => {
          serache['stateManual'] = v;
        },
      },
      colProps: { span: 6 },
    },
    // {
    //   field: 'score',
    //   label: '评分',
    //   component: 'Select',
    //   componentProps: {
    //     mode: 'multiple',
    //     placeholder: '请选择评分',
    //     options: scoreOptions,
    //     onChange: (v) => {
    //       serache['score'] = v;
    //     },
    //   },
    //   colProps: { span: 6 },
    // },
    {
      field: 'isExternalOnly',
      label: '是否仅异网节点',
      component: 'Select',
      componentProps: {
        placeholder: '请选择是否仅异网节点',
        options: isExternalOptions,
        onChange: (v) => {
          serache['isExternalOnly'] = v;
        },
      },
      colProps: { span: 6 },
    },
  ];
  // 获取同一节点编号的行范围
  const getNodeRowRange = (dataSource: any[], index: number) => {
    if (!dataSource || dataSource.length === 0) return { start: index, end: index };

    const currentRoomNo = dataSource[index].roomNo;
    let start = index;
    let end = index;

    // 向上查找相同节点编号的起始位置
    while (start > 0 && dataSource[start - 1].roomNo === currentRoomNo) {
      start--;
    }

    // 向下查找相同节点编号的结束位置
    while (end < dataSource.length - 1 && dataSource[end + 1].roomNo === currentRoomNo) {
      end++;
    }

    return { start, end };
  };

  // 计算基于节点编号的行合并函数
  const calculateRowSpanByNode = (
    dataSource: any[],
    index: number,
    key: string,
    compareFn?: (current: any, next: any) => boolean,
  ) => {
    try {
      if (!dataSource || dataSource.length === 0 || index < 0 || index >= dataSource.length) {
        return 1;
      }

      const { start, end } = getNodeRowRange(dataSource, index);
      const currentValue = dataSource[index][key];
      const currentRecord = dataSource[index];

      // 在同一节点编号范围内，查找连续相同值的范围
      let spanStart = index;
      let spanEnd = index;

      // 向上查找相同值（但不超出节点范围）
      while (spanStart > start) {
        const prevRecord = dataSource[spanStart - 1];
        const prevValue = prevRecord[key];

        // 如果提供了自定义比较函数，则使用它
        if (compareFn) {
          if (!compareFn(prevRecord, currentRecord)) break;
        } else if (prevValue !== currentValue) {
          break;
        }

        spanStart--;
      }

      // 向下查找相同值（但不超出节点范围）
      while (spanEnd < end) {
        const nextRecord = dataSource[spanEnd + 1];
        const nextValue = nextRecord[key];

        // 如果提供了自定义比较函数，则使用它
        if (compareFn) {
          if (!compareFn(nextRecord, currentRecord)) break;
        } else if (nextValue !== currentValue) {
          break;
        }

        spanEnd++;
      }

      const rowSpan = spanEnd - spanStart + 1;

      // 如果当前行是第一个相同值的行，返回总的行数；否则返回0（不显示）
      if (spanStart === index) {
        return rowSpan;
      }
      return 0;
    } catch (error) {
      console.warn('calculateRowSpanByNode error:', error);
      return 1;
    }
  };

  const beforeEditSubmit = async (data: { record: any; index: number; key: any; value: any }) => {
    const { record, key, value } = data;
    // 检查是否可以编辑
    if (!canEditRecord(record)) {
      throw new Error('只有当天数据或拥有历史数据编辑权限的用户才能修改');
    }
    let params: any = { id: record.id };

    // 根据编辑的字段设置对应的参数
    if (key === 'describe') {
      params.describe = value;
    } else if (key === 'planContent') {
      params.planContent = value;
    } else if (key === 'nonComplianceReason') {
      // 处理未达标原因，如果是数组转换为逗号分隔的字符串
      if (Array.isArray(value)) {
        params.nonComplianceReason = value.join(',');
      } else {
        params.nonComplianceReason = value;
      }
    } else {
      // 默认处理其他字段
      params[key] = value;
    }

    // 达标记录不允许编辑方案类型、方案内容、评分、未达标原因
    if (
      record.state === 1 &&
      ['planType', 'planContent', 'score', 'nonComplianceReason'].includes(key)
    ) {
      throw new Error('达标记录不允许编辑此字段');
    }

    // 暂无要求的记录不允许编辑方案内容、未达标原因
    if (record.state === 3 && ['planContent', 'nonComplianceReason'].includes(key)) {
      throw new Error('暂无要求记录不允许编辑此字段');
    }

    // 只有达标确认为"否"时才能编辑未达标原因
    if (record.stateManual !== 2 && key === 'nonComplianceReason') {
      throw new Error('只有达标确认为"否"时才能编辑未达标原因');
    }

    await editReports(params);
    // 不调用reload()
    message.success('更新成功');
  };

  // 安全获取数据源的函数
  const safeGetDataSource = () => {
    try {
      return getDataSource() || [];
    } catch (error) {
      console.warn('getDataSource error:', error);
      return [];
    }
  };

  const [registerTable, { getDataSource, reload }] = useTable({
    title: '利用率列表',
    api: unifiedApiCall,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
    },
    actionColumn: {
      width: 120,
      title: '操作',
      dataIndex: 'action',
      fixed: 'right',
    },
    useSearchForm: true,
    bordered: true,
    showIndexColumn: true,
    showTableSetting: true,
    canResize: true,
    rowKey: 'id',
    beforeEditSubmit,
  });
</script>
<style lang="scss" scoped>
  :deep(.ant-picker) {
    width: 100% !important;
  }

  :deep(.ant-tabs-nav-list) {
    margin-left: 20px !important;
  }

  .vben-basic-table-form-container {
    padding-top: 0 !important;
  }

  :deep(.vben-basic-table-form-container .ant-form) {
    margin-bottom: 8px !important;
  }

  :deep(.ant-dropdown-link-disabled) {
    color: rgb(0 0 0 / 25%);
    cursor: not-allowed;
    pointer-events: none;
  }

  :deep(.disabled-btn) {
    border-color: #d9d9d9 !important;
    background-color: #d9d9d9 !important;
    color: rgb(0 0 0 / 25%) !important;
  }
</style>
