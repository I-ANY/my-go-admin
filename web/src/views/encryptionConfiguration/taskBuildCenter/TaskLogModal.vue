<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="'任务构建日志'"
    @cancel="handleCancel"
  >
   <template #title>
    <div class="task-status">
      <span style="margin-right: 20px;">任务构建日志</span>
      <div style="font-weight: 400;font-size: 14px;" :class="status === 'completed' ? 'text-success' : 'text-danger'" v-if="status !== 'running'">{{ taskStatus[status] }}</div>
      <Button v-if="status === 'failed'" type="primary" style="margin-left: 10px;" @click="handleRebuild">重新构建</Button>
      <div 
        style="font-weight: 400;font-size: 14px;"
        class="text-progress building-status"
        v-if="status === 'running'"
      >
        <span class="building-text building-animated-text">
          <span>构建中</span>
          <span class="loading-dots-animated">
            <span class="dot dot-1">.</span>
            <span class="dot dot-2">.</span>
            <span class="dot dot-3">.</span>
          </span>
        </span>
      </div>
    </div>
   </template>
    <div class="console-log-container">
      <div class="console-log-content" ref="logContentRef">
        <div v-for="(line, index) in logLines" :key="index" :class="getLineClass(line)">
          {{ line }}
        </div>
      </div>
    </div>
    <div class="console-log-footer">
      <Space>
        <Tooltip title="回到底部">
          <Button shape="circle" @click="scrollToBottom" :icon="h(DownOutlined)" />
        </Tooltip>
        <Tooltip title="复制日志">
          <Button shape="circle" @click="handleCopy" :icon="h(CopyOutlined)" />
        </Tooltip>
      </Space>
    </div>
  </BasicModal>
</template>

<script lang="ts" setup>
import { ref, nextTick, h, computed, watch, onUnmounted } from 'vue';
import { BasicModal, useModalInner } from '@/components/Modal';
import { Space, Tooltip, Button, message } from 'ant-design-vue';
import { DownOutlined, CopyOutlined } from '@ant-design/icons-vue';
import { createTaskBuildCenter, getTaskBuildCenterDetail } from '@/linuxApi/config';

defineOptions({
  name: 'TaskLogModal',
});

const taskStatus = {
  'completed': '构建完成',
  'failed': '构建失败',
  'running': '构建中',
}

const emit = defineEmits(['register', 'success']);

const logContentRef = ref<HTMLElement>();
const logText = ref<string>('');
const status = ref<string>('');
const taskRecord = ref<any>({});
const taskId = ref<string>('');
const pollTimer = ref<NodeJS.Timeout | null>(null);

// 获取任务详情
async function fetchTaskDetail() {
  if (!taskId.value) return;
  
  try {
    const res = await getTaskBuildCenterDetail({ task_id: taskId.value });
    logText.value = res.logText;
    const newStatus = res.status
    status.value = newStatus;
    
    nextTick(() => {
      scrollToBottom();
    });
    
    // 如果状态不再是 running，停止轮询
    if (newStatus !== 'running') {
      stopPolling();
    }
  } catch (error) {
    stopPolling();
  }
}

// 开始轮询
function startPolling() {
  // 先清除可能存在的定时器
  stopPolling();
  
  if (status.value === 'running' && taskId.value) {
    pollTimer.value = setInterval(() => {
      fetchTaskDetail();
    }, 10000); // 10秒轮询
  }
}

// 停止轮询
function stopPolling() {
  if (pollTimer.value) {
    clearInterval(pollTimer.value);
    pollTimer.value = null;
  }
}

// 监听状态变化，如果变为 running 则开始轮询
watch(status, (newStatus) => {
  if (newStatus === 'running') {
    startPolling();
  } else {
    stopPolling();
  }
});

// 重新构建
async function handleRebuild() {
  const formData = {
    creator: JSON.parse(localStorage.getItem('userInfo') || '{}').nickname,
    iso_sign: taskRecord.value.iso_sign,
    basename: taskRecord.value.basename,
    os_type: taskRecord.value.os_type,
    tar_version: taskRecord.value.tar_version,
    sysdisk_minsize: taskRecord.value.sysdisk_minsize,
    pppoe_type: taskRecord.value.pppoe_type,
    create_pcdn_index_data: taskRecord.value.create_pcdn_index_data,
    auto_register: taskRecord.value.auto_register,
    syunhost_account: taskRecord.value.syunhost_account,
    yunhost_password: taskRecord.value.yunhost_password,
  }
  const res = await createTaskBuildCenter(formData);
  if (res.task_id) {
    message.success('重新构建成功');
    taskId.value = res.task_id;
    await fetchTaskDetail();
    startPolling();
  } else {
    message.error('重新构建失败');
  }
}

const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
  setModalProps({
    confirmLoading: false,
    footer: null,
    width: '70%',
  });
  
  // 先停止之前的轮询
  stopPolling();
  
  if (data.record.task_id) {
    taskRecord.value = data.record;
    taskId.value = data.record.task_id;
    await fetchTaskDetail();
    
    // 如果状态是 running，启动轮询
    if (status.value === 'running') {
      startPolling();
    }
  }
});



// 组件卸载时清理定时器
onUnmounted(() => {
  stopPolling();
});

// 将日志文本按行分割
const logLines = computed(() => {
  if (!logText.value) return [];
  return logText.value.split('\n');
});

// 根据日志行内容判断样式类
function getLineClass(line: string): string {
  const trimmedLine = line.trim();
  
  if (trimmedLine.startsWith('[ERROR]') || trimmedLine.includes('构建失败')) {
    return 'log-line log-error';
  }
  if (trimmedLine.startsWith('[STEP]')) {
    return 'log-line log-step';
  }
  if (trimmedLine.startsWith('[INFO]')) {
    return 'log-line log-info';
  }
  if (trimmedLine.startsWith('[PROGRESS]')) {
    return 'log-line log-progress';
  }
  if (trimmedLine.startsWith('任务ID:') || trimmedLine.startsWith('开始时间:') || trimmedLine.startsWith('结束时间:') || trimmedLine.startsWith('退出码:')) {
    return 'log-line log-header';
  }
  if (trimmedLine.startsWith('================================================================================')) {
    return 'log-line log-separator';
  }
  
  return 'log-line';
}

// 滚动到底部
function scrollToBottom() {
  nextTick(() => {
    if (logContentRef.value) {
      logContentRef.value.scrollTop = logContentRef.value.scrollHeight;
    }
  });
}

// 复制日志内容
async function handleCopy() {
  try {
    await navigator.clipboard.writeText(logText.value);
    message.success('日志已复制到剪贴板');
  } catch (error) {
    // 降级方案：使用传统方法
    const textarea = document.createElement('textarea');
    textarea.value = logText.value;
    textarea.style.position = 'fixed';
    textarea.style.opacity = '0';
    document.body.appendChild(textarea);
    textarea.select();
    try {
      document.execCommand('copy');
      message.success('日志已复制到剪贴板');
    } catch (err) {
      message.error('复制失败');
    }
    document.body.removeChild(textarea);
  }
}

function handleCancel() {
  stopPolling();
  logText.value = '';
  emit('success');
}
</script>

<style scoped lang="less">
.task-status {
  display: flex;
  align-items: center;
  .text-success {
    color: #52c41a;
  }
  .text-danger {
    color: #ff4d4f;
  }
  .text-progress {
    color: #0910e9;
  }
  .building-status {
    display: flex;
    align-items: center;
    .building-animated-text {
      display: flex;
      align-items: center;
    }
    .loading-dots-animated {
      display: inline-block;
      margin-left: 2px;
      .dot {
        display: inline-block;
        font-weight: bold;
        font-size: 20px;
        opacity: 0;
        animation: dotFlashing 1s infinite linear;
      }
      .dot-1 {
        animation-delay: 0s;
      }
      .dot-2 {
        animation-delay: 0.3s;
      }
      .dot-3 {
        animation-delay: 0.6s;
      }
    }
  }
}

@keyframes dotFlashing {
  0% { opacity: 0; }
  20% { opacity: 1; }
  100% { opacity: 0; }
}

.console-log-container {
  width: 100%;
  height: 70vh;
  overflow: hidden;
  border-radius: 4px;
  background-color: #1e1e1e;
  border: 1px solid #333;
}

.console-log-content {
  width: 100%;
  height: 100%;
  padding: 16px;
  overflow-y: auto;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.6;
  color: #d4d4d4;
  background-color: #1e1e1e;
  
  // 自定义滚动条样式
  &::-webkit-scrollbar {
    width: 8px;
    height: 8px;
  }
  
  &::-webkit-scrollbar-track {
    background: #252526;
  }
  
  &::-webkit-scrollbar-thumb {
    background: #424242;
    border-radius: 4px;
    
    &:hover {
      background: #4e4e4e;
    }
  }
}

.log-line {
  white-space: pre-wrap;
  word-break: break-word;
  margin-bottom: 2px;
  
  &.log-error {
    color: #f48771;
  }
  
  &.log-step {
    color: #4ec9b0;
    font-weight: 500;
  }
  
  &.log-info {
    color: #569cd6;
  }
  
  &.log-progress {
    color: #ce9178;
  }
  
  &.log-header {
    color: #dcdcaa;
    font-weight: 500;
  }
  
  &.log-separator {
    color: #808080;
  }
}

.console-log-footer {
  margin-top: 12px;
  display: flex;
  justify-content: flex-end;
  padding-top: 12px;
  border-top: 1px solid #d9d9d9;
}
</style>