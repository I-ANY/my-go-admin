<template>
  <BasicModal v-bind="$attrs" @register="registerModal" :title="title" @ok="handleSubmit">
    <CodeEditor class="code" :readonly="true" v-model:value="content" :mode="modeValue" />
  </BasicModal>
</template>
<script lang="ts">
  import { defineComponent, ref } from 'vue';
  import { useModalInner, BasicModal } from '@/components/Modal';
  import { CodeEditor, MODE } from '@/components/CodeEditor';

  export default defineComponent({
    name: 'JsonDataModal',
    components: { BasicModal, CodeEditor },
    emits: ['register'],
    setup(_) {
      const content = ref();
      const title = ref('');
      const modeValue = ref(MODE.JSON);

      const [registerModal, { setModalProps, closeModal }] = useModalInner(async (data) => {
        setModalProps({ width: 1200 });
        content.value = data.content;
        title.value = data.title;
      });

      async function handleSubmit() {
        try {
          setModalProps({ confirmLoading: true });
          closeModal();
          content.value = '';
        } finally {
          setModalProps({ confirmLoading: false });
        }
      }

      return { content, title, modeValue, registerModal, handleSubmit };
    },
  });
</script>
