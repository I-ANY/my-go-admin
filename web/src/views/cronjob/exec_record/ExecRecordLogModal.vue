<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="'查看日志'"
    ref="logModal"
    @ok="handleSubmit"
    @cancel="handleCancel"
  >
    <div class="scroll-wrap w-full h-450px">
      <ScrollContainer ref="scrollRef">
        <CodeEditor class="code" :readonly="true" v-model:value="text" :mode="modeValue" />
        <!-- <Affix :offset-bottom="5" :target="() => scrollRef as any"> aaa </Affix> -->
      </ScrollContainer>
    </div>
    <Space warp class="mt-5px">
      <Tooltip title="回到底部">
        <Button shape="circle" @click="handToBottom" :icon="h(DownOutlined)" />
      </Tooltip>
      <Spin v-if="loading" class="mt-1 ml-1" />
    </Space>
  </BasicModal>
</template>
<script lang="ts" setup>
  import { nextTick, ref, onBeforeUnmount, unref, onMounted, h } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { CodeEditor, MODE } from '@/components/CodeEditor';
  import { getExecRecordLog } from '@/api/cronjob/job';
  import { Spin, Button, Tooltip, Space } from 'ant-design-vue';
  import { DownOutlined } from '@ant-design/icons-vue';

  import { ScrollContainer, ScrollActionType } from '@/components/Container';

  defineOptions({
    name: 'ExecRecordLogModal',
  });
  const scrollRef = ref<Nullable<ScrollActionType>>(null);

  const text = ref('');
  const modeValue = ref(MODE.PLAIN);
  const logModal = ref(null);
  const loading = ref<boolean>(false);
  const needToBottom = ref<boolean>(true);

  let timer: IntervalHandle;
  let lastId = 0; // 默认第一行
  const emit = defineEmits(['success', 'register']);

  const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
    lastId = 0;
    setModalProps({
      confirmLoading: false,
      width: 1200,
      height: 500,
      footer: null,
    });
    loading.value = true;
    needToBottom.value = true;
    if (await refreshLog(data)) {
      // 有更多日志再启动定时器
      timer = setInterval(() => refreshLog(data), 1000);
    }
  });

  async function refreshLog(data: Recordable) {
    await sleep(300).then(() => {
      //code
    });
    const res = await getExecRecordLog(data.id, { lastId: lastId });
    if (res.logs && res.logs.length > 0) {
      res.logs.forEach((e) => {
        e.message = e.message.replaceAll('\n', '\r\n');
        text.value += `${e.logTime} [${e.level}] ${e.message}`;
        if (!e.message.endsWith('\r\n')) {
          text.value += '\r\n';
        }
      });
    }

    lastId = res.lastId;
    if (res.runStatus != 1) {
      loading.value = false;
      // 执行结束，没有日志了，清除定时器
      clearInterval(timer);
    }
    if (res.logs && res.logs.length > 0) {
      // console.log('logModal', unref(logModal));
      nextTick(() => {
        if (unref(needToBottom)) {
          scrollBottom();
        }
        // 滚动到最底部
        // const modalContent = (unref(logModal) as any)?.modalRef?.getContentRef();
        // modalContent.scrollTop = modalContent.scrollHeight;
      });
    }
    return res.runStatus == 1;
  }
  function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }
  async function handleSubmit() {
    try {
      setModalProps({ confirmLoading: true });
      closeModal();
      emit('success');
      clearInterval(timer);
    } finally {
      setModalProps({ confirmLoading: false });
    }
  }
  onBeforeUnmount(() => {
    clearInterval(timer);
  });

  async function handleCancel() {
    text.value = '';
    emit('success');
    clearInterval(timer);
  }

  const getScroll = () => {
    const scroll = unref(scrollRef);
    if (!scroll) {
      throw new Error('scroll is Null');
    }
    return scroll;
  };
  function scrollBottom() {
    getScroll().scrollBottom();
  }

  // 处理滚轮事件的方法
  const handleWheel = (event) => {
    const deltaY = event.deltaY;
    if (deltaY < 0) {
      needToBottom.value = false;
    } else if (deltaY > 0) {
      // scrollDirection.value = 'down'; // 向下滚动
    }
  };

  // 在组件挂载后添加监听器
  onMounted(() => {
    window.addEventListener('wheel', handleWheel);
  });

  // 在组件卸载前移除监听器
  onBeforeUnmount(() => {
    window.removeEventListener('wheel', handleWheel);
  });
  function handToBottom() {
    scrollBottom();
    needToBottom.value = true;
  }
</script>
<style scoped lang="less">
  /* light模式下，codemirror会使用idea主题，修改背景色 */
  ::v-deep(.cm-s-idea) {
    background-color: #f2f2f2 !important;
  }

  .scroll-wrap {
    // width: 100%;
    // height: 400px;
    background-color: @component-background;
  }
</style>
