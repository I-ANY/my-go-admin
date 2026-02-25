<template>
  <div class="test-container">
    <h1>去除末尾换行符测试</h1>

    <div class="test-section">
      <h2>修复说明</h2>
      <div class="fix-info">
        <p><strong>修复内容：</strong></p>
        <ul>
          <li>去除消息内容首尾的空白字符</li>
          <li>避免末尾出现多余的换行符</li>
          <li>确保段落标签正确包裹内容</li>
        </ul>
      </div>
    </div>

    <div class="test-section">
      <h2>测试按钮</h2>
      <div class="test-buttons">
        <button @click="testTrailingNewline">测试末尾换行</button>
        <button @click="testLeadingNewline">测试开头换行</button>
        <button @click="testBothNewlines">测试首尾换行</button>
        <button @click="testEmptyContent">测试空内容</button>
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
    botName: '换行符修复测试助手',
    placeholder: '输入消息进行测试...',
    replyDelay: 300,
    requireLogin: false,

    defaultReply: async (input: string) => {
      if (input.includes('末尾')) {
        return `这是第一行\n这是第二行\n这是第三行\n`;
      }

      if (input.includes('开头')) {
        return `\n这是第一行\n这是第二行\n这是第三行`;
      }

      if (input.includes('首尾')) {
        return `\n这是第一行\n这是第二行\n这是第三行\n`;
      }

      if (input.includes('空内容')) {
        return `\n\n\n`;
      }

      return `请测试以下格式：\n1. 输入"末尾"测试末尾换行\n2. 输入"开头"测试开头换行\n3. 输入"首尾"测试首尾换行\n4. 输入"空内容"测试空内容`;
    },
  };

  // 测试函数
  const testTrailingNewline = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示末尾换行', true);
    }
  };

  const testLeadingNewline = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示开头换行', true);
    }
  };

  const testBothNewlines = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示首尾换行', true);
    }
  };

  const testEmptyContent = () => {
    if (chatBotRef.value) {
      chatBotRef.value.addMessage('请显示空内容', true);
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

  .fix-info {
    padding: 15px;
    border-left: 4px solid #667eea;
    border-radius: 8px;
    background: #f8f9fa;
  }

  .fix-info ul {
    margin: 10px 0 0 20px;
  }

  .fix-info li {
    margin: 5px 0;
    line-height: 1.5;
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
