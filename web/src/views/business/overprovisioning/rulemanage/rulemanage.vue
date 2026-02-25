<template>
  <div>
    <BasicTable @register="registerTable">
      <template #toolbar>
        <a-button
          type="primary"
          v-auth="'business:overprovisioning:rule:create'"
          @click="handleCreateRule"
        >
          新建规则
        </a-button>
        <a-button
          type="success"
          v-auth="'business:overprovisioning:rule:execute'"
          @click="handleExecuteDetection"
          :loading="executeLoading"
        >
          超配检测
        </a-button>
      </template>
      <template #bodyCell="{ column, record }">
        <template v-if="column.key == 'ruleStatus'">
          <Tag :color="record.ruleStatus === 0 ? 'green' : 'red'">{{
            record.ruleStatus === 0 ? '启用' : '禁用'
          }}</Tag>
        </template>
        <template v-if="column.key == 'ruleDescription'">
          <Tooltip :title="record.ruleDescription">
            <span>{{ record.ruleDescription }}</span>
          </Tooltip>
        </template>
        <!-- <template v-if="column.key == 'businesses'">
          <Tooltip :title="record.ID">
            <span>{{ record.businesses.map((b) => b.businessName).join(' | ') }}</span>
          </Tooltip>
        </template> -->
        <!-- 操作按钮 -->
        <template v-if="column.key == 'action'">
          <TableAction
            :actions="[
              {
                label: '编辑',
                icon: 'clarity:note-edit-line',
                auth: 'business:overprovisioning:rule:update',
                onClick: handleEdit.bind(null, record),
              },
              {
                label: '删除',
                icon: 'ic:outline-delete-outline',
                auth: 'business:overprovisioning:rule:delete',
                onClick: handleDelete.bind(null, record),
              },
            ]"
          />
        </template>
      </template>
    </BasicTable>
    <RuleModal @register="registerRuleModal" @success="handleRuleSuccess" />
  </div>
</template>

<script lang="ts" setup>
  import { onMounted, ref, h } from 'vue';
  import { Tooltip, Tag, Modal, message, Radio } from 'ant-design-vue';
  import { BasicTable, useTable, TableAction } from '@/components/Table';
  import { useModal } from '@/components/Modal';
  import RuleModal from './RuleModal.vue';
  import { columns, searchFormSchema } from './data';
  import { getRuleList, deleteRule, executeDetection } from '@/api/business/overprovisioning';

  defineOptions({ name: 'Rulemanage' });

  const executeLoading = ref(false);
  const execType = ref('all');

  const [registerTable, { reload }] = useTable({
    title: '规则管理',
    api: getRuleList,
    columns,
    formConfig: {
      labelWidth: 120,
      schemas: searchFormSchema(),
    },
    useSearchForm: true,
    showTableSetting: true,
    bordered: true,
    showIndexColumn: false,
    pagination: {},
    rowKey: 'id',
    beforeFetch: (params) => {
      // parseParams(params);
      return params;
    },
  });

  const [registerRuleModal, { openModal }] = useModal();

  function handleCreateRule() {
    openModal(true, { isUpdate: false });
  }

  function handleRuleSuccess() {
    reload();
  }

  function handleEdit(record) {
    openModal(true, { isUpdate: true, record });
  }

  function handleDelete(record) {
    Modal.confirm({
      title: '确认删除',
      content: `确定要删除规则 "${record.ruleName || record.id}" 吗？`,
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        try {
          await deleteRule(record.id);
          message.success('删除成功');
          reload();
        } catch (error) {
          message.error('删除失败');
          console.error(error);
        }
      },
    });
  }

  // 执行超配检测
  async function handleExecuteDetection() {
    execType.value = 'all'; // 重置为默认值
    Modal.confirm({
      title: '执行超配检测',
      content: h({
        setup() {
          return () =>
            h('div', [
              h('p', { style: 'margin-bottom: 12px;' }, '请选择执行类型：'),
              h(
                Radio.Group,
                {
                  value: execType.value,
                  'onUpdate:value': (val: string) => {
                    execType.value = val;
                  },
                },
                {
                  default: () => [
                    h(Radio, { value: 'kvm' }, { default: () => '虚拟机' }),
                    h(Radio, { value: 'host' }, { default: () => '物理机' }),
                    h(Radio, { value: 'all' }, { default: () => '全部' }),
                  ],
                },
              ),
            ]);
        },
      }),
      okText: '确定',
      cancelText: '取消',
      onOk: async () => {
        executeLoading.value = true;
        try {
          await executeDetection(execType.value);
          message.success('超配检测执行中，请稍后查看结果');
          reload();
        } catch (error) {
          message.error('超配检测执行失败');
          console.error(error);
        } finally {
          executeLoading.value = false;
        }
      },
    });
  }

  onMounted(() => {
    // getForm().setFieldsValue({
    //   ruleName: '',
    //   mem: '',
    //   ssdSize: '',
    //   hddSize: '',
    //   // storageBwRatio: '', // 已移除，使用 storageBwRatioRead 和 storageBwRatioWrite
    //   businessBelong: '',
    // });
  });
</script>
