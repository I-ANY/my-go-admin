<template>
  <div class="assessment-container">
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="success"
          @click="handleImport"
          v-if="hasPermission('businessA:assessment:import')"
        >
          批量导入
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key === 'action'">
          <TableAction
            :actions="[
              {
                icon: 'clarity:note-edit-line',
                tooltip: '编辑考核规则',
                onClick: handleEdit.bind(null, record),
                ifShow: () => hasPermission('businessA:assessment:edit'),
              },
            ]"
          />
        </template>
        <template v-if="column.key === 'billingType'">
          <Tag :color="getBillingTypeColor(record.billingType)">
            {{ record.billingType }}
          </Tag>
        </template>
        <template v-if="column.key === 'roomType'">
          <Tag :color="getRoomTypeColor(record.roomType)">
            {{ getRoomTypeText(record.roomType) }}
          </Tag>
        </template>
        <template v-if="column.key === 'statType'">
          <Tag :color="getStatTypeColor(record.statType)">
            {{ getStatTypeText(record.statType) }}
          </Tag>
        </template>
        <template v-if="column.key === 'assessmentStandard'">
          <div class="assessment-standard">
            <div v-if="record.assessmentStandard">
              <div v-if="record.assessmentStandard.utilizationRateThreshold">
                日95&晚高峰95利用率:
                {{ (record.assessmentStandard.utilizationRateThreshold * 100).toFixed(1) }}%
              </div>
              <div v-if="record.assessmentStandard.nightPeakPointsThreshold">
                晚高峰点数: {{ record.assessmentStandard.nightPeakPointsThreshold }}
              </div>
              <div v-if="record.assessmentStandard.customRules">
                {{ record.assessmentStandard.customRules }}
              </div>
            </div>
            <span v-else>-</span>
          </div>
        </template>
      </template>
    </BasicTable>

    <!-- 编辑模态框 -->
    <AssessmentModal @register="registerModal" @success="handleSuccess" />
    <!-- 批量导入模态框 -->
    <ImportModal @register="registerImportModal" @success="handleSuccess" />
  </div>
</template>

<script lang="ts" setup>
  // import { ref } from 'vue';
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { useModal } from '@/components/Modal';
  import { Tag } from 'ant-design-vue';
  import { columns, searchFormSchema } from './data';
  import AssessmentModal from './AssessmentModal.vue';
  import ImportModal from './ImportModal.vue';
  import { getAssessmentRules } from '@/api/business/a';
  import { usePermission } from '@/hooks/web/usePermission';

  defineOptions({ name: 'AssessmentList' });

  const { hasPermission } = usePermission();
  const [registerImportModal, { openModal: openImportModal }] = useModal();

  // 模态框
  const [registerModal, { openModal }] = useModal();

  // 表格配置
  const [registerTable, { reload }] = useTable({
    title: '考核规则列表',
    api: getAssessmentRules,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema,
      autoSubmitOnEnter: true,
      showAdvancedButton: true,
      autoAdvancedLine: 4,
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    rowKey: 'id',
    showIndexColumn: false,
    actionColumn: {
      width: 120,
      title: '操作',
      dataIndex: 'action',
      fixed: 'right',
    },
  });

  // 计费类型颜色映射
  const getBillingTypeColor = (type: string) => {
    const colorMap = {
      日95: 'blue',
      月95: 'green',
      买断: 'orange',
    };
    return colorMap[type] || 'default';
  };

  // 机房类型颜色映射
  const getRoomTypeColor = (type: number) => {
    const colorMap = {
      1: 'blue',
      2: 'green',
      3: 'orange',
    };
    return colorMap[type] || 'default';
  };

  // 机房类型文本
  const getRoomTypeText = (type: number) => {
    const typeMap = { 1: 'IDC', 2: 'ACDN', 3: 'MCDN' };
    return typeMap[type] || type;
  };

  // 统计类型颜色映射
  const getStatTypeColor = (type: number) => {
    const colorMap = {
      1: 'blue',
      2: 'green',
      3: 'orange',
    };
    return colorMap[type] || 'default';
  };

  // 统计类型文本
  const getStatTypeText = (type: number) => {
    const typeMap = { 1: '机房总览', 2: '保底业务', 3: '削峰业务' };
    return typeMap[type] || type;
  };

  // 编辑
  function handleEdit(record: Recordable) {
    openModal(true, {
      record,
      isUpdate: true,
    });
  }

  // 操作成功回调
  function handleSuccess() {
    reload();
  }

  // 批量导入
  function handleImport() {
    openImportModal(true);
  }
</script>

<style lang="less" scoped>
  .assessment-container {
    padding: 16px;
  }

  .assessment-standard {
    font-size: 12px;
    line-height: 1.4;

    div {
      margin-bottom: 2px;
    }
  }
</style>
