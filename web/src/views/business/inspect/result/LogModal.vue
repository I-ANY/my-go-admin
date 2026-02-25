<template>
  <BasicModal
    v-bind="$attrs"
    @register="registerModal"
    :title="'查看详情'"
    ref="logModal"
    @cancel="handleCancel"
  >
    <Spin tip="加载中..." :spinning="loading">
      <div class="scroll-wrap w-full h-450px">
        <ScrollContainer ref="scrollRef">
          <CodeEditor class="code" :readonly="true" v-model:value="text" :mode="modeValue" />
          <!-- <Affix :offset-bottom="5" :target="() => scrollRef as any"> aaa </Affix> -->
        </ScrollContainer>
      </div>
    </Spin>
    <Space warp class="mt-5px">
      <Tooltip title="回到底部">
        <Button shape="circle" @click="handToBottom" :icon="h(DownOutlined)" />
      </Tooltip>
    </Space>
  </BasicModal>
</template>
<script lang="ts" setup>
  import { ref, h, unref, nextTick } from 'vue';
  import { BasicModal, useModalInner } from '@/components/Modal';
  import { CodeEditor, MODE } from '@/components/CodeEditor';
  import { Space, Tooltip, Button, Spin } from 'ant-design-vue';

  import { ScrollContainer, ScrollActionType } from '@/components/Container';
  import { GetResultField } from '@/api/business/inspect';
  import { DownOutlined } from '@ant-design/icons-vue';

  defineOptions({
    name: 'LogModal',
  });
  const scrollRef = ref<Nullable<ScrollActionType>>(null);

  const text = ref('');
  const modeValue = ref(MODE.PLAIN);
  const logModal = ref(null);
  const loading = ref<boolean>(false);

  const emit = defineEmits(['success', 'register']);

  const [registerModal, { setModalProps }] = useModalInner(async (data) => {
    loading.value = true;
    setModalProps({
      confirmLoading: false,
      width: 1200,
      height: 500,
      footer: null,
      title: data.name + ' 详情',
    });
    refreshLog(data);
  });

  async function refreshLog(data: Recordable) {
    try {
      await sleep(300).then(() => {
        //code
      });
      const res = await GetResultField(data.id);
      if (res.detail) {
        let log = res.detail.replaceAll('\n', '\r\n');
        text.value += log;
        if (!log.endsWith('\r\n')) {
          text.value += '\r\n';
        }
      }
      nextTick(() => {
        scrollBottom();
      });
    } finally {
      loading.value = false;
    }
  }
  function sleep(ms) {
    return new Promise((resolve) => setTimeout(resolve, ms));
  }

  async function handleCancel() {
    text.value = '';
    emit('success');
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
  function handToBottom() {
    scrollBottom();
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
