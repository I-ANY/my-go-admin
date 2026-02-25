<template>
  <div class="test-container">
    <h1>换行符修复测试</h1>

    <div class="test-section">
      <h2>测试按钮</h2>
      <div class="test-buttons">
        <button @click="testUserMessage">测试用户消息换行</button>
        <button @click="testBotMessage">测试机器人消息换行</button>
        <button @click="testComplexMessage">测试复杂消息</button>
      </div>
    </div>

    <div class="test-section">
      <h2>ChatBot 组件</h2>
      <ChatBot ref="chatBotRef" :options="chatBotOptions" />
    </div>
  </div>
</template>

<script setup lang="ts">
  import { ref } from 'vue';
  import ChatBot from './ChatBot.vue';
  import type { ChatBotOptions } from './types';

  const chatBotRef = ref();

  const chatBotOptions: ChatBotOptions = {
    botName: '换行测试助手',
    placeholder: '输入消息进行测试...',
    replyDelay: 300,
    requireLogin: false,

    defaultReply: async (input: string) => {
      if (input.includes('换行')) {
        return `这是第一行\n这是第二行\n这是第三行\n\n空行后的内容`;
      }

      if (input.includes('复杂')) {
        return `<think>\n嗯，用户问"你是谁"。根据知识信息，我只能基于提供的内容回答。这里没有关于用户身份的信息，所以无法提供具体帮助。需要告知用户暂时无法回答，并请提供更多信息。\n</think>\n您好！我是您的专属运维专家，但我暂时还未获取到足够的信息来回答您的问题。请您提供更多的信息，以便我更好地帮助您。`;
      }

      return `请测试以下格式：\n1. 输入"换行"测试换行符\n2. 输入"复杂"测试复杂消息`;
    },
  };

  // 测试函数
  const testUserMessage = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示带换行的文本', true);
    }
  };

  const testBotMessage = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示带换行的文本', true);
    }
  };

  const testComplexMessage = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示复杂消息', true);
    }
  };
</script>

<style scoped>
  .test-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 20px;
  }

  .test-section {
    margin-bottom: 30px;
  }

  .test-section h2 {
    margin-bottom: 15px;
    color: #333;
  }

  .test-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    margin-bottom: 20px;
  }

  .test-buttons button {
    padding: 10px 16px;
    transition: all 0.2s;
    border: 1px solid #ddd;
    border-radius: 6px;
    background: #f8f9fa;
    font-size: 14px;
    cursor: pointer;
  }

  .test-buttons button:hover {
    border-color: #667eea;
    background: #667eea;
    color: white;
  }
</style>
