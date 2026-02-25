<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    @cancel="handleCancel"
    class="modify-dscp-result-modal"
  >
    <div class="summary">
      <span class="summary-item success">
        <CheckCircleFilled class="icon" />
        成功：{{ successCount }} 台
      </span>
      <span class="summary-item failed">
        <CloseCircleFilled class="icon" />
        失败：{{ failedCount }} 台
      </span>
    </div>
    <div class="result-list">
      <div
        v-for="(item, index) in sortedResults"
        :key="index"
        :class="['result-item', item.status === ModifyResultStatus.SUCCESS ? 'success' : 'failed']"
      >
        <CheckCircleFilled v-if="item.status === ModifyResultStatus.SUCCESS" class="icon" />
        <CloseCircleFilled v-else class="icon" />
        <span class="label">主机名：</span>
        <span class="hostname">{{ item.hostname }}</span>
        <span v-if="item.message" class="label msg-label">Message：</span>
        <span v-if="item.message" class="message">{{ item.message }}</span>
      </div>
      <div v-if="data.results.length === 0" class="empty">暂无数据</div>
    </div>
  </BasicModal>
</template>

<script setup lang="ts">
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { reactive, computed } from 'vue';
  import { CheckCircleFilled, CloseCircleFilled } from '@ant-design/icons-vue';
  import { ModifyDscpResult, ModifyResultStatus } from './data';

  const emit = defineEmits(['success', 'register']);

  const data = reactive({
    results: [] as Array<ModifyDscpResult>,
  });

  // 成功数量
  const successCount = computed(() => {
    return data.results.filter((item) => item.status === ModifyResultStatus.SUCCESS).length;
  });

  // 失败数量
  const failedCount = computed(() => {
    return data.results.filter((item) => item.status === ModifyResultStatus.FAILED).length;
  });

  // 排序后的结果：失败在前，成功在后
  const sortedResults = computed(() => {
    return [...data.results].sort((a, b) => {
      if (a.status === b.status) return 0;
      return a.status === ModifyResultStatus.FAILED ? -1 : 1;
    });
  });

  defineOptions({ name: 'ModifyDscpResultModal' });
  const [registerModal, { setModalProps }] = useModalInner(async (d) => {
    setModalProps({
      title: '修改DSCP值结果',
      footer: null,
      width: 900,
    });
    data.results = (d.results as Array<ModifyDscpResult>) || [];
  });
  function handleCancel() {
    emit('success');
  }
</script>

<style scoped>
  .summary {
    display: flex;
    gap: 24px;
    margin-bottom: 16px;
    padding: 5px 16px 12px;
    border-radius: 6px;
    background-color: #fafafa;
  }

  .summary-item {
    display: flex;
    align-items: center;
    font-size: 15px;
    font-weight: 500;
  }

  .summary-item .icon {
    margin-right: 6px;
    font-size: 16px;
  }

  .summary-item.success {
    color: #52c41a;
  }

  .summary-item.failed {
    color: #ff4d4f;
  }

  .result-list {
    max-height: 450px;
    overflow-y: auto;
  }

  .result-item {
    display: flex;
    align-items: center;
    margin-bottom: 8px;
    padding: 6px 12px;
    border-radius: 6px;
    font-size: 14px;
  }

  .result-item.success {
    border: 1px solid #b7eb8f;
    background-color: #f6ffed;
    color: #52c41a;
  }

  .result-item.failed {
    border: 1px solid #ffccc7;
    background-color: #fff2f0;
    color: #ff4d4f;
  }

  .result-item .icon {
    margin-right: 10px;
    font-size: 18px;
  }

  .result-item .label {
    width: 70px;
    color: #888;
    font-size: 13px;
  }

  .result-item .label.msg-label {
    margin-left: 14px;
  }

  .result-item .hostname {
    color: #333;
    font-weight: 500;
  }

  .result-item .message {
    color: #666;
    font-size: 13px;
  }

  .result-item.failed .hostname {
    color: #cf1322;
  }

  .result-item.success .hostname {
    color: #389e0d;
  }

  .empty {
    padding: 40px;
    color: #999;
    text-align: center;
  }
</style>
